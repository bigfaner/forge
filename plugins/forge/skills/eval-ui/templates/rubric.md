# UI Design Evaluation Rubric

**Total: 100 points**
**Report template:** `plugins/forge/skills/eval-ui/templates/report.md`

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

### 1. Requirement Coverage (25 pts) — Product Manager Perspective

| Criterion | Points | What to check |
|-----------|--------|---------------|
| UI function coverage | 0-8 | Does every UI function from `prd-ui-functions.md` have a corresponding component in the design? Any gaps? |
| State requirement coverage | 0-8 | Are all states defined in `prd-ui-functions.md` (loading, empty, error, populated) addressed in the design's state tables? |
| Edge case handling | 0-9 | Are boundary conditions addressed: long text overflow, no data, permission denied, slow network, concurrent actions? Or does the design only show the happy path? |

### 2. User Experience (25 pts) — End User Perspective

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Information hierarchy | 0-8 | Can a user scan the page and immediately understand what matters most? Is visual weight correctly distributed — or is everything equally prominent? |
| Interaction intuitiveness | 0-8 | Are interaction patterns conventional (click button → action, scroll → more content)? Or does the user need to learn custom behaviors? |
| Accessibility | 0-9 | Contrast ratios adequate? Form labels present? Keyboard navigation considered? Are loading/error states communicated accessibly (not just visual)? |

### 3. Design Integrity (25 pts) — Designer Perspective

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Design system adherence | 0-8 | Does the design follow its referenced style system's rules (colors, typography, spacing, component shapes)? Or does it mix styles arbitrarily? |
| Visual coherence | 0-9 | Do all components look like they belong to the same product? Consistent border-radius, shadow depth, spacing rhythm across components? |
| State completeness | 0-8 | Does every interactive component cover all applicable states (Default/Loading/Empty/Error)? Are state transitions described — how does a component go from Loading to Empty? |

### 4. Implementability (25 pts) — Developer Perspective

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Layout specificity | 0-8 | Can a developer build the layout from the spec alone? Are widths, spacing, and positioning stated precisely — or left to "interpretation"? |
| Data binding explicit | 0-8 | Is every UI element mapped to a data field and source in the Data Binding table? No orphan elements with no data source? |
| Interaction unambiguity | 0-9 | Is every trigger → action → feedback chain explicit in the Interactions table? Or are there vague entries like "handle click" without stating what happens? |

## Deduction Rules

- **Missing required section**: 0 pts for that dimension
- **Vague language without quantification**: -2 pts per instance ("better UX", "faster", "improved")
- **Cross-section inconsistency**: -3 pts per conflict (e.g., interaction table contradicts layout, state missing from data binding)
- **Happy-path only design** (no error/empty/loading states): -5 pts from Design Integrity
- **PRD UI function gap**: -3 pts per unaddressed UI function (from Requirement Coverage)
- **Orphan UI elements** (no data binding): -3 pts per element (from Implementability)
- **Placeholder text ("TBD", "TODO", "lorem ipsum")**: -2 pts per instance
