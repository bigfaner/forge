---
created: 2026-05-19
prd: prd/prd-spec.md
status: Draft
---

# Technical Design: test-knowledge-convention-driven

## Overview

Three-layer Convention-driven architecture replacing the embedded Profile system. The change spans two codebases: **Go CLI** (forge-cli) and **Plugin Skills** (skills/). No new infrastructure — Convention files use the existing `docs/conventions/` directory and `domains` frontmatter loading mechanism.

**Approach**: Delete `pkg/profile/` entirely. Create a minimal `pkg/forgeconfig/` package for retained config types (`AutoConfig`, `ModeToggle`, `WorktreeConfig`). Rewrite all Profile consumers to either use `forgeconfig` (for config access) or remove language/framework dependencies entirely (for test generation, which moves to LLM + Convention).

**Key Decision** [Decision]: Create `pkg/forgeconfig/` as a minimal config package rather than keeping a hollowed-out `pkg/profile/`. This avoids confusion — any reference to "profile" in the codebase after this feature should be a bug, not a legacy path.

## Architecture

### Layer Placement

```
┌──────────────────────────────────────────────────────┐
│ Plugin Skills Layer (distributed to user environment)│
│                                                      │
│ gen-test-scripts  ←── Convention + Reconnaissance   │
│ run-e2e-tests     ←── Convention Result Format       │
│ init-justfile     ←── Convention + Reconnaissance   │
│ test-guide (NEW)  →── Convention file creation       │
│                                                      │
│ Skills read docs/conventions/ via LLM (no Go code)   │
├──────────────────────────────────────────────────────┤
│ Go CLI Layer (forge-cli binary)                      │
│                                                      │
│ pkg/forgeconfig/  — retained config types only       │
│ pkg/journey/      — simplified: remove Generate*     │
│ pkg/task/         — simplified: remove Language type │
│ pkg/e2e/          — simplified: test exec via just    │
│ pkg/just/         — simplified: remove ProjectType   │
│ internal/cmd/     — simplified init/config/index/add │
│                                                      │
│ CLI no longer knows about frameworks or languages    │
├──────────────────────────────────────────────────────┤
│ Convention Files (in user project docs/conventions/) │
│                                                      │
│ testing-go.md        domains: [testing, go]          │
│ testing-javascript   domains: [testing, javascript]  │
│                                                      │
│ User-editable markdown, fixed section structure       │
└──────────────────────────────────────────────────────┘
```

### Component Diagram

```
                    ┌─────────────────┐
                    │ docs/conventions │
                    │ testing-*.md    │
                    └────────┬────────┘
                             │ LLM reads
                    ┌────────▼────────┐
                    │  Plugin Skills   │
                    │ gen-test-scripts │
                    │ run-e2e-tests    │
                    │ init-justfile    │
                    │ test-guide (NEW) │
                    └────────┬────────┘
                             │ forge CLI + just
                    ┌────────▼────────┐
                    │   Go CLI Layer   │
                    │ forgeconfig      │── config.yaml (auto.*, worktree)
                    │ journey (slim)   │── output dir, smoke test name
                    │ task (slim)      │── task generation without language axis
                    │ e2e (slim)       │── test exec via just
                    │ just (slim)      │── recipe execution, gate sequences
                    └─────────────────┘
```

### Dependencies

**Internal**:
- `pkg/forgeconfig/` — new package, replaces `pkg/profile/config.go`
- `pkg/contract/` — unchanged
- `pkg/journey/` — simplified, drops profile dependency
- `pkg/task/` — simplified, drops profile types
- `pkg/e2e/` — simplified, drops profile dependency
- `pkg/just/` — simplified, drops project-type dependency

**External**: No new external dependencies.

## Interfaces

### I-1: pkg/forgeconfig — Config Access

```go
type Config struct {
    Auto     *AutoConfig     `yaml:"auto"`
    Worktree *WorktreeConfig `yaml:"worktree"`
}

type AutoConfig struct {
    E2eTest          ModeToggle `yaml:"e2eTest"`
    ConsolidateSpecs ModeToggle `yaml:"consolidateSpecs"`
    CleanCode        ModeToggle `yaml:"cleanCode"`
    GitPush          bool       `yaml:"gitPush"`
}

type ModeToggle struct {
    Quick bool `yaml:"quick"`
    Full  bool `yaml:"full"`
}

type WorktreeConfig struct {
    SourceBranch string   `yaml:"source-branch"`
    CopyFiles    []string `yaml:"copy-files"`
}

func ReadConfig(projectRoot string) (*Config, error)
func ReadAutoConfig(projectRoot string) (*AutoConfig, error)
func GetConfigValue(projectRoot, key string) (string, error)
func WriteConfig(projectRoot string, cfg *Config) error
```

