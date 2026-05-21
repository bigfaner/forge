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
| At least one Journey file in `docs/features/<slug>/testing/journeys/` | Run `/gen-journeys` first |

`<slug>` from `forge feature`.

## Step 0: Resolve Language and Interfaces

1. Load Convention files from `docs/conventions/` by `domains` frontmatter (match `testing`, `go`, `typescript`, etc.). Extract language from `Framework` section.
2. Fallback: scan existing source/test files (`go.mod`, `package.json`, `*_test.go`, etc.). Also check subdirectories for monorepo.
3. On failure: ask user.
4. **Detect interfaces**: Check `.forge/config.yaml`, `docs/conventions/`, project directory structure, and dependencies for interface types (cli, api, tui, web-ui, mobile).

<HARD-RULE>
Do NOT silently default to any language or interface.
</HARD-RULE>

## Convention Loading

After language/interface resolution and before entering the workflow steps, load project conventions into context.

**Resolution algorithm**:

1. **Project-wide conventions**: Read this skill's own frontmatter `conventions` field. For each filename, check `docs/conventions/{filename}` -- if it exists, read it into context; if missing, skip silently.
2. **Interface-specific conventions**: For each detected interface, check `docs/conventions/testing-{interface}.md` -- if it exists, read it into context; if missing, skip silently.

<HARD-RULE>
Convention loading is non-blocking. Missing convention files are silently skipped.
</HARD-RULE>

## Process Flow

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

Read source code to extract ground-truth values. Follow the full reconnaissance procedure per `rules/code-reconnaissance.md`, including generic and TUI-specific reconnaissance tables, Fact Table format, and source citation rules.

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

For each Outcome, declare all six dimensions per `rules/dimension-rules.md`: four mandatory (Preconditions, Input, Output, State) and two optional (Side-effect, Invariants).

#### 3.3 Semantic Descriptors

All dimension values use semantic descriptors per `rules/dimension-rules.md`. MUST NOT contain regex syntax. MUST be natural language expressing business intent.

#### 3.4 Multi-Outcome Preconditions Mutual Exclusivity

Each Outcome within a Step MUST have mutually exclusive Preconditions per `rules/dimension-rules.md`. Outcome count checkpoint: steps with > 5 Outcomes trigger a review.

#### 3.5 TUI Async Cmd Await Semantics

For TUI Steps involving async operations, declare `await` semantics per `rules/tui-async.md`, including timeout outcomes for async Cmds.

#### 3.6 State Verification Levels

Determine state verification level (full/partial/deferred) from Fact Table reconnaissance per `rules/tui-async.md`.

#### 3.7 Journey-Level Invariants

Every Contract file MUST end with a `## Journey Invariants` section per `rules/tui-async.md`. At least 1 invariant is mandatory.

#### 3.8 Batch Processing

Auto-split into batches when Contracts > 15 or tokens > 50k per `rules/tui-async.md`.

### Step 4: Validate Contracts

After generating all Contracts for a Journey, validate each one per `rules/validation.md`. Apply validation checks and failure handling as defined in the rules file.

### Step 5: Write Output

Write Contract files to `tests/<journey>/_contracts/`.

**File naming**: `step-<N>-<action-slug>.md` where:
- `<N>` is the 1-based step ordinal
- `<action-slug>` is a kebab-case summary of the Step's primary action

**Template**: Use `templates/contract.md` for the file structure. Use `templates/outcome-block.md` for each Outcome block.

**Create directories**: Create `tests/<journey>/_contracts/` if it does not exist.

<HARD-RULE>
Output path is strictly `tests/<journey>/_contracts/`. No other locations.
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
