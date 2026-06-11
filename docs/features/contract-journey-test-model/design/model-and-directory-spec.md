---
title: "Journey-Driven Test Model & Directory Specification"
status: canonical
serves:
  - gen-journeys
  - gen-contracts
  - gen-test-scripts
  - run-tests
---

# Journey-Driven Test Model & Directory Specification

This document defines the core concepts, data model, and directory conventions for the Journey-Driven testing pipeline. It is the authoritative reference for all pipeline skills (gen-journeys, gen-contracts, gen-test-scripts, run-tests).

## 1. Core Concepts

### 1.1 Journey

A **Journey** describes a user's real workflow to accomplish a goal. It is the primary organizational unit for tests -- one Journey maps to one cohesive user workflow.

**Structure**:

```
Journey "<name>"
  Risk: High | Medium | Low
  Steps:
    Step N: <user action> -> <expected outcome>
    Step N+1: ...
  Edge Cases:
    Step N variant: <alternative action> -> <alternative outcome>
```

**Properties**:

| Property | Description | Values |
|----------|-------------|--------|
| Name | Human-readable identifier, kebab-case (e.g., `task-lifecycle`) | String |
| Risk | Severity classification guiding test density | `High`, `Medium`, `Low` |
| Steps | Ordered sequence of user actions with expected outcomes | List of Step |
| Invariants | Cross-step properties that must hold throughout the Journey | List of invariant declarations |

**Risk Classification Criteria**:

| Risk | Criteria | Test Density Expectation |
|------|----------|--------------------------|
| High | Workflow involves state mutation, data loss risk, or irreversible operations | Edge case count >= happy path step count |
| Medium | Workflow involves multi-step interaction without irreversible side effects | Edge cases for each step with branching preconditions |
| Low | Workflow is read-only or purely observational | Happy path + critical error paths only |

**Journey Isolation**: Each Journey executes in an independent temporary working directory to prevent cross-Journey interference during parallel execution. The temp directory path includes the Journey name and a random suffix.

### 1.2 Step

A **Step** is a single user action within a Journey. Each Step maps to one Contract with one or more Outcomes.

**Properties**:

| Property | Description |
|----------|-------------|
| Sequence number | 1-based ordinal within the Journey |
| User action | What the user does (e.g., runs a command, clicks a button, sends a request) |
| Expected outcomes | One or more Outcome declarations, each with distinct Preconditions |

### 1.3 Contract

A **Contract** is the verification mechanism for a Journey Step. It defines the system's expected behavior across six dimensions, all declared at the Outcome level (with Invariants additionally declared at the Journey level).

#### 1.3.1 Six Dimensions

**Mandatory dimensions** (every Outcome MUST declare these four):

| Dimension | Description | CLI Example | API Example | TUI Example |
|-----------|-------------|-------------|-------------|-------------|
| Preconditions | State that must hold before the step executes | `task status is in_progress` | `user is authenticated` | `Model state is "idle"` |
| Input | What goes into the system | Command arguments + working directory state | Request schema + auth headers | Model state + event |
| Output | What comes out of the system | stdout/stderr semantic description + exit code | Response schema + status code | Model new state + View output |
| State | System state changes caused by the step | File system changes (declared as semantic descriptions) | Database/API state changes | In-memory Model field changes |

**Optional dimensions** (gen-contracts MAY omit these; omission = no constraint):

| Dimension | Description | CLI Example | API Example | TUI Example |
|-----------|-------------|-------------|-------------|-------------|
| Side-effect | External side effects triggered by the step | Hook trigger, external process invocation | Network side effects, message dispatch | Async Cmd execution |
| Invariants (step-level) | Properties that hold within a single step | `index.json remains valid JSON throughout` | Response time under threshold | No unhandled errors in Model update |

**Journey-level Invariants** (mandatory): Declared once per Journey, not per Step. These are cross-step properties that must hold from the first Step to the last.

