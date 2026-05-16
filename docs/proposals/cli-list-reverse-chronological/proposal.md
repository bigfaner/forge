---
created: 2026-05-16
author: faner
status: Draft
---

# Proposal: CLI List Commands — Reverse Chronological Sorting

## Problem

`forge proposal` and `forge feature list` display entries in filesystem (lexical) order, making recent items harder to find.

### Evidence

Both commands iterate `os.ReadDir()` results directly with no sorting. On typical filesystems this produces alphabetical order by slug name, unrelated to creation date. Users must scan the entire list to find recently added proposals or features.

### Urgency

Low urgency but high annoyance. Every invocation of these listing commands forces the user to visually scan for the newest entries.

## Proposed Solution

Sort both listing commands by date in descending order (newest first). Proposals use their existing `created` frontmatter field. Features use `manifest.md` modification time.

### Innovation Highlights

Straightforward improvement — no novel approach needed.

## Requirements Analysis

### Key Scenarios

- User runs `forge proposal` → sees most recently created proposals at the top
- User runs `forge feature list` → sees most recently active features at the top
- Proposals without `created` frontmatter fall back to file modification time (already implemented in `proposal.Discover()`)

### Non-Functional Requirements

- Sorting overhead is negligible (small in-memory slices)

### Constraints & Dependencies

- Proposals already expose `Created` field via `proposal.Discover()`
- Features need `manifest.md` stat during discovery — `discoverFeatures()` already reads manifest content

## Alternatives & Industry Benchmarking

### Industry Solutions

Most CLI tools with listing commands default to reverse chronological order (e.g., `git log`, `ls -lt`, `gh issue list`).

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | No code change | Poor UX, inconsistent with conventions | Rejected |
| Sort by slug (alphabetical) | — | Deterministic | Doesn't solve the problem | Rejected |
| **Sort by date descending** | Industry standard | Intuitive, matches user expectation | None | **Selected** |

## Feasibility Assessment

### Technical Feasibility

Trivial — add `sort.Slice()` calls after discovery in both commands. All required date data is already available.

### Resource & Timeline

Single developer, under 1 hour.

### Dependency Readiness

No external dependencies.

## Scope

### In Scope

- Sort `forge proposal` output by `created` date (descending), fallback to mtime
- Sort `forge feature list` output by `manifest.md` mtime (descending)
- Update existing tests to verify sort order

### Out of Scope

- Adding `created` field to feature manifests (could be a separate proposal)
- Configurable sort order flags
- Other `forge` subcommands

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| `manifest.mtime` doesn't reflect creation but last status update | M | L | Acceptable — "most recently active" is equally useful for features |

## Success Criteria

- [ ] `forge proposal` lists proposals newest-first by `created` date
- [ ] `forge feature list` lists features newest-first by manifest mtime
- [ ] Existing tests pass, new tests verify sort order

## Next Steps

- Proceed to `/quick-tasks` for task generation
