# 📊 如何查看服务日志

## 🎯 日志位置

所有日志文件都在：
```
~/blockchain-project/web3-knowledge/Oeasy-NFT/packages/services/logs/
```

---

## 📋 查看日志的方法

### 方法 1: 使用 tail 命令（实时查看）⭐ 推荐

**查看单个服务日志**:
```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT/packages/services

# 订单服务
tail -f logs/order-service.log

# 撮合引擎
tail -f logs/matching-engine.log

# 执行服务
tail -f logs/execution-service.log

# 索引服务
tail -f logs/indexer.log
```

**查看所有服务日志**:
```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT/packages/services
tail -f logs/*.log
```

**停止查看**: 按 `Ctrl+C`

---

### 方法 2: 使用 cat 查看历史日志

```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT/packages/services

# 查看最后 50 行
tail -50 logs/matching-engine.log

# 查看完整日志
cat logs/order-service.log

# 搜索特定内容
grep "ERROR" logs/execution-service.log
grep "匹配" logs/matching-engine.log
```

---

### 方法 3: 使用查看日志脚本

```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT
./scripts/view-logs.sh

# 会出现菜单让你选择
```

---

## 🎯 快速诊断命令

### 检查是否有错误

```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT/packages/services

# 查找所有错误
grep -i "error" logs/*.log | tail -20

# 查找匹配日志
grep "匹配" logs/matching-engine.log

# 查找交易提交
grep "transaction submitted" logs/execution-service.log
```

---

## 📁 日志文件说明

| 文件 | 内容 | 关键词 |
|------|------|--------|
| `order-service.log` | 订单创建、查询 | "订单创建", "http_request" |
| `matching-engine.log` | 撮合匹配 | "匹配订单对", "提交执行" |
| `execution-service.log` | 交易提交 | "executing trade", "txHash" |
| `indexer.log` | 事件监听 | "TradeExecuted", "更新订单" |

---

## 🛠️ 常用命令速查

```bash
# 必须先进入这个目录！
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT/packages/services

# 实时查看撮合日志
tail -f logs/matching-engine.log

# 查看最近的错误
grep ERROR logs/*.log | tail -10

# 查看最近 100 行
tail -100 logs/execution-service.log

# 查看所有服务状态
ps aux | grep "go run cmd"
```

---

## 🎯 当前日志位置

**完整路径**:
```
/Users/shuangfan/blockchain-project/web3-knowledge/Oeasy-NFT/packages/services/logs/

├── order-service.log      (571 字节)
├── matching-engine.log    (520 字节)
├── execution-service.log  (1345 字节)
└── indexer.log           (16541 字节)
```

**快速查看**:
```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT/packages/services
ls -lh logs/
```

---

## 💡 使用建议

**平时开发**:
```bash
# 启动服务
./scripts/start-all-services.sh

# 开一个终端实时查看撮合日志
cd packages/services
tail -f logs/matching-engine.log

# 需要时查看其他日志
tail -f logs/execution-service.log
```

**这样只需要 1-2 个终端，而不是 4-5 个！** ✅

---

**现在试试查看撮合引擎日志**:
```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT/packages/services
tail -f logs/matching-engine.log
```

🚀
