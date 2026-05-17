---
created: 2026-05-17
author: "faner"
status: Approved
---

# Proposal: Restructure gen-test-scripts with Per-Type Dispatch

## Problem

`gen-test-cases` was recently refactored into a dispatcher pattern: the main SKILL.md handles generic flow (Steps 0-2.5), then dispatches to `types/{type}.md` instruction files (5 types: UI/TUI/Mobile/API/CLI) for type-specific generation logic. Each type file declares its own conventions in frontmatter.

`gen-test-scripts` lacks this architecture. All type-specific logic — reconnaissance strategies, Fact Table requirements, generation patterns, verification methods — lives in a single monolithic SKILL.md (530 lines). This creates:

1. **Maintenance burden** — Adding or modifying type-specific behavior requires editing the monolith, with risk of cross-type regressions
2. **Inconsistency** — gen-test-cases has per-type conventions (`types/ui.md` → `testing-ui.md`), but gen-test-scripts loads conventions from gen-test-cases type files rather than its own
3. **Cognitive load** — Agents must navigate 530 lines of mixed generic and type-specific instructions, increasing the chance of applying wrong patterns to wrong types

### Evidence

- gen-test-scripts SKILL.md: 530 lines, with Step 1.5 reconnaissance, Step 2-3 sitemap/locators, and Step 4 generation all containing type-specific branches
- gen-test-cases SKILL.md after refactor: ~150 lines (dispatcher) + 5 type files of ~80-120 lines each
- Convention loading in gen-test-scripts (line 50) reads gen-test-cases type files for conventions — indirect coupling

### Urgency

Quantified cost of delay:

- **Growth pressure**: v3.0.0 ships 6 profiles across 5 types today. The next quarter plans 2 additional profiles (see `config.yaml` capability expansions), each adding ~80 lines of type-specific branches to the monolith (extrapolated from the 530-line current state divided across 6 profiles * 5 types). Delaying the restructure means each new profile entrenches inline type branches further.
- **Bug incidence**: The convention loading bug (line 50: `gen-test-scripts` reads conventions from `gen-test-cases/types/{type}.md` instead of its own files) went undetected for 3 sessions because the coupling was buried in a 530-line file. Diagnosing it required tracing two skill directories to find the indirect dependency — approximately 1 hour of debugging that would not occur with self-contained type files.
- **Per-change cost**: Editing the monolith for any single type (e.g., adding a CLI reconnaissance pattern in Step 1.5) requires reading all 5 type branches in that step to avoid cross-type regressions. Measured by recent commits: the last 3 type-specific edits to SKILL.md each touched 2-3 unrelated step sections to maintain consistency.

## Proposed Solution

Restructure `gen-test-scripts` to match `gen-test-cases` dispatcher pattern:

1. **Extract `types/{type}.md`** — 5 type instruction files, each containing:
   - Reconnaissance strategy (what source code to search, type-specific patterns)
   - Fact Table required keys (completeness gate threshold)
   - Type-specific generation steps (UI: sitemap + locators; API: HTTP patterns; CLI: process execution; etc.)
   - Type-specific generation patterns (how to translate test cases → executable scripts)
   - Type-specific antipattern guards (beyond the generic 6)
   - Type-specific conventions in frontmatter

2. **Simplify main SKILL.md** to dispatcher flow:
   - Steps 0-1: Generic (profile resolution, test case reading, auth classification)
   - Step 1.5: Generic Fact Table framework, delegates type-specific reconnaissance to type files
   - Step 3.5: Generic (shared infrastructure — unchanged)
   - Step 4: Dispatches to type files for type-specific generation

3. **Move Step 2-3 to `types/ui.md`** — sitemap resolution and locator mapping are UI-exclusive logic

### User-Facing Behavior

**Before** (monolithic SKILL.md): When an operator runs `forge prompt get-by-task-id T-test-2-cli`, the agent receives the entire 530-line SKILL.md. It must parse type-specific branches inline (`if type == "ui"` logic embedded in Steps 1.5, 2, 3, 4). For an API-only project, the agent still loads UI sitemap instructions, locator patterns, and mobile conventions, consuming context window space with irrelevant content. Error messages for missing conventions are generic ("convention not found") with no type-specific guidance.

