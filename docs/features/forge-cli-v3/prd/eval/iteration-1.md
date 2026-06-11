---
date: "2026-05-13"
doc_dir: "docs/features/forge-cli-v3/prd/"
iteration: "1"
target_score: "900"
scoring_mode: "MODE_B"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 1

**Score: 743/1000** (target: 900, mode: MODE_B)

```
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Background & Goals        │  125     │  150     │ ⚠️         │
│    Three elements            │  45/50   │          │            │
│    Goals quantified          │  32/40   │          │            │
│    Logical consistency       │  48/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Flow Diagrams             │  155     │  200     │ ⚠️         │
│    Mermaid diagram exists    │  65/70   │          │            │
│    Main path complete        │  55/70   │          │            │
│    Decision + error branches │  35/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3b. Flow Completeness (B)    │  110     │  200     │ ❌         │
│    Complete business process │  40/70   │          │            │
│    Data flow documented      │  55/70   │          │            │
│    Exception & edge cases    │  15/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. User Stories              │  215     │  300     │ ❌         │
│    Coverage per user type    │  60/70   │          │            │
│    Format correct            │  65/70   │          │            │
│    AC per story (G/W/T)      │  55/60   │          │            │
│    AC verifiability          │  35/100  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Scope Clarity             │  138     │  150     │ ✅         │
│    In-scope concrete         │  48/50   │          │            │
│    Out-of-scope explicit     │  38/40   │          │            │
│    Consistent with specs     │  52/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  743     │  1000    │ ❌         │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

> **Mode B** (prd-ui-functions.md absent): Dimension 3 evaluates Flow Completeness from prd-spec.md Flow Description.
> Sub-criteria: Flow steps describe complete business process /70, Data flow documented /70, Exception handling and edge cases /60.

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| prd-spec.md:Users table:27 | "Hooks/CI" listed as user persona — system actor, not a user | -5 |
| prd-spec.md:Goals:33 | Goal 4 "品牌一致性" has no numeric metric or measurement method | -8 |
| prd-spec.md:Goals:33 | Goal 1 measurement methodology underspecified ("LLM 命令选择测试" — which LLM? how many runs? what prompts?) | -12 |
| prd-spec.md:Flow Description:79-83 | Agent flow narrative only covers task execution lifecycle, omits developer e2e flow and CI flow in text | -30 |
| prd-spec.md:Business Flow Diagram:87-123 | Diagram shows refactoring process, not the feature's runtime business flow for all 3 user types | -15 |
| prd-spec.md:Agent Flow Diagram:129-145 | Primary entry point per Background is `forge prompt get-by-task-id`, but diagram starts with `forge task claim` | -15 |
| prd-spec.md:Agent Flow Diagram:129-145 | No error branch for invalid task ID in prompt/claim, no submit-failure branch, no profile-detection failure | -25 |
| prd-spec.md:Flow Description:69-83 | Zero narrative text describing exception handling, retry logic, or failure states | -45 |
| prd-user-stories.md:Story 1:17 | "自解释的描述" is subjective, not objectively testable | -5 |
| prd-user-stories.md:Story 4:51 | "As a Forge plugin hook 系统" — system actor, inconsistent with Background user types | -5 |
| prd-user-stories.md:Story 5:76 | "自动检测 profile" — how to verify auto-detection? No measurable assertion | -5 |
| prd-user-stories.md:Story 6:91 | "简短描述" — subjective, no measurable threshold | -5 |
| prd-user-stories.md:All stories | Zero error-case or edge-case ACs across all 6 stories (invalid ID, already-completed task, no profile, empty state) | -60 |
| prd-spec.md:In Scope:53 | "更新 23 个 skills" is specific but unnamed; Related Changes table says "23 个 skills" generically | -3 |
| prd-spec.md:Out of Scope:61 | "新增业务功能" is vague — what specific business features were considered? | -2 |
| prd-spec.md:In Scope:54 vs Related Changes:194 | In-scope names specific docs (OVERVIEW.md, WORKFLOW.md, 中文版); Related Changes says just "docs/" — specificity mismatch | -8 |

---

## Attack Points

### Attack 1: Dimension 4 (User Stories) — AC verifiability & boundary coverage catastrophic failure (35/100)

**Where**: All 6 user stories in prd-user-stories.md — zero error-case or edge-case Acceptance Criteria.

**Why it's weak**: Every single AC covers only the happy path. Not one story addresses what happens when things go wrong:
- Story 2: What if the task ID does not exist? What if the task has no prompt template?
- Story 3: What if the task is already completed? What if the submit payload is malformed?
- Story 4: What if cleanup runs but no tasks are in a terminal state? What if quality-gate creates a fix-task but that fix-task also fails?
- Story 5: What if no profile is configured? What if the profile config is invalid?
- Story 6: What if the task type registry is empty?

Additionally, two ACs use subjective language: Story 1's "自解释的描述" (self-explanatory description) and Story 6's "简短描述" (short description) cannot be verified without a subjective judgment call.

**What must improve**: For each story, add at least 2 error/edge-case ACs (Given invalid input / When command called / Then specific error behavior). Replace all subjective AC language with measurable assertions (e.g., "description length <= 80 characters" instead of "short description").

### Attack 2: Dimension 3 (Flow Completeness) — Exception handling and edge cases nearly absent (15/60)

**Where**: prd-spec.md Flow Description section, lines 69-145. Narrative text has zero sentences about failure handling.

**Why it's weak**: The flow description text ("命令迁移流程" and "Agent 使用流程") reads like a sequential checklist of happy-path steps. There is no mention of:
- What happens when a command fails mid-migration (e.g., e2e behavior equivalence check fails)
- What happens when `forge task claim` finds no available tasks (the diagram shows this but no text explains the behavior)
- What happens when profile detection logic encounters an unknown profile type
- What happens when `forge task submit` is called on an already-completed task
- Retry logic for any failure scenario
- Rollback strategy for the 4-phase migration if Phase 3 fails

The diagrams have diamond decision nodes but the accompanying text does not explain the failure semantics.

**What must improve**: Add a dedicated "Error Handling" subsection under Flow Description. Document at minimum: (1) failure behavior for each command, (2) retry/rollback strategy for the migration process, (3) state machine for task status transitions including invalid transitions. Each error path should describe the exit code, error message pattern, and recovery action.

### Attack 3: Dimension 3 (Flow Completeness) — Flow steps describe only 1 of 3 user-type business processes (40/70)

**Where**: prd-spec.md "Agent 使用流程" section (lines 78-83). Only the AI agent task-execution flow has narrative text.

**Why it's weak**: The Background section defines 3 user types with distinct typical scenarios:
- AI agent: claim → execute → submit (covered in narrative)
- Developer: `forge e2e run`, `forge task list-types`, `forge forensic search` (NOT covered in narrative)
- Hooks/CI: `forge cleanup`, `forge quality-gate`, `forge verify-task-done` (NOT covered in narrative)

The developer flow (how does a developer use `forge e2e run` end-to-end?) and the CI/hook flow (when exactly does each hook fire and what is the expected behavior?) have no narrative description. They appear only as named nodes in the Agent Task Execution Flow diagram, without any textual explanation of the business process.

**What must improve**: Add narrative flow descriptions for all 3 user types. For developers: describe the `forge e2e run` lifecycle (profile detection → test selection → execution → result reporting). For Hooks/CI: describe the hook trigger lifecycle (when each hook fires → what command runs → expected state changes → failure behavior). Each flow should cover trigger → processing → end state, matching the detail level of the existing Agent flow.

---

## Previous Issues Check

<!-- Only for iteration > 1 -->

N/A — Iteration 1.

---

## Verdict

- **Score**: 743/1000
- **Target**: 900/1000
- **Gap**: 157 points
- **Action**: Continue to iteration 2 — focus on AC verifiability (+65 pts potential), exception handling in flows (+45 pts potential), and covering all 3 user-type flows in narrative text (+30 pts potential)

SCORE: 743/1000
