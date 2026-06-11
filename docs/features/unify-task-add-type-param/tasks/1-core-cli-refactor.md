---
id: "1"
title: "Core CLI refactor: remove --template, add --type auto-discovery"
priority: "P0"
estimated_time: "1-2h"
dependencies: []
surface-key: ""
surface-type: ""
breaking: true
type: "coding.refactor"
mainSession: false
---

# 1: Core CLI refactor: remove --template, add --type auto-discovery

## Description
Remove the `--template` flag from `forge task add` and unify template loading under `--type`. Rename embedded template files to match their type values (`fix-task.md` → `coding.fix.md`, `cleanup-task.md` → `coding.cleanup.md`). When `--type` is specified, the system auto-discovers the matching template file; if found, loads content + defaults. Update `quality_gate.go` to use type values instead of template names.

## Reference Files
- `proposal.md#Proposed-Solution` — defines the unified --type model and template filename-as-type-value strategy
- `proposal.md#Requirements-Analysis` — Key Scenarios for type-based template discovery (with/without matching template)
- `proposal.md#Impact-Analysis` — affected Go source files with specific line numbers
- `proposal.md#Key-Risks` — grep for hardcoded template name strings and test fixtures

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/template/data/fix-task.md` | Rename to `coding.fix.md` |
| `forge-cli/pkg/template/data/cleanup-task.md` | Rename to `coding.cleanup.md` |
| `forge-cli/pkg/template/template.go` | Update defaults map keys from `fix-task`/`cleanup-task` to `coding.fix`/`coding.cleanup` |
| `forge-cli/internal/cmd/task/add.go` | Remove `--template` flag; in `executeAdd()`, when `--type` is set, check if matching template exists via `tmpl.Get(addType)` and apply defaults |
| `forge-cli/internal/cmd/quality_gate.go` | Change `addFixTask()` to use type values (`coding.fix`/`coding.cleanup`) instead of template names (`fix-task`/`cleanup-task`); update `handleGateFailure()` help message |
| `forge-cli/internal/cmd/task/add_cmd_test.go` | Replace `--template fix-task` with `--type coding.fix` in test cases |

### Delete
| File | Reason |
|------|--------|
| `forge-cli/pkg/template/data/fix-task.md` | Renamed to `coding.fix.md` |
| `forge-cli/pkg/template/data/cleanup-task.md` | Renamed to `coding.cleanup.md` |

## Acceptance Criteria
- [ ] `forge task add --type coding.fix --title "Fix X"` loads `coding.fix.md` template and applies its defaults (priority=P0, breaking=true, estimated_time=30min, id_prefix=fix)
- [ ] `forge task add --type coding.feature --title "Build Y"` works without template (no matching file, no error)
- [ ] `forge task add --template fix-task` returns an error (flag removed)
- [ ] `forge task add -h` shows no `--template` flag
- [ ] Quality gate auto-created fix tasks use type value `coding.fix` instead of template name `fix-task`
- [ ] All existing tests pass after rename (`go test ./...`)

## Hard Rules
- Template files are `//go:embed` — rename requires rebuild, no runtime path changes needed
- `quality_gate.go` line ~478: `tmplName` variable assignment must change from `"fix-task"` to the type value
- The `--var` flag behavior must remain unchanged — template variable injection still works via opts.Vars
- `add.go` line ~160: `opts.Template` field is repurposed: set to `addType` only if template file exists for that type

## Implementation Notes
- In `add.go`, the template defaults block (lines 168-190) changes from `if addTemplate != ""` to `if addType != ""`: try `tmpl.Get(addType)`, if found apply defaults, else just set type field without template
- In `quality_gate.go`, the `addFixTask()` function uses `fixTypeFromStep(step)` to determine type, then sets `opts.Template = taskType` (which now matches the template filename). The `tmplName` indirection (lines 478-481) can be removed entirely since type values match filenames
- Search for any other references to `"fix-task"` or `"cleanup-task"` as string literals in Go code
