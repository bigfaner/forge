---
status: "completed"
started: "2026-06-08 18:02"
completed: "2026-06-08 18:05"
time_spent: "~3m"
---

# Task Record: 2 gen-contracts 新增 Fixture Specification 维度

## Summary
Added Fixture Specification dimension to gen-contracts skill: created rules/fixture-spec.md with schema definition, minimal and complete examples, and backward compatibility rules; updated SKILL.md with fixture-spec validation check and section 4.2.1; updated dimension-rules.md with fixture_spec sub-dimension; updated contract.md and outcome-block.md templates with fixture_spec fields.

## Changes

### Files Created
- plugins/forge/skills/gen-contracts/rules/fixture-spec.md

### Files Modified
- plugins/forge/skills/gen-contracts/SKILL.md
- plugins/forge/skills/gen-contracts/rules/dimension-rules.md
- plugins/forge/skills/gen-contracts/templates/contract.md
- plugins/forge/skills/gen-contracts/templates/outcome-block.md

### Key Decisions
无

## Document Metrics
1 new rule file (~130 lines), 4 files modified, schema with 8 fields, 2 examples, backward compatibility section

## Referenced Documents
- docs/proposals/behavioral-test-accuracy/proposal.md
- plugins/forge/skills/gen-contracts/SKILL.md
- plugins/forge/skills/gen-contracts/rules/dimension-rules.md
- plugins/forge/skills/gen-contracts/templates/contract.md

## Review Status
final

## Acceptance Criteria
- [x] rules/fixture-spec.md defines Fixture Specification Schema (entities with entity_type/min_count/relationship_type/parent_entity/field_constraints, optional state_requirements)
- [x] rules/fixture-spec.md contains minimal example (single-entity CRUD) and complete example (parent-child entity relationship)
- [x] Contract template templates/contract.md Preconditions section includes fixture_spec field with complete schema structure
- [x] rules/fixture-spec.md declares fixture_spec as required field, entities must contain at least 1 entity declaration

## Notes
Backward compatibility HARD-RULE included per task requirement: downstream consumers must fallback to implicit inference with warning when fixture_spec is absent. Implementation Notes respected: declarative approach, generic relationship terms (belongs_to/has_many/has_one), any-type value field.
