---
id: "3"
title: "gen-test-scripts 断言深度 + seed data 丰富度规则"
priority: "P0"
estimated_time: "2h"
dependencies: [2]
type: "doc"
mainSession: false
---

# 3: gen-test-scripts 断言深度 + seed data 丰富度规则

## Description

当前 gen-test-scripts 生成的测试断言主要验证 HTTP 状态码和响应 schema，fixture 数据为空容器或最小数据集。本任务新增两个核心规则：(1) 断言深度规则——≥80% 的断言必须为行为性断言（验证业务语义），其中至少 30% 必须是深度断言（验证实体间关系或状态转换）；(2) seed data 丰富度规则——从 Contract 的 fixture_spec 读取并生成满足声明的丰富 fixture 数据。

断言分类判据表区分行为性断言（实体存在、状态正确、关系完整、业务规则满足）和结构性断言（HTTP 状态码、响应 schema、字段类型），80% 阈值允许 20% 结构性断言覆盖 health check 等合理场景。

## Reference Files
- `docs/proposals/behavioral-test-accuracy/proposal.md` — 断言分类判据, Proposed Solution, Key Risks
- `plugins/forge/skills/gen-test-scripts/SKILL.md` — 主 skill 定义，需新增规则引用 (ref: Proposed Solution)
- `plugins/forge/skills/gen-test-scripts/rules/quality-gates.md` — 质量门禁规则，需新增断言深度检查 (ref: 断言分类判据)
- `plugins/forge/skills/gen-test-scripts/types/_shared.md` — 共享测试类型规则，需新增 fixture 消费逻辑 (ref: Proposed Solution)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/gen-test-scripts/rules/assertion-depth.md` | 断言分类判据表 + 80% 行为性阈值 + 30% 深度断言规则 |
| `plugins/forge/skills/gen-test-scripts/rules/fixture-from-spec.md` | 从 Contract fixture_spec 生成丰富 fixture 的规则 |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | 新增 assertion-depth.md 和 fixture-from-spec.md 规则引用 |
| `plugins/forge/skills/gen-test-scripts/rules/quality-gates.md` | 新增断言深度质量门禁 |
| `plugins/forge/skills/gen-test-scripts/types/_shared.md` | 新增 fixture_spec 消费逻辑和 backward compatibility 处理 |

### Delete
| File | Reason |
|------|--------|
| |

## Acceptance Criteria
- [ ] `rules/assertion-depth.md` 包含完整断言分类判据表（行为性 vs 结构性，含边界案例说明）
- [ ] `rules/assertion-depth.md` 声明 ≥80% 断言为行为性断言的强制规则，并要求行为性断言中至少 30% 为深度断言
- [ ] `rules/fixture-from-spec.md` 规则声明从 Contract 的 fixture_spec.entities 读取并生成满足 min_count 的 fixture 数据
- [ ] `rules/fixture-from-spec.md` 规则声明当 fixture_spec 声明需要 N 个子实体时，必须创建 ≥N 个子实体（含 relationship 处理）
- [ ] `types/_shared.md` 新增 backward compatibility 处理：当 fixture_spec 不存在时回退到隐式推断模式并输出 warning

## Hard Rules

- 仅修改以下文件：SKILL.md, rules/assertion-depth.md, rules/fixture-from-spec.md, rules/quality-gates.md, types/_shared.md

## Implementation Notes

- 断言深度规则在生成时强制执行：规则要求 agent 统计行为性断言占比，若 <80% 则自动补充
- 深度断言定义：验证实体间关系（如 `assert milestone.map_id == map.id`）或状态转换（如 `assert response.data.status == "completed"`），区别于浅层断言（如 `assert name == input`）
- 混合断言（如 `assert response.status == 201 AND response.data.name == "milestone-1"`）中若包含至少一个业务字段验证，计为行为性
- fixture-from-spec 规则需处理 relationship_type（belongs_to, has_many）和 field_constraints
