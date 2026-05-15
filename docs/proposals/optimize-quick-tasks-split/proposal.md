---
status: draft
created: 2026-05-15
---

# 优化 quick-tasks 任务拆分机制

## Problem

quick-tasks 的任务模板从 breakdown-tasks 照搬了 `Affected Files` 区块（Create/Modify/Delete 三级表格），要求文件级精度。但 proposal 只提供能力级描述（如 "Add --type argument to gen-test-scripts"），不包含文件路径。

这个精度错配导致：

1. **Over-research**：SKILL.md Step 2 的 "Determine affected file paths" 指令触发 agent 研究下游技能文件（如 430 行的 gen-test-scripts/SKILL.md）来发现精确路径，浪费 context window
2. **不准确的任务分类**：agent 为填充表格而猜测文件路径，但 proposal 没有这个信息，结果依赖研究而非推断
3. **Scope 判定依赖不存在的输入**：Scope Assignment 算法基于文件路径分类 frontend/backend/undetermined，但路径本身是研究出来的

根因：**任务模板要求的信息精度与输入文档的信息密度不匹配**。

## Proposed Solution

### 1. 删除 implementation 模板的 Affected Files 区块

`templates/task.md` 去掉 `## Affected Files` 区块（Create/Modify/Delete 表格）。

Rationale：
- executor 在执行时从代码库自然发现文件路径，不需要创建时预填
- Reference Files（目录级路径）已经为 executor 提供了检索起点
- 消除了 agent 研究下游文件来填充表格的压力

**`templates/task-doc.md` 保留 Affected Files**。doc 任务的交付物就是文件本身，目标路径在创建时可确定，不需要研究。

### 2. Scope Inference 替代 Scope Assignment

删除 SKILL.md 中基于文件路径的 Scope Assignment 算法，替换为从任务描述语义推断：

```
Scope Inference:
- 描述涉及 UI、页面、组件、样式 → scope: "frontend"
- 描述涉及 API、服务端、数据库、CLI 工具 → scope: "backend"
- 混合或无法判断 → scope: "all"
```

Rationale：
- proposal 有语义信息但无文件路径，语义推断匹配输入的信息密度
- 对 quick-tasks 的精度需求足够（mixed 项目区分 frontend/backend/all）
- 无需读取任何文件

### 3. 按功能步骤拆分

当前规则：每个 In Scope 条目 → 一个任务。

新增拆分规则：如果一个 In Scope 条目包含多个可独立验证的功能步骤，按步骤拆分为多个任务。每个任务是可独立验证的功能单元。

**10 个任务硬限制不变**。拆分后总数仍 ≤ 10。如果 In Scope 需要超过 10 个任务，推荐走 full pipeline。

### 4. Reference Files 填充规则

Reference Files 从 proposal 的 In Scope 条目推断目录级路径，作为 executor 的检索起点。不需要精确文件路径。

示例：proposal 说 "Add --type argument to gen-test-scripts" → Reference Files 填 `plugins/forge/skills/gen-test-scripts/`。

## Scope

### 改动（3 个文件）

| 文件 | 改动内容 |
|------|----------|
| `quick-tasks/SKILL.md` | Step 2 重写：删除 "Determine affected file paths"（L54），删除 Scope Assignment 段（L60-69），新增 Scope Inference 规则 + 按功能步骤拆分规则；Step 3 删除 "Fill Affected Files"（L84）；Output Checklist 删除 Affected Files 检查项（L179） |
| `quick-tasks/templates/task.md` | 删除 `## Affected Files` 区块（L21-36） |
| `quick-tasks/templates/task-doc.md` | 不改（保留 Affected Files） |

### 不改

- **Go CLI**：无代码改动
- **breakdown-tasks**：输入（tech-design）有精确文件路径，当前机制正确
- **executor 命令**（execute-task、run-tasks）：无需变更

## Success Criteria

- [ ] quick-tasks 执行时 agent 不再读取下游技能文件来发现文件路径
- [ ] implementation 任务文件无 Affected Files 区块
- [ ] doc 任务文件保留 Affected Files 区块
- [ ] scope 从任务描述正确推断（frontend/backend/all）
- [ ] 多步骤 In Scope 条目正确拆分，总数 ≤ 10
- [ ] `forge task index` 正常生成 index.json
- [ ] executor 能基于 Description + Reference Files 完成实现

## Key Risks

| 风险 | 缓解 |
|------|------|
| Scope 推断不准（描述语义模糊） | fallback 到 "all"，全量质量门兜底 |
| 拆分粒度过细导致任务碎片化 | 10 个硬限制自然约束 |
| executor 缺少文件路径指引效率降低 | Reference Files 提供目录级起点；executor 本身有代码探索能力 |

## Related

- Lesson: `docs/lessons/gotcha-task-derivation-over-research.md`
