---
status: "completed"
started: "2026-05-24 20:39"
completed: "2026-05-24 20:54"
time_spent: "~15m"
---

# Task Record: 1 Add structural inference fallback and Sources map

## Summary
Added structural inference fallback (inferGoSurface, inferNodeSurface, inferPythonSurface) and Sources map to DetectResult. Inference fires only when ALL dependency signals return empty, determining ecosystem by manifest file presence.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/forgeconfig/detect_surface.go
- forge-cli/pkg/forgeconfig/detect_surface_test.go

### Key Decisions
- Sources map populated via new detectSurfaceAtDirWithSources() internal function; existing detectSurfaceAtDirWithConflicts() delegates to it for backward compatibility
- Ecosystem inference priority: go.mod > package.json > pyproject.toml/setup.py (matches dependency check order)
- Go dependency source uses LastIndex to extract short name (e.g., 'cobra' from 'github.com/spf13/cobra')
- Python library exclusion checks both [project.packages] and [tool.setuptools.packages.find] as spec requires
- scanSubdirsWithSources() is a new parallel to scanSubdirsWithConflicts() that also tracks Sources

## Test Results
- **Tests Executed**: Yes
- **Passed**: 259
- **Failed**: 0
- **Coverage**: 89.3%

## Acceptance Criteria
- [x] DetectResult struct has Sources map[string]string field; zero-value nil is backward-compatible
- [x] inferGoSurface: cmd/ subdirs -> cli; api/ or handler/ -> api; both -> api wins, cli discarded
- [x] inferNodeSurface: bin field in package.json -> cli; index.html at root -> web; no subdir scan
- [x] inferPythonSurface: [project.scripts] or setup.py entry_points -> cli; app.py/main.py -> cli with library exclusion
- [x] Priority chain: inference called ONLY when ALL dependency signals return empty
- [x] Each inference function returns (surfaceType, sourceAnnotation) with inference:<rule-id> pattern
- [x] Ecosystem determined by manifest file presence; only matching ecosystem inference function called
- [x] Sources map populated: inference:cmd-dir for inferred, dependency:cobra for detected
- [x] Filesystem error resilience: unreadable directory -> empty result, no panic
- [x] Malformed manifest handling: invalid JSON/TOML -> recover-and-return-empty, no crash
- [x] Performance: all inference functions use filesystem stat + directory listing only
- [x] All existing detection tests pass unchanged

## Notes
Pre-existing lint issue (unparam on writeGoMod helper) was not introduced by this change. 26 new test cases added across 8 top-level test functions.
