---
created: "2026-06-07"
tags: [testing, local-dev-deployment]
---

# Auto-generated review/doc tasks interleave with business tasks in claim order

## Problem
During `/run-tasks` dispatch, `T-review-doc` (auto-generated doc review) was claimed and completed between Phase 1 business tasks (1.4 → 1.3 → T-review-doc → 1.5). This was unexpected — the user expected all business tasks to complete before auto-generated pipeline tasks.

## Root Cause
1. `forge task index` auto-generates review, test, and validation tasks alongside business tasks
2. These auto-generated tasks have their own dependency chains but often have no explicit dependency on specific business tasks — only on phase gates
3. `forge task claim` picks the next "unblocked pending" task by priority/ID, which can interleave auto-generated tasks between business tasks
4. T-review-doc had no dependency on specific Phase 1 tasks (it only depended on the feature being set), so it was claimable as soon as earlier tasks in the queue were completed

## Solution
Fixed in `pipeline_validate.go`: added `ResolveAllBusinessTasks` to `T-test-gen-journeys`' DependsOn list. This ensures the test pipeline directly depends on **all** business tasks (not just the last one), so even when a middle task is blocked and a later independent task completes, the test pipeline stays blocked. Previously `ResolveIfGenerated` returned nil for ungenerated tasks, leaving zero dependencies; even the initial fix with `ResolveLastBusinessTask` was insufficient because it only depended on the highest-numbered task.

See also: [[gotcha-test-chain-not-linked-to-last-business-gate]] for the broader pattern and architectural recommendation.

## Reusable Pattern
When running `/run-tasks`:
- Auto-generated tasks (T-review-doc, clean-code, consolidate-specs, etc.) now depend on business tasks via `ResolveLastBusinessTask`
- The test pipeline chain starts only after all business tasks complete
- Check `forge task list --sort topo` before dispatching to verify the expected order

## Related Files
- `docs/features/milestone-map/tasks/index.json` — task dependency graph
