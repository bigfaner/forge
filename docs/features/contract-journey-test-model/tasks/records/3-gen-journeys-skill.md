---
status: "completed"
started: "2026-05-18 00:48"
completed: "2026-05-18 01:06"
time_spent: "~18m"
---

# Task Record: 3 gen-journeys skill

## Summary
Implemented gen-journeys skill: a Forge skill that extracts Journey narratives from PRD user stories. Created SKILL.md with 7-step workflow (read PRD -> identify workflows -> classify risk -> generate per-Journey files -> validate -> generate index -> commit) and journey.md template with structured format parseable by gen-contracts (Journey name, Risk level High/Medium/Low, Happy Path steps, Edge Cases, Setup preconditions, Journey Invariants). Skill follows all hard rules: no code reconnaissance, single Journey per generation with auto-batch, output structured for gen-contracts consumption. Added 15 e2e tests validating skill structure, template format, path conventions, and acceptance criteria coverage.

## Changes

### Files Created
- plugins/forge/skills/gen-journeys/SKILL.md
- plugins/forge/skills/gen-journeys/templates/journey.md
- tests/e2e/gen_journeys_skill_cli_test.go

### Files Modified
无

### Key Decisions
- Output path: docs/features/<slug>/testing/journeys/<journey-name>.md -- per-Journey files, not per-interface-type
- Journey template includes Setup preconditions section for Journey isolation support
- Risk classification inferred from PRD content: state mutation keywords -> High, multi-step interaction -> Medium, read-only -> Low
- Batch trigger: estimated content >50k tokens or step count >15, split within Journey (happy path batch 1, edge cases batch 2+)
- High-risk Journeys enforce edge case count >= happy path step count via HARD-RULE in SKILL.md

## Test Results
- **Tests Executed**: Yes
- **Passed**: 15
- **Failed**: 0
- **Coverage**: 88.5%

## Acceptance Criteria
- [x] gen-journeys from PRD user stories -> Journey Markdown output
- [x] Each Journey has: name, risk level (High/Medium/Low), >= 1 happy path step, >= 1 edge case step
- [x] Output per-Journey files (one user workflow per file), not per interface type
- [x] Format parseable by gen-contracts (Journey name + Step sequence + user action + expected result)
- [x] High-risk Journeys have edge case count >= happy path step count

## Notes
Task scope is backend (skill = Markdown documents, not executable code). Tests verify skill structure, template format, and content conformance to model-and-directory-spec.md. All 20 backend packages pass, 0 lint issues, e2e compile passes.