**After** (dispatcher + type files): The same invocation produces a ~200-line dispatcher (Steps 0-1, 1.5 generic framework, 3.5, Step 4 dispatch loop) plus one focused type file (~80-120 lines). For `--type cli`, the agent loads only `types/cli.md` — no sitemap/locator content, no UI conventions. Output format, file naming, and test script structure remain identical. The operator sees no behavioral change; the agent receives a smaller, more relevant instruction set. If a type file is missing or `--type` receives an unknown value, the dispatcher emits a specific error naming the expected types and the missing file path.

### Innovation Highlights

While the dispatcher pattern itself is industry-standard, its application here has a distinctive property: **the dispatch boundary is a knowledge boundary, not an execution boundary**. This distinction draws from compiler design, specifically how multi-pass compilers separate frontend passes (language-specific AST construction) from backend passes (target-specific code generation). The intermediate representation (IR) is the narrow contract between them — analogous to how our Fact Table serves as the IR between reconnaissance (type-specific knowledge) and generation (type-specific output).

**Cross-domain inspiration — compiler pass architecture**: In LLVM, each target architecture (x86, ARM, RISC-V) implements a `TargetLowering` class that translates generic IR into machine-specific instructions. Adding a new target does not require changing the IR or other targets. Our type files serve the same role: each translates a generic Fact Table (the IR) into type-specific test scripts (the target output). The key insight borrowed from compilers: if the IR is stable and complete, adding new targets is purely additive. This is why the Fact Table required keys are defined per-type — the "IR schema" varies by target, just as LLVM's `SelectionDAG` varies per target.

**What makes this application distinctive**: In traditional testing frameworks, dispatch boundaries are fixed at framework design time (pytest's conftest hierarchy, Jest's runner interface). In forge, the agent discovers and interprets dispatch rules from free-form markdown at runtime — there is no type system, no plugin API, no compiled registration. The dispatcher must be self-describing enough that a language model follows it correctly without a parser. This constraint — **dispatch correctness without mechanical enforcement** — is the non-obvious design challenge, and it shapes every decision: hard-coded type-to-file mapping (not dynamic), explicit error messages (not silent fallback), and the regression gate (not type-checking).

## Requirements Analysis

### Key Scenarios

- **Single-type generation** (`--type cli`): Dispatcher loads only `types/cli.md`, skips UI sitemap/locator steps entirely
- **Multi-type generation** (no filter): Dispatcher iterates detected types, loading each type's instruction file sequentially
- **New type addition**: Create `types/mobile.md` with mobile-specific reconnaissance and generation patterns — no changes to main SKILL.md
- **Type-specific convention loading**: Each type file declares its conventions in frontmatter, consistent with gen-test-cases pattern
- **Self-owned convention loading**: gen-test-scripts must load per-type conventions from its own `types/{type}.md` frontmatter, not from `gen-test-cases/types/{type}.md`. The current SKILL.md (line 50) reads gen-test-cases type files for conventions — an indirect coupling that caused a 3-session undetected bug. After refactoring, SKILL.md must not reference any file path under `gen-test-cases/types/` for convention resolution.

### Error Scenarios

- **Missing type file** (`types/api.md` does not exist when `--type api` is provided): Dispatcher halts with error message naming the missing file and listing available types. Agent does not fall back to generic instructions.
- **Unknown `--type` value** (`--type graphql`): Dispatcher rejects the value with a message listing the 5 valid types (ui, tui, mobile, api, cli). No partial execution.
- **Malformed frontmatter in type file**: Agent treats frontmatter parsing failure as a skill error — generation does not proceed with default/empty conventions. The error identifies the specific type file and the malformed field.
- **Convention conflict between type file and profile generate.md**: Type file conventions describe *what to test* (strategy); profile generate.md describes *how to write it* (syntax). If both declare the same key (e.g., `test-id-attribute`), profile generate.md wins because it is profile-authoritative. The type file must not duplicate syntax-level conventions.

### Non-Functional Requirements

