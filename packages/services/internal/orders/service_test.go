// Package orders contains unit tests for the order management service.
// Tests validate order creation with EIP-712 signature verification, order cancellation,
// and proper persistence to both PostgreSQL and Redis data layers.
//
// Test Strategy:
// - Uses in-memory SQLite for database isolation (no external Postgres needed)
// - Uses miniredis for Redis testing without external dependencies
// - Validates EIP-712 typed data signing matches contract behavior exactly
// - Ensures orders are cached in Redis and persisted in database atomically
package orders

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Oeasy-NFT/services/internal/config"
	redisutil "github.com/Oeasy-NFT/services/internal/redis"
	"github.com/alicebob/miniredis/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	mathhex "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestOrderService initializes an in-memory order service instance for testing.
// Uses SQLite in-memory database and miniredis to provide complete isolation.
// Returns the service instance and a cleanup function to tear down resources.
func setupTestOrderService(t *testing.T) (*Service, func()) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	// Initialize in-memory SQLite database for testing isolation
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&Order{}))

	// Start in-memory Redis server for test isolation
	srv, err := miniredis.Run()
	require.NoError(t, err)

	cfg := &config.Config{
		HTTPPort:         "8080",
		OrderServicePort: "8081",
		PostgresDSN:      "",
		RedisAddr:        srv.Addr(),
		MarketplaceAddr:  "0x0000000000000000000000000000000000000001",
		RPCURL:           "http://localhost",
		ChainID:          1,
	}

	redisClient := redisutil.New(cfg.RedisAddr, cfg.RedisPassword)

	// Construct order service with all required dependencies
	service := &Service{
		cfg:         cfg,
		repository:  NewRepository(db),
		engine:      gin.New(),
		chainID:     big.NewInt(int64(cfg.ChainID)),
		marketplace: common.HexToAddress(cfg.MarketplaceAddr),
		typedData:   buildTypedData(cfg),
		redisClient: redisClient,
	}

	service.engine.Use(gin.Recovery())
	service.registerRoutes()

	cleanup := func() {
		srv.Close()
	}
	return service, cleanup
}