```
## Journey Invariants
- feature_slug consistent across all steps
- task_id stable once assigned
```

#### 1.3.2 Multi-Outcome Contracts

Each Step supports multiple Outcomes. Outcomes are mutually exclusive -- at any given moment, only one Outcome's Preconditions can be satisfied. This prevents combinatorial explosion.

**Structure**:

```
Step N: <action>
  Outcome "success":
    Preconditions: "<precondition A>"
    Input: <input A>
    Output: <output A>
    State: <state change A>
    Side-effect: <side-effect A>  (optional)

  Outcome "<alternative>":
    Preconditions: "<precondition B>"  <-- mutually exclusive with A
    Input: <input B>
    Output: <output B>
    State: <state change B>
```

**Rules**:
- Each Outcome has independent Preconditions, Input, Output, State, and optional Side-effect/Invariants
- Preconditions across Outcomes MUST be mutually exclusive for a given system state
- Steps with more than 5 Outcomes trigger a review checkpoint (merge semantically similar Outcomes)

### 1.4 Outcome

An **Outcome** is a complete set of Contract dimension declarations for a specific scenario (success, error variant, edge case). Outcomes within a Step are distinguished by their Preconditions.

**Properties**:

| Property | Required | Description |
|----------|----------|-------------|
| Name | Yes | Descriptive label (e.g., `success`, `not-in-progress`, `no-tasks-available`) |
| Preconditions | Yes | State required for this Outcome to be the active one |
| Input | Yes | What the user provides to the system |
| Output | Yes | What the system produces (semantic descriptor) |
| State | Yes | How system state changes |
| Side-effect | No | External effects (default: `none`) |
| Invariants | No | Step-level invariants (default: none) |

### 1.5 Semantic Descriptors

**Semantic descriptors** are natural-language descriptions of expected values, used in the gen-contracts stage. They express *business intent* rather than precise matching patterns.

**Purpose**: gen-contracts operates without direct access to runtime output formats. By using semantic descriptors, it focuses on *what* the output should convey, not *how* it is formatted. The precise matching (regex, JSON schema, etc.) is deferred to gen-test-scripts, which has access to the Fact Table from code reconnaissance.

**Descriptor format**: Free-form natural language enclosed in quotes.

| Stage | Uses | Produces |
|-------|------|----------|
| gen-contracts | Business intent + code reconnaissance context | Semantic descriptors in Contract files |
| gen-test-scripts | Contract files + Fact Table from code reconnaissance | Precise matchers (regex, JSON schema assertions, etc.) |

**Conversion pipeline example**:

```
gen-contracts output:
  Output: "success confirmation containing feature-slug"

gen-test-scripts processing:
  1. Query Fact Table for the command's actual stdout sample
  2. Find: "Feature my-feature created successfully"
  3. Generate regex: /Feature\s+([\w-]+)\s+created successfully/
  4. Declare capture group: feature_slug: $1
```

**Constraint**: Semantic descriptors MUST NOT contain regex syntax or framework-specific assertion patterns. They are pure natural language.

## 2. Contract Specification File Format

Each Contract is stored as a structured Markdown file. This format is designed to be both human-readable and machine-parseable by gen-contracts and gen-test-scripts.

### 2.1 File Template

```markdown
# Contract: <journey-name> / Step <N>: <step-description>

## Outcome "<outcome-name>"
- Preconditions: "<semantic description>"
- Input: <semantic description of command/args/flags or request or event>
- Output: <semantic description of stdout/stderr + exit code, or response + status, or model state>
- State: <semantic description of state changes>
- Side-effect: <semantic description or "none">
- Invariants: <step-level invariants or omit>

## Outcome "<outcome-name-2>"
- Preconditions: "<semantic description>"
- Input: <semantic description>
- Output: <semantic description>
- State: <semantic description>
```

### 2.2 Journey Invariants Section

Every Contract file ends with a Journey Invariants section:

