# Oeasy-NFT 数据库文档

## 📊 数据库概览

Oeasy-NFT 使用 PostgreSQL 14+ 作为主数据库，存储订单、交易事件和索引状态。

### 数据库信息

- **数据库名**: `oeasy_nft`
- **字符集**: `utf8mb4`
- **排序规则**: `utf8mb4_unicode_ci`
- **表数量**: 3 张
- **视图数量**: 2 个
- **函数数量**: 2 个

---

## 🗂️ 数据表结构

### 1. orders (订单表)

存储用户创建的 NFT 买卖订单。

#### 字段说明

| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | BIGSERIAL | 主键 | PRIMARY KEY |
| maker | VARCHAR(66) | 订单创建者地址 | NOT NULL |
| nft_address | VARCHAR(66) | NFT 合约地址 | NOT NULL |
| token_id | NUMERIC(78,0) | NFT Token ID | NOT NULL |
| payment_token | VARCHAR(66) | 支付代币地址 | NOT NULL |
| price | NUMERIC(78,0) | 价格（wei） | NOT NULL |
| expiry | TIMESTAMP | 过期时间 | NOT NULL |
| nonce | NUMERIC(78,0) | 唯一 nonce | NOT NULL |
| side | VARCHAR(4) | 订单方向 (ask/bid) | NOT NULL, CHECK |
| status | VARCHAR(16) | 订单状态 | NOT NULL, CHECK, DEFAULT 'active' |
| signature | VARCHAR(132) | EIP-712 签名 | NOT NULL |
| hash | VARCHAR(66) | 订单哈希 | NOT NULL |
| created_at | TIMESTAMP | 创建时间 | DEFAULT NOW() |
| updated_at | TIMESTAMP | 更新时间 | DEFAULT NOW() |

#### 唯一约束

- `uk_orders_maker_nonce`: (maker, nonce) - 防止重复提交

#### 索引

- `idx_orders_maker`: maker
- `idx_orders_nft_token`: (nft_address, token_id)
- `idx_orders_side_status`: (side, status)
- `idx_orders_status`: status
- `idx_orders_expiry`: expiry
- `idx_orders_hash`: hash
- `idx_orders_created_at`: created_at DESC

---

### 2. trade_events (交易事件表)

记录链上执行的交易事件。

#### 字段说明

| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | BIGSERIAL | 主键 | PRIMARY KEY |
| transaction_hash | VARCHAR(66) | 交易哈希 | NOT NULL |
| log_index | INTEGER | 日志索引 | NOT NULL |
| block_number | BIGINT | 区块号 | NOT NULL |
| maker | VARCHAR(66) | 卖方地址 | NOT NULL |
| taker | VARCHAR(66) | 买方地址 | NOT NULL |
| nft_address | VARCHAR(66) | NFT 合约地址 | NOT NULL |
| token_id | NUMERIC(78,0) | NFT Token ID | NOT NULL |
| payment_token | VARCHAR(66) | 支付代币地址 | NOT NULL |
| price | NUMERIC(78,0) | 成交价格（wei） | NOT NULL |
| side | SMALLINT | 订单方向 (0/1) | NOT NULL, CHECK |
| fee | NUMERIC(78,0) | 平台手续费 | NOT NULL |
| created_at | TIMESTAMP | 创建时间 | DEFAULT NOW() |

#### 唯一约束

- `uk_trade_events_tx_log`: (transaction_hash, log_index) - 防止重复处理

#### 索引

- `idx_trade_events_block`: block_number
- `idx_trade_events_maker`: maker
- `idx_trade_events_taker`: taker
- `idx_trade_events_nft`: nft_address
- `idx_trade_events_created_at`: created_at DESC

---

### 3. indexer_status (索引器状态表)

记录索引服务的同步状态。

#### 字段说明

| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | INTEGER | 主键（固定为1） | PRIMARY KEY, CHECK (id=1) |
| last_processed_block | BIGINT | 最后处理的区块号 | NOT NULL, DEFAULT 0 |
| updated_at | TIMESTAMP | 更新时间 | DEFAULT NOW() |

