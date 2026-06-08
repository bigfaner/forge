---
id: "1"
title: "gen-journeys 新增 Golden Path Journey 强制要求"
priority: "P0"
estimated_time: "1.5h"
dependencies: []
type: "doc"
mainSession: false
---

# 1: gen-journeys 新增 Golden Path Journey 强制要求

## Description

gen-journeys skill 当前生成的 Journey 是操作级别的单步描述（如单个 CRUD 操作），不强制要求覆盖用户核心工作流。本任务在 gen-journeys 的 SKILL.md 和规则文件中新增 Golden Path Journey 强制要求，确保每个 feature 至少包含一个跨越多步骤的、覆盖核心业务语义的 Journey。

Golden Path 必须同时满足两个约束：(a) 跨越 3+ 步骤操作，(b) 覆盖 PRD/Design 文档中 primary user story 的核心领域动作序列（语义完整性——不是任意 3 步拼凑，而是用户真实工作流中的关键步骤链）。同时需要根据 Feature 复杂度分类启发式规则区分简单/复杂 feature 的差异化期望。

## Reference Files
- `docs/proposals/behavioral-test-accuracy/proposal.md` — Proposed Solution, 断言分类判据, Feature 复杂度分类启发式规则
- `plugins/forge/skills/gen-journeys/SKILL.md` — 主 skill 定义，需新增 Golden Path 规则引用 (ref: Proposed Solution)
- `plugins/forge/skills/gen-journeys/templates/journey.md` — Journey 模板，可能需新增 Golden Path 标记 (ref: Proposed Solution)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/gen-journeys/rules/golden-path.md` | Golden Path 强制规则 + Feature 复杂度分类启发式规则 |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-journeys/SKILL.md` | 新增 golden-path.md 规则文件引用和 Golden Path 相关硬约束 |
| `plugins/forge/skills/gen-journeys/templates/journey.md` | 新增 Golden Path Journey 模板标记（如 `golden_path: true` frontmatter 字段） |

### Delete
| File | Reason |
|------|--------|
| |

## Acceptance Criteria
- [ ] gen-journeys SKILL.md 引用 `rules/golden-path.md` 规则文件，声明每个 feature 必须至少生成一个 Golden Path Journey
- [ ] `rules/golden-path.md` 包含 Golden Path 双约束规则：(a) 跨越 3+ 步骤操作，(b) 步骤序列必须从 PRD/Design 的 primary user story 中提取核心领域动作序列
- [ ] `rules/golden-path.md` 包含 Feature 复杂度分类启发式规则表格（简单 vs 复杂 Feature 的判定判据和差异化期望）
- [ ] `rules/golden-path.md` 声明语义完整性代理指标：Golden Path 步骤描述必须引用领域术语（如"创建里程碑"）而非 API 术语（如"POST /milestones"）
- [ ] Journey 模板支持 Golden Path 标记（frontmatter 或结构化标记），使下游 gen-contracts 可识别

## Hard Rules

- Golden Path 规则适用于所有 surface type，不区分 surface-specific 文件
- 判定优先级：实体关系 > 工作流描述。存在父子实体关系即判定为复杂 feature

## Implementation Notes

- Golden Path 的语义完整性是关键设计点：规则必须明确禁止"用不相关操作凑数"的行为，要求 agent 从 PRD/Design 文档中提取用户操作序列
- 简单 feature（单实体 CRUD）的 Golden Path 可以是完整的 CRUD 循环（create → read → update → delete），不需要强制 5+ 步骤
- 复杂 feature（≥2 个实体类型 + 父子关系）的 Golden Path 期望 5+ 步骤，覆盖实体间交互
