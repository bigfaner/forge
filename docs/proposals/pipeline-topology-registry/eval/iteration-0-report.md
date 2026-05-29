---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Pre-Revision (Freeform Findings)

**Iteration**: 0 (pre-revision)
**Source**: Freeform review by Go Pipeline Integration & Type System Engineer
**Findings triaged**: 20 (accepted: 11, partially-accepted: 0, deferred: 0, skipped: 9)

## ATTACK_POINTS

### Factual Corrections

- **[medium]** 运行时任务协调部分遗漏了 doc-fix- 前缀 | quote: "The runtime task coordination section lists only `fix-` and `disc-` prefixes. But `InferType` in the current code also handles `doc-fix-` prefix." | improvement: 在 Runtime Task Coordination 部分和 InferType fallback 逻辑中加入 `doc-fix-` 前缀

- **[high]** PipelineNode 缺失 AutoGenTaskDef 的关键字段 Key | quote: "The `PipelineNode` struct is missing fields that `AutoGenTaskDef` currently uses: `Key` (the map key in index.json), `FileName`, `Breaking`, and `StrategyContent`. The `Key` field is particularly critical because it determines the filename and lookup key." | improvement: 在 PipelineNode 中补充 Key 字段（或 KeyFunc），说明 Key 如何从 registry entry 派生

### Structural/Architectural Issues

- **[high]** T-clean-code 的空 DependsOn 违背单一真相来源承诺 | quote: "The `T-clean-code` entry in the registry has `DependsOn: []DepRef{}` with the comment 'resolved by caller: depends on last business task.' This represents a fundamental escape hatch from the registry's promise of being the single source of truth." | improvement: 定义 ResolveLastBusinessTask resolver 并增加 PrependToFirstTestDep 或 ReverseDep 机制

- **[high]** T-review-doc 反向依赖注入未在 registry 中体现 | quote: "The `T-review-doc` injection logic -- prepending to the first test task's dependencies when both doc and coding tasks exist -- is not captured in the registry." | improvement: 在 registry 设计中明确 capture 这个反向注入，或显式文档化它为什么留在 registry 外

- **[high]** 单 surface 退化场景下 run-test ID 处理遗漏 | quote: "The registry completely omits the single-surface degenerate case for run-test tasks. In the current codebase, when `isSingleSurface(surfaces)` is true, the run-test task gets `ID: 'T-test-run'` (no suffix). But the registry entry uses `ID: 'T-test-run-{surface-key}'`." | improvement: 添加 ExpansionRules 子节说明单 surface 如何退化

- **[high]** InferType 重构未解释 surfaces map 依赖 | quote: "The `InferType` refactoring proposes wildcard pattern matching, but the current `testRunSurfaceKeyMatch` requires looking up the suffix in a `surfaces` map at runtime. The proposal does not explain how the registry iteration receives the `surfaces` map." | improvement: 明确 InferType 如何获取 surfaces map（接受 GenContext 或保持两阶段）

- **[high]** build.go 中关键依赖解析函数未纳入 registry 范围 | quote: "The proposal's scope includes refactoring `build.go` steps 7/7.5/7.6 but does not address `resolveTestDepsAndInjectReviewDoc`, `ResolveFirstTestDep`, `findHighestGateOrSummary`, or `findMaxBusinessTaskID`. The codebase will have a hybrid architecture." | improvement: 显式文档化这些函数与 registry 的关系，或将其纳入 scope

- **[medium]** GenerateCondition nil 默认语义与下游任务实际行为不一致 | quote: "The proposal's registry entries for downstream tasks (validation, consolidation, drift, clean-code) all lack an explicit `GenerateCondition` field. The document states: 'GenerateCondition: nil defaults to CondHasTestableTasks.'" | improvement: 为所有下游任务显式设置 GenerateCondition（CondAlways 或 CondHasTestableTasks），移除 nil 默认规则

- **[medium]** init-time 验证混淆了静态引用和动态 resolver | quote: "The proposal mentions init-time validation of `DepRef.Resolve` functions, but these are dynamic and produce different IDs depending on runtime state. Static init-time validation cannot check dynamic references." | improvement: 拆分为两阶段验证：静态（Ref 引用完整性）和动态（runtime 一致性检查）

- **[medium]** ValidateAutogenTemplates 与新 init-time 验证的关系未说明 | quote: "The existing `ValidateAutogenTemplates` function is not mentioned. Two separate validation passes at startup would be wasteful and confusing." | improvement: 统一到单一 ValidatePipeline 函数，或显式文档化两者的关系

- **[medium]** CondAlways 已定义但从未使用 | quote: "The proposal defines `CondAlways` but never uses it. This raises the question: when would it be used?" | improvement: 在 registry 中使用 CondAlways 或移除它

## BORDERLINE_FINDINGS

- claim_test.go 是否真正受影响需进一步验证 | 评审者指出 claim.go 不引用被修改函数，但可能通过 BuildIndex 间接影响

## SKIPPED_FINDINGS

- 建议：为 PipelineNode 定义显式 KeyDerivation 策略 → 结构建议，已通过 ATTACK_POINT 覆盖关键性
- 建议：用 ResolveLastBusinessTask 替代空 DependsOn → 已通过 ATTACK_POINT 覆盖
- 建议：增加 ExpansionRules 子节 → 已通过 ATTACK_POINT 覆盖
- 建议：加入 doc-fix-* 前缀 → 已通过 ATTACK_POINT 覆盖
- 建议：统一 init-time 验证 → 已通过 ATTACK_POINT 覆盖
- 建议：为下游任务显式设置 CondAlways → 已通过 ATTACK_POINT 覆盖
- 建议：解决 InferType surfaces map 依赖 → 已通过 ATTACK_POINT 覆盖
- 建议：文档化 registry 与现有函数关系 → 已通过 ATTACK_POINT 覆盖
- 建议：claim_test.go 范围精确性 → 低优先级，可在实现阶段验证

## Classification Audit

| Triage Layer | Count |
|---|---|
| Factual correction | 2 |
| Structural/architectural suggestion | 9 |
| Subjective preference | 9 (skipped) |
