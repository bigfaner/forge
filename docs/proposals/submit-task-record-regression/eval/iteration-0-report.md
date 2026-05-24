---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Pre-Revision Report (Iteration 0)

## Findings Triage Summary

6 findings triaged (4 accepted, 1 partially-accepted, 1 deferred, 7 skipped)

## ATTACK_POINTS

### Accepted (4)

1. **[high]** 标题和范围错配：标题暗示测试 submit-task skill 的类型分发逻辑，但方案只测试 Go CLI 渲染管线
   | quote: "无法确认 submit-task 为每种 task type 选择的模板、填充的字段、渲染的 markdown 是否都正确" (proposal.md line 11)
   | improvement: 重命名标题为反映实际测试目标，或在 scope 中明确声明 submit-task skill 层不在测试范围内

2. **[high]** 单向测试覆盖：方案只验证 Go→markdown 方向，但 Problem Evidence 的核心风险是 template→Go 方向
   | quote: "record-format 模板中的示例 JSON 可能与 Go 端实际接受的 schema 存在偏差" (proposal.md line 17)
   | improvement: 增加 record-format 模板 JSON 示例 → Go RecordData struct 的反向验证，或重新界定 Problem 使其与 Solution 对齐

3. **[high]** `-update` flag 缺乏机械执行机制：Risk 表中的 mitigation 是意图声明而非机制
   | quote: "确认是有意变更则更新 fixture" (proposal.md line 145)
   | improvement: 增加 diff gating 机制描述（如 CI 检测 .diff 文件即失败），或承认此为已知局限

4. **[medium]** `fix` (bare) 和 `coding.fix` 实际走同一模板路径，非独立分发路径
   | quote: "明确区分两种 type 的 fixture，验证各自走正确的模板" (proposal.md line 148)
   | improvement: 在 coverage matrix 中将 fix 标注为 coding.fix 的 alias，减少为 11 种独立 task type

### Partially Accepted (1)

5. **[medium]** Success criterion "新增 fixture 时只需复制文件 + 加一行测试用例" 低估了实际维护成本
   | quote: "新增 fixture 时只需复制文件 + 加一行测试用例" (proposal.md line 155)
   | improvement: 修正为更准确的描述（如 "从已通过校验的历史记录提取 fixture，新增一条 table-driven 测试用例"），删除 "只需" 的简化说法

### Deferred (1)

6. **[medium]** CI <30s 性能预算缺乏分解和基线测量
   | quote: "执行时间 < 30s" (proposal.md line 45)
   | improvement: 分解为可量化规格（~30 fixtures × 2 calls = ~60 function calls，预计 <10s）
   | defer reason: 合理的建议但不涉及内部不一致，留待 Scorer cycle 判断

## SKIPPED_FINDINGS (subjective preference)

- Fixture 选择的时间维度感知：合理的增强建议，但 proposal 并未做出矛盾声明
- 其他衍生建议（suggestion #9, #10, #12, #13）：已由上述 accepted findings 覆盖其核心关切

## Classification Audit

- Factual correction: 1 (finding #4 — fix/coding.fix same dispatch path)
- Structural suggestion: 4 (findings #1, #2, #3, #5 — internal inconsistencies between sections)
- Subjective preference: 2 (findings #2, #6 from review — valid improvements but no contradiction)

## RUBRIC

(all dimensions): N/A
