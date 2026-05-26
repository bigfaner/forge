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

## Surface Detection

Load: `rules/journey-contract-model.md` — core concepts (Journey, Step, Contract, Outcome), directory conventions, and tag-based promotion model used by downstream skills.

Before processing PRD sources, determine the project's configured surface types via the `forge surfaces` CLI command. Surface determines testing strategy, required Outcomes, and test level emphasis.

### Detection via CLI

Query the project's configured surface types using `forge surfaces <path>`:

```bash
forge surfaces <path>
```

**Exit code contract**:

| Exit Code | Meaning | Action |
|-----------|---------|--------|
| 0 | Surface types found. stdout contains one surface type per line (e.g., `web`, `api`). For monorepo configs, each configured surface appears on its own line. | Parse stdout to collect all surface type strings. Proceed to rule loading. |
| 1 | No surface configured for the given path. stderr contains an error message with configuration guidance. | **Pause pipeline**. Show the stderr message to the user and ask them to configure surfaces via `forge init`. |

**Examples**:

```bash
# Single-surface project (scalar form): any path returns the same type
forge surfaces .
# stdout: api  (exit 0)

# Multi-surface project: all configured surfaces are listed
forge surfaces .
# stdout:
# web
# api
# (exit 0)

# Monorepo with path-level surfaces: query specific path
forge surfaces frontend/src
# stdout: web  (exit 0)

# Path not configured
forge surfaces unknown-dir
# stderr: Error: no surface found for path "unknown-dir". Run `forge init` to configure surfaces.  (exit 1)
```

### Detection Flow

1. Run `forge surfaces .` (or the relevant source path for the feature being tested)
2. If exit code is 0: parse stdout to collect all surface type strings (one per line)
3. If exit code is 1: pause the pipeline and ask the user to configure surfaces
4. Load the corresponding rule files for each detected surface type

**Supported surface types**: `web`, `api`, `cli`, `tui`, `mobile`

**Surface rule files** (loaded dynamically based on detected types):
- `rules/surface-api.md`
- `rules/surface-cli.md`
- `rules/surface-mobile.md`
- `rules/surface-tui.md`
- `rules/surface-web.md`

<HARD-RULE>
Surface detection must complete before Journey generation begins. If the `forge surfaces` command returns exit code 1, the pipeline must pause and wait for user input. Never proceed with a guessed surface type. Do NOT scan project files independently for surface detection -- always use `forge surfaces <path>`.
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

#### API Surface

- Mandatory error outcomes: HTTP status code boundaries (4xx client errors, 5xx server errors)
- Test level emphasis: integration-heavy ratio (test strategy guidance from `rules/surface-api.md`)
- Edge case focus: authentication failures, rate limiting, payload validation

#### Web Surface

- Mandatory outcomes: page load states, navigation transitions, form validation feedback
- Test level emphasis: e2e-heavy ratio (test strategy guidance from `rules/surface-web.md`)
- Edge case focus: browser compatibility, responsive layout breaks, client-side validation

#### CLI Surface

- Mandatory outcomes: exit codes, stdout/stderr output, signal handling
- Test level emphasis: unit-heavy ratio (test strategy guidance from `rules/surface-cli.md`)
- Edge case focus: invalid flags, missing arguments, pipe failures

#### TUI Surface

- Mandatory outcomes: rendering states, keyboard navigation, screen transitions
- Test level emphasis: integration-heavy ratio (test strategy guidance from `rules/surface-tui.md`)
- Edge case focus: terminal resize, key binding conflicts, rendering artifacts

#### Mobile Surface

- Mandatory outcomes: screen transitions, gesture responses, offline/online state
- Test level emphasis: e2e-heavy ratio (test strategy guidance from `rules/surface-mobile.md`)
- Edge case focus: network interruptions, app lifecycle events, permission denials

### Journey Surface Coverage

Each Journey must declare which surface types it covers in its frontmatter `surface_types` field. This enables downstream skills (gen-contracts, gen-test-scripts) to generate surface-appropriate Contracts and test scripts.

**Coverage rules**:
- A Journey covering a cross-surface user workflow (e.g., "user submits form via web, backend API processes it") must list all involved surface types
- A Journey limited to a single surface's interaction lists only that surface type
- The union of all Journeys' `surface_types` must cover every configured surface type — no surface may be left uncovered

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

**Proposal Mode quality degradation**: When `## Key Scenarios` section is missing from proposal.md, generate smoke-level Journeys (happy path only) and annotate each Journey file with `quality: low` in frontmatter. Include a warning in the generated Journey:

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
sources:
  - docs/features/{{FEATURE_SLUG}}/prd/prd-user-stories.md
  - docs/features/{{FEATURE_SLUG}}/prd/prd-spec.md
generated: "{{DATE}}"
---
```

**Determining surface_types per Journey**:
- If the Journey covers a cross-surface workflow (e.g., user interacts via web frontend and backend API processes the request), list all involved surface types
- If the Journey is scoped to a single surface's interaction, list only that surface type
- When uncertain, err on the side of listing more surface types rather than fewer — downstream skills will generate surface-specific Contracts regardless

**One Journey = one directory**. Output is organized by Journey (user workflow), NOT by interface type (CLI, API, TUI, etc.).

<HARD-RULE>
gen-journeys output must be in a format that gen-contracts can directly consume. Each Journey file must contain:
- Journey name (in heading and frontmatter)
- Risk level (in frontmatter and body)
- Surface types (in frontmatter `surface_types` field) — list of covered surface types
- Step sequence with numbered steps, each containing:
  - User action (what the user does)
  - Expected result (what the system produces)
- Edge cases referencing happy path steps with divergent preconditions
- Journey Invariants (cross-step properties)

This structure directly maps to gen-contracts' input expectations (Section 5.1 of model-and-directory-spec.md).
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
| Happy path steps | At least 1 happy path step exists |
| Edge case steps | At least 1 edge case exists |
| High-risk density | For High-risk Journeys: edge case count >= happy path step count |
| Invariants | At least 1 Journey Invariant declared |
| User action | Every step has a User Action description |
| Expected result | Every step has an Expected Result description |
| PRD traceability | Every Journey traces back to specific PRD user story IDs or proposal section references |
| Surface coverage complete | The union of all Journeys' `surface_types` includes every configured surface type. No surface may be left uncovered |

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

## Related Skills

| Skill | Usage |
|-------|-------|
| `/write-prd` | Create PRD with user stories (input source) |
| `/gen-contracts` | Consume Journey documents to generate Contract specifications |
| `/gen-test-scripts` | Generate executable test code from Contracts |

## Reference

The authoritative model definition is at `docs/features/<slug>/design/model-and-directory-spec.md` (if it exists in the project). Key concepts used by this skill:

- **Journey**: User's real workflow to accomplish a goal (Section 1.1)
- **Step**: Single user action within a Journey (Section 1.2)
- **Risk Classification**: High/Medium/Low severity guiding test density (Section 1.1)
- **Journey Invariants**: Cross-step properties that must hold throughout (Section 1.3)
- **Semantic Descriptors**: Natural-language descriptions used in Contracts (Section 1.5)
