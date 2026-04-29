---
feature: "justfile-e2e-integration"
sources:
  - prd/prd-user-stories.md
  - prd/prd-spec.md
generated: "2026-04-29"
---

# Test Cases: justfile-e2e-integration

## Summary

| Type | Count |
|------|-------|
| UI   | 0     |
| API  | 0     |
| CLI  | 20    |
| **Total** | **20** |

---

## UI Test Cases

_None — this feature has no UI components._

---

## API Test Cases

_None — this feature has no API endpoints._

---

## CLI Test Cases

## TC-001: run-e2e-tests Step 1 uses just e2e-setup

- **Source**: Story 1 / AC-1
- **Type**: CLI
- **Target**: cli/run-e2e-tests
- **Test ID**: cli/run-e2e-tests/run-e2e-tests-step1-uses-just-e2e-setup
- **Pre-conditions**: `plugins/forge/skills/run-e2e-tests/SKILL.md` exists
- **Steps**:
  1. Read `plugins/forge/skills/run-e2e-tests/SKILL.md`
  2. Search for `cd tests/e2e && npm install` in Step 1
  3. Search for `npx playwright install chromium` in Step 1
  4. Search for `just e2e-setup` in Step 1
- **Expected**: `just e2e-setup` appears in Step 1; `cd tests/e2e && npm install` and `npx playwright install chromium` do not appear anywhere in the file
- **Priority**: P0

---

## TC-002: task-executor Step 3 uses just build and just test

- **Source**: Story 2 / AC-1
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/task-executor-step3-uses-just-build-and-just-test
- **Pre-conditions**: `plugins/forge/agents/task-executor/AGENT.md` (or equivalent) exists
- **Steps**:
  1. Read the task-executor agent file
  2. Locate Step 3 Full Verification section
  3. Search for `go test ./...`, `npm test`, `pytest` in Step 3
  4. Search for `just build && just test` in Step 3
- **Expected**: Step 3 explicitly contains `just build && just test`; language-specific commands (`go test ./...`, `npm test`, `pytest`) do not appear
- **Priority**: P0

---

## TC-003: just e2e-verify exits 1 when VERIFY markers present

- **Source**: Story 3 / AC-1
- **Type**: CLI
- **Target**: cli/e2e-verify
- **Test ID**: cli/e2e-verify/e2e-verify-exits-1-when-verify-markers-present
- **Pre-conditions**: A justfile with `e2e-verify` target exists; `tests/e2e/<slug>/` directory exists with at least one spec file containing `// VERIFY:` marker
- **Steps**:
  1. Create `tests/e2e/test-feature/sample.spec.ts` with a line containing `// VERIFY: check this`
  2. Run `just e2e-verify --feature test-feature`
  3. Capture exit code and stdout
- **Expected**: Exit code is 1; output includes the filename and line number of the residual `// VERIFY:` marker
- **Priority**: P0

---

## TC-004: just e2e-verify exits 0 when no VERIFY markers

- **Source**: Story 3 / AC-2
- **Type**: CLI
- **Target**: cli/e2e-verify
- **Test ID**: cli/e2e-verify/e2e-verify-exits-0-when-no-verify-markers
- **Pre-conditions**: A justfile with `e2e-verify` target exists; `tests/e2e/<slug>/` directory exists with spec files containing no `// VERIFY:` markers
- **Steps**:
  1. Create `tests/e2e/test-feature/sample.spec.ts` with no `// VERIFY:` markers
  2. Run `just e2e-verify --feature test-feature`
  3. Capture exit code and stdout
- **Expected**: Exit code is 0; output is `"OK: no unresolved // VERIFY: markers"`
- **Priority**: P0

---

## TC-005: fix-e2e template uses just test-e2e for post-fix verification

- **Source**: Story 4 / AC-1
- **Type**: CLI
- **Target**: cli/fix-e2e
- **Test ID**: cli/fix-e2e/fix-e2e-template-uses-just-test-e2e-for-post-fix-verification
- **Pre-conditions**: `plugins/forge/skills/breakdown-tasks/templates/fix-e2e.md` exists
- **Steps**:
  1. Read `plugins/forge/skills/breakdown-tasks/templates/fix-e2e.md`
  2. Locate Implementation Notes section
  3. Search for `npx tsx` in the file
  4. Search for `just test-e2e --feature <slug>` in Implementation Notes
