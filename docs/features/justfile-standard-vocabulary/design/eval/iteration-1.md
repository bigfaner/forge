---
date: "2026-04-30"
doc_dir: "docs/features/justfile-standard-vocabulary/design/"
iteration: "1"
target_score: "90"
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration 1

**Score: 76/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                     DESIGN QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Architecture Clarity      │  18      │  20      │ ✅         │
│    Layer placement explicit  │  7/7     │          │            │
│    Component diagram present │  6/7     │          │            │
│    Dependencies listed       │  5/6     │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 2. Interface & Model Defs    │  17      │  20      │ ✅         │
│    Interface signatures typed│  6/7     │          │            │
│    Models concrete           │  7/7     │          │            │
│    Directly implementable    │  4/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 3. Error Handling            │  12      │  15      │ ✅         │
│    Error types defined       │  4/5     │          │            │
│    Propagation strategy clear│  5/5     │          │            │
│    HTTP status codes mapped  │  3/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 4. Testing Strategy          │  11      │  15      │ ⚠️         │
│    Per-layer test plan       │  4/5     │          │            │
│    Coverage target numeric   │  4/5     │          │            │
│    Test tooling named        │  3/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 5. Breakdown-Readiness ★     │  13      │  20      │ ⚠️         │
│    Components enumerable     │  5/7     │          │            │
│    Tasks derivable           │  3/7     │          │            │
│    PRD AC coverage           │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 6. Security Considerations   │  5       │  10      │ N/A        │
│    Threat model present      │  N/A     │          │            │
│    Mitigations concrete      │  N/A     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ TOTAL                        │  76      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

★ Breakdown-Readiness < 12/20 blocks progression to `/breakdown-tasks` — PASSED (13/20, above gate)

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Component Diagram (line 64-76) | Component diagram is a simplified text sketch, not a true component diagram — it omits the relationship between `breakdown-tasks` and the scope field in index.json, and does not show the 14 skill/agent/command files being migrated | -1 pt (Architecture) |
| Dependencies (line 80-84) | Internal module dependency `task-cli/pkg/task/types.go` is referenced in prose but not listed in the dependency table; the `index.schema.json` file is not listed as a dependency | -1 pt (Architecture) |
| Interface 4 (line 143-156) | Scope Resolution Protocol is written as pseudocode steps, not a typed interface — there is no function signature, no input/output types, no concrete code block a developer could paste into a skill file | -1 pt (Interface Signatures) |
| Interface 2 (line 117-124) | `project-type` recipe contract lacks the recipe body for mixed projects — only shows `@echo "frontend"` but not `@echo "mixed"` or `@echo "backend"` examples with the surrounding recipe syntax | -1 pt (Interface Signatures) |
| Models (line 160-205) | TaskState struct shows `// ... existing fields ...` placeholder instead of listing actual fields — a developer cannot implement this without opening the source file | -2 pt (Directly Implementable) |
| Error Handling (line 229-233) | Error table lacks specificity for `just compile`/`just build` failures — says "enter error-fixer flow" but does not define what exit codes map to compile errors vs. runtime errors vs. missing-tool errors | -1 pt (Error Types) |
| Error Handling (line 229) | No exit code table for `project-type` edge cases beyond "exit != 0" — what about exit code 127 (command not found) vs exit code 1 (recipe error)? The distinction matters for fallback behavior | -2 pt (HTTP status / exit codes) |
| Testing (line 257-261) | Test tooling for skill migration is vague — "node:test" for content checking of skill files does not specify what assertion library or test runner flags are used | -2 pt (Test Tooling) |
| Testing (line 257-261) | No test plan for the scope resolution protocol (Interface 4) — this is the most complex behavioral logic in the design and has zero dedicated test cases | -1 pt (Per-Layer Test Plan) |
| Testing (line 272) | "100% PRD coverage" is stated but there is no coverage matrix mapping test scenarios back to specific PRD acceptance criteria | -1 pt (Coverage Target) |
| Breakdown-Readiness (line 13) | "14 处机械替换 + 1 处 prompt 工程修改" is stated in Overview but never enumerated as a numbered task list — a developer must cross-reference the PRD migration table to derive tasks | -2 pt (Components Enumerable) |
| Breakdown-Readiness | No explicit task-to-interface mapping: Interface 1 (recipe contract) does not map to specific implementation tasks (e.g., "write 15 recipe templates for 3 project types = 45 recipe variants") | -2 pt (Tasks Derivable) |
| Breakdown-Readiness | init-justfile template generation — the design says "3 种项目类型 × 15 个命令模板" but never shows or references the actual template structure. A developer cannot derive "write template for Go backend compile" vs "write template for Node frontend compile" tasks | -2 pt (Tasks Derivable) |
| PRD AC Coverage | Story 4 AC: "agent 连续執行 `just install`、`just compile`、`just test`，全部成功，每步退出码均为 0" — the design does not address sequential command execution or chaining behavior | -1 pt (PRD AC Coverage) |

