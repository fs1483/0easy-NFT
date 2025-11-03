package main

import (
	"github.com/Oeasy-NFT/services/internal/config"
	"github.com/Oeasy-NFT/services/internal/httpserver"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	server := httpserver.NewGateway(cfg)
	if err := server.Run(); err != nil {
		log.Fatalf("gateway server stopped: %v", err)
	}
}
