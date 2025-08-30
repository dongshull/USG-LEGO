package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"USG-LEGO/internal/config"
	"USG-LEGO/internal/files"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(router *gin.Engine, cfg *config.Config) {
	// 健康检查接口
	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	
	// 文件相关接口
	api := router.Group("/api")
	{
		api.GET("/files", files.GetFile)
	}
	
	// 原始文件访问接口
	router.GET("/raw/*path", func(c *gin.Context) {
		// 这里应该实现原始文件访问逻辑
		c.String(http.StatusOK, "Raw file access placeholder")
	})
}