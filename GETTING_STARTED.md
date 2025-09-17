# 快速开始指南

欢迎来到 SQL Review Learning Demo 项目！这是一个基于 Bytebase SQL 审查模块的学习实践项目。

## 🎯 项目目标

通过构建一个简化但功能完整的 SQL 审查系统，深度理解企业级规则引擎的设计思路和实现机制。

## 📚 背景知识

在开始编码之前，强烈建议先阅读以下文档：

1. **深度分析报告**: `docs/bytebase-sql-review-analysis.md`
   - 理解 Bytebase SQL 审查模块的完整架构
   - 掌握核心设计模式和技术亮点

2. **项目计划**: `docs/project-plan.md`
   - 了解整体实施计划和学习路径
   - 明确各阶段的目标和成功标准

3. **参考资源**: `REFERENCES.md`
   - 找到 Bytebase 源码中的关键文件
   - 了解相关技术文档和学习资源

## 🚀 开发环境设置

### 1. 检查前置条件

```bash
# 检查 Go 版本 (需要 1.24+)
go version

# 检查项目结构
ls -la
```

### 2. 安装依赖

```bash
# 安装项目依赖
make dev-setup

# 或者直接运行
go mod download
go mod tidy
```

### 3. 验证环境

```bash
# 尝试构建项目 (目前会失败，因为还没有实现代码)
make build

# 查看可用的 make 命令
make help
```

## 📋 开发顺序建议

### 第一步：实现核心接口
从最基础的接口开始：

1. **创建 `pkg/advisor/advisor.go`**
   ```go
   package advisor

   import "context"

   // 核心审查接口
   type Advisor interface {
       Check(ctx context.Context, checkCtx Context) ([]*Advice, error)
   }
   ```

2. **参考 Bytebase 实现**
   - 查看 `/Users/shenbo/goprojects/bytebase-3.5.2/backend/plugin/advisor/advisor.go`
   - 理解接口设计的思路
   - 简化但保持核心概念

### 第二步：集成解析器
1. **添加 ANTLR 依赖**
2. **实现基础 SQL 解析**
3. **封装解析器接口**

### 第三步：实现规则
从简单规则开始，逐步增加复杂度：

1. **表主键检查** - 最基础的规则
2. **命名约定检查** - 正则表达式应用
3. **SELECT 语句检查** - AST 遍历应用

## 🧪 测试驱动开发

使用提供的 SQL 样例文件进行测试：

```bash
# 测试正确的 SQL (应该通过所有检查)
cat examples/good_examples.sql

# 测试错误的 SQL (应该触发各种规则)
cat examples/bad_examples.sql

# 测试混合的 SQL (既有好的也有坏的)
cat examples/mixed_examples.sql
```

## 📖 学习方法建议

### 1. 对比学习法
- 先看 Bytebase 的实现
- 理解设计思路
- 用自己的方式简化实现

### 2. 迭代开发法
- 从最简单的功能开始
- 每次只增加一个小功能
- 确保每一步都能正常工作

### 3. 文档驱动法
- 及时更新 `docs/learning-notes.md`
- 记录遇到的问题和解决方案
- 分享设计思路和权衡考虑

## 🔧 开发工具使用

### Make 命令
```bash
make help          # 显示所有可用命令
make build         # 构建项目
make test          # 运行测试
make run           # 运行程序
make clean         # 清理构建文件
```

### Git 工作流
```bash
# 初始化 git (如果需要)
git init
git add .
git commit -m "Initial project setup"

# 为每个功能创建分支
git checkout -b feature/core-advisor-interface
git checkout -b feature/mysql-parser-integration
```

## 🎓 学习检查点

在每个阶段完成后，问自己这些问题：

### 阶段一检查点
- [ ] 我理解 Advisor 接口的设计思路吗？
- [ ] 我知道如何注册和执行规则吗？
- [ ] 我能解释插件化架构的优势吗？

### 阶段二检查点
- [ ] 我理解 ANTLR 如何解析 SQL 吗？
- [ ] 我能遍历 AST 并分析特定节点吗？
- [ ] 我知道如何实现不同类型的规则吗？

### 阶段三检查点
- [ ] 我的 CLI 工具用户体验如何？
- [ ] 错误信息是否清晰友好？
- [ ] 代码质量是否达到生产标准？

## 🆘 遇到问题时

### 1. 查看参考资源
- 回到 Bytebase 源码中找答案
- 查看相关的技术文档
- 阅读已有的分析报告

### 2. 简化问题
- 将复杂问题拆分为简单步骤
- 先实现最小可用版本
- 逐步增加功能

### 3. 记录过程
- 在 `docs/learning-notes.md` 中记录问题
- 分享解决方案和思考过程
- 为后续学习者提供参考

## 🎉 完成标志

当你能够成功运行以下命令时，说明项目基本完成：

```bash
# 检查好的 SQL (应该没有错误)
./bin/sql-review-demo check examples/good_examples.sql

# 检查坏的 SQL (应该发现多个问题)
./bin/sql-review-demo check examples/bad_examples.sql

# 显示规则列表
./bin/sql-review-demo rules list
```

---

**准备好了吗？让我们开始这个激动人心的学习之旅吧！** 🚀

**记住**: 这不仅仅是一个编程练习，更是深度理解企业级系统设计的绝佳机会。

**我们在代码中见！** 👋