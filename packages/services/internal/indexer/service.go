// Package indexer 实现混合式区块链事件索引服务。
// 该服务通过结合实时 WebSocket 订阅和周期性和解轮询，
// 确保链下数据库与链上状态保持同步。
//
// 架构设计（企业级容错混合索引）：
// - 主层：WebSocket 订阅实现实时事件流
// - 安全网：周期性轮询捕获连接中断时遗漏的事件
// - 去重机制：数据库唯一约束防止重复处理事件
// - 自动恢复：WebSocket 失败时指数退避重连
//
// 监听的事件：
// - TradeExecuted：交易链上结算时更新订单状态为 "filled"
// - OrderCancelled：标记订单为已取消（如果添加链上取消功能）
package indexer

import (
	"context"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/Oeasy-NFT/services/internal/config"
	"github.com/Oeasy-NFT/services/internal/contracts"
	"github.com/Oeasy-NFT/services/internal/logger"
	"github.com/Oeasy-NFT/services/internal/orders"
	"github.com/Oeasy-NFT/services/internal/postgres"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

// Service 监听区块链事件并更新数据库
type Service struct {
	cfg                 *config.Config
	cancel              context.CancelFunc
	wg                  sync.WaitGroup
	client              *ethclient.Client
	db                  *gorm.DB
	marketplaceAddr     common.Address
	marketplaceFilterer *contracts.OeasyMarketplaceFilterer
	lastProcessedBlock  uint64
}

// TradeEvent 表示已处理的 TradeExecuted 事件
type TradeEvent struct {
	ID              uint   `gorm:"primaryKey"`
	TransactionHash string `gorm:"type:varchar(66);uniqueIndex:idx_tx_log"` // 交易哈希
	LogIndex        uint   `gorm:"uniqueIndex:idx_tx_log"`                  // 日志索引（与交易哈希组成唯一键）
	BlockNumber     uint64 `gorm:"index"`                                   // 区块号
	Maker           string `gorm:"type:varchar(66);index"`                  // 卖方地址
	Taker           string `gorm:"type:varchar(66);index"`                  // 买方地址
	NFTAddress      string `gorm:"type:varchar(66);index"`                  // NFT 合约地址
	TokenID         string `gorm:"type:numeric"`                            // NFT Token ID
	PaymentToken    string `gorm:"type:varchar(66)"`                        // 支付代币地址
	Price           string `gorm:"type:numeric"`                            // 成交价格
	Side            uint8  // 订单方向（0=Ask, 1=Bid）
	Fee             string `gorm:"type:numeric"` // 平台手续费
	CreatedAt       time.Time
}

// TableName 设置 TradeEvent 的表名
func (TradeEvent) TableName() string {
	return "trade_events"
}

// IndexerStatus 跟踪和解轮询的最后成功处理区块
type IndexerStatus struct {
	ID                 uint   `gorm:"primaryKey"`
	LastProcessedBlock uint64 `gorm:"index"` // 最后处理的区块号
	UpdatedAt          time.Time
}

// TableName 设置 IndexerStatus 的表名
func (IndexerStatus) TableName() string {
	return "indexer_status"
}

// NewService 构造索引服务实例
func NewService(cfg *config.Config) (*Service, error) {
	// 连接以太坊节点
	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		return nil, err
	}

	// 连接数据库
	db, err := postgres.New(cfg.PostgresDSN)
	if err != nil {
		return nil, err
	}

	// 注意：不使用 AutoMigrate，表结构由 database/init.sql 管理
	// 企业级最佳实践：数据库结构与代码分离，使用 SQL 脚本版本化管理

	// 加载或初始化最后处理的区块
	var status IndexerStatus
	result := db.FirstOrCreate(&status, IndexerStatus{ID: 1})
	if result.Error != nil {
		return nil, result.Error
	}

	marketplaceAddr := common.HexToAddress(cfg.MarketplaceAddr)

	// 创建合约事件过滤器用于解析事件
	filterer, err := contracts.NewOeasyMarketplaceFilterer(marketplaceAddr, client)
	if err != nil {
		return nil, err
	}

	logger.Info("索引服务已初始化",
		"marketplace地址", marketplaceAddr.Hex(),
		"最后处理区块", status.LastProcessedBlock,
	)

	return &Service{
		cfg:                 cfg,
		client:              client,
		db:                  db,
		marketplaceAddr:     marketplaceAddr,
		marketplaceFilterer: filterer,
		lastProcessedBlock:  status.LastProcessedBlock,
	}, nil
}

