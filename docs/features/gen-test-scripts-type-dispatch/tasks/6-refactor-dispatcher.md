---
id: "6"
title: "Refactor main SKILL.md to dispatcher pattern"
priority: "P1"
estimated_time: "2-3h"
dependencies: ["1", "2", "3", "4", "5"]
type: "documentation"
mainSession: false
---

# 6: Refactor main SKILL.md to dispatcher pattern

## Description

Refactor the monolithic `gen-test-scripts` SKILL.md (530 lines) into a slim dispatcher (~200-250 lines). Remove all type-specific content (reconnaissance, Fact Table requirements, generation patterns, sitemap/locators) that was extracted into `types/{type}.md` files in tasks 1-5. Add a Step 4 dispatch loop that loads the appropriate type file based on the detected or filtered type.

## Reference Files
- `docs/proposals/gen-test-scripts-type-dispatch/proposal.md` — Source proposal
- `plugins/forge/skills/gen-test-cases/SKILL.md` — Reference dispatcher structure (~150 lines)
- `plugins/forge/skills/gen-test-scripts/types/*.md` — Type files created in tasks 1-5
- `plugins/forge/skills/gen-test-scripts/SKILL.md` — Current monolithic file to refactor

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | Refactor from 530-line monolith to ~200-250 line dispatcher |

## Acceptance Criteria

- [ ] SKILL.md line count is at most 250 lines (budget: Step 0 ~25 + Step 1 ~30 + Step 1.5 ~40 + Step 3.5 ~50 + Step 4 ~30 + frontmatter ~15 + conventions ~15 + errors ~20 = ~225 + 25 margin)
- [ ] `grep -c "gen-test-cases/types" SKILL.md` returns 0 — convention loading reads from own type files
- [ ] `grep -c "Sitemap" SKILL.md` returns 0 — Step 2 moved to types/ui.md
- [ ] `grep -c "Locator" SKILL.md` returns 0 — Step 3 moved to types/ui.md
- [ ] Convention loading section references `plugins/forge/skills/gen-test-scripts/types/{type}.md` frontmatter (not gen-test-cases type files)
- [ ] Step 4 contains a dispatch loop with explicit type-to-file mapping table (ui→types/ui.md, tui→types/tui.md, mobile→types/mobile.md, api→types/api.md, cli→types/cli.md)
- [ ] Error scenarios documented: unknown `--type` value, missing type file
- [ ] Step 3.5 (shared infrastructure) preserved verbatim — always runs regardless of type
- [ ] Steps 0-1 (profile resolution, test case reading, auth classification) preserved
- [ ] Step 1.5 retains generic Fact Table framework (build table, source citations, UNKNOWN handling) but delegates type-specific reconnaissance requirements to type files
- [ ] Frontmatter `conventions` field still lists project-wide conventions (`testing-isolation.md`)
- [ ] The `--type` filter documentation section is preserved (lines 57-88 in current SKILL.md)

## Hard Rules

- Step 3.5 (shared infrastructure) MUST NOT be moved to any type file — it is shared across all types
- The dispatcher must use hard-coded type-to-file mapping (not dynamic), matching the proposal's design decision
- Error messages must name the expected types and the missing file path — no silent fallback to generic instructions
- Preserve all existing HARD-RULEs that are type-agnostic (antipattern guard overview, traceability, empty result guard, VERIFY markers, post-generation checks)

## Implementation Notes

- This is the highest-risk task because it rewrites the main instruction file. All type-specific content must already exist in type files (tasks 1-5) before this task starts.
- Follow gen-test-cases SKILL.md as the structural template: frontmatter → title → core principle → HARD-GATE → Step 0 → Step 1 → Step 2.5 (interface detection) → Step 2.6 (convention loading) → Step 3 (per-type dispatch loop) → Step 4 (manifest)
- Key structural change: Steps 2-3 (sitemap/locators) are removed entirely from SKILL.md. Step 4 changes from "generate all types inline" to "load types/{type}.md and follow its instructions for each detected type"
- The convention loading section (currently lines 44-55) must be updated to read from own type files, not gen-test-cases type files
- Antipattern Guard section (lines 443-478): keep the generic 6-pattern table in SKILL.md; type-specific guards are in type files
- Integration Test Scripts section (lines 410-428): moved to types/ui.md (task 3)
- beforeAll Safety section (lines 430-440): keep in SKILL.md — applies to all types
