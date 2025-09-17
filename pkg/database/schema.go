package database

import (
	"database/sql"
	"fmt"
	"strings"
)

// Table 表结构信息
type Table struct {
	Name    string   `json:"name"`
	Engine  string   `json:"engine"`
	Comment string   `json:"comment"`
	Columns []Column `json:"columns"`
	Indexes []Index  `json:"indexes"`
}

// Column 列信息
type Column struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	IsNullable   bool   `json:"is_nullable"`
	DefaultValue string `json:"default_value"`
	Comment      string `json:"comment"`
	IsPrimaryKey bool   `json:"is_primary_key"`
	IsAutoIncr   bool   `json:"is_auto_increment"`
}

// Index 索引信息
type Index struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`    // PRIMARY, UNIQUE, INDEX
	Columns []string `json:"columns"`
}

// SchemaInfo 数据库schema信息
type SchemaInfo struct {
	DatabaseName string  `json:"database_name"`
	Tables       []Table `json:"tables"`
}

// SchemaManager schema管理器
type SchemaManager struct {
	db     *sql.DB
	engine string
}

// NewSchemaManager 创建schema管理器
func NewSchemaManager(db *sql.DB, engine string) *SchemaManager {
	return &SchemaManager{
		db:     db,
		engine: engine,
	}
}

// GetSchemaInfo 获取完整的schema信息
func (sm *SchemaManager) GetSchemaInfo(databaseName string) (*SchemaInfo, error) {
	tables, err := sm.GetTables(databaseName)
	if err != nil {
		return nil, err
	}

	schemaInfo := &SchemaInfo{
		DatabaseName: databaseName,
		Tables:       tables,
	}

	return schemaInfo, nil
}

// GetTables 获取所有表信息
func (sm *SchemaManager) GetTables(databaseName string) ([]Table, error) {
	switch sm.engine {
	case "mysql":
		return sm.getMySQLTables(databaseName)
	case "postgresql":
		return sm.getPostgreSQLTables(databaseName)
	default:
		return nil, fmt.Errorf("unsupported database engine: %s", sm.engine)
	}
}

// getMySQLTables 获取MySQL表信息
func (sm *SchemaManager) getMySQLTables(databaseName string) ([]Table, error) {
	query := `
		SELECT TABLE_NAME, ENGINE, TABLE_COMMENT
		FROM information_schema.TABLES
		WHERE TABLE_SCHEMA = ? AND TABLE_TYPE = 'BASE TABLE'
		ORDER BY TABLE_NAME`

	rows, err := sm.db.Query(query, databaseName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []Table
	for rows.Next() {
		var table Table
		if err := rows.Scan(&table.Name, &table.Engine, &table.Comment); err != nil {
			return nil, err
		}

		// 获取列信息
		columns, err := sm.getMySQLColumns(databaseName, table.Name)
		if err != nil {
			return nil, err
		}
		table.Columns = columns

		// 获取索引信息
		indexes, err := sm.getMySQLIndexes(databaseName, table.Name)
		if err != nil {
			return nil, err
		}
		table.Indexes = indexes

		tables = append(tables, table)
	}

	return tables, nil
}

// getMySQLColumns 获取MySQL列信息
func (sm *SchemaManager) getMySQLColumns(databaseName, tableName string) ([]Column, error) {
	query := `
		SELECT
			COLUMN_NAME,
			COLUMN_TYPE,
			IS_NULLABLE,
			IFNULL(COLUMN_DEFAULT, ''),
			IFNULL(COLUMN_COMMENT, ''),
			COLUMN_KEY = 'PRI',
			EXTRA LIKE '%auto_increment%'
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
		ORDER BY ORDINAL_POSITION`

	rows, err := sm.db.Query(query, databaseName, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []Column
	for rows.Next() {
		var column Column
		var isNullable string
		if err := rows.Scan(
			&column.Name,
			&column.Type,
			&isNullable,
			&column.DefaultValue,
			&column.Comment,
			&column.IsPrimaryKey,
			&column.IsAutoIncr,
		); err != nil {
			return nil, err
		}

		column.IsNullable = isNullable == "YES"
		columns = append(columns, column)
	}

	return columns, nil
}

// getMySQLIndexes 获取MySQL索引信息
func (sm *SchemaManager) getMySQLIndexes(databaseName, tableName string) ([]Index, error) {
	query := `
		SELECT
			INDEX_NAME,
			CASE
				WHEN INDEX_NAME = 'PRIMARY' THEN 'PRIMARY'
				WHEN NON_UNIQUE = 0 THEN 'UNIQUE'
				ELSE 'INDEX'
			END as INDEX_TYPE,
			GROUP_CONCAT(COLUMN_NAME ORDER BY SEQ_IN_INDEX) as COLUMNS
		FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
		GROUP BY INDEX_NAME, NON_UNIQUE
		ORDER BY INDEX_NAME`

	rows, err := sm.db.Query(query, databaseName, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexes []Index
	for rows.Next() {
		var index Index
		var columnsStr string
		if err := rows.Scan(&index.Name, &index.Type, &columnsStr); err != nil {
			return nil, err
		}

		index.Columns = strings.Split(columnsStr, ",")
		indexes = append(indexes, index)
	}

	return indexes, nil
}

// getPostgreSQLTables 获取PostgreSQL表信息 (简化实现)
func (sm *SchemaManager) getPostgreSQLTables(databaseName string) ([]Table, error) {
	query := `
		SELECT tablename, '' as engine, '' as comment
		FROM pg_tables
		WHERE schemaname = 'public'
		ORDER BY tablename`

	rows, err := sm.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []Table
	for rows.Next() {
		var table Table
		if err := rows.Scan(&table.Name, &table.Engine, &table.Comment); err != nil {
			return nil, err
		}

		// 简化实现，只获取表名
		// 完整实现需要查询pg_catalog获取详细信息
		tables = append(tables, table)
	}

	return tables, nil
}

// GenerateDDL 生成DDL语句
func (sm *SchemaManager) GenerateDDL(tableName string, databaseName string) (string, error) {
	if sm.engine != "mysql" {
		return "", fmt.Errorf("DDL generation currently only supports MySQL")
	}

	query := fmt.Sprintf("SHOW CREATE TABLE `%s`.`%s`", databaseName, tableName)
	var table, ddl string
	err := sm.db.QueryRow(query).Scan(&table, &ddl)
	if err != nil {
		return "", err
	}

	return ddl, nil
}