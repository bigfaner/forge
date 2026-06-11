---
id: "5"
title: "Add forge surfaces detect subcommand"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["1", "2"]
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 5: Add forge surfaces detect subcommand

## Description
Add a new `forge surfaces detect` subcommand that runs detection + structural inference independently of `forge init`. Default behavior is read-only: shows results with source annotations, exits without writing config. Use `--apply` flag to enable TUI confirmation and config writing. Non-interactive terminals print to stdout.

## Reference Files
- `proposal.md#Proposed-Solution` — forge surfaces detect spec: read-only default, --apply flag, non-interactive mode, TUI confirmation
- `proposal.md#Requirements-Analysis` — Key Scenarios 9-11 for detect subcommand behaviors (read-only, --apply, non-interactive)
- `proposal.md#Success-Criteria` — detect subcommand criteria (#16-18)
- `proposal.md#Constraints-Dependencies` — stdout format spec and unstable-format warning
- `forge-cli/internal/cmd/surfaces.go` — existing surfaces command (L24), surfacesCmd registration
- `forge-cli/internal/cmd/init_surfaces.go` — askSurfaceConfirmation TUI infrastructure to reuse for --apply mode

## Acceptance Criteria
- [ ] `forge surfaces detect` runs detection + inference in read-only mode, shows results with source annotations, exits without writing config
- [ ] `forge surfaces detect --apply`: shows TUI confirmation (same flow as init), writes to config on confirm; exit code 0
- [ ] Non-interactive terminal: prints results to stdout, no TUI, no config write (regardless of `--apply`); exit code 0 if detection succeeds, 1 if no surfaces found
- [ ] Stdout format: one line per surface `<path>=<type> (<source>)`, where `<source>` is `detected:<signal>` or `inferred:<rule-id>`
- [ ] Empty detection: prints nothing to stdout, exits with code 1
- [ ] After `--apply` confirm, config file on disk contains the detected surfaces (verified by reading config and asserting entries match)
- [ ] `--project-root` flag supported (consistent with existing `forge surfaces` command)

## Hard Rules
- Non-interactive stdout format is unstable — document in help text that scripted consumers should pin forge version
- No config write without explicit `--apply` flag, even in interactive mode

## Implementation Notes
- Add as subcommand to existing `surfacesCmd` in `surfaces.go` — `surfacesCmd.AddCommand(detectCmd)`
- Reuse `DetectSurfacesWithConflicts` for detection + inference pipeline
- Reuse `askSurfaceConfirmation` TUI infrastructure for `--apply` mode (Task 2 modifies its signature to return Sources info)
- Interactive terminal detection: use `golang.org/x/term.IsTerminal` or equivalent (check if already imported in codebase)
- Register `--apply` as bool flag on the detect subcommand
- Stdout format matches source annotation pattern: `inferred:<rule-id>` or `detected:<signal>`
- Error handling: config file not found → print warning, run detection anyway; malformed config → print error, exit code 1