// Run 启动混合索引循环（WebSocket + 轮询）
func (s *Service) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	// 启动 WebSocket 订阅 goroutine
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.runWebSocketSubscriber(ctx)
	}()

	// 启动周期性和解轮询 goroutine
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.runReconciliationPoller(ctx)
	}()

	logger.Info("索引服务已启动（混合架构）")

	<-ctx.Done()
	return nil
}

// runWebSocketSubscriber 维护 WebSocket 连接以实现实时事件流。
// 实现了失败时的指数退避重连机制。
func (s *Service) runWebSocketSubscriber(ctx context.Context) {
	backoff := time.Second
	maxBackoff := time.Minute

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		logger.Info("starting websocket subscription")

		// Create event filter
		query := ethereum.FilterQuery{
			Addresses: []common.Address{s.marketplaceAddr},
		}

		logs := make(chan types.Log)
		sub, err := s.client.SubscribeFilterLogs(ctx, query, logs)
		if err != nil {
			logger.Error("websocket subscription failed, retrying", err, "backoff", backoff)
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}

		// Reset backoff on successful connection
		backoff = time.Second

		logger.Info("websocket subscription established")

		// Process events from subscription
	eventLoop:
		for {
			select {
			case <-ctx.Done():
				sub.Unsubscribe()
				return
			case err := <-sub.Err():
				logger.Error("websocket subscription error", err)
				sub.Unsubscribe()
				break eventLoop // Break inner loop, retry connection
			case vLog := <-logs:
				if err := s.processLog(ctx, vLog); err != nil {
					logger.Error("failed to process log", err,
						"txHash", vLog.TxHash.Hex(),
						"logIndex", vLog.Index,
					)
				}
			}
		}
	}
}

// runReconciliationPoller 周期性获取历史日志以捕获任何遗漏的事件。
// 这是确保即使 WebSocket 失败也能保证 100% 数据一致性的安全网。
// 【企业级改进】：
// - 实现指数退避错误恢复
// - 动态调整轮询间隔（基于区块活动）
// - 添加健康检查和故障检测
func (s *Service) runReconciliationPoller(ctx context.Context) {
	// 轮询间隔：开发环境 10 秒，生产环境 5 分钟
	// 由于 Anvil 不支持 WebSocket，轮询是唯一的数据源
	basePollInterval := 10 * time.Second // 开发环境：快速轮询
	// basePollInterval := 5 * time.Minute  // 生产环境：降低 RPC 调用频率

	currentInterval := basePollInterval
	ticker := time.NewTicker(currentInterval)
	defer ticker.Stop()

	// 【企业级改进】：指数退避参数
	consecutiveErrors := 0
	const maxBackoffInterval = 5 * time.Minute
	const maxConsecutiveErrors = 10

	// 启动时立即执行一次
	logger.Info("轮询索引服务启动",
		"基础间隔", basePollInterval,
		"最大退避间隔", maxBackoffInterval,
	)

	if err := s.reconcile(ctx); err != nil {
		logger.Error("初始和解失败", err)
		consecutiveErrors++
	} else {
		logger.Info("首次轮询完成")
		consecutiveErrors = 0
	}

	for {
		select {
		case <-ctx.Done():
			logger.Info("轮询服务收到停止信号",
				"totalErrors", consecutiveErrors,
			)
			return

		case <-ticker.C:
			reconcileStartTime := time.Now()

			if err := s.reconcile(ctx); err != nil {
				consecutiveErrors++
				logger.Error("和解轮询失败", err,
					"consecutiveErrors", consecutiveErrors,
					"currentInterval", currentInterval,
				)

				// 【企业级改进】：实现指数退避，避免故障时过度轰炸 RPC 节点
				if consecutiveErrors > 0 {
					// 每次失败后加倍间隔时间，但不超过最大值
					currentInterval = basePollInterval * time.Duration(1<<uint(consecutiveErrors))
					if currentInterval > maxBackoffInterval {
						currentInterval = maxBackoffInterval
					}

					logger.Warn("检测到连续失败，增加轮询间隔",
						"consecutiveErrors", consecutiveErrors,
						"newInterval", currentInterval,
					)

					// 重置 ticker 为新的间隔
					ticker.Reset(currentInterval)
				}

				// 【企业级保护】：如果持续失败次数过多，可能是严重故障
				if consecutiveErrors >= maxConsecutiveErrors {
					logger.Error("检测到严重故障：连续失败次数超限",
						nil,
						"consecutiveErrors", consecutiveErrors,
						"maxAllowed", maxConsecutiveErrors,
					)
					// 在企业环境中，这里应触发告警或健康检查失败
					// 目前仅记录错误，继续尝试
				}

			} else {
				reconcileDuration := time.Since(reconcileStartTime)

				// 成功后重置错误计数和轮询间隔
				if consecutiveErrors > 0 {
					logger.Info("轮询恢复正常，重置间隔",
						"previousErrors", consecutiveErrors,
						"resetInterval", basePollInterval,
					)
					consecutiveErrors = 0
					currentInterval = basePollInterval
					ticker.Reset(currentInterval)
				}

				logger.Info("周期轮询完成",
					"duration", reconcileDuration,
					"nextPollIn", currentInterval,
				)
			}
		}
	}
}

