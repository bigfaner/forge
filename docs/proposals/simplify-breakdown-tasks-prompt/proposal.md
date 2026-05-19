---
name: simplify-breakdown-tasks-prompt
status: Draft
created: 2026-05-19
---

# Simplify breakdown-tasks Prompt

## Problem

`breakdown-tasks/SKILL.md` is 421 lines (~23KB) with 6 conditional inclusion tags (HAS_UI, NO_UI, UI_ONLY, HAS_PLACEMENT, RULE, HAS_DB), verbose algorithm descriptions, and redundant rule restatements. This causes:

1. **Execution instability** — observed in 3 of 8 recent task-generation runs (PR #117, #119, #121): LLMs applied UI-placement rules to backend-only features, omitted phase gates when PRD contained phase sections, and produced inconsistent scope assignments across runs with identical inputs. Root cause: the LLM evaluates 6 nested conditional tags inline and sometimes misclassifies which branch applies. Concrete example from PR #119: a backend-only feature (no `ui-design.md` artifact) generated tasks with `scope: frontend` and included a "Page Assembly" task — the LLM evaluated the HAS_UI tag as true despite the artifact being absent. Note: full failure output logs are not attached to this proposal; the PR references above contain the relevant task artifacts for verification, but reproducing the exact failure requires re-running the skill against the same input artifacts from those PRs.
2. **Maintenance cost** — modifying one rule requires understanding the entire conditional tag network
3. **Token waste** — 23KB base prompt + artifact files means high per-execution cost
4. **Learning curve** — new contributors cannot reason about the skill without reading the entire 421-line file

Evidence: the skill uses a meta-prompt structure where effective content changes based on which artifacts exist. Six tag pairs gate entire instruction sections, creating 2^3 = 8 possible effective prompt variants that are impractical to test exhaustively given the cost of each full validation run.

## Proposed Solution

**On-demand decomposition**: split SKILL.md into a self-contained skeleton (core rules inline) + conditionally-loaded rule files. The skeleton never references rules that don't apply to the current feature.

**Key principle**: every rule file has a single load condition checked against artifact presence. If the condition is false, the file is never read — zero token cost. The skeleton itself contains only rules that apply to ALL features.

### User-Facing Behavior

From the developer's perspective, nothing changes. Running `/breakdown-tasks` produces the same task files with the same structure, types, dependencies, and scopes. The refactor is purely internal — the prompt that instructs the LLM is reorganized, but the output contract is preserved. The only observable difference is faster execution on simple features (smaller prompt context) and reduced token consumption shown in session logs.

### Load Model

```
                    ┌─────────────────────────────┐
                    │  SKILL.md skeleton (~8KB)    │
                    │  ALWAYS loaded               │
                    │                              │
                    │  · Step 0–7 process flow     │
                    │  · Condition-Rule Matrix     │
                    │  · Element mapping (base)    │
                    │  · Scope algorithm           │
                    │  · Type assignment           │
                    │  · Template selection        │
                    │  · PRD coverage check        │
                    │  · Granularity basics (1-4h) │
                    └──────────┬──────────────────┘
                               │
              ┌────────────────┼────────────────┐
              ▼                ▼                ▼
   ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
   │ phase-       │  │ ui-          │  │ db-schema    │
   │ detection.md │  │ placement.md │  │ .md          │
   │ ~2KB         │  │ ~3KB         │  │ ~1KB         │
   │              │  │              │  │              │
   │ IF: PRD has  │  │ IF: ui/      │  │ IF: design/  │
   │ phase/gate   │  │ ui-design.md │  │ er-diagram.md│
   │ structure    │  │ exists       │  │ exists       │
   └──────────────┘  └──────────────┘  └──────────────┘
              │                │                │
              ▼                ▼                ▼
   phases detected   UI task chains     schema tasks
   inventory.json    integration rules   created
   gates created     placement valid.

   ┌──────────────────┐
   │ existing-code-   │
   │ split.md         │
   │ ~1.5KB           │
   │                  │
   │ IF: tech-design  │
   │ modifies shared  │
   │ existing code    │
   └──────────────────┘
```

### Token Savings by Scenario

| Scenario | Files Loaded | Total Size | vs Current (23KB) |
|----------|-------------|------------|-------------------|
| Greenfield backend, no phases, no DB | Skeleton only | ~8KB | **-65%** |
| Backend + phases + DB | Skeleton + phase + db | ~11KB | **-52%** |
| Backend + existing code modification | Skeleton + existing-code-split | ~9.5KB | **-59%** |
| Full-stack (UI + DB + phases + existing code) | Skeleton + all 4 | ~15.5KB | **-33%** |
| UI-only | Skeleton + ui-placement | ~11KB | **-52%** |

Simpler features get bigger savings. The full-stack worst case still saves 33%.

**Combinatorial note**: the proposed model has 4 independent conditional files, yielding 2^4 = 16 theoretical combinations — larger than the current 2^3 = 8. However, this is not a regression in practice because: (a) each combination is now a **deterministic** function of which artifacts exist (not of which conditional tags the LLM chose to evaluate), (b) the skeleton is the same for all 16 combinations — only the rule overlays change, and each rule file is independent (no interaction effects), (c) the current 8 variants are impossible to enumerate because tag evaluation is embedded in prompt text — the new 16 combinations are at least **explicit** and auditable via the condition-rule matrix. The real testing burden is lower because we validate the skeleton once and each rule file independently, rather than testing all cross-product combinations.

### Structure: Before vs After

```
BEFORE                                    AFTER
breakdown-tasks/                          breakdown-tasks/
├── SKILL.md              421L / 23KB    ├── SKILL.md              ~160L / ~8KB
│   ALL rules inline                      │   ALWAYS-LOADED rules only:
│   ALL conditional tags inline           │   · process flow + Condition Matrix
│   → agent loads 23KB every time        │   · element mapping (non-UI rows)
│                                         │   · scope / type / template selection
│                                         │   · PRD coverage + granularity
├── templates/                            │
│   ├── task.md                           ├── rules/                               ← NEW
│   ├── task-doc.md                       │   ├── phase-detection.md   ~2KB
│   └── manifest-update-tasks.md          │   │   ⚡ IF PRD has phase structure
│                                         │   │
                                          │   ├── ui-placement.md      ~3KB
                                          │   │   ⚡ IF ui/ui-design.md exists
                                          │   │   ⚡ IF prd/prd-ui-functions.md exists
                                          │   │
                                          │   ├── db-schema.md         ~1KB
                                          │   │   ⚡ IF design/er-diagram.md exists
                                          │   │
                                          │   └── existing-code-split.md ~1.5KB
                                          │       ⚡ IF tech-design modifies shared code
                                          │
                                          ├── templates/                          (unchanged)
                                          │   ├── task.md
                                          │   ├── task-doc.md
                                          │   └── manifest-update-tasks.md
```

## Industry Patterns & Prior Art

This approach aligns with established patterns in prompt engineering and LLM orchestration:

1. **Conditional prompt assembly** (LangChain's `PipelinePromptTemplate` — `langchain.prompts.PipelinePromptTemplate`, introduced in langchain-core 0.1.x; DSPy's `dspy.Module` composition with conditional `forward()` dispatch — dspy >= 2.4): these frameworks decompose monolithic prompts into composable modules loaded at runtime based on context. LangChain's pipeline allows composing a `PipelinePromptTemplate` from multiple named sub-prompts, each of which can be conditionally included by wrapping with a `RunnableLambda` that checks input keys. DSPy achieves similar decomposition by having each `Module.forward()` decide which sub-modules to invoke based on input signature. Our condition-rule matrix serves the same role — each rule file is included or excluded based on input artifacts rather than statically composed. The key difference: our checks are file-existence (deterministic boolean) rather than runtime Python conditionals.

2. **RAG-based rule retrieval** (LlamaIndex — `llama-index >= 0.10` with `VectorStoreIndex`; Microsoft Semantic Kernel — `semantic-kernel >= 0.9` with `TextMemoryPlugin`): these tools retrieve relevant instruction fragments from a knowledge store at query time using embedding similarity. Our approach is a simpler variant — instead of vector similarity retrieval, we use deterministic file-existence checks. This trades flexibility for reliability (no retrieval hallucination risk, no embedding quality dependency) while achieving the same goal of loading only relevant rules.

3. **Feature-flag driven configuration** (LaunchDarkly-style toggle routing): the pattern of gating behavior behind feature detection is standard in software engineering. LaunchDarkly's eval-on-read model (`variation()` returns a value per flag state) mirrors our file-existence check returning a load/don't-load decision. The analogy to compiler conditional compilation (`#ifdef`) is apt at the structural level — the skeleton is "common code" and rule files are "conditionally compiled modules" — but no specific compiler technique (dead code elimination, dependency graph analysis, symbol resolution) is adapted here. The analogy serves explanatory purposes only.

We chose deterministic file-existence checks over RAG-based retrieval because: (a) the rule set is small (4 files) and does not justify a retrieval index, (b) file existence is an unambiguous boolean check — no similarity threshold to tune, (c) the forge distribution model already uses file-based artifact detection, making this consistent with existing patterns.

**Acknowledgment**: this proposal does not introduce a novel technique. The core mechanism — conditional module loading based on input artifact detection — is a well-established pattern used by LangChain's `PipelinePromptTemplate` (langchain-core >= 0.1, `langchain.prompts.PipelinePromptTemplate`) and DSPy's module composition (`dspy.Module` with conditional `forward()` dispatch). The contribution here is applying this standard pattern to a concrete pain point in a skill-based CLI tool where it was previously absent, and doing so within the constraint that skills must remain self-contained (no external orchestration layer). If the skill were allowed to call out to a CLI pre-processor, the conditional dispatch could be made programmatic and fully reliable — that is the known better solution this proposal deliberately trades off against the self-containment constraint.

## Alternatives Considered

| Alternative | Pros | Cons | Verdict |
|-------------|------|------|---------|
| **Do nothing** | Zero risk of regression | All 4 pain points persist, growing worse with each addition | Rejected — status quo is unsustainable |
| **Rewrite from scratch** | Clean slate, no legacy baggage | High regression risk, requires re-validating all 8 prompt variants | Rejected — functional regression risk too high |
| **CLI-driven prompt assembly** (CLI pre-checks artifacts, assembles final prompt before LLM call) | Eliminates LLM unreliability in condition evaluation; guarantees correct loading | Requires modifying the CLI layer (currently prompt-agnostic); couples CLI to skill internals; breaks the "skill is self-contained" model | Rejected — violates forge's separation between CLI orchestration and skill instructions |
| **RAG-based rule retrieval** (embed rule files, retrieve by artifact similarity) | Scales to large rule sets; no manual condition matrix needed | Over-engineered for 4 rules; introduces embedding dependency; retrieval may fetch wrong rules | Rejected — complexity unjustified for current scale |
| **Incremental trim with on-demand loading** (this proposal) | Preserves all behavior, reduces each dimension, stays within skill boundary | Rules split across files, slightly less "one-file" readable; LLM must execute condition matrix correctly | **Selected** — best risk/reward ratio, consistent with forge's skill-is-self-contained model |

## Scope

### In Scope

- Restructure `breakdown-tasks/SKILL.md` into skeleton (core rules inline) + conditional rule files
- Merge 6 conditional tags into unified condition-rule matrix with on-demand loading
- Extract 4 conditional rule files: phase-detection, ui-placement, db-schema, existing-code-split
- Validate that generated task output remains identical for representative test cases (most common + most complex scenarios)

### Out of Scope

- `quick-tasks` skill (already lean at ~10KB)
- `forge task` CLI commands (Steps 0, 4b, 5, 6 delegate to CLI — unchanged)
- Task templates (`templates/task.md`, `templates/task-doc.md`)
- Changing any task generation logic, output format, or CLI interface

### Error Handling & Edge Cases

- **Rule file missing**: if a rule file referenced in the condition-rule matrix does not exist (e.g., file deleted accidentally), the LLM should proceed with skeleton-only rules. The skeleton's process steps are complete for all features — rule files are additive, not required. Output will be simpler but structurally valid. Note: the "functionally identical" success criterion only holds when all applicable rule files are present. Degraded output from missing files is an expected fallback, not a regression.
- **Rule file empty**: treated identically to "file missing" — no behavioral change.
- **Artifact exists but is empty or malformed** (e.g., `ui-design.md` has only a heading): the LLM loads the rule file but finds no applicable content. The rule file should include a guard clause: "If the referenced artifact has no parseable content, skip this rule and proceed."
- **LLM reads wrong file**: the condition-rule matrix specifies exact file paths. If the LLM reads an unintended file, the content will not match expected rule patterns, and the downstream validation (`forge task validate-index`) will catch structural deviations.

### Timeline & Resources

- **Duration**: 1-2 days (single developer)
- **Assignee**: project maintainer (familiar with current SKILL.md structure)
- **Breakdown**: extraction + matrix implementation (4-6 hours), validation against baseline (2-3 hours), buffer for iteration (2-4 hours)
- **Delivery**: single atomic commit on a feature branch, merged after validation passes

## Success Criteria

- [ ] SKILL.md reduced from ~23KB to <=8KB skeleton (core rules always needed stay inline)
- [ ] All 6 conditional tags replaced by a condition-rule matrix with on-demand loading
- [ ] No rule file is loaded when its condition is false — zero conditional token cost
- [ ] Generated tasks are functionally identical to current output, validated as follows:
  - **Test case 1 (most common)**: a backend+phases+DB feature (the most common scenario per usage logs) with an existing baseline output
  - **Test case 2 (most failure-prone)**: a full-stack feature loading all 4 rule files (UI + phases + DB + existing code). This is the most complex combination and the most likely to expose LLM non-compliance with the condition-rule matrix. If no full-stack baseline exists, create one with the current SKILL.md before refactoring.
  - **"Functionally identical" means**: same task count (+/-1 tolerance), same dependency graph structure, same type/scope assignments for each task, same PRD coverage (all user stories mapped). Wording may differ but structural elements must match.
  - **Validation protocol**: run `forge task validate-index` on both old and new output; diff the resulting JSON for structural equality; manually review any delta. Both test cases must pass.
- [ ] Each rule file is independently understandable — a contributor can understand a specific rule set without reading the main SKILL.md. Objective validation: given a rule file and the condition-rule matrix, a reviewer correctly identifies (a) the load condition, (b) the rules the file enforces, and (c) the expected output behavior when the file is loaded vs absent. For solo/small teams: perform a self-audit by reading only the rule file (no SKILL.md) and verifying these three points against the extraction plan in this proposal.
- [ ] All file paths use skill-relative references (`rules/X.md`), compatible with forge distribution model
- [ ] Execution stability improved: run the same backend+phases+DB input 3 times; all 3 runs produce structurally equivalent output (same task count, same types). Current baseline: 3/8 runs produce structural deviations.
- [ ] Learning curve reduced: a new contributor can locate and modify a specific rule (e.g., UI placement) within 5 minutes of opening the skill directory, compared to ~15 minutes in the current 421-line file
- [ ] Error handling validated: temporarily rename (remove) one rule file and run the skill — output must still be structurally valid (passes `forge task validate-index`) with skeleton-only rules, producing simpler but coherent tasks

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| LLM loads rule files unconditionally (ignoring conditions) | Medium | High — token savings lost, may exceed current 23KB | Condition-Rule Matrix is the FIRST section in the skeleton, before any step. Each condition is a single file-existence check with explicit IF/ELSE |
| LLM skips a rule file when its condition IS true | Low | High — missing rules produce incomplete tasks | Each step re-checks the condition inline: "IF `ui/ui-design.md` exists, read `rules/ui-placement.md`" — redundant with matrix but acts as safety net |
| Distribution path issues | Low | Medium — skill fails to load | All rule files inside `skills/breakdown-tasks/rules/`, accessed via relative paths — same pattern as `templates/` used by all skills |
| Subtle behavioral regression from wording changes | Medium | Medium — output quality degrades | Run `forge task validate-index` on generated output; compare task structure, dependencies, types, and scopes against baseline |
| Condition-rule matrix too complex for LLM to execute reliably | Medium | High — core mechanism fails | Matrix uses only 4 rows, each with a single file-existence check (no nested conditions). Validate via success criterion "execution stability" (3 consecutive structurally equivalent runs) |
| Rule files drift out of sync with skeleton as SKILL.md evolves | Medium | Medium — stale rules applied or new rules not reflected in files | Convention: any PR that modifies the skeleton's core process steps (Step 0-7 flow, element mapping, scope/type/template) must also check whether rule files need updating. Add a `## Maintenance Note` section at the top of each rule file listing which skeleton sections it depends on, so contributors can trace impact. If a rule file becomes stale, the fail-safe design (missing rules produce simpler output) prevents silent corruption — the output will be visibly incomplete rather than subtly wrong. |

### Rollback Plan

If the refactored skill produces worse task generation:
1. **Immediate rollback**: revert the commit (single commit containing all file changes) — the old SKILL.md is restored as a single file
2. **Validation**: re-run the baseline test case against the reverted version to confirm output quality is restored
3. **Post-mortem**: if rollback is needed, document which specific rule file or matrix instruction caused the regression before attempting a second iteration
4. **Commit strategy**: the entire refactor is a single atomic commit, making `git revert` a complete rollback with no partial state

### Enforcement Mechanism

The condition-rule matrix relies on the LLM executing conditional file reads correctly. To make this reliable:
- The matrix is the **first instruction block** in the skeleton, before any process steps — the LLM encounters it before it begins task generation
- Each condition is a **single file-existence check** (no boolean expressions, no nested logic)
- Each step that needs a rule file **re-prints the load instruction inline** as a safety net (redundant with matrix)
- The skeleton **never contains the rule content inline** — if the LLM skips the conditional read, the rule is absent (fail-safe: the LLM will produce a simpler output rather than a wrong output)

This is intentionally not a programmatic constraint (no CLI pre-processing). The trade-off: we accept a Medium-likelihood risk of LLM non-compliance in exchange for keeping the skill self-contained within the forge distribution model. If this proves unreliable in practice, the next iteration can add CLI-driven assembly as a fallback.

## Rule File Extraction Plan

**Principle**: only extract rules that are conditionally triggered. Always-needed rules stay in skeleton.

| File | Load Condition | Source Content | Size |
|------|---------------|----------------|------|
| `rules/phase-detection.md` | PRD contains phase/gate keywords or explicit section structure | Three-tier phase detection (PRD-defined → heuristic → fallback), `phase-inventory.json` format, phase naming. Depends on: Step 3 (Derive Phases), Step 4b (Phase Gate AC) | ~2KB |
| `rules/ui-placement.md` | `ui/ui-design.md` exists | UI element mapping row, new-page vs existing-page task chains, placement validation against sitemap.json, UI Reference Files requirements, Integration/Page Assembly task chains. Depends on: Step 2 (element mapping rows), Step 4a (task file creation) | ~3KB |
| `rules/db-schema.md` | `design/er-diagram.md` exists | Schema task creation rules, FK/index AC, breaking classification for ALTER vs CREATE. Depends on: Step 2 (element mapping rows), Step 4a (task file creation) | ~1KB |
| `rules/existing-code-split.md` | Tech-design references modifications to existing shared code (interfaces, models, API contracts) | Artifact-update + feature sub-task split, sub-ID convention, when-to-apply threshold (>5 files or cross-layer). Depends on: Step 4a (task file creation), Step 5 (task dependencies) | ~1.5KB |

**Stays in skeleton** (always loaded): element mapping base rows, scope algorithm, type assignment, intent propagation, template selection, PRD coverage check, granularity basics.

Total conditional: ~7.5KB. Skeleton: ~8KB. Worst case (all loaded): ~15.5KB vs current 23KB.

## Condition-Rule Matrix

The skeleton contains this dispatch table. Each row is evaluated independently — a file is loaded if and only if its condition is true.

```
Step 2: Map → Tasks
├─ Read rules/phase-detection.md   IF PRD has phase/gate structure
├─ Add UI mapping rows             IF rules/ui-placement.md loaded
└─ Add DB mapping rows             IF rules/db-schema.md loaded

Step 3: Derive Phases
└─ Apply rules/phase-detection.md  IF loaded (else: artifact-driven decomposition)

Step 4a: Create Task Files
├─ Apply rules/existing-code-split.md  IF loaded (shared code modifications detected)
├─ Apply rules/db-schema.md            IF loaded (DB schema tasks)
├─ Apply rules/ui-placement.md          IF loaded (UI task chains + reference files)
└─ Apply scope/type/template rules      (always — inline in skeleton)
```

No nested conditions, no overlapping tag scopes. Each condition maps to exactly one rule file.
