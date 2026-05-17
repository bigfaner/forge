---
id: "3"
title: "Rename DetectProfiles to DetectLanguages with language keys"
priority: "P0"
estimated_time: "1.5h"
dependencies: ["1"]
scope: "backend"
breaking: true
type: "refactor"
mainSession: false
---

# 3: Rename DetectProfiles to DetectLanguages with language keys

## Description
Rename `DetectProfiles()` to `DetectLanguages()`. Change return type from `[]string` (profile names like "go-test") to `[]Language` where `type Language string`. Output language keys (go, javascript, python, java, rust, mobile) instead of profile names. JavaScript detection only triggers when `package.json` has `@playwright/test` in devDependencies — a `package.json` without Playwright does NOT produce a javascript result.

## Reference Files
- `docs/proposals/simplify-testing-model/proposal.md` — Source proposal (D2: detection signals)
- `forge-cli/pkg/profile/detect.go` — Current detection logic
- `forge-cli/pkg/profile/detect_test.go` — Existing detection tests

## Acceptance Criteria
- `DetectLanguages()` function exists with return type `[]Language`
- `type Language string` type defined
- Detection signals produce correct language keys: go.mod→go, package.json+playwright→javascript, Cargo.toml→rust, pyproject.toml/requirements.txt+pytest→python, pom.xml/build.gradle→java, android|ios dir→mobile
- Existing detection tests updated to expect language keys instead of profile names
- All tests pass: `go test ./forge-cli/pkg/profile/...` (or `./forge-cli/pkg/testing/...` if package moved)
- `go build ./...` passes

## Hard Rules
- Detection reads only fixed filenames at project root (`os.ReadFile` on known paths) — no directory traversal, no glob patterns
- JavaScript detection requires `@playwright/test` in devDependencies (not just `package.json` existence)
- A `package.json` without Playwright dependency does NOT produce javascript detection — this prevents false positives from lint-only JS tooling
- Detection signals are additive (not exclusive): if both `go.mod` and `package.json`+playwright exist, both languages are returned

## Implementation Notes
- The detection logic itself is unchanged — only the return values change from profile names to language keys
- Consider defining `Language` constants: `LanguageGo Language = "go"`, etc. for type safety
- The function should live in the new `testing` package if Task 2 moved the package
