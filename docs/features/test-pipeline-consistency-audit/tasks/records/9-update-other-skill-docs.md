---
status: "completed"
started: "2026-05-27 15:38"
completed: "2026-05-27 15:53"
time_spent: "~15m"
---

# Task Record: 9 更新其他 Skill/Command 文档中的旧术语和路径

## Summary
批量更新 Skill/Command 文档中的旧术语和路径：fix-bug 路径更新、run-tasks 旧类型清理、submit-task record 删除 verify-regression、gen-sitemap 配置文件重命名、consolidate-specs 术语修正、init-justfile 路径更新、6 个 justfile 模板路径更新、run-tests/test-isolation 路径更新、test-guide build tag 对齐

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/fix-bug.md
- plugins/forge/commands/run-tasks.md
- plugins/forge/skills/run-tests/rules/test-isolation.md
- plugins/forge/skills/submit-task/data/record-format-test.md
- plugins/forge/skills/gen-sitemap/SKILL.md
- plugins/forge/skills/gen-sitemap/templates/test-config.yaml
- plugins/forge/skills/consolidate-specs/SKILL.md
- plugins/forge/skills/init-justfile/SKILL.md
- plugins/forge/skills/init-justfile/templates/go.just
- plugins/forge/skills/init-justfile/templates/rust.just
- plugins/forge/skills/init-justfile/templates/node.just
- plugins/forge/skills/init-justfile/templates/python.just
- plugins/forge/skills/init-justfile/templates/mixed.just
- plugins/forge/skills/init-justfile/templates/generic.just
- plugins/forge/skills/test-guide/rules/draft-generation.md
- plugins/forge/skills/test-guide/rules/pattern-extraction.md

### Key Decisions
无

## Document Metrics
16 files modified, 1 file renamed (e2e-config.yaml -> test-config.yaml)

## Referenced Documents
- docs/proposals/test-pipeline-consistency-audit/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] fix-bug.md tests/e2e/features/ -> tests/<journey>/
- [x] run-tasks.md T-test-verify-regression and e2e verification cleaned
- [x] test-guide draft-generation and pattern-extraction terminology fixed
- [x] submit-task record-format-test.md verify-regression removed, paths updated
- [x] gen-sitemap e2e-config.yaml renamed to test-config.yaml, SKILL.md updated
- [x] consolidate-specs SKILL.md e2e tests are promoted -> all tests pass
- [x] init-justfile SKILL.md tests/e2e/ example paths updated
- [x] 6 justfile templates tests/e2e/ paths updated
- [x] test-isolation.md 4 tests/e2e/ path references updated
- [x] Convention build tag table aligned with surface type

## Notes
All tests/e2e references eliminated from modified files. Build tags aligned to <surface>-<type> pattern.
