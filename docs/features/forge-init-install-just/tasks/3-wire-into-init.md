---
id: "3"
title: "Wire ensureJust into forge init + add --skip-just flag"
priority: "P1"
estimated_time: "1h"
dependencies: ["2"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 3: Wire ensureJust into forge init + add --skip-just flag

## Description

Integrate the `ensureJust()` function from task 2 into the `forge init` command sequence. Add the `--skip-just` CLI flag to allow users to skip the step. Update the `initAction` status reporting to support the new INSTALLED status.

The ensureJust step should run early in the init sequence — before the justfile update step (step 4), since that step depends on `just` being available.

## Reference Files
- `docs/proposals/forge-init-install-just/proposal.md` — Source proposal
- `forge-cli/internal/cmd/init.go` — Init command (main modification target)
- `forge-cli/pkg/just/ensure.go` — ensureJust function (from task 2)

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `forge-cli/internal/cmd/init.go` | Add ensureJust step to init sequence; add `--skip-just` flag; update `initAction` status to support INSTALLED |
| `forge-cli/internal/cmd/init_test.go` | Add tests for --skip-just flag and ensureJust integration |
| `scripts/version.txt` | Bump minor version (new feature) |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] `forge init` runs ensureJust as a step before the justfile update step
- [ ] `forge init --skip-just` skips the ensureJust step entirely, reporting SKIPPED
- [ ] `initAction.status` supports INSTALLED in addition to existing values
- [ ] Init summary output includes the just installation result (INSTALLED/SKIPPED/FAILED)
- [ ] When just is already installed and >= 1.40.0, the step reports SKIPPED with version detail
- [ ] When just is installed successfully, the step reports INSTALLED with method detail (brew/cargo/scoop/choco/embedded)
- [ ] When installation fails, the step reports FAILED but init continues (non-blocking)
- [ ] `forge init --help` documents the `--skip-just` flag

## Hard Rules

- The ensureJust step MUST run before the justfile update step (step 4), since justfile recipes require `just` to exist
- Installation failure is non-blocking — init continues even if just installation fails (prints WARNING)
- Do NOT change the order of existing steps 1-5, only insert the new step

## Implementation Notes

- Insert the ensureJust step between step 3 (gitignore) and step 4 (justfile)
- The `--skip-just` flag is a boolean: `initCmd.Flags().Bool("skip-just", false, "skip just installation check")`
- Convert `EnsureResult` to `initAction` with a helper function
- Update `printInitSummary` to handle the new INSTALLED status (already works generically via status string)
- Consider adding the initAction status "INSTALLED" as a valid constant if other code validates status strings
