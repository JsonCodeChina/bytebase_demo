# ANTLR计算器学习示例

## 🎯 什么是ANTLR？

**ANTLR** (ANother Tool for Language Recognition) 是一个强大的语法分析器生成器，广泛用于构建编译器、解释器和各种语言处理工具。

## 📖 这个例子演示什么？

通过一个简单的计算器，演示ANTLR的核心工作流程：

1. **词法分析 (Lexer)**: 将 `2 + 3 * 4` 分解为 `[2, +, 3, *, 4]`
2. **语法分析 (Parser)**: 根据优先级规则构建语法树
3. **语义处理 (Listener)**: 遍历语法树计算结果

## 📝 语法文件解析 (Calculator.g4)

### 语法规则 (Parser Rules)

```antlr
prog: expr EOF ;                 // 程序 = 表达式 + 文件结束

expr
    : expr ('*'|'/') expr        // 乘除运算（高优先级）
    | expr ('+'|'-') expr        // 加减运算（低优先级）
    | '(' expr ')'               // 括号表达式
    | NUMBER                     // 数字
    ;
```

### 词法规则 (Lexer Rules)

```antlr
NUMBER : [0-9]+('.'[0-9]+)? ;    // 匹配整数或小数
WS : [ \t\r\n]+ -> skip ;        // 跳过空白字符
```

## ⚡ 优先级说明

ANTLR通过**规则顺序**自动定义运算符优先级：

- `*` 和 `/` 在前面 → **高优先级**
- `+` 和 `-` 在后面 → **低优先级**
- 括号 `()` 可以**改变优先级**

## 🔄 示例解析过程

### 输入表达式: `2 + 3 * 4`

**第1步：词法分析结果**
```
NUMBER(2) → PLUS → NUMBER(3) → MULTIPLY → NUMBER(4) → EOF
```

**第2步：语法树构建**
```
         expr (AddSub)
        /              \
   NUMBER(2)         expr (MulDiv)
                    /              \
               NUMBER(3)        NUMBER(4)
```

**第3步：计算过程**
1. 先计算右子树：`3 * 4 = 12`
2. 再计算根节点：`2 + 12 = 14`

## ❌ 对比正则表达式的局限性

### 正则表达式的问题

```go
// ❌ 无法正确处理优先级和括号
result := regexp.MustCompile(`(\d+)\s*([+\-*/])\s*(\d+)`)
// "2 + 3 * 4" 会被错误地解析为 (2 + 3) * 4 = 20
```

### ✅ ANTLR的优势

- ✅ **自动处理运算符优先级**
- ✅ **支持任意嵌套括号**
- ✅ **生成结构化的语法树**
- ✅ **易于扩展新的运算符**

## 🚀 运行方式

```bash
cd learning/antlr-basics/calculator
go run main.go
```

---

## 🔍 详细执行分析：`(2 + 3) * 4`

这个例子展示了**括号如何改变运算优先级**，最终结果是 **20** 而不是 **14**。

### 第1步：词法分析 (tokenize函数)

```go
func tokenize(input string) []Token {
    // 输入: "(2 + 3) * 4"
    // 逐个字符扫描，识别Token类型
}
```

#### 🔄 扫描过程：

```
'(' → Token{LPAREN, "("}
'2' → Token{NUMBER, "2"}
'+' → Token{PLUS, "+"}
'3' → Token{NUMBER, "3"}
')' → Token{RPAREN, ")"}
'*' → Token{MULTIPLY, "*"}
'4' → Token{NUMBER, "4"}
EOF → Token{EOF, ""}
```

#### ✅ 词法分析结果：

```
[LPAREN, NUMBER(2), PLUS, NUMBER(3), RPAREN, MULTIPLY, NUMBER(4), EOF]
```

### 第2步：语法分析 - 递归调用追踪

#### 📞 调用1：parseExpression() - 开始解析

```go
func (p *Parser) parseExpression() ASTNode {
    left := p.parseTerm()  // 调用parseTerm获取左操作数
    // ... 后续处理
}
```

🎯 **当前位置：** `(2 + 3) * 4`

#### 📞 调用2：parseTerm() - 处理乘法

```go
func (p *Parser) parseTerm() ASTNode {
    left := p.parseFactor()  // 调用parseFactor处理括号
    // 当前Token: LPAREN "("
}
```

🎯 **当前位置：** `(2 + 3) * 4`

#### 📞 调用3：parseFactor() - 🔥 关键的括号处理！

```go
func (p *Parser) parseFactor() ASTNode {
    token := p.currentToken()  // LPAREN "("

    if token.Type == LPAREN {
        p.advance()                    // 跳过 '('，位置移到 '2'
        expr := p.parseExpression()    // 🔥递归调用！解析括号内的内容
        p.advance()                    // 跳过 ')'
        return expr
    }
}
```

