---
id: "2"
title: "Prompt 模板渲染 FEATURE_SLUG"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 2: Prompt 模板渲染 FEATURE_SLUG

## Description
6 个测试流水线 prompt 模板（`forge-cli/pkg/prompt/templates/`）在 context 中声明了 `FeatureSlug` 但未渲染到输出。Subagent 收到 prompt 时无法直接获取 slug，需要从 TASK_FILE 路径反推。在 `TASK_FILE` 行之后添加 `FEATURE_SLUG: {{.FeatureSlug}}`，使 agent 无需路径解析。

## Reference Files
- `forge-cli/pkg/prompt/templates/test-run.md`: 当前输出 TASK_ID + TASK_FILE + SURFACE_KEY，缺少 FEATURE_SLUG (source: proposal.md#改动-2)
- `forge-cli/pkg/prompt/templates/test-gen-journeys.md`: 同上，TASK_FILE 后需添加 FEATURE_SLUG 行 (source: proposal.md#改动-2)
- `forge-cli/pkg/prompt/templates/test-gen-contracts.md`: 同上 (source: proposal.md#改动-2)
- `forge-cli/pkg/prompt/templates/test-gen-scripts.md`: 同上 (source: proposal.md#改动-2)
- `forge-cli/pkg/prompt/templates/eval-journey.md`: 同上 (source: proposal.md#改动-2)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/prompt/templates/test-run.md` | TASK_FILE 行后添加 `FEATURE_SLUG: {{.FeatureSlug}}` |
| `forge-cli/pkg/prompt/templates/test-gen-journeys.md` | 同上 |
| `forge-cli/pkg/prompt/templates/test-gen-contracts.md` | 同上 |
| `forge-cli/pkg/prompt/templates/test-gen-scripts.md` | 同上 |
| `forge-cli/pkg/prompt/templates/eval-journey.md` | 同上 |
| `forge-cli/pkg/prompt/templates/eval-contract.md` | 同上 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] 6 个 prompt 模板均输出 `FEATURE_SLUG: {{.FeatureSlug}}` 行，位于 `TASK_FILE` 行之后
- [ ] `go build ./...` 和 `go test ./...` 通过

## Hard Rules
- 仅修改 `forge-cli/pkg/prompt/templates/` 下的 .md 文件，不修改任何 Go 代码

## Implementation Notes
- `{{.FeatureSlug}}` 由 `run-tasks` dispatcher 从 `index.json` 的 `feature` 字段传入，来源是 `forge task index` 时写入的值
- `FeatureSlug` 已在 prompt 模板 frontmatter 的 `context` 组声明，`promptTemplateData` Go 结构体已有此字段，不需修改 Go 代码
- 插入位置：`TASK_FILE: {{.TaskFile}}` 行之后、`{{if .SurfaceKey}}SURFACE_KEY` 行之前（若有 SURFACE_KEY 条件块）或直接在 TASK_FILE 之后
