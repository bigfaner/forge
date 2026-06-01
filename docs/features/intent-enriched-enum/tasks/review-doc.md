---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["1", "2", "3", "4", "5"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the intent-enriched-enum feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-update-brainstorm-intent-mapping
- [ ] brainstorm/SKILL.md Step 4.5 intent mapping table contains exactly 6 values: new-feature, enhancement, refactor, cleanup, fix, doc
- [ ] Fix heuristic logic removed entirely — `coding.fix` always maps to `fix` intent without runtime inference
- [ ] AskUserQuestion intent selection offers all 6 values as structured options
- [ ] brainstorm/templates/proposal.md intent valid values comment lists all 6 values
- [ ] `coding.feature` → `new-feature` and `coding.enhancement` → `enhancement` exist as independent mapping paths (not merged)


### 2-update-write-prd-pipeline-config
- [ ] write-prd/SKILL.md uses Pipeline Configuration table with 6 rows (one per intent) and 6 columns (PRD Format, User Stories, API Handbook, Test Pipeline, Security Review, at minimum)
- [ ] Override Signals table exists with 5 signal types: API 变更, 用户可见行为, 安全相关, 性能相关, 数据迁移
- [ ] Override trigger generates `<!-- Override: ... -->` comment in PRD output (e.g., `<!-- Override: API handbook enabled by signal "接口变更" -->`)
- [ ] Enhancement intent produces Simplified PRD format (Background + Goals + Test Pipeline), skipping User Stories
- [ ] Doc intent produces Minimal PRD format (title + goals + scope only)
- [ ] write-prd/rules/self-check.md intent-gated checks reference all 6 intent values
- [ ] Existing new-feature, refactor, cleanup pipeline artifacts unchanged from pre-modification behavior


### 3-update-tech-design-pipeline-config
- [ ] tech-design/SKILL.md uses Pipeline Configuration table with 6 rows (identical structure to write-prd)
- [ ] Override Signals table matches write-prd's exactly (5 signal types, same keywords, same override actions)
- [ ] Override trigger generates `<!-- Override: ... -->` comment in tech-design output
- [ ] tech-design/rules/design-quality-checks.md intent-gated checks reference all 6 intent values
- [ ] Existing new-feature, refactor, cleanup pipeline artifacts unchanged from pre-modification behavior


### 4-update-breakdown-tasks-intent-propagation
- [ ] Intent Propagation uses strict 1:1 mapping: new-feature→coding.feature, enhancement→coding.enhancement, refactor→coding.refactor, cleanup→coding.cleanup, fix→coding.fix, doc→doc
- [ ] Type Assignment table entry for `coding.fix` updated to: "可由 fix intent 自动映射，但不可通过 `forge task add` CLI 手动创建"
- [ ] `doc` intent resolves to `doc` task type without sub-type distinction (doc.consolidate/doc.drift unified under doc umbrella)


### 5-update-quick-tasks-intent-propagation
- [ ] Intent Propagation uses strict 1:1 mapping consistent with breakdown-tasks: new-feature→coding.feature, enhancement→coding.enhancement, refactor→coding.refactor, cleanup→coding.cleanup, fix→coding.fix, doc→doc
- [ ] Mapping table matches breakdown-tasks/SKILL.md exactly


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/intent-enriched-enum/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/intent-enriched-enum/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
