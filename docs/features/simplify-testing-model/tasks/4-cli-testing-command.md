---
id: "4"
title: "Implement forge testing CLI command group"
priority: "P0"
estimated_time: "2h"
dependencies: ["1", "2", "3"]
scope: "backend"
breaking: true
type: "feature"
mainSession: false
---

# 4: Implement forge testing CLI command group

## Description
Replace the `forge profile` command group with `forge testing`. Subcommands: `detect` (outputs detected languages), `get generate` (outputs generate.md), `get run` (outputs run.md), `get graduate` (outputs graduate.md), `get justfile` (outputs justfile-recipes), `get template <file>` (outputs specified template), `interfaces` (outputs interface types). The `get` subcommands auto-detect language when `--language` flag is not specified. When no language is detected and no `languages` override exists, exit with non-zero code and stderr containing "languages".

## Reference Files
- `docs/proposals/simplify-testing-model/proposal.md` — Source proposal (D4: CLI command spec)
- `forge-cli/internal/cmd/profile.go` — Current profile command implementation
- `forge-cli/internal/cmd/root.go` — Root command registration

## Acceptance Criteria
- `forge testing detect` outputs detected language(s) (plain text, one per line or comma-separated)
- `forge testing get generate` outputs generate.md content for auto-detected language
- `forge testing get generate --language go` outputs Go strategy specifically
- `forge testing get run`, `get graduate`, `get justfile` work for auto-detected language
- `forge testing get template <file>` returns specified template file content
- `forge testing interfaces` returns interface types (config.Interfaces > detected language defaults)
- No language detected + no `languages` override → exit code non-zero, stderr contains "languages"
- `forge profile` command removed entirely (not deprecated, not aliased)
- Output format is plain text with structured blocks, matching current `forge profile` format
- `go build ./...` and `go test ./...` pass

## Hard Rules
- CLI output format must be stable for skill parsing — same plain-text-with-structured-blocks pattern as current `forge profile`
- `forge testing detect` output must be parseable by skills (one language key per line or consistent format)
- Do not add backward-compat aliases for `forge profile` — v3.0 is a clean break

## Implementation Notes
- The resolve logic simplifies from a 3-step chain (config → detect → "none") to: read `languages` from config, or auto-detect
- Multi-language projects: `--language` flag selects which strategy to return; without it, return the first detected language
- This task depends on Tasks 1 (config), 2 (directory/embed), and 3 (detection) — all must be complete before implementation
