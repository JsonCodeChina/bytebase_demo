package advisor

import (
	"context"
	"database/sql"
)

// Level 建议级别
type Level string

const (
	LevelError   Level = "ERROR"
	LevelWarning Level = "WARNING"
	LevelInfo    Level = "INFO"
)

// Advice 审查建议
type Advice struct {
	Level   Level  `json:"level"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Line    int    `json:"line,omitempty"`
	Column  int    `json:"column,omitempty"`
	RuleID  string `json:"rule_id"`
}

// Context 审查上下文
type Context struct {
	SQL          string           `json:"sql"`
	Engine       string           `json:"engine"`
	DatabaseName string           `json:"database_name"`
	Rules        []string         `json:"rules"`        // 启用的规则ID列表
	Connection   *sql.DB          `json:"-"`            // 数据库连接
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

// Advisor SQL审查器接口
type Advisor interface {
	Check(ctx context.Context, checkCtx *Context) ([]*Advice, error)
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