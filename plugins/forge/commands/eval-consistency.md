---
name: eval-consistency
description: Evaluate and fix cross-document consistency (PRD, Design, UI, Tasks). Detects inconsistencies via 1000-point scoring, then auto-fixes downstream docs to align with PRD as source of truth. Supports --scope docs|full.
argument-hint: "[--target 900] [--iterations 3] [--scope docs|full]"
allowed-tools: Bash Skill
---

# /eval-consistency

Evaluate and fix cross-document consistency (PRD, Design, UI, Tasks). Detects inconsistencies via 1000-point scoring, then auto-fixes downstream docs to align with PRD as source of truth.

## Config Resolution

Resolve target and iterations from config, respecting CLI argument priority:

1. Check if `$ARGUMENTS` contains `--target` — if so, use the user-provided value.
2. Otherwise, run `forge config get eval.consistency.target 2>/dev/null` — if exit code 0 and non-empty, use the config value.
3. Same for `--iterations`: check `$ARGUMENTS` first, then `forge config get eval.consistency.iterations 2>/dev/null`.
4. If neither source provides a value, omit the argument (eval skill uses rubric default).

## Execution

```bash
# Resolve target
TARGET_ARG=""
if echo "$ARGUMENTS" | grep -qE '(^|\s)--target(\s|$)'; then
  TARGET_ARG=""
else
  TARGET=$(forge config get eval.consistency.target 2>/dev/null)
  if [ $? -eq 0 ] && [ -n "$TARGET" ]; then
    TARGET_ARG="--target $TARGET"
  fi
fi

# Resolve iterations
ITERATIONS_ARG=""
if echo "$ARGUMENTS" | grep -qE '(^|\s)--iterations(\s|$)'; then
  ITERATIONS_ARG=""
else
  ITERATIONS=$(forge config get eval.consistency.iterations 2>/dev/null)
  if [ $? -eq 0 ] && [ -n "$ITERATIONS" ]; then
    ITERATIONS_ARG="--iterations $ITERATIONS"
  fi
fi
```

Pass resolved config values alongside user arguments (including `--scope` if provided):

```
Skill(skill="forge:eval", args="--type consistency $TARGET_ARG $ITERATIONS_ARG $ARGUMENTS")
```
