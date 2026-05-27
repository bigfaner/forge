---
scale: 1000
target: 850
iterations: 3
type: contract
context:
  conventions: []
  business-rules: auto
---

# Contract Evaluation Rubric

**Total: 1000 points**

## Required Documents

| Document | Required |
|----------|----------|
| Contract files (`testing/<journey>/contracts/step-<N>-<action>.md`) | Yes |

## Scoring Dimensions

| Dimension | Points | Min Threshold |
|-----------|--------|---------------|
| 1. Completeness (完整性) | 200 | 120 |
| 2. Semantic Purity (语义纯度) | 200 | 120 |
| 3. Precondition Exclusivity (前置条件互斥性) | 150 | 90 |
| 4. Fact Alignment (事实依据) | 150 | 90 |
| 5. Surface Fitness (Surface 适配) | 150 | 90 |
| 6. Internal Consistency (一致性) | 150 | 90 |
| **Total** | **1000** | **sum >= 850** |

**Pass condition**: Total score >= 850 AND every dimension >= its min threshold.

## Dimensions

### 1. Completeness (完整性) -- 200 pts

Evaluates whether all required Contract fields and sections are present, with emphasis on six-dimension structural integrity.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Every Outcome has all four mandatory dimensions non-empty (Preconditions, Input, Output, State) | 0-70 | Each Outcome block must contain non-empty values for Preconditions, Input, Output, and State. Side-effect defaults to "none" when omitted. Invariants is optional per Outcome. Missing any mandatory dimension in any Outcome is a critical failure. |
| Journey Invariants section present with at least one entry in every Contract file | 0-60 | Every Contract file must have a `## Journey Invariants` section with at least one invariant. Missing or empty invariants section is a failure. |
| Outcomes cover happy path plus required derived scenarios per surface type | 0-70 | Beyond the happy path Outcome, boundary/error Outcomes must be present as mandated by the surface type's `required_outcomes` rules. For CLI: `not-found` + `already-exists` outcomes should be considered. For API: `unauthorized` per authenticated endpoint. Missing mandatory derived Outcomes is a significant gap. |

### 2. Semantic Purity (语义纯度) -- 200 pts

Evaluates whether all dimension values are natural language descriptions, free of regex patterns, framework-specific assertions, or implementation details.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Dimension values use natural language, not code/regex | 0-80 | No regex patterns (`\d`, `.*`, `[...]`, `(?:...)`, `\s`, `\w`, `\b`, `$`, `^`), CSS selectors, XPath expressions, or framework assertion calls. Dimension values describe *what* the system produces, not *how* to verify it. |
| Preconditions are declarative state descriptions, not procedural steps | 0-60 | Preconditions describe the system state that must hold, not setup instructions. E.g., "task 42 exists and is in pending state" not "run forge task add --id 42". |
| No implementation coupling in dimension values | 0-60 | Dimension values describe system-level behavior, not internal function calls, database queries, API endpoint paths, or file system paths. E.g., "success confirmation containing feature-slug" not "API returns 201 with {slug: ...}". |

### 3. Precondition Exclusivity (前置条件互斥性) -- 150 pts

Evaluates whether Outcomes within the same Step are distinguishable and mutually exclusive via their Preconditions. This is especially critical for Contracts because multiple Outcomes per Step must be unambiguously selectable.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Preconditions are distinct across Outcomes within each Step | 0-60 | For each Step's Contract, no two Outcomes share identical or semantically equivalent Preconditions. Overlapping Preconditions make it impossible to determine which Outcome applies. |
| Preconditions are sufficient to uniquely select an Outcome | 0-50 | Given the Preconditions and the current system state, exactly one Outcome should be applicable. No ambiguous scenarios where multiple Outcomes could match. |
| Error/boundary Outcomes state their triggering conditions explicitly | 0-40 | Error and boundary Outcomes must state what triggers them (e.g., "resource does not exist", "invalid input provided", "task is not in in_progress status"). Omitting Preconditions for non-happy-path Outcomes is a common failure. |

### 4. Fact Alignment (事实依据) -- 150 pts

Evaluates whether claims in Contract dimensions are grounded in verifiable facts from the Fact Table or clearly marked as inferred reasoning.

**Two categories of declarations**:

1. **Factual claims** (事实声明): Statements based on known facts from the Fact Table. Must be traceable to a specific `fact_id`. Unknown-origin claims must be marked `UNKNOWN`.

