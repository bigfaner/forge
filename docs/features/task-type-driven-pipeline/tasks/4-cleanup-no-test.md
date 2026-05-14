---
id: "4"
title: "Remove --no-test flag and update skill documentation"
priority: "P1"
estimated_time: "1h"
dependencies: ["2"]
scope: "all"
breaking: true
type: "implementation"
mainSession: false
---

# 4: Remove --no-test flag and update skill documentation

## Description
Remove the `--no-test` CLI flag and `NoTest` field from `BuildIndexOpts` since the auto-detection in task 2 makes manual control obsolete. Update all skill documentation that references `--no-test`. Bump the version.

## Reference Files
- `docs/proposals/task-type-driven-pipeline/proposal.md` — Source proposal (D4)

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `forge-cli/internal/cmd/index.go` | Remove `--no-test` flag variable, flag registration, and `NoTest` field in `BuildIndexOpts`; remove `NoTest` from summary output check |
| `forge-cli/pkg/task/build.go` | Remove `NoTest` field from `BuildIndexOpts` struct |
| `plugins/forge/skills/quick-tasks/SKILL.md` | Remove `--no-test` flag documentation and all references to it in Steps, Output Checklist |
| `plugins/forge/commands/quick.md` | Remove `--no-test` flag section and passthrough to quick-tasks |
| `forge-cli/scripts/version.txt` | Bump version (minor: new auto-detection feature) |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] `forge task index --no-test --feature <slug>` returns `unknown flag: --no-test` error
- [ ] `BuildIndexOpts` struct has no `NoTest` field
- [ ] `index.go` has no reference to `--no-test` or `indexNoTest`
- [ ] `quick-tasks/SKILL.md` has no mention of `--no-test`
- [ ] `quick.md` command has no `--no-test` flag section or passthrough
- [ ] `build.go` has no reference to `opts.NoTest`
- [ ] Version bumped in `forge-cli/scripts/version.txt`
- [ ] Existing tests in `index_test.go` pass after removing the flag
- [ ] Running `forge task index --feature task-type-driven-pipeline` (no --no-test) works correctly

## Implementation Notes
- This is a breaking change for any script/CI that passes `--no-test`. The auto-detection from task 2 replaces this entirely.
- The `build.go` reference to `opts.NoTest` was already replaced in task 2 with `isDocsOnlyFeature()`, so this task only needs to remove the struct field and the CLI flag wiring.
- For `quick-tasks/SKILL.md`, also update the "Flags" section and any conditional instructions that reference `--no-test`.
- For `quick.md`, update the "Flags" section and Step 3 (remove the conditional `args="--no-test"` passthrough).
- Version bump: this is a minor version bump (new feature: auto-detection of docs-only features).