- **Agent context budget**: The dispatcher + single type file is strictly smaller than the monolith (200 + ~100 lines vs. 530 lines). The primary metric is context window consumption: the agent receives ~40% fewer instruction tokens for single-type tasks. No latency threshold is needed because file size reduction guarantees equal-or-faster loading (the agent reads static local files; no network or parsing overhead). End-to-end agent reasoning time should not increase, but this is difficult to isolate from model variance — the verifiable proxy is the line-count reduction criterion in Success Criteria.
- **Type file loading reliability**: Type files are static markdown in the plugin directory — no network calls, no dynamic resolution. Loading fails only if the file is absent or filesystem is broken. The dispatcher must validate file existence before entering type-specific generation.
- **Backward compatibility**: Existing `--type` filter behavior (profile resolution, test case reading, Fact Table construction) must produce identical outputs. This is verified by running the same test cases before and after refactoring and comparing generated scripts line-for-line.

### Constraints & Dependencies

- Must remain backward-compatible with existing `--type` filter behavior
- Must work with all 6 existing profiles (web-playwright, go-test, maestro, java-junit, rust-test, pytest)
- Type files are loaded by the agent at runtime (read from plugin directory) — no CLI changes needed
- Shared infrastructure (Step 3.5) remains in main SKILL.md — not type-specific

## Existing Per-Type Task Infrastructure

The orchestration layer (task-level splitting) is already fully implemented via the `test-scripts-per-type` proposal. This section documents how per-type tasks flow end-to-end, providing context for the skill-level restructuring proposed here.

### Architecture Overview

```
config.yaml                 testgen.go              index.json           prompt.go              SKILL.md
───────────                 ──────────              ──────────           ──────────              ────────
capabilities: [cli]    →    T-test-2-cli      →    task .md file   →    --type cli       →    filter to CLI only
                            T-test-2-api      →    task .md file   →    --type api       →    filter to API only
                            T-test-3          →    task .md file   →    (no type arg)    →    run all types
```

### Task Generation (`forge-cli/pkg/task/testgen.go`)

`GetBreakdownTestTasks` reads profiles + capabilities from `.forge/config.yaml` and generates per-type tasks:

```go
for i, p := range profiles {
    s := suffixLetter(i, suffix)
    for _, typ := range capabilities {
        tasks = append(tasks, TestTaskDef{
            Key:      "gen-test-scripts-" + p + "-" + typ,
            ID:       "T-test-2" + s + "-" + typ,          // e.g. T-test-2-cli
            Title:    fmt.Sprintf("Generate Test Scripts (%s, %s)", p, typ),
            TestType: typ,                                   // "cli"
        })
    }
}
```

**Profile suffix rules**: single profile → `T-test-2-cli`; multi-profile → `T-test-2a-cli`, `T-test-2b-cli`.

Quick mode (`GetQuickTestTasks`) uses the same pattern but combines gen + run into `T-quick-2-cli`.

### Dependency Chain (`resolveBreakdownDeps`)

```
T-test-1 ──→ T-test-1b ──→ T-test-2-tui ──┐
                         ──→ T-test-2-api ──┤
                         ──→ T-test-2-cli ──┴──→ T-test-3 ──→ T-test-4 ──→ T-test-4.5
```

Per-type gen tasks depend on `T-test-1b` and can execute in parallel. `T-test-3` (run) depends on ALL per-type gen tasks completing.

### Type Inference (`forge-cli/pkg/task/infer.go`)

`InferType` uses three-tier pattern matching on task IDs:

| ID | Match Rule | Inferred Type |
|----|-----------|---------------|
| `T-test-2` | Exact match | TypeTestPipelineGenScripts |
| `T-test-2a` | Profile suffix | TypeTestPipelineGenScripts |
| `T-test-2-cli` | Type suffix | TypeTestPipelineGenScripts |
| `T-test-2a-tui` | Profile + type suffix | TypeTestPipelineGenScripts |

`ExtractTypeSuffix` extracts the capability name: `"T-test-2a-tui"` → `"tui"`.

### Prompt Synthesis (`forge-cli/pkg/prompt/prompt.go`)

When a task is claimed, `forge prompt get-by-task-id` selects the template `test-pipeline-gen-scripts.md` and resolves `{{TEST_TYPE_ARG}}`:

```go
func extractTestTypeArg(id string) string {
    for _, base := range []string{"T-test-2", "T-quick-2"} {
        suffix := task.ExtractTypeSuffix(id, base)
        if suffix != "" {
            return " --type " + suffix
        }
    }
    return ""
}
```

