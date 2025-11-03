// Package matching 实现 Oeasy NFT 市场的核心订单撮合引擎。
// 引擎持续扫描 Redis 订单簿，寻找兼容的 ask/bid 订单对，
// 并将匹配的订单转发给执行服务进行链上结算。
//
// 架构设计：
// - 从 Redis 读取活跃订单（由订单服务发布）
// - 实现价格-时间优先撮合算法
// - 发现匹配时通知执行服务
// - 执行成功后更新订单状态
package matching

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/Oeasy-NFT/services/internal/config"
	"github.com/Oeasy-NFT/services/internal/logger"
	redisutil "github.com/Oeasy-NFT/services/internal/redis"
	"github.com/redis/go-redis/v9"
)

// Order 表示从 Redis 获取用于撮合的订单
type Order struct {
	ID           uint       `json:"id"`
	Maker        string     `json:"maker"`
	NFTAddress   string     `json:"nftAddress"`
	TokenID      FlexString `json:"tokenId"` // 兼容字符串和数字
	PaymentToken string     `json:"paymentToken"`
	Price        FlexString `json:"price"`  // 兼容字符串和数字
	Expiry       CustomTime `json:"expiry"` // 自定义时间类型
	Nonce        FlexString `json:"nonce"`  // 兼容字符串和数字
	Side         string     `json:"side"`
	Status       string     `json:"status"`
	Signature    string     `json:"signature"`
	Hash         string     `json:"hash"`
}

// FlexString 灵活字符串类型，兼容 JSON 中的字符串和数字
type FlexString string

// UnmarshalJSON 自定义解析，支持字符串和数字
func (fs *FlexString) UnmarshalJSON(b []byte) error {
	var s string
	// 尝试作为字符串解析
	if err := json.Unmarshal(b, &s); err == nil {
		*fs = FlexString(s)
		return nil
	}

	// 尝试作为数字解析
	var n json.Number
	if err := json.Unmarshal(b, &n); err == nil {
		*fs = FlexString(n.String())
		return nil
	}

	return fmt.Errorf("无法解析为字符串或数字")
}

// String 转换为字符串
func (fs FlexString) String() string {
	return string(fs)
}

// CustomTime 自定义时间类型，支持多种格式解析
type CustomTime struct {
	time.Time
}

// UnmarshalJSON 自定义 JSON 解析，支持多种时间格式
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1] // 去掉引号

	// 尝试多种时间格式
	formats := []string{
		time.RFC3339,                 // 2006-01-02T15:04:05Z07:00
		"2006-01-02T15:04:05.999999", // PostgreSQL 格式（无时区）
		"2006-01-02T15:04:05",
		time.RFC3339Nano,
	}

	for _, format := range formats {
		t, err := time.Parse(format, s)
		if err == nil {
			ct.Time = t
			return nil
		}
	}

	return fmt.Errorf("无法解析时间: %s", s)
}

// MatchPair 表示成功匹配的 ask 和 bid 订单对
type MatchPair struct {
	Ask Order
	Bid Order
}

// ExecuteTradeRequest 表示提交给执行服务的交易请求
type ExecuteTradeRequest struct {
	MakerOrder     OrderData `json:"makerOrder"`
	TakerOrder     OrderData `json:"takerOrder"`
	MakerSignature string    `json:"makerSignature"`
}

// OrderData 表示符合智能合约格式的订单数据
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

// ExecuteTradeResponse 表示执行服务的响应
type ExecuteTradeResponse struct {
	TxHash string `json:"txHash"`
	Status string `json:"status"`
}

// Engine 表示撮合引擎的运行时实例
type Engine struct {
	cfg         *config.Config
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	redisClient *redis.Client
}

