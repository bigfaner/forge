---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["5", "2", "3", "4"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the forge-cli-codebase-standards feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 2-write-package-organization
- [ ] `docs/conventions/package-organization.md` 存在，包含目标态定义（非描述性）和偏差分析表
- [ ] 依赖方向规则明确：`cmd → internal → pkg`（严格单向），`pkg/` 内三层模型定义清晰（leaf/基础设施/领域）
- [ ] 偏差分析表引用 Task 1 依赖图的具体数据（如 `pkg/infocmd` 被 4 个领域包导入的事实）
- [ ] PR review checklist 包含：包结构变更需 review 确认符合依赖方向规则和包职责定义
- [ ] 开发者工作流描述完整：新增命令时应在 `internal/cmd/<command-group>/` 下创建文件


### 3-write-naming
- [ ] `docs/conventions/naming.md` 存在，覆盖文件名、函数名、常量名、包名四类命名规则
- [ ] 包含目标态定义（规范性，如"包名使用单个单词、小写、无下划线"）
- [ ] 包含模块级偏差摘要（如 `forgeconfig` 应为 `config` 或保持 `forgeconfig` 的取舍说明）
- [ ] 规则可执行：每条规则可通过 `grep` 或 `go vet` 验证，或明确标注为人工 review 项


### 4-write-constants-conventions
- [ ] `docs/conventions/constants.md` 存在，覆盖分类规则（路径、颜色、超时、哨兵数、权限值）和提取规则（何时提取、集中管理位置）
- [ ] `docs/conventions/enum-constants.md` 已扩展，增加路径常量、超时值、颜色值的非枚举常量管理规则
- [ ] 两个文件均包含目标态定义和偏差分析（引用 Evidence 中的具体魔法值案例）
- [ ] 常量集中管理位置明确（如建议每个包内的 `constants.go` 或专门的常量文件）


### 5-write-deadcode-conventions
- [ ] `docs/conventions/dead-code.md` 存在，覆盖识别标准（deprecated 字段、重复定义、构建产物）、deprecation 策略、清理流程
- [ ] `docs/conventions/code-structure.md` 已扩展，增加包组织相关的结构规则（引用 package-organization.md 中的依赖方向）
- [ ] dead-code.md 明确区分三类：纯粹死代码（可直接删除）、test-bridge 别名（需评估后处理）、deprecated 保留字段（需迁移计划）
- [ ] 包含目标态定义和模块级偏差摘要


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/forge-cli-codebase-standards/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/forge-cli-codebase-standards/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.
