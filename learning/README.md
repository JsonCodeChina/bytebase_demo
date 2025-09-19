# ANTLR学习实验室

这个目录包含了理解和学习ANTLR语法分析的完整示例，解释为什么我们需要从正则表达式升级到ANTLR来解析SQL。

## 📚 学习路径

### 1. ANTLR基础概念 (`antlr-basics/`)

#### 🧮 计算器示例 (`antlr-basics/calculator/`)
- **目标**: 理解ANTLR的核心概念
- **内容**:
  - `Calculator.g4` - ANTLR语法文件
  - `main.go` - 完整的词法/语法分析演示
  - `README.md` - 详细的概念解释

**运行方式**:
```bash
cd learning/antlr-basics/calculator
go run main.go
```

**学到什么**:
- 词法分析 (Lexer): 文本 → Token
- 语法分析 (Parser): Token → AST
- 语义处理 (Listener): AST → 结果
- 运算符优先级的自动处理

#### 🗄️ 简化SQL示例 (`antlr-basics/sql-simple/`)
- **目标**: 将ANTLR概念应用到SQL解析
- **内容**:
  - `SimpleSQL.g4` - 简化的SQL语法
  - `main.go` - SQL解析演示
  - `README.md` - SQL解析特有的问题

**运行方式**:
```bash
cd learning/antlr-basics/sql-simple
go run main.go
```

**学到什么**:
- SQL语法比算术表达式复杂得多
- 如何精确识别CREATE TABLE结构
- 为什么需要语法感知而不是文本匹配

### 2. 深度对比分析 (`parsing-comparison/`)

#### ⚔️ 正则表达式 vs ANTLR (`parsing-comparison/`)
- **目标**: 直观对比两种方法的优劣
- **内容**:
  - `regex_vs_antlr.go` - 详细的准确性对比
  - `README.md` - 误判案例分析

**运行方式**:
```bash
cd learning/parsing-comparison
go run regex_vs_antlr.go
```

**关键发现**:
- 正则表达式准确率: **60%**
- ANTLR方法准确率: **100%**

## 🎯 核心问题演示

### 当前项目的问题
我们的 `pkg/rules/mysql/table_require_pk.go` 使用简单字符串匹配：
```go
if strings.Contains(upperDef, "PRIMARY KEY") {
    return true  // 这里会误判！
}
```

### 误判案例

1. **注释误判** ❌:
```sql
CREATE TABLE logs (
    id INT,
    message TEXT -- Should add PRIMARY KEY later
);
-- 正则表达式认为有主键，但实际没有！
```

2. **字符串值误判** ❌:
```sql
CREATE TABLE messages (
    id INT,
    content TEXT DEFAULT 'Add PRIMARY KEY to table'
);
-- 默认值包含关键字，被误认为有主键！
```

3. **列名混淆** (幸好没误判):
```sql
CREATE TABLE settings (
    id INT,
    primary_config VARCHAR(100)  -- 列名包含'primary'
);
-- 这种情况正则表达式碰巧没有误判
```

## 🚀 ANTLR的解决方案

### 核心优势

1. **语法感知**: 理解SQL的语法结构，不只是文本匹配
2. **上下文区分**: 自动区分关键字、注释、字符串、标识符
3. **结构化输出**: 生成AST，便于复杂分析
4. **易于扩展**: 添加新语法规则简单直观

### 工作原理

```
SQL文本 → 词法分析 → Token流 → 语法分析 → AST → Listener处理 → 结果
```

## 📊 对比总结

| 方面 | 正则表达式 | ANTLR |
|------|------------|--------|
| **准确性** | 60% | 100% |
| **可扩展性** | 困难 | 简单 |
| **维护性** | 复杂 | 清晰 |
| **性能** | 快 | 稍慢但可接受 |
| **学习成本** | 低 | 中等 |

## 🛣️ 下一步计划

基于这些学习，我们将：

1. **集成真正的MySQL解析器**
   - 添加 `github.com/bytebase/mysql-parser` 依赖
   - 参考Bytebase的实现方式

2. **重构advisor架构**
   - 修改 `Context` 支持AST输入
   - 将解析阶段和规则检查阶段分离

3. **重写规则实现**
   - 用Listener模式重写 `table_require_pk` 规则
   - 大幅提升准确性

4. **添加更多规则**
   - 命名规范检查
   - 语句安全检查
   - 性能优化建议

## 🎓 学习价值

通过这个实验室，我们深入理解了：

- 为什么企业级数据库工具都使用ANTLR
- 正则表达式在复杂语法解析中的局限性
- 如何从"文本匹配"升级到"语法理解"
- 真正的解析器是如何工作的

这些知识将帮助我们构建更加准确和强大的SQL审查系统！

## 🔗 参考资源

- [ANTLR官方文档](https://www.antlr.org/)
- [Bytebase项目](https://github.com/bytebase/bytebase) - 企业级应用示例
- [mysql-parser](https://github.com/bytebase/mysql-parser) - 我们将要集成的解析器