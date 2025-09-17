# SQL Review Learning Demo - 开发待办事项

## 📋 当前进度

### ✅ 已完成任务
- [x] **创建 Demo 程序的基础项目结构**
  - 项目目录结构完整创建
  - 文档和配置文件齐全
  - 测试样例 SQL 文件准备完毕

### 🔄 进行中任务

### ⏳ 待开始任务

1. **实现核心 Advisor 接口和规则引擎**
   - [ ] 定义 `pkg/advisor/advisor.go` 核心接口
   - [ ] 实现 `pkg/advisor/registry.go` 规则注册器
   - [ ] 创建 `pkg/advisor/context.go` 审查上下文
   - [ ] 参考: `bytebase/backend/plugin/advisor/advisor.go`

2. **集成 ANTLR MySQL 解析器**
   - [ ] 添加 ANTLR MySQL 语法文件
   - [ ] 生成 Go 语言解析器代码
   - [ ] 封装解析器接口 `pkg/parser/mysql/parser.go`
   - [ ] 实现 AST 遍历和节点分析
   - [ ] 参考: Bytebase 的解析器集成方式

3. **实现 5-10 个典型的 SQL 审查规则**
   - [ ] **表结构规则**: `pkg/rules/mysql/table_require_pk.go` - 表必须有主键
   - [ ] **命名规范**: `pkg/rules/mysql/naming_convention.go` - 表名/列名规范
   - [ ] **数据类型**: `pkg/rules/mysql/column_type_check.go` - 禁止特定类型
   - [ ] **语句安全**: `pkg/rules/mysql/statement_safety.go` - 危险操作检查
   - [ ] **性能优化**: `pkg/rules/mysql/select_performance.go` - SELECT 优化建议
   - [ ] 参考: `bytebase/backend/plugin/advisor/mysql/` 目录下的规则实现

4. **开发简单的 CLI 工具展示功能**
   - [ ] 实现 `cmd/demo/main.go` 主程序入口
   - [ ] 添加 Cobra CLI 框架
   - [ ] 实现 `check` 命令用于检查 SQL 文件
   - [ ] 添加配置文件支持
   - [ ] 实现友好的结果展示格式

5. **编写测试用例和使用文档**
   - [ ] 为每个规则编写单元测试
   - [ ] 添加集成测试用例
   - [ ] 完善 README.md 使用说明
   - [ ] 更新学习笔记和心得体会

## 🎯 详细实施计划

### 阶段一：核心引擎开发 (2-3天)
**目标**: 建立基础架构，实现规则引擎框架

#### Day 1-2: 核心接口设计
```go
// pkg/advisor/advisor.go
type Advisor interface {
    Check(ctx context.Context, checkCtx Context) ([]*Advice, error)
}

type Context struct {
    SQL        string
    Engine     Engine
    Rules      []*Rule
    // 其他上下文信息
}
```

#### Day 2-3: 解析器集成
- 集成 ANTLR4 Go runtime
- 实现基础的 MySQL SQL 解析
- 提供统一的 AST 接口

### 阶段二：规则实现 (3-4天)
**目标**: 实现核心 SQL 审查规则

#### 优先实现的规则（按难度递增）:
1. **表主键检查** - 最基础，容易理解
2. **命名约定检查** - 正则表达式匹配
3. **SELECT * 检查** - AST 节点分析
4. **WHERE 子句检查** - 复杂语句分析
5. **数据类型检查** - 表结构分析

### 阶段三：工具完善 (2-3天)
**目标**: 用户友好的 CLI 工具和完整文档

#### CLI 功能:
- `demo check <file.sql>` - 检查单个文件
- `demo check <directory>` - 批量检查
- `demo config` - 配置管理
- `demo rules` - 列出可用规则

## 🔗 核心参考资源

### Bytebase 关键文件
- **核心接口**: `/Users/shenbo/goprojects/bytebase-3.5.2/backend/plugin/advisor/advisor.go`
- **MySQL 规则**: `/Users/shenbo/goprojects/bytebase-3.5.2/backend/plugin/advisor/mysql/`
- **API 集成**: `/Users/shenbo/goprojects/bytebase-3.5.2/backend/api/v1/sql_service.go:1606`

### 学习资源
- 详细分析报告: `docs/bytebase-sql-review-analysis.md`
- 项目计划: `docs/project-plan.md`
- 学习笔记: `docs/learning-notes.md`

### 测试资源
- 正确示例: `examples/good_examples.sql`
- 错误示例: `examples/bad_examples.sql`
- 混合示例: `examples/mixed_examples.sql`

## 🚀 开发命令快捷方式

```bash
# 开发环境设置
make dev-setup

# 构建和测试
make build
make test
make run

# 示例运行
make example-good    # 测试正确的 SQL
make example-bad     # 测试错误的 SQL
make example-mixed   # 测试混合的 SQL

# 参考 Bytebase
make ref-bytebase    # 显示 Bytebase 参考路径
make ref-docs        # 打开文档
```

## 📊 成功标准

### 阶段一完成标准
- [ ] 基础 Advisor 接口能够运行
- [ ] 简单的 MySQL 语句能够解析
- [ ] 规则注册机制工作正常

### 阶段二完成标准
- [ ] 至少 5 个规则正确实现
- [ ] 所有测试样例能正确识别问题
- [ ] 错误信息清晰友好

### 阶段三完成标准
- [ ] CLI 工具功能完备
- [ ] 使用文档清晰完整
- [ ] 代码质量达到生产标准

## 🎓 学习目标检查清单

- [ ] 深度理解插件化架构设计
- [ ] 掌握 ANTLR 在实际项目中的应用
- [ ] 学会企业级规则引擎实现
- [ ] 提升 Go 语言高级编程技巧
- [ ] 获得可重用的架构模板

---

**创建时间**: 2025-09-16
**预计完成**: 7-10 天
**学习重点**: 从 Bytebase 学习企业级 SQL 审查系统的设计和实现

**祝你学习愉快！我们在 demo 项目中见！** 🎉