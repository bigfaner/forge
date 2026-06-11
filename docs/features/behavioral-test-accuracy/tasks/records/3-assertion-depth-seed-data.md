---
status: "completed"
started: "2026-06-08 18:06"
completed: "2026-06-08 18:12"
time_spent: "~6m"
---

# Task Record: 3 gen-test-scripts 断言深度 + seed data 丰富度规则

## Summary
Created assertion depth rule (assertion-depth.md) with classification criteria table, >=80% behavioral threshold, and >=30% deep assertion requirement. Created fixture-from-spec rule (fixture-from-spec.md) with entity count, relationship, field constraints, and multi-step fixture passing rules. Modified SKILL.md to add rule references in Step 2.5.4 and post-compile checks in Step 4.4. Modified quality-gates.md to add assertion depth gate with enforcement flow. Modified types/_shared.md to add fixture_spec consumption logic with backward compatibility handling.

## Changes

### Files Created
- plugins/forge/skills/gen-test-scripts/rules/assertion-depth.md
- plugins/forge/skills/gen-test-scripts/rules/fixture-from-spec.md

### Files Modified
- plugins/forge/skills/gen-test-scripts/SKILL.md
- plugins/forge/skills/gen-test-scripts/rules/quality-gates.md
- plugins/forge/skills/gen-test-scripts/types/_shared.md

### Key Decisions
无

## Document Metrics
assertion-depth: ~170 lines, fixture-from-spec: ~160 lines, SKILL.md: +20 lines, quality-gates: +25 lines, _shared.md: +25 lines

## Referenced Documents
- docs/proposals/behavioral-test-accuracy/proposal.md
- plugins/forge/skills/gen-test-scripts/SKILL.md
- plugins/forge/skills/gen-test-scripts/rules/quality-gates.md
- plugins/forge/skills/gen-test-scripts/types/_shared.md

## Review Status
final

## Acceptance Criteria
- [x] rules/assertion-depth.md contains complete assertion classification table (behavioral vs structural, with edge case explanations)
- [x] rules/assertion-depth.md declares >=80% behavioral assertion mandatory rule and >=30% deep assertion requirement
- [x] rules/fixture-from-spec.md declares reading from Contract fixture_spec.entities and generating fixture satisfying min_count
- [x] rules/fixture-from-spec.md declares >=N child entity creation with relationship handling when fixture_spec requires N
- [x] types/_shared.md adds backward compatibility: fallback to implicit inference with warning when fixture_spec absent

## Notes
All 5 acceptance criteria met. Hard Rules respected: only modified SKILL.md, rules/assertion-depth.md, rules/fixture-from-spec.md, rules/quality-gates.md, types/_shared.md. Cross-references between new rules and existing rules are consistent.