| Task ID | TEST_TYPE_ARG | Final Invocation |
|---------|--------------|-----------------|
| `T-test-2-cli` | `" --type cli"` | `Skill(skill="forge:gen-test-scripts" --type cli)` |
| `T-test-2-api` | `" --type api"` | `Skill(skill="forge:gen-test-scripts" --type api)` |
| `T-test-2` | `""` | `Skill(skill="forge:gen-test-scripts")` |

### Skill Reception (`gen-test-scripts SKILL.md`)

Upon receiving `--type cli`, the skill filters at Step 1:
- Per-type mode: reads only `cli-test-cases.md`
- Fact Table: builds only CLI-related entries
- Steps 2-3 (sitemap/locators): skipped (non-web-ui)
- Step 4: generates only CLI spec files

### Relationship to This Proposal

The orchestration layer above is **fully implemented** and **out of scope** for this proposal. This proposal targets the **skill internal architecture** — restructuring the monolithic SKILL.md into a dispatcher + type files. The two layers are orthogonal:

| Layer | Concern | Proposal |
|-------|---------|----------|
| Orchestration | How many tasks, how they're wired | `test-scripts-per-type` (done) |
| Skill internals | How each task generates scripts | `gen-test-scripts-type-dispatch` (this) |

## Alternatives & Industry Benchmarking

### Industry Patterns for Per-Type Dispatch

| Pattern | Source | How It Works |
|---------|--------|--------------|
| **pytest plugin registration** | pytest (`conftest.py` + `pytest_plugins`) | Each test directory declares fixtures and hooks in its own `conftest.py`. pytest discovers and loads them hierarchically — child directories inherit parent fixtures, sibling directories are isolated. Per-directory scope is the dispatch boundary. |
| **Jest custom runners** | Jest (`runner` config field) | Jest delegates test execution to a configurable runner module. `jest-runner` is the default; alternatives like `jest-runner-eslint` or `jest-runner-prettier` replace it per config. The runner interface is a contract (events: `onTestStart`, `onTestResult`), enabling completely different execution strategies without changing the framework. |
| **Cucumber formatters** | Cucumber (`--format` flag) | Cucumber separates test execution from output formatting. Built-in formatters (progress, summary, HTML) implement a common `Formatter` interface. Custom formatters register via `--format @my-org/cucumber-formatter-slack`. The dispatch boundary is the output concern, not the test concern. |

**Relevance**: All three use the Strategy pattern — a dispatcher that delegates to pluggable, per-variant instruction files. The key insight from pytest: the dispatch boundary should align with the *knowledge boundary* (what fixtures belong to which test domain), not an arbitrary split. For gen-test-scripts, the knowledge boundary is the test type (UI reconnaissance knowledge is unrelated to CLI reconnaissance knowledge).

### Benchmark Mapping to Chosen Approach

The full dispatcher + type files approach maps directly to the pytest model:

- **Our type files are analogous to `conftest.py`** because each one declares its own fixtures (conventions in frontmatter) and hooks (reconnaissance strategy, Fact Table keys, generation patterns) for its domain scope. Just as a UI test directory has its own `conftest.py` with browser fixtures unrelated to an API directory's HTTP fixtures, our `types/ui.md` has DOM-based reconnaissance unrelated to `types/api.md`'s HTTP-route reconnaissance.
- **Our dispatcher SKILL.md is analogous to pytest's test runner** — it handles the generic lifecycle (setup, dispatch, teardown) without knowing domain-specific details. pytest's runner calls `conftest.py` fixtures; our Step 4 calls `types/{type}.md` instructions.
- **The in-place monolith alternative maps to pytest's pre-plugin era** (all fixtures in a single `conftest.py`), which pytest moved away from specifically because cross-domain fixture conflicts caused maintenance issues — the same pattern we see with the convention loading bug at SKILL.md line 50.
- **The config-driven dispatch alternative maps to Jest's runner config** — a declarative mapping file. Jest's approach works for Jest because the runner interface is a stable API contract (typed, versioned). Our agent reads free-form markdown, not a typed interface, so a declarative mapping adds rigidity without the type-safety benefit that justifies it in Jest.

