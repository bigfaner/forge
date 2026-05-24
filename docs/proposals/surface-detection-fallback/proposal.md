---
created: 2026-05-24
author: faner
status: Draft
---

# Proposal: Surface Detection Fallback & Init Display Improvement

## Problem

Surface detection relies solely on framework dependency signals. When a project's `go.mod`/`package.json`/`pyproject.toml` lacks recognized frameworks, detection returns empty and silently falls back to TUI manual entry — with no explanation to the user about why detection failed. Additionally, `forge init` summary shows cryptic output like `CREATED surfaces (1 mappings)` instead of the actual detected types.

### Evidence

- forge-cli itself: `go.mod` has cobra + bubbletea → detection works, but the init summary just says "(1 mappings)"
- Pure stdlib Go project (only `net/http`, no framework): detection returns empty, user must manually type surface type
- User reported confusion at `forge init` output: no indication of what was actually detected or why detection failed

### Urgency

Surface detection is the first interactive experience in `forge init`. A confusing or silent detection step undermines trust in the tool. The fix is small-scope (two functions + summary format) and unblocks better onboarding.

## Proposed Solution

Two changes:

1. **Structural inference fallback**: when dependency signals return empty but manifest files exist (go.mod, package.json, pyproject.toml), infer surface type from project directory structure and conventions instead of immediately falling back to manual entry.

2. **Init summary display**: show actual detected surface types in the init summary instead of opaque "N mappings". When detection falls back to structural inference, annotate the output so the user knows the source.

### Innovation Highlights

No special innovation — this is standard convention-over-configuration. The insight is that **project structure is a reliable secondary signal**: a Go project with `cmd/` subdirectories is almost certainly a CLI, and a Python project with `[project.scripts]` in `pyproject.toml` is a CLI tool. These structural patterns are ecosystem conventions, not guesswork.

## Requirements Analysis

### Key Scenarios

1. **Go stdlib project**: `go.mod` exists but no framework deps → `cmd/` with subdirs → `cli`; `main.go` with `net/http` import → `api`
2. **Node.js minimal project**: `package.json` exists but no framework deps → `bin` field → `cli`; `index.html` at root → `web`
3. **Python minimal project**: `pyproject.toml` exists but no framework deps → `[project.scripts]` → `cli`
4. **Subdir detection**: `forge init` at project root detects `forge-cli/cli` → summary shows `forge-cli=cli`, not "(1 mappings)"
5. **Detection with signals (existing)**: cobra detected → summary shows `cli (from cobra)` — keeps working
6. **No manifest at all**: detection returns empty → manual TUI entry (unchanged behavior)

### Non-Functional Requirements

- Structural inference must be fast (no external network calls, no AST parsing)
- Inference confidence is lower than dependency signals — annotate output to distinguish inference vs detection
- Backward compatible: existing configs and workflows unchanged

### Constraints & Dependencies

- Inference rules are per-ecosystem (Go, Node.js, Python) — each needs its own heuristic set
- No new dependencies — pure filesystem checks
- `forge init` summary format change must not break scripted consumers (if any)

## Alternatives & Industry Benchmarking

### Industry Solutions

Similar tools use a layered approach: dependency signals → file structure → ask the user. For example, `npm init` infers project type from existing files before prompting.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero cost | User confusion persists, onboarding friction | Rejected: user-reported pain point |
| Expand signal maps only | — | Simple | Doesn't solve stdlib case, endless whack-a-mole | Partial: add common ones but not sufficient alone |
| Full AST/import analysis | Static analysis tools | Very accurate | Heavy, fragile, overkill for init UX | Rejected: complexity disproportionate to value |
| **Structural inference + display fix** | Convention-over-configuration | Covers stdlib case, improves all detection UX | Inference is probabilistic (annotated in output) | **Selected: best ROI** |

## Feasibility Assessment

### Technical Feasibility

All inference is filesystem-based (check directory structure, read specific config fields). No new dependencies. The inference functions live in `detect_surface.go` alongside existing signal maps.

### Resource & Timeline

2-3 coding tasks: inference rules, summary display, tests. Small scope.

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

- Structural inference fallback for Go, Node.js, Python when dependency signals are empty
- `forge init` summary format improvement: show detected types with source annotation
- Inference source annotation in `DetectResult` (e.g., `Source: "dependency"` vs `Source: "inference"`)

### Out of Scope

- Expanding dependency signal maps (separate effort)
- Rust/Cargo.toml structural inference
- Import-level analysis (parsing `net/http` from Go source)
- Changes to `forge surfaces` CLI command output format
- Mobile/TUI project structural inference (too rare to heuristic)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Structural inference gives wrong type | M | L | Annotate source in output; user can always edit in TUI |
| Inference rules become stale as ecosystems evolve | L | L | Rules are simple filesystem checks, easy to update |
| Init summary format change breaks scripted consumers | L | M | Summary is human-readable text, not machine-parsed |

## Success Criteria

- [ ] Go project with only `cmd/` subdirs → detected as `cli` without manual entry
- [ ] Go project with no signals at all → init summary shows "(inferred: cli)" or similar annotation
- [ ] `forge init` summary shows actual surface types: `CREATED surfaces cli` or `CREATED surfaces forge-cli=cli` instead of `(1 mappings)`
- [ ] Existing detection tests (27 tests) still pass unchanged
- [ ] New inference tests cover Go, Node.js, Python structural patterns

## Next Steps

- Proceed to `/quick` (this is a small-scope improvement, no PRD needed)
