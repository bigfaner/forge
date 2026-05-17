---
name: gen-contracts
description: Generate Contract specifications (six dimensions + semantic descriptors + multi-Outcome + Invariants) from Journey documents and code reconnaissance (Fact Table). Formalizes TUI async Cmd await semantics.
conventions:
  - testing-isolation.md
---

# Gen Contracts

Generate Contract specifications from Journey documents, enriched with code reconnaissance (Fact Table).

**Core principle**: Every Step gets a Contract with six dimensions. gen-contracts uses *semantic descriptors* (natural language) -- precise regex is deferred to gen-test-scripts. This skill is the most technically complex in the pipeline because it bridges narrative Journeys with verifiable Contracts using code reconnaissance.

<HARD-GATE>
This skill ONLY writes to `tests/<journey>/_contracts/`. It does NOT generate test scripts or execute tests (handled by downstream skills).

**FORBIDDEN output paths**: Any path outside `tests/<journey>/_contracts/`. The `_contracts/` directory is the sole output target.
</HARD-GATE>

## Pipeline Position

```
gen-journeys -> gen-contracts -> gen-test-scripts -> run-tests
     (narrative)     (technical)     (code gen)        (execution)
```

gen-contracts is the second step. It reads Journey documents and source code, producing Contract specifications that gen-test-scripts consumes.

| Step | Input | Output | Reads Code |
|------|-------|--------|------------|
| gen-contracts | Journey documents + code reconnaissance | Contract specifications (six dimensions, semantic descriptors) | Yes (Fact Table) |

## Prerequisites

Check previous stage artifacts. Abort and prompt user if missing:

| Artifact | Missing prompt |
|----------|----------------|
| `docs/features/<slug>/testing/journeys/manifest.md` | Run `/gen-journeys` first |
| At least one Journey file in `docs/features/<slug>/testing/journeys/` | Run `/gen-jneys` first |

`<slug>` from `forge feature`.

## Step 0: Resolve Language and Interfaces

1. **Detect language**: Run `forge test detect` to auto-detect the project's test language(s) from file signals.
2. **On failure** (no language detected): ask the user to add `languages` to `.forge/config.yaml`.
3. **Detect interfaces**: Run `forge test interfaces` to discover which interface types (cli, api, tui, web-ui, mobile) the project exposes.

<HARD-RULE>
Do NOT silently default to any language or interface. If detection fails and the user cannot configure, abort the skill.
</HARD-RULE>

## Convention Loading

After language/interface resolution and before entering the workflow steps, load project conventions into context.

**Resolution algorithm**:

1. **Project-wide conventions**: Read this skill's own frontmatter `conventions` field. For each filename, check `docs/conventions/{filename}` -- if it exists, read it into context; if missing, skip silently.
2. **Interface-specific conventions**: For each detected interface, check `docs/conventions/testing-{interface}.md` -- if it exists, read it into context; if missing, skip silently.

<HARD-RULE>
Convention loading is non-blocking. Missing convention files are silently skipped.
</HARD-RULE>

## Workflow

```
0. Resolve language + interfaces -> 1. Read Journeys -> 2. Code Reconnaissance (Fact Table) -> 3. Generate Contracts -> 4. Validate -> 5. Write Output
```

### Step 1: Read Journey Documents

1. Read `docs/features/<slug>/testing/journeys/manifest.md` to discover all Journey files.
2. For each Journey listed in the manifest, read the Journey document.
3. Parse each Journey's structure:
   - Journey name and Risk level (from frontmatter)
   - Happy path steps (sequence number, user action, expected result)
   - Edge cases (referenced step, precondition, user action, expected result)
   - Journey Invariants (cross-step properties)

<HARD-RULE>
Every Journey in the manifest MUST be processed. Do not skip Journeys based on Risk level or step count.
</HARD-RULE>

### Step 2: Code Reconnaissance (Build Fact Table)

Read source code to extract ground-truth values for enriching Contracts with real context.

**This step is REQUIRED** -- gen-contracts needs code context to produce accurate State dimensions, Input schemas, and Side-effect declarations.

**Generic reconnaissance reads**:

| Source | What to extract |
|--------|-----------------|
| CLI entry points | Command names, flag names, flag types, output patterns |
| API handlers | Request/response schemas, status codes, middleware |
| TUI model files | Model struct fields, Cmd definitions, Msg types, View rendering |
| Config files | Port numbers, base paths, timeout values, auth mechanisms |
| State storage | File paths, JSON schemas, database tables |
| Hook definitions | Hook names, trigger conditions, parameter schemas |

**TUI-specific reconnaissance** (when `tui` interface detected):

| Source | What to extract |
|--------|-----------------|
| Cmd definitions | Cmd function names, async behavior (do they return Msg?) |
| Batch usage | `tea.Batch()` calls and their Cmd arguments |
| Timeout configurations | Any timeout constants, default wait durations |
| Model transitions | Init -> Idle -> Processing -> Result states |

Build Fact Table with source citations:

