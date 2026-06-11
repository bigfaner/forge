---
id: "4"
title: "config_test.go 更新适配 bool 格式"
priority: "P1"
estimated_time: "1h"
complexity: "medium"
dependencies: [1]
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 4: config_test.go 更新适配 bool 格式

## Description
更新 `forge-cli/pkg/forgeconfig/config_test.go` 中的 eval 相关测试，适配 EvalConfig 从 ModeToggle 到 bool 的变更。包括默认值验证、get/set roundtrip、旧格式兼容读取等测试用例。

## Reference Files
- `forge-cli/pkg/forgeconfig/config_test.go`: TestEvalConfigDefaults (line 1976)、TestGetConfigValue_EvalFullRoundtrip (line 2179)、TestDetectPipelineMode (line 1783) 等需适配 (source: proposal.md#Scope > In Scope)
- `forge-cli/pkg/forgeconfig/config_test.go`: 新增旧格式 ModeToggle → bool 兼容读取测试 (source: proposal.md#Key-Risks)

## Acceptance Criteria
- [ ] `TestEvalConfigDefaults` 更新验证 bool 默认值（proposal:true, prd:false, uiDesign:true, techDesign:false）
- [ ] 新增旧格式兼容读取测试：ModeToggle YAML（如 `{quick: true, full: true}`）→ 解析为 bool `true`
- [ ] `TestGetConfigValue_EvalFullRoundtrip` 适配扁平 bool 格式（4 个 key 而非 8 个）
- [ ] 移除或更新 `TestDetectPipelineMode` 中 eval 相关的 mode 判断逻辑

## Implementation Notes
主要受影响测试：
- `TestGetConfigValue_EvalConfig` (line 1292)：get 路径验证
- `TestSetConfigValue_EvalConfig` (line 1358)：set 路径验证
- `TestParseAutoRaw_EvalConfig` (line 1433)：raw YAML tracking
- `TestEvalConfigDefaults` (line 1976)：默认值断言
- `TestAutoConfigIsZero_IncludesEval` (line 1999)：IsZero 计数
- `TestSetConfigValue_EvalPersistence` (line 2015)：write-read roundtrip
- `TestGetConfigValue_EvalFullRoundtrip` (line 2179)：全量 key roundtrip（从 8 个减为 4 个）