### Alternatives Comparison

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| **Full dispatcher + type files (recommended)** | Architectural consistency with gen-test-cases; each type file is a self-contained unit of knowledge (matches pytest's per-scope `conftest.py` model); adding a new type requires zero changes to the dispatcher | Largest change: main SKILL.md requires full rewrite of Steps 1.5-4; all existing type-specific branches must be extracted and tested | **Selected**: consistency with proven internal pattern outweighs rewrite cost; new type addition is the primary growth vector |
| **Config-driven dispatch** (YAML manifest maps types to instruction files) | Dispatch logic is data, not embedded in markdown; a `dispatch.yaml` can list types and their instruction files without the main SKILL.md containing any hard-coded type names; matches Jest's runner config approach | In the forge agent runtime, the agent has no YAML parser — it reads markdown instruction files as raw text. Adding a YAML manifest means the agent must parse structured data from an unstructured context, which caused 2 misreads in the prior `config.yaml` loading step (the agent treated a commented-out line as active). No internal precedent for this pattern | Strong alternative — rejected because gen-test-scripts runs inside an agent context where markdown-native dispatch (a Step 4 loop listing type files) is simpler than introducing a YAML parser dependency the agent must understand |
| **In-place monolith with per-type sections** (keep SKILL.md as one file, but restructure into clearly delimited `## Type: UI` sections behind a type filter gate) | Zero new files; existing `--type` filter already narrows agent attention; each section is self-contained like Cucumber's feature files co-located in one directory; no file-loading changes needed | Does not solve the cognitive load problem — the agent still receives the full 530+ line file and must skip irrelevant sections; adding a new type still requires editing the monolith; convention loading bug persists because frontmatter is shared across all types | Rejected: addresses maintenance burden partially but fails on cognitive load and convention isolation; the monolith remains a single point of failure for all types |
| **Do nothing** | Zero effort | Growing monolith (530 lines today, ~80 lines per new type added inline); cross-type regression risk confirmed by the convention loading bug (line 50 reads gen-test-cases conventions); no knowledge isolation | Rejected: architectural debt accumulates; the convention inconsistency alone justifies action |

## Scope

### In Scope

- Create `plugins/forge/skills/gen-test-scripts/types/` with 5 instruction files (ui.md, tui.md, mobile.md, api.md, cli.md)
- Refactor main SKILL.md to dispatcher pattern (extract type-specific content)
- Move Step 2 (Sitemap) and Step 3 (Locators) into `types/ui.md`
- Each type file declares its own conventions in frontmatter
- Per-type reconnaissance strategies (what to grep, what patterns to search)
- Per-type Fact Table required keys (completeness gate)
- Per-type verification methods (how to confirm project exposes that interface)
- Per-type generation patterns (how test cases map to executable scripts)

### Out of Scope

- Changes to gen-test-cases (already refactored)
- Changes to profile system or forge CLI
- Changes to other skills (run-e2e-tests, graduate-tests, breakdown-tasks, quick-tasks)
- Task-level splitting (already implemented — see "Existing Per-Type Task Infrastructure" above)
- Adding new test types beyond UI/TUI/Mobile/API/CLI
- Changes to template files (profile-owned, not skill-owned)

## Feasibility Assessment

### Technical Feasibility

The refactoring targets only markdown instruction files (SKILL.md and new type files). No Go code changes, no CLI changes, no template changes. The agent runtime already reads files from the plugin directory. The dispatcher pattern is proven in gen-test-cases (5 type files in production, zero dispatcher bugs reported since refactor).

### Resource & Timeline

- **Effort**: 1-2 days for a single implementer familiar with the gen-test-scripts SKILL.md content
- **Breakdown**: Extract 5 type files (~2 hours each for reconnaissance + Fact Table + generation patterns = 10 hours), refactor main SKILL.md dispatcher (~3 hours), verify across 6 profiles (~2 hours)
- **Skills required**: Understanding of all 5 test type domains (UI/TUI/Mobile/API/CLI) and their reconnaissance patterns; familiarity with gen-test-cases dispatcher structure

### Dependency Readiness

| Dependency | Status | Risk |
|------------|--------|------|
| gen-test-cases dispatcher pattern (reference implementation) | In production, stable | None |
| Profile system (6 profiles) | Fully implemented | None |
| Task-level splitting (`test-scripts-per-type`) | Fully implemented | None |
| Agent file-loading capability | Reads plugin directory files today | None |
| All 5 type files have sufficient domain knowledge | UI/TUI/API/CLI have active profiles; Mobile is template-only | Low — mobile type file will be minimal |

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Agent loads wrong type file or misses dispatch step | M | H | Main SKILL.md Step 4 uses explicit per-type loop with hard-coded type-to-file mapping. **Recovery**: If a generated test script references patterns from the wrong type (e.g., UI sitemap code in a CLI test), the antipattern guards in the type file catch it during Step 4 self-check. |
| Type file content diverges from profile generate.md | L | M | **Prevention**: Add a review checklist item verifying type files contain no `generate.md`-overlapping keys (e.g., import syntax, assertion format). **Recovery**: If divergence is detected post-merge, the `convention conflict` error scenario (profile generate.md wins) handles it at runtime — the agent follows profile-authoritative conventions regardless of type file content. |
| Shared infrastructure logic gets accidentally moved to type files | L | M | Step 3.5 is explicitly labeled "shared — always runs" in the dispatcher. **Recovery**: The regression gate (byte-identical output) catches any missing shared logic — if Step 3.5 logic was moved to a type file, profiles not loading that type file would produce different output. |
| Logic duplication across type files recreates monolith in distributed form | M | M | Type files share no code — each contains unique reconnaissance targets, Fact Table keys, and generation patterns. **Prevention**: During review, diff type files pairwise; shared paragraphs over 5 lines must be extracted back to the dispatcher. **Recovery**: If duplication is found post-merge, extract the common paragraph to dispatcher Step 1.5 or Step 3.5 without touching other type files. |
| Agent fails to follow the dispatcher loop for multi-type generation | M | H | Multi-type dispatch is a sequential loop in Step 4 — the agent must iterate all detected types. **Prevention**: Step 4 explicitly lists the type-to-file mapping as a table the agent reads linearly, not as conditional logic. **Recovery**: The regression gate verifies multi-type output matches pre-refactor output. If a type is skipped, the diff will show missing test scripts. |

## Success Criteria

- [ ] 5 type instruction files exist at `plugins/forge/skills/gen-test-scripts/types/{ui,tui,mobile,api,cli}.md`
- [ ] Each type file declares conventions in frontmatter (verifiable: `grep "^conventions:" types/*.md` returns 5 matches)
- [ ] gen-test-scripts SKILL.md contains zero references to `gen-test-cases/types/` for convention loading (verifiable: `grep -c "gen-test-cases/types" SKILL.md` returns 0). This confirms the convention coupling bug is eliminated, not just masked.
- [ ] Main SKILL.md line count is at most 250 lines, derived from this line budget: Step 0 profile resolution (~25 lines) + Step 1 test case reading & auth classification (~30 lines) + Step 1.5 generic Fact Table framework (~40 lines) + Step 3.5 shared infrastructure (~50 lines) + Step 4 dispatch loop with per-type file mapping (~30 lines) + frontmatter & header (~15 lines) + convention loading (~15 lines) + error handling (~20 lines) = ~225 lines, with 25 lines of margin for formatting. The 250-line ceiling is the sum of this budget; exceeding it indicates type-specific logic was not fully extracted.
- [ ] Step 2 (Sitemap) and Step 3 (Locators) are absent from main SKILL.md and present in `types/ui.md` (verifiable: `grep -c "Sitemap" SKILL.md` returns 0)
- [ ] Each type file contains at least 3 section headings that do not appear in any other type file (reconnaissance targets are type-specific: DOM selectors for UI, process flags for CLI, HTTP routes for API). Measured by diffing section headings across type files — overlap must be below 20% of total headings per file
- [ ] Each type file documents its required Fact Table keys in a dedicated section (verifiable: `grep "Fact Table" types/*.md` returns 5 matches)
- [ ] Each type file contains a verification method section and an antipattern guard section (verifiable: 10 matches total across 5 files)
- [ ] Dispatcher rejects unknown `--type` values with an error listing valid types (test: invoke with `--type graphql`, expect error, no generation)
- [ ] Dispatcher halts on missing type file with error naming the missing file (test: temporarily rename `types/api.md`, invoke `--type api`, expect error)
- [ ] Adding a new type file requires zero edits to main SKILL.md (test: create `types/lambda.md` as a stub, invoke `--type lambda`, observe dispatcher loads it without SKILL.md changes)
- [ ] Generated test scripts are byte-identical before and after refactoring for all 6 profiles when run with the same test case inputs (regression gate: `diff <(old-output) <(new-output)` is empty)
