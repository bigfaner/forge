---
date: "2026-04-30"
doc_dir: "docs/features/justfile-standard-vocabulary/design/"
iteration: "2"
target_score: "90"
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration 2

**Score: 88/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                     DESIGN QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 1. Architecture Clarity      │  20      │  20      │ ✅         │
│    Layer placement explicit  │  7/7     │          │            │
│    Component diagram present │  7/7     │          │            │
│    Dependencies listed       │  6/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 2. Interface & Model Defs    │  18      │  20      │ ✅         │
│    Interface signatures typed│  7/7     │          │            │
│    Models concrete           │  7/7     │          │            │
│    Directly implementable    │  4/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 3. Error Handling            │  14      │  15      │ ✅         │
│    Error types defined       │  5/5     │          │            │
│    Propagation strategy clear│  5/5     │          │            │
│    HTTP status codes mapped  │  4/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 4. Testing Strategy          │  14      │  15      │ ✅         │
│    Per-layer test plan       │  5/5     │          │            │
│    Coverage target numeric   │  4/5     │          │            │
│    Test tooling named        │  5/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 5. Breakdown-Readiness ★     │  16      │  20      │ ✅         │
│    Components enumerable     │  7/7     │          │            │
│    Tasks derivable           │  4/7     │          │            │
│    PRD AC coverage           │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 6. Security Considerations   │  6       │  10      │ ✅         │
│    Threat model present      │  3/5     │          │            │
│    Mitigations concrete      │  3/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ TOTAL                        │  88      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

★ Breakdown-Readiness 16/20 — PASSED (above 12/20 gate)

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Interface 4 (lines 204-222) | Scope Resolution Protocol is embedded as a prompt text block in Markdown — this is appropriate for skill files, but the "worked examples" table does not cover the full 6-row scenario from the PRD flowchart (specifically missing: scope=backend, project-type returns unexpected string like "unknown" with exit 0) | -1 pt (Interface Signatures) |
| Model 5 / Recipe Template table (lines 354-374) | The recipe table uses shorthand like "bash case: `npm run build` / `go build ./...`" instead of showing the full recipe body for every cell — a developer must mentally reconstruct the full bash case block from the pattern, which is a direct implementability gap | -2 pt (Directly Implementable) |
| Error Handling (lines 386-387) | `just project-type` exit code table mentions "exit != 0" generically but does not distinguish between exit code 1 (recipe failed), exit code 127 (command not found), and exit code 2 (just argument error) — these are different failure modes that a skill implementing fallback logic should handle distinctly | -1 pt (Exit Codes) |
| Testing (line 411) | "100% PRD coverage" is claimed but no test-to-AC traceability matrix exists — the "Key Test Scenarios" list 5 scenarios but the PRD has 12+ acceptance criteria across 5 stories; no mapping shows which test covers which AC | -1 pt (Coverage Target) |
| Breakdown-Readiness (lines 354-374) | While the recipe table now shows all 15 commands across 3 project types, no task decomposition maps these into implementation units — e.g., "Task A: write backend-only justfile template (15 recipes)", "Task B: write frontend-only template (15 recipes)", "Task C: write mixed template (15 recipes with scope)" — a developer must infer the task structure | -3 pt (Tasks Derivable) |
| PRD AC Coverage | Story 4 AC: "agent 连续执行 `just install`、`just compile`、`just test`，全部成功，每步退出码均为 0，agent 无需人工介入" — sequential command chaining is not addressed in the design. No recipe or protocol specifies how sequential execution works or how partial failure (e.g., install succeeds, compile fails) is reported back to the agent | -1 pt (PRD AC Coverage) |
| Security (lines 430-440) | Threat model lists only 2 threats, one of which is low-risk. The "中" risk (init-justfile overwrite) mitigation is "提示确认" but the design does not specify the confirmation mechanism — is it a CLI prompt? A skill prompt? Since agent files cannot use interactive prompts (PRD requirement), this is an unresolved contradiction | -2 pt (Threat Model) |
| Security (lines 438-439) | "init-justfile 检测已有 justfile 时提示确认" — no specification of what happens if the user does not confirm (abort? overwrite anyway? merge?) and no specification of how the `# --- forge standard recipes ---` boundary marker is used during overwrite (preserve user recipes? overwrite everything?) | -2 pt (Mitigations) |

---

## Attack Points

### Attack 1: Breakdown-Readiness — recipe template implementation tasks are still implicit

**Where**: Model 5 recipe table (lines 354-374) and "各项目类型生成的 recipe 清单"
**Why it's weak**: The design provides a comprehensive recipe table showing what each of the 15 commands looks like across 3 project types, but never decomposes this into implementable task units. The table shows 15 rows x 3 columns = 45 recipe variants, but the design does not specify: (a) are these templates stored as string literals in `init-justfile.md`, or as separate template files, or generated programmatically? (b) what is the task breakdown — one task per project type (3 tasks), or one task per recipe (15 tasks), or one monolithic task? The line "模板以字符串字面量形式内嵌于 init-justfile 命令文件中" gives storage format but not task decomposition. A developer using `/breakdown-tasks` must guess how to chunk this work.
**What must improve**: Add an explicit implementation task list like: "Task 1: Implement pure-backend recipe templates (15 recipes as string literals in init-justfile)", "Task 2: Implement pure-frontend recipe templates (15 recipes)", "Task 3: Implement mixed recipe templates (15 recipes with scope dispatch)", "Task 4: Implement project-type detection logic in init-justfile". Each task should reference the specific model/interface it implements.

