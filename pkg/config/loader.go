package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Loader 配置加载器
type Loader struct {
	configDir string
	env       string
}

// NewLoader 创建配置加载器
func NewLoader(configDir string) *Loader {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	return &Loader{
		configDir: configDir,
		env:       env,
	}
}

// GetEnv 获取当前环境
func (l *Loader) GetEnv() string {
	return l.env
}

// Load 加载配置
func (l *Loader) Load() (*Config, error) {
	// 从默认配置开始
	config := GetDefaultConfig()

	// 加载主配置文件
	if err := l.loadYAMLFile("app.yaml", config); err != nil {
		return nil, fmt.Errorf("failed to load app.yaml: %w", err)
	}

	// 加载环境特定配置
	envFile := fmt.Sprintf("%s.yaml", l.env)
	if err := l.loadYAMLFile(envFile, config); err != nil {
		// 环境配置文件不存在不是错误
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load %s: %w", envFile, err)
		}
	}

	// 加载规则配置
	if err := l.loadRulesConfig(config); err != nil {
		return nil, fmt.Errorf("failed to load rules config: %w", err)
	}

	// 应用环境变量覆盖
	if err := l.applyEnvOverrides(config); err != nil {
		return nil, fmt.Errorf("failed to apply environment overrides: %w", err)
	}

	// 验证配置
	if err := l.validateConfig(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// loadYAMLFile 加载YAML文件
func (l *Loader) loadYAMLFile(filename string, config *Config) error {
	filepath := filepath.Join(l.configDir, filename)

	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, config)
}

// loadRulesConfig 加载规则配置
func (l *Loader) loadRulesConfig(config *Config) error {
	rulesFile := filepath.Join(l.configDir, "rules.yaml")

	data, err := os.ReadFile(rulesFile)
	if err != nil {
		if os.IsNotExist(err) {
			// 规则配置文件不存在时使用默认配置
			return nil
		}
		return err
	}

	// 只更新规则部分
	var rulesConfig struct {
		Rules RulesConfig `yaml:"rules"`
	}

	if err := yaml.Unmarshal(data, &rulesConfig); err != nil {
		return err
	}

	config.Rules = rulesConfig.Rules
	return nil
}

// applyEnvOverrides 应用环境变量覆盖
func (l *Loader) applyEnvOverrides(config *Config) error {
	return l.applyEnvToStruct(reflect.ValueOf(config).Elem())
}

// applyEnvToStruct 递归应用环境变量到结构体
func (l *Loader) applyEnvToStruct(v reflect.Value) error {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// 跳过非导出字段
		if !field.CanSet() {
			continue
		}

		// 处理嵌套结构体
		if field.Kind() == reflect.Struct {
			if err := l.applyEnvToStruct(field); err != nil {
				return err
			}
			continue
		}

		// 获取env标签
		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			continue
		}

		// 获取环境变量值
		envValue := os.Getenv(envTag)
		if envValue == "" {
			continue
		}

		// 根据字段类型设置值
		if err := l.setFieldValue(field, envValue); err != nil {
			return fmt.Errorf("failed to set field %s from env %s: %w", fieldType.Name, envTag, err)
		}
	}

	return nil
}

// setFieldValue 设置字段值
func (l *Loader) setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(intVal))
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.String {
			// 处理字符串切片，用逗号分隔
			values := strings.Split(value, ",")
			for i, v := range values {
				values[i] = strings.TrimSpace(v)
			}
			field.Set(reflect.ValueOf(values))
		}
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}

// validateConfig 验证配置
func (l *Loader) validateConfig(config *Config) error {
	// 验证服务器配置
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	if config.Server.Mode != "debug" && config.Server.Mode != "release" {
		return fmt.Errorf("invalid server mode: %s (must be 'debug' or 'release')", config.Server.Mode)
	}

	// 验证数据库配置
	if config.Database.Pool.MaxOpenConns <= 0 {
		return fmt.Errorf("invalid max_open_conns: %d", config.Database.Pool.MaxOpenConns)
	}

	if config.Database.Pool.MaxIdleConns < 0 {
		return fmt.Errorf("invalid max_idle_conns: %d", config.Database.Pool.MaxIdleConns)
	}

	if config.Database.ConnectionTimeout <= 0 {
		return fmt.Errorf("invalid connection_timeout: %v", config.Database.ConnectionTimeout)
	}

	// 验证日志配置
	validLevels := []string{"debug", "info", "warn", "error"}
	validLevel := false
	for _, level := range validLevels {
		if config.Logging.Level == level {
			validLevel = true
			break
		}
	}
	if !validLevel {
		return fmt.Errorf("invalid log level: %s (must be one of: %v)", config.Logging.Level, validLevels)
	}

	return nil
}

// LoadFromFile 从指定文件加载配置（用于测试）
func LoadFromFile(filename string) (*Config, error) {
	config := GetDefaultConfig()

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}