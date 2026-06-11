---
scale: 1000
target: 900
iterations: 3
type: design
context:
  conventions: [api, error-handling]
  business-rules: auto
---

# Design Evaluation Rubric

**Total: 1000 points**

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

### 1. Architecture Clarity (170 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Layer placement explicit | 0-60 | Does the doc state which layer (API/service/repo/etc.) this belongs to? |
| Component diagram present | 0-60 | Is there an ASCII or text diagram showing components and relationships? |
| Dependencies listed | 0-50 | Are internal modules and external packages named? |

### 2. Interface & Model Definitions (170 pts)

**When `er-diagram.md` exists (db-schema: "yes"):**

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Interface signatures typed | 0-40 | Do all interfaces have typed params and return values (not prose)? |
| Inline models concrete | 0-40 | Are all non-DB model fields named with types and constraints? |
| ER diagram complete | 0-30 | Does er-diagram.md have Mermaid erDiagram with all entities, relationships, and cardinality? |
| SQL DDL directly usable | 0-30 | Can schema.sql be executed as-is? Inline COMMENT syntax, all FKs, indexes, defaults present? |
| Cross-layer consistency | 0-30 | Do field names in Cross-Layer Data Map match er-diagram.md entity column names? |

**When `er-diagram.md` absent (db-schema: "no"):**

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Interface signatures typed | 0-60 | Do all interfaces have typed params and return values (not prose)? |
| Models concrete | 0-60 | Are all model fields named with types and constraints (not just described)? |
| Directly implementable | 0-50 | Can a developer code from this without guessing any types or shapes? |

### 3. Error Handling (130 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Error types defined | 0-45 | Are custom error types or error codes explicitly defined? |
| Propagation strategy clear | 0-45 | Is there a stated strategy for how errors flow between layers? |
| HTTP status codes mapped | 0-40 | If API: are error types mapped to HTTP status codes? |

### 4. Testing Strategy (130 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Per-layer test plan | 0-45 | Does each layer have a stated test approach (unit/integration/e2e)? |
| Coverage target numeric | 0-45 | Is there a numeric coverage target (e.g., 80%)? |
| Test tooling named | 0-40 | Are specific test libraries/frameworks named? |

### 5. Breakdown-Readiness ★ (180 pts — critical gate)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Components enumerable | 0-65 | Can you list and count all components/modules? Or are they described vaguely? |
| Tasks derivable | 0-65 | Does each interface → at least one impl task? Each model → at least one schema task? |
| PRD AC coverage | 0-50 | If PRD exists: are all acceptance criteria addressed somewhere in the design? |

★ This dimension is the direct gate to `/breakdown-tasks`. A score below 160/180 blocks progression.

### 6. Security Considerations (80 pts)

Only scored if PRD has auth, data privacy, or multi-user requirements. Mark N/A (full credit) otherwise.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Threat model present | 0-40 | Are specific threats identified (not just "we'll add auth")? |
| Mitigations concrete | 0-40 | Is each threat paired with a specific countermeasure? |

### 7. Implementation Feasibility (140 pts)

This dimension uses injected context (project conventions, existing dependencies, architectural patterns) to assess whether the design can be implemented as described. The scorer should reference the injected conventions and business-rules to detect contradictions with project reality.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Dependencies available | 0-50 | Do referenced libraries/packages exist in the project's dependency manifest? If injected conventions declare a preferred library for a concern (e.g., error-handling), does the design use it? Deduct 15 pts per dependency that conflicts with project conventions or does not exist in the project |
| Architecture fits project structure | 0-50 | Does the proposed architecture align with the project's existing layer structure and module boundaries (per injected conventions)? Does it introduce patterns that contradict established project patterns (e.g., proposing MVC in a project with established hexagonal architecture)? |
| Technical claims grounded | 0-40 | Are performance claims, capacity estimates, and technology choices grounded in the project's actual tech stack (per injected context)? No speculative claims about capabilities that the current stack cannot deliver? |

## Deduction Rules

- **Missing required section**: 0 pts for that dimension
- **Prose-only (no code/diagram where expected)**: -50 pts from that dimension
- **PRD AC gap**: -30 pts per unaddressed acceptance criterion (from Breakdown-Readiness)
- **Vague language without quantification**: -20 pts per instance ("better performance", "faster", "improved")
- **Cross-section inconsistency**: -30 pts per conflict (e.g., interface contradicts data model, error handling conflicts with architecture)
- **Placeholder text ("TBD", "TODO")**: -20 pts per instance
