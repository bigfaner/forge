---
status: "completed"
started: "2026-05-17 14:00"
completed: "2026-05-17 14:00"
time_spent: ""
---

# Task Record: 1 Config schema: replace TestProfiles/Capabilities with Interfaces/Languages

## Summary
Replaced TestProfiles/Capabilities with Interfaces/Languages in ForgeConfig. Added ReadLanguages() (config override -> DetectProfiles fallback) and ReadInterfaces() (config override -> languageCapabilities map union). Global rename: capabilities/capability -> interfaces/interface in all Go exported symbols (ValidateCapabilities->ValidateInterfaces, ValidTestTypes->ValidInterfaceTypes, GetProfileCapabilities->GetProfileInterfaces, UnionCapabilities->UnionLanguageInterfaces, TestCapabilities->TestInterfaces). Added languageCapabilities hardcoded map and WriteLanguages helper. Updated all callers across pkg/profile, pkg/e2e, pkg/task, internal/cmd.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/profile/config.go
- forge-cli/pkg/profile/embed.go
- forge-cli/pkg/profile/config_test.go
- forge-cli/pkg/profile/embed_test.go
- forge-cli/pkg/e2e/e2e.go
- forge-cli/pkg/e2e/e2e_test.go
- forge-cli/pkg/e2e/actions_test.go
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/build_test.go
- forge-cli/pkg/task/testgen.go
- forge-cli/pkg/task/testgen_test.go
- forge-cli/internal/cmd/profile.go
- forge-cli/internal/cmd/config.go
- forge-cli/internal/cmd/config_test.go
- forge-cli/internal/cmd/index.go
- forge-cli/internal/cmd/add.go
- forge-cli/internal/cmd/init.go
- forge-cli/internal/cmd/init_test.go
- forge-cli/tests/e2e/features/forge-info-commands/forge_info_commands_cli_test.go

### Key Decisions
- ReadLanguages() falls back to DetectProfiles() when config.Languages is empty, keeping auto-detect as the primary path
- languageCapabilities hardcoded map provides default interfaces per profile without reading manifests, with manifest fallback for unknown profiles
- UnionLanguageInterfaces takes explicit profiles slice (no projectRoot), keeping auto-detect separate in ReadInterfaces/defaultInterfaces
- BuildIndexOpts.TestProfiles name kept unchanged since it represents resolved profiles/languages from caller
- profileManifest.Capabilities field kept as-is since it's unexported and maps to the YAML field in manifest.yaml files (profiles/ not renamed per hard rules)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 20
- **Failed**: 0
- **Coverage**: 79.0%

## Acceptance Criteria
- [x] ForgeConfig struct has Interfaces and Languages fields; TestProfiles and Capabilities removed
- [x] ReadLanguages() function exists: returns config.Languages if set, otherwise calls DetectProfiles()
- [x] ReadInterfaces() function exists: returns config.Interfaces if set, otherwise defaults to union of detected languages' capabilities
- [x] Zero Go exported symbols contain 'capability' or 'Capability'
- [x] Config YAML field names: 'interfaces', 'languages' (snake_case)
- [x] go build ./... passes

## Notes
Hard rules respected: no changes to detection logic (detect.go unchanged), no changes to embed paths (profiles/ dir unchanged), interfaces valid values unchanged from v2 capabilities set.
