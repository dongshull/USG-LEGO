package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"USG-LEGO/internal/auth"
	"USG-LEGO/internal/config"
	"USG-LEGO/internal/logger"
	"USG-LEGO/internal/routes"
)

func main() {
	// 初始化配置
	cfg := config.LoadConfig()
	
	// 初始化日志
	logger.InitLogger(cfg.LogLevel)
	
	// 初始化认证模块
	auth.InitAuth(cfg.RootDir)
	
	// 创建 Gin 引擎
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	// 添加日志中间件
	router.Use(logger.LoggingMiddleware())
	
	// 添加恢复中间件
	router.Use(gin.Recovery())
	
	// 注册路由
	routes.RegisterRoutes(router, cfg)
	
	// 添加一个简单的静态文件服务
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "<h1>Welcome to USG-LEGO</h1><p>This is a single file hosting service for iOS proxy tools.</p>")
	})
	
	// 启动服务
	log.Info().Str("addr", cfg.Listen).Msg("Starting server")
	if err := router.Run(cfg.Listen); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}