# Go Test Staging/Graduation Friction: Package Isolation vs Nested Staging

## Problem

go-test profile 的测试文件放在 `tests/e2e/features/<slug>/` 子目录后，无法引用父目录 `tests/e2e/helpers.go` 中的函数。Go 的 package 按目录隔离编译，子目录是独立的编译单元。

## Root Cause

forge 的 staging/graduation 两阶段流程是为 TypeScript 设计的：

- **TypeScript**: `import { x } from '../../helpers.js'` -- 相对路径，搬文件后改路径即可
- **Go**: `package e2e` -- 同一 package 的所有文件必须在同一目录，子目录看不到父目录的符号

`gen-test-scripts` SKILL.md 的 HARD-GATE 硬编码输出到 `tests/e2e/features/<feature>/`，但 Go 在该子目录编译时找不到 `helpers.go`。

## Workaround

当前每个 feature 子目录必须：
1. 复制一份 helpers（`tsptRunCLI`、`tsptRunCLIRaw` 等）
2. 加 `init()` 用 `runtime.Caller` 找到项目根目录

这导致 `runCLI`、`withRetry`、`parseBlock` 等函数在项目中存在多个副本。

## Solution

go-test profile 应使用 `flat` staging 模式：直接生成到 `tests/e2e/` 根目录，跳过 `features/` 子目录。

**关键决策**：flat 模式下不生成 graduation 任务（T-test-4 / T-quick-4）。测试文件已在最终位置，run-e2e-tests 通过后自动写 marker。Graduation 任务的存在前提（文件搬移 + import rewrite）在 flat 模式下不成立，生成空转任务只会增加 pipeline 开销。

详见: `docs/proposals/go-flat-staging/proposal.md`

## Key Takeaway

| 语言特性 | Staging 模式 | Graduation 任务 |
|----------|-------------|----------------|
| 相对路径 import (TS) | `nested` -- 需要 features/ staging 做 import rewrite | 生成 -- 需要搬移 + rewrite |
| Module path / 同 package 共享 (Go, Python, Rust) | `flat` -- 直接放目标位置，无需 rewrite | **不生成** -- 无事可做 |

Profile 的 staging 策略应该匹配语言的 import/package 模型，而非一刀切。当 staging 已经是最终位置时，graduation 任务不应存在。
