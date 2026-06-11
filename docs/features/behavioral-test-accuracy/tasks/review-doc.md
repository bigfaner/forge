---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["1", "2", "3", "4"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the behavioral-test-accuracy feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 1-golden-path-journey
- [ ] gen-journeys SKILL.md 引用 `rules/golden-path.md` 规则文件，声明每个 feature 必须至少生成一个 Golden Path Journey
- [ ] `rules/golden-path.md` 包含 Golden Path 双约束规则：(a) 跨越 3+ 步骤操作，(b) 步骤序列必须从 PRD/Design 的 primary user story 中提取核心领域动作序列
- [ ] `rules/golden-path.md` 包含 Feature 复杂度分类启发式规则表格（简单 vs 复杂 Feature 的判定判据和差异化期望）
- [ ] `rules/golden-path.md` 声明语义完整性代理指标：Golden Path 步骤描述必须引用领域术语（如"创建里程碑"）而非 API 术语（如"POST /milestones"）
- [ ] Journey 模板支持 Golden Path 标记（frontmatter 或结构化标记），使下游 gen-contracts 可识别


### 2-fixture-specification
- [ ] `rules/fixture-spec.md` 定义 Fixture Specification Schema（entities 含 entity_type/min_count/relationship_type/parent_entity/field_constraints，可选 state_requirements）
- [ ] `rules/fixture-spec.md` 包含最小合法示例（单实体 CRUD）和完整示例（父子实体关系）
- [ ] Contract 模板 `templates/contract.md` 的 Preconditions 部分包含 `fixture_spec` 字段及其完整 schema 结构
- [ ] `rules/fixture-spec.md` 声明 fixture_spec 为必需字段，entities 至少包含 1 个实体声明


### 3-assertion-depth-seed-data
- [ ] `rules/assertion-depth.md` 包含完整断言分类判据表（行为性 vs 结构性，含边界案例说明）
- [ ] `rules/assertion-depth.md` 声明 ≥80% 断言为行为性断言的强制规则，并要求行为性断言中至少 30% 为深度断言
- [ ] `rules/fixture-from-spec.md` 规则声明从 Contract 的 fixture_spec.entities 读取并生成满足 min_count 的 fixture 数据
- [ ] `rules/fixture-from-spec.md` 规则声明当 fixture_spec 声明需要 N 个子实体时，必须创建 ≥N 个子实体（含 relationship 处理）
- [ ] `types/_shared.md` 新增 backward compatibility 处理：当 fixture_spec 不存在时回退到隐式推断模式并输出 warning


### 4-eval-rubrics-update
- [ ] Journey eval rubric 新增 "Workflow Coverage" 维度（150 分），评分标准包含 Golden Path 存在性子项和多步覆盖度子项
- [ ] Workflow Coverage 维度最低通过阈值 ≥90/150（60%），Golden Path 存在性子项不得为 0 分（一票否决）；eval prompt 要求评审者验证步骤序列是否对应 PRD/Design 中的具体用户故事
- [ ] Contract eval rubric 新增 "Fixture Specification" 维度（100 分），评分标准包含前置数据声明完整性和实体关系覆盖度
- [ ] Fixture Specification 维度最低通过阈值 ≥60/100（60%），entities 必须包含 Contract 涉及的所有实体类型（完整性子项一票否决）；eval prompt 要求评审者验证 entity_type 是否与 Design 中的领域模型一致


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/behavioral-test-accuracy/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/behavioral-test-accuracy/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.

## Acceptance Criteria

- [ ] All acceptance criteria met
