# Forge Skill Gap Analysis: write-prd & ui-design

> 从 train-recorder 项目暴露的原型质量问题，反推 forge skill 的具体改进点。

## 问题→根因→Skill 缺陷映射

| 问题 | 根因 | 出在哪个 Skill 的哪个步骤 |
|------|------|--------------------------|
| 6 个页面 Tab 栏不一致 | 没有导航契约，各页面独立生成 | **ui-design** prototype.md 没有「共享导航」指令 |
| app.js 覆盖页面函数 | app.js 含页面特定逻辑 | **ui-design** prototype.md app.js 描述边界模糊 |
| PRD 缺少统计页/动作详情页 | 原型演变后未回写 PRD | **ui-design** 没有对齐步骤 |
| PRD 没有 Page Composition 导航列 | 模板没有导航架构节 | **write-prd** 模板缺 Navigation Architecture |
| 链接指向不存在的 body.html | 生成时无交叉验证 | **ui-design** 没有集成验证步骤 |

## 核心洞察

**问题的本质是"页面缝隙"问题，不是"页面内部"问题。**

write-prd 按 UF（页面）逐个定义需求，ui-design 按 UF 逐个生成设计。两者都没有一个步骤专门处理**页面之间的关系**——导航、共享代码边界、链接一致性。

这就像盖楼：每间房都设计得很好，但走廊、楼梯、门牌号没人管。

---

## 改进方案 1: write-prd 模板增加 Navigation Architecture

**文件**: `skills/write-prd/templates/prd-ui-functions.md`

**现状**: 模板只有 `## UI Scope`（一句话列表）和重复的 `## UI Function N` 块。没有地方定义页面之间的关系。

**改进**: 在 `## UI Scope` 之后、第一个 UF 之前，新增 `## Navigation Architecture` 节：

```markdown
## Navigation Architecture

<!-- 在定义任何 UF 之前，先确定所有页面的关系和导航结构 -->

### Tab Navigation (全局底部导航栏)

| 序号 | 标签 | 目标页面 | 图标关键词 |
|------|------|---------|-----------|
| 1    |      |         |           |
| 2    |      |         |           |
| ...  |      |         |           |

### Push Pages (二级页面，有返回按钮)

| 页面 | 入口 | 返回目标 |
|------|------|---------|
|      |      |         |

### Navigation Rules

- Tab 页面共享底部导航栏，当前页面高亮
- Push 页面左上角返回按钮指向入口页面
- 不可在 Tab 栏中链接到不存在的页面
- 所有页面之间的跳转必须在此表中定义
```

**为什么有效**:
- write-prd 阶段就确定导航结构，而非留到 ui-design 时各页面自行发挥
- 这个表格成为 ui-design 的消费契约——生成 Tab 栏时直接从这个表复制
- 消除了"各页面独立想象导航"的根因

---

## 改进方案 2: ui-design SKILL.md 增加导航提取步骤

**文件**: `skills/ui-design/SKILL.md`

**现状**: Step 2 "Read UI Functions" 之后直接进入 Step 3 "Select Design Style"。没有步骤专门提取和确认导航架构。

**改进**: 在 Step 2 和 Step 3 之间插入 Step 2.5:

```markdown
## Step 2.5: Extract Navigation Architecture

Read the `## Navigation Architecture` section from `prd-ui-functions.md`.

If it exists:
- Use it as the single source of truth for all inter-page navigation
- Every tab page must use the exact same tab bar defined in this section
- Every push page must use the entry/return targets defined here

If it does NOT exist:
- Prompt the user: "PRD 缺少 Navigation Architecture 节。建议先运行 /write-prd 补充，或在此定义："
- Use AskUserQuestion to collect: tab pages, push pages, navigation rules
- Write the collected structure back into prd-ui-functions.md (append at UI Scope section)
```

**为什么有效**:
- 明确了导航的单一数据源
- 即使 write-prd 没有生成导航节，ui-design 也能在此补齐
- 后续所有页面生成步骤都引用这个提取结果

---

## 改进方案 3: ui-design prototype.md 增加共享导航规则

**文件**: `skills/ui-design/templates/prototype.md`

### 改进 3a: 新增「Navigation Contract」节

在 `## File Structure` 之前新增：

```markdown
## Navigation Contract

Before generating any page, define the navigation structure from the PRD's Navigation Architecture:

### Tab Bar (shared HTML snippet)

Extract the tab bar from the PRD's Navigation Architecture. Generate the HTML ONCE, then
copy it identically into every tab page. The only difference between pages is which tab
has the `active` class.

Rules:
- Tab bar HTML must be byte-identical across all tab pages (except active class position)
- Tab labels must exactly match the PRD's Navigation Architecture table
- Tab href targets must exactly match the filenames in this prototype
- NEVER generate different tab bars for different pages

### Code Layer Separation

| Layer | File | Contents | Prohibited |
|-------|------|----------|------------|
| Shared | app.js | Nav highlight, toast, modal, dropdown, accordion, slider, generic form validation | Any function that only one page uses |
| Page-specific | inline `<script>` in each HTML | Page data, page-specific event handlers, page-specific DOM manipulation | Re-implementing shared behaviors |

Rules:
- If a function is used by only ONE page, it goes in that page's inline script, NOT in app.js
- Shared functions must not accept page-specific parameters that only one page provides
- When in doubt, put it in the page's inline script
```

