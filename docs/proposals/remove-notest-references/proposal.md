---
created: 2026-05-23
author: "faner"
status: Draft
---

# Proposal: 移除 noTest 残留引用

## Problem

`noTest` 字段已从 Go 运行时（`Task`、`TaskState`、`FrontmatterData` 结构体）和 CLI（`--no-test` 标志）中完全移除，系统已切换到基于 `task.IsTestableType(t.Type)` 的类型前缀判断。但代码库中仍残留约 75 处 `"noTest": true`（index.json）、9 处 frontmatter `noTest: true`、以及多处技能文档和模板中的引用，形成活代码与死数据的割裂。

### Evidence

- `forge-cli/pkg/task/build.go:433` — `IsTestableType()` 已替代 `NoTest` 字段
- `forge-cli/pkg/task/autogen_test.go:162` — 测试验证 frontmatter 不再包含 `noTest`
- `forge-cli/tests/task-lifecycle/task_stage_gates_test.go:739` — 测试验证 `--no-test` 标志已移除
- ~75 个 `index.json` 文件仍含 `"noTest": true`（Go JSON 解析器静默忽略未知字段）

### Urgency

低紧迫性但高价值。残留数据不影响运行时，但会误导新开发者认为 `noTest` 仍是有效配置项。清理成本低、风险极低。

## Proposed Solution

批量删除活代码和配置中所有 `noTest` 引用，保持历史文档（proposals、forensics、self-evolution 报告）不动。

### Innovation Highlights

纯清理任务，无创新点。标准的废弃字段清理流程。

## Requirements Analysis

### Key Scenarios

- 开发者阅读任务文件时不被废弃字段误导
- eval-forge 审计不再检查已移除的 noTest 绕过路径
- index.json 数据模型与 Go 结构体保持一致

### Non-Functional Requirements

- 零功能回归：所有删除项必须是运行时已不读取的死数据
- 可验证：清理后 grep `noTest` 在活代码中零命中

### Constraints & Dependencies

- Go 运行时已完成迁移（`IsTestableType` 替代 `NoTest`），无依赖约束

## Alternatives & Industry Benchmarking

### Industry Solutions

N/A — 标准代码卫生实践。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 误导开发者，代码库不整洁 | Rejected: 活代码与数据不一致 |
| **Batch cleanup** | Standard practice | 彻底清理，低风险 | 需修改约 50 个文件 | **Selected: 成本低、收益明确** |

## Feasibility Assessment

### Technical Feasibility

完全可行。Go 运行时已不读取 `noTest`，所有删除都是安全的。

### Resource & Timeline

单次提交，预计 1 小时内完成。

### Dependency Readiness

无外部依赖。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| noTest 残留仅是 index.json | 5 Whys | Confirmed: 还存在于 frontmatter、skill 文档、eval 模板、测试代码中 |
| 清理不会影响运行时 | Stress Test | Confirmed: Go JSON 解析器忽略未知字段，YAML 解析器忽略未知键 |

## Scope

### In Scope

- 从 ~40 个 `index.json` 中删除 `"noTest": true` 字段
- 从 ~9 个任务 `.md` frontmatter 中删除 `noTest: true` 行
- 更新 `plugins/forge/hooks/guide.md` 移除 noTest 引用
- 更新 `plugins/forge/skills/consolidate-specs/SKILL.md` 移除 noTest 引用
- 更新 `.claude/skills/eval-forge/templates/` 下模板移除 noTest 审计项
- 更新 `tests/test-generation/quick_test_slim_test.go` 移除本地 NoTest 结构体字段
- 更新 `docs/lessons/gotcha-docs-only-needs-code-audit.md` 中的 noTest 引用

### Out of Scope

- 历史文档（proposals/、forensics/、self-evolution/）中的 noTest 引用保留不动
- Go 运行时代码（已完成迁移，无需改动）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 漏删某些引用 | L | L | 清理后 grep 验证零命中 |
| 编辑 index.json 破坏 JSON 格式 | L | L | 批量脚本确保 JSON 合法性 |

## Success Criteria

- [ ] `grep -r "noTest" --include="*.json" --include="*.md" --include="*.go" plugins/ forge-cli/ tests/` 在活代码和配置中零命中（排除 docs/proposals/、docs/forensics/、docs/self-evolution/）
- [ ] 所有 `index.json` 文件 JSON 格式合法
- [ ] 现有测试全部通过

## Next Steps

- Proceed to `/quick-tasks` to generate cleanup tasks
