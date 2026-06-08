# Six-Dimension Declaration Rules

For each Outcome, declare all six dimensions. Four are mandatory (non-empty), two are optional:

**Mandatory** (every Outcome MUST have non-empty values):

| Dimension | Content | Source |
|-----------|---------|--------|
| Preconditions | State that must hold before execution, including `fixture_spec` (see below) | Journey edge case preconditions + Fact Table state info |
| Input | What goes into the system | Journey user action + Fact Table command/flag/endpoint info |
| Output | What the system produces | Journey expected result + Fact Table output patterns |
| State | How system state changes | Fact Table state storage info + output inference |

**Optional** (may be omitted; omission = no constraint):

| Dimension | Content | When to include |
|-----------|---------|-----------------|
| Side-effect | External effects (hooks, network calls, async Cmds) | When Fact Table reveals side effects |
| Invariants (step-level) | Properties within the step | When the step has internal consistency requirements |

**Side-effect defaults**: When omitted, Side-effect defaults to `none`.
**Step-level Invariants defaults**: When omitted, no step-level invariant constraint.

## Preconditions Sub-Dimension: fixture_spec

Every Outcome's Preconditions MUST include a `fixture_spec` field per `rules/fixture-spec.md`. This field is a structured declaration of the pre-existing data state required before the Step executes.

**fixture_spec structure** (full schema in `rules/fixture-spec.md`):

```yaml
fixture_spec:
  entities:
    - entity_type: string
      min_count: integer
      relationship_type: string    # optional
      parent_entity: string        # optional
      field_constraints:           # optional
        - field: string
          value: any
  state_requirements:              # optional
    - description: string
      prerequisite_entity: string
```

**Rules**:
- `fixture_spec` is REQUIRED — a Contract without it is schema-invalid
- `entities` MUST contain at least 1 entity declaration
- `entity_type` should match domain model naming from PRD/Design documents
- `relationship_type` uses generic terms: `belongs_to`, `has_many`, `has_one`
- `field_constraints[].value` is `any` type — allows string, number, boolean, or natural language constraint descriptions

# Semantic Descriptors

All dimension values use semantic descriptors -- natural language descriptions of expected behavior.

**Rules**:
- MUST NOT contain regex syntax (`\d`, `.*`, `[^...]`, `(?:...)`, `\s`, `\w`, `\b`, `$`, `^` as anchor, etc.)
- MUST NOT contain framework-specific assertion patterns
- MUST be natural language expressing business intent

**Good examples**:
- `"success confirmation containing feature-slug"`
- `"task status changed from pending to in_progress"`
- `"stderr contains error message about missing feature"`

**Bad examples** (these belong in gen-test-scripts):
- `"Feature\s+([\w-]+)\s+created"` (regex)
- `"assert.Equal(t, 0, exitCode)"` (framework assertion)
- `"matches pattern /task_\d+/"` (regex reference)

<HARD-RULE>
Semantic descriptors MUST NOT contain regex syntax. gen-contracts stage does not generate regex. If you find yourself writing a pattern match, replace it with a natural language description of what the pattern matches.
</HARD-RULE>

# Preconditions Mutual Exclusivity

Each Outcome within a Step MUST have Preconditions that are mutually exclusive with all other Outcomes in the same Step.

**Mutual exclusivity rule**: For any given system state, at most one Outcome's Preconditions can be satisfied. This prevents combinatorial explosion.

**Validation**: Before writing a Contract, verify that no two Outcomes in the same Step have identical or overlapping Preconditions. If overlap is detected:
1. Differentiate the Preconditions (add a distinguishing condition)
2. If impossible, merge the Outcomes into a single Outcome with disjunctive Preconditions
3. Never write Outcomes whose Preconditions can be simultaneously satisfied

<HARD-RULE>
Outcomes MUST be mutually exclusive by Preconditions. If two Outcomes' Preconditions can both be true for the same system state, the Contract is invalid and must be fixed before writing.
</HARD-RULE>

**Outcome count checkpoint**: Steps with more than 5 Outcomes trigger a review. Consider merging semantically similar Outcomes. Do not automatically exceed 5 Outcomes without explicit justification.
