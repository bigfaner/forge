---
id: "8"
title: "Update prompt templates for new recipe names"
priority: "P1"
estimated_time: "30min"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 8: Update prompt templates for new recipe names

## Description

更新 3 个 prompt 模板文件中的 recipe 引用，将 `just test` 替换为 `just unit-test`（用于 gate 场景的快速反馈提示）。

## Reference Files
- `proposal.md#Impact-Analysis` — Tier 2 lists prompt template files: gate.md, fix-record-missed.md, validation-code.md
- `proposal.md#Proposed-Solution` — defines recipe name changes: test → unit-test for gate-level usage

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `pkg/prompt/data/gate.md` | `just test` → `just unit-test` |
| `pkg/prompt/data/fix-record-missed.md` | `just test` → `just unit-test` |
| `pkg/prompt/data/validation-code.md` | `just test` → `just unit-test` |

## Acceptance Criteria
- All 3 prompt templates reference `just unit-test` instead of `just test` for per-task gate scenarios
- No residual `just test` references in gate/fix/validation prompt contexts

## Implementation Notes
- These are .md files embedded via `go:embed` — text changes only
- Distinguish between gate context (unit-test) and all-completed context (test) if applicable
