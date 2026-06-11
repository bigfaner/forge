---
created: 2026-05-15
author: faner
status: Draft
---

# Proposal: noTest 字段驱动的 docs-only 检测

## Problem

`isDocsOnlyFeature()` 仅依赖 `type` 字段判断是否生成测试管线，但 `type` 同时承担两个不相关的职责：

1. **描述性**：任务在做什么（implementation / documentation / fix）
2. **控制流**：是否需要测试管线（implementation/fix → 生成，documentation → 不生成）

对于 forge-internal feature（修改 SKILL.md、模板、rubric），这两个职责冲突：修改 SKILL.md 是在**实现** pipeline 功能（描述性 = implementation），但 e2e 测试管线**无法测试** markdown 变更（控制流应 = documentation）。

Agent 被迫用一个标签满足两个需求，导致误判。

### Evidence

`tui-ui-design` feature 的 6 个业务任务全部修改 markdown 文件（`plugins/forge/skills/` 下），全部标注 `type: "implementation"`。`isDocsOnlyFeature()` 返回 false，生成了 5 个无意义的 T-quick 测试任务。测试管线（go-test profile）无法测试任何 markdown 变更。

## Proposed Solution

复用已有的 `Task.NoTest` 字段，扩展 `isDocsOnlyFeature()` 的检测逻辑。`noTest` 已在 task frontmatter 中支持、已由 `BuildIndex` 解析、已跳过 quality gate —— 仅需让它同时影响 feature-level 的测试管线生成决策。

### 核心改动

`forge-cli/pkg/task/build.go` 中 `isDocsOnlyFeature()` 增加 `noTest` 检查：

```go
func isDocsOnlyFeature(tasks map[string]Task) bool {
    for _, t := range tasks {
        if isAutoGenTaskID(t.ID) {
            continue
        }
        if !t.NoTest && (t.Type == TypeImplementation || t.Type == TypeFix) {
            return false
        }
    }
    return true
}
```

**判定逻辑**：业务任务**同时满足** `NoTest == false` 且 `type` 为 implementation/fix 时，才算非 docs-only。当所有业务任务都有 `noTest: true`（或 type 非 implementation/fix）→ feature 为 docs-only → 生成 `T-eval-doc`。

### 配套改动

quick-tasks SKILL.md 的 Type Assignment 增加 `noTest` 引导：

> Set `noTest: true` when the task only modifies non-code files (markdown, YAML, JSON under `plugins/forge/skills/` or `docs/`). These tasks implement forge pipeline functionality but produce no artifacts the e2e test infrastructure can verify.

### 设计决策

#### D1: 复用 noTest 而非新增字段

**选择**: 扩展现有 `NoTest` 字段的语义

**理由**: `NoTest` 已存在且语义高度相关（"这个任务不需要测试" ↔ "docs-only feature 不需要测试管线"）。避免新增字段带来的认知负担和迁移成本。

**否决方案**: 新增 `testable` 字段。功能等价但增加 schema 复杂度和 agent 决策负担（需要同时设置 type + testable 两个字段）。

#### D2: ALL 语义而非 ANY 语义

**选择**: 当**所有**业务任务都有 `noTest: true` 时才判定为 docs-only

**理由**: 测试管线是 feature-level 的，只要有任何一个任务产生可测试产物，管线就应该生成。Mixed feature（部分 code + 部分 docs）不应被错误地判定为 docs-only。

**否决方案**: ANY 语义（任何一个任务有 `noTest: true` 就判定为 docs-only）。会错误跳过 mixed feature 的测试管线。

#### D3: 不改变 type 字段的语义

**选择**: `type: "implementation"` 继续表示"实现了功能"

**理由**: type 的描述性语义是正确的——修改 SKILL.md 确实是在实现功能。只是 type 不应独自承担测试管线生成的决策。通过 `noTest` 解耦后，agent 可以自由使用语义正确的 type 值。

## Requirements Analysis

### Key Scenarios

