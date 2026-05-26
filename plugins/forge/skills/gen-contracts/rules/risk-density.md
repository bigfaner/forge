# Risk-Driven Outcome Density

This rule defines how Journey `risk_level` (set by gen-journeys) controls Outcome count per Step and total test count per Journey. gen-contracts reads the Journey's risk_level and applies the corresponding density target when generating Outcomes.

## Density Targets

| Risk Level | Outcomes per Step | Total Test Count (Journey) | Boundary/Exception Outcomes |
|------------|-------------------|----------------------------|----------------------------|
| High       | 3-5               | 13-20                      | Must include surface-required + inferred boundary Outcomes |
| Medium     | 2-3               | 8-12                       | Must include surface-required Outcomes; inferred Outcomes optional |
| Low        | 1-2               | 4-7                        | Surface-required Outcomes only; no inferred boundaries |

**Total test count** = sum of all Outcomes across all Step Contracts in a single Journey.

## Applying Density Rules

### Step 1: Read Journey Risk Level

Read the Journey document's frontmatter `risk_level` field (set by gen-journeys). This determines the density target:

- `risk_level: High` -- apply High density
- `risk_level: Medium` -- apply Medium density (default if field is missing)
- `risk_level: Low` -- apply Low density

### Step 2: Generate Outcomes per Step

For each Step in the Journey:

1. **Happy path Outcome** is always generated (counts as 1 Outcome)
2. **Surface-required Outcomes** are generated based on the project's surface type (read from `.forge/config.yaml`). Surface-required Outcomes come from the surface rule's "Required Outcome Reference" section (e.g., `surface-cli.md`, `surface-api.md`).
3. **Inferred boundary Outcomes** (only for High/Medium risk) are LLM-derived based on:
   - Project Fact Table (known error codes, edge cases from code reconnaissance)
   - Journey narrative (what could go wrong at this step)
   - Each inferred Outcome MUST be annotated with `source: inferred` and a reasoning explanation

### Step 3: Check Density Target

After generating Outcomes for all Steps in the Journey:

- Count total Outcomes across all Step Contracts
- If total is below the target range for the risk level, review Steps for missing boundary scenarios
- If total is above the target range, merge semantically similar Outcomes
- Per-Step count must also stay within range (3-5 for High, 2-3 for Medium, 1-2 for Low)

**Exception**: If a Step has more Outcomes than the per-Step target due to surface-required Outcomes, this is acceptable -- surface requirements override density targets. Document the override with a comment.

## Surface-Required Outcome Derivation

gen-contracts reads the surface type's required_outcomes from the corresponding surface rule file in the gen-journeys skill (resolve relative to the gen-journeys skill directory):

| Surface | Rule File | Required Outcomes |
|---------|-----------|-------------------|
| CLI     | `rules/surface-cli.md` | `not-found`, `already-exists` |
| API     | `rules/surface-api.md` | `unauthorized` (for authenticated endpoints) |
| TUI     | `rules/surface-tui.md` | `timeout` (for async Cmds) |
| Web   | `rules/surface-web.md` | `validation-error`, `session-expired` |
| Mobile  | `rules/surface-mobile.md` | (best-effort, no mandatory Outcomes) |

### Derivation Rules

1. Load the surface rule file for the project's detected surface type
2. For each Step, check if the Step's context matches a required Outcome:
   - CLI: every Step that accesses a resource -> add `not-found`; every Step that creates a resource -> add `already-exists`
   - API: every Step targeting an authenticated endpoint -> add `unauthorized`
   - TUI: every Step with async Cmd -> add `timeout` (per `rules/tui-async.md`)
   - Web: every Step with form submission -> add `validation-error`; every Step with session-dependent state -> add `session-expired`
3. Required Outcomes count toward the density target but do NOT count against it (i.e., they are minimum guarantees, not maximum limits)

## Inferred Boundary Outcome Annotation

LLM-derived boundary Outcomes (beyond surface-required ones) must include:

```markdown
## Outcome "resource-conflict"
<!-- source: inferred -->
<!-- reasoning: Fact Table shows forge.task.add checks for existing task ID (pkg/task/add.go:42), duplicate submission is a realistic boundary -->
- Preconditions: "task with same ID already exists in the index"
- Input: "forge task add with duplicate task ID"
- Output: "error message indicating task ID conflict"
- State: "task index unchanged, no new task file created"
```

**Annotation rules**:
- `source: inferred` is mandatory for all LLM-derived boundary Outcomes
- `reasoning` must cite Fact Table sources or explain the inference logic
- Surface-required Outcomes do NOT need `source: inferred` annotation (they come from rule files)

## Density Checkpoint

After generating all Outcomes for a Journey, output a density checkpoint:

```
Risk Density Checkpoint:
  Journey: <name>
  Risk Level: <High/Medium/Low>
  Target: <N-M> Outcomes per Step, <X-Y> total
  Actual: <counts per Step>, <total> total
  Status: ON_TARGET / ABOVE_TARGET / BELOW_TARGET
```

If BELOW_TARGET: review Steps for missed boundary scenarios before proceeding to validation.
If ABOVE_TARGET: merge similar Outcomes until within range.
