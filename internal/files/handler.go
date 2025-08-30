package files

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GetFile handles GET /api/files?path=... requests
// @Summary Get file content
// @Description Get file content by path, either as JSON metadata or as file stream
// @Tags files
// @Accept json
// @Produce json
// @Param path query string true "File path"
// @Param api query string false "API Key"
// @Success 200 {object} FileInfo
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/files [get]
func GetFile(c *gin.Context) {
	// 获取路径参数
	path := c.Query("path")
	if path == "" {
		log.Warn().Str("path", path).Msg("Path parameter is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path parameter is required"})
		return
	}

	// 清理路径并检查是否跳出根目录
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		log.Warn().Str("path", path).Str("cleanPath", cleanPath).Msg("Path traversal attempt detected")
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid path"})
		return
	}

	// 获取根目录配置
	rootDir := os.Getenv("ROOT_DIR")
	if rootDir == "" {
		log.Error().Msg("ROOT_DIR environment variable not set")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
		return
	}

	// 构造完整文件路径
	fullPath := filepath.Join(rootDir, cleanPath)

	// 检查文件是否存在
	fileInfo, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		log.Warn().Str("path", fullPath).Msg("File not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	if err != nil {
		log.Error().Err(err).Str("path", fullPath).Msg("Failed to stat file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to access file"})
		return
	}

	// 检查是否为公共路径
	isPublic := IsPublic(cleanPath)

	// 如果不是公共路径，需要验证 API Key
	if !isPublic {
		// 获取 API Key
		apiKey := c.Query("api")
		if apiKey == "" {
			log.Warn().Str("path", cleanPath).Msg("API key required for private file access")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			return
		}

		// 验证 API Key
		if !ValidateAPIKey(apiKey) {
			log.Warn().Str("apiKey", apiKey).Msg("Invalid API key")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}
	}

	// 如果是目录，返回目录信息
	if fileInfo.IsDir() {
		// 读取目录内容
		entries, err := os.ReadDir(fullPath)
		if err != nil {
			log.Error().Err(err).Str("path", fullPath).Msg("Failed to read directory")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read directory"})
			return
		}

		// 构造目录内容信息
		var files []FileInfo
		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			
			files = append(files, FileInfo{
				Name:    entry.Name(),
				IsDir:   entry.IsDir(),
				Size:    info.Size(),
				ModTime: info.ModTime(),
			})
		}

		// 返回 JSON 格式的目录信息
		c.JSON(http.StatusOK, gin.H{
			"path":  cleanPath,
			"isDir": true,
			"files": files,
		})
		return
	}

	// 如果是文件，判断返回方式
	acceptHeader := c.GetHeader("Accept")
	if strings.Contains(acceptHeader, "application/json") {
		// 返回 JSON 格式的文件信息
		fileInfoResult := FileInfo{
			Name:    fileInfo.Name(),
			IsDir:   false,
			Size:    fileInfo.Size(),
			ModTime: fileInfo.ModTime(),
		}
		
		c.JSON(http.StatusOK, gin.H{
			"path":     cleanPath,
			"isDir":    false,
			"fileInfo": fileInfoResult,
		})
		return
	}

	// 返回文件流
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileInfo.Name()))
	c.File(fullPath)
}

// FileInfo 文件信息结构体
type FileInfo struct {
	Name    string    `json:"name"`
	IsDir   bool      `json:"isDir"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
}

// IsPublic 检查路径是否为公共路径
func IsPublic(path string) bool {
	// 简单实现：检查路径是否以 /public 开头
	cleanPath := filepath.Clean(path)
	return strings.HasPrefix(cleanPath, "/public")
}

// ValidateAPIKey 验证 API Key
func ValidateAPIKey(apiKey string) bool {
	// 简单实现：检查 API Key 是否为 "test-key"
	// 实际实现应该查询数据库验证 API Key
	return apiKey == "test-key"
}