```markdown
## Journey Invariants
- <invariant description 1>
- <invariant description 2>
```

Journey-level Invariants are mandatory and appear in every Contract file within the Journey. They are identical across all Step Contract files in the same Journey.

### 2.3 State Verification Levels

When a project does not expose a state query interface, the State dimension degrades gracefully:

| Level | Meaning | Declaration |
|-------|---------|-------------|
| `full` | All state fields can be independently verified | Default |
| `partial` | State fields inferred from Output only | `state-verification: partial` |
| `deferred` | Some state fields cannot be inferred from Output | `state-verification: deferred` + `limitations` section |

```markdown
## Outcome "success"
- Preconditions: "feature exists with slug matching arg"
- Input: feature-slug as positional arg
- Output: "success confirmation containing feature-slug", exit code 0
- State: feature directory created with manifest.md and tasks/index.json
- Side-effect: none

<!-- state-verification: partial -->
<!-- State fields inferred from Output: feature-slug in stdout -->
<!-- State fields deferred: internal task index ordering -->
```

## 3. Directory Convention

### 3.1 Structure

```
tests/
  <journey-name>/                     <-- Journey directory (domain-oriented, stable)
    _contracts/                       <-- Contract specification directory
      step-1-<action-slug>.md         <-- Contract for Step 1
      step-2-<action-slug>.md         <-- Contract for Step 2
      step-N-<action-slug>.md         <-- Contract for Step N
    <test-file-1>                     <-- Generated test files (directly in final location)
    <test-file-2>
    ...
```

### 3.2 Rules

1. **Journey directory** (`tests/<journey-name>/`): Named after the user workflow (kebab-case). This is the permanent location for all test files related to this Journey. Tests are generated directly here -- no staging area.

2. **Contract directory** (`tests/<journey-name>/_contracts/`): Contains Contract specification files, one per Step. Files are named `step-<N>-<action-slug>.md` where `<N>` is the 1-based step ordinal and `<action-slug>` is a kebab-case summary of the user action.

3. **Test files**: Generated directly into the Journey directory by gen-test-scripts. File naming follows the project's test framework conventions. No intermediate directories or staging areas.

4. **Lifecycle management**: Tests are managed via tags, not file movement:
   - New tests: injected with `@feature` tag
   - Promoted tests: `@feature` -> `@regression` via `forge test promote <journey>`
   - CI selection: `forge test run --tags regression` or `--tags feature`

### 3.3 Example

```
tests/
  task-lifecycle/
    _contracts/
      step-1-feature-create.md
      step-2-task-claim.md
      step-3-task-submit.md
    claim_submit_test.go              <-- @feature (newly generated)
    task_record_test.go               <-- @regression (promoted)

  session-diagnostics/
    _contracts/
      step-1-open-session.md
      step-2-browse-call-tree.md
      step-3-expand-entry.md
      step-4-open-diagnosis-panel.md
    session_test.go                   <-- @feature
```

### 3.4 Tag Format by Test Framework

Tags are embedded using the test framework's native mechanism:

| Framework | Tag Syntax |
|-----------|-----------|
| Go testing | `//go:build feature` or `//go:build regression` |
| Python pytest | `@pytest.mark.feature` or `@pytest.mark.regression` |
| JavaScript (mocha/jest) | `describe("@feature", ...)` or `describe("@regression", ...)` |
| Playwright | `test.describe("@feature", ...)` or `test.describe("@regression", ...)` |
| Rust | `#[cfg(feature = "test-feature")]` or `#[cfg(feature = "test-regression")]` |

Projects declare their test framework in `.forge/config.yaml` (see Section 4), and gen-test-scripts uses the corresponding tag syntax.

## 4. Configuration Schema

The configuration-driven framework is declared in `.forge/config.yaml`. Forge never hardcodes language or framework names -- all framework selection is driven by configuration with built-in templates as defaults.

### 4.1 Schema

