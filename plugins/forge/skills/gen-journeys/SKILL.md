---
name: gen-journeys
description: Extract Journey narratives (user workflows with Risk classification) from PRD user stories. Each Journey is a per-file Markdown document describing happy path + edge cases + invariants.
---

# Gen Journeys

Extract Journey narratives from PRD user stories, outputting per-Journey Markdown files.

**Core principle**: Pure narrative extraction -- no code reconnaissance required. Each user story from the PRD maps to a Journey describing the user's real workflow to accomplish a goal.

<HARD-GATE>
This skill only generates Journey narrative documents (per-Journey Markdown files). It does NOT generate Contracts, test scripts, or executable code. Those are handled by downstream skills:
- `/gen-contracts` -- generates Contract specifications from Journey documents
- `/gen-test-scripts` -- generates executable test code from Contracts
</HARD-GATE>

## Core Concepts

<!-- INLINE:origin=gen-contracts/rules/journey-contract-model.md -->

The Forge test pipeline organizes around user workflows (Journeys), defines expected behavior through six-dimension Contract declarations, and manages test lifecycles via Tag-Based Promotion.

### Journey

A Journey describes a real user workflow for achieving a goal. It is the primary organizational unit for testing.

| Property | Description |
|----------|-------------|
| Name | kebab-case identifier (e.g., `task-lifecycle`) |
| Risk | `High` (state changes / data loss risk), `Medium` (multi-step interactions without irreversible side effects), `Low` (read-only operations) |
| Steps | Ordered sequence of user actions, each with expected outcomes |
| Invariants | Cross-step constraints that must hold throughout the entire Journey |

Each Journey executes in its own temporary working directory to prevent cross-contamination during parallel execution.

### Step

A Step is a single user action within a Journey. Each Step maps to a Contract containing one or more Outcomes.

| Property | Description |
|----------|-------------|
| Sequence number | 1-based index within the Journey |
| User action | The operation the user performs (running a command, clicking a button, sending a request, etc.) |
| Expected outcomes | One or more Outcome declarations, each with independent Preconditions |

### Contract

A Contract is the verification mechanism for a Step, defining expected system behavior through six-dimension declarations. All dimensions are declared at the Outcome level; Invariants are additionally declared at the Journey level.

### Outcome

An Outcome is a complete set of Contract dimension declarations for a specific scenario (success, error variant, edge case). Outcomes within the same Step are distinguished by Preconditions and must be mutually exclusive — at most one Outcome's Preconditions can be satisfied for any given system state.

| Property | Required | Description |
|----------|----------|-------------|
| Name | Yes | Descriptive label (e.g., `success`, `not-in-progress`) |
| Preconditions | Yes | System state required for this Outcome to become active |
| Input | Yes | Input provided by the user to the system |
| Output | Yes | Output produced by the system (semantic descriptors) |
| State | Yes | System state changes |
| Side-effect | No | External side effects (default: `none`) |
| Invariants | No | Step-level invariants (default: no constraints) |

### Semantic Descriptors

All dimension values use semantic descriptors -- natural language descriptions of expected behavior that express business intent rather than precise matching patterns. Regex is prohibited in descriptors.

### Contract File Format

Each Contract is stored as a structured Markdown file:

```markdown
# Contract: <journey-name> / Step <N>: <step-description>

## Outcome "<outcome-name>"
- Preconditions: "<semantic description>"
- Input: <semantic description>
- Output: <semantic description>
- State: <semantic description>
- Side-effect: <semantic description or "none">
- Invariants: <step-level invariants or omit>

## Journey Invariants
- <invariant description 1>
```

**Parseable Structure Rules**:
1. Journey name extracted from file path: `docs/features/<slug>/testing/<journey>/contracts/step-N-*.md`
2. Step sequence extracted from filename: `step-<N>-<slug>.md`
3. Outcome sections: `## Outcome "<name>"` headings declare new Outcome blocks
4. Dimension format: `- <DimensionName>: <value>`
5. `## Journey Invariants` section MUST appear exactly once in each Contract file

### Directory Convention

