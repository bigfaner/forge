---
created: 2026-05-14
author: "faner"
status: Draft
---

# Proposal: Establish justfile-recipes as Single Source of Truth for E2E Commands

## Problem

Test commands are duplicated across three layers: profile `justfile-recipes`, Go `e2e/actions.go` switch/case, and `manifest.yaml` `run.*`/`graduate.*` fields. They drift apart silently — actions.go only supports 2/6 profiles, and manifest.yaml test fields are never parsed by any code.

### Evidence

- `e2e/actions.go` hardcodes `go test` for go-test and `npx playwright test` for web-playwright. The other 4 profiles (pytest, java-junit, rust-test, maestro) return "unsupported profile" errors.
- `manifest.yaml` defines 8 command fields (`run.command`, `run.compile`, `run.result-format`, `graduate.target-directory`, `graduate.merge-strategy`, `graduate.import-rewrite`, `graduate.compile-check`, `graduate.list-tests`) across 6 profiles = 48 field values — zero of which are parsed by Go code or referenced by name in any skill.
- The primary consumer (quality-gate) already bypasses `pkg/e2e` entirely, calling `just test-e2e` directly. The Go switch/case only serves 5 secondary CLI commands.
- Adding a new profile requires changes in 3 places (manifest.yaml, actions.go, justfile-recipes) with no enforcement that they stay in sync.

### Urgency

The migrate-e2e-to-go feature is converting test infrastructure. This is the right time to clean up the command duplication before it entrenches further. Delay means the next profile addition (e.g., rust-test becoming supported) requires coordinated changes across 3 layers.

## Proposed Solution

Make `justfile-recipes` the single source of truth for e2e test execution. Remove dead command fields from manifest.yaml. Delegate `e2e/actions.go` functions to `just` recipes.

**End state:**

| Layer | Role |
|-------|------|
| `justfile-recipes` | Canonical execution commands (run, setup, compile, discover) |
| `manifest.yaml` | Profile metadata only (name, display, language, file-extension, test-directory, capabilities, templates) |
| Strategy `.md` files | Skill guidance (run.md, generate.md, graduate.md) |
| `e2e/actions.go` | Thin CLI wrappers: `exec.Command("just", "test-e2e")` + generic Verify() |

### Innovation Highlights

**Delegation abstraction: "task runner as single source of truth."** The pattern — a thin CLI wrapper delegating all execution to a canonical task runner — is used by Cargo (`cargo test` delegates to the Rust test harness), Gradle (`gradle test` delegates to JVM test frameworks), and Bazel (`bazel test //...` delegates to language-specific test rules). In each case, the build tool does not embed language-specific test logic; it defines a task graph and delegates execution. Forge can apply the same pattern: `forge e2e run` becomes a task-graph entry that delegates to `just test-e2e`, just as `cargo test` delegates to `libtest`.

**Cross-domain inspiration: the CRI analogy.** The novel twist in forge's context is that `forge e2e run` dispatches to fundamentally different language ecosystems (Go, TypeScript, Python, Java, Rust) through a single task runner interface. This is structurally analogous to Kubernetes' Container Runtime Interface (CRI): kubelet exposes a single `RuntimeService` API, and delegates to containerd, CRI-O, or gVisor depending on the node's runtime configuration. Forge's justfile-recipes play the same role as the CRI: a stable, versioned interface that decouples the orchestrator (`forge e2e`) from the runtime implementation (Go test harness, Playwright, pytest, JUnit). Each profile "implements" the recipe interface the way each CRI runtime implements the CRI API — the orchestrator neither knows nor cares which runtime executes. Other cross-domain analogues: Git's `difftool` / `mergetool` mechanism delegates diff/merge display to any configured tool via a stable interface; Unix shell PATH resolution routes a single command name to whichever executable appears first. The common principle: a thin routing layer with a stable contract, backed by pluggable implementations. Forge's contribution is applying this principle to a CLI that manages multi-language test profiles — a context where the build-tool analogues (Cargo, Gradle) are mono-runtime by design.

**Reusable insight for forge subsystems.** This delegation pattern is not specific to e2e. Any forge subsystem that currently embeds command logic in Go switch/case blocks (e.g., `forge build`, `forge lint`) could adopt the same pattern: define the canonical command in justfile-recipes, replace Go logic with `exec.Command("just", <recipe>)`, and remove profile-specific Go code. The metadata-code gap (manifest.yaml fields that no code reads) is a symptom of violating this pattern. By establishing it for e2e first, forge gains a repeatable template: *if a task runner recipe exists, delegate to it; do not duplicate the command in Go or YAML.*

