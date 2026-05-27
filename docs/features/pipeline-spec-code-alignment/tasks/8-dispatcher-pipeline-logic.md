---
id: "8"
title: "Fix dispatcher and pipeline logic in docs"
priority: "P0"
estimated_time: "2h"
dependencies: [7]
type: "doc"
mainSession: false
---

# 8: Fix dispatcher and pipeline logic in docs

## Description

Fix dispatcher output handling, pipeline logic, and conditional execution in skill/command documentation. Covers Cluster 3 (issues C1-C12):

1. **run-tasks.md conditional messages** (line 117): Post-loop message hardcodes `T-test-run`/`T-test-verify-regression` task names which don't exist in quick mode. Fix: conditionally reference actual test task names based on mode.

2. **run-tasks.md summary format** (lines 37, 109): "print summary" has no defined format. Define a structured summary format.

3. **run-tasks.md timeout mechanism** (line 110): Agent timeout "Mark blocked" has no specified mechanism. Define the blocking procedure.

4. **quick.md knowledge extraction** (lines 151, 165): Claims run-tasks has "knowledge extraction" but run-tasks.md has no such step. Remove the false claim.

5. **execute-task.md summary format** (line 104): "Output your final summary" has no format. Define the expected format.

6. **execute-task.md status distinction** (line 66): All STATUS≠completed are handled the same, missing `in_progress`→record-missing recovery. Add explicit status branches.

7. **execute-task.md subagent_type** (line 113): Agent call omits `subagent_type="forge:task-executor"`. Add it.

8. **execute-task.md MAIN_SESSION** (line 35): Inconsistent instructions handling vs run-tasks.md (lines 62-65). Unify.

9. **task-executor.md DONE format** (lines 62, 68): Field positions inconsistent — commit-hash and status occupy the same position. Unify the DONE output format.

10. **gen-test-scripts SKIP_EVAL_GATE** (SKILL.md line 29): Missing SKIP_EVAL_GATE mode that gen-contracts has (lines 46-52). Quick mode gets blocked by eval gate in gen-test-scripts. Add conditional SKIP_EVAL_GATE.

11. **run-tests.md surface detection** (line 66): Passes task file path to `forge surfaces --json` which expects source directory paths. Fix to use correct paths.

12. **run-tests.md Step 0 reference** (line 187): References "Convention loaded in Step 0" but Step 0 is Stale State Recovery, not convention loading.

13. **run-tests.md Chinese text** (lines 122, 165): Contains Chinese text ("编排序", "测试环境异常") while rest is English. Translate to English.

## Reference Files
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Problem` — Evidence C1-C12 (all dispatcher and pipeline issues)
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Proposed-Solution` — Cluster 3 description
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Success-Criteria` — SC for conditional messages, SKIP_EVAL_GATE, format definitions

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/commands/run-tasks.md` | Conditional messages, summary format, timeout mechanism |
| `plugins/forge/commands/quick.md` | Remove false knowledge extraction claim |
| `plugins/forge/commands/execute-task.md` | Summary format, status distinction, subagent_type, MAIN_SESSION |
| `plugins/forge/agents/task-executor.md` | Unified DONE format |
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | Add SKIP_EVAL_GATE condition |
| `plugins/forge/skills/run-tests/SKILL.md` | Surface detection path, Step 0 ref, translate Chinese |

## Acceptance Criteria
- [ ] Post-loop message in run-tasks.md reflects actual task names (conditional on mode)
- [ ] Summary format defined in run-tasks.md
- [ ] Timeout/blocking mechanism specified in run-tasks.md
- [ ] quick.md does not claim run-tasks has knowledge extraction
- [ ] execute-task.md has explicit status branches (completed, blocked, in_progress)
- [ ] execute-task.md includes `subagent_type="forge:task-executor"` in agent call
- [ ] task-executor.md DONE format is consistent (no ambiguous field positions)
- [ ] gen-test-scripts/SKILL.md has SKIP_EVAL_GATE for Quick mode
- [ ] run-tests/SKILL.md uses source directory paths for surface detection
- [ ] No Chinese text in run-tests/SKILL.md

## Hard Rules
- Do not change Go code — this task is doc-only
- SKIP_EVAL_GATE must only apply in Quick mode, not in full pipeline mode

## Implementation Notes
- For SKIP_EVAL_GATE: mirror the pattern used in gen-contracts/SKILL.md lines 46-52
- For DONE format: choose one field order and document it clearly, then update all DONE examples in task-executor.md
- For Chinese text: translate to English to match the rest of the file's language
