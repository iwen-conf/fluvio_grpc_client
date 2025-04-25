package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config 定义了应用程序的配置结构
type Config struct {
	Server ServerConfig `json:"server"`
}

// ServerConfig 定义了 gRPC 服务器的配置
type ServerConfig struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	ConnectTimeout int    `json:"connectTimeout"` // 连接超时时间（秒）
	CallTimeout    int    `json:"callTimeout"`    // 调用超时时间（秒）
}

// Load 从指定路径加载 JSON 配置文件
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败 %s: %w", path, err)
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败 %s: %w", path, err)
	}
	return &cfg, nil
}
