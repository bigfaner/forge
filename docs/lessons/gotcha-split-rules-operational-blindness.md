---
created: "2026-05-28"
tags: [architecture]
---

# Split Rules 只检查功能粒度，遗漏操作粒度

## Problem

尽管昨天已优化了 quick-tasks 的拆分规则（AC ceiling ≤6、multi-verb detection、inline Reference Files、complexity 字段），Task 2（精简全部 prompt 模板）仍然被生成为单个任务，导致 task-executor 执行 65 次 Edit 调用、修改 21 个文件、总耗时 6.7 分钟后被中断。

## Root Cause

因果链（3 层）：

1. **表面现象**：Task 2 通过了全部 3 条拆分规则，但实际操作量（65 edits × 21 files）远超单次 task-executor dispatch 的合理上限
2. **直接原因**：拆分规则只评估"功能粒度"——AC 数量、动词数量、是否可独立验证——但没有评估"操作粒度"——需要修改多少文件、每文件需要多少种编辑操作
3. **根因**：拆分规则的设计假设是"功能复杂度 ≈ 操作复杂度"，但这两者是正交维度。一个功能简单的任务（5 条 AC、单一动词）可能需要大量的机械重复编辑（N 个文件 × M 种编辑类型），而拆分规则对这种"功能简单但操作繁重"的情况完全盲视

### 为什么昨天的优化没覆盖

| 昨天的优化 | 解决的问题 | 未覆盖的维度 |
|-----------|-----------|------------|
| AC ceiling ≤6 | 功能范围过大 | 不反映文件数量和编辑次数 |
| Multi-verb detection | 多步骤任务 | Task 2 是单一动词"slim" |
| Inline Reference Files | scope 越界（读 proposal） | scope 越界（编辑 scope 外文件）仍存在 |
| complexity 字段 | 标记复杂度 | 只是 metadata，不触发拆分 |

## Solution

在拆分规则中增加第 4 条 **Operational Ceiling**：

> **Operational Ceiling**: If a task requires modifying >8 files with the same pattern, split by file group (e.g., by complexity tier, by feature area, by directory). Each sub-task targets ≤8 files.

判定方式：在 Step 2 Derive Tasks 时，如果 In Scope bullet 明确提及"N 个文件/模板/模块"且 N > 8，按文件的复杂度分层拆分。

示例（本次 Task 2 应拆为）：
- **2a**: Slim 5 coding-* 模板（最复杂：role + CODING_PRINCIPLES + AC block + Record Fields）
- **2b**: Slim gate/doc/doc-review + validation-* 模板（中等：role + AC block + Record Fields）
- **2c**: Slim test-* + code-quality 模板（最简：role + Step 2 desc）

### scope 越界的补充修复

除 operational ceiling 外，还需：
1. **枚举 scope 内文件**：当 task 涉及"N 个文件"时，在 Implementation Notes 中明确列出文件名（如 "coding-feature.md, coding-enhancement.md, ..."），不使用"全部"等模糊措辞
2. **排除 scope 外文件**：在 Hard Rules 中增加 "仅修改以下文件：" 列表

### Directory-driven scope creep（第二种越界路径）

已有的 `gotcha-task-reference-files-scope-creep` 解决了 **spec-driven scope creep**：inline Reference Files 阻止 executor 读 proposal 从而发现其他 task 的需求。但 Task 2 暴露了一条全新的越界路径：

| 越界类型 | 触发路径 | 现有防御 | 缺失防御 |
|---------|---------|---------|---------|
| **spec-driven** | 读 proposal → 发现额外需求 → 越界修复 | inline Reference Files ✓ | — |
| **directory-driven** | ls 目录 → 看到全部文件 → 读全部 → 编辑全部 | 无 | 枚举 scope 文件 + Hard Rules 写边界 |

证据链（Task 2 agent 行为）：
- Line 7: `ls forge-cli/pkg/prompt/templates/` → 看到 21 个文件
- Line 14-44: 分 3 批读完全部 21 个文件（包括 scope 外的 doc-drift/doc-review/doc-summary/doc-consolidate/fix-record-missed/eval-contract/eval-journey）
- Line 85+: 对全部 21 个文件执行 Edit

**为什么 inline Reference Files 没阻止这次越界**：Reference Files 列了 5 个文件作为代表性示例（coding-feature, coding-enhancement, gate, test-run, code-quality-simplify），agent 将其理解为"这些是需要精简的文件类型示例"而非"只能修改这些文件"。Agent ls 目录后看到 21 个文件，认为"这些都有类似冗余，应一并处理"。

**核心区别**：inline Reference Files 解决的是 **读越界**（不读外部 spec），但不解决 **写越界**（编辑 scope 外文件）。两者的防御机制完全不同——前者控制信息输入边界，后者需要控制文件操作边界。

## Reusable Pattern

- **功能粒度 ≠ 操作粒度**：AC 数量和动词数量衡量的是"任务要达成什么"，不衡量"任务需要做多少次机械操作"。两者必须独立评估
- **批量重复操作是拆分信号**：当一个任务的核心模式是"对 N 个文件做相同类型的编辑"，N > 8 时应按文件分组拆分
- **枚举 > 模糊**：当 task 描述涉及具体文件列表时，枚举文件名比"全部"/"所有"更安全——枚举即是 scope 边界声明

## Related Files

- `skills/quick-tasks/SKILL.md` — 拆分规则需增加 operational ceiling
- [[gotcha-prompt-template-complexity-agnostic]] — prompt 模板复杂度盲视（同属复杂度评估维度缺失）
- [[gotcha-task-reference-files-scope-creep]] — Reference Files scope 越界（scope 边界声明问题）
- [[gotcha-task-executor-thinking-overhead]] — thinking 开销叠加（操作量大的直接后果）