```yaml
# .forge/config.yaml

# Existing fields (backward compatible)
project-type: backend | frontend | fullstack
interfaces:
  - cli        # CLI commands
  - api        # HTTP endpoints
  - tui        # Terminal UI
  - web-ui     # Browser UI
  - mobile     # Mobile app
languages:
  - go | javascript | python | java | rust | ...

# New fields for Journey-Driven testing
test-framework: <framework-name>       # e.g., go-testing, pytest, mocha, playwright
test-command: <execution-command>       # e.g., "go test ./...", "pytest tests/", "npx playwright test"
capabilities:                          # Optional: declare project-specific test capabilities
  state-query: true | false            # Whether the project exposes state query interfaces
  auth-strategy: none | token | cookie | api-key  # Auth mechanism for API/Web-UI Journeys
  tui-await-timeout: <milliseconds>    # Default timeout for TUI async Cmd await (default: 3000)
```

### 4.2 Field Semantics

| Field | Required | Default | Description |
|-------|----------|---------|-------------|
| `test-framework` | No | Auto-detected from `languages` + project files | The test framework used for generating test code. When omitted, Forge auto-detects from project files. |
| `test-command` | No | Derived from `test-framework` | The command to execute tests. When omitted, derived from the framework's standard runner. |
| `capabilities` | No | All capabilities default to conservative values | Project-specific capabilities that affect Contract generation and test script behavior. |

### 4.3 Auto-Detection

When `test-framework` is not declared, Forge auto-detects from project files:

| Signal File | Detected Framework |
|-------------|-------------------|
| `go.mod` + `*_test.go` | `go-testing` |
| `package.json` + `@playwright/test` | `playwright` |
| `package.json` + `mocha` | `mocha` |
| `package.json` + `jest` | `jest` |
| `pytest.ini` / `pyproject.toml` with pytest | `pytest` |
| `Cargo.toml` + `tests/` | `rust-testing` |

Auto-detection is a convenience default. Explicit `test-framework` declaration always overrides auto-detection.

### 4.4 Framework-to-Tag Mapping

The framework determines how tags are embedded in generated test files. This mapping is defined in built-in templates and can be overridden by project configuration.

| Framework | Feature Tag | Regression Tag | CLI Filter |
|-----------|------------|----------------|------------|
| `go-testing` | `//go:build feature` | `//go:build regression` | `-tags feature` / `-tags regression` |
| `pytest` | `@pytest.mark.feature` | `@pytest.mark.regression` | `-m feature` / `-m regression` |
| `mocha` | `describe("@feature", ...)` | `describe("@regression", ...)` | `--grep @feature` / `--grep @regression` |
| `playwright` | `test.describe("@feature", ...)` | `test.describe("@regression", ...)` | `--grep @feature` / `--grep @regression` |
| `rust-testing` | `#[cfg(test_feature)]` | `#[cfg(test_regression)]` | `--cfg test_feature` / `--cfg test_regression` |

### 4.5 Backward Compatibility

The `languages` and `interfaces` fields in `.forge/config.yaml` continue to work as before. The new `test-framework` and `test-command` fields extend the configuration without breaking existing projects.

When both `languages` and `test-framework` are present:
- `test-framework` determines code generation templates and tag syntax
- `languages` is used for auto-detection fallback and project-level language configuration

When neither is present, auto-detection from project files applies.

## 5. gen-contracts Parseable Structure

The Contract specification files must be parseable by gen-contracts for validation and by gen-test-scripts for code generation. The following structural rules ensure machine-parseability:

### 5.1 Structural Rules

1. **Journey name**: Extracted from the Contract file path: `tests/<journey-name>/_contracts/step-N-*.md` -> Journey name = directory name.

2. **Step sequence**: Extracted from the filename: `step-<N>-<slug>.md` -> Step number = N.

3. **Outcome sections**: Each `## Outcome "<name>"` heading declares a new Outcome block. All dimension declarations for that Outcome follow until the next `## Outcome` heading or `## Journey Invariants` heading.