---

## Attack Points

### Attack 1: Breakdown-Readiness — init-justfile template generation is a black box

**Where**: "3 种项目类型 × 15 个命令模板" (Overview, line 13) and Model 4 (lines 209-222)
**Why it's weak**: The design claims the core deliverable is "3 project types x 15 command templates" but never shows a single concrete template. Model 4 defines project detection signals, but there is no template specification for what a generated justfile actually looks like for each project type. A developer cannot derive implementation tasks like "write the Go backend compile recipe" or "write the Node frontend test recipe" because the design provides zero template examples beyond the generic Interface 1 pattern. The phrase "自适应生成逻辑" is a label, not a specification.
**What must improve**: Add a concrete template for at least one project type (e.g., Go backend) showing all 15 recipes with actual command bodies. Specify whether templates are string literals in the init-justfile command file, separate template files, or generated programmatically. Without this, the single largest implementation effort (45 recipe variants) is entirely unspecified.

### Attack 2: Breakdown-Readiness — no task decomposition for the 14+1 migration

**Where**: "14 处机械替换 + 1 处 prompt 工程修改（breakdown-tasks scope 标注）" (Overview, line 15)
**Why it's weak**: The design identifies 14 mechanical replacements and 1 prompt engineering change as a key deliverable, but never decomposes this into concrete implementation tasks. The PRD Section 5.4 lists specific files and their current/target commands, but the design does not reference this list or translate it into tasks. The "1 处 prompt 工程修改" for `breakdown-tasks` is never elaborated — what is the prompt change? What does the modified prompt text look like? How does it determine file-to-scope mapping? A developer must reverse-engineer the implementation scope from the PRD, which defeats the purpose of a design document.
**What must improve**: (1) Enumerate the 14 files as a numbered migration checklist with before/after command pairs. (2) Specify the prompt engineering change for `breakdown-tasks` — at minimum, show the scope assignment logic or the modified prompt excerpt. (3) Map each migration item to a task description suitable for `/breakdown-tasks` consumption.

### Attack 3: Interface & Model — Scope Resolution Protocol is untyped pseudocode

**Where**: Interface 4: Skill Scope Resolution Protocol (lines 143-156)
**Why it's weak**: This is the most complex behavioral contract in the design — it involves conditional branching across three layers (skill reads index.json, executes `just project-type`, decides whether to pass scope). Yet it is written as numbered prose steps with no typed interface. Compare to Interfaces 1-3 which have code blocks with concrete syntax. A developer implementing a skill file (which is a Markdown prompt, not code) needs to know: what is the exact text block to insert? The design provides an algorithmic description but not the concrete prompt text or template that would appear in a skill file. This is the one interface that is not "directly implementable" from the spec alone.
**What must improve**: Provide the concrete scope resolution block as it would appear in a skill file — either as a prompt template or as a code block with exact variable names. Show at least one worked example (e.g., "task scope=frontend, project-type=mixed -> execute `just build frontend`"). This bridges the gap between "algorithm described" and "code/prompt I can paste."

---

## Previous Issues Check

<!-- First iteration — no previous issues -->

---

## Verdict

- **Score**: 76/100
- **Target**: 90/100
- **Gap**: 14 points
- **Breakdown-Readiness**: 13/20 — can proceed to `/breakdown-tasks` (above 12/20 gate), but tasks will require significant clarification during breakdown
- **Action**: Continue to iteration 2 — primary gaps are in Breakdown-Readiness (task decomposition, template specification) and Testing (scope resolution protocol has no test plan)
