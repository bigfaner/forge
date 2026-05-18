---
created: 2026-05-18
author: faner
status: Draft
---

# Proposal: Expert Template Architecture for Eval Pipeline

## Problem

doc-scorer and doc-reviser agent definitions embed ALL domain expert personas in a single monolithic prompt. When evaluating a PRD, the scorer sees persona descriptions for proposals, UI, QA, harness, etc. — creating noise that reduces evaluation precision. Each invocation should only see the expert persona relevant to the document type being evaluated.

### Evidence

- `doc-scorer.md` contains a 10-row persona selection table spanning proposal, PRD, design, UI, QA, consistency, harness, validation types
- The scorer must parse and filter this table every invocation, adding cognitive load
- Domain-specific failure patterns from unrelated domains appear in the context window, diluting focus
- Current single-persona approach cannot benefit from multi-perspective analysis (e.g., a PRD evaluated by both PM and QA viewpoints)

### Urgency

The eval pipeline is a core quality gate. Imprecise scoring leads to either: (a) documents passing with hidden flaws, or (b) reviser wasting iterations on manufactured issues. As forge scales to more document types, the persona table grows and noise increases.

## Proposed Solution

Eliminate `doc-scorer.md` and `doc-reviser.md` agent definitions entirely. Split into three layers: **protocol files** (generic workflow), **expert files** (persona + domain knowledge), and **eval dispatch** (composition + orchestration). Eval composes protocol + expert into a prompt and spawns `general-purpose` agents directly.

### Architecture

```
agents/experts/
  protocol/
    scorer-protocol.md    # Three-phase adversarial scoring protocol
    reviser-protocol.md   # Attack-point-driven revision workflow
  scorer/
    cto.md                # CTO persona for proposals
    pm.md                 # PM persona for PRDs
    architect.md          # Architect persona for designs
    ux-engineer.md        # UX persona for UI designs
    qa.md                 # QA persona for test cases
    editor.md             # Editor persona for consistency
    harness-engineer.md   # Harness persona
    code-reviewer.md      # Code reviewer for validate-code
    ux-auditor.md         # UX auditor for validate-ux
  reviser/
    cto.md                # CTO revision strategies
    pm.md                 # PM revision strategies
    architect.md          # Architect revision strategies
    ux-engineer.md        # UX revision strategies
    qa.md                 # QA revision strategies
    editor.md             # Editor revision strategies
    harness-engineer.md   # Harness revision strategies
    code-reviewer.md      # Code reviewer revision strategies
    ux-auditor.md         # UX auditor revision strategies
```

### Dispatch Table (in eval SKILL.md)

| type | scorer experts | reviser expert |
|------|---------------|----------------|
| proposal | [cto] | cto |
| prd | [pm, qa] | pm |
| design | [architect] | architect |
| ui-web, ui-mobile, ui-tui | [ux-engineer] | ux-engineer |
| test-cases, *-test-cases | [qa] | qa |
| consistency | [editor] | editor |
| harness | [harness-engineer] | harness-engineer |
| validate-code | [code-reviewer] | code-reviewer |
| validate-ux | [ux-auditor] | ux-auditor |

### Scoring Flow (multi-expert example: prd)

```
Iteration N:
  eval reads protocol/scorer-protocol.md + scorer/pm.md + scorer/qa.md
  eval composes two prompts: [protocol + pm] and [protocol + qa]

  [PM scorer agent]    ──┐
  [QA scorer agent]    ──┼── eval merges attacks, averages scores → gate
                        │
  if score < target:
    eval reads protocol/reviser-protocol.md + reviser/pm.md
    eval composes prompt: [reviser-protocol + pm + merged-attacks]
    [PM reviser agent] ──→ edits docs
    → next iteration
```

### Innovation Highlights

- **Protocol-persona separation**: Scoring protocol (three-phase adversarial, verification stance, output format) is a single file shared by all experts. Expert files contain ONLY persona + domain-specific failure patterns (~20 lines each). Protocol changes: update one file.
- **No agent definition files**: `doc-scorer.md` and `doc-reviser.md` are deleted. Eval spawns `general-purpose` agents with composed prompts, passing `model: "sonnet"`.
- **Multi-expert parallel scoring**: Types like PRD get scored by PM + QA in parallel, producing diverse attack angles that a single persona cannot.
- **Merge-dedup-rank pipeline**: Overlapping attack points from different experts are merged before reaching the reviser.

## Requirements Analysis

### Key Scenarios

1. **Single-expert eval (harness, validate-code, validate-ux)**: One expert scores via general-purpose agent. Same behavior as today but with cleaner, focused prompt.
2. **Multi-expert eval (prd)**: PM + QA experts score in parallel. Eval merges attacks, averages scores for gate decision.
3. **Iteration loop**: Same as today — score → gate → revise — but reviser receives merged attacks from all experts plus its own type-specific revision strategies.
4. **Backward compat**: Types not in the dispatch table use a generic fallback expert. Existing behavior preserved.

