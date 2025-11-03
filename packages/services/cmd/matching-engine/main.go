package main

import (
	"log"

	"github.com/Oeasy-NFT/services/internal/config"
	"github.com/Oeasy-NFT/services/internal/matching"
	"github.com/joho/godotenv"
)

func main() {
	// 自动加载 .env 文件
	_ = godotenv.Load()

	log.Println("🚀 正在启动撮合引擎...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v", err)
	}

	log.Printf("✅ 配置加载成功 - Redis: %s\n", cfg.RedisAddr)

	engine, err := matching.NewEngine(cfg)
	if err != nil {
		log.Fatalf("❌ 初始化撮合引擎失败: %v", err)
	}

	log.Println("✅ 撮合引擎初始化完成，开始扫描订单簿（每 5 秒）...")

	if err := engine.Run(); err != nil {
		log.Fatalf("❌ 撮合引擎停止: %v", err)
	}
}
