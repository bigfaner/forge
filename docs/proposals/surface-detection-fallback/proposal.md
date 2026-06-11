---
created: 2026-05-24
author: faner
status: Draft
---

# Proposal: Surface Detection Fallback & Init Display Improvement

## Problem

Two distinct but related problems:

**Sub-problem 1: Silent fallback on stdlib projects.** Surface detection relies solely on framework dependency signals. When a project's `go.mod`/`package.json`/`pyproject.toml` lacks recognized frameworks, detection returns empty and silently falls back to TUI manual entry — with no explanation to the user about why detection failed. `forge init` summary shows cryptic output like `CREATED surfaces (1 mappings)` instead of the actual detected types.

**Sub-problem 2: No explicit re-detection entry point.** Users who want to re-run surface detection after project changes (e.g., adding a `cmd/` directory) must re-run the entire `forge init` flow. There is no standalone command for detection.

This proposal addresses both problems across two touchpoints: (1) the `forge init` flow — adding structural inference as a fallback when dependency signals return empty, and improving the init summary to show actual detected types with source annotations; (2) a new `forge surfaces detect` subcommand — giving users an explicit entry point for re-detection without re-running the full init flow. Both touchpoints share the same detection + inference pipeline and TUI confirmation infrastructure, making them natural to build together.

### Evidence

- forge-cli itself: `go.mod` has cobra + bubbletea → detection works, but the init summary just says "(1 mappings)"
- Pure stdlib Go project (only `net/http`, no framework): detection returns empty, user must manually type surface type
- Internal testing (forge v2.x onboarding walkthrough, 2026-04): 3 of 5 test subjects with stdlib Go projects asked "what did it detect?" after seeing "(1 mappings)" with no indication of what was inferred or why detection failed. Node.js and Python stdlib projects exhibit the same UX pattern (empty detection → silent manual fallback); evidence is inferred from the Go testing rather than separately reproduced per ecosystem.

### Urgency

Surface detection is the first interactive experience in `forge init`. A confusing or silent detection step undermines trust in the tool. In the last 10 `forge init` runs on internal projects, 4 produced empty detection results and required manual surface entry. Each failed detection adds friction that risks user abandonment at the onboarding stage.

## Proposed Solution

Three changes:

1. **Structural inference fallback**: when dependency signals return empty but manifest files exist (go.mod, package.json, pyproject.toml), infer surface type from project directory structure and conventions instead of immediately falling back to manual entry. Ecosystem is determined by which manifest file exists (`go.mod` → Go, `package.json` → Node.js, `pyproject.toml` → Python); only the matching ecosystem's inference function is called. Internally, `detectSurfaceAtDirWithConflicts` calls ecosystem-specific inference functions (e.g., `inferGoSurface`, `inferNodeSurface`, `inferPythonSurface`) only after all dependency signal maps return empty. Each inference function returns `(surfaceType, sourceAnnotation)` and is stateless — pure filesystem reads with no side effects. Results are merged into the existing `DetectResult.Sources` map keyed by path. Key inference rules per ecosystem: Go — `cmd/` subdirectories → `cli`, `api/` or `handler/` directory → `api`, both present → `api` wins; Node.js — `bin` field in `package.json` → `cli`, `index.html` at project root → `web`; Python — `[project.scripts]` in `pyproject.toml` or `setup.py` `entry_points` → `cli`, `app.py`/`main.py` at root → `cli` (with library exclusion checks).

2. **Init summary display**: show actual detected surface types in the init summary instead of opaque "N mappings". When detection falls back to structural inference, annotate the output so the user knows the source.

   Before:
   ```
   CREATED surfaces (1 mappings)
   ```
   After:
   ```
   CREATED surfaces cli (inferred from cmd/ directory structure)
   ```

3. **TUI confirmation flow update**: `askSurfaceConfirmation` will display the inference source in the Description field alongside each surface type (e.g., `"cli (inferred from cmd/ directory structure)"` vs `"cli (detected from cobra dependency)"`). This gives users context to judge whether to accept or override the inference. When the user edits or overrides an inferred value, the override replaces the value directly — no source annotation is persisted to config. The TUI shows a brief hint: `"This was inferred from project structure. Edit to correct if needed."` for inferred entries. Source information is display-only and discarded after TUI confirmation completes. Multi-surface example — when subdir detection finds multiple surfaces, the TUI shows a map-form list:
   ```
   Detected surfaces:
     forge-cli/cli  →  cli  (inferred from cmd/ directory structure)
     forge-cli/api  →  api  (inferred from api/ directory)
   [Accept] [Edit] [Skip]
   ```

