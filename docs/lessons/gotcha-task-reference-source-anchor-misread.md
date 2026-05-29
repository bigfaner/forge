---
created: "2026-05-29"
tags: [tooling, interface, architecture]
---

# quick-tasks 生成的 source 引注存在路径错误和章节伪造

## Problem

Task 1 的 task-executor 执行记录中报告 "No proposal.md found"，最终只能依赖 AC 和 Implementation Notes 作为规格来源。task 文件中的 Reference Files 包含 `(source: proposal.md#Core-Data-Structure)` 等引注，但这些引注有三层错误：路径不存在、章节名伪造、格式本身不适合 agent 消费。

## Root Cause

因果链（4 层）：

1. **L1 症状**: task-executor 无法找到 proposal.md 文件，所有 `source:` 引注失效
2. **L2 三重错误**:
   - **路径错误**: `proposal.md` 不是有效路径。实际路径是 `docs/proposals/pipeline-topology-registry/proposal.md`
   - **章节名伪造**: `proposal.md#Intent-Gate-Functions` 和 `proposal.md#Dependency-Resolver-Functions` 在 proposal 中**根本不存在**。实际章节名是 `Predefined Gate, Condition & Resolver Functions`。LLM 在生成 source 注解时编造了看似合理的章节名，而非从实际 proposal 中提取
   - **格式歧义**: `(source: ...)` 是给人类看的引注格式，不是文件路径，但 task-executor 尝试将其作为路径读取
3. **L3 设计缺陷**: quick-tasks skill 的 "Reference Files Generation" 规则要求 LLM 生成 `(source: proposal.md#Section-Title)` 格式的注解，但：
   - 没有规定使用 proposal 的完整路径 `docs/proposals/<slug>/proposal.md`
   - 没有要求验证章节名是否真实存在于 proposal 中
   - 没有区分"人类引注格式"和"agent 可读文件路径"
4. **L4 根因**: source traceability 被设计为 LLM 自由编写的元数据，缺乏与实际文档的结构化验证。LLM 倾向于生成看似合理但虚构的引用（hallucinated references），这在 source annotation 场景中尤为隐蔽——因为它看起来很"正确"

## Solution

两阶段修复：

**短期（task-executor 侧）**: task-executor 应忽略 `(source: ...)` 注解中的文件路径部分，将实际 `<file-path>`（冒号前的部分）作为唯一可读取的文件引用。`source:` 仅作为人类可读的溯源信息。

**长期（quick-tasks skill 侧）**: 修改 Reference Files Generation 规则：
1. 使用 proposal 完整路径：`docs/proposals/<slug>/proposal.md`
2. 生成 source 注解前，先提取 proposal 的实际 markdown headers，确保章节名与真实标题匹配
3. 或者：完全移除 `(source: ...)` 注解格式，改为在 Reference Files 中直接引用 proposal 路径加具体行号范围

## Reusable Pattern

当 LLM 生成"溯源注解"类元数据时（source references、provenance tags、section pointers）：

1. **路径必须经过文件系统验证** —— 不允许使用缩写路径（`proposal.md` → 必须是完整相对路径）
2. **章节名必须从实际文件提取** —— 不允许 LLM 推断/编造章节标题。先 grep headers，再用实际值
3. **区分人类注解与 agent 可读路径** —— 如果格式同时服务于两个消费者，必须明确标注哪些部分可被机器解析

## Example

```
# 错误（当前生成的）
- `forge-cli/pkg/task/pipeline.go`: New file (source: proposal.md#Core-Data-Structure)
  → proposal.md 路径不存在，Core-Data-Structure 是编造的章节名

# 正确方案 A：完整路径 + 真实章节名
- `forge-cli/pkg/task/pipeline.go`: New file (source: docs/proposals/pipeline-topology-registry/proposal.md#core-data-structure)
  → 路径可读，章节名匹配实际 header "#### Core Data Structure"

# 正确方案 B：省略 source，只用行号
- `forge-cli/pkg/task/pipeline.go`: New file — 见 proposal line 45-133
  → 无歧义，agent 可直接定位
```

## Related Files

- `plugins/forge/skills/quick-tasks/` — quick-tasks skill 的 Reference Files Generation 规则
- `plugins/forge/agents/task-executor/` — task-executor prompt 定义
- `docs/features/pipeline-topology-registry/tasks/1-define-pipeline-registry.md` — 问题 task 文件

## References

- quick-tasks skill 的 "Reference Files Generation" 规则，source traceability 格式定义
- Task 1 执行记录：`docs/features/pipeline-topology-registry/tasks/records/1-define-pipeline-registry.md`
