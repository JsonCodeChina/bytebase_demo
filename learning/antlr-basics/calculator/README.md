# Calculator Demo with ANTLR

这个目录展示了 ANTLR 的完整工作流程，从语法文件到实际运行的计算器。

## 文件说明

### 核心文件
- `Calculator.g4` - ANTLR 语法文件（手工编写）
- `main_antlr.go` - 使用 ANTLR 生成代码的实现（**推荐**）
- `main_recursive.go` - 手工递归实现版本（仅用于对比理解）

### ANTLR 生成的文件
- `calculator_parser.go` - 语法分析器
- `calculator_lexer.go` - 词法分析器
- `calculator_listener.go` - 监听器接口
- `calculator_base_listener.go` - 基础监听器实现

### 工具和配置
- `tools/antlr-4.13.1-complete.jar` - ANTLR 工具
- `tools/antlr4` - 便捷脚本
- `go.mod` - Go 模块依赖

## 标准 ANTLR 工作流程

### 1. 编写语法文件 (Calculator.g4)
```antlr
grammar Calculator;

prog: expr EOF ;

expr
    : expr ('*'|'/') expr   # MulDiv      // 乘除，高优先级
    | expr ('+'|'-') expr   # AddSub      // 加减，低优先级
    | '(' expr ')'          # Parens      // 括号
    | NUMBER                # Number      // 数字
    ;

NUMBER : [0-9]+('.'[0-9]+)? ;
WS : [ \t\r\n]+ -> skip ;
```

### 2. 生成 Go 代码
```bash
./tools/antlr4 -Dlanguage=Go Calculator.g4
```
这会生成：
- `calculator_parser.go` - 解析器
- `calculator_lexer.go` - 词法分析器
- `calculator_listener.go` - 监听器接口
- `calculator_base_listener.go` - 基础监听器

### 3. 实现业务逻辑
创建自定义监听器来处理语法树：
```go
type CalculatorEvaluator struct {
    *BaseCalculatorListener
    stack []float64
}

func (c *CalculatorEvaluator) ExitNumber(ctx *NumberContext) {
    // 处理数字节点
}

func (c *CalculatorEvaluator) ExitMulDiv(ctx *MulDivContext) {
    // 处理乘除节点
}
```

### 4. 运行解析
```go
// 1. 创建输入流
input := antlr.NewInputStream(expression)

// 2. 词法分析
lexer := NewCalculatorLexer(input)
stream := antlr.NewCommonTokenStream(lexer, 0)

// 3. 语法分析
parser := NewCalculatorParser(stream)
tree := parser.Prog()

// 4. 语义处理
evaluator := NewCalculatorEvaluator()
antlr.ParseTreeWalkerDefault.Walk(evaluator, tree)
```

## 运行示例

```bash
# 运行 ANTLR 版本（推荐）
go run main_antlr.go calculator_parser.go calculator_lexer.go calculator_listener.go calculator_base_listener.go

# 或者运行递归版本（仅对比用）
go run main_recursive.go
```

## 示例输出

```
表达式: 2 + 3 * 4
结果: 14.00
语法树结构:
  Prog
    AddSub: +
      Number: 2
        2
      +
      MulDiv: *
        Number: 3
          3
        *
        Number: 4
          4
    <EOF>
```

## ANTLR 优势

✅ **自动优先级处理**: 通过语法规则顺序定义优先级
✅ **错误处理**: 自动生成语法错误检测
✅ **类型安全**: 生成强类型的语法树节点
✅ **可扩展性**: 通过修改 .g4 文件轻松扩展语法
✅ **工具支持**: 丰富的调试和可视化工具

## 重新生成代码

如果修改了 `Calculator.g4` 文件，需要重新生成代码：

```bash
# 重新生成
./tools/antlr4 -Dlanguage=Go Calculator.g4

# 修正包名（生成的代码默认是 parser 包）
sed -i '' 's/package parser/package main/g' calculator_*.go

# 运行
go run main_antlr.go calculator_parser.go calculator_lexer.go calculator_listener.go calculator_base_listener.go
```