4. **`forge surfaces detect` command**: new subcommand that runs detection + inference and displays results, independent of `forge init`. Gives users an explicit entry point for re-detection without re-running the full init flow. Default behavior is **read-only**: runs `DetectSurfacesWithConflicts` with inference, shows results with source annotations, exits without writing. Use `--apply` flag to enable config writing: shows TUI confirmation (same flow as init) and writes to config on confirm. Non-interactive terminals print results to stdout and exit (no TUI, no config write regardless of flags).

### Innovation Highlights

No special innovation — this is standard convention-over-configuration. The insight is that **project structure is a reliable secondary signal**: a Go project with `cmd/` subdirectories (but without `internal/` and `pkg/` coexisting — a pattern typical of microservices, not CLIs) is almost certainly a CLI, and a Python project with `[project.scripts]` in `pyproject.toml` is a CLI tool. These structural patterns are ecosystem conventions, not guesswork.

## Requirements Analysis

### Key Scenarios

1. **Go stdlib project**: `go.mod` exists but no framework deps → `cmd/` with subdirs → `cli`; `api/` or `handler/` directory present → `api`; if both `cmd/` and `api/` exist → `api` wins and only `api` is returned (the losing candidate is discarded entirely, not stored)
2. **Node.js minimal project**: `package.json` exists but no framework deps → `bin` field → `cli`; `index.html` at project root (same directory as `package.json`) → `web`; `inferNodeSurface` does not scan subdirectories for `index.html`
3. **Python minimal project**: `pyproject.toml` exists but no framework deps → `[project.scripts]` or `setup.py` `entry_points` → `cli`; `app.py`/`main.py` at root → `cli` **only when no `setup.py` with `name` matching the directory AND no `[project.packages]`/`[tool.setuptools.packages.find]` library-only markers** (avoids false positives on library packages). **Known false-negative boundary**: the exclusion condition assumes `[project.packages]`/`[tool.setuptools.packages.find]` implies library-only, but many legitimate CLI tools also declare package discovery for `pip install` — these projects will not be inferred as `cli` and will fall back to manual entry. This is an acceptable trade-off: over-exclusion (missing a CLI) is safer than over-inclusion (labeling a library as CLI).
4. **Subdir detection**: `forge init` at project root detects `forge-cli/cli` → summary shows `forge-cli=cli`, not "(1 mappings)"
5. **Detection with signals (existing)**: cobra detected → summary shows `cli (from cobra)` — keeps working
6. **No manifest at all**: detection returns empty → manual TUI entry (unchanged behavior)
7. **Re-run with existing config**: `config.yaml` already has `surfaces: cli` → TUI shows `"Surfaces already configured: cli. Re-detect?"` with Confirm (keep existing) / Re-detect / Edit (enter `manualSurfaceEntry` flow to modify manually) options. Confirm skips detection; Re-detect runs full detection + inference pipeline and shows TUI confirmation as usual; Edit opens the manual surface entry flow
8. **User overrides inference in TUI**: user sees `"cli (inferred from cmd/ directory structure)"`, selects Edit, types `api` → config receives `surfaces: api` with no source annotation (source is TUI-only, not persisted)
9. **`forge surfaces detect`**: runs detection + inference in read-only mode (shows results with source annotations, no config write). Interactive terminals show TUI confirmation with option to apply; non-interactive terminals print to stdout
10. **`forge surfaces detect --apply`**: runs detection + inference, shows results with source annotations, TUI confirmation, writes to config on confirm. Without `--apply`, the command is read-only by default
11. **`forge surfaces detect` in non-interactive terminal**: prints results to stdout in human-readable format (no TUI), no config write, exits with code 0

### Non-Functional Requirements

- Structural inference must be fast (no external network calls, no AST parsing)
- Inference confidence is lower than dependency signals — annotate output to distinguish inference vs detection
- Backward compatible: existing configs and workflows unchanged
- Inference functions return empty on any filesystem error (permission denied, unreadable directory), never crash `forge init`
- Inference performance budget: all inference functions combined must complete in <50ms (filesystem stat + directory listing only, no file content reads beyond manifest parsing already done by dependency detection)

### Constraints & Dependencies

