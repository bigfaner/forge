---
id: "2"
title: "Implement config-driven task generation with mode scoping"
priority: "P0"
estimated_time: "2h"
dependencies: ["1"]
scope: "backend"
breaking: true
type: "enhancement"
mainSession: false
---

# 2: Implement config-driven task generation with mode scoping

## Description

Modify `forge task index` to read the `auto` config block and conditionally generate tasks based on the current mode (quick/full). Rename T-test-5 → T-specs-1 and T-quick-5 → T-quick-specs-1. Add T-clean-code-1 task type for auto code cleanup.

## Reference Files
- `docs/proposals/auto-behavior-config/proposal.md` — Source proposal

## Acceptance Criteria
- [ ] `auto.e2eTest.quick=false` → quick mode generates zero T-quick test tasks
- [ ] `auto.e2eTest.full=false` → full mode generates zero T-test tasks
- [ ] `auto.consolidateSpecs.quick=false` → no T-quick-specs-1 generated
- [ ] `auto.consolidateSpecs.full=false` → no T-specs-1 generated
- [ ] `auto.cleanCode.quick=true` → T-clean-code-1 generated between last business task and first quick test task
- [ ] `auto.cleanCode.full=true` → T-clean-code-1 generated between last business task and first full test task
- [ ] Projects without `auto` config behave identically to before (backward compat)
- [ ] T-test-5 renamed to T-specs-1 in all Go code
- [ ] T-quick-5 renamed to T-quick-specs-1 in all Go code
- [ ] New type constant + prompt template for T-clean-code-1 (calls /simplify)
- [ ] T-clean-code-1 depends on last business task; first test task depends on T-clean-code-1 (when both exist)

## Hard Rules
- **Backward compat is critical**: when `auto` block is missing from config, ALL defaults must produce the exact same behavior as before the change. This means e2eTest defaults to true for both modes, consolidateSpecs defaults to true, cleanCode defaults to false.
- The Go config struct must use `ModeToggle{Quick: true, Full: true}` as zero-value default, NOT `false`. Use pointer types or explicit default-filling if needed.
- Keep existing capability gating: if capabilities are empty, return nil (no test/maintenance tasks).
- Existing e2e tests (`tests/e2e/`) that create temporary projects WITHOUT `.forge/config.yaml` must continue to pass — they rely on default behavior.

## Implementation Notes
- **Config struct** (`forge-cli/pkg/profile/config.go`): add `AutoConfig` struct with `ModeToggle` sub-structs
- **Config reading** (`forge-cli/internal/cmd/index.go`): read `auto` block, resolve mode (quick/full from existing mode detection), extract per-mode bools
- **Task generation** (`forge-cli/pkg/task/testgen.go`):
  - `GetBreakdownTestTasks()` — gate on `auto.e2eTest.full` for T-test-1~4.5, gate on `auto.consolidateSpecs.full` for T-specs-1 (renamed from T-test-5)
  - `GetQuickTestTasks()` — gate on `auto.e2eTest.quick` for T-quick-1~4, gate on `auto.consolidateSpecs.quick` for T-quick-specs-1 (renamed from T-quick-5)
  - Add `T-clean-code-1` generation gated on `auto.cleanCode.{mode}`
- **Types** (`forge-cli/pkg/task/types.go`): add `TypeCleanCode = "code-quality.simplify"`
- **Inference** (`forge-cli/pkg/task/infer.go`): add T-specs-1, T-quick-specs-1, T-clean-code-1 patterns
- **Prompt** (`forge-cli/pkg/prompt/data/`): add `code-quality-simplify.md` template; update consolidate/drift templates for renamed IDs
- **Breaking change**: existing index.json with T-test-5/T-quick-5 IDs incompatible. Document in CHANGELOG.