**注意**: 这是一个单行表，只有一条记录 (id=1)。

---

## 📈 视图

### 1. v_active_orders (活跃订单视图)

快速查询所有活跃且未过期的订单。

```sql
SELECT * FROM v_active_orders;
```

### 2. v_trade_statistics (交易统计视图)

按日期统计交易数据。

```sql
SELECT * FROM v_trade_statistics;
```

返回字段：
- trade_date: 交易日期
- total_trades: 总交易数
- unique_sellers: 唯一卖家数
- unique_buyers: 唯一买家数
- unique_collections: 唯一NFT集合数
- total_volume: 总成交额
- avg_price: 平均价格
- min_price: 最低价格
- max_price: 最高价格
- total_fees: 总手续费

---

## 🔧 函数

### 1. get_user_order_stats(user_address)

获取指定用户的订单统计信息。

```sql
SELECT * FROM get_user_order_stats('0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266');
```

返回字段：
- total_orders: 总订单数
- active_orders: 活跃订单数
- filled_orders: 已成交订单数
- cancelled_orders: 已取消订单数
- total_ask_orders: 总卖单数
- total_bid_orders: 总买单数

### 2. cleanup_expired_orders()

清理过期订单，将过期的活跃订单标记为已取消。

```sql
SELECT cleanup_expired_orders();
```

返回：受影响的行数

---

## 🚀 快速开始

### 1. 初始化数据库

```bash
# 方式 1: 直接使用 psql
psql -U postgres -f database/init.sql

# 方式 2: 使用 Docker
docker exec -i oeasy-nft-postgres psql -U postgres < database/init.sql
```

### 2. 插入测试数据

```bash
# 方式 1: 直接使用 psql
psql -U postgres -d oeasy_nft -f database/seed.sql

# 方式 2: 使用 Docker
docker exec -i oeasy-nft-postgres psql -U postgres -d oeasy_nft < database/seed.sql
```

### 3. 验证安装

```sql
-- 连接到数据库
\c oeasy_nft

-- 查看所有表
\dt

-- 查看所有视图
\dv

-- 查看所有函数
\df

-- 查询订单数量
SELECT COUNT(*) FROM orders;

-- 查看活跃订单
SELECT * FROM v_active_orders;
```

---

## 📝 常用查询

### 订单相关

```sql
-- 1. 查询所有活跃订单
SELECT * FROM orders WHERE status = 'active';

-- 2. 查询某个 NFT 的所有订单
SELECT * FROM orders 
WHERE nft_address = '0xe7f1725e7734ce288f8367e1bb143e90bb3f0512' 
  AND token_id = 1;

-- 3. 查询某用户的所有订单
SELECT * FROM orders WHERE maker = '0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266';

-- 4. 查询卖单（Ask）
SELECT * FROM orders WHERE side = 'ask' AND status = 'active';

-- 5. 查询买单（Bid）
SELECT * FROM orders WHERE side = 'bid' AND status = 'active';

-- 6. 按价格排序的活跃订单
SELECT * FROM orders 
WHERE status = 'active' 
ORDER BY price ASC;
```

### 交易事件相关

```sql
-- 1. 查询所有交易
SELECT * FROM trade_events ORDER BY created_at DESC;

-- 2. 查询某个 NFT 的交易历史
SELECT * FROM trade_events 
WHERE nft_address = '0xe7f1725e7734ce288f8367e1bb143e90bb3f0512' 
  AND token_id = 1;

-- 3. 查询某用户作为卖方的交易
SELECT * FROM trade_events WHERE maker = '0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266';

-- 4. 查询某用户作为买方的交易
SELECT * FROM trade_events WHERE taker = '0x70997970c51812dc3a010c7d01b50e0d17dc79c8';

-- 5. 查询某个区块的交易
SELECT * FROM trade_events WHERE block_number = 12345;

-- 6. 计算总交易量
SELECT 
    COUNT(*) AS total_trades,
    SUM(price) AS total_volume,
    SUM(fee) AS total_fees
FROM trade_events;
```

