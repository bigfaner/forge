---
created: "2026-06-01"
author: "faner"
status: Draft
intent: "enhancement"
---

# Proposal: Sitemap Surface Guard

## Problem

`gen-web-sitemap` skill 及其产出物 `sitemap.json` 是纯 web 专属能力（依赖 agent-browser 探索路由、抓取无障碍树），但 skill 定义和 5 个下游消费方均未做 surface 类型校验，导致非 web 项目（CLI、API、TUI、Mobile）中被误调用或误引用。

### Evidence

- `gen-web-sitemap` SKILL.md 无 surface 类型检查，可在任何项目中被直接调用
- 5 个下游 skill 无条件引用 `sitemap.json`：
  - `write-prd/SKILL.md`（Step 1 读取）无条件读取 sitemap.json
  - `write-prd/rules/ui-functions.md` 无条件要求读取 sitemap 作为 placement 依据
  - `write-prd/rules/self-check.md` 对所有 new-feature intent 执行 sitemap 路由校验，不区分 surface
  - `breakdown-tasks/rules/ui-placement.md` 仅依赖 artifact 存在性（ui-functions.md）而非 surface 类型来判断是否加载
  - `eval/rules/validate-ux-pipeline.md` 虽然表格按类型分行，但无显式守卫阻止非 web 项目访问 sitemap
- 仅 `gen-test-scripts/types/ui.md` 已有正确的守卫：`Only execute when the project has web-ui interface AND UI-type test cases exist`

### Urgency

已发生误用：非 web 项目中 agent 被引导执行 sitemap 相关操作，浪费 token 且可能产生误导性输出。随着 Forge 支持的 surface 类型增多（CLI、API、TUI、Mobile），误用频率会持续上升。

## Proposed Solution

为 `gen-web-sitemap` 和所有下游 sitemap 消费方添加基于 `forge surfaces --json` 的 surface 类型守卫：

1. **gen-web-sitemap 增加前置检查**：Step 0 检测 project surface，无 `web` 类型时中止并提示
2. **下游 skill 条件化 sitemap 引用**：读取 sitemap.json 前，先检查项目是否有 `web` surface，无则跳过相关步骤
3. **统一守卫模式**：以 `gen-test-scripts/types/ui.md` 已有守卫为参考模板，保持一致性

### Innovation Highlights

这是 Forge skill 生态中 surface-aware 条件执行的系统性加固。核心洞察是：skill 的隐式约束（名称中带 "web"）不如显式约束（代码检查 surface 类型）可靠。此模式可推广到其他 surface 专属 skill。

## Requirements Analysis

### Key Scenarios

- **纯 web 项目**：`forge surfaces --json` 返回包含 `web` 类型的 surface → gen-web-sitemap 正常执行 → 下游 skill 正常引用 sitemap.json
- **纯 CLI/API/TUI/Mobile 项目**：无 `web` surface → gen-web-sitemap 中止提示 → 下游 skill 跳过 sitemap 相关步骤
- **Monorepo 混合 surface**：`forge surfaces --json` 返回多个 surface（如 web + api）→ 检测到 `web` → gen-web-sitemap 正常执行，下游 web 相关路径正常引用
- **无 surface 配置的项目**：`forge surfaces --json` 返回空或未配置 → gen-web-sitemap 中止，下游 skill 跳过 sitemap 步骤

### Constraints & Dependencies

- 守卫必须使用 `forge surfaces --json`（而非 `--types`），因为 monorepo 可能包含多个 surface key
- 守卫逻辑是自然语言指令（写入 SKILL.md / rules 文件），不是代码 — 由 LLM agent 在执行时读取并遵循
- `gen-test-scripts/types/ui.md` 已有的守卫模式作为参考基准

## Alternatives & Industry Benchmarking

### Industry Solutions

Plugin/技能系统中对能力适用范围做显式声明是常见模式（如 VS Code extension 的 `activationEvents`、npm 的 `engines` 字段）。Forge 的 surface 类型体系提供了天然的适用范围标识。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 无开发成本 | 误用持续，token 浪费，误导输出 | Rejected: 已发生实际误用 |
| 文件存在性检查 | — | 简单 | 治标不治本：非 web 项目可能因遗留文件存在 sitemap.json | Rejected: 不检查类型，仅检查文件 |
| Skill description 标注 | VS Code activationEvents | 声明式，清晰 | LLM 可能忽略 description 中的约束 | Rejected: 隐式约束已被证明不可靠 |
| **运行时 surface 检测守卫** | gen-test-scripts 已有模式 | 显式、可验证、与现有模式一致 | 每次执行多一次 CLI 调用（< 50ms） | **Selected: 直接、一致、可靠** |

## Feasibility Assessment

### Technical Feasibility

