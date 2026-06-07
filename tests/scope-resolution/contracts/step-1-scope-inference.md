# Contract: scope-resolution / Step 1: Scope Inference

## Outcome "mixed-project-scope-field"
- Preconditions: "breakdown-tasks skill SKILL.md exists"
- Input: "Read breakdown-tasks SKILL.md content"
- Output: "content contains 'scope' field description"
- State: "no state changes"
- Side-effect: none

## Outcome "frontend-only-scope"
- Preconditions: "breakdown-tasks skill SKILL.md exists"
- Input: "Read breakdown-tasks SKILL.md content"
- Output: "content contains 'frontend' and 'scope' (scope=frontend logic)"
- State: "no state changes"
- Side-effect: none

## Outcome "cross-scope-all"
- Preconditions: "breakdown-tasks skill SKILL.md exists"
- Input: "Read breakdown-tasks SKILL.md content"
- Output: "content contains 'all' and 'scope' (scope=all for cross-scope)"
- State: "no state changes"
- Side-effect: none

## Outcome "non-mixed-default-all"
- Preconditions: "breakdown-tasks skill SKILL.md exists"
- Input: "Read breakdown-tasks SKILL.md content"
- Output: "content describes default scope=all for non-mixed projects"
- State: "no state changes"
- Side-effect: none

## Journey Invariants
- breakdown-tasks SKILL.md is the authoritative scope inference source
- scope values are one of: frontend, backend, all
