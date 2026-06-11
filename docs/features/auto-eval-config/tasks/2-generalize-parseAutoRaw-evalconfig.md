---
id: "2"
title: "Generalize parseAutoRaw and add EvalConfig struct"
priority: "P0"
estimated_time: "1h"
dependencies: ["1"]
surface-key: ""
surface-type: ""
breaking: false
type: "coding.feature"
mainSession: false
---

# 2: Generalize parseAutoRaw and add EvalConfig struct

## Description
泛化 `parseAutoRaw` 为递归扫描 auto 子树，不再硬编码 `modeFields` 列表。新增 `EvalConfig` 嵌套结构体（proposal/prd/uiDesign/techDesign 各为 ModeToggle），更新 `AutoConfigDefaults()` 和 `applyDefaults()`。raw 数据结构使用扁平路径 key（`map[string]map[string]bool`），如 `"eval.proposal"` → `{"quick": true, "full": true}`。

## Reference Files
- `docs/proposals/auto-eval-config/proposal.md#Part-2-Auto-Eval-Configuration` — EvalConfig 结构体定义、4 个 ModeToggle 字段、默认值
- `docs/proposals/auto-eval-config/proposal.md#Scope` — parseAutoRaw 泛化规格：flat-path raw key 格式、applyDefaults 调用方式
- `docs/proposals/auto-eval-config/proposal.md#Key-Risks` — parseAutoRaw raw tracking 精度变化风险（M/M）

## Acceptance Criteria
- [ ] `parseAutoRaw` 对 `auto.eval.*` 字段生成正确的 flat-path raw map（如 `map["eval.proposal"]["quick"]=true`）
- [ ] `applyDefaults` 仅补充 YAML 中未显式出现的子键，不覆盖用户显式设置的值
- [ ] `parseAutoRaw` 对现有 auto 字段（test、consolidateSpecs、gitPush）在泛化后仍产生正确的 raw 数据
- [ ] EvalConfig 包含 4 个 ModeToggle 字段：Proposal、Prd、UiDesign、TechDesign
- [ ] AutoConfigDefaults 正确设置默认值：proposal {true,true}, prd {false,false}, uiDesign {true,true}, techDesign {false,false}
- [ ] `AutoConfig.IsZero()` 更新以包含 Eval 字段

## Hard Rules
- raw map 外层 key 使用扁平路径（`"eval.proposal"` 而非嵌套 `"eval"` → `"proposal"`）
- applyDefaults 调用改为 `applyModeDefault(&a.Eval.Proposal, a.raw, "eval.proposal", d.Eval.Proposal)` 的 flat-path 格式

## Implementation Notes
- AutoConfig 新增 `Eval EvalConfig` 字段，yaml tag 为 `"eval"`
- parseAutoRaw 递归遍历 auto 下的 YAML mapping node，遇到 ModeToggle 类型的叶子节点时记录 flat-path key
- applyDefaults 中对 EvalConfig 的每个子字段调用 applyModeDefault，使用 `"eval.proposal"` 等 flat-path key
