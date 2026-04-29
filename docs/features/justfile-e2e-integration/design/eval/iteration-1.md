---
date: "2026-04-29"
doc_dir: "docs/features/justfile-e2e-integration/design"
iteration: "1"
target_score: "N/A"
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration 1

**Score: 85/100** (target: N/A)

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
│ 2. Interface & Model Defs    │  15      │  20      │ ⚠️          │
│    Interface signatures typed│   5/7    │          │            │
│    Models concrete           │   6/7    │          │            │
│    Directly implementable    │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Error Handling            │  13      │  15      │ ⚠️          │
│    Error types defined       │   3/5    │          │            │
│    Propagation strategy clear│   5/5    │          │            │
│    HTTP status codes mapped  │   5/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Testing Strategy          │  10      │  15      │ ⚠️          │
│    Per-layer test plan       │   3/5    │          │            │
│    Coverage target numeric   │   5/5    │          │            │
│    Test tooling named        │   2/5    │          │            │
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
│ TOTAL                        │  85      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

★ Breakdown-Readiness 20/20 — can proceed to `/breakdown-tasks`

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Architecture / Component Diagram | File tree shows what changes, not component relationships or execution flow | -2 pts |
| Architecture / Dependencies | `[arg("feature", long)]` syntax implies a specific `just` version for argument parsing; 1.50.0 is stated but not verified against this syntax requirement | -1 pt |
| Interface / Signatures | Interface section describes `just e2e-verify` behavior in pseudo-bash; actual justfile recipe syntax never shown in the Interface section | -2 pts |
| Interface / Directly Implementable | Recipe code deferred to Phase 1 prose; Interface section alone is insufficient to write the recipe | -2 pts |
| Error Handling / Error Types | Missing failure modes: `npm install` network error during `e2e-setup`, `playwright install chromium` failure, `just e2e-verify` with non-existent slug directory (exits 0 silently — false pass) | -2 pts |
| Testing Strategy / Per-Layer Test Plan | All verification is grep-based text matching; no test ever runs `just e2e-setup` or `just e2e-verify` to confirm the recipes function correctly | -2 pts |
| Testing Strategy / Test Tooling | Only `grep` named; no test framework; rubric requires specific test libraries/frameworks | -3 pts |

---

## Attack Points

### Attack 1: Testing Strategy — grep verification proves text changed, not that anything works

**Where**: "Happy path: After all changes, `grep -r 'npx tsx\|cd tests/e2e && npm\|project-test-command' plugins/forge/` returns 0 results"

**Why it's weak**: Every single test scenario in the Testing Strategy section is a grep command checking that old strings are absent or new strings are present. This proves text replacement happened — it does not prove the resulting documentation is correct, that the just recipes execute without error, or that agents following the updated docs will behave as intended. A typo in the recipe template, a wrong flag name, or a broken idempotency check would all pass every test in this plan.

**What must improve**: Add at least one functional verification step: scaffold a minimal test project, run `just e2e-setup` against it, and assert exit 0 + expected output. Run `just e2e-verify --feature nonexistent` and assert exit 1 with the usage message. These are the only tests that actually validate the interface contracts stated in the Interface section.

---

### Attack 2: Interface & Model Definitions — recipe syntax absent from the Interface section

**Where**: Interface 2 shows `just e2e-verify --feature <slug>` with pseudo-bash behavior (`if [ ! -d ... ]`, `Scans tests/e2e/<slug>/**/*.spec.ts`). Phase 1 mentions `[arg("feature", long)]` syntax in passing.

**Why it's weak**: The Interface section is the contract an implementer reads to write the recipe. It describes what the command does in pseudo-bash but never shows the actual justfile recipe syntax. `[arg("feature", long)]` is a just-specific annotation that controls how arguments are parsed — its placement, the full recipe structure, and whether it requires `set positional-arguments` or other justfile directives are all left to the implementer to figure out. The interface is incomplete without the recipe skeleton.

**What must improve**: Add the actual justfile recipe template to each Interface definition — the same template that will appear in `init-justfile.md`. The Interface section should be self-contained: read it, write the recipe, done.

---

### Attack 3: Error Handling — silent false-pass when slug directory does not exist

**Where**: Error table row: "`// VERIFY:` markers remain | `just e2e-verify` | 1 | Count + file:line list"

**Why it's weak**: The behavior description says `e2e-verify` "Scans `tests/e2e/<slug>/**/*.spec.ts` for lines matching `// VERIFY:`" and "if count > 0 → exit 1." The inverse — if count = 0 — exits 0 with "OK: no unresolved markers." But if `<slug>` is misspelled or the directory does not exist, the glob matches nothing, count = 0, and the command exits 0 with a success message. This is a silent false-pass: an agent with a wrong slug gets a green light and proceeds to `run-e2e-tests` with unresolved markers still in the actual spec files. The error table has no row for this case.

**What must improve**: Add an explicit check: if `tests/e2e/<slug>/` does not exist, exit 1 with `Error: tests/e2e/<slug>/ not found`. Add this as a row in the error table and in the Interface 2 exit code list.

---

### Attack 4: Architecture — "Component Diagram" is a file tree, not a diagram

**Where**: "### Component Diagram" — the ASCII block is a `plugins/forge/` directory tree with `← EDIT` and `← ADD` annotations.

**Why it's weak**: The rubric criterion is "components and relationships." The tree shows which files change but not how they relate to each other at runtime. There is no diagram showing: which agents consume which skill files, what triggers `e2e-setup` vs `e2e-verify`, how the hard gate in `gen-test-scripts` connects to `run-e2e-tests`, or where the justfile sits relative to the skill invocation chain. For a reader unfamiliar with the forge agent architecture, the diagram provides no structural understanding.

**What must improve**: Add a flow diagram (even ASCII arrows) showing the agent execution path: `gen-test-scripts` → `just e2e-verify` (hard gate) → `run-e2e-tests` → `just test-e2e`. This is the core behavioral change the design introduces and it is invisible in the current diagram.

---

## Previous Issues Check

N/A — iteration 1.

---

## Verdict

- **Score**: 85/100
- **Target**: N/A
- **Gap**: N/A
- **Breakdown-Readiness**: 20/20 — can proceed to `/breakdown-tasks`
- **Action**: Document is implementable as-is. Three issues should be addressed before implementation to avoid rework: the silent false-pass in `e2e-verify` (Attack 3) is a correctness bug in the interface spec; the missing recipe syntax (Attack 2) will force the implementer to reverse-engineer the just argument syntax; and the grep-only test plan (Attack 1) will not catch a broken recipe.
