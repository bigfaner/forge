# Cross-Document Consistency Rubric

**Total: 100 points**
**Report template:** `plugins/forge/skills/eval-consistency/templates/report.md`

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
| 1. PRD-Design Alignment | 25 |
| 2. PRD-UI Consistency | 15 |
| 3. Design-Task Coverage | 20 |
| 4. Terminology Consistency | 15 |
| 5. Data Model Consistency | 15 |
| 6. Traceability Completeness | 10 |
| **Total** | **100** |

> N/A rule: Dimensions 2, 3, 5 grant full points when the relevant documents don't exist (e.g., no UI docs → Dimension 2 = 15/15). Dimension 2 requires `prd-ui-functions.md` OR `ui-design.md`; Dimension 3 requires `tasks-index.json`; Dimension 5 requires `tech-design.md` with data models.

### Mode: full — Document-Code Consistency

Activated when `--scope full` is specified. Includes code-snapshot.md in the bundle.

| Dimension | Points |
|-----------|--------|
| 1. PRD-Design Alignment | 15 |
| 2. PRD-UI Consistency | 10 |
| 3. Design-Task Coverage | 15 |
| 4. Terminology Consistency | 15 |
| 5. Data Model Consistency | 10 |
| 6. Traceability Completeness | 10 |
| 7. Interface-Code Alignment | 15 |
| 8. Data Model-Code Alignment | 10 |
| **Total** | **100** |

## Dimensions

### 1. PRD-Design Alignment (25 pts / 15 pts)

Evaluates whether PRD and Tech Design describe the same feature.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Functional coverage: every PRD functional section has a corresponding design section | 0-10 / 0-6 | Compare PRD flow/functional specs against design architecture. No PRD function should be missing from design. |
| No orphan design: every design component traces back to a PRD requirement | 0-8 / 0-5 | Design sections with no PRD justification indicate scope creep. |
| Error handling alignment: PRD error cases covered by design | 0-7 / 0-4 | PRD edge cases and error flows should appear in design error handling. |

### 2. PRD-UI Consistency (15 pts / 10 pts)

> Skip and award full points if neither `prd-ui-functions.md` nor `ui-design.md` exists.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| UI component coverage: every UI function in PRD has a matching UI design component | 0-6 / 0-4 | Compare prd-ui-functions.md entries against ui-design.md components. |
| Data binding alignment: UI data fields match PRD data requirements | 0-5 / 0-3 | Field names and types in UI design should match PRD data requirements tables. |
| Validation rules: UI validation matches PRD validation specs | 0-4 / 0-3 | Compare validation rules between PRD and UI design. |

### 3. Design-Task Coverage (20 pts / 15 pts)

> Skip and award full points if `tasks-index.json` does not exist.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Design module coverage: every design section maps to at least one task | 0-8 / 0-6 | Cross-reference tech-design.md sections against task files. |
| PRD AC coverage: every acceptance criterion in user stories is testable via tasks | 0-7 / 0-5 | Map prd-user-stories.md ACs to task descriptions. |
| Task source references: tasks reference correct design sections | 0-5 / 0-4 | Tasks should cite their originating design/PRD sections. |

### 4. Terminology Consistency (15 pts / 15 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Entity naming: same entity uses the same name across all documents | 0-6 | No synonyms for the same concept (e.g., "user" in PRD vs "account" in design). Check manifest.md traceability table for entity names. |
| Field naming: same field uses the same name across all documents | 0-5 | Field names in PRD data tables, design data models, and API params must be identical. |
| Status/enum consistency: status values and enums match across documents | 0-4 | Compare status values in PRD flows, design state diagrams, and UI states. |

### 5. Data Model Consistency (15 pts / 10 pts)

> Skip and award full points if tech-design.md has no data models (no DB feature).

| Criterion | Points | What to check |
|-----------|--------|---------------|
| ER diagram matches PRD data descriptions | 0-6 / 0-4 | Entities and relationships in er-diagram.md should match PRD data flow and data requirements. |
| API parameters match data model fields | 0-5 / 0-4 | API request/response fields in api-handbook.md should use the same names and types as data models. |
| Schema.sql consistent with design data models | 0-4 / 0-2 | Column names, types, and constraints in schema.sql match design data model definitions. |

### 6. Traceability Completeness (10 pts / 10 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Traceability table in manifest.md is complete | 0-5 | Every PRD section row has non-empty Design Section cell (if design exists). |
| No orphan rows: traceability references point to existing sections | 0-5 | Cross-check that referenced section names/IDs actually exist in the linked documents. |

### 7. Interface-Code Alignment (15 pts) — Mode: full only

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Function signatures match: design interface functions exist in code with correct parameters | 0-6 | Compare api-handbook.md or tech-design.md interface definitions against code-snapshot.md. |
| Return types consistent | 0-5 | Return types in code match design interface specifications. |
| Error types consistent | 0-4 | Error types/codes in code match design error handling specifications. |

### 8. Data Model-Code Alignment (10 pts) — Mode: full only

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Struct/class fields match design data models | 0-6 | Field names, types in code structs/classes match tech-design.md data model definitions. |
| Validation constraints match | 0-4 | Code validation logic (required fields, min/max, patterns) matches PRD/design constraints. |

## Deduction Rules

- **Cross-document inconsistency**: -3 pts per conflict (a term used differently, a function missing from downstream doc, a data field mismatch)
- **Missing traceability link**: -2 pts per empty cell in traceability table
- **Orphan content** (downstream doc has content with no upstream justification): -3 pts per instance
- **Terminology synonym** (same concept, different name): -2 pts per pair