```
docs/features/<slug>/testing/
  <journey-name>/                     # Journey directory (kebab-case)
    journey.md                        # Journey narrative document
    contracts/                        # Contract specification directory
      step-1-<action-slug>.md         # Contract for Step 1
      step-N-<action-slug>.md

tests/                                # Generated test files (by gen-test-scripts)
  <journey-name>/                     # Single surface: flat structure
  <surfaceKey>/<journey-name>/        # Multi surface: partitioned by surface key
```

**Test output path rules** (applied by gen-test-scripts, not gen-journeys):
- **Single surface**: `tests/<journey-name>/` — no surface-key directory layer
- **Multi surface**: `tests/<surfaceKey>/<journey-name>/` — surface-key partitions test files

### Tag-Based Promotion

Tests manage their lifecycle via tags rather than file movement:

| Stage | Tag | Action |
|-------|-----|--------|
| New | `@feature` | Automatically injected into newly generated tests |
| Promoted | `@regression` | Automatically upgraded by `/run-tests` |

<!-- END INLINE -->

## Surface Detection

1. Run `forge surfaces` to list all configured surfaces
2. Parse stdout using the unified text parsing rule (see Forge Guide → Surface Output Parsing)
3. If exit code is 1: **pause pipeline** and ask user to configure surfaces via `forge init`
4. Load `rules/surface-<type>.md` for each detected surface type

**Supported surface types**: `web`, `api`, `cli`, `tui`, `mobile`

<HARD-RULE>
Never guess surface types — always use `forge surfaces`. Do NOT pass `.` or arbitrary paths; path-based lookup requires a surface key prefix.
</HARD-RULE>

## Multi-Surface Rules Loading

When the project has multiple configured surface types (e.g., `web` + `api`), load ALL detected surface rule files to inform Journey generation. gen-journeys is a narrative extraction skill — surface rules serve as reference guidance for downstream stages, not as primary input. Loading multiple rule files does not significantly increase context noise.

### Loading Strategy

1. After surface detection, load `rules/surface-<type>.md` for each detected surface type
2. Organize rule loading by surface type — load each rule file sequentially
3. Collect the union of:
   - Required Outcomes from each surface's "Required Outcome Reference" section
   - Test level emphasis ratios from each surface's "Test Strategy Guidance" section
   - Risk-level Outcome density targets adjusted by each surface's guidance

### Per-Surface Rule Application

When generating Journeys, apply the following per-surface guidance:

#### API — HTTP status boundaries (4xx/5xx), auth failures, rate limiting, payload validation. Emphasis: 50/50 Contract/Journey.

#### Web — Page load states, navigation transitions, form validation. Emphasis: 50/50 Contract/Journey.

#### CLI — Exit codes, stdout/stderr, signal handling. Emphasis: unit-heavy.

#### TUI — Rendering states, keyboard navigation, screen transitions. Emphasis: 80/20 Contract/Journey.

#### Mobile — Screen transitions, gestures, offline/online state. Emphasis: e2e-heavy.

### Journey Surface Coverage

Each Journey must declare which surfaces it covers in its frontmatter, using both `surface_types` (technical classification) and `surface_keys` (user-defined identifiers). This enables downstream skills (gen-contracts, gen-test-scripts) to generate surface-appropriate Contracts and test scripts, and to correctly partition test output directories by surface key.

**Coverage rules**:
- A Journey covering a cross-surface user workflow (e.g., "user submits form via web, backend API processes it") must list all involved surface types and keys
- A Journey limited to a single surface's interaction lists only that surface's type and key
- The union of all Journeys' `surface_types` must cover every configured surface type — no surface may be left uncovered
- The union of all Journeys' `surface_keys` must cover every configured surface key — no surface key may be left uncovered

### Extensibility

New surface types can be added by creating a new `rules/surface-<type>.md` file following the same 4-section structure (Detection Signals, General Testing Principles, Test Strategy Guidance, Required Outcome Reference). No pipeline code changes are needed. Add the corresponding subsection above following the same format.

## Prerequisites

gen-journeys supports two input modes, determined by file availability:

### PRD Mode (default)

| Artifact | Required | Missing prompt |
|----------|----------|----------------|
| `docs/features/<slug>/prd/prd-user-stories.md` | Yes | Run `/write-prd` first |
| `docs/features/<slug>/prd/prd-spec.md` | Yes | Run `/write-prd` first |

### Proposal Mode (Quick mode fallback)

