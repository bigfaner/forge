## Eval-journey Complete
**Final Score**: 466/1000 (target: 850)
**Iterations Used**: 1/1

### Score Progression
| Iteration | Score | Delta |
|-----------|-------|-------|
| 1         | 466   | —     |

### Dimension Breakdown (final)
| Dimension                   | Score     | Min  | Pass? |
|-----------------------------|-----------|------|-------|
| 1. Completeness             | 135/200   | 120  | YES   |
| 2. Semantic Purity          | 111/200   | 120  | NO    |
| 3. Precondition Exclusivity | 130/150   | 90   | YES   |
| 4. Fact Alignment           | 19/150    | 90   | NO    |
| 5. Surface Fitness          | 11/150    | 90   | NO    |
| 6. Internal Consistency     | 60/150    | 90   | NO    |
| **Total**                   | **466/1000** | **850** | **NO** |

### Outcome
Target NOT reached — 1 iteration exhausted. 3 of 6 dimensions below threshold (Semantic Purity, Fact Alignment, Surface Fitness, Internal Consistency).

### Key Issues
1. **Hallucinated constants**: status-migration Step 1 lists `StatusCancelled`, `StatusFailed`, `StatusReview` — actual constants are `StatusSuspended`, `StatusSkipped`, `StatusRejected`
2. **Missing CLI mandatory outcomes**: none of the 4 journeys include `not-found` or `already-exists` derived outcomes
3. **Wrong string value**: full-verification uses `"in-progress"` (hyphen) but actual value is `"in_progress"` (underscore)
4. **Implementation coupling**: steps describe code-editing actions, not CLI subprocess invocations

### Report Location
- Iteration 1 detail: `testing/eval/iteration-1.md`
