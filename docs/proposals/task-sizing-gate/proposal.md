---
created: 2026-05-31
author: faner
status: Draft
intent: refactor
---

# Proposal: Task Sizing Gate

## Problem

Task 11（reorganize internal/cmd/）违反了 Forge 自身的多条任务粒度规范（multi-verb、8 文件上限），导致子代理实际工作 26 分钟后因 macOS Idle Sleep 挂起 9.2 小时。任务粒度过大扩大了 sleep 杀死连接的风险窗口。

### Evidence

- Task 11 合并了两个独立操作（subpackage 化 15 个文件 + 拆分 5 个大文件），违反 multi-verb 规则
- 实际涉及 13 个文件，超出 8 文件上限
- `breakdown-tasks` 和 `quick-tasks` 已有完整的拆分规则，但 LLM 在 Task 11 上完全无视了这些规则
- 取证报告确认：30 分钟 timeout 只是约定，Agent tool 不强制执行

### Urgency

任务粒度过大不仅增加 sleep 风险，还降低任务执行质量——Task 11 的子代理做了 133 次工具调用、93% 时间消耗在 thinking，且因 scope 过大跳过了增量验证。

## Proposed Solution

在两个层面强化 task sizing 执行力：

1. **CLI 层**：将 `forge task validate-index` 重命名为 `forge task validate`，新增 AC 数量校验（≤ 6 且 ≥ 1），程序化强制执行，不依赖 LLM 自觉
2. **Skill 层**：在 `breakdown-tasks` 和 `quick-tasks` 的 task 文件生成后、`forge task index` 之前，插入一个独立的 task sizing audit step，让 LLM 对每个 task 做聚焦自审（multi-verb、AC 跨不相关领域），发现问题自动拆分并输出报告

### Innovation Highlights

无特别创新。核心思路是将"规则写在文档里等 LLM 自觉遵守"升级为"程序化校验 + 聚焦自审"的双重防线。

## Requirements Analysis

### Key Scenarios

1. **AC 超标**：task 生成后 `forge task validate` 检测到某个 task 有 7 条 AC，返回 exit 1 + 具体错误信息，skill 的 Step 6 失败，LLM 被迫修复
2. **Multi-verb 违规**：audit step 发现 task title 为 "Reorganize cmd structure and split large files"，自动拆分为两个 task 并输出报告
3. **AC 跨不相关领域**：audit step 发现一个 task 的 AC 同时覆盖"重命名文件"和"更新测试 fixture"，拆分为独立 task
4. **正常任务**：所有校验通过，流程无中断

### Non-Functional Requirements

- 向后兼容：`validate-index` 作为 alias 保留一段时间，或直接替换（breaking change 可接受，属于内部 CLI）
- 性能：AC 校验只需解析 markdown 文本，无显著开销

### Constraints & Dependencies

- AC 数量校验依赖 task .md 文件的 `## Acceptance Criteria` section 格式一致性（`- [ ]` 前缀）
- Coding task 模板无 `## Affected Files` section，文件数量上限暂不纳入 CLI 校验
- Multi-verb 检测依赖 LLM 语义判断，无法机械化

## Alternatives & Industry Benchmarking

### Industry Solutions

CI/CD 中的 lint gate 是标准实践——在 pipeline 中加入质量门禁，不依赖开发者自觉。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零改动 | 规则形同虚设 | Rejected: Task 11 已证明无效 |
| 加 caffeinate 防休眠 | macOS 最佳实践 | 解决 sleep 问题 | 治标不治本，不解决任务粒度问题 | Rejected: 根因未消除 |
| **CLI 校验 + Skill audit** | CI lint gate 模式 | 程序化强制 + 语义兜底 | 需改 Go 代码 + skill 文件 | **Selected: 双重防线** |
| 只强化 skill 指令 | — | 最简单 | LLM 会无视 | Rejected: Task 11 已证明无效 |

## Feasibility Assessment

### Technical Feasibility

完全可行。CLI 侧扩展 `validate-index` 增加 AC 解析是字符串处理。Skill 侧插入 audit step 是 markdown 编辑。所有改动在 Forge 自身范围内，不依赖上游 Claude Code 变更。

### Resource & Timeline

小型变更：Go 代码 ~50 行（AC 解析 + 重命名），2 个 skill 文件各 ~30 行（audit step + 编号顺延）。预计 2-3 个任务可完成。

### Dependency Readiness

`forge task validate-index` 已有成熟的校验框架。Skill 文件的 step 编号是纯文档变更。无需新增基础设施。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "macOS sleep 是根因，需要基础设施层面的修复" | 5 Whys | Overturned: sleep 是触发器，任务粒度过大是根因。粒度合理的任务（5-8 分钟）大幅缩小 sleep 风险窗口 |
| "规则写在文档里，LLM 会遵守" | Assumption Flip | Overturned: Task 11 证明 multi-verb 和 8 文件上限规则被完全无视 |
| "文件数量也应该在 CLI 层校验" | Occam's Razor | Refined: coding task 模板无 Affected Files section，强行加会增加模板改动范围。AC 数量校验 + skill audit 已足够 |
| "caffeinate 是必要的兜底" | XY Detection | Confirmed: caffeinate 防的是 sleep，但根因是任务太大。正确的兜底是 Agent tool timeout（上游能力） |

## Scope

### In Scope

- CLI：`validate-index` → `validate`，新增 AC 数量校验（≤ 6 且 ≥ 1）
- CLI：更新所有 `validate-index` 调用点为 `validate`
- Skill：`breakdown-tasks/SKILL.md` 插入 task sizing audit step（Step 4 之后）
- Skill：`quick-tasks/SKILL.md` 插入 task sizing audit step（Step 3 之后）
- Skill：两个 skill 中 `validate-index` 引用更新为 `validate`
- Skill：两个 skill 的后续 step 编号顺延

### Out of Scope

- Agent tool timeout（上游 Claude Code 能力缺口）
- API client read timeout（上游 Claude Code 能力缺口）
- Caffeinate 防休眠机制（根因已通过任务粒度解决）
- Coding task 模板增加 `## Affected Files` section
- 跨平台休眠防护（Linux/Windows）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Audit step 的自动拆分引入新错误 | M | M | 拆分后重新跑 validate 校验；输出报告便于人工审查 |
| LLM 在 audit step 仍忽略 multi-verb | L | M | CLI 层的 AC 校验作为兜底——粒度合理的 multi-verb task 通常 AC 也多，会被 AC ≤ 6 拦截 |
| `validate-index` 重命名破坏外部脚本 | L | L | Forge CLI 是内部工具，无外部消费者 |
| AC 解析误判（非标准格式） | L | M | 解析 `## Acceptance Criteria` 下的 `- [ ]` 行，格式由模板保证 |

## Success Criteria

- [ ] `forge task validate docs/features/<slug>/tasks/index.json` 成功校验 index.json 结构 + 所有 task 文件的 AC 数量
- [ ] AC > 6 时返回 exit 1 + 错误信息（包含 task 文件名和 AC 数量）
- [ ] AC = 0 时返回 exit 1 + 错误信息
- [ ] `forge task validate-index` 不再存在或已别名化
- [ ] `breakdown-tasks` 在写 task 文件后、`forge task index` 前执行 task sizing audit
- [ ] `quick-tasks` 同上
- [ ] audit 发现 multi-verb 或跨域 AC 时，自动拆分并输出报告
- [ ] 所有 step 编号连续无跳跃

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
