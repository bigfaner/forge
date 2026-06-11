---
id: "5"
title: "eval-ui multi-platform rubrics and selection logic"
priority: "P2"
estimated_time: "1.5h"
dependencies: ["3"]
scope: "all"
breaking: false
type: "implementation"
mainSession: false
---

# 5: eval-ui multi-platform rubrics and selection logic

## Description
Create 3 independent platform rubrics (rename existing to web, create new mobile and tui) and modify eval-ui SKILL.md to select the correct rubric based on platform. Each rubric has 4 dimensions x 250 points = 1000 total, with platform-specific scoring criteria and deduction rules.

This task resolves the current situation where a single rubric.md attempts to cover all platforms with different dimensions.

## Reference Files
- `docs/proposals/tui-ui-design/proposal.md` — Source proposal (D8, D9, D10 sections)
- `plugins/forge/skills/eval-ui/SKILL.md` — Current eval-ui skill logic
- `plugins/forge/skills/eval-ui/templates/rubric.md` — Current single rubric (to be renamed to rubric-web.md)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/eval-ui/templates/rubric-web.md` | Renamed from rubric.md — existing web evaluation rubric |
| `plugins/forge/skills/eval-ui/templates/rubric-mobile.md` | Mobile rubric: Requirement Coverage, Touch Experience, Adaptive Layout, Implementability (proposal D10) |
| `plugins/forge/skills/eval-ui/templates/rubric-tui.md` | TUI rubric: Requirement Coverage, Terminal Experience, Visual Specification, Implementability (proposal D9) |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval-ui/SKILL.md` | Add platform detection and rubric path selection: web → rubric-web.md, mobile → rubric-mobile.md, tui → rubric-tui.md (proposal D8) |

### Delete
| File | Reason |
|------|--------|
| — | — |

## Acceptance Criteria
- [ ] `rubric-web.md` contains the existing web rubric content (renamed from rubric.md)
- [ ] `rubric-tui.md` has 4 dimensions: Requirement Coverage (250), Terminal Experience (250), Visual Specification (250), Implementability (250) per proposal D9
- [ ] `rubric-tui.md` deduction rules: missing ASCII mockup → panel Visual Specification = 0; "待定" characters → -30/instance; missing mandatory edge case → -50/each; vague dimensions → -20/instance
- [ ] `rubric-mobile.md` has 4 dimensions: Requirement Coverage (250), Touch Experience (250), Adaptive Layout (250), Implementability (250) per proposal D10
- [ ] `rubric-mobile.md` deduction rules: touch targets without size → -30/instance; missing landscape/portrait → -50; missing safe area → -40
- [ ] `eval-ui/SKILL.md` detects platform from ui-design document and selects matching rubric file
- [ ] Multi-platform features evaluate each platform's ui-design file with its respective rubric

## Implementation Notes
- The existing `rubric.md` becomes `rubric-web.md` — rename only, no content changes to the web rubric
- Study how eval-ui SKILL.md currently reads the rubric to understand where to add the platform-based path selection
- The TUI rubric's Visual Specification dimension directly enforces the 5 structural requirements from the lesson — this is the quality gate that prevents incomplete TUI specs
- Each rubric is self-contained (1000 points, 4 dimensions, own deduction rules) — no cross-referencing between rubrics