2. **Reasonable inference claims** (合理推理声明): LLM-derived boundary Outcomes not present in the Fact Table, but generated based on surface type `required_outcomes` rules. Must include:
   - The reasoning basis (which `required_outcomes` rule triggered the derivation)
   - `source: inferred` annotation

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Factual claims traceable to fact_id or marked UNKNOWN | 0-60 | Every Outcome that asserts specific system behavior (exit codes, error messages, response formats, state changes) should reference a fact or be marked as unverified. Unverified factual claims without `UNKNOWN` marking are a failure. |
| Inferred claims have required_outcomes rule support and source: inferred | 0-50 | Derived boundary Outcomes (e.g., `not-found`, `validation-error`) must cite the `required_outcomes` rule from the surface type configuration that mandated their generation. Must be annotated `source: inferred`. |
| No hallucinated claims without classification | 0-40 | Claims that are neither factual (with traceability) nor inferred (with rule support) represent hallucinations. Zero tolerance for unclassified claims. |

### 5. Surface Fitness (Surface 适配) -- 150 pts

Evaluates whether the Contract follows the surface type's `required_outcomes` rules and uses surface-appropriate dimension descriptions.

**Dynamic adaptation**: This dimension's evaluation is parameterized by the detected surface type. Check against the corresponding surface rule file in `skills/gen-journeys/rules/surface-<type>.md`.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Mandatory derived Outcomes from surface `required_outcomes` are present | 0-60 | Check the surface type's `required_outcomes` rules. For CLI: `not-found` + `already-exists` must be considered. For Web: `validation-error` + `session-expired`. For API: `unauthorized` per authenticated endpoint. Score 0 if mandatory Outcomes are completely absent. |
| Dimension descriptions use surface-appropriate language | 0-50 | CLI: subprocess commands, exit codes, stdout/stderr. TUI: keyboard input, view descriptions, await semantics. Web: user interactions, page elements, async operations. API: HTTP methods, status codes, request/response bodies. Mobile: tap gestures, screen navigation, deep links. Using inappropriate surface language (e.g., DOM selectors in CLI) indicates poor adaptation. |
| TUI async Contracts include timeout Outcomes with await semantics | 0-40 | For TUI surface: any Contract with async Cmds must include a timeout Outcome with proper `await <N>ms` semantics, and the timeout Outcome must report the timed-out Cmd name. Missing timeout Outcomes for async steps is a critical failure for TUI. Non-TUI surfaces should score full points here. |

### 6. Internal Consistency (一致性) -- 150 pts

Evaluates whether Journey Invariants hold across all Step Contracts and cross-Contract references are consistent.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Invariants hold in every Step Contract | 0-60 | If the Journey declares invariants (e.g., "user is authenticated throughout", "working directory is /tmp"), verify that no Step Contract violates them. An invariant violated in any Contract is a critical failure. |
| Cross-Contract state references are consistent | 0-50 | If Contract for Step 3 references state created in Step 1 (e.g., "task created in step 1"), verify that Step 1's Contract actually describes creating that state and the reference is unambiguous. Dangling or contradictory cross-Contract references are failures. |
| Outcome Preconditions are consistent with preceding Steps' State changes | 0-40 | Step N's Outcome Preconditions should be achievable based on Step N-1's State changes. A Step claiming "task status is in_progress" as a precondition when the previous Step sets it to "completed" is an inconsistency. |

## Deduction Rules

- **Missing mandatory dimension in any Outcome**: 0 pts for Completeness
- **Hallucinated unclassified claim**: -30 pts per instance (Fact Alignment)
- **Surface type violation** (e.g., Web selectors in CLI Contract): -25 pts per instance (Surface Fitness)
- **Invariant violation in any Contract**: -40 pts per violation (Internal Consistency)
- **Precondition overlap across Outcomes**: -20 pts per ambiguous pair (Precondition Exclusivity)
- **Missing `## Journey Invariants` section**: 0 pts for Completeness Journey Invariants criterion

## Eval Failure Handling

When the scorer cannot parse the document (e.g., malformed Contract file):
1. Retry scoring once with a simplified parsing prompt
2. If still failing, mark the Contract as `eval-skipped` with confidence `LOW`
3. Write `testing/<journey>/.eval-status.json`:
   ```json
   {"status": "eval-skipped", "confidence": "LOW", "reason": "<parse failure reason>"}
   ```
4. Record in eval report: parse failure reason, raw output snippet

When eval fails after all iterations:
1. Output the final score and per-dimension breakdown
2. List all dimensions below threshold with specific gaps
3. Return the evaluation results for pipeline orchestration (PAUSE_C handling is in gen-contracts SKILL.md, not here)
