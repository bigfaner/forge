---
scale: 1000
target: 900
iterations: 3
type: design
---

# Design Evaluation Rubric

**Total: 1000 points**
**Report template:** `plugins/forge/skills/eval-design/templates/report.md`

## Required Sections

Mark missing required sections as 0 pts for that dimension:

| Section | Required |
|---------|----------|
| Overview + tech stack | ✓ |
| Architecture (layer + diagram) | ✓ |
| Interfaces | ✓ |
| Data Models | ✓ |
| Error Handling | ✓ |
| Testing Strategy | ✓ |
| Security Considerations | ○ (required if PRD has auth/data requirements) |

## Dimensions

### 1. Architecture Clarity (200 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Layer placement explicit | 0-70 | Does the doc state which layer (API/service/repo/etc.) this belongs to? |
| Component diagram present | 0-70 | Is there an ASCII or text diagram showing components and relationships? |
| Dependencies listed | 0-60 | Are internal modules and external packages named? |

### 2. Interface & Model Definitions (200 pts)

**When `er-diagram.md` exists (db-schema: "yes"):**

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Interface signatures typed | 0-50 | Do all interfaces have typed params and return values (not prose)? |
| Inline models concrete | 0-50 | Are all non-DB model fields named with types and constraints? |
| ER diagram complete | 0-30 | Does er-diagram.md have Mermaid erDiagram with all entities, relationships, and cardinality? |
| SQL DDL directly usable | 0-40 | Can schema.sql be executed as-is? Inline COMMENT syntax, all FKs, indexes, defaults present? |
| Cross-layer consistency | 0-30 | Do field names in Cross-Layer Data Map match er-diagram.md entity column names? |

**When `er-diagram.md` absent (db-schema: "no"):**

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Interface signatures typed | 0-70 | Do all interfaces have typed params and return values (not prose)? |
| Models concrete | 0-70 | Are all model fields named with types and constraints (not just described)? |
| Directly implementable | 0-60 | Can a developer code from this without guessing any types or shapes? |

### 3. Error Handling (150 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Error types defined | 0-50 | Are custom error types or error codes explicitly defined? |
| Propagation strategy clear | 0-50 | Is there a stated strategy for how errors flow between layers? |
| HTTP status codes mapped | 0-50 | If API: are error types mapped to HTTP status codes? |

### 4. Testing Strategy (150 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Per-layer test plan | 0-50 | Does each layer have a stated test approach (unit/integration/e2e)? |
| Coverage target numeric | 0-50 | Is there a numeric coverage target (e.g., 80%)? |
| Test tooling named | 0-50 | Are specific test libraries/frameworks named? |

### 5. Breakdown-Readiness ★ (200 pts — critical gate)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Components enumerable | 0-70 | Can you list and count all components/modules? Or are they described vaguely? |
| Tasks derivable | 0-70 | Does each interface → at least one impl task? Each model → at least one schema task? |
| PRD AC coverage | 0-60 | If PRD exists: are all acceptance criteria addressed somewhere in the design? |

★ This dimension is the direct gate to `/breakdown-tasks`. A score below 180/200 blocks progression.

### 6. Security Considerations (100 pts)

Only scored if PRD has auth, data privacy, or multi-user requirements. Mark N/A (full credit) otherwise.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Threat model present | 0-50 | Are specific threats identified (not just "we'll add auth")? |
| Mitigations concrete | 0-50 | Is each threat paired with a specific countermeasure? |

## Deduction Rules

- **Missing required section**: 0 pts for that dimension
- **Prose-only (no code/diagram where expected)**: -50 pts from that dimension
- **PRD AC gap**: -30 pts per unaddressed acceptance criterion (from Breakdown-Readiness)
- **Vague language without quantification**: -20 pts per instance ("better performance", "faster", "improved")
- **Cross-section inconsistency**: -30 pts per conflict (e.g., interface contradicts data model, error handling conflicts with architecture)
- **Placeholder text ("TBD", "TODO")**: -20 pts per instance
