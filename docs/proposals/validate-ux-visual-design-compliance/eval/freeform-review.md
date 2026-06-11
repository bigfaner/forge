# Freeform Review: validate-ux Visual Design Compliance

**Reviewer**: Web Visual Compliance & Accessibility Architecture Expert
**Date**: 2026-06-08
**Document**: `/docs/proposals/validate-ux-visual-design-compliance/proposal.md`

---

## Section 1: Background Assessment

### Problem as I Understand It

The proposal identifies a blind spot in forge's validate-ux pipeline: it verifies code-level structural compliance (component names, class names in source code) but never validates whether the **rendered output** actually matches the design specification in `ui-design.md`. The evidence is concrete -- a pm-work-tracker feature passed all pipeline checks yet had 4 visual defects caught only by manual QA, including mismatched filter components across pages, incorrect button placement, and a broken calendar popup. The core claim is that "all 7 UI functions conform to design spec" was a false positive because no rendering validation existed.

### Core Technical Approach

The proposal introduces a "Design Compliance" collection phase inside the existing validate-ux pipeline. For each route in `sitemap.json`, it would: (1) navigate to the route via agent-browser, (2) capture a screenshot and extract the accessibility tree, (3) compare the accessibility tree against the structured requirements in `ui-design.md`, and (4) produce a pass/fail report per requirement. The key innovation is using **accessibility tree semantic matching** rather than pixel-level screenshot comparison -- avoiding baseline image maintenance and producing structured, explainable verdicts.

### Assumptions the Proposal Rests On

1. `ui-design.md` contains sufficiently structured and specific information to serve as a comparison baseline.
2. The browser accessibility tree captures enough semantic information to validate design compliance.
3. A mapping algorithm can reliably translate natural-language design requirements into accessibility tree node assertions.
4. The existing `ui-design.md` template format is standardized enough for automated parsing.

---

## Section 2: Key Risk Identification

### Risk 1: Accessibility Tree 表达力不足以覆盖核心视觉属性

**风险：** accessibility tree 的核心设计目标是辅助技术（screen reader），不是视觉回归检测。它天然缺失大量视觉属性。

提案原文：

> "将 accessibility tree 中的实际渲染结构与 ui-design.md 的要求逐条映射，标记 pass/fail"

以及：

> "能检测组件缺失、错误变体、布局偏差"

问题在于，accessibility tree 中 **不存在** 以下关键视觉属性：

- **颜色值**：accessibility tree 不包含 CSS color、background-color 信息。提案证据中的"列表页和详情页使用不同风格的筛选组件"如果是颜色差异，accessibility tree 完全无法检测。
- **字体大小/粗细**：font-size、font-weight 不在 accessibility tree 中（除非触发了 WCAG 对比度/字号违规）。
- **间距和定位**：margin、padding、position、flex/grid 布局信息完全缺失。提案证据中的"刷新按钮位置与设计稿不符"是 **布局错误**，但 accessibility tree 只能告诉按钮 **存在**（role=button），无法判断它在页面上的位置是否正确。
- **阴影、圆角、边框样式**：这些视觉属性同样不在 accessibility tree 中。

这意味着提案声称能检测的"布局偏差"在 accessibility tree 层面实际上只能做到"组件是否存在"，对于"组件在错误位置"这一核心场景几乎无能为力。

### Risk 2: ui-design.md 模板的结构化程度不足以支撑自动化逐条对比

**风险：** 当前 `ui-design.md` 模板大量使用自由文本和注释占位符，不是机器可解析的结构化格式。

查看实际的 `ui-design.md` 模板，`### Layout Structure` 部分是一个空注释 `<!-- Component hierarchy, grid/flex layout description -->`，这意味着每个项目的 layout 描述格式完全取决于生成时 LLM 的自由发挥。同样，`### States` 表格虽然有结构，但 `Visual` 列的内容是自由文本。

提案原文：

> "读取 ui-design.md：提取每个页面/路由的 UI 要求（组件类型、布局规范、样式属性、交互元素）"

问题：从自由文本中"提取"结构化要求本质上是一个 NLU 问题，不是确定性解析。不同项目、不同 LLM 生成轮次产生的 `ui-design.md` 格式会有显著差异。没有 schema 约束，"逐条映射"的精度无法保证。

### Risk 3: 语义匹配算法的鲁棒性未经讨论

