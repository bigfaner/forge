---
scale: 1000
target: 900
iterations: 3
type: consistency
context:
  conventions: []
  business-rules: auto
---

# Cross-Document Consistency Rubric

**Total: 1000 points**

## Required Files

| File | Required |
|------|----------|
| manifest.md | ✓ |
| prd-spec.md | ✓ |
| prd-user-stories.md | If exists |
| prd-ui-functions.md | If exists (UI features) |
| tech-design.md | If exists |
| api-handbook.md | If exists |
| ui-design.md | If exists |
| tasks-index.json | If exists |

## Scoping Rule

The scorer evaluates consistency **across** documents, not individual document quality.
Individual document quality is the job of eval-prd, eval-design, eval-ui, etc.
This rubric measures alignment between documents.

## Scoring Modes

### Mode: docs — Cross-Document Consistency (default)

Activated when `--scope docs` or no scope specified.

| Dimension | Points |
|-----------|--------|
| 1. PRD-Design Alignment | 250 |
| 2. PRD-UI Consistency | 150 |
| 3. Design-Task Coverage | 200 |
| 4. Terminology Consistency | 150 |
| 5. Data Model Consistency | 150 |
| 6. Traceability Completeness | 100 |
| **Total** | **1000** |

> N/A rule: Dimensions 2, 3, 5 grant full points when the relevant documents don't exist (e.g., no UI docs → Dimension 2 = 150/150). Dimension 2 requires `prd-ui-functions.md` OR `ui-design.md`; Dimension 3 requires `tasks-index.json`; Dimension 5 requires `tech-design.md` with data models.

### Mode: full — Document-Code Consistency

Activated when `--scope full` is specified. Includes code-snapshot.md in the bundle.

| Dimension | Points |
|-----------|--------|
| 1. PRD-Design Alignment | 150 |
| 2. PRD-UI Consistency | 100 |
| 3. Design-Task Coverage | 150 |
| 4. Terminology Consistency | 150 |
| 5. Data Model Consistency | 100 |
| 6. Traceability Completeness | 100 |
| 7. Interface-Code Alignment | 150 |
| 8. Data Model-Code Alignment | 100 |
| **Total** | **1000** |

## Dimensions

### 1. PRD-Design Alignment (250 pts / 150 pts)

Evaluates whether PRD and Tech Design describe the same feature.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Functional coverage: every PRD functional section has a corresponding design section | 0-100 / 0-60 | Compare PRD flow/functional specs against design architecture. No PRD function should be missing from design. |
| No orphan design: every design component traces back to a PRD requirement | 0-80 / 0-50 | Design sections with no PRD justification indicate scope creep. |
| Error handling alignment: PRD error cases covered by design | 0-70 / 0-40 | PRD edge cases and error flows should appear in design error handling. |

### 2. PRD-UI Consistency (150 pts / 100 pts)

> Skip and award full points if neither `prd-ui-functions.md` nor `ui-design.md` exists.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| UI component coverage: every UI function in PRD has a matching UI design component | 0-60 / 0-40 | Compare prd-ui-functions.md entries against ui-design.md components. |
| Data binding alignment: UI data fields match PRD data requirements | 0-50 / 0-30 | Field names and types in UI design should match PRD data requirements tables. |
| Validation rules: UI validation matches PRD validation specs | 0-40 / 0-30 | Compare validation rules between PRD and UI design. |

### 3. Design-Task Coverage (200 pts / 150 pts)

> Skip and award full points if `tasks-index.json` does not exist.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Design module coverage: every design section maps to at least one task | 0-80 / 0-60 | Cross-reference tech-design.md sections against task files. |
| PRD AC coverage: every acceptance criterion in user stories is testable via tasks | 0-70 / 0-50 | Map prd-user-stories.md ACs to task descriptions. |
| Task source references: tasks reference correct design sections | 0-50 / 0-40 | Tasks should cite their originating design/PRD sections. |

### 4. Terminology Consistency (150 pts / 150 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Entity naming: same entity uses the same name across all documents | 0-60 | No synonyms for the same concept (e.g., "user" in PRD vs "account" in design). Check manifest.md traceability table for entity names. |
| Field naming: same field uses the same name across all documents | 0-50 | Field names in PRD data tables, design data models, and API params must be identical. |
| Status/enum consistency: status values and enums match across documents | 0-40 | Compare status values in PRD flows, design state diagrams, and UI states. |

### 5. Data Model Consistency (150 pts / 100 pts)

> Skip and award full points if tech-design.md has no data models (no DB feature).

| Criterion | Points | What to check |
|-----------|--------|---------------|
| ER diagram matches PRD data descriptions | 0-60 / 0-40 | Entities and relationships in er-diagram.md should match PRD data flow and data requirements. |
| API parameters match data model fields | 0-50 / 0-40 | API request/response fields in api-handbook.md should use the same names and types as data models. |
| Schema.sql consistent with design data models | 0-40 / 0-20 | Column names, types, and constraints in schema.sql match design data model definitions. |

### 6. Traceability Completeness (100 pts / 100 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Traceability table in manifest.md is complete | 0-50 | Every PRD section row has non-empty Design Section cell (if design exists). |
| No orphan rows: traceability references point to existing sections | 0-50 | Cross-check that referenced section names/IDs actually exist in the linked documents. |

### 7. Interface-Code Alignment (150 pts) — Mode: full only

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Function signatures match: design interface functions exist in code with correct parameters | 0-60 | Compare api-handbook.md or tech-design.md interface definitions against code-snapshot.md. |
| Return types consistent | 0-50 | Return types in code match design interface specifications. |
| Error types consistent | 0-40 | Error types/codes in code match design error handling specifications. |

### 8. Data Model-Code Alignment (100 pts) — Mode: full only

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Struct/class fields match design data models | 0-60 | Field names, types in code structs/classes match tech-design.md data model definitions. |
| Validation constraints match | 0-40 | Code validation logic (required fields, min/max, patterns) matches PRD/design constraints. |

## Deduction Rules

- **Cross-document inconsistency**: -30 pts per conflict (a term used differently, a function missing from downstream doc, a data field mismatch)
- **Missing traceability link**: -20 pts per empty cell in traceability table
- **Orphan content** (downstream doc has content with no upstream justification): -30 pts per instance
- **Terminology synonym** (same concept, different name): -20 pts per pair
