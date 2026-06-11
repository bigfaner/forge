# Lesson: TUI Feature 必须在 Tech Design 阶段提供 ASCII Layout Mockup

## 问题

deep-drill-analytics feature 在 Forge pipeline 完成后，产生了 **16 个 fix/style 提交**（占后功能提交的 70%）用于修 bug 和调样式。根因之一是 tech design 没有提供精确的视觉规格，agent 只能猜测布局细节。

### 具体表现

5 个 style 提交是逐次试错收敛的：
```
compact bar chart → narrower bars → half-height bars → underscore for unfilled
```
11 个 fix 提交是布局边界条件：
```
percentage wrapping → column misalignment → scrollbar width → content overflow → screen garbling
```

如果 tech design 提前定义了这些细节，大部分提交可以避免。

**Why:** TUI 没有可视化设计工具，agent 只能靠文字描述想象布局。没有精确规格时，agent 会选择"看起来合理"的默认值，然后用户在运行时逐个修正。每次修正都是一个 fix/style 提交。

## 规则

对于涉及 TUI 渲染的 feature，`/tech-design` 阶段**必须**包含 ASCII layout mockup，覆盖以下要素：

### 1. 面板布局结构

用 box-drawing 字符画出每个面板的精确结构：

```
┌─ Tool Statistics ──────────────────────┐
│                                        │
│  ▄▄▄▄▄▄▄▄▄___  67%  Read              │
│  ▄▄▄▄▄___         43%  Write           │
│  ▄▄▄____________  12%  Bash            │
│                                        │
├─ File Operations ─────────────────────┤
│                                        │
│  ▪▪▪▪▪  auth.go          R×12  E×3    │
│  ▪▪▪     handler.go       R×8   E×1   │
│  ▪▪      main.go          R×5          │
│                                        │
└────────────────────────────────────────┘
```

### 2. 尺寸与比例

明确写出数值，不要用"大约"、"适当"：

```
- 面板总宽: m.width - 2 (border)
- 内容区: m.width - 5 (border + scrollbar)
- 两列布局: (contentWidth - 3) / 2
- 柱状图宽度: barWidth = colWidth - labelWidth - 6
- 路径截断: maxLen = contentWidth - suffixWidth - indent
```

### 3. 边界条件场景

**必须覆盖**以下 5 个场景（不可省略）：

| # | Scenario | Expected |
|---|----------|----------|
| 1 | 窄终端 80×24 | 布局不溢出，柱状图缩至 4 格 |
| 2 | 宽终端 140+ col | 布局不变形，柱状图最大 15 格 |
| 3 | 混合数字宽度 (1 vs 100) | 右 pad 到最宽值，列对齐 |
| 4 | 路径超长 (>50 chars) | ".../parent/file.go" 保留尾部段 |
| 5 | 无数据 | 显示 "No data" 居中 |

feature 涉及 CJK 文本时，增加第 6 个场景：CJK 字符串宽度计算正确。

### 4. 字符调色板

**禁止**将视觉元素的字符选择留空或"待定"。必须明确指定每个视觉元素使用的 Unicode 字符：

```
| 元素 | 字符 | Unicode | 选择理由 |
|------|------|---------|----------|
| 柱状图填充 | ▄ | U+2584 | 半高，行间无空行时紧凑 |
| 柱状图空白 | _ | U+005F | 与 ▄ 等高，视觉连续 |
| 文件操作条 | ▪ | U+25AA | 小方块，不喧宾夺主 |
| 分隔线 | ─ | U+2500 | box-drawing 标准 |
| 滚动条轨道 | │ | U+2502 | 窄，不占空间 |
| 滚动条滑块 | ┃ | U+2503 | 粗于轨道，视觉突出 |
```

> **教训**: bar chart 字符经历了 3 次迭代（`▎` → `▓`/`░` → `▄`/`_`），每次都是 style 提交。在 tech design 阶段确定字符选择可避免此浪费。

### 5. 颜色映射

从 `tui-layout-ui.md` 颜色表中选取，不要引入新色号：

```
| 元素 | 字符 | 前景色 | 背景色 |
|------|------|--------|--------|
| 柱状图填充 | ▄ | 82 (green) | - |
| 柱状图空白 | _ | 239 (dim) | - |
| R 计数 | - | 82 (green) | - |
| E 计数 | - | 196 (red) | - |
| 光标行 | - | 15 (white) | 55 (purple) |
```

## How to apply

1. `/tech-design` skill 的 prompt 中加入 TUI mockup 要求（当 feature scope 包含 frontend/tui 时）
2. mockup 审查时，reviewer（人或 `/eval-design`）检查：
   - 每个面板是否有 ASCII 图示
   - 尺寸是否有具体数值（非模糊描述）
   - 是否覆盖 5 个强制边界场景
   - 字符调色板是否完整（每个视觉元素都指定了 Unicode 字符）
   - 颜色是否从现有色表选取
3. `/breakdown-tasks` 时，将 mockup 中的边界条件直接转化为 task 的 verify criteria

## 预期效果

| 指标 | 改进前 | 改进后（预期） |
|------|--------|---------------|
| 后功能 style 提交 | 5 | 0-1 |
| 后功能 fix 提交 | 11 | 2-3 |
| Vibe coding 迭代轮数 | ~16 | ~3 |

## Tags

`tech-design`, `tui`, `mockup`, `process`, `forge-improvement`
