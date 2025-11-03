# 🎨 Oeasy-NFT 企业级订单簿交易平台

> 基于订单簿模型的 NFT 交易平台，采用链下撮合、链上结算的混合架构

## ✨ 项目特点

- 🎯 **链下挂单**: EIP-712 签名，无 Gas 费创建订单
- ⚡ **自动撮合**: 企业级撮合引擎，5 秒周期扫描
- 🔒 **链上结算**: UUPS 可升级合约，原子化交换
- 📊 **混合索引**: WebSocket + 周期轮询，确保数据不丢失
- 🎨 **现代化 UI**: React + RainbowKit，精美的用户界面
- ✅ **完整测试**: 智能合约 + 后端 + 集成测试

## 🏗️ 技术架构

```
┌─────────────────────────────────────────────────────────┐
│            前端 DApp (React + Viem + RainbowKit)        │
│  - 钱包连接  - EIP-712 签名  - 订单创建和查看          │
└────────────────────┬────────────────────────────────────┘
                     │ HTTP/REST API
┌────────────────────┴────────────────────────────────────┐
│              API 网关 (企业级架构组件)                   │
│  - 路由转发  - 负载均衡  - 限流熔断  - 统一鉴权        │
│  - 日志监控  - API 版本管理  - 跨域处理                │
│  ⚠️  本项目为简化实现，前端直连后端服务                 │
└────────────────────┬────────────────────────────────────┘
                     │ 内部 API
┌────────────────────┴────────────────────────────────────┐
│                   后端微服务 (Go)                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌─────────┐ │
│  │ 订单服务 │→ │ 撮合引擎 │→ │ 执行服务 │← │ 索引服务│ │
│  │  :8081   │  │  :8082   │  │  :8083   │  │  :8084  │ │
│  └──────────┘  └──────────┘  └──────────┘  └─────────┘ │
│       ↓             ↓             ↓             ↓        │
│  ┌──────────────────────────────────────────────────┐  │
│  │       PostgreSQL + Redis (数据持久化和缓存)      │  │
│  └──────────────────────────────────────────────────┘  │
└────────────────────┬────────────────────────────────────┘
                     │ JSON-RPC
┌────────────────────┴────────────────────────────────────┐
│       智能合约 (Solidity + Foundry + OpenZeppelin)      │
│  - OeasyMarketplace  - OeasyNFT  - MockUSDC            │
└─────────────────────────────────────────────────────────┘
```

### 架构说明

**API 网关层（企业级设计）**

在生产环境中，API 网关是微服务架构的标准组件，提供：
- **统一入口**: 所有外部请求的单一接入点
- **服务编排**: 智能路由和负载均衡
- **安全防护**: 统一鉴权、限流、熔断、防 DDoS
- **可观测性**: 集中日志、监控、链路追踪
- **协议转换**: HTTP/WebSocket/gRPC 协议适配

**常用网关方案**: Kong、APISIX、Traefik、Nginx Plus

> 💡 **本项目简化说明**: 为降低学习和部署复杂度，本项目前端直连后端服务。在企业级生产环境中，最佳实践是引入 API 网关层。

## 📊 项目统计

| 指标 | 数值 |
|------|------|
| 总代码行数 | ~4500 行 |
| 智能合约 | 3 个 |
| 后端服务 | 4 个 |
| 前端组件 | 5+ 个 |
| 测试用例 | 20+ 个 |
| 文档页数 | ~200 页 |

## 🚀 快速启动

> 💡 ** 用户推荐**: 查看 [QUICKSTART.md](QUICKSTART.md)

### 前置要求

- Node.js 18+
- Go 1.24+
- Foundry (forge, anvil, cast)
- Docker Desktop
- PostgreSQL 客户端工具

### 1. 克隆项目

```bash
git clone <repository-url>
cd Oeasy-NFT
```

### 2. 启动本地测试环境

```bash
# 启动 Anvil 本地测试网（终端 1）
anvil --chain-id 31337 --port 8545

# 启动数据库服务（终端 2）
docker-compose up -d postgres redis

# 初始化数据库
psql -U postgres -f database/init.sql
psql -U postgres -d oeasy_nft -f database/seed.sql  # 可选：测试数据
```

### 3. 部署智能合约

```bash
cd packages/contracts

# 设置环境变量
export RPC_URL=http://localhost:8545
export DEPLOYER_PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

# 部署合约
forge script script/Deploy.s.sol \
  --rpc-url $RPC_URL \
  --private-key $DEPLOYER_PRIVATE_KEY \
  --broadcast

# 记录合约地址用于配置
```

