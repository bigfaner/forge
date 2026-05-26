---
id: "3"
title: "Update InferType for prefix matching"
priority: "P0"
estimated_time: "1h"
dependencies: ["1"]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 3: Update InferType for prefix matching

## Description
将 `infer.go` 中 `InferType` 对 `T-test-run` 的精确匹配改为 `T-test-run-` 前缀匹配，提取后缀作为 surface-key 在 surfaces map 中查找对应 type。未命中任何已知 key 时回退到原有精确匹配逻辑。更新 `infer_test.go` 测试用例覆盖三种场景。

## Reference Files
- `proposal.md#Proposed-Solution` — surface-key 命名策略，run-tests 使用 surface-key 后缀
- `proposal.md#Feasibility-Assessment` — InferType 变更详情：前缀匹配、回退逻辑、测试用例要求
- `proposal.md#Key-Risks` — 前缀匹配引入歧义的风险及缓解（surface-key 查找 + 精确匹配回退）

## Acceptance Criteria
- [ ] `InferType("T-test-run-backend")` 返回正确的 surface type（`api`），通过前缀匹配而非精确匹配
- [ ] 测试覆盖：已知 surface-key → 正确 type；未知 key → 回退精确匹配；单 surface 退化 → 无后缀场景
- [ ] 单 surface 项目（`surfaces: api`）`InferType("T-test-run")` 行为不变

## Implementation Notes
- 前缀匹配 `T-test-run-` 后的片段作为 surface-key 查找 surfaces map；若未命中任何已知 key，回退到原有精确匹配逻辑
- 新增 InferType 单元测试覆盖三种场景