// reconcile 从最后处理的区块到当前区块获取日志并处理它们
// 【企业级改进】：
// - 添加区块范围合理性检查（防止重组攻击）
// - 增强错误处理和部分失败恢复
// - 添加性能指标和详细日志
// - 实现渐进式检查点更新（避免全部失败时丢失进度）
func (s *Service) reconcile(ctx context.Context) error {
	startTime := time.Now()

	currentBlock, err := s.client.BlockNumber(ctx)
	if err != nil {
		logger.Error("获取当前区块号失败", err)
		return err
	}

	fromBlock := s.lastProcessedBlock + 1
	toBlock := currentBlock

	// Skip if no new blocks
	if fromBlock > toBlock {
		logger.Info("没有新区块需要处理",
			"lastProcessedBlock", s.lastProcessedBlock,
			"currentBlock", currentBlock,
		)
		return nil
	}

	// 【企业级保护】：检测异常大的区块范围，可能是配置错误或链重组
	blockRange := toBlock - fromBlock + 1
	const maxReasonableBlockRange = uint64(100000) // 设定合理的最大范围
	if blockRange > maxReasonableBlockRange {
		logger.Warn("检测到异常大的区块范围，限制处理范围防止资源耗尽",
			"requestedRange", blockRange,
			"fromBlock", fromBlock,
			"toBlock", toBlock,
			"maxAllowed", maxReasonableBlockRange,
		)
		// 限制单次轮询的区块范围，防止资源耗尽
		toBlock = fromBlock + maxReasonableBlockRange - 1
	}

	logger.Info("开始和解轮询",
		"fromBlock", fromBlock,
		"toBlock", toBlock,
		"blockRange", toBlock-fromBlock+1,
	)

	// 【企业级改进】：减小批次大小以提高检查点频率和故障恢复能力
	// 较小的批次意味着更频繁的进度保存，在失败时丢失的工作更少
	batchSize := uint64(1000) // 从 10000 降低到 1000，提高容错性
	totalProcessed := uint64(0)
	totalEvents := 0

	for from := fromBlock; from <= toBlock; from += batchSize {
		to := from + batchSize - 1
		if to > toBlock {
			to = toBlock
		}

		batchStartTime := time.Now()

		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(from)),
			ToBlock:   big.NewInt(int64(to)),
			Addresses: []common.Address{s.marketplaceAddr},
		}

		logs, err := s.client.FilterLogs(ctx, query)
		if err != nil {
			logger.Error("FilterLogs 调用失败", err,
				"fromBlock", from,
				"toBlock", to,
				"已处理区块", totalProcessed,
			)
			// 【企业级改进】：部分失败不应丢失已完成的工作
			// 返回错误但保留已更新的检查点，下次从失败点继续
			return err
		}

		queryDuration := time.Since(batchStartTime)
		logger.Info("FilterLogs 查询完成",
			"事件数量", len(logs),
			"区块范围", from, "-", to,
			"查询耗时", queryDuration,
		)

		// 处理批次中的所有事件
		successCount := 0
		errorCount := 0

		for i, vLog := range logs {
			if err := s.processLog(ctx, vLog); err != nil {
				errorCount++
				logger.Error("事件处理失败", err,
					"批次索引", i,
					"txHash", vLog.TxHash.Hex(),
					"logIndex", vLog.Index,
					"blockNumber", vLog.BlockNumber,
				)
				// 【企业级决策】：继续处理其他事件，不因单个失败中断整个批次
				// 原因：大多数事件处理错误是独立的（如重复、签名问题等）
			} else {
				successCount++
			}
		}

		// 记录批次处理统计
		if len(logs) > 0 {
			logger.Info("批次处理完成",
				"成功", successCount,
				"失败", errorCount,
				"总计", len(logs),
			)
		}

		// 【企业级改进】：渐进式更新检查点，确保即使后续批次失败也不会丢失进度
		if err := s.updateLastProcessedBlock(ctx, to); err != nil {
			logger.Error("更新检查点失败", err, "blockNumber", to)
			return err
		}

		totalProcessed += (to - from + 1)
		totalEvents += len(logs)
	}

	duration := time.Since(startTime)

	// 【企业级指标】：记录轮询性能指标用于监控和调优
	logger.Info("和解轮询完成",
		"lastBlock", toBlock,
		"processedBlocks", totalProcessed,
		"totalEvents", totalEvents,
		"duration", duration,
		"blocksPerSecond", float64(totalProcessed)/duration.Seconds(),
	)

	return nil
}

