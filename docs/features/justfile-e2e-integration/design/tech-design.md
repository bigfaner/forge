---
created: 2026-04-29
prd: prd/prd-spec.md
status: Draft
---

# Technical Design: Justfile E2E Integration

## Overview

This feature is a pure documentation change: edit 13 markdown files in `plugins/forge/` to replace raw shell commands with just targets, and extend the `init-justfile` command template with two new recipes (`e2e-setup`, `e2e-verify`). No application code is written. No new dependencies are introduced.

The implementation is a sequence of targeted text replacements across skill, agent, command, and task-template files.

## Architecture

### Layer Placement

Single-layer feature. All changes are confined to `plugins/forge/` documentation files (SKILL.md, command .md, agent .md, task template .md). No application code, no APIs, no data models.

### Component Diagram

```
plugins/forge/
├── commands/
│   ├── init-justfile.md        ← ADD: e2e-setup + e2e-verify recipe templates
│   ├── fix-bug.md              ← EDIT: replace raw test commands
│   ├── run-tasks.md            ← EDIT: replace raw test commands
│   └── execute-task.md         ← EDIT: replace raw test commands
├── skills/
│   ├── gen-test-scripts/SKILL.md  ← EDIT: Step 4 + Step 5
│   ├── run-e2e-tests/SKILL.md     ← EDIT: Step 1 + Step 2 + Error table
│   ├── record-task/SKILL.md       ← EDIT: Metrics Collection
│   └── improve-harness/SKILL.md   ← EDIT: Step 4.3
├── agents/
│   ├── task-executor.md        ← EDIT: Step 3
│   └── error-fixer.md          ← EDIT: Step 4
└── skills/breakdown-tasks/templates/
    ├── run-e2e-tests.md         ← EDIT: Implementation Notes
    ├── gen-test-scripts.md      ← EDIT: Implementation Notes
    └── fix-e2e.md               ← EDIT: add post-fix verification step
```

### Dependencies

No new dependencies. `just` (>= 1.50.0) is already required by `init-justfile`.

## Interfaces

### Interface 1: `just e2e-setup`

```bash
# Signature
just e2e-setup

# Behavior
if [ ! -d tests/e2e/node_modules ]; then
    npm install --prefix tests/e2e
fi
npx --prefix tests/e2e playwright install chromium

# Exit codes
# 0 — success, outputs: "OK: e2e dependencies ready"
# 1 — tests/e2e/package.json not found, outputs: "Error: tests/e2e/package.json not found"
```

### Interface 2: `just e2e-verify --feature <slug>`

```bash
# Signature
just e2e-verify --feature <slug>

# Parameters
# --feature <slug>  Required. Subdirectory name under tests/e2e/ (lowercase, hyphens).
#                   Agent obtains via: task feature

# Behavior
# Scans tests/e2e/<slug>/**/*.spec.ts for lines matching "// VERIFY:"
# Counts matches; if count > 0 → exit 1 with file:line list

# Exit codes
# 0 — no unresolved markers, outputs: "OK: no unresolved // VERIFY: markers in tests/e2e/<slug>/"
# 1 — markers found, outputs: count + file:line list for each marker
# 1 — --feature omitted, outputs: "Usage: just e2e-verify --feature <slug>"
```

## Data Models

Single-layer documentation feature. No data models.

## Error Handling

### Error Types

| Situation | Command | Exit Code | Output |
|-----------|---------|-----------|--------|
| `package.json` missing | `just e2e-setup` | 1 | `Error: tests/e2e/package.json not found` |
| `// VERIFY:` markers remain | `just e2e-verify` | 1 | Count + file:line list |
| `--feature` arg omitted | `just e2e-verify` | 1 | `Usage: just e2e-verify --feature <slug>` |
| `justfile` not found in project | any `just` call | — | Skill checks `ls justfile` first; prompts user to run `/init-justfile` |

### Propagation Strategy