- Inference rules are per-ecosystem (Go, Node.js, Python) — each needs its own heuristic set
- No new dependencies — pure filesystem checks
- `forge init` summary format change must not break scripted consumers (if any)
- Re-running `forge init` with existing surfaces config: `runSurfaceConfig` detects existing surfaces in config, shows TUI prompt `"Surfaces already configured: <current>. Re-detect?"` with Confirm (keep existing) / Re-detect / Edit options. Confirm skips detection and returns `SKIPPED surfaces (already configured)`. Re-detect runs the full detection + inference pipeline and presents the standard TUI confirmation flow. Edit enters `manualSurfaceEntry` flow. This eliminates the need for a `--force` flag — the user makes the choice interactively.
- Source persistence via YAML comments: when writing surfaces config, append a comment line like `surfaces: api  # source: inference:cmd-dir` for human context. This preserves source metadata across re-runs without modifying the SurfacesMap schema — the comment is purely informational and ignored by the YAML parser. On re-detection, the existing comment (if present) is displayed alongside the current value for context

## Alternatives & Industry Benchmarking

### Industry Solutions

Similar tools use a layered approach: dependency signals → file structure → ask the user. `npm init` infers project type from existing `package.json` fields before prompting — it reads `main`, `bin`, and `scripts` to pre-fill the package type, falling back to defaults when these fields are absent. Its failure mode is producing a generic `index.js` entry point that may not match the project's actual purpose. `cargo init` distinguishes `--bin` (CLI) from `--lib` (library) using a single flag, defaulting to `--bin` when a `main.rs` exists in the target directory. This binary choice avoids ambiguity but cannot detect multi-target projects (e.g., a workspace with both bins and libs). VS Code's workspace detection reads `.vscode/settings.json` for explicit configuration, then falls back to probing directory conventions — `manage.py` signals Django, `requirements.txt` suggests a Python project, `go.mod` indicates Go. The failure mode is misidentifying a monorepo subdirectory as the project root or detecting stale markers from a migrated project. Spring Boot's `spring.factories` auto-configuration uses a layered classpath scan that degrades gracefully when explicit annotations are absent, falling back to conditional `@Configuration` classes. Its failure mode is silent misconfiguration when multiple auto-configurations conflict — a problem analogous to our multi-surface detection ambiguity. The common pattern across all four tools is a strict priority chain: explicit configuration wins over convention, convention wins over defaults, and the user is prompted only when all automated signals fail. Forge's proposed approach follows this same chain: dependency detection → structural inference → manual entry. Inference signals **only activate when ALL dependency signals are empty** — they never mix with dependency results. This ensures a clean priority chain: dependency detection wins outright, structural inference is a pure fallback, and manual entry is the last resort.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero cost | User confusion persists, onboarding friction | Rejected: user-reported pain point |
| Expand signal maps only | — | Simple | Doesn't solve stdlib case, endless whack-a-mole | Partial: add common ones but not sufficient alone |
| Regex-based import scanning | grep/sed-style scan of source files | Covers more cases than filesystem checks (e.g., detects `import http` usage) | Requires reading source files, fragile across languages, false positives from unused imports | Rejected: couples init to source code parsing; breaks "no AST" constraint |
| **Structural inference + display fix** | Convention-over-configuration | Covers stdlib case, improves all detection UX | Inference is probabilistic (annotated in output) | **Selected: best ROI** — ~60-80 LOC total across 3 inference functions yields estimated coverage for ~70% of Go stdlib projects (estimate based on common Go project layouts observed in forge usage: most stdlib Go projects follow `cmd/` or `api/` layout conventions), ~50% of Node stdlib projects (estimate: `bin` field or `index.html` are common in CLI/web projects), and ~40% of Python stdlib projects (estimate: `[project.scripts]` is standard for CLI tools, but many Python projects are libraries with no structural CLI markers). Investment is low: each rule is a single filesystem check with clear true/false semantics, testable with directory fixtures, and maintainable by anyone familiar with ecosystem conventions. Compare to regex scanning: ~150 LOC across 3 languages for import syntax handling alone, with no clear true/false boundary (unused imports cause false positives), higher test complexity (must mock source file contents), and ongoing maintenance as import syntax evolves. The "do nothing" alternative has a hidden cost: every stdlib project user must manually enter surface types, adding ~30 seconds of friction per `forge init` run at the critical onboarding stage. |

## Feasibility Assessment

### Technical Feasibility

All inference is filesystem-based (check directory structure, read specific config fields). No new dependencies. The inference functions live in `detect_surface.go` alongside existing signal maps.

### Resource & Timeline

4-5 coding tasks: (1) inference rules per ecosystem, (2) Sources map integration into DetectResult, (3) TUI confirmation flow update, (4) init summary display + tests, (5) `forge surfaces detect` command with `--apply` flag. Task dependencies: tasks 1-4 can be parallelized within the same PR (shared `detect_surface.go` file but independent functions); task 5 depends on tasks 1-3 (reuses detection + inference pipeline and TUI infrastructure). Critical path: tasks 1-3 → task 5, with task 4 in parallel.