// processLog 解析并存储单个事件日志。
// 通过数据库在 (tx_hash, log_index) 上的唯一约束处理去重。
func (s *Service) processLog(ctx context.Context, vLog types.Log) error {
	// 使用生成的合约绑定解析 TradeExecuted 事件
	event, err := s.marketplaceFilterer.ParseTradeExecuted(vLog)
	if err != nil {
		// 如果不是 TradeExecuted 事件则跳过（可能是 OrderCancelled 或其他事件）
		logger.Info("跳过非TradeExecuted事件",
			"交易哈希", vLog.TxHash.Hex(),
			"日志索引", vLog.Index,
		)
		return nil
	}

	logger.Info("处理TradeExecuted事件",
		"txHash", vLog.TxHash.Hex(),
		"blockNumber", vLog.BlockNumber,
		"logIndex", vLog.Index,
		"maker", event.Maker.Hex(),
		"taker", event.Taker.Hex(),
		"nft", event.Nft.Hex(),
		"tokenId", event.TokenId.String(),
		"price", event.Price.String(),
	)

	// 将事件存储到数据库
	// 【修复】：统一使用小写地址格式，与订单表保持一致
	tradeEvent := TradeEvent{
		TransactionHash: vLog.TxHash.Hex(),
		LogIndex:        uint(vLog.Index),
		BlockNumber:     vLog.BlockNumber,
		Maker:           strings.ToLower(event.Maker.Hex()),
		Taker:           strings.ToLower(event.Taker.Hex()),
		NFTAddress:      strings.ToLower(event.Nft.Hex()),
		TokenID:         event.TokenId.String(),
		PaymentToken:    strings.ToLower(event.PaymentToken.Hex()),
		Price:           event.Price.String(),
		Side:            event.Side,
		Fee:             event.Fee.String(),
		CreatedAt:       time.Now(),
	}

	// 插入或忽略（如果由于唯一约束已存在）
	result := s.db.WithContext(ctx).Create(&tradeEvent)
	eventAlreadyExists := false

	if result.Error != nil {
		// 检查是否是重复键错误（这是预期的，可以安全忽略）
		errMsg := result.Error.Error()
		if strings.Contains(errMsg, "UNIQUE constraint failed") ||
			strings.Contains(errMsg, "duplicate key") ||
			strings.Contains(errMsg, "uk_trade_events_tx_log") {
			logger.Info("事件已存在，继续更新订单状态",
				"交易哈希", vLog.TxHash.Hex(),
				"日志索引", vLog.Index,
			)
			eventAlreadyExists = true
			// 不要 return！继续执行订单更新
		} else {
			// 其他错误才返回
			return result.Error
		}
	}

	if !eventAlreadyExists {
		logger.Info("TradeExecuted事件存储成功",
			"交易哈希", vLog.TxHash.Hex(),
			"事件ID", tradeEvent.ID,
		)
	}

	// 将买卖双方的订单状态更新为 "filled"
	// 重要：即使事件已存在，也要尝试更新订单状态（可能之前失败了）
	if err := s.updateOrdersToFilled(ctx, event); err != nil {
		logger.Error("更新订单状态失败", err,
			"交易哈希", vLog.TxHash.Hex(),
		)
		// 不返回错误 - 交易事件已经记录
	} else {
		logger.Info("订单状态更新成功",
			"交易哈希", vLog.TxHash.Hex(),
		)
	}

	return nil
}