## Requirements Analysis

### Key Scenarios

- `forge e2e run` — delegates to `just test-e2e`; passes `--feature` flag through as justfile argument
- `forge e2e setup` — delegates to `just e2e-setup`
- `forge e2e compile` — delegates to `just e2e-compile`
- `forge e2e discover` — delegates to `just e2e-discover`
- `forge e2e verify` — stays as-is (generic file scanner, no justfile equivalent)
- `forge quality-gate` — unchanged (already calls just directly)
- Skills (gen-test-scripts, run-e2e-tests, graduate-tests) — unchanged (they load strategy .md files, not manifest.yaml command fields)

### Constraints & Dependencies

- Requires `just` installed on the system. Already a requirement for quality-gate, so no new dependency.
- `justfile-recipes` must be injected into project justfile via `init-justfile` skill before `forge e2e *` commands work. Already the case for any forge-managed project.

### Error Scenarios

| Scenario | Trigger | Expected behavior |
|----------|---------|-------------------|
| `just` not installed | User runs `forge e2e run` without `just` on PATH | Exit code 1; stderr prints: `error: 'just' is required but not found on PATH. Install: https://github.com/casey/just` |
| Recipe missing from justfile | User runs `forge e2e compile` but justfile lacks `e2e-compile` recipe | `just` exits non-zero with "just: no recipe found for 'e2e-compile'"; forge propagates the exit code and stderr |
| Recipe exits non-zero | Test command in recipe fails (compilation error, test failure) | forge propagates the recipe's exit code; stderr/stdout from `just` passes through unchanged |
| No project justfile | User runs `forge e2e run` in a directory without a justfile | `just` exits non-zero with "just: no justfile found"; forge propagates the exit code and stderr |
| No profile configured | User runs `forge e2e run` without a test profile in `.forge/config.yaml` | Unchanged from current behavior — the justfile recipe reads profile at runtime; if no profile is set, the recipe prints a profile-selection prompt or exits with a message |

### Non-Functional Requirements

**Performance.** Delegating via `exec.Command("just", ...)` adds one process spawn per invocation. On modern systems, `just` startup is ~10ms. This overhead is negligible compared to test execution time (seconds to minutes). No performance regression expected; the current Go switch/case has the same subprocess overhead for profiles that call `go test` or `npx`.

**Compatibility.** All 6 existing profiles (go-test, web-playwright, pytest, java-junit, rust-test, maestro) must continue to work after the change. The `just` minimum version requirement is `just >=1.5.0` (for recipe parameters and `--feature` flag support), which is already satisfied by the version specified in forge's installation docs. No OS-specific behavior changes — `just` and `exec.Command` behave consistently across macOS, Linux, and Windows.

**Security.** No new attack surface. The `just` recipes are injected by forge's `init-justfile` skill and are part of the project's source tree (version-controlled). Delegating to `just` does not introduce arbitrary command execution beyond what the justfile already contains — the recipes are the same commands that were previously hardcoded in Go.

## Alternatives & Industry Benchmarking

### Industry Patterns for Command-Source Unification

The core problem — keeping a single canonical source for build/test commands across a multi-language project — is solved across the industry with varying trade-offs:

**Bazel (BUILD rules as single source of truth).** Bazel centralizes all build/test actions in `BUILD` files. Rules are language-agnostic declarations; Bazel's execution engine dispatches to language-specific toolchains. The key insight: the BUILD file is the only place that defines *what* runs; the toolchains define *how*. Applied to forge: the justfile is the BUILD file, and `just` recipes are the rules. Unlike Bazel, forge does not need hermetic builds or dependency graphs — only command dispatch — so the full Bazel machinery would be vast overkill.

**Taskfile / GitHub Actions (YAML task runner as command source).** Both systems define all tasks in a single YAML file (`Taskfile.yml` or `.github/workflows/*.yaml`). Tasks are shell commands with structured metadata (dependencies, conditions). The YAML is both human-readable and machine-parseable. Applied to forge: manifest.yaml *could* serve this role, but it already fails at it — 48 command fields are defined and zero are parsed by code. Reanimating manifest.yaml as a task runner requires building a YAML executor in Go, reproducing what `just` already provides.

**Cargo / Gradle (language-native task graph).** Cargo defines `cargo test` as the canonical test entry point for Rust projects; Gradle's `build.gradle` defines test tasks for JVM projects. Both eliminate the need for a separate command registry — the build tool *is* the command source. The pattern: if your ecosystem has a dominant task runner, delegate to it rather than maintaining a parallel registry.

