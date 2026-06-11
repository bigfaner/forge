---
created: 2026-05-19
author: faner
status: Draft
---

# Proposal: Task Type/ID Redesign

## Problem

三个相关问题叠加：

1. **Type 名不传达 code/docs 属性**：`enhancement` 可能改代码也可能只改 markdown，agent 和运行时无法从类型名判断是否需要 quality-gate。`IsTestableType` 手动维护 map，新增类型必须同步更新。
2. **自动生成任务 ID 不可读**：`T-test-2` 是什么？必须查表。`T-quick-4.5` 更糟。
3. **Skill 模板目录 23 个死文件**：CLI 已程序化生成所有测试/gate/fix/validate 任务的 `.md`，skill 模板无人读取。
4. **validate-code/validate-ux 未接入 CLI 管线**：eval 能力已完整实现，但无 type 常量、无 CLI 生成逻辑。

### Evidence

- `types.go` 的 `testableTypes` 手动维护 `{feature, enhancement, fix, cleanup, refactor}`，漏了就错
- `T-test-4.5` 中的 `.5` 后缀对人类不友好
- `breakdown-tasks/templates/` 中 15 个 `.md` 文件无 SKILL.md 或 Go 代码引用
- `quick-tasks/templates/` 中 10 个 `.md` 文件同样无引用
- `index.json` / `index.schema.json` 各 2 个，不被任何代码消费
- `validate-code-task.md` 指定 `type: "gate"`，但无 `TypeValidationCode` 常量

### Urgency

中——不会导致错误产出，但浪费时间（quality-gate 空转、改 ID/Type 时多处分发），且 validate-code/validate-ux 作为 eval-reality-validation 的关键防御层未自动生效。

## Proposed Solution

### Part A：前缀式 Type 重命名

```
coding.* → 涉及代码 → quality-gate 运行
doc.*    → 纯文档   → quality-gate 跳过
test.*   → 测试管线 → 特殊处理
validation.* → 验证任务 → 特殊处理
其他      → 元类型   → 特殊处理
```

| 新 type | 旧 type | 分类 |
|---------|---------|------|
| `coding.feature` | `feature` | code |
| `coding.enhancement` | `enhancement` | code |
| `coding.cleanup` | `cleanup` | code |
| `coding.refactor` | `refactor` | code |
| `coding.fix` | `fix` | code |
| `coding.clean` | `code-quality.simplify` | code |
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
| `gate` | `gate` | meta (不变) |
| `validation.code` | _(新增)_ | validation |
| `validation.ux` | _(新增)_ | validation |

`IsTestableType` 简化为前缀判断：

```go
func IsTestableType(typ string) bool {
    return strings.HasPrefix(typ, "coding.")
}
```

`isDocsOnlyType` 同理：

```go
func isDocsOnlyType(typ string) bool {
    return strings.HasPrefix(typ, "doc") // covers "doc", "doc.eval", etc.
}
```

### Part B：可读式 ID 重命名

| 旧 ID | 新 ID | type |
|-------|-------|------|
| `T-test-1` | `T-test-gen-cases` | `test.gen-cases` |
| `T-test-1b` | `T-test-eval-cases` | `test.eval-cases` |
| `T-test-2` | `T-test-gen-scripts` | `test.gen-scripts` |
| `T-test-3` | `T-test-run` | `test.run` |
| `T-test-4` | `T-test-graduate` | `test.graduate` |
| `T-test-4.5` | `T-test-verify-regression` | `test.verify-regression` |
| `T-specs-1` | `T-specs-consolidate` | `doc.consolidate` |
| `T-quick-1` | `T-quick-gen-cases` | `test.gen-cases` |
| `T-quick-2` | `T-quick-gen-and-run` | `test.gen-and-run` |
| `T-quick-3` | `T-quick-graduate` | `test.graduate` |
| `T-quick-4` | `T-quick-verify-regression` | `test.verify-regression` |
| `T-quick-specs-1` | `T-quick-doc-drift` | `doc.drift` |
| `T-clean-code-1` | `T-clean-code` | `coding.clean` |
| `T-eval-doc` | `T-eval-doc` _(不变)_ | `doc.eval` |
| _(新增)_ | `T-validate-code` | `validation.code` |
| _(新增)_ | `T-validate-ux` | `validation.ux` |

Profile 后缀保留：`T-test-gen-scripts-a`（profile a）、`T-test-gen-scripts-api`（type api）。

`InferType` 仍为映射表（ID 后缀 ≠ 类型名），但 key 从数字变为可读字符串。

`genScriptBases` 更新：`"T-test-gen-scripts"`, `"T-quick-gen-and-run"`。

### Part C：删除死模板（23 个文件）

