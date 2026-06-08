---
name: gen-contracts
description: Generate Contract specifications (six dimensions, semantic descriptors, risk-driven Outcomes) from Journey documents and code reconnaissance.
---

# Gen Contracts

**Core principle**: Every Step gets a Contract with six dimensions. gen-contracts uses *semantic descriptors* (natural language) -- precise regex is deferred to gen-test-scripts. When handbooks exist, Contract frontmatter also includes *technical anchors* (endpoint, command, page, screen) extracted from design documents, bridging design intent to test code.

**Risk-driven density**: Outcome count per Step and total test count per Journey are driven by the Journey's `risk_level` field (set by gen-journeys). Density targets are defined in `rules/risk-density.md`. Boundary/error Outcomes are auto-derived from surface rules' required_outcomes and the project Fact Table.

<HARD-GATE>
This skill ONLY writes to `docs/features/<slug>/testing/<journey>/contracts/`. It does NOT generate test scripts or execute tests (handled by downstream skills).

**FORBIDDEN output paths**: Any path outside `docs/features/<slug>/testing/<journey>/contracts/`. The `contracts/` directory is the sole output target.
</HARD-GATE>

## Pipeline Position

```
gen-journeys -> gen-contracts -> gen-test-scripts -> run-tests
     (narrative)     (technical)     (code gen)        (execution)
```

gen-contracts is the second step. It reads Journey documents and source code, producing Contract specifications that gen-test-scripts consumes.

| Step | Input | Output | Reads Code |
|------|-------|--------|------------|
| gen-contracts | Journey documents + code reconnaissance + handbooks | Contract specifications (six dimensions, semantic descriptors, technical anchors) | Yes (Fact Table) |

## Prerequisites

Check previous stage artifacts. Abort and prompt user if missing:

| Artifact | Missing prompt |
|----------|----------------|
| At least one Journey directory under `docs/features/<slug>/testing/` with `journey.md` | Run `/gen-journeys` first |
| Eval report for all Journeys (`testing/<journey>/.eval-report.md`) | Run `/eval --type journey` first. **Blocker**: do not proceed if any Journey scored below target. |

`<slug>` from `forge feature`.

### SKIP_EVAL_GATE Mode

When the task context contains `SKIP_EVAL_GATE=true` (injected by Quick mode task templates), the eval report prerequisite is **conditionally waived**:

- **Skip**: eval-journey report check (`testing/<journey>/.eval-report.md`) is bypassed entirely
- **Proceed directly**: move to Step 1 (Read Journeys) and Step 2 (Code Reconnaissance) without eval verification
- **Mark output**: every Contract file generated under SKIP_EVAL_GATE MUST include a frontmatter field `skip_eval: true` and a comment at the top of the body: `> **Note**: Contracts generated without eval-journey verification (SKIP_EVAL_GATE=true). Review with extra scrutiny.`

**When SKIP_EVAL_GATE is NOT set** (Breakdown mode or manual `/gen-contracts` invocation): the eval report Blocker remains mandatory. Behavior is unchanged.

## Step 0: Resolve Language and Surfaces

Load: `rules/journey-contract-model.md` — core concepts (Journey, Step, Contract, Outcome), directory conventions, and tag-based promotion model.

### 0.1 Surface Detection

<!-- INLINE:origin=gen-journeys/SKILL.md#Surface Detection -->

1. Run `forge surfaces` to list all configured surfaces
2. Parse stdout using the unified text parsing rule (see Forge Guide → Surface Output Parsing)
3. If exit code is 1: **pause pipeline** and ask user to configure surfaces via `forge init`
4. Load `rules/surface-<type>.md` for each detected surface type

**Supported surface types**: `web`, `api`, `cli`, `tui`, `mobile`

<!-- END INLINE:origin=gen-journeys/SKILL.md#Surface Detection -->

<HARD-RULE>
Never guess surface types — always use `forge surfaces`. Do NOT pass `.` or arbitrary paths; path-based lookup requires a surface key prefix.
</HARD-RULE>

