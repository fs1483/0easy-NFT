#!/bin/bash

# ============================================
# 一键启动所有后端服务脚本
# ============================================
# 解决多终端管理混乱的问题
# 所有服务在后台运行，日志输出到独立文件

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_header() {
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

# 进入服务目录
cd "$(dirname "$0")/../packages/services"

# 创建日志目录
mkdir -p logs

print_header "停止旧服务"

# 停止所有旧的服务进程
pkill -f "go run cmd/order-service" 2>/dev/null && print_info "已停止订单服务" || true
pkill -f "go run cmd/matching-engine" 2>/dev/null && print_info "已停止撮合引擎" || true
pkill -f "go run cmd/execution-service" 2>/dev/null && print_info "已停止执行服务" || true
pkill -f "go run cmd/indexer" 2>/dev/null && print_info "已停止索引服务" || true

sleep 2

print_header "启动后端服务"

# 启动订单服务
print_info "启动订单服务..."
nohup go run cmd/order-service/main.go > logs/order-service.log 2>&1 &
ORDER_PID=$!
echo $ORDER_PID > logs/order-service.pid
print_success "订单服务已启动 (PID: $ORDER_PID, 日志: logs/order-service.log)"

sleep 2

# 启动撮合引擎
print_info "启动撮合引擎..."
nohup go run cmd/matching-engine/main.go > logs/matching-engine.log 2>&1 &
MATCHING_PID=$!
echo $MATCHING_PID > logs/matching-engine.pid
print_success "撮合引擎已启动 (PID: $MATCHING_PID, 日志: logs/matching-engine.log)"

sleep 1

# 启动执行服务
print_info "启动执行服务..."
nohup go run cmd/execution-service/main.go > logs/execution-service.log 2>&1 &
EXEC_PID=$!
echo $EXEC_PID > logs/execution-service.pid
print_success "执行服务已启动 (PID: $EXEC_PID, 日志: logs/execution-service.log)"

sleep 1

# 启动索引服务
print_info "启动索引服务..."
nohup go run cmd/indexer/main.go > logs/indexer.log 2>&1 &
INDEXER_PID=$!
echo $INDEXER_PID > logs/indexer.log.pid
print_success "索引服务已启动 (PID: $INDEXER_PID, 日志: logs/indexer.log)"

sleep 3

print_header "验证服务状态"

# 检查服务是否运行
check_service() {
    if ps -p $1 > /dev/null 2>&1; then
        print_success "$2 运行中 (PID: $1)"
    else
        print_error "$2 启动失败"
    fi
}

check_service $ORDER_PID "订单服务"
check_service $MATCHING_PID "撮合引擎"
check_service $EXEC_PID "执行服务"
check_service $INDEXER_PID "索引服务"

# 检查端口
echo ""
print_info "端口监听状态:"
lsof -i :8081 > /dev/null 2>&1 && print_success "端口 8081 (订单服务)" || print_error "端口 8081 未监听"
lsof -i :8083 > /dev/null 2>&1 && print_success "端口 8083 (执行服务)" || print_error "端口 8083 未监听"

print_header "服务启动完成"

echo "📊 查看实时日志:"
echo ""
echo -e "  ${YELLOW}订单服务:${NC}   tail -f logs/order-service.log"
echo -e "  ${YELLOW}撮合引擎:${NC}   tail -f logs/matching-engine.log"
echo -e "  ${YELLOW}执行服务:${NC}   tail -f logs/execution-service.log"
echo -e "  ${YELLOW}索引服务:${NC}   tail -f logs/indexer.log"
echo ""
echo -e "  ${YELLOW}所有日志:${NC}   tail -f logs/*.log"
echo ""
echo "🛑 停止所有服务:"
echo -e "  ${YELLOW}运行:${NC}        ./scripts/stop-all-services.sh"
echo ""
echo "✅ 所有服务已在后台运行"