When both PRD files (`prd-user-stories.md` and `prd-spec.md`) do not exist, gen-journeys automatically switches to Proposal Mode, using `proposal.md` as the input source.

| Artifact | Required | Missing prompt |
|----------|----------|----------------|
| `docs/proposals/<slug>/proposal.md` | Yes (conditional) | Proposal file not found. Cannot generate Journeys without PRD or Proposal input. |

**Proposal Mode minimum information check**: Before generating Journeys from proposal.md, verify these mandatory fields exist:

| Field | Requirement | Abort behavior |
|-------|-------------|----------------|
| `## Scope` (or `### Scope`) | **Must exist** — provides feature boundary for Journey extraction | Abort with diagnostic: "proposal.md missing Scope section — cannot determine feature boundaries for Journey generation" |
| `## Success Criteria` (or `### Success Criteria`) | **Must exist** — provides verifiable outcomes for Journey expected results | Abort with diagnostic: "proposal.md missing Success Criteria section — cannot derive expected results for Journeys" |

If either mandatory field is missing, abort immediately and output the diagnostic message. Do NOT attempt to generate Journeys.

**Proposal Mode quality degradation**: When `## Key Scenarios` section (or any heading matching `key scenarios`, case-insensitive, with or without `##` / `###` prefix) is missing from proposal.md, generate smoke-level Journeys (happy path only) and annotate each Journey file with `quality: low` in frontmatter. Include a warning in the generated Journey:

```
> **Quality Notice**: This Journey was generated without Key Scenarios from the proposal.
> Only happy path is covered. Edge cases and invariants are inferred at minimal level.
```

### Optional inputs (both modes)

| Artifact | Usage |
|----------|-------|
| `docs/features/<slug>/prd/prd-ui-functions.md` | Map UI interactions to Journey steps |

<HARD-RULE>
Mode detection is automatic based on file existence. Do NOT ask the user which mode to use — if PRD files exist, use PRD Mode; if not, use Proposal Mode. If neither PRD files nor proposal.md exist, abort with an error message listing all missing files.
</HARD-RULE>

## When to Use

- User asks to "generate journeys" or "extract journeys"
- User provides `/gen-journeys` command
- After PRD is finalized, before gen-contracts
- Automated pipeline task (`T-test-gen-journeys`) generated by `forge task index`

## Pipeline Position

```
gen-journeys -> gen-contracts -> gen-test-scripts -> run-tests
     (narrative)     (technical)     (code gen)        (execution)
```

gen-journeys is the first step. It reads PRD user stories and produces Journey documents that gen-contracts consumes as input.

| Step | Input | Output | Reads Code |
|------|-------|--------|------------|
| gen-journeys | PRD user stories or proposal.md | Journey narrative documents | No |

## Process Flow

```
Read input sources (PRD or Proposal) -> Identify user workflows -> Classify Risk -> Generate per-Journey files -> Validate output -> Generate index
```

## Step 1: Read Input Sources

Read input documents based on the detected mode (see Prerequisites):

### PRD Mode

1. `docs/features/<slug>/prd/prd-user-stories.md` -- primary source: user stories with Given/When/Then acceptance criteria
2. `docs/features/<slug>/prd/prd-spec.md` -- functional specs, scope boundaries, user roles
3. `docs/proposals/<slug>/proposal.md` -- Key Scenarios section (if available, enriches Journey candidates)
4. `docs/features/<slug>/prd/prd-ui-functions.md` -- UI interaction details (if available)

### Proposal Mode

1. `docs/proposals/<slug>/proposal.md` -- sole source. Extract:
   - **Scope** section → defines feature boundaries and what is in/out of scope
   - **Success Criteria** section → defines verifiable outcomes (maps to Journey expected results)
   - **Key Scenarios** section (optional) → provides concrete user workflow descriptions (maps to Journey candidates)
   - **Requirements Analysis** section → provides additional context for edge cases and constraints

<HARD-RULE>
gen-journeys does NOT need code reconnaissance. Do not read source code, test files, or implementation files. The extraction is purely narrative, based on PRD content only.
</HARD-RULE>

## Step 2: Identify User Workflows

From the input sources, identify distinct user workflows. Each workflow becomes one Journey.

**Extraction rules (PRD Mode)**:

