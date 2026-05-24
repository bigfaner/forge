---
id: "4"
title: "Improve quick-tasks SKILL.md Reference Files generation"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "doc"
mainSession: false
---

# 4: Improve quick-tasks SKILL.md Reference Files generation

## Description

Modify `plugins/forge/skills/quick-tasks/SKILL.md` to add guidance for generating precise section-level Reference Files in task files. Currently, quick-tasks has no instructions for populating `## Reference Files` — the template's `{{REFERENCE_FILES}}` placeholder defaults to bare `proposal.md` without section anchors.

## Reference Files
- `docs/proposals/spec-authority-enforcement/proposal.md#Scope` — quick-tasks scope item: "improve quick-tasks SKILL.md Reference Files generation"
- `docs/proposals/spec-authority-enforcement/proposal.md#Requirements-Analysis` — Scenario 4: quick-tasks with only proposal.md input

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/quick-tasks/SKILL.md` | Add Reference Files generation guidance in Step 2 "Derive Tasks" section |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] SKILL.md Step 2 "Derive Tasks" section includes explicit instructions for generating Reference Files with section-level precision
- [ ] Instructions require format: `proposal.md#Section-Title — brief description of what the section defines`
- [ ] Instructions specify that when proposal.md is the only input (quick-tasks has no tech-design.md), Reference Files must reference 2-5 specific sections from proposal.md relevant to each task
- [ ] Instructions specify extraction logic: for each task, identify which proposal sections are most relevant based on the task's description and affected files
- [ ] If proposal.md references external design documents that exist on disk, include those document sections as additional Reference Files
- [ ] Each generated coding task must have >=1 precise section reference (not just bare file path `proposal.md`)

## Hard Rules
- MUST load `docs/conventions/forge-distribution.md` before modifying SKILL.md
- MUST NOT change the overall SKILL.md structure — only add Reference Files guidance within existing Step 2

## Implementation Notes

Add a subsection under Step 2 "Derive Tasks" after the existing content, titled something like "Reference Files Generation". The guidance should cover:

1. For each derived task, identify 2-5 relevant sections from proposal.md
2. Use format: `proposal.md#Section-Title — description of what this section defines for the task`
3. Match sections by task description keywords and affected file paths
4. If proposal.md references existing design docs (e.g., `docs/lessons/` or `docs/conventions/`), include relevant sections from those docs too
5. Ensure each coding task has at least 1 section-level reference — bare file paths are insufficient
