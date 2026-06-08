---
id: "2"
title: "gen-contracts 新增 Fixture Specification 维度"
priority: "P0"
estimated_time: "1.5h"
dependencies: []
type: "doc"
mainSession: false
---

# 2: gen-contracts 新增 Fixture Specification 维度

## Description

当前 Contract 的 Preconditions 仅描述前置条件（如"用户已认证"），不声明前置数据状态（需要哪些实体、实体间关系、最小数据量）。本任务在 gen-contracts 的规则和模板中新增 Fixture Specification 维度，让每个 Contract 明确声明其前置数据需求。

Fixture Specification Schema 包含 entities（entity_type, min_count, relationship_type, parent_entity, field_constraints）和可选的 state_requirements。gen-test-scripts 将在下游消费此 specification 生成丰富 fixture，消除"猜测需要什么数据"的问题。

## Reference Files
- `docs/proposals/behavioral-test-accuracy/proposal.md` — Fixture Specification Schema, Constraints & Dependencies, Backward compatibility 说明
- `plugins/forge/skills/gen-contracts/SKILL.md` — 主 skill 定义，需新增 fixture-spec 规则引用 (ref: Fixture Specification Schema)
- `plugins/forge/skills/gen-contracts/rules/dimension-rules.md` — 现有维度规则，需新增 fixture-spec 维度 (ref: Fixture Specification Schema)
- `plugins/forge/skills/gen-contracts/templates/contract.md` — Contract 模板，需新增 fixture_spec 字段 (ref: Fixture Specification Schema)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/gen-contracts/rules/fixture-spec.md` | Fixture Specification 生成规则 + schema 定义 + 最小合法示例 |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-contracts/SKILL.md` | 新增 fixture-spec.md 规则文件引用 |
| `plugins/forge/skills/gen-contracts/rules/dimension-rules.md` | 新增 fixture_spec 维度到维度列表 |
| `plugins/forge/skills/gen-contracts/templates/contract.md` | Contract 模板 Preconditions 部分新增 fixture_spec 字段 |

### Delete
| File | Reason |
|------|--------|
| |

## Acceptance Criteria
- [ ] `rules/fixture-spec.md` 定义 Fixture Specification Schema（entities 含 entity_type/min_count/relationship_type/parent_entity/field_constraints，可选 state_requirements）
- [ ] `rules/fixture-spec.md` 包含最小合法示例（单实体 CRUD）和完整示例（父子实体关系）
- [ ] Contract 模板 `templates/contract.md` 的 Preconditions 部分包含 `fixture_spec` 字段及其完整 schema 结构
- [ ] `rules/fixture-spec.md` 声明 fixture_spec 为必需字段，entities 至少包含 1 个实体声明

## Hard Rules

- Backward compatibility：`rules/fixture-spec.md` 必须说明当 fixture_spec 不存在时下游 gen-test-scripts 应回退到隐式推断模式并输出 warning，不阻断管线

## Implementation Notes

- Fixture Specification 是声明式的（在 Contract 层声明"需要什么"），不是命令式的（不说"怎么创建"）。生成逻辑由 gen-test-scripts 的 fixture 消费规则负责
- schema 中的 relationship_type 应使用通用关系术语（belongs_to, has_many, has_one），与 ORM 无关
- field_constraints 的 value 字段为 any 类型，允许字符串、数字、布尔值或约束描述
