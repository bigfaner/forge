---
id: "5"
title: "Improve breakdown-tasks SKILL.md Reference Files generation for non-UI tasks"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "doc"
mainSession: false
---

# 5: Improve breakdown-tasks SKILL.md Reference Files generation for non-UI tasks

## Description

Modify `plugins/forge/skills/breakdown-tasks/SKILL.md` to add Reference Files filling guidance for non-UI tasks. Currently, breakdown-tasks has `rules/ui-placement.md` for UI Reference Files but no guidance for non-UI tasks. Each generated task should include precise tech-design.md section references.

## Reference Files
- `docs/proposals/spec-authority-enforcement/proposal.md#Scope` — breakdown-tasks scope item: "improve breakdown-tasks SKILL.md Reference Files filling"
- `docs/proposals/spec-authority-enforcement/proposal.md#Requirements-Analysis` — Scenario 5: breakdown-tasks with tech-design.md input

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Add Reference Files filling guidance in Step 4a "Business Tasks" section |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] SKILL.md Step 4a "Business Tasks" section includes explicit instructions for populating Reference Files with section-level precision from tech-design.md
- [ ] Instructions require format: `path/to/tech-design.md#Section-Title — brief description`
- [ ] Instructions specify a 4-step extraction heuristic for each task: (1) extract file paths from `## Affected Files`, (2) search tech-design.md sections mentioning those paths, (3) extract architecture decision sections matching task description keywords, (4) merge and keep 2-5 most relevant sections
- [ ] Instructions include a checklist: each generated task must have >=1 design-level Reference File
- [ ] The guidance is described as heuristic strategy (not algorithm), since the executor is an LLM agent
- [ ] UI tasks continue to use existing `rules/ui-placement.md` Reference File requirements — no conflict

## Hard Rules
- MUST load `docs/conventions/forge-distribution.md` before modifying SKILL.md
- MUST NOT change the overall SKILL.md structure — only add Reference Files guidance within existing Step 4a
- MUST NOT modify or conflict with existing `rules/ui-placement.md` UI Reference File requirements

## Implementation Notes

Add a subsection under Step 4a "Business Tasks", after the existing content about Hard Rules and before Scope Assignment. Title: "Reference Files Population" or similar.

The guidance should cover:

1. For each business task, determine relevant tech-design.md sections:
   - Extract file paths from task's `## Affected Files`
   - Search tech-design.md for sections that mention these file paths
   - Also find architecture decision sections matching task description keywords
   - Merge, deduplicate, keep 2-5 most relevant sections
2. Format: `design/tech-design.md#Section-Title — brief description of what the section defines`
3. Checklist: every task must have >=1 design-level Reference File
4. For tasks without clear tech-design.md matches, fall back to the most relevant architecture overview section
5. This is heuristic guidance for LLM execution, not a deterministic algorithm