### 改进 3b: 修改 app.js 描述

**现状**:
```js
// Mobile nav toggle
// Modal open/close (backdrop click + Escape)
// Dropdown toggle (click-outside to close)
// Tab switching
// Toast notifications (auto-dismiss 3s)
// Form validation with inline errors
// Active page highlight in navigation
```

**改进**: 明确增加「禁止」列表：

```js
// Mobile nav toggle
// Modal open/close (backdrop click + Escape)
// Dropdown toggle (click-outside to close)
// Tab/segmented switching (generic, not page-specific)
// Toast notifications (auto-dismiss 3s)
// Form validation with inline errors
// Active page highlight in navigation
//
// PROHIBITED in app.js:
// - Page-specific data objects (e.g., dayData, exerciseData)
// - Page-specific event handlers (e.g., selectDay, toggleEditMode)
// - Any function referenced by onclick in only one HTML file
// - DOM queries for elements that only exist on one page
```

### 改进 3c: 增加集成验证清单

在 `## Output` 之前新增：

```markdown
## Post-Generation Verification

After generating ALL files, perform these checks:

### Navigation Consistency
- [ ] All tab pages have identical tab bar HTML (except active class)
- [ ] All tab bar href targets correspond to existing files
- [ ] All push pages have a back button linking to the correct parent page
- [ ] No href points to a non-existent file

### Code Separation
- [ ] app.js does NOT contain any function that is only called from one page
- [ ] No function name collision between app.js and any inline script
- [ ] Inline scripts load AFTER app.js (or use unique function names)

### Cross-Page Links
- [ ] Every `<a href="...">` in every file points to an existing file
- [ ] Every `onclick="location.href='...'"` points to an existing file
- [ ] index.html links to all generated pages
```

---

## 改进方案 4: ui-design SKILL.md 增加 PRD 对齐步骤

**文件**: `skills/ui-design/SKILL.md`

**现状**: Step 9 生成原型后直接结束。如果原型与 PRD 有差异（新增页面、交互改变），无人知道。

**改进**: 在 Step 9 之后增加 Step 10:

```markdown
## Step 10: Reconcile PRD (Optional)

After prototype generation, compare the generated prototype against prd-ui-functions.md.

Check for discrepancies:
- Pages in prototype but NOT in PRD → suggest adding a new UF
- Interactions in prototype but different from PRD → suggest updating UF interaction flow
- Data fields in prototype but not in PRD → suggest updating UF data requirements
- Navigation structure in prototype but not in PRD → suggest adding Navigation Architecture

If discrepancies found:
- List them with specific file references
- Ask user: "原型与 PRD 存在 N 处差异，是否回写 prd-ui-functions.md？"
- If yes → update prd-ui-functions.md to match the prototype (prototype is the source of truth)
- If no → note discrepancies in ui-design.md for future reference
```

**为什么有效**:
- 原型是需求的最终表达，是用户看到和确认的东西
- 不对齐的 PRD 比没有 PRD 更危险（给人错误预期）
- 这个步骤使 PRD 始终反映真实的 UI 状态

---

## 优先级排序

| 优先级 | 改进 | 影响范围 | 实现难度 |
|--------|------|---------|---------|
| **P0** | 3a: prototype.md 增加 Navigation Contract | 消除 Tab 栏不一致的根因 | 低（加一段文档） |
| **P0** | 3b: prototype.md app.js 增加禁止列表 | 消除函数冲突的根因 | 低（加几行注释） |
| **P0** | 3c: prototype.md 增加验证清单 | 消除断链的根因 | 低（加一个 checklist） |
| **P1** | 1: write-prd 模板增加 Navigation Architecture | 从源头定义导航契约 | 低（加一个模板节） |
| **P1** | 2: ui-design 增加导航提取步骤 | 确保消费契约 | 中（改 SKILL.md 流程） |
| **P2** | 4: ui-design 增加 PRD 对齐步骤 | 保持 PRD 与原型同步 | 中（加一个步骤） |

P0 三项都是改 prototype.md 模板（同一个文件），可以一次 PR 完成。
P1 两项分别改 write-prd 和 ui-design 的模板/SKILL.md。
P2 是锦上添花。

## 总结

一句话：**prototype.md 是杠杆率最高的改进点。** 它直接控制原型生成行为，3 个 P0 改进都是改这一个文件，就能消除我们遇到的全部 3 类问题。
