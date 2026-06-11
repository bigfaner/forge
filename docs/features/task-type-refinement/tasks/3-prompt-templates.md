---
id: "3"
title: "Add prompt templates for feature, enhancement, cleanup, and refactor types"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 3: Add prompt templates for feature, enhancement, cleanup, and refactor types

## Description

Create four new prompt templates that replace `implementation.md` with type-specific execution strategies. Each template defines the appropriate workflow for its type: feature/enhancement follow TDD, cleanup follows improvement-then-verify, refactor follows behavior-preservation-verify.

## Reference Files
- `docs/proposals/task-type-refinement/proposal.md` — Source proposal (D3: Execution strategy)
- `forge-cli/pkg/prompt/data/implementation.md` — Current template (to deprecate)
- `forge-cli/pkg/prompt/data/fix.md` — Reference for fix workflow structure
- `forge-cli/pkg/prompt/prompt.go` — `typeToTemplate` map (line 22-38)

## Acceptance Criteria
- [ ] `data/feature.md` created: implement functionality → quality gate
- [ ] `data/enhancement.md` created: enhance existing behavior → quality gate
- [ ] `data/cleanup.md` created: improve technical debt → quality gate (no TDD requirement)
- [ ] `data/refactor.md` created: restructure code → quality gate (behavior preservation check)
- [ ] `typeToTemplate` map updated with 4 new entries mapping to new templates
- [ ] `implementation.md` kept but marked deprecated (header comment)
- [ ] All templates support standard variables: `{{TASK_ID}}`, `{{TASK_FILE}}`, `{{SCOPE}}`, `{{FEATURE_SLUG}}`

## Hard Rules
- Each template must be self-contained — no shared includes or imports between templates.
- Templates must follow the same variable substitution pattern as existing templates (handled by `renderTemplate()`).

## Implementation Notes
- `feature.md` and `enhancement.md` can be similar to `implementation.md` (TDD workflow) but with type-specific language.
- `cleanup.md` should NOT require writing failing tests first — its workflow is: read task → make improvements → run quality gate.
- `refactor.md` should emphasize: make structural changes → verify all existing tests still pass (behavior unchanged).
- The proposal D3 table specifies the key steps for each type.
