grammar Calculator;

// 语法规则（Parser Rules）- 以小写字母开头
// 程序入口：一个表达式后跟文件结束符
prog: expr EOF ;

// 表达式规则：支持四则运算和括号
// ANTLR 4 左递归规则优先级：从上到下优先级递减
expr
    : expr ('*'|'/') expr   # MulDiv      // 乘除运算，最高优先级
    | expr ('+'|'-') expr   # AddSub      // 加减运算，较低优先级
    | '(' expr ')'          # Parens      // 括号表达式
    | NUMBER                # Number      // 数字字面量
    ;

// 词法规则（Lexer Rules）- 以大写字母开头
// 数字：支持整数和小数
NUMBER : [0-9]+('.'[0-9]+)? ;

// 空白字符：跳过处理
WS : [ \t\r\n]+ -> skip ;