- **Expected**: Implementation Notes contains `just test-e2e --feature <slug>` for post-fix verification; `npx tsx` does not appear anywhere in the file
- **Priority**: P0

---

## TC-006: fix-bug uses just test not project-test-command placeholder

- **Source**: Story 5 / AC-1
- **Type**: CLI
- **Target**: cli/fix-bug
- **Test ID**: cli/fix-bug/fix-bug-uses-just-test-not-project-test-command
- **Pre-conditions**: `plugins/forge/commands/fix-bug` (or equivalent) file exists
- **Steps**:
  1. Read the fix-bug command file
  2. Search for `<project-test-command>` in the file
  3. Search for language-specific commands: `go test`, `npm test`, `pytest`
  4. Search for `just test` in the test verification step
- **Expected**: `just test` appears in the test verification step; `<project-test-command>` placeholder and language-specific test commands do not appear
- **Priority**: P0

---

## TC-007: run-tasks Breaking Gate uses just test

- **Source**: Story 5 / AC-2
- **Type**: CLI
- **Target**: cli/run-tasks
- **Test ID**: cli/run-tasks/run-tasks-breaking-gate-uses-just-test
- **Pre-conditions**: `plugins/forge/skills/run-tasks` (or equivalent) file exists
- **Steps**:
  1. Read the run-tasks skill/command file
  2. Locate Breaking Gate section
  3. Search for `npm test`, `go test` in Breaking Gate section
  4. Search for `just test` in Breaking Gate section
- **Expected**: Breaking Gate section explicitly contains `just test`; `npm test` and `go test` do not appear in that section
- **Priority**: P0

---

## TC-008: record-task Metrics Collection uses just test

- **Source**: Story 5 / AC-3
- **Type**: CLI
- **Target**: cli/record-task
- **Test ID**: cli/record-task/record-task-metrics-collection-uses-just-test
- **Pre-conditions**: `plugins/forge/skills/record-task/SKILL.md` exists
- **Steps**:
  1. Read `plugins/forge/skills/record-task/SKILL.md`
  2. Locate Metrics Collection section
  3. Search for `go test -cover ./...`, `npm test -- --coverage`, `pytest --cov=` in the file
  4. Search for `just test` in Metrics Collection section
- **Expected**: Language examples in Metrics Collection are unified as `just test`; `go test -cover ./...`, `npm test -- --coverage`, `pytest --cov=...` do not appear
- **Priority**: P1

---

## TC-009: just e2e-setup exits 1 when package.json missing

- **Source**: Spec Section 5.1
- **Type**: CLI
- **Target**: cli/e2e-setup
- **Test ID**: cli/e2e-setup/e2e-setup-exits-1-when-package-json-missing
- **Pre-conditions**: A justfile with `e2e-setup` target exists; `tests/e2e/package.json` does not exist
- **Steps**:
  1. Ensure `tests/e2e/package.json` does not exist
  2. Run `just e2e-setup`
  3. Capture exit code and stdout
- **Expected**: Exit code is 1; output is `"Error: tests/e2e/package.json not found"`
- **Priority**: P0

---

## TC-010: just e2e-setup exits 0 with OK message when deps ready

- **Source**: Spec Section 5.1
- **Type**: CLI
- **Target**: cli/e2e-setup
- **Test ID**: cli/e2e-setup/e2e-setup-exits-0-with-ok-message-when-deps-ready
- **Pre-conditions**: A justfile with `e2e-setup` target exists; `tests/e2e/package.json` exists; `tests/e2e/node_modules` exists (deps already installed)
- **Steps**:
  1. Ensure `tests/e2e/package.json` and `tests/e2e/node_modules` exist
  2. Run `just e2e-setup`
  3. Capture exit code and stdout
- **Expected**: Exit code is 0; output includes `"OK: e2e dependencies ready"`
- **Priority**: P0

---

## TC-011: just e2e-verify exits 1 when feature flag missing

- **Source**: Spec Section 5.1
- **Type**: CLI
- **Target**: cli/e2e-verify
- **Test ID**: cli/e2e-verify/e2e-verify-exits-1-when-feature-flag-missing
- **Pre-conditions**: A justfile with `e2e-verify` target exists
- **Steps**:
  1. Run `just e2e-verify` (without `--feature` argument)
  2. Capture exit code and stdout
- **Expected**: Exit code is 1; output includes usage hint indicating `--feature <slug>` is required
- **Priority**: P1

