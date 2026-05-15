---
id: "2"
title: "Update task ID inference and prompt template for per-type tasks"
priority: "P0"
estimated_time: "1-2h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 2: Update task ID inference and prompt template for per-type tasks

## Description

Update the task type inference system (`InferType()`) and the prompt template to recognize and handle per-type task IDs. Currently, `InferType()` maps `T-test-2` and `T-quick-2` (with optional profile suffix like `a`, `b`) to `TypeTestPipelineGenScripts`. This needs to also recognize type-suffixed variants like `T-test-2-api`, `T-test-2-tui`, `T-quick-2-cli`, etc.

The prompt template (`test-pipeline-gen-scripts.md`) must pass the `--type` argument to the skill invocation when the task ID contains a type suffix.

## Reference Files
- `docs/proposals/test-scripts-per-type/proposal.md` — Source proposal
- `forge-cli/pkg/task/infer.go` — `InferType()` function to modify
- `forge-cli/pkg/prompt/prompt.go` — `typeToTemplate` mapping
- `forge-cli/pkg/prompt/data/test-pipeline-gen-scripts.md` — Prompt template to modify
- `forge-cli/internal/cmd/prompt_get.go` — Prompt rendering with runtime values

## Acceptance Criteria
- [ ] `InferType("T-test-2-api")` returns `TypeTestPipelineGenScripts`
- [ ] `InferType("T-test-2-tui")` returns `TypeTestPipelineGenScripts`
- [ ] `InferType("T-test-2-cli")` returns `TypeTestPipelineGenScripts`
- [ ] `InferType("T-quick-2-api")` returns `TypeTestPipelineGenScripts`
- [ ] Same for profile-suffixed + type-suffixed: `T-test-2a-api`, `T-quick-2b-tui`
- [ ] Existing patterns still work: `T-test-2`, `T-test-2a`, `T-quick-2`, `T-quick-2a`
- [ ] Prompt template passes `--type <capability>` to `Skill(skill="forge:gen-test-scripts")` when task ID contains a type suffix
- [ ] Prompt template omits `--type` when task ID has no type suffix (backward compatible)

## Hard Rules
- MUST NOT break existing `InferType()` behavior for non-type-suffixed IDs
- The type suffix pattern must be distinct from profile suffix (profile: single lowercase letter; type: hyphen + lowercase word like `-api`, `-tui`, `-cli`, `-web-ui`)

## Implementation Notes
- Current `profileSuffixedID()` allows a single trailing lowercase letter. Type suffix adds a hyphen + capability name after the profile suffix (or directly after the number if no profile suffix)
- ID pattern: `T-(test|quick)-2[a-z]?-(<capability>)` where capability matches profile capabilities (tui, api, cli, web-ui, etc.)
- The prompt template can extract the type from the task ID using a regex or by having the Go code pass it as a template variable
- Key risk: naming collision between profile suffix and type suffix — mitigated by hyphen separator (profile suffix has no hyphen)