```markdown
## Fact Table
| Key | Value | Source |
|-----|-------|--------|
| CLI_COMMAND_FEATURE | forge feature | cmd/feature.go:15 |
| TUI_AWAIT_TIMEOUT | 3000ms | internal/tui/config.go:8 |
```

<HARD-RULE>
- Every Fact Table value must cite source file and line number. Unknown sources -> `UNKNOWN`. Do not fabricate.
- Fact Table values inform State dimension and Input dimension declarations. Use them to ground semantic descriptors in real code.
- When the project does not expose a state query interface, set `state-verification: partial` or `state-verification: deferred` in the Contract.
</HARD-RULE>

### Step 3: Generate Contracts

For each Journey, generate one Contract file per Step.

#### 3.1 Step-to-Contract Mapping

Each Journey Step (happy path or edge case) maps to an Outcome within the Step's Contract:

- **Happy path steps** become the `"success"` Outcome (or similar positive label) in their respective Step Contract
- **Edge cases** become additional Outcomes in the corresponding Step's Contract, with Preconditions that differ from the happy path

**Outcome grouping rule**: All Outcomes for a single Step are grouped into one Contract file. The Step number determines the file: `step-<N>-<action-slug>.md`.

**Example mapping**:

```
Journey "task-lifecycle":
  Happy Path:
    Step 1: forge feature my-feature   -> Contract step-1-feature-create.md, Outcome "success"
    Step 2: forge task claim           -> Contract step-2-task-claim.md, Outcome "success"
    Step 3: forge task submit          -> Contract step-3-task-submit.md, Outcome "success"
  Edge Cases:
    Step 2b: claim with no tasks       -> Contract step-2-task-claim.md, Outcome "no-tasks-available"
    Step 3b: submit without claim      -> Contract step-3-task-submit.md, Outcome "not-in-progress"
```

#### 3.2 Six-Dimension Declaration

For each Outcome, declare all six dimensions. Four are mandatory (non-empty), two are optional:

**Mandatory** (every Outcome MUST have non-empty values):

| Dimension | Content | Source |
|-----------|---------|--------|
| Preconditions | State that must hold before execution | Journey edge case preconditions + Fact Table state info |
| Input | What goes into the system | Journey user action + Fact Table command/flag/endpoint info |
| Output | What the system produces | Journey expected result + Fact Table output patterns |
| State | How system state changes | Fact Table state storage info + output inference |

**Optional** (may be omitted; omission = no constraint):

| Dimension | Content | When to include |
|-----------|---------|-----------------|
| Side-effect | External effects (hooks, network calls, async Cmds) | When Fact Table reveals side effects |
| Invariants (step-level) | Properties within the step | When the step has internal consistency requirements |

**Side-effect defaults**: When omitted, Side-effect defaults to `none`.
**Step-level Invariants defaults**: When omitted, no step-level invariant constraint.

#### 3.3 Semantic Descriptors

All dimension values use semantic descriptors -- natural language descriptions of expected behavior.

**Rules**:
- MUST NOT contain regex syntax (`\d`, `.*`, `[^...]`, `(?:...)`, `\s`, `\w`, `\b`, `$`, `^` as anchor, etc.)
- MUST NOT contain framework-specific assertion patterns
- MUST be natural language expressing business intent

**Good examples**:
- `"success confirmation containing feature-slug"`
- `"task status changed from pending to in_progress"`
- `"stderr contains error message about missing feature"`

**Bad examples** (these belong in gen-test-scripts):
- `"Feature\s+([\w-]+)\s+created"` (regex)
- `"assert.Equal(t, 0, exitCode)"` (framework assertion)
- `"matches pattern /task_\d+/"` (regex reference)

<HARD-RULE>
Semantic descriptors MUST NOT contain regex syntax. gen-contracts stage does not generate regex. If you find yourself writing a pattern match, replace it with a natural language description of what the pattern matches.
</HARD-RULE>

#### 3.4 Multi-Outcome Preconditions Mutual Exclusivity

Each Outcome within a Step MUST have Preconditions that are mutually exclusive with all other Outcomes in the same Step.

**Mutual exclusivity rule**: For any given system state, at most one Outcome's Preconditions can be satisfied. This prevents combinatorial explosion.

**Validation**: Before writing a Contract, verify that no two Outcomes in the same Step have identical or overlapping Preconditions. If overlap is detected:
1. Differentiate the Preconditions (add a distinguishing condition)
2. If impossible, merge the Outcomes into a single Outcome with disjunctive Preconditions
3. Never write Outcomes whose Preconditions can be simultaneously satisfied

<HARD-RULE>
Outcomes MUST be mutually exclusive by Preconditions. If two Outcomes' Preconditions can both be true for the same system state, the Contract is invalid and must be fixed before writing.
</HARD-RULE>

**Outcome count checkpoint**: Steps with more than 5 Outcomes trigger a review. Consider merging semantically similar Outcomes. Do not automatically exceed 5 Outcomes without explicit justification.

#### 3.5 TUI Async Cmd Await Semantics

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

#### 3.6 State Verification Levels

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

#### 3.7 Journey-Level Invariants

