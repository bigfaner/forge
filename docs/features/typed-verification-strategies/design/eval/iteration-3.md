---
date: "2026-05-14"
doc_dir: "docs/features/typed-verification-strategies/design/"
iteration: 3
target_score: "900"
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration 3

**Score: 883/1000** (target: 900)

```
+---------------------------------------------------------------+
|                     DESIGN QUALITY SCORECARD                   |
+------------------------------+----------+----------+----------+
| Dimension                    | Score    | Max      | Status   |
+------------------------------+----------+----------+----------+
| 1. Architecture Clarity      |  175     |  200     | PASS     |
|    Layer placement explicit  |  65/70   |          |          |
|    Component diagram present |  60/70   |          |          |
|    Dependencies listed       |  50/60   |          |          |
+------------------------------+----------+----------+----------+
| 2. Interface & Model Defs    |  174     |  200     | PASS     |
|    Interface signatures typed|  62/70   |          |          |
|    Models concrete           |  62/70   |          |          |
|    Directly implementable    |  50/60   |          |          |
+------------------------------+----------+----------+----------+
| 3. Error Handling            |  146     |  150     | PASS     |
|    Error types defined       |  48/50   |          |          |
|    Propagation strategy clear|  48/50   |          |          |
|    HTTP status codes mapped  |  50/50   |          |          | (N/A)
+------------------------------+----------+----------+----------+
| 4. Testing Strategy          |  125     |  150     | WARNING  |
|    Per-layer test plan       |  40/50   |          |          |
|    Coverage target numeric   |  45/50   |          |          |
|    Test tooling named        |  40/50   |          |          |
+------------------------------+----------+----------+----------+
| 5. Breakdown-Readiness       |  181     |  200     | PASS     |
|    Components enumerable     |  65/70   |          |          |
|    Tasks derivable           |  63/70   |          |          |
|    PRD AC coverage           |  53/60   |          |          |
+------------------------------+----------+----------+----------+
| 6. Security Considerations   |  82      |  100     | WARNING  |
|    Threat model present      |  42/50   |          |          |
|    Mitigations concrete      |  40/50   |          |          |
+------------------------------+----------+----------+----------+
| TOTAL                        |  883     |  1000    |          |
+------------------------------+----------+----------+----------+
```

Breakdown-Readiness >= 180/200 — can proceed to `/breakdown-tasks`

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Architecture:Component Diagram | Mermaid diagram uses process step nodes (Step 2.7, Step 3, Step 4) inside component subgraphs, conflating pipeline stages with structural components — a reader cannot distinguish which boxes are "components" vs "process steps" | -10 pts |
| Architecture:Dependencies | Missing dependency on existing SKILL.md files that are being modified (gen-test-cases/SKILL.md, gen-test-scripts/SKILL.md) and on manifest.yaml format as a coupling point | -10 pts |
| Interface:Signatures | Interface 1 validation rules are prose ("Each ## section must contain ### Verification Dimensions >=3 items") rather than a formal schema — the concrete example compensates but the canonical format rule is still prose | -8 pts |
| Interface:Models | `StrategyFile` model uses `Map<string, VerificationStrategy>` which is a runtime data structure concept, but this is a Skills/LLM layer where the "map" is markdown sections — model conflates runtime with file structure | -8 pts |
| Interface:Implementable | `StrategyMetadata.tokenCount: int` — no specification of which tokenizer computes this count or how it is measured during LLM prompt construction | -10 pts |
| Error:Types | Error codes (STRATEGY_NOT_FOUND, etc.) are well-defined but their representation mechanism remains ambiguous: are these enum values in Go code, LLM prompt directives in SKILL.md prose, or constants in a config file? The Skills-layer context makes this non-obvious | -2 pts |
| Error:Propagation | Propagation strategy is stated globally ("fail-fast on validation errors, best-effort on generation warnings") but not mapped to specific pipeline steps — does Step 2.7 fail-fast on STRATEGY_INVALID? | -2 pts |
| Testing:Per-layer | gen-test-scripts Level-based branching (Interface 3) has no independent test row — it is bundled with gen-test-cases in "Feature E2E Pipeline test" despite being a distinct skill with distinct behavior | -10 pts |
| Testing:Coverage | "Strategy file format edge case coverage >= 80%" has no defined denominator — what are the edge cases? Without enumeration, the 80% target is unmeasurable | -5 pts |
| Testing:Tooling | "Shell harness script" is named but not specified (no file path, no script name, no invocation command); grep/diff/jq are generic Unix commands, not test frameworks. No test runner or harness script location is given | -10 pts |
| Breakdown:Components | The 6 profile strategy files are listed as one monolithic task ("Create verification-strategies.md for each of 6 profiles") rather than individually trackable items | -5 pts |
| Breakdown:Tasks | Three of 9 tasks map to the same file (`plugins/forge/skills/gen-test-cases/SKILL.md`) without specifying which step or section of the SKILL.md corresponds to which task | -7 pts |
| Breakdown:PRD AC | Error messages in design don't exactly match PRD AC: S5 AC says "Strategy/manifest mismatch. Missing in strategy: cli. Extra in strategy: api." but design says "gen-test-cases: strategy/manifest mismatch..."; S3 AC uses Chinese "验证维度 ≥3, 边界场景 ≥2" but design uses English "Verification Dimensions >=3, Boundary Scenarios >=2" | -7 pts |
| Security:Threat Model | LLM prompt injection threat is identified with Medium likelihood but no justification for the likelihood rating, no discussion of attack surface size, and no real-world precedent referenced | -8 pts |
| Security:Mitigations | Strategy file provenance check mitigation ("detected via git blame or file hash comparison") implies runtime code that contradicts the architecture constraint of "Skills layer only — no runtime code changes" | -10 pts |