> **🚨 关键时刻！** 这里发生了递归调用，开始解析括号内的 `2 + 3`

#### 📞 调用4：parseExpression() - 🔁 递归处理括号内容

```go
// 🎯 现在解析的是括号内的 "2 + 3"
func (p *Parser) parseExpression() ASTNode {
    left := p.parseTerm()  // 获取 "2"

    // left = NumberNode(2)
    // 🔍 发现 PLUS "+"
    operator := "+"
    right := p.parseTerm()  // 获取 "3"
    // right = NumberNode(3)

    // ✅ 创建加法节点
    return BinaryOpNode(NumberNode(2) + NumberNode(3))
}
```

#### 🔙 回到调用3：parseFactor()

```go
func (p *Parser) parseFactor() ASTNode {
    // 📥 递归调用返回了 BinaryOpNode(2 + 3)
    expr := BinaryOpNode(2 + 3)  // 括号内的结果
    p.advance()  // 跳过 ')'，现在位置在 '*'
    return expr  // 📤 返回括号内的表达式
}
```

#### 🔙 回到调用2：parseTerm() - 处理乘法

```go
func (p *Parser) parseTerm() ASTNode {
    left := BinaryOpNode(2 + 3)  // 📥 从parseFactor得到的结果

    // 🔍 检查下一个token
    currentToken = "*"  // 发现乘号！
    p.advance()  // 跳过 '*'

    right := p.parseFactor()  // 获取 "4"
    right = NumberNode(4)

    // ✅ 创建乘法节点
    left = BinaryOpNode(BinaryOpNode(2 + 3) * NumberNode(4))
    return left
}
```

#### 🔙 回到调用1：parseExpression()

```go
func (p *Parser) parseExpression() ASTNode {
    left := BinaryOpNode((2 + 3) * 4)  // 📥 从parseTerm得到的结果

    // 🔍 检查是否还有 + 或 -
    // 已经到达EOF，没有更多操作符

    return left  // 📤 返回最终的AST
}
```

### 第3步：生成的AST结构

#### 🌳 `(2 + 3) * 4` 的AST：

```
         BinaryOpNode(*)
        /               \
  BinaryOpNode(+)    NumberNode(4)
   /           \
NumberNode(2)  NumberNode(3)
```

#### 🆚 对比没有括号的情况 `2 + 3 * 4`：

```
        BinaryOpNode(+)
       /              \
NumberNode(2)    BinaryOpNode(*)
                 /              \
           NumberNode(3)    NumberNode(4)
```

### 第4步：计算过程 (Evaluate)

#### 🔄 递归计算算法

```go
// 根节点：BinaryOpNode(*).Evaluate()
func (n *BinaryOpNode) Evaluate() float64 {
    left := n.Left.Evaluate()   // 📞 计算左子树
    right := n.Right.Evaluate() // 📞 计算右子树
    return left * right         // ✅ 执行运算
}
```

#### 📊 计算步骤：

1. **🔢 左子树** `BinaryOpNode(+).Evaluate()`:
   - `NumberNode(2).Evaluate() = 2.0`
   - `NumberNode(3).Evaluate() = 3.0`
   - `2.0 + 3.0 = 5.0` ✅

2. **🔢 右子树** `NumberNode(4).Evaluate()`:
   - `= 4.0` ✅

3. **🎯 根节点计算**:
   - `5.0 * 4.0 = 20.0` 🎉

---

## 🎯 关键洞察

### 💡 为什么括号改变了优先级？

1. **没有括号时**：`parseTerm` 先处理乘法，`parseExpression` 后处理加法
2. **有括号时**：`parseFactor` 中的递归调用**强制先计算括号内容**！

### ✨ 递归的魔法

```go
// 🔑 关键代码：
if token.Type == LPAREN {
    p.advance()                    // 跳过 '('
    expr := p.parseExpression()    // 🔥强制递归，优先处理括号内容
    p.advance()                    // 跳过 ')'
    return expr
}
```

> **💎 核心机制**：这个递归调用打破了正常的优先级顺序，让括号内的加法先于括号外的乘法执行！

### 🆚 与正则表达式的根本区别

#### ❌ 正则表达式无法处理嵌套结构：

```go
// 正则表达式完全无法理解括号的含义
regexp.MustCompile(`\((.*)\) \* (\d+)`)  // 能匹配，但无法正确计算
```

#### ✅ 递归解析自然处理任意嵌套：

```go
// 任意复杂的嵌套都能处理
((2 + 3) * (4 - 1)) + 5  // 完全没问题！
```

### 🚀 总结

这就是为什么复杂语法解析必须使用**递归算法**，而不能依赖简单的**字符串匹配**！

> **💡 重要启示**：这个原理直接适用于SQL解析 - 这是我们需要从正则表达式升级到ANTLR的根本原因。
