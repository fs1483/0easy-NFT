package orders

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/Oeasy-NFT/services/internal/config"
	"github.com/Oeasy-NFT/services/internal/logger"
	"github.com/Oeasy-NFT/services/internal/postgres"
	redisutil "github.com/Oeasy-NFT/services/internal/redis"
	"github.com/ethereum/go-ethereum/common"
	mathhex "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
)

// Service represents the order management microservice runtime.
type Service struct {
	cfg         *config.Config
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	repository  *Repository
	engine      *gin.Engine
	chainID     *big.Int
	marketplace common.Address
	typedData   apitypes.TypedData
	redisClient *redis.Client
}

// NewService constructs the order service wiring data stores.
func NewService(cfg *config.Config) (*Service, error) {
	db, err := postgres.New(cfg.PostgresDSN)
	if err != nil {
		return nil, err
	}
	repo := NewRepository(db)

	// 注意：不使用 AutoMigrate，表结构由 database/init.sql 管理
	// 这是企业级最佳实践：数据库结构与代码分离
	// 优点：
	// - 表结构版本化（Git 管理）
	// - 支持复杂的索引、视图、函数
	// - 避免 ORM 自动修改表结构的风险

	redisClient := redisutil.New(cfg.RedisAddr, cfg.RedisPassword)
	if err := redisutil.Ping(context.Background(), redisClient); err != nil {
		return nil, err
	}

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(logger.HTTPLogger())

	// 配置 CORS 允许前端跨域访问
	engine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	chainID := new(big.Int).SetUint64(cfg.ChainID)
	marketplaceAddr := common.HexToAddress(cfg.MarketplaceAddr)

	typedData := apitypes.TypedData{
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
			VerifyingContract: marketplaceAddr.Hex(),
		},
	}

	service := &Service{
		cfg:         cfg,
		repository:  repo,
		engine:      engine,
		chainID:     chainID,
		marketplace: marketplaceAddr,
		typedData:   typedData,
		redisClient: redisClient,
	}
	service.registerRoutes()

	return service, nil
}

// registerRoutes attaches handlers to the engine.
func (s *Service) registerRoutes() {
	api := s.engine.Group("/api")
	RegisterRoutes(api, s)
}

func (s *Service) cacheOrder(ctx context.Context, ord *Order) error {
	payload, err := json.Marshal(ord)
	if err != nil {
		return err
	}
	key := "orders:active:" + ord.Side
	return s.redisClient.HSet(ctx, key, ord.Hash, payload).Err()
}

// Run starts background routines (event listeners, HTTP handlers, etc.).
func (s *Service) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	srv := &http.Server{
		Addr:              ":" + s.cfg.OrderServicePort,
		Handler:           s.engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("order service http server error", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	return srv.Shutdown(shutdownCtx)
}

// Shutdown gracefully stops the service.
func (s *Service) Shutdown() {
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
}