1. **纯 markdown feature**: 所有任务 `type: "implementation"` + `noTest: true` → docs-only → 生成 T-eval-doc
2. **纯 code feature**: 所有任务 `type: "implementation"` + `noTest: false`（默认）→ 非 docs-only → 生成 T-quick-1~5（行为不变）
3. **Mixed feature**: 部分任务 `noTest: true`，部分 `noTest: false` → 非 docs-only → 生成 T-quick-1~5（行为正确）
4. **传统 documentation feature**: 所有任务 `type: "documentation"` → docs-only → 生成 T-eval-doc（行为不变）
5. **向后兼容**: 旧 index.json 无 `noTest` 字段（Go 零值 `false`）→ 退化为原有 type-based 逻辑（行为不变）

### Constraints & Dependencies

- `NoTest` 字段已存在于 `Task` struct（types.go:88-89），已由 frontmatter 解析（build.go:115），无需 schema 变更
- 改动限于 `isDocsOnlyFeature()` 函数（build.go:400-410），3 行改动
- quick-tasks SKILL.md 的 Type Assignment 段落追加引导文字

### Non-Functional Requirements

| NFR Category | Requirement | Verification Method |
|-------------|-------------|-------------------|
| 兼容性 | 无 `noTest` 字段的旧 index.json 行为完全不变 | 单元测试：NoTest=false + type=implementation → 非 docs-only |
| 正确性 | 纯 markdown feature 不生成测试管线 | 单元测试：所有任务 NoTest=true → docs-only |
| 正确性 | Mixed feature 正确生成测试管线 | 单元测试：混合 NoTest → 非 docs-only |

## Alternatives & Industry Benchmarking

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| **SKILL.md 引导（仅改文档）** | 零代码改动 | 依赖 agent 正确分类 type，type 语义被曲解 | Rejected |
| **新增 `testable` 字段** | 关注点分离清晰 | 新增字段 + agent 需设两个字段 | Rejected |
| **复用 `noTest`（本提案）** | 零 schema 变更，复用已有语义 | noTest 语义从 per-task 扩展到 feature-level | **Selected** |
| **自动检测文件扩展名** | 全自动，零 agent 判断 | 解析自由格式 markdown 脆弱 | Rejected |

## Feasibility Assessment

### Technical Feasibility

改动范围极小：`build.go` 一个函数 3 行改动 + SKILL.md 一段文字。`NoTest` 字段已有完整管线支持（frontmatter 解析、index.json 序列化、quality gate 跳过）。

### Resource & Timeline

| Step | Task | 预计时间 |
|------|------|---------|
| 1 | 修改 `isDocsOnlyFeature()` + 单元测试 | 30min |
| 2 | 更新 quick-tasks SKILL.md Type Assignment | 15min |
| 3 | 验证 tui-ui-design feature 的 `forge task index` 输出 | 15min |
| **Total** | | **~1h** |

## Scope

### In Scope

- `forge-cli/pkg/task/build.go` — `isDocsOnlyFeature()` 增加 `noTest` 检查
- `forge-cli/pkg/task/build_test.go` — 新增 docs-only 检测单元测试（5 个场景）
- `plugins/forge/skills/quick-tasks/SKILL.md` — Type Assignment 增加 `noTest` 引导

### Out of Scope

- `isDocsOnlyFeature()` 以外的逻辑（stage-gate、quality gate 等不受影响）
- 新增 frontmatter 字段
- 自动检测文件扩展名（可作为未来增强）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Agent 未设置 noTest: true，仍然误判 | Medium | Low — 回退到当前行为（type-based 检测） | SKILL.md 引导 + 未来可加 lint warning |
| noTest 语义扩展造成概念混淆 | Low | Low | 文档明确说明 noTest 的两层含义 |

## Success Criteria

- [ ] 纯 markdown feature（所有任务 `noTest: true`）生成 `T-eval-doc`，不生成 T-quick-1~5
- [ ] 纯 code feature（默认 `noTest: false`）行为完全不变
- [ ] Mixed feature（部分 `noTest: true`，部分 `noTest: false`）正确生成测试管线
- [ ] 旧 index.json（无 `noTest` 字段）行为完全不变
- [ ] 单元测试覆盖 5 个关键场景

## Next Steps

- Proceed to `/quick-tasks` to generate tasks
