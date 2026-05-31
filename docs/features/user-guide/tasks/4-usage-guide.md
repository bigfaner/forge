---
id: "4"
title: "编写使用指南 usage-guide.md"
priority: "P1"
estimated_time: "2h"
dependencies: []
type: "doc"
mainSession: false
---

# 4: 编写使用指南 usage-guide.md

## Description
创建 `docs/user-guide/usage-guide.md`，提供 Full Mode 和 Quick Mode 的端到端实战示例、单命令场景和常见问题排错指引。当前 README.md 的"5分钟体验"仅两条命令，缺少完整的端到端实战和常见问题排查。

## Reference Files
- `README.md`: 5分钟体验部分和命令快速参考，作为端到端示例的基础 (source: proposal.md#Problem)
- `docs/ARCHITECTURE.md`: Full Mode 和 Quick Mode 的流程描述 (source: proposal.md#Problem)
- `docs/business-rules/task-lifecycle.md`: 任务生命周期和状态转换，用于常见问题排查 (source: proposal.md#Proposed-Solution)
- `docs/business-rules/quality-gate.md`: Quality Gate 规则，帮助用户理解工作流中的检查点 (source: proposal.md#Proposed-Solution)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `docs/user-guide/usage-guide.md` | 使用指南：Full/Quick Mode 端到端实战、单命令场景、常见问题与排错 |

### Modify
| File | Changes |
|------|---------|
| (无) | |

### Delete
| File | Reason |
|------|--------|
| (无) | |

## Acceptance Criteria
- [ ] 包含 Full Mode 至少一个端到端实战示例（从 brainstorm 到任务执行完成）
- [ ] 包含 Quick Mode 至少一个端到端实战示例（从 /quick 到任务执行完成）
- [ ] 包含至少 2 个单命令场景示例（如 /learn、/consolidate-specs）
- [ ] 包含 5 条以上常见问题及排错指引（涵盖安装失败、配置错误、工作流异常、任务阻塞、测试失败）
- [ ] 所有代码示例可直接复制执行，无需额外修改

## Implementation Notes
- 文档使用中文
- 端到端示例需使用具体的命令序列（如 `/forge:write-prd` → `/forge:tech-design` → `/forge:breakdown-tasks` → `/forge:run-tasks`），从 README 和 ARCHITECTURE.md 中提取实际命令名
- 常见问题需覆盖实际用户可能遇到的场景，从 docs/business-rules/ 中提取排错信息
- 文档顶部标注"最后更新"日期和对应版本（v3.0.0）
