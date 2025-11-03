# 📜 脚本说明

## 🎯 核心脚本（4 个）

### 1. `setup-test-environment.sh`
**用途**: 准备完整的测试环境（一次性执行）

**功能**:
- 铸造 NFT #1, #2 给 User1
- 铸造 USDC 给 User2, User3
- 授权 NFT 和 USDC 给 Marketplace
- 验证链上状态

**何时使用**:
- ✅ 第一次启动项目
- ✅ 重置 Anvil 后
- ✅ 需要测试完整交易流程

**运行**:
```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT
./scripts/setup-test-environment.sh
```

---

### 2. `start-all-services.sh`
**用途**: 一键启动所有后端服务（每次开发）

**功能**:
- 启动订单服务
- 启动撮合引擎
- 启动执行服务
- 启动索引服务
- 所有日志保存到 logs/ 目录

**何时使用**:
- ✅ 每次开发时
- ✅ 代替手动开 4 个终端

**运行**:
```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT
./scripts/start-all-services.sh
```

---

### 3. `stop-all-services.sh`
**用途**: 停止所有后端服务

**功能**:
- 停止所有 go run 进程
- 清理 PID 文件

**何时使用**:
- ✅ 开发结束
- ✅ 需要重启服务

**运行**:
```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT
./scripts/stop-all-services.sh
```

---

### 4. `view-logs.sh`
**用途**: 交互式查看日志

**功能**:
- 菜单选择要查看的服务
- 实时显示日志

**何时使用**:
- ✅ 需要查看特定服务日志
- ✅ 调试问题

**运行**:
```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT
./scripts/view-logs.sh
```

---

## 📊 完整的开发流程

### 第一次启动（完整流程）

```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT

# 1. 启动 Anvil（新终端保持运行）
anvil --chain-id 31337

# 2. 启动数据库（新终端）
docker-compose up -d

# 3. 初始化数据库
docker exec -i oeasy-nft-postgres psql -U postgres < database/init.sql
docker exec -i oeasy-nft-postgres psql -U postgres -d oeasy_nft < database/seed.sql

# 4. 部署合约
cd packages/contracts
forge script script/Deploy.s.sol --broadcast --rpc-url http://localhost:8545 --private-key 0xac0974...
# 记录合约地址

# 5. 配置后端
cd ../services
# 编辑 .env 填入合约地址

# 6. 准备测试资产
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT
./scripts/setup-test-environment.sh  ← 只需这一个

# 7. 启动所有后端服务
./scripts/start-all-services.sh  ← 只需这一个

# 8. 启动前端（新终端）
cd packages/frontend
npm run dev
```

---

### 日常开发（简化流程）

```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT

# 1. 启动服务
./scripts/start-all-services.sh

# 2. 查看日志（如果需要）
cd packages/services
tail -f logs/matching-engine.log

# 3. 停止服务
./scripts/stop-all-services.sh
```

---

## 🗑️ 已删除的重复脚本

- ❌ `prepare-test-data.sh` - 功能并入 setup-test-environment.sh
- ❌ `sync-orders-to-redis.sh` - 不再需要（API 创建时自动缓存）
- ❌ `start-dev.sh` - 重复 start-all-services.sh

---

## 🎯 现在只有 4 个核心脚本

**清晰、简洁、必要！** ✅

查看详细说明：`scripts/README.md`

