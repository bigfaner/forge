---
status: "completed"
started: "2026-05-17 14:29"
completed: "2026-05-17 14:39"
time_spent: "~10m"
---

# Task Record: 3 Rename DetectProfiles to DetectLanguages with language keys

## Summary
Renamed DetectProfiles to DetectLanguages with Language type. Added Language string type with constants (LanguageGo, LanguageJavaScript, etc.). Removed package.json-without-playwright JavaScript fallback and playwright.config.* glob detection per Hard Rules. Updated all callers in config.go and profile.go CLI command. Updated tests to use Language type and new detection semantics.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/profile/detect.go
- forge-cli/pkg/profile/detect_test.go
- forge-cli/pkg/profile/config.go
- forge-cli/internal/cmd/profile.go

### Key Decisions
- Defined Language type as `type Language string` with typed constants for compile-time safety
- Removed playwright.config.* glob detection entirely — Hard Rule requires only os.ReadFile on fixed filenames, no glob patterns
- Removed package.json-without-playwright fallback — Hard Rule: package.json without @playwright/test does NOT produce javascript detection
- Added languagesToStrings helper that preserves nil (nil Language slice -> nil string slice) to maintain backward-compatible behavior for ReadLanguages callers
- Added detectLanguagesAsStrings helper to encapsulate DetectLanguages + conversion for config.go callers

## Test Results
- **Tests Executed**: Yes
- **Passed**: 18
- **Failed**: 0
- **Coverage**: 76.6%

## Acceptance Criteria
- [x] DetectLanguages() function exists with return type []Language
- [x] type Language string type defined
- [x] Detection signals produce correct language keys: go.mod→go, package.json+playwright→javascript, Cargo.toml→rust, pyproject.toml/requirements.txt+pytest→python, pom.xml/build.gradle→java, android|ios dir→mobile
- [x] Existing detection tests updated to expect language keys instead of profile names
- [x] All tests pass: go test ./forge-cli/pkg/profile/...
- [x] go build ./... passes

## Notes
Coverage at 76.6% for profile package — gaps are in pre-existing AutoConfig/ReadInterfaces/GetConfigValue functions unrelated to detection logic. Full suite (all packages) passes. Language constants (LanguageGo, etc.) defined in detect.go alongside the type.
