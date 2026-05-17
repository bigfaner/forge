---
id: "2"
title: "Migrate CLI template and fix wiring"
priority: "P0"
estimated_time: "30m"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "feature"
mainSession: false
---

# 2: Migrate CLI template and fix wiring

## Description

删除旧的 `code-quality-simplify.md` prompt template，新建 `code-quality-clean-code.md` 调用 `forge:clean-code` skill。修复 `typeToTemplate` 映射，使 `TypeCleanCode` 指向新模板。Bump version。

## Reference Files

- `docs/proposals/clean-code-skill/proposal.md` — Source proposal
- `forge-cli/pkg/prompt/data/code-quality-simplify.md` — Old template to delete
- `forge-cli/pkg/prompt/prompt.go` — typeToTemplate mapping to fix
- `forge-cli/pkg/task/types.go` — TypeCleanCode constant reference
- `scripts/version.txt` — Version to bump

## Acceptance Criteria

- [ ] `forge-cli/pkg/prompt/data/code-quality-simplify.md` deleted
- [ ] `forge-cli/pkg/prompt/data/code-quality-clean-code.md` created, calls `Skill(skill="forge:clean-code")`
- [ ] `forge-cli/pkg/prompt/prompt.go` `typeToTemplate` has entry: `task.TypeCleanCode: "data/code-quality-clean-code.md"`
- [ ] `forge prompt get-by-task-id T-clean-code-1` returns valid prompt (no "unknown type" error)
- [ ] `scripts/version.txt` bumped (patch)
- [ ] `go test -race -cover ./...` passes from `forge-cli/` directory

## Hard Rules

- Prompt template 文件通过 `//go:embed` 嵌入，新文件必须放在 `pkg/prompt/data/` 目录下
- `typeToTemplate` 的 key 必须使用 `task.TypeCleanCode` 常量，不要用字符串字面量

## Implementation Notes

- `TypeCleanCode = "code-quality.simplify"` 定义在 `pkg/task/types.go` — 类型名不变，只是模板文件更换
- Go embed 在 `pkg/prompt/prompt.go` 顶部声明，确保新文件名在 embed pattern 内
- 参考现有 prompt template 格式（如 `data/test-pipeline-gen-cases.md`）了解模板结构
