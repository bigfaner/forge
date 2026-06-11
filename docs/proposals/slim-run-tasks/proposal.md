---
created: 2026-05-16
author: faner
status: Draft
---

# Proposal: Slim Down run-tasks for Token Efficiency

**Summary**: Compress the `run-tasks.md` skill document (~250 lines) by ~40% and eliminate verbose test output from context, reducing per-iteration token waste without changing CLI behavior.

## Problem

`run-tasks.md` (~250 lines) is loaded into the main session context every iteration, consuming tokens proportional to its verbosity. The Breaking Gate and E2E Gate also pipe full test output (potentially hundreds of lines) into context. Combined, this wastes ~40-60% of per-iteration tokens on content that doesn't drive decisions.

### Evidence

- `run-tasks.md` has 250 lines, much of it verbose explanations and code examples that the AI doesn't need after the first iteration.
- `forge task query` outputs 4 fields (TASK_ID, STATUS, SCOPE, BREAKING) when only STATUS is needed for verification — `forge task status` outputs just 2 fields.
- `just test` output can be 100+ lines of test runner output, all flowing into the main session context even on success.
- **Estimated savings**: ~100 lines of skill text (~40% reduction) + ~100 lines of eliminated test output per gate execution, compounding across iterations.

### Urgency

Token cost scales linearly with iteration count. A 10-task feature run with the current flow wastes significant tokens. Quick fix, high payback.

## Proposed Solution

Three-pronged optimization of `run-tasks.md` only (no CLI changes):

1. **Replace `forge task query` → `forge task status`** in Step 2b — outputs only TASK_ID + STATUS
2. **Silent test execution** — redirect Breaking Gate and E2E Gate output to file, only check exit code; tail key lines on failure
3. **Compress skill document** — keep mermaid diagram and structure, remove verbose explanations and redundant examples; target ~250 → ~150 lines

## Requirements Analysis

### Key Scenarios

- Happy path: task completes, tests pass, minimal output
- Test failure: fix task created with enough context to diagnose (tail output file for key error lines)
- Record missing: recovery subagent spawned (unchanged)

### Constraints & Dependencies

- Only changes `plugins/forge/commands/run-tasks.md` — no CLI code changes
- Must preserve all existing functionality (claim, dispatch, verify, fix-task, main-session routing)

### References

- `plugins/forge/commands/run-tasks.md` — target file (current ~250 lines)
- `forge-cli/internal/cmd/status.go` — `forge task status` command (2-field output)
- `forge-cli/internal/cmd/query.go` — `forge task query` command (4-field output, being replaced in skill)

## Alternatives & Industry Benchmarking

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | No risk | Token waste continues | Rejected: clear ROI |
| Only replace query→status | Minimal change | Misses biggest savings (test output + doc size) | Rejected: incomplete |
| **Full slim-down** | Maximum token savings | Requires careful rewrite to preserve semantics | **Selected** |

## Scope

### In Scope

- Slim `run-tasks.md` text (keep mermaid diagram, remove verbose explanations)
- Replace `forge task query` → `forge task status` in Step 2b
- Silent Breaking Gate (Step 3a) — redirect test output to file, check exit code
- Silent E2E Gate (Step 3b) — same treatment
- Condense Error Handling table

### Out of Scope

- Changes to CLI commands (forge task status, query, etc.)
- Changes to task-executor agent or other skills
- Changes to execute-task.md (separate skill)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Silent failures harder to debug | L | M | On failure, tail last N lines from output file before creating fix-task |
| Over-condensed instructions lose nuance | M | M | Keep mermaid diagram; preserve EXTREMELY-IMPORTANT blocks verbatim |
| Edge case lost in compression | L | L | Preserve all error handling paths, just compress the prose |

## Success Criteria

- [ ] `run-tasks.md` reduced from ~250 to ~150 lines
- [ ] Step 2b uses `forge task status` instead of `forge task query`
- [ ] Breaking Gate and E2E Gate produce zero output lines on success (exit code only)
- [ ] On test failure, output file is tailed for diagnostic context in fix-task
- [ ] All existing functionality preserved: claim, dispatch, verify, fix-task, main-session routing

## Next Steps

- Proceed to `/quick-tasks` to generate tasks from this proposal
