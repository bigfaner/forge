---
id: "1"
title: "Inline decision-logging.md into consuming skills"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "refactor"
mainSession: false
---

# 1: Inline decision-logging.md into consuming skills

## Description
Replace `${CLAUDE_SKILL_DIR}/../../references/shared/decision-logging.md` read instructions in 3 skill files with the full inlined protocol content. The decision-logging reference defines a domain-to-file mapping and an archiving flow used by consolidate-specs, tech-design, and learn.

## Reference Files
- `docs/proposals/remove-references-dir/proposal.md` — Source proposal
- `plugins/forge/references/shared/decision-logging.md` — Content to inline

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/consolidate-specs/SKILL.md` | Replace line 207 reference with inlined domain-to-decision-file mapping (Section 1) |
| `plugins/forge/skills/tech-design/SKILL.md` | Replace line 191 reference with inlined decision archiving flow (Section 2) |
| `plugins/forge/skills/learn/SKILL.md` | Replace line 105 reference with inlined authoritative format |

### Delete
| File | Reason |
|------|--------|
| (none — directory deletion is task 6) | |

## Acceptance Criteria
- [ ] No occurrence of `references/shared/decision-logging` in any of the 3 modified files
- [ ] Each file contains the relevant decision-logging protocol section verbatim (not a summary)
- [ ] `${CLAUDE_SKILL_DIR}` path references to decision-logging.md are fully removed

## Hard Rules
- Inline the full content of the relevant sections, not a paraphrased version
- Do not modify any logic or flow — this is a structural copy-paste only

## Implementation Notes
- decision-logging.md has 2 sections: Section 1 (domain-to-file mapping, used by consolidate-specs) and Section 2 (archiving flow, used by tech-design and learn)
- The learn skill references the "authoritative format" — inline the complete decision-logging protocol
