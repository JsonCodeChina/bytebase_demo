// Package advisor defines the interface for analyzing SQL statements.
// This is a simplified but functional implementation inspired by Bytebase.
package advisor

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

// Engine represents the database engine type.
type Engine string

const (
	// MySQL represents MySQL database engine.
	MySQL Engine = "MYSQL"
	// PostgreSQL represents PostgreSQL database engine.
	PostgreSQL Engine = "POSTGRESQL"
)

// Level represents the severity level of a rule.
type Level string

const (
	LevelError   Level = "ERROR"
	LevelWarning Level = "WARNING"
	LevelInfo    Level = "INFO"
)

// Status represents the status of advice (maps to Level for compatibility).
type Status = Level

// Position represents a position in the SQL text.
type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// Advice represents a piece of advice from SQL review.
type Advice struct {
	Status        Status    `json:"status"`
	Code          int32     `json:"code"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	StartPosition *Position `json:"startPosition,omitempty"`
	// Legacy fields for compatibility
	Level   Level  `json:"level"`
	Message string `json:"message"`
	Line    int    `json:"line,omitempty"`
	Column  int    `json:"column,omitempty"`
	RuleID  string `json:"rule_id"`
}

// Context encapsulates the checking context and configuration.
type Context struct {
	// Core fields
	SQL    string  `json:"sql"`
	Engine Engine  `json:"engine"`
	Rule   *Rule   `json:"rule"`
	AST    any     `json:"-"` // Parsed abstract syntax tree

	// Legacy compatibility fields
	DatabaseName string            `json:"database_name"`
	Rules        []string          `json:"rules"`        // 启用的规则ID列表
	Connection   *sql.DB           `json:"-"`            // 数据库连接
	Metadata     *DatabaseMetadata `json:"metadata"`     // 数据库元数据
}

// DatabaseMetadata 数据库元数据
type DatabaseMetadata struct {
	Tables map[string]*TableMetadata `json:"tables"`
}

// TableMetadata 表元数据
type TableMetadata struct {
	Name    string                    `json:"name"`
	Columns map[string]*ColumnMetadata `json:"columns"`
	Indexes map[string]*IndexMetadata  `json:"indexes"`
}

// ColumnMetadata 列元数据
type ColumnMetadata struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	IsNullable   bool   `json:"is_nullable"`
	IsPrimaryKey bool   `json:"is_primary_key"`
	IsAutoIncr   bool   `json:"is_auto_increment"`
}

// IndexMetadata 索引元数据
type IndexMetadata struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Columns []string `json:"columns"`
}

// Advisor is the interface for SQL statement analysis.
// Each rule should implement this interface.
type Advisor interface {
	Check(ctx context.Context, checkCtx Context) ([]*Advice, error)
}

// Rule 规则接口
type Rule interface {
	ID() string
	Name() string
	Description() string
	Level() Level
	Check(ctx context.Context, checkCtx *Context) ([]*Advice, error)
}

// BaseRule 基础规则实现
type BaseRule struct {
	RuleID      string
	RuleName    string
	RuleDesc    string
	RuleLevel   Level
}

func (r *BaseRule) ID() string {
	return r.RuleID
}

func (r *BaseRule) Name() string {
	return r.RuleName
}

func (r *BaseRule) Description() string {
	return r.RuleDesc
}

func (r *BaseRule) Level() Level {
	return r.RuleLevel
}

// DefaultAdvisor 默认审查器实现
type DefaultAdvisor struct {
	rules map[string]Rule
}

// NewDefaultAdvisor 创建默认审查器
func NewDefaultAdvisor() *DefaultAdvisor {
	return &DefaultAdvisor{
		rules: make(map[string]Rule),
	}
}

// RegisterRule 注册规则
func (a *DefaultAdvisor) RegisterRule(rule Rule) {
	a.rules[rule.ID()] = rule
}

// GetRule 获取规则
func (a *DefaultAdvisor) GetRule(id string) (Rule, bool) {
	rule, exists := a.rules[id]
	return rule, exists
}

// ListRules 列出所有规则
func (a *DefaultAdvisor) ListRules() []Rule {
	rules := make([]Rule, 0, len(a.rules))
	for _, rule := range a.rules {
		rules = append(rules, rule)
	}
	return rules
}

// Check 执行SQL审查
func (a *DefaultAdvisor) Check(ctx context.Context, checkCtx *Context) ([]*Advice, error) {
	var allAdvices []*Advice

	// 如果没有指定规则，使用所有规则
	rulesToCheck := checkCtx.Rules
	if len(rulesToCheck) == 0 {
		for ruleID := range a.rules {
			rulesToCheck = append(rulesToCheck, ruleID)
		}
	}

	// 执行每个规则
	for _, ruleID := range rulesToCheck {
		rule, exists := a.rules[ruleID]
		if !exists {
			continue
		}

		advices, err := rule.Check(ctx, checkCtx)
		if err != nil {
			// 记录错误但继续执行其他规则
			continue
		}

		allAdvices = append(allAdvices, advices...)
	}

	return allAdvices, nil
}

// Type represents the type identifier for an advisor.
type Type string

// Rule represents a SQL review rule configuration.
type Rule struct {
	Type    Type   `yaml:"type"`
	Level   Level  `yaml:"level"`
	Engine  Engine `yaml:"engine"`
	Payload any    `yaml:"payload,omitempty"`
}

// Global registry for advisors (Bytebase-style).
var (
	advisorMu sync.RWMutex
	advisors  = make(map[Engine]map[Type]Advisor)
)

// Register makes an advisor available by the provided engine and type.
// If Register is called twice with the same combination or if advisor is nil,
// it panics.
func Register(engine Engine, advType Type, advisor Advisor) {
	advisorMu.Lock()
	defer advisorMu.Unlock()

	if advisor == nil {
		panic("advisor: Register advisor is nil")
	}

	engineAdvisors, ok := advisors[engine]
	if !ok {
		advisors[engine] = map[Type]Advisor{
			advType: advisor,
		}
	} else {
		if _, dup := engineAdvisors[advType]; dup {
			panic(fmt.Sprintf("advisor: Register called twice for advisor %v for %v", advType, engine))
		}
		engineAdvisors[advType] = advisor
	}
}

// CheckByType runs the specified advisor and returns the advice list.
func CheckByType(ctx context.Context, engine Engine, advType Type, checkCtx Context) (adviceList []*Advice, err error) {
	// Panic recovery for safer execution
	defer func() {
		if panicErr := recover(); panicErr != nil {
			if panicError, ok := panicErr.(error); ok {
				err = errors.Errorf("advisor check PANIC RECOVER, type: %v, err: %v", advType, panicError)
			} else {
				err = errors.Errorf("advisor check PANIC RECOVER, type: %v, err: %v", advType, panicErr)
			}
		}
	}()

	advisorMu.RLock()
	engineAdvisors, ok := advisors[engine]
	defer advisorMu.RUnlock()

	if !ok {
		return nil, errors.Errorf("advisor: unknown engine type %v", engine)
	}

	advisor, ok := engineAdvisors[advType]
	if !ok {
		return nil, errors.Errorf("advisor: unknown advisor %v for %v", advType, engine)
	}

	return advisor.Check(ctx, checkCtx)
}

// GetRegisteredAdvisors returns all registered advisor types for a given engine.
func GetRegisteredAdvisors(engine Engine) []Type {
	advisorMu.RLock()
	defer advisorMu.RUnlock()

	engineAdvisors, ok := advisors[engine]
	if !ok {
		return nil
	}

	var types []Type
	for advType := range engineAdvisors {
		types = append(types, advType)
	}
	return types
}

// NewStatusByRuleLevel converts rule level to advice status.
func NewStatusByRuleLevel(level Level) Status {
	return Status(level)
}