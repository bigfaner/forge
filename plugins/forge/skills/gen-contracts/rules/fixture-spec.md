# Fixture Specification Dimension Rules

Every Contract's Preconditions MUST include a `fixture_spec` field that declaratively specifies the pre-existing data state required before the Step executes. This specification is consumed by gen-test-scripts to generate rich fixtures, eliminating guesswork about what data a test needs.

## Why Fixture Specification

Without fixture_spec, gen-test-scripts has no authoritative source for what entities and relationships must exist before a test runs. This leads to empty-container testing: tests pass on CRUD correctness but miss business failures caused by missing entity relationships or insufficient data volume.

Fixture Specification is **declarative** — it declares "what data is needed", not "how to create it". The creation logic is the responsibility of gen-test-scripts' fixture consumption rules.

## Fixture Specification Schema

```yaml
fixture_spec:
  entities:                          # Required, at least 1 entity declaration
    - entity_type: string            # Required, entity type name (e.g., "Project", "Milestone")
      min_count: integer             # Required, minimum creation count (>= 1)
      relationship_type: string      # Optional, relationship to parent entity (e.g., "belongs_to", "has_many", "has_one")
      parent_entity: string          # Optional, parent entity type name (required when establishing entity relationships)
      field_constraints:             # Optional, specific field value constraints
        - field: string              # Field name
          value: any                 # Expected value or constraint description (string, number, boolean, or constraint text)
  state_requirements:                # Optional, pre-existing system state
    - description: string            # State description
      prerequisite_entity: string    # Dependent entity type
```

### Field Semantics

| Field | Required | Description |
|-------|----------|-------------|
| `entities` | Yes | List of entity declarations. Must contain at least 1 entry. |
| `entities[].entity_type` | Yes | The domain entity type name. Should match the naming used in PRD/Design documents (e.g., "Map", "Milestone", "Project"). |
| `entities[].min_count` | Yes | Minimum number of instances that must exist before the Step executes. Must be >= 1. |
| `entities[].relationship_type` | No | The relationship between this entity and its parent. Uses generic relational terminology independent of any ORM: `belongs_to`, `has_many`, `has_one`. |
| `entities[].parent_entity` | No | The entity type that this entity is related to. Required when `relationship_type` is specified. Must reference another `entity_type` in the same fixture_spec or a previously established entity. |
| `entities[].field_constraints` | No | List of field-level constraints that fixtures must satisfy. Used when a Step requires entities in a specific state (e.g., status = "pending"). |
| `entities[].field_constraints[].field` | Yes (within constraint) | The field name on the entity. |
| `entities[].field_constraints[].value` | Yes (within constraint) | The expected value or a constraint description. Type is `any` — may be a string, number, boolean, or natural language constraint (e.g., "any non-empty string", "> 0"). |
| `state_requirements` | No | System-level state that must hold before the Step. Used for non-entity prerequisites (e.g., "authentication enabled", "feature flag active"). |
| `state_requirements[].description` | Yes (within state_req) | Natural language description of the required state. |
| `state_requirements[].prerequisite_entity` | Yes (within state_req) | The entity type this state depends on. |

### Relationship Terminology

`relationship_type` uses ORM-independent, generic relational terms:

| Term | Meaning | Example |
|------|---------|---------|
| `belongs_to` | This entity is a child of the parent entity | Milestone belongs_to Map |
| `has_many` | This entity is a parent with multiple children (declared on the child side) | Project has_many Tasks |
| `has_one` | This entity has exactly one associated entity | User has_one Profile |

<HARD-RULE>
fixture_spec is a REQUIRED field in every Contract's Preconditions. The `entities` array MUST contain at least 1 entity declaration. A Contract without fixture_spec is schema-invalid.
</HARD-RULE>

## Minimal Valid Example (Single-Entity CRUD)

For Contracts that operate on a single entity type without relationships:

```yaml
## Outcome "success"
- Preconditions: "User is authenticated and authorized"
  fixture_spec:
    entities:
      - entity_type: "Project"
        min_count: 1
- Input: ...
- Output: ...
- State: ...
```

## Complete Example (Parent-Child Entity Relationship)

For Contracts that involve entity relationships and field constraints:

```yaml
## Outcome "success"
- Preconditions: "User is authenticated and authorized"
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
- Input: ...
- Output: ...
- State: ...
```

## Determining fixture_spec from Source Documents

When generating fixture_spec, use these sources in priority order:

1. **PRD/Design documents**: Identify entity types, relationships, and field values mentioned in the user story or domain model.
2. **Code reconnaissance (Fact Table)**: Verify entity types and relationships against source code models/database schemas.
3. **Journey document context**: Infer from the Step's user action what entities must pre-exist (e.g., "add milestone to map" implies both Map and Milestone entities).

### Heuristic: Step Action → Entity Inference

| Step Action Pattern | Inferred fixture_spec |
|---------------------|----------------------|
| Create a child entity (e.g., "add milestone to map") | Parent entity exists (min_count: 1) |
| List children of parent (e.g., "view map milestones") | Parent entity (min_count: 1) + children (min_count: >= 1) |
| Update entity state (e.g., "mark milestone complete") | Entity exists with specific initial state (field_constraints) |
| Delete entity (e.g., "remove milestone") | Entity exists (min_count: 1) |
| Simple CRUD on single entity | Entity exists (min_count: 1) |

## Backward Compatibility

When gen-test-scripts encounters a Contract whose Preconditions do NOT contain a `fixture_spec` field:

1. **Fallback**: Proceed using implicit fixture inference (current behavior) — infer minimal fixture needs from the Step's Input and Preconditions text.
2. **Warning**: Output a warning: `"Contract [file] has no fixture_spec. Falling back to implicit fixture inference. Consider regenerating Contracts with /gen-contracts to include fixture specifications."`
3. **Non-blocking**: The pipeline MUST NOT halt or fail when fixture_spec is missing. Existing Contracts without fixture_spec continue to function.

This ensures backward compatibility with Contracts generated before this dimension was introduced. New Contracts generated via `/gen-contracts` will automatically include fixture_spec.

<HARD-RULE>
When fixture_spec is absent from a Contract, downstream consumers MUST fall back to implicit inference and output a warning. They MUST NOT fail or block the pipeline.
</HARD-RULE>
