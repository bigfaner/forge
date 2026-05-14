# Forge Plugin Consistency Rubric

**Total: 1000 points**
**Report template:** `.claude/skills/eval-forge/templates/report.md`

## What This Rubric Measures

Structural consistency of the forge plugin — not individual skill quality, but whether components (skills, commands, agents, templates, hooks, CLI integration) are correctly wired together.

## Scoring Dimensions

| Dimension | Points |
|-----------|--------|
| 1. Directory-Name Alignment | 40 |
| 2. Agent Reference Integrity | 100 |
| 3. Reference Integrity (Templates + Cross-Skill) | 80 |
| 4. Frontmatter Completeness | 110 |
| 5. Eval Template Convention | 100 |
| 6. Orchestrator / Safety Marker Convention | 40 |
| 7. Task CLI Alignment | 240 |
| 8. Hook Wiring Integrity | 70 |
| 9. Guide Coverage | 70 |
| 10. Command Metadata Completeness | 60 |
| 11. Plugin Metadata Consistency | 40 |
| 12. Safety Marker Consistency | 50 |
| **Total** | **1000** |

## Deduction Tiers

| Severity | Penalty | Examples |
|----------|---------|---------|
| Low | -5 | Missing frontmatter field, guide.md extra reference |
| Medium | -15 | Dangling reference, missing template, wrong CLI flag |
| High | -25 | State machine violation, missing hook script, safety marker missing |

## Dimensions

### 1. Directory-Name Alignment (40 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Skill `name` matches directory | 0-25 | Every `SKILL.md` frontmatter `name` equals its parent directory name. Mismatch per skill = -10 (Medium). Missing `name` = auto-fail for this criterion. |
| Command `name` matches filename | 0-15 | Every command file frontmatter `name` equals its filename stem. Missing `name` per file = -10 (Medium). |

### 2. Agent Reference Integrity (100 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Referenced agents exist | 0-70 | Every `forge:<agent-name>` or `subagent_type` reference in skills/commands points to an existing file in `plugins/forge/agents/`. Each dangling = -15 (Medium). |
| No orphan agents | 0-30 | Every agent file is referenced by at least one skill or command. Each orphan = -15 (Medium). |

### 3. Reference Integrity (80 pts)

Merged from Template + Cross-Skill reference checks. All "does the target exist" checks unified.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Template references valid | 0-25 | Every template path referenced in a SKILL.md points to an existing file. Each dangling = -15 (Medium). |
| Cross-skill references valid | 0-30 | Every `invoke /<name>` or `Skill tool` reference points to an existing skill directory or command file. Each dangling = -15 (Medium). |
| No orphan templates | 0-15 | Every file in `skills/*/templates/` is referenced (directly or via rubric→report chain). Each orphan = -5 (Low). |
| No cross-file duplication | 0-10 | No factual information is copy-pasted across 3+ files when a canonical location exists. Each instance = -5 (Low). Exception: autonomous agents that cannot read other files at runtime may duplicate essential facts. |

### 4. Frontmatter Completeness (110 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Skill frontmatter: `name` + `description` | 0-45 | Every `SKILL.md` has both fields. Missing each = -10 (Medium) per file. |
| Command frontmatter: `name` + `description` | 0-35 | Every command file has both fields. Missing `name` = -10, missing `description` = -5 per file. |
| Agent frontmatter: `name` + `description` + `model` | 0-30 | Every agent file has all three. Missing each = -5 (Low) per file. |

### 5. Eval Template Convention (100 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| `eval-*` has rubric.md | 0-30 | Every directory matching `eval-*` contains `templates/rubric.md`. Missing = -15 (Medium) per skill. |
| `eval-*` has report.md | 0-30 | Every `eval-*` contains `templates/report.md`. Missing = -15 (Medium) per skill. |
| Rubric → report chain valid | 0-20 | Every rubric.md references its report.md via `Report template:` line. Missing = -5 (Low) each. |
| Rubric totals correct | 0-20 | Every rubric.md dimension point values sum to the declared total. Wrong total = -10 (Medium). |

### 6. Orchestrator / Safety Marker Convention (40 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| `eval-*` has Iron Laws | 0-25 | Every `eval-*` skill (except `eval-harness`) has `<EXTREMELY-IMPORTANT>` block with "Main session controls the loop". Missing = -15 (Medium) per skill. |
| `eval-*` has Hard Gate | 0-15 | Every `eval-*` skill has `<HARD-GATE>` section. Missing = -5 (Low) per skill. |

### 7. Task CLI Alignment (240 pts)

This dimension requires reading task CLI source code to verify behavioral alignment.

**Source files to read:**
- `task-cli/internal/cmd/claim.go` — task claiming and priority scheduling
- `task-cli/internal/cmd/record.go` — recording, validation, auto-downgrade rules
- `task-cli/internal/cmd/status.go` — state machine transitions and guards
- `task-cli/internal/cmd/all_completed.go` — all-completed hook behavior
- `task-cli/internal/cmd/add.go` — dynamic task addition and ID generation
- `task-cli/pkg/task/types.go` — data model and status/priority enums