**风险：** 提案完全没有讨论核心的匹配算法如何工作。

提案原文：

> "渲染结果 vs 设计文档逐条对比：将 accessibility tree 中的实际渲染结构与 ui-design.md 的要求逐条映射"

这个"逐条映射"是整个方案的核心，但提案没有回答以下关键问题：

- ui-design.md 说"筛选器组件"，accessibility tree 里是 `role="combobox"` + `role="listbox"` 的组合节点。这个映射规则从哪里来？是硬编码的规则引擎，还是依赖 LLM 做语义判断？
- 如果依赖 LLM 判断，那整个 pipeline 就变成了"LLM 读设计文档 → LLM 读 accessibility tree → LLM 判断是否匹配"，pass/fail 的可重复性和确定性都无法保证。
- 如果用规则引擎，那面对设计文档中的自由文本描述，规则的覆盖率会非常有限。

### Risk 4: 提案证据中的问题案例实际上暴露了方案能力边界

**问题：** 回看提案给出的 4 个实际问题案例，accessibility tree 方案可能只能检测其中 1 个。

| 问题案例 | Accessibility Tree 能检测？ |
|---------|--------------------------|
| 列表页和详情页使用不同风格的筛选组件（样式不一致） | 部分可能 -- 如果 DOM role 不同可以检测；如果只是 CSS 样式不同（同 role），无法检测 |
| 刷新按钮位置与设计稿不符（布局错误） | 大概率不能 -- accessibility tree 无位置信息 |
| 日历弹出框无法显示（交互组件失效） | 可能能检测 -- 如果触发后 accessibility tree 中没有 dialog/datepicker role 出现 |
| 视觉偏差（视觉回归未检测） | 取决于偏差类型，纯视觉偏差（颜色、间距、字号）无法检测 |

提案选取的 evidence 恰恰说明 accessibility tree 方案的覆盖范围有限，但提案没有诚实地量化这个覆盖率。

### Risk 5: pass/fail 判定的可解释性承诺可能与实际实现矛盾

**问题：** 提案要求每条 fail 判断包含"期望值（来自 ui-design.md）、实际值（来自 accessibility tree）"，但期望值来源是自由文本。

提案原文：

> "design compliance 报告中的每一条 fail 判断必须包含：失败原因、期望值（来自 ui-design.md）、实际值（来自 accessibility tree）、截图路径"

如果 ui-design.md 中写的是"表格页包含筛选区域，位于页面顶部"，"期望值"就是这段自然语言。而"实际值"是 accessibility tree 的 JSON 片段。这两者的格式差异使得"逐条对比"的输出不是结构化的 diff，而是一段自然语言解释。这降低了报告的可操作性 -- 用户仍然需要人工判断 LLM 的 pass/fail 判定是否合理。

### Risk 6: 跨页面一致性检测被划出范围，但这是证据中最常见的实际问题

**风险：** 提案 evidence 中的第一个问题是"列表页和详情页使用不同风格的筛选组件"，这本质上是跨页面一致性问题，不需要设计文档参照就能检测。

提案原文：

> "跨页面一致性审计（无设计文档参照的对比）——未来增强方向"

以及 out of scope 列表明确排除了"跨页面一致性审计"。但 gotcha evidence 中的第一个问题恰恰是跨页面一致性问题。如果 ui-design.md 中列表页和详情页的筛选器描述分别写在不同 section，且措辞不完全一致（这是大概率事件），逐条对比不会检测到这个问题。

把最容易自动化检测的问题类型排除在外，而主攻"是否符合设计文档"这个更难的问题，优先级值得商榷。

### Risk 7: performance 约束可能在实际项目中难以满足

**问题：** 提案要求"design compliance 采集不应使 validate-ux 的执行时间增加超过 50%"。

对于一个有 20 个路由的项目，每个路由需要：启动浏览器 tab → 等待页面加载 → 等待渲染稳定 → 截图 → 提取 accessibility tree → LLM 对比分析。即使并行处理，浏览器渲染和 LLM 调用都有不可压缩的延迟。如果现有 validate-ux 对 web 项目需要 3 分钟，50% 增量只有 1.5 分钟来处理所有路由的采集 + 分析，这在实际项目中可能不够。

---

## Section 3: Improvement Suggestions

