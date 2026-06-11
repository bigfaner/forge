---
date: "2026-05-14"
doc_dir: "docs/features/typed-verification-strategies/design/"
iteration: 2
target_score: "900"
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration 2

**Score: 845/1000** (target: 900)

```
+---------------------------------------------------------------+
|                     DESIGN QUALITY SCORECARD                   |
+------------------------------+----------+----------+----------+
| Dimension                    | Score    | Max      | Status   |
+------------------------------+----------+----------+----------+
| 1. Architecture Clarity      |  190     |  200     | PASS     |
|    Layer placement explicit  |  70/70   |          |          |
|    Component diagram present |  65/70   |          |          |
|    Dependencies listed       |  55/60   |          |          |
+------------------------------+----------+----------+----------+
| 2. Interface & Model Defs    |  180     |  200     | PASS     |
|    Interface signatures typed|  65/70   |          |          |
|    Models concrete           |  65/70   |          |          |
|    Directly implementable    |  50/60   |          |          |
+------------------------------+----------+----------+----------+
| 3. Error Handling            |  130     |  150     | WARNING  |
|    Error types defined       |  45/50   |          |          |
|    Propagation strategy clear|  45/50   |          |          |
|    HTTP status codes mapped  |  40/50   |          |          |
+------------------------------+----------+----------+----------+
| 4. Testing Strategy          |  100     |  150     | WARNING  |
|    Per-layer test plan       |  35/50   |          |          |
|    Coverage target numeric   |  35/50   |          |          |
|    Test tooling named        |  30/50   |          |          |
+------------------------------+----------+----------+----------+
| 5. Breakdown-Readiness       |  185     |  200     | PASS     |
|    Components enumerable     |  65/70   |          |          |
|    Tasks derivable           |  65/70   |          |          |
|    PRD AC coverage           |  55/60   |          |          |
+------------------------------+----------+----------+----------+
| 6. Security Considerations   |  60      |  100     | WARNING  |
|    Threat model present      |  30/50   |          |          |
|    Mitigations concrete      |  30/50   |          |          |
+------------------------------+----------+----------+----------+
| TOTAL                        |  845     |  1000    |          |
+------------------------------+----------+----------+----------+
```

Breakdown-Readiness >= 180/200 — can proceed to `/breakdown-tasks`

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Architecture:Component Diagram | Mermaid diagram is an improvement but mixes process flow steps (Step 2.7, Step 3) with component boxes, making it a hybrid that conflates structure and flow | -5 pts |
| Architecture:Dependencies | Missing dependency on existing SKILL.md files being modified (gen-test-cases/SKILL.md, gen-test-scripts/SKILL.md) as the host system; only lists new files and CLI tools | -5 pts |
| Interface:Signatures | Interface 1 (Strategy File Format) validation rules are prose ("Each ## section must contain ### Verification Dimensions >=3 items") rather than a formal schema or type definition — though the concrete example compensates significantly | -5 pts |
| Interface:Models | `StrategyFile` model uses `Map<string, VerificationStrategy>` which implies a runtime map, but this is a Skills/LLM layer — the map is actually just sections in a Markdown file. The model conflates runtime data structures with file structure | -5 pts |
| Interface:Implementable | The `LevelTag` model says "Derivation rule: capability -> (level, interface) from the static lookup table" but the lookup table (Interface 2) only covers 5 capabilities. What happens for capabilities not in the table (e.g., a future `grpc` or `graphql` capability)? No fallback rule specified | -10 pts |
| Error:Types | Error codes (STRATEGY_NOT_FOUND, STRATEGY_INVALID, etc.) are well-defined but no error type structure or grouping is provided — are these codes enum values? LLM prompt directives? Constants in a file? The format is clear but the representation mechanism is ambiguous | -5 pts |
| Error:Propagation | The propagation strategy statement is clear ("fail-fast on validation errors, best-effort on generation warnings") but the Error Types table includes GOLDEN_STALE as "Info" with "manual judgment" — this is a different layer (test execution, not generation) and its placement in the gen-test-cases error table is inconsistent | -5 pts |
| Error:HTTP Status | HTTP status codes: N/A for Skills layer, but GOLDEN_STALE ("Test failure + diff, manual judgment") is semantically an exit code or test result, not an error type. Its inclusion muddles the error model. | -10 pts |
| Testing:Per-layer | Test plan now has 3 rows (Strategy File Format, Feature E2E Pipeline, Feature E2E Eval) but "Strategy File Format" test type is "Manual LLM invocation" — no definition of what that means mechanically. Is it prompt comparison? Golden prompt? Human review? | -15 pts |
| Testing:Coverage | Coverage target is "8/8 PRD ACs passing" and "3/3 format scenarios" and "New dimension >= 150 pts" — these are pass/fail gates, not coverage targets. No percentage-based coverage metric is provided for any layer | -15 pts |
| Testing:Tooling | Test tooling names "Manual LLM invocation with crafted strategy files" and "forge gen-test-cases + forge gen-test-scripts CLI" and "forge eval-test-cases CLI" — these are the system under test, not test tools/frameworks. No assertion library, test runner, or comparison tool is named | -20 pts |
| Breakdown:Components | Components are enumerable but the "6 profiles" are listed as one task ("Create verification-strategies.md for each of 6 profiles") rather than individually tracked | -5 pts |
| Breakdown:Tasks | Interface-to-task mapping table is excellent, but 3 of the 8 rows map to the same location (`plugins/forge/skills/gen-test-cases/SKILL.md`). While technically correct, this conflates multiple tasks into a single location without specifying which step or section of SKILL.md corresponds to which task | -5 pts |
| Breakdown:PRD AC | S8 coverage is now detailed with scoring formula and 3 check items. However, S2 AC specifies "Level field coverage rate >= 95%" — the design maps this to "LevelTag.level" but does not explain how 95% coverage is measured or enforced. The >=95% metric from PRD is not addressed in any test scenario | -5 pts |
| Security:Threat Model | Threat model is one sentence: "Strategy files are authored by profile maintainers and read as LLM prompt context. No user-supplied input flows into strategy files." This is still thin — it does not address: (1) compromised profile repos, (2) malicious strategy content that manipulates LLM behavior (prompt injection), (3) strategy content influencing test generation to skip security tests | -20 pts |
| Security:Mitigations | Four bullet points that essentially say the same thing 4 ways: no shell injection, no eval, validated format, PRD says no shell injection. No mitigation for the actual novel threat: malicious strategy content affecting LLM test generation quality or causing the LLM to generate insecure test code | -20 pts |

