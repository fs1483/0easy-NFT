-- ============================================
-- Oeasy-NFT 数据库初始化脚本
-- ============================================
-- 描述: 创建 Oeasy NFT 订单簿交易平台所需的数据库、表和索引
-- 版本: 1.0.0
-- 创建日期: 2025-10-14
-- ============================================

-- 创建数据库
-- 注意：PostgreSQL 的 CREATE DATABASE 不支持 IF NOT EXISTS
-- 开发环境使用：先删除再创建（简单粗暴，会清空数据）
-- 生产环境需要：使用数据库迁移工具（如 Flyway、Liquibase）
DROP DATABASE IF EXISTS oeasy_nft;
CREATE DATABASE oeasy_nft;

-- 切换到目标数据库
\c oeasy_nft;

-- ============================================
-- 表 1: orders (订单表)
-- ============================================
-- 功能: 存储用户创建的 NFT 买卖订单
-- 关键字段:
--   - maker: 订单创建者地址
--   - nonce: 防重放攻击的唯一标识
--   - side: 订单方向 (ask=卖单, bid=买单)
--   - status: 订单状态 (active=活跃, filled=已成交, cancelled=已取消)
-- ============================================

CREATE TABLE IF NOT EXISTS orders (
    -- 主键
    id BIGSERIAL PRIMARY KEY,
    
    -- 订单创建者信息
    maker VARCHAR(66) NOT NULL,                    -- 订单创建者地址 (0x...)
    
    -- NFT 信息
    nft_address VARCHAR(66) NOT NULL,              -- NFT 合约地址
    token_id NUMERIC(78, 0) NOT NULL,              -- NFT Token ID (支持大整数)
    
    -- 支付信息
    payment_token VARCHAR(66) NOT NULL,            -- 支付代币地址 (如 USDC)
    price NUMERIC(78, 0) NOT NULL,                 -- 价格 (wei 单位)
    
    -- 订单元数据
    expiry TIMESTAMP NOT NULL,                     -- 过期时间
    nonce NUMERIC(78, 0) NOT NULL,                 -- 唯一 nonce (防重放)
    side VARCHAR(4) NOT NULL CHECK (side IN ('ask', 'bid')),  -- 订单方向
    status VARCHAR(16) NOT NULL DEFAULT 'active'   -- 订单状态
        CHECK (status IN ('active', 'filled', 'cancelled')),
    
    -- 签名和哈希
    signature VARCHAR(132) NOT NULL,               -- EIP-712 签名 (0x + 130 字符)
    hash VARCHAR(66) NOT NULL,                     -- 订单哈希
    
    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- 唯一约束: maker + nonce 必须唯一 (防止重复提交)
    CONSTRAINT uk_orders_maker_nonce UNIQUE (maker, nonce)
);

-- 创建索引以优化查询性能
CREATE INDEX idx_orders_maker ON orders(maker);                          -- 按创建者查询
CREATE INDEX idx_orders_nft_token ON orders(nft_address, token_id);     -- 按 NFT 查询
CREATE INDEX idx_orders_side_status ON orders(side, status);            -- 按类型和状态查询
CREATE INDEX idx_orders_status ON orders(status);                       -- 按状态查询
CREATE INDEX idx_orders_expiry ON orders(expiry);                       -- 按过期时间查询
CREATE INDEX idx_orders_hash ON orders(hash);                           -- 按哈希查询
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);          -- 按创建时间倒序查询

-- 添加表注释
COMMENT ON TABLE orders IS 'NFT 订单表 - 存储用户创建的买卖订单';
COMMENT ON COLUMN orders.maker IS '订单创建者钱包地址';
COMMENT ON COLUMN orders.nft_address IS 'NFT 合约地址';
COMMENT ON COLUMN orders.token_id IS 'NFT Token ID';
COMMENT ON COLUMN orders.payment_token IS '支付代币合约地址';
COMMENT ON COLUMN orders.price IS '订单价格（wei 单位）';
COMMENT ON COLUMN orders.expiry IS '订单过期时间';
COMMENT ON COLUMN orders.nonce IS '唯一 nonce，防止重放攻击';
COMMENT ON COLUMN orders.side IS '订单方向: ask=卖单, bid=买单';
COMMENT ON COLUMN orders.status IS '订单状态: active=活跃, filled=已成交, cancelled=已取消';
COMMENT ON COLUMN orders.signature IS 'EIP-712 签名';
COMMENT ON COLUMN orders.hash IS '订单哈希值';

