package mysql

import (
	"context"
	"regexp"
	"strings"

	"github.com/shenbo/sql-review-learning-demo/pkg/advisor"
)

const (
	MySQLRuleTableRequirePK = "mysql.table.require-pk"
)

// TableRequirePKRule 表必须有主键规则
type TableRequirePKRule struct {
	*advisor.BaseRule
}

// NewTableRequirePKRule 创建表主键检查规则
func NewTableRequirePKRule() *TableRequirePKRule {
	return &TableRequirePKRule{
		BaseRule: &advisor.BaseRule{
			RuleID:    MySQLRuleTableRequirePK,
			RuleName:  "表必须有主键",
			RuleDesc:  "每个表都应该有主键，以确保数据唯一性和复制一致性",
			RuleLevel: advisor.LevelError,
		},
	}
}

// Check 执行规则检查
func (r *TableRequirePKRule) Check(ctx context.Context, checkCtx *advisor.Context) ([]*advisor.Advice, error) {
	var advices []*advisor.Advice

	// 解析SQL语句，查找CREATE TABLE语句
	sql := strings.ToUpper(strings.TrimSpace(checkCtx.SQL))

	// 简单的正则表达式匹配CREATE TABLE语句
	createTableRegex := regexp.MustCompile(`CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?(?:` + "`" + `)?(\w+)(?:` + "`" + `)?\s*\((.*?)\)`)
	matches := createTableRegex.FindAllStringSubmatch(sql, -1)

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		tableName := match[1]
		tableDefinition := match[2]

		// 检查是否有主键定义
		if !r.hasPrimaryKey(tableDefinition) {
			advice := &advisor.Advice{
				Level:   r.Level(),
				Title:   r.Name(),
				Message: "表 '" + tableName + "' 缺少主键定义。建议添加主键以确保数据唯一性。",
				RuleID:  r.ID(),
			}
			advices = append(advices, advice)
		}
	}

	return advices, nil
}

// hasPrimaryKey 检查表定义中是否包含主键
func (r *TableRequirePKRule) hasPrimaryKey(tableDefinition string) bool {
	upperDef := strings.ToUpper(tableDefinition)

	// 检查PRIMARY KEY关键字
	if strings.Contains(upperDef, "PRIMARY KEY") {
		return true
	}

	// 检查列定义中的PRIMARY KEY
	if strings.Contains(upperDef, " PRIMARY") {
		return true
	}

	// 检查AUTO_INCREMENT列（通常暗示主键）
	lines := strings.Split(upperDef, ",")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "AUTO_INCREMENT") {
			// AUTO_INCREMENT列通常是主键
			return true
		}
	}

	return false
}