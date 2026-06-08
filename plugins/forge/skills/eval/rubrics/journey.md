---
scale: 1150
target: 975
iterations: 3
type: journey
context:
  conventions: []
  business-rules: auto
---

# Journey Evaluation Rubric

**Total: 1150 points**

## Required Documents

| Document | Required |
|----------|----------|
| Journey file (`testing/<journey>/journey.md`) | Yes |

## Scoring Dimensions

| Dimension | Points | Min Threshold |
|-----------|--------|---------------|
| 1. Completeness (完整性) | 200 | 120 |
| 2. Semantic Purity (语义纯度) | 200 | 120 |
| 3. Precondition Exclusivity (前置条件互斥性) | 150 | 90 |
| 4. Fact Alignment (事实依据) | 150 | 90 |
| 5. Surface Fitness (Surface 适配) | 150 | 90 |
| 6. Internal Consistency (一致性) | 150 | 90 |
| 7. Workflow Coverage (工作流覆盖度) | 150 | 90 |
| **Total** | **1150** | **sum ≥ 975** |

**Pass condition**: Total score ≥ 975 AND every dimension ≥ its min threshold.

## Dimensions

### 1. Completeness (完整性) — 200 pts

Evaluates whether all required Journey fields and sections are present and populated.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Journey metadata complete (name in kebab-case, risk_level as High/Medium/Low) | 0-50 | Name follows kebab-case convention. Risk level is one of the three valid values and justified by the content. |
| Steps are complete with required fields (step name, action, expected outcomes) | 0-80 | Every Step has a clear action description and at least one Outcome covering the happy path. Steps form a coherent ordered sequence. |
| Outcomes cover happy path plus required derived scenarios | 0-70 | Beyond the happy path, are boundary/error Outcomes present as required by the surface type's `required_outcomes` rules? Check that mandatory derived Outcomes (e.g., `not-found`, `already-exists` for CLI) are considered. |

### 2. Semantic Purity (语义纯度) — 200 pts

Evaluates whether all Outcome descriptions are in natural language, free of regex patterns, framework-specific assertions, or implementation details.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Outcome descriptions use natural language, not code/regex | 0-80 | No regex patterns, CSS selectors, XPath expressions, or framework assertion calls (e.g., `expect(...)`, `assertEqual`). Outcomes describe *what* the user/system observes, not *how* to verify it programmatically. |
| Preconditions are declarative statements, not procedural code | 0-60 | Preconditions describe the state or condition that must hold, not the steps to set it up. E.g., "task 42 exists and is in pending state" not "run forge task add --id 42". |
| No implementation coupling in Step descriptions | 0-60 | Steps describe user-level or system-level actions, not internal function calls, database queries, or API endpoint details. |

### 3. Precondition Exclusivity (前置条件互斥性) — 150 pts

Evaluates whether Outcomes within the same Step are distinguishable and mutually exclusive via their Preconditions.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Preconditions are distinct across Outcomes within each Step | 0-60 | For each Step, no two Outcomes share identical or semantically equivalent Preconditions. Overlapping Preconditions make it impossible to determine which Outcome applies. |
| Preconditions are sufficient to uniquely select an Outcome | 0-50 | Given the Preconditions and the current state, exactly one Outcome should be applicable. No ambiguous scenarios where multiple Outcomes could match. |
| No missing Preconditions for error/boundary Outcomes | 0-40 | Error and boundary Outcomes must state what triggers them (e.g., "resource does not exist", "invalid input provided"). Omitting Preconditions for non-happy-path Outcomes is a common failure. |

### 4. Fact Alignment (事实依据) — 150 pts

Evaluates whether claims in the Journey are grounded in verifiable facts or clearly marked as inferred reasoning.

**Two categories of declarations**:

1. **Factual claims** (事实声明): Statements based on known facts from the Fact Table. Must be traceable to a specific `fact_id`. Unknown-origin claims must be marked `UNKNOWN`.

2. **Reasonable inference claims** (合理推理声明): LLM-derived boundary Outcomes not present in the Fact Table, but generated based on surface type `required_outcomes` rules. Must include:
   - The reasoning basis (which `required_outcomes` rule triggered the derivation)
   - `source: inferred` annotation

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Factual claims traceable to fact_id or marked UNKNOWN | 0-60 | Every Outcome that asserts specific system behavior (exit codes, error messages, response formats) should reference a fact or be marked as unverified. Unverified factual claims without `UNKNOWN` marking are a failure. |
| Inferred claims have required_outcomes rule support and source: inferred | 0-50 | Derived boundary Outcomes (e.g., `not-found`, `validation-error`) must cite the `required_outcomes` rule from the surface type configuration that mandated their generation. Must be annotated `source: inferred`. |
| No hallucinated claims without classification | 0-40 | Claims that are neither factual (with traceability) nor inferred (with rule support) represent hallucinations. Zero tolerance for unclassified claims. |

### 5. Surface Fitness (Surface 适配) — 150 pts

Evaluates whether the Journey follows the surface type's `required_outcomes` rules and testing strategy proportions.