### Non-Functional Requirements

- **Performance**: Multi-expert parallel scoring should take similar wall-clock time as current single-scorer (parallelism absorbs the overhead)
- **Cost**: 2-3x token cost per iteration for multi-expert types (acceptable for quality gate)
- **Distribution**: All files under `agents/experts/` are distributed with the plugin

### Constraints & Dependencies

- Protocol and expert files live in `agents/experts/` and are distributed with the plugin
- Eval SKILL.md owns the type→experts dispatch table (no rubric modifications)
- Eval composes full prompt and spawns `general-purpose` agents via Agent tool with `model: "sonnet"`
- `doc-scorer.md` and `doc-reviser.md` are deleted

## Alternatives & Industry Benchmarking

### Industry Solutions

Multi-perspective evaluation is common in peer review systems. Academic peer review uses 2-3 reviewers per paper with different expertise. Code review best practices recommend at least 2 reviewers. LLM-based evaluation frameworks increasingly use multi-model or multi-prompt ensembles to reduce single-evaluator bias.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero effort | Noise persists, single perspective limited | Rejected: root cause is precision |
| Expert-as-agent | Custom agent per type | Simple dispatch | Duplicated protocol across 10+ agent files | Rejected: maintenance burden |
| Keep agents + inject expert | Expert as agent input | Backward compat | Agent definitions still exist as unnecessary wrapper; two-layer indirection | Rejected: unnecessary abstraction |
| **Protocol + expert files** | Inspired by academic peer review | Clean three-layer separation; no agent defs; protocol shared | Slightly heavier eval orchestration | **Selected: simplest architecture** |

## Feasibility Assessment

### Technical Feasibility

Fully achievable. The Agent tool supports `subagent_type: "general-purpose"` with `model: "sonnet"` and a custom `prompt`. Eval already spawns subagents — this just changes what prompt it constructs and removes the agent definition layer.

### Resource & Timeline

~6-8 tasks. Expert files are short (~20 lines each). Protocol files are ~80 lines (extracted from current agent defs). Eval SKILL.md changes are the most complex part. Well-scoped for quick mode.

### Dependency Readiness

No external dependencies. All changes are within the forge plugin.

## Scope

### In Scope

- Extract scorer protocol from `doc-scorer.md` → `agents/experts/protocol/scorer-protocol.md`
- Extract reviser protocol from `doc-reviser.md` → `agents/experts/protocol/reviser-protocol.md`
- Create expert scorer files: `agents/experts/scorer/*.md` (9 experts: cto, pm, architect, ux-engineer, qa, editor, harness-engineer, code-reviewer, ux-auditor)
- Create expert reviser files: `agents/experts/reviser/*.md` (9 matching experts)
- Delete `agents/doc-scorer.md` and `agents/doc-reviser.md`
- Update eval SKILL.md:
  - Add type→experts dispatch table
  - Compose prompt from protocol + expert file content
  - Spawn `general-purpose` agents with `model: "sonnet"` and composed prompt
  - Parallel multi-expert spawning
  - Merge + dedup + rank attack points
  - Average scores across experts for gate decision
  - Compose reviser prompt from reviser-protocol + reviser expert + merged attacks
- Update forge-distribution.md (remove doc-scorer/doc-reviser from agent listing, add experts/ to distribution tree)

### Out of Scope

- Changing rubric scoring dimensions or criteria
- Adding new eval types
- Changes to non-eval skills or commands
- CLI/UI changes to eval commands

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Merge logic produces contradictory attacks | M | M | Merge preserves all unique attacks; dedup only combines overlapping ones from different experts |
| Token cost increase (2-3x per iteration) | H | L | Only applies to multi-expert types; single-expert types unchanged |
| Protocol files drift from expert expectations | M | M | Expert files are persona-only; protocol changes automatically apply to all experts |
| Parallel spawning hits concurrency limits | L | M | Fallback to sequential if parallel fails; max 3 experts per type |
| Agent tool prompt size limits | L | L | Composed prompt (protocol ~80 lines + expert ~20 lines) is well within limits |

## Success Criteria

- [ ] Each scorer invocation receives exactly one expert persona — no cross-domain noise
- [ ] Multi-expert types (prd) produce attack points from at least 2 distinct perspectives
- [ ] Merged attack list has zero duplicates (same quote + same issue)
- [ ] Gate decision uses averaged score across all experts
- [ ] Single-expert types produce equivalent output to current behavior
- [ ] `doc-scorer.md` and `doc-reviser.md` are deleted; no agent definition files remain
- [ ] All expert files are distributed with the plugin package

## Next Steps

- Proceed to `/quick-tasks` to generate tasks from this proposal