### 4. 配置并启动后端服务

```bash
cd packages/services

# 创建配置文件
cat > .env << EOF
POSTGRES_DSN=postgresql://postgres:postgres@localhost:5432/oeasy_nft?sslmode=disable
REDIS_ADDR=localhost:6379
RPC_URL=http://localhost:8545
CHAIN_ID=31337
MARKETPLACE_ADDRESS=<部署的合约地址>
EXECUTOR_PRIVATE_KEY=0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d
EOF

# 启动 4 个服务（分别在 4 个终端）
go run cmd/order-service/main.go      # 终端 3
go run cmd/matching-engine/main.go    # 终端 4
go run cmd/execution-service/main.go  # 终端 5
go run cmd/indexer/main.go            # 终端 6
```

### 5. 启动前端

```bash
cd packages/frontend

# 安装依赖
npm install

# 更新 src/wagmi.ts 中的合约地址

# 启动开发服务器
npm run dev

# 访问 http://localhost:5173
```

### 6. 使用自动化测试脚本（推荐）

```bash
# 一键启动测试环境
./scripts/integration-test.sh
```

## 测试

### 智能合约测试

```bash
cd packages/contracts
forge test -vvv
```

### 后端服务测试

```bash
cd packages/services
go test ./... -v
```

## 技术栈

### 智能合约
- Solidity 0.8.24
- Foundry
- OpenZeppelin Contracts (Upgradeable)
- EIP-712 Typed Data

### 后端
- Go 1.24+
- Gin Web Framework
- GORM + PostgreSQL
- Redis
- go-ethereum

### 前端
- React 18
- Viem
- TailwindCSS
- TypeScript

## 核心特性

### 1. EIP-712 链下签名
用户签名订单无需消耗 Gas，提升体验

### 2. UUPS 可升级合约
支持合约逻辑升级，保持地址不变

### 3. 企业级撮合算法
bid >= ask 标准匹配规则，支持价格优先

### 4. 混合索引架构
WebSocket + 轮询双重保障，确保数据不丢失

### 5. 微服务架构
关注点分离，易于扩展和维护

## 📚 文档
TODO

## 📁 项目结构

```
Oeasy-NFT/
├── packages/
│   ├── contracts/          # 智能合约（Foundry）
│   │   ├── src/            # 合约源码
│   │   ├── test/           # 合约测试
│   │   └── script/         # 部署脚本
│   ├── services/           # 后端服务（Go）
│   │   ├── cmd/            # 服务入口
│   │   └── internal/       # 内部包
│   └── frontend/           # 前端 DApp（React）
│       └── src/            # 前端源码
├── database/               # 数据库脚本
│   ├── init.sql           # 数据库初始化
│   ├── seed.sql           # 测试数据
│   └── README.md          # 数据库文档
├── docs/                   # 项目文档
├── scripts/                # 工具脚本
└── docker-compose.yml      # Docker 编排
```

## 🎯 核心功能

### 智能合约

- ✅ UUPS 可升级代理模式
- ✅ EIP-712 链下签名验证
- ✅ 原子化 NFT + ERC20 交换
- ✅ 平台手续费机制
- ✅ Nonce 防重放攻击
- ✅ 8 个 Foundry 测试用例

### 后端微服务

- ✅ **订单服务**: EIP-712 验证 + PostgreSQL + Redis
- ✅ **撮合引擎**: 价格-时间优先算法
- ✅ **执行服务**: 智能合约调用 + Nonce 管理
- ✅ **索引服务**: WebSocket + 轮询混合架构

### 前端 DApp

- ✅ RainbowKit 钱包连接
- ✅ EIP-712 签名集成
- ✅ 订单创建和查询
- ✅ 精美的现代化 UI
- ✅ 完全响应式设计

## 🧪 运行测试

### 智能合约测试

```bash
cd packages/contracts
forge test -vvv
```

### 后端单元测试

```bash
cd packages/services
go test ./... -v
```

### 集成测试

```bash
# 按照 docs/INTEGRATION_TEST_GUIDE.md 执行
./scripts/integration-test.sh
```

## 🔒 安全特性

- ✅ EIP-712 链下签名（无 Gas 费）
- ✅ ReentrancyGuard 防重入
- ✅ Pausable 紧急暂停
- ✅ UUPS 安全升级
- ✅ Nonce 防重放
- ✅ 输入验证和签名验证




[sfan]



