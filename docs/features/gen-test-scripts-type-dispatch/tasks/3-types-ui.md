---
id: "3"
title: "Create types/ui.md instruction file (includes Step 2-3 Sitemap/Locators)"
priority: "P1"
estimated_time: "2-3h"
dependencies: []
type: "documentation"
mainSession: false
---

# 3: Create types/ui.md instruction file (includes Step 2-3 Sitemap/Locators)

## Description

Extract UI-specific generation logic from the monolithic `gen-test-scripts` SKILL.md into a dedicated `types/ui.md` type instruction file. This is the most complex type file because it includes Step 2 (Sitemap resolution) and Step 3 (Locator mapping) — UI-exclusive logic that is currently in the main SKILL.md.

Modeled after `gen-test-cases/types/ui.md` structure.

## Reference Files
- `docs/proposals/gen-test-scripts-type-dispatch/proposal.md` — Source proposal
- `plugins/forge/skills/gen-test-cases/types/ui.md` — Reference structure (gen-test-cases UI type file)
- `plugins/forge/skills/gen-test-scripts/SKILL.md` — Source of UI-specific content: Step 2 (lines 292-304), Step 3 (lines 306-327), UI reconnaissance (lines 258), UI verification (line 389), integration tests (lines 410-428)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/gen-test-scripts/types/ui.md` | UI type instruction file with conventions frontmatter, including sitemap + locator logic |

## Acceptance Criteria

- [ ] `plugins/forge/skills/gen-test-scripts/types/ui.md` exists
- [ ] Frontmatter declares `type: ui` and `conventions: [testing-ui.md, frontend.md]`
- [ ] Contains a **Reconnaissance Strategy** section with UI-specific search patterns (grep data-testid, component files, route configs, dynamic testid patterns)
- [ ] Contains a **Fact Table Required Keys** section listing minimum keys for UI type (FRONTEND_BASE or TESTID_* entries)
- [ ] Contains a **Sitemap Resolution** section (Step 2 equivalent): reading sitemap.json, matching routes, collecting page structure, handling missing routes
- [ ] Contains a **Locator Mapping** section (Step 3 equivalent): locator priority chain (Fact Table > Sitemap > Semantic inference), dynamic testid handling, integration test component locators
- [ ] Contains a **Verification Method** section describing how to confirm the project exposes a UI (grep data-testid in frontend source, TESTID_* or FRONTEND_* Fact Table entries)
- [ ] Contains a **Generation Patterns** section describing how UI test cases translate to executable scripts (DOM interaction, navigation, form submission, visibility assertions, screenshot capture)
- [ ] Contains an **Integration Test Scripts** section for existing-page placement tests (locate page by route, locate component, assert visibility + data rendering)
- [ ] Contains a **UI Antipattern Guards** section (CSS class selectors, fixed delays/sleeps, debug output, missing testid fallback rules)
- [ ] At least 3 section headings are unique to this file (Sitemap Resolution, Locator Mapping, Integration Test Scripts are UI-only)

## Hard Rules

- Sitemap and locator content must be moved verbatim from current SKILL.md Steps 2-3, not paraphrased — this is a structural move, not a rewrite
- The locator priority chain (Fact Table > Sitemap > Semantic inference) must be preserved exactly
- `data-testid` derivation rules: every testid locator must come from Fact Table, never guessed

## Implementation Notes

- This is the largest type file (~120-150 lines) because it absorbs two full steps from the current SKILL.md
- Step 2 (Sitemap, lines 292-304): move entirely, including the sitemap field descriptions and WARNING behavior for missing routes
- Step 3 (Locators, lines 306-327): move entirely, including locator priority chain, dynamic testid handling, integration test locators, and HARD-RULEs about source-code-first derivation
- Integration Test Scripts section (lines 410-428): move to this file since it's UI-specific
- UI-specific HARD-RULEs (lines 319-324): move these into the type file's antipattern guards
