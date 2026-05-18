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

## Prerequisites

| Artifact | Missing prompt |
|----------|----------------|
| `docs/features/<slug>/prd/prd-user-stories.md` | Run `/write-prd` first |
| `docs/features/<slug>/prd/prd-spec.md` | Run `/write-prd` first |

Optional inputs (used for richer context when available):

| Artifact | Usage |
|----------|-------|
| `docs/proposals/<slug>/proposal.md` | Extract Key Scenarios as Journey candidates |
| `docs/features/<slug>/prd/prd-ui-functions.md` | Map UI interactions to Journey steps |

## When to Use

- User asks to "generate journeys" or "extract journeys"
- User provides `/gen-journeys` command
- After PRD is finalized, before gen-contracts

## Pipeline Position

```
gen-journeys -> gen-contracts -> gen-test-scripts -> run-tests
     (narrative)     (technical)     (code gen)        (execution)
```

gen-journeys is the first step. It reads PRD user stories and produces Journey documents that gen-contracts consumes as input.

| Step | Input | Output | Reads Code |
|------|-------|--------|------------|
| gen-journeys | PRD user stories | Journey narrative documents | No |

## Process Flow

```
Read PRD sources -> Identify user workflows -> Classify Risk -> Generate per-Journey files -> Validate output -> Generate index
```

## Step 1: Read PRD Sources

Read all available PRD documents to understand user workflows:

1. `docs/features/<slug>/prd/prd-user-stories.md` -- primary source: user stories with Given/When/Then acceptance criteria
2. `docs/features/<slug>/prd/prd-spec.md` -- functional specs, scope boundaries, user roles
3. `docs/proposals/<slug>/proposal.md` -- Key Scenarios section (if available, enriches Journey candidates)
4. `docs/features/<slug>/prd/prd-ui-functions.md` -- UI interaction details (if available)

<HARD-RULE>
gen-journeys does NOT need code reconnaissance. Do not read source code, test files, or implementation files. The extraction is purely narrative, based on PRD content only.
</HARD-RULE>

## Step 2: Identify User Workflows

From the PRD sources, identify distinct user workflows. Each workflow becomes one Journey.

**Extraction rules**:

1. **One user story = one Journey candidate**. Each `Given/When/Then` block in the user stories file describes a cohesive user workflow.
2. **Merge related stories**. If two stories describe sequential steps toward the same goal (e.g., "claim a task" and "submit a task"), merge them into a single Journey.
3. **Preserve PRD references**. Every Journey must trace back to specific PRD user story IDs (e.g., "Story 1", "Story 2").

**Structuring each Journey**:

For each identified workflow, extract:

| Element | Source | Description |
|---------|--------|-------------|
| Journey name | User story title or synthesized from actions | kebab-case identifier (e.g., `task-lifecycle`) |
| Happy path steps | `When` clauses from Given/When/Then | Sequential user actions that achieve the goal |
| Edge cases | `Given` clauses that describe error/alternative states | Variations where preconditions differ from happy path |
| Expected results | `Then` clauses from Given/When/Then | System responses for each step |
| Setup preconditions | `Given` clauses from the first story | Environment state needed before the Journey starts |

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

For each Journey, generate a Markdown file using `${CLAUDE_SKILL_DIR}/templates/journey.md`.

**Output location**: `docs/features/<slug>/testing/journeys/<journey-name>.md`

**One Journey = one file**. Output is organized by Journey (user workflow), NOT by interface type (CLI, API, TUI, etc.).

<HARD-RULE>
gen-journeys output must be in a format that gen-contracts can directly consume. Each Journey file must contain:
- Journey name (in heading and frontmatter)
- Risk level (in frontmatter and body)
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
| Happy path steps | At least 1 happy path step exists |
| Edge case steps | At least 1 edge case exists |
| High-risk density | For High-risk Journeys: edge case count >= happy path step count |
| Invariants | At least 1 Journey Invariant declared |
| User action | Every step has a User Action description |
| Expected result | Every step has an Expected Result description |
| PRD traceability | Every Journey traces back to specific PRD user story IDs |

If validation fails, fix the Journey file before proceeding.

## Step 6: Generate Index

After all Journey files are written, generate a manifest index at `docs/features/<slug>/testing/journeys/manifest.md`:

```yaml
---
feature: "{{FEATURE_SLUG}}"
generated: "{{DATE}}"
journey-count: {{COUNT}}
---
```

Include a **Journeys Summary** table:

| Journey | Risk | Happy Path Steps | Edge Cases | Source Stories | File |
|---------|------|-----------------|------------|---------------|------|

This index allows gen-contracts to discover all Journey files.

## Step 7: Review & Commit

<HARD-RULE>
Do NOT commit documents automatically. Present all generated Journey files to the user for review and wait for explicit approval before committing.
</HARD-RULE>

1. Present all generated Journey files and the manifest to the user
2. Wait for the user to review and approve (or request changes)
3. Only commit after explicit user approval:

```bash
git add docs/features/<slug>/testing/journeys/
git commit -m "docs: generate journeys for <feature-slug>"
```

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
