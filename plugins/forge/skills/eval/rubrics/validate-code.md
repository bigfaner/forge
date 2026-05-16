---
scale: 1000
target: 700
iterations: 1
type: validate-code
context:
  conventions: []
  business-rules: auto
---

# validate-code Evaluation Rubric

**Total: 1000 points | iterations: 1 (single-pass, no revise loop)**

## Purpose

This rubric drives static code tracing: for each PRD user scenario, trace through git diff and implementation code to verify a complete implementation path exists. The output is a problem report, not a revised document.

## Required Input

| Input | Source |
|-------|--------|
| PRD user scenarios | `prd/prd-spec.md` + `prd/prd-user-stories.md` |
| Git diff | `git diff <base-branch>...HEAD` |
| Changed file list | From git diff |
| Implementation code | Files referenced by diff hunks |

## Scoring Dimensions

| Dimension | Points |
|-----------|--------|
| 1. Scenario Traceability | 400 |
| 2. Path Completeness | 350 |
| 3. Code-PRD Consistency | 250 |
| **Total** | **1000** |

## Dimensions

### 1. Scenario Traceability (400 pts)

For each PRD user scenario, verify that a traceable path exists from the scenario description to concrete code changes.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Scenario-to-diff mapping | 0-200 | For each PRD scenario, identify which git diff hunks relate to it. Score based on coverage: what fraction of scenarios have at least one mapped diff hunk? 100% coverage = full score. |
| Trace chain clarity | 0-120 | Each mapped scenario has a clear chain: PRD scenario -> acceptance criteria -> implementation code path. Ambiguous or broken chains lose points. |
| Unmapped scenarios identified | 0-80 | Scenarios with no diff mapping are explicitly listed (not silently skipped). The report must state which scenarios have no implementation evidence. |

### 2. Path Completeness (350 pts)

For each scenario that has diff mapping, assess whether the implementation path is complete.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Happy path completeness | 0-150 | For each scenario, does the implementation cover the full happy path from trigger to end state? Score by fraction of scenarios with complete happy paths. |
| Error handling coverage | 0-100 | For scenarios that describe error/failure cases in the PRD, does the code implement corresponding error handling? Missing error paths lose points proportionally. |
| Edge case handling | 0-100 | Boundary conditions mentioned in PRD acceptance criteria (empty input, max values, concurrent access) have corresponding code guards. |

### 3. Code-PRD Consistency (250 pts)

Verify that implementation code does not contradict or deviate from PRD specifications.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Behavioral alignment | 0-100 | Implementation behavior matches PRD-described behavior. Divergences (different flow, different output, different side effects) lose points. |
| Interface contract fidelity | 0-80 | Function signatures, CLI flags, API endpoints, and data shapes match what the PRD specifies. Naming deviations that change semantics lose points. |
| Scope compliance | 0-70 | Implementation does not silently include out-of-scope behavior or omit in-scope requirements. Check against PRD Scope section. |

## Report Format

The scorer produces a structured problem report:

```markdown
# validate-code Report

## Summary
- Scenarios traced: X/Y
- Paths complete: X/Y
- Issues found: N

## Per-Scenario Results

### Scenario: <scenario-name>
- **Status**: pass / partial / fail
- **Trace chain**: <PRD section> -> <diff hunk> -> <code file:function>
- **Happy path**: complete / incomplete (missing: ...)
- **Error handling**: covered / missing (unhandled: ...)
- **Issues**: (list if any)

### Scenario: <scenario-name>
...

## Unmapped Scenarios
(list scenarios with no diff mapping)

## Consistency Issues
(list code-PRD contradictions found)
```

## Scoring Guide

| Score Range | Interpretation |
|-------------|----------------|
| 900-1000 | All scenarios fully traced, complete paths, no inconsistencies |
| 700-899 | Most scenarios traced, minor gaps in error handling or edge cases |
| 500-699 | Significant gaps: multiple scenarios unmapped or partially implemented |
| Below 500 | Major implementation gaps: majority of scenarios lack code evidence |

## Deduction Rules

- **Scenario not attempted in trace**: -50 per scenario
- **False positive trace** (claims mapping but code is unrelated): -30 per instance
- **Ignored PRD acceptance criterion**: -20 per criterion
