---
scale: 1100
target: 935
iterations: 3
type: contract
context:
  conventions: []
  business-rules: auto
---

# Contract Evaluation Rubric

**Total: 1100 points**

## Required Documents

| Document | Required |
|----------|----------|
| Contract files (`testing/<journey>/contracts/step-<N>-<action>.md`) | Yes |

## Scoring Dimensions

| Dimension | Points | Min Threshold |
|-----------|--------|---------------|
| 1. Completeness (完整性) | 150 | 90 |
| 2. Semantic Purity (语义纯度) | 200 | 120 |
| 3. Precondition Exclusivity (前置条件互斥性) | 150 | 90 |
| 4. Fact Alignment (事实依据) | 150 | 90 |
| 5. Surface Fitness (Surface 适配) | 100 | 60 |
| 6. Internal Consistency (一致性) | 150 | 90 |
| 7. Anchor Integrity (锚点完整性) | 100 | 60 |
| 8. Fixture Specification (前置数据声明) | 100 | 60 |
| **Total** | **1100** | **sum >= 935** |

**Pass condition**: Total score >= 935 AND every dimension >= its min threshold.

## Dimensions

### 1. Completeness (完整性) -- 150 pts

Evaluates whether all required Contract fields and sections are present, with emphasis on six-dimension structural integrity.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Every Outcome has all four mandatory dimensions non-empty (Preconditions, Input, Output, State) | 0-50 | Each Outcome block must contain non-empty values for Preconditions, Input, Output, and State. Side-effect defaults to "none" when omitted. Invariants is optional per Outcome. Missing any mandatory dimension in any Outcome is a critical failure. |
| Journey Invariants section present with at least one entry in every Contract file | 0-50 | Every Contract file must have a `## Journey Invariants` section with at least one invariant. Missing or empty invariants section is a failure. |
| Outcomes cover happy path plus required derived scenarios per surface type | 0-50 | Beyond the happy path Outcome, boundary/error Outcomes must be present as mandated by the surface type's `required_outcomes` rules. For CLI: `not-found` + `already-exists` outcomes should be considered. For API: `unauthorized` per authenticated endpoint. Missing mandatory derived Outcomes is a significant gap. |

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

### 5. Surface Fitness (Surface 适配) -- 100 pts

Evaluates whether the Contract follows the surface type's `required_outcomes` rules and uses surface-appropriate dimension descriptions.

**Dynamic adaptation**: This dimension's evaluation is parameterized by the detected surface type. Check against the corresponding surface rule file in gen-journeys skill's `rules/surface-<type>.md` (resolve relative to the gen-journeys skill directory).

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Mandatory derived Outcomes from surface `required_outcomes` are present | 0-40 | Check the surface type's `required_outcomes` rules. For CLI: `not-found` + `already-exists` must be considered. For Web: `validation-error` + `session-expired`. For API: `unauthorized` per authenticated endpoint. Score 0 if mandatory Outcomes are completely absent. |
| Dimension descriptions use surface-appropriate language | 0-35 | CLI: subprocess commands, exit codes, stdout/stderr. TUI: keyboard input, view descriptions, await semantics. Web: user interactions, page elements, async operations. API: HTTP methods, status codes, request/response bodies. Mobile: tap gestures, screen navigation, deep links. Using inappropriate surface language (e.g., DOM selectors in CLI) indicates poor adaptation. |
| TUI async Contracts include timeout Outcomes with await semantics | 0-25 | For TUI surface: any Contract with async Cmds must include a timeout Outcome with proper `await <N>ms` semantics, and the timeout Outcome must report the timed-out Cmd name. Missing timeout Outcomes for async steps is a critical failure for TUI. Non-TUI surfaces should score full points here. |

### 6. Internal Consistency (一致性) -- 150 pts

Evaluates whether Journey Invariants hold across all Step Contracts and cross-Contract references are consistent.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Invariants hold in every Step Contract | 0-60 | If the Journey declares invariants (e.g., "user is authenticated throughout", "working directory is /tmp"), verify that no Step Contract violates them. An invariant violated in any Contract is a critical failure. |
| Cross-Contract state references are consistent | 0-50 | If Contract for Step 3 references state created in Step 1 (e.g., "task created in step 1"), verify that Step 1's Contract actually describes creating that state and the reference is unambiguous. Dangling or contradictory cross-Contract references are failures. |
| Outcome Preconditions are consistent with preceding Steps' State changes | 0-40 | Step N's Outcome Preconditions should be achievable based on Step N-1's State changes. A Step claiming "task status is in_progress" as a precondition when the previous Step sets it to "completed" is an inconsistency. |

### 7. Anchor Integrity (锚点完整性) -- 100 pts

Evaluates whether Contracts contain correct technical anchor fields that link to the corresponding handbook, and whether the handbook itself is internally consistent. This dimension bridges design intent and test implementation.

**Activation condition**: This dimension applies only when the corresponding handbook file exists in the design directory (e.g., `design/api-handbook.md` for API surface). When no handbook exists, all criteria in this dimension score full marks (backward-compatible — no penalty for missing handbook).

**Anchor field mapping by surface type**:

