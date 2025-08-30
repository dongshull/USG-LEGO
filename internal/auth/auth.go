package auth

import (
	"path/filepath"
	"strings"
)

// IsPublic 检查路径是否为公共路径
func IsPublic(path string) bool {
	// 简单实现：检查路径是否以 /public 开头
	cleanPath := filepath.Clean(path)
	return strings.HasPrefix(cleanPath, "/public")
}

// InitAuth 初始化认证模块
func InitAuth(rootDir string) {
	// 这里应该读取 rootDir/.usg-lego.yml 文件并解析 public/private 路径规则
	// 为简化测试，我们暂时不实现这个功能
}