**breakdown-tasks/templates/** 删除 14 个：

- `gen-test-cases.md`, `eval-test-cases.md`, `gen-test-scripts.md`, `run-e2e-tests.md`, `graduate-tests.md`, `verify-regression.md`, `consolidate-specs.md`
- `gate-task.md`, `phase-summary-task.md`, `fix-task.md`
- `validate-code-task.md`, `validate-ux-task.md`
- `index.json`, `index.schema.json`

保留：`task.md`, `task-doc.md`, `manifest-update-tasks.md`

**quick-tasks/templates/** 删除 9 个：

- `quick-test-cases.md`, `quick-gen-scripts.md`, `quick-run-tests.md`, `quick-graduate.md`, `quick-verify-regression.md`
- `validate-code-task.md`, `validate-ux-task.md`
- `index.json`, `index.schema.json`

保留：`task.md`, `task-doc.md`, `manifest-quick.md`

### Part D：接入 Validation 任务到 CLI 管线

1. `testgen.go`：`GetBreakdownTestTasks()` 和 `GetQuickTestTasks()` 新增 `T-validate-code` 和 `T-validate-ux` 生成
2. `build.go`：验证任务依赖链接线 + `auto.Validation` 配置
3. `prompt/data/`：新增 `validation-code.md` 和 `validation-ux.md` prompt 模板
4. `prompt.go`：`typeToTemplate` 新增映射
5. `infer.go`：`InferType()` 新增条目

Validation 任务定位：
- `T-validate-code`：在最后业务任务之后、test pipeline 之前；`noTest: true`, `mainSession: false`
- `T-validate-ux`：在 test pipeline 之后、all-completed 之前；`noTest: true`, `mainSession: true`
- 新增 `auto.Validation.Full` / `auto.Validation.Quick` 配置项，默认 `false`

### Part E：Prompt 模板文件重命名

`pkg/prompt/data/` 中 17 个文件重命名以匹配新 type，3 个不变，2 个新增：

| 旧文件名 | 新文件名 |
|---------|---------|
| `feature.md` | `coding-feature.md` |
| `enhancement.md` | `coding-enhancement.md` |
| `cleanup.md` | `coding-cleanup.md` |
| `refactor.md` | `coding-refactor.md` |
| `fix.md` | `coding-fix.md` |
| `documentation.md` | `doc.md` |
| `doc-evaluation.md` | `doc-eval.md` |
| `doc-generation-summary.md` | `doc-summary.md` |
| `doc-generation-consolidate.md` | `doc-consolidate.md` |
| `doc-generation-drift.md` | `doc-drift.md` |
| `test-pipeline-gen-cases.md` | `test-gen-cases.md` |
| `test-pipeline-eval-cases.md` | `test-eval-cases.md` |
| `test-pipeline-gen-scripts.md` | `test-gen-scripts.md` |
| `test-pipeline-run.md` | `test-run.md` |
| `test-pipeline-gen-and-run.md` | `test-gen-and-run.md` |
| `test-pipeline-graduate.md` | `test-graduate.md` |
| `test-pipeline-verify-regression.md` | `test-verify-regression.md` |

不变：`gate.md`，`fix-record-missed.md`，`clean-code.md`。

新增：`validation-code.md`，`validation-ux.md`。

## Requirements Analysis

### Key Scenarios

1. `forge task index` 在 coding feature 中生成新 ID（如 `T-test-gen-scripts`）和新 type（如 `test.gen-scripts`）
2. `forge task index` 在 docs-only feature 中只生成 `T-eval-doc`
3. `forge task index` 在 `auto.Validation` 启用时生成 `T-validate-code` 和 `T-validate-ux`
4. `IsTestableType("coding.feature")` → true；`IsTestableType("doc")` → false；`IsTestableType("validation.code")` → false
5. `run-tasks` 通过新 type 名查找重命名后的 prompt 模板
6. 旧 index.json 中的旧 type 名需兼容（InferType 兜底或一次性迁移）

### Constraints & Dependencies

| 约束 | 说明 |
|------|------|
| 旧 index.json 兼容 | 已有 feature 的 index.json 使用旧 type/ID。`forge task index` 重新生成即可迁移 |
| Profile 后缀解析 | `profileSuffixedID` / `typeSuffixedID` 的 base 参数需更新为新 ID 前缀 |
| `isAutoGenTaskID` | 需覆盖新 ID 前缀：`T-test-`, `T-quick-`, `T-specs-`, `T-clean-`, `T-validate-`, `T-eval-` |
| `isTestTaskID` | 需覆盖新 ID 前缀 |
| forge 版本 bump | Go 代码变更需 bump `scripts/version.txt`（minor：新功能 + breaking type name 变更） |

## Alternatives & Industry Benchmarking

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | 零成本 | 所有 4 个问题持续存在 | **Rejected** |
| 仅删死模板 | 减少混淆 | type/ID 问题不解决 | **Rejected**: 只解决 1/4 |
| 仅改 type/ID | 解决核心问题 | 死文件仍存在，validation 仍孤儿 | **Rejected**: 解决 2/4 |
| **全部合并** | 一次性彻底解决 | 变更范围大（~10 个文件） | **Selected**: 四个问题共享同一套代码，分开改反而多倍工作量 |

## Feasibility Assessment

### Technical Feasibility

完全可行。所有变更都是机械性的重命名 + 删除 + 新增，无架构变更。Modeled on existing patterns。

### Resource & Timeline

| Step | 工作量 |
|------|--------|
| types.go：重命名常量 + 新增 validation types | 20min |
| prompt/data/：重命名 18 个文件 + 新增 2 个 | 15min |
| prompt.go：更新 typeToTemplate + genScriptBases | 10min |
| testgen.go：更新所有 ID + 新增 validation 任务 | 40min |
| infer.go：更新所有 ID pattern + 新增 validation | 15min |
| build.go：IsTestableType 前缀化 + isAutoGenTaskID + 依赖链 | 30min |
| 删除 23 个死模板文件 | 5min |
| submit-task SKILL.md：quality-gate 检查 type | 10min |
| guide.md：quality-gate 协议 | 10min |
| Go 单元测试 | 40min |
| **Total** | **~3.5h** |

## Scope

### In Scope

**Go 代码（forge-cli/）：**
- `pkg/task/types.go`：重命名所有 type 常量 + 新增 `TypeValidationCode`, `TypeValidationUx` + 更新 Registry/ValidTypes
- `pkg/task/testgen.go`：重命名所有任务 ID + 新增 validation 任务生成 + 更新依赖链解析
- `pkg/task/infer.go`：更新所有 ID pattern + 新增 validation 条目 + 更新 profileSuffixedID/typeSuffixedID 的 base
- `pkg/task/build.go`：`IsTestableType` → 前缀判断 + `isDocsOnlyType` → 前缀判断 + 更新 `isTestTaskID`/`isAutoGenTaskID` + validation auto config
- `pkg/prompt/prompt.go`：更新 `typeToTemplate` + `genScriptBases`
- `pkg/prompt/data/`：重命名模板文件 + 新增 validation prompt 模板
- Go 单元测试覆盖上述变更
- `scripts/version.txt`：minor bump

**Skill 文档（plugins/forge/）：**
- 删除 `breakdown-tasks/templates/` 中 14 个死文件
- 删除 `quick-tasks/templates/` 中 9 个死文件
- `references/shared/type-assignment.md`：更新类型表 + 前缀规则
- `hooks/guide.md`：quality-gate 协议更新
- `skills/submit-task/SKILL.md`：quality-gate 检查更新

### Out of Scope

- eval SKILL.md 变更（已实现）
- rubric / expert 文件变更（已实现）
- validate-code/validate-ux 功能本身（已完成）
- `noTest` 字段废弃
- 已有 feature 目录批量迁移（`forge task index` 重新生成即可）
- `SKILL.md` 对业务任务模板的引用（已正确，无需改）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 旧 index.json 类型名不兼容 | High | M | `forge task index` 重新生成即可；InferType 兜底 |
| Profile 后缀解析 base 更新遗漏 | M | H | 更新 `genScriptBases` 和所有 `profileSuffixedID`/`typeSuffixedID` base 字符串 |
| `isAutoGenTaskID` 遗漏新 ID 前缀 | M | M | 用前缀判断替代硬编码列表 |
| Dependency 链硬编码旧 ID | M | H | 全文搜索旧 ID 并替换 |
| Prompt 模板文件名与新 type 不匹配 | L | H | `typeToTemplate` map 编译时校验（lookup miss 即报错） |

## Success Criteria

- [ ] `IsTestableType("coding.feature")` → true, `IsTestableType("doc")` → false, `IsTestableType("test.gen-cases")` → false
- [ ] `InferType("T-test-gen-cases")` → `"test.gen-cases"`, `InferType("T-validate-code")` → `"validation.code"`
- [ ] `forge task index` 生成新 ID（如 `T-test-gen-scripts-go-api`）和新 type（如 `test.gen-scripts`）
- [ ] `prompt.Synthesize()` 通过新 type 查找到重命名后的 prompt 模板
- [ ] `breakdown-tasks/templates/` 只剩 3 个文件：`task.md`, `task-doc.md`, `manifest-update-tasks.md`
- [ ] `quick-tasks/templates/` 只剩 3 个文件：`task.md`, `task-doc.md`, `manifest-quick.md`
- [ ] `forge task index` 在 `auto.Validation` 启用时生成 `T-validate-code` 和 `T-validate-ux`
- [ ] Go 单元测试覆盖新 type 常量、InferType、IsTestableType 前缀判断
