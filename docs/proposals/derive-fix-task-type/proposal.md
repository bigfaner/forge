---
created: "2026-05-29"
author: fanhuifeng
status: Draft
---

# Proposal: Derive Fix Task Type from Source Task Category

## Problem

The task dispatcher hardcodes `--type coding.fix` when creating fix tasks for all error recovery paths (agent timeout, blocked status, missing instructions), regardless of whether the source task is a doc, coding, or test type. This causes doc task failures to spawn `coding.fix` tasks that run irrelevant code-level quality gates (golangci-lint, go test), leading to cascading failures and infinite fix-task chains.

### Evidence

- 17 occurrences of hardcoded `coding.fix` across 7 plugin files (run-tasks.md, execute-task.md, task-executor.md, submit-task/SKILL.md, breakdown-tasks/SKILL.md, quick-tasks/SKILL.md, record-format-coding.md)
- Lesson recorded in `docs/lessons/gotcha-fix-task-type-hardcoded.md` with root cause analysis
- `forge task claim` already outputs `TYPE` and `TASK_CATEGORY` fields, but no skill file uses them for fix type derivation

### Urgency

Doc tasks in `/quick` and `/run-tasks` pipelines produce false failures when quality gates meant for code run against markdown-only changes. Each false failure spawns another `coding.fix` task, creating an infinite chain that blocks pipeline completion.

## Proposed Solution

1. Add `doc.fix` type to the Go type taxonomy (alongside existing `coding.fix`)
2. Define a category-based fix type derivation rule in all error-handling instructions
3. Update all 7 affected files to extract `TYPE` from claim output and derive the correct fix type

The derivation rule:

| Source Task Category | Fix Task Type | Rationale |
|----------------------|---------------|-----------|
| `doc`                | `doc.fix`     | 纯文档操作，修复也是改 `.md` 文件 |
| `eval`               | `doc.fix`     | 评估/修复的都是 `.md` spec 文件（journey.md、contract .md），无需改代码 |
| `coding`             | `coding.fix`  | 改代码 |
| `test`               | `coding.fix`  | 测试生成/运行失败需改代码 |
| `validation`         | `coding.fix`  | 验证失败需改代码 |
| `gate`               | `coding.fix`  | 门禁检查含编译/单元测试，失败需改代码 |

分类逻辑：`doc` 和 `eval` 类别的任务只操作 `.md` 文件，`IsTestableType()` 对它们返回 false，质量门禁自动跳过。其余类别的任务涉及代码，修复走 `coding.fix`。

### Innovation Highlights

Straightforward category-to-type mapping. Two insights:
1. `test`/`validation`/`gate` failures are essentially coding work → `coding.fix`
2. `eval` tasks only modify `.md` spec files → same bucket as `doc` → `doc.fix`

## Requirements Analysis

### Key Scenarios

- Doc review task (type: `doc.review`) gets blocked due to AC failure → fix task type is `doc.fix`
- Eval journey task (type: `eval.journey`) scores below threshold → fix task type is `doc.fix`
- Coding feature task (type: `coding.feature`) times out → fix task type is `coding.fix`
- Test generation task (type: `test.gen-scripts`) fails → fix task type is `coding.fix`
- Gate task (type: `gate`) unit test fails → fix task type is `coding.fix`

### Constraints & Dependencies

- `doc.fix` must be registered in Go type system (`types.go`) and validated by `IsValidType()`
- `doc-fix.md` template must exist for prompt synthesis
- `InferType()` must recognize `fix-` ID prefix for doc fix tasks
- `fixTypeFromStep()` in quality_gate.go already maps gate steps to fix types — only needs review for consistency

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | Zero cost | Infinite fix-task chains for doc pipelines | Rejected: active production issue |
| Reuse `doc` type (no new type) | No Go code changes | Loses traceability — can't distinguish original doc tasks from fix doc tasks; `fix-` ID prefix would conflict with InferType | Rejected: semantic ambiguity |
| **Add `doc.fix` + instruction-level derivation** | Clean type taxonomy; minimal Go changes; agent reads TYPE from claim output | 7 markdown files need updating | **Selected: best trade-off** |

## Feasibility Assessment

### Technical Feasibility

All infrastructure exists: `TYPE` field in claim output, `CategoryForType()` in Go, `fix-` ID prefix convention. Changes are additive (new type constant, new template, instruction updates).

### Resource & Timeline

Small scope: ~4 files of Go code changes, ~7 markdown instruction updates. Single developer, < 1 day.

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| All fix tasks are coding tasks | 5 Whys | Overturned: doc task failures need doc-type fixes, not code fixes |
| `eval` failures need code fixes | Codebase check | Overturned: eval tasks only modify `.md` spec files, `IsTestableType()` returns false → `doc.fix` |
| `test` failures need a `test.fix` type | XY Detection | Overturned: test failures require code changes → `coding.fix` is correct |
| `forge task claim` output lacks type info | Codebase check | Overturned: `TYPE` and `TASK_CATEGORY` fields exist but are unused |

## Scope

### In Scope

- Add `doc.fix` type constant to `types.go`
- Add `doc-fix.md` task template
- Update `InferType()` to handle doc fix IDs
- Update all 7 markdown files with hardcoded `coding.fix` to use derivation rule
- Document `TYPE` and `TASK_CATEGORY` as extractable fields in claim output references

### Out of Scope

- Adding `test.fix`, `eval.fix`, or other fix type variants
- Changes to quality gate logic
- Changes to `fixTypeFromStep()` in quality_gate.go (already correct for gate contexts)
- Template changes beyond adding `doc-fix.md`

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Agent misinterprets derivation rule in markdown | L | M | Use explicit if/else table, not prose description |
| `doc.fix` not recognized by downstream tooling | L | H | Register in `SystemTypes` map, `ValidTypes` slice, and `TaskTemplateDefaults` |
| Missed hardcoded `coding.fix` in a file | M | M | Grep for all occurrences before closing |

## Success Criteria

- [ ] `forge task add --type doc.fix --title "Fix: ..." --source-task-id <doc-task-id>` creates a valid task without validation errors
- [ ] When a `doc.review` task fails in `/run-tasks`, the spawned fix task has type `doc.fix` (not `coding.fix`)
- [ ] When an `eval.journey` task fails in `/run-tasks`, the spawned fix task has type `doc.fix` (not `coding.fix`)
- [ ] When a `coding.feature` task fails in `/run-tasks`, the spawned fix task has type `coding.fix` (unchanged behavior)
- [ ] When a `gate` task (unit test failure) fails in `/run-tasks`, the spawned fix task has type `coding.fix` (unchanged behavior)
- [ ] `forge task claim` output documentation in skill files lists `TYPE` and `TASK_CATEGORY` as extractable fields
- [ ] Zero remaining hardcoded `--type coding.fix` in error-handling instructions that should derive type dynamically

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
