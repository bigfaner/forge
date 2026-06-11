---
id: "4"
title: "Write constants.md and extend enum-constants.md"
priority: "P1"
estimated_time: "2h"
dependencies: [1]
type: "doc"
mainSession: false
---

# 4: Write constants.md and extend enum-constants.md

## Description
新增 `docs/conventions/constants.md`（魔法值全面管理策略：分类、提取规则、集中管理位置）并扩展 `docs/conventions/enum-constants.md`（增加非枚举常量管理规则：路径常量、超时值、颜色值）。两个文件均包含目标态定义和偏差分析。

## Reference Files
- forge-cli/internal/cmd/quality_gate.go: `"tests/results/raw-output.txt"` 出现 2 次 (source: proposal.md#Evidence)
- forge-cli/internal/cmd/init.go: 颜色值 `#7DCFFF` 硬编码 (source: proposal.md#Evidence)
- forge-cli/internal/cmd/task/list.go: 哨兵数 `99999` (source: proposal.md#Evidence)
- forge-cli/internal/cmd/quality_gate.go: 重试参数 `3` 次、`5*time.Second` 内联 (source: proposal.md#Evidence)

## Affected Files

### Create
| File | Description |
|------|-------------|
| docs/conventions/constants.md | 魔法值管理策略 |

### Modify
| File | Changes |
|------|---------|
| docs/conventions/enum-constants.md | 增加非枚举常量管理规则 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] `docs/conventions/constants.md` 存在，覆盖分类规则（路径、颜色、超时、哨兵数、权限值）和提取规则（何时提取、集中管理位置）
- [ ] `docs/conventions/enum-constants.md` 已扩展，增加路径常量、超时值、颜色值的非枚举常量管理规则
- [ ] 两个文件均包含目标态定义和偏差分析（引用 Evidence 中的具体魔法值案例）
- [ ] 常量集中管理位置明确（如建议每个包内的 `constants.go` 或专门的常量文件）

## Implementation Notes
- 偏差分析优先覆盖常量管理领域，逐文件审计 forge-cli/internal/ 和 forge-cli/pkg/ 中的魔法值
- 参考 `goconst` linter 的检测逻辑作为分类依据