---

## Attack Points

### Attack 1: Testing Strategy — test tooling is SUT, not test tooling

**Where**: Testing Strategy table, Tool column: "Manual LLM invocation with crafted strategy files (valid, malformed, key-mismatched)", "forge gen-test-cases + forge gen-test-scripts CLI with go-test profile"
**Why it's weak**: The "Tool" column names the system under test, not the test execution tooling. "Manual LLM invocation" is not a test tool — it is a manual process with no defined assertion mechanism. How does the tester know whether gen-test-cases "correctly accepts valid strategy files"? By reading the output? By running a diff? By checking error codes? No test runner, assertion library, or automated comparison tool is named. For a Skills-layer feature, the test tooling should specify: (1) what invocation mechanism runs the test (CLI, script, harness), (2) how pass/fail is determined (golden diff, exit code, rubric score threshold), and (3) what automates the comparison. "Manual LLM invocation" tells a developer nothing about how to reproduce or automate the test.
**What must improve**: Replace "Manual LLM invocation" with a concrete test mechanism (e.g., "Shell script invoking `forge gen-test-cases` with crafted profile, asserting exit code 0/1 and stderr content against expected error messages"). Name the comparison/assertion tool (e.g., `diff`, `grep`, test harness script).

### Attack 2: Security — mitigations are circular, threat model is superficial

**Where**: Security Considerations section: "Strategy content is LLM prompt context only -- no shell command injection path", "No eval or command interpolation from strategy file content", "Input validation: strategy file format and key correctness are validated before use as generation instructions", "Per PRD security requirement: strategy content is never spliced into shell commands or eval"
**Why it's weak**: All four mitigations address the same non-threat (shell injection) while ignoring the primary threat surface: LLM prompt injection via malicious strategy content. A compromised or malicious verification-strategies.md could instruct the LLM to generate test code that (1) skips security-critical test scenarios, (2) introduces false-positive tests that mask real bugs, or (3) includes insecure code patterns in generated tests. The threat model dismisses this with "No user-supplied input flows into strategy files" — but profile repos may be shared, vendored, or community-sourced. The PRD's security requirement ("strategy content only as LLM prompt context, no shell injection") addresses shell injection only, which is a low-risk path. The higher-risk path (LLM output manipulation) is unaddressed.
**What must improve**: (1) Add a threat for "malicious strategy content influencing LLM test generation" with a specific mitigation (e.g., strategy file content review checklist, template constraints that limit what strategies can request, or output validation that checks generated tests for common insecure patterns). (2) Differentiate the mitigations — currently all 4 say the same thing.

### Attack 3: Testing Strategy — coverage targets are pass/fail gates, not coverage metrics

