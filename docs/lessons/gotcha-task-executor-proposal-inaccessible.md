---
created: "2026-05-29"
tags: [architecture, testing, interface]
---

# Task executor 读不到 proposal 导致实现严重偏离设计

## Problem

Task 1（定义 Pipeline Topology Registry）的实现与提案产生 4 个关键行为偏差 + 5 个类型签名偏差，下游 task 2-6 无法在错误基础上构建。

## Root Cause

因果链（5 层）：

1. **L1 症状**: Task 1 实现与提案存在 9 处偏差（IntentGate 分布错误、缺少 Mode 字段、缺少 per-surface-type 展开、CondHasTestableTasks 逻辑错误、5 个类型签名偏差）

2. **L2 直接原因**: Task executor 无法读取 proposal.md。执行记录明确写道 "No proposal.md found — task AC and implementation notes used as authoritative spec"

3. **L3 中间原因**: Task 的 Reference Files 使用了不存在的路径 `proposal.md`（应为 `docs/proposals/pipeline-topology-registry/proposal.md`）和伪造的章节名（如 `#Intent-Gate-Functions` 不存在于 proposal 中），导致 executor 找不到 proposal 文件（详见 `gotcha-task-reference-source-anchor-misread.md`）

4. **L4 设计缺陷**: Task executor 被迫从两个不完整的来源推导设计：
   - **AC（验收标准）是验证清单，不是规格说明**：AC 检查特定点（"T-review-doc has IntentGate: GateAllowAll"），但不定义完整行为（"其余 11 个节点的 IntentGate 默认为 GateBlockSkipTest"）
   - **Implementation Notes 是提示，不是定义**：提到 "per-surface-key expansion"，但没提 "per-surface-type expansion"
   - Executor 在信息不足时用自身判断填补空白，产生了看似合理但偏离设计的实现

5. **L5 根因**: Proposal 中包含提案的**完整规格**——精确的 struct 定义、函数签名、行为表、12 个节点每个字段的精确值。AC 和 Implementation Notes 是从 proposal 中**摘要**出来的验证点，天然不完整。当 executor 读不到 proposal 时，等于在没有完整规格的情况下做实现设计——这不是 executor 的问题，是信息传递链的断裂。

## Solution

短期：Task 文件的 Reference Files 必须包含 proposal 的完整有效路径，确保 executor 可以读取完整的规格定义。

长期（quick-tasks skill 改进）：
1. 自动在每个 task 的 Reference Files 中加入 proposal 完整路径作为第一个条目
2. source traceability 注解使用完整路径 + 真实章节名
3. 或：在 task 文件中内联关键设计细节（struct 定义、行为表），减少对 proposal 的依赖

## Reusable Pattern

**AC ≠ Spec**。当 task executor 需要实现复杂的数据结构或行为逻辑时：

1. **Proposal 是权威规格**：包含精确的 struct 定义、函数签名、行为表。Executor 必须能读到它
2. **AC 是验证清单**：只检查关键点，不定义完整行为。不足以独立指导实现
3. **信息传递链不可断裂**：如果 Reference Files 中 proposal 路径无效，executor 等于"盲飞"
4. **偏差会级联**：Task 1 的偏差会传播到所有依赖 task（2-6），修复成本随依赖深度指数增长

**判断标准**：当 task 涉及定义新的数据结构、接口或核心算法时，executor 必须能访问完整的设计规格（proposal/tech-design）。如果只能读 AC + Implementation Notes，风险等级应标记为 "critical" 而非 "high"。

## Example

```
# 提案定义了 PipelineNode 的完整字段
type PipelineNode struct {
    Mode string      // "quick", "breakdown", "" — AC 未提及此字段
    Expansion string  // "per-surface-key", "per-surface-type", "" — AC 只提到 per-surface-key
    IntentGate IntentGateFunc // 每个节点都有精确值 — AC 只检查 T-review-doc
    ...
}

# AC 只检查了 3 个点：
- [ ] T-review-doc has IntentGate: GateAllowAll
- [ ] PipelineRegistry slice contains all 12 nodes
- [ ] per-surface-key expansion creates serial chains

# Executor 读不到提案 → 不知道 Mode 字段存在、不知道 per-surface-type、
# 不知道其他节点应该是 GateBlockSkipTest → 用自己的判断填补 → 偏差
```

## Related Files

- `docs/lessons/gotcha-task-reference-source-anchor-misread.md` — 上游 lesson：Reference Files 路径错误导致 proposal 不可达
- `docs/features/pipeline-topology-registry/tasks/1-define-pipeline-registry.md` — 问题 task 文件
- `docs/features/pipeline-topology-registry/tasks/records/1-define-pipeline-registry.md` — 执行记录，证实 "No proposal.md found"
- `forge-cli/pkg/task/pipeline.go` — 含偏差的实现文件
- `docs/proposals/pipeline-topology-registry/proposal.md` — 完整规格（executor 读不到）