-- ============================================
-- 表 2: trade_events (交易事件表)
-- ============================================
-- 功能: 记录链上执行的交易事件
-- 数据源: 智能合约发出的 TradeExecuted 事件
-- 用途: 审计追踪、数据分析、订单状态同步
-- ============================================

CREATE TABLE IF NOT EXISTS trade_events (
    -- 主键
    id BIGSERIAL PRIMARY KEY,
    
    -- 区块链信息
    transaction_hash VARCHAR(66) NOT NULL,         -- 交易哈希
    log_index INTEGER NOT NULL,                    -- 日志索引（同一交易可能有多个事件）
    block_number BIGINT NOT NULL,                  -- 区块号
    
    -- 交易参与方
    maker VARCHAR(66) NOT NULL,                    -- 卖方地址
    taker VARCHAR(66) NOT NULL,                    -- 买方地址
    
    -- NFT 信息
    nft_address VARCHAR(66) NOT NULL,              -- NFT 合约地址
    token_id NUMERIC(78, 0) NOT NULL,              -- NFT Token ID
    
    -- 支付信息
    payment_token VARCHAR(66) NOT NULL,            -- 支付代币地址
    price NUMERIC(78, 0) NOT NULL,                 -- 成交价格
    
    -- 交易元数据
    side SMALLINT NOT NULL CHECK (side IN (0, 1)), -- 订单方向 (0=Ask, 1=Bid)
    fee NUMERIC(78, 0) NOT NULL,                   -- 平台手续费
    
    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- 唯一约束: transaction_hash + log_index 必须唯一
    -- 防止同一事件被重复处理（去重机制）
    CONSTRAINT uk_trade_events_tx_log UNIQUE (transaction_hash, log_index)
);

-- 创建索引
CREATE INDEX idx_trade_events_block ON trade_events(block_number);      -- 按区块号查询
CREATE INDEX idx_trade_events_maker ON trade_events(maker);             -- 按卖方查询
CREATE INDEX idx_trade_events_taker ON trade_events(taker);             -- 按买方查询
CREATE INDEX idx_trade_events_nft ON trade_events(nft_address);         -- 按 NFT 查询
CREATE INDEX idx_trade_events_created_at ON trade_events(created_at DESC); -- 按时间倒序

-- 添加表注释
COMMENT ON TABLE trade_events IS '交易事件表 - 记录链上执行的 TradeExecuted 事件';
COMMENT ON COLUMN trade_events.transaction_hash IS '交易哈希';
COMMENT ON COLUMN trade_events.log_index IS '事件在交易中的索引';
COMMENT ON COLUMN trade_events.block_number IS '区块号';
COMMENT ON COLUMN trade_events.maker IS '卖方地址';
COMMENT ON COLUMN trade_events.taker IS '买方地址';
COMMENT ON COLUMN trade_events.nft_address IS 'NFT 合约地址';
COMMENT ON COLUMN trade_events.token_id IS 'NFT Token ID';
COMMENT ON COLUMN trade_events.payment_token IS '支付代币地址';
COMMENT ON COLUMN trade_events.price IS '成交价格（wei）';
COMMENT ON COLUMN trade_events.side IS '订单方向: 0=Ask, 1=Bid';
COMMENT ON COLUMN trade_events.fee IS '平台手续费';

-- ============================================
-- 表 3: indexer_status (索引器状态表)
-- ============================================
-- 功能: 记录索引服务的同步状态
-- 用途: 断点续传，确保不漏掉任何区块的事件
-- ============================================

