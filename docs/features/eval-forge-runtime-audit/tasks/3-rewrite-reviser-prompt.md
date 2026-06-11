---
id: "3"
title: "Rewrite reviser-prompt.md: two-layer fix strategy (safe-fix + guided-fix)"
priority: "P0"
estimated_time: "1h"
dependencies: [1]
type: "documentation"
mainSession: false
---

# 3: Rewrite reviser-prompt.md: two-layer fix strategy (safe-fix + guided-fix)

## Description
Replace the current single-layer mechanical reviser with a two-layer strategy: safe-fix (mechanical, no semantic change) and guided-fix (rule-based fixes with clear precedence rules). The current reviser can only fix frontmatter, references, CLI flags, and status values. The new reviser needs to handle instruction conflicts, content deduplication, and bypass hardening.

## Reference Files
- `docs/proposals/eval-forge-runtime-audit/proposal.md` — Source proposal (Section 3: Reviser Two-Layer Fix)
- `.claude/skills/eval-forge/templates/reviser-prompt.md` — Current reviser prompt (to be rewritten)
- `.claude/skills/eval-forge/templates/rubric.md` — New rubric (rewritten by Task 1)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `.claude/skills/eval-forge/templates/reviser-prompt.md` | Complete rewrite: two-layer fix strategy with 3 rules |

## Acceptance Criteria
- [ ] Reviser prompt defines two fix layers: safe-fix and guided-fix
- [ ] safe-fix covers: frontmatter missing fields, name-directory mismatch, CLI flag corrections, dead reference removal
- [ ] guided-fix implements 3 rules:
  - Rule 1: Instruction conflicts → guide.md wins; SKILL.md changes to reference; if guide.md lacks the concept, migrate most complete version to guide.md
  - Rule 2: Content dedup → keep authoritative version; eval loop protocol extracts to `references/shared/eval-loop-protocol.md`
  - Rule 3: Bypass hardening → add minimal HARD-RULE with specific consequences ("If you skip X, Y will fail because Z"), not empty prohibitions
- [ ] HARD-RULE block preserved: only fix listed attack points, no refactoring, preserve existing intent
- [ ] Output format: FIXES APPLIED + FIXES SKIPPED (requires human judgment)
- [ ] Reviser reads the audit report at `docs/self-evolution/{seq}/iteration-{N}.md` and rubric at `.claude/skills/eval-forge/templates/rubric.md`

## Hard Rules
- safe-fix must never change semantics — only mechanical corrections
- guided-fix Rule 1 must always prefer guide.md as authority
- guided-fix Rule 3 must never add empty prohibitions — every HARD-RULE addition must include specific failure consequence
- Do NOT create new files except `references/shared/eval-loop-protocol.md` when Rule 2 triggers eval protocol extraction

## Implementation Notes
- Current reviser has 12 fix rule categories — consolidate into the two layers
- The guided-fix rules map to the new rubric's D2 (bypass), D3 (precision), D4 (dedup) dimensions
- Command metadata fixes and plugin metadata fixes from current reviser are no longer needed (D6 is only 50 pts, low priority)
