---
created: 2026-05-16
author: "faner + Claude"
status: Draft
---

# Proposal: Skill Rationalization — Trim Duplication, Preserve UX

## Problem

Forge carries 24 skills and 11 commands, but the maintenance surface is larger than necessary: 7 eval skills duplicate identical scorer→gate→revise orchestration (1,400 lines total), one skill is a dead alias, and one command is a rarely-used meta-tool.

### Evidence

- **Eval duplication**: `eval-proposal`, `eval-prd`, `eval-design`, `eval-ui`, `eval-test-cases`, `eval-consistency`, `eval-harness` each contain ~200 lines of SKILL.md. The scorer→gate→revise loop is copy-pasted across all 7. Any orchestration fix must be applied 7 times.
- **Dead alias**: `record-task` was functionally renamed to `submit-task` but the old skill directory remains (29 lines, just delegates).
- **Meta-tool**: `simplify-skill` command (66 lines) is a one-time refactoring utility with no ongoing use.

### Urgency

Each eval skill duplication means: bugs fixed in one copy remain in 6 others; context window carries 7 near-identical skill descriptions; new contributors must understand 7 files that do the same thing.

## Proposed Solution

**Merge 7 eval skills into 1 generic `eval` skill + rubric files**, surfaced through 7 thin command wrappers so user-facing slash commands (`/eval-proposal`, `/eval-prd`, etc.) remain unchanged.

**Remove dead weight**: `record-task` skill and `simplify-skill` command.

### Innovation Highlights

The "thin command wrapper over generic skill" pattern is a standard plugin architecture technique (used by VS Code extensions, Cargo subcommands). The key insight is separating *routing* (which eval type?) from *execution* (the scorer→gate→revise loop). This is straightforward decomposition, not novel — but it eliminates the maintenance trap of N-way duplication.

## Requirements Analysis

### Key Scenarios

- User types `/eval-prd` → command wrapper invokes `Skill("eval", "--type prd")` → generic skill loads `rubrics/prd.md` → scorer→gate→revise loop runs → identical behavior to today
- User types `/eval-ui` → same flow, loads `rubrics/ui-web.md` or `rubrics/ui-mobile.md` or `rubrics/ui-tui.md` based on manifest context
- User types `/eval-harness` → same flow, loads `rubrics/harness.md` (100-point scale instead of 1000-point)
- Existing pipeline invocations (e.g., brainstorm → eval-proposal, write-prd → eval-prd) work without changes since they call the same slash commands

### Non-Functional Requirements

- **Backward compatibility**: All existing `/eval-*` slash commands must work identically after migration
- **Context window**: System prompt should list fewer skills (24→17), reducing token consumption
- **Rubric discoverability**: Each rubric file is self-contained and can be audited independently

### Constraints & Dependencies

- Command wrappers must invoke `Skill("eval", args)` — this is the standard command-to-skill delegation pattern already used by `/quick` and other commands
- The `eval` skill must auto-detect UI platform (web/mobile/tui) for `/eval-ui` — this logic currently lives in `eval-ui/SKILL.md`
- `eval-harness` uses a 100-point scale (not 1000-point) and single-pass scoring (no reviser loop) — the generic skill must handle both modes

## Alternatives & Industry Benchmarking

### Industry Solutions

Monorepo tooling systems (Bazel, Nx, Turborepo) solve similar problems with "target pattern + configuration" rather than "one binary per target". The eval consolidation applies the same principle.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero migration cost | 7× duplicated logic persists; context bloat remains | Rejected: pain points are real and growing |
| Auto-detect eval type | — | Single `/eval` command, simplest UX | User loses explicit control; auto-detection can guess wrong | Rejected: user preferred explicit commands |
| **Keep commands, DRY impl** | VS Code extension pattern | UX unchanged; maintenance halved; rubric files auditable | Slightly more files (7 wrappers); one-time migration effort | **Selected: best UX/maintenance tradeoff** |
| Keep all separate | — | No coupling | 7× duplication; bug propagation risk | Rejected: status quo is the problem |

## Feasibility Assessment

### Technical Feasibility

High. The 7 eval skills already share identical agent invocations (`doc-scorer`, `doc-reviser`). The generic skill extracts the shared loop and parameterizes only the rubric path, target score, and max iterations.

### Resource & Timeline

One migration session. Each eval skill has a clear 1:1 mapping to a rubric file + command wrapper. No design ambiguity.

### Dependency Readiness

No external dependencies. All logic exists in current skill files. The command-to-skill invocation pattern is already used by `commands/quick.md`.

## Scope

### In Scope

1. **Merge eval skills**: Create `skills/eval/SKILL.md` with generic scorer→gate→revise loop
2. **Extract rubrics**: Create `skills/eval/rubrics/` with 8 rubric files (proposal, prd, design, ui-web, ui-mobile, ui-tui, test-cases, consistency, harness)
3. **Create command wrappers**: Create 7 thin command files (`commands/eval-proposal.md` through `commands/eval-harness.md`) that invoke `Skill("eval", "--type <type>")`
4. **Remove eval skill directories**: Delete `skills/eval-proposal/` through `skills/eval-harness/` (7 directories)
5. **Remove record-task skill**: Delete `skills/record-task/` (superseded by `submit-task`)
6. **Remove simplify-skill command**: Delete `commands/simplify-skill.md` (rarely-used meta-tool)
7. **Update forge guide**: Update CLAUDE.md system prompt section to reflect new eval skill structure

### Out of Scope

- Core pipeline skills (brainstorm, write-prd, ui-design, tech-design, breakdown-tasks, quick-tasks)
- Test pipeline skills (gen-test-cases, gen-test-scripts, run-e2e-tests, graduate-tests)
- Utility skills (consolidate-specs, init-justfile, extract-design-md, improve-harness, forensic, learn-lesson, submit-task)
- gen-sitemap (stays as command)
- CLI binary changes (forge command structure unchanged)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Command wrapper invocation syntax differs from current skill invocation | L | H | Verify command-to-skill delegation pattern works with `args` parameter before migration |
| eval-harness 100-point scale breaks when merged into 1000-point generic | M | M | Generic skill reads `scale` field from rubric frontmatter; harness rubric declares `scale: 100` |
| eval-ui multi-platform rubric selection lost | M | M | Command wrapper passes `--type ui` and generic skill resolves platform from manifest/config |
| Pipeline invocations break if slash command name changes | L | H | Slash commands names stay identical — only implementation moves from skill to command wrapper |

## Success Criteria

- [ ] All 7 `/eval-*` slash commands produce identical behavior to pre-migration
- [ ] `skills/eval/` contains exactly 1 SKILL.md + rubric directory (no eval-specific orchestration outside generic skill)
- [ ] `skills/record-task/` no longer exists
- [ ] `commands/simplify-skill.md` no longer exists
- [ ] Total SKILL.md lines for eval reduced from ~1,400 to ~200 (generic) + ~700 (rubrics)
- [ ] `eval-forge` runtime audit passes with new structure
- [ ] System prompt skill count: 24 → 17

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks directly from this proposal
