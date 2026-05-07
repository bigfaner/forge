# PRD Evaluation Rubric

**Total: 100 points**
**Report template:** `plugins/forge/skills/eval-prd/templates/report.md`

## Required Sections (prd-spec.md)

| Section | Required |
|---------|----------|
| Background (Reason/Target/Users) | ✓ |
| Goals + Quantified Metrics | ✓ |
| Scope (In + Out) | ✓ |
| Flow Description + Mermaid Diagram | ✓ |
| Functional Specs (reference to prd-ui-functions.md) | ✓ |

## Required Sections (prd-user-stories.md)

| Section | Required |
|---------|----------|
| User Stories | ✓ |
| Acceptance Criteria (Given/When/Then) | ✓ |

## Required Sections (prd-ui-functions.md)

> **Mandatory** when the feature has a UI surface. Skip for backend/API/CLI-only features.

| Section | Required |
|---------|----------|
| UI Functions with Placement | ✓ |
| User Interaction Flow | ✓ |
| Data Requirements | ✓ |
| Validation Rules | ✓ |

## Scoring Modes

### Mode A: Feature WITH UI (prd-ui-functions.md present)

| Dimension | Points |
|-----------|--------|
| 1. Background & Goals | 15 |
| 2. Flow Diagrams | 20 |
| 3. Functional Specs (evaluates prd-ui-functions.md) | 20 |
| 4. User Stories | 30 |
| 5. Scope Clarity | 15 |
| **Total** | **100** |

### Mode B: Feature WITHOUT UI (prd-ui-functions.md absent)

| Dimension | Points |
|-----------|--------|
| 1. Background & Goals | 15 |
| 2. Flow Diagrams | 20 |
| 3. Flow Completeness (evaluates prd-spec.md Flow Description) | 20 |
| 4. User Stories | 30 |
| 5. Scope Clarity | 15 |
| **Total** | **100** |

> Detection: if `prd-ui-functions.md` exists in the PRD directory → Mode A; otherwise → Mode B.

## Dimensions

### 1. Background & Goals (15 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Background has three elements (Reason/Target/Users) | 0-5 | Are all three present and specific? |
| Goals are quantified | 0-4 | Is there at least one numeric target (%, count, time)? |
| Background and goals are logically consistent | 0-6 | Does the goal follow from the stated problem? |

### 2. Flow Diagrams (20 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Mermaid diagram exists | 0-7 | Is there at least one Mermaid flowchart? |
| Main path complete (start → end) | 0-7 | Does the diagram cover the full happy path? |
| Decision points + error branches covered | 0-6 | Are there diamond nodes and at least one error/exception branch? |

### 3. Functional Specs (20 pts) — Mode A: evaluates prd-ui-functions.md

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Placement & Interaction completeness | 0-7 | Does every UI Function have Placement? Does User Interaction Flow cover the full path? |
| Data Requirements & States clarity | 0-7 | Are field tables and state tables filled completely? Are sources and triggers explicit? |
| Validation Rules explicit | 0-6 | Does every UI Function have validation rules that are actionable (not just "validate input")? |

### 3. Flow Completeness (20 pts) — Mode B: evaluates prd-spec.md Flow Description

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Flow steps describe complete business process | 0-7 | Does the text cover all steps from trigger to end state, including state transitions? |
| Data flow documented (if multi-system) | 0-7 | For multi-system features: is the Data Flow table complete? For single-system: auto-full-score if N/A |
| Exception handling and edge cases covered | 0-6 | Are error paths, retry logic, and failure states documented? |

### 4. User Stories (30 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Coverage: one story per target user | 0-7 | Does every user type from the background section have at least one story? |
| Format correct (As a / I want / So that) | 0-7 | Do all stories follow the format? Are actions concrete (not "manage", "handle")? |
| AC per story (Given/When/Then) | 0-6 | Does every story have at least one AC in Given/When/Then format? |
| AC verifiability & boundary coverage | 0-10 | Are ACs objectively testable? Do they cover happy path, error cases, and edge conditions? Can each "Then" be verified without subjective judgment? |

### 5. Scope Clarity (15 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| In-scope items are concrete deliverables | 0-5 | Each item is a specific feature/screen/API, not a vague area |
| Out-of-scope explicitly lists deferred items | 0-4 | Are deferred items named, not just implied by absence? |
| Scope consistent with functional specs and user stories | 0-6 | Do the in-scope items match what's described in Functional Specs and user stories? |

## Deduction Rules

- **Missing required section**: 0 pts for that dimension
- **Vague language without quantification**: -2 pts per instance ("better UX", "faster", "improved")
- **Cross-file/section inconsistency**: -3 pts per conflict (e.g., scope says X is out but prd-ui-functions.md describes X; user story references a role not in Background)
- **Placeholder text ("TBD", "TODO")**: -2 pts per instance
