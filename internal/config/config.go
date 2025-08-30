package config

import (
	"os"
)

// Config 保存应用配置
type Config struct {
	RootDir   string
	Listen    string
	LogLevel  string
	JWTSecret string
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
	cfg := &Config{
		RootDir:  os.Getenv("ROOT_DIR"),
		Listen:   ":8080",
		LogLevel: "info",
	}
	
	if cfg.RootDir == "" {
		// 默认使用 configs 目录作为根目录
		cfg.RootDir = "./configs"
	}
	
	if listen := os.Getenv("LISTEN"); listen != "" {
		cfg.Listen = listen
	}
	
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}
	
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	
	return cfg
}