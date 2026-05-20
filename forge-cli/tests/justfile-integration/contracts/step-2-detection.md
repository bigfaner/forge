# Contract: justfile-integration / Step 2: Forge Detection

## Outcome "project-type-via-probe"
- Preconditions: "forge CLI binary built and available"
- Input: "forge probe command"
- Output: "exit code 0, non-empty project type string"
- State: "project type detected"
- Side-effect: none

## Outcome "ten-scoped-recipes"
- Preconditions: "forge project justfile with standard recipes section"
- Input: "Extract standard section between boundary markers"
- Output: "10 scoped recipes with scope parameter (compile, build, run, dev, test, lint, fmt, check, clean, install)"
- State: "scoped recipe vocabulary verified"
- Side-effect: none

## Outcome "unscoped-recipes"
- Preconditions: "forge project justfile with standard recipes section"
- Input: "Extract standard section between boundary markers"
- Output: "ci recipe present without scope parameter"
- State: "unscoped recipe verified"
- Side-effect: none

## Outcome "boundary-markers"
- Preconditions: "forge project justfile exists"
- Input: "Read justfile content"
- Output: "Start and end boundary markers present"
- State: "boundary structure verified"
- Side-effect: none

## Outcome "compile-backend-command"
- Preconditions: "standard section extracted from justfile"
- Input: "Inspect compile recipe"
- Output: "Compile has scope parameter, contains go vet"
- State: "compile recipe verified for Go project"
- Side-effect: none

## Outcome "build-backend-command"
- Preconditions: "standard section extracted from justfile"
- Input: "Inspect build recipe"
- Output: "Build has scope parameter, contains go build"
- State: "build recipe verified for Go project"
- Side-effect: none

## Outcome "test-backend-command"
- Preconditions: "standard section extracted from justfile"
- Input: "Inspect test recipe"
- Output: "Test has scope parameter, contains go test"
- State: "test recipe verified for Go project"
- Side-effect: none

## Outcome "shebang-pipefail"
- Preconditions: "standard section extracted from justfile"
- Input: "Check each scoped recipe body"
- Output: "All scoped recipes have #!/usr/bin/env bash and set -euo pipefail"
- State: "error propagation structure verified"
- Side-effect: none

## Outcome "ci-chain"
- Preconditions: "standard section extracted from justfile"
- Input: "Inspect ci recipe"
- Output: "ci chains: just install, just compile, just build, just test, just lint"
- State: "ci recipe structure verified"
- Side-effect: none

## Outcome "custom-recipes-preserved"
- Preconditions: "forge project justfile exists with custom recipes"
- Input: "Read justfile content"
- Output: "claude: and claude-c: recipes present (outside boundary markers)"
- State: "custom recipe preservation verified"
- Side-effect: none

## Outcome "compile-scope-dispatch"
- Preconditions: "just command available in PATH"
- Input: "runJust with compile backend, compile frontend, compile (empty), compile invalidscope"
- Output: "No scope errors for any scope value"
- State: "scope dispatch behavior verified"
- Side-effect: none

## Outcome "detection-signals"
- Preconditions: "init-justfile SKILL.md exists"
- Input: "Read SKILL.md content"
- Output: "Detection checks package.json, go.mod, Cargo.toml, pyproject.toml"
- State: "detection signal mapping verified"
- Side-effect: none

## Outcome "classification"
- Preconditions: "init-justfile SKILL.md exists"
- Input: "Read SKILL.md content"
- Output: "Classification produces mixed, frontend, backend, error cases"
- State: "classification logic verified"
- Side-effect: none

## Outcome "template-selection"
- Preconditions: "init-justfile templates exist (go.just, node.just, mixed.just)"
- Input: "Read each template"
- Output: "Backend has go vet, frontend has npx tsc, mixed has frontend_dir/backend_dir"
- State: "template selection verified"
- Side-effect: none

## Outcome "force-flag"
- Preconditions: "init-justfile SKILL.md exists"
- Input: "Read SKILL.md content"
- Output: "--force flag documented, agent/non-interactive skip described"
- State: "force flag behavior verified"
- Side-effect: none

## Journey Invariants
- boundary markers format consistent with Step 1
- testkit.ProjectRoot resolves to same root across all steps