Removed fields (vs current `profile.ForgeConfig`): `ProjectType`, `Languages`, `Interfaces`, `TestFramework`, `TestCommand`.

### I-2: pkg/journey — Simplified Journey Model

```go
type TestGenerationOpts struct {
    Journey   string
    Contracts []contract.Contract
    Facts     string // Fact Table content from Code Reconnaissance
}

type GeneratedTest struct {
    Filename    string
    Content     string
    IsSmokeTest bool
}

func TestOutputDir(featureSlug string) string
func SmokeTestName(journey string) string
```

**Removed**: `FrameworkInfo` field from `TestGenerationOpts`, `GenerateGoTest()`, `GeneratePythonTest()`, `GenerateJSTest()`, `GenerateDispatched()`, `FeatureTag()`, `RegressionTag()`.

**Tag handling** [Decision]: `FeatureTag`/`RegressionTag` removed from Go code. `forge test promote` discovers tag syntax by regex-scanning the existing test file for known tag patterns (`//go:build \w+`, `@pytest.mark.\w+`, `describe\w+`, etc.) rather than from a framework database. This makes promote framework-agnostic without Convention access.

### I-3: pkg/task — Simplified Task Generation

```go
type TestTaskDef struct {
    ID       string
    Key      string
    Title    string
    TestType string // "e2e" only
}

type BuildIndexOpts struct {
    AutoConfig *forgeconfig.AutoConfig
    // No Languages, TestInterfaces, or StrategyResolver
}

func GetBreakdownTestTasks(auto forgeconfig.AutoConfig) []TestTaskDef
func GetQuickTestTasks(auto forgeconfig.AutoConfig) []TestTaskDef
```

**Simplification**: `auto.E2eTest.Full/Quick` decides whether to generate e2e tasks. No per-language task generation — single e2e task per Journey, no embedded strategy content.

**Task content change**: Generated e2e task files no longer contain `forge test get generate` output (strategy content). Instead, they reference Convention path: "Read `docs/conventions/testing-*.md` for framework knowledge."

### I-4: pkg/e2e — Simplified E2E Execution

```go
type RunOpts struct {
    ProjectRoot string
    Feature     string
    Force       bool
}
```

**Removed**: `ResolveProfile()`, all `profile.*` imports. E2E execution goes through `just e2e-test` directly.

### I-5: pkg/just — Simplified Scope Resolution

```go
func ResolveScope(projectRoot string, verb string, scope string) []string
```

**Change**: Remove `profile.ReadConfig` dependency for `ProjectType`. Instead, try-based resolution: attempt `just <verb> <scope>` — if recipe not found, fall back to `just <verb>` without scope argument.

### I-6: Plugin Skills — Convention-Driven Interfaces

**gen-test-scripts**:
- Step 0: List `docs/conventions/` → load files matching `domains: [testing, ...]`
- Step 1: Code Reconnaissance → scan existing test files for framework patterns
- Step 2: Generate code using Convention + Reconnaissance + Skill methodology
- Step 3: Compile gate `just e2e-compile` → retry on failure (max 2)
- Output: `tests/<journey>/` with test files
- **Removed**: All `forge test detect`, `forge test get generate`, `forge test framework`, `forge test get template` calls

**run-e2e-tests**:
- Step 0: Load Convention Result Format section
- Step 1: Execute `just e2e-test`
- Step 2: Parse results using Convention-declared format type (json-stream, json-report, text-verbose)
- Output: Structured test results report
- **Removed**: `forge test detect`, `forge test get run` calls; no hardcoded framework-specific parsing logic

**init-justfile**:
- Step 0: List `docs/conventions/` → load Framework section
- Step 1: Generate e2e-compile/e2e-test/e2e-setup recipes from Convention + LLM knowledge
- Output: justfile with e2e recipes
- **Removed**: `forge test detect`, `forge test get justfile` calls

