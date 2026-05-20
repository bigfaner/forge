# TUI Async Cmd Await Semantics

For TUI Steps involving asynchronous operations, declare `await` semantics in the Contract:

**Declaration format**:
```
- Input: key "d" await 3000ms
```

**Rules**:
- `await <N>ms` = wait for all pending Cmds to complete, up to N milliseconds
- Default timeout: `tui-await-timeout` from `.forge/config.yaml` capabilities, or 3000ms if not configured
- Timeout behavior: fail-fast, report the timed-out Cmd name
- `tea.Batch(cmd1, cmd2)`: all concurrent Cmds must complete before proceeding

**Outcome modeling for async TUI Steps**:

```
Outcome "diagnosis-loaded":
  Preconditions: "session loaded, call tree visible, entry expanded"
  Input: key "d" await 3000ms
  Output: "view contains diagnosis summary panel"
  State: "Model.diagnosis_panel field set to visible"

Outcome "diagnosis-timeout":
  Preconditions: "session loaded, call tree visible, entry expanded, async Cmd takes longer than 3000ms"
  Input: key "d" await 3000ms
  Output: "error message containing timed-out Cmd name, fail-fast"
  State: "unchanged from pre-Cmd state"
```

<HARD-RULE>
TUI async Steps MUST include a timeout Outcome when the Step has async Cmds. The timeout Outcome's Preconditions must include the timeout condition (async Cmd exceeds await duration). The timeout Outcome must report the timed-out Cmd name in its Output.
</HARD-RULE>

# State Verification Levels

When a project does not expose a state query interface, the State dimension degrades gracefully:

| Level | Declaration | When to use |
|-------|-------------|-------------|
| `full` | (default, no annotation needed) | Project exposes state query (CLI flag, API endpoint, file read) |
| `partial` | `<!-- state-verification: partial -->` | State fields can be inferred from Output only |
| `deferred` | `<!-- state-verification: deferred -->` + `limitations` section | Some state fields cannot be inferred |

gen-contracts determines the level automatically from Fact Table reconnaissance:
1. If the Fact Table contains state query interfaces (file paths, CLI flags, API endpoints) -> `full`
2. If state fields appear in output patterns -> `partial`
3. Otherwise -> `deferred`

# Journey-Level Invariants

Every Contract file MUST end with a `## Journey Invariants` section containing at least one cross-step invariant.

**Source**: Journey-level Invariants come from the Journey document's `## Journey Invariants` section. Copy them verbatim into every Contract file for that Journey.

**At least 1 invariant is mandatory**. If the Journey document has no declared Invariants, generate at least one from the workflow analysis (e.g., "feature_slug consistent across all steps" or "working directory unchanged between steps").

# Batch Processing

When a single Journey has more than 15 Contracts or the estimated token count exceeds 50k, automatically split into multiple batches:

- **Batch 1**: Happy path Outcomes (all steps' success Outcomes)
- **Batch 2+**: Edge case Outcomes grouped by semantic similarity

Split batches are merged back into complete Contract files (one per Step). The merged result must be structurally identical to a single-batch generation.

<HARD-RULE>
Batch splitting occurs within a Journey, not across Journeys. Each Journey is always processed completely before moving to the next. Merged Contract files must not lose or duplicate any Outcome or dimension.
</HARD-RULE>
