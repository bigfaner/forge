# Confidence Rating Rules

## Overview

Confidence rating evaluates the trustworthiness of generated tests based on
how many Contract Outcomes are backed by confirmed runtime facts from the
Fact Table.

## Rating Thresholds

| Level  | Condition                          | Mark   |
|--------|------------------------------------|--------|
| HIGH   | confirmed_fact_ratio >= 0.80       | VERIFY |
| MEDIUM | 0.40 <= confirmed_fact_ratio < 0.80| VERIFY |
| LOW    | confirmed_fact_ratio < 0.40        | REVIEW |

## Forced Downgrade

Regardless of the computed ratio, the level is forced to **LOW** when:

- **eval-skipped**: the eval-journey or eval-contract step was skipped
  (e.g. no rubric configured, user chose to skip)
- **eval-bypassed**: the user explicitly bypassed the eval gate

Users may clear eval-skipped markers manually; after clearing, confidence
is recalculated from the Fact Table coverage ratio.

## Calculation

```
confirmed_fact_ratio = (outcomes covered by runtime+confirmed facts) / (total outcomes)
```

Where:
- "outcomes" are the unique subject strings derived from Contract Outcomes
- "runtime+confirmed facts" come from the Fact Table via `forge fact summary`
- An outcome is "covered" if at least one FactEntry with source=runtime and
  confidence=confirmed exists for that subject

## Test File Annotation

After computing the confidence rating, embed it in the test file header:

```go
// confidence: HIGH
// confirmed_fact_ratio: 0.85
// total_outcomes: 20
// confirmed_outcomes: 17
```

Or for non-Go test files, use the language's comment syntax.

## Report Integration

The test report must include:

### Confidence Distribution

| Level  | Count | Percentage |
|--------|-------|------------|
| HIGH   | N     | X%         |
| MEDIUM | N     | X%         |
| LOW    | N     | X%         |

### Verification Mark Summary

| Mark   | Count | Description              |
|--------|-------|--------------------------|
| VERIFY | N     | HIGH or MEDIUM confidence|
| REVIEW | N     | LOW confidence, needs human review |

## Behavior

- LOW confidence tests **still execute** normally
- LOW confidence tests are **not blocked or skipped**
- LOW confidence tests are marked as **REVIEW** in the report to flag for
  manual verification
- The confidence rating is informational, not a gate
