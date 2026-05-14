---
id: "1"
title: "Config commands and schema extension"
priority: "P1"
estimated_time: "1.5h"
dependencies: []
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 1: Config commands and schema extension

## Description

Extend `.forge/config.yaml` schema with `project-type` and `capabilities` fields. Implement `forge config init` for interactive setup and `forge config get <key>` for agent-friendly queries.

Currently, `ForgeConfig` only has `TestProfiles`. The config command group will be a new parent (`forge config`) with two subcommands: `init` (interactive wizard) and `get` (programmatic query).

## Reference Files
- `docs/proposals/forge-info-commands/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|
| `forge-cli/internal/cmd/config.go` | `forge config` parent + `init` + `get` subcommands |
| `forge-cli/internal/cmd/config_test.go` | Unit tests for config commands |

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/profile/config.go` | Add `ProjectType` and `Capabilities` fields to `ForgeConfig`; add `ReadConfig` and `GetConfigValue` functions |
| `forge-cli/pkg/profile/config_test.go` | Test new config fields and get-value logic |
| `forge-cli/internal/cmd/root.go` | Register `configCmd` in root |

## Acceptance Criteria

- [ ] `forge config init` interactively collects project-type (frontend/backend/mixed), test-profiles (multi-select from KnownProfiles), and capabilities (union from selected profiles, user picks subset)
- [ ] `forge config init` writes `.forge/config.yaml` with all three fields
- [ ] `forge config init` prompts to reconfigure if config already exists
- [ ] `forge config get project-type` returns the value as plain text
- [ ] `forge config get capabilities` returns each value on a new line
- [ ] `forge config get <key>` exits with code 1 and no output when key doesn't exist
- [ ] `ForgeConfig` struct has `ProjectType string`, `TestProfiles []string`, `Capabilities []string` fields
- [ ] Test coverage ≥ 80% for new and modified code

## Hard Rules

- `forge config get` output must be plain text (no formatting blocks) — agents parse this with subprocess capture
- Interactive init uses stdin reading only (no bubbletea/heavy deps)

## Implementation Notes

- `config get` should support arbitrary keys for forward compatibility — if the key exists as a YAML field, return it
- For array values, output one item per line; for scalar values, output the raw string
- The `init` wizard for capabilities: read selected profiles' manifests via `profile.GetManifest()`, extract `capabilities` arrays, take union, present as multi-select
- If `.forge/config.yaml` already exists during `init`, prompt "Config already exists. Reconfigure? [y/N]" before proceeding
