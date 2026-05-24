---
id: "2"
title: "Strengthen templates with Reference Files declaration + AC validation"
priority: "P0"
estimated_time: "1h"
dependencies: [1]
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 2: Strengthen templates with Reference Files declaration + AC validation

## Description

For each template identified by Task 1's audit as needing strengthening, insert two `<IMPORTANT>` blocks:
1. **Reference Files authority declaration** in Step 1 (after reading task file, before existing Hard Rules block)
2. **AC per-item validation** in the Verify/Self-Check step (as first sub-step before static checks)

After all template modifications, run `go build` to re-embed templates into the binary and verify the build succeeds.

## Reference Files
- `docs/proposals/spec-authority-enforcement/proposal.md#Proposed-Solution` — Exact `<IMPORTANT>` template text for Reference Files declaration and AC validation
- `docs/proposals/spec-authority-enforcement/proposal.md#Priority-Rules` — Conflict priority: Hard Rules > Reference Files > existing code
- `docs/proposals/spec-authority-enforcement/proposal.md#Edge-Cases-&-Degradation` — Edge cases: empty Reference Files, missing sections, file not found
- `docs/conventions/forge-distribution.md` — Forge distribution model, embed.FS compilation requirement

## Acceptance Criteria
- [ ] Each template identified by Task 1 has a `<IMPORTANT>` Reference Files authority declaration inserted in Step 1 after "read task file" and before the existing `<IMPORTANT>` Hard Rules block
- [ ] Declaration text matches the proposal's exact template (4 MUST items + 2 fallback outputs)
- [ ] Each template has AC per-item validation inserted in its Verify/Self-Check step as the first sub-step (before static checks)
- [ ] AC validation text matches the proposal's exact template (per-item PASS/FAIL + skip condition)
- [ ] Use `<IMPORTANT>` tag (not `<EXTREMELY-IMPORTANT>`) to avoid marker dilution with existing blocks
- [ ] `go build ./...` succeeds after all modifications
- [ ] Existing `<IMPORTANT>` Hard Rules blocks in templates are NOT modified or removed

## Hard Rules
- MUST use `<IMPORTANT>` tag for new declarations, NOT `<EXTREMELY-IMPORTANT>` — avoid diluting existing EXTREMELY-IMPORTANT markers
- MUST NOT modify existing `<IMPORTANT>` Hard Rules blocks — only insert new blocks
- MUST NOT change the overall workflow structure (step count, step names) of any template
- MUST load `docs/conventions/forge-distribution.md` before modifying any template files

## Implementation Notes

### Insertion Points (per template type)

**coding-feature.md / coding-enhancement.md** (3-step workflow):
- Reference Files declaration: between line 28 (`Output: Step 1/3...`) and line 30 (`<IMPORTANT>` Hard Rules)
- AC validation: at the beginning of Step 3 "Static Checks + Targeted Tests", before the static checks code block

**coding-refactor.md** (4-step workflow):
- Reference Files declaration: between line 37 (`Output: Step 1/4...`) and line 39 (`<IMPORTANT>` Hard Rules)
- AC validation: at the beginning of Step 4 "Static Checks + Targeted Tests", before the static checks code block

**coding-fix.md** (4-step workflow):
- Reference Files declaration: between line 32 (`Output: Step 1/4...`) and line 34 (`<IMPORTANT>` Hard Rules)
- AC validation: at the beginning of Step 4 "Static Checks + Targeted Tests", before the static checks code block

**coding-cleanup.md** (3-step workflow):
- Reference Files declaration: between line 26 (`Output: Step 1/3...`) and line 28 (`<IMPORTANT>` Hard Rules)
- AC validation: at the beginning of Step 3 "Static Checks + Targeted Tests", before the static checks code block

**Other templates identified by Task 1**: Follow the same pattern — insert Reference Files declaration after Step 1 output line, insert AC validation at the beginning of the verify step.

### Build Verification

After all modifications, run `go build ./...` from `forge-cli/` to verify templates are correctly embedded. The templates are embedded via `embed.FS` in the Go binary — a successful build confirms correct embedding.

### SYNC NOTICE

coding-enhancement.md has a sync notice: "This template shares ~90% structure with coding-feature.md. When modifying this file, review coding-feature.md for equivalent changes." — ensure both receive identical modifications.
