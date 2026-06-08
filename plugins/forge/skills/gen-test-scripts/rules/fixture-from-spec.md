---
name: fixture-from-spec
description: Generate rich fixture data from Contract fixture_spec declarations, handling entity relationships and field constraints
---

# Fixture From Specification Rule

Generated test code MUST create fixture data that satisfies the `fixture_spec` declared in the Contract's Preconditions. The fixture_spec is the authoritative source for test data requirements — test code reads and fulfills these declarations rather than guessing data needs.

## Source: Contract fixture_spec

Each Contract's Preconditions section contains a `fixture_spec` field with this schema:

```yaml
fixture_spec:
  entities:                          # Required, at least 1 entity declaration
    - entity_type: string            # Required, entity type name (e.g., "Project", "Milestone")
      min_count: integer             # Required, minimum creation count (>=1)
      relationship_type: string      # Optional, relationship to parent ("belongs_to", "has_many")
      parent_entity: string          # Optional, parent entity type name (required when relationship_type is set)
      field_constraints:             # Optional, specific field value constraints
        - field: string              # Field name
          value: any                 # Expected value or constraint description
  state_requirements:                # Optional, prerequisite system state
    - description: string            # State description
      prerequisite_entity: string    # Dependent entity type
```

## Fixture Generation Rules

### R1: Entity Count Rule

When `fixture_spec.entities` declares `min_count: N` for an entity type, the generated test code MUST create >= N instances of that entity type.

**Violation example**: `min_count: 3` but test only creates 1 entity.
**Correct behavior**: Generate a loop or repeated create calls to produce at least N entities, each with distinct identifiable values (e.g., `"milestone-1"`, `"milestone-2"`, `"milestone-3"`).

### R2: Relationship Rule

When an entity declares `relationship_type` and `parent_entity`:

| relationship_type | Fixture generation behavior |
|-------------------|---------------------------|
| `belongs_to` | Create the parent entity first. Set the child entity's foreign key field to reference the parent's ID. Example: `milestone.map_id = created_map.id` |
| `has_many` | Create the parent entity first. Create >= `min_count` child entities, each referencing the parent's ID. The parent entity MUST have >= `min_count` associated children after fixture setup. |

**Ordering constraint**: Parent entities MUST be created before child entities in the fixture setup code. The generation order follows the dependency graph implied by `relationship_type` declarations.

**Circular dependency guard**: If the entity dependency graph contains cycles (entity A belongs_to B, B belongs_to A), the generation MUST abort with an error: `"Circular fixture dependency detected: A -> B -> A. Review Contract fixture_spec."`

### R3: Field Constraints Rule

When `field_constraints` declares field values, the generated fixture data MUST use those values when creating entities.

```yaml
field_constraints:
  - field: "status"
    value: "pending"
```

Generates:
```go
payload := map[string]interface{}{
    "name":   "milestone-1",
    "status": "pending",  // from field_constraints
}
```

**Default value generation**: For fields NOT listed in `field_constraints`, generate sensible deterministic values:
- `name`/`title` fields: `"{entity_type_lower}-{index}"` (e.g., `"milestone-1"`, `"milestone-2"`)
- `description` fields: `"Test {entity_type} description"`
- Numeric fields: Sequential integers starting from 1
- Boolean fields: `true` (positive default)
- Date/time fields: Fixed deterministic timestamp (e.g., `"2025-01-15T10:00:00Z"`)

### R4: State Requirements Rule

When `state_requirements` declares prerequisite system state:

1. Read the `description` to understand the required state
2. Read the `prerequisite_entity` to identify which entity must exist in a specific state
3. Generate fixture setup code that establishes the required state BEFORE the test action executes
4. Assert that the prerequisite state was established (verification assertion)

### R5: Multi-Step Fixture Passing

When a Journey has multiple steps, fixture data created in earlier steps (via API responses) MUST be passed to subsequent steps via test variables. Do NOT recreate entities that earlier steps have already created.

```go
// Step 1: Create map
mapResp := createMap(t, mapPayload)
createdMapID := mapResp.Data.ID

// Step 2: Create milestones referencing the map
milestonePayload["map_id"] = createdMapID  // Pass fixture from step 1
milestoneResp := createMilestone(t, milestonePayload)
```

## Fixture Consumption in Test Code

### Reading fixture_spec from Contract

The agent generating test code MUST:

1. Parse the Contract file's Preconditions section
2. Extract the `fixture_spec` field (YAML parsing)
3. For each entity in `fixture_spec.entities`, generate setup code according to R1-R5
4. If `fixture_spec` is absent, apply the backward compatibility handling defined in `types/_shared.md`

### Integration with Assertion Depth Rule

Fixture data created from `fixture_spec` enables deep assertions (see `rules/assertion-depth.md`):

- Entity relationships declared in `fixture_spec` (`relationship_type`, `parent_entity`) provide the targets for cross-entity assertions
- `min_count` declarations enable count-based assertions (e.g., `assert len(milestones) >= 3`)
- `field_constraints` enable state verification assertions (e.g., `assert milestone.status == "pending"`)

When supplementing assertions to meet the 80% behavioral threshold, prioritize assertions that verify fixture data was correctly established.

## Relationship to Other Rules

| Rule | Relationship |
|------|-------------|
| `rules/assertion-depth.md` | fixture_spec provides the entity data that enables deep assertions. Without fixture_spec, assertion depth is limited to shallow field matching. |
| `types/_shared.md` | `_shared.md` defines backward compatibility handling when fixture_spec is absent. This rule defines the forward behavior when fixture_spec is present. |
| `rules/quality-gates.md` | Fixture generation failures (e.g., API returns error during entity creation) are handled by the error handling table in quality-gates.md. |

## Examples

### Minimal fixture_spec (Single Entity CRUD)

Contract precondition:
```yaml
fixture_spec:
  entities:
    - entity_type: "Project"
      min_count: 1
```

Generated fixture setup:
```go
func setupProjectFixture(t *testing.T, baseURL string) map[string]interface{} {
    payload := map[string]interface{}{
        "name": "project-1",
    }
    resp := createEntity(t, baseURL, "/projects", payload)
    return resp
}
```

### Rich fixture_spec (Parent-Child with Constraints)

Contract precondition:
```yaml
fixture_spec:
  entities:
    - entity_type: "Map"
      min_count: 1
    - entity_type: "Milestone"
      min_count: 3
      relationship_type: "belongs_to"
      parent_entity: "Map"
      field_constraints:
        - field: "status"
          value: "pending"
```

Generated fixture setup:
```go
func setupMilestoneMapFixture(t *testing.T, baseURL string) (map[string]interface{}, []map[string]interface{}) {
    // Parent entity first (R2: ordering constraint)
    mapPayload := map[string]interface{}{
        "name": "map-1",
    }
    mapResp := createEntity(t, baseURL, "/maps", mapPayload)
    createdMapID := mapResp["id"].(string)

    // Child entities with relationship and field constraints (R1, R2, R3)
    var milestones []map[string]interface{}
    for i := 1; i <= 3; i++ {
        milestonePayload := map[string]interface{}{
            "name":    fmt.Sprintf("milestone-%d", i),
            "map_id":  createdMapID,   // R2: belongs_to relationship
            "status":  "pending",       // R3: field constraint
        }
        resp := createEntity(t, baseURL, "/milestones", milestonePayload)
        milestones = append(milestones, resp)
    }

    return mapResp, milestones
}
```