1. **One user story = one Journey candidate**. Each `Given/When/Then` block in the user stories file describes a cohesive user workflow.
2. **Merge related stories**. If two stories describe sequential steps toward the same goal (e.g., "claim a task" and "submit a task"), merge them into a single Journey.
3. **Preserve PRD references**. Every Journey must trace back to specific PRD user story IDs (e.g., "Story 1", "Story 2").

**Extraction rules (Proposal Mode)**:

1. **One Key Scenario = one Journey candidate**. Each scenario in the Key Scenarios section describes a distinct user workflow.
2. **Derive from Success Criteria when Key Scenarios are absent**. Parse each success criterion as a distinct goal-oriented workflow. Generate smoke-level Journeys (happy path only, `quality: low`).
3. **Derive from Scope when both Key Scenarios and Success Criteria are high-level**. Extract each "In Scope" bullet as a Journey candidate.
4. **Preserve proposal references**. Every Journey must trace back to specific proposal sections (e.g., "Key Scenario 1", "Success Criterion 3").

**Structuring each Journey**:

For each identified workflow, extract:

| Element | Source (PRD Mode) | Source (Proposal Mode) | Description |
|---------|-------------------|------------------------|-------------|
| Journey name | User story title or synthesized from actions | Key Scenario title or synthesized from Scope item | kebab-case identifier (e.g., `task-lifecycle`) |
| Happy path steps | `When` clauses from Given/When/Then | Key Scenario narrative or inferred from Success Criteria | Sequential user actions that achieve the goal |
| Edge cases | `Given` clauses that describe error/alternative states | Constraints/Dependencies section, or inferred from Risk section | Variations where preconditions differ from happy path |
| Expected results | `Then` clauses from Given/When/Then | Success Criteria mapped to the workflow | System responses for each step |
| Setup preconditions | `Given` clauses from the first story | Scope boundaries and constraints | Environment state needed before the Journey starts |

## Step 3: Classify Risk

Assign a Risk level to each Journey based on the workflow's characteristics:

| Risk | Criteria | Edge Case Density |
|------|----------|-------------------|
| **High** | Workflow involves state mutation, data loss risk, or irreversible operations | Edge case count MUST be >= happy path step count |
| **Medium** | Workflow involves multi-step interaction without irreversible side effects | Edge cases for each step with branching preconditions |
| **Low** | Workflow is read-only or purely observational | Happy path + critical error paths only |

**Risk inference from PRD**:
- Look for keywords in acceptance criteria: "create", "delete", "update", "modify" -> signals state mutation -> High
- Look for "view", "list", "search", "filter" -> signals read-only -> Low
- Multi-step workflows with validation but no destructive operations -> Medium
- When the PRD explicitly mentions severity or failure impact, use that as the primary signal

<HARD-RULE>
High-risk Journeys MUST have edge case count >= happy path step count. If extracting from the PRD yields fewer edge cases than happy path steps for a High-risk Journey, generate additional edge cases by considering: (1) invalid inputs, (2) missing preconditions, (3) concurrent/overlapping operations, (4) boundary conditions.
</HARD-RULE>

## Step 4: Generate Per-Journey Files

For each Journey, generate a directory and Markdown file using `templates/journey.md`.

**Output location**: `docs/features/<slug>/testing/<journey-name>/journey.md`

Create the directory `docs/features/<slug>/testing/<journey-name>/` if it does not exist.

**Proposal Mode quality annotation**: When generating smoke-level Journeys (Key Scenarios absent from proposal.md), add `quality: low` to the Journey file's frontmatter.

**Surface type annotation**: Every Journey file MUST include a `surface_types` field in its frontmatter listing the surface types this Journey covers (e.g., `[web, api]`). This field is derived from the detected surface types and the workflow's scope:

```yaml
---
feature: "{{FEATURE_SLUG}}"
journey: "{{JOURNEY_NAME}}"
risk_level: "{{RISK_LEVEL}}"
surface_types: ["web", "api"]
surface_keys: ["frontend", "backend"]
sources:
  - docs/features/{{FEATURE_SLUG}}/prd/prd-user-stories.md
  - docs/features/{{FEATURE_SLUG}}/prd/prd-spec.md
generated: "{{DATE}}"
---
```

