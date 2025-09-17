# 项目参考资源

## 🔗 Bytebase 源码参考

### 核心文件位置
- **项目路径**: `/Users/shenbo/goprojects/bytebase-3.5.2/`
- **许可证**: 查看 `LICENSE` 文件了解使用限制

### 关键代码位置

#### SQL 审查核心引擎
- `backend/plugin/advisor/advisor.go` - Advisor 接口定义
- `backend/plugin/advisor/builtin_rules.go` - 内置规则定义
- `backend/plugin/advisor/change_type.go` - 变更类型定义

#### MySQL 规则实现
- `backend/plugin/advisor/mysql/` - MySQL 专用规则 (97个文件)
  - `advisor_table_require_pk.go` - 表主键检查
  - `advisor_naming_*.go` - 命名规范检查
  - `advisor_column_*.go` - 列约束检查
  - `advisor_statement_*.go` - 语句规范检查

#### API 服务层
- `backend/api/v1/review_config_service.go` - 审查配置服务
- `backend/api/v1/sql_service.go` - SQL 服务 (审查集成点在 1606 行)

#### 数据存储层
- `backend/store/review_config.go` - 审查配置存储
- `backend/store/policy.go` - 策略相关存储

#### 前端组件
- `frontend/src/components/IssueV1/components/SQLCheckSection/SQLCheckButton.vue`
- `frontend/src/components/SQLCheck/` - SQL 检查相关组件

## 📚 技术文档参考

### ANTLR4
- [官方文档](https://github.com/antlr/antlr4/doc/index.md)
- [Go Target](https://github.com/antlr/antlr4/blob/master/doc/go-target.md)
- [MySQL Grammar](https://github.com/antlr/grammars-v4/tree/master/sql/mysql)

### Go 语言
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go 语言规范](https://golang.org/ref/spec)
- [设计模式 in Go](https://github.com/tmrts/go-patterns)

### SQL 标准
- [MySQL 8.0 参考手册](https://dev.mysql.com/doc/refman/8.0/en/)
- [SQL 标准文档](https://www.iso.org/standard/63555.html)

## 🎯 学习重点

### 从 Bytebase 学习的核心概念
1. **插件化架构设计**
2. **规则引擎的可扩展性**
3. **ANTLR 在企业项目中的实际应用**
4. **多数据库支持的抽象设计**
5. **企业级功能的集成方式**

### 关键设计模式
- **Strategy Pattern**: 规则实现
- **Visitor Pattern**: AST 遍历
- **Factory Pattern**: 规则创建
- **Observer Pattern**: 结果通知

## ⚖️ 使用说明

### Bytebase 代码参考原则
- **仅用于学习目的**：理解架构设计和实现思路
- **不直接复制代码**：避免许可证问题
- **学习设计模式**：重点理解设计思路而非具体实现
- **简化实现**：去除企业级复杂性，专注核心概念

### 参考方式
1. **接口设计参考**：学习如何定义清晰的接口
2. **错误处理参考**：学习企业级的错误处理模式
3. **代码组织参考**：学习大型项目的目录结构
4. **命名规范参考**：学习一致的命名风格

## 🔄 更新记录

- **2025-09-16**: 初始创建，添加 Bytebase v3.5.2 参考信息
- 后续会根据学习进度更新更多参考资源

---

**注意**: 本项目是学习性质的演示项目，所有参考均遵循相应的开源协议和使用条款。