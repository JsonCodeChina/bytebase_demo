# 解析方法对比：正则表达式 vs ANTLR

## 对比目标

通过具体的测试用例，展示正则表达式和ANTLR在SQL解析方面的差异，说明为什么我们需要升级解析方法。

## 核心问题演示

### 当前项目的问题
我们的 `table_require_pk.go` 使用这样的逻辑：
```go
if strings.Contains(upperDef, "PRIMARY KEY") {
    return true
}
```

这种方法无法区分：
- SQL关键字 vs 注释内容
- SQL关键字 vs 字符串值
- SQL关键字 vs 列名标识符

### 误判案例

1. **注释误判**:
```sql
CREATE TABLE logs (
    id INT,
    message TEXT -- Should add PRIMARY KEY later
);
-- 正则表达式会误认为有主键！
```

2. **字符串值误判**:
```sql
CREATE TABLE messages (
    id INT,
    content TEXT DEFAULT 'Add PRIMARY KEY to table'
);
-- 默认值包含关键字，但不是真的主键定义
```

## ANTLR的解决方案

### 词法分析阶段
ANTLR首先将SQL分解为有类型的Token：
```
'PRIMARY KEY' → PRIMARY_SYMBOL + KEY_SYMBOL (在约束上下文中)
'primary_config' → IDENTIFIER (列名)
'-- comment' → COMMENT (被跳过)
'string value' → STRING_LITERAL (字符串值)
```

### 语法分析阶段
然后根据SQL语法规则构建AST：
```antlr
columnConstraint
    : PRIMARY KEY              # 只有这里的PRIMARY KEY才是约束
    | NOT NULL
    | AUTO_INCREMENT
    ;
```

### 语义处理阶段
Listener只响应真正的语法节点：
```go
func (checker *tableRequirePKChecker) EnterColumnConstraint(ctx *ColumnConstraintContext) {
    // 只有真正的约束定义才会触发这个方法
    // 注释、字符串、列名都不会到达这里
}
```

## 运行对比

```bash
cd learning/parsing-comparison
go run regex_vs_antlr.go
```

## 学习价值

这个对比帮助我们理解：
1. 为什么简单的字符串匹配不足以处理复杂的语法
2. ANTLR如何通过语法感知解决误判问题
3. 真正的企业级SQL工具为什么都使用ANTLR

## 下一步

基于这个理解，我们将：
1. 集成Bytebase的MySQL解析器
2. 重构advisor架构支持AST
3. 用Listener模式重写规则检查
4. 大幅提升规则的准确性和可扩展性