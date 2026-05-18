# Testing Decisions

| Date | Feature | Decision | Rationale | Source |
|------|---------|----------|-----------|--------|
| 2026-05-18 | contract-journey-test-model | Semantic descriptors → regex via Fact Table | Natural-language descriptors (e.g. "stdout contains \"claimed task\"") are converted to regex patterns using Fact Table entries as ground truth, avoiding brittle hardcoded patterns. When no fact matches, descriptor is marked unresolved rather than producing a false-positive regex. | task 5, `pkg/descriptor/` |
| 2026-05-18 | contract-journey-test-model | Tag-based promotion replaces graduate-tests | Replaced the separate `graduate-tests` skill/workflow with `forge test promote --tag` command — simpler, composable, no special skill needed. | task 8 |
