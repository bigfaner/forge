---
created: 2026-05-19
author: faner
status: Draft
---

# Proposal: Inline Expert Files into Eval Skill

## Problem

Expert files under `agents/experts/` (9 scorer personas + 2 protocol files) are auto-discovered by Claude Code's plugin system and registered as standalone `forge:experts:*` agents. These files are not independent agents — they are scoring templates read by the eval skill at runtime to compose prompts for `general-purpose` agents. Exposing 11 phantom agents pollutes the `/agents` panel and confuses users.

Additionally, the eval SKILL.md has several logical inconsistencies and sequencing gaps that should be fixed alongside the file move.

### Evidence

- Running `/agents` shows `forge:experts:cto`, `forge:experts:pm`, etc. — none of which are meant to be invoked directly
- The expert files contain role descriptions + domain failure patterns, not agent execution instructions
- The original `expert-template-eval` proposal placed them in `agents/` without considering the plugin auto-discovery behavior

### Urgency

Users see 11 spurious agents every time they open the agents panel. Low severity but degrades plugin UX. The eval mechanism issues are low urgency but cause confusion when they manifest (especially the multi-expert reviser path ambiguity).

## Proposed Solution

### Part A: Move Expert Files

Move `agents/experts/` → `skills/eval/experts/`. Since expert files now live within the eval skill's own directory, SKILL.md references them via relative paths (same convention as `rubrics/<type>.md`) — no `${CLAUDE_SKILL_DIR}` needed.

#### Subagent Dispatch Mechanism (unchanged)

**Scorer dispatch** (per iteration, per expert):
1. Read protocol file → `experts/protocol/scorer-protocol.md`
2. Read expert file → `experts/scorer/<expert>.md`
3. Read rubric → `rubrics/<type>.md`
4. Compose prompt: protocol content (variables replaced) + expert content + context injection
5. Spawn `general-purpose` agent via Agent tool (`model: "sonnet"`) with composed prompt

Multi-expert types (e.g., `prd` → `[pm, qa]`): spawn in parallel, each with its own composed prompt.

**Reviser dispatch** (when score < target):
1. Read protocol file → `experts/protocol/reviser-protocol.md`
2. Compose prompt: protocol content (variables replaced) + merged attack points + context injection
3. Spawn `general-purpose` agent via Agent tool (`model: "sonnet"`) with composed prompt

### Part B: Fix Eval Mechanism Issues

Five issues found in SKILL.md logic:

#### Fix 1: Multi-expert reviser EVAL_REPORT_PATH ambiguity

**Problem**: Step 4.1 passes `EVAL_REPORT_PATH` to the reviser, but for multi-expert types (e.g., `prd` → `[pm, qa]`) there are multiple reports (`iteration-N-pm.md`, `iteration-N-qa.md`). No single report path exists.

**Fix**: After Step 2.3's LLM merge, the main session writes a merged report to `<doc_dir>/eval/iteration-{{N}}-merged.md` (containing merged attacks + averaged scores). This merged report serves as `EVAL_REPORT_PATH` for the reviser. Single-expert types continue using `iteration-{{N}}.md` directly.

#### Fix 2: Iteration counter initialization

**Problem**: Step 4.2 says "Increment iteration counter, return to Step 2" but Step 1 never initializes `ITERATION = 1`.

**Fix**: Add explicit initialization after Step 1: "Set `ITERATION = 1`, `MAX_ITERATIONS = resolved value from rubric or CLI`."

#### Fix 3: Remove ambiguous "continue" override in gate

**Problem**: Step 3b says "On 'continue'/'keep going': run scorer again (Step 2), then re-evaluate this gate." This adds an unclear manual override to an automated gate decision. The gate should be purely score-driven.

**Fix**: Remove the "continue/keep going" line from Step 3b. If the user wants additional iterations after reaching the target, they re-run `/eval`.

#### Fix 4: Context injection for reviser

**Problem**: Step 2.1 injects `CONTEXT_CONTENT` (project conventions, business rules) into every scorer prompt, but Step 4.1 does not inject it into the reviser. If attacks reference convention violations, the reviser lacks context to fix them properly.

**Fix**: Apply the same context injection block from Step 2.1 to Step 4.1. Append `<injected-context>...</injected-context>` to the reviser prompt when `CONTEXT_CONTENT` was loaded.

#### Fix 5: Score extraction robustness

**Problem**: Step 2.3 says "extract directly" the `SCORE: X/{{scale}}` format. But `general-purpose` agents may add preamble/postamble text, causing extraction to fail.

**Fix**: Add extraction instruction to Step 2.3: "Extract score using regex `/SCORE:\s*(\d+)\/(\d+)/`. If pattern not found, scan the scorer agent's output for the last line matching a `number/number` pattern. If still not found, report error and stop."

### File Changes

