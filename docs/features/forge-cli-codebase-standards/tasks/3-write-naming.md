---
id: "3"
title: "Write naming.md convention"
priority: "P1"
estimated_time: "1.5h"
dependencies: [1]
type: "doc"
mainSession: false
---

# 3: Write naming.md convention

## Description
新增 `docs/conventions/naming.md`，定义文件名、函数名、常量名、包名命名规范。包含目标态定义和模块级偏差分析摘要（非逐文件，覆盖命名领域的主要偏差模式）。

## Reference Files
- forge-cli/internal/cmd/: 顶层文件命名模式（如 `quality_gate.go`、`init_surfaces.go`）与子包命名模式对比 (source: proposal.md#Scope item 4)
- forge-cli/pkg/: 包命名模式分析（如 `infocmd`、`gitx` → `git`、`forgeconfig` 等） (source: proposal.md#Scope item 4)

## Affected Files

### Create
| File | Description |
|------|-------------|
| docs/conventions/naming.md | 命名规范，含目标态和模块级偏差摘要 |

### Modify
| File | Changes |
|------|---------|

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] `docs/conventions/naming.md` 存在，覆盖文件名、函数名、常量名、包名四类命名规则
- [ ] 包含目标态定义（规范性，如"包名使用单个单词、小写、无下划线"）
- [ ] 包含模块级偏差摘要（如 `forgeconfig` 应为 `config` 或保持 `forgeconfig` 的取舍说明）
- [ ] 规则可执行：每条规则可通过 `grep` 或 `go vet` 验证，或明确标注为人工 review 项

## Implementation Notes
- 偏差分析以模块级摘要覆盖，不逐文件扫描
- 参考 Go 社区命名惯例（effective go, go code review comments）
