---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["3", "4", "5", "6", "7", "1", "2"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the surface-first-testing feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-create-surface-templates
- [ ] `templates/surfaces/cli.md` 包含 7 个必要段落：文件位置、隔离模型、断言重点、超时策略、生命周期、Contract/Journey 比例、反模式，以及断言偏好表
- [ ] `templates/surfaces/api.md` 包含同样 7 个段落 + 断言偏好表
- [ ] `templates/surfaces/web.md` 包含同样 7 个段落 + 断言偏好表
- [ ] `templates/surfaces/tui.md` 包含同样 7 个段落 + 断言偏好表
- [ ] `templates/surfaces/mobile.md` 包含同样 7 个段落 + 断言偏好表
- [ ] `references/test-type-model.md` 包含完整分类标准、Surface → Test Type 映射表、e2e 术语约束和语义定义


### 2-rewrite-test-guide
- [ ] SKILL.md 包含读取 `.forge/config.yaml` surfaces 配置的步骤
- [ ] SKILL.md 包含从 `templates/surfaces/*.md` 生成 per-surface convention 文件（index.md + core.md）的步骤
- [ ] SKILL.md 包含生成顶层 `docs/conventions/testing/index.md` 速查表的步骤
- [ ] 框架检测（signal-detection）重构为辅助步骤，结果仅用于填充 core.md 断言偏好表
- [ ] 旧的"框架检测 → 单文件生成"流程已移除


### 3-update-guide-hook
- [ ] guide.md 新增 Testing section，增量 <= 20 行
- [ ] 包含 Surface → Test Type 映射表（cli/api/web/tui/mobile 5 种 Surface）
- [ ] 包含 e2e 术语约束说明（"e2e" 仅用于 Web/Mobile）
- [ ] 包含测试文件位置规则（`tests/{surface}/` 或 `tests/e2e/`）
- [ ] 包含引导运行 `/test-guide` 获取完整策略的提示


### 4-update-gen-test-scripts
- [ ] SKILL.md Convention 加载路径改为 `testing/{surface}/core.md` surface 目录遍历
- [ ] 检测到旧结构 convention 文件时输出迁移提示而非静默失败
- [ ] 生成的测试代码使用 per-surface build tag 命名（如 `cli_functional` 而非 `e2e`）


### 5-update-run-tests
- [ ] SKILL.md Convention 读取路径改为 `testing/{surface}/core.md`
- [ ] 检测到旧结构 convention 文件时输出迁移提示而非静默失败


### 6-update-init-justfile
- [ ] SKILL.md Test recipe 根据新目录结构 `testing/{surface}/` 生成
- [ ] Test recipe 命名遵循 Surface type 约定（如 `test-cli`、`test-api`）


### 7-cleanup-regenerate
- [ ] `docs/conventions/testing/` 下 6 个旧框架文件（ginkgo/go/junit/pytest/rust/vitest）已删除
- [ ] `docs/reference/test-type-model.md` 已删除
- [ ] Forge 项目自身 `docs/conventions/testing/cli/` 已用新 test-guide 重新生成（含 index.md + core.md）
- [ ] 顶层 `docs/conventions/testing/index.md` 已重新生成


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/surface-first-testing/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/surface-first-testing/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.

## Acceptance Criteria
- [ ] All doc task deliverables reviewed against their acceptance criteria
