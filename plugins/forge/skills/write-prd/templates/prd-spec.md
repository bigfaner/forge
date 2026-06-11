---
feature: "{{FEATURE_NAME}}"
status: Draft
db-schema: "{{DB_SCHEMA}}"
---

# {{FEATURE_NAME}} — PRD Spec

> PRD Spec: defines WHAT the feature is and why it exists.

## Background

### Why (Reason)
<!-- Describe the root cause that triggered this requirement -->

### What (Target)
<!-- Describe the functionality or system to be implemented -->

### Who (Users)
<!-- Describe the target user roles -->

<!-- Example: Currently, unreachable delivery areas are configured only at the courier-type level, but each business system has its own special needs — a courier service may be reachable in general, yet a business system may designate it as unreachable based on its own business scenario. Therefore, a per-business-system unreachable area configuration is needed. -->

## Goals

<!-- Goals or benefits, quantified wherever possible -->

| Goal | Metric | Notes |
|------|--------|-------|
| <!-- Goal 1 --> | <!-- e.g., efficiency improvement 10% --> | <!-- Notes --> |
| <!-- Goal 2 --> | <!-- --> | <!-- --> |

## Scope

### In Scope
- [ ] <!-- Feature point 1 -->
- [ ] <!-- Feature point 2 -->

### Out of Scope
- <!-- Excluded item 1 -->
- <!-- Excluded item 2 -->

## Flow Description

### Business Flow Description

<!-- Detailed description of each step and state transition in the business flow -->
<!-- Include: main business flow steps, key decision points and branching logic, exception handling flow, state machine transitions -->

### Business Flow Diagram

> **Required**: Use Mermaid to draw the business flow diagram. Text-only descriptions are not acceptable.

```mermaid
flowchart TD
    Start([Start]) --> Step1[Step 1]
    Step1 --> Decision{Decision}
    Decision -->|Yes| Step2[Step 2]
    Decision -->|No| Error[Exception Handling]
    Step2 --> End([End])
    Error --> End
```

<!-- The flow diagram must include: complete main flow path, key decision points (diamond nodes), exception branches, user interaction nodes -->

### Data Flow Description

<!-- Required for multi-system interaction; remove this section for single-system scenarios -->

| Data Flow ID | Source System | Target System | Data Content | Transport | Frequency | Format | Notes |
|-----------|--------|----------|----------|----------|------|------|------|
| DF001 | <!-- --> | <!-- --> | <!-- --> | <!-- REST API / Message Queue --> | <!-- --> | <!-- JSON / XML --> | <!-- --> |

## Functional Specs

> UI 功能规格详见 [prd-ui-functions.md](./prd-ui-functions.md)。

### Related Changes

<!-- Fill this section if changes affect other modules/systems -->

| # | Project | Module | Change Point | Updated Logic |
|------|----------|----------|------------|----------------|
| 1 | <!-- --> | <!-- --> | <!-- --> | <!-- --> |

## Other Notes

### Performance Requirements
- Response time: <!-- -->
- Concurrency: <!-- -->
- Data storage: <!-- -->
- Compatibility: <!-- Browsers, resolutions, mobile devices -->

### Data Requirements
- Data tracking: <!-- -->
- Data initialization: <!-- -->
- Data migration: <!-- -->

### Monitoring Requirements
- <!-- API or service monitoring, alerting mechanisms -->

### Security Requirements
- Transport encryption: <!-- -->
- Storage encryption: <!-- -->
- Display masking: <!-- -->
- Rate limiting: <!-- -->

---

## Quality Checklist

- [ ] Is the requirement title accurate and descriptive
- [ ] Does the background include all three elements: reason, target, users
- [ ] Are the goals quantified
- [ ] Is the flow description complete
- [ ] Does the business flow diagram exist (Mermaid format)
- [ ] Is prd-ui-functions.md referenced and UI specs complete
- [ ] Are related changes thoroughly analyzed
- [ ] Are non-functional requirements considered (performance / data / monitoring / security)
- [ ] Are all tables filled completely
- [ ] Is there any ambiguous or vague wording
- [ ] Is the spec actionable and verifiable