4. **Dimension format**: Each dimension is a single line starting with `- <DimensionName>:` followed by the semantic descriptor value.

5. **Journey Invariants**: The `## Journey Invariants` section contains a bulleted list of cross-step invariant declarations. This section MUST appear exactly once per Contract file.

### 5.2 Completeness Validation

gen-contracts validates each generated Contract against these rules:

| Check | Rule |
|-------|------|
| Mandatory dimensions | Each Outcome MUST have: Preconditions, Input, Output, State |
| Semantic descriptor purity | No dimension value may contain regex syntax (`\d`, `.*`, `[^...]`, etc.) |
| Outcome name uniqueness | Outcome names within a Step MUST be unique |
| Preconditions mutual exclusivity | Different Outcomes' Preconditions MUST be distinguishable (not identical) |
| Journey Invariants presence | Every Contract file MUST have a `## Journey Invariants` section with at least 1 entry |
| Side-effect default | When omitted, Side-effect defaults to `none` |
| Step-level Invariants default | When omitted, step-level Invariants default to no constraint |

### 5.3 Batch Processing

When a single Journey has more than 15 Contracts or the estimated token count exceeds 50k, gen-contracts automatically splits into multiple batches:

- **Batch 1**: Happy path Outcomes (all steps' success Outcomes)
- **Batch 2+**: Edge case Outcomes grouped by semantic similarity

Split batches are merged back into a single complete Contract document per Step. The merged result is structurally identical to a single-batch generation.

## 6. TUI Async Await Semantics

For TUI Journey Steps involving asynchronous operations, Contract specifications use `await` semantics:

- `await` = wait for all pending Cmds to complete (up to `tui-await-timeout` milliseconds from config, default 3000ms)
- Timeout behavior: fail-fast, report the name of the timed-out Cmd
- `tea.Batch(cmd1, cmd2)`: all concurrent Cmds must complete before proceeding to the next Step

**Contract declaration example**:

```markdown
## Outcome "diagnosis-loaded"
- Preconditions: "session loaded, call tree visible, entry expanded"
- Input: key "d" await 3000ms
- Output: "view contains diagnosis summary panel"
- State: Model.diagnosis_panel = visible
```

## 7. Pipeline Integration

### 7.1 Four-Step Pipeline

```
gen-journeys -> gen-contracts -> gen-test-scripts -> run-tests
     (narrative)     (technical)     (code gen)        (execution)
```

| Step | Input | Output | Reads Code |
|------|-------|--------|------------|
| gen-journeys | PRD user stories | Journey narrative documents (user workflows + Risk classification) | No |
| gen-contracts | Journey documents + code reconnaissance | Contract specifications (six dimensions, semantic descriptors) | Yes (Fact Table) |
| gen-test-scripts | Contract specifications + templates + code reconnaissance | Executable test code + Journey smoke test | Yes (Fact Table) |
| run-tests | Test code | Results report | No |

### 7.2 Test Layer Hierarchy

```
Unit (TDD)                          -- Developer hand-written, pure functions / state machines
    |
Contract (Journey Step validation)  -- Forge generated, CLI / API / TUI
    |
E2E (Web-UI / Mobile-UI)           -- Frontend projects only, full user experience
```

Projects with only CLI interfaces: Unit -> Contract (no E2E).
Projects with Web/Mobile UI: Unit -> Contract -> E2E.

### 7.3 Smoke Tests

Each Journey has exactly one smoke test that end-to-end executes the happy path. The smoke test:
- Runs all happy path Steps in sequence within the Journey's isolated temp directory
- Validates that runtime interactions (environment setup, inter-process communication, timing behavior) work correctly
- Is generated by gen-test-scripts, reusing the Journey's happy path Steps
- Retries up to 3 times on failure (5s interval); after 3 failures, marks as `@flaky` with a fix suggestion