Every Contract file MUST end with a `## Journey Invariants` section containing at least one cross-step invariant.

**Source**: Journey-level Invariants come from the Journey document's `## Journey Invariants` section. Copy them verbatim into every Contract file for that Journey.

**At least 1 invariant is mandatory**. If the Journey document has no declared Invariants, generate at least one from the workflow analysis (e.g., "feature_slug consistent across all steps" or "working directory unchanged between steps").

#### 3.8 Batch Processing

When a single Journey has more than 15 Contracts or the estimated token count exceeds 50k, automatically split into multiple batches:

- **Batch 1**: Happy path Outcomes (all steps' success Outcomes)
- **Batch 2+**: Edge case Outcomes grouped by semantic similarity

Split batches are merged back into complete Contract files (one per Step). The merged result must be structurally identical to a single-batch generation.

<HARD-RULE>
Batch splitting occurs within a Journey, not across Journeys. Each Journey is always processed completely before moving to the next. Merged Contract files must not lose or duplicate any Outcome or dimension.
</HARD-RULE>

### Step 4: Validate Contracts

After generating all Contracts for a Journey, validate each one:

| Check | Rule |
|-------|------|
| Mandatory dimensions | Each Outcome MUST have non-empty: Preconditions, Input, Output, State |
| Semantic descriptor purity | No dimension value may contain regex syntax |
| Outcome name uniqueness | Outcome names within a Step MUST be unique |
| Preconditions mutual exclusivity | Different Outcomes' Preconditions MUST be distinguishable |
| Journey Invariants | Every Contract file MUST have a `## Journey Invariants` section with at least 1 entry |
| Side-effect default | When Side-effect is omitted or empty, it defaults to `none` |
| Outcome count checkpoint | Steps with > 5 Outcomes trigger a review warning |
| Unclassified validation points | Any validation point that cannot be mapped to a dimension MUST go to Invariants with `dimension: unclassified` annotation |

**Validation failure handling**:
- If mandatory dimensions are empty: fix the Contract (add content from Journey + Fact Table)
- If semantic descriptors contain regex: rewrite as natural language
- If Preconditions are not mutually exclusive: differentiate or merge Outcomes
- If Journey Invariants are missing: generate from workflow analysis

<HARD-RULE>
- Semantic descriptors MUST NOT contain regex syntax.
- Outcome Preconditions MUST be mutually exclusive.
- Steps with > 5 Outcomes trigger an LLM review checkpoint.
- Validation points that cannot be classified into existing dimensions MUST go to Invariants with `dimension: unclassified` annotation.
</HARD-RULE>

### Step 5: Write Output

Write Contract files to `tests/<journey>/_contracts/`.

**File naming**: `step-<N>-<action-slug>.md` where:
- `<N>` is the 1-based step ordinal
- `<action-slug>` is a kebab-case summary of the Step's primary action

**Template**: Use `${CLAUDE_SKILL_DIR}/templates/contract.md` for the file structure. Use `${CLAUDE_SKILL_DIR}/templates/outcome-block.md` for each Outcome block.

**Create directories**: Create `tests/<journey>/_contracts/` if it does not exist.

<HARD-RULE>
Output path is strictly `tests/<journey>/_contracts/`. No other locations.
</HARD-RULE>

## Error Handling

| Situation | Action |
|-----------|--------|
| Journey manifest missing | Abort with prompt to run `/gen-journeys` |
| Journey file not found | Abort with error listing the missing file path |
| Language detection fails | Ask user to configure `languages` in config.yaml |
| Interface detection fails | Ask user to configure `interfaces` in config.yaml |
| Source files not found for Fact Table | Mark as `UNKNOWN`, do not fabricate values |
| State verification level ambiguous | Default to `partial`, annotate with comment |
| Mandatory dimension empty after generation | Fix using Journey + Fact Table context, retry once |
| Semantic descriptor contains regex | Rewrite as natural language |
| Preconditions not mutually exclusive | Differentiate or merge Outcomes |
| Journey Invariants missing | Generate from workflow analysis |

## Related Skills

| Skill | Usage |
|-------|-------|
| `/gen-journeys` | Generate Journey documents (input source) |
| `/gen-test-scripts` | Generate executable test code from Contracts |
| `/run-e2e-tests` | Execute test scripts and report results |

## Reference

The authoritative model definition is at `docs/features/<slug>/design/model-and-directory-spec.md` (if it exists in the project). Key concepts used by this skill:

- **Contract**: Six-dimension verification mechanism for Journey Steps (Section 1.3)
- **Outcome**: A complete set of dimension declarations for a specific scenario (Section 1.4)
- **Semantic Descriptors**: Natural-language descriptions used in gen-contracts, converted to regex by gen-test-scripts (Section 1.5)
- **TUI Await Semantics**: Async Cmd wait specification with fail-fast timeout (Section 6)
- **State Verification Levels**: Full / partial / deferred degradation path (Section 2.3)
- **Batch Processing**: Auto-split when Contracts > 15 or tokens > 50k (Section 5.3)