// NewEngine 创建新的撮合引擎实例
func NewEngine(cfg *config.Config) (*Engine, error) {
	redisClient := redisutil.New(cfg.RedisAddr, cfg.RedisPassword)
	if err := redisutil.Ping(context.Background(), redisClient); err != nil {
		return nil, err
	}

	return &Engine{
		cfg:         cfg,
		redisClient: redisClient,
	}, nil
}

// Run 启动撮合引擎循环，持续扫描兼容的订单
func (e *Engine) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	e.cancel = cancel

	logger.Info("matching engine started", "interval", "5s")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := e.matchOrders(ctx); err != nil {
				logger.Error("matching cycle failed", err)
			}
		}
	}
}

// matchOrders 从 Redis 获取活跃的 ask 和 bid 订单并尝试撮合
func (e *Engine) matchOrders(ctx context.Context) error {
	// 从 Redis 获取所有活跃的卖单（ask）
	asks, err := e.fetchOrders(ctx, "ask")
	if err != nil {
		return err
	}

	// 从 Redis 获取所有活跃的买单（bid）
	bids, err := e.fetchOrders(ctx, "bid")
	if err != nil {
		return err
	}

	// 寻找兼容的订单匹配
	matches := e.findMatches(asks, bids)

	// 将匹配的订单对发送到执行服务
	if len(matches) > 0 {
		logger.Info("发现订单匹配", "数量", len(matches))
		for _, match := range matches {
			logger.Info("匹配订单对",
				"NFT地址", match.Ask.NFTAddress,
				"TokenID", match.Ask.TokenID,
				"价格", match.Ask.Price,
				"卖方", match.Ask.Maker,
				"买方", match.Bid.Maker,
			)

			// 提交到执行服务进行链上结算
			if err := e.submitToExecution(ctx, match); err != nil {
				logger.Error("提交执行失败", err,
					"NFT", match.Ask.NFTAddress,
					"TokenID", match.Ask.TokenID,
				)
				// 继续处理其他匹配，不因单个失败而中断
				continue
			}

			logger.Info("订单对已提交执行",
				"卖方", match.Ask.Maker,
				"买方", match.Bid.Maker,
			)

			// 从 Redis 删除已提交的订单，避免重复撮合
			// 注意：订单状态最终由索引服务更新，这里只是从缓存中移除
			askKey := "orders:active:ask"
			bidKey := "orders:active:bid"
			e.redisClient.HDel(ctx, askKey, match.Ask.Hash)
			e.redisClient.HDel(ctx, bidKey, match.Bid.Hash)

			logger.Info("已从 Redis 缓存中移除匹配的订单",
				"askHash", match.Ask.Hash[:10]+"...",
				"bidHash", match.Bid.Hash[:10]+"...",
			)
		}
	}

	return nil
}

// fetchOrders 从 Redis 检索指定方向的所有活跃订单
func (e *Engine) fetchOrders(ctx context.Context, side string) ([]Order, error) {
	key := "orders:active:" + side
	ordersMap, err := e.redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	orders := make([]Order, 0, len(ordersMap))
	for _, payload := range ordersMap {
		var ord Order
		if err := json.Unmarshal([]byte(payload), &ord); err != nil {
			logger.Error("failed to unmarshal order from redis", err, "payload", payload)
			continue
		}

		// 跳过已过期订单
		if ord.Expiry.Time.Before(time.Now()) {
			continue
		}

		orders = append(orders, ord)
	}

	return orders, nil
}

// findMatches 实现简单的撮合算法，将 ask 和 bid 配对
// 撮合条件：
// - 相同的 NFT 合集和 token ID
// - 相同的支付代币
// - Ask 价格 <= Bid 价格（买方愿意支付至少卖方要价）
// - 两个订单都未过期
//
// TODO: [性能优化] - 实现价格-时间优先队列，达到 O(log n) 撮合复杂度
// TODO: [可扩展性] - 支持部分成交和多数量撮合
func (e *Engine) findMatches(asks []Order, bids []Order) []MatchPair {
	matches := make([]MatchPair, 0)

	// MVP 采用简单的 O(n*m) 撮合算法
	// 生产环境应使用排序价格队列实现 O(log n) 复杂度
	for _, ask := range asks {
		for _, bid := range bids {
			if e.isMatch(ask, bid) {
				matches = append(matches, MatchPair{Ask: ask, Bid: bid})
				break // 在这个简单实现中，每个 ask 最多匹配一个 bid
			}
		}
	}

	return matches
}

