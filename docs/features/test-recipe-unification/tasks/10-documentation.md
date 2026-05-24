---
id: "10"
title: "Update project documentation for new test model"
priority: "P2"
estimated_time: "1h"
dependencies: ["6", "8", "9"]
type: "doc"
mainSession: false
---

# 10: Update project documentation for new test model

## Description

更新项目文档中的测试相关引用，包括 CLI docs（OVERVIEW.md, WORKFLOW.md 及其中文版）、ARCHITECTURE.md、business-rules/quality-gate.md、以及 conventions 文件。

## Reference Files
- `proposal.md#Impact-Analysis` — Tier 5 lists documentation files: CLI docs (4), ARCHITECTURE.md (1), quality-gate.md (1), conventions (2)
- `proposal.md#Proposed-Solution` — overall two-layer test model context for documentation updates
- `proposal.md#Success-Criteria` — no residual e2e-test/e2e-setup/e2e-verify in Go source, prompts, skill markdown

## Affected Files

### Modify
| File | Changes |
|------|---------|
| CLI docs (OVERVIEW.md, WORKFLOW.md, zh versions) | Update test command references, gate descriptions |
| `docs/ARCHITECTURE.md` | Update gate sequence description |
| `docs/business-rules/quality-gate.md` | Update gate steps and recipe names |
| `docs/conventions/testing/go.md` | Update test recipe conventions |
| `docs/conventions/forge-distribution.md` | Update any test-related references |

## Acceptance Criteria
- CLI docs reference `unit-test`, `test` (not `e2e-test`)
- ARCHITECTURE.md describes FullGateSequence, UnitGateSequence, NonBreakingGateSequence
- quality-gate.md reflects new gate steps and two-layer model
- No residual `e2eTest` or `e2e-test` references in documentation (historical lessons/proposals excluded)

## Implementation Notes
- 历史 lessons/proposals 中的 `e2eTest` 引用不在范围内（不影响功能）
- Low-priority: 可最后处理，不影响功能正确性
