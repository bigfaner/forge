---
id: "8"
title: "Migrate forge-cli/tests/e2e/ Journeys: forge-commands + skill-ops"
priority: "P1"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 8: Migrate forge-cli/tests/e2e/ Journeys: forge-commands + skill-ops

## Description

Migrate 8 test files into 2 Journey directories.

**forge-commands Journey:**
- `discovery_cli_test.go` → `forge-cli/tests/forge-commands/discovery_test.go`
- `e2e_cli_test.go` → `forge-cli/tests/forge-commands/e2e_commands_test.go`
- `features/forge-info-commands/` → `forge-cli/tests/forge-commands/forge_info_commands_test.go`
- `features/forge-init-install-just/` (3 files) → `forge-cli/tests/forge-commands/forge_init_install_just_test.go` (merge)
- Contracts: `step-1-discovery.md`, `step-2-info-commands.md`, `step-3-e2e-runner.md`

**skill-ops Journey:**
- `plugin_content_cli_test.go` → `forge-cli/tests/skill-ops/plugin_content_test.go`
- `clean_code_skill_cli_test.go` → `forge-cli/tests/skill-ops/clean_code_skill_test.go`
- `forensic_cli_test.go` → `forge-cli/tests/skill-ops/forensic_test.go`
- `prompt_cli_test.go` → `forge-cli/tests/skill-ops/prompt_test.go`
- Contracts: `step-1-plugin-validation.md`, `step-2-forensic.md`, `step-3-prompt.md`

## Reference Files
- `forge-cli/tests/e2e/testkit/` — Already existing testkit
- All source files listed above

## Acceptance Criteria
- [ ] 2 Journey directories with 8 migrated test files
- [ ] Each Journey has `contracts/` with spec files and `main_test.go` via testkit
- [ ] Tests pass: `go test ./forge-cli/tests/forge-commands/... ./forge-cli/tests/skill-ops/... -tags=e2e -count=1`
- [ ] forge-init-install-just 3 files (api/cli/tui) properly merged

## Hard Rules
- Package names: `forgecommands`, `skillops`
- All test functions use `//go:build e2e` build tag
- forge-init-install-just has 3 separate test files (api, cli, tui) — merge into single file by test interface

## Implementation Notes
- forge-init-install-just tests validate 3 interfaces (API, CLI, TUI) — can merge into one file since they test the same feature through different interfaces
- plugin_content_cli_test.go has a single large test that validates all skill/agent/command files
