---
id: "4"
title: "Delete dead templates and bump version"
priority: "P2"
estimated_time: "7min"
dependencies: []
scope: "all"
breaking: false
type: "cleanup"
mainSession: false
---

# 4: Delete dead templates and bump version

## Description
Delete 14 dead template files from `breakdown-tasks/templates/` and 9 from `quick-tasks/templates/`. These files are never read by any SKILL.md or Go code — the CLI generates all test/gate/fix/validate task content programmatically. Also bump minor version in `scripts/version.txt`.

### breakdown-tasks/templates/ — delete 14 files
- `gen-test-cases.md`, `eval-test-cases.md`, `gen-test-scripts.md`, `run-e2e-tests.md`, `graduate-tests.md`, `verify-regression.md`, `consolidate-specs.md`
- `gate-task.md`, `phase-summary-task.md`, `fix-task.md`
- `validate-code-task.md`, `validate-ux-task.md`
- `index.json`, `index.schema.json`

Keep: `task.md`, `task-doc.md`, `manifest-update-tasks.md`

### quick-tasks/templates/ — delete 9 files
- `quick-test-cases.md`, `quick-gen-scripts.md`, `quick-run-tests.md`, `quick-graduate.md`, `quick-verify-regression.md`
- `validate-code-task.md`, `validate-ux-task.md`
- `index.json`, `index.schema.json`

Keep: `task.md`, `task-doc.md`, `manifest-quick.md`

## Reference Files
- `docs/proposals/task-type-id-redesign/proposal.md` — Source proposal
- `plugins/forge/skills/breakdown-tasks/templates/` — Dead template files
- `plugins/forge/skills/quick-tasks/templates/` — Dead template files
- `scripts/version.txt` — Version file

## Acceptance Criteria
- [ ] `breakdown-tasks/templates/` contains exactly 3 files: `task.md`, `task-doc.md`, `manifest-update-tasks.md`
- [ ] `quick-tasks/templates/` contains exactly 3 files: `task.md`, `task-doc.md`, `manifest-quick.md`
- [ ] `scripts/version.txt` minor version bumped

## Hard Rules
- Verify no Go code or SKILL.md references deleted files before deletion
- Version bump follows semver: minor bump for new features + breaking type name changes

## Implementation Notes
- These templates were superseded when the CLI began generating task content programmatically
- `index.json` and `index.schema.json` in template dirs are not consumed by any code
