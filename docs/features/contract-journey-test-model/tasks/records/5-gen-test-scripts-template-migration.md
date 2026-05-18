---
status: "completed"
started: "2026-05-18 01:31"
completed: "2026-05-18 01:58"
time_spent: "~27m"
---

# Task Record: 5 gen-test-scripts 重写 + 内置模板迁移

## Summary
Rewrote gen-test-scripts skill for Journey-Driven model: added descriptor package (semantic descriptor to regex conversion via Fact Table), journey package (Contract-based test code generation with @feature tags for Go/Python/JS), Journey smoke test templates for all 3 languages, updated SKILL.md to Contract/Journey workflow, added built-in Journey templates to language profiles, bumped version to 3.22.0.

## Changes

### Files Created
- forge-cli/pkg/descriptor/descriptor.go
- forge-cli/pkg/descriptor/descriptor_test.go
- forge-cli/pkg/journey/journey.go
- forge-cli/pkg/journey/journey_test.go
- forge-cli/pkg/profile/languages/go/templates/journey-test.go
- forge-cli/pkg/profile/languages/python/templates/journey_test.py
- forge-cli/pkg/profile/languages/javascript/templates/journey-smoke.spec.ts

### Files Modified
- plugins/forge/skills/gen-test-scripts/SKILL.md
- forge-cli/pkg/profile/embed_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- Semantic descriptor to regex uses keyword matching with hyphen-aware partial matching (min 2 keyword overlap threshold)
- Journey package dispatches to per-framework generators (Go/Python/JS) via GenerateDispatched for extensibility
- Tests go directly to tests/<journey>/ with @feature tags (no staging area), matching proposal tag-based lifecycle
- Sensitive field detection uses pattern matching (token/secret/password/key/credential) to enforce no-hardcoded-secrets rule
- Built-in Journey templates embedded in language profiles alongside existing templates, preserving backward compatibility

## Test Results
- **Tests Executed**: Yes
- **Passed**: 92
- **Failed**: 0
- **Coverage**: 84.5%

## Acceptance Criteria
- [x] gen-test-scripts converts semantic descriptors to precise regex via Fact Table, generates compilable test code
- [x] Generated tests carry @feature tags (Go //go:build feature, Python @pytest.mark.feature, JS describe('@feature'))
- [x] Tests go directly into tests/<journey>/ directory (no staging)
- [x] At least 1 Journey smoke test generated per Journey (happy path end-to-end)
- [x] Smoke test output matches Contract success Outcome Output/State declarations
- [x] Zero-config output uses built-in templates; custom template path override supported
- [x] 6 language profile generate.md/run.md work as built-in default templates

## Notes
Test data safety enforced: IsSensitiveField/BuildSensitiveFieldPlaceholder ensure no hardcoded secrets. Batch generation rule (one Journey at a time) enforced in SKILL.md HARD-RULE.
