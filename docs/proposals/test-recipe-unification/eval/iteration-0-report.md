iteration: 0
title: "Pre-Revision (Freeform Findings)"

# ATTACK_POINTS

- **[high]** RunProjectTests 探测链与新 recipe 模型语义冲突，proposal 混淆了两条独立测试调用路径 | quote: "This means there are two distinct code paths that invoke tests, and the proposal conflates them." | improvement: 明确区分 gate sequence 的 RunGate 和 RunProjectTests 两条路径的行为差异

- **[high]** addFixTask 硬编码映射 step=="unit-test" 到 "just test"，重命名后映射将不正确 | quote: "the current `addFixTask` function (line 430-431) has a hardcoded mapping where `step == \"unit-test\"` maps to `testScript = \"just test\"`. This mapping becomes incorrect after the rename." | improvement: 移除旧映射，让 step=="unit-test" 自然映射到 "just unit-test"

- **[medium]** DefaultGateSequence 与 UnitGateSequence 职责划分不清，submit.go 中的调用逻辑存在歧义 | quote: "If both exist, what is `DefaultGateSequence` used for after the migration? Is it renamed to `UnitGateSequence` and a new `DefaultGateSequence` is created for all-completed? The proposal is ambiguous." | improvement: 定义迁移后三个 sequence 函数的精确内容，考虑将 DefaultGateSequence 重命名为 FullGateSequence

- **[medium]** auto.e2eTest 重命名后旧脚本调用会静默失败，无弃用警告 | quote: "`SetConfigValue(\"auto.e2eTest\", ...)` from any existing script would silently fail with \"unknown config key\" -- there is no deprecation warning." | improvement: 在 parseAutoRaw 中检测旧键名并输出迁移提示

- **[medium]** journey_isolation.go 迁移后 just test 缺少标准参数签名约定 | quote: "the proposal also defines `just test` as \"Surface 级高级测试\" without mentioning a journey name parameter. The Key Scenarios section says \"just test <journeyName>\" but the justfile templates in Tier 3 do not show this signature." | improvement: 在 justfile template contract 中定义 test recipe 接受可选的 journey 过滤参数

- **[medium]** 影响范围评估遗漏了 internal/cmd/test/test.go 中的 e2e-test 引用 | quote: "Grep results show `internal/cmd/test/test.go` line 36 contains the string \"Runs just e2e-test from the project root with the journey name as filter.\" This file is not listed in any of the five tiers of the Impact Analysis." | improvement: 将 internal/cmd/test/test.go 补充到 Tier 1 影响范围

- **[medium]** RunProjectTests 的五级 fallback 链与「无 Fallback」声明矛盾 | quote: "`RunProjectTests` in `testrunner.go` has a five-level fallback chain. The \"无 Fallback\" claim is misleading." | improvement: 明确「无 Fallback」仅适用于 gate sequence 的 RunGate，不适用于 RunProjectTests 探测链

# BORDERLINE_FINDINGS

- test recipe 语义过载（迁移前后含义完全反转） | 评审者指出用户习惯会受影响，但这是提案的有意设计决策（Test Pyramid 对齐），属于 UX 主观判断，非内部不一致

# SKIPPED_FINDINGS

- auto.test YAML tag 可能与未来顶层 test 键名冲突 | 理由：假设性问题，无实际证据支持，属于主观偏好

# rubric
(all dimensions): N/A
