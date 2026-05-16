---
id: "3"
title: "Extend test-cases rubric with antipattern detection dimension"
priority: "P2"
estimated_time: "45m"
dependencies: []
scope: "all"
breaking: false
type: "documentation"
mainSession: false
---

# 3: Extend test-cases rubric with antipattern detection dimension

## Description

Add a new scoring dimension to `plugins/forge/skills/eval/rubrics/test-cases.md` that penalizes test scripts containing known antipatterns. This prevents `/gen-test-scripts` from producing the same low-quality patterns identified in the lesson.

The current rubric has 5 dimensions totaling 1000 points:
- PRD Traceability (250)
- Step Actionability (250)
- Interface Accuracy (200)
- Completeness (200)
- Structure & ID Integrity (100)

Add a 6th dimension: **Test Code Quality** that checks for the 6 antipatterns.

## Reference Files
- `docs/proposals/e2e-test-quality-cleanup/proposal.md` — Source proposal
- `docs/lessons/gotcha-recursive-go-test-process-explosion.md` — Recursion antipattern
- `docs/lessons/gotcha-e2e-test-quality-antipatterns.md` — All antipatterns

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/rubrics/test-cases.md` | Add "Test Code Quality" dimension with antipattern checks |

## Acceptance Criteria
- [ ] `rubrics/test-cases.md` includes a "Test Code Quality" dimension
- [ ] The dimension checks for all 6 antipatterns:
  1. Recursive test invocation (e.g., `exec.Command("go", "test"` without recursion guard)
  2. Unconditional `t.Skip` without environment-detection rationale
  3. Vacuous assertions (`if cond { assert }` where assertion may never execute)
  4. Conditional skip without self-contained fixture setup
  5. Duplicate test function names across packages
  6. Reading static source files and asserting text content (not testing runtime behavior)
- [ ] Total rubric points remain 1000 (redistribute from existing dimensions, e.g., reduce Completeness from 200 to 150 and Structure from 100 to 50, add Test Code Quality at 200)

## Hard Rules
- Do NOT change the rubric file path or overall structure
- Keep existing dimension definitions intact (only adjust point values)

## Implementation Notes
- The 6th dimension acts as a blocking gate similar to "Step Actionability < 200". If Test Code Quality score is below a threshold, the generated scripts should be rejected.
- Reference the lesson documents as the authoritative source for antipattern definitions.
