---
title: "Code Structure Conventions"
domains: [nesting, indentation, early-return, flat-control-flow, readability, package-organization, structural-rules]
---

# Code Structure Conventions

本文档定义 forge-cli 的代码结构规范，涵盖控制流风格和包组织结构规则。这是**规范性（normative）**文档，描述目标状态而非当前状态。所有新增代码必须符合此规范；已有代码的偏差记录在偏差分析表中，按计划逐步收敛。

## 1. 目标状态

| 属性 | 目标 |
|------|------|
| 控制流嵌套深度 | 最大 2 层，guard clause 优先 |
| 函数职责 | 每个函数单一职责，不超过 50 行 |
| 包间别名 | 零 test-bridge 别名；生产代码直接导入目标包 |
| 重复定义 | 相同签名函数仅存在一处定义 |
| 依赖方向 | 严格 `cmd → internal → pkg`，pkg 内部 `domain → infrastructure → leaf` |

## 2. 控制流规范

### TECH-code-structure-001: Prefer Flat Control Flow Over Deep Nesting

**Requirement**: Avoid deep nesting (3+ levels). Use early returns, guard clauses, and extracted helper functions to keep the main logic at minimal indentation. Each tab level adds cognitive load — flatten aggressively.
**Scope**: [CROSS]
**Source**: /learn entry 2026-05-18

**Pattern to avoid** (4 levels):
```
func process() {
    if cond1 {
        if cond2 {
            if cond3 {
                if cond4 {
                    // actual logic buried here
                }
            }
        }
    }
}
```

**Preferred pattern** (1 level, guard clauses):
```
func process() {
    if !cond1 { return ... }
    if !cond2 { return ... }
    if !cond3 { return ... }
    if !cond4 { return ... }
    // actual logic at minimal indentation
}
```

**For switch/case**: prefer flat switch bodies. Avoid nested if/else inside case arms — extract to helpers when logic exceeds 5-10 lines.

## 3. 包组织结构规则

以下规则引用 [package-organization.md](./package-organization.md) 中的依赖方向定义，补充结构层面的具体约束。

### TECH-code-structure-002: No Duplicate Function Definitions Across Packages

**Requirement**: 禁止在不同包中定义签名相同的导出函数。如果多个包需要同一能力，应在最底层的合适包中定义一次，消费方统一导入。
**Scope**: `internal/cmd/`, `internal/cmd/base/`, `pkg/`

**违规示例**：`internal/cmd/output.go` 和 `internal/cmd/base/output.go` 均定义了 `func Debugf(verbose bool, format string, args ...any)` -- 签名完全相同，行为完全相同。

**正确模式**：
```
// 仅在 base/output.go 定义
package base
func Debugf(verbose bool, format string, args ...any) { ... }

// 消费方统一导入
package cmd
import "forge-cli/internal/cmd/base"
// 使用 base.Debugf(...)
```

### TECH-code-structure-003: No Package-Level Variable Aliases for Production Code

**Requirement**: 禁止使用包级变量别名（`var funcName = pkg.SomeFunc`）来桥接 `internal/cmd` 与 `pkg/` 之间的调用。生产代码必须直接导入并调用目标包的函数。

**例外**：如果别名的唯一目的是测试可替换性（且无生产调用），可以暂时保留但必须标注 `// Kept as alias for test mockability` 并在偏差表中记录。

**Scope**: `internal/cmd/` 下所有子包

**违规示例**：
```go
// internal/cmd/task/claim.go
var getTaskPhase = task.GetTaskPhase  // 避免此模式
```

**正确模式**：
```go
// internal/cmd/task/validate_index.go
import "forge-cli/pkg/task"
// 直接调用
phase := task.GetTaskPhase(g.id)
```

### TECH-code-structure-004: Enforce Dependency Direction in cmd Subpackages

**Requirement**: `internal/cmd/<subpackage>/` 中的代码只能导入 `pkg/` 中的包和 `internal/cmd/base/`。禁止导入 `internal/cmd/` 顶层包或其他 `internal/cmd/<other-subpackage>/`。

**Scope**: `internal/cmd/` 下所有子包

**依赖方向**：
```
internal/cmd/<subpkg>/ → internal/cmd/base/
internal/cmd/<subpkg>/ → pkg/*
internal/cmd/<subpkg>/ ✗→ internal/cmd/ (顶层)
internal/cmd/<subpkg>/ ✗→ internal/cmd/<other-subpkg>/
```

### TECH-code-structure-005: Deprecated Fields Must Document Migration Path

**Requirement**: 标记为 `Deprecated` 的字段必须在注释中说明替代方案和迁移检测函数。Deprecated 字段不应有新的写入方；只允许存在迁移相关的读取逻辑。

**Scope**: `pkg/` 下所有结构体定义

**正确模式**：
```go
// Deprecated: Scope is the legacy scope field replaced by SurfaceKey/SurfaceType.
// Retained solely for migration detection via CheckLegacyScope. Do not use in new code.
Scope string `yaml:"scope"`
```

## 4. 模块级偏差摘要

| 编号 | 偏差项 | 违反规则 | 当前状态 | 目标状态 |
|:----:|--------|----------|----------|----------|
| ~~CS-1~~ | `cmd.Debugf` 重复定义 | TECH-code-structure-002 | **已修复** — `cmd/output.go` 中的重复定义已删除，仅保留 `base.Debugf` | N/A |
| CS-2 | `getTaskPhase` 别名有生产调用 | TECH-code-structure-003 | 5 处生产调用通过别名而非直接调用 | 迁移为 `task.GetTaskPhase` 直接调用 |
| CS-3 | `checkExistingTaskState` 别名 | TECH-code-structure-003 | 1 处生产调用通过别名 | 迁移为 `task.CheckExistingTaskState` 直接调用 |
| CS-4 | `compareVersionIDs` 别名 | TECH-code-structure-003 | 1 处生产调用通过别名 | 迁移为 `task.CompareVersionIDs` 直接调用 |
| ~~CS-5~~ | `Scope` deprecated 字段 | TECH-code-structure-005 | **已修复** — `FrontmatterData.Scope` 已从 `frontmatter.go` 移除，`CheckLegacyScope` 保留用于迁移检测 | N/A |

## 5. 参考

- [package-organization.md](./package-organization.md) -- 包组织规范：三层模型、依赖方向、扇入分析
- [dead-code.md](./dead-code.md) -- 死代码识别标准和清理流程
- [pkg-dependency-graph.md](../../features/forge-cli-codebase-standards/pkg-dependency-graph.md) -- pkg/ 依赖图事实基线
