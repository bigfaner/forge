---
id: "1"
title: "L1 Pilot Audit: README.md"
priority: "P0"
estimated_time: "1h"
dependencies: []
type: "doc"
complexity: "medium"
mainSession: false
---

# 1: L1 Pilot Audit: README.md

## Description
Perform pilot audit on README.md to validate audit methodology and produce accuracy baseline report. This is a prerequisite for all subsequent audit tasks — if pilot accuracy is below threshold, audit flow must be adjusted before proceeding.

Follow the L1/L2 audit methodology: declaration extraction → code location → item-by-item comparison → gap detection → result recording. Apply all three comparison types: path/file reference, behavior/process description, state/config assertion.

## Reference Files
- `docs/proposals/global-doc-code-audit/proposal.md` — Audit Execution Flow, Technical Feasibility, Constraints & Dependencies
- README.md: pilot audit target (ref: Audit Execution Flow)
- docs/conventions/forge-distribution.md: path resolution rules for validating file path claims (ref: Constraints & Dependencies)

## Affected Files

### Create
| File | Description |
|------|-------------|
| docs/features/global-doc-code-audit/audit/l1-pilot-report.md | Pilot audit report with accuracy baseline |

### Modify
| File | Changes |
|------|---------|
| (none) | Audit reads only, no modifications |

### Delete
| File | Reason |
|------|--------|
| (none) | No deletions |

## Acceptance Criteria
- [ ] All factual claims in README.md extracted (code paths, command names, config values, behavior descriptions)
- [ ] Each claim verified against actual codebase — path existence via `find`/`grep`, behavior via code reading
- [ ] Every inconsistency recorded with: file path, line range, severity (P0-P3), suggested action
- [ ] Accuracy baseline report produced: total claims examined, correct identifications, misses, false positives
- [ ] Miss rate < 20% (if ≥ 20%, report methodology adjustment recommendations and stop)

## Hard Rules
- Do NOT modify any code or documentation file — audit only
- Record the baseline commit hash (first audit task start time: latest commit on v3.0.0 branch) in the report header
- All audit output files (reports, task templates) must be written in English

## Implementation Notes
- Use the proposal's severity classification: P0 (leads to destructive AI action) → P3 (style only)
- Focus on three dimensions: outdated/incorrect, missing (code does X but doc doesn't say, or vice versa), redundant
- The pilot serves as quality gate — its results determine whether the full audit proceeds
