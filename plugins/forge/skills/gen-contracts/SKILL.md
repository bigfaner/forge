---
name: gen-contracts
description: Generate Contract specifications (six dimensions + semantic descriptors + multi-Outcome + Invariants) from Journey documents and code reconnaissance (Fact Table). Formalizes TUI async Cmd await semantics. Risk-driven Outcome density based on Journey risk_level, auto-derived boundary/error Outcomes from surface rules and Fact Table.
---

# Gen Contracts

Generate Contract specifications from Journey documents, enriched with code reconnaissance (Fact Table).

**Core principle**: Every Step gets a Contract with six dimensions. gen-contracts uses *semantic descriptors* (natural language) -- precise regex is deferred to gen-test-scripts. This skill is the most technically complex in the pipeline because it bridges narrative Journeys with verifiable Contracts using code reconnaissance.

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
| gen-contracts | Journey documents + code reconnaissance | Contract specifications (six dimensions, semantic descriptors) | Yes (Fact Table) |

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

1. Read `docs/conventions/testing/index.md` to discover available Convention files. Select the Convention matching the project's language/framework based on index descriptions and project context. Load the selected Convention from `docs/conventions/testing/<convention>.md`.
2. Fallback: scan existing source/test files (`go.mod`, `package.json`, `*_test.go`, etc.). Also check subdirectories for monorepo.
3. On failure: ask user.
4. **Detect surfaces**: Check `.forge/config.yaml` `surfaces` field, `docs/conventions/`, project directory structure, and dependencies for surface types (cli, api, tui, web, mobile).

<HARD-RULE>
Do NOT silently default to any language or surface.
</HARD-RULE>

## Process Flow

```
0. Resolve language + surfaces -> 1. Read Journeys -> 2. Code Reconnaissance (Fact Table) -> 3. Generate Contracts (risk-driven density + boundary derivation) -> 4. Validate (schema + retry) -> 5. Write Output + Fact Table
```

### Step 1: Read Journey Documents

1. Enumerate subdirectories under `docs/features/<slug>/testing/` — each subdirectory represents one Journey.
2. For each Journey directory, read `journey.md` to get the Journey document.
3. Parse each Journey's structure:
   - Journey name and Risk level (from frontmatter `risk_level`: High/Medium/Low)
   - Happy path steps (sequence number, user action, expected result)
   - Edge cases (referenced step, precondition, user action, expected result)
   - Journey Invariants (cross-step properties)
4. Load the project's surface type from `.forge/config.yaml` and read the corresponding surface rule from gen-journeys skill's `rules/surface-<type>.md` (resolve relative to the gen-journeys skill directory) to identify required_outcomes.

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

### Step 3: Generate Contracts

For each Journey, generate one Contract file per Step. Apply risk-driven Outcome density per `rules/risk-density.md`.

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

For each Outcome, declare all six dimensions per `rules/dimension-rules.md`: four mandatory (Preconditions, Input, Output, State) and two optional (Side-effect, Invariants).

#### 3.3 Semantic Descriptors

All dimension values use semantic descriptors per `rules/dimension-rules.md`. MUST NOT contain regex syntax. MUST be natural language expressing business intent.

#### 3.4 Multi-Outcome Preconditions Mutual Exclusivity

Each Outcome within a Step MUST have mutually exclusive Preconditions per `rules/dimension-rules.md`. Outcome count checkpoint: steps with > 5 Outcomes trigger a review.

#### 3.5 Risk-Driven Outcome Density

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
- WebUI: `validation-error` (form submission Steps), `session-expired` (session-dependent Steps)
- Mobile: best-effort, no mandatory Outcomes
</HARD-RULE>

<HARD-RULE>
LLM-derived boundary Outcomes MUST be annotated with `source: inferred` and include a reasoning explanation citing Fact Table sources or inference logic.
</HARD-RULE>

#### 3.6 Density Checkpoint

After generating Outcomes for all Steps in a Journey, output a density checkpoint (format in `rules/risk-density.md`). If actual total is below target, review Steps for missed boundary scenarios. If above target, merge semantically similar Outcomes.

#### 3.7 TUI Async Cmd Await Semantics

For TUI Steps involving async operations, declare `await` semantics per `rules/tui-async.md`, including timeout outcomes for async Cmds.

#### 3.8 State Verification Levels

Determine state verification level (full/partial/deferred) from Fact Table reconnaissance per `rules/tui-async.md`.

#### 3.9 Journey-Level Invariants

Every Contract file MUST end with a `## Journey Invariants` section per `rules/tui-async.md`. At least 1 invariant is mandatory.

#### 3.10 Batch Processing

Auto-split into batches when Contracts > 15 or tokens > 50k per `rules/tui-async.md`.

### Step 4: Validate Contracts (Schema Validation + Retry)

After generating all Contracts for a Journey, validate each one per `rules/validation.md`. Apply validation checks and failure handling as defined in the rules file.

**Schema validation** checks (6-dimension structural completeness):

| Check | Rule |
|-------|------|
| Mandatory dimensions | Each Outcome has non-empty: Preconditions, Input, Output, State |
| Semantic descriptor purity | No dimension value contains regex syntax |
| Outcome name uniqueness | Outcome names within a Step are unique |
| Preconditions mutual exclusivity | Different Outcomes' Preconditions are distinguishable and non-overlapping |
| Journey Invariants present | Every Contract file has `## Journey Invariants` section with >= 1 entry |
| Side-effect default | When Side-effect is omitted or empty, defaults to `none` |

**Retry logic**:

1. If schema validation fails, record the non-compliance items (which Contract, which Outcome, which dimension/check failed)
2. Automatically regenerate the failed Contracts **once**, injecting the schema errors as feedback into the generation prompt
3. Re-validate the regenerated Contracts
4. If validation still fails after retry: **pause the pipeline** and output the non-compliance items for manual correction

<HARD-RULE>
Schema validation MUST be executed after generation. Validation failure triggers exactly 1 automatic retry with error feedback. If the retry also fails, the pipeline pauses -- it does NOT silently continue with invalid Contracts.
</HARD-RULE>

### Step 5: Write Output + Fact Table

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

## Related Skills

| Skill | Usage |
|-------|-------|
| `/gen-journeys` | Generate Journey documents (input source) |
| `/gen-test-scripts` | Generate executable test code from Contracts |
| `/run-tests` | Execute test scripts and report results |

## Reference

The authoritative model definition is at `docs/features/<slug>/design/model-and-directory-spec.md` (if it exists in the project). Key concepts used by this skill:

- **Contract**: Six-dimension verification mechanism for Journey Steps (Section 1.3)
- **Outcome**: A complete set of dimension declarations for a specific scenario (Section 1.4)
- **Semantic Descriptors**: Natural-language descriptions used in gen-contracts, converted to regex by gen-test-scripts (Section 1.5)
- **TUI Await Semantics**: Async Cmd wait specification with fail-fast timeout (Section 6)
- **State Verification Levels**: Full / partial / deferred degradation path (Section 2.3)
- **Batch Processing**: Auto-split when Contracts > 15 or tokens > 50k (Section 5.3)
