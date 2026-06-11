---
id: "6"
title: "L2 Conventions Audit Batch 2"
priority: "P1"
estimated_time: "1.5h"
dependencies: [1]
type: "doc"
complexity: "high"
mainSession: false
---

# 6: L2 Conventions Audit Batch 2

## Description
Audit remaining docs/conventions/ files (7 top-level + 3 testing subdirectory = 10 files): naming.md, package-organization.md, prompt-template-hierarchy.md, skill-self-containment.md, skill-structure.md, surface-cli.md, surface-rules.md, testing/index.md, testing/cli/index.md, testing/cli/core.md.

## Reference Files
- `docs/proposals/global-doc-code-audit/proposal.md` — Audit Execution Flow, Constraints & Dependencies, Key Risks
- docs/conventions/naming.md: naming conventions vs actual code (ref: Audit Execution Flow)
- docs/conventions/skill-structure.md: skill structure rules vs actual SKILL.md files (ref: Audit Execution Flow)
- docs/conventions/testing/: test-related conventions vs actual test infrastructure (ref: Key Risks)

## Affected Files

### Create
| File | Description |
|------|-------------|
| docs/features/global-doc-code-audit/audit/l2-conventions-batch2-report.md | L2 conventions batch 2 audit report |

### Modify
| File | Changes |
|------|---------|
| (none) | Audit reads only |

### Delete
| File | Reason |
|------|--------|
| (none) | No deletions |

## Acceptance Criteria
- [ ] All 10 target files audited with declaration extraction
- [ ] Each convention verified against codebase: naming patterns vs actual code, skill structure vs actual files, test conventions vs actual test setup
- [ ] Every inconsistency recorded: file path, line range, severity (P0-P3), suggested action
- [ ] Cross-layer influence items recorded for L3 reference
- [ ] Audit report follows unified template

## Hard Rules
- Do NOT modify any code or documentation — audit only
- Only audit the following 10 files: naming.md, package-organization.md, prompt-template-hierarchy.md, skill-self-containment.md, skill-structure.md, surface-cli.md, surface-rules.md, testing/index.md, testing/cli/index.md, testing/cli/core.md
- All audit output written in English

## Implementation Notes
- naming.md is known to have discrepancies (proposal evidence: constant name mismatches with CLI code)
- skill-structure.md and skill-self-containment.md: verify rules against actual SKILL.md files in plugins/forge/skills/
- testing/ conventions: verify test path references (tests/<journey>/ vs legacy tests/e2e/)