`forge surfaces --json` 已是成熟的 CLI 命令，输出结构化 JSON。所有修改点均为 markdown 文件中的自然语言指令变更，不涉及代码逻辑。

### Resource & Timeline

7 个 markdown 文件的局部修改，工作量极小。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| Skill 名称中带 "web" 足以防止误用 | 5 Whys | Overturned: LLM 在 pipeline 中按流程执行，不一定校验 skill 名称与项目类型的匹配关系。追问后发现下游 skill 的无条件引用才是主要误用来源 |
| 下游 skill 只需检查 sitemap.json 是否存在即可 | Assumption Flip | Overturned: 非法遗留的 sitemap.json 会导致错误校验。应检查 surface 类型而非文件存在性 |
| ui-placement.md 的 load condition（依赖 artifact 存在性）足以守卫 | Stress Test | Overturned: TUI/Mobile 项目可能有 ui-design.md 或 prd-ui-functions.md，但仍不应访问 web 专属的 sitemap |

## Scope

### In Scope

- `plugins/forge/skills/gen-web-sitemap/SKILL.md`：新增 Step 0 Surface Check，无 web surface 时中止
- `plugins/forge/skills/write-prd/SKILL.md`（line 226）：Step 1 中读取 sitemap.json 前增加 surface 检查
- `plugins/forge/skills/write-prd/rules/ui-functions.md`（line 7）：读取 sitemap 前增加 surface 检查
- `plugins/forge/skills/write-prd/rules/self-check.md`（lines 19-20）：Placement consistency 和 Sitemap availability 检查增加 surface 前置条件
- `plugins/forge/skills/breakdown-tasks/rules/ui-placement.md`（lines 26-28）：route 校验前增加 surface 检查
- `plugins/forge/skills/eval/rules/validate-ux-pipeline.md`（line 33）：Web 行的 sitemap 引用增加显式 surface 守卫
- `plugins/forge/skills/gen-test-scripts/types/ui.md`：审查已有守卫，确认是否需要加固

### Out of Scope

- 新增 `forge surfaces` 命令或修改其输出格式
- `sitemap.json` schema 变更
- `hooks/guide.md` 中 sitemap 的文档描述更新
- 其他 web 专属 skill 的 surface 守卫（可作为后续 follow-up）
- 非 web surface 的布局结构化数据源（TUI screen layout、Mobile view hierarchy 等）——非 web 项目跳过 sitemap 后的 placement 信息空白是已知 gap，通过 follow-up proposal 解决

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| LLM agent 忽略守卫指令仍执行 sitemap 操作 | L | M | 守卫指令放在执行步骤的最早位置（Step 0 / 读取前），配合 STOP 关键词增强可见性 |
| `forge surfaces --json` 在未配置 surface 的项目中返回空导致误判 | M | L | 空结果视为无 web surface，中止/跳过是安全的（无 web surface 时 sitemap 无意义） |
| 下游 skill 在非 web 项目中跳过 sitemap 步骤后，丢失必要的 placement 校验 | L | M | 非 web surface 的 placement 校验不应依赖 sitemap（web 专属数据），跳过是正确行为 |

## Success Criteria

- [ ] `gen-web-sitemap` 在无 `web` surface 的项目中执行时，于 Step 0 中止并输出明确提示
- [ ] `gen-web-sitemap` 在有 `web` surface 的项目中执行时，Step 0 通过，正常进入后续步骤
- [ ] `write-prd` 在无 `web` surface 的项目中，跳过 sitemap.json 读取（Step 1、ui-functions、self-check）
- [ ] `breakdown-tasks` 在无 `web` surface 的项目中，跳过 sitemap route 校验
- [ ] `gen-test-scripts/types/ui.md` 已有守卫经审查确认充分，无需额外修改（或有修改点并记录）
- [ ] Monorepo 项目（多 surface 含 web）中所有 sitemap 相关功能正常工作

## Follow-up: Surface Layout Generalization

本 proposal 仅做 guard（非 web 项目跳过 sitemap 步骤），不补偿跳过后的 placement 信息空白。这是一个已知 gap：

- **TUI**：有界面布局（screen、panel、widget tree）但无结构化数据源，agent 只能从代码和设计稿推断
- **Mobile**：有 view hierarchy 和导航流但无对应抓取工具（web 有 agent-browser，TUI/Mobile 无等价物）
- **API**：无 UI 布局需求，不涉及

后续 proposal 应将 sitemap 概念泛化为 `surface-layout`：按 surface type 选择对应数据源（web → 无障碍树、TUI → 终端布局快照、Mobile → view hierarchy），下游 skill（`ui-functions.md`、`ui-placement.md`）根据 surface type 读取对应数据，而非简单跳过。

## Next Steps

- Proceed to `/quick` for streamlined implementation (task type: `coding.enhancement`)