```
# Move
agents/experts/protocol/scorer-protocol.md  →  skills/eval/experts/protocol/scorer-protocol.md
agents/experts/protocol/reviser-protocol.md →  skills/eval/experts/protocol/reviser-protocol.md
agents/experts/scorer/architect.md          →  skills/eval/experts/scorer/architect.md
agents/experts/scorer/code-reviewer.md      →  skills/eval/experts/scorer/code-reviewer.md
agents/experts/scorer/cto.md                →  skills/eval/experts/scorer/cto.md
agents/experts/scorer/editor.md             →  skills/eval/experts/scorer/editor.md
agents/experts/scorer/harness-engineer.md   →  skills/eval/experts/scorer/harness-engineer.md
agents/experts/scorer/pm.md                 →  skills/eval/experts/scorer/pm.md
agents/experts/scorer/qa.md                 →  skills/eval/experts/scorer/qa.md
agents/experts/scorer/ux-auditor.md         →  skills/eval/experts/scorer/ux-auditor.md
agents/experts/scorer/ux-engineer.md        →  skills/eval/experts/scorer/ux-engineer.md

# Edit
skills/eval/SKILL.md  — path updates (Part A) + mechanism fixes (Part B)

# Delete
agents/experts/ (entire directory)

# Update
docs/conventions/forge-distribution.md  — directory tree and component descriptions
```

### Path Updates in SKILL.md (Part A)

| Location | Old | New |
|----------|-----|-----|
| Dispatch table header | `${CLAUDE_SKILL_DIR}/../../agents/experts/scorer/` | `experts/scorer/` |
| Step 2.1 (scorer protocol) | `${CLAUDE_SKILL_DIR}/../../agents/experts/protocol/scorer-protocol.md` | `experts/protocol/scorer-protocol.md` |
| Step 2.1 (expert file example) | `${CLAUDE_SKILL_DIR}/../../agents/experts/scorer/pm.md` | `experts/scorer/pm.md` |
| Step 4.1 (reviser protocol) | `${CLAUDE_SKILL_DIR}/../../agents/experts/protocol/reviser-protocol.md` | `experts/protocol/reviser-protocol.md` |

### Mechanism Fixes in SKILL.md (Part B)

| Fix | Location | Change |
|-----|----------|--------|
| Fix 1 | Step 2.3 + Step 4.1 | Multi-expert: write merged report, use as EVAL_REPORT_PATH |
| Fix 2 | After Step 1 | Add "Set ITERATION = 1, MAX_ITERATIONS = resolved value" |
| Fix 3 | Step 3b | Remove "On continue/keep going" line |
| Fix 4 | Step 4.1 | Add context injection block (same as Step 2.1) |
| Fix 5 | Step 2.3 | Add regex extraction with fallback |

### Distribution Doc Updates

`docs/conventions/forge-distribution.md`:
- Directory tree: move `agents/experts/` subtree under `skills/eval/`, remove from `agents/`
- Component table: remove `agents/experts/` from agents row, add `experts/` under skills row description
- Section "3. 核心依赖 → agents/experts/": update paths to reflect new location

## Alternatives

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| **Move + fix all** | Clean state; no deferred issues | Slightly larger scope | **Selected** |
| Move only, log issues separately | Minimal scope | Issues persist until separate fix | Rejected: cheap to fix now |
| Inline into SKILL.md | Zero extra files | SKILL.md already 354 lines; would exceed 500+ | Rejected: readability |

## Scope

### In Scope

- Move 11 expert/protocol files from `agents/experts/` to `skills/eval/experts/`
- Update 4 path references in `skills/eval/SKILL.md` (relative paths, no `${CLAUDE_SKILL_DIR}`)
- Fix 5 eval mechanism issues in `skills/eval/SKILL.md`
- Update `docs/conventions/forge-distribution.md`
- Delete `agents/experts/` directory

### Out of Scope

- Changing scorer or reviser protocol content
- Changing dispatch table expert assignments
- Adding new expert types
- Changes to non-eval skills or commands

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Cached plugin still reads old path | L | L | User runs `/plugin` to update cache |
| Merged report format breaks reviser | L | M | Merged report follows same structure as scorer output |
| Context injection increases reviser token cost | H | L | Only activates when rubric declares `context` — same as scorer |

## Success Criteria

- [ ] `agents/experts/` directory no longer exists
- [ ] No `forge:experts:*` agents appear in `/agents` panel
- [ ] All 11 files exist under `skills/eval/experts/`
- [ ] Multi-expert reviser receives a single merged report as EVAL_REPORT_PATH
- [ ] Iteration counter explicitly initialized before first scoring
- [ ] Reviser prompt includes context injection when rubric declares `context`
- [ ] Score extraction has regex fallback for malformed output
- [ ] No "continue/keep going" override in automated gate
