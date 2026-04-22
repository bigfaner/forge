# Proposal Evaluation Rubric

**Total: 100 points**
**Report template:** `plugins/zcode/skills/eval-proposal/templates/report.md`

## Dimensions

### 1. Problem Definition (20 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Problem stated clearly | 0-7 | Is the core problem unambiguous? Could two readers interpret it differently? |
| Evidence provided | 0-7 | Is there data, user feedback, or concrete examples backing the problem? Not just "we think..." |
| Urgency justified | 0-6 | Why solve this now? What happens if we don't? |

### 2. Solution Clarity (20 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Approach is concrete | 0-7 | Can a reader explain back what will be built? Or is it vague hand-waving? |
| User-facing behavior described | 0-7 | What does the end user experience? Not internals — the observable behavior |
| Distinguishes from alternatives | 0-6 | Is it clear why this approach over others? What's the differentiator? |

### 3. Alternatives Analysis (15 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| At least 2 alternatives listed | 0-5 | Including "do nothing" as a valid alternative |
| Pros/cons for each | 0-5 | Are trade-offs honest? Not straw-man arguments? |
| Rationale for chosen approach | 0-5 | Is the verdict justified against the alternatives? |

### 4. Scope Definition (15 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| In-scope items are concrete | 0-5 | Each item is a deliverable, not a vague area |
| Out-of-scope explicitly listed | 0-5 | Are deferred items named, not just implied? |
| Scope is bounded | 0-5 | Can a team execute this in a defined timeframe? Or is it open-ended? |

### 5. Risk Assessment (15 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Risks identified | 0-5 | At least 3 meaningful risks, not trivial ones |
| Likelihood + impact rated | 0-5 | Is the assessment honest? Not all "low likelihood, high impact"? |
| Mitigations are actionable | 0-5 | Can someone act on the mitigation? Or is it "we'll handle it"? |

### 6. Success Criteria (15 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Criteria are measurable | 0-5 | Can you objectively verify each criterion? "Works well" is not measurable |
| Coverage is complete | 0-5 | Do criteria cover all in-scope items? Any gaps? |
| Criteria are testable | 0-5 | Could you write a test or checklist for each criterion? |

## Deduction Rules

- **Vague language penalty**: -2 per instance of "better", "improved", "enhanced" without quantification
- **Missing section penalty**: 0 points for that dimension
- **Inconsistency penalty**: -3 if scope contradicts solution or success criteria don't cover scope
