---
id: "3"
title: "Create shared knowledge extraction routine"
priority: "P0"
estimated_time: "2h"
dependencies: []
type: "feature"
mainSession: false
---

# 3: Create shared knowledge extraction routine

## Description

Create a reusable prompt section (`knowledge-extraction.md`) that all 4 auto-extract triggers will include. This routine defines the knowledge identification, extraction, and summarization logic used when a pipeline step completes. It reads the feature's artifacts (PRD, tech-design, task outcomes), identifies notable knowledge, and presents it for user confirmation.

## Reference Files
- `docs/proposals/knowledge-accumulation-loop/proposal.md` — Source proposal (Part 2: Auto-Extract Triggers)
- `plugins/forge/references/shared/decision-logging.md` — Decision format for extracted decisions
- `plugins/forge/skills/learn-lesson/templates/template.md` — Lesson format for extracted lessons

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/references/shared/knowledge-extraction.md` | Shared extraction routine for auto-extract triggers |

### Modify
| File | Changes |
|------|---------|
| (none) | |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] File defines a reusable prompt section that can be included by any trigger point
- [ ] Defines the extraction flow: scan artifacts → identify knowledge → extract & summarize → present for confirmation → write on confirm
- [ ] Knowledge identification covers 4 types: decisions, lessons, conventions, business rules
- [ ] Defines heuristics for "notable knowledge" vs "routine changes" (to achieve <30% false-positive rate)
- [ ] Silent when no notable knowledge detected — produces no output
- [ ] Extracted knowledge presented for user confirmation before writing (AskUserQuestion)
- [ ] Reuses same file formats as `/learn` skill (decision rows, lesson template, convention/business-rule entries)
- [ ] Uses auto-generated vocabulary (from consolidate-specs) when available for classification suggestions
- [ ] Parameterizable by trigger context: what artifacts to scan (varies by trigger point)
- [ ] Defines the artifact scanning scope per trigger type:
  - run-tasks: task outcomes, code changes, manifest
  - fix-bug: root cause analysis, fix approach
  - write-prd: PRD content
  - tech-design: design document

## Hard Rules
- This is a shared reference file — trigger points include it, not copy-paste it
- Extraction logic must be conservative: only extract genuinely non-obvious knowledge
- Must not write to knowledge directories without explicit user confirmation
- Format compatibility with `/learn` skill output and `/consolidate-specs` overlap detection

## Implementation Notes
- The "notable knowledge" heuristics are the key design challenge. Consider:
  - Decisions: non-obvious choices where alternatives existed (not "used standard library")
  - Lessons: root causes that weren't immediately apparent (not "typo in variable name")
  - Conventions: patterns that should be repeated (not one-off choices)
  - Business rules: constraints that apply across features (not feature-specific logic)
- The routine should be a self-contained prompt section with clear sections: scanning, identification, extraction, confirmation, writing
