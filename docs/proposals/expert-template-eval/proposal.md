---
created: 2026-05-18
author: faner
status: Approved
---

# Proposal: Expert File Architecture for Eval Pipeline

## Problem

doc-scorer agent definition embeds ALL domain expert role descriptions in a single monolithic prompt. When evaluating a PRD, the scorer sees role descriptions for proposals, UI, QA, harness, etc. — creating noise that reduces evaluation precision. Each invocation should only see the expert role relevant to the document type being evaluated.

### Evidence

- `doc-scorer.md` contains a 10-row persona selection table spanning proposal, PRD, design, UI, QA, consistency, harness, validation types
- The scorer must parse and filter this table every invocation, adding cognitive load
- Domain-specific failure patterns from unrelated domains appear in the context window, diluting focus
- Current single-persona approach cannot benefit from multi-perspective analysis (e.g., a PRD evaluated by both PM and QA viewpoints)

### Urgency

The eval pipeline is a core quality gate. Imprecise scoring leads to either: (a) documents passing with hidden flaws, or (b) reviser wasting iterations on manufactured issues. As forge scales to more document types, the persona table grows and noise increases.

## Proposed Solution

Eliminate `doc-scorer.md` and `doc-reviser.md` agent definitions entirely. Split into two layers: **protocol files** (generic workflow) and **expert files** (scorer role + domain knowledge). Reviser uses only a protocol file + merged attacks — it doesn't need the rubric because attack points already prescribe what to fix; structural issues are caught by the scorer. Eval composes protocol + expert into per-expert prompts and spawns `general-purpose` agents directly.

### Architecture

```
agents/experts/
  protocol/
    scorer-protocol.md    # Three-phase adversarial scoring protocol
    reviser-protocol.md   # Generic attack-point-driven revision workflow (no expert file needed)
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
```

No `reviser/` subdirectory. The reviser is a generic executor — it reads attack points (which already contain domain-informed prescriptions) and edits documents. It doesn't need the rubric; if document structure is wrong, the scorer will flag it in attacks. Domain knowledge lives only in scorer expert files.

### Dispatch Table (in eval SKILL.md)

| type | scorer experts |
|------|---------------|
| proposal | [cto] |
| prd | [pm, qa] |
| design | [architect] |
| ui-web, ui-mobile, ui-tui | [ux-engineer] |
| test-cases, *-test-cases | [qa] |
| consistency | [editor] |
| harness | [harness-engineer] |
| validate-code | [code-reviewer] |
| validate-ux | [ux-auditor] |

### Scoring Flow (multi-expert example: prd)

```
Iteration N:
  eval reads protocol/scorer-protocol.md + scorer/pm.md + scorer/qa.md
  eval composes two prompts: [protocol + pm] and [protocol + qa]

  [PM scorer agent]    ──┐
  [QA scorer agent]    ──┼── eval LLM-merges attacks, averages scores → gate
                          │
  if score < target:
    eval reads protocol/reviser-protocol.md
    eval composes prompt: [reviser-protocol + merged-attacks]
    [reviser agent] ──→ edits docs
    → next iteration
```

### Innovation Highlights

- **Protocol–expert separation**: Scoring protocol (three-phase adversarial, verification stance, output format) is a single file shared by all experts. Expert files contain ONLY role description + domain-specific failure patterns (~20 lines each). Scorer protocol changes propagate to all scorer experts via one file.
- **No agent definition files**: `doc-scorer.md` and `doc-reviser.md` are deleted. Eval spawns `general-purpose` agents with composed prompts, passing `model: "sonnet"`.
- **Multi-expert parallel scoring**: Types like PRD get scored by PM + QA in parallel, producing diverse attack angles that a single persona cannot.
- **LLM-based semantic dedup**: Overlapping attack points from different experts are merged via LLM in the main session, not string matching.
- **Reviser is domain-agnostic**: Attack points already contain domain-informed prescriptions ("rewrite as Given/When/Then"). Reviser executes fixes without needing the rubric; structural issues are caught by the scorer.

## Requirements Analysis

### Key Scenarios

