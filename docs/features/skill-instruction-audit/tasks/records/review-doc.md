---
status: "blocked"
started: "2026-05-28 23:12"
completed: "N/A"
time_spent: ""
---

# Task Record: T-review-doc Review Documentation Quality

## Summary
Reviewed all 7 doc task ACs against actual skill/command files. 36/38 ACs pass. 2 ACs fail in gen-contracts/SKILL.md: (1) description has 3 sentences, AC requires <=1; (2) Reference section still has 'Section X.Y' references. Cannot fix because SCOPE CONSTRAINT limits modifications to docs/ directory only.

## Changes

### Files Created
无

### Files Modified
无

### Key Decisions
无

## Document Metrics
36/38 ACs pass (94.7%), 2 ACs fail in gen-contracts/SKILL.md

## Referenced Documents
- docs/proposals/skill-instruction-audit/proposal.md
- docs/features/skill-instruction-audit/tasks/review-doc.md

## Review Status
final

## Acceptance Criteria
- [x] [1.1] submit-task/SKILL.md has no 'What .* Does' section; command and exit code remain
- [x] [1.2] breakdown-tasks/SKILL.md and quick-tasks/SKILL.md have no 'Auto-generated tasks by forge task index' block; Step 5 command remains
- [x] [1.3] execute-task.md retains field name list but removes per-field semantic explanations and example values
- [x] [1.4] run-tasks.md retains field name list; no 'Subagent calls forge prompt get-by-task-id internally'
- [x] [1.5] No output contract field names lost (SURFACE_KEY, SURFACE_TYPE, TASK_ID, FILE, MAIN_SESSION remain)
- [x] [2.1] gen-journeys/SKILL.md has no example output blocks for forge surfaces; Exit Code table remains
- [x] [2.2] run-tests/SKILL.md has no 'segment prefix matching'; command remains
- [x] [2.3] eval/SKILL.md has no repeated tool usage explanations
- [x] [2.4] forensic/SKILL.md has no go build or ~/.zcode-forge-cli references
- [x] [2.5] ui-design/SKILL.md config check is natural language, not bash script
- [x] [2.6] quick.md has no behavioral descriptions of run-tasks or brainstorm internals
- [x] [3.1] test-guide/SKILL.md E-I has <=4 items; each passes constraint-level audit
- [x] [3.2] eval/SKILL.md has 1 concise E-I rule
- [x] [3.3] init-justfile/SKILL.md Notes has no E-I overlap
- [x] [3.4] clean-code/SKILL.md has 1 'preserve scope' statement
- [x] [3.5] deep-research/SKILL.md has no key points duplicating template
- [x] [4.1] 4 commands body doesn't start with frontmatter sentence
- [x] [4.2] learn/SKILL.md body doesn't duplicate frontmatter
- [ ] [4.3] gen-contracts/SKILL.md description <=1 sentence
- [x] [4.4] run-tests/SKILL.md has single merged opening
- [x] [4.5] extract-design-md/SKILL.md has no Overview paragraphs; only rules reference
- [x] [5.1] tech-design Process Flow: 0->1->...->11 no gaps
- [x] [5.2] run-tests Step 5 does not reference 'Convention loaded in Step 0'
- [x] [5.3] write-prd has no decimal step numbers
- [x] [5.4] quick-tasks has no decimal step numbers
- [x] [5.5] breakdown-tasks Step 4b is informational note, not numbered step
- [ ] [5.6] gen-contracts has no 'Section X.Y' references
- [x] [6.1] quick.md Non-zero fallback = 'show gate'
- [x] [6.2] execute-task.md defines MAIN_SESSION; Step 1.5 verify is self-contained
- [x] [6.3] run-tasks.md defines successful = STATUS completed; defines T-test-run; has slug failure path
- [x] [6.4] gen-journeys has case-insensitive 'scenario' matching rule
- [x] [6.5] submit-task has objective type reclassification criteria
- [x] [6.6] fix-bug has no duplicate between E-I and HARD-GATE
- [x] [7.1] consolidate-specs Steps 9-11 specify scan directories
- [x] [7.2] gen-sitemap has no @latest vs pinning contradiction
- [x] [7.3] clean-code has clear config file exception or none
- [x] [7.4] deep-research has explicit 'wait for user review' between report and conversion ask
- [x] [7.5] ui-design has web+mobile output rule; eval-skip option is clearly labeled
- [x] [7.6] simplify-skill states whether it targets user skills or plugin skills
- [x] [7.7] forensic explains how to obtain project-hash path

## Notes
2 failing ACs both in gen-contracts/SKILL.md: (1) frontmatter description has 3 sentences instead of <=1; (2) Reference section still contains 'Section X.Y' style references (e.g., 'Section 1.1', 'Section 1.2'). These are in plugins/ directory which is outside the SCOPE CONSTRAINT (docs/ only). A follow-up task is needed to fix gen-contracts/SKILL.md.
