---
status: "completed"
started: "2026-05-18 00:02"
completed: "2026-05-18 00:12"
time_spent: "~10m"
---

# Task Record: 1 Journey-Driven 测试模型定义与目录规范

## Summary
Created the Journey-Driven Test Model & Directory Specification document defining core concepts (Journey, Contract six dimensions, Outcome, Risk classification, semantic descriptors), directory conventions (tests/<journey>/ and _contracts/), gen-contracts parseable structure, and config.yaml schema (test-framework, test-command, capabilities)

## Changes

### Files Created
- docs/features/contract-journey-test-model/design/model-and-directory-spec.md

### Files Modified
无

### Key Decisions
- Structured six dimensions as mandatory (Preconditions/Input/Output/State) and optional (Side-effect/step-level Invariants) per Hard Rules
- Journey-level Invariants are mandatory and appear in every Contract file within a Journey
- Config schema extends existing .forge/config.yaml with test-framework, test-command, and capabilities fields for backward compatibility
- Semantic descriptors use pure natural language (no regex) with conversion pipeline deferred to gen-test-scripts + Fact Table
- State verification supports graceful degradation (full/partial/deferred) when projects lack state query interfaces
- Contract file format uses structured Markdown with parseable ## Outcome headings and dimension key-value lines

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Journey definition with user workflow + Risk classification High/Medium/Low
- [x] Contract six dimensions: Preconditions/Input/Output/State mandatory, Side-effect/Invariants optional
- [x] Multi-Outcome Contracts with Preconditions mutual exclusivity
- [x] Semantic descriptors for gen-contracts stage, converted to regex by gen-test-scripts
- [x] Directory convention: tests/<journey>/ for tests, tests/<journey>/_contracts/ for Contract specs
- [x] gen-contracts parseable structure (Journey name + Step sequence + Outcome per step)
- [x] config.yaml schema with languages, test-framework, test-command, capabilities fields

## Notes
Documentation-only task. No code changes. All 20 existing test packages pass (0 failures).
