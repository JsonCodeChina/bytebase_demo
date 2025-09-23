package advisor

// MySQL advisor types inspired by Bytebase.
const (
	// MySQLTableRequirePK is an advisor type for MySQL table require primary key.
	MySQLTableRequirePK Type = "mysql.table.require-pk"

	// MySQLNamingConvention is an advisor type for MySQL naming convention.
	MySQLNamingConvention Type = "mysql.naming.convention"

	// MySQLColumnTypeCheck is an advisor type for MySQL column type checking.
	MySQLColumnTypeCheck Type = "mysql.column.type-check"

	// MySQLStatementSafety is an advisor type for MySQL statement safety.
	MySQLStatementSafety Type = "mysql.statement.safety"

	// MySQLSelectPerformance is an advisor type for MySQL SELECT performance.
	MySQLSelectPerformance Type = "mysql.select.performance"
)

// Error codes for advisor checks.
const (
	// Success codes
	CodeOK int32 = 0

	// Table related error codes (800 range)
	CodeTableNoPrimaryKey   int32 = 801
	CodeTableNamingInvalid  int32 = 802
	CodeTableCreateRequired int32 = 803

	// Column related error codes (900 range)
	CodeColumnTypeDisallowed int32 = 901
	CodeColumnNamingInvalid  int32 = 902
	CodeColumnRequireDefault int32 = 903

	// Statement related error codes (1000 range)
	CodeStatementSelectAll        int32 = 1001
	CodeStatementNoWhere          int32 = 1002
	CodeStatementUnsafeOperation  int32 = 1003
	CodeStatementPerformanceIssue int32 = 1004
)