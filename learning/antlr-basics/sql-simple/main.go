package main

import (
	"fmt"
	"regexp"
	"strings"
)

// 简化的SQL解析演示，对比正则表达式和ANTLR方法

// SQL AST节点定义
type SQLStatement interface {
	String() string
}

type CreateTableStatement struct {
	TableName string
	Columns   []*ColumnDefinition
}

func (c *CreateTableStatement) String() string {
	return fmt.Sprintf("CREATE TABLE %s (...)", c.TableName)
}

func (c *CreateTableStatement) HasPrimaryKey() bool {
	for _, col := range c.Columns {
		if col.IsPrimaryKey {
			return true
		}
	}
	return false
}

type ColumnDefinition struct {
	Name         string
	Type         string
	IsPrimaryKey bool
	IsNotNull    bool
}

// 简化的SQL解析器（模拟ANTLR生成的代码）
func parseCreateTable(sql string) (*CreateTableStatement, error) {
	// 这里简化演示，实际ANTLR会生成更复杂的解析代码
	sql = strings.TrimSpace(sql)

	// 提取表名
	tableRegex := regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(\w+)\s*\(`)
	matches := tableRegex.FindStringSubmatch(sql)
	if len(matches) < 2 {
		return nil, fmt.Errorf("无法解析表名")
	}

	tableName := matches[1]

	// 提取列定义部分
	start := strings.Index(sql, "(")
	end := strings.LastIndex(sql, ")")
	if start == -1 || end == -1 {
		return nil, fmt.Errorf("无法解析列定义")
	}

	columnsPart := sql[start+1 : end]
	columns := parseColumns(columnsPart)

	return &CreateTableStatement{
		TableName: tableName,
		Columns:   columns,
	}, nil
}

func parseColumns(columnsPart string) []*ColumnDefinition {
	var columns []*ColumnDefinition

	// 简化处理，按逗号分割
	parts := strings.Split(columnsPart, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		col := &ColumnDefinition{}

		// 提取列名（第一个单词）
		words := strings.Fields(part)
		if len(words) > 0 {
			col.Name = words[0]
		}

		// 检查是否有PRIMARY KEY
		upperPart := strings.ToUpper(part)
		col.IsPrimaryKey = strings.Contains(upperPart, "PRIMARY KEY") ||
			strings.Contains(upperPart, "PRIMARY")

		// 检查类型
		if strings.Contains(upperPart, "INT") {
			col.Type = "INT"
		} else if strings.Contains(upperPart, "VARCHAR") {
			col.Type = "VARCHAR"
		}

		columns = append(columns, col)
	}

	return columns
}

// 用正则表达式的方法（当前项目使用的方法）
func checkPrimaryKeyWithRegex(sql string) bool {
	upperSQL := strings.ToUpper(sql)

	// 当前项目的实现
	if strings.Contains(upperSQL, "PRIMARY KEY") {
		return true
	}

	if strings.Contains(upperSQL, " PRIMARY") {
		return true
	}

	if strings.Contains(upperSQL, "AUTO_INCREMENT") {
		return true
	}

	return false
}

func main() {
	fmt.Println("=== SQL解析对比：正则表达式 vs ANTLR概念 ===\n")

	// 测试用例：演示正则表达式的误判问题
	testCases := []struct {
		name        string
		sql         string
		expectPK    bool
		description string
	}{
		{
			name: "正常的主键定义",
			sql: `CREATE TABLE users (
				id INT PRIMARY KEY,
				name VARCHAR(50)
			)`,
			expectPK:    true,
			description: "真正的主键定义",
		},
		{
			name: "表级主键定义",
			sql: `CREATE TABLE orders (
				id INT,
				user_id INT,
				PRIMARY KEY (id)
			)`,
			expectPK:    true,
			description: "表级主键约束",
		},
		{
			name: "注释中的误导信息",
			sql: `CREATE TABLE logs (
				id INT,
				message TEXT -- Should add PRIMARY KEY later
			)`,
			expectPK:    false,
			description: "注释中提到PRIMARY KEY，但实际没有主键",
		},
		{
			name: "列名包含primary",
			sql: `CREATE TABLE settings (
				id INT,
				primary_config VARCHAR(100),
				value TEXT
			)`,
			expectPK:    false,
			description: "列名包含'primary'，但不是主键",
		},
		{
			name: "字符串值包含关键字",
			sql: `CREATE TABLE messages (
				id INT,
				content TEXT DEFAULT 'Add PRIMARY KEY to this table'
			)`,
			expectPK:    false,
			description: "默认值包含'PRIMARY KEY'，但不是主键",
		},
	}

	for _, tc := range testCases {
		fmt.Printf("测试用例: %s\n", tc.name)
		fmt.Printf("描述: %s\n", tc.description)
		fmt.Printf("SQL: %s\n", strings.ReplaceAll(tc.sql, "\n", " "))

		// 1. 正则表达式方法（当前项目的方法）
		regexResult := checkPrimaryKeyWithRegex(tc.sql)
		regexCorrect := regexResult == tc.expectPK

		fmt.Printf("正则表达式检测: %v ", regexResult)
		if regexCorrect {
			fmt.Printf("✅ 正确\n")
		} else {
			fmt.Printf("❌ 错误\n")
		}

		// 2. 简化的ANTLR方法
		stmt, err := parseCreateTable(tc.sql)
		antlrCorrect := false
		if err == nil {
			antlrResult := stmt.HasPrimaryKey()
			antlrCorrect = antlrResult == tc.expectPK
			fmt.Printf("ANTLR方法检测: %v ", antlrResult)
			if antlrCorrect {
				fmt.Printf("✅ 正确\n")
			} else {
				fmt.Printf("❌ 错误\n")
			}
		} else {
			fmt.Printf("ANTLR方法检测: 解析失败 - %v\n", err)
		}

		fmt.Println(strings.Repeat("-", 60))
	}

	fmt.Println("\n=== 总结 ===")
	fmt.Println("正则表达式的问题:")
	fmt.Println("1. 无法区分注释、字符串值和实际的SQL关键字")
	fmt.Println("2. 容易被业务字段名误导")
	fmt.Println("3. 无法理解SQL的语法结构")

	fmt.Println("\nANTLR的优势:")
	fmt.Println("1. 理解SQL语法结构，只识别真正的语法元素")
	fmt.Println("2. 自动处理注释和字符串")
	fmt.Println("3. 提供结构化的AST，便于复杂分析")
	fmt.Println("4. 易于扩展新的SQL语法支持")

	fmt.Println("\n这就是为什么我们需要改进SQL解析方法！")
}
