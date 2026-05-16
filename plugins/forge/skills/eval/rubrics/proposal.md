---
scale: 1000
target: 900
iterations: 3
type: proposal
context:
  conventions: []
  business-rules: []
---

# Proposal Evaluation Rubric

**Total: 1000 points**

## Dimensions

### 1. Problem Definition (110 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Problem stated clearly | 0-40 | Is the core problem unambiguous? Could two readers interpret it differently? |
| Evidence provided | 0-40 | Is there data, user feedback, or concrete examples backing the problem? Not just "we think..." |
| Urgency justified | 0-30 | Why solve this now? What happens if we don't? What's the cost of delay? |

### 2. Solution Clarity (120 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Approach is concrete | 0-40 | Can a reader explain back what will be built? Or is it vague hand-waving? |
| User-facing behavior described | 0-45 | What does the end user experience? Not internals — the observable behavior |
| Technical direction clear | 0-35 | Is there enough technical hint to know the general implementation approach? |

### 3. Industry Benchmarking (120 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Industry solutions referenced | 0-40 | Are real-world solutions/patterns for this type of problem cited? Product names, open-source projects, or published patterns? Not just self-invented options |
| At least 3 meaningful alternatives | 0-30 | Including "do nothing". Each alternative must be a genuinely different approach, not a straw man. At least one must be an industry-validated solution |
| Honest trade-off comparison | 0-25 | Are pros/cons based on actual project constraints? Not cherry-picked? |
| Chosen approach justified against benchmarks | 0-25 | Why this approach over industry standards? Or why adopt the standard? Clear rationale required |

### 4. Requirements Completeness (110 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Scenario coverage | 0-40 | Happy path + edge cases + error scenarios all identified? Or only the golden path? |
| Non-functional requirements | 0-40 | Performance, security, compatibility, accessibility — are relevant NFRs called out? |
| Constraints & dependencies | 0-30 | External systems, tech constraints, prerequisite conditions named? |

### 5. Solution Creativity (100 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Novelty over industry baseline | 0-40 | Is the proposal just copying an industry solution, or does it innovate beyond? Is the differentiation from the benchmark clearly articulated? |
| Cross-domain inspiration | 0-35 | Does it borrow ideas from other domains/products that face similar problems? |
| Simplicity of insight | 0-25 | Is the creative leap elegant ("why didn't I think of that") or forced/overengineered? |

### 6. Feasibility (100 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Technical feasibility | 0-40 | Can current tech stack support this? Any showstopper dependencies? |
| Resource & timeline feasibility | 0-30 | Does the team have skills/bandwidth? Is the scope realistic for the proposed timeline? |
| Dependency readiness | 0-30 | Are external APIs/services available? Are upstream prerequisites in place? |

### 7. Scope Definition (80 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| In-scope items are concrete | 0-30 | Each item is a deliverable, not a vague area |
| Out-of-scope explicitly listed | 0-25 | Are deferred items named, not just implied? |
| Scope is bounded | 0-25 | Can a team execute this in a defined timeframe? Or is it open-ended? |

### 8. Risk Assessment (90 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Risks identified | 0-30 | At least 3 meaningful risks, not trivial ones |
| Likelihood + impact rated | 0-30 | Is the assessment honest? Not all "low likelihood, high impact"? |
| Mitigations are actionable | 0-30 | Can someone act on the mitigation? Or is it "we'll handle it"? |

### 9. Success Criteria (80 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Criteria are measurable and testable | 0-55 | Can you objectively verify each criterion? Could you write a test or checklist? "Works well" is neither measurable nor testable |
| Coverage is complete | 0-25 | Do criteria cover all in-scope items? Any gaps? |

### 10. Logical Consistency (90 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Solution addresses the stated problem | 0-35 | Does the proposed solution actually solve the problem defined in section 1? Or is there a gap? |
| Scope ↔ Solution ↔ Success Criteria aligned | 0-30 | No contradictions between what we're building, what's in scope, and how we measure success |
| Requirements ↔ Solution coherent | 0-25 | Do the identified requirements map cleanly to the proposed solution? No orphan requirements or solution features with no requirement? |

## Section → Dimension Mapping

All sections are required; missing sections score 0 for their corresponding dimension.

| Section | Primary Dimension(s) |
|---------|---------------------|
| Problem (Evidence, Urgency) | 1. Problem Definition |
| Proposed Solution | 2. Solution Clarity |
| Proposed Solution > Innovation Highlights | 5. Solution Creativity |
| Requirements Analysis | 4. Requirements Completeness |
| Alternatives & Industry Benchmarking | 3. Industry Benchmarking |
| Feasibility Assessment | 6. Feasibility |
| Scope (In/Out) | 7. Scope Definition |
| Key Risks | 8. Risk Assessment |
| Success Criteria | 9. Success Criteria |
| *(cross-section)* | 10. Logical Consistency |

## Deduction Rules

- **Missing required section**: 0 pts for that dimension
- **Vague language without quantification**: -20 pts per instance ("better", "improved", "enhanced", "faster performance")
- **Placeholder text ("TBD", "TODO")**: -20 pts per instance
- **Straw-man alternative**: -20 pts per instance (an alternative presented only to be rejected)
