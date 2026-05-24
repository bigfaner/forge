---
status: "completed"
started: "2026-05-24 13:44"
completed: "2026-05-24 14:52"
time_spent: "~1h 8m"
---

# Task Record: 1 Config struct refactor: Surfaces dual-form + remove Interfaces

## Summary
Replaced Interfaces []string with Surfaces SurfacesMap in Config struct, implemented dual-form YAML marshal/unmarshal (scalar + map), added ReadSurfaces/SurfaceTypes helpers, and implemented auto-migration logic (single-interface auto-convert, multi-interface error guidance).

## Changes

### Files Created
- forge-cli/pkg/forgeconfig/surfaces_test.go
- docs/lessons/gotcha-breaking-change-quality-gate-deadlock.md

### Files Modified
- forge-cli/pkg/forgeconfig/config.go
- forge-cli/pkg/forgeconfig/detect.go
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/build_test.go

### Key Decisions
- Created SurfacesMap custom type (map[string]string) with UnmarshalYAML/MarshalYAML for dual-form YAML support, rather than overriding the entire Config's YAML handling
- Scalar form uses '.': value' key internally; MarshalYAML detects single-dot-key to emit scalar
- Nil SurfacesMap marshals as 'surfaces: {}' to avoid omitempty silent-skip bug
- SurfaceTypes() extracts deduplicated values from surfaces map, keeping autogen.go/extract.go interfaces unchanged ([]string)
- MigrateInterfacesToSurfaces() is a standalone function for CLI startup invocation, not embedded in ReadConfig
- Included minimal build.go caller update in Task 1 scope to maintain compilation

## Test Results
- **Tests Executed**: Yes
- **Passed**: 896
- **Failed**: 0
- **Coverage**: 87.6%

## Acceptance Criteria
- [x] Interfaces []string field removed from Config struct
- [x] Surfaces map[string]string yaml:'surfaces' added (no omitempty)
- [x] Custom UnmarshalYAML: scalar 'api' -> map[string]string{'.': 'api'}; map form used as-is
- [x] Custom MarshalYAML: single entry with key '.' -> scalar; otherwise -> map
- [x] Auto-migration: single interfaces -> auto-write surfaces scalar + console prompt
- [x] Auto-migration: multi interfaces -> error exit with guidance
- [x] Empty map (0 entries) serializes as surfaces: {}, never omitted
- [x] Existing tests pass after struct change

## Notes
build.go updated to use ReadSurfaces + SurfaceTypes instead of ReadInterfaces. build_test.go helper writeForgeConfig updated from 'interfaces: - api' to 'surfaces: api'. The uiInterfaces map in autogen.go was NOT changed per task scope (task 4 handles that). Integration test TestTC_005_ConfigGetInterfacesArrayOutput in tests/forge-commands/ will need update in a follow-up task since it tests the removed 'interfaces' field via CLI.