Skills check `ls justfile` before invoking any just command. If missing, skill aborts with a user-facing message. No silent failures.

`just e2e-verify` exit 1 is a hard gate in `gen-test-scripts`: skill marks itself incomplete and does not proceed to `run-e2e-tests`.

## Cross-Layer Data Map

Single-layer feature — not applicable.

## Implementation Plan

Changes are grouped into 4 phases. Each phase is independently verifiable.

### Phase 1: New Targets (init-justfile)

**File**: `plugins/forge/commands/init-justfile.md`

Changes:
1. Add `e2e-setup` and `e2e-verify` rows to the Standard Target Contract table
2. Add `e2e-setup` recipe template (all languages share the same recipe)
3. Add `e2e-verify` recipe template with `[arg("feature", long)]` syntax
4. Update Step 4 Output Confirmation to include the two new targets

**Verification**: `grep -c 'e2e-setup\|e2e-verify' plugins/forge/commands/init-justfile.md` >= 4

---

### Phase 2: E2E Skill Files

**File**: `plugins/forge/skills/gen-test-scripts/SKILL.md`

| Location | Before | After |
|----------|--------|-------|
| Step 4 post-generation check (code block) | `grep -r '// VERIFY:' tests/e2e/<feature>/` | `just e2e-verify --feature <slug>` |
| Step 4 post-generation check (prose) | "run `grep -r '// VERIFY:'`..." | "run `just e2e-verify --feature <slug>`; exit 1 = skill incomplete" |
| Step 5 deps install (code block) | `cd tests/e2e && npm install` | `just e2e-setup` |

**File**: `plugins/forge/skills/run-e2e-tests/SKILL.md`

| Location | Before | After |
|----------|--------|-------|
| Step 1 deps install (code block) | `cd tests/e2e && npm install` | `just e2e-setup` |
| Step 1 playwright install (code block) | `npx playwright install chromium` | (removed — covered by `just e2e-setup`) |
| Step 1 verify install (code block) | `npx playwright --version` | (removed — covered by `just e2e-setup`) |
| Step 1 prose | "Install dependencies (if `node_modules` doesn't exist):" | "Run `just e2e-setup` (idempotent — installs deps and Playwright browser):" |
| Step 2 CLI spec (code block) | `npx tsx <slug>/cli.spec.ts 2>&1 \| tee ...` | `just test-e2e --feature <slug> 2>&1 \| tee results/<slug>-output.txt` |
| Step 2 API spec (code block) | `npx tsx <slug>/api.spec.ts 2>&1 \| tee ...` | (removed — test-e2e runs all specs) |
| Step 2 UI spec (code block) | `npx tsx <slug>/ui.spec.ts 2>&1 \| tee ...` | (removed — test-e2e runs all specs) |
| Error table: playwright not installed | "Run `npx playwright install chromium`, retry" | "Run `just e2e-setup`, retry" |
| Error table: node_modules missing | "Run `npm install`, retry" | "Run `just e2e-setup`, retry" |

**Verification**:
- `grep -c 'npx tsx\|npx playwright\|npm install' plugins/forge/skills/run-e2e-tests/SKILL.md` = 0
- `grep -c 'just e2e-setup\|just test-e2e' plugins/forge/skills/run-e2e-tests/SKILL.md` >= 2

---

### Phase 3: Build/Test Files (7 files)

**`plugins/forge/commands/fix-bug.md`**

| Location | Before | After |
|----------|--------|-------|
| Step 2 code block | `<project-test-command>   # e.g. npm test, go test ./..., pytest` | `just test` |
| Step 2 prose | "Run existing tests to establish baseline" | "Run `just test` to establish baseline" |
| Step 3a code block | `<project-test-command> --testNamePattern "bug:"` | `just test` |
| Step 3a prose | "Run the new test — it must fail before the fix:" | "Run `just test` — it must fail before the fix:" |
| Step 3b code block | `npx tsx <spec-file> 2>&1` | `just test-e2e --feature <slug>` |
| Step 3b prose | "Run the e2e test — it must fail before the fix:" | "Run `just test-e2e --feature <slug>` — it must fail before the fix:" |
| Step 5 code block | `<project-test-command>` | `just test` |
| Step 5 code block | `npx tsx <spec-file> 2>&1` | `just test-e2e --feature <slug>` |
| Step 5 prose | "Run the full test suite. All tests must pass." | "Run `just build && just test`. All must pass." |