| Surface | Required Anchor Fields | Handbook File |
|---------|----------------------|---------------|
| API | `endpoint`, `method` | `design/api-handbook.md` |
| CLI | `command` | `design/cli-handbook.md` |
| TUI | `command` | `design/cli-handbook.md` |
| Web | `page` | `design/page-map.md` |
| Mobile | `screen` | `design/screen-map.md` |

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Anchor field completeness: Contract frontmatter contains required anchor fields matching handbook entries | 0-40 | When handbook exists for the detected surface type, each Contract file must contain the required anchor fields (e.g., `endpoint` + `method` for API). Check that every endpoint/command/page/screen defined in the handbook has a corresponding Contract with matching anchor values. Missing anchor fields on any Contract when handbook is present is a deduction. Report missing fields grouped by surface type. |
| Anchor values match handbook definitions | 0-30 | Anchor field values in Contract frontmatter must match the corresponding handbook entry exactly (endpoint paths, HTTP methods, command names, page/screen identifiers). Mismatches indicate the Contract is out of sync with the design. Each mismatch is a deduction. |
| Handbook internal consistency: no duplicate or conflicting endpoint/command definitions | 0-30 | Within the handbook itself, verify no conflicting definitions exist: (1) For API: same endpoint path with different HTTP methods (method conflict) or different paths mapping to the same logical operation (path conflict). (2) For CLI: same command with conflicting subcommand definitions. (3) For Web/Mobile: same page/screen with different routes/navigation paths. Conflicts indicate design ambiguity that will propagate to Contracts. |

**Scoring guidance**:
- When handbook does not exist for the surface type: score full marks (100/100) and note in report "handbook not found — anchor integrity check skipped (backward-compatible)"
- When handbook exists but contains no entries for the surface type: score full marks and note "handbook exists but empty for this surface type"
- Missing anchor field: -10 pts per missing field (from the 0-40 criterion)
- Anchor value mismatch: -10 pts per mismatch (from the 0-30 criterion)
- Handbook conflict: -15 pts per conflict (from the 0-30 criterion)

### 8. Fixture Specification (前置数据声明) -- 100 pts

Evaluates whether the Contract's Preconditions include a complete `fixture_spec` that declares all prerequisite entities, relationships, and data constraints needed for test execution.

**Activation condition**: This dimension applies to all Contracts. For Contracts generated before this rubric version (where `fixture_spec` is absent), apply backward-compatible scoring: note in report "fixture_spec absent — backward-compatible, no penalty" and score full marks. For newly generated Contracts, `fixture_spec` is mandatory and scored normally.

**Veto item**: If the entity completeness criterion (first row) scores 0, the entire Fixture Specification dimension scores 0 regardless of other sub-scores. A Contract that does not declare the entities it operates on cannot guarantee test data sufficiency.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| **Entity completeness** (veto item) | 0-40 | `fixture_spec.entities` must include every entity type the Contract's operations interact with (create, read, update, delete, or reference). **Semantic verification required**: the reviewer must verify that each `entity_type` matches a domain model entity defined in the Design document, not merely check that the `entities` field is non-empty. Cross-reference `entity_type` values against the Design's domain model / entity-relationship diagram. Score 0 if any entity type referenced in the Contract's Preconditions, Input, or State changes is missing from `fixture_spec.entities` — this triggers the veto. |
| Relationship and constraint coverage | 0-35 | For Contracts involving parent-child or associated entities: verify that `relationship_type` and `parent_entity` are declared correctly. For entities with specific field requirements: verify `field_constraints` captures necessary constraints (e.g., status values, role types). For simple single-entity Contracts, score full marks if the single entity is properly declared. |
| Minimum data quantity declarations | 0-25 | Each entity's `min_count` must be sufficient for the Contract's test scenarios. If the Contract creates N items and then performs a list/pagination operation, `min_count` must be ≥N. If the Contract tests "delete one of many", `min_count` must be ≥2. Under-declared `min_count` that would cause test scenarios to be infeasible is a deduction. |

**Scoring guidance**:
- When `fixture_spec` is absent (legacy Contracts): score full marks (100/100) and note "fixture_spec absent — backward-compatible, no penalty"
- Entity completeness veto triggered: 0 pts for entire dimension
- Missing `relationship_type`/`parent_entity` for multi-entity Contract: -10 pts per missing relationship (from 0-35 criterion)
- `min_count` too low for test scenario: -8 pts per under-declared entity (from 0-25 criterion)

## Deduction Rules

- **Missing mandatory dimension in any Outcome**: 0 pts for Completeness
- **Hallucinated unclassified claim**: -30 pts per instance (Fact Alignment)
- **Surface type violation** (e.g., Web selectors in CLI Contract): -25 pts per instance (Surface Fitness)
- **Invariant violation in any Contract**: -40 pts per violation (Internal Consistency)
- **Precondition overlap across Outcomes**: -20 pts per ambiguous pair (Precondition Exclusivity)
- **Missing `## Journey Invariants` section**: 0 pts for Completeness Journey Invariants criterion
- **Missing anchor field when handbook exists**: -10 pts per missing field (Anchor Integrity)
- **Anchor value mismatch with handbook**: -10 pts per mismatch (Anchor Integrity)
- **Handbook endpoint/command conflict**: -15 pts per conflict (Anchor Integrity)
- **Entity completeness veto triggered** (entity_type missing from fixture_spec): 0 pts for entire Fixture Specification dimension
- **Missing relationship declaration for multi-entity Contract**: -10 pts per missing relationship (Fixture Specification)
- **Under-declared min_count**: -8 pts per under-declared entity (Fixture Specification)

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
