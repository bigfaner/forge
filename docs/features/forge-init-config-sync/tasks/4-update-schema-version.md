---
id: "4"
title: "Update example YAML, JSON schema, and bump version"
priority: "P2"
estimated_time: "30m"
dependencies: ["3"]
scope: "backend"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 4: Update example YAML, JSON schema, and bump version

## Description

Remove `test-command` from the example config and JSON schema to match the code removal in Task 3. Also bump the CLI version per semver (minor bump for removed config field).

## Reference Files
- `docs/proposals/forge-init-config-sync/proposal.md` — Source proposal
- `forge-cli/internal/cmd/testdata/forge-config.example.yaml` — `test-command` at line 27
- `forge-cli/internal/cmd/testdata/forge-config.schema.json` — `test-command` at lines 48-51
- `forge-cli/scripts/version.txt` — Current version 4.4.3

## Acceptance Criteria
- [ ] `test-command` removed from `forge-config.example.yaml` (line 27, commented out section)
- [ ] `test-command` removed from `forge-config.schema.json` (lines 48-51)
- [ ] Version bumped in `forge-cli/scripts/version.txt` following semver (minor: 4.4.3 → 4.5.0, since removing a config field is a backward-compatible removal)

## Hard Rules
- Only touch the three listed files
- Verify JSON schema is valid after edit

## Implementation Notes
- The example YAML has `test-command` commented out — remove the comment line entirely.
- The JSON schema entry for `test-command` is a `{"type": "string", "description": "..."}` object — remove it from the `properties` object.
