---
id: "2"
title: "Info commands (proposal, feature, lesson)"
priority: "P1"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 2: Info commands (proposal, feature, lesson)

## Description

Implement three read-only query commands for project artifacts: `forge proposal` (list + detail), `forge feature list` + `forge feature status` (extend existing), and `forge lesson` (list + detail). These commands let agents and users discover project state without manual `ls` + file reading.

## Reference Files
- `docs/proposals/forge-info-commands/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|
| `forge-cli/internal/cmd/proposal.go` | `forge proposal` list and detail commands |
| `forge-cli/internal/cmd/proposal_test.go` | Tests for proposal commands |
| `forge-cli/internal/cmd/lesson.go` | `forge lesson` list and detail commands |
| `forge-cli/internal/cmd/lesson_test.go` | Tests for lesson commands |
| `forge-cli/pkg/proposal/proposal.go` | Proposal discovery and parsing logic |
| `forge-cli/pkg/proposal/proposal_test.go` | Tests for proposal package |
| `forge-cli/pkg/lesson/lesson.go` | Lesson discovery and parsing logic |
| `forge-cli/pkg/lesson/lesson_test.go` | Tests for lesson package |

### Modify
| File | Changes |
|------|---------|
| `forge-cli/internal/cmd/feature.go` | Add `list` and `status` subcommands to existing feature command |
| `forge-cli/internal/cmd/feature_test.go` | Tests for new feature subcommands |
| `forge-cli/internal/cmd/root.go` | Register `proposalCmd` and `lessonCmd` |

## Acceptance Criteria

- [ ] `forge proposal` lists all proposals in table format: Slug | Created | Status | PRD | Feature
- [ ] `forge proposal <slug>` shows detail: metadata, content summary, linked artifacts, file path
- [ ] Created date reads from frontmatter `created` field, falls back to file birth time
- [ ] PRD column checks `docs/features/{slug}/prd/prd-spec.md` existence
- [ ] Feature column reads `docs/features/{slug}/manifest.md` status field
- [ ] `forge feature list` lists all features: Slug | Status | Progress | PRD(score) | Design(score) | UI(score) | Tests(score)
- [ ] Progress shows completed/total from `tasks/index.json`
- [ ] Scores read from frontmatter `score` field; show `—` when missing
- [ ] `forge feature status <slug>` shows manifest summary, task counts by status, artifacts with scores
- [ ] `forge lesson` lists all lessons: Name | Created | Tags | Category
- [ ] Category inferred from file prefix (gotcha-/arch-/pattern-/tool-/lesson-/hook-)
- [ ] `forge lesson <name>` shows metadata and file path (not full content)
- [ ] Test coverage ≥ 80% for new and modified code

## Hard Rules

- Output for list commands uses PrintBlock/PrintField format consistent with existing commands
- `forge feature` with no args keeps existing behavior (display current feature); `list` and `status` are new subcommands

## Implementation Notes

- Proposal discovery: walk `docs/proposals/*/proposal.md`, parse frontmatter for metadata
- Feature discovery: walk `docs/features/*/manifest.md`, parse frontmatter
- Lesson discovery: walk `docs/lessons/*.md`, parse frontmatter, infer category from filename prefix
- For `feature list` scores: read frontmatter of prd-spec.md, tech-design.md, ui-design.md to find `score` field
- Use existing `pkg/feature` constants for path construction
- Feature command extension: `forge feature` (no args) = show current; `forge feature <slug>` = set current; `forge feature list` = list all; `forge feature status <slug>` = detail
