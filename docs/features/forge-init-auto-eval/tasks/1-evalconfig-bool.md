---
id: "1"
title: "EvalConfig 扁平化为 bool 并简化默认值"
priority: "P0"
estimated_time: "1.5h"
complexity: "medium"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: true
type: "coding.refactor"
mainSession: false
---

# 1: EvalConfig 扁平化为 bool 并简化默认值

## Description
将 `EvalConfig` 的 4 个字段从 `ModeToggle` 改为 `bool`，简化 `AutoConfigDefaults()` 中的 eval 默认值，并确保旧格式 config.yaml（ModeToggle map）兼容读取。

proposal 的核心缺陷：brainstorm 执行时 mode="none"，`auto.eval.proposal.$MODE` key 不存在，导致 auto.eval.proposal 配置从未生效。改为 bool 后消除 mode 维度，直接 get/set `auto.eval.proposal`。

## Reference Files
- `forge-cli/pkg/forgeconfig/config.go`: EvalConfig struct (line 35-40) ModeToggle → bool；AutoConfigDefaults() (line 62-78) 简化 eval 默认值；ReadConfig 处理旧格式兼容 (source: proposal.md#Part-1)
- `forge-cli/pkg/forgeconfig/config.go`: ModeToggle → bool 迁移时，ReadConfig 对 bool 字段遇到 map 值取 `full` 子键值 (source: proposal.md#Key-Risks)

## Acceptance Criteria
- [ ] EvalConfig 的 4 个字段（Proposal, Prd, UiDesign, TechDesign）类型为 `bool`，非 `ModeToggle`
- [ ] `AutoConfigDefaults()` 中 eval 默认值简化：proposal:true, prd:false, uiDesign:true, techDesign:false
- [ ] `forge config get auto.eval.proposal` 返回 `"true"` 或 `"false"`（非 ModeToggle 格式）
- [ ] `forge config set auto.eval.prd true` 正确写入并持久化到 config.yaml
- [ ] 已有 config.yaml 中 ModeToggle 格式（如 `proposal: {quick: true, full: true}`）兼容读取：取 `full` 子键值

## Implementation Notes
兼容读取策略：YAML 反序列化时，若 EvalConfig 的 bool 字段收到 map 值，取 `full` 子键作为 bool 值。这确保已有用户的 config.yaml 无需手动修改。

### Test Impact
- Affected test suite(s): `forge-cli/pkg/forgeconfig/`
- Expected fixture changes: eval 相关测试用例需适配 bool 格式
- Risk level: medium
