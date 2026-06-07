# Contract: justfile-integration / Step 4: Mixed-Scope & CLI

## Outcome "mixed-template-identity"
- Preconditions: "mixed.just template exists"
- Input: "Read mixed template"
- Output: "Template has frontend_dir and backend_dir variables"
- State: "mixed template identity verified"
- Side-effect: none

## Outcome "scoped-recipes-bash-case"
- Preconditions: "mixed.just template exists"
- Input: "Read mixed template, check each scoped recipe"
- Output: "11 scoped recipes (compile, build, run, dev, test, lint, fmt, check, clean, install) have scope parameter with bash case dispatch"
- State: "scoped recipe structure verified"
- Side-effect: none

## Outcome "star-branch-error"
- Preconditions: "mixed.just template exists"
- Input: "Count occurrences of invalid scope error message"
- Output: "At least 10 occurrences of *) error branch with stderr output"
- State: "error handling for invalid scope verified"
- Side-effect: none

## Outcome "empty-branch-both"
- Preconditions: "mixed.just template exists"
- Input: "Read mixed template"
- Output: "Empty scope branch chains both frontend npm and backend placeholder commands"
- State: "dual-scope dispatch verified"
- Side-effect: none

## Outcome "bash-pipefail"
- Preconditions: "mixed.just template exists"
- Input: "Count bash recipes and pipefail occurrences"
- Output: "pipefail count >= bash recipe count"
- State: "error propagation in all bash recipes verified"
- Side-effect: none

## Outcome "npm-backend-placeholders"
- Preconditions: "mixed.just template exists"
- Input: "Read mixed template"
- Output: "Frontend uses npm commands, backend uses BACKEND_* placeholders"
- State: "command structure verified"
- Side-effect: none

## Outcome "unscoped-recipes"
- Preconditions: "mixed.just template exists"
- Input: "Check probe, e2e-test, ci, e2e-setup, e2e-verify recipes"
- Output: "None of these recipes have scope parameter"
- State: "unscoped recipe structure verified"
- Side-effect: none

## Outcome "boundary-markers"
- Preconditions: "mixed.just template exists"
- Input: "Read mixed template"
- Output: "Forge boundary markers present in template"
- State: "boundary structure verified"
- Side-effect: none

## Outcome "all-15-recipes"
- Preconditions: "mixed.just template exists"
- Input: "Read mixed template"
- Output: "All 15 recipes present (compile, build, run, dev, test, e2e-test, lint, fmt, check, clean, install, ci, e2e-setup, probe, e2e-verify)"
- State: "recipe completeness verified"
- Side-effect: none

## Outcome "ci-chains-standard"
- Preconditions: "mixed.just template exists"
- Input: "Read ci recipe line"
- Output: "ci chains install, compile, build, test, lint"
- State: "ci recipe structure verified"
- Side-effect: none

## Outcome "cli-e2e-setup-usage"
- Preconditions: "forge skills and commands exist"
- Input: "Read run-e2e-tests/SKILL.md, execute-task.md, run-tasks.md, fix-bug.md, submit-task/SKILL.md"
- Output: "All use 'just test' or 'just compile', none use hardcoded language-specific commands"
- State: "CLI integration with justfile verified"
- Side-effect: none

## Outcome "e2e-setup-verify-recipes"
- Preconditions: "mixed.just template exists, init-justfile SKILL.md exists"
- Input: "Read mixed template and SKILL.md"
- Output: "e2e-setup has package.json check, e2e-verify has VERIFY marker scanning and --feature flag"
- State: "e2e recipe structure verified"
- Side-effect: none

## Outcome "e2e-setup-live-execution"
- Preconditions: "just command available, tests/e2e/package.json and node_modules exist"
- Input: "runJust e2e-setup"
- Output: "Exit code 0, output contains 'OK: e2e dependencies ready'"
- State: "e2e-setup live execution verified"
- Side-effect: none

## Outcome "e2e-setup-idempotent"
- Preconditions: "just command available, tests/e2e/package.json and node_modules exist"
- Input: "runJust e2e-setup twice"
- Output: "Both runs exit 0 with 'OK: e2e dependencies ready'"
- State: "e2e-setup idempotency verified"
- Side-effect: none

## Outcome "init-justfile-e2e-targets"
- Preconditions: "init-justfile SKILL.md and templates exist"
- Input: "Read SKILL.md and generic.just template"
- Output: "SKILL.md references e2e-setup and e2e-verify with --feature and VERIFY markers; generic.just has playwright install"
- State: "e2e target generation verified"
- Side-effect: none

## Journey Invariants
- mixed template consistent across Step 4 tests
- CLI files reference just commands, never hardcoded language commands
- runJust helper used for live execution tests
