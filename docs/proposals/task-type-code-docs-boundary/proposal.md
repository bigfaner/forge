---
created: 2026-05-18
author: faner
status: Draft
---

# Proposal: Task Type Code/Docs Boundary

## Problem

Task `type` field 有两个问题：

1. **类型名不传达 code/docs 属性**：`enhancement` 可能是改代码也可能只是改 SKILL.md，agent 和运行时都无法从类型名判断是否需要 quality-gate
2. **自动生成任务的 ID 不可读**：`T-quick-5` 是什么？必须查表才知道是 doc drift 检测

### Evidence

- `eval-adversarial-scorer` 的 2 个任务只改 `.md` 文件，但 `type: "enhancement"` + `scope: "backend"` → 走了完整 quality-gate
- `tui-ui-design` 的 6 个任务只改 `plugins/forge/skills/` 下 markdown，标注 `type: "implementation"` → 同样空转 quality-gate
- `IsTestableType` 手动维护 map（只有 feature/enhancement/fix），cleanup/refactor 不在列表中 → 逻辑不一致
- `T-test-4.5` 这类 ID 对人类不友好

### Urgency

中——不会导致错误产出，但浪费时间（quality-gate 在无代码变更时全部空转），且 agent 在 `forge profile` 等步骤上浪费注意力。

## Proposed Solution

### 1. 前缀式类型命名

```
coding.* → 涉及代码 → quality-gate 运行
doc.*    → 纯文档   → quality-gate 跳过
其他     → 元类型   → 特殊处理
```

| 新类型 | 旧类型 | 分类 |
|--------|--------|------|
| `coding.feature` | `feature` | code |
| `coding.enhancement` | `enhancement` | code |
| `coding.cleanup` | `cleanup` | code |
| `coding.refactor` | `refactor` | code |
| `coding.fix` | `fix` | code |
| `coding.simplify` | `code-quality.simplify` | code |
| `doc` | `documentation` | doc |
| `doc.eval` | `doc-evaluation` | doc |
| `doc.summary` | `doc-generation.summary` | doc |
| `doc.consolidate` | `doc-generation.consolidate` | doc |
| `doc.drift` | `doc-generation.drift` | doc |
| `test.gen-cases` | `test-pipeline.gen-cases` | test |
| `test.eval-cases` | `test-pipeline.eval-cases` | test |
| `test.gen-scripts` | `test-pipeline.gen-scripts` | test |
| `test.run` | `test-pipeline.run` | test |
| `test.gen-and-run` | `test-pipeline.gen-and-run` | test |
| `test.graduate` | `test-pipeline.graduate` | test |
| `test.verify-regression` | `test-pipeline.verify-regression` | test |
| `gate` | `gate` | meta |

### 2. `IsTestableType` 简化为前缀判断

```go
// Before: 手动维护 map
var testableTypes = map[string]bool{
    TypeFeature: true, TypeEnhancement: true, TypeFix: true,
}

// After: 前缀即规则，零维护
func IsTestableType(typ string) bool {
    return strings.HasPrefix(typ, "coding.")
}
```

`needsDocEval` 同理简化：

```go
func isDocsOnlyType(typ string) bool {
    return strings.HasPrefix(typ, "doc.")
}
```

### 3. `submit.go` 增加类型检查

```go
// Before: 只看 noTest
if rd.Status == "completed" && !submitForce && !t.NoTest {

// After: noTest 或非 coding 类型均跳过
if rd.Status == "completed" && !submitForce && !t.NoTest && task.IsTestableType(t.Type) {
```

### 4. 可读的自动生成任务 ID

| 旧 ID | 新 ID | 类型 |
|--------|-------|------|
| `T-quick-1` | `T-quick-gen-cases` | `test.gen-cases` |
| `T-quick-2` | `T-quick-gen-and-run` | `test.gen-and-run` |
| `T-quick-3` | `T-quick-graduate` | `test.graduate` |
| `T-quick-4` | `T-quick-verify-regression` | `test.verify-regression` |
| `T-quick-5` / `T-quick-specs-1` | `T-quick-doc-drift` | `doc.drift` |
| `T-test-1` | `T-test-gen-cases` | `test.gen-cases` |
| `T-test-1b` | `T-test-eval-cases` | `test.eval-cases` |
| `T-test-2` | `T-test-gen-scripts` | `test.gen-scripts` |
| `T-test-3` | `T-test-run` | `test.run` |
| `T-test-4` | `T-test-graduate` | `test.graduate` |
| `T-test-4.5` | `T-test-verify-regression` | `test.verify-regression` |
| `T-test-5` / `T-specs-1` | `T-specs-consolidate` | `doc.consolidate` |
| `T-eval-doc` | `T-eval-doc` | `doc.eval` |
| `T-clean-code-1` | `T-clean-code` | `coding.simplify` |

Profile 后缀保留：`T-test-gen-scripts-a`（profile a）、`T-test-gen-scripts-api`（type api）。

`InferType` 仍为映射表（ID 后缀不等于类型名，如 `gen-cases` ≠ `test.gen-cases`，`clean-code` ≠ `coding.simplify`），但 key 从数字变为可读字符串。

### 5. 分类规则：按产出物分类

```
看 Affected Files 的文件扩展名：
  - 全部是 .md/.yaml/.json（非构建配置） → type = doc
  - 包含任何可编译/可执行文件          → type = coding.* (按意图选后缀)
```

### 6. Agent 指导更新

