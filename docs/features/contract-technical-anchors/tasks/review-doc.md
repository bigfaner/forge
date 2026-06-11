---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["2", "3", "4", "5", "1"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the contract-technical-anchors feature (quick mode).

## Acceptance Criteria
- [ ] 所有 doc task 的 AC 在最终交付文档中得到满足
- [ ] 交付文档内部无矛盾或重复
- [ ] 变更仅涉及 allowlist 目录下的 .md 文件

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-add-contract-anchor-fields
- [ ] Contract 模板 frontmatter 包含 API 锚点字段：`endpoint` (string)、`method` (string)，及可选字段 `content_type`、`auth_required`
- [ ] Contract 模板 frontmatter 包含 CLI/TUI 锚点字段：`command` (string)，及可选字段 `subcommand`、`flags`、`aliases`
- [ ] Contract 模板 frontmatter 包含 Web 锚点字段：`page` (string)，及可选字段 `route`、`requires_auth`、`layout`
- [ ] Contract 模板 frontmatter 包含 Mobile 锚点字段：`screen` (string)，及可选字段 `navigation_path`、`deeplink`、`platform`
- [ ] `last_anchor_sync` 时间戳字段包含在 frontmatter 中


### 2-add-handbook-generation
- [ ] tech-design skill 根据项目 surface 配置（`forge surfaces`）自动选择并生成对应 handbook
- [ ] cli-handbook 模板覆盖 CLI/TUI surface 的命令、子命令、参数、别名等信息
- [ ] page-map 模板覆盖 Web surface 的页面名称、路由、布局、认证要求
- [ ] screen-map 模板覆盖 Mobile surface 的屏幕名称、导航路径、deeplink、平台
- [ ] 新 handbook 模板格式与 api-handbook 保持一致的 frontmatter 和 section 结构


### 3-gen-contracts-anchor-filling
- [ ] gen-contracts 读取 api-handbook 自动填充 API Contract 的 endpoint/method 字段
- [ ] gen-contracts 读取 cli-handbook 自动填充 CLI/TUI Contract 的 command/subcommand 字段
- [ ] gen-contracts 读取 page-map/screen-map 自动填充 Web/Mobile Contract 的 page/screen 字段
- [ ] Handbook 新鲜度检查：handbook 生成时间早于 tech-design 最后修改时间时，提示用户"handbook 可能过期，建议重新生成"
- [ ] 缺少 handbook 时跳过锚点填充并提示用户"缺少 handbook，建议运行 tech-design 生成"
- [ ] `last_anchor_sync` 时间戳在填充时自动更新


### 4-eval-contract-anchor-checks
- [ ] eval-contract 评分包含锚点字段完整性检查：当 handbook 存在时，Contract 缺少对应锚点字段扣分
- [ ] handbook 内部一致性检查：检测同一 endpoint 的 method 冲突、路径冲突
- [ ] 评分结果报告明确列出缺失的锚点字段（按 surface 类型分组）


### 5-gen-test-scripts-cross-validation
- [ ] 交叉验证比对 Fact Table 与 Contract frontmatter 锚点，结果分类为高置信度/低置信度/无法验证
- [ ] 不匹配时以 handbook 为权威源生成建议修复，展示 diff 供用户确认后写入 Contract
- [ ] 设计文档（handbook）与代码实现不一致时，生成明确的代码 bug 标记报告
- [ ] 输出 surface 覆盖报告，明确列出已验证和未验证的 surface 类型
- [ ] 缺少 handbook 或锚点字段时，降级为 Fact Table 推断（向后兼容），并提示用户
- [ ] 能捕获 lesson 场景（POST vs PUT 不匹配），建议修复为 handbook 定义的 PUT


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/contract-technical-anchors/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/contract-technical-anchors/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