### 建议 1: 引入 computed style 采集作为 accessibility tree 的补充

**针对风险 1。**

accessibility tree 不包含视觉属性，但 `window.getComputedStyle()` API 可以获取任意 DOM 节点的全部计算样式。具体改动：

在 agent-browser 采集步骤中，除了提取 accessibility tree，还对关键 DOM 节点（通过 accessibility tree 的 `computedRole` 和可见性筛选确定）执行 `getComputedStyle()`，采集一组预定义的视觉属性（color, background-color, font-size, font-weight, margin, padding, display, flex/grid 相关属性）。

这样 design compliance 报告可以对比：组件存在性（accessibility tree）、布局属性（computed style）、视觉样式（computed style），覆盖面从单纯的"语义存在"扩展到"视觉呈现"。

### 建议 2: 对 ui-design.md 模板增加 machine-readable 标注格式

**针对风险 2。**

在现有 ui-design.md 模板的自由文本区域增加可选的结构化标注块。例如：

```markdown
### Layout Structure
<!-- @design-spec
{
  "components": [
    {"name": "filter-bar", "role": "search", "required": true, "position": "top"},
    {"name": "data-table", "role": "table", "required": true, "position": "main"},
    {"name": "pagination", "role": "navigation", "required": true, "position": "bottom"}
  ]
}
-->
```

这样设计文档既保持人类可读性，又提供机器可解析的标注。design compliance 步骤优先解析标注块，不存在标注块时降级为自由文本 NLU 解析（并标注低置信度）。这解决了 ui-design.md 格式不确定的核心问题，同时通过 `@design-spec` 标注的存在性可以度量设计文档的"结构化覆盖度"。

### 建议 3: 将跨页面一致性检测提升为 P0 scope

**针对风险 6。**

跨页面一致性检测不依赖 ui-design.md 的质量，只需要对比同一项目内不同路由的 accessibility tree。具体做法：对所有路由的 accessibility tree 做结构聚类，自动识别 role/composition 异常的路由（例如两个页面都应该有筛选器，但一个用 combobox 一个用 select）。这比"与设计文档对比"更简单、更可靠、且直接命中 evidence 中的实际痛点。

建议将此功能从 out-of-scope 移到 in-scope 的第一优先级，将"与 ui-design.md 对比"作为第二优先级。先实现不依赖设计文档质量的自动化检测，再叠加设计文档参照。

### 建议 4: 明确匹配算法的实现策略 -- 规则优先 + LLM 兜底

**针对风险 3 和风险 5。**

将匹配过程分为两层：

1. **规则层（高置信度，确定性）**：预定义常见组件类型的 role 映射规则（table → role=table/grid/treegrid, date-picker → role=combobox/dialog + aria-haspopup, filter → role=search + input, pagination → role=navigation + aria-label）。规则命中的条目直接 pass/fail，无需 LLM 介入。
2. **LLM 兜底层（低置信度，标注后供人工审核）**：规则无法覆盖的设计要求（自由文本描述的布局关系、组件嵌套语义等）交给 LLM 判断，但在报告中标注为"LLM 判定"，置信度标记为 low。

这确保了常见组件类型（表格、表单、导航）的判定是确定性的和可解释的，同时不丢失对复杂场景的覆盖能力。

### 建议 5: 增加 design compliance 能力边界的显式声明

**针对风险 4。**

在提案中增加一个"Capability Boundary" section，诚实列出 design compliance **可以** 和 **不可以** 检测的问题类型。例如：

| 能检测 | 不能检测 |
|--------|---------|
| 组件存在/缺失 | 颜色渐变、阴影精确值 |
| 组件 role 类型正确性 | 间距的像素级精确度 |
| 交互元素可聚焦/可操作 | 字体粗细/大小细微差异 |
| 必填表单字段的 aria 属性 | 动画效果 |
| 页面结构层次（heading 层级） | 视觉对齐的亚像素偏差 |

这防止了用户对 design compliance 产生不切实际的期望，也使得 success criteria 的验证更公平。

### 建议 6: 将 per-route 采集结果缓存，支持增量验证

**针对风险 7。**

对于未变更的路由（通过 git diff 检测路由对应的前端代码是否变更），复用上一轮的采集结果，只重新采集和对比有变更的路由。这可以显著降低大项目的增量验证时间，使 50% 性能约束在实际项目中可达成。
