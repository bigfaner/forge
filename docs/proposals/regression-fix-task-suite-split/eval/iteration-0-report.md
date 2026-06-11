---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Iteration 0: Pre-Revision (Freeform Findings)

## Findings Triage Summary

14 findings triaged (9 accepted, 0 partially-accepted, 2 deferred, 3 skipped)

| Finding | Severity | Status | Edit Summary |
|---------|----------|--------|-------------|
| SourceTaskID 与 countFixTasks 实际实现不一致 | high | accepted | Verify against current code, update 前置修复 description |
| 非测试文件路径输出行处理未定义 | medium | accepted | Clarify extractFileLineMap behavior for non-test-file lines |
| fallback 到 addFixTask 并非零功能损失 | high | accepted | Add fallback observability, update Key Risks |
| 一行匹配多测试文件导致输出污染 | high | accepted | Add containment: only create task for files with direct FAIL entries |
| 10 个并发 fix task 缺乏并发分析 | high | accepted | Add concurrent execution budget to Risks/NFR |
| 共享 createFixTask helper 行为漂移风险 | medium | accepted | Add test coverage requirement for shared helper |
| 误报 task 生命周期成本低估 | medium | accepted | Update cost estimate from "成本可控" to full lifecycle cost |
| Phase 0 无对应 success criteria | medium | accepted | Add SC for Phase 0 |
| RELATED_TASKS field 承诺但未列入 Scope | medium | accepted | Add to In Scope or remove from mitigation |
| 多语言模式无界误报面 | medium | deferred | MVP is Go-only, defer to multi-language extension |
| 不识别输出语言导致误报增长 | medium | deferred | Same — defer to multi-language extension |
| 建议：fallback 增加结构化日志 | low | skipped | Incorporated into fallback observability finding |
| 建议：按行号精确匹配 | low | skipped | Incorporated into output pollution containment |
| 建议：并发执行预算 | low | skipped | Incorporated into concurrent analysis finding |

**Skipped Findings Detail**:
- 建议：仅为有直接 FAIL 条目的文件生成 fix task — incorporated into output pollution containment finding
- 建议：实现前验证 countFixTasks — incorporated into SourceTaskID discrepancy finding

**Borderline Findings**:
- Phase 0 result may make Phase 1 unnecessary, but Phase 1 is fully designed with no decision gate. This is a structural observation about phasing strategy, not a factual error. Defer to Scorer cycle.

**Classification Audit**:
- Total findings by triage layer: 2 factual correction / 7 structural suggestion / 2 deferred / 3 incorporated-into-accepted

## ATTACK_POINTS

- **[high]** SourceTaskID 与 countFixTasks 实际实现不一致，前置修复效果存疑 | quote: "Either the lesson references an older version of the code, or the proposal is describing a bug that has already been partially fixed. This discrepancy needs resolution before implementation: if `countFixTasks` doesn't actually filter by SourceTaskID, then setting `opts.SourceTaskID = \"quality-gate\"` won't make the cap work differently." | improvement: Verify against current code whether countFixTasks uses SourceTaskID filtering or title-based matching. Update proposal's 前置修复 description to match actual code behavior.

- **[medium]** 非测试文件路径的输出行处理未定义，可能导致错误信息不完整 | quote: "The proposal does not specify what happens to output lines that reference only non-test files -- are they lost? Or are they included in the fallback task? This gap could cause incomplete error information in the per-file fix tasks." | improvement: Clarify in extractFileLineMap description whether non-test-file lines are folded into fallback task, included as context, or discarded with justification.

- **[high]** fallback 到 addFixTask 并非零功能损失，而是与待修复行为相同的静默回退 | quote: "The fallback is therefore not truly 'zero functionality loss' -- it is 'identical to the broken behavior we're trying to fix.' There is no observability into when fallback triggers." | improvement: Add fallback observability requirement: log structured warning when isTestFile returns zero matches. Update Key Risks mitigation from "零功能损失" to acknowledge fallback produces same monolithic task.

- **[high]** 一行匹配多测试文件导致输出污染，utils_test.go 场景下重新引入 broad-scope 问题 | quote: "If `utils_test.go` appears in 15 different test files' stack traces, the `utils_test.go` fix task receives output from all 15 test files' contexts, effectively recreating the broad-scope problem for that task." | improvement: Add containment strategy — only create fix task for files with direct FAIL entries, not files appearing only as stack trace references. Update Key Scenarios dismissal from "成本可控" to acknowledge the broad-scope recreation risk.

- **[high]** 移除 maxFixTasksPerStep 后 10 个并发 fix task 缺乏资源和并发冲突分析 | quote: "The proposal mentions RELATED_TASKS as a mitigation, but this is an informational hint to the agent, not a concurrency control mechanism." | improvement: Add concurrent execution budget analysis to Key Risks or Non-Functional Requirements. Specify max concurrent active tasks (e.g., 3) vs total created tasks (10).

- **[medium]** 共享 createFixTask helper 存在行为漂移风险，缺乏测试覆盖要求 | quote: "The proposal does not specify test coverage requirements for the shared helper to guard against this." | improvement: Add test coverage requirement for createFixTask helper in In Scope or Success Criteria.

- **[medium]** 误报 task 的完整生命周期成本被低估 | quote: "The cost is one full agent execution cycle plus another quality-gate run. In a system where quality-gate is a stop hook, this extends session time linearly with the number of false-positive tasks." | improvement: Update Key Scenarios cost estimate from "成本可控" to acknowledge full lifecycle cost (agent claim + execution + quality-gate re-run).

- **[medium]** Phase 0 无对应 success criteria，In Scope 列为 deliverable 但 SC 全部针对 Phase 1 | improvement: Add SC for Phase 0: "after Phase 0, agent sees complete --- FAIL line list in task description without reading raw-output.txt"

- **[medium]** RELATED_TASKS field 在 mitigation 中承诺但未出现在 In Scope 或 Out of Scope | improvement: Either add RELATED_TASKS field generation to In Scope, or remove from mitigation.

## Rubric Scores

All dimensions: N/A (freeform findings, no rubric scoring)
