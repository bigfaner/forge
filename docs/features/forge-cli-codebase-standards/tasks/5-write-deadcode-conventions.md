---
id: "5"
title: "Write dead-code.md and extend code-structure.md"
priority: "P1"
estimated_time: "1.5h"
dependencies: [1]
type: "doc"
mainSession: false
---

# 5: Write dead-code.md and extend code-structure.md

## Description
新增 `docs/conventions/dead-code.md`（死代码识别标准、deprecation 策略、清理流程）并扩展 `docs/conventions/code-structure.md`（增加包组织相关的结构规则）。两个文件均包含目标态定义和偏差分析。

## Reference Files
- forge-cli/pkg/task/frontmatter.go: deprecated `Scope` 字段 (source: proposal.md#Evidence)
- forge-cli/internal/cmd/output.go + internal/cmd/base/output.go: 重复 `Debugf` 定义 (source: proposal.md#Evidence)
- forge-cli/internal/cmd/task/claim.go: test-bridge 别名函数 `checkExistingTaskState`、`getTaskPhase`、`compareVersionIDs` (source: proposal.md#Evidence)

## Affected Files

### Create
| File | Description |
|------|-------------|
| docs/conventions/dead-code.md | 死代码识别标准和清理流程 |

### Modify
| File | Changes |
|------|---------|
| docs/conventions/code-structure.md | 增加包组织相关的结构规则 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] `docs/conventions/dead-code.md` 存在，覆盖识别标准（deprecated 字段、重复定义、构建产物）、deprecation 策略、清理流程
- [ ] `docs/conventions/code-structure.md` 已扩展，增加包组织相关的结构规则（引用 package-organization.md 中的依赖方向）
- [ ] dead-code.md 明确区分三类：纯粹死代码（可直接删除）、test-bridge 别名（需评估后处理）、deprecated 保留字段（需迁移计划）
- [ ] 包含目标态定义和模块级偏差摘要

## Implementation Notes
- 死代码分类须特别标注 `getTaskPhase`（有 5 处生产调用，非纯粹死代码）避免误删
