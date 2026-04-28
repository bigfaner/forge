---
name: graduate-tests
description: Migrate feature test scripts to the regression suite (tests/e2e/). Agent-driven: reads scripts, analyzes content, decides classification, splits/merges as needed, rewrites imports, creates graduation marker.
---

# /graduate-tests

Migrate feature test scripts from `tests/e2e/<feature>/` to the project-wide regression suite at `tests/e2e/<target>/`.

**Core principle**: Graduation is a decision, not a file copy. The agent reads and understands each spec before deciding where it belongs.

<HARD-GATE>
- Do NOT overwrite existing files in `tests/e2e/` without merging
- Do NOT graduate if marker already exists (idempotent)
- Do NOT modify the source scripts in `tests/e2e/<slug>/`
</HARD-GATE>

## Prerequisites

Check before running. Abort and prompt user if missing:

| Artifact | Missing prompt |
|----------|----------------|
| `tests/e2e/<slug>/` directory | Run `/gen-test-scripts` first |
| At least one `.spec.ts` file | Run `/gen-test-scripts` first |
| `tests/e2e/` graduation marker absent | Already graduated — skip |

```bash
task feature   # get current slug
ls tests/e2e/<slug>/
cat tests/e2e/.graduated/<slug> 2>/dev/null && echo "already graduated"
```

## When to Use

- After all tasks for a feature are completed and e2e tests pass
- User invokes `/graduate-tests` manually
- `/run-tasks` orchestrator suggests it post-completion

## Workflow

```
1. Check marker → 2. Read scripts → 3. Analyze structure → 4. Decide classification → 5. Migrate → 6. Write marker
```

### Step 1: Check Graduation Marker

```bash
cat tests/e2e/.graduated/<slug>
```

If marker exists: print "Already graduated on <timestamp>" and stop.

### Step 2: Read Source Scripts

Read all files in `tests/e2e/<slug>/`:

```bash
ls tests/e2e/<slug>/
```

Read each `.spec.ts` file and `helpers.ts`. Understand:
- What each `describe`/`test` block tests
- Which routes, APIs, or CLI commands are covered
- Whether a single spec mixes multiple functional domains

### Step 3: Analyze Existing Structure

If `tests/e2e/` exists, read its directory structure:

```bash
ls -R tests/e2e/
```

Understand the existing classification convention (by type, by route, by feature domain). New specs must follow the same convention.

### Step 4: Decide Classification

For each spec file, answer:

| Question | Decision |
|----------|----------|
| What functional domains does this spec cover? | One domain → keep as-is; multiple → split |
| Which `tests/e2e/<category>/` does it belong to? | Match existing convention or create new category |
| Does a spec file already exist at the target path? | Yes → merge tests, not overwrite |
| Does `tests/e2e/helpers.ts` already exist? | Yes → shared helpers already available, no action needed |

**Classification examples**:

```
# Input: tests/e2e/<slug>/
ui.spec.ts    # contains login, dashboard, profile tests
api.spec.ts   # all auth-related
cli.spec.ts   # general CLI commands

# Agent decides:
ui.spec.ts → split:
  tests/e2e/ui/login/login.spec.ts
  tests/e2e/ui/dashboard/dashboard.spec.ts
  tests/e2e/ui/profile/profile.spec.ts

api.spec.ts → tests/e2e/api/auth/auth.spec.ts

cli.spec.ts → tests/e2e/cli/cli.spec.ts  (no split — all general)
```

### Step 5: Execute Migration

For each target file:

1. **Create directory** if it doesn't exist
2. **Write spec file** (or merge into existing)

<HARD-RULE>
When moving spec files within `tests/e2e/`, the `'../helpers.js'` import path remains correct as long as the spec stays one level below `tests/e2e/`. No import rewriting is needed for single-level moves (e.g., `tests/e2e/<feature>/` → `tests/e2e/<target>/`). If a spec moves deeper (e.g., `tests/e2e/ui/login/`), adjust the import depth accordingly.
</HARD-RULE>

Shared infrastructure (`helpers.ts`, `package.json`, `tsconfig.json`) already exists at `tests/e2e/` — no merging or copying needed.

### Step 6: Create Graduation Marker

```bash
mkdir -p tests/e2e/.graduated
echo "$(date -u +%Y-%m-%dT%H:%M:%SZ)" > tests/e2e/.graduated/<slug>
```

## Output

Report to user:

```
Graduated <slug>:
  ui.spec.ts → tests/e2e/ui/login/login.spec.ts
  ui.spec.ts → tests/e2e/ui/dashboard/dashboard.spec.ts
  api.spec.ts → tests/e2e/api/auth/auth.spec.ts
  cli.spec.ts → tests/e2e/cli/cli.spec.ts

Marker: tests/e2e/.graduated/<slug>
```

## Related Skills

| Skill | Usage |
|-------|-------|
| `/gen-test-scripts` | Generate source scripts before graduation |
| `/run-e2e-tests` | Execute scripts before graduating |
| `/run-tasks` | Suggests graduation after all tasks complete |
