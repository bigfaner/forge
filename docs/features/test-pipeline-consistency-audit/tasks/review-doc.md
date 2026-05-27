---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["9"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the test-pipeline-consistency-audit feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 10-fix-architecture-md
- [ ] Quick 模式流程图不含 gen-contracts/gen-scripts
- [ ] 任务 ID 为描述性名称（非 `T-test-1~5`）
- [ ] 不含 `T-test-promote` 条目
- [ ] "profile type" 已替换为 "Convention" / "surface type"
- [ ] "profile 路由" 已替换为 "Convention 路由"
- [ ] 第 305 行并行执行描述与 `autogen.go` 实际代码一致


### 11-sync-overview-workflow
- [ ] `OVERVIEW.md` 和 `OVERVIEW.zh.md` 中 "e2e" 泛用替换为 surface-specific 术语
- [ ] `WORKFLOW.md` 和 `WORKFLOW.zh.md` 中 "e2e" 泛用替换为 surface-specific 术语
- [ ] graduation/staging 描述已更新为 tag-based promotion
- [ ] "profile" 旧术语引用已移除
- [ ] `grep -rn "tests/e2e" forge-cli/docs/OVERVIEW.zh.md forge-cli/docs/WORKFLOW.zh.md` 返回 0 结果


### 12-update-surface-test-type-model
- [ ] 第 73 行 recipe 命名已更新为 `<surface-key>-test`
- [ ] 第 85 行 recipe 命名已更新为 `<surface-key>-test`
- [ ] 第 107 行多 surface recipe 命名已同步更新
- [ ] NFR1 向后兼容要求标记为 v3.0.0 已覆盖
- [ ] 仅修改 recipe 命名部分，测试类型映射和术语定义不变


### 7-update-gen-scripts-docs
- [ ] `gen-contracts/` 下所有 Skill 文档中 "e2e 测试管道" 替换为 "Forge 测试管道"
- [ ] `gen-test-scripts/` 下所有 Skill 文档中 `tests/e2e/` 旧路径替换为 `tests/<journey>/`
- [ ] `gen-test-scripts/rules/step-1-contract-loading.md` 中 `tests/e2e/step1_test.go` 示例路径已更新
- [ ] `gen-test-scripts/rules/convention-guide.md` 中 "e2e tests" 引用已替换
- [ ] `gen-contracts/rules/journey-contract-model.md` 第 159 行 "language profile" 未被修改


### 8-fix-quick-tasks-order
- [ ] `quick-tasks/SKILL.md` 执行顺序已修正为 `gen-journeys → run-test`
- [ ] `breakdown-tasks/SKILL.md` 和 `quick-tasks/SKILL.md` 中 "Integration Test Impact Assessment" → "Test Impact Assessment"
- [ ] `breakdown-tasks/SKILL.md` 执行顺序未被修改（已经正确）


### 9-update-other-skill-docs
- [ ] `commands/fix-bug.md` 中 `tests/e2e/features/` → `tests/<journey>/`
- [ ] `commands/run-tasks.md` 中 `T-test-verify-regression` 和 "e2e verification" 引用已清理
- [ ] test-guide（含 `rules/draft-generation.md` 和 `rules/pattern-extraction.md`）术语修正
- [ ] `submit-task/data/record-format-test.md` 删除 `test.verify-regression` 类型，更新 `tests/e2e/` 示例路径
- [ ] `gen-sitemap` 配置文件从 `e2e-config.yaml` 重命名为 `test-config.yaml`，SKILL.md 路径引用更新
- [ ] `consolidate-specs/SKILL.md` 第 22 行 "e2e tests are promoted" → "all tests pass"
- [ ] `init-justfile/SKILL.md` 中 `tests/e2e/` 示例路径更新
- [ ] `init-justfile/templates/` 下 6 个 justfile 模板中 `tests/e2e/` 路径更新
- [ ] `run-tests/rules/test-isolation.md` 中 4 处 `tests/e2e/` 路径引用更新
- [ ] Convention 文件中 build tag 与 surface 类型对齐


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/test-pipeline-consistency-audit/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/test-pipeline-consistency-audit/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
