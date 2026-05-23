---
id: "3"
title: "Update 13 embed template files with strategy-based content"
priority: "P1"
estimated_time: "1h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 3: Update 13 embed template files with strategy-based content

## Description

Replace all 13 embed template files in `forge-cli/pkg/task/data/` with concise feature-context content based on the proposal's three-strategy model. The current templates contain full workflow instructions that duplicate what prompt templates already provide. Replace them with per-instance context that only the task file body can carry.

## Reference Files
- `docs/proposals/auto-task-main-instructions/proposal.md` — Template content for each strategy
- `forge-cli/pkg/task/data/*.md` — 13 template files to replace

## Acceptance Criteria

- [ ] All 13 template files replaced with strategy-based content from the proposal
- [ ] Strategy A templates (7 files) contain: feature slug, scope, interfaces placeholders; skill invocation hints
- [ ] Strategy B templates (2 files) contain: acceptance criteria placeholders; quality gate commands
- [ ] Strategy C templates (4 files) contain: discovery strategy steps (git diff, directory scan)
- [ ] No workflow duplication between template body and prompt templates — templates only provide feature-specific context
- [ ] Templates are concise (5-15 lines each, per proposal)

## Hard Rules

- MUST NOT include full workflow instructions (those belong in prompt templates)
- MUST use `{{PLACEHOLDER}}` syntax for dynamic content (resolved by renderBody)
- MUST preserve the strategy distinction: A provides context, B provides criteria + discovery, C provides discovery method only

## Implementation Notes

### Strategy mapping:

**Strategy A — Feature Context (7 files):**
- `test-gen-cases.md`: "Generate structured test cases for the {{FEATURE_SLUG}} feature ({{MODE}} mode)." + scope + interfaces
- `test-eval-cases.md`: "Evaluate generated test cases for executability for the {{FEATURE_SLUG}} feature." + scope
- `test-gen-scripts.md`: "Generate executable test scripts for the {{FEATURE_SLUG}} feature. Test type: {{TEST_TYPE}}." + read approved test cases
- `test-gen-and-run.md`: "Generate and run test scripts for the {{FEATURE_SLUG}} feature. Test type: {{TEST_TYPE}}." + Phase 1/2
- `test-run.md`: "Execute staged e2e test scripts for the {{FEATURE_SLUG}} feature." + scope
- `test-graduate.md`: "Promote feature test scripts to the project regression suite for the {{FEATURE_SLUG}} feature."
- `test-verify-regression.md`: "Run full e2e regression suite to verify no regressions after the {{FEATURE_SLUG}} feature." + scope + `just test-e2e`

**Strategy B — Static + Discovery (2 files):**
- `validation-code.md`: "Validate code quality for the {{FEATURE_SLUG}} feature." + `{{ACCEPTANCE_CRITERIA}}` + quality gate commands
- `validation-ux.md`: "Validate user experience for the {{FEATURE_SLUG}} feature." + `{{ACCEPTANCE_CRITERIA}}` + UI design check

**Strategy C — Discovery Strategy (4 files):**
- `doc-eval.md`: Scan directories for ALL documents, score with 8-dimension rubric
- `doc-consolidate.md`: Scan feature docs, extract rules/specs, compare against existing, auto-integrate
- `doc-drift.md`: `git diff --name-only main...HEAD` → match domains → narrow scope → auto-fix
- `code-quality-simplify.md`: `git diff --name-only main...HEAD` → focus on changed files only

### Template format (example for Strategy A):

```markdown
Generate structured test cases for the {{FEATURE_SLUG}} feature ({{MODE}} mode).

## Feature Context
- Scope: {{SCOPE}}
- Test interfaces: {{INTERFACES}}

Read the PRD/proposal to extract acceptance criteria and generate test cases with full traceability.
```