// TestCreateOrder_InvalidPayload ensures malformed order requests are rejected with 400 status.
func TestCreateOrder_InvalidPayload(t *testing.T) {
	service, cleanup := setupTestOrderService(t)
	defer cleanup()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/orders", strings.NewReader(`{"maker":"bad"}`))
	service.engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

// TestCreateOrder_Success validates end-to-end order creation flow:
// 1. Client generates EIP-712 signature matching contract domain
// 2. Server validates signature against maker address
// 3. Order is persisted to PostgreSQL with status "active"
// 4. Order is cached in Redis for fast matching engine access
func TestCreateOrder_Success(t *testing.T) {
	service, cleanup := setupTestOrderService(t)
	defer cleanup()

	// Generate test wallet and sign order
	makerKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	body, nonce := buildSignedOrderRequest(t, service, makerKey)

	// Submit order creation request
	w := httptest.NewRecorder()
	httpReq := httptest.NewRequest(http.MethodPost, "/api/orders", strings.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	service.engine.ServeHTTP(w, httpReq)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp orderResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	// Verify order was persisted to database
	stored, err := service.repository.FindByID(context.Background(), resp.ID)
	require.NoError(t, err)
	require.Equal(t, nonce.String(), stored.Nonce)
	require.Equal(t, OrderStatusActive, stored.Status)

	// Verify order was cached in Redis for matching engine
	redisKey := "orders:active:" + stored.Side
	cacheVal, err := service.redisClient.HGet(context.Background(), redisKey, stored.Hash).Result()
	require.NoError(t, err)
	require.NotEmpty(t, cacheVal)
}

// TestCancelOrder_Success validates order cancellation workflow:
// 1. Create an active order first
// 2. Submit signed cancel request with matching maker and nonce
// 3. Verify order status updated to "cancelled" in database
// 4. Verify order removed from Redis active order book
func TestCancelOrder_Success(t *testing.T) {
	service, cleanup := setupTestOrderService(t)
	defer cleanup()

	makerKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	// Create order first
	body, nonce := buildSignedOrderRequest(t, service, makerKey)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/orders", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	service.engine.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var created orderResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))

	// Build signed cancel request
	cancelBody := buildSignedCancelRequest(t, service, makerKey, nonce)

	w = httptest.NewRecorder()
	cancelReq := httptest.NewRequest(
		http.MethodPost,
		"/api/orders/"+strconv.Itoa(int(created.ID))+"/cancel",
		strings.NewReader(cancelBody),
	)
	cancelReq.Header.Set("Content-Type", "application/json")
	service.engine.ServeHTTP(w, cancelReq)

	require.Equal(t, http.StatusOK, w.Code)

	// Verify status updated in database
	stored, err := service.repository.FindByID(context.Background(), created.ID)
	require.NoError(t, err)
	require.Equal(t, OrderStatusCancelled, stored.Status)

	// Verify removed from Redis order book
	redisKey := "orders:active:" + stored.Side
	_, err = service.redisClient.HGet(context.Background(), redisKey, stored.Hash).Result()
	require.Error(t, err) // Should not exist
}

// bigIntFromUint is a helper to convert uint64 to *big.Int.
func bigIntFromUint(value uint64) *big.Int {
	return new(big.Int).SetUint64(value)
}

// buildTypedData constructs the EIP-712 typed data structure matching contract configuration.
// This must exactly match the domain and types defined in OeasyMarketplace.sol.
func buildTypedData(cfg *config.Config) apitypes.TypedData {
	marketplace := common.HexToAddress(cfg.MarketplaceAddr).Hex()
	return apitypes.TypedData{
		Types: map[string][]apitypes.Type{
			"EIP712Domain": {
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"Order": {
				{Name: "maker", Type: "address"},
				{Name: "nft", Type: "address"},
				{Name: "tokenId", Type: "uint256"},
				{Name: "paymentToken", Type: "address"},
				{Name: "price", Type: "uint256"},
				{Name: "expiry", Type: "uint256"},
				{Name: "nonce", Type: "uint256"},
				{Name: "side", Type: "uint8"},
			},
			"Cancel": {
				{Name: "maker", Type: "address"},
				{Name: "nonce", Type: "uint256"},
			},
		},
		PrimaryType: "Order",
		Domain: apitypes.TypedDataDomain{
			Name:              "Oeasy Marketplace",
			Version:           "1",
			ChainId:           mathhex.NewHexOrDecimal256(int64(cfg.ChainID)),
			VerifyingContract: marketplace,
		},
	}
}

// buildSignedOrderRequest creates a properly signed order request payload for testing.
// Generates EIP-712 signature using provided private key, matching contract expectations.
// Returns JSON-encoded request body and the nonce used (for subsequent cancellation tests).
func buildSignedOrderRequest(t *testing.T, service *Service, makerKey *ecdsa.PrivateKey) (string, *big.Int) {
	makerAddr := crypto.PubkeyToAddress(makerKey.PublicKey)
	expiry := time.Now().Add(time.Hour).Unix()
	nonce := big.NewInt(1)
	tokenID := big.NewInt(1)
	price := big.NewInt(1_000_000_000_000_000_000) // 1 ether

	nftAddr := common.HexToAddress("0x0000000000000000000000000000000000000002")
	paymentAddr := common.HexToAddress("0x0000000000000000000000000000000000000003")

	// Construct EIP-712 message with all fields lowercase hex for addresses
	message := apitypes.TypedDataMessage{
		"maker":        strings.ToLower(makerAddr.Hex()),
		"nft":          strings.ToLower(nftAddr.Hex()),
		"tokenId":      tokenID,
		"paymentToken": strings.ToLower(paymentAddr.Hex()),
		"price":        price,
		"expiry":       big.NewInt(expiry),
		"nonce":        nonce,
		"side":         big.NewInt(0), // Ask = 0
	}

	td := service.typedData
	td.PrimaryType = "Order"
	td.Message = message

	digest, _, err := apitypes.TypedDataAndHash(td)
	require.NoError(t, err)

	// Sign the digest with maker's private key
	sig, err := crypto.Sign(digest, makerKey)
	require.NoError(t, err)

	// Build HTTP request payload
	req := createOrderRequest{
		Maker:        strings.ToLower(makerAddr.Hex()),
		NFTAddress:   strings.ToLower(nftAddr.Hex()),
		TokenID:      tokenID.String(),
		PaymentToken: strings.ToLower(paymentAddr.Hex()),
		Price:        price.String(),
		Expiry:       expiry,
		Nonce:        nonce.String(),
		Side:         "ask",
		Signature:    strings.ToLower(hexutil.Encode(sig)),
	}

	body, err := json.Marshal(req)
	require.NoError(t, err)
	return string(body), nonce
}

// buildSignedCancelRequest creates a signed cancel request for testing order cancellation.
// Uses EIP-712 Cancel message type with maker address and nonce.
func buildSignedCancelRequest(t *testing.T, service *Service, makerKey *ecdsa.PrivateKey, nonce *big.Int) string {
	makerAddr := crypto.PubkeyToAddress(makerKey.PublicKey)

	// Construct Cancel typed message
	message := apitypes.TypedDataMessage{
		"maker": strings.ToLower(makerAddr.Hex()),
		"nonce": nonce,
	}

	td := service.typedData
	td.PrimaryType = "Cancel"
	td.Message = message

	digest, _, err := apitypes.TypedDataAndHash(td)
	require.NoError(t, err)

	sig, err := crypto.Sign(digest, makerKey)
	require.NoError(t, err)

	payload := cancelOrderRequest{
		Maker:     strings.ToLower(makerAddr.Hex()),
		Nonce:     nonce.String(),
		Signature: strings.ToLower(hexutil.Encode(sig)),
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)
	return string(body)
}
