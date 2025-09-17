package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// ConnectionConfig 数据库连接配置
type ConnectionConfig struct {
	ID       string `json:"id" yaml:"id"`
	Name     string `json:"name" yaml:"name"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Database string `json:"database" yaml:"database"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Engine   string `json:"engine" yaml:"engine"` // mysql, postgresql
}

// DatabaseManager 数据库连接管理器
type DatabaseManager struct {
	connections     map[string]*sql.DB
	configs         map[string]*ConnectionConfig
	poolConfig      PoolConfig
	connTimeout     time.Duration
}

// PoolConfig 连接池配置
type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// NewDatabaseManager 创建数据库管理器
func NewDatabaseManager() *DatabaseManager {
	return &DatabaseManager{
		connections: make(map[string]*sql.DB),
		configs:     make(map[string]*ConnectionConfig),
		poolConfig: PoolConfig{
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: time.Hour,
		},
		connTimeout: 5 * time.Second,
	}
}

// DatabaseConfigInterface 数据库配置接口
type DatabaseConfigInterface interface {
	GetPoolConfig() (int, int, time.Duration)
	GetConnectionTimeout() time.Duration
}

// NewDatabaseManagerWithConfig 使用配置创建数据库管理器
func NewDatabaseManagerWithConfig(dbConfig DatabaseConfigInterface) *DatabaseManager {
	dm := NewDatabaseManager()

	maxOpen, maxIdle, maxLifetime := dbConfig.GetPoolConfig()
	dm.poolConfig = PoolConfig{
		MaxOpenConns:    maxOpen,
		MaxIdleConns:    maxIdle,
		ConnMaxLifetime: maxLifetime,
	}
	dm.connTimeout = dbConfig.GetConnectionTimeout()

	return dm
}

// TestConnection 测试数据库连接
func (dm *DatabaseManager) TestConnection(config *ConnectionConfig) error {
	dsn, err := dm.buildDSN(config)
	if err != nil {
		return err
	}

	db, err := sql.Open(config.Engine, dsn)
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}
	defer db.Close()

	// 设置连接超时
	db.SetConnMaxLifetime(dm.connTimeout)
	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// AddConnection 添加数据库连接
func (dm *DatabaseManager) AddConnection(config *ConnectionConfig) error {
	if err := dm.TestConnection(config); err != nil {
		return err
	}

	dsn, err := dm.buildDSN(config)
	if err != nil {
		return err
	}

	db, err := sql.Open(config.Engine, dsn)
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}

	// 配置连接池
	db.SetMaxOpenConns(dm.poolConfig.MaxOpenConns)
	db.SetMaxIdleConns(dm.poolConfig.MaxIdleConns)
	db.SetConnMaxLifetime(dm.poolConfig.ConnMaxLifetime)

	dm.connections[config.ID] = db
	dm.configs[config.ID] = config

	return nil
}

// GetConnection 获取数据库连接
func (dm *DatabaseManager) GetConnection(id string) (*sql.DB, error) {
	db, exists := dm.connections[id]
	if !exists {
		return nil, fmt.Errorf("connection not found: %s", id)
	}
	return db, nil
}

// GetConfig 获取连接配置
func (dm *DatabaseManager) GetConfig(id string) (*ConnectionConfig, error) {
	config, exists := dm.configs[id]
	if !exists {
		return nil, fmt.Errorf("config not found: %s", id)
	}
	return config, nil
}

// ListConnections 列出所有连接
func (dm *DatabaseManager) ListConnections() []*ConnectionConfig {
	configs := make([]*ConnectionConfig, 0, len(dm.configs))
	for _, config := range dm.configs {
		configs = append(configs, config)
	}
	return configs
}

// RemoveConnection 移除连接
func (dm *DatabaseManager) RemoveConnection(id string) error {
	if db, exists := dm.connections[id]; exists {
		db.Close()
		delete(dm.connections, id)
	}
	delete(dm.configs, id)
	return nil
}

// Close 关闭所有连接
func (dm *DatabaseManager) Close() error {
	for _, db := range dm.connections {
		if err := db.Close(); err != nil {
			return err
		}
	}
	dm.connections = make(map[string]*sql.DB)
	dm.configs = make(map[string]*ConnectionConfig)
	return nil
}

// buildDSN 构建数据库连接字符串
func (dm *DatabaseManager) buildDSN(config *ConnectionConfig) (string, error) {
	switch config.Engine {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Username, config.Password, config.Host, config.Port, config.Database), nil
	case "postgresql":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.Host, config.Port, config.Username, config.Password, config.Database), nil
	default:
		return "", fmt.Errorf("unsupported database engine: %s", config.Engine)
	}
}