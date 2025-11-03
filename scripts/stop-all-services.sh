#!/bin/bash

# ============================================
# 停止所有后端服务（优化版）
# ============================================

GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🛑 停止所有后端服务...${NC}"
echo ""

# 函数：停止服务
stop_service() {
    local service_name=$1
    local process_pattern=$2
    
    # 先尝试 pkill
    if pkill -f "$process_pattern" 2>/dev/null; then
        sleep 1
        # 检查是否还在运行
        if pgrep -f "$process_pattern" > /dev/null 2>&1; then
            # 如果还在，强制杀死
            pkill -9 -f "$process_pattern" 2>/dev/null
            sleep 0.5
        fi
        echo -e "${GREEN}[✓]${NC} $service_name 已停止"
    else
        echo -e "${RED}[✗]${NC} $service_name 未运行"
    fi
}

# 停止各个服务
stop_service "订单服务" "go run cmd/order-service"
stop_service "撮合引擎" "go run cmd/matching-engine"
stop_service "执行服务" "go run cmd/execution-service"
stop_service "索引服务" "go run cmd/indexer"

# 额外清理：强制杀死任何残留的 go run 进程
sleep 1
remaining=$(pgrep -f "go run cmd" | wc -l | tr -d ' ')
if [ "$remaining" -gt 0 ]; then
    echo -e "${BLUE}[INFO]${NC} 清理残留进程..."
    pkill -9 -f "go run cmd" 2>/dev/null
    sleep 1
fi

# 清理端口占用
for port in 8081 8083 8084; do
    if lsof -ti :$port > /dev/null 2>&1; then
        echo -e "${BLUE}[INFO]${NC} 清理端口 $port..."
        lsof -ti :$port | xargs kill -9 2>/dev/null || true
    fi
done

# 删除 PID 文件
cd "$(dirname "$0")/../packages/services"
rm -f logs/*.pid 2>/dev/null

echo ""
echo -e "${GREEN}✅ 所有服务已彻底停止${NC}"
echo ""

# 验证
remaining=$(pgrep -f "go run cmd" | wc -l | tr -d ' ')
if [ "$remaining" -eq 0 ]; then
    echo -e "${GREEN}✓${NC} 确认：无残留进程"
else
    echo -e "${RED}⚠${NC} 警告：仍有 $remaining 个进程未停止"
    echo "可以手动执行: pkill -9 -f 'go run cmd'"
fi

