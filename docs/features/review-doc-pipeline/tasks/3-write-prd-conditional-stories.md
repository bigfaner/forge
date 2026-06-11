---
id: "3"
title: "Add conditional user stories to write-prd Step 7"
priority: "P1"
estimated_time: "30min"
dependencies: []
type: "doc"
mainSession: false
---

# 3: Add conditional user stories to write-prd Step 7

## Description

Modify `write-prd/SKILL.md` Step 7 so it only generates `prd-user-stories.md` when the feature involves code. For doc-only features, user stories have no downstream consumer (gen-journeys → test scripts only serve testable code), so generating them adds noise.

## Reference Files
- `docs/proposals/review-doc-pipeline/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/write-prd/SKILL.md` | Add conditional logic to Step 7 |

## Acceptance Criteria

- [ ] Step 7 includes a condition check: if In Scope contains only non-compilable/non-runnable file paths (`.md`, `.yaml` under `docs/`, `skills/`, etc.), skip user story generation
- [ ] Step 7 generates `prd-user-stories.md` normally when In Scope contains any compilable or runnable file paths (`.go`, `.ts`, `.py`, etc.)
- [ ] The condition logic is consistent with task type assignment rules (same heuristic as classifying `doc` vs `coding.*` types)
- [ ] When skipping, Step 7 outputs a brief note explaining why user stories are skipped for doc-only features

## Hard Rules

- The determination must be based on file paths/artifacts in In Scope, not on the user's stated intent
- Do not remove or restructure existing Step 7 content — only add the conditional gate at the beginning

## Implementation Notes

The condition should mirror the Type Assignment logic from quick-tasks/breakdown-tasks:
- If all In Scope items target non-compilable, non-runnable artifacts → doc-only → skip Step 7
- If any In Scope item involves compilable/runnable files → generate user stories

Add a brief preamble to Step 7 like:
"**Gate**: If all In Scope items are non-compilable artifacts (markdown, specs, templates), skip this step and note that user stories are not needed for doc-only features. User stories serve gen-journeys → test script generation, which requires testable code."
