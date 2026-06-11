---
created: 2026-05-15
author: "fanhuifeng"
status: Draft
---

# Proposal: Split gen-test-scripts by Test Type

## Problem

The `gen-test-scripts` task generates all test types (UI, API, CLI) in a single task execution. This creates two issues:

1. **Task too large** — A single task generates all spec files, shared helpers, and infrastructure, making it slow and hard to review
2. **No independent retry** — If one test type fails generation, the entire task is blocked. You can't retry just the API scripts without regenerating UI scripts too

### Evidence

- gen-test-scripts produces UI specs + API specs + CLI specs + helpers + auth-setup + config in one pass
- When API script generation fails (e.g., missing endpoint), the UI scripts that already succeeded must be regenerated
- Task review is difficult — reviewing 3+ spec files + infrastructure in a single task context

### Urgency

As projects grow with more test types and profiles, the single-task bottleneck worsens. Splitting now prevents accumulating technical debt in the test pipeline.

## Proposed Solution

Split `gen-test-scripts` into separate tasks by test type. Each detected interface type (UI, API, CLI) gets its own task:

- `T-test-2-ui`: generates UI test scripts only
- `T-test-2-api`: generates API test scripts only
- `T-test-2-cli`: generates CLI test scripts only

Only types detected by the profile capabilities and present in the project are created. `T-test-3` (run-e2e-tests) depends on ALL gen tasks completing, then runs all scripts together.

Same split applies to quick mode (`T-quick-2-ui`, `T-quick-2-api`, `T-quick-2-cli`).

### Innovation Highlights

Straightforward decomposition of a monolithic task into type-scoped tasks. The key insight is using the existing profile capability detection to determine which per-type tasks to create, ensuring no empty tasks are generated for unsupported types.

## Requirements Analysis

### Key Scenarios

- **Single type project**: Only API endpoints, no UI. Only `T-test-2-api` is created. `T-test-3` runs after it.
- **Multi-type project**: Web app with UI + API. `T-test-2-ui` and `T-test-2-api` created in parallel. `T-test-3` waits for both.
- **Type generation failure**: API scripts fail. UI scripts already generated. Only `T-test-2-api` is retried.
- **Multi-profile project**: Two profiles active. Per-profile tasks get per-type sub-tasks: `T-test-2a-ui`, `T-test-2a-api`, `T-test-2b-ui`, `T-test-2b-api`.

### Non-Functional Requirements

- Backward compatible: existing test tasks continue to work for projects that don't opt in
- Shared infrastructure (helpers, config, auth-setup) generated once by the first gen task, reused by subsequent ones
- No increase in total generation time — parallel execution offsets task overhead

### Constraints & Dependencies

- Depends on existing profile capability detection mechanism
- Must work with both single-profile and multi-profile configurations
- Shared infrastructure files must be idempotent (safe to generate multiple times)

## Alternatives & Industry Benchmarking

### Industry Solutions

Standard practice in CI/CD pipelines is to split test generation and execution by type (unit, integration, e2e, API, UI) into separate jobs/stages. This is the norm in GitHub Actions, Jenkins, and GitLab CI.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | No change cost | Task remains monolithic, no independent retry | Rejected: doesn't solve stated problems |
| Split gen only (recommended) | CI/CD standard | Independent retry, parallel gen, focused review | More tasks in index, shared infra needs coordination | **Selected: best cost/benefit ratio** |
| Split gen + run | CI/CD standard | Maximum parallelism | More complex dependency graph, run task needs type awareness | Rejected: diminishing returns, adds complexity |

## Feasibility Assessment

### Technical Feasibility

All required mechanisms exist: profile capabilities determine types, task index supports custom types, gen-test-scripts can accept a type filter. No new infrastructure needed.

### Resource & Timeline

Small scope — modifying 4 existing skills (gen-test-scripts, breakdown-tasks, quick-tasks, task executor). Estimated 4-6 tasks.

### Dependency Readiness

No external dependencies. All changes are within the forge plugin.

## Scope

### In Scope

- Add `--type` argument to gen-test-scripts skill to filter generation by type (ui/api/cli)
- Update breakdown-tasks to create per-type gen-scripts tasks (T-test-2-ui, T-test-2-api, T-test-2-cli)
- Update quick-tasks to create per-type gen-scripts tasks (T-quick-2-ui, T-quick-2-api, T-quick-2-cli)
- Update task executor to map new task types to gen-test-scripts with type filter
- Update T-test-3 / T-quick-3 dependencies to depend on ALL per-type gen tasks
- Handle shared infrastructure generation (first gen task creates, subsequent reuse)

### Out of Scope

- Splitting run-e2e-tests (T-test-3) by type
- Splitting graduate-tests (T-test-4) by type
- Changes to gen-test-cases (T-test-1) or eval-test-cases (T-test-1b)
- Changes to the task CLI core (forge binary)
- Multi-profile naming scheme changes (keep existing profile suffix convention)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Shared infra race condition (two gen tasks write helpers simultaneously) | M | H | First-gen-wins with idempotent writes; gen-test-scripts already has "create only if missing" logic |
| Naming collision between profile suffix and type suffix | L | M | Use hyphen-separated type suffix (-ui, -api, -cli) distinct from profile lowercase letters (a, b, c) |
| Tasks generated for types with no test cases | L | L | Only create per-type task when test-cases.md contains cases of that type |

## Success Criteria

- [ ] gen-test-scripts accepts `--type <ui|api|cli>` and generates only scripts for that type
- [ ] breakdown-tasks creates separate tasks per detected test type instead of one T-test-2
- [ ] quick-tasks creates separate tasks per detected test type instead of one T-quick-2
- [ ] T-test-3 depends on all per-type T-test-2-* tasks completing successfully
- [ ] Failed gen task can be independently retried without affecting other type tasks
- [ ] Shared infrastructure files are generated idempotently (no conflicts from parallel gen)

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks from this proposal
