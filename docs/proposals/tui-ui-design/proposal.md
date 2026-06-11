---
created: 2026-05-14
author: faner
status: Draft
---

# Proposal: ui-design 支持 TUI 平台

## Problem

当前 `/ui-design` 只支持 web 和 mobile 两个平台。TUI（Terminal User Interface）feature 的视觉规格没有专门的 UI 设计环节，只能靠 `/tech-design` 阶段补充。

### Evidence

**deep-drill-analytics 经验**：该 TUI feature 在 pipeline 完成后产生了 16 个 fix/style 提交（占后功能提交的 70%），根因之一是 tech design 没有提供精确的视觉规格。5 个 style 提交是逐次试错收敛的字符选择，11 个 fix 提交是布局边界条件。

**lesson 总结**（`docs/lessons/lesson-tui-tech-design-mockup.md`）：提出了 TUI 视觉规格的 5 个结构性要求（ASCII mockup、尺寸数值、边界场景、字符调色板、颜色映射），但这些要求没有系统化地融入 pipeline。

### Urgency

每次 TUI feature 都依赖 `/tech-design` agent 的自觉性来补充视觉规格，不可靠。随着 TUI feature 增多，缺少系统化支持会导致每次都重复相同的试错过程。

## Proposed Solution

### 核心机制

1. **TUI 作为第三平台**：与 web/mobile 并列加入 `/ui-design`，ui-design 负责所有"用户看到什么、怎么交互"的决策
2. **TUI 专属模板**：ui-design.md 的 TUI 组件模板包含 lesson 的 5 个结构性要求，确保每个面板都有精确视觉规格
3. **HTML 模拟终端原型**：用 HTML + CSS 模拟终端外观（黑底等宽字体），供人类直观确认布局效果
4. **场景化评估 rubric**：eval-ui 按 platform 选择独立 rubric 文件，TUI 评估维度聚焦终端特有问题
5. **多平台并行支持**：同一 feature 声明多个 platform 时，单次运行按平台拆分输出文件

### Developer Workflow — Before & After

#### Before（当前行为）

1. `/write-prd` — PRD 声明 platform=tui，但 Navigation Architecture 表格无 TUI 适配
2. 跳过 `/ui-design`（只支持 web/mobile）
3. `/tech-design` — agent 需要自行补充 ASCII mockup、字符选择、尺寸规格，质量依赖 prompt 理解
4. 实现 → 大量 style/fix 提交试错收敛

#### After（本提案实施后）

1. `/write-prd` — PRD 声明 platform=tui，Navigation Architecture 自动渲染 TUI 专用结构（Keymap + Panel Layout + Mode）
2. `/ui-design` — 选择 TUI 主题（Modern Dark / Minimal ASCII / DESIGN.md 自定义），输出包含 5 个结构性要求的完整 TUI 视觉规格
3. 原型 — 生成 HTML 模拟终端，人类在浏览器中确认布局后批准
4. `/eval-ui` — 使用 TUI 专属 rubric 评估，确保 ASCII mockup 完整、字符调色板精确、边界场景覆盖
5. `/tech-design` — 直接消费精确视觉规格，无需猜测

### 设计决策

#### D1: TUI 视觉规格归属 — 放在 `/ui-design`

**选择**: TUI 作为第三个平台加入 `/ui-design`

**理由**: ui-design 负责所有"用户看到什么、怎么交互"的决策。TUI 的 ASCII mockup、字符调色板、颜色映射本质上是 UI 设计决策（视觉元素选择），不是技术实现细节。tech-design 只需要关心"用什么 TUI 框架实现"。

**否决方案**: 将 TUI mockup 保留在 `/tech-design`。这是 lesson 的原始建议，但当时没有 TUI 的 ui-design 支持。

#### D2: TUI 原型格式 — HTML 模拟终端

**选择**: HTML + CSS 模拟终端外观，浏览器查看

**理由**: 复用现有 HTML prototype 基础设施。ASCII mockup 在 ui-design.md 中供 agent 精确读取，HTML 原型供人类直观确认，两者互补。模拟的"失真"由精确数值规格兜底。

**否决方案**: Go/Python 可运行脚本（运行时依赖）、ANSI 文本文件（无导航）、纯 Markdown（不够直观）。

#### D3: TUI 设计风格 — 2 个内置主题 + DESIGN.md 自定义

| 主题 | 颜色空间 | 字符集 | 配色 | 密度 | 适用场景 |
|------|---------|--------|------|------|---------|
| Modern Dark | 256-color | box-drawing + block elements（▄▪─│┃） | 深色底，高对比，绿/红/蓝语义色 | 紧凑 | 大多数 CLI 工具 |
| Minimal ASCII | 16-color | 纯 ASCII（`#=-\|*+.`） | 默认终端底色，靠粗细和间距区分 | 宽松 | 极简工具、最大兼容性 |

#### D4: PRD 导航模板 — 按 platform 条件渲染

TUI Navigation Architecture 使用独立结构：

