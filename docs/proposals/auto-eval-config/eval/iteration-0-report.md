---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Pre-Revision: Freeform Expert Review Findings

## ATTACK_POINTS

1. **[high]** Part 1（泛化 key resolution）与 Part 2（eval 配置）耦合交付，无法单独回滚 | quote: "预计 3-4 小时：泛化路由重写（2h）+ eval 配置 + 4 个 skill（1h）+ 测试（1h）" — Part 1 是 Part 2 前置依赖，但 Part 1 影响所有现有 config get/set 路径的回归风险。如果 Part 1 有 bug，eval 配置无法单独回滚。提案缺少增量交付策略。 | improvement: 增加增量交付策略说明，或将拆分 PR 作为推荐的交付方式

2. **[high]** 反射遍历对 `CoverageConfig.ByType`（`map[string]CoverageStrategy`）的序列化复杂度被低估 | quote: "反射遍历对 map 类型（SurfacesMap、CoverageConfig.ByType）的处理边界" — `coverage.coding.feature` 返回策略值需要理解 `CoverageStrategy` struct 内部结构（Type + Percentage），不是简单的 map key 查找。反射遍历器需要类型特化注册表。 | improvement: 在 Key Risks 中提升 CoverageConfig map 处理的风险评级，或在 Feasibility 中补充 map 序列化的具体方案

3. **[high]** mode 检测依赖的 manifest mode 字段在当前代码中不存在 | quote: "manifest 文件格式需包含 `mode` 字段（已存在于 `feature_complete.go`）" — 实际代码中 `feature_complete.go` 通过检查 proposal.md 是否存在于 feature 目录来判断 quick mode，而非读取 manifest.md 的 mode 字段。整个 eval quick/full 区分依赖尚未实现的 mode 检测 API。 | improvement: 修正 Constraints 中关于 manifest mode 字段的错误声明；在 Scope 中新增 mode 检测机制的实现项，或采用 CLI 级别的 API

4. **[medium]** YAML Node set 路径的序列化保真度：set 路径是走 struct marshal 还是 YAML Node 修改未明确选择 | quote: "Set 路径（YAML Node）：将 `key` 按 `.` 拆分为路径段，沿 YAML Node 树递归走" — 当前所有 set 路径通过 Go struct marshal 写回文件。如果改用 YAML Node 操作，键顺序和注释可能丢失。两种方案有不同的格式保真度特性。 | improvement: 在提案中明确选择 set 路径的实现方式（struct marshal vs YAML Node），并说明选择理由

5. **[medium]** 反射 get 路径的 YAML tag 匹配机制与指针解引用约定未定义 | quote: "将 `key` 按 `.` 拆分为路径段，沿 Go Config struct 树递归走 reflect.Value...struct 按 YAML tag 匹配字段" — `WorktreeConfig` 使用 `yaml:"source-branch"`（含连字符），Config struct 的字段是指针类型（`*AutoConfig`），nil 指针的 get/set 行为未定义。 | improvement: 补充边界行为定义：YAML tag 匹配优先、nil 指针的 get 返回 errKeyNotFound、set 自动初始化

6. **[medium]** `parseAutoRaw` 泛化后 raw map 的 key 格式未定义 | quote: "递归扫描保持相同的叶子节点追踪粒度" — 当前 raw 类型是 `map[string]map[string]bool`，对于嵌套 `auto.eval.proposal`，raw 的 key 格式是扁平路径还是嵌套结构未定义。applyDefaults 的实现依赖这个数据结构。 | improvement: 定义泛化后的 raw 数据结构（建议使用扁平路径 `"eval.proposal"` → `{"quick": true, "full": true}`）

7. **[medium]** `forge config get auto.eval`（中间节点）的行为未定义 | quote: 成功标准列出了 `auto.eval.proposal` 和 `auto.eval.proposal.quick` 的行为，但没有定义 `auto.eval` 本身 — `forge config get auto.eval` 应返回什么？`forge config set auto.eval true` 是否等同于设置所有子字段？ | improvement: 在 Constraints 或成功标准中定义中间节点的查询和设置行为

## BORDERLINE_FINDINGS

- **[medium]** 建议将 Part 1 和 Part 2 拆分为独立可交付的 PR — 这是交付策略建议而非提案内容缺陷，但与 ATTACK_POINTS #1（耦合交付风险）直接相关。考虑在提案中作为推荐的交付策略提及。 | quote: "建议先交付 Part 1 并确保所有现有 config_test.go 和 config_schema_test.go 的回归测试通过"

- **[low]** 建议为 mode 检测提供 CLI 级别的 API（`forge config get mode`） — 这是新功能建议，会扩大 scope，但直接解决 ATTACK_POINTS #3（mode 检测不存在）的问题。 | quote: "建议在提案的 Scope 中新增一个 `forge config get mode` 或 `forge mode` 命令"

- **[low]** 建议补充 `getStructValueByPath` 和 `setYAMLValueByPath` 的完整签名与边界行为 — 实现细节建议，但直接关联 ATTACK_POINTS #5。 | quote: "建议补充：指针 nil 的处理方式、map 类型的 key 匹配策略、YAML tag 与 Go field name 的优先级"

## SKIPPED_FINDINGS

- **[low]** 建议重新考虑 `proposal` 的默认值为 `quick: false, full: true` — 主观偏好，提案已通过 5 Whys 分析选择了 `quick: true, full: true`，有明确的论证链。 | classification: subjective preference | rationale: 默认值选择是设计决策，提案已提供合理的论证

- **[low]** proposal 默认 `quick:true` 可能因低质量评估导致反复迭代 — 与上述 skipped finding 同源的主观判断，假设 "brainstorm 产出的 proposal 通常比较粗糙" 是未经证实的前提。 | classification: subjective preference | rationale: 假设 proposal 质量低但无证据支持

## Rubric Scores

All dimensions: N/A (freeform review, no rubric scoring)
