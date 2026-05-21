---
id: "4"
title: "Integrate types/ into SKILL.md with loading logic"
priority: "P0"
estimated_time: "1-2h"
dependencies: ["1", "2", "3"]
scope: "all"
breaking: false
type: "doc"
mainSession: false
---

# 4: Integrate types/ into SKILL.md with loading logic

## Description

Modify gen-test-scripts SKILL.md to load type files from `types/` directory. Add explicit step for type rule loading between Step 2 (Read Contract Specifications) and Step 3 (Generate Test Code). Declare the principle/implementation layering between types/ and Convention.

## Reference Files
- `docs/proposals/gen-test-scripts-golden-rules/proposal.md` — Source proposal
- `plugins/forge/skills/gen-test-scripts/SKILL.md` — Target file
- `plugins/forge/skills/gen-test-scripts/types/_shared.md` — Shared principles (created in task 1)
- `plugins/forge/skills/gen-test-scripts/types/cli.md` — Example restructured type file (from task 2)
- `plugins/forge/skills/gen-test-cases/SKILL.md` — Reference: how gen-test-cases loads its types/

## Affected Files

### Create
| File | Description |
|------|-------------|
| — | — |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | Add type loading step + priority declaration |

### Delete
| File | Reason |
|------|--------|
| — | — |

## Acceptance Criteria

- [ ] New step added between Step 2 and Step 3: "Step 2.5: Load Type Rules"
- [ ] Step 2.5 reads Contract files, extracts all interface types referenced
- [ ] Step 2.5 loads `_shared.md` (always) + matching type files from `types/` (only for detected interface types)
- [ ] Step 2.5 emits WARNING if >3 types detected (token budget risk)
- [ ] `<HARD-RULE>` added declaring types/ vs Convention priority: "types/ Golden Rules define non-overridable principle constraints. Convention provides framework implementation details. When both cover the same aspect, Golden Rules' principles take precedence, Convention's implementation details supplement areas Golden Rules don't cover."
- [ ] `<HARD-RULE>` added: "Reconnaissance Hints in type files are discovery aids only. Information discovered via Hints must be converted to Fact Table values, not used directly in generation instructions."
- [ ] Step 2.5 is consistent with gen-test-cases's type loading pattern (per-type dispatch, not bulk load)

## Hard Rules

- Loading timing MUST be after Step 2 (Contract read, which provides interface type information) and before Step 3 (code generation)
- `_shared.md` is ALWAYS loaded regardless of detected types
- Only type files matching detected interface types are loaded — no speculative bulk loading
- The HARD-RULE for types/ vs Convention priority MUST appear in SKILL.md, not just in type files

## Implementation Notes

- Reference gen-test-cases SKILL.md to see how it implements type dispatch — follow the same pattern for consistency
- The Domain Reconnaissance table in Step 1.2 of SKILL.md already mentions "CLI entry points" and "TUI components" — these stay as-is, they're about code reconnaissance, not type rule loading
- Step 1.2's Domain Reconnaissance and Step 2.5's type loading are complementary: Step 1 discovers project facts, Step 2.5 loads principle constraints for the discovered types