### 0.2 Convention Loading (Surface-First)

Load surface-specific Convention files for each detected surface type:

1. **Legacy detection**: Check if `docs/conventions/testing/` contains any `.md` files that are NOT inside a subdirectory (i.e., flat files like `go.md`, `vitest.md`). If legacy files are detected:
   - Output migration prompt: "Legacy Convention structure detected in `docs/conventions/testing/` (framework-first files). Run `/test-guide` to regenerate with the new surface-first structure (`testing/{surface}/core.md`)."
   - Proceed without loading the legacy files.
2. **Load surface Convention**: For each detected surface type, load `docs/conventions/testing/{surface}/core.md`.
3. **No Convention found**: Proceed with LLM defaults. Output hint: "No test Convention files found for surface `{surface}` in `docs/conventions/testing/{surface}/core.md`. Generation will use LLM defaults. Run `/test-guide` to create one."
4. **Resolve framework**: If `core.md` was loaded, read its assertion preference table to identify the target framework. Otherwise scan existing source/test files (`go.mod`, `package.json`, `*_test.go`, etc.) for auto-detection. On failure: ask user.

<HARD-RULE>
Convention loading is surface-driven, not framework-driven. The `{surface}` segment comes from Step 0.1 surface detection via `forge surfaces`. Do NOT fall back to loading framework-specific flat files.
</HARD-RULE>

<HARD-RULE>
Do NOT silently default to any language or surface.
</HARD-RULE>

## Process Flow

```
0. Resolve surfaces (forge surfaces CLI) + Convention loading (surface-first) -> 1. Read Journeys -> 2. Code Reconnaissance (Fact Table) -> 3. Load Handbooks + Anchor Filling -> 4. Generate Contracts (risk-driven density + boundary derivation + anchors) -> 5. Validate (schema + retry) -> 6. Write Output + Fact Table
```

### Step 1: Read Journey Documents

1. Enumerate subdirectories under `docs/features/<slug>/testing/` — each subdirectory represents one Journey.
2. For each Journey directory, read `journey.md` to get the Journey document.
3. Parse each Journey's structure:
   - Journey name and Risk level (from frontmatter `risk_level`: High/Medium/Low)
   - Happy path steps (sequence number, user action, expected result)
   - Edge cases (referenced step, precondition, user action, expected result)
   - Journey Invariants (cross-step properties)
4. For each detected surface type, load the surface-required Outcomes from `rules/risk-density.md` (Surface-Required Outcome Derivation table) to identify required_outcomes for the detected surface types.

<HARD-RULE>
Every Journey in the manifest MUST be processed. Do not skip Journeys based on Risk level or step count.
</HARD-RULE>

### Step 2: Code Reconnaissance (Build Fact Table)

Read source code to extract ground-truth values. Follow the full reconnaissance procedure per `rules/code-reconnaissance.md`, including generic and TUI-specific reconnaissance tables, Fact Table format, and source citation rules.

**Static Fact Table output**: After completing reconnaissance, write the Fact Table to `.forge/fact-table.json` in the project root. Each fact entry follows the canonical schema (defined in `forge-cli/pkg/facttable/facttable.go`):

```json
{
  "fact_id": "<FACT_ID>",
  "source": "static",
  "subject": "<what this fact describes>",
  "kind": "<signature | output_format | error_code | side_effect | precondition>",
  "value": "<extracted value or JSON object>",
  "confidence": "inferred",
  "updated_at": "<ISO8601 timestamp>"
}
```

All entries use `"source": "static"` to distinguish from runtime facts (added by Run-to-Learn with `"source": "runtime"`). Static entries default to `"confidence": "inferred"` (confirmed at runtime by R2L).

### Step 3: Load Handbooks + Anchor Filling

For each detected surface type, attempt to load the corresponding handbook. Handbooks provide authoritative technical anchor values (endpoint, command, page, screen) that are written into Contract frontmatter during generation (Step 4).

