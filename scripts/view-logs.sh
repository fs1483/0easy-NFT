#!/bin/bash

# ============================================
# 查看服务日志工具
# ============================================

BLUE='\033[0;34m'
NC='\033[0m'

cd "$(dirname "$0")/../packages/services"

if [ ! -d "logs" ]; then
    echo "日志目录不存在，请先运行服务"
    exit 1
fi

echo -e "${BLUE}请选择要查看的日志:${NC}"
echo ""
echo "  1) 订单服务 (order-service.log)"
echo "  2) 撮合引擎 (matching-engine.log)"
echo "  3) 执行服务 (execution-service.log)"
echo "  4) 索引服务 (indexer.log)"
echo "  5) 所有服务（并排显示）"
echo "  6) 退出"
echo ""
read -p "选择 [1-6]: " choice

case $choice in
    1)
        echo -e "${BLUE}=== 订单服务日志 ===${NC}"
        tail -f logs/order-service.log
        ;;
    2)
        echo -e "${BLUE}=== 撮合引擎日志 ===${NC}"
        tail -f logs/matching-engine.log
        ;;
    3)
        echo -e "${BLUE}=== 执行服务日志 ===${NC}"
        tail -f logs/execution-service.log
        ;;
    4)
        echo -e "${BLUE}=== 索引服务日志 ===${NC}"
        tail -f logs/indexer.log
        ;;
    5)
        echo -e "${BLUE}=== 所有服务日志 ===${NC}"
        tail -f logs/*.log
        ;;
    6)
        exit 0
        ;;
    *)
        echo "无效选择"
        exit 1
        ;;
esac