---

## TC-012: just e2e-verify outputs file and line number for residual markers

- **Source**: Spec Section 5.1
- **Type**: CLI
- **Target**: cli/e2e-verify
- **Test ID**: cli/e2e-verify/e2e-verify-outputs-file-and-line-for-residual-markers
- **Pre-conditions**: A justfile with `e2e-verify` target exists; `tests/e2e/my-feature/login.spec.ts` exists with `// VERIFY:` on line 42
- **Steps**:
  1. Create `tests/e2e/my-feature/login.spec.ts` with `// VERIFY: implement this` on a known line
  2. Run `just e2e-verify --feature my-feature`
  3. Capture stdout
- **Expected**: Output includes the filename (`login.spec.ts`) and the line number of the `// VERIFY:` marker
- **Priority**: P1

---

## TC-013: Skills prompt to run init-justfile when justfile missing

- **Source**: Spec Section 5.3
- **Type**: CLI
- **Target**: cli/run-e2e-tests
- **Test ID**: cli/run-e2e-tests/skills-prompt-init-justfile-when-justfile-missing
- **Pre-conditions**: `plugins/forge/skills/run-e2e-tests/SKILL.md` exists
- **Steps**:
  1. Read `plugins/forge/skills/run-e2e-tests/SKILL.md`
  2. Search for `ls justfile` check or equivalent justfile existence check
  3. Search for `/init-justfile` prompt instruction
- **Expected**: Skill includes a check for justfile existence and instructs the agent to prompt the user to run `/init-justfile` if not found, then stop
- **Priority**: P1

---

## TC-014: gen-test-scripts Step 4 uses just e2e-verify

- **Source**: Spec Section 5.2 / Story 3
- **Type**: CLI
- **Target**: cli/gen-test-scripts
- **Test ID**: cli/gen-test-scripts/gen-test-scripts-step4-uses-just-e2e-verify
- **Pre-conditions**: `plugins/forge/skills/gen-test-scripts/SKILL.md` exists
- **Steps**:
  1. Read `plugins/forge/skills/gen-test-scripts/SKILL.md`
  2. Locate Step 4 (VERIFY check section)
  3. Search for `grep -r '// VERIFY:'` in the file
  4. Search for `just e2e-verify --feature` in Step 4
- **Expected**: Step 4 uses `just e2e-verify --feature <slug>`; raw `grep -r '// VERIFY:'` command does not appear
- **Priority**: P0

---

## TC-015: error-fixer uses just build and just test

- **Source**: Spec Section 5.2
- **Type**: CLI
- **Target**: cli/error-fixer
- **Test ID**: cli/error-fixer/error-fixer-uses-just-build-and-just-test
- **Pre-conditions**: `plugins/forge/agents/error-fixer` (or equivalent) file exists
- **Steps**:
  1. Read the error-fixer agent file
  2. Search for `go build ./...`, `go vet ./...`, `go test -race -cover ./...` in the file
  3. Search for `npm run build && npm test` in the file
  4. Search for `pytest --cov` in the file
  5. Search for `just build && just test` in the verification step
- **Expected**: Verification step contains `just build && just test`; language-specific build/test commands do not appear
- **Priority**: P0

---

## TC-016: execute-task Step 3 uses just build and just test

- **Source**: Spec Section 5.2
- **Type**: CLI
- **Target**: cli/execute-task
- **Test ID**: cli/execute-task/execute-task-step3-uses-just-build-and-just-test
- **Pre-conditions**: `plugins/forge/skills/execute-task` (or equivalent) file exists
- **Steps**:
  1. Read the execute-task skill/command file
  2. Locate Step 3 description
  3. Search for language-specific verification commands in Step 3
  4. Search for `just build && just test` in Step 3
- **Expected**: Step 3 description contains `just build && just test`; language-specific commands do not appear
- **Priority**: P0

---

## TC-017: improve-harness uses just test

- **Source**: Spec Section 5.2
- **Type**: CLI
- **Target**: cli/improve-harness
- **Test ID**: cli/improve-harness/improve-harness-uses-just-test
- **Pre-conditions**: `plugins/forge/skills/improve-harness/SKILL.md` exists
- **Steps**:
  1. Read `plugins/forge/skills/improve-harness/SKILL.md`
  2. Locate Step 4.3 description
  3. Search for "Run project test suite" or similar raw description
  4. Search for `just test` in Step 4.3
