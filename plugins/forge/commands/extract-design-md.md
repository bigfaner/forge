---
name: extract-design-md
description: Analyze a web app's visual style and generate a DESIGN.md for use with ui-design skill.
argument-hints:
  - name: url
    description: 要分析的 web 应用 URL（如 https://stripe.com）
    required: false
---

# /extract-design-md

从 web 应用自动提取视觉风格，生成符合 forge 规范的 `DESIGN.md`，供 `ui-design` skill 直接消费。

**核心原则**：观察真实产品的视觉语言，提炼为可复用的设计系统规范。

## Process Flow

```
1. Get URL → 2. Analyze visual style → 3. Match strategy → 4. Build design tokens → 5. Write DESIGN.md → 6. Confirm
```

## Step 1: Get URL

从命令参数或用户消息中提取目标 URL。若未提供，使用 `AskUserQuestion` 询问：

> 请提供要分析的 web 应用 URL（例如 https://stripe.com）

检查项目根目录是否已存在 `DESIGN.md`：

```bash
ls DESIGN.md
```

若已存在，使用 `AskUserQuestion` 询问是否覆盖：

> 项目根目录已存在 `DESIGN.md`，是否覆盖？

- **Yes** → 继续
- **No** → 中止，提示用户当前 `DESIGN.md` 将被 `ui-design` skill 自动使用

## Step 2: Analyze Visual Style

目标维度：

| 维度 | 提取内容 |
|------|---------|
| Color palette | 背景色、主色/强调色、文字色（primary/secondary/tertiary）、边框色、语义色（success/warning/error） |
| Typography | 字体族（font-family）、字重层级、字号层级（display/h1-h3/body/caption）、行高、字间距 |
| Components | 按钮（形状/颜色/hover）、卡片（背景/边框/圆角/阴影）、输入框（边框/focus 状态）、导航（布局/active 状态） |
| Layout | 最大内容宽度、网格列数、间距系统、section padding |
| Depth & elevation | 阴影层级、模糊值、透明度用法 |
| Design philosophy | 整体风格关键词（minimal/bold/elegant/playful/corporate 等） |

### 提取策略（按优先级）

SPA（React/Vue 等）的样式不在 HTML 源码中，需按以下层次逐步提取，**每层成功则停止，不必执行后续层**：

#### 层次 1：追踪 CSS bundle

用 `WebFetch` 抓取页面 HTML，从 `<link rel="stylesheet">` 提取 CSS 文件 URL（React 构建产物通常为 `/static/css/main.xxxxxx.css`），再 fetch 该 CSS 文件。

CSS bundle 包含完整样式规则，是最直接的信息源。

#### 层次 2：提取 CSS 自定义属性（设计 token）

在 CSS bundle 中搜索 `:root` 块，现代设计系统几乎都用 CSS variables 存储 token：

```css
:root {
  --color-primary: #635bff;
  --font-size-base: 16px;
  --radius-md: 8px;
}
```

这些变量直接映射到 DESIGN.md 的各字段，是最高质量的信息源。

#### 层次 3：多页面采样

单页可能只有 landing，缺少 form、card、table 等组件样式。额外 fetch 以下路径（若可公开访问）：

- `/login` — 输入框、按钮
- `/dashboard` 或 `/app` — 卡片、导航、数据展示
- `/settings` — 表单、开关、分组

对比多页面提取结果，补全缺失的组件规范。

#### 层次 4：agent-browser 运行时提取（本地应用）

若目标为本地运行的 SPA（如 `http://localhost:3000`），用 agent-browser 执行 JS 提取运行时计算样式：

```
ab('open <url>')
ab('wait --load networkidle')
// 提取 CSS 自定义属性
ab('eval document.documentElement 获取 getComputedStyle 中所有 --color-* --font-* --radius-* 变量')
// 提取关键组件的计算样式
ab('eval 查找 button[class*="primary"] 并获取 backgroundColor borderRadius padding')
ab('eval 查找 [class*="card"] 并获取 background border boxShadow borderRadius')
```

#### 层次 5：视觉推断（兜底）

前四层均无法获取具体值时，从页面的视觉描述、截图或 HTML 结构中推断。

记录不确定的值，在输出中标注 `(estimated)`。

## Step 3: Match Strategy

使用 `AskUserQuestion` 让用户选择生成策略：

| 选项 | 说明 |
|------|------|
| 匹配最接近的内置风格，基于它生成 | 识别最接近的内置风格（vercel/shadcn/tailwind-ui/stripe/apple），用提取的实际 token 覆盖差异 |
| 完全从 web 应用提取，生成全新自定义风格 | 不参考内置风格，完全基于分析结果生成独立 DESIGN.md |

