---
created: 2026-05-18
author: faner
status: Approved
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

**2. Go 代码：`testableTypes` 扩展**

`forge-cli/pkg/task/build.go` 的 `testableTypes` 当前只有 `{feature, enhancement, fix}`，缺少 `cleanup` 和 `refactor`。导致 cleanup/refactor 任务被当作 docs-only，不走 quality-gate 也不生成 test pipeline。

修正：将 `cleanup` 和 `refactor` 加入 `testableTypes`，使 `IsTestableType` 和 `needsTestPipeline` 正确识别它们。

**3. Go 代码：submit-task quality-gate 改用 type 驱动**

`forge-cli/internal/cmd/submit.go` 当前仅用 `t.NoTest` 决定是否跳过 quality-gate。修正为：`!t.NoTest && IsTestableType(t.Type)` 才触发 quality-gate。同时，coverage auto-set（`coverage = -1.0`）也应对非 testable type 生效。

**4. quality-gate 执行点（guide.md）增加 type 检查**

guide.md quality-gate 协议统一规则：

// After: noTest 或非 coding 类型均跳过
if rd.Status == "completed" && !submitForce && !t.NoTest && task.IsTestableType(t.Type) {
```

**5. task 生成器（quick-tasks、breakdown-tasks）强化分类**

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

- 修改 `plugins/forge/` 下的文件前必须先加载 `docs/conventions/forge-distribution.md`
- `type-assignment.md` 是 shared reference，被 quick-tasks 和 breakdown-tasks 引用
- `guide.md` 是全局 hook，所有 task-executing 工作流读取
- Go 代码变更涉及 `forge-cli/`，需遵循 TDD（RED → GREEN → REFACTOR），版本号需 bump
- `testableTypes` 扩展后，`isDocsOnly()`（quality_gate.go）和 `needsTestPipeline()`（build.go）自动受益——它们已使用 `IsTestableType`

### D4: `noTest` 保留作为显式覆盖

**选择**: `coding.*` 隐式走 quality-gate，`doc`/`test.*`/`gate` 隐式跳过。`noTest: true` 保留给 edge case（coding task 但不需要测试）。

## Requirements Analysis

### Key Scenarios

8 个文件变更：5 个 markdown（skill/reference 文档）+ 2 个 Go 文件 + 1 个测试文件。Go 变更范围小（修改 map + 条件判断），测试用 table-driven 覆盖。

### Constraints & Dependencies

| Step | Task | 预计时间 |
|------|------|---------|
| 1 | 更新 `type-assignment.md` 增加判定规则 | 15min |
| 2 | 更新 `guide.md` quality-gate 协议增加 type 检查 | 15min |
| 3 | 更新 `submit-task` 检查 type=documentation | 15min |
| 4 | 更新 quick-tasks / breakdown-tasks 的 docs-only 检测引导 | 15min |
| 5 | Go: 扩展 `testableTypes`，修改 submit.go quality-gate skip 逻辑 | 30min |
| 6 | Go: 单元测试覆盖 | 30min |
| **Total** | | **~2h** |

## Scope

### In Scope

- `plugins/forge/references/shared/type-assignment.md` — 增加 code/docs 判定规则
- `plugins/forge/hooks/guide.md` — quality-gate 协议增加 `type: "documentation"` 跳过规则
- `plugins/forge/skills/submit-task/SKILL.md` — type 检查跳过 quality-gate
- `plugins/forge/skills/quick-tasks/SKILL.md` — 强化 docs-only 分类引导
- `plugins/forge/skills/breakdown-tasks/SKILL.md` — 强化 docs-only 分类引导
- `forge-cli/pkg/task/build.go` — `testableTypes` 增加 cleanup 和 refactor
- `forge-cli/internal/cmd/submit.go` — quality-gate skip 改用 `IsTestableType` 驱动
- Go 单元测试覆盖上述变更

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

- [ ] `type-assignment.md` 明确 "按产出物分类" 规则，包含 code types / doc type 分类表
- [ ] `guide.md` quality-gate 协议明确 `type: "documentation"` 跳过 quality-gate
- [ ] `submit-task` skill 检查 type 并跳过 quality-gate
- [ ] quick-tasks / breakdown-tasks 的 docs-only 分类引导更新
- [ ] 纯 markdown feature（如 eval-adversarial-scorer 风格）的任务自动标 `type: "documentation"`
- [ ] `testableTypes` 包含 cleanup 和 refactor，`IsTestableType` 正确返回 true
- [ ] submit.go quality-gate skip 同时检查 `noTest` 和 `IsTestableType`
- [ ] Go 单元测试覆盖 `testableTypes` 扩展和 submit type-based skip
