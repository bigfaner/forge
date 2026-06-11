---
date: "2026-05-29"
feature: intent-driven-pipeline-branching
session: ae1022de-6363-4077-be27-ab96dab4180f/subagents/agent-a43d0562633565f25
trigger: "T-review-doc task executor stuck on read operations, user-interrupted after ~42s"
---

# Forensic Report: T-review-doc Agent Stuck on Reads

## Timing Breakdown

| Action | Duration | Detail |
|--------|----------|--------|
| Thinking | 46.5s | Across multiple turns, agent reasoning without clear direction |
| Read | 0.6s | 3 reads (review-doc.md × 2, proposal.md, manifest.md) |
| Bash | 2.4s | 6 commands (find, ls, forge prompt) |
| TaskCreate | 0.1s | Internal task tracking |
| TaskUpdate | 0.02s | Set status in_progress |

Total tool time: 3.1s / Session: 42s / **Thinking:tool ratio: 15:1**

## Symptom

Agent dispatched for T-review-doc, executed 11 tool calls in 42 seconds, then was interrupted by user. User perceived agent as "stuck on read operations" — reading files without making progress on actual review.

## Root Cause: Discovery Strategy Mismatch

**Classification: context-starvation (template limitation)**

The `doc-review.md` embed template (`forge-cli/pkg/task/templates/doc-review.md`) hardcodes a Discovery Strategy that only scans two directories:

```
- docs/features/<slug>/ (prd/, design/, testing/)
- docs/proposals/<slug>/
```

But for features whose doc tasks modify files **outside** these directories, the deliverables are invisible to the review agent. In this case:

| Task | Actual Deliverable | In Discovery Scope? |
|------|-------------------|-------------------|
| Task 1 | `plugins/forge/skills/brainstorm/templates/proposal.md` | **No** |
| Task 1 | `plugins/forge/skills/brainstorm/SKILL.md` | **No** |
| Task 2 | `plugins/forge/skills/write-prd/SKILL.md` | **No** |
| Task 3 | `plugins/forge/skills/tech-design/SKILL.md` | **No** |

The agent found **zero deliverables** in the allowed scan paths (feature dir had only empty prd/design/ui dirs, manifest.md was excluded). With nothing to review, the agent entered aimless exploration:
1. Read proposal.md (not a review target)
2. Read manifest.md (explicitly excluded by template)
3. Then interrupted

## Causal Chain

1. **Symptom**: Agent stuck reading irrelevant files, no progress on review
2. **Direct cause**: Agent found zero deliverables in Discovery Strategy's allowlist, defaulted to reading anything available
3. **Root cause**: `doc-review.md` template hardcodes `docs/features/` and `docs/proposals/` as the only scan targets — no mechanism to include deliverable paths from doc tasks' `Affected Files` sections

## Secondary Issue: Wrong Initial Path

Agent's first read attempted `forge-cli/docs/features/.../review-doc.md` — a non-existent path. This added one wasted read + one find command before locating the correct path. Minor impact (2s), but indicates the agent lacked clear project root context.

## Fix Recommendations

### Option A: Enrich Discovery Strategy with Affected Files (Recommended)

In `build.go`, when generating the review-doc task body, inject the actual file paths from doc tasks' `## Affected Files` sections into the template. Add a new template variable:

```
{{.DocTaskFiles}}
```

Which would contain:
```
- plugins/forge/skills/brainstorm/templates/proposal.md
- plugins/forge/skills/brainstorm/SKILL.md
- plugins/forge/skills/write-prd/SKILL.md
- plugins/forge/skills/tech-design/SKILL.md
```

This ensures the review agent always knows where the actual deliverables are.

### Option B: Expand Discovery Scope

Add `plugins/forge/skills/` to the allowlist. Simpler but less precise — may cause the review agent to read unrelated skill files.

### Option C: Include Source Task Files

Add a `{{.DocTaskPaths}}` variable that lists the task .md file paths, so the agent can read each task's `## Affected Files` to discover deliverables. Two-step discovery but more generic.
