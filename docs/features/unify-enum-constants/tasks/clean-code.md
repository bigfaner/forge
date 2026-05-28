---
id: "T-clean-code"
title: "Simplify and Clean Code"
priority: "P2"
estimated_time: "20min"
dependencies: ["5.1"]
type: "code-quality.simplify"
surface-key: ""
surface-type: ""
---

Simplify and clean up code for the unify-enum-constants feature.

## Discovery Strategy
1. Run `git diff --name-only main...HEAD` to identify files changed by this feature
2. Focus cleanup on changed files only
3. The skill resolves scope: git diff > feature context > user-specified paths

Do NOT clean up files outside this feature's scope.
