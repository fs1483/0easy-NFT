// Package execution implements the blockchain transaction execution service.
// This service is responsible for submitting matched order pairs to the smart contract,
// managing nonce sequencing, handling gas price estimation, and ensuring transaction reliability.
//
// Architecture:
// - Maintains a hot wallet for signing and submitting transactions
// - Implements robust nonce management to prevent transaction stuck/replacement issues
// - Monitors pending transactions and handles resubmission on failure
// - Updates order status after successful on-chain settlement
//
// Security:
// - Private key loaded from environment variable (should use KMS in production)
// - Only submits pre-validated matched orders from matching engine
package execution

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Oeasy-NFT/services/internal/config"
	"github.com/Oeasy-NFT/services/internal/contracts"
	"github.com/Oeasy-NFT/services/internal/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidTokenID = errors.New("invalid token ID format")
	ErrInvalidPrice   = errors.New("invalid price format")
	ErrInvalidNonce   = errors.New("invalid nonce format")
)

// ExecuteTradeRequest represents a matched order pair ready for on-chain execution.
type ExecuteTradeRequest struct {
	MakerOrder     OrderData `json:"makerOrder" binding:"required"`
	TakerOrder     OrderData `json:"takerOrder" binding:"required"`
	MakerSignature string    `json:"makerSignature" binding:"required"`
}

// OrderData represents the order struct matching the smart contract.
type OrderData struct {
	Maker        string `json:"maker"`
	NFT          string `json:"nft"`
	TokenID      string `json:"tokenId"`
	PaymentToken string `json:"paymentToken"`
	Price        string `json:"price"`
	Expiry       int64  `json:"expiry"`
	Nonce        string `json:"nonce"`
	Side         uint8  `json:"side"`
}

// Service handles submission of transactions to the blockchain.
type Service struct {
	cfg          *config.Config
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	client       *ethclient.Client
	privateKey   *ecdsa.PrivateKey
	fromAddress  common.Address
	marketplace  *contracts.OeasyMarketplace
	engine       *gin.Engine
	nonceMu      sync.Mutex
	pendingNonce uint64
}

// NewService constructs the execution service.
func NewService(cfg *config.Config) (*Service, error) {
	// Connect to Ethereum node
	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		return nil, err
	}

	// Load executor private key from environment
	// TODO: [Security] - Replace with KMS or hardware wallet for production environments
	if cfg.PrivateKeyHex == "" {
		logger.Info("no executor private key provided, execution service will run in read-only mode")
		return &Service{cfg: cfg, client: client}, nil
	}

	// 去除 0x 前缀（如果有）
	privateKeyHex := strings.TrimPrefix(cfg.PrivateKeyHex, "0x")

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, err
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Initialize nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	// Create marketplace contract instance
	marketplaceAddr := common.HexToAddress(cfg.MarketplaceAddr)
	marketplace, err := contracts.NewOeasyMarketplace(marketplaceAddr, client)
	if err != nil {
		return nil, err
	}

	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()
	ginEngine.Use(gin.Recovery())

	svc := &Service{
		cfg:          cfg,
		client:       client,
		privateKey:   privateKey,
		fromAddress:  fromAddress,
		marketplace:  marketplace,
		engine:       ginEngine,
		pendingNonce: nonce,
	}

	svc.registerRoutes()

	logger.Info("execution service initialized",
		"executor", fromAddress.Hex(),
		"marketplace", cfg.MarketplaceAddr,
		"chainId", cfg.ChainID,
	)

	return svc, nil
}

// registerRoutes sets up internal API endpoints for execution requests.
func (s *Service) registerRoutes() {
	api := s.engine.Group("/internal")
	api.POST("/execute", s.handleExecuteTrade)
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "executor": s.fromAddress.Hex()})
	})
}

// handleExecuteTrade receives matched order pairs and submits them on-chain.
func (s *Service) handleExecuteTrade(c *gin.Context) {
	var req ExecuteTradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	txHash, err := s.executeTrade(c.Request.Context(), &req)
	if err != nil {
		logger.Error("failed to execute trade", err,
			"askMaker", req.MakerOrder.Maker,
			"bidMaker", req.TakerOrder.Maker,
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"txHash": txHash.Hex(),
		"status": "submitted",
	})
}

