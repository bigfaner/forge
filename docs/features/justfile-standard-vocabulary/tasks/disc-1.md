---
id: "disc-1"
title: "Fix: cli.spec.ts test expectations after task 1.5 migration"
priority: "P1"
dependencies: []
status: 
breaking: true
---

# disc-1: Fix: cli.spec.ts test expectations after task 1.5 migration

After task 1.5 migrated just build -> just compile in task-executor.md, error-fixer.md, execute-task.md, the cli.spec.ts tests still expect 'just build && just test'. Also TC-005 references non-existent file fix-e2e.md. Tests TC-002, TC-005, TC-015, TC-016 fail. Need to update cli.spec.ts assertions to match the migrated commands.