1. **Single-expert eval (harness, validate-code, validate-ux)**: One expert scores via general-purpose agent. Behavior equivalent to current single-scorer approach with cleaner prompt.
2. **Multi-expert eval (prd)**: PM + QA experts score in parallel — this is a new behavior (currently PRD uses single PM scorer). Eval LLM-merges attacks, averages scores for gate decision.
3. **Iteration loop**: Same as today — score → gate → revise — but reviser receives only merged attacks (no rubric). Reviser is generic, no expert file needed.
4. **Fallback**: Types not in the dispatch table use a generic fallback expert file (equivalent to current `*(unmapped)` persona row).

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
| Keep agents + inject expert | Expert as agent input | Backward compat | Agent definitions still exist as unnecessary wrapper | Rejected: unnecessary abstraction |
| **Protocol + expert files** | Inspired by academic peer review | Clean separation; no agent defs; protocol shared | Slightly heavier eval orchestration | **Selected: simplest architecture** |

## Feasibility Assessment

### Technical Feasibility

Fully achievable. The Agent tool supports `subagent_type: "general-purpose"` with `model: "sonnet"` and a custom `prompt`. Claude Code supports parallel agent spawning (multiple Agent tool calls in a single message). Eval already spawns subagents — this just changes what prompt it constructs.

### Resource & Timeline

~6-8 tasks. Expert files are short (~20 lines each). Protocol files are ~80 lines (extracted from current agent defs). Eval SKILL.md changes are the most complex part. Well-scoped for quick mode.

### Dependency Readiness

No external dependencies. All changes are within the forge plugin.

## Scope

### In Scope

- Extract scorer protocol from `doc-scorer.md` → `agents/experts/protocol/scorer-protocol.md`
- Extract reviser protocol from `doc-reviser.md` → `agents/experts/protocol/reviser-protocol.md`
- Create scorer expert files: `agents/experts/scorer/*.md` (9 experts: cto, pm, architect, ux-engineer, qa, editor, harness-engineer, code-reviewer, ux-auditor)
- Delete `agents/doc-scorer.md` and `agents/doc-reviser.md`
- Update eval SKILL.md:
  - Add type→experts dispatch table
  - Compose scorer prompt from protocol + expert file
  - Compose reviser prompt from protocol + merged attacks only (no rubric; structural issues caught by scorer)
  - Spawn `general-purpose` agents with `model: "sonnet"` and composed prompts
  - Parallel multi-expert spawning
  - LLM-based semantic dedup of attack points (main session)
  - Average scores across experts for gate decision
- Update forge-distribution.md (remove doc-scorer/doc-reviser from agent listing, add experts/ to distribution tree)

### Out of Scope

- Changing rubric scoring dimensions or criteria
- Adding new eval types
- Changes to non-eval skills or commands
- CLI/UI changes to eval commands
- Type-specific reviser strategies (reviser stays generic)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| `general-purpose` agent adds preamble to output, breaking score parsing | M | M | Protocol file ends with `<HARD-RULE>` enforcing exact format; eval uses loose regex (`/SCORE:\s*(\d+)\/(\d+)/`) |
| LLM dedup misses semantic duplicates or over-merges | M | L | Reviser processes attacks sequentially — if two attacks target same issue, second fix is a no-op |
| Token cost increase (2-3x per iteration for multi-expert types) | H | L | Only applies to multi-expert types (prd); single-expert types unchanged |
| Protocol file changes break all experts simultaneously | L | H | Protocol changes are rare; expert files are persona-only so unaffected by workflow changes |
| `${CLAUDE_SKILL_DIR}/../../agents/experts/` path fragility | L | M | Standard cross-skill reference pattern documented in forge-distribution.md |

## Success Criteria

- [ ] Each scorer invocation loads exactly one expert file — no cross-domain content in prompt
- [ ] Multi-expert types (prd) produce attack points from at least 2 distinct expert perspectives
- [ ] Gate decision uses averaged score across all experts
- [ ] Single-expert types score within ±50 points of current behavior on the same document
- [ ] Reviser receives protocol + merged attacks only (no rubric, no expert file)
- [ ] `doc-scorer.md` and `doc-reviser.md` are deleted; no agent definition files remain
- [ ] All expert/protocol files are distributed with the plugin package

## Next Steps

- Proceed to `/quick-tasks` to generate tasks from this proposal