**Where**: Testing Strategy, Coverage Target column: "3/3 format scenarios", "All 8 ACs", "New dimension >= 150 pts"
**Why it's weak**: The rubric requires "a numeric coverage target (e.g., 80%)". The design provides pass/fail gates: "3/3 format scenarios" means "all 3 must pass", "All 8 ACs" means "all 8 must pass". These are acceptance gates, not coverage metrics. A coverage metric measures the proportion of a space that is tested (e.g., "80% of strategy file edge cases covered", "95% of capability types tested"). The PRD itself has a quantitative target ("Level field coverage rate >= 95%") that could serve as a coverage target, but the design does not adopt it. The "Overall Coverage Target" is stated as "8/8 PRD ACs passing" which is binary (all or nothing), not a coverage percentage.
**What must improve**: Add at least one percentage-based coverage target. For example: "Level field auto-tagging accuracy >= 95% (matching PRD S2 metric)" or "Strategy file format edge case coverage >= 80% (including empty sections, malformed markdown, extra subsections)".

### Attack 4: Interface & Model — LevelTag has no fallback for unknown capabilities

**Where**: Data Models section, LevelTag model: "Derivation rule: capability -> (level, interface) from the static lookup table in Interface 2. No capability or no strategy -> both values empty."
**Why it's weak**: The lookup table in Interface 2 defines exactly 5 mappings (tui, web-ui, mobile-ui, api, cli). The design states "No capability or no strategy -> both values empty" as the fallback, but this means any capability not in the table silently produces empty Level/Interface. This is inconsistent with the error handling section's `UNKNOWN_CAPABILITY` warning, which says "warning + generate generic TC (no Level)". The model definition says "both values empty" without referencing the UNKNOWN_CAPABILITY warning path. A developer implementing this must reconcile two different descriptions of the same fallback behavior. Additionally, the PRD's monitoring requirement ("test-cases.md header comment includes strategy metadata") is not reflected in any model or interface definition.
**What must improve**: (1) Explicitly state the fallback derivation for unknown capabilities in the LevelTag model, cross-referencing the UNKNOWN_CAPABILITY error. (2) Add the PRD's monitoring metadata (Applied strategy comment) to either the test-cases.md interface or a model.

### Attack 5: Error Handling — GOLDEN_STALE is out of scope for this pipeline stage

**Where**: Error Types table, GOLDEN_STALE row: "gen-test-scripts: golden file mismatch for {test-id}. Diff: {path}" with behavior "Test failure + diff, manual judgment"
**Why it's weak**: GOLDEN_STALE is a test-execution-time error, not a generation-time error. It occurs when tests are run against golden files, not during gen-test-cases or gen-test-scripts generation. Including it in the error handling table for a design that covers test generation (not test execution) creates confusion about the error boundary. The error code is also prefixed with `gen-test-scripts` in its message, implying it occurs during script generation, but the described behavior (test failure, diff) only occurs during test execution. This is a cross-section inconsistency between the error table and the feature scope.
**What must improve**: Either remove GOLDEN_STALE from the error table (since test execution is out of scope per PRD) or clearly mark it as a runtime error that occurs outside the generation pipeline, with a note that it is documented here for completeness only.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: S8 eval-test-cases AC incompletely addressed | YES | PRD Coverage Map now includes full S8 entry: "Scoring formula: score = (covered_items / total_items) * weight. Check items: (1) dimension coverage, (2) boundary scenario coverage, (3) test data compliance. Fallback: when no strategy file exists, dimension score rate >= 0.8." |
| Attack 2: Test tooling unexplained and incomplete | PARTIAL | Opaque "T-test-2, T-test-3" references are gone. Replaced with "Manual LLM invocation with crafted strategy files" and "forge gen-test-cases + forge gen-test-scripts CLI". Still not actual test tools — these are SUT invocations, not test tooling. |
| Attack 3: No concrete strategy file example | YES | Full concrete example added for go-test profile TUI capability (lines 99-120): 4 verification dimensions, 5 boundary scenarios, 2 test data requirements. |
| Attack 4: No formal error types or propagation strategy | YES | Error codes defined (STRATEGY_NOT_FOUND, STRATEGY_INVALID, KEY_MISMATCH, UNKNOWN_CAPABILITY, GOLDEN_STALE). Propagation strategy explicitly stated: "fail-fast on validation errors, best-effort on generation warnings". Error messages now follow `<context>: <detail>` format with `gen-test-cases:` prefix. |
| Attack 5: No explicit interface-to-task mapping | YES | New "Interface-to-Task Mapping" table added with 8 rows mapping each interface/model/error to a concrete implementation task and file location. |
| Attack 6: Component diagram is informal ASCII | YES | Replaced with proper Mermaid flowchart with subgraphs, labeled nodes, and clear data flow arrows. |

---

## Verdict

- **Score**: 845/1000
- **Target**: 900/1000
- **Gap**: 55 points
- **Breakdown-Readiness**: 185/200 — can proceed to `/breakdown-tasks`
- **Action**: Continue to iteration 3. Priority fixes: (1) Name concrete test tools and assertion mechanisms, (2) Add LLM prompt injection threat + mitigation to security, (3) Add percentage-based coverage targets, (4) Remove or scope GOLDEN_STALE error, (5) Add LevelTag fallback rule for unknown capabilities.
