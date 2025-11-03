package main

import (
	"log"

	"github.com/Oeasy-NFT/services/internal/config"
	"github.com/Oeasy-NFT/services/internal/indexer"
	"github.com/joho/godotenv"
)

func main() {
	// 自动加载 .env 文件
	_ = godotenv.Load()

	log.Println("🚀 正在启动索引服务...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v", err)
	}

	log.Printf("✅ 配置加载成功 - RPC: %s, Marketplace: %s\n",
		cfg.RPCURL,
		cfg.MarketplaceAddr)

	svc, err := indexer.NewService(cfg)
	if err != nil {
		log.Fatalf("❌ 初始化索引服务失败: %v", err)
	}

	log.Println("✅ 索引服务初始化完成，开始监听区块链事件...")

	if err := svc.Run(); err != nil {
		log.Fatalf("❌ 索引服务停止: %v", err)
	}
}
