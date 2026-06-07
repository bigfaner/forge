---
created: "2026-06-07"
tags: [testing, data-model]
---

# SQL schema tasks must use type "doc", not "coding.feature"

## Problem
Task 1.1 and 1.2 (SQL DDL changes: CREATE TABLE + ALTER TABLE) were classified as `coding.feature`. This triggered the forge CLI quality gate for `coding.*` tasks, which runs `go build + go fmt + go lint + go test`. Since `.sql` files are not Go source code, these checks are meaningless — they verify Go code that wasn't changed by the task.

## Root Cause
1. During `/breakdown-tasks`, SQL DDL tasks were typed as `coding.feature` because they produce "new runtime behavior" (new database tables)
2. The breakdown-tasks skill classifies by output artifact: "If output is only `.md` files, type must be `doc`". But `.sql` files are not `.md`, so the rule didn't catch this case
3. The `coding.feature` type triggered the Go quality gate (compile + test), which is irrelevant for `.sql`-only output
4. The task executor agent saw the Go quality gate and assumed it needed to also modify Go model code to make tests pass — crossing into task 1.3's scope

## Solution
**SQL DDL tasks (`.sql` output only) must use type `doc`** — not `coding.feature`.

The quality gate mapping:
| Type | Quality-gate | Correct for SQL? |
|------|-------------|-----------------|
| `coding.*` | Run (compile + fmt + lint + test) | ❌ — `.sql` not testable by Go |
| `doc` | Skip | ✅ — no Go compilation needed |

SQL DDL verification happens at migration time (running against actual DB), not at Go compile/test time.

## Reusable Pattern
When classifying task type during `/breakdown-tasks`:
1. Ask: "Does the output include any `.go`, `.ts`, `.tsx` files?" → If NO → type should be `doc`
2. The `.md` rule is a special case. The general rule is: **if no compilable runtime code is produced, use `doc`**
3. `.sql`, `.yaml`, `.json`, `.toml` are all non-compilable → `doc` type
4. Exception: if a task modifies both `.sql` AND `.go` files (e.g., schema + model in one task), use `coding.feature` because the Go quality gate IS relevant for the `.go` part

## Example
```
# Task 1.1 — CREATE TABLE (SQL only)
Type: doc (NOT coding.feature)
Files: MySql-schema.sql, SQLite-schema.sql
Quality gate: skip (SQL verified at migration time)

# Task 1.3 — Go model structs
Type: coding.feature
Files: model/milestone_map.go, model/milestone.go
Quality gate: go build + go test
```

## References
- breakdown-tasks skill → Type Assignment rules
- Task 1.2 execution record — ran 17 Go tests for a SQL-only change
