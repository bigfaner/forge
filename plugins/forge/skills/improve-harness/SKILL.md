---
name: improve-harness
description: Dynamically implement harness improvements from eval-harness report. Reads P0/P1/P2 priorities and fixes each finding.
---

# Improve Harness

Read the latest `/eval-harness` report and implement improvements for each finding.

## Prerequisites

| Artifact | Missing prompt |
|----------|----------------|
| `docs/harness-reports/YYYY-MM-DD.md` (eval report) | Run `/eval-harness` first |

## When to Use

**Trigger:**
- After running `/eval-harness`
- User asks to "fix harness issues" or "implement improvements"
- User provides `/improve-harness` command

**Skip:**
- No evaluation report exists (run `/eval-harness` first)
- All findings already resolved

## Document Context

| Document | Path | Purpose |
|----------|------|---------|
| Rubric (scoring criteria) | `plugins/forge/skills/eval-harness/templates/rubric.md` | Defines the 4 dimensions and 12 criteria that findings are scored against |
| Eval report | `docs/harness-reports/YYYY-MM-DD.md` | Scored report with Priority Improvements table |
| Snapshot (raw evidence) | `docs/harness-reports/YYYY-MM-DD-snapshot.md` | Original context the scorer evaluated; useful when investigating findings |
| P0/P1/P2 classification | Defined in `eval-harness/SKILL.md` Step 4 | P0 = score 0 on any criterion; P1 = < 50%; P2 = < 80% |

## Workflow

```
1. Read report → 2. Extract findings → 3. Init record → 4. Fix one by one (append record after each) → 5. Finalize
```

## Step 1: Read Latest Report

```bash
latest=$(ls -t docs/harness-reports/????-??-??.md 2>/dev/null | grep -v -e snapshot -e improvements | head -1)
```

If not found, prompt user to run `/eval-harness` first.

Read the report and extract all rows from the **Priority Improvements** table.

## Step 2: Extract Findings

From the report, parse each priority item into:

| Field | Source |
|-------|--------|
| Priority | P0 / P1 / P2 |
| Dimension | Which rubric dimension |
| Criterion | Which specific criterion |
| Finding | What's wrong |
| Suggested Fix | What the report recommends |

Sort by priority (P0 first, then P1, then P2).

### 2.5: Skip Already-Completed Findings

If `docs/harness-reports/YYYY-MM-DD-improvements.md` already exists (from a previous interrupted run), read it and extract all completed findings from the "Completed" table. Match findings by the triple `(Priority, Dimension, Criterion)` — these three fields uniquely identify each finding. Do not match by Finding text, as wording may differ. Remove matched findings from the list.

Report to user:
```
Found existing improvement record with N completed fixes.
Remaining: X P0, Y P1, Z P2 findings to address.
```

## Step 3: Initialize Improvement Record

Create the improvement record file with header only, before starting any fixes. This ensures partial progress survives interruption.

**Path:** `docs/harness-reports/YYYY-MM-DD-improvements.md`

**Template:** See `templates/improvements.md`

Fill in the header (date, baseline score, report link) and leave tables empty.

## Step 4: Fix Findings

For each finding in priority order:

### 4.1 Present to User

```
## [P0] Architectural Boundaries > Boundaries mechanically enforced

**Finding:** No linter or CI check enforces dependency direction.
**Fix plan:** Create a lint script that validates import/dependency rules.

Execute? [Y/n/e(xplain)/s(kip)]
```

- **Y (default)**: Proceed with fix
- **n**: Skip this finding, add to Skipped table with reason "user declined"
- **e(xplain)**: Explain the fix plan in more detail, then re-ask the same question
- **s(kip)**: Skip permanently, add to Skipped table with reason "user skipped"

### 4.2 Implement Fix

Based on the finding, dynamically design and implement the fix. Common fix patterns:

| Finding Pattern | Fix Pattern |
|----------------|-------------|
| No doc validation | Create freshness detection script or CI check |
| No boundary enforcement | Create architecture lint script |
| No shared tools | Extract reusable skill/agent/script |
| No execution records | Set up task record schema and templates |
| Errors lack remediation | Add fix hints to linter/test output |
| No docs index | Create docs/README.md with catalog |
| Ad-hoc patterns | Centralize into shared utility |

### 4.3 Verify Fix

After implementing, verify the fix addresses the specific finding:

1. **For scripts**: Run the script, confirm it detects violations
2. **For docs**: Check file exists, content is correct
3. **For config**: Parse config, confirm setting is applied
4. **For project tests**: Apply **Scope Resolution** (see Forge Guide), then run `just test` to ensure nothing broke

### 4.4 Append to Record Immediately

After each verified fix (or skip), append a row to the improvement record's "Completed" or "Skipped" table. Do NOT wait until all fixes are done — this ensures partial progress survives interruption.

## Step 5: Finalize Record

After all findings are processed (or user stops), update the improvement record:

- Fill in the **Verification** section with all check results
- Fill in the **Files Changed** section
- Fill in the **Follow-up** section

Report to user:
```
Improvements complete. N fixes applied, M skipped.
Record: docs/harness-reports/YYYY-MM-DD-improvements.md
Re-run `/eval-harness` to verify score improvement.
```

## Guidelines

- **Fix the finding, not the symptom.** If the report says "no boundary enforcement," don't just create a file — create a script that actually detects violations.
- **Prefer mechanical enforcement over documentation.** A linter beats a comment. A CI check beats a README note.
- **Keep fixes project-appropriate.** Don't create Go scripts for a Python project. Use the project's language and tooling.
- **One finding = one atomic fix.** Don't bundle multiple fixes into one step. Each should be independently verifiable.

## Related

- `/eval-harness` - Generate evaluation report
- `docs/harness-reports/` - Reports and improvement records
