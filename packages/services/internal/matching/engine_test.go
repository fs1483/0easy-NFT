// Package matching contains unit tests for the matching engine.
// Tests validate order matching logic, Redis integration, and match detection algorithms.
//
// Test Strategy:
// - Uses miniredis for isolated Redis testing
// - Validates matching algorithm correctly pairs compatible ask/bid orders
// - Ensures expired orders are filtered out
package matching

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Oeasy-NFT/services/internal/config"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

// setupTestMatchingEngine creates an in-memory matching engine for testing.
func setupTestMatchingEngine(t *testing.T) (*Engine, *redis.Client, func()) {
	t.Helper()

	srv, err := miniredis.Run()
	require.NoError(t, err)

	cfg := &config.Config{
		RedisAddr:       srv.Addr(),
		MarketplaceAddr: "0x0000000000000000000000000000000000000001",
		RPCURL:          "http://localhost",
		ChainID:         1,
	}

	engine, err := NewEngine(cfg)
	require.NoError(t, err)

	cleanup := func() {
		srv.Close()
	}

	return engine, engine.redisClient, cleanup
}

// TestFindMatches_ExactPriceMatch validates that compatible ask/bid pairs are matched.
func TestFindMatches_ExactPriceMatch(t *testing.T) {
	engine, _, cleanup := setupTestMatchingEngine(t)
	defer cleanup()

	asks := []Order{
		{
			ID:           1,
			Maker:        "0xaaa",
			NFTAddress:   "0x0000000000000000000000000000000000000002",
			TokenID:      "1",
			PaymentToken: "0x0000000000000000000000000000000000000003",
			Price:        "1000000000000000000",
			Side:         "ask",
		},
	}

	bids := []Order{
		{
			ID:           2,
			Maker:        "0xbbb",
			NFTAddress:   "0x0000000000000000000000000000000000000002",
			TokenID:      "1",
			PaymentToken: "0x0000000000000000000000000000000000000003",
			Price:        "1000000000000000000",
			Side:         "bid",
		},
	}

	matches := engine.findMatches(asks, bids)
	require.Len(t, matches, 1)
	require.Equal(t, asks[0].ID, matches[0].Ask.ID)
	require.Equal(t, bids[0].ID, matches[0].Bid.ID)
}

// TestFindMatches_DifferentNFT ensures orders for different NFTs don't match.
func TestFindMatches_DifferentNFT(t *testing.T) {
	engine, _, cleanup := setupTestMatchingEngine(t)
	defer cleanup()

	asks := []Order{
		{
			NFTAddress:   "0x0000000000000000000000000000000000000002",
			TokenID:      "1",
			PaymentToken: "0x0000000000000000000000000000000000000003",
			Price:        "1000000000000000000",
		},
	}

	bids := []Order{
		{
			NFTAddress:   "0x0000000000000000000000000000000000000002",
			TokenID:      "2", // Different token ID
			PaymentToken: "0x0000000000000000000000000000000000000003",
			Price:        "1000000000000000000",
		},
	}

	matches := engine.findMatches(asks, bids)
	require.Len(t, matches, 0)
}

// TestFindMatches_BidPriceTooLow ensures bid below ask price prevents matching.
func TestFindMatches_BidPriceTooLow(t *testing.T) {
	engine, _, cleanup := setupTestMatchingEngine(t)
	defer cleanup()

	asks := []Order{
		{
			NFTAddress:   "0x0000000000000000000000000000000000000002",
			TokenID:      "1",
			PaymentToken: "0x0000000000000000000000000000000000000003",
			Price:        "2000000000000000000", // Ask: 2 ETH
		},
	}

	bids := []Order{
		{
			NFTAddress:   "0x0000000000000000000000000000000000000002",
			TokenID:      "1",
			PaymentToken: "0x0000000000000000000000000000000000000003",
			Price:        "1000000000000000000", // Bid: 1 ETH (too low)
		},
	}

	matches := engine.findMatches(asks, bids)
	require.Len(t, matches, 0, "bid price below ask should not match")
}

// TestFindMatches_BidExceedsAsk ensures bid >= ask results in a valid match.
func TestFindMatches_BidExceedsAsk(t *testing.T) {
	engine, _, cleanup := setupTestMatchingEngine(t)
	defer cleanup()

	asks := []Order{
		{
			ID:           1,
			NFTAddress:   "0x0000000000000000000000000000000000000002",
			TokenID:      "1",
			PaymentToken: "0x0000000000000000000000000000000000000003",
			Price:        "1000000000000000000", // Ask: 1 ETH
		},
	}

	bids := []Order{
		{
			ID:           2,
			NFTAddress:   "0x0000000000000000000000000000000000000002",
			TokenID:      "1",
			PaymentToken: "0x0000000000000000000000000000000000000003",
			Price:        "2000000000000000000", // Bid: 2 ETH (willing to pay more)
		},
	}

	matches := engine.findMatches(asks, bids)
	require.Len(t, matches, 1, "bid >= ask should match")
	require.Equal(t, asks[0].ID, matches[0].Ask.ID)
	require.Equal(t, bids[0].ID, matches[0].Bid.ID)
}

// TestFetchOrders_FiltersExpired validates that expired orders are excluded from matching.
func TestFetchOrders_FiltersExpired(t *testing.T) {
	engine, redisClient, cleanup := setupTestMatchingEngine(t)
	defer cleanup()

	ctx := context.Background()

	// Insert expired order
	expiredOrder := Order{
		ID:           1,
		Maker:        "0xaaa",
		NFTAddress:   "0x0000000000000000000000000000000000000002",
		TokenID:      "1",
		PaymentToken: "0x0000000000000000000000000000000000000003",
		Price:        "1000000000000000000",
		Expiry:       CustomTime{time.Now().Add(-1 * time.Hour)}, // Expired
		Side:         "ask",
		Hash:         "0xhash1",
	}

	payload, _ := json.Marshal(expiredOrder)
	redisClient.HSet(ctx, "orders:active:ask", expiredOrder.Hash, payload)

	// Fetch orders - should filter out expired
	orders, err := engine.fetchOrders(ctx, "ask")
	require.NoError(t, err)
	require.Len(t, orders, 0)
}

// TestFetchOrders_ReturnsActive validates active orders are retrieved correctly.
func TestFetchOrders_ReturnsActive(t *testing.T) {
	engine, redisClient, cleanup := setupTestMatchingEngine(t)
	defer cleanup()

	ctx := context.Background()

	// Insert active order
	activeOrder := Order{
		ID:           1,
		Maker:        "0xaaa",
		NFTAddress:   "0x0000000000000000000000000000000000000002",
		TokenID:      "1",
		PaymentToken: "0x0000000000000000000000000000000000000003",
		Price:        "1000000000000000000",
		Expiry:       CustomTime{time.Now().Add(1 * time.Hour)}, // Active
		Side:         "ask",
		Hash:         "0xhash1",
	}

	payload, _ := json.Marshal(activeOrder)
	redisClient.HSet(ctx, "orders:active:ask", activeOrder.Hash, payload)

	orders, err := engine.fetchOrders(ctx, "ask")
	require.NoError(t, err)
	require.Len(t, orders, 1)
	require.Equal(t, activeOrder.Maker, orders[0].Maker)
}
