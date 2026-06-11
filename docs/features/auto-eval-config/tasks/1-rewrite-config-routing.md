---
id: "1"
title: "Rewrite config key resolution with reflection"
priority: "P0"
estimated_time: "2.5h"
dependencies: []
surface-key: ""
surface-type: ""
breaking: true
type: "coding.feature"
mainSession: false
---

# 1: Rewrite config key resolution with reflection

## Description
重写 `GetConfigValue`/`SetConfigValue` 为基于任意深度 key 的反射遍历器。实现 `getStructValueByPath`（反射 get）和 `setStructValueByPath`（反射 set + struct marshal 写回），替换当前所有硬编码分发器。删除 `autoModeField`、`getAutoKeyValue`/`getWorktreeKeyValue`/`getCoverageKeyValue`、`setAutoConfigValue`/`setWorktreeConfigValue`/`setCoverageConfigValue`（保留函数体作为 SurfacesMap 等自定义类型的 fallback）。

## Reference Files
- `docs/proposals/auto-eval-config/proposal.md#Part-1-Generic-Config-Key-Resolution` — Get/Set 路径的完整规格：字段匹配规则（YAML tag 优先）、指针处理、叶子节点格式化、错误行为、inline tag 处理
- `docs/proposals/auto-eval-config/proposal.md#Feasibility-Assessment` — getStructValueByPath/setStructValueByPath 的函数签名和行为规格
- `docs/proposals/auto-eval-config/proposal.md#Key-Risks` — inline tag 风险（M/H）、SurfacesMap 自定义类型 fallback（L/M）、注释格式保真度（L/L）

## Acceptance Criteria
- [ ] `forge config get auto.eval.proposal` 返回 `quick:true full:true`（三层深度，需先添加 EvalConfig 结构体定义）
- [ ] `forge config get auto.eval.proposal.quick` 返回 `true`（四层深度）
- [ ] `forge config get auto.eval` 返回 eval 子字段汇总（`proposal: quick:true full:true` 等格式）
- [ ] `forge config get auto` 返回混合类型汇总（ModeToggle、bool、嵌套 struct）
- [ ] `forge config set auto.eval` 被拒绝（"cannot set non-leaf key"）
- [ ] `forge config set auto.eval.proposal true` 被拒绝（"cannot set ModeToggle directly"）
- [ ] `forge config set auto.eval.prd.full true` 正确写入嵌套 config
- [ ] `forge config get auto.eval.proposal.quick.extra` 返回 errKeyNotFound
- [ ] `forge config get auto.nonexistent` 返回 errKeyNotFound
- [ ] `forge config get coverage.coding.feature` 保持现有行为（inline tag 兼容）
- [ ] `forge config get worktree.source-branch` 保持现有行为（回归）

## Hard Rules
- 字段匹配规则：YAML tag 优先于 Go field name；遇到 `yaml:",inline"` 的 map 字段时跳过字段名层级
- SurfacesMap 实现了 `yaml.Unmarshaler`，遇到自定义 YAML 类型时回退到硬编码路径
- Set 路径通过 `readOrCreateConfig` + 反射 set + `writeConfig`（yaml.Marshal），不操作 YAML Node
- 指针 nil 时 get 返回 errKeyNotFound，set 自动初始化（`reflect.New` + `Set`）

## Implementation Notes
- `autoModeField` switch/case 删除，用反射替代
- `getAutoKeyValue`/`getWorktreeKeyValue`/`getCoverageKeyValue` 合并为 `getByPath`
- `setAutoConfigValue`/`setWorktreeConfigValue`/`setCoverageConfigValue` 合并为 `setByPath`
- 保留旧函数体作为 SurfacesMap fallback（新路由返回 errUnsupportedType 时回退）
- CoverageConfig.ByType 使用 `yaml:",inline"` tag，map key 直接出现在 `coverage:` 下
- 非叶子节点 get 的输出格式：ModeToggle → `<name>: quick:<bool> full:<bool>`，bool → `<name>: <bool>`，嵌套 struct → 递归缩进 2 空格
