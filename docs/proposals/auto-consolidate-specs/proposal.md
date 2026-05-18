---
created: 2026-05-18
author: "fanhuifeng"
status: Draft
---

# Proposal: Auto-consolidate-specs in Task Pipeline

## Problem

consolidate-specs 在 task pipeline 中无法真正自动化：Full pipeline 遇到 `[CROSS]` 规则时阻塞等待用户确认，Quick pipeline 根本不生成该任务。这破坏了 pipeline 的自动执行连续性。

### Evidence

- `SKILL.md` Step 6: 非交互模式下遇到 `[CROSS]` 项时将任务标记为 `blocked`，附注 "User review required for integration"
- `.forge/config.yaml`: `auto.consolidateSpecs.quick: false`，Quick pipeline 不生成此任务
- quick-tasks 的 `index.json` 模板中没有 consolidate-specs 任务槽位

### Urgency

每次 `/run-tasks` 执行到 consolidate-specs 时要么跳过（Quick）、要么中断（Full），导致 spec 知识无法持续积累。随着 feature 增多，project-level specs 的 drift 会越来越严重。

## Proposed Solution

将 consolidate-specs 改为"全自动保存 + 事后 git diff 审查"模式：pipeline 中自动执行全部流程（提取、集成、drift 检测、修复），变更单独提交并带 `[auto-specs]` 标记，用户事后通过 git diff 审查并可轻松 revert。

### Innovation Highlights

无特别创新——本质是移除不必要的交互阻塞点，遵循"commit early, review later"的 CI 理念。核心洞察：spec 集成的风险远低于代码变更，因为 spec 错误不会导致运行时故障，且 git revert 可完美回退。

## Requirements Analysis

### Key Scenarios

- **Happy path**: `/run-tasks` 执行到 consolidate-specs 任务时，自动提取规则、集成到 project-level、检测并修复 drift，生成带标记的单独 commit
- **No feature docs**: Quick 模式可能没有 PRD/design，此时走 drift-only 路径（现有行为）
- **冲突/重叠**: 与已有 project-level spec 有 >50% 重叠时，仍执行合并但 commit message 中标注 warning
- **用户审查**: 用户事后 `git log --grep="[auto-specs]"` 查看所有自动集成，`git revert` 回滚不需要的

### Non-Functional Requirements

- 兼容性：不影响用户手动调用 `/consolidate-specs` 的交互行为
- 可追溯性：每次自动集成都必须有独立的、可 revert 的 commit

### Constraints & Dependencies

- 依赖现有的 consolidate-specs skill 和 task pipeline 机制
- 需要修改 skill 模板、config、quick-tasks 模板

## Alternatives & Industry Benchmarking

### Industry Solutions

类似 CI 中的 auto-formatting（prettier、gofmt）——自动应用变更，开发者事后 review diff。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零风险 | pipeline 持续中断，specs 持续 drift | Rejected: 问题会恶化 |
| 保存为草稿 | — | 安全，用户手动确认每条 | 增加审查负担，违背自动化初衷 | Rejected: 多了一步手动操作 |
| **直接集成 + 事后审查** | CI auto-format | 完全自动化，git revert 保险 | 可能集成错误规则 | **Selected: spec 错误无运行时风险，回退容易** |

## Feasibility Assessment

### Technical Feasibility

完全可行。变更集中在 SKILL.md 的 Step 6（移除 block 逻辑）和 quick-tasks 模板（加任务槽位）。

### Resource & Timeline

小型变更，涉及 4-5 个文件，工作量可控。

### Dependency Readiness

无外部依赖，所有组件已就绪。

## Scope

### In Scope

- 修改 `SKILL.md` Step 6：非交互时自动集成所有 `[CROSS]` 项，不阻塞
- 修改 `SKILL.md` Step 11：生成带 `[auto-specs]` 标记的单独 commit
- 修改 `quick-tasks` 模板：添加 T-quick-specs-1 任务槽位
- 修改 `config`：`auto.consolidateSpecs.quick: true`
- 修改 consolidate-specs task template：移除用户确认步骤
- 修改 `doc-generation-consolidate.md` prompt template：适配全自动模式

### Out of Scope

- 用户手动调用 `/consolidate-specs` 的交互行为（保持不变）
- 新增 review 工具或命令（用 git diff/revert 即可）
- drift 检测逻辑的修改（已经是自动的）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 自动集成错误规则到 project-level | L | M | `[auto-specs]` 标记使 revert 简单；spec 错误无运行时风险 |
| drift 自动修复误删正确规则 | L | M | 单独 commit，git revert 可恢复 |
| 用户忘记审查 | M | L | spec drift 影响有限，下次 consolidate-specs 会重新检测 |

## Success Criteria

- [ ] `/run-tasks` 在 Quick 模式下自动生成并执行 consolidate-specs 任务
- [ ] `/run-tasks` 在 Full 模式下 consolidate-specs 不再阻塞等待用户确认
- [ ] 所有自动集成的变更在单独的 commit 中，commit message 包含 `[auto-specs]`
- [ ] 用户手动 `/consolidate-specs` 的交互行为不受影响
- [ ] `git log --grep="[auto-specs]"` 能找到所有自动集成的 commit

## Next Steps

- Proceed to task generation via `/quick-tasks`