**`plugins/forge/commands/run-tasks.md`**

| Location | Before | After |
|----------|--------|-------|
| Step 5 code block comment | `# go test ./... OR npm test OR the testCommand from index.json` | `just test` |
| Step 5 prose | "Run project-level full test suite" | "Run `just test`" |

**`plugins/forge/agents/task-executor.md`**

| Location | Before | After |
|----------|--------|-------|
| Step 3 code block | `# Go: go build ./... && go vet ./... && go test -race -cover ./...` | `just build && just test` |
| Step 3 code block | `# Node: npm run build && npm test` | (removed) |
| Step 3 code block | `# Python: pytest --cov` | (removed) |
| Step 3 prose | "Run complete verification suite for your project:" | "Run `just build && just test`:" |

**`plugins/forge/agents/error-fixer.md`** — same changes as task-executor Step 3/Step 4

**`plugins/forge/commands/execute-task.md`**

| Location | Before | After |
|----------|--------|-------|
| Step 3 prose | "Run project-specific verification commands." | "Run `just build && just test`." |

**`plugins/forge/skills/record-task/SKILL.md`**

| Location | Before | After |
|----------|--------|-------|
| Metrics Collection code block | `Go: go test -cover ./changed/package/...` | `just test` |
| Metrics Collection code block | `TypeScript: npm test -- --coverage --watchAll=false` | (removed) |
| Metrics Collection code block | `Python: pytest --cov=<module> --cov-report=term-missing` | (removed) |

**`plugins/forge/skills/improve-harness/SKILL.md`**

| Location | Before | After |
|----------|--------|-------|
| Step 4.3 prose | "Run project test suite to ensure nothing broke" | "Run `just test` to ensure nothing broke" |

**Verification** (run after all Phase 3 changes):
- `grep -rn 'project-test-command\|npx tsx\|go test\|npm test\|pytest' plugins/forge/commands/fix-bug.md plugins/forge/commands/run-tasks.md plugins/forge/agents/task-executor.md plugins/forge/agents/error-fixer.md plugins/forge/commands/execute-task.md plugins/forge/skills/record-task/SKILL.md plugins/forge/skills/improve-harness/SKILL.md` = 0 lines

---

### Phase 4: Breakdown-Tasks Templates

**`plugins/forge/skills/breakdown-tasks/templates/run-e2e-tests.md`**

Add to Implementation Notes:
```
Run: `just test-e2e --feature <slug>` (replace <slug> with current feature name from `task feature`)
```

**`plugins/forge/skills/breakdown-tasks/templates/gen-test-scripts.md`**

Add to Implementation Notes:
```
After generating spec files, run: `just e2e-verify --feature <slug>`
If exit 1 (unresolved // VERIFY: markers): task is incomplete — resolve markers before proceeding to T-test-3.
```

**`plugins/forge/skills/breakdown-tasks/templates/fix-e2e.md`**

Add post-fix verification step to Implementation Notes:
```
After fixing, verify with: `just test-e2e --feature <slug>`
All tests must pass before marking this task completed.
```

**Verification**: `grep -c 'just ' plugins/forge/skills/breakdown-tasks/templates/run-e2e-tests.md plugins/forge/skills/breakdown-tasks/templates/gen-test-scripts.md plugins/forge/skills/breakdown-tasks/templates/fix-e2e.md` >= 3 total

## Testing Strategy

### Per-Layer Test Plan