---

## Attack Points

### Attack 1: Testing Strategy — test tooling lacks specificity and harness definition

**Where**: Testing Strategy table, Tool column: "Shell harness script invoking `forge gen-test-cases` with crafted profile directories", "assertion via exit code (`$?` = 0 for success, 1 for abort) + `grep` on stderr output"
**Why it's weak**: "Shell harness script" is a description of a concept, not a name of a tool. There is no file path, no script name, no invocation command. A developer implementing tests cannot find or run this script because it doesn't exist yet and has no specified location. The assertion tools (exit code, grep, diff, jq, test -d) are generic Unix commands — they work but are not a test framework. There is no test runner, no test file naming convention, and no specification of where test scripts live. For iteration 3, the tooling should name at minimum: (1) the test script file path (e.g., `tests/harness/strategy-validation.sh`), (2) the test runner mechanism (e.g., `make test-design` or `bash tests/harness/run.sh`), and (3) how test results are collected and reported.
**What must improve**: Provide a concrete test script file path and invocation command. Specify where test harness scripts live and how they are executed.

### Attack 2: Security — provenance check mitigation contradicts architecture constraint

**Where**: Security Considerations, Mitigations: "Strategy file provenance check: For community-sourced profiles, gen-test-cases emits a warning when the strategy file is not present in the profile's original source repository (detected via git blame or file hash comparison against known-good versions)."
**Why it's weak**: The architecture section explicitly states "Skills layer only. All changes are within `plugins/forge/skills/` and `forge-cli/pkg/profile/profiles/`. No runtime code changes." But "git blame" and "file hash comparison against known-good versions" are runtime operations requiring Go code to execute git commands, compute file hashes, and maintain a known-good version registry. This mitigation describes runtime behavior that cannot exist in the Skills/LLM layer. A developer implementing this would either (1) skip the mitigation as impossible in the Skills layer, leaving a security gap, or (2) add runtime code, violating the architecture constraint. The mitigation is therefore either unimplementable or architecturally invalid.
**What must improve**: Replace the provenance check with a Skills-layer-appropriate mitigation. Options: (1) add a review checklist item in SKILL.md instructing the LLM to flag suspicious strategy content during parsing, (2) specify that strategy files are trusted by default and document the trust assumption explicitly, (3) add a convention-level mitigation (e.g., profile maintainers must sign off on strategy files via PR review).

### Attack 3: Breakdown-Readiness — PRD AC error messages differ from design error messages

**Where**: PRD Coverage Map, S5 row: "gen-test-cases SKILL.md Step 2.7 key check" with "KEY_MISMATCH error". Error Types table: `KEY_MISMATCH` message is `gen-test-cases: strategy/manifest mismatch. Missing in strategy: {keys}. Extra in strategy: {keys}`.
**Why it's weak**: PRD User Story 5 AC specifies: "中止执行并输出错误：'Strategy/manifest mismatch. Missing in strategy: cli. Extra in strategy: api.'" The design's error message prepends `gen-test-cases:` and uses `{keys}` placeholder. This is a format mismatch — the PRD AC specifies an exact error message string that the test would assert against, and the design adds a context prefix not present in the AC. Similarly, S3 AC specifies error message "Strategy file validation failed for profile: X. Section <capability-key> missing required subsections (验证维度 ≥3, 边界场景 ≥2)" but the design says "strategy file validation failed for profile {name}. Section {key} missing required subsections (Verification Dimensions >=3, Boundary Scenarios >=2)" — the AC uses Chinese field names while the design uses English. A developer implementing against both the PRD AC and the design will face conflicting specifications.
**What must improve**: Either (1) align design error messages exactly with PRD AC error messages, or (2) add a note in the PRD Coverage Map explaining the intentional deviation (e.g., "Error messages follow `docs/conventions/error-handling.md` format with context prefix; PRD AC messages shown without prefix for brevity").

