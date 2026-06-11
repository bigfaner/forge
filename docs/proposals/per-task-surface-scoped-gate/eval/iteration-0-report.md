---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Pre-Revision Report — Iteration 0

## ATTACK_POINTS

### Factual Corrections (direct edit)

- **[high]** RunGate()共用风险：缺少scope守卫伪代码 | quote: "需确保 feature-level gate 传空 scope 时不触发 prefixed 解析" | improvement: 在 Proposed Solution 中增加 RunGate() prefixed resolution 的伪代码，明确触发条件为 `scope != ""`

- **[medium]** 部署原子性 Impact 评级不准确 | quote: "Likelihood=L, Impact=L" | improvement: 将 Key Risks 表中 "两个改动的部署原子性" 的 Impact 从 L 修正为 M，并在 Scope 或新增 Migration Notes 中文档化中间态行为

### Structural/Architectural Suggestions (edit where internal inconsistency found)

- **[medium]** resolveRecipe() isomorphism 差异未声明 | quote: "复用 quality_gate_lifecycle.go 中 resolveRecipe() 的模式，但改用 surface key（非 type）作为前缀" | improvement: 在 Innovation Highlights 中显式声明 prefixed resolution 在 RunGate() 中是单次选择（pick one by key），lifecycle 层是全量遍历（iterate all by type），两者 dispatcher isomorphism 有本质差异

- **[high]** HasRecipe() Windows 性能分析缺失 | quote: "just.HasRecipe() 可探测 prefixed recipe 是否存在" | improvement: 在 Feasibility Assessment 或 Key Risks 中增加 N×2 HasRecipe() 探测的性能分析（4 steps × 2 probes × ~50ms/fork on Windows ≈ 400ms overhead），并说明 resolveRecipe 模式实际上最坏 2N 但平均 < N+2

- **[medium]** Surface rule stub LLM 填充一致性机制缺失 | quote: "Surface rule 的 stub recipe 未被 LLM 正确填充" | improvement: 在 Key Risks 的 Mitigation 列中增加约束注释要求：stub recipe 模板中应注明 "This recipe compiles ONLY the <key> surface code"

- **[medium]** 失败模式分析不完整 | quote: "执行 just backend-compile → just backend-fmt → just backend-lint → just backend-unit-test" | improvement: 在 Key Scenarios 中增加 scenario 6：prefixed recipe 存在但执行失败，说明 onFail 收到的 step name 包含 surface 上下文

- **[medium]** surface key 命名冲突边界未讨论 | quote: "不需要从 index.json 读取" | improvement: 在 Constraints 中增加 NormalizeSurfaceKey 输出字符集约束，作为 prefixed recipe 命名安全性的论据

## BORDERLINE_FINDINGS

(none)

## SKIPPED_FINDINGS

(none — all findings are factual or structural with verifiable internal inconsistencies)

## Rubric

All dimensions: N/A (pre-revision)
