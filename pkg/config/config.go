package config

import (
	"time"
)

// Config 应用总配置
type Config struct {
	Server   ServerConfig   `yaml:"server" mapstructure:"server"`
	Database DatabaseConfig `yaml:"database" mapstructure:"database"`
	Logging  LoggingConfig  `yaml:"logging" mapstructure:"logging"`
	Rules    RulesConfig    `yaml:"rules" mapstructure:"rules"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int        `yaml:"port" mapstructure:"port" env:"SERVER_PORT"`
	Mode string     `yaml:"mode" mapstructure:"mode" env:"GIN_MODE"`
	CORS CORSConfig `yaml:"cors" mapstructure:"cors"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowOrigins []string `yaml:"allow_origins" mapstructure:"allow_origins"`
	AllowMethods []string `yaml:"allow_methods" mapstructure:"allow_methods"`
	AllowHeaders []string `yaml:"allow_headers" mapstructure:"allow_headers"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Pool            PoolConfig    `yaml:"pool" mapstructure:"pool"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout" mapstructure:"connection_timeout"`
}

// GetPoolConfig 实现DatabaseConfigInterface接口
func (d DatabaseConfig) GetPoolConfig() (int, int, time.Duration) {
	return d.Pool.MaxOpenConns, d.Pool.MaxIdleConns, d.Pool.ConnMaxLifetime
}

// GetConnectionTimeout 实现DatabaseConfigInterface接口
func (d DatabaseConfig) GetConnectionTimeout() time.Duration {
	return d.ConnectionTimeout
}

// PoolConfig 连接池配置
type PoolConfig struct {
	MaxOpenConns    int           `yaml:"max_open_conns" mapstructure:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level  string `yaml:"level" mapstructure:"level" env:"LOG_LEVEL"`
	Format string `yaml:"format" mapstructure:"format" env:"LOG_FORMAT"`
}

// RulesConfig 规则配置
type RulesConfig struct {
	MySQL MySQLRulesConfig `yaml:"mysql" mapstructure:"mysql"`
}

// MySQLRulesConfig MySQL规则配置
type MySQLRulesConfig struct {
	TableRequirePK     RuleConfig `yaml:"table_require_pk" mapstructure:"table_require_pk"`
	NamingConvention   RuleConfig `yaml:"naming_convention" mapstructure:"naming_convention"`
	StatementSafety    RuleConfig `yaml:"statement_safety" mapstructure:"statement_safety"`
	ColumnTypeCheck    RuleConfig `yaml:"column_type_check" mapstructure:"column_type_check"`
	SelectPerformance  RuleConfig `yaml:"select_performance" mapstructure:"select_performance"`
}

// RuleConfig 单个规则配置
type RuleConfig struct {
	Enabled bool                   `yaml:"enabled" mapstructure:"enabled"`
	Level   string                 `yaml:"level" mapstructure:"level"`
	Options map[string]interface{} `yaml:"options" mapstructure:"options"`
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: 8080,
			Mode: "debug",
			CORS: CORSConfig{
				AllowOrigins: []string{"*"},
				AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowHeaders: []string{"Content-Type", "Authorization"},
			},
		},
		Database: DatabaseConfig{
			Pool: PoolConfig{
				MaxOpenConns:    10,
				MaxIdleConns:    5,
				ConnMaxLifetime: time.Hour,
			},
			ConnectionTimeout: 5 * time.Second,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
		},
		Rules: RulesConfig{
			MySQL: MySQLRulesConfig{
				TableRequirePK: RuleConfig{
					Enabled: true,
					Level:   "ERROR",
					Options: make(map[string]interface{}),
				},
				NamingConvention: RuleConfig{
					Enabled: true,
					Level:   "WARNING",
					Options: map[string]interface{}{
						"table_pattern":  "^[a-z][a-z0-9_]*$",
						"column_pattern": "^[a-z][a-z0-9_]*$",
					},
				},
				StatementSafety: RuleConfig{
					Enabled: true,
					Level:   "ERROR",
					Options: make(map[string]interface{}),
				},
				ColumnTypeCheck: RuleConfig{
					Enabled: false,
					Level:   "WARNING",
					Options: map[string]interface{}{
						"forbidden_types": []string{"text", "blob"},
					},
				},
				SelectPerformance: RuleConfig{
					Enabled: true,
					Level:   "WARNING",
					Options: make(map[string]interface{}),
				},
			},
		},
	}
}