```
Platform: tui
Keymap: [Key | Action | Context/Mode]
Panel Layout: [Panel | View | Position | Size Hint]
Modes: [Mode | Description | Default Keybindings]
Navigation Rules: [...]
```

#### D5: ui-design.md TUI 组件模板

```
### Panel: {PanelName}
#### Panel Placement (View + Position)
#### ASCII Layout Mockup (box-drawing 图示)
#### Dimensions (具体数值，无模糊描述)
#### Character Palette (每元素指定 Unicode + 选择理由)
#### Color Mapping (前景/背景色号)
#### Edge Cases (5 个强制场景)
#### States
#### Key Bindings
#### Data Binding
```

**lesson 要素映射**: ASCII Mockup → 面板布局结构, Dimensions → 尺寸与比例, Character Palette → 字符调色板, Color Mapping → 颜色映射, Edge Cases → 边界条件场景。

#### D6: HTML 原型交互 — 单页面多面板 + 模拟按键

单个 `index.html`，所有 panel 渲染在一个黑色"终端窗口" div 中，底部模拟按键按钮（`[Tab]` `[1]` `[q]` `[:command]`）切换面板。

#### D7: 多平台 feature — 单次运行，按平台拆文件

```
docs/features/{slug}/ui/
  ui-design-web.md
  ui-design-tui.md
  prototype/web/
  prototype/tui/
```

#### D8: eval-ui 按平台选择独立 rubric

```
eval-ui/templates/rubric-web.md       # 原 rubric.md 重命名
eval-ui/templates/rubric-mobile.md    # 新增
eval-ui/templates/rubric-tui.md       # 新增
```

#### D9: TUI rubric 评估维度

| 维度 | 分值 | 说明 |
|------|------|------|
| Requirement Coverage | 250 | UI function 覆盖、导航覆盖、状态覆盖、边界处理 |
| Terminal Experience | 250 | 键盘可达性、信息扫描效率、面板切换直觉性、模式切换一致性 |
| Visual Specification | 250 | ASCII mockup 完整性、字符调色板精确性、颜色映射合规性、尺寸数值具体性 |
| Implementability | 250 | 尺寸有具体数值、5 个强制边界场景覆盖、字符 Unicode 编码明确 |

扣分规则：缺少 ASCII mockup → 该 panel Visual Specification = 0；字符有"待定" → -30/处；缺少强制边界场景 → -50/个；尺寸模糊描述 → -20/处。

#### D10: Mobile rubric 评估维度

| 维度 | 分值 | 说明 |
|------|------|------|
| Requirement Coverage | 250 | 含横竖屏、前后台、推送中断、弱网/离线边界场景 |
| Touch Experience | 250 | 触控目标尺寸（≥44pt）、拇指可达性、手势直觉性、平台规范遵循度 |
| Adaptive Layout | 250 | 多屏幕尺寸适配、横竖屏布局、安全区域、平台原生感 |
| Implementability | 250 | 自适应断点明确、平台交互模式指定、导航模式明确 |

扣分规则：触控目标未标注尺寸 → -30/处；缺少横竖屏适配 → -50；缺少安全区域处理 → -40。

## Requirements Analysis

### Key Scenarios

1. **纯 TUI feature**: PRD platform=tui → ui-design 使用 TUI 模板和主题 → 生成 HTML 模拟终端原型 → eval-ui 使用 TUI rubric
2. **纯 web feature**: 行为不变，eval-ui 使用 web rubric（原 rubric.md 重命名）
3. **纯 mobile feature**: 行为不变，eval-ui 使用 mobile rubric（新增）
4. **多平台 feature（web + TUI）**: ui-design 单次运行，输出 `ui-design-web.md` + `ui-design-tui.md`，各自独立原型和评估
5. **TUI theme 自定义**: 用户提供 `DESIGN.md` → ui-design 跳过内置主题选择，使用用户定义的风格
6. **Feature 无 UI**: platform 声明为 api/cli → 不触发 ui-design，行为不变
7. **原型审批拒绝**: 人类在 browser 中查看 TUI 原型 → 发现面板布局不合理 → 修改 ui-design.md 的 ASCII mockup → 重新生成原型
8. **eval-ui TUI rubric 不达标**: Visual Specification 得分低（缺少字符 Unicode 或尺寸模糊） → reviser 补充精确规格 → 重新评分

### Constraints & Dependencies

- 依赖 PRD 的 platform 字段（由 `/write-prd` 产出）
- TUI 原型是 HTML 模拟，不等同于真实终端渲染，最终效果需在终端验证
- lesson 的 5 个结构性要求必须通过模板强制（非可选），否则会退化到当前状态
- 移动端的 rubric 独立创建，不与 web 共用，避免评估标准混淆

### Non-Functional Requirements

