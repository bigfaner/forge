# Forge Plugin Evaluation

> Date: 2026-04-30 | Evaluator: Claude Opus 4.7 (independent subagent)

## Overview

Forge is a well-designed AI workflow plugin that encodes the complete software engineering lifecycle (brainstorm → PRD → eval → design → breakdown → execute → test) as a repeatable AI workflow.

**Scale:** 18 skills, 12 commands, 4 agents, 18+ templates

---

## Strengths

1. **Exceptional workflow design** — The skill pipeline covers the full engineering lifecycle with clear prerequisite chains and automated remediation messages.

2. **Robust guardrails** — `<HARD-GATE>`, `<HARD-RULE>`, and `<EXTREMELY-IMPORTANT>` tags provide strong behavioral constraints for AI agents. The orchestrator "Iron Laws" in eval skills prevent infinite loops and uncontrolled delegation.

3. **Well-engineered task-cli** — Go CLI with 26 test files, 9 validation passes (circular deps, gate integrity, phase ordering, file existence, etc.), structured error messages with remediation hints.

4. **Traceability throughout** — Manifest system creates bidirectional traceability from PRD acceptance criteria through design → tasks → test cases → execution records.

5. **Adversarial evaluation system** — doc-scorer / doc-reviser separation pattern mirrors real-world review processes with main-session gate control.

6. **Cross-platform support** — Handles Windows (PowerShell) and Unix (bash) installation. Justfile system supports Go, Rust, Node.js, Python, and mixed projects.

7. **Hook-driven automation** — SessionStart, PostToolUse, SessionEnd, Stop hooks provide automatic context injection, index validation, and test execution on completion.

---

## Weaknesses / Risks

### P0 — Critical

1. **No plugin-level tests** — The entire skill system (18 skills, 12 commands, 4 agents) has zero automated tests. Changes can only be verified by manual execution.

### P1 — High

2. **Eval skill duplication** — The four eval skills (eval-proposal, eval-prd, eval-design, eval-ui) share ~80% identical structure. Changes to the orchestrator pattern must be replicated 4 times.

3. **Fragile bash JSON escaping** — Session-start hook constructs JSON via bash string manipulation. Content with unescaped characters could produce invalid JSON and break session initialization silently.

### P2 — Medium

4. **Language inconsistency in task-cli** — `fillRecordTemplate` in `record.go` uses Chinese characters ("无") for empty list placeholders, inconsistent with the English-only codebase.

5. **Single-user assumption** — Task state management (single `process/state.json`) assumes one agent per feature. Concurrent work would cause state corruption.

6. **Hardcoded toolchain assumptions** — Test lifecycle skills assume TypeScript/Playwright/node:test. Projects using Cypress, Selenium, etc. cannot use them without significant adaptation.

7. **No rollback mechanism** — Failed mid-skill execution leaves partial artifacts that may confuse subsequent invocations.

---

## Recommendations

### 1. Create a skill integration test framework (P0)

Design a lightweight harness that exercises the skill pipeline end-to-end:
- Verify prerequisite checks are correct
- Validate templates render with valid placeholders
- Confirm manifest traceability chains are correct
- Ensure task-cli validate catches known error patterns

### 2. Extract shared eval orchestration (P1)

Create `references/shared/eval-orchestration.md` containing the common orchestrator pattern (Iron Laws, Steps 1-5, mermaid diagram, gate logic). Have each eval skill reference this instead of duplicating it. Reduces maintenance from 4 copies to 1.

### 3. Fix bash JSON escaping (P1)

Replace manual string manipulation with `jq` (already used in `validate-index.sh`) or pre-encode the guide.md content at build time.

### 4. Task-cli template English localization (P2)

Change "无" → "None" in `fillRecordTemplate` (`internal/cmd/record.go`). One-line fix.

### 5. Add dry-run mode to breakdown-tasks (P2)

Allow preview of task decomposition without writing files. Reduces rework for large features.

---

## Structural Notes

- Plugin directory: `plugins/forge/`
- Skills: `plugins/forge/skills/` (18 directories)
- Commands: `plugins/forge/commands/` (12 files)
- Agents: `plugins/forge/agents/` (4 files)
- Hooks: `plugins/forge/hooks/`
- Shared references: `plugins/forge/references/shared/`
- Task-cli: `task-cli/` (Go, separate from plugin directory)

## Anti-patterns Found

- `record-task/SKILL.md` instructs agents to use `echo '...' > process/record.json` for JSON writing — fragile with special characters
- `validateFilesExist` computes project root by going up exactly 4 directory levels from index.json — depth-dependent and fragile
- `all-completed.go` silently returns nil on errors (intentional for not-all-done case, but risks masking real errors)
