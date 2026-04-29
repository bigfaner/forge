---
date: "2026-04-29"
doc_dir: "docs/features/justfile-e2e-integration/design"
iteration: "3"
target_score: "N/A"
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration 3

**Score: 95/100** (target: N/A)

```
┌─────────────────────────────────────────────────────────────────┐
│                     DESIGN QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Architecture Clarity      │  18      │  20      │ ⚠️          │
│    Layer placement explicit  │   7/7    │          │            │
│    Component diagram present │   6/7    │          │            │
│    Dependencies listed       │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Interface & Model Defs    │  19      │  20      │ ⚠️          │
│    Interface signatures typed│   7/7    │          │            │
│    Models concrete           │   7/7    │          │            │
│    Directly implementable    │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Error Handling            │  14      │  15      │ ⚠️          │
│    Error types defined       │   4/5    │          │            │
│    Propagation strategy clear│   5/5    │          │            │
│    HTTP status codes mapped  │   5/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Testing Strategy          │  14      │  15      │ ⚠️          │
│    Per-layer test plan       │   5/5    │          │            │
│    Coverage target numeric   │   5/5    │          │            │
│    Test tooling named        │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Breakdown-Readiness ★     │  20      │  20      │ ✅          │
│    Components enumerable     │   7/7    │          │            │
│    Tasks derivable           │   7/7    │          │            │
│    PRD AC coverage           │   6/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Security Considerations   │  10      │  10      │ N/A        │
│    Threat model present      │   5/5    │          │            │
│    Mitigations concrete      │   5/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  95      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

★ Breakdown-Readiness 20/20 — can proceed to `/breakdown-tasks`

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Architecture / Component Diagram | "Component Diagram" section is still a file tree with `← EDIT` annotations; execution flow is a separate subsection rather than integrated — the section heading promises a diagram but delivers an inventory | -1 pt |
| Architecture / Dependencies | Version inconsistency: Dependencies section states `just (>= 1.50.0)` but recipe comment states "Requires just >= 1.40.0 (version that introduced [arg] attribute support)" — two different minimum versions in the same document | -1 pt |
| Interface / Directly Implementable | The 1.40.0 vs 1.50.0 contradiction means an implementer cannot determine the actual minimum just version from this document alone | -1 pt |
| Error Handling / Error Types | `npm install` network failure and `npx playwright install chromium` failure during `e2e-setup` still absent from error table — flagged in both iteration 1 and iteration 2, unaddressed for three iterations | -1 pt |
| Testing Strategy / Test Tooling | Document states "Each test below is a `@test` block runnable with `bats t-func.bats`" but the code shown is raw bash, not bats-core `@test` blocks — the tooling claim and the actual test format do not match | -1 pt |

---

## Attack Points

### Attack 1: Architecture — version contradiction between Dependencies section and recipe comment

**Where**: Dependencies section: `just (>= 1.50.0) is already required by init-justfile` vs Interface 2 recipe comment: `# Requires just >= 1.40.0 (version that introduced [arg] attribute support)`

**Why it's weak**: The document states two different minimum just versions for the same feature. The Dependencies section requires 1.50.0; the recipe comment requires 1.40.0. These cannot both be correct. An implementer with just 1.45.0 would read the recipe comment and believe they are compliant, then read the Dependencies section and believe they are not. The fallback note in the recipe ("If your just version predates 1.40.0, replace...") compounds the confusion — it implies 1.40.0 is the real floor, making the 1.50.0 requirement in Dependencies look like a stale copy-paste. The reviser introduced the 1.40.0 citation to address the iter2 attack about unverified syntax, but did not update the Dependencies section to match, creating a new inconsistency.

**What must improve**: Resolve to a single minimum version. If `[arg]` was introduced in 1.40.0, the Dependencies section should state `just >= 1.40.0`. If 1.50.0 is required for other reasons (e.g., a different feature in `init-justfile`), state that reason explicitly and keep both requirements with separate justifications.

---

