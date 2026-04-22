# PRD Evaluation Rubric

**Total: 100 points**
**Report template:** `plugins/zcode/skills/eval-prd/templates/report.md`

## Required Sections (prd-spec.md)

| Section | Required |
|---------|----------|
| 需求背景（原因/对象/人员） | ✓ |
| 需求目标 + 量化指标 | ✓ |
| Scope（In + Out） | ✓ |
| 流程说明 + Mermaid 流程图 | ✓ |
| 功能描述 | ✓ |

## Required Sections (prd-user-stories.md)

| Section | Required |
|---------|----------|
| User Stories | ✓ |
| Acceptance Criteria (Given/When/Then) | ✓ |

## Dimensions

### 1. Background & Goals (20 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Background has three elements (原因/对象/人员) | 0-7 | Are all three present and specific? |
| Goals are quantified | 0-7 | Is there at least one numeric target (%, count, time)? |
| Background and goals are logically consistent | 0-6 | Does the goal follow from the stated problem? |

### 2. Flow Diagrams (20 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Mermaid diagram exists | 0-7 | Is there at least one Mermaid flowchart? |
| Main path complete (start → end) | 0-7 | Does the diagram cover the full happy path? |
| Decision points + error branches covered | 0-6 | Are there diamond nodes and at least one error/exception branch? |

### 3. Functional Specs (20 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Tables complete (list page 7 elements, button 4 elements, form 2 elements) | 0-7 | Are all required table columns filled in? |
| Field descriptions clear | 0-7 | Is each field's purpose, type, and source stated? |
| Validation rules explicit | 0-6 | Are validation rules stated per field/button (not just "validate input")? |

### 4. User Stories (20 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Coverage: one story per target user | 0-7 | Does every user type from the background section have at least one story? |
| Format correct (As a / I want / So that) | 0-7 | Do all stories follow the format? Are actions concrete (not "manage", "handle")? |
| AC per story (Given/When/Then) | 0-6 | Does every story have at least one AC in Given/When/Then format? |

### 5. Scope Clarity (20 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| In-scope items are concrete deliverables | 0-7 | Each item is a specific feature/screen/API, not a vague area |
| Out-of-scope explicitly lists deferred items | 0-7 | Are deferred items named, not just implied by absence? |
| Scope consistent with functional specs and user stories | 0-6 | Do the in-scope items match what's described in 功能描述 and user stories? |

## Deduction Rules

- **Missing required section**: 0 pts for that dimension
- **Vague language without quantification**: -2 pts per instance ("better UX", "faster", "improved")
- **Inconsistency between sections**: -3 pts per conflict (e.g., scope says X is out but functional spec describes X)
- **Placeholder text ("TBD", "TODO")**: -2 pts per instance
