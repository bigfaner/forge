---
id: "5"
title: "gen-test-scripts 增加交叉验证和 surface 覆盖报告"
priority: "P1"
estimated_time: "2h"
dependencies: [1]
type: "doc"
mainSession: false
---

# 5: gen-test-scripts 增加交叉验证和 surface 覆盖报告

## Description
gen-test-scripts Step 1 代码侦察后增加交叉验证：将 Fact Table（代码侦察结果）与 Contract frontmatter 锚点比对，不匹配时分类为高置信度/低置信度/无法验证，以设计文档（handbook）为准生成建议修复（用户确认后写入 Contract）。设计文档与代码不一致时标记为代码 bug。输出 surface 覆盖报告。

## Reference Files
- `docs/proposals/contract-technical-anchors/proposal.md` — Proposed Solution, Key Scenarios, Success Criteria
- `plugins/forge/skills/gen-test-scripts/SKILL.md`: 增加交叉验证步骤 (ref: Proposed Solution)
- `plugins/forge/skills/gen-test-scripts/rules/step-0.5-validation.md`: 验证步骤参考 (ref: Key Scenarios)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | 增加交叉验证逻辑和 surface 覆盖报告 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] 交叉验证比对 Fact Table 与 Contract frontmatter 锚点，结果分类为高置信度/低置信度/无法验证
- [ ] 不匹配时以 handbook 为权威源生成建议修复，展示 diff 供用户确认后写入 Contract
- [ ] 设计文档（handbook）与代码实现不一致时，生成明确的代码 bug 标记报告
- [ ] 输出 surface 覆盖报告，明确列出已验证和未验证的 surface 类型
- [ ] 缺少 handbook 或锚点字段时，降级为 Fact Table 推断（向后兼容），并提示用户
- [ ] 能捕获 lesson 场景（POST vs PUT 不匹配），建议修复为 handbook 定义的 PUT

## Implementation Notes
- 交叉验证以设计文档（handbook）为 authority source，设计-实现不一致时定位为代码 bug
- 用户确认环节作为最终防线，修复前展示 diff 供审阅
- 结果分类中低置信度和无法验证不自动处理，仅提示用户人工确认
- 静态分析无法覆盖所有路由注册模式（动态加载、反射等），侦察结果不完整时标记为低置信度
