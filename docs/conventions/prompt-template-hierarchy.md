---
title: "Prompt Template Instruction Hierarchy"
domains: [prompt, template, CODING_PRINCIPLES, EXTREMELY-IMPORTANT, HARD-GATE, TASK-CONSTRAINTS, instruction]
---

# Prompt Template Instruction Hierarchy

Forge 的提示词模板采用三级指令层次，由弱到强：

1. **`<CODING_PRINCIPLES>`** — 行为指南（自律遵循）。定义 agent 在执行任务时应遵循的行为准则，如"先思考再编码""最小修改范围"。用大写 XML 标签包裹以提高 LLM 关注度。目前仅在 `commands/fix-bug.md` 中使用。
2. **`<EXTREMELY-IMPORTANT>`** — 任务级硬约束（必须遵循）。定义不可违反的规则，如"Hard Rules override your default approach"。广泛用于 skills、commands 和 agents 中。
3. **`<HARD-GATE>`** — 流程级强制检查点（不可绕过）。定义必须通过的验证条件，如"If the bug cannot be reproduced, STOP"。用于需要前置条件验证的 skill 入口。

### 位置规则

- `<CODING_PRINCIPLES>` 置于角色描述之后、工作流步骤（`## Workflow`）之前。
- `<EXTREMELY-IMPORTANT>` 置于工作流步骤内部，靠近其约束的步骤。
- `<HARD-GATE>` 置于需要强制检查的步骤内部。

### 设计原则

- 无时序冲突：行为守则不与工作流步骤产生执行顺序矛盾。
- 无语义重叠：与现有规则重叠时合并为统一表述，不并存。
- 指令层次清晰：三级标签各有明确语义，不混用。

### 模板级约束与任务级约束的区分

- **`<TASK-CONSTRAINTS>`** — 模板级工作流约束。定义模板自身对 agent 执行流程的强制要求（如"必须通过 Skill 调用，禁止直接执行"）。用于 test.* 模板，由模板设计者设定。
- **Hard Rules** — 任务级约束。定义具体任务文件中的不可违反规则（如文件范围限制、命令限制），由任务创建者设定，模板通过 `<IMPORTANT>` 标签注入 agent 上下文。

两者命名不同以避免混淆：TASK-CONSTRAINTS 是模板设计约束，Hard Rules 是任务执行约束。
