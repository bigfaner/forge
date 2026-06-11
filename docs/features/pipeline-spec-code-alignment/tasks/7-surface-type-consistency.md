---
id: "7"
title: "Fix surface and type consistency in docs"
priority: "P1"
estimated_time: "1.5h"
dependencies: [6]
type: "doc"
mainSession: false
---

# 7: Fix surface and type consistency in docs

## Description

Fix surface resolution, type naming, and field consistency in documentation. Covers Cluster 2 doc changes (issues B1-B3, B5-B6):

1. **Two-layer surface resolution**: Update quick-tasks/SKILL.md, breakdown-tasks/SKILL.md, and scope-to-surface-key.md to document the two-layer resolution strategy: (a) project-level shortcut — single-surface projects skip `forge surfaces` calls entirely, (b) file-level query — multi-surface projects use path prefix first, `forge surfaces` only for ambiguous files.

2. **task-doc.md surface fields**: The doc task template (`quick-tasks/templates/task-doc.md`) lacks `surface-key`/`surface-type` fields present in `task.md`. Add these fields or add an exemption note explaining why doc tasks don't need them.

3. **webui → web**: `gen-test-scripts/rules/step-0.5-validation.md` (line 20) uses `webui` as surface type, but canonical name is `web` (per Go code `KnownSurfaceTypes` and gen-journeys).

4. **record-format-doc.md ghost types**: Remove non-existent `doc.eval` type and add valid `doc.review` type.

## Reference Files
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Problem` — Evidence B1-B3 (surface resolution), B5 (webui→web), B6 (doc.eval ghost type)
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Proposed-Solution` — Cluster 2 description for surface consistency
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Success-Criteria` — SC for zero forge surfaces calls in single-surface projects, web not webui, doc.review not doc.eval

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/quick-tasks/SKILL.md` | Two-layer surface resolution docs |
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Two-layer surface resolution docs |
| `plugins/forge/skills/breakdown-tasks/rules/scope-to-surface-key.md` | Two-layer resolution strategy |
| `plugins/forge/skills/quick-tasks/templates/task-doc.md` | Add surface fields or exemption note |
| `plugins/forge/skills/gen-test-scripts/rules/step-0.5-validation.md` | `webui` → `web` |
| `plugins/forge/skills/submit-task/data/record-format-doc.md` | Remove `doc.eval`, add `doc.review` |

## Acceptance Criteria
- [ ] Surface resolution docs describe two-layer strategy (project-level shortcut + file-level query)
- [ ] `task-doc.md` either has surface fields or documents why they're absent
- [ ] No reference to `webui` surface type — all use canonical `web`
- [ ] `record-format-doc.md` lists `doc.review` and does not list `doc.eval`
- [ ] Single-surface project surface-type is non-empty (not left blank as placeholder)

## Hard Rules
- Do not change Go code behavior — this task is doc-only

## Implementation Notes
- For task-doc.md: the simplest approach is adding a comment in the template explaining that doc tasks don't interact with the quality gate, so surface fields are unnecessary. Alternatively, add the fields with empty defaults.
