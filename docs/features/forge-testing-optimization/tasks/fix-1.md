---
id: "fix-1"
title: "Fix: init-justfile test specs reference moved file path"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
---

# Fix: init-justfile test specs reference moved file path

## Root Cause

49 test failures: all tests reference 'plugins/forge/commands/init-justfile.md' but the file was moved to 'plugins/forge/skills/init-justfile/SKILL.md'. The getInitJustfileContent() helper in 3 spec files needs to point to the new path. Also check if the file content structure changed (was commands/*.md, now skills/*/SKILL.md).

## Reference Files

- Source: tests/e2e/init-justfile/init-justfile.spec.ts,tests/e2e/justfile-e2e-integration/cli.spec.ts,tests/e2e/justfile-e2e-integration/detection-assembly.spec.ts
- Test script: tests/e2e/init-justfile/init-justfile.spec.ts,tests/e2e/justfile-e2e-integration/cli.spec.ts,tests/e2e/justfile-e2e-integration/detection-assembly.spec.ts
- Test results: tests/e2e/results/

## Verification

After fixing, verify the fix works:
1. `just test [scope]` — must pass
2. If UI/page related: `just test-e2e --feature <slug>` — must also pass

When this task is recorded as completed via `task record`, the source task T-quick-5 is automatically restored to pending if all its dependencies are completed.
