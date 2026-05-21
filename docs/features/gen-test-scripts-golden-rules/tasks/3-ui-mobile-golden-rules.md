---
id: "3"
title: "Restructure ui.md and rewrite mobile.md with golden rules"
priority: "P0"
estimated_time: "1-2h"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 3: Restructure ui.md and rewrite mobile.md with golden rules

## Description

Restructure ui.md and rewrite mobile.md into Golden Rules + Reconnaissance Hints dual-zone structure. mobile.md requires the most extensive changes — decoupling from Maestro YAML syntax to framework-agnostic principles. Add missing golden rules identified by expert evaluation.

## Reference Files
- `docs/proposals/gen-test-scripts-golden-rules/proposal.md` — Source proposal
- `plugins/forge/skills/gen-test-scripts/types/_shared.md` — Cross-type shared principles (created in task 1)
- `plugins/forge/skills/gen-test-scripts/types/ui.md` — Current UI type
- `plugins/forge/skills/gen-test-scripts/types/mobile.md` — Current Mobile type
- `plugins/forge/skills/gen-test-scripts/SKILL.md` — For step number reference alignment

## Affected Files

### Create
| File | Description |
|------|-------------|
| — | — |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-scripts/types/ui.md` | Restructure + add golden rules + fix references |
| `plugins/forge/skills/gen-test-scripts/types/mobile.md` | Rewrite to framework-agnostic principles |

### Delete
| File | Reason |
|------|--------|
| — | — |

## Acceptance Criteria

### Structural (both files)
- [ ] Each file has `## Golden Rules` section (framework-agnostic constraints)
- [ ] Each file has `## Reconnaissance Hints` section with discovery-only annotation
- [ ] Golden Rules reference `_shared.md` principles instead of redefining
- [ ] Output path changed from `tests/e2e/features/<feature>/` to `tests/<journey>/`
- [ ] All step number references corrected to match SKILL.md

### UI-specific additions
- [ ] Session Reuse: authentication scenarios must login once and reuse session context across tests; forbid per-test login
- [ ] Network Interception: for tests depending on external services, recommend intercepting network requests and returning fixed responses (framework-agnostic principle, not Playwright/Cypress-specific)
- [ ] Viewport Management: default viewport size for tests; responsive testing uses explicit viewport switching as pre-step

### Mobile-specific (rewrite)
- [ ] App State Reset: every test must clean app state before execution (kill + clearState or reinstall); forbid state leakage between tests
- [ ] Permission Handling: pre-authorize permissions or handle system permission dialogs as pre-step
- [ ] Deep Link Pattern: test opening app via URL scheme as a supported navigation entry
- [ ] Framework-agnostic generation patterns: describe touch/gesture/navigation principles (tap, swipe, navigate) without binding to Maestro YAML syntax
- [ ] Maestro YAML examples moved to Reconnaissance Hints as reference example, marked `<!-- Reference example for Maestro — not generation instructions -->`
- [ ] Element Location Strategy defined as priority chain: accessibility ID > resource ID > text — framework-agnostic, no Maestro `tapOn` syntax in Golden Rules

## Hard Rules

- mobile.md Golden Rules must contain ZERO Maestro-specific syntax (`tapOn`, `appId`, `onFlowStart`, YAML flow skeleton)
- Maestro examples are allowed ONLY in Reconnaissance Hints, clearly marked as reference examples
- ui.md must not contain Playwright-specific code (`page.locator`, `expect()`) in Golden Rules — those belong in Convention
- Element location strategies are defined as abstract priority chains, not concrete selector syntax

## Implementation Notes

- mobile.md is the heaviest rewrite — the existing file is essentially a Maestro tutorial (lines 83-161 are all Maestro YAML). The rewrite needs to extract the PRINCIPLES from these patterns (e.g., "locate elements by accessibility label before resource ID before text" from the `tapOn` priority table) and express them framework-agnostically
- ui.md is lighter — the existing Locator Mapping and Antipattern Guards sections are already well-structured, mainly need to add Session Reuse, Network Interception, Viewport Management to Golden Rules
- Keep existing Classification Indicators, Fact Table Required Keys, and Verification Method sections as-is
