# 🎯 端到端测试完整指南

> 本文档提供从环境准备到完整业务流程测试的详细步骤

---

## 📋 目录

1. [功能清单](#功能清单)
2. [快速测试（5分钟）](#快速测试5分钟)
3. [完整测试流程](#完整测试流程)
4. [验证与调试](#验证与调试)
5. [常见问题](#常见问题)

---

## 功能清单

### ✅ 业务闭环完备性

#### 前端功能
- ✅ 铸造 NFT（自动授权）
- ✅ 获取测试 USDC（仅铸造）
- ✅ 创建卖单（自动检查并授权 NFT）
- ✅ 创建买单（自动检查并授权 USDC，用户选择额度）
- ✅ 查看市场订单
- ✅ 我的订单管理
- ✅ 取消订单
- ✅ 交易历史
- ✅ 点击"出价"快捷创建买单

#### 后端服务
- ✅ 订单服务（API）
- ✅ 撮合引擎（自动匹配）
- ✅ 执行服务（提交链上交易）
- ✅ 索引服务（监听事件，更新状态）

#### 智能合约
- ✅ OeasyMarketplace（交易执行）
- ✅ OeasyNFT（NFT 合约）
- ✅ MockUSDC（测试代币）

**所有功能完备！** ✅

---

## 快速测试（5分钟）

> 适用于快速验证系统是否正常工作

### 前置条件

```bash
# 确保基础服务运行
✅ Anvil (端口 8545)
✅ PostgreSQL + Redis (Docker)
✅ 合约已部署
```

### 步骤 1: 检查服务状态

```bash
# 检查服务
lsof -i :8081  # 订单服务
lsof -i :8083  # 执行服务
ps aux | grep indexer | grep -v grep  # 索引服务
```

**如果服务未运行，启动它们**:

```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT/packages/services

# 终端 1: 订单服务
go run cmd/order-service/main.go

# 终端 2: 执行服务
go run cmd/execution-service/main.go

# 终端 3: 索引服务
go run cmd/indexer/main.go

# 终端 4: 撮合引擎
go run cmd/matching-engine/main.go
```

### 步骤 2: 一键系统检查

创建并运行检查脚本：

```bash
#!/bin/bash
# 保存为 check-system.sh

echo "🔍 系统状态检查..."
echo ""

# 服务检查
echo "📊 服务状态:"
lsof -i :8081 > /dev/null && echo "✅ 订单服务" || echo "❌ 订单服务"
lsof -i :8083 > /dev/null && echo "✅ 执行服务" || echo "❌ 执行服务"
ps aux | grep "indexer" | grep -v grep > /dev/null && echo "✅ 索引服务" || echo "❌ 索引服务"
ps aux | grep "matching-engine" | grep -v grep > /dev/null && echo "✅ 撮合引擎" || echo "❌ 撮合引擎"

# Redis 检查
echo ""
echo "📦 Redis 订单:"
echo "  Ask: $(docker exec oeasy-nft-redis redis-cli HLEN orders:active:ask 2>/dev/null || echo 0)"
echo "  Bid: $(docker exec oeasy-nft-redis redis-cli HLEN orders:active:bid 2>/dev/null || echo 0)"

# 数据库检查
echo ""
echo "💾 数据库订单:"
echo "  Active: $(docker exec oeasy-nft-postgres psql -U postgres -d oeasy_nft -t -c "SELECT COUNT(*) FROM orders WHERE status='active';" 2>/dev/null | tr -d ' ')"
echo "  Filled: $(docker exec oeasy-nft-postgres psql -U postgres -d oeasy_nft -t -c "SELECT COUNT(*) FROM orders WHERE status='filled';" 2>/dev/null | tr -d ' ')"

# 链上资产检查
echo ""
echo "🎨 链上资产:"
NFT_OWNER=$(cast call 0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0 "ownerOf(uint256)" 2 --rpc-url http://localhost:8545 2>&1 | grep "0x")
if [ -n "$NFT_OWNER" ]; then
    echo "  ✅ NFT #2 存在"
else
    echo "  ❌ NFT #2 不存在 - 需要准备测试数据"
fi

echo ""
echo "✅ 检查完成"
```

运行检查：
```bash
chmod +x check-system.sh
./check-system.sh
```

### 步骤 3: 准备测试数据（如需要）

如果链上资产不存在，运行准备脚本：

```bash
./scripts/setup-test-environment.sh
```

### 步骤 4: 同步订单到 Redis

```bash
# 使用缓存脚本
/tmp/cache_orders.sh

# 验证
docker exec oeasy-nft-redis redis-cli HLEN "orders:active:ask"
docker exec oeasy-nft-redis redis-cli HLEN "orders:active:bid"
# 应该都显示 2
```

### 步骤 5: 观察撮合

撮合引擎应该在 5-10 秒内自动匹配订单，观察日志：

```
✅ 发现订单匹配 数量=1
✅ 匹配订单对 NFT地址=0x9fe... TokenID=2
✅ 交易已提交到链上 txHash=0x45b3...
✅ 订单对已提交执行
✅ 已从 Redis 缓存中移除
```

### 步骤 6: 验证结果

```bash
# 检查交易状态
cast receipt <交易哈希> --rpc-url http://localhost:8545 | grep status
# 应该显示 "status 1 (success)"

# 检查订单状态
docker exec oeasy-nft-postgres psql -U postgres -d oeasy_nft \
  -c "SELECT id, status FROM orders WHERE id IN (20, 22);"
# 应该显示 status = 'filled'
```

---

## 完整测试流程

> 适用于完整的用户体验测试，包括前端交互

### 前置准备（一次性）

#### 1. 确保基础服务运行

```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT

# 检查 Anvil
lsof -i :8545

# 检查数据库
docker ps | grep postgres

# 启动所有后端服务
./scripts/start-all-services.sh

# 验证服务
lsof -i :8081,8083
```

#### 2. MetaMask 导入测试账户

**账户 1 - User1（卖家）**:
```
私钥: 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
地址: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
重命名: User1 - Seller
```

**账户 2 - User3（买家）**:
```
私钥: 0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a
地址: 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC
重命名: User3 - Buyer
```

#### 3. MetaMask 网络配置

```
网络名称: Anvil Local
RPC URL: http://localhost:8545
Chain ID: 31337
货币符号: ETH
```

---

### 阶段 1: User1 卖家流程

#### 步骤 1: 铸造 NFT

```
1. MetaMask 切换到 User1
2. 确保连接 Anvil Local 网络
3. 访问 http://localhost:5173
4. 点击"Connect Wallet" → 连接
5. 点击"🎨 铸造 NFT"标签
6. 输入 Token ID: 100
7. 点击"开始铸造"
8. MetaMask 弹出 → 签名（铸造）✅
9. 等待 2-3 秒
10. MetaMask 再次弹出 → 签名（授权 NFT）✅
11. 看到提示："NFT #100 铸造成功并已授权"
```

**预期结果**:
- ✅ 拥有 NFT #100
- ✅ NFT 已授权给 Marketplace

#### 步骤 2: 创建卖单

```
12. 点击"➕ 创建订单"标签
13. 选择"卖单 (Ask)"
14. 填写:
    - NFT 地址: (已填充) 0x9fE...
    - Token ID: 100
    - 价格: 50
    - 有效期: 7 天
15. 点击"创建订单"
16. 因为已授权，直接签名订单 ✅
17. 看到："订单创建成功 ID: XX"
```

**预期结果**:
- ✅ 卖单出现在"市场订单"
- ✅ 显示"💰 出价"按钮

---

### 阶段 2: User3 买家流程

#### 步骤 3: 获取测试 USDC

```
18. MetaMask 切换到 User3
19. 刷新页面 → 重新连接钱包
20. 点击"💰 测试代币"标签
21. 看到:
    - 当前余额: 0 USDC
    - Marketplace 授权额度: 未授权
22. 点击"🎁 免费获取 1000 USDC"
23. MetaMask 弹出 → 签名（铸造 USDC）✅
24. 等待 2-3 秒
25. 看到余额更新为 1000 USDC ✅
```

**预期结果**:
- ✅ 拥有 1000 USDC
- 🟡 授权额度仍显示"未授权"（符合预期）

**💡 关于 MetaMask 显示 USDC**:

USDC 不会自动显示在 MetaMask 中，需要手动添加代币：

```
1. MetaMask → 资产 → 导入代币
2. 代币合约地址: 0x5FC8d32690cc91D4c39d9d3abcBD16989F875707
3. 代币符号: USDC
4. 小数位数: 6
5. 添加

然后会显示：
USDC
1000 USDC
```

**这是 MetaMask 的标准行为，所有 ERC20 代币都需要手动添加！** ✅

#### 步骤 4: 创建买单并触发授权

```
26. 点击"📊 市场订单"
27. 找到 User1 的卖单（NFT #100, 50 USDC）
28. 点击"💰 出价"按钮
29. 自动跳转到创建订单，已填充:
    - 订单类型: 买单 ✅
    - NFT 地址: 0x9fE... ✅
    - Token ID: 100 ✅
    - 价格: 50.00 ✅
30. 点击"创建订单"
31. 系统检测到 USDC 未授权 → 弹出授权弹窗 ⭐
    ┌─────────────────────────┐
    │ 🔐 需要授权 USDC         │
    │                          │
    │ ○ 仅本次 (50 USDC)      │
    │ ● 无限授权 (推荐) ✅     │
    │                          │
    │ [取消] [授权]           │
    └─────────────────────────┘
32. 选择"无限授权" → 点击"授权"
33. MetaMask 弹出 → 签名（授权 USDC）✅
34. 授权完成 → 自动继续创建订单
35. MetaMask 再次弹出 → 签名（订单签名）✅
36. 看到："订单创建成功 ID: YY"
```

**预期结果**:
- ✅ USDC 已授权给 Marketplace（无限额度）
- ✅ 买单创建成功

---

### 阶段 3: 自动撮合和成交

#### 步骤 5: 观察撮合

```
37. 等待 5-10 秒（撮合引擎周期）
38. 查看日志（可选）:
    cd ~/blockchain-project/web3-knowledge/Oeasy-NFT/packages/services
    tail -f logs/matching-engine.log
    
应该看到：
[INFO] 发现订单匹配 数量=1
[INFO] 匹配订单对 NFT地址=0x9fe... TokenID=100
[INFO] 交易已提交到链上 txHash=0x45b3...
```

**完整日志输出示例**:

```
撮合引擎:
  ✅ 发现订单匹配 数量=1
  ✅ 匹配订单对 NFT地址=0x9fe... TokenID=100
  ✅ 交易已提交到链上 txHash=0x45b3...
  ✅ 订单对已提交执行
  ✅ 已从 Redis 缓存中移除

执行服务:
  ✅ executing trade maker=0xf39... taker=0x3c44...
  ✅ trade transaction submitted txHash=0x45b3... nonce=0

索引服务:
  ✅ 处理TradeExecuted事件 txHash=0x45b3...
  ✅ TradeExecuted事件存储成功
  ✅ 已更新maker订单为已成交
  ✅ 已更新taker订单为已成交
```

#### 步骤 6: 验证成交

```
39. 刷新前端页面
40. 点击"✅ 交易历史"
41. 应该看到新的成交记录 ✅
    - NFT #100
    - 价格: 50 USDC
    - 卖方: 0xf39F...
    - 买方: 0x3C44...
```

**预期结果**:
- ✅ 订单状态: active → filled
- ✅ 交易历史显示成交
- ✅ NFT 所有权: User1 → User3
- ✅ USDC: User3 → User1

---

## 验证与调试

### 链上状态验证

```bash
# 验证 NFT 所有权转移
cast call 0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0 \
  "ownerOf(uint256)" 100 \
  --rpc-url http://localhost:8545
# 应该显示 User3 地址: 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC

# 验证 User1 收到 USDC
cast call 0x5FC8d32690cc91D4c39d9d3abcBD16989F875707 \
  "balanceOf(address)" \
  0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 \
  --rpc-url http://localhost:8545
# 应该显示 50000000 (50 USDC)

# 验证 User3 的 USDC 减少
cast call 0x5FC8d32690cc91D4c39d9d3abcBD16989F875707 \
  "balanceOf(address)" \
  0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC \
  --rpc-url http://localhost:8545
# 应该显示 950000000 (950 USDC)
```

### 数据库验证

```bash
# 检查订单状态
docker exec oeasy-nft-postgres psql -U postgres -d oeasy_nft \
  -c "SELECT id, token_id, status, updated_at FROM orders WHERE status='filled' ORDER BY updated_at DESC LIMIT 5;"

# 检查交易事件
docker exec oeasy-nft-postgres psql -U postgres -d oeasy_nft \
  -c "SELECT * FROM trade_events ORDER BY created_at DESC LIMIT 3;"
```

### 如果交易失败

#### 查看失败原因

```bash
# 1. 查看交易收据的完整信息
cast receipt <交易哈希> --rpc-url http://localhost:8545

# 2. 检查执行服务日志
tail -50 logs/execution-service.log

# 3. 检查索引服务日志
tail -50 logs/indexer.log
```

#### 常见失败原因

| 错误 | 原因 | 解决 |
|------|------|------|
| `ERC721NonexistentToken` | NFT 不存在 | 运行 setup-test-environment.sh |
| `ERC721InsufficientApproval` | NFT 未授权 | 前端重新授权或使用 cast send setApprovalForAll |
| `ERC20InsufficientAllowance` | USDC 未授权 | 前端重新授权或使用 cast send approve |
| `InvalidSignature` | 签名错误 | 检查签名数据和 EIP-712 格式 |
| `NonceConsumed` | Nonce 已使用 | 使用新的 nonce 创建订单 |

---

## 常见问题

### Q1: USDC 不显示在 MetaMask 中？

**A**: 需要手动添加代币（这是标准流程）

```
MetaMask → 资产 → 导入代币
代币地址: 0x5FC8d32690cc91D4c39d9d3abcBD16989F875707
符号: USDC
小数: 6
```

### Q2: 授权弹窗没出现？

**A**: 可能已经授权过了

查看浏览器控制台（F12）会有日志说明。可以在前端查看当前授权额度。

### Q3: 撮合没有发生？

**A**: 检查以下几点：

1. **撮合引擎是否运行**
   ```bash
   ps aux | grep matching-engine | grep -v grep
   ```

2. **Redis 中是否有订单**
   ```bash
   docker exec oeasy-nft-redis redis-cli HLEN orders:active:ask
   docker exec oeasy-nft-redis redis-cli HLEN orders:active:bid
   ```

3. **查看撮合引擎日志**
   ```bash
   tail -f logs/matching-engine.log
   ```

### Q4: 订单状态没有更新为 filled？

**A**: 检查索引服务

1. **索引服务是否运行**
   ```bash
   ps aux | grep indexer | grep -v grep
   ```

2. **查看索引服务日志**
   ```bash
   tail -f logs/indexer.log
   ```

3. **检查是否有地址格式问题**（已在 BUG_FIX_SUMMARY.md 中修复）

### Q5: 前端显示"连接钱包失败"？

**A**: 检查 MetaMask 配置

1. 确保 MetaMask 已安装
2. 确保已添加 Anvil Local 网络
3. 确保 Chain ID 是 31337
4. 刷新页面重试

---

## 🎊 完整测试检查清单

### 准备阶段 ✅

- [ ] Anvil 运行中（端口 8545）
- [ ] PostgreSQL 运行中
- [ ] Redis 运行中
- [ ] 4 个后端服务运行中
  - [ ] 订单服务 (8081)
  - [ ] 撮合引擎 (8082)
  - [ ] 执行服务 (8083)
  - [ ] 索引服务 (8084)
- [ ] 前端运行中（端口 5173）
- [ ] 2 个测试账户已导入 MetaMask
- [ ] MetaMask 已配置 Anvil Local 网络

### User1 测试 ✅

- [ ] 铸造 NFT #100 成功
- [ ] NFT 自动授权成功
- [ ] 创建卖单成功
- [ ] 卖单显示在市场订单

### User3 测试 ✅

- [ ] 获取 1000 USDC 成功
- [ ] USDC 显示在 MetaMask（需手动添加代币）
- [ ] 点击"出价"跳转并预填充
- [ ] 创建买单触发授权弹窗
- [ ] 选择授权额度（无限/精确）
- [ ] 授权成功
- [ ] 买单创建成功

### 撮合成交 ✅

- [ ] 5-10 秒后自动撮合
- [ ] 交易提交到链上
- [ ] 交易执行成功（status 1）
- [ ] 索引服务监听到事件
- [ ] 订单状态更新为 filled
- [ ] 交易历史显示成交记录

### 链上验证 ✅

- [ ] NFT 所有权转移（User1 → User3）
- [ ] USDC 转账（User3 → User1）

---

## 🎯 测试步骤总结

### 快速测试（5分钟）

1. **检查服务**: 运行 `check-system.sh`
2. **准备数据**: 运行 `setup-test-environment.sh`（如需要）
3. **同步 Redis**: 运行 `/tmp/cache_orders.sh`
4. **启动撮合**: 启动撮合引擎
5. **观察日志**: 查看撮合和执行日志
6. **验证结果**: 检查订单状态和交易历史

### 完整测试（10分钟）

1. **User1 铸造 NFT**: 前端操作
2. **User1 创建卖单**: 前端操作
3. **User3 获取 USDC**: 前端操作
4. **User3 创建买单**: 前端操作（触发授权）
5. **等待自动撮合**: 5-10 秒
6. **验证交易历史**: 前端查看

**预计耗时**: 
- 快速测试: 5 分钟
- 完整测试: 10 分钟

---

## 💡 最佳实践

### 测试前

1. ✅ 确保所有服务都在运行
2. ✅ 清理旧的测试数据（可选）
3. ✅ 检查 Anvil 区块高度是否正常

### 测试中

1. ✅ 打开浏览器控制台查看详细日志
2. ✅ 同时监控后端服务日志
3. ✅ 每个步骤完成后验证状态

### 测试后

1. ✅ 验证链上状态
2. ✅ 验证数据库状态
3. ✅ 记录任何异常情况

---

## 🎊 总结

**业务闭环完备性**: ✅ **100% 完备**


