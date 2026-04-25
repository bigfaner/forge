---
name: breakdown-tasks
description: Use when design.md is finalized to break down into executable tasks. Creates task files based on technical design.
---

# Breakdown Tasks

Break a technical design into executable tasks (1-4h each, clear dependencies, testable acceptance criteria).

## Prerequisites

**Conditional Tags**: `<HAS_UI>`, `<NO_UI>`, `<UI_ONLY>`, `<RULE>` are inclusion markers. If `ui/ui-design.md` exists, include `<HAS_UI>` and `<UI_ONLY>` blocks and exclude `<NO_UI>` blocks. If not, include `<NO_UI>` blocks and exclude `<HAS_UI>`/`<UI_ONLY>`/`<RULE>` blocks.

| Artifact                | Missing? Run                    |
| ----------------------- | ------------------------------- |
| `prd/prd-spec.md`       | `/write-prd`                    |
| `design/tech-design.md` | `/tech-design` → `/eval-design` |

## Step 1: Read All Documents

Read `manifest.md` to locate documents, then read all available files:

- `prd/prd-spec.md` — WHAT to build
- `design/tech-design.md` — HOW to build it
- `design/api-handbook.md` — interfaces (if exists)
- `prd/prd-user-stories.md` — user scenarios with Given/When/Then AC (if exists)
- `prd/prd-ui-functions.md` — UI function requirements (if exists)
- `ui/ui-design.md` — UI component specs (if exists)

<HAS_UI>
If `ui/ui-design.md` exists, also list `ui/prototype/` files and read `ui/prototype/index.html` for page inventory (skip if no prototype directory).
</HAS_UI>

## Step 2: Map → Tasks

### Element Mapping

| Design Element       | Source         | Task Type                |
| -------------------- | -------------- | ------------------------ |
| Interface definition | tech-design.md | Interface task           |
| Data model           | tech-design.md | Model task               |
| Backend component    | tech-design.md | Implementation (Backend) |
| Error type           | tech-design.md | Error handling task      |

<UI_ONLY>
| UI Component (Layout + States + Interactions + Binding) | ui-design.md | Implementation (UI) |
</UI_ONLY>

<RULE>
- Each `ui-design.md` Component → **one** UI task (split only if >4h)
</RULE>

### PRD Coverage Verification

Read the **PRD Coverage Map** from `tech-design.md`. Every PRD acceptance criterion must map to at least one task. UI-facing requirements → UI tasks, not generic Implementation.

Fallback: if Coverage Map is incomplete, use `prd/prd-user-stories.md` acceptance criteria directly.

## Step 3: Task Order & Dependencies

<NO_UI>

```
1.x Interfaces → 2.x Models → 3.x Implementation → 4.x Integration → 5.x Tests
```

</NO_UI>

<HAS_UI>

```
1.x Interfaces → 2.x Models → 3.x Implementation (Backend) → 4.x Implementation (UI) → 5.x Integration → 6.x Tests
```

**Dependency rules:**

- Phase 3.x (Backend) → depends on 1.x + 2.x
- Phase 4.x (UI) → depends on 1.x interfaces + 2.x models (can mock backend; does NOT need 3.x)
- Phase 5.x (Integration) → depends on 3.x + 4.x
- Phase 6.x (Tests) → depends on 5.x
</HAS_UI>

## Step 4: Create Task Files

<HARD-RULE>
Read `templates/task.md` before writing any task file.
Naming: `<sequence>.<sub-sequence>-<slug>.md`
</HARD-RULE>

<HAS_UI>

### UI Task Reference Files

For each Phase 4.x task, **Reference Files** must include:

1. Matching `ui-design.md` Component section
2. Corresponding `ui/prototype/<page>.html` (or note "No HTML prototype available")
3. Data binding table for this component
4. Relevant `tech-design.md` interfaces

Example:

```
- ui/ui-design.md Component "Dashboard" — layout, states, interactions
- ui/prototype/dashboard.html — interactive prototype
- design/tech-design.md Interfaces — data contracts
```

</HAS_UI>

### User Stories

For each task, populate the **User Stories** section with matching stories from `prd/prd-user-stories.md`. Include full Given/When/Then acceptance criteria — this enables direct test generation. If no match, note "No direct user story mapping."

## Step 4b: Append Standard Test Tasks

After all business tasks are created, append two fixed test tasks to ensure e2e coverage.

<HARD-RULE>
Read the templates and create task files:
- `templates/gen-test-cases.md` — Generate e2e test cases
- `templates/gen-test-scripts.md` — Generate e2e test scripts

Replace `{{LAST_BUSINESS_TASK_ID}}` in gen-test-cases.md with the ID of the last business task.
Add both tasks to `index.json` in Step 5.
</HARD-RULE>

**T-test-1** (gen-test-cases.md): Calls `/gen-test-cases` skill, depends on last business task
**T-test-2** (gen-test-scripts.md): Calls `/gen-test-scripts` skill, depends on T-test-1

## Step 5: Create index.json

<HARD-RULE>Read `templates/index.json` before writing. Paths (`file`, `record`) are relative to `tasks/` directory. Populate `dependencies` per Step 3 rules.
</HARD-RULE>

Reference: [templates/index.json](templates/index.json) | Schema: [templates/index.schema.json](templates/index.schema.json)

## Step 6: Validate

```bash
task validate -file docs/features/<slug>/tasks/index.json
```

## Step 7: Update Manifest

- Fill traceability table (4-column: PRD Section | Design Section | UI Component | Tasks); use "—" for UI Component when no UI
- Advance status to `tasks`

## Output Checklist

- [ ] All task files created with `<id>-<slug>.md` naming
- [ ] `index.json` valid per schema, `task validate` passes
- [ ] Every PRD AC covered by ≥1 task
- [ ] Dependencies follow Step 3 rules (no cycles)
- [ ] UI tasks reference prototype files (if applicable)
- [ ] User Stories populated from `prd-user-stories.md`
- [ ] index.json 末尾包含 T-test-1 和 T-test-2，依赖关系正确
- [ ] `manifest.md` updated with traceability + `status: tasks`
