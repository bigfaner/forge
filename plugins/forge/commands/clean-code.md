---
name: clean-code
description: Simplify and clean up code. Optionally specify paths to scope cleanup.
allowed-tools: Bash Read Edit Write Glob Grep
---

Simplify and clean up code. Supports optional paths to scope cleanup to specific files or directories.

Without arguments, scope is determined automatically (git diff or feature context).

Invoke the skill:

```
Skill(skill="forge:clean-code")
```

With specific paths:

```
Skill(skill="forge:clean-code", args="pkg/service/handler.go forge-cli/internal/cmd/")
```
