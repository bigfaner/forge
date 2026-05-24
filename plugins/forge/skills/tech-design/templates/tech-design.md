---
created: "{{DATE}}"
prd: prd/prd-spec.md
status: Draft
---

# Technical Design: {{FEATURE_NAME}}

## Overview
<!-- High-level technical approach -->

## Architecture

### Layer Placement
<!-- Which layer(s) this feature belongs to -->

### Component Diagram
```
+----------+     +----------+
| Component|---->| Component|
+----------+     +----------+
```

### Dependencies
<!-- Internal and external dependencies -->

## Interfaces

### Interface 1: <!-- Name -->

<!-- Each interface must have typed parameters and return values.
   A developer should be able to implement it without guessing any types or shapes. -->
```
// Full typed signature
methodName(param1: Type, param2: Type): ReturnType
```

## Data Models

<!-- CONDITIONAL: Choose ONE format based on db-schema value.
     Delete this comment, the unused format, and UNWRAP the chosen format
     (remove its surrounding <!-- --> delimiters) so it renders as visible content. -->

<!-- DB-Schema: yes — USE THIS FORMAT (delete the "no" block, then unwrap this block):
> Full database design in separate files.

**ER Diagram**: design/er-diagram.md
**SQL Schema**: design/schema.sql

### Field Quick Reference
| Model | Key Fields | Notes |
|-------|------------|-------|
| ModelName | field1, field2 | brief note |
-->

<!-- DB-Schema: no — USE THIS FORMAT (delete the "yes" block, then unwrap this block):

### Model 1: Name

ModelName = {
    fieldName: Type       // constraint or note
}
-->

## Error Handling

### Error Types & Codes

| Error Code | Name | Description | HTTP Status |
|------------|------|-------------|-------------|
| <!-- e.g. ERR_NOT_FOUND --> | <!-- e.g. ResourceNotFoundError --> | <!-- when it occurs --> | <!-- e.g. 404 --> |

### Propagation Strategy
<!-- How errors flow between layers. Who catches what? How are internal errors translated to API responses? -->

## Cross-Layer Data Map

<!-- Required when the feature spans multiple architectural layers (e.g., database ↔ API ↔ UI).
     Skip with "Single-layer feature — not applicable." when the feature is confined to one layer.
     Every data field that crosses layer boundaries MUST appear here. -->

| Field Name | Storage Layer | Backend Model | API/DTO | Frontend Type | Validation Rule |
|------------|---------------|---------------|---------|---------------|-----------------|
| <!-- e.g. user_id --> | <!-- e.g. UUID, NOT NULL --> | <!-- e.g. User.ID uuid.UUID --> | <!-- e.g. json:"userId" --> | <!-- e.g. userId: string --> | <!-- e.g. required, UUID format --> |

## Integration Specs

<!-- Required when any UI Function has placement: existing-page.
     Skip with "No existing-page integrations — not applicable." otherwise. -->

### Integration: {{Component Name}} → {{Target Page}}

- **Target File**: {{file path of the existing page component}}
- **Insertion Point**: {{where to add the component in the page, e.g., "above SubItemsTable"}}
- **Data Source**: {{API call or data hook needed, e.g., "decision logs API by mainItemId"}}

<!-- Repeat for each existing-page integration -->

## Testing Strategy

### Per-Layer Test Plan

| Layer | Test Type | Tool | What to Test | Coverage Target |
|-------|-----------|------|--------------|-----------------|
| <!-- e.g. API --> | <!-- Integration --> | <!-- e.g. supertest --> | <!-- endpoint contracts --> | <!-- e.g. 80% --> |

### Key Test Scenarios
<!-- Critical scenarios that MUST pass: happy path, error paths, edge cases -->

### Overall Coverage Target
<!-- Single numeric target, e.g., 80% -->

## Security Considerations

### Threat Model
<!-- Potential security risks -->

### Mitigations
<!-- How risks are addressed -->

## PRD Coverage Map

<!-- Map each PRD acceptance criterion to a design component.
   Ensures /breakdown-tasks can derive tasks for every requirement. -->

| PRD Requirement / AC | Design Component | Interface / Model |
|----------------------|------------------|-------------------|
| <!-- from prd-spec or user-stories --> | <!-- module/component --> | <!-- specific interface or model --> |

## Open Questions
- [ ] Question 1
- [ ] Question 2

## Appendix

### Alternatives Considered
| Approach | Pros | Cons | Why Not Chosen |
|----------|------|------|----------------|
| Approach A | ... | ... | ... |

### References
- Link to relevant docs
