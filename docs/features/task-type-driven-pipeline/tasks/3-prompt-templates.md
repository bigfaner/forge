---
id: "3"
title: "Create documentation and doc-evaluation prompt templates"
priority: "P1"
estimated_time: "2h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 3: Create documentation and doc-evaluation prompt templates

## Description
Create two new prompt templates for the typed-task-dispatch system. The `documentation.md` template guides agents through document creation/modification with a self-check step. The `doc-evaluation.md` template implements a 1000-point rubric (8 dimensions x 125 points) with a 3-round score-revise-rescore iteration cycle for docs-only features.

## Reference Files
- `docs/proposals/task-type-driven-pipeline/proposal.md` — Source proposal (D2, D3)
- `forge-cli/pkg/prompt/data/implementation.md` — Existing template for structural reference

## Affected Files

### Create
| File | Description |
|------|-------------|
| `forge-cli/pkg/prompt/data/documentation.md` | Prompt template for `documentation` type tasks |
| `forge-cli/pkg/prompt/data/doc-evaluation.md` | Prompt template for `doc-evaluation` type tasks (T-eval-doc) |

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/prompt/prompt.go` | Add `TypeDocumentation` and `TypeDocEvaluation` entries to `typeToTemplate` map |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] `typeToTemplate` maps `TypeDocumentation` → `"data/documentation.md"` and `TypeDocEvaluation` → `"data/doc-evaluation.md"`
- [ ] `documentation.md` template follows 4-step structure: (1) read task description + refs, (2) execute document work, (3) self-check (format, cross-references, terminology consistency), (4) submit via `forge task submit`
- [ ] `documentation.md` template uses standard placeholders: `{{TASK_ID}}`, `{{TASK_FILE}}`, `{{FEATURE_SLUG}}`
- [ ] `doc-evaluation.md` template implements 1000-point rubric with 8 dimensions x 125 points each: structural completeness, logical consistency, traceability, accuracy, completeness, terminology consistency, formatting standards, language quality
- [ ] `doc-evaluation.md` implements iteration cycle: score → if < 900 and round < 3, revise and re-score → final score reported
- [ ] `doc-evaluation.md` uses `{{TASK_FILE}}` to read T-eval-doc task which contains the list of documents to evaluate
- [ ] Existing tests in `prompt_test.go` pass; new test covers `Synthesize` for both new types returning non-empty result

## Implementation Notes
- Follow the existing `implementation.md` template for placeholder patterns and section structure. Key placeholders: `{{TASK_ID}}`, `{{TASK_FILE}}`, `{{SCOPE}}`, `{{FEATURE_SLUG}}`.
- The `doc-evaluation.md` rubric is hardcoded in the template (not externalized). Future rubric changes require editing this template file.
- For `doc-evaluation.md`, the agent needs to: (1) read all documentation files produced by the feature's business tasks, (2) evaluate each against the 8-dimension rubric, (3) produce a scored report, (4) if below 900, revise documents and re-evaluate up to 2 more times, (5) submit final report via `forge task submit`.
- The doc-evaluation task's dependencies (set in task 2) point to the last business task, ensuring all documents exist before evaluation.
