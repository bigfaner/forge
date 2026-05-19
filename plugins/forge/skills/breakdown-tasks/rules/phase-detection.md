# Phase Detection Rules

**Load condition**: load this file IF the PRD has phase/gate structure (detected by explicit or heuristic patterns below).

**Guard clause**: if the PRD and tech-design contain no parseable phase or gate structure after applying all three detection tiers, skip this rule and proceed with artifact-driven decomposition (Step 3 fallback).

## Three-Tier Detection

Detect phases using three priority tiers. Stop at the first tier that produces results.

### Tier 1: Explicit Detection (highest priority)

Look for unambiguous phase markers in `prd/prd-spec.md` and `design/tech-design.md`:

- Flow diagrams with diamond decision nodes that represent quality gates or phase transitions
- PRD sections explicitly named "Round 1/2", "Phase 1/2", "Stage 1/2"

If found, use these as the definitive phase boundaries. Do not apply heuristic detection.

### Tier 2: Heuristic Detection

When no explicit structure is defined, scan `prd/prd-spec.md` and `design/tech-design.md` for these patterns:

- **Sequential markers**: "Round 1/2/3", "Phase/Stage 1/2/3", "Step 1/2/3", "第X阶段/轮", "第一轮/第二轮"
- **Conditional transitions**: "after X passes", "once X is verified", "X通过后", "确认X后再进行"
- **Go/no-go checkpoints**: "verify all tests pass", "confirm X before proceeding", "全部通过"
- **Gated prose**: "第一阶段...第二阶段...", "first pass...second pass...", "先X再Y"

Both English and Chinese patterns are recognized. If heuristic patterns are found, construct phases from the sequence implied by the markers.

### Tier 3: Fallback

When neither explicit nor heuristic patterns are detected, the Phase Inventory source is "fallback". In this case, Step 3 will use artifact-driven decomposition: group design elements into dependency layers based on what builds on what.

## Phase Inventory Format

Collect detection results and write to `tasks/phase-inventory.json`:

```json
[
  {"phase": 1, "name": "...", "source": "PRD-explicit|PRD-heuristic|design|fallback", "gates": [{"afterPhase": 1, "description": "..."}]},
  {"phase": 2, "name": "...", "source": "...", "gates": []}
]
```

Fields:
- `source`: how this phase was detected. One of `PRD-explicit`, `PRD-heuristic`, `design`, `fallback`.
- `gates`: array of gate objects. Each gate has `afterPhase` (the phase number it follows) and `description` (what must pass before proceeding). Empty array if no gate at this boundary.

This file persists the planning output for cross-step reference and later review.

## Maintenance Note

This rule file depends on the following sections in the skeleton SKILL.md:

- **Step 2: Map -> Tasks** — Phase & Gate Detection (detection trigger and inventory creation)
- **Step 3: Derive Phases & Dependencies** — consumes the Phase Inventory to structure phases and gates

If either of these sections changes in the skeleton, verify that the detection logic and inventory format in this file remain consistent.
