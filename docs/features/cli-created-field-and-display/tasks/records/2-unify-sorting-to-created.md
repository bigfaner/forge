---
status: "completed"
started: "2026-05-21 00:48"
completed: "2026-05-21 00:59"
time_spent: "~11m"
---

# Task Record: 2 统一 feature/lesson 排序为 created 降序

## Summary
Unified feature and lesson sorting to use frontmatter created field descending with mtime fallback. Feature list now parses created from manifest.md frontmatter. Lesson Discover now sorts by Created field instead of modTime. Proposal sorting unchanged (already used created).

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/feature.go
- forge-cli/internal/cmd/feature_test.go
- forge-cli/pkg/lesson/lesson.go
- forge-cli/pkg/lesson/lesson_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- Created field left empty when absent from frontmatter (no synthetic mtime-based value) — display shows empty but sorting still works via mtime fallback
- Sorting uses three-tier comparison: both have created -> lexicographic desc, only one has created -> it wins, neither has created -> mtime desc
- Same sorting pattern applied to both feature and lesson for consistency

## Test Results
- **Tests Executed**: Yes
- **Passed**: 18
- **Failed**: 0
- **Coverage**: 84.6%

## Acceptance Criteria
- [x] forge feature list 按 manifest frontmatter created 字段降序排列（缺少 created 时 fallback 到 mtime）
- [x] forge lesson 按 frontmatter created 字段降序排列（缺少 created 时 fallback 到 mtime）
- [x] forge proposal 排序逻辑不变（验证仍为 created 降序）
- [x] 缺少 created 的旧文档 fallback 到 mtime，不报错、不跳过

## Notes
Existing TestDiscover_NoFrontmatterFallsBackToModTime renamed to TestDiscover_NoFrontmatterCreatedIsEmpty to reflect that Created is now empty when no frontmatter exists. Replaced TestFeatureList_SortedByManifestMtime with three new tests covering created-descending, mtime-fallback, and created-over-mtime-priority.
