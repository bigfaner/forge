---
date: "2026-05-14"
doc_dir: "docs/features/typed-verification-strategies/design/"
iteration: 1
target_score: "900"
evaluator: Claude (automated, adversarial)
---

# Design Eval -- Iteration 1

**Score: 780/1000** (target: 900)

```
+---------------------------------------------------------------+
|                     DESIGN QUALITY SCORECARD                   |
+------------------------------+----------+----------+----------+
| Dimension                    | Score    | Max      | Status   |
+------------------------------+----------+----------+----------+
| 1. Architecture Clarity      |  170     |  200     | WARNING  |
|    Layer placement explicit  |  65/70   |          |          |
|    Component diagram present |  55/70   |          |          |
|    Dependencies listed       |  50/60   |          |          |
+------------------------------+----------+----------+----------+
| 2. Interface & Model Defs    |  160     |  200     | WARNING  |
|    Interface signatures typed|  55/70   |          |          |
|    Models concrete           |  60/70   |          |          |
|    Directly implementable    |  45/60   |          |          |
+------------------------------+----------+----------+----------+
| 3. Error Handling            |  120     |  150     | WARNING  |
|    Error types defined       |  35/50   |          |          |
|    Propagation strategy clear|  35/50   |          |          |
|    HTTP status codes mapped  |  50/50   |          | N/A      |
+------------------------------+----------+----------+----------+
| 4. Testing Strategy          |  100     |  150     | WARNING  |
|    Per-layer test plan       |  30/50   |          |          |
|    Coverage target numeric   |  45/50   |          |          |
|    Test tooling named        |  25/50   |          |          |
+------------------------------+----------+----------+----------+
| 5. Breakdown-Readiness       |  155     |  200     | FAIL     |
|    Components enumerable     |  60/70   |          |          |
|    Tasks derivable           |  50/70   |          |          |
|    PRD AC coverage           |  45/60   |          |          |
+------------------------------+----------+----------+----------+
| 6. Security Considerations   |  75      |  100     | WARNING  |
|    Threat model present      |  35/50   |          |          |
|    Mitigations concrete      |  40/50   |          |          |
+------------------------------+----------+----------+----------+
| TOTAL                        |  780     |  1000    |          |
+------------------------------+----------+----------+----------+
```

Breakdown-Readiness < 180/200 blocks progression to `/breakdown-tasks`

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Architecture:Component Diagram | ASCII diagram is informal, mixes flow with structure; PRD uses Mermaid but design does not | -15 pts |
| Architecture:Dependencies | Missing dependency on gen-test-cases/gen-test-scripts skill system itself as the host for modifications | -10 pts |
| Interface:Signatures | Interface 1 (strategy file format) uses prose validation rules instead of formal schema definition | -15 pts |
| Interface:Implementable | No filled example of verification-strategies.md -- template is structural but has no concrete instance | -15 pts |
| Error:Types | No formal error codes or error type constants defined -- scenarios described in prose only | -15 pts |
| Error:Propagation | No explicit statement of overall error propagation strategy (fail-fast vs best-effort); only per-scenario behaviors listed | -15 pts |
| Testing:Per-layer | Test plan has only "Feature E2E" rows; no separate test layer for strategy file format validation or profile system testing | -20 pts |
| Testing:Tooling | "T-test-2, T-test-3" referenced without explanation; no test frameworks or libraries named for the skill layer | -25 pts |
| Breakdown:Tasks | Interface-to-task mapping is implicit -- no explicit table mapping each interface to one or more implementation tasks | -20 pts |
| Breakdown:PRD AC | S8 AC specifies scoring formula (covered items / total items * weight) and 3 check items; design only says "New dimension (150 pts)" without formula or checks | -15 pts |
| Security:Threat Model | Threat model is very thin -- does not consider compromised profile repos or malicious strategy content affecting LLM behavior | -15 pts |
| Security:Mitigations | Mitigations restate the same point 4 ways (no shell injection); no mitigation for LLM prompt injection via malicious strategy content | -10 pts |

---

## Attack Points

### Attack 1: Breakdown-Readiness -- S8 eval-test-cases AC incompletely addressed

**Where**: PRD Coverage Map line: `| S8: Typed verification eval | eval-test-cases rubric | New dimension (150 pts) |`
**Why it's weak**: PRD User Story 8 acceptance criteria specifies a concrete scoring formula: "completeness score = covered items / total items * weight score" and three explicit check items: (1) dimension coverage >=3/3, (2) boundary scenarios >=2 from strategy, (3) test data satisfies strategy requirements. The design maps S8 to "New dimension (150 pts)" with zero detail on the scoring formula, check items, or how eval-test-cases should compute typed verification completeness. This is a gap between the AC specificity and the design response.
**What must improve**: Add a subsection under eval-test-cases changes that defines the scoring formula, the three check items from S8 AC, and the fallback scoring (>=0.8 when no strategy file exists).

### Attack 2: Testing Strategy -- test tooling unexplained and incomplete

