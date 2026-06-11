---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["9"]
type: "doc.review"
scope: "all"
---

Review documentation quality for the pipeline-spec-code-alignment feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 10-submit-task-path
- [ ] Record format file path can be resolved from the subagent's working directory
- [ ] Path follows the forge distribution model conventions
- [ ] Both `record-format-coding.md` and `record-format-doc.md` paths resolve correctly


### 11-architecture-conventions
- [ ] ARCHITECTURE.md does not describe doc-scorer/doc-reviser as independent agents
- [ ] ARCHITECTURE.md scope resolution algorithm matches actual Go code
- [ ] No `forge forge` duplicate in ARCHITECTURE.md
- [ ] dispatcher-quality.md uses `just` abstractions and mentions `coding.cleanup`
- [ ] gen-contracts docs use "surfaces" terminology (not "interfaces")
- [ ] clean-code/SKILL.md references `just unit-test`
- [ ] fix-bug.md references correct just target for running tests
- [ ] execute-task.md frontmatter description is accurate


### 12-templates-guidance
- [ ] manifest-quick.md uses single unified slug placeholder
- [ ] manifest-quick.md does not reference non-existent `testing/test-cases.md`
- [ ] quick-tasks/SKILL.md documents all template placeholder mappings
- [ ] quick-tasks/SKILL.md Output Checklist is accurate for quick mode (no stage-gate claims)
- [ ] breakdown-tasks/SKILL.md has a Commit step
- [ ] gen-journeys/SKILL.md error messages match actual CLI output
- [ ] run-tests/SKILL.md `BIZ-error-reporting-001` has resolvable path
- [ ] forge-distribution.md references `/learn` not `/record-decision`/`/learn-lesson`
- [ ] prompt-template-hierarchy.md documents `<HARD-RULE>` as fourth level
- [ ] journey-contract-model.md has single canonical copy (not duplicated)


### 6-ghost-fields-cleanup
- [ ] `scope-assignment.md` deleted
- [ ] No skill doc references `interfaces` config field (grep confirms zero hits)
- [ ] No skill doc extracts `SCOPE` from claim output (use `SURFACE_KEY`/`SURFACE_TYPE`)
- [ ] No skill doc references `decision-logging.md` (use `decision-entry.md`)
- [ ] No skill doc references `test-template-dir` config field
- [ ] `db-schema.md` uses surface-type fields instead of scope
- [ ] `existing-code-split.md` references surface inference, not scope assignment


### 7-surface-type-consistency
- [ ] Surface resolution docs describe two-layer strategy (project-level shortcut + file-level query)
- [ ] `task-doc.md` either has surface fields or documents why they're absent
- [ ] No reference to `webui` surface type — all use canonical `web`
- [ ] `record-format-doc.md` lists `doc.review` and does not list `doc.eval`
- [ ] Single-surface project surface-type is non-empty (not left blank as placeholder)


### 8-dispatcher-pipeline-logic
- [ ] Post-loop message in run-tasks.md reflects actual task names (conditional on mode)
- [ ] Summary format defined in run-tasks.md
- [ ] Timeout/blocking mechanism specified in run-tasks.md
- [ ] quick.md does not claim run-tasks has knowledge extraction
- [ ] execute-task.md has explicit status branches (completed, blocked, in_progress)
- [ ] execute-task.md includes `subagent_type="forge:task-executor"` in agent call
- [ ] task-executor.md DONE format is consistent (no ambiguous field positions)
- [ ] gen-test-scripts/SKILL.md has SKIP_EVAL_GATE for Quick mode
- [ ] run-tests/SKILL.md uses source directory paths for surface detection
- [ ] No Chinese text in run-tests/SKILL.md


### 9-fix-task-template-vars
- [ ] Every fix-task creation point in run-tasks.md includes `--var SOURCE_FILES=... --var TEST_SCRIPT=... --var TEST_RESULTS=...`
- [ ] task-executor.md fix-task creation includes all three `--var` parameters
- [ ] execute-task.md fix-task creation includes all three `--var` + `--description`
- [ ] submit-task/SKILL.md recovery includes all three `--var`
- [ ] quick-tasks and breakdown-tasks SKILL.md have breaking task IT impact assessment guidance
- [ ] Fix-task grouping guidance specifies by test suite (directory), not problem type


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/pipeline-spec-code-alignment/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/pipeline-spec-code-alignment/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