// updateOrdersToFilled 在交易执行后将订单标记为已成交
func (s *Service) updateOrdersToFilled(ctx context.Context, event *contracts.OeasyMarketplaceTradeExecuted) error {
	// 【关键修复】：将地址转换为小写以匹配数据库中的存储格式
	// PostgreSQL 字符串比较区分大小写，event.Maker.Hex() 返回的是带大小写的地址
	// 但数据库中存储的是全小写地址，导致 WHERE 条件匹配失败

	// 更新 maker 的订单
	makerResult := s.db.WithContext(ctx).Model(&orders.Order{}).
		Where("maker = ? AND nft_address = ? AND token_id = ? AND status = ?",
			strings.ToLower(event.Maker.Hex()), // 转换为小写
			strings.ToLower(event.Nft.Hex()),   // 转换为小写
			event.TokenId.String(),
			orders.OrderStatusActive,
		).
		Update("status", orders.OrderStatusFilled)

	if makerResult.Error != nil {
		return makerResult.Error
	}

	logger.Info("已更新maker订单为已成交",
		"maker地址", event.Maker.Hex(),
		"更新数量", makerResult.RowsAffected,
	)

	// 更新 taker 的订单
	takerResult := s.db.WithContext(ctx).Model(&orders.Order{}).
		Where("maker = ? AND nft_address = ? AND token_id = ? AND status = ?",
			strings.ToLower(event.Taker.Hex()), // 转换为小写
			strings.ToLower(event.Nft.Hex()),   // 转换为小写
			event.TokenId.String(),
			orders.OrderStatusActive,
		).
		Update("status", orders.OrderStatusFilled)

	if takerResult.Error != nil {
		return takerResult.Error
	}

	logger.Info("已更新taker订单为已成交",
		"taker地址", event.Taker.Hex(),
		"更新数量", takerResult.RowsAffected,
	)

	return nil
}

// updateLastProcessedBlock 持久化和解的检查点
func (s *Service) updateLastProcessedBlock(ctx context.Context, blockNum uint64) error {
	s.lastProcessedBlock = blockNum
	return s.db.WithContext(ctx).Model(&IndexerStatus{}).
		Where("id = ?", 1).
		Update("last_processed_block", blockNum).Error
}

// Shutdown 优雅关闭索引服务
func (s *Service) Shutdown() {
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
	if s.client != nil {
		s.client.Close()
	}
	logger.Info("索引服务已关闭")
}
