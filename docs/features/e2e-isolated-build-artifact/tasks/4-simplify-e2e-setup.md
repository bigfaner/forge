---
id: "4"
title: "Simplify e2e-setup in justfile to optional cache optimization"
priority: "P2"
estimated_time: "30min"
dependencies: ["1", "2", "3"]
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 4: Simplify e2e-setup in justfile to optional cache optimization

## Description

With all E2E test modules now auto-building via TestMain, the `e2e-setup` recipe's build step is no longer a prerequisite for running tests. Update the justfile to reflect this: mark the build step as optional cache optimization, and update comments/documentation to clarify that developers can run E2E tests directly without `e2e-setup`.

## Reference Files
- `docs/proposals/e2e-isolated-build-artifact/proposal.md` — Source proposal
- `justfile` — Current e2e-setup recipe

## Acceptance Criteria
- `e2e-setup` recipe's build step commented as optional/cache optimization
- Running E2E tests without `e2e-setup` first works correctly (tests auto-build)
- No recipe that is a prerequisite for E2E tests requires `e2e-setup` build step

## Hard Rules
- Do not remove `e2e-setup` entirely — it's still useful as a cache warm-up
- Do not change the `e2e-setup` recipe's name or external interface

## Implementation Notes
- The build step in e2e-setup can be kept but should be clearly documented as optional
- Consider adding a comment like "# Optional: pre-builds forge binary for faster test startup"
