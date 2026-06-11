---
id: "1"
title: "Generate pkg/ dependency graph as factual baseline"
priority: "P0"
estimated_time: "2h"
complexity: "medium"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 1: Generate pkg/ dependency graph as factual baseline

## Description
Phase 1 的第一个产出：用 `go list -json ./pkg/...` 自动生成 `pkg/` 层的完整依赖图作为事实基线。产出须包含每个包的导入关系（哪些 pkg 包导入了哪些其他 pkg 包），分类为三层：leaf（零内部依赖）、基础设施层（仅依赖 types）、领域层（依赖基础设施层或 types）。同时标注领域包之间的横向依赖（如 `pkg/infocmd` 被 4 个领域包导入）。此依赖图将作为后续所有规范文档和包重组的事实依据。

## Reference Files
- forge-cli/pkg/: 所有 17 个子包，需逐一分析导入关系 (source: proposal.md#Proposed-Solution Phase 1)
- forge-cli/go.mod: module path 为 `forge-cli`，无外部可引用路径 (source: proposal.md#Non-Functional-Requirements)

## Acceptance Criteria
- [ ] `go list -json ./pkg/...` 输出已解析，每个 pkg 子包的 ImportPath 和 Imports 已记录
- [ ] 依赖图以 Markdown 格式保存在 `docs/features/forge-cli-codebase-standards/pkg-dependency-graph.md`，包含：(1) 完整导入关系表，(2) 三层分类（leaf/基础设施/领域），(3) 横向依赖标注
- [ ] 每个包标注为 leaf / 基础设施层 / 领域层之一，且标注依据可追溯（列出了该包导入的内部包列表）
- [ ] 发现的横向依赖（领域包之间的互相导入）已逐一列出并标注方向

## Implementation Notes
- 使用 `go list -json ./pkg/... | jq '{ImportPath, Imports}'` 提取导入关系
- 仅关注 `pkg/` 内部的互相导入（过滤掉标准库和第三方库）
- 三层分类规则：leaf（不导入任何 forge-cli pkg）、基础设施（仅导入 `pkg/types/`）、领域（导入其他 pkg 子包）
- 如果发现双向耦合（A 导入 B 且 B 导入 A），标记为"待解耦"并在风险表中记录
