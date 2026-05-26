---
id: "T-clean-code"
title: "Simplify and Clean Code"
priority: "P2"
estimated_time: "20min"
dependencies: []
type: "code-quality.simplify"
scope: "all"
---

Simplify and clean up code for the surface-aware-justfile feature.

## Discovery Strategy
1. Run `git diff --name-only main...HEAD` to identify files changed by this feature
2. Focus cleanup on changed files only
3. The skill resolves scope: git diff > feature context > user-specified paths

Do NOT clean up files outside this feature's scope.