**Where**: Testing Strategy table line: `| Feature E2E | Pipeline test | T-test-2, T-test-3 | gen-test-cases produces Level/Interface fields; gen-test-scripts produces differentiated code | All 8 ACs |`
**Why it's weak**: "T-test-2, T-test-3" are opaque identifiers with no explanation. Are they test scripts? Test harness tools? Test plan references? No test framework or library is named. For a Skills-layer feature, the test tooling should reference the specific mechanism (e.g., manual LLM invocation, prompt comparison, golden file diff) rather than unexplained abbreviations. Additionally, the per-layer test plan has no entry for validating the strategy file format parsing or the profile system interaction.
**What must improve**: (1) Explain what T-test-2 and T-test-3 are, or replace with actual test tool names. (2) Add a test row for strategy file format validation. (3) Name the concrete test execution mechanism for the Skills layer.

### Attack 3: Interface & Model -- no concrete strategy file example

**Where**: Interface 1 (lines 70-81): The strategy file format template shows `## <capability-key>` with `### Verification Dimensions` and `### Boundary Scenarios` but contains only placeholder text like "Dimension description (>=3 required)".
**Why it's weak**: A developer implementing the go-test profile's verification-strategies.md has no concrete example to follow. The PRD's type verification matrix (TUI golden file + dimension checks + CJK boundary scenarios) provides the content, but the design does not bridge from that matrix to an actual filled-in file. Without a concrete instance, the `dimensions: string[]` model field is abstract -- what does a real dimension look like? "Golden file comparison" or "Verify output matches golden file at {path} for terminal dimensions {W}x{H}"?
**What must improve**: Add at least one fully filled verification-strategies.md example for a real capability (e.g., TUI for go-test profile) showing actual dimension descriptions, boundary scenarios, and test data requirements.

### Attack 4: Error Handling -- no formal error types or error propagation strategy

**Where**: Error Handling section (lines 158-167): Error scenarios listed in prose table with message templates but no error codes, type names, or formal propagation strategy.
**Why it's weak**: The error table lists 5 scenarios with message templates and behaviors (Warning/Abort), but there are no error type definitions (e.g., `ErrStrategyNotFound`, `ErrStrategyInvalid`, `ErrKeyMismatch`). The `<context>: <detail>` format is referenced from a convention doc but not demonstrated in the actual error messages. The messages like `Strategy file validation failed for profile: {name}. Section {key} missing required subsections...` do not follow the `<context>: <detail>` format -- they use a different pattern. Additionally, there is no explicit statement of the overall error propagation philosophy (fail-fast? collect-all-errors? partial-continue?).
**What must improve**: (1) Define named error types or codes for each error scenario. (2) Add an explicit error propagation strategy statement (e.g., "gen-test-cases uses fail-fast on validation errors, best-effort on generation warnings"). (3) Verify error messages conform to the referenced `<context>: <detail>` convention.

### Attack 5: Breakdown-Readiness -- no explicit interface-to-task mapping

**Where**: PRD Coverage Map (lines 221-230) maps PRD requirements to design components, but no equivalent table maps interfaces to implementation tasks.
**Why it's weak**: The rubric requires "Does each interface -> at least one impl task? Each model -> at least one schema task?" The design defines 3 interfaces and 3 models but never explicitly states "Interface 1 (Strategy File Format) -> Task: Create verification-strategies.md for each of 6 profiles" or "Interface 3 (Level-based Code Gen) -> Task: Modify gen-test-scripts SKILL.md Step 4". The mapping is derivable by a careful reader but not explicitly provided, making breakdown more error-prone.
**What must improve**: Add an explicit Interface-to-Task mapping table: each interface and model should list the concrete implementation task(s) it requires.

### Attack 6: Architecture -- component diagram is informal

**Where**: Component Diagram (lines 28-52): ASCII text diagram showing Profile System, gen-test-cases, gen-test-scripts, and eval-test-cases.
**Why it's weak**: The PRD itself uses a proper Mermaid flowchart (lines 71-95 in prd-spec.md), establishing an expectation for diagram quality. The design's ASCII diagram is informal: arrows use mixed styles (`|`, `v`, `------->`), layout is ragged, and it mixes component structure with step numbers and flow direction. It is not immediately clear which boxes are components vs. process steps. A Mermaid diagram (flowchart or C4-style) would be clearer and consistent with the PRD's convention.
**What must improve**: Replace the ASCII diagram with a Mermaid diagram (e.g., `flowchart LR` or component diagram) that cleanly separates components from process flow. Consider using subgraphs to group related components.

---

## Previous Issues Check

N/A -- first iteration.

---

## Verdict

- **Score**: 780/1000
- **Target**: 900/1000
- **Gap**: 120 points
- **Breakdown-Readiness**: 155/200 -- cannot proceed to `/breakdown-tasks`
- **Action**: Continue to iteration 2. Priority fixes: (1) Flesh out S8 eval scoring formula and check items, (2) Add concrete strategy file example, (3) Clarify test tooling, (4) Add interface-to-task mapping table, (5) Define formal error types and propagation strategy.
