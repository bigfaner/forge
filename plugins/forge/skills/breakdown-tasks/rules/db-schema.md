# DB Schema Task Rules

**Load condition**: load this file IF `design/er-diagram.md` exists.

**Guard clause**: if `design/er-diagram.md` exists but contains no parseable entity definitions, skip this rule and proceed without schema tasks.

## Schema Task Creation

For each entity in `design/er-diagram.md`, create one Schema task.

### Input Artifacts

Each schema task references:
- `design/schema.sql` — DDL source
- `design/er-diagram.md` — entity relationship diagram

### Acceptance Criteria

Every schema task must include these acceptance criteria:

1. "DDL executes without error"
2. "All FK references resolve"
3. "Indexes created"

### Breaking Classification

Determine the `breaking` field based on the operation type:

- **ALTER existing table** -> `breaking: true`
- **All CREATE TABLE are new** -> `breaking: false`

Mixed scenarios (some ALTER, some CREATE): set `breaking: true` on tasks that ALTER existing tables, `breaking: false` on tasks that only CREATE new tables.

### Surface-Key/Type Assignment

All schema tasks typically resolve to `surface-type: "api"` or `"cli"` depending on the project. Use `forge surfaces --json <migration-file-path>` to determine the correct `surface-key` and `surface-type`. If no surfaces are configured, set both fields to empty strings.

### Dependency Rule

Schema tasks depend on interface tasks (if any), since the migration may need type information from the interface definitions. If no interface tasks exist, schema tasks have no intra-feature dependencies beyond phase structure.

### Type Assignment

Schema tasks that produce only `.sql` files (DDL: CREATE TABLE, ALTER TABLE, indexes) **must** use `type: "doc"`. SQL files are non-compilable and non-runnable — they should not trigger the language quality gate (compile + fmt + lint + test). Schema correctness is verified at migration time against an actual database, not at Go/TS compile time.

**Exception**: If a schema task also modifies compilable source files (e.g., a `.go` model struct in the same task), use `coding.feature` because the quality gate IS relevant for the compilable portion.

**Template selection**: SQL-only → `templates/task-doc.md`; mixed SQL + compilable code → `templates/task.md`.

## Maintenance Note

This rule file depends on the following sections in the skeleton SKILL.md:

- **Step 2: Map -> Tasks** — Element Mapping table (adds DB schema row)
- **Step 4a: Create Task Files** — task file creation, breaking classification, surface-key/type assignment

If either of these sections changes in the skeleton, verify that the schema task rules in this file remain consistent.
