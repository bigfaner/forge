---
created: "2026-05-28"
tags: [architecture, testing]
---

# 纯 .md 文件的任务不应使用 coding.* 类型

## Problem

Feature `slim-task-prompt-templates` 的任务 1、2、2a-2c、3 只涉及修改 `.md` prompt template 文件，但任务类型被设为 `coding.enhancement`。这导致：

- Agent 收到 TDD/coverage 指令，对纯 markdown 无意义
- Submit 时强求测试证据（testsPassed > 0 或 coverage）
- Quality gate 执行编译+lint+test，对非代码文件纯属浪费
- 自动生成不必要的测试流水线任务

## Root Cause

4层因果链：

1. **L1 表面**：任务 1-3、2a-2c 被标为 `coding.enhancement`，但它们只修改 `.md` 文件
2. **L2 直接**：LLM agent 在创建任务时，没有遵循 SKILL.md 中的分类规则（"If all affected files are non-compilable, type must be `doc`"）
3. **L3 规则歧义**：分类规则说"看 affected files 是否为 non-compilable"，但 `.md` 文件位于 `pkg/prompt/templates/` 这种代码风格的路径下，agent 被 `pkg/` 目录语境误导，将"改进现有 prompt 模板"理解为"改进现有行为" → `coding.enhancement`，而非跳到 `task-doc.md` 模板
4. **L4 机制缺陷**：
   - **模板层**：`task.md` 硬编码 `type: "coding.feature"`，agent 手动改成 `coding.enhancement` 但从未考虑应该用 `task-doc.md` 模板（硬编码 `type: "doc"`）
   - **CLI 层**：`BuildIndex` 只校验 `type` 非空和是否在 `ValidTypes` 中，不校验 `coding.*` 类型任务是否实际引用了可编译文件。没有任何 warning 机制兜底

**附带发现**：task 4（Refactor Go metadata parsing）涉及 `.go` 代码重构，用 `coding.enhancement` 也不准确——应该是 `coding.refactor`。这说明即使涉及代码的任务，分类也不够精确，agent 倾向于默认 `coding.enhancement` 而非精确匹配 `coding.refactor`/`coding.cleanup`

## Solution

- 纯 `.md` 文件编辑任务应使用 `doc` 类型
- `doc` 类型：跳过 quality gate，使用文档导向 prompt 模板，record 格式要求 referencedDocs/docMetrics
- 仅当任务涉及可编译/可执行文件（`.go`, `.ts`, `.py` 等）时才使用 `coding.*` 类型
- 例外：如果 `.md` 文件是代码生成目标（如生成的代码文档），且修改它需要运行验证，才考虑 `coding.*`

### 修复点

| 层级 | 修复 | 效果 |
|------|------|------|
| SKILL.md | 在 Type Assignment 规则中增加：".md 文件即使在 `pkg/`/`src/` 下，只要不是可编译/可执行的，类型就必须是 `doc`" | 消除 L3 路径歧义 |
| Go CLI | `BuildIndex` 中增加警告：`coding.*` 类型任务如果所有 affected files 都是 `.md`/`.yaml`/`.json`，emit warning | 捕获 L4 遗漏，兜底防护 |

## Reusable Pattern

创建任务时，用"操作语义"而非"文件位置"判断类型：

1. **该任务产出的文件是否可编译/可执行？** → 是 → `coding.*`
2. **该任务是否只修改文档/规范/markdown？** → 是 → `doc`
3. **该任务是否涉及测试生成/执行？** → 是 → `test.*`
4. **不确定时，选影响最小的类型**（如 `doc` 可以后续升级为 `coding.*`，反过来则浪费 agent 时间在无意义的测试上）

判断标准：Forge 的 `IsTestableType()` 只对 `coding.*` 和 `code-quality.simplify` 返回 true。如果一个任务的产出不触发 quality gate 反而更合理，那它就不应该是 `coding.*` 类型。

## Related Files

- `forge-cli/pkg/task/types.go` — 类型常量定义
- `forge-cli/pkg/task/category.go` — CategoryForType() 分类函数
- `forge-cli/pkg/task/build.go` — IsTestableType() 判断逻辑
- `docs/features/slim-task-prompt-templates/tasks/index.json` — 当前任务索引