### Attack 4: Testing Strategy — gen-test-scripts has no independent test coverage

**Where**: Testing Strategy table: only "Feature E2E Pipeline test" covers gen-test-cases + gen-test-scripts together. No independent row for gen-test-scripts Level-based branching.
**Why it's weak**: Interface 3 defines a critical behavior: Level-based code generation branching (e2e vs integration templates). This is a distinct skill modification with distinct outputs (different directories, different imports, different assertion styles). The test plan bundles testing this behavior into the same row as gen-test-cases testing, meaning there is no test specifically targeting the Level -> directory + template selection logic in isolation. If gen-test-cases produces correct Level fields but gen-test-scripts ignores them, the combined test might catch it -- or might not if the test only checks gen-test-cases output. The Key Test Scenarios section includes scenario 5 ("Mixed capability profile -> generates e2e + integration code") but this is listed as a scenario, not a structured test with tool and assertion columns.
**What must improve**: Add a separate test row for gen-test-scripts Level-based branching: "gen-test-scripts produces differentiated e2e/integration code in correct directories with distinct imports/assertions" as its own test entry with specific tool and assertion columns.

### Attack 5: Interface & Model — tokenCount computation is unspecified

**Where**: Data Models, StrategyMetadata: `tokenCount: int  // Strategy file token count for monitoring`
**Why it's weak**: The model defines `tokenCount: int` but provides no specification of how this value is computed. Token count depends on the tokenizer used (Claude's tokenizer? tiktoken? character count?). The PRD says "gen-test-cases 完成时输出 token 计数日志（策略文件占比）" but the design does not define: (1) which tokenizer, (2) whether it counts the entire strategy file or just the injected dimensions/boundary sections, (3) what "策略文件占比" means as a ratio (strategy tokens / total prompt tokens?). A developer implementing StrategyMetadata must guess how to compute tokenCount, making the model not directly implementable for this field.
**What must improve**: Specify the tokenization method (e.g., "approximate token count via character count / 4 for English text", or "exact count using tiktoken cl100k_base", or "character count as proxy"). Define whether tokenCount covers the entire strategy file or just the sections used for the active capabilities.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Iter2 Attack 1: Test tooling is SUT, not test tooling | PARTIAL | Replaced "Manual LLM invocation" with "Shell harness script" + specific Unix tools (grep, diff, jq, exit code). Still no test script file path or harness location specified. |
| Iter2 Attack 2: Security mitigations circular, threat model superficial | YES | Threat model now has 2 differentiated threats including LLM prompt injection. Mitigations are differentiated (no shell path, format validation, output review, provenance check). However, provenance check introduces a new issue (contradicts Skills-only architecture). |
| Iter2 Attack 3: Coverage targets are pass/fail gates | YES | Now includes "Level field auto-tagging accuracy >= 95%" and "Strategy file format edge case coverage >= 80%" as percentage-based targets. |
| Iter2 Attack 4: LevelTag no fallback for unknown capabilities | YES | LevelTag model now has explicit fallback: "emit UNKNOWN_CAPABILITY warning, set both level and interface to empty string, and generate a generic test case without Level/Interface fields." Cross-references UNKNOWN_CAPABILITY error. |
| Iter2 Attack 5: GOLDEN_STALE out of scope | YES | GOLDEN_STALE removed from error table. Explicit note added: "Golden file staleness is a test-execution-time concern... out of scope for this design." |

---

## Verdict

- **Score**: 883/1000
- **Target**: 900/1000
- **Gap**: 17 points
- **Breakdown-Readiness**: 181/200 — can proceed to `/breakdown-tasks`
- **Action**: Continue to iteration 4. Priority fixes: (1) Specify test harness script path and invocation, (2) Replace provenance check with Skills-layer-appropriate mitigation, (3) Align error messages with PRD AC or document deviation, (4) Add independent gen-test-scripts test row, (5) Define tokenCount computation method.
