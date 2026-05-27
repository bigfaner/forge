---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["6"]
type: "doc.review"
scope: "all"
---

Review documentation quality for the surface-test-type-model feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-test-type-concept-doc

- [ ] 文档包含 5 种 surface 的 Test Type 名称（EN + CN）、验证维度和执行模型
- [ ] 文档包含分类标准声明：一级分类键（Surface）和二级属性（测试范围：功能测试/端到端测试）
- [ ] 文档包含每种测试类型的语义定义（如 CLI 功能测试 vs Web 端到端测试的区别）
- [ ] 文档包含 "e2e" 术语的使用约束：仅用于 Web/Mobile surface 的端到端测试上下文
- [ ] 文档 frontmatter 的 `domains` 字段包含 `testing`、`surface`、`test-type` 关键词


### 2-update-core-docs-terminology

- [ ] ARCHITECTURE.md 中不再出现将所有生成测试统称为 "e2e 测试" 或 "高级测试" 的表述
- [ ] guide.md Terminology 部分新增 **Test Type** 条目，包含 Surface → Test Type 映射简要说明
- [ ] task-lifecycle.md 保留类型列表已包含新的 test 类型名（如 test.gen-scripts.{surfaceType}、test.run.{surfaceKey}），与 coding task 的命名变更同步
- [ ] 所有术语变更引用 `docs/reference/test-type-model.md` 作为权威定义来源


### 3-update-gen-scripts-journeys

- [ ] gen-test-scripts/SKILL.md 中不再使用 "e2e" 作为统称，改为引用 `docs/reference/test-type-model.md`
- [ ] gen-test-scripts/types/ 5 个文件各自使用对应 surface 的测试类型名称和语义定义
- [ ] gen-journeys/rules/surface-*.md 5 个文件使用对应的测试类型名称
- [ ] 所有 rules 文件引用概念文档 `docs/reference/test-type-model.md` 中的测试类型定义
- [ ] 生成的测试代码注释/标签中使用 surface-specific 测试类型名称（如 `@cli-functional` 而非 `@e2e`）


### 4-update-run-tests

- [ ] run-tests/SKILL.md 中 "e2e" 仅出现在 promote 功能描述中，且标注为 "Web/Mobile 端到端测试"
- [ ] rules/surfaces/ 5 个文件各自使用对应 surface 的测试类型名称
- [ ] 测试执行输出中的 suite 名称使用 surface-specific 测试类型（如 `cli-functional/journey-name` 而非 `e2e/journey-name`）
- [ ] Web surface 的 Journey filter 标签从 `@e2e` 更新为 `@web-e2e`（或等效的 surface-specific 标签）
- [ ] 所有 rules 文件引用概念文档 `docs/reference/test-type-model.md`


### 5-update-init-justfile

- [ ] 每个 surface 的 justfile 规则文件包含向后兼容 alias（`alias test-e2e := <surface>-test-<type>`）
- [ ] alias 行带有 `# DEPRECATED: removed after v{current+2}` 注释
- [ ] recipe 描述使用 surface-specific 测试类型名称（如 "Run CLI functional tests" 而非 "Run e2e tests"）
- [ ] `just --list` 输出中 recipe 名称和描述清晰区分测试类型
- [ ] 聚合 recipe（`test`）的描述更新为 "Run all surface tests"
- [ ] 所有 rules 文件引用概念文档 `docs/reference/test-type-model.md`


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/surface-test-type-model/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/surface-test-type-model/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
