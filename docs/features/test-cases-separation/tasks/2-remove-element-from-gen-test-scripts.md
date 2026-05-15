---
id: "2"
title: "Remove Element handling from gen-test-scripts and enforce source-code-first"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["1"]
type: "documentation"
mainSession: false
---

# 2: Remove Element handling from gen-test-scripts and enforce source-code-first

## Description
gen-test-scripts currently reads the Element field from test-cases.md and uses it as the primary locator source, with `sitemap-missing` as a fallback trigger to switch to source-code reconnaissance. This creates a dependency on gen-test-cases providing accurate Element IDs — which historically fails (provisional testids, unverified references).

Remove all Element field processing logic. Make source-code reconnaissance (Step 1.5 Fact Table) the sole locator source. The sitemap is retained as a secondary reference for page structure only, not as a locator source driven by test-case Element IDs.

## Reference Files
- `docs/proposals/test-cases-separation/proposal.md` — Source proposal
- `plugins/forge/skills/gen-test-scripts/` — Target skill directory

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | Remove Element field references in Sitemap section, Step 2 sitemap-missing handling, Step 3 locator priority "Sitemap element data" as primary. Add HARD-RULE enforcing source-code-first locator derivation. Update Integration Tests locator strategy reference to be source-code-driven. |

## Acceptance Criteria
- [ ] SKILL.md contains no reference to Element field from test-cases.md as a locator source
- [ ] `sitemap-missing` fallback logic is removed (Fact Table is always built, not just when sitemap is missing)
- [ ] Locator priority starts with Fact Table data from Step 1.5 Code Reconnaissance, sitemap is supplementary
- [ ] A HARD-RULE exists: "Derive all locators from source code (Fact Table). Do NOT reference any testid, selector, or locator from test-cases.md — test-cases provides scenario context only"
- [ ] Integration Tests locator strategy references are updated to source-code-driven (no Element field dependency)

## Implementation Notes
- Step 1.5 Code Reconnaissance and Fact Table Completeness Gate already exist — these become the primary locator mechanism
- The sitemap (Step 2) remains useful for page structure and route matching, but is no longer a locator source keyed by Element IDs from test-cases
- Remove the `// VERIFY: sitemap-missing` comment convention since it's no longer needed
- The Step 3 locator priority should be reordered: Fact Table first, sitemap as secondary confirmation only
- The proposal's In Scope item 5 (web-playwright generate.md) has no separate target file — the locator strategy is embedded in gen-test-scripts SKILL.md, which is covered by this task