### Attack 2: Breakdown-Readiness — no task mapping from migration checklist to implementation

**Where**: Migration Checklist (lines 23-39)
**Why it's weak**: The migration checklist now enumerates all 15 items with file names, current/target commands, and types (mechanical replacement / verification / prompt engineering). This is a significant improvement over iteration 1. However, the "verification" items (rows 5-14) are 10 items that require no code changes — they are validation tasks, not implementation tasks. The design does not specify how validation is performed (manual grep? automated test? CI check?). The 3 mechanical replacements (rows 1-4) and 1 prompt engineering change (row 15) are the actual work items, but these are mixed into the same table without clear separation. A developer cannot derive "do these 4 tasks" from the current presentation.
**What must improve**: Split the migration table into two sections: (1) "Implementation Tasks" with the 4 items that require actual changes (rows 1, 2, 3, 4, 15), each mapped to a concrete task description; (2) "Validation Checklist" with the 10 verification items, each with a specified validation method (e.g., "grep -r 'just build && just test' plugins/ should return 0 results").

### Attack 3: Security — init-justfile overwrite threat has unresolved contradiction with agent requirements

**Where**: Security Considerations (lines 430-440) and PRD Agent Requirements (PRD line 253)
**Why it's weak**: The security section identifies "init-justfile 覆盖" as a medium-risk threat and proposes "提示确认" as mitigation. But the PRD explicitly states "所有 recipe 不使用交互式输入（如 read、select），agent 无法处理交互". If `init-justfile` is invoked by an agent (which is a stated persona), the confirmation prompt creates a contradiction — the agent cannot respond to an interactive prompt. The design does not resolve this: does init-justfile skip confirmation when invoked by an agent? Is there a `--force` flag? Does the boundary marker (`# --- forge standard recipes ---`) enable merge-mode instead of overwrite? The mitigation "提示确认" is stated but its mechanism is unspecified, leaving the medium-risk threat without a concrete, implementable countermeasure.
**What must improve**: Specify the confirmation mechanism: (a) when invoked interactively by a human, prompt for confirmation; (b) when invoked by an agent, either use `--force` flag or implement merge-mode (preserve sections outside the `# --- forge standard recipes ---` boundary). The boundary marker is already defined but never tied to the overwrite behavior. Add this to Interface 1 or as a new subsection.

---

## Previous Issues Check

| Previous Attack (Iteration 1) | Addressed? | Evidence |
|-------------------------------|------------|----------|
| Attack 1: init-justfile template generation is a black box | ✅ Partially | Model 5 now shows concrete recipe templates for mixed and backend projects (lines 316-351), plus a full 15-command x 3-type table (lines 354-374). Still uses shorthand "bash case:" notation rather than full recipe bodies for every cell. |
| Attack 2: no task decomposition for 14+1 migration | ✅ Yes | Migration Checklist (lines 23-39) now enumerates all 15 items with file paths, current/target commands, and change types. The breakdown-tasks scope annotation (line 41-77) now shows the complete prompt text with algorithm, examples, and edge cases. |
| Attack 3: Scope Resolution Protocol is untyped pseudocode | ✅ Yes | Interface 4 (lines 204-222) now contains the full scope resolution text block as it would appear in a skill file, with a 6-row worked examples table covering all major scenarios including fallback cases. |

**Other iteration-1 deductions addressed:**

| Iteration 1 Deduction | Addressed? | Evidence |
|----------------------|------------|----------|
| Component diagram omits breakdown-tasks/scope relationship | ✅ Yes | Component diagram (lines 124-136) now shows `breakdown-tasks ──generates──▶ index.json` and `tasks[].scope ──reads──▶ skill execution` |
| Dependencies table missing `task-cli/pkg/task/types.go` and `index.schema.json` | ✅ Yes | Dependencies table (lines 141-146) now lists both `task-cli/pkg/task/types.go` and `index.schema.json` |
| Interface 2 lacks mixed/backend project-type examples | ✅ Yes | Model 5 (lines 341-351) now shows all three project-type recipe variants |
| TaskState struct had `// ... existing fields ...` placeholder | ✅ Yes | Model 2 (lines 258-272) now lists all fields explicitly |
| No test plan for scope resolution protocol | ✅ Yes | Key Test Scenarios (lines 417-422) now includes scenario 4: "scope 与 project-type 不匹配：skill 层警告 + 回退" |
| "100% PRD coverage" with no traceability matrix | ❌ No | Still states "100% PRD 覆盖" (line 425) without a test-to-AC mapping table |
| Sequential command execution (Story 4 AC) not addressed | ❌ No | Still not addressed — no design for how agent chains `just install && just compile && just test` with partial failure handling |

---

## Verdict

- **Score**: 88/100
- **Target**: 90/100
- **Gap**: 2 points
- **Breakdown-Readiness**: 16/20 — can proceed to `/breakdown-tasks` (well above 12/20 gate)
- **Action**: Continue to iteration 3 — remaining gap is 2 points. Primary issues: (1) recipe template task decomposition is still implicit, not explicit; (2) sequential command execution AC remains unaddressed; (3) init-justfile overwrite mitigation has an agent-contradiction. All are fixable in one pass.