// executeTrade submits the matched order pair to the marketplace contract.
// Returns the transaction hash if successfully broadcast.
func (s *Service) executeTrade(ctx context.Context, req *ExecuteTradeRequest) (*common.Hash, error) {
	logger.Info("executing trade",
		"maker", req.MakerOrder.Maker,
		"taker", req.TakerOrder.Maker,
		"nft", req.MakerOrder.NFT,
		"tokenId", req.MakerOrder.TokenID,
		"price", req.MakerOrder.Price,
	)

	// Convert OrderData to contract Order struct
	makerOrder, err := s.convertToContractOrder(&req.MakerOrder)
	if err != nil {
		return nil, err
	}

	takerOrder, err := s.convertToContractOrder(&req.TakerOrder)
	if err != nil {
		return nil, err
	}

	// Decode maker signature
	makerSig := common.FromHex(req.MakerSignature)

	// 从链上重新获取最新 nonce（防止 nonce 不同步）
	// 企业级最佳实践：每次交易前都从链上获取，避免 nonce 冲突
	currentNonce, err := s.client.PendingNonceAt(ctx, s.fromAddress)
	if err != nil {
		return nil, fmt.Errorf("获取 nonce 失败: %w", err)
	}

	// 更新内存中的 nonce（同步链上状态）
	s.nonceMu.Lock()
	if currentNonce > s.pendingNonce {
		s.pendingNonce = currentNonce
	}
	nonce := s.pendingNonce
	s.pendingNonce++
	s.nonceMu.Unlock()

	// Create transaction options
	chainID := big.NewInt(int64(s.cfg.ChainID))
	auth, err := bind.NewKeyedTransactorWithChainID(s.privateKey, chainID)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Context = ctx

	// TODO: [Gas] - Implement dynamic gas price estimation
	// For MVP, use suggested gas price from node
	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	auth.GasPrice = gasPrice
	auth.GasLimit = 500000 // Conservative gas limit, should estimate in production

	// Call marketplace.executeTrade
	tx, err := s.marketplace.ExecuteTrade(auth, makerOrder, takerOrder, makerSig)
	if err != nil {
		// Rollback nonce on failure
		s.nonceMu.Lock()
		s.pendingNonce--
		s.nonceMu.Unlock()
		return nil, fmt.Errorf("合约调用失败: %w", err)
	}

	txHash := tx.Hash()
	logger.Info("trade transaction submitted",
		"txHash", txHash.Hex(),
		"nonce", nonce,
		"gasPrice", gasPrice.String(),
	)

	return &txHash, nil
}

// convertToContractOrder converts API OrderData to contract IMarketplace.Order struct.
func (s *Service) convertToContractOrder(order *OrderData) (contracts.IMarketplaceOrder, error) {
	tokenID := new(big.Int)
	tokenID, ok := tokenID.SetString(order.TokenID, 10)
	if !ok {
		return contracts.IMarketplaceOrder{}, ErrInvalidTokenID
	}

	price := new(big.Int)
	price, ok = price.SetString(order.Price, 10)
	if !ok {
		return contracts.IMarketplaceOrder{}, ErrInvalidPrice
	}

	nonce := new(big.Int)
	nonce, ok = nonce.SetString(order.Nonce, 10)
	if !ok {
		return contracts.IMarketplaceOrder{}, ErrInvalidNonce
	}

	return contracts.IMarketplaceOrder{
		Maker:        common.HexToAddress(order.Maker),
		Nft:          common.HexToAddress(order.NFT),
		TokenId:      tokenID,
		PaymentToken: common.HexToAddress(order.PaymentToken),
		Price:        price,
		Expiry:       big.NewInt(order.Expiry),
		Nonce:        nonce,
		Side:         order.Side,
	}, nil
}

// getNextNonce returns and increments the pending nonce in a thread-safe manner.
// This prevents nonce collisions when submitting multiple transactions concurrently.
func (s *Service) getNextNonce() uint64 {
	s.nonceMu.Lock()
	defer s.nonceMu.Unlock()
	nonce := s.pendingNonce
	s.pendingNonce++
	return nonce
}

// Run starts the execution service HTTP server to receive execution requests.
func (s *Service) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	srv := &http.Server{
		Addr:              ":" + s.cfg.ExecutionServicePort,
		Handler:           s.engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("execution service http server error", err)
		}
	}()

	logger.Info("execution service listening", "port", s.cfg.ExecutionServicePort)

	<-ctx.Done()

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	return srv.Shutdown(shutdownCtx)
}

// Shutdown gracefully stops the execution service.
func (s *Service) Shutdown() {
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
	if s.client != nil {
		s.client.Close()
	}
	logger.Info("execution service shutdown complete")
}
