---
id: "3"
title: "Update 7 eval-* commands to read config and pass to skill"
priority: "P1"
estimated_time: "1.5h"
complexity: "medium"
dependencies: [1]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 3: Update 7 eval-* commands to read config and pass to skill

## Description

Update all 7 eval-* commands (eval-proposal, eval-prd, eval-design, eval-ui, eval-journey, eval-contract, eval-consistency) to resolve target/iterations from config before invoking the eval skill. Each command calls `forge config get eval.<type>.target` and `eval.<type>.iterations` via Bash, falls back to omitting the arg if not configured, and passes resolved values to the Skill invocation. CLI arguments (`--target`/`--iterations`) take priority over config values.

## Reference Files
- `docs/proposals/eval-target-iterations-config/proposal.md` — Feasibility Assessment (eval-* command change example), Scope > In Scope (item 6)
- `docs/conventions/forge-distribution.md` — Plugin distribution model (required before modifying plugin files)
- `plugins/forge/commands/eval-proposal.md` — Current command pattern (Skill invocation)
- `plugins/forge/commands/eval-consistency.md` — Current command pattern (additional arg: --scope)

## Acceptance Criteria
- [ ] Each eval-* command reads `eval.<type>.target` and `eval.<type>.iterations` via `forge config get`, passing as `--target`/`--iterations` args when configured
- [ ] When config not set for a parameter, the corresponding `--target`/`--iterations` arg is omitted (eval skill uses rubric default)
- [ ] CLI `--target`/`--iterations` arguments take priority over config values — user-passed args are not overridden
- [ ] All 7 eval-* commands (proposal, prd, design, ui, journey, contract, consistency) follow the same config resolution pattern
- [ ] Existing `auto.eval.*` auto-run behavior unchanged (commands only change parameter resolution, not trigger logic)

## Implementation Notes
- Config resolution template per command:
  ```
  TARGET=$(forge config get eval.<type>.target 2>/dev/null)
  if [ $? -eq 0 ] && [ -n "$TARGET" ]; then TARGET_ARG="--target $TARGET"; fi
  ITERATIONS=$(forge config get eval.<type>.iterations 2>/dev/null)
  if [ $? -eq 0 ] && [ -n "$ITERATIONS" ]; then ITERATIONS_ARG="--iterations $ITERATIONS"; fi
  Skill(skill="forge:eval", args="--type <type> $TARGET_ARG $ITERATIONS_ARG [other args]")
  ```
- CLI arg priority: if user passes `--target 800` to the command, do NOT override with config value. The command template should check for existing user-provided args before falling back to config.
- eval-consistency has an additional `--scope` arg — preserve this alongside the config resolution logic.
- eval-ui may need special handling for platform detection (already handled by the eval skill, not the command).
