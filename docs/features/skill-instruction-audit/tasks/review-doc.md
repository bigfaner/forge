---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["1", "2", "3", "4", "5", "6", "7"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the skill-instruction-audit feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-delete-cli-descriptions-pipeline

- [ ] `submit-task/SKILL.md` has no "What .* Does" section; `forge task submit` command and exit code contract remain
- [ ] `breakdown-tasks/SKILL.md` and `quick-tasks/SKILL.md` have no "Auto-generated tasks by forge task index" block; Step 5 command remains
- [ ] `execute-task.md` retains field name list (TASK_ID, FILE, etc.) but removes per-field semantic explanations and example values
- [ ] `run-tasks.md` retains field name list; no "Subagent calls forge prompt get-by-task-id internally"
- [ ] No output contract field names lost (SURFACE_KEY, SURFACE_TYPE, TASK_ID, FILE, MAIN_SESSION remain)


### 2-delete-cli-descriptions-remaining

- [ ] `gen-journeys/SKILL.md` has no example output blocks for forge surfaces; Exit Code table remains
- [ ] `run-tests/SKILL.md` has no "segment prefix matching"; command remains
- [ ] `eval/SKILL.md` has no repeated tool usage explanations
- [ ] `forensic/SKILL.md` has no go build or ~/.zcode-forge-cli references
- [ ] `ui-design/SKILL.md` config check is natural language, not bash script
- [ ] `quick.md` has no behavioral descriptions of run-tasks or brainstorm internals


### 3-remove-ei-redundancy

- [ ] `test-guide/SKILL.md` E-I has ≤4 items; each passes constraint-level audit
- [ ] `eval/SKILL.md` has 1 concise E-I rule
- [ ] `init-justfile/SKILL.md` Notes has no E-I overlap
- [ ] `clean-code/SKILL.md` has 1 "preserve scope" statement
- [ ] `deep-research/SKILL.md` has no key points duplicating template


### 4-remove-other-redundancy

- [ ] 4 commands (clean-code, git-commit, git-checkout, init-forge) body doesn't start with frontmatter sentence
- [ ] `learn/SKILL.md` body doesn't duplicate frontmatter
- [ ] `gen-contracts/SKILL.md` description ≤1 sentence
- [ ] `run-tests/SKILL.md` has single merged opening
- [ ] `extract-design-md/SKILL.md` has no Overview paragraphs; only rules reference


### 5-fix-numbering-references

- [ ] `tech-design` Process Flow: 0→1→...→8→9→10→11 no gaps
- [ ] `run-tests` Step 5 does not reference "Convention loaded in Step 0"
- [ ] `write-prd` has no decimal step numbers
- [ ] `quick-tasks` has no decimal step numbers
- [ ] `breakdown-tasks` Step 4b is informational note, not numbered step
- [ ] `gen-contracts` has no "Section X.Y" references


### 6-fix-pipeline-ambiguity

- [ ] `quick.md` Non-zero fallback = "show gate"
- [ ] `execute-task.md` defines MAIN_SESSION; Step 1.5 verify is self-contained
- [ ] `run-tasks.md` defines successful = STATUS completed; defines T-test-run; has slug failure path
- [ ] `gen-journeys` has case-insensitive "scenario" matching rule
- [ ] `submit-task` has objective type reclassification criteria
- [ ] `fix-bug` has no duplicate between E-I and HARD-GATE


### 7-fix-remaining-ambiguity

- [ ] `consolidate-specs` Steps 9-11 specify scan directories (docs/business-rules/, docs/conventions/)
- [ ] `gen-sitemap` has no @latest vs pinning contradiction
- [ ] `clean-code` has clear config file exception or none
- [ ] `deep-research` has explicit "wait for user review" between report presentation and proposal conversion ask
- [ ] `ui-design` has web+mobile output rule; eval-skip option is clearly labeled
- [ ] `simplify-skill` states whether it targets user skills or plugin skills
- [ ] `forensic` explains how to obtain project-hash path


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/skill-instruction-audit/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/skill-instruction-audit/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
