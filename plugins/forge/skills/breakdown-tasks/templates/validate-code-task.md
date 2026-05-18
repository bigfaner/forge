---
id: "T-validate-code"
title: "validate-code: Static Code Tracing"
priority: "P0"
estimated_time: "1h"
dependencies: [{{ALL_IMPLEMENTATION_TASK_IDS}}]
breaking: false
type: "gate"
mainSession: false
---

# T-validate-code: validate-code: Static Code Tracing

## Description

Run static code tracing validation against the PRD. For each PRD user scenario, trace through git diff and implementation code to verify a complete implementation path exists. This produces a problem report, not a revised document.

Position: after all implementation tasks, before T-test-1 (gen-sitemap).

## Instructions

### Step 1: Run validate-code eval

```bash
forge eval --type validate-code
```

This command:
1. Reads the PRD to extract user scenarios
2. Runs `git diff <base-branch>...HEAD` to get changed files
3. Compiles the changed file list
4. Passes all to the scorer for scenario tracing

### Step 2: Review the report

Read the generated report at `docs/features/<slug>/eval/validate-code.md`.

- If all scenarios are "pass" or "partial" with minor gaps → proceed, note gaps in record
- If any scenario is "fail" → investigate and either fix or document as known limitation

## Reference Files

- `docs/features/<slug>/prd/prd-spec.md` — PRD with user scenarios
- `docs/features/<slug>/prd/prd-user-stories.md` — User stories with acceptance criteria
- `../../../eval/rubrics/validate-code.md` — Rubric definition

## Acceptance Criteria

- [ ] `forge eval --type validate-code` executed successfully
- [ ] Report generated at `docs/features/<slug>/eval/validate-code.md`
- [ ] Every PRD scenario appears in the report with a traceability status
- [ ] Any "fail" or "partial" scenarios are investigated and documented
- [ ] Record created via `/submit-task`

## Hard Rules

- MUST NOT modify implementation code — this is verification only
- MUST NOT skip scenarios — every PRD user scenario must appear in the trace report

## Implementation Notes

- The eval runs with `iterations: 1`, so there is no revise loop
- If the feature is docs-only (no runtime code), this task can be skipped — mark as completed with note "docs-only feature, no code to validate"
- The base branch for git diff is determined by the eval skill automatically
