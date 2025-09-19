package main

import (
	"fmt"
	"strconv"
	"strings"
)

// 这是一个简化的ANTLR概念演示，展示ANTLR的工作原理
// 在真实项目中，这些代码会由ANTLR工具自动生成

// Token类型 - 对应ANTLR的词法分析结果
type TokenType int

const (
	NUMBER TokenType = iota
	PLUS
	MINUS
	MULTIPLY
	DIVIDE
	LPAREN
	RPAREN
	EOF
)

// Token - 词法分析的最小单位
type Token struct {
	Type  TokenType
	Value string
}

// 简化的词法分析器 - 对应ANTLR的Lexer
func tokenize(input string) []Token {
	var tokens []Token
	input = strings.ReplaceAll(input, " ", "") // 移除空格

	for i := 0; i < len(input); i++ {
		char := input[i]
		switch char {
		case '+':
			tokens = append(tokens, Token{PLUS, "+"})
		case '-':
			tokens = append(tokens, Token{MINUS, "-"})
		case '*':
			tokens = append(tokens, Token{MULTIPLY, "*"})
		case '/':
			tokens = append(tokens, Token{DIVIDE, "/"})
		case '(':
			tokens = append(tokens, Token{LPAREN, "("})
		case ')':
			tokens = append(tokens, Token{RPAREN, ")"})
		default:
			if char >= '0' && char <= '9' {
				// 收集完整的数字
				num := ""
				for i < len(input) && ((input[i] >= '0' && input[i] <= '9') || input[i] == '.') {
					num += string(input[i])
					i++
				}
				i-- // 回退一步
				tokens = append(tokens, Token{NUMBER, num})
			}
		}
	}
	tokens = append(tokens, Token{EOF, ""})
	return tokens
}

// AST节点接口 - 对应ANTLR的语法树节点
type ASTNode interface {
	Evaluate() float64
	String() string
}

// 数字节点
type NumberNode struct {
	Value float64
}

func (n *NumberNode) Evaluate() float64 {
	return n.Value
}

func (n *NumberNode) String() string {
	return fmt.Sprintf("%.2f", n.Value)
}

// 二元运算节点
type BinaryOpNode struct {
	Left     ASTNode
	Operator string
	Right    ASTNode
}

func (n *BinaryOpNode) Evaluate() float64 {
	left := n.Left.Evaluate()
	right := n.Right.Evaluate()

	switch n.Operator {
	case "+":
		return left + right
	case "-":
		return left - right
	case "*":
		return left * right
	case "/":
		return left / right
	default:
		return 0
	}
}

func (n *BinaryOpNode) String() string {
	return fmt.Sprintf("(%s %s %s)", n.Left.String(), n.Operator, n.Right.String())
}

// 简化的语法分析器 - 对应ANTLR的Parser
type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func (p *Parser) currentToken() Token {
	if p.pos >= len(p.tokens) {
		return Token{EOF, ""}
	}
	return p.tokens[p.pos]
}

func (p *Parser) advance() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

// 解析表达式 - 处理加减法（低优先级）
func (p *Parser) parseExpression() ASTNode {
	left := p.parseTerm()

	for p.currentToken().Type == PLUS || p.currentToken().Type == MINUS {
		op := p.currentToken().Value
		p.advance()
		right := p.parseTerm()
		left = &BinaryOpNode{Left: left, Operator: op, Right: right}
	}

	return left
}

// 解析项 - 处理乘除法（高优先级）
func (p *Parser) parseTerm() ASTNode {
	left := p.parseFactor()

	for p.currentToken().Type == MULTIPLY || p.currentToken().Type == DIVIDE {
		op := p.currentToken().Value
		p.advance()
		right := p.parseFactor()
		left = &BinaryOpNode{Left: left, Operator: op, Right: right}
	}

	return left
}

// 解析因子 - 处理数字和括号
func (p *Parser) parseFactor() ASTNode {
	token := p.currentToken()

	if token.Type == NUMBER {
		p.advance()
		value, _ := strconv.ParseFloat(token.Value, 64)
		return &NumberNode{Value: value}
	}

	if token.Type == LPAREN {
		p.advance() // 跳过 '('
		expr := p.parseExpression()
		p.advance() // 跳过 ')'
		return expr
	}

	panic(fmt.Sprintf("Unexpected token: %s", token.Value))
}

// 演示ANTLR概念的主函数
func main() {
	// 演示不同的表达式
	expressions := []string{
		"2 + 3",
		"2 + 3 * 4",
		"(2 + 3) * 4",
		"2 + 3 * (4 - 1)+1+3",
		"10 / 2 - 3",
	}

	fmt.Println("=== ANTLR概念演示：计算器 ===\n")

	for _, expr := range expressions {
		fmt.Printf("表达式: %s\n", expr)

		// 1. 词法分析 (Lexer)
		tokens := tokenize(expr)
		fmt.Print("词法分析 (Tokens): ")
		for _, token := range tokens {
			if token.Type != EOF {
				fmt.Printf("%s ", token.Value)
			}
		}
		fmt.Println()

		// 2. 语法分析 (Parser) - 构建AST
		parser := NewParser(tokens)
		ast := parser.parseExpression()
		fmt.Printf("语法树 (AST): %s\n", ast.String())

		// 3. 语义处理 (Listener/Visitor) - 计算结果
		result := ast.Evaluate()
		fmt.Printf("计算结果: %.2f\n", result)

		fmt.Println(strings.Repeat("-", 40))
	}

	// 演示为什么需要语法分析
	fmt.Println("\n=== 为什么需要语法分析？ ===")
	fmt.Println("表达式: 2 + 3 * 4")
	fmt.Println("错误理解: (2 + 3) * 4 = 20")
	fmt.Println("正确理解: 2 + (3 * 4) = 14")
	fmt.Println("ANTLR通过语法规则自动处理运算符优先级！")

	fmt.Println("\n=== 对比正则表达式的局限性 ===")
	fmt.Println("正则表达式: 只能做文本匹配，无法理解语法结构")
	fmt.Println("ANTLR: 理解语法规则，生成结构化的语法树")
	fmt.Println("这就是为什么SQL解析需要ANTLR而不是正则表达式！")
}
