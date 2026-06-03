---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["14", "2", "5", "6", "1", "11", "15", "4", "7", "12", "3", "10", "13", "8", "9"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the global-doc-code-audit feature (quick mode).

## Acceptance Criteria
- [ ] All doc task deliverables reviewed against their acceptance criteria
- [ ] Missing or incomplete deliverables flagged with specific task ID and gap description
- [ ] Review summary produced with pass/fail per task

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-l1-pilot-audit
- [ ] All factual claims in README.md extracted (code paths, command names, config values, behavior descriptions)
- [ ] Each claim verified against actual codebase — path existence via `find`/`grep`, behavior via code reading
- [ ] Every inconsistency recorded with: file path, line range, severity (P0-P3), suggested action
- [ ] Accuracy baseline report produced: total claims examined, correct identifications, misses, false positives
- [ ] Miss rate < 20% (if ≥ 20%, report methodology adjustment recommendations and stop)


### 10-l3-lessons-batch4-audit
- [ ] All 20 target items classified (code-reference, process-standard, experience-summary)
- [ ] Each item's validity assessed using structured rules (tool change -> outdated, process contradiction -> needs-update, path invalid -> outdated, generalized conclusion -> valid, partial aging -> needs-update)
- [ ] Duplicate detection performed via topic clustering
- [ ] Every item marked as valid/outdated/duplicate/needs-update with justification
- [ ] Cross-layer influence items from L1/L2 reports checked against relevant items
- [ ] Audit report follows unified template


### 11-l3-lessons-batch5-audit
- [ ] All 20 target items classified (code-reference, process-standard, experience-summary)
- [ ] Each item's validity assessed using structured rules (tool change -> outdated, process contradiction -> needs-update, path invalid -> outdated, generalized conclusion -> valid, partial aging -> needs-update)
- [ ] Duplicate detection performed via topic clustering
- [ ] Every item marked as valid/outdated/duplicate/needs-update with justification
- [ ] Cross-layer influence items from L1/L2 reports checked against relevant items
- [ ] Audit report follows unified template


### 12-l3-lessons-batch6-audit
- [ ] All 20 target items classified (code-reference, process-standard, experience-summary)
- [ ] Each item's validity assessed using structured rules (tool change -> outdated, process contradiction -> needs-update, path invalid -> outdated, generalized conclusion -> valid, partial aging -> needs-update)
- [ ] Duplicate detection performed via topic clustering
- [ ] Every item marked as valid/outdated/duplicate/needs-update with justification
- [ ] Cross-layer influence items from L1/L2 reports checked against relevant items
- [ ] Audit report follows unified template


### 13-l3-final-batch-audit
- [ ] All 23 target items classified (code-reference, process-standard, experience-summary)
- [ ] Each item's validity assessed using structured rules (tool change -> outdated, process contradiction -> needs-update, path invalid -> outdated, generalized conclusion -> valid, partial aging -> needs-update)
- [ ] Duplicate detection performed via topic clustering
- [ ] Every item marked as valid/outdated/duplicate/needs-update with justification
- [ ] Cross-layer influence items from L1/L2 reports checked against relevant items
- [ ] Audit report follows unified template


### 14-cross-layer-verification-and-consolidation
- [ ] Cross-layer influence lists verified: every L1/L2 finding checked against relevant L3 items, every L3 finding checked against relevant L2 conventions
- [ ] Unified report produced with all findings sorted by severity (P0 → P1 → P2 → P3), each with file path, line range, severity, suggested action
- [ ] Severity counts reported: P0/P1/P2/P3 counts + L3 validity counts (valid/outdated/duplicate/needs-update)
- [ ] P0 issues flagged as release-blocking for v3.0.0; P0 report extractable within 1 working day
- [ ] All output written in English


### 15-generate-fix-tasks
- [ ] All findings converted to executable fix tasks using appropriate template: fix-type, review-type, or cross-layer-verification-type
- [ ] Knowledge base cleanup tasks (deletion/merge recommendations) marked as requiring human confirmation
- [ ] Fix tasks are self-contained: include full context, do not depend on other fix tasks
- [ ] All output written in English


### 2-l1-core-docs-audit
- [ ] All 6 target files audited with complete declaration extraction
- [ ] Each claim verified against codebase (paths via `find`/`grep`, behaviors via code reading)
- [ ] Every inconsistency recorded: file path, line range, severity (P0-P3), suggested action
- [ ] Cross-layer influence items identified and recorded for L3 reference (e.g., hook names, module paths mentioned in docs)
- [ ] Audit report follows unified template: baseline commit, issue summary, issue details, quality review


### 3-l1-official-refs-audit
- [ ] All 5 target files audited with complete declaration extraction
- [ ] Each claim verified: hook names/parameters vs code, plugin structure vs actual templates, skill definitions vs actual SKILL.md files
- [ ] Every inconsistency recorded: file path, line range, severity (P0-P3), suggested action
- [ ] Cross-layer influence items recorded for L3 reference
- [ ] Audit report follows unified template


### 4-l2-business-rules-audit
- [ ] All 4 business-rules files + CLAUDE.md audited with declaration extraction
- [ ] Each business rule claim verified against actual code enforcement (e.g., naming rules vs code constants)
- [ ] CLAUDE.md claims verified against actual project structure, file paths, and conventions
- [ ] Every inconsistency recorded: file path, line range, severity (P0-P3), suggested action
- [ ] Cross-layer influence items recorded for L3 reference
- [ ] Audit report follows unified template


### 5-l2-conventions-batch1-audit
- [ ] All 8 target files audited with declaration extraction
- [ ] Each convention claim verified: file paths via `find`, code constants via `grep`, structural rules vs actual codebase
- [ ] Every inconsistency recorded: file path, line range, severity (P0-P3), suggested action
- [ ] Cross-layer influence items recorded for L3 reference
- [ ] Audit report follows unified template


### 6-l2-conventions-batch2-audit
- [ ] All 10 target files audited with declaration extraction
- [ ] Each convention verified against codebase: naming patterns vs actual code, skill structure vs actual files, test conventions vs actual test setup
- [ ] Every inconsistency recorded: file path, line range, severity (P0-P3), suggested action
- [ ] Cross-layer influence items recorded for L3 reference
- [ ] Audit report follows unified template


### 7-l3-lessons-batch1-audit
- [ ] All 20 target items classified (code-reference, process-standard, experience-summary)
- [ ] Each item's validity assessed using structured rules (tool change -> outdated, process contradiction -> needs-update, path invalid -> outdated, generalized conclusion -> valid, partial aging -> needs-update)
- [ ] Duplicate detection performed via topic clustering
- [ ] Every item marked as valid/outdated/duplicate/needs-update with justification
- [ ] Cross-layer influence items from L1/L2 reports checked against relevant items
- [ ] Audit report follows unified template


### 8-l3-lessons-batch2-audit
- [ ] All 20 target items classified (code-reference, process-standard, experience-summary)
- [ ] Each item's validity assessed using structured rules (tool change -> outdated, process contradiction -> needs-update, path invalid -> outdated, generalized conclusion -> valid, partial aging -> needs-update)
- [ ] Duplicate detection performed via topic clustering
- [ ] Every item marked as valid/outdated/duplicate/needs-update with justification
- [ ] Cross-layer influence items from L1/L2 reports checked against relevant items
- [ ] Audit report follows unified template


### 9-l3-lessons-batch3-audit
- [ ] All 20 target items classified (code-reference, process-standard, experience-summary)
- [ ] Each item's validity assessed using structured rules (tool change -> outdated, process contradiction -> needs-update, path invalid -> outdated, generalized conclusion -> valid, partial aging -> needs-update)
- [ ] Duplicate detection performed via topic clustering
- [ ] Every item marked as valid/outdated/duplicate/needs-update with justification
- [ ] Cross-layer influence items from L1/L2 reports checked against relevant items
- [ ] Audit report follows unified template


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/global-doc-code-audit/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/global-doc-code-audit/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
