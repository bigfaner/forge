---
created: 2026-05-20
author: "fanhuifeng"
status: Approved
---

# Proposal: Post-loop Artifact Auto-commit

## Problem

After `/run-tasks` completes all tasks and finishes knowledge extraction, uncommitted artifacts remain in the working tree. The user must manually run `/git-commit` to capture these leftovers.

### Evidence

Lesson `arch-post-loop-artifact-commit-gap` documents this gap.

**Concrete case (forge-init-config-sync feature)**: After 4 coding tasks + 1 drift detection task completed, `git status` showed:

```
Changes not staged for commit:
  modified:   docs/features/forge-init-config-sync/tasks/index.json

Untracked files:
  docs/features/forge-init-config-sync/tasks/fix-1.md
  docs/features/forge-init-config-sync/tasks/fix-2.md
  docs/lessons/gotcha-quick-tasks-stale-detect-command.md
```

Three artifact types left uncommitted: fix task definitions (fix-*.md), index.json status updates, and a knowledge entry (lesson). The user had to manually run `/git-commit` to capture these.

### Urgency

Medium. Not blocking, but every run-tasks invocation produces a manual cleanup step, reducing the promise of autonomous execution. The gap becomes more noticeable as fix task loops and knowledge extraction produce more artifacts.

## Proposed Solution

Add a "commit remaining artifacts" step at the end of the run-tasks post-completion flow (after knowledge extraction). This step unconditionally runs `git status`, filters to feature-scope files, and commits any detected artifacts. If nothing is uncommitted, it skips silently.

### Innovation Highlights

Straightforward adoption of the "final sweep" pattern — common in CI/CD pipelines and build systems. No creative insight needed; the gap is a missing step, not a design flaw.

## Requirements Analysis

### Key Scenarios

- **Happy path**: Knowledge extraction produces entries → user confirms → commit step detects and commits all feature-related artifacts
- **No knowledge extracted**: Knowledge extraction exits silently → commit step still detects any uncommitted index.json / manifest / records updates → commits if present
- **Nothing to commit**: All artifacts already committed by subagents → git status shows nothing → step skips silently
- **3 consecutive failures**: Loop stops without knowledge extraction → no commit step runs (matches existing behavior)

### Non-Functional Requirements

- **Safety**: Only commit files within `docs/features/<slug>/` and knowledge directories (docs/decisions/, docs/lessons/, docs/conventions/, docs/business-rules/). Never commit files outside these paths.
- **Idempotency**: Running the step multiple times is harmless — `git status` returns empty if nothing changed.

### Constraints & Dependencies

- Must follow Conventional Commits format (existing project convention)
- Must not interfere with the existing per-task commit flow inside task-executor subagents
- Only runs after knowledge extraction, not during the task loop

## Alternatives & Industry Benchmarking

### Industry Solutions

CI/CD pipelines commonly use "artifact sweep" or "final commit" steps after job completion. Git-based workflow tools (GitHub Actions, GitLab CI) routinely commit generated artifacts as a post-job step.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | No implementation cost | Every run requires manual cleanup | Rejected: defeats autonomous execution goal |
| Explicit path commit | Lesson suggestion | Simple, predictable | May miss artifacts if new types are added | Rejected: fragile to future changes |
| **git status detection + feature scope filter** | CI/CD sweep pattern | Catches all artifacts, future-proof | Slightly more complex logic | **Selected: comprehensive and safe** |

## Feasibility Assessment

### Technical Feasibility

Pure markdown edit to `plugins/forge/commands/run-tasks.md`. No code changes, no new CLI commands. The commit is performed by the agent following the skill instructions.

### Resource & Timeline

Single-file change. Estimated 1 coding task.

### Dependency Readiness

No external dependencies. All required tools (git, forge CLI) are already available.

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| Only knowledge extraction produces uncommitted artifacts | 5 Whys | Overturned: fix task records, index.json updates, and manifest changes can also remain uncommitted when subagent commits fail or when updates happen in the main session |
| Subagents always commit their artifacts successfully | Stress Test | Refined: subagents should commit, but edge cases (agent timeout, commit failure) can leave artifacts uncommitted |
| User wants to review artifacts before auto-commit | Assumption Flip | Refined: knowledge extraction already has user confirmation. Other artifacts (records, index.json) are machine-generated and safe to auto-commit |

## Scope

### In Scope

- Add "commit remaining artifacts" step to `plugins/forge/commands/run-tasks.md` Post-Completion section
- Step logic: `git status --porcelain` → filter to feature-scope paths → `git add` + `git commit` if artifacts exist → silent skip if not

### Out of Scope

- Changes to submit-task, git-commit, or task-executor agent
- Changes to knowledge extraction logic
- New CLI commands or scripts
- Changes to other pipelines (quick-tasks, breakdown-tasks)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Commits files the user didn't intend | L | M | Strict path filtering: only docs/features/<slug>/ and knowledge dirs |
| Commit fails due to pre-commit hook | L | L | Step is non-blocking; user can manually retry |
| git status picks up unrelated changes | L | L | Filter by feature slug path prefix |

## Success Criteria

- [ ] After run-tasks completes with knowledge extraction, running `git status` shows no uncommitted feature artifacts
- [ ] After run-tasks completes without knowledge extraction (silent exit), uncommitted index.json/manifest/records are still committed
- [ ] When no artifacts are uncommitted, the step produces no output and does not fail
- [ ] Commit message follows Conventional Commits format: `chore(<slug>): commit post-loop artifacts`

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
