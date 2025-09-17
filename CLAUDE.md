# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a SQL Review Learning Demo project that implements a simplified but functional SQL review system based on Bytebase's SQL review module. The project serves as a learning exercise to understand enterprise-level rule engine design and implementation.

## Key Commands

### Development Commands
```bash
# Setup development environment
make dev-setup

# Build the project
make build

# Run the demo application
make run

# Run tests
make test
make test-verbose

# Code quality
make fmt
make vet
make lint  # requires golangci-lint installation

# Clean build artifacts
make clean
```

### Example Commands
```bash
# Test with provided SQL examples
make example-good    # Run with good SQL examples
make example-bad     # Run with bad SQL examples
make example-mixed   # Run with mixed SQL examples

# Direct execution (after build)
./bin/sql-review-demo check examples/good_examples.sql
./bin/sql-review-demo check examples/bad_examples.sql
```

### Reference Commands
```bash
make ref-bytebase    # Show Bytebase reference paths
make ref-docs        # Display documentation paths
```

## Architecture

### Core Components Structure
```
pkg/
├── advisor/         # Core review engine interfaces and registry
├── parser/          # SQL parser wrappers (ANTLR integration)
├── rules/           # Specific rule implementations
└── config/          # Configuration management

cmd/
└── demo/           # CLI application entry point

examples/           # Test SQL files for validation
├── good_examples.sql    # SQL that should pass all rules
├── bad_examples.sql     # SQL that should trigger violations
└── mixed_examples.sql   # Mixed good/bad SQL
```

### Key Design Patterns
The project follows a plugin-based architecture inspired by Bytebase:

1. **Advisor Interface Pattern**: Core `Advisor` interface for rule checking
2. **Registry Pattern**: Dynamic rule registration and execution
3. **Context Pattern**: Encapsulates check context and configuration
4. **Visitor Pattern**: AST traversal for SQL analysis

### Dependencies
- **ANTLR4 Go Runtime**: SQL parsing (`github.com/antlr4-go/antlr/v4`)
- **Cobra**: CLI framework (`github.com/spf13/cobra`)
- **YAML**: Configuration management (`gopkg.in/yaml.v3`)

## Development Guidelines

### Reference Implementation
The project references Bytebase v3.5.2 located at `/Users/shenbo/goprojects/bytebase-3.5.2/`. Key reference files:
- `backend/plugin/advisor/advisor.go` - Core interfaces
- `backend/plugin/advisor/mysql/` - MySQL-specific rules
- `backend/api/v1/sql_service.go` - API integration patterns

### Implementation Phases
1. **Core Engine**: Implement basic Advisor interfaces and registry
2. **Parser Integration**: Add ANTLR MySQL parser support
3. **Rule Development**: Implement specific SQL review rules
4. **CLI Tool**: Build user-friendly command-line interface

### Rule Implementation Priority
Start with these rules in order of complexity:
1. Table primary key requirements (basic AST analysis)
2. Naming conventions (regex-based)
3. Column type restrictions (schema analysis)
4. Statement safety checks (complex AST traversal)
5. Performance optimization suggestions (advanced analysis)

### Testing Strategy
Use the provided SQL examples for validation:
- `examples/good_examples.sql`: Should pass all rules
- `examples/bad_examples.sql`: Should trigger various rule violations
- `examples/mixed_examples.sql`: Contains both valid and invalid SQL

### Code Organization
- Place interfaces in `pkg/advisor/`
- Implement MySQL-specific rules in `pkg/rules/mysql/`
- Keep parser logic in `pkg/parser/mysql/`
- Use table-driven tests for rule validation

## Project Status

Currently in planning/setup phase. The directory structure exists but implementation is pending. The project follows a learning-driven development approach with detailed documentation in:
- `docs/bytebase-sql-review-analysis.md` - Deep analysis of Bytebase architecture
- `docs/project-plan.md` - Implementation roadmap
- `TODO.md` - Detailed task breakdown
- `GETTING_STARTED.md` - Development setup guide

## Important Notes

- This is a learning project, not production software
- Focus on understanding architectural patterns over feature completeness
- Maintain compatibility with Go 1.24+
- Follow Go best practices and naming conventions
- Document learning insights in `docs/learning-notes.md`
- 生成关键代码的时候要先写测试用例
- 没有必要生成预期失败的测试用例