**Dynamic adaptation**: This dimension's evaluation is parameterized by the detected surface type. Check against the corresponding surface rule file in gen-journeys skill's `rules/surface-<type>.md` (resolve relative to the gen-journeys skill directory).

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Mandatory derived Outcomes from surface `required_outcomes` are present | 0-60 | Check the surface type's `required_outcomes` rules. For CLI: `not-found` + `already-exists` must be considered. For Web: `validation-error` + `session-expired`. For API: `unauthorized` per authenticated endpoint. Score 0 if mandatory Outcomes are completely absent. |
| Test strategy proportions match surface guidance | 0-50 | CLI/TUI: Contract 80% / Journey 20%. Web/API: balanced 50/50. Mobile: Journey skeleton + deep link focus. Verify that Outcome density and depth reflect these proportions. |
| Surface-specific environment and execution assumptions are realistic | 0-40 | CLI: subprocess isolation, exit code checks. TUI: terminal I/O, timeout handling. Web: browser interaction, async operations. API: HTTP methods, auth headers. Mobile: Maestro YAML, deep links. Unrealistic assumptions (e.g., asserting DOM structure for CLI) indicate poor surface adaptation. |

### 6. Internal Consistency (一致性) — 150 pts

Evaluates whether Journey Invariants hold across all Steps and cross-Step references are consistent.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Invariants hold in every Step | 0-60 | If the Journey declares invariants (e.g., "user is authenticated throughout", "working directory is /tmp"), verify that no Step violates them. An invariant violated in any Step is a critical failure. |
| Cross-Step references are consistent | 0-50 | If Step 3 references the output of Step 1 (e.g., "the task created in step 1"), verify that Step 1 indeed creates a task and the reference is unambiguous. Dangling or contradictory cross-Step references are failures. |
| Risk level is consistent with Journey content | 0-40 | A Journey marked `High` risk should involve security, data loss, or irreversible operations. A `Low` risk Journey should be read-only. Inconsistency between risk level and actual content indicates a classification error. |

### 7. Workflow Coverage (工作流覆盖度) — 150 pts

Evaluates whether the Journey includes a Golden Path and provides sufficient multi-step workflow coverage grounded in PRD/Design user stories.

**Veto item**: If the Golden Path existence criterion (first row) scores 0, the entire Workflow Coverage dimension scores 0 regardless of other sub-scores. This is a non-negotiable requirement — a Journey without a Golden Path cannot demonstrate that the feature's core workflow is testable.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| **Golden Path existence** (veto item) | 0-60 | The Journey must contain at least one Golden Path: a contiguous sequence of 3+ steps that covers the primary user story from PRD/Design. **Semantic verification required**: the reviewer must verify that the Golden Path step sequence corresponds to a specific user story or core workflow described in the PRD/Design document, not merely check that 3+ steps exist. Steps must reference domain-level user operations (e.g., "create a milestone", "transition task status"), not bare HTTP methods or API calls. Score 0 if no Golden Path exists or if the step sequence is semantically unrelated to any user story — this triggers the veto. |
| Multi-step coverage depth | 0-50 | Beyond the Golden Path, evaluate whether the Journey covers additional meaningful workflow variations: state transitions, entity lifecycle operations (create → update → delete), cross-entity interactions (parent-child relationship operations). Shallow coverage (only CRUD on a single entity) scores low. Deep coverage (multi-entity workflows, state machines, error recovery paths) scores high. |
| Workflow completeness against PRD/Design scope | 0-40 | Cross-reference the Journey's step coverage against the PRD/Design document's described user workflows. Identify any user-facing workflow mentioned in PRD/Design that has no corresponding Journey coverage. The reviewer should list uncovered workflows as gaps. Minor auxiliary workflows may be absent without penalty, but primary workflows must be covered. |

## Deduction Rules

- **Missing required field/section**: 0 pts for the affected dimension
- **Hallucinated unclassified claim**: -30 pts per instance (Fact Alignment)
- **Surface type violation** (e.g., Web assertions in CLI Journey): -25 pts per instance (Surface Fitness)
- **Invariant violation**: -40 pts per violation (Internal Consistency)
- **Precondition overlap across Outcomes**: -20 pts per ambiguous pair (Precondition Exclusivity)
- **Golden Path veto triggered** (Golden Path scores 0): 0 pts for entire Workflow Coverage dimension regardless of other sub-scores
- **Golden Path steps use API-level descriptions instead of domain-level user operations**: -15 pts per step (Workflow Coverage)

## Eval Failure Handling

When the scorer cannot parse the document (e.g., malformed Journey file):
1. Retry scoring once with a simplified parsing prompt
2. If still failing, mark the Journey as `eval-skipped` with confidence `LOW`
3. Write `testing/<journey>/.eval-status.json`:
   ```json
   {"status": "eval-skipped", "confidence": "LOW", "reason": "<parse failure reason>"}
   ```
4. Record in eval report: parse failure reason, raw output snippet

When eval fails after all iterations:
1. Output the final score and per-dimension breakdown
2. List all dimensions below threshold with specific gaps
3. Return the evaluation results for pipeline orchestration (PAUSE_J handling is in gen-journeys SKILL.md, not here)
