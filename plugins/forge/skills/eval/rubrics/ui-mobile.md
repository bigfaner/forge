---
scale: 1000
target: 950
iterations: 3
type: ui-mobile
context:
  conventions: [ux, frontend]
  business-rules: auto
---

# Mobile Design Evaluation Rubric

**Total: 1000 points**
**Platform: mobile**

## Perspectives

Each dimension represents an independent stakeholder perspective. The scorer must evaluate from that role's standpoint — not from a generic "quality" viewpoint.

| Perspective | Role | Core Question |
|-------------|------|---------------|
| Requirement Coverage | Product Manager | Are all PRD UI requirements covered? Edge cases? |
| Touch Experience | End User | Is it usable with touch? Ergonomic? Platform-native feel? |
| Adaptive Layout | Designer | Does it adapt across screen sizes, orientations, and safe areas? |
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
| Touch target sizes (minimum 44pt) | yes |
| Adaptive layout breakpoints | yes |
| Safe area handling | yes |

## Dimensions

### 1. Requirement Coverage (250 pts) — Product Manager Perspective

| Criterion | Points | What to check |
|-----------|--------|---------------|
| UI function coverage | 0-80 | Does every UI function from `prd-ui-functions.md` have a corresponding screen/component in the design? Any gaps? |
| Navigation coverage | 0-40 | If PRD defines `## Navigation Architecture`, does the design cover all navigation flows? Are tab bars, navigation stacks, and modal presentations consistent? |
| State coverage | 0-80 | Are all states (loading, empty, error, populated) addressed? Are mobile-specific states covered: background/foreground transitions, push notification interruption, offline/weak network? |
| Edge case handling | 0-50 | Are boundary conditions addressed: landscape/portrait transitions, long text overflow, no data, permission denied, slow network, incoming call interruption? Or does the design only show the happy path? |

### 2. Touch Experience (250 pts) — End User Perspective

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Touch target sizing | 0-70 | Are all interactive elements at least 44pt in both dimensions? Are touch targets sized explicitly in the spec? Touch targets without size specification → -30 per instance. |
| Thumb reachability | 0-60 | Are primary actions placed within thumb-reach zones (bottom of screen)? Are destructive actions placed away from frequent touch zones? |
| Gesture intuitiveness | 0-60 | Are gesture patterns conventional (swipe to dismiss, pull to refresh, pinch to zoom)? Or does the user need to learn custom gestures? |
| Platform convention adherence | 0-60 | Does the design follow platform conventions (iOS HIG / Material Design)? Are platform-native interaction patterns used (e.g., swipe-back on iOS, back gesture on Android)? |

### 3. Adaptive Layout (250 pts) — Designer Perspective

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Multi-screen adaptation | 0-70 | Does the design specify how layout adapts across phone sizes (compact/regular)? Are breakpoints defined explicitly? Missing landscape/portrait adaptation → -50 pts. |
| Orientation handling | 0-60 | Does the design address both portrait and landscape orientations? Are layout differences specified, not just "responsive"? |
| Safe area handling | 0-60 | Are safe areas accounted for (notch, home indicator, status bar)? Does content avoid overlapping system UI? Missing safe area handling → -40 pts. |
| Platform-native feel | 0-60 | Do components look and feel native to the target platform(s)? Are navigation patterns, typography, and spacing consistent with platform expectations? |

### 4. Implementability (250 pts) — Developer Perspective

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Adaptive breakpoints clarity | 0-80 | Are layout breakpoints stated explicitly (e.g., "compact width < 600pt")? Can a developer implement responsive layout rules directly from the spec? |
| Platform interaction specificity | 0-80 | Is every interaction pattern tied to a specific platform behavior (e.g., "iOS: swipe-back gesture", "Android: back button")? Or are there vague entries like "go back"? |
| Navigation pattern clarity | 0-90 | Is the navigation structure explicit (tab bar, stack navigation, modal sheet)? Are transition animations specified? Can a developer wire navigation without guessing? |

## Deduction Rules

- **Touch targets without size specification**: -30 pts per instance (from Touch Experience)
- **Missing landscape/portrait adaptation**: -50 pts (from Adaptive Layout)
- **Missing safe area handling**: -40 pts (from Adaptive Layout)
- **Missing required section**: 0 pts for that dimension
- **Vague language without quantification**: -20 pts per instance ("better UX", "faster", "improved")
- **Cross-section inconsistency**: -30 pts per conflict (e.g., interaction table contradicts layout, state missing from data binding)
- **Happy-path only design** (no error/empty/loading states): -50 pts from Touch Experience
- **Navigation Architecture gap**: -20 pts per PRD navigation entry not covered in design (from Requirement Coverage)
- **PRD UI function gap**: -30 pts per unaddressed UI function (from Requirement Coverage)
- **Orphan UI elements** (no data binding): -30 pts per element (from Implementability)
- **Placeholder text ("TBD", "TODO", "lorem ipsum")**: -20 pts per instance
