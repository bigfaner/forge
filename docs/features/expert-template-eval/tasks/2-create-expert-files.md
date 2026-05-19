---
id: "2"
title: "Create 9 scorer expert files"
priority: "P0"
estimated_time: "1.5h"
dependencies: []
type: "documentation"
mainSession: false
---

# 2: Create 9 scorer expert files

## Description

Create one expert file per scorer persona extracted from the persona selection table in `doc-scorer.md`. Each file contains ONLY the role description and domain-specific failure patterns (~20 lines). Expert files are composed with the scorer protocol by eval SKILL.md at invocation time.

## Reference Files
- `docs/proposals/expert-template-eval/proposal.md` — Source proposal
- `plugins/forge/agents/doc-scorer.md` — Source for persona table (lines 39-50)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/agents/experts/scorer/cto.md` | CTO persona for proposal eval |
| `plugins/forge/agents/experts/scorer/pm.md` | PM persona for PRD eval |
| `plugins/forge/agents/experts/scorer/architect.md` | Architect persona for design eval |
| `plugins/forge/agents/experts/scorer/ux-engineer.md` | UX persona for UI eval |
| `plugins/forge/agents/experts/scorer/qa.md` | QA persona for test-cases eval |
| `plugins/forge/agents/experts/scorer/editor.md` | Editor persona for consistency eval |
| `plugins/forge/agents/experts/scorer/harness-engineer.md` | Harness persona for harness eval |
| `plugins/forge/agents/experts/scorer/code-reviewer.md` | Code reviewer for validate-code |
| `plugins/forge/agents/experts/scorer/ux-auditor.md` | UX auditor for validate-ux |

### Modify
| File | Changes |
|------|---------|
| (none) | |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] 9 expert files created under `agents/experts/scorer/`
- [ ] Each file contains: role description + domain-specific failure patterns from the persona table
- [ ] No file exceeds ~30 lines
- [ ] No file duplicates workflow logic (that belongs in scorer-protocol.md)
- [ ] Each file is self-contained — no cross-references to other expert files
- [ ] Dispatch table mapping matches proposal: cto→proposal, pm→prd, architect→design, ux-engineer→ui-*, qa→test-cases, editor→consistency, harness-engineer→harness, code-reviewer→validate-code, ux-auditor→validate-ux

## Hard Rules

- Expert files must NOT contain scoring workflow steps — only persona role + failure patterns
- Keep each file ~20 lines. The goal is minimal domain injection, not comprehensive guides.

## Implementation Notes

- Source content: `doc-scorer.md` lines 39-50 contain the persona table with role descriptions and failure patterns
- The *(unmapped)* fallback persona ("Senior Technical Reviewer") does NOT get a file — eval SKILL.md handles unmapped types inline
- Key risk: Expert files that are too long dilute the protocol's focus, negating the separation benefit