#### 3.1 Handbook Location

Handbooks are generated by `/tech-design` and reside under the feature directory:

| Surface | Handbook File | Location Pattern |
|---------|--------------|------------------|
| API | `api-handbook.md` | `docs/features/<slug>/api-handbook.md` |
| CLI | `cli-handbook.md` | `docs/features/<slug>/cli-handbook.md` |
| TUI | `cli-handbook.md` | `docs/features/<slug>/cli-handbook.md` (shared with CLI) |
| Web | `page-map.md` | `docs/features/<slug>/page-map.md` |
| Mobile | `screen-map.md` | `docs/features/<slug>/screen-map.md` |

<HARD-RULE>
Handbook paths use the surface key from `forge surfaces` (Step 0.1). Do NOT guess or invent paths — construct from the slug and handbook type using the table above.
</HARD-RULE>

#### 3.2 Handbook Freshness Check

For each handbook that exists, compare its freshness against the tech-design document:

1. Extract the `generated` (or `updated`) timestamp from the handbook's frontmatter
2. Find the tech-design file at `docs/features/<slug>/tech-design.md`
3. Compare: if `handbook.generated < tech-design.last_modified`, the handbook is **stale**
4. On stale detection: output a warning — `"⚠ Handbook [name] may be stale (generated before tech-design update). Consider re-running /tech-design to regenerate."`
5. Proceed with the stale handbook — do NOT abort the pipeline. The warning informs the user but does not block generation

#### 3.3 Missing Handbook Graceful Degradation

When a handbook file does not exist for a surface type:

1. Skip anchor filling for that surface type entirely
2. Output a hint — `"Missing handbook for surface [type]. Anchor fields will be empty. Consider running /tech-design to generate."`
3. Continue pipeline execution — do NOT abort. Contracts are still generated without anchor values

<HARD-RULE>
Missing handbooks MUST NOT cause pipeline failure. Anchor fields default to empty strings/arrays per `templates/contract.md`. The pipeline degrades gracefully to existing behavior (no anchors filled).
</HARD-RULE>

#### 3.4 Anchor Extraction

For each handbook that exists and is fresh (or stale but usable), extract anchor values:

| Surface | Extraction Rules |
|---------|-----------------|
| API | Parse `api-handbook.md` entries. For each endpoint, extract `endpoint` (path pattern), `method` (HTTP verb), and optional `content_type`, `auth_required` into the `anchors.api` section |
| CLI | Parse `cli-handbook.md` entries. For each command, extract `command` (top-level), `subcommand` (nested), and optional `flags`, `aliases` into the `anchors.cli` section |
| TUI | Parse `cli-handbook.md` entries. For interactive terminal commands, extract `command` (entry command), `interactive_prompt`, and `keybindings` into the `anchors.tui` section |
| Web | Parse `page-map.md` entries. For each page, extract `page` (name), `route` (URL path), `requires_auth`, and `layout` into the `anchors.web` section |
| Mobile | Parse `screen-map.md` entries. For each screen, extract `screen` (name), `navigation_path`, `deeplink`, and `platform` into the `anchors.mobile` section |

**Anchor-to-Step matching**: Anchors are matched to Contract Steps by correlating the Journey Step's user action with handbook entries. When no clear match exists, the anchor fields for that Step remain empty (not guessed).

<HARD-RULE>
Anchors come from handbooks (design documents) as the authority source. Do NOT reverse-engineer anchors from source code. When a handbook entry cannot be matched to a Step, leave the anchor empty rather than guessing.
</HARD-RULE>

#### 3.5 Anchor Sync Timestamp

When anchors are successfully filled for any Step, set `last_anchor_sync` in the Contract frontmatter to the current ISO-8601 timestamp. This timestamp is updated every time anchors are filled or refreshed, enabling downstream tools to detect anchor staleness.

### Step 4: Generate Contracts

For each Journey, generate one Contract file per Step. Apply risk-driven Outcome density per `rules/risk-density.md`. Anchor values from Step 3 are written into each Contract's frontmatter `anchors` section per `templates/contract.md`.

