grammar SimpleSQL;

// 简化的SQL语法，只支持CREATE TABLE语句
// 用于演示如何准确解析SQL结构

// 语法规则 (Parser Rules)
statement : createTableStatement EOF ;

createTableStatement
    : CREATE TABLE tableName '(' columnDefinitionList ')'
    ;

columnDefinitionList
    : columnDefinition (',' columnDefinition)*
    ;

columnDefinition
    : columnName columnType columnConstraint*
    ;

columnConstraint
    : PRIMARY KEY              # PrimaryKeyConstraint
    | NOT NULL                 # NotNullConstraint
    | AUTO_INCREMENT           # AutoIncrementConstraint
    ;

columnType
    : INT
    | VARCHAR '(' NUMBER ')'
    | TEXT
    | TIMESTAMP
    ;

tableName : IDENTIFIER ;
columnName : IDENTIFIER ;

// 词法规则 (Lexer Rules)
CREATE : 'CREATE' | 'create' ;
TABLE : 'TABLE' | 'table' ;
PRIMARY : 'PRIMARY' | 'primary' ;
KEY : 'KEY' | 'key' ;
NOT : 'NOT' | 'not' ;
NULL : 'NULL' | 'null' ;
AUTO_INCREMENT : 'AUTO_INCREMENT' | 'auto_increment' ;

INT : 'INT' | 'int' ;
VARCHAR : 'VARCHAR' | 'varchar' ;
TEXT : 'TEXT' | 'text' ;
TIMESTAMP : 'TIMESTAMP' | 'timestamp' ;

IDENTIFIER : [a-zA-Z_][a-zA-Z0-9_]* ;
NUMBER : [0-9]+ ;

WS : [ \t\r\n]+ -> skip ;

// 注释：这个语法演示如何准确识别：
// 1. 表名和列名（通过IDENTIFIER）
// 2. PRIMARY KEY约束（通过专门的语法规则）
// 3. 其他约束和类型
//
// 对比正则表达式：
// - 正则表达式: if strings.Contains(sql, "PRIMARY KEY")
// - ANTLR: 只有真正的PRIMARY KEY约束才会触发规则