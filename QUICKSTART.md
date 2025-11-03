# 🚀 5 分钟快速启动（macOS）

最简化的启动流程，让你快速体验 Oeasy-NFT。

---

## 📋 前置要求（一次性安装）

```bash
# 安装所有依赖（复制整段运行）
brew install node go postgresql@14
curl -L https://foundry.paradigm.xyz | bash && foundryup
npm install -g pnpm
brew install --cask docker
```

启动 Docker Desktop 应用。

---

## ⚡ 快速启动（每次使用）

### 第 1 步：启动 Anvil（终端 1）

```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT
anvil --chain-id 31337 --port 8545
```

**看到账户列表就成功了！** ✅

---

### 第 2 步：启动数据库（终端 2）

```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT

# 启动数据库
docker-compose up -d

# 初始化（首次）
docker exec -i oeasy-nft-postgres psql -U postgres < database/init.sql
```

**看到 "数据库初始化完成" 就成功了！** ✅

---

### 第 3 步：部署合约（终端 2 继续）

```bash
cd packages/contracts

# 一键部署
forge script script/Deploy.s.sol \
  --rpc-url http://localhost:8545 \
  --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
  --broadcast
```

**📝 记录输出中的三个合约地址！**

---

### 第 4 步：配置后端（终端 2 继续）

```bash
cd ../services

# 创建配置文件
cat > .env << EOF
POSTGRES_DSN=postgresql://postgres:postgres@localhost:5432/oeasy_nft?sslmode=disable
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
RPC_URL=http://localhost:8545
CHAIN_ID=31337
MARKETPLACE_ADDRESS=你的Marketplace地址
NFT_ADDRESS=你的NFT地址
USDC_ADDRESS=你的USDC地址
EXECUTOR_PRIVATE_KEY=0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d
ORDER_SERVICE_PORT=8081
EXECUTION_SERVICE_PORT=8083
EOF

# 🔴 用文本编辑器打开 .env，替换三个合约地址
```

---

### 第 5 步：启动后端（4 个新终端）

**终端 3**:
```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT/packages/services
go run cmd/order-service/main.go
```

**终端 4**:
```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT/packages/services
go run cmd/matching-engine/main.go
```

**终端 5**:
```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT/packages/services
go run cmd/execution-service/main.go
```

**终端 6**:
```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT/packages/services
go run cmd/indexer/main.go
```

**每个都看到 "service initialized" 就成功了！** ✅

> 💡 **提示**: 不需要 `source .env`，程序会自动加载配置文件！

---

### 第 6 步：启动前端（终端 7）

```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT/packages/frontend

# 首次需要安装依赖
npm install

# 🔴 创建配置文件（推荐使用环境变量，符合企业级标准）
cp env.example .env.local

# 编辑 .env.local，填入实际的合约地址
vim .env.local
# 或
code .env.local

# 修改这三行：
# VITE_MARKETPLACE_ADDRESS=你的Marketplace地址
# VITE_NFT_ADDRESS=你的NFT地址
# VITE_USDC_ADDRESS=你的USDC地址

# 启动（会自动读取 .env.local）
npm run dev
```

**访问**: http://localhost:5173

> 💡 **企业级最佳实践**: 使用 `.env.local` 配置文件而不是修改代码！

---

### 第 7 步：准备测试数据（终端 8）

```bash
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT
./scripts/prepare-test-data.sh
```

**看到 "测试数据准备完成" 就成功了！** ✅

---

## 🦊 配置 MetaMask（一次性）

### 1. 添加 Anvil 网络

- 网络名称: `Anvil Local`
- RPC URL: `http://localhost:8545`
- Chain ID: `31337`
- 货币符号: `ETH`

### 2. 导入测试账户

**User1（卖家）**:
```
私钥: 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
```

**User2（买家）**:
```
私钥: 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d
```

---

## 🎮 开始测试！

### 测试流程

1. **连接钱包** → 选择 User1
2. **创建卖单** → NFT #1, 价格 100 USDC
3. **切换账户** → 选择 User2
4. **创建买单** → NFT #1, 价格 100 USDC
5. **等待 5 秒** → 自动撮合和执行
6. **验证结果** → NFT 转移给 User2

---

## 📊 系统检查

### 所有服务运行检查

```bash
# 新开一个终端运行
cat > check-services.sh << 'EOF'
#!/bin/bash
echo "🔍 检查所有服务..."
echo ""

check_port() {
    if lsof -Pi :$1 -sTCP:LISTEN -t >/dev/null; then
        echo "✅ 端口 $1: $2"
    else
        echo "❌ 端口 $1: $2 未运行"
    fi
}

check_port 8545 "Anvil"
check_port 5432 "PostgreSQL"
check_port 6379 "Redis"
check_port 8081 "订单服务"
check_port 8082 "撮合引擎"
check_port 8083 "执行服务"
check_port 8084 "索引服务"
check_port 5173 "前端"

echo ""
echo "💾 数据库检查:"
docker exec oeasy-nft-postgres psql -U postgres -d oeasy_nft -t -c "SELECT COUNT(*) FROM orders;" 2>/dev/null && echo "✅ 数据库连接正常" || echo "❌ 数据库连接失败"

echo ""
echo "📊 Redis 检查:"
docker exec oeasy-nft-redis redis-cli ping 2>/dev/null && echo "✅ Redis 连接正常" || echo "❌ Redis 连接失败"
EOF

chmod +x check-services.sh
./check-services.sh
```

---

## 🛑 停止所有服务

```bash
# 在各个终端按 Ctrl+C

# 或运行停止脚本
pkill -f anvil
pkill -f "go run"
pkill -f vite
docker-compose down
```

---

## 💡 提示

### 如果遇到问题

1. **查看详细指南**: `docs/MACOS_QUICKSTART.md`
2. **查看日志**: 在各个终端查看输出
3. **重启服务**: 按 Ctrl+C 停止，重新运行命令
4. **重置环境**: 
   ```bash
   docker-compose down -v  # 清理数据库
   pkill -f anvil          # 重启 Anvil
   ```

### 常用命令

```bash
# 查看订单
curl http://localhost:8081/api/orders | jq

# 查看数据库
docker exec -it oeasy-nft-postgres psql -U postgres -d oeasy_nft

# 查看 Redis
docker exec -it oeasy-nft-redis redis-cli

# 查看区块链状态
cast block latest --rpc-url http://localhost:8545
```

---