// isMatch 判断 ask 和 bid 订单是否兼容可执行
func (e *Engine) isMatch(ask Order, bid Order) bool {
	// 必须是相同的 NFT
	if ask.NFTAddress != bid.NFTAddress || ask.TokenID.String() != bid.TokenID.String() {
		return false
	}

	// 必须使用相同的支付代币
	if ask.PaymentToken != bid.PaymentToken {
		return false
	}

	// Bid 价格必须达到或超过 ask 价格（买方愿意支付 >= 卖方要价）
	// 这是标准订单簿撮合规则: bid >= ask
	// MVP 中我们按 ask 价格执行（maker 获得其要求的价格）
	// TODO: [撮合优化] - 实现价格改善逻辑或中间价执行
	askPrice, askOk := new(big.Int).SetString(ask.Price.String(), 10)
	bidPrice, bidOk := new(big.Int).SetString(bid.Price.String(), 10)
	if !askOk || !bidOk {
		return false
	}

	// Bid 必须 >= ask 才能形成有效匹配
	if bidPrice.Cmp(askPrice) < 0 {
		return false
	}

	return true
}

// submitToExecution 将匹配的订单对提交给执行服务进行链上结算
func (e *Engine) submitToExecution(ctx context.Context, match MatchPair) error {
	// 构建执行请求
	// 在标准订单簿中：ask 是 maker（挂单方），bid 是 taker（吃单方）
	req := ExecuteTradeRequest{
		MakerOrder: OrderData{
			Maker:        match.Ask.Maker,
			NFT:          match.Ask.NFTAddress,
			TokenID:      match.Ask.TokenID.String(),
			PaymentToken: match.Ask.PaymentToken,
			Price:        match.Ask.Price.String(),
			Expiry:       match.Ask.Expiry.Time.Unix(),
			Nonce:        match.Ask.Nonce.String(),
			Side:         0, // 0 = Ask (卖单)
		},
		TakerOrder: OrderData{
			Maker:        match.Bid.Maker,
			NFT:          match.Bid.NFTAddress,
			TokenID:      match.Bid.TokenID.String(),
			PaymentToken: match.Bid.PaymentToken,
			Price:        match.Bid.Price.String(),
			Expiry:       match.Bid.Expiry.Time.Unix(),
			Nonce:        match.Bid.Nonce.String(),
			Side:         1, // 1 = Bid (买单)
		},
		MakerSignature: match.Ask.Signature,
	}

	// 序列化请求
	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("序列化执行请求失败: %w", err)
	}

	// 构建执行服务 URL
	executionURL := fmt.Sprintf("http://localhost:%s/internal/execute", e.cfg.ExecutionServicePort)

	// 创建 HTTP 请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", executionURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("创建HTTP请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("调用执行服务失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("执行服务返回错误: HTTP %d", resp.StatusCode)
	}

	// 解析响应
	var execResp ExecuteTradeResponse
	if err := json.NewDecoder(resp.Body).Decode(&execResp); err != nil {
		return fmt.Errorf("解析执行响应失败: %w", err)
	}

	logger.Info("交易已提交到链上",
		"交易哈希", execResp.TxHash,
		"状态", execResp.Status,
	)

	return nil
}

// Shutdown 优雅关闭撮合引擎
func (e *Engine) Shutdown() {
	if e.cancel != nil {
		e.cancel()
	}
	e.wg.Wait()
	logger.Info("撮合引擎已关闭")
}
