# Agent Copies Examples, Not Schema Docs

## Problem
Task execution records (records/*.md) had empty sections for key decisions, testing stats (0/0/0%), and acceptance criteria, even though the agent clearly ran tests and made design decisions.

## Root Cause
The `echo` command example in both `task-executor.md` (Step 4) and `record-task/SKILL.md` only demonstrated 3 fields:

```bash
echo '{"summary":"...","filesCreated":[...],"filesModified":[...]}' > process/record.json
```

Agents copy-pasted this example literally, omitting `keyDecisions`, `testsPassed`, `coverage`, and `acceptanceCriteria`. The full JSON schema documentation existed above the example but agents followed the runnable example, not the reference table.

Evidence: Record 2.3 summary says "20 tests with 95.5% coverage" but `testsPassed: 0`, `coverage: 0.0%` — the data was in prose, not in structured fields.

## Solution
1. Update the `echo` example in both `task-executor.md` and `record-task/SKILL.md` to include ALL fields with placeholder values
2. Consider making the CLI validate required fields and reject incomplete record.json

## Key Takeaway
**Examples > Documentation.** When writing agent instructions, the example command IS the contract. Agents will copy the example verbatim and ignore field reference tables. Always include every required field in examples, even if it makes them longer.
