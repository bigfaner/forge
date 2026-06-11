---
id: "disc-1"
title: "Fix: 5 pre-existing flaky tests in feature_set_command"
priority: "P2"
dependencies: []
status: 
---

# disc-1: Fix: 5 pre-existing flaky tests in feature_set_command

5 tests in feature_set_command fail because forge feature fallback chain (state.json -> git worktree -> features-dir) is only partially implemented. TC-009/011/012 expect fallback when state.json is absent/corrupt/points to nonexistent dir. TC-014 expects solo-feature detection. TC-018 expects features-dir source. All return (none) instead of the expected fallback. Fix the fallback chain implementation or remove these tests.
