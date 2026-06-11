---
name: assertion-depth
description: Assertion classification criteria, behavioral threshold (>=80%), and deep assertion requirement (>=30% of behavioral)
---

# Assertion Depth Rule

Generated tests MUST satisfy two quantitative thresholds for assertion quality:

1. **Behavioral threshold**: >=80% of all assertions must be behavioral assertions (not structural)
2. **Deep assertion requirement**: >=30% of behavioral assertions must be deep assertions

This rule is enforced at generation time. The agent generating test code MUST count and classify assertions, and automatically supplement if either threshold is not met.

## Assertion Classification Criteria

### Behavioral Assertions

An assertion is **behavioral** when it verifies business semantics — entity existence, correct state, relationship completeness, business rule satisfaction, or observable side effects.

| Category | Definition | Example |
|----------|-----------|---------|
| Entity existence | Verifies a specific entity was created/exists | `assert milestone.map_id == map.id` |
| State correctness | Verifies entity state matches expected value | `assert response.data.status == "completed"` |
| Relationship completeness | Verifies entity relationships are correct | `assert project.owner_id == user.id` |
| Business rule satisfaction | Verifies business rules hold | `assert created_count == 3` |
| Observable side effect | Verifies a side effect occurred | `assert list.length > 0` (proves data creation) |

### Structural Assertions

An assertion is **structural** when it verifies transport/format layer correctness — HTTP status codes, response schema, field types.

| Category | Definition | Example |
|----------|-----------|---------|
| HTTP status code | Verifies response status | `assert response.status == 200` |
| Response schema | Verifies response structure | `assert typeof response.data.id == "string"` |
| Field existence | Verifies field is present (without value check) | `assert response.body contains "id"` |
| Content type | Verifies response format | `assert response.headers["content-type"] == "application/json"` |

### Edge Cases

| Assertion | Classification | Rationale |
|-----------|---------------|-----------|
| `assert response.status == 201 AND response.data.name == "milestone-1"` | Behavioral | Mixed assertion containing at least one business field verification counts as behavioral |
| `assert response.body contains "id"` | Structural | Only verifies field existence without value check |
| `assert list.length > 0` | Behavioral | Verifies non-empty collection, proving data creation succeeded |
| Health check / readiness endpoint test | Structural | Legitimate structural assertion, exempt from 80% threshold |

### Classification Rule for Mixed Assertions

When a single assertion statement contains multiple conditions connected by AND/OR:
- If **any** condition verifies a business field value (not just existence/type), the entire assertion is classified as **behavioral**
- If **all** conditions are structural (status code, schema, field existence, type), the assertion is classified as **structural**

## Behavioral Threshold (>=80%)

### Calculation

```
behavioral_ratio = count(behavioral_assertions) / count(all_assertions)
```

The ratio MUST be >= 0.80 (80%).

The 20% allowance for structural assertions covers legitimate scenarios:
- Health check / readiness endpoint tests
- API discovery / schema validation tests
- Response format conformance tests
- Smoke tests that verify connectivity before deeper assertions

### Enforcement

At generation time, after writing all test functions for a Journey:

1. Count all assertions across all generated test functions for the Journey
2. Classify each as behavioral or structural using the criteria above
3. Calculate `behavioral_ratio`
4. If `behavioral_ratio < 0.80`:
   - Identify test functions with the fewest behavioral assertions
   - Supplement by adding behavioral assertions derived from Contract Outcomes (State, Side-effect, Invariants dimensions)
   - Re-count and re-classify until threshold is met
5. Output a summary: `"Assertion depth: {behavioral_count}/{total_count} behavioral ({ratio}%), {deep_count}/{behavioral_count} deep ({deep_ratio}%)"`

### Per-Function Minimum

Each test function MUST contain at least 1 behavioral assertion. Test functions with only structural assertions are vacuous and violate the Vacuous Assertions antipattern guard (see `quality-gates.md`).

## Deep Assertion Requirement (>=30%)

### Definition

A behavioral assertion is a **deep assertion** when it verifies:
- **Entity relationships**: assertions that cross entity boundaries (e.g., `assert milestone.map_id == map.id`, `assert order.items.length == 3`)
- **State transitions**: assertions that verify an entity changed from one state to another (e.g., `assert response.data.status == "completed"` when the previous state was "pending")

Shallow behavioral assertions verify single-entity field value matches (e.g., `assert name == input`, `assert response.data.title == "expected"`). These are behavioral but not deep.

### Calculation

```
deep_ratio = count(deep_assertions) / count(behavioral_assertions)
```

The ratio MUST be >= 0.30 (30% of behavioral assertions must be deep).

### Enforcement

At generation time, after satisfying the behavioral threshold:

1. Among behavioral assertions, classify each as deep or shallow
2. Calculate `deep_ratio`
3. If `deep_ratio < 0.30`:
   - Identify entity relationships from Contract Preconditions and fixture_spec
   - Add assertions that verify cross-entity relationships (e.g., child entity references parent ID, parent entity contains expected child count)
   - Add assertions that verify state transitions (compare initial state to final state)
   - Re-count until threshold is met
4. Include deep assertion count in the summary output

## Relationship to Other Rules

| Rule | Relationship |
|------|-------------|
| `quality-gates.md` — Antipattern Guard | Assertion depth complements the Vacuous Assertions antipattern guard. Vacuous Assertions catches zero-assertion functions; this rule catches insufficient behavioral depth. |
| `quality-gates.md` — Error Handling | Assertion depth enforcement failure (cannot reach threshold) is a generation-time error, not a compile error. If threshold cannot be met (e.g., health check only Journey), document the exception in the test file header. |
| `fixture-from-spec.md` | Rich fixture data from fixture_spec provides the entities needed for deep assertions (entity relationships, state transitions). Without fixture_spec, deep assertion opportunities are limited. |

## Examples

### Journey: Milestone Map (Parent-Child Entity)

```go
// Behavioral assertion — entity existence
assert.Equal(t, map.ID, createdMap.Data.ID)

// Behavioral assertion — relationship (DEEP)
assert.Equal(t, createdMap.Data.ID, milestones[0].MapID)

// Behavioral assertion — state transition (DEEP)
assert.Equal(t, "completed", updatedMilestone.Data.Status)

// Behavioral assertion — business rule (shallow)
assert.Equal(t, 3, len(milestones))

// Structural assertion (allowed within 20%)
assert.Equal(t, 201, resp.StatusCode)
```

Summary: 5 assertions, 4 behavioral (80%), 2 deep (50% of behavioral). Both thresholds met.

### Journey: Health Check (Exception)

```go
// Structural only — this Journey is exempt
assert.Equal(t, 200, resp.StatusCode)
assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
```

Summary: 0 behavioral assertions. This Journey tests infrastructure readiness, not business behavior. Exempt from threshold — document in test file header: `// ASSERTION_DEPTH_EXEMPT: health check / readiness Journey`.
