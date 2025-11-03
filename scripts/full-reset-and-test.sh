#!/bin/bash

# ============================================
# 完整重置并准备测试环境
# ============================================
# 用途：清理所有状态，从头开始，确保一致性

set -e

BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Oeasy-NFT 完整重置${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 进入项目根目录
cd ~/blockchain-project/web3-knowledge/Oeasy-NFT

# 1. 停止所有服务
echo -e "${YELLOW}1. 停止所有服务...${NC}"
./scripts/stop-all-services.sh > /dev/null 2>&1

# 2. 清理数据库
echo -e "${YELLOW}2. 清理数据库...${NC}"
docker exec oeasy-nft-postgres psql -U postgres -d oeasy_nft -c "DELETE FROM orders; DELETE FROM trade_events; UPDATE indexer_status SET last_processed_block = 0 WHERE id = 1;" > /dev/null 2>&1
echo -e "${GREEN}[✓]${NC} 数据库已清理"

# 3. 清理 Redis
echo -e "${YELLOW}3. 清理 Redis...${NC}"
docker exec oeasy-nft-redis redis-cli FLUSHALL > /dev/null 2>&1
echo -e "${GREEN}[✓]${NC} Redis 已清理"

# 4. 重新部署合约（使用最新代码）
echo -e "${YELLOW}4. 重新部署智能合约...${NC}"
cd packages/contracts
export DEPLOYER_PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
DEPLOY_OUTPUT=$(forge script script/Deploy.s.sol --rpc-url http://localhost:8545 --broadcast 2>&1)

MARKETPLACE=$(echo "$DEPLOY_OUTPUT" | grep "OeasyMarketplace" | grep "0x" | awk '{print $NF}')
NFT=$(echo "$DEPLOY_OUTPUT" | grep "OeasyNFT" | grep "0x" | awk '{print $NF}')
USDC=$(echo "$DEPLOY_OUTPUT" | grep "MockUSDC" | grep "0x" | awk '{print $NF}')

echo -e "${GREEN}[✓]${NC} 合约已部署"
echo "  Marketplace: $MARKETPLACE"
echo "  NFT:         $NFT"
echo "  USDC:        $USDC"

# 5. 自动更新配置
echo -e "${YELLOW}5. 更新配置文件...${NC}"
cd ../services
cat > .env << ENVEOF
POSTGRES_DSN=postgresql://postgres:postgres@localhost:5432/oeasy_nft?sslmode=disable
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
RPC_URL=http://localhost:8545
CHAIN_ID=31337
MARKETPLACE_ADDRESS=$MARKETPLACE
NFT_ADDRESS=$NFT
USDC_ADDRESS=$USDC
EXECUTOR_PRIVATE_KEY=0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d
ORDER_SERVICE_PORT=8081
MATCHING_SERVICE_PORT=8082
EXECUTION_SERVICE_PORT=8083
INDEXER_SERVICE_PORT=8084
ENVEOF

cd ../frontend
cat > .env.local << ENVEOF
VITE_CHAIN_ID=31337
VITE_RPC_URL=http://localhost:8545
VITE_MARKETPLACE_ADDRESS=$MARKETPLACE
VITE_NFT_ADDRESS=$NFT
VITE_USDC_ADDRESS=$USDC
VITE_API_BASE_URL=http://localhost:8081
VITE_ENABLE_DEBUG=true
ENVEOF

echo -e "${GREEN}[✓]${NC} 配置已更新"

cd ../..

# 6. 启动服务
echo -e "${YELLOW}6. 启动后端服务...${NC}"
./scripts/start-all-services.sh > /dev/null 2>&1
sleep 5
echo -e "${GREEN}[✓]${NC} 所有服务已启动"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  重置完成！${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "📝 新的合约地址："
echo "  Marketplace: $MARKETPLACE"
echo "  NFT:         $NFT"
echo "  USDC:        $USDC"
echo ""
echo "🎯 下一步："
echo "  1. 硬刷新前端页面（Cmd+Shift+R）"
echo "  2. MetaMask 更新 USDC 代币地址为：$USDC"
echo "  3. 开始测试："
echo "     - User1: 铸造 NFT #100 → 创建卖单 50 USDC"
echo "     - User3: 获取 USDC → 创建买单 55 USDC"
echo "     - 等待 5-10 秒自动撮合"
echo "     - 查看交易历史 ✅"
echo ""
echo "✅ 环境已完全重置，可以开始干净的测试了！"

