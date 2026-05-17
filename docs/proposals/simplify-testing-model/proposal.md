---
created: 2026-05-17
author: faner
status: Draft
---

# Proposal: 简化测试概念模型 — 去除 Profile，统一为 project-type + interfaces

## Problem

Forge v2 的测试配置概念混杂，用户需要理解三个独立概念才能配置测试：

```yaml
# v2 config.yaml — 3 个字段，概念重叠
project-type: backend
test-profiles: [go-test]        # 什么是 profile？
capabilities: [api, cli]        # capability 和 profile 是什么关系？
```

### Evidence

Internal friction points observed during v2 development and onboarding. External validation: during onboarding of a new contributor (session 2026-05-10), the contributor asked "what profile should I use for a Go project?" — the question that `test-profiles` forces every new user to answer, even though `detect.go` already knows the answer.

| Problem | Concrete Evidence |
|---------|-------------------|
| Overlapping config layers | `config.go:ForgeConfig` has `TestProfiles` + `Capabilities` as independent fields; `testgen.go:GetBreakdownTestTasks` takes both and produces a profile x capability matrix, requiring callers to merge them |
| Profile concept is opaque | `profile.go:runProfileResolve` has a 3-step fallback chain (config -> detect -> "none") with a stderr hint telling the AI to ask the user — the system cannot self-serve |
| Capability naming is confusing | `ValidTestTypes` in `embed.go` uses "web-ui", "tui", "mobile-ui", "api", "cli" — these are system interfaces, not testing capabilities, yet the field is named `capabilities` |
| Detection already works | `detect.go:DetectProfiles` already scans for `go.mod`, `package.json`+playwright, `Cargo.toml`, etc. — the auto-detection infrastructure exists but is treated as a fallback rather than the primary path |
| Code complexity cost | `GetBreakdownTestTasks` accepts `profiles []string` + `capabilities []string` and must compute per-profile-per-type task variants, creating a combinatorial expansion that is hard to reason about |

### Urgency