**test-guide (NEW)**:
- Step 1: Scan file signals (go.mod, package.json, Cargo.toml, etc.)
- Step 2: Scan existing test files → extract patterns (imports, tags, naming)
- Step 3: Present findings / ask user to select framework
- Step 4: Write `docs/conventions/testing-<scope>.md` with minimal structure
- Output: Convention file with Framework + Assertion + Tags + Result Format sections

## Data Models

### Convention File (user project level)

```yaml
# docs/conventions/testing-<scope>.md
frontmatter:
  title: string       # "Test Conventions"
  domains: []string   # [testing, go] or [testing, javascript, web-ui]

sections:
  Framework:
    name: string           # "Go testing package + testify/assert"
    file-pattern: string   # "*_test.go"
    package: string        # "e2e"
    required: true

  Assertion:
    library: string        # "assert (not require)"
    key-functions: []string  # ["assert.NoError", "assert.Contains"]
    required: true

  Tags:
    build-tag: string      # "//go:build e2e"
    required: true

  Result Format:
    output-flags: string    # "-json"
    format-type: string     # "json-stream" | "json-report" | "text-verbose"
    required: true

  Helpers:         # optional
  Import Patterns: # optional
  Code Style:      # optional
  Anti-patterns:   # optional
```

### Config (simplified)

```yaml
# .forge/config.yaml — after cleanup
auto:
  e2eTest: { quick: false, full: true }
  consolidateSpecs: { quick: true, full: true }
  cleanCode: { quick: false, full: true }
  gitPush: true
worktree:
  source-branch: v3.0.0
```

## Error Handling

### Error Types

| Error Context | Behavior | User-Facing Output |
|---------------|----------|--------------------|
| Convention file missing `domains` frontmatter | Skip file, log warning | "Convention file `<path>` has no domains frontmatter. Skipping." |
| Convention file missing required section | Proceed with LLM defaults for that area | "Convention file `<path>` is missing sections: Assertion, Tags. Using LLM defaults for those sections." |
| Convention file unreadable (permissions, encoding) | Skip file, log warning | "Cannot read Convention file `<path>`: `<error>`. Skipping." |
| No Convention files found | Proceed with LLM defaults + Reconnaissance | "No test Convention files found in `docs/conventions/`. Generation will use LLM defaults. Run `/forge:test-guide` to create one." |
| Convention vs Reconnaissance conflict | Convention wins, log conflict | "Convention declares `<X>` but existing tests use `<Y>`. Using Convention value." |
| `just e2e-compile` recipe missing | Block generation | "Missing justfile `e2e-compile` recipe. Run `/forge:init-justfile` first, or add a recipe manually." |
| Compile gate failed (all retries) | Block task, preserve generated files | Output compile error + generated file path + recovery actions: check Convention, run test-guide, or manually edit |
| `forge test promote` finds no tag pattern | Skip tag promotion | "No test tag pattern found in `<file>`. Skipping promotion." |

### Propagation Strategy

- **Skills (LLM layer)**: Errors are output as text messages. Skills do not throw — they report and suggest recovery. Compile gate failures block the task (agent stops).
- **Go CLI layer**: Errors use existing error patterns (exit code 1, descriptive stderr). No new error types needed — removed commands are simply gone, remaining commands have fewer error paths.

## Cross-Layer Data Map

Single-layer feature (CLI tool + skill prompts). Cross-Layer Data Map not applicable.

## Integration Specs

No existing-page integrations — not applicable.

## Testing Strategy

### Regression-First Approach

The 126+ existing e2e tests in forge-cli are the primary validation mechanism. No new unit tests for the removed code — the testing strategy is "everything that worked before must still work."

| Phase | Test Type | Tool | What to Test | Gate |
|-------|-----------|------|--------------|------|
| Phase 0 (POC) | Compile rate | `just e2e-compile` | gen-test-scripts output compiles on 3 frameworks | >= 70% first-pass rate |
| Phase 1 (Profile removal) | Regression | `go build ./...`, `go test ./...` | Go CLI compiles and all unit tests pass | 100% pass |
| Phase 2 (Skill rewrites) | End-to-end | gen-test-scripts → `just e2e-compile` → `just e2e-test` | Full pipeline on forge-cli's 126+ Journeys | >= 85% first-pass compile rate; 100% test pass |
| Phase 3 (test-guide) | Functional | `/forge:test-guide` → gen-test-scripts → compile | Convention generation → code generation pipeline | Convention file valid; generated code compiles |
| Validation | Diff equivalence | diff Convention-generated vs Profile-generated | Framework core patterns (import, assertion, tag) identical | 0 diff on core patterns |

