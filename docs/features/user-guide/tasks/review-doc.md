---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["4", "5", "1", "2", "3"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the user-guide feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-environment-setup
- [ ] 文档覆盖 3 种安装方式：Marketplace 安装、本地构建安装、开发模式安装
- [ ] 包含完整前置条件清单（操作系统、Go 版本、Claude Code CLI 版本及验证命令）
- [ ] 包含安装后验证步骤（`forge --version`、环境检查命令）
- [ ] 包含至少 3 条常见安装问题及解决方案（如 Go 版本不兼容、Claude Code 未安装、权限问题）
- [ ] 所有代码示例可直接复制执行，无需额外修改


### 2-initialization
- [ ] 包含 `forge init` 的完整流程说明（从命令执行到项目就绪）
- [ ] 包含 config.yaml 全字段表格，至少 8 个配置项，每个字段有名称、类型、默认值、说明
- [ ] 包含 Surface 检测机制说明（`forge surfaces detect` 的使用和结果解读）
- [ ] 包含首个项目设置的端到端示例（从 init 到可以开始使用 Forge）
- [ ] 所有代码示例可直接复制执行，无需额外修改


### 3-architecture-overview
- [ ] 包含插件机制说明（Claude Code 插件加载方式和 Forge 的定位）
- [ ] 包含四大组件角色表格（skill、command、agent、hook），每个有名称、用途、触发方式
- [ ] 包含数据流向图解（从用户输入 → Forge 处理 → 文件系统变更的可视化说明）
- [ ] 包含目录约定说明（`.forge/` 目录结构、`docs/features/` 结构、`manifest.md` 作用）
- [ ] 不包含 Go 包结构、CLI 内部命令注册、ResolveScope 等开发者内部实现细节


### 4-usage-guide
- [ ] 包含 Full Mode 至少一个端到端实战示例（从 brainstorm 到任务执行完成）
- [ ] 包含 Quick Mode 至少一个端到端实战示例（从 /quick 到任务执行完成）
- [ ] 包含至少 2 个单命令场景示例（如 /learn、/consolidate-specs）
- [ ] 包含 5 条以上常见问题及排错指引（涵盖安装失败、配置错误、工作流异常、任务阻塞、测试失败）
- [ ] 所有代码示例可直接复制执行，无需额外修改


### 5-readme-index-update
- [ ] README.md 文档索引表包含 `docs/user-guide/environment-setup.md` 链接
- [ ] README.md 文档索引表包含 `docs/user-guide/initialization.md` 链接
- [ ] README.md 文档索引表包含 `docs/user-guide/architecture-overview.md` 链接
- [ ] README.md 文档索引表包含 `docs/user-guide/usage-guide.md` 链接
- [ ] 链接格式与现有文档索引表条目一致


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/user-guide/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/user-guide/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