### Dependency Readiness

No external dependencies. Uses existing `DetectSurfacesWithConflicts` infrastructure.

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| Dependency signals are sufficient for all projects | 5 Whys: what about stdlib projects? → users hit empty detection on real projects | Overturned: structural inference needed as fallback |
| "1 mappings" in init summary is clear enough | XY Detection: user wants to know WHAT was detected, not how many | Refined: show actual type values and inference source |
| Inference accuracy must be 100% | Occam's Razor: annotated 80% confidence beats silent 0% | Confirmed: annotation is sufficient for user trust |

## Scope

### In Scope

- Structural inference fallback for Go, Node.js, Python when dependency signals are empty (note: Python coverage is limited to CLI detection via `[project.scripts]`/`setup.py entry_points`/`app.py` with library exclusion; Python API-type stdlib projects without Flask/FastAPI remain undetectable by structural conventions alone)
- `forge init` summary format improvement: show detected types with source annotation
- Per-path inference source annotation in `DetectResult` via `Sources map[string]string` field (e.g., `{"forge-cli/cli": "dependency", "forge-cli/api": "inference"}`), since a single `Source` field cannot express mixed sources across multiple surfaces
- `forge surfaces detect` subcommand: runs detection + inference independently of `forge init`, default read-only with `--apply` flag for config writing, and non-interactive mode

### Out of Scope

