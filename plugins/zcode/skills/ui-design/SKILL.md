---
name: ui-design
description: Use after PRD ui-functions are defined to create UI design specifications. Supports design style selection and HTML prototype generation.
---

# UI Design

## Overview

从 PRD 的 UI 功能需求产出 UI 设计规格文档，并可选生成 HTML 原型。

**核心原则**：定义 HOW 界面呈现和交互，与 PRD 的 WHAT（需求层）分离。

<HARD-GATE>
Do NOT write any implementation code (except HTML prototypes). This skill produces design specification documents and optionally HTML/CSS/JS prototypes.
</HARD-GATE>

## Prerequisites

检查上一阶段产物，缺失则中止并提示用户：

```bash
ls docs/features/<slug>/prd/prd-ui-functions.md
```

| 产物                      | 缺失时提示                             |
| ------------------------- | -------------------------------------- |
| `prd/prd-ui-functions.md` | 先执行 `/write-prd` 并补充 UI 功能要点 |

## When to Use

**Trigger conditions:**

- PRD with `prd/prd-ui-functions.md` exists
- User asks to design the UI

**Skip when:**

- No UI functions defined (backend/API/CLI features)
- UI design already exists

## Process Flow

```
1. Read manifest → 2. Read UI functions → 3. Select design style → 4. Draft design → 5. Write design → 6. Update manifest → 7. Eval prompt → 8. Prototype prompt → 9. Generate prototype (optional)
```

## Step 1: Read Manifest

Read `manifest.md` to locate `prd/prd-ui-functions.md`.

## Step 2: Read UI Functions

Read `prd/prd-ui-functions.md` to understand UI requirements.

## Step 3: Select Design Style

按优先级选择设计风格：

### Priority 1: User-provided DESIGN.md

检查以下位置，找到则直接使用：

- 项目根目录 `DESIGN.md`
- Feature 目录 `docs/features/<slug>/ui/style.md`

如果用户指定了自定义路径，优先使用该路径。

### Priority 2: Built-in Styles

若用户未提供 DESIGN.md，让用户从 5 种内置风格中选择：

| Style           | Vibe                     | 适合场景                       |
| --------------- | ------------------------ | ------------------------------ |
| **Vercel**      | 黑白极简，开发者工具感   | 开发者平台、CLI 工具、技术文档 |
| **Shadcn**      | Zinc 中性色，功能极简    | SaaS、管理后台、工具类应用     |
| **Tailwind UI** | Indigo 主色，专业温暖    | 通用 SaaS、营销页、企业后台    |
| **Stripe**      | 紫色渐变，轻量优雅       | 金融科技、品牌官网、支付产品   |
| **Apple**       | 大留白，影像驱动，高级感 | 消费产品、品牌官网、营销页     |

使用 `AskUserQuestion` 工具让用户选择，选项附简要说明。

Built-in style files located at: `templates/styles/{vercel,shadcn,tailwind-ui,stripe,apple}.md`

### Priority 3: Clone More Styles from Repo

如果内置 5 种不满足，支持从 awesome-design-md 仓库获取更多：

```bash
# Clone to a temp directory (not into project)
git clone --depth 1 git@github.com:VoltAgent/awesome-design-md.git /tmp/awesome-design-md

# Then use npx to fetch a specific site's DESIGN.md:
npx getdesign@latest add <site-name>
```

> Note: The repo's DESIGN.md files are hosted externally at getdesign.md. The `npx getdesign@latest` CLI fetches them.

### Using the Selected Style

将选中的 style 内容作为 `Design System` 部分写入 `ui-design.md`，后续组件设计必须遵循该 style 的色彩、排版、组件规范。

## Step 4: Draft UI Design

For each UI function, define:

- Layout structure (component hierarchy)
- States (loading, empty, error, populated)
- Interactions (triggers, actions, feedback)
- Data binding (UI element → data field)

所有设计决策必须遵循 Step 3 中选定的 design style。

## Step 5: Write UI Design

Save to `ui/ui-design.md` using `templates/ui-design.md`.

Design System 部分引用 Step 3 选中的 style（内联核心 token，不引用外部文件）。

## Step 6: Update Manifest

Update `manifest.md`:

- Add UI Design row to Documents table
- Add traceability links from UI Functions to UI Design sections
- Advance status to `design` if `/design-tech` already completed

Use `templates/manifest-update-ui.md` for the update pattern.

## Step 7: Adversarial Eval Prompt

After updating manifest, use `AskUserQuestion` to ask:

> Run `/eval-ui` for adversarial evaluation of the UI design? (default 80/3)

- **Yes** → invoke `/eval-ui` via `Skill` tool
- **Custom** → invoke `/eval-ui --target X --iterations Y` via `Skill` tool
- **No** → proceed to prototype prompt

## Step 8: Prototype Prompt

Use `AskUserQuestion`:

> Generate HTML/CSS/JS interactive prototype from the UI design?

- **Yes** → proceed to Step 9
- **No** → done

## Step 9: Generate Prototype (Optional)

根据 `ui-design.md` 和选定的 design style 生成 HTML 原型：

- 遵循 `templates/prototype.md` 的规范
- 多文件结构：共享 CSS/JS + 每页独立 HTML + index.html 导航
- 实现所有组件的所有状态（loading/empty/error/populated）
- 交互行为可用（导航切换、模态框、标签页等）
- 响应式布局

Save to: `docs/features/<slug>/ui/prototype/`

## Integration

Works well with:

- `/write-prd` — Produces `prd/prd-ui-functions.md` input
- `/design-tech` — Parallel skill; both must complete before breakdown
- `/eval-ui` — Adversarial evaluation after UI design is created
