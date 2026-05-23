---
id: "1"
title: "创建 pkg/infocmd/ 共享工具包"
priority: "P0"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: true
type: "coding.feature"
mainSession: false
---

# 1: 创建 pkg/infocmd/ 共享工具包

## Description

创建 `forge-cli/pkg/infocmd/` 包，用 Go 泛型封装三个 info-command（research/proposal/lesson）的重复模式。提供泛型 Discover[T]、FindByID[T]、排序工具、共享 parseFrontmatter 函数。

这是整个重构的基础设施任务。下游三个迁移任务（task 2-4）将依赖此包。

## Reference Files
- `docs/proposals/infocmd-shared-pkg/proposal.md` — Source proposal
- `forge-cli/pkg/research/research.go` — 现有 Discover/FindBySlug/parseFrontmatter 实现
- `forge-cli/pkg/proposal/proposal.go` — 现有 Discover/FindBySlug/parseFrontmatter 实现
- `forge-cli/pkg/lesson/lesson.go` — 现有 Discover/FindByName/parseFrontmatter 实现
- `forge-cli/internal/cmd/base/output.go` — 现有输出工具函数

## Acceptance Criteria
- [ ] `parseFrontmatter` 仅存在于 `pkg/infocmd/` 一处
- [ ] `Discover[T]` 泛型函数支持两种目录模式：flat files（research/lesson）和 subdirectory+fixed-file（proposal）
- [ ] `FindByID[T]` 泛型函数支持配置化的标识符提取（Slug 或 Name）
- [ ] 排序逻辑统一：按 Created 降序，mtime 回退
- [ ] 包通过 `go vet` 和编译
- [ ] 现有三个命令的行为不受影响（此 task 仅创建新包，不修改现有代码）

## Hard Rules
- 不修改现有 `pkg/research/`、`pkg/proposal/`、`pkg/lesson/` 中的任何代码
- 泛型约束使用 `any`，不引入不必要的接口约束
- 新包不得导入 `pkg/research`、`pkg/proposal`、`pkg/lesson`（避免循环依赖）

## Implementation Notes

### 目录扫描模式差异
- research/lesson：扫描 `docs/<dir>/*.md`，slug = filename minus `.md`
- proposal：扫描 `docs/proposals/*/proposal.md`，slug = subdirectory name

建议通过配置结构体封装差异：
```go
type ScanConfig[T any] struct {
    BaseDir    string                          // e.g. "docs/research"
    IsSubdir   bool                            // true for proposal-style scanning
    FileName   string                          // only used when IsSubdir=true, e.g. "proposal.md"
    IDKey      func(T) string                  // extract identifier (Slug or Name)
    ParseEntry func(name, path string, content []byte, modTime time.Time) (T, error)
}
```

### 排序逻辑
research 和 lesson 的 Discover 内部排序逻辑一致（三段式比较：both have Created → lex desc; one has → it wins; neither → mtime desc）。统一到 Discover 内部。

### parseFrontmatter
三个实现逐字节相同。直接提取到新包，变为 exported 函数 `ParseFrontmatter`。
