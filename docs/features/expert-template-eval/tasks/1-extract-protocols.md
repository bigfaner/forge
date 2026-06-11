---
id: "1"
title: "Extract scorer and reviser protocols"
priority: "P0"
estimated_time: "1h"
dependencies: []
type: "documentation"
mainSession: false
---

# 1: Extract scorer and reviser protocols

## Description

Extract the generic workflow portions from `doc-scorer.md` and `doc-reviser.md` into standalone protocol files under `agents/experts/protocol/`. The scorer protocol contains the three-phase adversarial scoring workflow (reasoning audit, rubric scoring, blindspot hunt). The reviser protocol contains the attack-point-driven revision workflow.

## Reference Files
- `docs/proposals/expert-template-eval/proposal.md` — Source proposal
- `plugins/forge/agents/doc-scorer.md` — Current scorer agent (source for extraction)
- `plugins/forge/agents/doc-reviser.md` — Current reviser agent (source for extraction)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/agents/experts/protocol/scorer-protocol.md` | Three-phase adversarial scoring protocol |
| `plugins/forge/agents/experts/protocol/reviser-protocol.md` | Generic attack-point-driven revision workflow |

### Modify
| File | Changes |
|------|---------|
| (none) | |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] `scorer-protocol.md` contains the three-phase workflow (Reasoning Audit → Rubric Scoring → Blindspot Hunt) extracted from `doc-scorer.md`
- [ ] `scorer-protocol.md` does NOT contain persona selection logic or domain-specific failure patterns
- [ ] `scorer-protocol.md` ends with the `<HARD-RULE>` output format specification (SCORE/DIMENSIONS/ATTACKS)
- [ ] `reviser-protocol.md` contains the attack-point-driven revision workflow extracted from `doc-reviser.md`
- [ ] `reviser-protocol.md` does NOT reference rubric path as an input (reviser receives only protocol + merged attacks)
- [ ] Both protocol files are self-contained — no references to `doc-scorer.md` or `doc-reviser.md`

## Hard Rules

- Do NOT copy persona selection or domain-specific failure patterns into protocol files — those belong in expert files (Task 2)
- Protocol files must use template variables (e.g., `{{RUBRIC_PATH}}`, `{{DOC_DIR}}`) compatible with eval SKILL.md prompt composition

## Implementation Notes

- The scorer protocol corresponds to lines 27-128 of `doc-scorer.md` (everything between the frontmatter and the output format), minus the "Persona Selection" section (lines 36-52)
- The reviser protocol corresponds to lines 22-99 of `doc-reviser.md`, minus rubric-related references (reviser no longer reads rubric per proposal design)
- Key risk: Protocol file changes affect all experts simultaneously. Keep the protocol stable and generic.