v3.0.0 is the designated breaking-change window. Every new feature added to v2 (e.g., `testgen.go`'s per-profile-per-type task expansion in the `task` package) compounds the profile x capability matrix. Delaying this refactor increases the number of call sites that must migrate later.

## Proposed Solution

**Eliminate the Profile concept; replace with language auto-detection.** User config shrinks from 3 fields to 2:

```yaml
# v3 config.yaml — 2 fields, intuitive
project-type: backend
interfaces: [api, cli]          # system-facing interface types
```

Forge auto-detects language from project files (`go.mod` -> Go, `package.json` + playwright -> JavaScript) and selects the corresponding test strategy package internally. Profiles are removed entirely — not hidden, not renamed.

**Why the new override chain is cleaner than v2**: In v2, `test-profiles` in config.yaml declares which profiles to use, but `capabilities` in config.yaml overrides the capabilities that the profile itself declares — two orthogonal axes that both affect the same output (which tests to generate). In v3, there is a single axis: language determines the strategy, and `interfaces` narrows which interface types to test. The `languages` field is a simple override of the auto-detect result (not a second axis), and `interfaces` is a subset filter (not a cross-product). This is a 1D narrowing model, not a 2D matrix.

### Innovation Highlights

**Language auto-detect replaces explicit selection** — mainstream test frameworks map 1:1 to languages (Go -> go test, Rust -> cargo test, Python -> pytest), making user profile selection redundant. Industry practice (CI systems, IDEs, zero-config frameworks) favors inferring from project structure over requiring explicit declaration.

**Multi-framework language handling**: The 1:1 assumption breaks for JavaScript, which has Jest, Playwright, Cypress, Vitest, and Mocha. For v3.0, Forge handles this by detecting the specific framework from project files (Playwright via `@playwright/test` dependency) rather than detecting "JavaScript" generically. The language key is actually a framework key: `javascript` maps to Playwright strategies because that is the only JS framework currently supported. Adding Jest support in the future would add a `javascript-jest` language key with its own detection signal (`jest` in `package.json` devDependencies).

**Cross-domain inspiration — type inference**: The auto-detect-then-override pattern mirrors how statically typed languages perform local type inference (e.g., Rust's `let x = HashMap::new()` infers the type from usage context, with an explicit `let x: HashMap<String, i32>` override available). Forge infers the "type" (language) from project file signals and provides a narrow escape hatch (`languages` override) when inference is wrong. This is the same design trade-off: inference handles 90% of cases correctly, and the escape hatch exists for the remaining 10% — the user never needs to spell out what the system can infer.

## Requirements Analysis

### Key Scenarios

1. **Single-language backend** (Go CLI) — zero config: `project-type: backend`, interfaces optional, language auto-detected from `go.mod`
2. **Frontend** (React + Playwright) — `project-type: frontend`, `interfaces: [web-ui, api]`
3. **Multi-language** (Go backend + JS frontend) — both languages auto-detected and activated; user can narrow via `languages: [go]`
4. **No language detected** — `forge testing detect` returns empty; `forge testing get generate` prints an error message directing the user to add `languages` to config.yaml. Skills surface this as "no test strategy found, please configure languages"
5. **Conflicting signals** (e.g., `go.mod` + `package.json` for lint tooling) — `detect.go` returns both languages; if the `package.json` has no Playwright dependency and no test scripts, `detect.go` does not return `javascript` as a detected language (Playwright is the only supported JS framework, and it is detected via the `@playwright/test` dependency — absence of that dependency means no JS language is detected). User can override with `languages: [go]` to suppress any remaining false positive
6. **Monorepo with language-specific subdirectories** — detection scans the project root only (matching current `detect.go` behavior). Subdirectory-only languages are not detected automatically; user must specify via `languages` override. This is a documented limitation.

### Non-Functional Requirements

| NFR | Requirement | Verification |
|-----|-------------|-------------|
| Performance | Language detection must complete in < 100ms for projects with < 1000 files (current `detect.go` reads at most 6 files from the project root) | Benchmark `DetectProfiles` against repos of varying sizes |
| Compatibility | `forge testing` CLI output format must be stable so skills can parse it reliably; output is plain text with structured blocks, matching current `forge profile` format | E2E tests for each `forge testing` subcommand |
| Security | Detection reads only well-known filenames at the project root (`go.mod`, `package.json`, `Cargo.toml`, etc.) — no directory traversal, no execution of project code | Code review confirms `detect.go` uses `os.ReadFile` on fixed paths only |
| Extensibility | Adding a new language requires: (1) add detection case in `detect.go`, (2) add language directory under `languages/`, (3) add language key to Go constant — no schema migration, no manifest files | Checklist for adding a hypothetical 7th language |

### Constraints & Dependencies

- v3.0.0 breaking change, no backward compatibility with v2 config format
- Strategy files (generate.md, run.md, graduate.md) remain embedded in Go binary via `embed.FS`, exposed through CLI to skills
- 6 existing strategy sets migrated as-is (file content unchanged, only directory structure and embedding paths change)
- Go 1.22+ required (current baseline; `embed.FS` and `os.ReadFile` already in use)
- `gopkg.in/yaml.v3` for config parsing (already a dependency in `profile/config.go`)
- `github.com/spf13/cobra` CLI framework (already in use; `forge testing` is a new command group, same pattern as existing `forge profile`)

## Alternatives & Industry Benchmarking

### Industry References

| Pattern | Reference | How it works | What Forge borrows |
|---------|-----------|-------------|--------------------|
| CI auto-detect | [GitHub Actions setup-go](https://github.com/actions/setup-go) scans for `go.mod` to determine Go version; [setup-node](https://github.com/actions/setup-node) reads `.nvmrc`/`package.json` | File-signal detection as primary path, user override as fallback | Forge's `detect.go` already uses the same signals (`go.mod`, `package.json`, `Cargo.toml`); the pattern is proven for CI and applies directly to test configuration |
| Framework auto-discovery | [Jest configuration](https://jestjs.io/docs/configuration) auto-discovers test files via glob patterns; [pytest](https://docs.pytest.org/en/stable/explanation/goodpractices.html#conventions) auto-discovers `test_*.py` files by convention | Convention over configuration — the tool infers what to test from project structure | Forge does not need to discover individual test files, but the same principle applies: infer the testing strategy from project signals rather than requiring explicit declaration |
| Explicit adapter selection | [Terraform provider configuration](https://developer.hashicorp.com/terraform/language/providers/configuration) requires explicit provider blocks; [ESLint `overrideConfig`](https://eslint.org/docs/latest/use/configure/configuration-files) requires explicit config path | User declares which adapter to use; no detection, full user control | This is the v2 model (`test-profiles` in config.yaml). It works but forces users to learn profile names — acceptable for infrastructure tools, but unnecessary overhead for a test automation tool where language is easily inferred |
| Zero-config frameworks | [Next.js](https://nextjs.org/docs/getting-started/project-structure) auto-detects TypeScript from `tsconfig.json`, testing from `jest.config.*`; [Vite](https://vitejs.dev/config/) auto-detects framework plugins from dependencies | Zero user configuration for common cases; escape hatch available | Forge's ideal: user runs `forge testing get generate` with no config and gets correct output because the language was detected from project files |

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing (keep v2 model) | Current state | Zero migration cost; all skills already work | Conceptual model stays confusing (profile + capability overlap); `testgen.go` matrix complexity persists; every new skill must learn two config axes | Rejected: v3.0 window makes this the right time; the 2-axis model is the root cause of complexity, not a surface issue |
| Gradual deprecation (support both models) | Python 2→3 migration pattern | Backward compatible; users can migrate at their pace | Two parallel codepaths in `profile.go` and `testgen.go`; `ForgeConfig` struct grows new fields alongside deprecated ones; skills must handle both old and new CLI output formats — migration burden shifts from users to maintainers. Python 2→3 itself is widely regarded as a cautionary tale: the decade-long dual-support period fragmented the ecosystem and ultimately required a hard cutoff anyway (PEP 373 sunsetting, Python 2 EOL 2020). Forge's scope (10 skills, single developer) is far smaller, but the lesson holds: dual paths delay rather than avoid the inevitable break | Rejected: the dual-path maintenance cost exceeds the one-time migration cost; Python 2→3 showed that gradual deprecation often ends in a forced cutoff regardless; v3.0 is a breaking-change release by design |
| Rename only (profile → language, capability → interface) | Surface refactor | Lower risk — no directory restructuring, no embed path changes, no `manifest.yaml` removal; all existing tests pass unchanged; no migration guide needed for internal paths | Does not address the core problem: `DetectProfiles()` still returns profile names (renamed to "language" but still explicit), and the user must still set `languages: [go-test]` in config.yaml because auto-detect remains a fallback behind `ReadTestProfiles()`. `GetBreakdownTestTasks()` still accepts two independent axes (`profiles` + `capabilities`, now `languages` + `interfaces`), preserving the combinatorial expansion in `testgen.go` that this proposal aims to eliminate | Rejected: renaming without restructuring leaves the 2-axis complexity intact; the UX improvement is cosmetic |
| **Remove profile + auto-detect language** | CI auto-detect, zero-config frameworks | Users configure only `project-type` + `interfaces`; language detected from file signals that `detect.go` already reads; single-axis config model eliminates the profile × capability matrix in `testgen.go` | All 10 skills (D6) must migrate CLI calls; `embed.go` directory structure changes from `profiles/<name>/` to `languages/<key>/`; `manifest.yaml` files removed, metadata hardcoded in Go | **Selected: matches the zero-config pattern proven in CI/IDE tools; auto-detect is the primary path, explicit override (`languages` field) is the escape hatch** |

### Why Auto-detect Fits Forge Specifically

Forge operates in a constrained domain: the test profile is determined by the project's language and framework, which are discoverable from a small set of well-known file signals (`go.mod`, `package.json`+playwright, `Cargo.toml`). This is unlike Terraform where the provider is an arbitrary choice, or ESLint where config files can be in non-standard locations. The detection signals in Forge are:

1. **Reliable**: `go.mod` present means the project uses Go; `package.json` with `@playwright/test` in devDependencies means the project uses Playwright. False positives are rare — a `package.json` without Playwright in devDependencies does not trigger JavaScript detection, eliminating the most common false positive (lint-only JS tooling in a non-JS project).
2. **Already implemented**: `detect.go:DetectProfiles()` already performs all 6 detections. The change is making it the primary path instead of a fallback behind `ReadTestProfiles()`.
3. **Low-impact failure mode**: If detection fails or returns the wrong language, the user adds `languages: [go]` to config.yaml — a single-line override, simpler than the current model where the user must know the profile name "go-test".

**Limitation acknowledged**: JavaScript has multiple test frameworks (Jest, Playwright, Cypress, Vitest). Forge currently detects only Playwright (via `@playwright/test` dependency); a `package.json` without that dependency does not produce a JavaScript detection result. For v3.0, this is acceptable because all existing Forge JS projects use Playwright. Multi-framework JavaScript support is a future enhancement that would require detecting the test runner from `package.json` scripts or config files.

## Feasibility Assessment

### Technical Feasibility

Verified by auditing the existing codebase:

**1. Detection logic (`forge-cli/pkg/profile/detect.go`)**
`DetectProfiles()` already performs all 6 signal detections. The change: rename return values from profile names ("go-test") to language keys ("go"). The function signature changes from returning `[]string` of profile names to `[]Language` where `type Language string`. Detection logic (file existence checks, package.json parsing) is unchanged.

**2. Config reading (`forge-cli/pkg/profile/config.go`)**
`ForgeConfig` struct changes:
- Remove fields: `TestProfiles []string`, `Capabilities []string`
- Add fields: `Interfaces []string` (replacing capabilities), `Languages []string` (optional override)
`ReadTestProfiles()` is replaced by `ReadLanguages()` which returns the `languages` field if set, otherwise calls `DetectLanguages()` (renamed from `DetectProfiles()`).

**3. Embedded strategy files (`forge-cli/pkg/profile/embed.go`)**
Current: `//go:embed all:profiles` embeds `profiles/<profile-name>/` directories.
Change: `//go:embed all:languages` embeds `languages/<language-key>/` directories. Functions change from `GetStrategy(profileName, kind)` to `GetStrategy(language, kind)` — same logic, different directory prefix. `manifest.yaml` files are removed; capabilities metadata for each language is hardcoded as a Go map:
```go
var languageCapabilities = map[Language][]string{
    "go":         {"api", "cli"},
    "javascript": {"web-ui", "api"},
    "python":     {"api", "cli"},
    "java":       {"api", "cli"},
    "rust":       {"api", "cli"},
    "mobile":     {"mobile-ui"},
}
// Trade-off acknowledged: this hardcodes capability metadata in Go source,
// meaning adding a new language requires a code change and recompilation.
// This is less flexible than manifest.yaml but eliminates a config parsing
// layer and keeps the total surface area small (6 languages, unlikely to grow
// rapidly). The Extensibility NFR's "3 steps" include editing this map as
// step 3 — it is a code change, not zero-code extensibility.
```

**4. CLI command (`forge-cli/internal/cmd/profile.go`)**
New command group `testing` replaces `profile`. The resolve logic in `runProfileResolve` simplifies from a 3-step chain (config -> detect -> none) to: read `languages` from config, or auto-detect. The `get` subcommand no longer takes a profile name argument — it auto-detects language (or uses `--language` flag for multi-language projects).

**5. Task generation (`forge-cli/pkg/task/testgen.go`)**
`GetBreakdownTestTasks(profiles, capabilities, auto)` changes to `GetBreakdownTestTasks(languages, interfaces, auto)`. The per-profile-per-type expansion simplifies to per-language-per-interface. The `TestTaskDef` struct replaces `ProfileName string` with `Language Language`.

### Resource & Timeline

| Work item | Effort estimate | Dependencies |
|-----------|----------------|--------------|
| Restructure `profiles/` to `languages/` (rename directories, remove manifest.yaml, update `embed.go`) | 2 days | None |
| Modify `config.go` (ForgeConfig struct, new ReadLanguages/ReadInterfaces functions) | 1 day | None |
| Rename `DetectProfiles` to `DetectLanguages`, change return type | 1 day | config.go changes |
| Implement `forge testing` CLI command group in `cmd/` | 2 days | config.go + detect.go changes |
| Update `testgen.go` to use Language instead of ProfileName | 1 day | config.go changes |
| Migrate 10 skill files (D6 table) | 2 days | CLI stable |
| Update config.yaml examples and documentation | 1 day | All above |
| **Total** | **~10 person-days (2 weeks)** | |

Single developer can execute; no cross-team coordination needed. The only external dependency is the v3.0.0 release branch being open for merges.

### Dependency Readiness

| Dependency | Status | Risk |
|-----------|--------|------|
| 6 existing strategy file sets under `profiles/` | Verified: all 6 profiles (go-test, web-playwright, rust-test, pytest, java-junit, maestro) have generate.md, run.md, graduate.md, justfile-recipes, and templates/ — straightforward directory rename | None |
| `embed.FS` mechanism | Already in use at `embed.go:14` (`//go:embed all:profiles`); changing to `all:languages` is a path change only | None |
| `gopkg.in/yaml.v3` config parsing | Already in use for `ForgeConfig`; new fields `interfaces` and `languages` are standard YAML string arrays | None |
| `cobra` CLI framework | Already in use; `forge testing` is a new `cobra.Command` with subcommands, same pattern as `forge profile` | None |
| Strategy file content compatibility | Content is pure markdown/text; no code executes within strategy files. Audited: strategy files reference templates via relative paths like `test_file.go` which are resolved by the CLI relative to the strategy directory, not as absolute paths — so renaming `profiles/go-test/` to `languages/go/` does not break template resolution. However, if any strategy file contains a hardcoded reference to `profiles/` (e.g., in an example command), that text would become stale — grep for `profiles/` across all strategy files before rename to catch this | Low (grep audit required) |

## Scope

**Timeline**: 2 weeks (10 person-days), single developer, targeting v3.0.0 release branch.

### In Scope

- Config schema: remove `test-profiles` and `capabilities`, add `interfaces`; optional `languages` override field
- Global rename: `capabilities` -> `interfaces` (Go code, skills, CLI, documentation)
- Internal refactoring: `profiles/<profile-name>/` -> `languages/<language-key>/`, remove `manifest.yaml`, remove `KnownProfiles`
- CLI: `forge profile` -> `forge testing` (detect, get generate/run/graduate/justfile/template)
- Detection logic: output language keys (go, javascript, python, java, rust, mobile)
- All 10 skills consuming profile/capability updated (see D6 table)
- Config schema and example files updated
- Migration guide: v2-to-v3 config mapping document (delete `test-profiles` -> add `interfaces`; map profile names to auto-detected languages)

### Out of Scope

- New strategy packages (migrate existing 6 only)
- Per-feature language narrowing (deferred to future release)
- Mobile strategy deep design (migrate `mobile/` as-is, no new mobile-specific logic)
- Strategy content changes (only change directory organization, not generate.md etc. file content)
- Multi-framework JavaScript support beyond Playwright (Jest, Cypress, Vitest — future release)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Multi-language false positive (e.g., Go project has `package.json` for lint tooling) | M | M | Detection returns both; user adds `languages: [go]` to config.yaml to suppress. Migration guide documents this as a common override. Detection priority is not configurable — signals are additive, not exclusive. |
| Skill migration gap causes runtime errors | M | H | Rated M (not H) because: (1) the 10 skills are in a single monorepo with no external consumers — every migration point is grep-able; (2) CI smoke tests will run all 10 skills against the v3 config format before merge, catching any remaining `forge profile` calls; (3) `grep -r "forge profile" plugins/` is a definitive completeness check. Task: run this grep before closing v3.0 milestone. Each match must be migrated to `forge testing` or `interfaces`. |
| CI pipeline breakage for existing v2 users | H | M | v3.0.0 is a major version bump with documented breaking changes. Migration guide (see Out of Scope) maps each v2 field to its v3 equivalent. No automated migration tool — manual config edit is a 2-line change (delete `test-profiles`, add `interfaces`). |
| File scanning performance on large monorepos | L | M | `detect.go` reads only fixed filenames at the project root (6 `os.ReadFile` calls maximum). No recursive directory walking. Performance is bounded by filesystem stat calls, not repo size. No mitigation needed beyond the current design. |
| Strategy file format incompatibility after directory rename | L | H | Strategy files are pure markdown/text with no path references to their own location. Verified by auditing all 6 strategy sets — none contain relative paths to sibling files. Templates use language-relative paths (e.g., `test_file.go`) which are resolved by the CLI, not embedded in strategy content. |
| Mobile detection does not fit the "language" model (`android/` + `ios/` are platform dirs, not language indicators) | M | L | `mobile` is treated as a special language key in the `Language` type. Detection logic in `detect.go` checks `dirExists(root, "android") || dirExists(root, "ios")` — same as current implementation. The strategies under `languages/mobile/` contain Maestro-specific content. This is a semantic inconsistency but has no user-facing impact since mobile projects are detected correctly. |
| Hardcoded `languageCapabilities` map drifts from actual strategy capabilities | M | M | The Go map must be manually updated whenever a strategy directory adds or changes its supported interfaces. Mitigation: CI test validates that `languageCapabilities[lang]` matches the actual interface types present in `languages/<lang>/` strategy files. This test runs on every PR that touches the `languages/` directory or `embed.go`. |

## Success Criteria

- [ ] `.forge/config.yaml` no longer contains `test-profiles` or `capabilities` fields; `ForgeConfig` struct has `Interfaces` and `Languages` fields only
- [ ] A Go project with `project-type: backend` and no other config fields generates and runs tests end-to-end via `forge testing get generate` -> skill execution -> `forge testing get run` -> skill execution
- [ ] `forge testing detect` correctly identifies all 6 language/platform types: go (via `go.mod`), javascript (via `package.json` + `@playwright/test`), rust (via `Cargo.toml`), python (via `pyproject.toml` or `requirements.txt` containing "pytest"), java (via `pom.xml` or `build.gradle`), mobile (via `android/` or `ios/` directory)
- [ ] Multi-language project (e.g., `go.mod` + `package.json`): `forge testing detect` returns 2 languages; `forge testing get generate --language go` returns the Go strategy; `forge testing get generate --language javascript` returns the Playwright strategy
- [ ] `languages` override works: when `languages: [go]` is set in config.yaml and the project has both `go.mod` and `package.json`, `forge testing detect` returns only `go`
- [ ] `interfaces` defaulting works: when `interfaces` is omitted, `forge testing interfaces` returns the union of all capabilities for the detected language(s) (e.g., Go project returns `api, cli`; JavaScript project returns `web-ui, api`)
- [ ] `forge testing get generate/run/graduate` output is functionally equivalent to the corresponding v2 `forge profile get <name> --generate/run/graduate` output for each of the 6 languages — verified by normalizing both outputs (strip trailing whitespace, normalize line endings, ignore version strings in headers) and comparing the result
- [ ] All skills call `forge testing` CLI, with zero remaining `forge profile` calls (verified by `grep -r "forge profile" plugins/` returning empty)
- [ ] Zero occurrences of `profile` or `capability` in Go exported symbols, CLI help text, and skill instruction text (verified by `grep -rE "(profile|capability)" forge-cli/pkg/ forge-cli/internal/ plugins/` excluding test files and comments)
- [ ] Migration guide exists at `docs/proposals/simplify-testing-model/migration-v2-to-v3.md` covering: field mapping, profile-name-to-language mapping, common override patterns, and troubleshooting detection failures
- [ ] Benchmark: `DetectLanguages()` completes in < 100ms on a project with all 6 signal files present (measured via Go benchmark test in `detect_test.go`)
- [ ] E2E: each `forge testing` subcommand (`detect`, `get generate`, `get run`, `get graduate`, `get justfile`, `interfaces`) passes against a test fixture project with known language (Go) — verified by CLI integration test
- [ ] Error behavior: when no language is detected and no `languages` override is set in config.yaml, `forge testing get generate` exits with non-zero code and stderr contains the string "languages" (directing user to add the override)
- [ ] Extensibility: a developer can add a hypothetical 7th language by completing 3 steps — (1) add detection case in `detect.go`, (2) add `languages/<key>/` directory with strategy files, (3) add entry to `languageCapabilities` map — without modifying any existing language's code or strategy files (verified by code review checklist)

## Core Design Decisions

### D1. 用户配置模型

```yaml
project-type: backend          # 保留，用于 scope resolution
interfaces: [api, cli]         # 替代 capabilities，系统对外接口类型
# languages: [go]              # 可选覆盖，默认自动检测
```

`interfaces` 可选值（封闭枚举，与 v2 capabilities 值一致）：web-ui, tui, mobile-ui, api, cli。

省略 `interfaces` 时，默认使用检测到的语言所支持的全部接口类型。

### D2. 自动检测 → 语言 key

| 检测信号 | 语言 key | 对应 v2 profile |
|----------|----------|-----------------|
| `go.mod` | go | go-test |
| `package.json` + playwright | javascript | web-playwright |
| `Cargo.toml` | rust | rust-test |
| `pyproject.toml`/`requirements.txt` + pytest | python | pytest |
| `pom.xml`/`build.gradle` | java | java-junit |
| `android/` 或 `ios/` 目录 | mobile | maestro |

检测到的语言 key 是内部策略包的唯一标识。

### D3. 内部策略目录

```
forge-cli/pkg/testing/
  detect.go              # 语言检测（从 detect.go 迁移）
  embed.go               # 嵌入策略文件（从 embed.go 迁移）
  config.go              # 配置读取（从 config.go 迁移）
  languages/
    go/
      generate.md
      run.md
      graduate.md
      justfile-recipes
      templates/
        test-file.go
        helpers.go
    javascript/
      ...
    python/
      ...
    java/
      ...
    rust/
      ...
    mobile/
      generate.md
      run.md
      graduate.md
      justfile-recipes
      templates/
        ...
```

每个语言目录包含完整的测试策略文件，无需 manifest.yaml。语言名即 key，元数据（file-extension, test-directory 等）直接硬编码在 Go 代码或策略文件中。

### D4. CLI 命令：forge testing

```
forge testing detect              → 输出检测到的语言（可多个）
forge testing get generate        → generate.md（自动检测语言）
forge testing get run             → run.md
forge testing get graduate        → graduate.md
forge testing get justfile        → justfile-recipes
forge testing get template <file> → 指定模板文件
forge testing interfaces          → 当前项目的 interfaces（config > 检测到的语言的默认值）
```

多语言项目：`forge testing get generate --language go` 指定语言，不指定则返回检测到的第一个。

### D5. Multi-language Handling

- Default: all auto-detected languages are activated
- User override: `languages: [go]` to explicitly specify
- Per-feature narrowing: deferred to a future release (not in v3.0 scope)

### D6. 对 Skills 的影响

| Skill | 改动 |
|-------|------|
| gen-test-scripts | `forge profile get <name> --generate` → `forge testing get generate` |
| run-e2e-tests | `forge profile get <name> --run` → `forge testing get run` |
| graduate-tests | `forge profile get <name> --graduate` → `forge testing get graduate` |
| gen-test-cases | `capabilities` → `interfaces` |
| eval-test-cases | 动态评分维度从 capabilities → interfaces |
| breakdown-tasks | 不再展开 per-profile 任务，改为 per-language |
| quick-tasks | 同 breakdown-tasks |
| init-justfile | `forge profile get <name> --justfile` → `forge testing get justfile` |
| tech-design | 删除 profile 选择步骤，改为自动检测 |
| /quick | 同 tech-design |

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
