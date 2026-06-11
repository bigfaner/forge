---
feature: "gen-test-scripts-golden-rules"
created: "2026-05-21"
status: completed
mode: quick
---

# Feature (Quick): gen-test-scripts-golden-rules

<!-- Status flow: tasks -> in-progress -> completed -->

## Overview

Restructure gen-test-scripts `types/` directory into a three-layer model (`_shared.md` -> type Golden Rules -> Convention) and integrate type loading into SKILL.md. Addresses three root problems: types/ not loaded by SKILL.md, principle/implementation mixing in type files, and missing cross-type golden rules (timeout, determinism, isolation).

**Target files:**
- `plugins/forge/skills/gen-test-scripts/types/_shared.md` (new)
- `plugins/forge/skills/gen-test-scripts/types/{cli,tui,ui,mobile,api}.md` (restructured)
- `plugins/forge/skills/gen-test-scripts/SKILL.md` (loading logic added)

## Documents

| Document | Path |
|----------|------|
| Proposal | ../../proposals/gen-test-scripts-golden-rules/proposal.md |

## Tasks

| ID | Title | Priority | Dependencies | Status | File |
|----|-------|----------|--------------|--------|------|
| 1 | Create types/_shared.md cross-type golden rules | P0 | — | pending | 1-shared-golden-rules.md |
| 2 | Restructure cli.md, tui.md, api.md with golden rules | P0 | 1 | pending | 2-cli-tui-api-golden-rules.md |
| 3 | Restructure ui.md and rewrite mobile.md with golden rules | P0 | 1 | pending | 3-ui-mobile-golden-rules.md |
| 4 | Integrate types/ into SKILL.md with loading logic | P0 | 1, 2, 3 | pending | 4-skill-md-integration.md |

## Acceptance Criteria Summary

- [ ] SKILL.md loads types/ files with `_shared.md` always + per-type selective loading
- [ ] Each type file split into `## Golden Rules` (framework-agnostic) + `## Reconnaissance Hints` (discovery-only)
- [ ] `_shared.md` defines 5 universal principles: Isolation, Determinism, Timeout Protection, Idempotency, Resource Cleanup
- [ ] All Output paths unified to `tests/<journey>/`, all step references corrected
- [ ] Zero framework-specific code in Golden Rules sections