- `type-assignment.md`：更新类型表 + 前缀规则 + 产出物分类规则
- `guide.md`：quality-gate 协议明确 `coding.*` 运行 / 其他跳过
- `submit-task` SKILL.md：检查 `IsTestableType` 而非仅 `noTest`
- `quick-tasks` / `breakdown-tasks` SKILL.md：docs-only 检测引导
- Task templates (`task.md` / `task-doc.md`)：更新默认 type 值

## Design Decisions

### D1: 前缀驱动而非 noTest 驱动

**选择**: `coding.*` 前缀作为 quality-gate 运行的唯一条件

**理由**: 前缀判断零维护（`strings.HasPrefix`），不需要手动维护 `testableTypes` map。Agent 看到类型名就知道是否涉及代码。

**否决方案**: 扩展 `noTest` 语义（`notest-docs-only-detection` 提案）。需要 agent 设置两个字段，增加遗漏风险。

### D2: 按产出物分类而非意图分类

**选择**: 只改 `.md` 文件的任务 type = `doc`，无论意图是"实现功能"还是"写文档"

**理由**: 意图模糊（"改进 agent" 是 coding 还是 doc？），产出物客观。Quality-gate 关心的也是产出物——compile/fmt/lint/test 对 .md 没有意义。

### D3: InferType 仍为映射表

**选择**: ID 和 type 之间保持显式映射，不从 ID 推导 type

**理由**: ID 后缀不完全等同于类型名。映射表的 key 变得更可读，但逻辑结构不变。

### D4: `noTest` 保留作为显式覆盖

**选择**: `coding.*` 隐式走 quality-gate，`doc`/`test.*`/`gate` 隐式跳过。`noTest: true` 保留给 edge case（coding task 但不需要测试）。

## Requirements Analysis

### Key Scenarios

1. **纯 markdown feature**: 所有任务 `type: "doc"` → `IsTestableType` 返回 false → 跳过 quality-gate → 生成 `T-eval-doc`
2. **纯 code feature**: 任务 `type: "coding.feature"` → `IsTestableType` 返回 true → 走 quality-gate → 生成测试管线
3. **Mixed feature**: 部分 `doc`，部分 `coding.enhancement` → 各自按 type 决定 → 整体有 coding 任务 → 生成测试管线
4. **向后兼容**: 旧 index.json（`type: "feature"`）→ `IsTestableType` 返回 false → 需重新生成 index.json

### Constraints & Dependencies

- 修改 `plugins/forge/` 下的文件前必须先加载 `docs/conventions/forge-distribution.md`
- 类型常量重命名影响 `types.go`、`build.go`、`infer.go`、`prompt.go`、`submit.go`、`quality_gate.go`
- Task ID 重命名影响 `infer.go`、`testgen.go`
- Prompt template 文件名需与新类型名对齐

## Scope

### In Scope

**Go 代码变更**：
- `forge-cli/pkg/task/types.go` — 重命名 19 个类型常量
- `forge-cli/pkg/task/build.go` — `testableTypes` map → `IsTestableType` 前缀判断
- `forge-cli/pkg/task/infer.go` — 映射 key 从数字更新为可读字符串
- `forge-cli/pkg/task/testgen.go` — 更新生成的 task ID
- `forge-cli/pkg/prompt/prompt.go` — 更新 `typeToTemplate` keys
- `forge-cli/pkg/prompt/data/*.md` — 文件重命名对齐新类型名
- `forge-cli/internal/cmd/submit.go` — 增加 `IsTestableType` 检查
- `forge-cli/internal/cmd/quality_gate.go` — `isDocsOnly` 用前缀判断

**Agent 指导变更**：
- `plugins/forge/references/shared/type-assignment.md` — 更新类型表 + 分类规则
- `plugins/forge/hooks/guide.md` — quality-gate 协议更新
- `plugins/forge/skills/submit-task/SKILL.md` — type 检查逻辑
- `plugins/forge/skills/quick-tasks/SKILL.md` — docs-only 检测引导
- `plugins/forge/skills/breakdown-tasks/SKILL.md` — docs-only 检测引导
- `plugins/forge/skills/quick-tasks/templates/task*.md` — 更新默认 type 值
- `plugins/forge/skills/breakdown-tasks/templates/task*.md` — 更新默认 type 值

### Out of Scope

- `noTest` 字段废弃
- `docs-only-fast-path` 提案的范围
- 已完成的 feature 目录批量迁移（按需手动更新）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 旧 index.json 类型名不兼容 | High | M — validate-index 报错 | `forge task index` 重新生成即可 |
| Agent 仍按意图分类而非产出物 | M | M — type 错误导致空转 quality-gate | type-assignment.md 明确规则 + 示例 |
| ID 重命名破坏 profile/type 后缀解析 | M | H — test pipeline 无法正确推断类型 | 保留 `profileSuffixedID` / `typeSuffixedID` 逻辑，更新 base strings |
| Prompt template 文件重命名遗漏 | L | M — `forge prompt get-by-task-id` 报错 unknown type | `typeToTemplate` map 编译时校验 |

## Success Criteria

- [ ] Go 类型常量全部使用 `coding.`/`doc`/`test.`/`gate` 命名
- [ ] `IsTestableType` 使用 `strings.HasPrefix(typ, "coding.")`，不再依赖 map
- [ ] `submit.go` 检查 `IsTestableType` 决定是否走 quality-gate
- [ ] 自动生成任务 ID 可读（`T-quick-doc-drift` 而非 `T-quick-5`）
- [ ] `infer.go` 映射表 key 更新为可读字符串
- [ ] 现有测试全部通过
- [ ] 纯 markdown feature 的任务自动标 `type: "doc"`