| NFR Category | Requirement | Verification Method |
|-------------|-------------|-------------------|
| **兼容性** | 无 platform 字段或 platform=web 时行为完全不变 | 回归测试 |
| **扩展性** | 新增 platform 只需加 platforms/{name}.md + rubric-{name}.md + styles/*-tui.md，不改 SKILL.md 核心逻辑 | 未来验证 |
| **质量** | TUI ui-design.md 的 5 个结构性要求覆盖率 100%（每个 panel 必须包含全部 5 项） | eval-ui rubric 自动检查 |
| **性能** | 多平台单次运行不显著增加耗时（共享 PRD 读取） | 对比单平台/多平台执行时间 |

## Alternatives & Industry Benchmarking

### Alternatives Analysis

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| **TUI mockup 保留在 tech-design** | 零改动 | lesson 已证明不可靠，agent 自觉性不可控 | Rejected |
| **可运行 Go/Python TUI 原型** | 最真实终端体验 | 运行时依赖，agent 生成复杂 TUI 代码容易出错 | Rejected |
| **单一 rubric + 条件评分项** | 文件少 | rubric 膨胀，web/TUI 评分维度差异大难以统一 | Rejected |
| **HTML 模拟终端（本提案）** | 复用现有基础设施，浏览器可查看 | 非真实终端渲染 | **Selected** |
| **按平台拆 rubric（本提案）** | 各平台评估维度独立精确 | 3 个 rubric 文件需维护 | **Selected** |

## Feasibility Assessment

### Technical Feasibility

所有改动在 skill 文件层面（SKILL.md、模板文件），不涉及编译代码。ui-design 和 eval-ui 已有成熟的 platform 分支模式（web/mobile），TUI 遵循相同模式。

### Resource & Timeline

| Step | Task | 预计时间 |
|------|------|---------|
| 1 | 创建 `platforms/tui.md` + 2 个 TUI 主题文件 | 1.5h |
| 2 | 修改 `ui-design.md` 模板 + `SKILL.md` 增加 TUI 分支 | 2h |
| 3 | 修改 `prototype.md` 增加 HTML 模拟终端规则 | 1h |
| 4 | 修改 `prd-ui-functions.md` TUI 导航 + `write-prd/SKILL.md` | 1h |
| 5 | 创建 3 个 rubric 文件 + 修改 `eval-ui/SKILL.md` | 1.5h |
| 6 | 修改 `manifest-update-ui.md` 多平台支持 | 0.5h |
| 7 | 端到端验证 | 1.5h |
| **Total** | | **~9h** |

## Scope

### In Scope

- `ui-design` SKILL.md 增加 platform 判断和 TUI 分支逻辑
- `ui-design.md` 模板增加 TUI 组件模板（Panel 结构含 lesson 5 要素）
- `platforms/tui.md` TUI 平台导航规则
- 2 个 TUI 主题文件（Modern Dark、Minimal ASCII）
- `prototype.md` 增加 TUI HTML 模拟终端规则
- `prd-ui-functions.md` Navigation Architecture TUI 分支
- `eval-ui` 3 个独立 rubric（web、mobile、tui）
- `eval-ui` SKILL.md rubric 路径选择逻辑
- `manifest-update-ui.md` 多平台文件输出

### Out of Scope

- `tech-design` 改动（现有 TUI mockup 要求保留作为 fallback）
- `eval-design` 改动
- `breakdown-tasks` 改动
- `gen-test-cases` / `gen-test-scripts` TUI 类型化验证（由 typed-verification-strategies proposal 覆盖）
- 真实终端渲染验证基础设施
- 新增 profile

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| HTML 模拟终端与真实终端渲染差异 | Medium | Medium — 原型审批通过但实际效果不同 | ASCII mockup 的精确数值规格作为 agent 实现依据（非原型）；原型仅用于人类确认布局关系 |
| TUI 主题定义不够灵活 | Low | Low — 无法覆盖所有 TUI 风格偏好 | DESIGN.md 自定义入口覆盖 |
| Mobile rubric 与 web rubric 漂移 | Low | Low — 评估标准不一致 | 共享基础扣分规则，仅维度定义不同 |
| 多平台 feature 单次运行 token 开销增加 | Medium | Low | 共享 PRD 读取，各平台模板独立 |

## Success Criteria

- [ ] TUI feature 走完 `/write-prd` → `/ui-design` → `/eval-ui` 全流程，产出包含 5 个结构性要求的 `ui-design-tui.md`
- [ ] HTML 模拟终端原型可在浏览器中正常显示，面板布局与 ASCII mockup 一致
- [ ] eval-ui 使用 TUI rubric 评估 TUI 设计文档，Visual Specification 维度正确检查字符 Unicode、尺寸数值、边界场景
- [ ] eval-ui 使用 mobile rubric 评估 mobile 设计文档，Touch Experience 和 Adaptive Layout 维度生效
- [ ] 多平台 feature（web + tui）单次运行产出 2 个独立 ui-design 文件和 2 套原型
- [ ] 现有 web-only feature 的 ui-design 和 eval-ui 行为完全不变（回归验证）

## Next Steps

- Proceed to `/write-prd` to formalize requirements
