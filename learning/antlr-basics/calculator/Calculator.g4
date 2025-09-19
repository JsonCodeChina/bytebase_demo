grammar Calculator;

// 语法规则（Parser Rules）- 以小写字母开头
prog: expr EOF ;

expr
    : expr ('*'|'/') expr   # MulDiv
    | expr ('+'|'-') expr   # AddSub
    | '(' expr ')'          # Parens
    | NUMBER                # Number
    ;

// 词法规则（Lexer Rules）- 以大写字母开头
NUMBER : [0-9]+('.'[0-9]+)? ;

WS : [ \t\r\n]+ -> skip ;

// 注释：
// prog: 程序入口，一个表达式加文件结束符
// expr: 表达式规则，支持四则运算和括号
//   - MulDiv: 乘除运算，优先级高
//   - AddSub: 加减运算，优先级低
//   - Parens: 括号表达式
//   - Number: 数字
// NUMBER: 匹配整数或小数
// WS: 空白字符，跳过不处理