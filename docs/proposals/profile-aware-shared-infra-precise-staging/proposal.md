---
created: 2026-05-17
author: faner
status: Draft
---

# Proposal: Profile-Aware Shared Infrastructure + Precise Git Staging

## Problem

gen-test-scripts SKILL.md Step 3.5 auth section hardcodes `.ts` filenames (`helpers.ts`) that mislead AI agents into creating/referencing TypeScript files even when the active profile is `go-test`. Separately, git-commit.md lacks an explicit prohibition on `git add -A`, allowing agents to stage unrelated files.

### Evidence

- SKILL.md line 363: "Verify `helpers.ts` exports all auth-related symbols" — hardcodes `.ts` regardless of profile
- SKILL.md line 502: "verify `helpers.ts` exports cover all imports" — same issue
- SKILL.md line 180: "auth-setup.ts, playwright.config.ts, helpers.ts" — descriptive but reinforces `.ts` assumption
- git-commit.md: template says `git add <changed-files>` but never forbids `git add -A` / `git add .` / `git add --all`
- Historical incident: fix-2 task (d7f8a13) committed 169 files instead of 2, because Step 3.5 regenerated .ts residue + agent used `git add -A`

### Urgency

The `.ts` hardcoded references will cause the same residue problem every time gen-test-scripts runs with a non-playwright profile. The `git add -A` gap means any agent that can't determine its exact file list will stage everything — a recurring risk during fix tasks.

## Proposed Solution

1. Replace all hardcoded `.ts` filename references in SKILL.md Step 3.5 with profile-manifest-derived references (e.g., "the profile's helpers file from `manifest.templates.helpers`")
2. Add a `HARD-RULE` to git-commit.md explicitly forbidding `git add -A`, `git add .`, and `git add --all`, requiring explicit file paths

### Innovation Highlights

Straightforward fix — no innovation needed. The profile system already provides the correct file names; the issue is that a few lines in the SKILL.md bypass it with hardcoded values.

## Requirements Analysis

### Key Scenarios

- go-test profile runs gen-test-scripts → Step 3.5 should reference `helpers.go`, never `helpers.ts`
- web-playwright profile runs gen-test-scripts → Step 3.5 should reference `helpers.ts` as before
- error-fixer commits a 2-file fix → commit contains exactly those 2 files
- task-executor completes a task → commit contains only task-related files

### Non-Functional Requirements

- Backward compatible: web-playwright profile behavior unchanged
- No runtime code changes — only SKILL.md instruction text and git-commit.md rules

### Constraints & Dependencies

- Profile manifest `templates.helpers` field must exist (already true for all 6 profiles)
- git-commit.md is the single chokepoint for all agent commits

## Alternatives & Industry Benchmarking

### Industry Solutions

Precise staging is standard practice in CI/CD. Many tools enforce it via pre-commit hooks or commit linting.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero cost | Recurring .ts residue + staging amplification | Rejected: proven incident at d7f8a13 |
| Fix SKILL.md + git-commit.md | Internal lesson | Targeted, single-point fix | Relies on agents reading rules | **Selected: minimal scope, high impact** |
| Fix all agents individually | Defense in depth | Redundant enforcement | 4+ files to maintain, rules drift apart | Rejected: over-engineering for this scope |

## Feasibility Assessment

### Technical Feasibility

2 markdown files, ~15 lines changed. No code, no tests, no build impact.

### Resource & Timeline

Single session.

### Dependency Readiness

All profile manifests already define `templates.helpers`. git-commit.md is the sole commit instruction file.

## Scope

### In Scope

- SKILL.md: Replace 3 hardcoded `.ts` references with profile-manifest-derived references (lines 180, 363, 502)
- git-commit.md: Add `HARD-RULE` prohibiting `git add -A` / `git add .` / `git add --all`

### Out of Scope

- Cleanup of existing .ts residue files (manual git operation)
- Changes to agent definitions (error-fixer.md, task-executor.md)
- Changes to other SKILL.md steps beyond Step 3.5 auth section and post-generation helper merge

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Agent ignores the new HARD-RULE and uses git add -A anyway | L | M | The rule is explicit; if agent still violates, it's an LLM compliance issue beyond document-level fixes |
| Profile manifest missing `templates.helpers` field | L | L | All 6 existing profiles define this field; new profiles follow the authoring guide |

## Success Criteria

- [ ] SKILL.md contains zero hardcoded `.ts` filename references in Step 3.5 auth section and post-generation helper merge
- [ ] git-commit.md contains an explicit `HARD-RULE` forbidding `git add -A`, `git add .`, `git add --all`
- [ ] web-playwright profile behavior is unchanged (helpers.ts still generated)

## Next Steps

- Proceed to `/quick-tasks` to generate tasks
