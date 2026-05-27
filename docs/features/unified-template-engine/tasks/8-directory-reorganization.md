---
id: "8"
title: "目录重组：合并、拆分、重命名"
priority: "P1"
estimated_time: "2h"
complexity: "medium"
dependencies: [7]
surface-key: "."
surface-type: "cli"
breaking: true
type: "coding.refactor"
mainSession: false
---

# 8: 目录重组：合并、拆分、重命名

## Description
执行目录结构重组：`pkg/template/` 整包合并入 `pkg/task/`（2 个模板迁入 `pkg/task/templates/`，代码逻辑迁入 `pkg/task/` 对应文件）；`pkg/task/data/` 拆分为 `pkg/task/templates/`（14 个 autogen+任务创建模板）和 `pkg/task/records/`（6 个 record 模板，去掉 `record-` 前缀）；`pkg/prompt/data/` 重命名为 `pkg/prompt/templates/`。更新所有 `//go:embed` 路径、`templatePath()` 函数和 import 路径。删除空包 `pkg/template/`。

## Reference Files
- `forge-cli/pkg/template/`: 整包合并入 pkg/task/，删除空包 (source: proposal.md#目录重组)
- `forge-cli/pkg/task/data/`: 拆分为 templates/ 和 records/ 子目录 (source: proposal.md#目录重组)
- `forge-cli/pkg/prompt/data/`: 重命名为 pkg/prompt/templates/ (source: proposal.md#目录重组)

## Acceptance Criteria
- [ ] `pkg/template/` 包已删除，其模板文件（coding-fix.md, coding-cleanup.md）迁入 `pkg/task/templates/`，代码逻辑合并到 `pkg/task/` 对应文件
- [ ] `pkg/task/data/` 已拆分为 `pkg/task/templates/`（14 个模板）和 `pkg/task/records/`（6 个模板，无 `record-` 前缀）
- [ ] `pkg/prompt/data/` 已重命名为 `pkg/prompt/templates/`
- [ ] 所有 `//go:embed` 路径、`templatePath()` 函数和 import 路径已更新
- [ ] `go build ./...` 通过，无编译错误

## Implementation Notes
- 目录重组是纯机械性文件移动 + import 更新，不涉及逻辑变更
- 先完成引擎迁移和占位符替换（Task 1-7），再执行目录重组——两步分离降低风险
- 注意 record 文件名去掉 `record-` 前缀（如 `record-coding.md` → `coding.md`）
- `//go:embed` 路径更新需覆盖所有相关 Go 文件

### Test Impact
- Affected test suite(s): `forge-cli/pkg/prompt/`, `forge-cli/pkg/task/`, `forge-cli/pkg/template/`（merged）
- Expected fixture changes: embed path fixtures, import path references
- Risk level: medium
