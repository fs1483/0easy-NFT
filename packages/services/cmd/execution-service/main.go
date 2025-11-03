package main

import (
	"log"

	"github.com/Oeasy-NFT/services/internal/config"
	"github.com/Oeasy-NFT/services/internal/execution"
	"github.com/joho/godotenv"
)

func main() {
	// 自动加载 .env 文件
	_ = godotenv.Load()

	log.Println("🚀 正在启动执行服务...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v", err)
	}

	log.Printf("✅ 配置加载成功 - 端口: %s, Marketplace: %s\n",
		cfg.ExecutionServicePort,
		cfg.MarketplaceAddr)

	execSvc, err := execution.NewService(cfg)
	if err != nil {
		log.Fatalf("❌ 初始化执行服务失败: %v", err)
	}

	log.Printf("✅ 执行服务初始化完成，开始监听端口 %s...\n", cfg.ExecutionServicePort)

	if err := execSvc.Run(); err != nil {
		log.Fatalf("❌ 执行服务停止: %v", err)
	}
}
