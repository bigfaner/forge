# Contract: justfile-integration / Step 1: Init-Justfile

## Outcome "frontend-scope-free"
- Preconditions: "init-justfile skill and templates exist"
- Input: "Read frontend template (node.just)"
- Output: "Template has no scope parameter, contains npx tsc --noEmit"
- State: "template content verified as scope-free"
- Side-effect: none

## Outcome "backend-scope-free"
- Preconditions: "init-justfile skill and templates exist"
- Input: "Read backend template (go.just)"
- Output: "Template has no scope parameter, contains go vet"
- State: "template content verified as scope-free"
- Side-effect: none

## Outcome "mixed-scope-aware"
- Preconditions: "init-justfile skill and templates exist"
- Input: "Read mixed template (mixed.just)"
- Output: "Template has scope parameter with case/esac dispatch"
- State: "template content verified as scope-aware"
- Side-effect: none

## Outcome "standard-commands-present"
- Preconditions: "forge project justfile exists"
- Input: "Read project justfile"
- Output: "All 11 standard commands present (compile, build, run, dev, test, lint, fmt, check, clean, install, ci)"
- State: "justfile vocabulary complete"
- Side-effect: none

## Outcome "no-markers-error"
- Preconditions: "init-justfile SKILL.md exists"
- Input: "Read SKILL.md content"
- Output: "Error handling description for no project markers detected"
- State: "error path documented"
- Side-effect: none

## Outcome "existing-justfile-confirm"
- Preconditions: "init-justfile SKILL.md exists"
- Input: "Read SKILL.md content"
- Output: "User confirmation mechanism described (confirm/prompt/overwrite/--force)"
- State: "interactive path documented"
- Side-effect: none

## Outcome "boundary-marker-idempotent"
- Preconditions: "forge project justfile exists with boundary markers"
- Input: "Read project justfile"
- Output: "Boundary markers present, custom recipes outside markers"
- State: "idempotent merge structure verified"
- Side-effect: none

## Journey Invariants
- project root consistent across all steps (testkit.ProjectRoot)
- boundary markers format: "# --- forge standard recipes ---" / "# --- end forge standard recipes ---"
