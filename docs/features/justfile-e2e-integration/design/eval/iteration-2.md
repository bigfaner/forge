---
date: "2026-04-29"
doc_dir: "docs/features/justfile-e2e-integration/design"
iteration: "2"
target_score: "N/A"
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration 2

**Score: 92/100** (target: N/A)

```
┌─────────────────────────────────────────────────────────────────┐
│                     DESIGN QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Architecture Clarity      │  17      │  20      │ ⚠️          │
│    Layer placement explicit  │   7/7    │          │            │
│    Component diagram present │   5/7    │          │            │
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
│ 4. Testing Strategy          │  12      │  15      │ ⚠️          │
│    Per-layer test plan       │   4/5    │          │            │
│    Coverage target numeric   │   5/5    │          │            │
│    Test tooling named        │   3/5    │          │            │
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
│ TOTAL                        │  92      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

★ Breakdown-Readiness 20/20 — can proceed to `/breakdown-tasks`

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Architecture / Component Diagram | File tree still shows which files change, not execution flow or component relationships — Attack 4 from iteration 1 not addressed | -2 pts |
| Architecture / Dependencies | `[arg("feature", long)]` syntax requires a specific just version; 1.50.0 stated but not verified against this attribute syntax | -1 pt |
| Interface / Directly Implementable | `[arg("feature", long)]` — the `long` parameter in just's `[arg]` attribute is unverified syntax; if wrong, the recipe silently fails to parse and the verbatim template is incorrect | -1 pt |
| Error Handling / Error Types | `npm install` network failure and `playwright install chromium` failure during `e2e-setup` still absent from error table — flagged in iteration 1, not addressed | -1 pt |
| Testing Strategy / Per-Layer Test Plan | T-func-1 through T-func-6 use `assert exit_code = 0` and `assert output contains "..."` — `assert` is not a POSIX bash command; tests are pseudo-code that cannot run as written | -1 pt |
| Testing Strategy / Test Tooling | Tooling listed as "bash (POSIX shell assertions)" but the assertion syntax used (`assert exit_code = 0`) does not exist in bash; no real test framework or assertion library named | -2 pts |

---

## Attack Points

### Attack 1: Testing Strategy — functional tests are pseudo-code; `assert` is not a bash command

**Where**: `assert exit_code = 0` / `assert output contains "OK: e2e dependencies ready"` (T-func-1 through T-func-6)

**Why it's weak**: Every functional test uses `assert` as if it were a bash builtin. It is not. POSIX bash has no `assert` command. Running these tests verbatim produces `bash: assert: command not found` and exits non-zero on the first assertion — meaning the tests fail for the wrong reason and provide no signal about recipe correctness. The test tooling section names "bash (POSIX shell assertions)" but the assertion syntax shown does not exist in bash. An implementer who tries to run these tests gets immediate failures unrelated to the recipes under test.

**What must improve**: Replace pseudo-assertions with real bash idioms: `[[ $exit_code -eq 0 ]] || { echo "FAIL T-func-1: expected exit 0, got $exit_code"; exit 1; }`. Alternatively, name a real test framework — `bats-core` is the standard choice for bash recipe testing and would make the test file directly runnable.

---

### Attack 2: Architecture — component diagram is still a file tree; Attack 4 from iteration 1 not addressed

**Where**: "### Component Diagram" — the ASCII block is a `plugins/forge/` directory tree with `← EDIT` and `← ADD` annotations

**Why it's weak**: The iteration-1 report explicitly called this out: "The tree shows which files change but not how they relate to each other at runtime." Nothing changed. The diagram still has zero arrows, zero execution flow, and zero indication of how `gen-test-scripts` → `just e2e-verify` (hard gate) → `run-e2e-tests` → `just test-e2e` connects. The core behavioral change this design introduces — a hard gate that blocks skill progression — is invisible in the diagram. A reader unfamiliar with the forge agent architecture learns nothing about component relationships from this tree.

**What must improve**: Add a flow diagram (ASCII arrows are sufficient) showing the agent execution path and where the hard gate fires. The file tree can stay as a change inventory, but it does not substitute for a component relationship diagram.

---

### Attack 3: Interface — `[arg("feature", long)]` syntax is unverified and may be incorrect

**Where**: `[arg("feature", long)]` in the verbatim `e2e-verify` recipe template (Interface 2, line 108)

**Why it's weak**: The `[arg]` attribute in just controls CLI argument parsing. The syntax `[arg("feature", long)]` implies `long` is a valid named parameter to the `[arg]` attribute — but just's attribute syntax uses positional or keyword arguments specific to its own grammar, not clap's. If `long` is not a recognized parameter, just will either reject the justfile at parse time or silently ignore the attribute and fall back to positional argument handling, breaking the `--feature <slug>` interface contract. The dependency section states `just >= 1.50.0` but does not cite which version introduced `[arg]` or verify that `long` is valid syntax. T-func-6 tests the missing-arg case, but if the recipe fails to parse, T-func-6 exits 1 for the wrong reason and masks the bug.

**What must improve**: Cite the just version that introduced `[arg]` support and link to the relevant changelog or docs. Verify that `[arg("feature", long)]` is the correct syntax by running `just --version` against a real justfile containing this attribute. If the syntax is wrong, correct it before the template ships.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: grep-only tests prove text changed, not that recipes work | ✅ Partially | T-func-1 through T-func-6 added; however tests use non-existent `assert` bash command — functional in intent, broken in execution |
| Attack 2: recipe syntax absent from Interface section | ✅ | Verbatim justfile recipe blocks added to both Interface 1 and Interface 2 |
| Attack 3: silent false-pass when slug directory does not exist | ✅ | `if [ ! -d "tests/e2e/{{feature}}" ]` check added to recipe; error table row added |
| Attack 4: component diagram is a file tree, not a diagram | ❌ | Diagram unchanged — still a directory tree with `← EDIT` annotations, no flow or relationship arrows |

---

## Verdict

- **Score**: 92/100
- **Target**: N/A
- **Gap**: N/A
- **Breakdown-Readiness**: 20/20 — can proceed to `/breakdown-tasks`
- **Action**: Three attacks from iteration 1 were addressed cleanly. The remaining gap is concentrated in two areas: the functional tests are pseudo-code that cannot run as written (the single highest-risk issue — it means the test plan looks complete but provides no actual validation), and the component diagram was not updated. The `[arg("feature", long)]` syntax risk is lower severity but should be verified before the template ships.
