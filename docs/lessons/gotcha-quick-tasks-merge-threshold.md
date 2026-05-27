---
created: "2026-05-27"
tags: [architecture, testing]
---

# Quick-Tasks 合并阈值导致粗粒度任务

## Problem

Task 1 合并了常量重命名、物理目录扁平化、路径引用更新、testrunner 路径确认 4 个独立步骤，导致 executor 耗时 18 分钟，且越界覆盖了 tasks 2-4 的部分范围。验收时 AC5/AC6（目录扁平化）未完成但被误判为已完成。

## Root Cause

因果链（3 层）：

1. **表面现象**：task 1 的 scope 包含多个独立关注点，验收标准 11 条中 2 条未满足却被标记 completed
2. **直接原因**：quick-tasks Step 2 规则 "merge if <30min" 将 4 个各 <30min 的独立步骤合并为一个 2h 任务。合并的依据是时间估算而非功能边界
3. **根因**：quick-tasks SKILL.md 的拆分规则存在两个缺陷：
   - **合并维度错误**：用时间估算（<30min）作为合并标准，但正确的合并标准应该是"是否可以独立验证"。常量重命名和目录扁平化各自有独立的 grep 验证命令，应该分开
   - **Proposal In Scope 条目粒度不均**：proposal 的 In Scope 用一个条目覆盖了 "Go 代码层的术语、常量、路径、函数名统一"，quick-tasks 遵循 "one task per In Scope bullet" 但没有要求对过宽的 bullet 做进一步拆分

## Solution

修改 quick-tasks SKILL.md 的拆分规则：

1. 合并标准改为"是否可以独立验证"而非时间估算
2. 添加规则：当 In Scope bullet 包含多个独立动词（如"重命名 + 扁平化 + 确认"）时，按功能边界拆分
3. 每个任务的验收标准不超过 6 条——超过意味着 scope 过大

## Reusable Pattern

- **独立验证性 > 时间估算**：判断是否拆分任务时，问"这个任务是否有独立的 grep/build/测试验证命令？"而非"这个任务是否 <30min？"
- **验收标准计数**：如果任务的验收标准超过 6 条，大概率 scope 过大，应该拆分
- **多动词检测**：任务描述中出现"并""以及""同时"等连接多个独立动作时，按动词拆分

## Related Files

- `plugins/forge/skills/quick-tasks/SKILL.md` — quick-tasks 拆分规则
