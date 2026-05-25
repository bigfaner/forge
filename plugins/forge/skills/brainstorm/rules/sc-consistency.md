# SC Consistency Check

Success Criteria (SC) and In Scope entries must be internally consistent before a proposal is committed. This rule detects logical contradictions within SC entries (SC-to-SC) and between SC and In Scope entries (SC-to-InScope) through clustering and bidirectional satisfiability proof.

## Zero-Output Principle

If the SC set contains no contradictions, this check produces **zero output** — an empty report. Any non-empty output indicates a detected contradiction or ambiguity. There is no "informational" or "summary" output for clean sets.

## Protocol

### Phase 1: Clustering

Group all SC entries and In Scope entries by affected area. The affected area is the file path, directory, or module each entry touches.

**Clustering heuristic**: Extract file/directory/module references from each SC and InScope entry. Entries sharing the same root area belong to the same cluster. An entry referencing multiple areas is added to each relevant cluster.

```
Cluster assignment:
- For each SC/InScope entry, extract area references (file paths, directory names, module names)
- Group entries that share at least one area reference into the same cluster
- Entries with no discernible area reference form a "global" cluster
```

### Phase 2: Intra-Group Satisfiability Check

For each cluster, check every pair of entries within the cluster using **bidirectional proof**. Do NOT ask "are these contradictory?" directly — this biases the LLM toward agreement.

**Bidirectional proof structure for each pair (A, B)**:

```
Step A→B: Assume entry A is TRUE and fully satisfied.
  - Derive the resulting state of the affected area.
  - Check whether entry B can be satisfied in that derived state.
  - Record: CAN_SATISFY | CANNOT_SATISFY | AMBIGUOUS

Step B→A: Assume entry B is TRUE and fully satisfied.
  - Derive the resulting state of the affected area.
  - Check whether entry A can be satisfied in that derived state.
  - Record: CAN_SATISFY | CANNOT_SATISFY | AMBIGUOUS
```

**Contradiction判定**:

| A→B | B→A | Verdict |
|-----|-----|---------|
| CANNOT_SATISFY | CANNOT_SATISFY | **Mutual exclusion** — both cannot be true simultaneously |
| CANNOT_SATISFY | CAN_SATISFY | **Direction conflict** — A blocks B but not vice versa |
| CAN_SATISFY | CANNOT_SATISFY | **Direction conflict** — B blocks A but not vice versa |
| CANNOT_SATISFY | AMBIGUOUS | **Probable conflict** — leans toward contradiction |
| AMBIGUOUS | CANNOT_SATISFY | **Probable conflict** — leans toward contradiction |
| AMBIGUOUS | AMBIGUOUS | **Ambiguous** — requires user confirmation |
| CAN_SATISFY | CAN_SATISFY | No conflict — skip |

### Phase 3: Fallback Cross-Group Direction Check

After all intra-group checks complete, run a lightweight all-pair scan across clusters for **ADD vs SUBTRACT on the same symbol**. This catches cross-group contradictions that clustering may miss.

**Scope of fallback check**:
- Only detect direction-type conflicts: one entry ADDs a symbol (file, function, variable, behavior) while another entry SUBTRACTs the same symbol
- Non-directional cross-group contradictions (e.g., performance vs completeness trade-offs across modules) are NOT covered by this fallback — they are left to the eval layer's full-pair scan

```
Cross-group scan:
- For each entry across all clusters, classify its direction: ADD, SUBTRACT, or NEUTRAL
- For each ADD-SUBTRACT pair targeting the same symbol, flag as a direction conflict
- Symbols are identified by name (function name, file path, module identifier)
```

### Token Overflow Protection

For proposals with more than 25 SC entries, process clusters **serially** — check one cluster at a time rather than loading all pairs at once. Each cluster check is independent and results can be accumulated.

## Output Format

For each detected contradiction, output a structured conflict report:

```
CONFLICT #N:
  Pair: [Entry-A reference] <-> [Entry-B reference]
  Type: mutual_exclusion | direction_conflict | resource_competition
  Cluster: [cluster name or "cross-group"]
  Derivation:
    A->B: [brief derivation of B's state assuming A is true]
    B->A: [brief derivation of A's state assuming B is true]
  Suggestion: [delete one | declare mutual exclusion zone | rewrite as compatible | mark ambiguous]
```

### Entry Reference Format

Reference entries by their section and number:
- SC entries: `SC-N` (e.g., SC-3)
- InScope entries: `InScope-N` (e.g., InScope-2)

### Suggestion Guidelines

| Verdict | Suggestion Options |
|---------|--------------------|
| Mutual exclusion | Delete one entry, or declare the conflicting area as an explicit exclusion zone, or rewrite both to compatible formulations |
| Direction conflict | Restrict the blocking entry's scope, or reorder the entries to establish precedence, or rewrite to eliminate the directional conflict |
| Resource competition | Quantify the competition (e.g., performance budget), or split the resource between entries, or relax one entry's constraint |
| Ambiguous | Mark as "ambiguous — requires user confirmation". Do NOT force a binary choice. Present both interpretations and ask the user to clarify |

## Example: Pipeline-Integration-Stitch Contradiction

The `pipeline-integration-stitch` proposal contained this contradiction (scored 897/1000 in adversarial eval yet passed):

- **SC-3**: `grep -r "gen-and-run" forge-cli/` must return zero results (SUBTRACT: remove all references to gen-and-run)
- **InScope-2**: Preserve migration error prompts that reference gen-and-run (ADD: keep gen-and-run references in migration messages)

**Bidirectional proof**:

```
A->B (assume SC-3 true): grep returns zero results, meaning no file in forge-cli/
  contains "gen-and-run". Migration error prompts cannot reference "gen-and-run"
  without creating a grep match. InScope-2 CANNOT_SATISFY.

B->A (assume InScope-2 true): migration error prompts exist that reference
  "gen-and-run". Running grep -r "gen-and-run" forge-cli/ will match these
  prompts. SC-3 CANNOT_SATISFY.

Verdict: mutual_exclusion
```

**Output**:

```
CONFLICT #1:
  Pair: SC-3 <-> InScope-2
  Type: mutual_exclusion
  Cluster: forge-cli/
  Derivation:
    A->B: Zero grep results means no "gen-and-run" text exists, so migration prompts cannot reference it
    B->A: Preserved migration prompts contain "gen-and-run", so grep will return non-zero results
  Suggestion: Delete SC-3 (allow gen-and-run in migration prompts only), or change InScope-2 to "preserve migration guidance in documentation only (no code references)"
```

## Structured Output Field

To ensure this check is actually executed (hard protection against agent skipping the rule), include a `consistency_check_result` field in the SC section of the proposal:

```
consistency_check_result:
  status: pass | fail | ambiguous
  pairs_checked: <number>
  conflicts_found: <number>
```

If this field is missing from the proposal output, the SC section is considered incomplete and must be revised before proceeding.
