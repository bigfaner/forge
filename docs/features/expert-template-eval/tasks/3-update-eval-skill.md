---
id: "3"
title: "Update eval SKILL.md for expert-based orchestration"
priority: "P0"
estimated_time: "2h"
dependencies: ["1", "2"]
type: "documentation"
mainSession: false
breaking: true
---

# 3: Update eval SKILL.md for expert-based orchestration

## Description

Rewrite the eval SKILL.md to use protocol+expert composition instead of spawning `doc-scorer`/`doc-reviser` agent types. The eval skill now reads protocol and expert files, composes full prompts, and spawns `general-purpose` agents directly.

## Reference Files
- `docs/proposals/expert-template-eval/proposal.md` — Source proposal
- `plugins/forge/skills/eval/SKILL.md` — Current eval skill (being modified)
- `plugins/forge/agents/experts/protocol/scorer-protocol.md` — Created in Task 1
- `plugins/forge/agents/experts/protocol/reviser-protocol.md` — Created in Task 1
- `plugins/forge/agents/experts/scorer/*.md` — Created in Task 2

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/SKILL.md` | Replace agent spawning with protocol+expert composition |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] Step 2 (Invoke Scorer) reads protocol + expert file(s), composes prompt, spawns `general-purpose` agents with `model: "sonnet"`
- [ ] Dispatch table maps each eval type to its scorer expert(s) per proposal table
- [ ] Multi-expert types (prd → [pm, qa]) spawn parallel agents in a single message
- [ ] Gate decision (Step 3b) averages scores across all experts for multi-expert types
- [ ] Attack points from multiple experts are LLM-merged (semantic dedup) in main session before passing to reviser
- [ ] Reviser (Step 4) receives reviser-protocol + merged attacks only — no rubric, no expert file
- [ ] Reviser spawned as `general-purpose` agent, not `forge:doc-reviser`
- [ ] Fallback for unmapped types: use generic inline prompt (equivalent to current *(unmapped)* persona)
- [ ] No references to `doc-scorer` or `doc-reviser` agent types remain in SKILL.md
- [ ] Context injection (from rubric `context` frontmatter) still works — appended after expert content in composed prompt

## Hard Rules

- Path to expert files: `${CLAUDE_SKILL_DIR}/../../agents/experts/` (standard cross-skill reference per forge-distribution.md)
- All scorer agents must use `model: "sonnet"` via Agent tool parameter
- The output format parsing (`SCORE: X/Y`, `DIMENSIONS:`, `ATTACKS:`) must remain unchanged — downstream consumers depend on it
- Do NOT change the iteration loop structure (Step 2 → Step 3 → Step 4). Only change HOW agents are spawned and WHAT prompts they receive.

## Implementation Notes

- Most complex task. The core change is in Step 2 and Step 4 of SKILL.md.
- Step 2 currently says: "Spawn `doc-scorer` via Agent tool (subagent_type: `forge:doc-scorer` or `general-purpose`)" — change to read protocol file, read expert file(s), compose prompt, spawn `general-purpose` with composed prompt
- For multi-expert: spawn multiple agents in parallel (multiple Agent tool calls in single message), then LLM-merge attack points in main session using a prompt like "Merge overlapping attack points from N expert evaluations. Keep unique attacks from each. Combine duplicates into single attacks preserving the strongest prescription."
- Step 4 currently says: "Spawn `doc-reviser` via Agent tool (subagent_type: `forge:doc-reviser` or `general-purpose`)" — change to read reviser protocol, compose prompt with merged attacks, spawn `general-purpose`
- Key risk: `general-purpose` agent may add preamble to output, breaking score parsing. Mitigate with `<HARD-RULE>` in scorer protocol enforcing exact format, and use loose regex for parsing.
