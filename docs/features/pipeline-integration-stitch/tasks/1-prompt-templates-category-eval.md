---
id: "1"
title: "创建 prompt 模板 + CategoryEval + eval record/validation 系统"
priority: "P0"
estimated_time: "3h"
dependencies: []
scope: "backend"
breaking: true
type: "coding.enhancement"
mainSession: false
---

# 1: 创建 prompt 模板 + CategoryEval + eval record/validation 系统

## Description

解决 P0（4 个 prompt 模板缺失）和 P1 核心（eval 类型分类错误 + submit 验证 + record 渲染）。当前 `eval.journey`/`eval.contract` 落入 `CategoryCoding`，导致 eval 任务被要求提供测试证据；`prompt/data/` 缺少 4 个执行阶段模板导致 `Synthesize()` 必定失败。

## Reference Files
- `proposal.md#P0-—-Pipeline-执行必定失败` — 定义 4 个缺失模板文件的列表和影响
- `proposal.md#P1-—-CategoryEval-+-依赖加固` — CategoryEval 分类、submit 验证、RecordData 字段、record 渲染的完整规格
- `proposal.md#模板内容模式` — test-gen 和 eval 两种模板内容模式的详细定义
- `proposal.md#Key-Risks` — 模板内容准确性风险和 CategoryEval 验证字段匹配风险

## Acceptance Criteria

- [ ] `prompt/data/` 包含 `test-gen-journeys.md`、`test-gen-contracts.md`、`eval-journey.md`、`eval-contract.md`
- [ ] `Synthesize()` 对 `test.gen-journeys`、`test.gen-contracts`、`eval.journey`、`eval.contract` 返回有效 prompt
- [ ] `CategoryForType("eval.journey")` 返回 `CategoryEval`（非 `CategoryCoding`）
- [ ] `CategoryForType("unknown.type")` 返回 `CategoryCoding` 并通过 `log.Printf` 输出警告
- [ ] `forge submit-task` 对 eval 任务接受含 summary/findings 的提交，拒绝仅含 testsPassed/coverage 的提交
- [ ] `RenderRecord` 对 CategoryEval 使用 eval 专用 record 模板，渲染输出包含 ScoreFormatted/FindingsFormatted/SeverityFormatted/PassedFormatted
- [ ] `plugins/forge/skills/submit-task/data/record-format-eval.md` 存在且包含 score/findings/severity/passed 字段定义
- [ ] `category_test.go`: CategoryEval 正向/负向/边界测试用例
- [ ] `submit_test.go`: CategoryEval 验证分支测试
- [ ] 所有现有测试通过

## Hard Rules

- RecordData eval 字段使用无前缀命名（`Score`/`Findings`/`Severity`/`Passed`），遵循现有惯例
- test-gen 模板遵循 test-gen-scripts.md 模式（skill 委托 + 2-step workflow）
- eval 模板使用全新模式（quality evaluation 角色 + `forge:eval` skill 委托），参考 eval skill rubric 输出契约
- CategoryEval 验证分支接受 review 字段（summary/findings/severity），拒绝纯测试字段（testsPassed/coverage）

## Implementation Notes

**模板创建顺序**：先创建 test-gen 模板（模式明确，参考 test-gen-scripts.md），再创建 eval 模板（全新模式）。

**CategoryForType 修改**：在 `eval.` 前缀分支之前插入，确保 eval 类型不再落入 default。default 分支添加 `log.Printf` 警告。

**RecordData 字段设计**：
- `Score float64 json:"score,omitempty"` — eval 评分（0-1000）
- `Findings []string json:"findings,omitempty"` — eval 发现问题列表
- `Severity string json:"severity,omitempty"` — 问题严重程度（critical/major/minor）
- `Passed bool json:"passed,omitempty"` — eval 是否通过质量门控

**eval record 模板**：字段设计基于 eval skill 输出契约（`plugins/forge/skills/eval/rubrics/` 的 score/findings/severity/passed 四维度）。

**render.go 需同步更新**：RecordTemplateData 添加 eval 格式化字段 + NewRecordTemplateData 填充逻辑。
