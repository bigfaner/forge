# Pre-Revision Synthetic Eval Report (Iteration 0)

**Title**: Pre-Revision (Freeform Findings)
**Iteration**: 0

## ATTACK_POINTS

### Factual Corrections (direct edits)

1. **[high]** 单 surface scalar 形式产生 dot-key 哨兵值，导致任务 ID 无效 | quote: "Scalar form `surfaces: api` produces `{\".\": \"api\"}`. Using surface-key as suffix would create `T-test-run-.`, which is invalid." | improvement: 明确声明 scalar 形式（单 surface）退化为无后缀 T-test-run，并添加 success criterion 验证此行为

2. **[high]** InferType 对 T-test-run 使用精确匹配，需改为 typeSuffixedID | quote: "The proposal does not mention this change. This is a **critical omission**." | improvement: 在 In Scope 中显式列出 InferType 和 infer_test.go 的改动

3. **[high]** 已有 fix-tasks 的 SourceTaskID 引用旧 T-test-run，迁移后指向不存在任务 | quote: "Existing fix-tasks with `SourceTaskID: \"T-test-run\"` would reference a non-existent task after the rename. The auto-restore logic would silently fail." | improvement: 在 In Scope 中添加 SourceTaskID 迁移步骤，或在 Constraints 中声明此为已知迁移风险

4. **[high]** YAML map key 可能含非法字符，未定义归一化规则 | quote: "Arbitrary YAML map keys may contain characters invalid for task IDs/filenames (e.g., `/`, spaces, uppercase). No normalization rules are defined." | improvement: 在 Requirements Analysis 中添加 surface-key 合法性约束或归一化规则

5. **[medium]** Quick 模式的 gen-journeys 合并未显式说明 | quote: "The proposal focuses on breakdown mode but quick mode has the same per-type gen-journeys loop (autogen.go line 229). The merge must apply to both modes." | improvement: 在 Solution 中明确声明"两种模式均受影响"

6. **[medium]** gen-journeys 合并后 TestType 字段无意义 | quote: "Single gen-journeys task has no meaningful TestType. The `renderBody` function's behavior with empty TestType differs from per-type tasks." | improvement: 在 In Scope 中说明 TestType 字段处理方式

### Structural Suggestions (edit when verifiable inconsistency found)

7. **[high]** 提案使用 surface-key 后缀与现有 surface-type 后缀命名不一致 | quote: "The existing system passes `capabilities []string` (deduplicated surface types) to both GetBreakdownTestTasks and GetQuickTestTasks" | improvement: 决定命名策略（surface-key vs surface-type）并在 Solution 中明确声明，或解释为何需要两种策略

8. **[high]** index.json 中已有 T-test-run 条目变更后成为孤儿 | quote: "Existing `T-test-run` entries in index.json would become orphans when replaced with `T-test-run-{key}`. Runtime state (status, blocked-reason) would be lost." | improvement: 在 Key Risks 中添加迁移风险，在 In Scope 或 Out of Scope 中明确处理策略

9. **[medium]** execution-order 验证时机未指定 | quote: "Where and when should invalid surface-key references be caught? Config load time (fail fast) or build time (contextual)?" | improvement: 在 Requirements Analysis Constraints 中明确验证时机（推荐 config load time）

10. **[medium]** 依赖拓扑变更的具体实现未展示 | quote: "Show the concrete chain for a 3-surface project in both breakdown and quick modes" | improvement: 在 Solution 中添加 3-surface 依赖链示意图

11. **[low]** 默认优先级未覆盖所有组合的回退行为 | quote: "doesn't cover all pairwise combinations. What about projects with tui + cli surfaces?" | improvement: 明确回退策略（alphabetical 或 config order）

12. **[low]** verify-regression 应只依赖最后一个 run-test 而非全部 | quote: "verify-regression should only depend on the LAST run-test task, not all of them, for correct scheduling semantics" | improvement: 修正 In Scope 中的描述为"依赖最后一个 run-test 子任务"

### Skipped (subjective preference)

None — all findings are backed by concrete code evidence.

### BORDERLINE_FINDINGS

- Finding #7 (naming inconsistency): partially structural — the concern is legitimate but choosing naming convention is a design decision, not a verifiable inconsistency. Marked as partially-accepted: apply the "declare the decision" portion.

## Rubric Data

All dimensions: N/A (pre-revision synthetic report)
