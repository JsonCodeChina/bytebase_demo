// Package mysql implements MySQL-specific SQL review rules.
package mysql

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/shenbo/sql-review-learning-demo/pkg/advisor"
)

func init() {
	// Register this advisor with the global registry
	advisor.Register(advisor.MySQL, advisor.MySQLTableRequirePK, &TableRequirePKAdvisor{})
}

// TableRequirePKAdvisor checks that tables have a primary key.
type TableRequirePKAdvisor struct{}

// TableRequirePKRule 表必须有主键规则 (legacy compatibility)
type TableRequirePKRule struct {
	*advisor.BaseRule
}

// NewTableRequirePKRule 创建表主键检查规则 (legacy compatibility)
func NewTableRequirePKRule() *TableRequirePKRule {
	return &TableRequirePKRule{
		BaseRule: &advisor.BaseRule{
			RuleID:    string(advisor.MySQLTableRequirePK),
			RuleName:  "表必须有主键",
			RuleDesc:  "每个表都应该有主键，以确保数据唯一性和复制一致性",
			RuleLevel: advisor.LevelError,
		},
	}
}

// Check implements the advisor.Advisor interface.
func (a *TableRequirePKAdvisor) Check(ctx context.Context, checkCtx advisor.Context) ([]*advisor.Advice, error) {
	var advices []*advisor.Advice

	// Get the rule level for advice severity
	level := advisor.NewStatusByRuleLevel(checkCtx.Rule.Level)

	// TODO: This will be replaced with ANTLR-based parsing
	// For now, we use the existing regex-based approach
	if isPossibleCreateTable(checkCtx.SQL) {
		tableName := extractTableName(checkCtx.SQL)
		if tableName != "" && !hasPrimaryKey(checkCtx.SQL) {
			advices = append(advices, &advisor.Advice{
				Status:  level,
				Code:    advisor.CodeTableNoPrimaryKey,
				Title:   "Table requires primary key",
				Content: fmt.Sprintf("Table `%s` requires PRIMARY KEY", tableName),
				// TODO: Add proper position information from AST
				StartPosition: &advisor.Position{Line: 1, Column: 1},
				// Legacy fields for compatibility
				Level:   advisor.Level(level),
				Message: fmt.Sprintf("Table `%s` requires PRIMARY KEY", tableName),
				RuleID:  string(advisor.MySQLTableRequirePK),
			})
		}
	}

	return advices, nil
}

// Check 执行规则检查 (legacy compatibility)
func (r *TableRequirePKRule) Check(ctx context.Context, checkCtx *advisor.Context) ([]*advisor.Advice, error) {
	var advices []*advisor.Advice

	// 解析SQL语句，查找CREATE TABLE语句
	sql := strings.ToUpper(strings.TrimSpace(checkCtx.SQL))

	// 简单的正则表达式匹配CREATE TABLE语句，使用(?s)标志允许.匹配换行符
	createTableRegex := regexp.MustCompile(`(?s)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?(?:` + "`" + `)?(\w+)(?:` + "`" + `)?\s*\((.*?)\)`)
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

// Helper functions for the new Advisor interface (will be replaced with ANTLR)

// isPossibleCreateTable checks if the SQL might be a CREATE TABLE statement.
func isPossibleCreateTable(sql string) bool {
	upperSQL := strings.ToUpper(strings.TrimSpace(sql))
	return strings.HasPrefix(upperSQL, "CREATE TABLE")
}

// extractTableName extracts the table name from a CREATE TABLE statement.
// This is a simplified implementation that will be replaced by AST traversal.
func extractTableName(sql string) string {
	// Very basic regex-like extraction - will be replaced with AST node analysis
	upperSQL := strings.ToUpper(sql)

	// Find "CREATE TABLE" and extract the next identifier
	if idx := strings.Index(upperSQL, "CREATE TABLE"); idx != -1 {
		remaining := strings.TrimSpace(sql[idx+12:]) // Skip "CREATE TABLE"

		// Extract first word as table name (simplified)
		parts := strings.Fields(remaining)
		if len(parts) > 0 {
			tableName := parts[0]
			// Remove backticks if present
			tableName = strings.Trim(tableName, "`")
			return tableName
		}
	}

	return ""
}

// hasPrimaryKey checks if the CREATE TABLE statement has a primary key.
// This will be replaced with proper AST analysis.
func hasPrimaryKey(sql string) bool {
	upperSQL := strings.ToUpper(sql)

	// Check for PRIMARY KEY keyword
	if strings.Contains(upperSQL, "PRIMARY KEY") {
		return true
	}

	// Check for column-level PRIMARY
	if strings.Contains(upperSQL, " PRIMARY") {
		return true
	}

	// Check for AUTO_INCREMENT (often implies primary key)
	if strings.Contains(upperSQL, "AUTO_INCREMENT") {
		return true
	}

	return false
}

// TODO: ANTLR Integration Plan
//
// When ANTLR integration is complete, this file will be restructured as follows:
//
// 1. Import ANTLR-generated MySQL parser types
// 2. Implement a Listener that extends BaseMySQLParserListener
// 3. Use the Listener pattern to traverse the AST and collect table information
//
// Example structure:
//
// type tableRequirePKChecker struct {
//     *mysql.BaseMySQLParserListener
//     baseLine   int
//     advices    []*advisor.Advice
//     level      advisor.Status
//     title      string
//     tables     map[string]bool // table name -> has primary key
//     positions  map[string]*advisor.Position
// }
//
// func (checker *tableRequirePKChecker) EnterCreateTable(ctx *mysql.CreateTableContext) {
//     tableName := extractTableNameFromContext(ctx)
//     checker.tables[tableName] = false
//     checker.positions[tableName] = extractPositionFromContext(ctx)
//
//     // Analyze table elements for primary key constraints
//     for _, element := range ctx.TableElementList().AllTableElement() {
//         if hasPrimaryKeyConstraint(element) {
//             checker.tables[tableName] = true
//             break
//         }
//     }
// }
//
// func (checker *tableRequirePKChecker) generateAdvices() []*advisor.Advice {
//     for tableName, hasPK := range checker.tables {
//         if !hasPK {
//             advice := &advisor.Advice{
//                 Status:        checker.level,
//                 Code:          advisor.CodeTableNoPrimaryKey,
//                 Title:         checker.title,
//                 Content:       fmt.Sprintf("Table `%s` requires PRIMARY KEY", tableName),
//                 StartPosition: checker.positions[tableName],
//             }
//             checker.advices = append(checker.advices, advice)
//         }
//     }
//     return checker.advices
// }