| Criterion | Points | What to check |
|-----------|--------|---------------|
| 7a. Command existence | 0-25 | Every `forge <cmd>` referenced in skills/commands/agents exists in `forge -h` output. Each unknown = -15 (Medium). |
| 7b. Flag correctness | 0-25 | Every `forge` flag used in skills matches the CLI's actual flags. Each unknown = -15 (Medium). |
| 7c. Output field parsing | 0-15 | Skills parsing CLI output use correct field names (ACTION, TASK_ID, etc.). Each wrong field = -5 (Low). |
| 7d. Status machine alignment | 0-35 | Read `status.go`. Verify skills respect the state machine: `completed`/`rejected` are terminal, `in_progress → completed` blocked (must use `forge task submit`), `rejected` does NOT satisfy deps. Each violation = -25 (High). |
| 7e. Claim scheduling alignment | 0-35 | Read `claim.go`. Verify skills describe the correct claim priority: (1) all deps met, (2) P0 > P1 > P2, (3) semantic version ordering. Also verify skills handle `ACTION: CONTINUE` correctly. Each violation = -25 (High). |
| 7f. Record validation alignment | 0-35 | Read `record.go`. Verify skills match: auto-downgrade `completed + testsFailed > 0 → blocked` (non-overridable), test evidence required (overridable with `--force`), all AC must be met (overridable with `--force`), quality gate before completed. Each violation = -25 (High). |
| 7g. Dynamic task addition alignment | 0-25 | Read `add.go`. Verify: `--template fix-task` for fix tasks, `--source-task-id` auto-injects dependency, ID format `disc-N`, pre-add pattern (`forge task status blocked` → `forge task add` → `forge task claim`). Each violation = -15 (Medium). |
| 7h. Schema-code alignment | 0-20 | `index.schema.json` fields match Go `Task`/`TaskIndex` struct fields. Verify enum values match. Each mismatch = -5 (Low). |
| 7i. All-completed hook alignment | 0-10 | Read `all_completed.go`. Verify guide.md description matches actual behavior. Each mismatch = -5 (Low). |
| 7j. Template existence | 0-10 | `fix-task` template referenced by skills exists on disk. Missing = -10 (Medium). |

### 8. Hook Wiring Integrity (70 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| `hooks.json` is valid JSON | 0-10 | Parse `plugins/forge/hooks/hooks.json`. Invalid JSON = -10 (Medium). |
| Hook script files exist | 0-25 | Every file path referenced in `hooks.json` (e.g., `${CLAUDE_PLUGIN_ROOT}/hooks/run-hook.cmd`, `scripts/validate-index.sh`) exists on disk. Each missing = -25 (High). |
| Hook CLI commands valid | 0-15 | Every CLI command referenced in `hooks.json` (e.g., `forge cleanup`, `forge quality-gate`) exists in `forge -h`. Each unknown = -15 (Medium). |
| Hook event names valid | 0-20 | Every event name in `hooks.json` (e.g., `SessionStart`, `PostToolUse`, `Stop`, `SessionEnd`, `SubagentStop`) is a Claude Code supported hook event. Each unknown = -15 (Medium). |

### 9. Guide Coverage and Conciseness (70 pts)

Bidirectional check: guide.md references must be valid, workflow-critical skills must be documented, and guide.md must stay concise.

**Workflow-critical skills** are those appearing in the guide.md workflow diagrams (Skill Workflow, Quick Mode) or referenced in the Quality Gate Protocol and Task-CLI sections. Utility/setup commands (init-forge, init-justfile, git-checkout, simplify-skill, extract-design-md, forensic, improve-harness, learn-lesson, record-decision) are NOT required to appear in guide.md.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Guide references are valid | 0-30 | Every `/name` pattern in `plugins/forge/hooks/guide.md` points to an existing skill or command. Each dangling = -15 (Medium). |
| Core workflow skills documented | 0-25 | Every workflow-critical skill (those in Mermaid diagrams, Quality Gate, or Task-CLI sections) has at least a mention in guide.md. Each undocumented workflow-critical skill = -5 (Low). |
| Conciseness / no redundancy | 0-15 | guide.md contains only workflow rules and conventions — no CLI output format tables, no API reference, no information that belongs in `forge -h` or individual SKILL.md files. Each instance of misplaced reference material = -5 (Low). Duplicated information across sections = -5 (Low). |

> Note: guide.md is a workflow guide, not a registry. CLI output field tables belong in individual SKILL.md or `forge -h`. Setup/utility commands need not be mentioned.

### 10. Command Metadata Completeness (60 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| `allowed_tools` declared where needed | 0-35 | Commands that invoke Bash, Edit, or other tools must declare `allowed_tools`. Check each command's content for tool usage and verify `allowed_tools` is declared. Missing = -15 (Medium) per file. |
| `argument-hints` declared where needed | 0-25 | Commands that accept parameters must declare `argument-hints`. Check each command's content for parameter handling. Missing = -5 (Low) per file. |

### 11. Plugin Metadata Consistency (40 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| `keywords` coverage | 0-25 | `plugin.json` keywords cover the core capability areas of all skills/commands. Check for major gaps (e.g., missing "design" when `tech-design` skill exists). Each major gap = -5 (Low). |
| `description` accurate | 0-15 | `plugin.json` description accurately represents the plugin's scope. Grossly inaccurate = -15 (Medium). |

### 12. Safety Marker Consistency (50 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Command/agent markers valid | 0-30 | For each `<EXTREMELY-IMPORTANT>`, `<HARD-RULE>`, `<HARD-GATE>` marker in commands and agents, verify the rule content is actionable and non-contradictory. Non-actionable marker (e.g., "be careful") = -5 (Low). Contradiction between files = -25 (High). |
| Marker coverage for dispatch commands | 0-20 | Commands that dispatch subagents (execute-task, fix-bug, run-tasks, quick) have `<EXTREMELY-IMPORTANT>` blocks with safety constraints. Missing = -15 (Medium) per file. |

## Known Acceptable Discrepancies

These should be noted as INFO, not deducted:

- Schema marks `prd`/`design` as required, but Go allows `Proposal` as alternative (quick mode)
- `sourceTaskID` exists in Go struct but not in JSON schema (auto-managed field)
- `eval-harness` lacks Iron Laws/Hard Gate by design (single-pass, no adversarial loop)
