---
created: 2026-05-15
author: "fanhuifeng"
status: Draft
---

# Proposal: Document Docs-Only Fast Path in Skill Docs

## Problem

Skills that generate tasks (`quick-tasks`, `breakdown-tasks`) and the global guide (`guide.md`) do not document that docs-only features skip profile resolution, test task generation, and quality gates. An agent executing these skills for a docs-only feature cannot determine the skip behavior from reading the skill docs alone.

### Evidence

- `quick-tasks/SKILL.md` Step 0 (lines 20-32): mandates `forge profile` with HARD-RULE "Do NOT silently default to any profile" — no mention that docs-only features skip this entirely
- `breakdown-tasks/SKILL.md` Step 0 (lines 25-35): same mandate, same omission
- `guide.md` Quality Gate Protocol (lines 109-113): states "All task-executing workflows MUST pass the quality gate" — no exception for documentation tasks
- `guide.md` All-Completed Hook (lines 126-133): describes full test pipeline — no mention that `forge quality-gate` (the actual CLI) already skips docs-only features
- `guide.md` Quick mode differences (line 99): mentions "Docs-only features auto-detected: no test tasks, generates T-eval-doc instead" — but this is the only reference, and it's buried in a bullet

### Urgency

Low urgency — the runtime already handles docs-only correctly (`noTest: true`, `forge quality-gate` skips docs-only features, `BREAKING` not set). The gap is purely documentation: agents waste time on `forge profile` and test steps they don't need, and developers reading the docs get a misleading picture of the workflow.

## Proposed Solution

Add a `## Docs-Only Fast Path` section at the top of `quick-tasks/SKILL.md` and `breakdown-tasks/SKILL.md`, listing which steps to skip when all business tasks are documentation type. Update `guide.md` Quality Gate Protocol and All-Completed Hook sections with explicit docs-only exceptions.

### Innovation Highlights

Straightforward documentation alignment — no new patterns. The innovation is making implicit runtime behavior explicit in the agent-facing docs, so agents can self-correct without exploring the codebase.

## Requirements Analysis

### Key Scenarios

- Agent runs `/quick-tasks` for a docs-only feature: reads the fast path section, skips Step 0 (profile) and Step 4 (test tasks), proceeds to generate documentation-type tasks
- Agent runs `/breakdown-tasks` for a docs-only feature: reads the fast path section, skips Step 0 and Step 4b (test tasks)
- Agent reads `guide.md` during task execution: sees explicit docs-only exception in Quality Gate Protocol, knows `noTest: true` tasks skip the gate
- Developer reads skill docs: understands the full pipeline and when steps are optional

### Constraints & Dependencies

- No runtime code changes — this is documentation only
- Must stay in sync with actual runtime behavior (which is already correct)

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero cost | Agent confusion, wasted steps on docs-only features | Rejected: doc gap causes real agent inefficiency |
| Per-step conditional notes | — | Minimal change, stays near relevant steps | Scattered, easy to miss one | Rejected: user explicitly chose centralized section |
| **Centralized fast path section** | — | Agent reads once at top, knows full skip list; clear separation | Must keep in sync when steps change | **Selected: clearer for agents** |

## Feasibility Assessment

### Technical Feasibility

Pure markdown changes. No code, no build, no tests.

### Resource & Timeline

Single pass. 3 files, ~30 lines of additions total.

## Scope

### In Scope

- `plugins/forge/skills/quick-tasks/SKILL.md` — add `## Docs-Only Fast Path` section
- `plugins/forge/skills/breakdown-tasks/SKILL.md` — add `## Docs-Only Fast Path` section
- `plugins/forge/hooks/guide.md` — add docs-only exceptions to Quality Gate Protocol and All-Completed Hook sections

### Out of Scope

- Other profile-referencing skills (gen-test-cases, gen-test-scripts, etc.) — docs-only features never reach these
- Runtime code changes
- Task templates (task.md, task-doc.md)
- The existing `simplify-skill-task-docs` proposal (different scope)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Fast path section goes stale when steps change | M | L | Section references specific step numbers — easy to verify |
| Agent ignores fast path section and runs all steps anyway | L | L | Runtime already handles correctly; doc is self-correcting guidance |

## Success Criteria

- [ ] `quick-tasks/SKILL.md` has a `## Docs-Only Fast Path` section that lists Step 0 and Step 4 as skippable for docs-only features
- [ ] `breakdown-tasks/SKILL.md` has a `## Docs-Only Fast Path` section that lists Step 0 and Step 4b as skippable for docs-only features
- [ ] `guide.md` Quality Gate Protocol explicitly states that documentation tasks (`noTest: true`) skip the quality gate
- [ ] `guide.md` All-Completed Hook explicitly states that `forge quality-gate` skips docs-only features
- [ ] An agent reading only these 3 files can determine the complete docs-only workflow without exploring the codebase