**Cross-cutting insight:** Every successful system above shares one principle: a single artifact defines the commands, and all other layers delegate to it. No system maintains a parallel command registry alongside the canonical one. Forge currently violates this principle — the justfile works, but actions.go and manifest.yaml maintain ghost copies.

### Comparison Table

| Approach | Industry Analog | How it solves command unification | Pros | Cons | Verdict |
|----------|----------------|-----------------------------------|------|------|---------|
| Do nothing | — | Does not | Zero cost | 3-layer drift continues; new profiles need 3-file coordination; 48 unused YAML fields accumulate | Rejected: cost of inaction grows with each profile |
| Declarative task registry (manifest.yaml parsed at runtime) | Bazel BUILD rules, GitHub Actions YAML | manifest.yaml becomes machine-readable config; Go code parses and executes commands from it | Flexible; plugins could define custom commands without touching Go code | Requires building a YAML command executor in Go; over-engineering for 6 profiles; manifest.yaml fields already unused — reanimating dead code | Rejected: would need to build and maintain a mini task-runner inside forge |
| Go-native task interface (strategy pattern in actions.go) | Gradle task graph, cargo test dispatch | Define a `TaskRunner` interface with profile-specific implementations; each profile registers its own commands | Type-safe; no subprocess overhead; IDE-completable | Duplicates the profile logic that justfile-recipes already encode; every new profile requires Go code changes and a forge release | Rejected: justfile-recipes already work for all 6 profiles without Go changes |
| **justfile-canonical** | Cargo/Makefile convention, Bazel BUILD-file delegation | justfile-recipes is the single artifact; Go code delegates via `exec.Command("just", ...)` | Already working in production (quality-gate); all 6 profiles work immediately; adding a profile requires only YAML + justfile changes, no Go release | Requires `just` installed; subprocess spawn overhead (~10ms per invocation) | **Selected: follows industry principle of single-source delegation; minimal change to achieve it** |

## Feasibility Assessment

### Technical Feasibility

All changes are in forge-cli Go code and profile YAML files. No external dependencies. The delegation pattern (`exec.Command("just", ...)`) is already used throughout the codebase.

### Resource & Timeline

~3-5 tasks: manifest cleanup, actions.go delegation, test updates, version bump.

## Scope

### In Scope

- Remove `run.*` fields (3) and `graduate.*` fields (5) from all 6 manifest.yaml files
- Simplify `e2e/actions.go`: Run/Setup/Compile/Discover delegate to `just` via exec.Command
- Keep Verify() as-is (generic file scanner)
- Update `e2e/actions_test.go` for new delegation behavior
- Version bump in scripts/version.txt

### Out of Scope

- Strategy .md files (run.md, generate.md, graduate.md) — no changes
- Skill markdown files — no changes (none reference the removed fields)
- `justfile-recipes` files — no changes (already canonical)
- CLI command structure — no changes (same flags, same interface)
- Quality-gate flow — already uses just directly, no changes needed

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| `just` not installed when user calls `forge e2e run` | L | M | Clear error message: "just is required. Install: https://github.com/casey/just" |
| Feature filtering differs between old Go code and justfile recipe | M | M | Add integration test comparing old directory-based vs. new test-name-based filtering results for each existing profile. Document behavior change in CHANGELOG and release notes with migration note for users relying on directory-scoped test selection. If any profile shows divergent test selection, add a compatibility note to the profile's strategy .md file. |
| actions_test.go tests need significant rework | M | L | Tests already use mockable runner interface. Replace mock expectations with "just" command assertions. Same coverage. |

## Success Criteria

- [ ] For each of `run`, `setup`, `compile`, `discover`: `actions_test.go` asserts `exec.Command` is called with args `["just", "<recipe-name>"]` (e.g. `["just", "test-e2e"]` for run, `["just", "e2e-setup"]` for setup). When `--feature` is passed, the feature value is appended as a justfile argument.
- [ ] Subprocess exit code propagation: `actions_test.go` asserts that a non-zero exit from `just` produces a non-nil error from the action function, and a zero exit produces nil error. The forge CLI exit code matches the `just` exit code in both cases.
- [ ] `forge e2e verify` remains unchanged and passes existing tests
- [ ] All 6 profile manifest.yaml files contain zero `run.*` or `graduate.*` fields
- [ ] `go test ./...` passes with 80%+ coverage on modified files
- [ ] `forge quality-gate` behavior unchanged (already uses just directly)
- [ ] Adding a new profile requires changes in 2 files (justfile-recipes + manifest.yaml metadata), not 3

## Next Steps

- Proceed to `/quick-tasks` to generate tasks directly from this proposal