### 统计分析

```sql
-- 1. 查看今日交易统计
SELECT * FROM v_trade_statistics WHERE trade_date = CURRENT_DATE;

-- 2. 查看最近7天的交易统计
SELECT * FROM v_trade_statistics 
WHERE trade_date >= CURRENT_DATE - INTERVAL '7 days'
ORDER BY trade_date DESC;

-- 3. 查询最活跃的卖家
SELECT 
    maker,
    COUNT(*) AS trade_count,
    SUM(price) AS total_volume
FROM trade_events
GROUP BY maker
ORDER BY trade_count DESC
LIMIT 10;

-- 4. 查询最热门的 NFT
SELECT 
    nft_address,
    COUNT(*) AS trade_count,
    AVG(price) AS avg_price,
    MAX(price) AS max_price
FROM trade_events
GROUP BY nft_address
ORDER BY trade_count DESC;

-- 5. 按价格区间统计订单
SELECT 
    CASE 
        WHEN price < 50000000 THEN '< 50 USDC'
        WHEN price < 100000000 THEN '50-100 USDC'
        WHEN price < 200000000 THEN '100-200 USDC'
        ELSE '> 200 USDC'
    END AS price_range,
    COUNT(*) AS order_count
FROM orders
GROUP BY price_range
ORDER BY MIN(price);
```

---

## 🛠️ 维护任务

### 定期清理过期订单

建议使用 cron 或 pg_cron 定期执行：

```sql
-- 每小时执行一次
SELECT cleanup_expired_orders();
```

### 数据库备份

```bash
# 备份整个数据库
pg_dump -U postgres oeasy_nft > backup_$(date +%Y%m%d).sql

# 仅备份数据（不含结构）
pg_dump -U postgres -a oeasy_nft > data_backup_$(date +%Y%m%d).sql

# 仅备份结构（不含数据）
pg_dump -U postgres -s oeasy_nft > schema_backup_$(date +%Y%m%d).sql
```

### 数据库恢复

```bash
# 恢复数据库
psql -U postgres -d oeasy_nft < backup.sql
```

### 性能优化

```sql
-- 1. 分析表统计信息
ANALYZE orders;
ANALYZE trade_events;

-- 2. 重建索引
REINDEX TABLE orders;
REINDEX TABLE trade_events;

-- 3. 清理死元组
VACUUM FULL orders;
VACUUM FULL trade_events;

-- 4. 查看表大小
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

---

## 🔐 安全建议

### 1. 创建只读用户

```sql
-- 创建只读用户（用于前端查询）
CREATE USER readonly_user WITH PASSWORD 'secure_password';
GRANT CONNECT ON DATABASE oeasy_nft TO readonly_user;
GRANT USAGE ON SCHEMA public TO readonly_user;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO readonly_user;
```

### 2. 创建应用用户

```sql
-- 创建应用用户（用于后端服务）
CREATE USER app_user WITH PASSWORD 'secure_password';
GRANT CONNECT ON DATABASE oeasy_nft TO app_user;
GRANT USAGE ON SCHEMA public TO app_user;
GRANT SELECT, INSERT, UPDATE ON orders TO app_user;
GRANT SELECT, INSERT ON trade_events TO app_user;
GRANT SELECT, UPDATE ON indexer_status TO app_user;
GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO app_user;
```

### 3. 备份策略

- **每日备份**: 完整数据库备份
- **每周备份**: 长期归档备份
- **实时备份**: 配置 WAL 归档（生产环境）

---

## 📚 参考资料

- [PostgreSQL 官方文档](https://www.postgresql.org/docs/)
- [PostgreSQL 性能优化](https://wiki.postgresql.org/wiki/Performance_Optimization)
- [GORM 文档](https://gorm.io/docs/)

---

**最后更新**: 2025-10-14

