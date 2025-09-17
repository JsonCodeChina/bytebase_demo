# SQL Review Learning Demo

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![GitHub Issues](https://img.shields.io/github/issues/JsonCodeChina/bytebase_demo.svg)](https://github.com/JsonCodeChina/bytebase_demo/issues)

> 🎯 **学习目标**: 基于 Bytebase 架构的企业级 SQL 审查系统学习演示项目

通过构建一个简化但功能完整的 SQL 审查系统，深度理解企业级规则引擎的设计思路和实现机制。

## ✨ 项目亮点

- 🏗️ **插件化架构**: 基于 Bytebase 设计的规则系统
- 🔌 **多数据库支持**: MySQL、PostgreSQL 连接管理
- 🌐 **RESTful API**: 完整的 HTTP 服务接口
- ⚙️ **企业级配置**: 多环境配置管理系统
- 📊 **Schema 分析**: 数据库结构读取和分析
- 🔍 **实时审查**: SQL 语句实时审查功能

## 🚀 快速开始

### 前置要求

- Go 1.24 或更高版本
- MySQL/PostgreSQL 数据库（可选，用于测试）

### 安装运行

```bash
# 克隆项目
git clone https://github.com/JsonCodeChina/bytebase_demo.git
cd bytebase_demo

# 安装依赖
make deps

# 构建项目
make build-all

# 启动 API 服务器
make run-server
```

服务器将在 `http://localhost:8080` 启动

### 测试 API

```bash
# 健康检查
curl http://localhost:8080/health

# 查看可用规则
curl http://localhost:8080/api/rules

# 测试 SQL 审查
curl -X POST http://localhost:8080/api/sql/review \
  -H "Content-Type: application/json" \
  -d '{
    "sql": "CREATE TABLE users (name VARCHAR(50))",
    "connection_id": "demo",
    "rules": ["mysql.table.require-pk"]
  }'
```

## 📁 项目结构

```
sql-review-learning-demo/
├── cmd/
│   └── server/              # HTTP API 服务器
├── pkg/
│   ├── advisor/             # SQL 审查核心引擎
│   ├── api/                 # HTTP API 处理器
│   ├── config/              # 配置管理系统
│   ├── database/            # 数据库连接管理
│   └── rules/               # SQL 审查规则实现
│       └── mysql/           # MySQL 特定规则
├── config/                  # 配置文件
│   ├── app.yaml            # 主配置
│   ├── rules.yaml          # 规则配置
│   ├── development.yaml    # 开发环境配置
│   └── production.yaml     # 生产环境配置
├── examples/               # SQL 测试样例
└── docs/                   # 项目文档
```

## 🎨 核心架构

### 插件化规则系统

```go
// 审查器接口
type Advisor interface {
    Check(ctx context.Context, checkCtx *Context) ([]*Advice, error)
}

// 规则接口
type Rule interface {
    ID() string
    Check(ctx context.Context, checkCtx *Context) ([]*Advice, error)
}
```

### 配置管理

支持多层配置加载：**默认配置** → **文件配置** → **环境变量**

```bash
# 开发环境
APP_ENV=development make run-server

# 生产环境
APP_ENV=production make run-server

# 自定义端口
SERVER_PORT=9000 make run-server
```

## 🔧 API 接口

| 端点 | 方法 | 描述 |
|------|------|------|
| `/health` | GET | 健康检查 |
| `/api/rules` | GET | 列出所有规则 |
| `/api/connections/test` | POST | 测试数据库连接 |
| `/api/connections` | GET/POST | 管理数据库连接 |
| `/api/schema/:id` | GET | 获取数据库 schema |
| `/api/sql/review` | POST | 执行 SQL 审查 |

### 请求示例

**测试数据库连接**:
```json
POST /api/connections/test
{
  "host": "localhost",
  "port": 3306,
  "database": "test",
  "username": "root",
  "password": "password",
  "engine": "mysql"
}
```

**SQL 审查**:
```json
POST /api/sql/review
{
  "sql": "CREATE TABLE users (id INT, name VARCHAR(50))",
  "connection_id": "demo",
  "rules": ["mysql.table.require-pk"]
}
```

## 📋 已实现的规则

- ✅ **表主键检查** (`mysql.table.require-pk`): 确保每个表都有主键
- 🔄 **命名规范检查** (规划中): 表名和列名命名约定
- 🔄 **语句安全检查** (规划中): 危险 SQL 操作检查
- 🔄 **性能优化建议** (规划中): SELECT 语句优化建议

## 🛠️ 开发指南

### 添加新规则

1. 在 `pkg/rules/mysql/` 创建新规则文件
2. 实现 `Rule` 接口
3. 在 `cmd/server/main.go` 中注册规则

示例：
```go
type MyRule struct {
    *advisor.BaseRule
}

func (r *MyRule) Check(ctx context.Context, checkCtx *advisor.Context) ([]*advisor.Advice, error) {
    // 实现规则逻辑
    return advices, nil
}
```

### 运行测试

```bash
# 运行所有测试
make test

# 详细测试输出
make test-verbose

# 测试特定规则
./bin/sql-review-server &
curl -X POST http://localhost:8080/api/sql/review -d @examples/bad_examples.sql
```

### 代码格式化

```bash
# 格式化代码
make fmt

# 代码检查
make vet

# 运行 linter
make lint
```

## 📚 学习资源

- **架构分析**: `docs/bytebase-sql-review-analysis.md` - Bytebase 架构深度分析
- **项目计划**: `docs/project-plan.md` - 详细的实施计划
- **学习笔记**: `docs/learning-notes.md` - 开发过程中的学习心得
- **参考文档**: `REFERENCES.md` - 相关技术文档链接

## 🎓 学习价值

### 技术收获
- 🏗️ **企业级架构设计**: 插件化、可扩展的系统架构
- 🔧 **Go 高级编程**: 接口设计、依赖注入、配置管理
- 🌐 **API 服务开发**: RESTful API 设计和实现
- 🗄️ **数据库交互**: 连接池管理、Schema 分析

### 设计模式
- **Strategy Pattern**: 规则策略模式
- **Plugin Architecture**: 插件化架构
- **Dependency Injection**: 依赖注入
- **Configuration Management**: 配置管理模式

## 🤝 贡献指南

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/awesome-rule`)
3. 提交更改 (`git commit -am 'Add awesome rule'`)
4. 推送到分支 (`git push origin feature/awesome-rule`)
5. 创建 Pull Request

## 📄 许可证

本项目基于 MIT 许可证开源 - 查看 [LICENSE](LICENSE) 文件了解详情

## 🙏 致谢

- **Bytebase 团队**: 提供了优秀的开源 SQL 审查系统作为学习参考
- **Go 社区**: 提供了丰富的开源库和工具

---

**项目创建**: 2025-09-17
**基于**: Bytebase v3.5.2 架构
**学习重点**: SQL 审查引擎设计与实现

🌟 **如果这个项目对你有帮助，请给个 Star！**