### Key Test Scenarios

1. **forge-cli project (Go + go-testing + testify)**: Convention-declared `assert` (not `require`), `//go:build e2e` tag, `go test -json` result format. All 126+ tests pass.
2. **Cold start (no Convention)**: LLM generates Go testing code from defaults. `just e2e-compile` passes.
3. **Malformed Convention**: Missing Assertion section. Skill proceeds with LLM defaults, outputs warning.
4. **Compile gate failure**: Generated code fails `just e2e-compile`. Skill retries with error feedback, then blocks with recovery guidance on exhaustion.
5. **Multi-framework**: `testing-go.md` and `testing-javascript.md` coexist. gen-test-scripts loads only the relevant one per Journey.
6. **forge test promote**: Discovers `//go:build e2e` tag in existing file, promotes to `//go:build e2e && regression`. No Profile needed.

### Import Audit Gate

Pre-deletion verification: `grep -r "pkg/profile" forge-cli/` must return zero results before `pkg/profile/` is deleted. This is an automated check in the build/test pipeline.

## Security Considerations

### Threat Model

- **Convention file injection**: Malicious Convention content could instruct LLM to generate harmful code. Mitigated by: (1) Convention files are user-editable markdown in the user's own project, (2) generated code is constrained by compile gate, (3) no runtime execution of Convention content — it's LLM context.

### Mitigations

- Convention files live in `docs/conventions/` (user project, version controlled) — not downloaded from external sources.
- Compile gate acts as a sandbox — generated code must compile as valid test code.
- No new attack surface vs current Profile system (which also provides LLM instructions via generate.md).

## PRD Coverage Map

| PRD Requirement / AC | Design Component | Interface / Model |
|----------------------|------------------|-------------------|
| FS-1: Convention file structure | Convention File data model | Data Models section |
| FS-2: Convention loading mechanism | gen-test-scripts Skill I-6, error handling table | Skill prompt (Step 0) |
| FS-3: Code Reconnaissance extension | gen-test-scripts Skill I-6, TestGenerationOpts.Facts | I-2: TestGenerationOpts |
| FS-4: Compile gate | gen-test-scripts Skill I-6, error handling table | Skill prompt (Step 3) |
| FS-5: /forge:test-guide command | test-guide Skill I-6 | New skill: plugins/forge/skills/test-guide/ |
| FS-6: config.yaml cleanup | I-1: forgeconfig.Config | pkg/forgeconfig/ |
| FS-7: Profile removal | Architecture section, import audit gate | pkg/profile/ deleted |
| FS-8: Silent migration | I-1: forgeconfig.ReadConfig ignores unknown fields | yaml.Unmarshal behavior |
| FS-9: consolidate-specs integration | Convention files use standard domains frontmatter | docs/conventions/ (existing) |
| Story 1: Non-default framework | Convention File + gen-test-scripts I-6 | docs/conventions/testing-<scope>.md |
| Story 2: Convention bootstrap | test-guide Skill I-6 | New skill |
| Story 3: Cold start | gen-test-scripts I-6 (no Convention → LLM defaults) | Skill prompt |
| Story 4: Backward compatibility | forgeconfig silent ignore + e2e test regression | I-1 + Testing Strategy |
| Story 5: Multi-framework | Convention domains filtering | Skill prompt (Step 0) |
| Story 6: Compile gate recovery | Error handling: compile gate failed | Error Handling table |

## Open Questions

- [ ] None — all design decisions resolved during PRD and tech design review.

## Appendix

### Alternatives Considered

| Approach | Pros | Cons | Why Not Chosen |
|----------|------|------|----------------|
| Keep hollowed-out `pkg/profile/` with config types only | Fewer import changes | Misleading name — any "profile" reference after this feature should be a bug | New `pkg/forgeconfig/` package makes the clean break explicit |
| Hardcode tag patterns in `forge test promote` | Simpler than regex scanning | Couples promote to known frameworks, same problem as Profile | Regex scan is more flexible and requires no framework knowledge |
| Store test-command in config.yaml | Allows override for non-just users | Adds config surface for an already-standardized path (justfile) | Users modify justfile directly; no config needed |