### Attack 2: Error Handling — `e2e-setup` network failure modes absent from error table for three iterations

**Where**: Error Handling table — five rows cover: `package.json` missing, slug directory not found, `// VERIFY:` markers remain, `--feature` omitted, justfile not found. No row for `npm install` failure or `npx playwright install chromium` failure.

**Why it's weak**: `e2e-setup` calls `npm install --prefix tests/e2e` and `npx --prefix tests/e2e playwright install chromium` — both are network operations that fail in CI environments with restricted egress, behind corporate proxies, or during npm registry outages. The recipe uses `set -euo pipefail`, so any non-zero exit from either command will abort with no diagnostic output beyond the raw npm/playwright error. An agent following the skill that hits a network failure gets no guidance from the error table on what to do. This was flagged in iteration 1, flagged again in iteration 2, and is now in its third iteration unaddressed.

**What must improve**: Add two rows to the error table: one for `npm install` failure (exit code from npm, output: npm error log path, remediation: check network/proxy) and one for `playwright install chromium` failure (exit code, output: playwright error, remediation: `just e2e-setup` with `PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1` or manual install). Alternatively, add a `|| true` with a diagnostic echo if the intent is to treat these as non-fatal.

---

### Attack 3: Testing Strategy — bats-core named as framework but tests are not written in bats format

**Where**: "Test tooling: `bash` + `bats-core` (v1.x)" and "Each test below is a `@test` block runnable with `bats t-func.bats`. If bats-core is unavailable, the inline bash idioms shown are equivalent and directly runnable."

**Why it's weak**: The document claims the tests are `@test` blocks, but the code shown is plain bash scripts with no `@test` declarations, no `setup`/`teardown` functions, and no bats shebang (`#!/usr/bin/env bats`). A developer who copies the test code into `t-func.bats` and runs `bats t-func.bats` will get a parse error because bats requires `@test "name" { ... }` syntax. The "inline bash idioms shown are equivalent" claim is true for the assertion logic, but the framing that these are bats tests is false. The iter2 attack was about pseudo-code `assert` commands — that was fixed. But the fix introduced a new inconsistency: the tests are now runnable as bash but are described as bats.

**What must improve**: Either rewrite the test blocks as proper bats `@test` blocks (add `#!/usr/bin/env bats`, wrap each test in `@test "T-func-1: e2e-setup happy path" { ... }`, use `run just e2e-setup` + `assert_success` + `assert_output --partial "OK: e2e dependencies ready"`), or drop the bats claim and describe the tooling accurately as "bash scripts with inline assertions." The current state misleads an implementer about what they can actually run.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1 (iter2): functional tests use non-existent `assert` bash command | ✅ | Replaced with real bash idioms: `[[ $exit_code -eq 0 ]] \|\| { echo "FAIL T-func-1: ..."; exit 1; }` — tests are now runnable |
| Attack 2 (iter2): component diagram is a file tree with no flow or relationship arrows | ✅ | Execution Flow subsection added with ASCII arrows showing hard gate: `exit 1 ──► HARD GATE: skill marks incomplete` |
| Attack 3 (iter2): `[arg("feature", long)]` syntax unverified, no version citation | ✅ Partially | Version citation added (1.40.0) with fallback syntax for older versions — but introduced a version contradiction with the Dependencies section (1.50.0 vs 1.40.0) |

---

## Verdict

- **Score**: 95/100
- **Target**: N/A
- **Gap**: N/A
- **Breakdown-Readiness**: 20/20 — can proceed to `/breakdown-tasks`
- **Action**: All three iter2 attacks were addressed, with one introducing a new version inconsistency. The remaining 5-point gap is spread across three issues: the 1.40.0/1.50.0 version contradiction (new, introduced by the iter2 revision), the persistent absence of network failure modes from the error table (three iterations unaddressed), and the bats-core/bash format mismatch (new, introduced by the iter2 revision). None of these block implementation, but the version contradiction is the highest-priority fix — it directly affects the minimum just version requirement that will appear in the project's setup documentation.
