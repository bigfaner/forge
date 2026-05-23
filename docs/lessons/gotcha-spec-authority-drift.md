---
created: "2026-05-23"
tags: [architecture, testing]
---

# 规范文档权威性漂移：任务文档已指引但仍偏离

## Problem

修改 Forge 测试流水线的 6 个 skill 文件时，产生了 43 处与 tech-design.md 不一致的偏差。包括：

- 路径错误：`_contracts/` 应为 `contracts/`、`testing/journeys/<name>.md` 应为 `testing/<journey>/journey.md`
- 加载机制错误：用 `docs/conventions/` + `domains` frontmatter 过滤，应为 `docs/conventions/testing/index.md` 二级索引
- Convention section 名称错误：`Framework/Assertion/Tags/Result Format` 应为 `framework/discovery/structure/assertions`
- 旧代码中过时的模式被传播到新代码中

用户质问："为什么这么多地方没有遵循 tech-design.md？"

## Root Cause

### Level 0: 合成 prompt 模板未强制读取 Reference Files

task-executor agent 的执行流程是：`forge prompt get-by-task-id` → 合成模板 → agent 执行。检查了所有相关模板后发现：

| 模板 | Step 1 对 Reference Files 的处理 | Acceptance Criteria 验证 |
|------|----------------------------------|--------------------------|
| `coding-enhancement.md` | 只说"读任务文件"，未要求读 Reference Files 列出的文档 | 无显式验证步骤 |
| `coding-refactor.md` | 同上 | 无显式验证步骤 |
| `coding-feature.md` | 同上 | 无显式验证步骤 |
| `doc.md` | 说"Identify all reference files listed in the task and read them"（较好） | Self-Check 有 Completeness 但未逐条对 AC 验收 |

**关键发现**：coding.* 模板对 Hard Rules 用了 `<IMPORTANT>` 强调，但对 Reference Files 和 Implementation Notes 没有同等对待。agent 读到任务文件后，Reference Files 被当作"参考信息"跳过，直接看 AC 开始编码。

**需要的模板改进**（doc* 类任务优先）：
1. Step 1 强制读取 `## Reference Files` 列出的所有文档，声明为权威来源
2. 最后一步增加"逐条验证 Acceptance Criteria"的显式验收步骤

### Level 1: 任务文档已提供充分指引，但未被遵循

任务文档（如 3.1、2.6）的 `Reference Files` 明确指向 `tech-design.md`，`Implementation Notes` 写明了具体机制（如 "Convention 加载机制为两级：先读 index.md，再按需加载"），`Hard Rules` 写了禁止事项（如 "迁移后 Convention 文件不使用 `domains` frontmatter 过滤"）。

**问题不是"任务没说"，而是"说了但执行时没照做"。** 规范已经被提炼到任务级别了，仍然被忽略。

### Level 2: 直接原因——以现有代码为参照而非以任务文档+规范为参照

修复时的工作模式是"读目标文件 → 发现问题 → 修复"，把正在编辑的文件内容当作修改的基准。没有先读任务文档中的 Reference Files 和 Implementation Notes，再对照 tech-design.md 做全局验证。

后果：旧文件中已有的错误模式（如 `_contracts/`、`domains` frontmatter）被当作"正确现状"传播到新修改中。

### Level 3: 过程原因——缺少"自顶向下验证"步骤

修复流程是 **反应式局部修复**：看到一个问题就修一个，修完就认为完成。缺少 **主动式全局验证**：先列举规范的所有要求，再逐条检查每个文件的合规性，最后一次性修复所有偏差。

这导致：
- 在一个 section 更新了路径，但同文件的另一个 section 仍引用旧路径
- 更新了 SKILL.md 但忘记更新同 skill 下的 rules/ 和 templates/ 文件
- 新增内容（如 eval task types）使用了旧路径模式

### Level 4: 认知陷阱——局部一致性优于全局一致性

LLM agent 在编辑文件 X 时，天然倾向于让 X 内部自洽（local consistency），而非让 X 对齐外部规范 Y（global consistency）。这是一种**局部优化陷阱**：

- 编辑 `gen-test-scripts/SKILL.md` 时，注意力在该文件的内部结构上
- tech-design.md 是"另一个文件"，需要额外加载和对照
- 每次修改都是"我已经修了这部分"的增量思维，而非"我需要修所有不一致"的全局思维

**本质**：规范驱动的修改需要一个"列举所有规范要求 → 逐条检查合规性"的验证步骤。这个步骤不会自然出现在迭代编辑的流程中——它必须被显式要求。

## Solution

### 立即修复

对 6 个 skill 文件做了全面审计：以 tech-design.md 为权威，列举所有路径、目录结构、加载机制的要求，逐条检查每个 skill 文件（包括 rules/ 和 templates/ 子文件），一次性修复 43 处偏差。

### 模板层改进（待实施）

对 `forge-cli/pkg/prompt/data/` 下的 doc* 模板做两项改进：

1. **强制读取 Reference Files**：Step 1 增加 `<EXTREMELY-IMPORTANT>` 块，声明 `## Reference Files` 列出的文档为权威来源，必须完整读取后才能开始实现
2. **Acceptance Criteria 逐条验收**：最后一步增加显式的 AC checklist 验证——逐条检查每个 `[ ]` checkbox，未通过则修复后重验

涉及的模板文件：
- `doc.md` — 优先（Step 1 已有较好基础，需加强 Step 3 的 AC 验收）
- `coding-enhancement.md` — Step 1 缺少 Reference Files 强调
- `coding-refactor.md` — 同上
- `coding-feature.md` — 同上

## Reusable Pattern

**规范驱动的修改流程**（适用于任何有 tech-design / spec 的场景）：

```
0. 读任务文档的 Reference Files、Implementation Notes、Hard Rules
1. 按任务文档指引加载权威规范文档（tech-design.md）
2. 提取硬性约束清单（路径、schema、名称、机制、字段名）
3. 对每个受影响文件执行"规范 → 文件"对照审计
4. 输出偏差清单（文件:行号:预期值:实际值）
5. 按清单批量修复
6. 再次 grep 验证零残留
7. 逐条验证 Acceptance Criteria
```

**反面模式**（本次教训）：
```
1. 直接读目标代码文件
2. 发现一个不一致
3. 修完就继续
4. 重复直到"看起来没问题"
（全程没有回头读任务文档和规范文档）
```

## Related Files

- `docs/features/test-capability-v2/design/tech-design.md` — 权威规范
- `forge-cli/pkg/prompt/data/doc.md` — doc 类型模板（需改进）
- `forge-cli/pkg/prompt/data/coding-enhancement.md` — coding enhancement 模板（需改进）
- `forge-cli/pkg/prompt/data/coding-refactor.md` — coding refactor 模板（需改进）
- `forge-cli/pkg/prompt/data/coding-feature.md` — coding feature 模板（需改进）
- `plugins/forge/agents/task-executor.md` — task-executor agent 定义
- `plugins/forge/commands/execute-task.md` — 任务分发命令
