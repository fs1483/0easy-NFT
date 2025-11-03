// Package httpserver 实现 API 网关服务
// 网关作为系统统一入口，负责路由分发、CORS 配置、速率限制等横切关注点
package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/Oeasy-NFT/services/internal/config"
	"github.com/Oeasy-NFT/services/internal/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Gateway HTTP 服务器，将外部请求路由到对应的下游微服务
type Gateway struct {
	cfg    *config.Config
	engine *gin.Engine
	srv    *http.Server
}

// NewGateway 构造 API 网关实例，配置默认中间件和 CORS
func NewGateway(cfg *config.Config) *Gateway {
	gin.SetMode(gin.ReleaseMode)
	eng := gin.New()
	eng.Use(gin.Recovery())
	eng.Use(logger.HTTPLogger())

	// 配置 CORS 允许前端跨域访问
	eng.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	g := &Gateway{
		cfg:    cfg,
		engine: eng,
	}
	g.registerRoutes()
	return g
}

// registerRoutes 注册公开的 API 路由
// 在 MVP 阶段，网关主要提供健康检查，订单服务可直接暴露端口
// TODO: [Gateway] - 添加反向代理逻辑将请求转发到各微服务（使用 httputil.ReverseProxy）
func (g *Gateway) registerRoutes() {
	// 健康检查端点
	g.engine.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok", "service": "gateway"})
	})

	// API 版本信息
	g.engine.GET("/api/version", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"version": "1.0.0",
			"services": gin.H{
				"orders":    "http://localhost:" + g.cfg.OrderServicePort,
				"matching":  "running",
				"execution": "http://localhost:" + g.cfg.ExecutionServicePort,
				"indexer":   "running",
			},
		})
	})

	// TODO: [Gateway] - 实现订单服务代理
	// api := g.engine.Group("/api")
	// orderProxy := httputil.NewSingleHostReverseProxy(orderServiceURL)
	// api.Any("/orders/*path", gin.WrapH(orderProxy))
}

// Run 启动 HTTP 服务器
func (g *Gateway) Run() error {
	g.srv = &http.Server{
		Addr:              ":" + g.cfg.HTTPPort,
		Handler:           g.engine,
		ReadHeaderTimeout: 5 * time.Second,
	}
	logger.Info("API gateway listening", "port", g.cfg.HTTPPort)
	return g.srv.ListenAndServe()
}

// Shutdown 优雅关闭服务器
func (g *Gateway) Shutdown(ctx context.Context) error {
	if g.srv == nil {
		return nil
	}
	logger.Info("API gateway shutting down")
	return g.srv.Shutdown(ctx)
}
