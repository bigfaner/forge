---
id: "1"
title: "Embed 模板添加 ## Feature Paths 区域"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Embed 模板添加 ## Feature Paths 区域

## Description
6 个测试流水线 embed 模板（`forge-cli/pkg/task/templates/`）缺少 feature 级路径上下文。Subagent 读 task file 时无法获知 journeys 和 contracts 的目录位置。统一添加 `## Feature Paths` 区域，包含 journeys 和 contracts 两个 discovery `ls` 命令，让 agent 可预探索目录结构。

## Reference Files
- `forge-cli/pkg/task/templates/test-run.md`: 薄模板，几乎无路径信息，需完整添加 `## Feature Paths` (source: proposal.md#改动-1)
- `forge-cli/pkg/task/templates/test-gen-scripts.md`: 薄模板，同上 (source: proposal.md#改动-1)
- `forge-cli/pkg/task/templates/eval-journey.md`: 薄模板，同上 (source: proposal.md#改动-1)
- `forge-cli/pkg/task/templates/eval-contract.md`: 薄模板，同上 (source: proposal.md#改动-1)
- `forge-cli/pkg/task/templates/test-gen-journeys.md`: 富模板，已有较完整路径上下文，仅当路径引用不足时补充 (source: proposal.md#薄模板与富模板的差异化策略)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/task/templates/test-run.md` | 添加 `## Feature Paths` 区域 |
| `forge-cli/pkg/task/templates/test-gen-scripts.md` | 添加 `## Feature Paths` 区域 |
| `forge-cli/pkg/task/templates/eval-journey.md` | 添加 `## Feature Paths` 区域 |
| `forge-cli/pkg/task/templates/eval-contract.md` | 添加 `## Feature Paths` 区域 |
| `forge-cli/pkg/task/templates/test-gen-journeys.md` | 检查现有路径引用，不足时补充 |
| `forge-cli/pkg/task/templates/test-gen-contracts.md` | 检查现有路径引用，不足时补充 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] 6 个 embed 模板均包含 `## Feature Paths` 区域，含 journeys (`ls docs/features/{{.FeatureSlug}}/testing/`) 和 contracts (`ls docs/features/{{.FeatureSlug}}/testing/<journey>/contracts/`) 两个 discovery 命令
- [ ] 富模板（test-gen-journeys、test-gen-contracts）若已有等价路径引用则不重复添加
- [ ] `go build ./...` 和 `go test ./...` 通过

## Hard Rules
- 仅修改 `forge-cli/pkg/task/templates/` 下的 .md 文件，不修改任何 Go 代码

## Implementation Notes
- `{{.FeatureSlug}}` 由 CLI 在 `forge task index` 时从目录路径填充，生成时即确定
- Discovery 命令是信息参考，供 agent 了解目录布局，不要求 agent 在 Step 1 执行
- 对富模板（test-gen-journeys、test-gen-contracts），逐一检查现有内容，仅在路径上下文不足时添加
