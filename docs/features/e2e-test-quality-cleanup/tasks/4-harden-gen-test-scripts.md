---
id: "4"
title: "Harden gen-test-scripts SKILL.md to forbid antipatterns"
priority: "P2"
estimated_time: "30m"
dependencies: ["3"]
scope: "all"
breaking: false
type: "documentation"
mainSession: false
---

# 4: Harden gen-test-scripts SKILL.md to forbid antipatterns

## Description

Update `plugins/forge/skills/gen-test-scripts/SKILL.md` to explicitly prohibit the 6 antipatterns during test script generation. This ensures agents generating test scripts never produce recursive tests, dead skips, vacuous assertions, etc.

## Reference Files
- `docs/proposals/e2e-test-quality-cleanup/proposal.md` — Source proposal
- `plugins/forge/skills/eval/rubrics/test-cases.md` — Updated rubric (Task 3)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | Add "Forbidden Patterns" section with 6 antipatterns |

## Acceptance Criteria
- [ ] `gen-test-scripts/SKILL.md` has a "Forbidden Patterns" or "Antipattern Guard" section
- [ ] Each of the 6 antipatterns is listed with: what it is, why it's harmful, what to do instead
- [ ] The section is referenced in the generation flow (e.g., "Before writing each test, verify it does not match any forbidden pattern")
- [ ] References the lesson documents as authoritative sources

## Hard Rules
- Do NOT change the generation flow structure — only add validation rules
- Do NOT add new files — only modify SKILL.md

## Implementation Notes
- The 6 antipatterns mirror the rubric dimension added in Task 3:
  1. Recursive test invocation → use recursion guard env var
  2. Unconditional t.Skip → either implement with fixture or don't generate
  3. Vacuous assertions → assertions must always execute, no conditional guards
  4. Conditional skip without fixture → every test must set up its own world via `t.TempDir()`
  5. Duplicate tests → scan existing test files before generating
  6. Static file text checks → only test runtime behavior, not source file content
