---
id: "2"
title: "更新核心文档测试术语"
priority: "P0"
estimated_time: "1.5h"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 2: 更新核心文档测试术语

## Description

更新 ARCHITECTURE.md、guide.md（Terminology 部分）和 task-lifecycle.md 中的测试相关术语。将泛化的 "e2e 测试" / "高级测试" 替换为 surface-specific 的测试类型名称，并在 guide.md 中补充 Surface Type → Test Type 映射说明。

## Reference Files
- `docs/proposals/surface-test-type-model/proposal.md#Problem` — 列出了当前术语混用的 6 个证据点
- `docs/proposals/surface-test-type-model/proposal.md#Evidence` — ARCHITECTURE.md 中"高级测试"和"e2e 测试"交替使用的具体问题
- `docs/proposals/surface-test-type-model/proposal.md#Success-Criteria` — SC2（guide.md）和 SC7（task-lifecycle.md 保留类型列表）

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `docs/ARCHITECTURE.md` | 测试相关章节：将 "e2e 测试"/"高级测试" 替换为 surface-specific 测试类型名称 |
| `plugins/forge/hooks/guide.md` | Terminology 部分：补充 Test Type 定义条目，说明 Surface → Test Type 映射关系 |
| `docs/business-rules/task-lifecycle.md` | 保留类型列表（BIZ-task-lifecycle-003）：新增 surface-specific test 类型名（如 test.gen-scripts.cli 等） |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria

- [ ] ARCHITECTURE.md 中不再出现将所有生成测试统称为 "e2e 测试" 或 "高级测试" 的表述
- [ ] guide.md Terminology 部分新增 **Test Type** 条目，包含 Surface → Test Type 映射简要说明
- [ ] task-lifecycle.md 保留类型列表已包含新的 test 类型名（如 test.gen-scripts.{surfaceType}、test.run.{surfaceKey}），与 coding task 的命名变更同步
- [ ] 所有术语变更引用 `docs/reference/test-type-model.md` 作为权威定义来源

## Hard Rules

- guide.md Terminology 条目保持简短（≤3 行定义 + 1 行指向概念文档的链接），详细定义在概念文档中
- task-lifecycle.md 保留类型列表的更新方式：新增通配模式（如 `test.gen-scripts.*`）而非逐一列举 5 种 surface

## Implementation Notes

guide.md 是 forge plugin hooks 入口，每次 agent 启动时加载。新增的 Test Type 条目应紧接在 Surface Type 条目之后，便于 agent 建立关联。
