---
id: "1"
title: "Add structural inference fallback and Sources map"
priority: "P0"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: true
type: "coding.feature"
mainSession: false
---

# 1: Add structural inference fallback and Sources map

## Description
When dependency signals return empty but manifest files exist (go.mod, package.json, pyproject.toml), surface detection silently falls back to manual TUI entry with no explanation. Add structural inference as a secondary signal: inspect project directory structure and conventions to infer surface type. Add a `Sources map[string]string` field to `DetectResult` to track how each surface was determined (dependency vs inference).

## Reference Files
- `proposal.md#Proposed-Solution` — defines the three inference functions (inferGoSurface, inferNodeSurface, inferPythonSurface), their rules, and the Sources map schema
- `proposal.md#Requirements-Analysis` — Key Scenarios 1-3 for Go/Node.js/Python inference rules and priority chain
- `proposal.md#Constraints-Dependencies` — ecosystem determination logic, performance budget, and no-new-dependencies constraint
- `proposal.md#Key-Risks` — filesystem error resilience, malformed manifest handling, inference suppression risk
- `forge-cli/pkg/forgeconfig/detect_surface.go` — existing detection infrastructure (DetectResult struct at L138, detectSurfaceAtDirWithConflicts at L323)

## Acceptance Criteria
- [ ] `DetectResult` struct has `Sources map[string]string` field; zero-value (`nil`) is backward-compatible with all existing construction sites
- [ ] `inferGoSurface`: `cmd/` subdirectories → `cli`; `api/` or `handler/` directory → `api`; both present → `api` wins, cli candidate discarded entirely
- [ ] `inferNodeSurface`: `bin` field in `package.json` → `cli`; `index.html` at project root (same dir as package.json) → `web`; does NOT scan subdirectories for `index.html`
- [ ] `inferPythonSurface`: `[project.scripts]` in `pyproject.toml` or `setup.py` `entry_points` → `cli`; `app.py`/`main.py` at root → `cli` ONLY when no `setup.py` with matching name AND no `[project.packages]`/`[tool.setuptools.packages.find]` library markers
- [ ] Priority chain enforced: inference functions called ONLY when ALL dependency signals return empty; when dependency signals return non-empty, inference is never called
- [ ] Each inference function returns `(surfaceType string, sourceAnnotation string)`; annotation follows pattern `inference:<rule-id>` (e.g., `inference:cmd-dir`, `inference:api-dir`, `inference:bin-field`)
- [ ] Ecosystem determined by manifest file presence: `go.mod` → Go, `package.json` → Node.js, `pyproject.toml` → Python; only matching ecosystem's inference function is called
- [ ] Sources map populated: `{"forge-cli/cli": "inference:cmd-dir"}` for inferred, `{"forge-cli/cli": "dependency:cobra"}` for detected
- [ ] Filesystem error resilience: unreadable directory (permission denied) → empty result, no panic; `forge init` completes with manual TUI fallback
- [ ] Malformed manifest handling: invalid JSON/TOML in package.json/pyproject.toml → recover-and-return-empty, no crash
- [ ] Performance: all inference functions combined complete in <50ms on a project with 50+ directories (filesystem stat + directory listing only)
- [ ] All 27 existing detection tests pass unchanged

## Hard Rules
- No external dependencies — pure filesystem checks (os.Stat, os.ReadDir, manifest parsing)
- No AST parsing or source code file content reads beyond manifest fields
- Inference functions must be stateless — pure filesystem reads with no side effects

## Implementation Notes
- Inference functions live in `detect_surface.go` alongside existing signal maps
- `detectSurfaceAtDirWithConflicts` at L323 calls ecosystem-specific inference only after dependency signals return empty
- Wrap all manifest parsing in `defer recover()` handlers that return empty on parse errors
- Python exclusion logic: before inferring `cli` from `app.py`/`main.py`, check for `setup.py` with `name` matching directory AND `[project.packages]`/`[tool.setuptools.packages.find]` — these indicate library-only packages
- For Node.js `bin` field: parse `package.json`, check both string form (`"bin": "foo"`) and object form (`"bin": {"foo": "./bin/foo"}`)
- Sources map annotations: `inference:cmd-dir`, `inference:api-dir`, `inference:bin-field`, `inference:index-html`, `inference:py-scripts`, `inference:py-main`; `dependency:<framework>` for detected paths
- Each inference function signature: `func inferXxxSurface(dir string) (surfaceType string, source string)`
