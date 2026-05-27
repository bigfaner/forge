---
id: "4"
title: "Unit tests for generic routing and eval config"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["1", "2", "3"]
surface-key: ""
surface-type: ""
breaking: false
type: "coding.feature"
mainSession: false
---

# 4: Unit tests for generic routing and eval config

## Description
为泛化路由、EvalConfig、parseAutoRaw 泛化、mode 检测编写全面的单元测试。包含回归测试确保现有 config get/set 路径行为不变。新增测试用例覆盖 inline tag、自定义 YAML 类型 fallback、flat-path raw tracking。

## Reference Files
- `docs/proposals/auto-eval-config/proposal.md#Key-Risks` — 各风险的 mitigation 中指定的测试用例名：`TestGetByPath_InlineMap`、`TestParseAutoRaw_EvalConfig`、`TestParseAutoRaw_ExistingFields_Regression`
- `docs/proposals/auto-eval-config/proposal.md#Success-Criteria` — PR-1 和 PR-2 的所有可自动化验证的 SC

## Acceptance Criteria
- [ ] `TestGetByPath_InlineMap`: `coverage.coding.feature` 通过 inline tag 正确解析
- [ ] `TestParseAutoRaw_EvalConfig`: raw map 包含 `"eval.proposal"` flat-path key，applyDefaults 正确补充/不覆盖
- [ ] `TestParseAutoRaw_ExistingFields_Regression`: 现有 auto 字段（test、consolidateSpecs、gitPush）raw 数据不变
- [ ] `TestGetStructValueByPath_*`: 多层路径 get（三层、四层）、中间节点 get、nil 指针 get、不存在字段 get
- [ ] `TestSetStructValueByPath_*`: 多层路径 set、ModeToggle 直接 set 拒绝、非 leaf set 拒绝、nil 指针自动初始化
- [ ] `TestGetConfigValue_SurfacesMap_Fallback`: SurfacesMap 字段回退到硬编码路径
- [ ] `TestDetectPipelineMode`: quick/full/none 三种场景
- [ ] 现有 config_test.go + config_schema_test.go 全部通过

## Hard Rules
- 遵循项目 TDD 规范：table-driven tests，覆盖 >80%
- 测试文件：`forge-cli/pkg/forgeconfig/config_test.go`（路由测试）、`forge-cli/pkg/forgeconfig/config_schema_test.go`（schema 测试）

## Implementation Notes
- parseAutoRaw 回归测试需要对比泛化前后对相同 YAML 输入产生的 raw map
- SurfacesMap fallback 测试需要构造包含 scalar 和 map 形式 surfaces 的 config YAML
- mode detection 测试需要 mock 不同目录结构（feature 内/外、有/无 proposal.md）
