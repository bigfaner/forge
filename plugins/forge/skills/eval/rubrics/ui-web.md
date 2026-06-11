---
scale: 1000
target: 950
iterations: 3
type: ui-web
context:
  conventions: [ux, frontend]
  business-rules: auto
---

# UI Design Evaluation Rubric

**Total: 1000 points**

## Perspectives

Each dimension represents an independent stakeholder perspective. The scorer must evaluate from that role's standpoint — not from a generic "quality" viewpoint.

| Perspective | Role | Core Question |
|-------------|------|---------------|
| Requirement Coverage | Product Manager | Are all PRD UI requirements covered? Edge cases? |
| User Experience | End User | Is it usable? Intuitive? Accessible to everyone? |
| Design Integrity | Designer | Does it follow the chosen design system? Visually consistent? States complete? |
| Implementability | Developer | Can I code from this without guessing? |

## Required Sections

Mark missing required sections as 0 pts for the relevant dimension:

| Section | Required |
|---------|----------|
| Design System Reference | yes |
| Component definitions (layout, states, interactions, data binding) | yes |
| States table (Default/Loading/Empty/Error) | yes |
| Interactions table (Trigger/Action/Feedback) | yes |
| Data Binding table (UI Element/Field/Source) | yes |

## Dimensions

### 1. Requirement Coverage (250 pts) — Product Manager Perspective

| Criterion | Points | What to check |
|-----------|--------|---------------|
| UI function coverage | 0-80 | Does every UI function from `prd-ui-functions.md` have a corresponding component in the design? Any gaps? |
| Navigation Architecture coverage | 0-40 | If PRD defines `## Navigation Architecture`, does the design cover all primary navigation entries and Secondary Pages? Are navigation targets consistent with page names? |
| State requirement coverage | 0-80 | Are all states defined in `prd-ui-functions.md` (loading, empty, error, populated) addressed in the design's state tables? |
| Edge case handling | 0-50 | Are boundary conditions addressed: long text overflow, no data, permission denied, slow network, concurrent actions? Or does the design only show the happy path? |

### 2. User Experience (250 pts) — End User Perspective

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Information hierarchy | 0-80 | Can a user scan the page and immediately understand what matters most? Is visual weight correctly distributed — or is everything equally prominent? |
| Interaction intuitiveness | 0-80 | Are interaction patterns conventional (click button → action, scroll → more content)? Or does the user need to learn custom behaviors? |
| Accessibility | 0-90 | Contrast ratios adequate? Form labels present? Keyboard navigation considered? Are loading/error states communicated accessibly (not just visual)? |

### 3. Design Integrity (250 pts) — Designer Perspective

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Design system adherence | 0-80 | Does the design follow its referenced style system's rules (colors, typography, spacing, component shapes)? Or does it mix styles arbitrarily? |
| Visual coherence | 0-90 | Do all components look like they belong to the same product? Consistent border-radius, shadow depth, spacing rhythm across components? Cross-page navigation elements identical where shared (same tab bar, same nav bar layout)? |
| State completeness | 0-80 | Does every interactive component cover all applicable states (Default/Loading/Empty/Error)? Are state transitions described — how does a component go from Loading to Empty? |

### 4. Implementability (250 pts) — Developer Perspective

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Layout specificity | 0-80 | Can a developer build the layout from the spec alone? Are widths, spacing, and positioning stated precisely — or left to "interpretation"? |
| Data binding explicit | 0-80 | Is every UI element mapped to a data field and source in the Data Binding table? No orphan elements with no data source? |
| Interaction unambiguity | 0-90 | Is every trigger → action → feedback chain explicit in the Interactions table? Or are there vague entries like "handle click" without stating what happens? |

## Deduction Rules

- **Missing required section**: 0 pts for that dimension
- **Vague language without quantification**: -20 pts per instance ("better UX", "faster", "improved")
- **Cross-section inconsistency**: -30 pts per conflict (e.g., interaction table contradicts layout, state missing from data binding)
- **Happy-path only design** (no error/empty/loading states): -50 pts from Design Integrity
- **Navigation Architecture gap**: -20 pts per PRD navigation entry not covered in design (from Requirement Coverage)
- **Cross-page inconsistency**: -30 pts per inconsistency in shared navigation elements (from Design Integrity)
- **PRD UI function gap**: -30 pts per unaddressed UI function (from Requirement Coverage)
- **Orphan UI elements** (no data binding): -30 pts per element (from Implementability)
- **Placeholder text ("TBD", "TODO", "lorem ipsum")**: -20 pts per instance
