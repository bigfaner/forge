---
name: journey-contract-model
description: Core concepts of the Journey-Contract test model, directory conventions, tag-based promotion, and migration guide from old test model
---

# Journey-Contract Test Model

The Forge test pipeline organizes around user workflows (Journeys), defines expected behavior through six-dimension Contract declarations, and manages test lifecycles via Tag-Based Promotion.

## Core Concepts

### Journey

A Journey describes a real user workflow for achieving a goal. It is the primary organizational unit for testing — one Journey corresponds to one coherent user workflow.

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

## Semantic Descriptors

All dimension values use semantic descriptors -- natural language descriptions of expected behavior that express business intent rather than precise matching patterns.

See `rules/dimension-rules.md` "Semantic Descriptors" section for the full rules (regex prohibition, good/bad examples) and the HARD-RULE.

## Contract File Format

Each Contract is stored as a structured Markdown file that is both human-readable and machine-parseable.

### Template

```markdown
# Contract: <journey-name> / Step <N>: <step-description>

## Outcome "<outcome-name>"
- Preconditions: "<semantic description>"
- Input: <semantic description>
- Output: <semantic description>
- State: <semantic description>
- Side-effect: <semantic description or "none">
- Invariants: <step-level invariants or omit>

## Outcome "<outcome-name-2>"
- Preconditions: "<semantic description>"
- Input: <semantic description>
- Output: <semantic description>
- State: <semantic description>

## Journey Invariants
- <invariant description 1>
- <invariant description 2>
```

### Parseable Structure Rules

1. **Journey name**: Extracted from the file path: `docs/features/<slug>/testing/<journey>/contracts/step-N-*.md`
2. **Step sequence**: Extracted from the filename: `step-<N>-<slug>.md`
3. **Outcome sections**: `## Outcome "<name>"` headings declare new Outcome blocks
4. **Dimension format**: Each line follows `- <DimensionName>: <value>`
5. **Journey Invariants**: The `## Journey Invariants` section MUST appear exactly once in each Contract file

## Directory Convention

```
docs/features/<slug>/testing/
  <journey-name>/                     # Journey directory (kebab-case)
    journey.md                        # Journey narrative document
    contracts/                        # Contract specification directory
      step-1-<action-slug>.md         # Contract for Step 1
      step-2-<action-slug>.md         # Contract for Step 2
      step-N-<action-slug>.md

tests/
  <journey-name>/                     # Generated test files
    <test-file-1>
    <test-file-2>
```

### Rules

1. **Journey directory**: `docs/features/<slug>/testing/<journey>/` is named after the user workflow (kebab-case), containing journey.md and contracts/
2. **Contract directory**: `testing/<journey>/contracts/` contains Contract specification files for each Step, named `step-<N>-<action-slug>.md`
3. **Test files**: Generated directly into `tests/<journey>/` by gen-test-scripts, following project test framework naming conventions
4. **No staging area**: Tests are generated directly to their final location, without intermediate directories

## Tag-Based Promotion

Tests manage their lifecycle via tags rather than file movement:

| Stage | Tag | Action |
|-------|-----|--------|
| New | `@feature` | Automatically injected into newly generated tests |
| Promoted | `@regression` | Automatically upgraded by `/run-tests` |
| CI selection | -- | Use test framework's native tag filter (e.g., `go test -tags=regression`, `pytest -m regression`) |

Tags use the test framework's native mechanism (Go build tags, pytest markers, describe groups, etc.).

## Migration Guide: Old Test Model -> Journey-Contract Model

The old model classified tests by interface type, used single-step TC format, and relied on a staging+graduation lifecycle. This guide provides the complete mapping and steps for migrating to the Journey-Contract model.

### Concept Mapping

| Old Model Concept | New Model Concept | Change Description |
|-------------------|-------------------|-------------------|
| Classification by interface type (CLI/API/TUI/Web/Mobile) | Organization by Journey (user workflow) | Organization axis shifted from "technical interface" to "user scenario" |
| Single-step TC (one command / one endpoint) | Step + Contract (one step within a workflow) | TC becomes a Step within a Journey, preserving complete workflow context |
| TC Steps (freeform operation list) | Contract six-dimension declarations (structured verification) | Upgraded from unstructured descriptions to six-dimension specifications |
| TC Expected (freeform text expectations) | Outcome (Preconditions-mutually-exclusive multi-result) | Supports multiple results per step (success, failure, edge cases) |
| Type field (CLI/API/TUI/Web/Mobile) | Interface-specific descriptions within the Journey dimension | Type information is embedded in Contract dimensions rather than top-level classification |
| Setup declarations (API data preparation) | Contract Preconditions + Fact Table | Data preparation expressed as Preconditions constraints |
| `{stepN.field}` references | Semantic descriptors within Contracts + precise matching via gen-test-scripts | Declaration phase uses natural language; code generation phase uses Fact Table for precise matching |
| staging directory + graduation process | Tag-Based Promotion (`@feature` -> `@regression`) | Tag-driven rather than file movement |
| 6 hardcoded language profiles | Convention-driven (`docs/conventions/testing/<surface>/core.md`) | Extensible, editable Convention files |

### Directory Restructuring

#### Old Structure

```
tests/e2e/
  features/
    <feature-name>/
      *_test.go            # Tests grouped by feature
  *_test.go                # Top-level scattered tests
```

#### New Structure

```
docs/features/<slug>/testing/
  <journey-name>/               # Grouped by user workflow
    journey.md
    contracts/
      step-1-<action>.md
      step-N-<action>.md

tests/
  <journey-name>/               # Generated test files
    <test-files>
```

### Tag-Based Promotion

| Process | Description |
|---------|-------------|
| Tests generated directly to `tests/<journey>/` | No staging directory |
| Automatic `@feature` tag | No explicit staging required |
| `/run-tests` (tag promotion) | Tag changes from `@feature` to `@regression` |
| Files always in their final location | No file movement |
| Uses test framework's native tag filtering mechanism (e.g., `go test -tags=regression`) | Select by tag |
