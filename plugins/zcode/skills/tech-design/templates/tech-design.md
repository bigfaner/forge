---
created: YYYY-MM-DD
prd: prd/prd-spec.md
status: Draft
---

# Technical Design: <Feature Name>

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

### Model 1: <!-- Name -->

<!-- Every field must have name, type, and constraint. No prose-only descriptions. -->
```
ModelName = {
    fieldName: Type       // constraint or note
    fieldName: Type       // constraint or note
}
```

## Error Handling

### Error Types & Codes

| Error Code | Name | Description | HTTP Status |
|------------|------|-------------|-------------|
| <!-- e.g. ERR_NOT_FOUND --> | <!-- e.g. ResourceNotFoundError --> | <!-- when it occurs --> | <!-- e.g. 404 --> |

### Propagation Strategy
<!-- How errors flow between layers. Who catches what? How are internal errors translated to API responses? -->

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