**若选"匹配内置"：**

根据 Step 2 的分析结果，对照以下特征识别最接近的内置风格：

| 内置风格 | 识别特征 |
|---------|---------|
| Vercel | 黑色背景、Geist 字体、无阴影、边框深度 |
| Shadcn | Zinc 中性色、CSS 变量、深色模式支持、Tailwind 间距 |
| Tailwind UI | Indigo 主色、白色背景、shadow-sm 系统、Inter 字体 |
| Stripe | 紫色渐变按钮、浅灰背景(#f6f9fc)、weight-300 display |
| Apple | 纯白背景、极大留白、SF Pro、圆角胶囊按钮 |

读取对应的内置 style 文件：`plugins/forge/skills/ui-design/templates/styles/<name>.md`

## Step 4: Build Design Tokens

**若选"匹配内置"：**

以内置 style 文件为基础，用 Step 2 提取的实际值覆盖差异项：
- 替换颜色值为实际提取的 hex/rgb
- 替换字体族为实际使用的字体
- 替换圆角、间距等具体数值
- 保留内置风格的 Do's/Don'ts 和 Signature Patterns 结构（可补充网站特有模式）
- 在文件顶部注明：`Based on: <内置风格名> (customized from <URL>)`

**若选"全新自定义"：**

完全基于 Step 2 的分析结果构建所有 section，遵循内置 style 文件的结构规范。

## Step 5: Write DESIGN.md

将生成的设计系统写入项目根目录 `DESIGN.md`，格式如下：

```markdown
# Design System: {{App Name or Domain}}

> Extracted from: {{URL}}
> Date: {{YYYY-MM-DD}}
> Based on: {{内置风格名 or "Custom"}}

## Visual Theme & Atmosphere

{{2-3 句描述整体视觉风格、氛围和设计哲学}}

## Color Palette

| Role | Value | Usage |
|------|-------|-------|
| Background | #... | 页面主背景 |
| Surface | #... | 卡片、面板背景 |
| Border | #... | 分割线、输入框边框 |
| Text Primary | #... | 主要文字 |
| Text Secondary | #... | 次要文字、说明文字 |
| Text Tertiary | #... | 占位符、禁用文字 |
| Accent | #... | 主要交互色、CTA 按钮 |
| Success | #... | 成功状态 |
| Warning | #... | 警告状态 |
| Error | #... | 错误状态 |

## Typography

| Role | Font | Weight | Size | Line Height |
|------|------|--------|------|-------------|
| Display | ... | ... | ...px | ... |
| H1 | ... | ... | ...px | ... |
| H2 | ... | ... | ...px | ... |
| H3 | ... | ... | ...px | ... |
| Body | ... | ... | ...px | ... |
| Caption | ... | ... | ...px | ... |
| Mono | ... | ... | ...px | ... |

## Components

### Buttons

- **Primary**: {{背景色、文字色、圆角、padding、hover 效果}}
- **Secondary**: {{样式描述}}
- **Ghost**: {{样式描述}}
- **Sizes**: sm / md / lg 对应的 padding 和字号

### Cards

- Background: {{值}}
- Border: {{值}}
- Border Radius: {{值}}
- Padding: {{值}}
- Shadow: {{值}}
- Hover: {{hover 效果}}

### Inputs

- Background: {{值}}
- Border: {{值}}
- Border Radius: {{值}}
- Height: {{值}}
- Padding: {{值}}
- Focus: {{focus 状态描述}}

### Navigation

- Layout: {{顶部导航 / 侧边栏 / 其他}}
- Background: {{值}}
- Active state: {{active 项样式}}
- Responsive: {{移动端行为}}

## Layout

- Max content width: {{值}}
- Grid: {{列数}} columns, {{gap}} gap
- Section padding: {{desktop}} desktop / {{mobile}} mobile
- Component spacing: {{间距系统描述}}

## Depth & Elevation

| Level | Shadow | Usage |
|-------|--------|-------|
| 0 | none | 平面元素 |
| 1 | {{值}} | 卡片、下拉菜单 |
| 2 | {{值}} | 模态框、弹出层 |

## Do's and Don'ts

| Do | Don't |
|----|-------|
| {{正确做法}} | {{错误做法}} |

## Responsive Behavior

| Breakpoint | Behavior |
|-----------|---------|
| Mobile (<768px) | {{描述}} |
| Tablet (768-1024px) | {{描述}} |
| Desktop (>1024px) | {{描述}} |

## Signature Patterns

{{该设计系统的标志性视觉模式，2-5 条}}
```

## Step 6: Confirm & Next Step

输出完成提示：

```
✓ DESIGN.md 已写入项目根目录
```