| Phase | Verification Method | Pass Condition |
|-------|--------------------|----|
| Phase 1 | `grep -c 'e2e-setup\|e2e-verify' plugins/forge/commands/init-justfile.md` | >= 4 |
| Phase 2 | `grep -c 'npx tsx\|npx playwright\|npm install' plugins/forge/skills/run-e2e-tests/SKILL.md` | = 0 |
| Phase 2 | `grep -c 'npx playwright install\|cd tests/e2e' plugins/forge/skills/gen-test-scripts/SKILL.md` | = 0 |
| Phase 3 | `grep -rn 'project-test-command\|npx tsx' plugins/forge/commands/fix-bug.md` | = 0 lines |
| Phase 3 | `grep -rn 'go test\|npm test\|pytest' plugins/forge/agents/task-executor.md plugins/forge/agents/error-fixer.md plugins/forge/skills/record-task/SKILL.md` | = 0 lines |
| Phase 4 | `grep -c 'just ' plugins/forge/skills/breakdown-tasks/templates/run-e2e-tests.md` | >= 1 |
| All | `grep -rn 'just e2e-setup\|just e2e-verify\|just test-e2e\|just test\|just build' plugins/forge/` | >= 20 lines total |

### Key Test Scenarios

1. **Happy path**: After all changes, `grep -r 'npx tsx\|cd tests/e2e && npm\|project-test-command' plugins/forge/` returns 0 results
2. **New targets present**: `init-justfile` template contains both `e2e-setup` and `e2e-verify` recipes with correct syntax
3. **Hard gate wired**: `gen-test-scripts` SKILL.md Step 4 contains `just e2e-verify` with explicit "exit 1 = skill incomplete" note
4. **Template coverage**: All 3 breakdown-tasks templates reference just commands in Implementation Notes

### Overall Coverage Target

100% — every in-scope file must pass its grep verification. No partial completion.

## Security Considerations

No security surface. All changes are to documentation files read by AI agents. No user input, no network requests, no secrets.

## PRD Coverage Map

| PRD AC | Design Component | Implementation |
|--------|-----------------|----------------|
| Story 1 AC: run-e2e-tests SKILL.md Step 1 shows only `just e2e-setup` | Phase 2: run-e2e-tests.md | Replace `cd tests/e2e && npm install` + `npx playwright install` |
| Story 2 AC: task-executor Step 3 shows `just build && just test` | Phase 3: task-executor.md | Replace language-specific examples |
| Story 3 AC: `just e2e-verify` exit 1 blocks run-e2e-tests | Phase 2: gen-test-scripts SKILL.md Step 4 | Add hard gate note after `just e2e-verify` call |
| Story 4 AC: fix-e2e template shows `just test-e2e --feature <slug>` | Phase 4: fix-e2e.md | Add post-fix verification step |
| Story 5 AC: fix-bug shows `just test`, no `<project-test-command>` | Phase 3: fix-bug.md | Replace all `<project-test-command>` occurrences |
| Story 5 AC: run-tasks Breaking Gate shows `just test` | Phase 3: run-tasks.md | Replace comment command |
| Story 5 AC: record-task Metrics shows `just test` | Phase 3: record-task SKILL.md | Replace three language examples |

## Open Questions

- [x] Should `just test-e2e` output be tee'd per-spec or combined? → Combined into one output file per feature run (run-e2e-tests parses the combined output)

## Appendix

### Alternatives Considered

| Approach | Pros | Cons | Why Not Chosen |
|----------|------|------|----------------|
| Add justfile check to every skill individually | Granular control | Duplicates the check logic 13 times | Centralize in one place: skills check once at entry point |
| Keep language examples as comments alongside just commands | Familiar to humans | Confuses agents — two commands for same action | Agents follow the first executable command; comments create ambiguity |
| Single `just e2e` mega-target (setup + verify + run) | One command | Loses granularity; can't run verify without setup | Separate concerns: setup is idempotent, verify is a gate, run is the test |
