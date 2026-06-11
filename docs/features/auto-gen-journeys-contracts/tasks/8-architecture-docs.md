---
id: "8"
title: "更新 ARCHITECTURE.md 和相关文档"
priority: "P2"
estimated_time: "0.5h"
dependencies: ["3", "4", "5"]
type: "doc"
mainSession: false
---

# 8: 更新 ARCHITECTURE.md 和相关文档

## Description

更新项目文档以反映新的自动生成测试流水线架构。

## Reference Files
- `docs/proposals/auto-gen-journeys-contracts/proposal.md` — Source proposal
- `docs/ARCHITECTURE.md` — 系统架构文档

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `docs/ARCHITECTURE.md` | 更新测试流水线描述、任务类型列表、依赖拓扑图 |
| `forge-cli/pkg/task/data/test-gen-and-run.md` | 文件头部添加 deprecated 注释 |

## Acceptance Criteria

- [ ] ARCHITECTURE.md 的测试流水线章节包含完整的 Breakdown 模式链路：gen-journeys → eval-journey → gen-contracts → eval-contract → gen-scripts → run → verify
- [ ] ARCHITECTURE.md 的测试流水线章节包含 Quick 模式链路：gen-journeys → gen-contracts → gen-scripts → run → verify（无 eval 中间关卡）
- [ ] `data/test-gen-and-run.md` 文件头部添加 `<!-- DEPRECATED: Replaced by test.gen-journeys + test.gen-contracts in v3.0.0 -->` 注释
- [ ] 文档准确描述 staged across types 拓扑（各 interface type 的 gen-journeys 并行 → 汇聚到 gen-contracts）

## Hard Rules

- 不删除 test-gen-and-run.md 文件（embed.FS 编译需要它存在）
- 仅添加 deprecated 注释，不修改模板内容

## Implementation Notes

- 检查 ARCHITECTURE.md 是否有现有的测试流水线章节；如有则更新，如无则新增
- 确保 Quick 模式和 Breakdown 模式的描述与 proposal.md 中的 Proposed Solution 一致