CREATE TABLE IF NOT EXISTS indexer_status (
    -- 主键 (单行表，只有 ID=1 的记录)
    id INTEGER PRIMARY KEY CHECK (id = 1),
    
    -- 同步状态
    last_processed_block BIGINT NOT NULL DEFAULT 0, -- 最后处理的区块号
    
    -- 时间戳
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 插入初始记录
INSERT INTO indexer_status (id, last_processed_block)
VALUES (1, 0)
ON CONFLICT (id) DO NOTHING;

-- 添加表注释
COMMENT ON TABLE indexer_status IS '索引器状态表 - 记录区块同步进度';
COMMENT ON COLUMN indexer_status.last_processed_block IS '最后处理的区块号，用于断点续传';

-- ============================================
-- 触发器: 自动更新 updated_at
-- ============================================

-- 创建更新时间戳的函数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 为 orders 表创建触发器
CREATE TRIGGER trg_orders_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- 为 indexer_status 表创建触发器
CREATE TRIGGER trg_indexer_status_updated_at
BEFORE UPDATE ON indexer_status
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 视图: 活跃订单视图
-- ============================================
-- 功能: 快速查询所有活跃且未过期的订单
-- 用途: 前端订单列表展示、撮合引擎读取
-- ============================================

CREATE OR REPLACE VIEW v_active_orders AS
SELECT 
    id,
    maker,
    nft_address,
    token_id,
    payment_token,
    price,
    expiry,
    nonce,
    side,
    status,
    signature,
    hash,
    created_at,
    updated_at
FROM orders
WHERE status = 'active'
  AND expiry > CURRENT_TIMESTAMP
ORDER BY created_at DESC;

COMMENT ON VIEW v_active_orders IS '活跃订单视图 - 仅显示活跃且未过期的订单';

-- ============================================
-- 视图: 交易统计视图
-- ============================================
-- 功能: 统计交易数据
-- 用途: 数据分析、报表展示
-- ============================================

CREATE OR REPLACE VIEW v_trade_statistics AS
SELECT 
    DATE(created_at) AS trade_date,
    COUNT(*) AS total_trades,
    COUNT(DISTINCT maker) AS unique_sellers,
    COUNT(DISTINCT taker) AS unique_buyers,
    COUNT(DISTINCT nft_address) AS unique_collections,
    SUM(price) AS total_volume,
    AVG(price) AS avg_price,
    MIN(price) AS min_price,
    MAX(price) AS max_price,
    SUM(fee) AS total_fees
FROM trade_events
GROUP BY DATE(created_at)
ORDER BY trade_date DESC;

COMMENT ON VIEW v_trade_statistics IS '交易统计视图 - 按日期统计交易数据';

-- ============================================
-- 函数: 获取用户订单统计
-- ============================================

CREATE OR REPLACE FUNCTION get_user_order_stats(user_address VARCHAR(66))
RETURNS TABLE (
    total_orders BIGINT,
    active_orders BIGINT,
    filled_orders BIGINT,
    cancelled_orders BIGINT,
    total_ask_orders BIGINT,
    total_bid_orders BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(*) AS total_orders,
        COUNT(*) FILTER (WHERE status = 'active') AS active_orders,
        COUNT(*) FILTER (WHERE status = 'filled') AS filled_orders,
        COUNT(*) FILTER (WHERE status = 'cancelled') AS cancelled_orders,
        COUNT(*) FILTER (WHERE side = 'ask') AS total_ask_orders,
        COUNT(*) FILTER (WHERE side = 'bid') AS total_bid_orders
    FROM orders
    WHERE maker = user_address;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_user_order_stats IS '获取指定用户的订单统计信息';

-- ============================================
-- 函数: 清理过期订单（定期任务）
-- ============================================

CREATE OR REPLACE FUNCTION cleanup_expired_orders()
RETURNS INTEGER AS $$
DECLARE
    affected_rows INTEGER;
BEGIN
    -- 将过期的活跃订单标记为取消
    UPDATE orders
    SET status = 'cancelled'
    WHERE status = 'active'
      AND expiry <= CURRENT_TIMESTAMP;
    
    GET DIAGNOSTICS affected_rows = ROW_COUNT;
    
    RETURN affected_rows;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION cleanup_expired_orders IS '清理过期订单 - 将过期的活跃订单标记为已取消';

-- ============================================
-- 初始化完成提示
-- ============================================

DO $$
BEGIN
    RAISE NOTICE '========================================';
    RAISE NOTICE 'Oeasy-NFT 数据库初始化完成！';
    RAISE NOTICE '========================================';
    RAISE NOTICE '已创建:';
    RAISE NOTICE '  - 3 张表: orders, trade_events, indexer_status';
    RAISE NOTICE '  - 2 个视图: v_active_orders, v_trade_statistics';
    RAISE NOTICE '  - 2 个函数: get_user_order_stats, cleanup_expired_orders';
    RAISE NOTICE '  - 多个索引以优化查询性能';
    RAISE NOTICE '========================================';
END $$;

-- ============================================
-- 使用示例
-- ============================================

-- 示例 1: 查询所有活跃订单
-- SELECT * FROM v_active_orders LIMIT 10;

-- 示例 2: 查询用户订单统计
-- SELECT * FROM get_user_order_stats('0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266');

-- 示例 3: 查询今日交易统计
-- SELECT * FROM v_trade_statistics WHERE trade_date = CURRENT_DATE;

-- 示例 4: 清理过期订单
-- SELECT cleanup_expired_orders();

-- 示例 5: 查询某个 NFT 的所有订单
-- SELECT * FROM orders WHERE nft_address = '0x...' AND token_id = 1;

-- 示例 6: 查询某用户的所有成交记录
-- SELECT * FROM trade_events WHERE maker = '0x...' OR taker = '0x...';

