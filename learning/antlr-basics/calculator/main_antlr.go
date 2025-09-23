package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

// CalculatorEvaluator 实现计算器的计算逻辑
// 使用 Listener 模式遍历 ANTLR 生成的语法树
type CalculatorEvaluator struct {
	*BaseCalculatorListener
	stack []float64
}

// NewCalculatorEvaluator 创建新的计算器求值器
func NewCalculatorEvaluator() *CalculatorEvaluator {
	return &CalculatorEvaluator{
		BaseCalculatorListener: &BaseCalculatorListener{},
		stack:                  make([]float64, 0),
	}
}

// push 将值压入栈
func (c *CalculatorEvaluator) push(value float64) {
	c.stack = append(c.stack, value)
}

// pop 从栈中弹出值
func (c *CalculatorEvaluator) pop() float64 {
	if len(c.stack) == 0 {
		return 0
	}
	value := c.stack[len(c.stack)-1]
	c.stack = c.stack[:len(c.stack)-1]
	return value
}

// GetResult 获取最终计算结果
func (c *CalculatorEvaluator) GetResult() float64 {
	if len(c.stack) > 0 {
		return c.stack[len(c.stack)-1]
	}
	return 0
}

// ExitNumber 当退出数字节点时，将数字值压入栈
func (c *CalculatorEvaluator) ExitNumber(ctx *NumberContext) {
	value, _ := strconv.ParseFloat(ctx.GetText(), 64)
	c.push(value)
}

// ExitMulDiv 当退出乘除节点时，执行乘除运算
func (c *CalculatorEvaluator) ExitMulDiv(ctx *MulDivContext) {
	right := c.pop()
	left := c.pop()

	op := ctx.GetChild(1).GetPayload().(*antlr.CommonToken).GetText()

	var result float64
	switch op {
	case "*":
		result = left * right
	case "/":
		if right != 0 {
			result = left / right
		} else {
			fmt.Println("警告：除零操作")
			result = 0
		}
	}

	c.push(result)
}

// ExitAddSub 当退出加减节点时，执行加减运算
func (c *CalculatorEvaluator) ExitAddSub(ctx *AddSubContext) {
	right := c.pop()
	left := c.pop()

	op := ctx.GetChild(1).GetPayload().(*antlr.CommonToken).GetText()

	var result float64
	switch op {
	case "+":
		result = left + right
	case "-":
		result = left - right
	}

	c.push(result)
}

// ExitParens 括号表达式不需要特殊处理，值已经在栈中

// ErrorListener 自定义错误监听器
type ErrorListener struct {
	*antlr.DefaultErrorListener
}

// SyntaxError 处理语法错误
func (e *ErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, ex antlr.RecognitionException) {
	fmt.Printf("语法错误 [行%d:列%d]: %s\n", line, column, msg)
}

// calculate 使用 ANTLR 解析和计算表达式
func calculate(expression string) (float64, error) {
	// 1. 创建输入流
	input := antlr.NewInputStream(expression)

	// 2. 创建词法分析器
	lexer := NewCalculatorLexer(input)

	// 3. 创建 Token 流
	stream := antlr.NewCommonTokenStream(lexer, 0)

	// 4. 创建语法分析器
	parser := NewCalculatorParser(stream)

	// 5. 设置错误监听器
	parser.RemoveErrorListeners()
	parser.AddErrorListener(&ErrorListener{})

	// 6. 从起始规则开始解析
	tree := parser.Prog()

	// 7. 创建计算器求值器
	evaluator := NewCalculatorEvaluator()

	// 8. 使用 Listener 遍历语法树
	antlr.ParseTreeWalkerDefault.Walk(evaluator, tree)

	return evaluator.GetResult(), nil
}

// printTree 打印语法树结构（用于调试）
func printTree(tree antlr.ParseTree, indent string) {
	if tree == nil {
		return
	}

	// 打印当前节点
	switch ctx := tree.(type) {
	case *ProgContext:
		fmt.Printf("%sProg\n", indent)
	case *NumberContext:
		fmt.Printf("%sNumber: %s\n", indent, ctx.GetText())
	case *MulDivContext:
		op := ctx.GetChild(1).(*antlr.TerminalNodeImpl).GetText()
		fmt.Printf("%sMulDiv: %s\n", indent, op)
	case *AddSubContext:
		op := ctx.GetChild(1).(*antlr.TerminalNodeImpl).GetText()
		fmt.Printf("%sAddSub: %s\n", indent, op)
	case *ParensContext:
		fmt.Printf("%sParens\n", indent)
	default:
		fmt.Printf("%s%s\n", indent, tree.GetText())
	}

	// 递归打印子节点
	for i := 0; i < tree.GetChildCount(); i++ {
		if child := tree.GetChild(i); child != nil {
			if parseTreeChild, ok := child.(antlr.ParseTree); ok {
				printTree(parseTreeChild, indent+"  ")
			}
		}
	}
}

func main() {
	fmt.Println("=== ANTLR 真实工作流程演示：计算器 ===\n")

	// 测试表达式
	expressions := []string{
		"2 + 3",
		//"2 + 3 * 4",
		//"(2 + 3) * 4",
		//"10 / 2 - 3",
		//"2.5 * (3 + 4.5)",
	}

	for _, expr := range expressions {
		fmt.Printf("表达式: %s\n", expr)

		// 使用 ANTLR 解析并计算
		result, err := calculate(expr)
		if err != nil {
			fmt.Printf("错误: %v\n", err)
		} else {
			fmt.Printf("结果: %.2f\n", result)
		}

		// 显示语法树结构
		fmt.Println("语法树结构:")
		input := antlr.NewInputStream(expr)
		lexer := NewCalculatorLexer(input)
		stream := antlr.NewCommonTokenStream(lexer, 0)
		parser := NewCalculatorParser(stream)
		tree := parser.Prog()
		printTree(tree, "  ")

		fmt.Println(strings.Repeat("-", 50))
	}

	fmt.Println("\n=== ANTLR 工作流程说明 ===")
	fmt.Println("1. Grammar File (.g4)    → 定义语法规则")
	fmt.Println("2. ANTLR Tool            → 生成 Lexer, Parser, Listener")
	fmt.Println("3. Input Stream          → 将文本转换为字符流")
	fmt.Println("4. Lexer                 → 词法分析，产生 Token 流")
	fmt.Println("5. Parser                → 语法分析，构建语法树 (AST)")
	fmt.Println("6. Listener/Visitor      → 遍历语法树，执行语义处理")
	fmt.Println("7. Result                → 获得最终结果")

	fmt.Println("\n=== 与手工实现的对比 ===")
	fmt.Println("✅ ANTLR 自动处理运算符优先级")
	fmt.Println("✅ ANTLR 自动生成错误处理")
	fmt.Println("✅ ANTLR 支持复杂的语法结构")
	fmt.Println("✅ ANTLR 生成的代码经过严格测试")
	fmt.Println("✅ ANTLR 易于扩展和维护")
}
