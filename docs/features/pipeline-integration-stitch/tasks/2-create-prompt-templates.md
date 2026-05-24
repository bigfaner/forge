---
id: "2"
title: "创建 4 个缺失的执行阶段 prompt 模板"
priority: "P0"
estimated_time: "2h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 2: 创建缺失的 prompt 模板

## Description

创建 4 个缺失的 prompt 模板文件：`test-gen-journeys.md`、`test-gen-contracts.md`、`eval-journey.md`、`eval-contract.md`。这些是**执行阶段模板**（agent 运行时指令），不是 autogen 规划阶段模板（生成 .md 文件）。放置在 `forge-cli/pkg/prompt/data/` 目录下，由 Task 1 的自动发现机制识别。

## Reference Files
- `docs/proposals/pipeline-integration-stitch/proposal.md` — Source proposal
- `forge-cli/pkg/prompt/data/code-quality-simplify.md` — 参考模板（执行阶段结构）
- `forge-cli/pkg/prompt/data/test-gen-scripts.md` — 参考模板（test 类型结构）
- `forge-cli/pkg/prompt/data/test-run.md` — 参考模板（test 运行结构）
- `forge-cli/pkg/task/autogen.go` — autogen 模板内容参考（仅参考类型信息，不照搬结构）
- `plugins/forge/skills/gen-journeys/SKILL.md` — gen-journeys skill（理解 agent 期望行为）
- `plugins/forge/skills/gen-contracts/SKILL.md` — gen-contracts skill
- `plugins/forge/skills/eval/SKILL.md` — eval skill

## Acceptance Criteria
- [ ] `forge prompt get-by-task-id` 对 `test.gen-journeys`、`test.gen-contracts`、`eval.journey`、`eval.contract` 均返回有效 prompt
- [ ] 每个模板包含：任务上下文说明、输入格式描述（从 task .md 文件读取）、期望输出格式、质量标准
- [ ] eval 模板包含 MainSession 标记和 subagent 调度指令
- [ ] Task 1 的 init-time 校验通过（4 个新模板文件均被识别）
- [ ] 所有现有测试通过

## Hard Rules
- 不可照搬 autogen.go 的规划阶段模板结构
- 模板使用 `{{TASK_ID}}`、`{{TASK_FILE}}`、`{{FEATURE_SLUG}}` 等占位符

## Implementation Notes
- 执行阶段模板 vs 规划阶段模板的区别：规划模板生成 .md 任务文件内容，执行模板生成 agent 运行时指令
- eval 模板需要包含 rubric 评分指令和 subagent 调度逻辑（eval.journey/eval.contract 是 MainSession 任务）