#### 4.1 Step-to-Contract Mapping

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

#### 4.2 Six-Dimension Declaration

For each Outcome, declare all six dimensions per `rules/dimension-rules.md`: four mandatory (Preconditions, Input, Output, State) and two optional (Side-effect, Invariants).

#### 4.2.1 Fixture Specification (Preconditions Sub-Dimension)

Each Outcome's Preconditions MUST include a `fixture_spec` field per `rules/fixture-spec.md`. This field declaratively specifies the pre-existing data state (entities, relationships, minimum counts) required before the Step executes.

**fixture_spec is required** in every Outcome's Preconditions. A Contract without fixture_spec is schema-invalid. When fixture_spec is absent from a legacy Contract, downstream consumers must fall back to implicit inference (see Backward Compatibility in `rules/fixture-spec.md`).

#### 4.3 Semantic Descriptors

All dimension values use semantic descriptors per `rules/dimension-rules.md`. MUST NOT contain regex syntax. MUST be natural language expressing business intent.

#### 4.4 Multi-Outcome Preconditions Mutual Exclusivity

Each Outcome within a Step MUST have mutually exclusive Preconditions per `rules/dimension-rules.md`. Outcome count checkpoint: steps with > 5 Outcomes trigger a review.

#### 4.5 Risk-Driven Outcome Density

Apply density targets from `rules/risk-density.md` based on the Journey's `risk_level`:

| Risk Level | Outcomes per Step | Total (Journey) |
|------------|-------------------|-----------------|
| High       | 3-5               | 13-20           |
| Medium     | 2-3               | 8-12            |
| Low        | 1-2               | 4-7             |

