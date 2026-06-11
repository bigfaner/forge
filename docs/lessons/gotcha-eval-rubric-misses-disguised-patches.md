---
created: "2026-05-18"
tags: [architecture, interface, testing]
---

# Eval Rubric Misses Disguised Patches in Refactoring Designs

## Problem

The design eval rubric (1000-point scale, 7 dimensions) scored the Forge Architecture Simplification tech design at 945/1000 (iteration 1) and 958/1000 (iteration 2). Yet a subsequent expert review caught 5 patch patterns that the eval completely missed:

1. `ViaSubmit bool` — a disguised flag-driven branch, the same structural defect the refactoring aims to eliminate
2. `TransitionOpts` coupling 3 concerns (override, identity, dependency resolution) in one struct
3. `CanAutoUnblock` as standalone interface when it's internal guard logic
4. `SaveStateAtomic` / `SaveIndexAtomic` — two functions doing the same temp+rename operation
5. Exit code 2 semantics incorrectly assumed to apply to all commands (only applies to hook execution)

The eval gave D2 (Interface & Model Definitions) 161/170 and D3 (Error Handling) 130/130 (max). It praised the design as "directly implementable" and "sound."

## Root Cause

**Level 1**: The rubric evaluates *completeness* (are signatures typed? are models concrete?) and *correctness* (do exit codes match conventions?), but not *structural quality* (does the design introduce new forms of the problems it aims to solve?).

**Level 2**: The rubric's criteria are **document-quality metrics**, not **design-quality metrics**. "Interface signatures typed" checks whether Go function signatures exist, not whether the interface has good separation of concerns. "Directly implementable" checks whether types are specified, not whether the contract is cohesive.

**Level 3**: The rubric has no dimension for **"self-consistency of refactoring intent vs. design output"**. A refactoring proposal that says "eliminate 4-path-4-rule" and then introduces `ViaSubmit bool` (which is just a 5th path with a flag) passes every rubric criterion because each criterion measures the design in isolation, not against its own stated goals.

This is the deepest blindspot: **the eval measures whether the design document is well-formed, not whether the design itself is well-formed for its stated purpose**.

## Solution

For refactoring-focused designs, add an adversarial check that the current eval rubric cannot perform:

1. **Patch regression audit**: For each interface, ask "does this introduce the same structural pattern it aims to eliminate?" This requires reading the PRD's problem statement alongside each interface definition — a cross-reference the rubric doesn't encode.

2. **Cohesion check on parameter objects**: When a struct has 5+ fields serving 3+ concerns (like `TransitionOpts`), the rubric's "directly implementable" criterion should flag this as a cohesion smell, not reward it.

3. **Execution context verification**: When a design claims external semantics (e.g., "Claude Code treats exit 2 as blocking"), the rubric should require evidence — not just a reference to documentation, but verification that the semantics apply to the specific execution context (hook vs. Bash tool).

## Reusable Pattern

When evaluating a **refactoring design** (as opposed to a greenfield design), supplement the standard eval rubric with:

- **Intent-output consistency**: Does each interface actually embody the refactoring's stated goals, or does it replicate the old pattern in new clothes?
- **Cohesion per parameter**: If a function needs a parameter object with 5+ fields, the design should justify why the concerns can't be separated.
- **Context-specific claims**: Any claim about external system behavior must specify the *execution context*, not just cite documentation.

The eval rubric works well for greenfield designs. For refactoring designs, it needs a "don't re-introduce what you're removing" dimension.

## Example

```
ViaSubmit bool in TransitionOpts:
- Eval saw: "typed parameter, clear contract" ✓ → 58/60
- Guru saw: same flag-driven branching the PRD says to eliminate → patch
- Why eval missed: no criterion measures "intent-output consistency"
```

## Related Files

- `docs/features/forge-architecture-simplification/design/tech-design.md`
- `docs/features/forge-architecture-simplification/eval/iteration-2.md`

## References

- Design eval rubric: `plugins/forge/skills/eval/rubrics/design.md`
- PRD problem statement: BC-1~BC-3 (4 paths with 4 different rules)
