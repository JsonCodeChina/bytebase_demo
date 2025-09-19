# 简化SQL解析示例

## 目标

通过一个简化的SQL解析示例，演示为什么我们需要从正则表达式升级到ANTLR。

## 核心问题

当前项目使用正则表达式检查SQL：
```go
if strings.Contains(upperDef, "PRIMARY KEY") {
    return true
}
```

这种方法的问题：
1. **无法区分上下文** - 注释、字符串值、列名都被当作SQL关键字
2. **容易误判** - `primary_config` 这样的列名会被误认为主键
3. **无法扩展** - 添加新规则很困难

## 测试用例

### 会被正则表达式误判的情况

1. **注释中的关键字**:
```sql
CREATE TABLE logs (
    id INT,
    message TEXT -- Should add PRIMARY KEY later
);
-- 正则表达式会误认为有主键，但实际没有
```

2. **列名包含关键字**:
```sql
CREATE TABLE settings (
    id INT,
    primary_config VARCHAR(100)  -- 列名包含'primary'
);
-- 正则表达式会误判为有主键
```

3. **字符串值包含关键字**:
```sql
CREATE TABLE messages (
    id INT,
    content TEXT DEFAULT 'Add PRIMARY KEY to table'
);
-- 默认值包含关键字，但不是真正的主键定义
```

## ANTLR的解决方案

ANTLR通过语法规则精确识别：
```antlr
columnConstraint
    : PRIMARY KEY              # 只有真正的PRIMARY KEY约束
    | NOT NULL
    | AUTO_INCREMENT
    ;
```

### 关键优势

1. **语法感知** - 理解SQL的语法结构
2. **上下文区分** - 自动处理注释、字符串、标识符
3. **结构化输出** - 生成AST，便于复杂分析
4. **易于扩展** - 添加新语法规则很简单

## 运行演示

```bash
cd learning/antlr-basics/sql-simple
go run main.go
```

这个演示会运行多个测试用例，对比正则表达式和ANTLR方法的准确性。

## 与Bytebase的对比

Bytebase使用相同的思路：
- **输入**: AST而不是原始SQL字符串
- **处理**: Listener模式遍历语法树
- **输出**: 结构化的检查结果

这就是我们需要学习和采用的方法！