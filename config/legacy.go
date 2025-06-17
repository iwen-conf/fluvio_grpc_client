package config

import (
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/internal/config"
)

// FromLegacyConfig 从旧的配置格式转换为新的配置格式
func FromLegacyConfig(legacyCfg *config.Config) *Config {
	if legacyCfg == nil {
		return DefaultConfig()
	}

	cfg := DefaultConfig()

	// 转换服务器配置
	cfg.Server.Host = legacyCfg.Server.Host
	cfg.Server.Port = legacyCfg.Server.Port

	// 转换连接配置
	if legacyCfg.Server.ConnectTimeout > 0 {
		cfg.Connection.ConnectTimeout = time.Duration(legacyCfg.Server.ConnectTimeout) * time.Second
	}
	if legacyCfg.Server.CallTimeout > 0 {
		cfg.Connection.CallTimeout = time.Duration(legacyCfg.Server.CallTimeout) * time.Second
	}

	return cfg
}

// ToLegacyServerConfig 转换为旧的服务器配置格式
func (c *Config) ToLegacyServerConfig() *config.ServerConfig {
	return &config.ServerConfig{
		Host:           c.Server.Host,
		Port:           c.Server.Port,
		ConnectTimeout: int(c.Connection.ConnectTimeout.Seconds()),
		CallTimeout:    int(c.Connection.CallTimeout.Seconds()),
	}
}

// LoadLegacyConfig 加载旧格式的配置文件
func LoadLegacyConfig(path string) (*Config, error) {
	legacyCfg, err := config.Load(path)
	if err != nil {
		return nil, err
	}
	return FromLegacyConfig(legacyCfg), nil
}
