---
id: "1"
title: "Add execution-order config with surface-key validation"
priority: "P0"
estimated_time: "2h"
dependencies: []
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.feature"
mainSession: false
---

# 1: Add execution-order config with surface-key validation

## Description
在 `config.go` 中新增 `ExecutionOrder []string` 字段，实现 surface-key 归一化（空格/特殊字符→连字符、大写→小写）、合法性校验（`[a-z][a-z0-9-]*`）、同类型冲突检测（多个同 type surface 报错提示显式配置）、execution-order 引用校验（key 必须存在于 surfaces map）、默认优先级约定（api > web > cli > tui > mobile，未覆盖组合按 YAML map key 顺序）。所有校验在 config load time 执行（fail fast）。

## Reference Files
- `proposal.md#Proposed-Solution` — 定义 execution-order 配置语义和 surface-key 命名策略
- `proposal.md#Constraints-&-Dependencies` — surface-key 合法性约束、验证时机、与 justfile 提案的归一化函数对齐要求
- `proposal.md#Key-Risks` — 默认优先级不覆盖所有场景的风险、surface-key 与 justfile 提案定义冲突

## Acceptance Criteria
- [ ] `execution-order` 引用不存在的 surface-key 时在 config load time 报错
- [ ] 同类型冲突场景（2 个 api surface）在 config load time 报错，提示配置 `execution-order`
- [ ] Surface-key 校验：`surfaces: { "ADMIN PANEL": web }` 归一化为 `admin-panel` 通过；`surfaces: { "123bad": web }` 在 config load time 报错
- [ ] 默认优先级：`surfaces: { mobile: mobile, cli: cli, web: web, api: api }` 无 `execution-order` 时，执行顺序为 api → web → cli → mobile

## Hard Rules
- 所有校验在 config load time 执行，不推迟到 build time
- surface-key 归一化函数需与 surface-aware-justfile 提案共用同一实现（`/` → `-`，大写 → 小写）

## Implementation Notes
- 默认优先级未覆盖的组合（如 tui + cli）按 YAML map key 声明顺序排列
- 与 justfile 提案对齐：统一为 config load time 归一化后的值，两提案共用同一归一化函数