- Expanding dependency signal maps (separate effort)
- Rust/Cargo.toml structural inference (Rust's `[[bin]]`/`[lib]` targets in Cargo.toml could distinguish CLI vs library, but Rust project surface detection is not yet implemented in the existing dependency signal maps; adding inference for an ecosystem with no baseline detection support would be premature)
- Import-level analysis (parsing source code imports — excluded to keep inference purely filesystem-based)
- Changes to `forge surfaces` CLI command output format
- Mobile/TUI project structural inference (too rare to heuristic)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Structural inference gives wrong type | M | M | Annotate source in output; add integration test suite covering each inference rule with known-correct fixtures; ship with verbose logging (`FORGE_DEBUG=1`) so mis-inference is diagnosable |
| Inference rules become stale as ecosystems evolve | L | L | Rules are simple filesystem checks, easy to update; document each rule in code comments with the convention it targets |
| Init summary format change breaks scripted consumers | L | M | Summary is human-readable text, not machine-parsed; grep test in CI verifies no programmatic parser exists in codebase |
| `Sources map[string]string` schema migration | L | M | Zero-value of `map[string]string` is `nil`, so existing `DetectResult` construction sites are backward-compatible without changes; add compile-time check that all code paths populate `Sources` when inference fires |
| Re-run prompt adds friction for power users | L | L | Prompt only fires when surfaces already configured; first-run experience unchanged; one extra keypress for re-runs is acceptable tradeoff vs accidental overwrite |
| YAML comment stripping during config round-trip | M | L | Forge uses `gopkg.in/yaml.v3` which preserves comments when using `yaml.Node` round-trip API (unmarshal to Node, modify value, marshal back). Mitigation: write the comment via Node manipulation, not string concatenation. Fallback: if comment is stripped, detection still works correctly — the comment is informational only. Verification: integration test reads config file after `forge init` and asserts comment is present |
| Inference suppression masks valid surfaces | M | M | When dependency signals detect one surface (e.g., cobra → CLI) but structural inference would find another (e.g., `api/` directory → API), inference is suppressed because the priority chain requires ALL dependency signals to be empty. The user sees only the dependency-detected surface and misses the second surface. Likelihood upgraded to M: multi-surface Go projects (a CLI that also serves an API) are common enough that this pattern will occur in practice. This is an acceptable trade-off for v1: mixing inference with dependency results risks conflicting signals. Document this as a known limitation; a future iteration may support partial inference (per-surface priority chain) |
| Re-run Edit flow diverges from first-run manual entry | L | L | The Edit option in the re-run prompt calls the same `manualSurfaceEntry` function used by first-run manual entry, so the UX is identical. Risk is that future changes to one path are not applied to the other. Mitigation: shared function ensures single code path; integration test verifies Edit produces same output as first-run manual entry |
| Malformed manifest file causes inference crash | M | M | Inference functions parse manifest files (`package.json`, `pyproject.toml`) to extract fields like `bin` or `[project.scripts]`. A malformed or encoding-corrupted manifest could cause parse errors that propagate up. Mitigation: wrap all manifest parsing in recover-and-return-empty handlers; add test fixtures with malformed manifests (invalid JSON, truncated TOML, BOM-prefixed files) to verify graceful degradation |

## Success Criteria

- [ ] Go project with only `cmd/` subdirs → detected as `cli` without manual entry
- [ ] Go project with both `cmd/` and `api/` directories → only `api` returned (cli candidate discarded); `Sources` map contains `{"<path>": "inference:api-dir"}` with no cli entry
- [ ] Go project with no dependency signals but `cmd/` subdirs present → init summary shows `CREATED surfaces cli (inferred:cmd-dir)` matching the pattern `inferred:<rule-id>`
- [ ] Node.js project with `bin` field in `package.json` → detected as `cli`; project with `index.html` at root → detected as `web`
- [ ] Python project with `[project.scripts]` in `pyproject.toml` → detected as `cli`; Python project with `app.py` but also `setup.py` with `[options.packages.find]` → NOT inferred as `cli` (exclusion condition)
- [ ] `forge init` summary shows actual surface types: `CREATED surfaces cli` or `CREATED surfaces forge-cli=cli` instead of `(1 mappings)`
- [ ] Subdir detection (Key Scenario #4): `forge init` at project root with `forge-cli/cli` subdir → summary shows `CREATED surfaces forge-cli=cli (inferred:cmd-dir)`, not `(1 mappings)`
- [ ] TUI confirmation shows inference source: `"cli (inferred from cmd/ directory structure)"` vs `"cli (detected from cobra dependency)"`
- [ ] TUI hint text: when an inferred surface is displayed, the TUI shows the hint `"This was inferred from project structure. Edit to correct if needed."`; hint is absent for dependency-detected surfaces
- [ ] `DetectResult.Sources` map correctly populated: `{"forge-cli/cli": "inference:cmd-dir"}` for inferred paths, `{"forge-cli/cli": "dependency:cobra"}` for detected paths
- [ ] Priority chain verified: when dependency signals return non-empty, inference functions are never called
- [ ] Existing detection tests (27 tests) still pass unchanged
- [ ] Re-run behavior: `config.yaml` with `surfaces: cli` present → TUI shows `"Surfaces already configured: cli. Re-detect?"` with Confirm / Re-detect / Edit; Confirm returns `SKIPPED surfaces (already configured)`, Re-detect runs full detection + inference pipeline, Edit enters `manualSurfaceEntry` flow
- [ ] User override: after editing an inferred `cli` to `api` in TUI, config contains `surfaces: api` with no source field; source annotation does not appear in serialized config
- [ ] YAML comment persistence: after `forge init` infers `cli` from `cmd/` directory, config file contains `surfaces: cli  # source: inference:cmd-dir`; the comment is present in the file but not parsed by the YAML unmarshaler
- [ ] `forge surfaces detect` runs detection + inference in read-only mode and shows results with source annotations (same format as init TUI); no config write without `--apply`
- [ ] `forge surfaces detect --apply` runs detection + inference, shows TUI confirmation, writes to config on confirm; exit code 0. After confirm, config file on disk contains the detected surfaces (verified by reading the config file and asserting the surface entries match what was detected)
- [ ] `forge surfaces detect` in non-interactive terminal: prints results to stdout, no TUI, no config write (regardless of `--apply`); exit code 0 if detection succeeds, 1 if no surfaces found. Stdout format: one line per surface, `<path>=<type> (<source>)`, where `<source>` is `detected:<signal>` or `inferred:<rule-id>` matching the `inferred:<rule-id>` pattern from init summary; empty detection prints nothing to stdout and exits with code 1. **Note**: the non-interactive stdout format is unstable and subject to change between minor versions; scripted consumers should not rely on it without pinning the forge version
- [ ] Performance budget: all inference functions combined complete in <50ms on a project with 50+ directories (benchmark test with synthetic fixture)
- [ ] Filesystem error resilience: inference function called on a directory with permission denied (unreadable `cmd/` subdirectory) returns empty result without panicking; `forge init` completes successfully with manual TUI fallback
- [ ] Python exclusion boundary: project with `app.py` + `setup.py` where `name` matches directory name → NOT inferred as `cli`; project with `app.py` + `[tool.setuptools.packages.find]` → NOT inferred as `cli`; project with `app.py` alone (no `setup.py`, no library markers) → inferred as `cli`

## Next Steps

- Proceed to `/quick` (this is a small-scope improvement, no PRD needed)
