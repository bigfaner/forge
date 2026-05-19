---
id: "3"
title: "Update prompt pipeline"
priority: "P1"
estimated_time: "25min"
dependencies: ["1"]
scope: "backend"
breaking: true
type: "refactor"
mainSession: false
---

# 3: Update prompt pipeline

## Description
In `prompt.go`, update `typeToTemplate` map to use new type names, update `genScriptBases` with new ID prefixes. In `pkg/prompt/data/`, rename 17 template files to match new type names, add `validation-code.md` and `validation-ux.md` prompt templates.

### File renames (Part E)
| Old filename | New filename |
|---|---|
| `feature.md` | `coding-feature.md` |
| `enhancement.md` | `coding-enhancement.md` |
| `cleanup.md` | `coding-cleanup.md` |
| `refactor.md` | `coding-refactor.md` |
| `fix.md` | `coding-fix.md` |
| `documentation.md` | `doc.md` |
| `doc-evaluation.md` | `doc-eval.md` |
| `doc-generation-summary.md` | `doc-summary.md` |
| `doc-generation-consolidate.md` | `doc-consolidate.md` |
| `doc-generation-drift.md` | `doc-drift.md` |
| `test-pipeline-gen-cases.md` | `test-gen-cases.md` |
| `test-pipeline-eval-cases.md` | `test-eval-cases.md` |
| `test-pipeline-gen-scripts.md` | `test-gen-scripts.md` |
| `test-pipeline-run.md` | `test-run.md` |
| `test-pipeline-gen-and-run.md` | `test-gen-and-run.md` |
| `test-pipeline-graduate.md` | `test-graduate.md` |
| `test-pipeline-verify-regression.md` | `test-verify-regression.md` |

Unchanged: `gate.md`, `fix-record-missed.md`, `clean-code.md`.

New: `validation-code.md`, `validation-ux.md`.

## Reference Files
- `docs/proposals/task-type-id-redesign/proposal.md` — Source proposal
- `forge-cli/pkg/prompt/prompt.go` — typeToTemplate, genScriptBases
- `forge-cli/pkg/prompt/data/` — Template files

## Acceptance Criteria
- [ ] `prompt.Synthesize()` resolves new type names to correct template files
- [ ] All 17 template files renamed per mapping
- [ ] `validation-code.md` and `validation-ux.md` prompt templates created
- [ ] `genScriptBases` updated to new ID prefixes

## Hard Rules
- Template filename must exactly match the type name with `.` replaced by `-` (e.g., `coding.feature` → `coding-feature.md`)

## Implementation Notes
- `typeToTemplate` map key changes from old type name to new type name, value changes from old filename to new filename
- Use `git mv` for renames to preserve history