**Determining surface_types and surface_keys per Journey**:
- If the Journey covers a cross-surface workflow (e.g., user interacts via web frontend and backend API processes the request), list all involved surface types and their corresponding keys
- If the Journey is scoped to a single surface's interaction, list only that surface's type and key
- When uncertain, err on the side of listing more surface types rather than fewer — downstream skills will generate surface-specific Contracts regardless
- `surface_keys` values come directly from `forge surfaces` output (the key portion of `key=type` pairs)

**One Journey = one directory**. Output is organized by Journey (user workflow), NOT by interface type (CLI, API, TUI, etc.).

<HARD-RULE>
gen-journeys output must be in a format that gen-contracts can directly consume. Each Journey file must contain:
- Journey name (in heading and frontmatter)
- Risk level (in frontmatter and body)
- Surface types (in frontmatter `surface_types` field) — list of covered surface types
- Surface keys (in frontmatter `surface_keys` field) — list of covered surface keys (user-defined identifiers from `forge surfaces`)
- Step sequence with numbered steps, each containing:
  - User action (what the user does)
  - Expected result (what the system produces)
- Edge cases referencing happy path steps with divergent preconditions
- Journey Invariants (cross-step properties)

This structure directly maps to gen-contracts' input expectations (see Core Concepts above).
</HARD-RULE>

### Batch Processing

Generate one Journey at a time. If a single Journey's content (happy path + all edge cases) exceeds the context window:

1. **Auto-batch trigger**: When estimated content exceeds ~50k tokens or step count exceeds 15
2. **Batch strategy**: Split within the Journey -- happy path steps as batch 1, edge case groups as subsequent batches
3. **Merge**: Combine all batches into a single Journey file per the template

Do NOT split across Journeys -- each Journey is always a single cohesive document.

## Step 5: Validate Output

After generating all Journey files, validate each one:

| Check | Rule |
|-------|------|
| Name present | Journey has a non-empty name in heading and frontmatter |
| Risk level valid | Risk is one of: High, Medium, Low |
| Surface types present | Journey has a non-empty `surface_types` array in frontmatter, listing valid surface type strings |
| Surface keys present | Journey has a non-empty `surface_keys` array in frontmatter, listing valid surface key strings matching `forge surfaces` output |
| Happy path steps | At least 1 happy path step exists |
| Edge case steps | At least 1 edge case exists |
| High-risk density | For High-risk Journeys: edge case count >= happy path step count |
| Invariants | At least 1 Journey Invariant declared |
| User action | Every step has a User Action description |
| Expected result | Every step has an Expected Result description |
| PRD traceability | Every Journey traces back to specific PRD user story IDs or proposal section references |
| Surface coverage complete | The union of all Journeys' `surface_types` includes every configured surface type. No surface may be left uncovered |
| Surface key coverage complete | The union of all Journeys' `surface_keys` includes every configured surface key. No surface key may be left uncovered |

If validation fails, fix the Journey file before proceeding.

## Step 6: Review & Commit

<HARD-RULE>
Do NOT commit documents automatically. Present all generated Journey files to the user for review and wait for explicit approval before committing.
</HARD-RULE>

<HARD-RULE>
AUTO_COMMIT mode MUST still execute Step 5 validation before committing. Skipping user review does NOT skip output validation.
</HARD-RULE>

### Interactive Mode (default)

When the skill is invoked manually via `/gen-journeys` (no `AUTO_COMMIT=true` in task context):

1. Present all generated Journey files to the user
2. Wait for the user to review and approve (or request changes)
3. Only commit after explicit user approval:

```bash
git add docs/features/<slug>/testing/
git commit -m "docs: generate journey <journey-name> for <feature-slug>"
```

### AUTO_COMMIT Mode (non-interactive)

When the task context contains `AUTO_COMMIT=true` (injected by automated task templates):

1. Skip user review — do NOT wait for approval
2. Ensure Step 5 validation has passed (this is mandatory, not skippable)
3. Proceed directly to commit:

```bash
git add docs/features/<slug>/testing/
git commit -m "docs: generate journey <journey-name> for <feature-slug>"
```

AUTO_COMMIT mode is intended for automated pipeline execution where human review is deferred to downstream eval stages (e.g., `eval-journey`).
