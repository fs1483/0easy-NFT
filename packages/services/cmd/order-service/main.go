package main

import (
	"log"

	"github.com/Oeasy-NFT/services/internal/config"
	"github.com/Oeasy-NFT/services/internal/orders"
	"github.com/joho/godotenv"
)

func main() {
	// 自动加载 .env 文件
	_ = godotenv.Load()

	log.Println("🚀 正在启动订单服务...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v", err)
	}

	log.Printf("✅ 配置加载成功 - 端口: %s, 数据库: oeasy_nft\n", cfg.OrderServicePort)

	svc, err := orders.NewService(cfg)
	if err != nil {
		log.Fatalf("❌ 初始化订单服务失败: %v", err)
	}

	log.Printf("✅ 订单服务初始化完成，开始监听端口 %s...\n", cfg.OrderServicePort)

	if err := svc.Run(); err != nil {
		log.Fatalf("❌ 订单服务停止: %v", err)
	}
}