**Outcome generation priority** (in order):
1. Happy path Outcome (always, counts as 1)
2. Surface-required Outcomes (from surface rule's Required Outcome Reference section)
3. Fact Table-informed boundary Outcomes (based on code reconnaissance findings, annotated with `source: inferred` + reasoning)
4. LLM-inferred edge case Outcomes (for High/Medium risk only, annotated with `source: inferred` + reasoning)

<HARD-RULE>
Surface-required Outcomes MUST be derived for every matching Step. They are not optional. Required Outcomes are defined in the surface rule's "Required Outcome Reference" section:
- CLI: `not-found` (resource access Steps), `already-exists` (resource creation Steps)
- API: `unauthorized` (authenticated endpoint Steps)
- TUI: `timeout` (async Cmd Steps, per rules/tui-async.md)
- Web: `validation-error` (form submission Steps), `session-expired` (session-dependent Steps)
- Mobile: best-effort, no mandatory Outcomes
</HARD-RULE>

<HARD-RULE>
LLM-derived boundary Outcomes MUST be annotated with `source: inferred` and include a reasoning explanation citing Fact Table sources or inference logic.
</HARD-RULE>

#### 4.6 Density Checkpoint

After generating Outcomes for all Steps in a Journey, output a density checkpoint (format in `rules/risk-density.md`). If actual total is below target, review Steps for missed boundary scenarios. If above target, merge semantically similar Outcomes.

#### 4.7 TUI Async Cmd Await Semantics *(TUI-specific)*

**Applies to: TUI surface only.** For TUI Steps involving async operations, declare `await` semantics per `rules/tui-async.md`, including timeout outcomes for async Cmds. Skip this step for non-TUI surfaces.

#### 4.8 State Verification Levels *(applies to all surface types)*

**Applies to: all surface types.** Determine state verification level (full/partial/deferred) from Fact Table reconnaissance per `rules/tui-async.md` (State Verification Levels section).

#### 4.9 Journey-Level Invariants *(applies to all surface types)*

**Applies to: all surface types.** Every Contract file MUST end with a `## Journey Invariants` section per `rules/tui-async.md` (Journey-Level Invariants section). At least 1 invariant is mandatory.

#### 4.10 Batch Processing *(applies to all surface types)*

**Applies to: all surface types.** Auto-split into batches when Contracts > 15 or tokens > 50k per `rules/tui-async.md` (Batch Processing section).

### Step 5: Validate Contracts (Schema Validation + Retry)

After generating all Contracts for a Journey, validate each one per `rules/validation.md`. Apply validation checks and failure handling as defined in the rules file.

**Schema validation** checks (6-dimension structural completeness):

| Check | Rule |
|-------|------|
| Mandatory dimensions | Each Outcome has non-empty: Preconditions, Input, Output, State |
| Fixture Specification present | Each Outcome's Preconditions includes `fixture_spec` with at least 1 entity declaration (see `rules/fixture-spec.md`) |
| Semantic descriptor purity | No dimension value contains regex syntax |
| Outcome name uniqueness | Outcome names within a Step are unique |
| Preconditions mutual exclusivity | Different Outcomes' Preconditions are distinguishable and non-overlapping |
| Journey Invariants present | Every Contract file has `## Journey Invariants` section with >= 1 entry |
| Side-effect default | When Side-effect is omitted or empty, defaults to `none` |
| Anchor frontmatter structure | When handbook exists, the `anchors` section for the matching surface type is populated with at least the required fields per `templates/contract.md` |
| Anchor sync timestamp | When any anchor is filled, `last_anchor_sync` is set to an ISO-8601 timestamp |

**Retry logic**:

1. If schema validation fails, record the non-compliance items (which Contract, which Outcome, which dimension/check failed)
2. Automatically regenerate the failed Contracts **once**, injecting the schema errors as feedback into the generation prompt
3. Re-validate the regenerated Contracts
4. If validation still fails after retry: **pause the pipeline** and output the non-compliance items for manual correction

<HARD-RULE>
Schema validation MUST be executed after generation. Validation failure triggers exactly 1 automatic retry with error feedback. If the retry also fails, the pipeline pauses -- it does NOT silently continue with invalid Contracts.
</HARD-RULE>

### Step 6: Write Output + Fact Table

Write Contract files to `docs/features/<slug>/testing/<journey>/contracts/`.

**File naming**: `step-<N>-<action-slug>.md` where:
- `<N>` is the 1-based step ordinal
- `<action-slug>` is a kebab-case summary of the Step's primary action

**Template**: Use `templates/contract.md` for the file structure. Use `templates/outcome-block.md` for each Outcome block.

**Create directories**: Create `docs/features/<slug>/testing/<journey>/contracts/` if it does not exist.

**Static Fact Table**: Write the Fact Table (from Step 2 reconnaissance) to `.forge/fact-table.json` in the project root. All entries have `"source": "static"`. If `.forge/fact-table.json` already exists, merge new entries by `fact_id` (do not delete existing runtime entries).

<HARD-RULE>
Output path is strictly `docs/features/<slug>/testing/<journey>/contracts/`. No other locations. Static Fact Table must be written to `.forge/fact-table.json`.
</HARD-RULE>

## Error Handling

See `rules/validation.md` for the complete error handling table.

## Key Concepts

- **Contract**: Six-dimension verification mechanism for Journey Steps
- **Outcome**: A complete set of dimension declarations for a specific scenario
- **Semantic Descriptors**: Natural-language descriptions used in gen-contracts, converted to regex by gen-test-scripts
- **Technical Anchors**: Authoritative technical metadata (endpoint, command, page, screen) extracted from handbooks into Contract frontmatter, bridging design intent to test code
- **Handbook**: Design document generated by `/tech-design` containing surface-specific technical details (api-handbook, cli-handbook, page-map, screen-map)
- **Anchor Sync**: Process of filling anchor fields from handbooks; tracked via `last_anchor_sync` timestamp
- **TUI Await Semantics**: Async Cmd wait specification with fail-fast timeout
- **State Verification Levels**: Full / partial / deferred degradation path
- **Batch Processing**: Auto-split when Contracts > 15 or tokens > 50k
