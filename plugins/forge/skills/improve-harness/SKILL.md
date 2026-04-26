---
name: improve-harness
description: Improve project harness based on evaluation report priorities. Implements P0/P1/P2 improvements from docs/harness-reports/.
---

# Improve Harness

Systematically improve project harness based on `/eval-harness` evaluation reports.

## When to Use

**Trigger:**
- After running `/eval-harness`
- User asks to "fix harness issues" or "implement P0 improvements"
- User provides `/improve-harness` command

**Skip:**
- No evaluation report exists (run `/eval-harness` first)
- All improvements already implemented

## Workflow

```
1. Read latest report → 2. Parse priorities → 3. Execute one by one → 4. Verify improvement
   docs/harness-reports/     P0→P1→P2       Confirm then execute    Test confirmation
```

### Step 1: Read Latest Report

```bash
latest=$(ls -t docs/harness-reports/*.md 2>/dev/null | head -1)
```

If not found, prompt user to run `/eval-harness` first.

### Step 2: Parse Priority Items

Extract the "Priority Improvements" table from the report:

| Priority | Tasks |
|----------|-------|
| P0 | Blocking improvements |
| P1 | High priority |
| P2 | Medium priority |

### Step 3: Execute Improvements

For each item (P0 → P1 → P2 order):

1. **Show task description**
2. **Ask confirmation**: `Execute <TASK_ID>? [Y/n/e(xplain)]`
3. **Implement the improvement**
4. **Verify with tests**

### Step 4: Verify

After each improvement, run project-specific verification:

| Language | Verification Command |
|----------|---------------------|
| Go | `go build ./... && go test ./...` |
| Node.js | `npm run build && npm test` |
| Python | `pytest` |
| Rust | `cargo build && cargo test` |
| Java | `mvn test` |

## Common Improvement Tasks

### P0 - Blocking

| ID | Task | Output |
|----|------|--------|
| P0.1 | Document freshness detection | `scripts/check-doc-freshness.sh` |
| P0.2 | Duplicate code detection | `scripts/check-duplicates.sh` |

### P1 - High Priority

| ID | Task | Output |
|----|------|--------|
| P1.1 | Knowledge base index | `docs/README.md` |
| P1.2 | Principle enforcement mapping | Update project rules |

### P2 - Medium Priority

| ID | Task | Output |
|----|------|--------|
| P2.1 | Architecture lint in CI | Update CI config |
| P2.2 | Lint error fix hints | `docs/LINT-FIXES.md` |

## Output

After completion, create improvement record:

**Path**: `docs/harness-reports/YYYY-MM-DD-improvements.md`

**Template**: See `templates/improvements.md`

## Related

- `/eval-harness` - Generate evaluation report
- `docs/HARNESS-EVALUATION.md` - Current evaluation summary
- `docs/harness-reports/` - Historical reports
