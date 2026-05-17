---
id: "2"
title: "Restructure profiles/ to languages/ with strategy files"
priority: "P0"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: true
type: "refactor"
mainSession: false
---

# 2: Restructure profiles/ to languages/ with strategy files

## Description
Rename `forge-cli/pkg/profile/profiles/<profile-name>/` directories to `forge-cli/pkg/testing/languages/<language-key>/`. Remove all `manifest.yaml` files from each profile directory. Update `embed.go` embed directive from `//go:embed all:profiles` to `//go:embed all:languages`. Add hardcoded `languageCapabilities` map in Go code to replace manifest metadata. Remove `KnownProfiles` constant/variable.

## Reference Files
- `docs/proposals/simplify-testing-model/proposal.md` — Source proposal (D3: Internal strategy directory)
- `forge-cli/pkg/profile/profiles/` — Current profile directories (6 profiles)
- `forge-cli/pkg/profile/embed.go` — Current embed directive and strategy lookup

## Acceptance Criteria
- `profiles/` directory no longer exists under `forge-cli/pkg/profile/` (or `forge-cli/pkg/testing/`)
- `languages/` directory contains 6 subdirectories: go, javascript, python, java, rust, mobile
- Each language subdir has: generate.md, run.md, graduate.md, justfile-recipes, templates/
- Zero `manifest.yaml` files remain in the codebase
- `embed.go` uses `//go:embed all:languages`
- `languageCapabilities` Go map defines capabilities for all 6 languages matching proposal D2 table
- `KnownProfiles` constant/variable removed
- `go build ./...` passes

## Hard Rules
- Do NOT modify strategy file content (generate.md, run.md, etc.) — only directory organization changes
- Grep all strategy files for hardcoded `profiles/` references before rename; fix any found
- Profile-to-language mapping (strict): go-test→go, web-playwright→javascript, rust-test→rust, pytest→python, java-junit→java, maestro→mobile

## Implementation Notes
- The package may need to move from `profile/` to `testing/` — coordinate with Task 3 (detection) and Task 4 (CLI) on the final package path
- `languageCapabilities` map values: go→{api,cli}, javascript→{web-ui,api}, python→{api,cli}, java→{api,cli}, rust→{api,cli}, mobile→{mobile-ui}
- Strategy files use language-relative paths internally (e.g., `test_file.go`) resolved by CLI — renaming directories should not break this