- **Expected**: Step 4.3 contains `just test`; generic "run project test suite" description without just command does not appear
- **Priority**: P1

---

## TC-018: init-justfile generates e2e-setup target

- **Source**: Spec Section 5.1
- **Type**: CLI
- **Target**: cli/init-justfile
- **Test ID**: cli/init-justfile/init-justfile-generates-e2e-setup-target
- **Pre-conditions**: `plugins/forge/skills/init-justfile` (or equivalent) file exists
- **Steps**:
  1. Read the init-justfile skill/command file and its templates
  2. Search for `e2e-setup` recipe definition in the justfile template
  3. Verify the recipe includes idempotent npm install logic and playwright install
- **Expected**: The generated justfile template contains an `e2e-setup` recipe that idempotently installs npm deps and playwright chromium
- **Priority**: P0

---

## TC-019: init-justfile generates e2e-verify target

- **Source**: Spec Section 5.1
- **Type**: CLI
- **Target**: cli/init-justfile
- **Test ID**: cli/init-justfile/init-justfile-generates-e2e-verify-target
- **Pre-conditions**: `plugins/forge/skills/init-justfile` (or equivalent) file exists
- **Steps**:
  1. Read the init-justfile skill/command file and its templates
  2. Search for `e2e-verify` recipe definition in the justfile template
  3. Verify the recipe accepts `--feature <slug>` parameter and scans for `// VERIFY:` markers
- **Expected**: The generated justfile template contains an `e2e-verify` recipe that accepts `--feature <slug>`, scans `tests/e2e/<slug>/` for `// VERIFY:` markers, exits 0 if none found, exits 1 with file/line info if found
- **Priority**: P0

---

## TC-020: just e2e-setup is idempotent

- **Source**: Spec Section 5.1
- **Type**: CLI
- **Target**: cli/e2e-setup
- **Test ID**: cli/e2e-setup/e2e-setup-is-idempotent
- **Pre-conditions**: A justfile with `e2e-setup` target exists; `tests/e2e/package.json` exists; `tests/e2e/node_modules` already exists
- **Steps**:
  1. Ensure `tests/e2e/node_modules` already exists (deps previously installed)
  2. Run `just e2e-setup` twice in succession
  3. Capture exit codes and outputs for both runs
- **Expected**: Both runs exit 0 with `"OK: e2e dependencies ready"`; second run does not re-run `npm install` (node_modules check skips install)
- **Priority**: P1

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Story 1 / AC-1 | CLI | cli/run-e2e-tests | P0 |
| TC-002 | Story 2 / AC-1 | CLI | cli/task-executor | P0 |
| TC-003 | Story 3 / AC-1 | CLI | cli/e2e-verify | P0 |
| TC-004 | Story 3 / AC-2 | CLI | cli/e2e-verify | P0 |
| TC-005 | Story 4 / AC-1 | CLI | cli/fix-e2e | P0 |
| TC-006 | Story 5 / AC-1 | CLI | cli/fix-bug | P0 |
| TC-007 | Story 5 / AC-2 | CLI | cli/run-tasks | P0 |
| TC-008 | Story 5 / AC-3 | CLI | cli/record-task | P1 |
| TC-009 | Spec Section 5.1 | CLI | cli/e2e-setup | P0 |
| TC-010 | Spec Section 5.1 | CLI | cli/e2e-setup | P0 |
| TC-011 | Spec Section 5.1 | CLI | cli/e2e-verify | P1 |
| TC-012 | Spec Section 5.1 | CLI | cli/e2e-verify | P1 |
| TC-013 | Spec Section 5.3 | CLI | cli/run-e2e-tests | P1 |
| TC-014 | Spec Section 5.2 / Story 3 | CLI | cli/gen-test-scripts | P0 |
| TC-015 | Spec Section 5.2 | CLI | cli/error-fixer | P0 |
| TC-016 | Spec Section 5.2 | CLI | cli/execute-task | P0 |
| TC-017 | Spec Section 5.2 | CLI | cli/improve-harness | P1 |
| TC-018 | Spec Section 5.1 | CLI | cli/init-justfile | P0 |
| TC-019 | Spec Section 5.1 | CLI | cli/init-justfile | P0 |
| TC-020 | Spec Section 5.1 | CLI | cli/e2e-setup | P1 |
