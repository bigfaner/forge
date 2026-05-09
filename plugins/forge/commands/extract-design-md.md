---
name: extract-design-md
description: Analyze a web app's visual style and generate a DESIGN.md for use with ui-design skill.
allowed_tools: ["Bash", "Read", "Write", "WebFetch"]
argument-hints:
  - name: url
    description: Web application URL to analyze (e.g. https://stripe.com)
    required: false
---

# /extract-design-md

Auto-extract visual style from a web application and generate a forge-compatible `DESIGN.md` for direct consumption by the `ui-design` skill.

**Core principle**: Observe a real product's visual language and distill it into a reusable design system specification.

## Process Flow

```
1. Get URL → 2. Analyze visual style → 3. Match strategy → 4. Build design tokens → 5. Write DESIGN.md → 6. Confirm
```

## Step 1: Get URL

Extract the target URL from command arguments or user message. If not provided, use `AskUserQuestion`:

> Please provide the web application URL to analyze (e.g. https://stripe.com)

Check if `DESIGN.md` already exists in the project root:

```bash
ls DESIGN.md
```

If it exists, use `AskUserQuestion` to ask whether to overwrite:

> DESIGN.md already exists in the project root. Overwrite?

- **Yes** → continue
- **No** → abort, inform user that the current `DESIGN.md` will be used automatically by the `ui-design` skill

## Step 2: Analyze Visual Style

Target dimensions:

| Dimension | What to Extract |
|-----------|----------------|
| Color palette | Background, primary/accent, text (primary/secondary/tertiary), border, semantic colors (success/warning/error) |
| Typography | Font family, weight scale, size scale (display/h1-h3/body/caption), line height, letter spacing |
| Components | Buttons (shape/color/hover), cards (background/border/radius/shadow), inputs (border/focus state), navigation (layout/active state) |
| Layout | Max content width, grid columns, spacing system, section padding |
| Depth & elevation | Shadow levels, blur values, opacity usage |
| Design philosophy | Overall style keywords (minimal/bold/elegant/playful/corporate, etc.) |

### Extraction Strategy (by priority)

SPAs (React/Vue, etc.) have styles not in HTML source. Extract layer by layer — **stop at the first successful layer, no need to proceed further**:

#### Layer 1: Trace CSS bundle

Use `WebFetch` to get page HTML, extract CSS file URLs from `<link rel="stylesheet">` (React build output is typically `/static/css/main.xxxxxx.css`), then fetch that CSS file.

The CSS bundle contains complete style rules — the most direct information source.

#### Layer 2: Extract CSS custom properties (design tokens)

Search for `:root` blocks in the CSS bundle. Modern design systems almost always use CSS variables for tokens:

```css
:root {
  --color-primary: #635bff;
  --font-size-base: 16px;
  --radius-md: 8px;
}
```

These variables map directly to DESIGN.md fields — the highest quality information source.

#### Layer 3: Multi-page sampling

A single page may only have landing content, missing form, card, table component styles. Fetch additional paths (if publicly accessible):

- `/login` — inputs, buttons
- `/dashboard` or `/app` — cards, navigation, data display
- `/settings` — forms, toggles, grouped sections

Compare multi-page extraction results to fill in missing component specs.

#### Layer 4: agent-browser runtime extraction (local apps)

If the target is a locally running SPA (e.g. `http://localhost:3000`), use agent-browser to execute JS for runtime computed styles:

```
ab('open <url>')
ab('wait --load networkidle')
// Extract CSS custom properties
ab('eval document.documentElement getComputedStyle --color-* --font-* --radius-* variables')
// Extract key component computed styles
ab('eval find button[class*="primary"] get backgroundColor borderRadius padding')
ab('eval find [class*="card"] get background border boxShadow borderRadius')
```

#### Layer 5: Visual inference (fallback)

When none of the above layers can obtain specific values, infer from page visual descriptions, screenshots, or HTML structure.

Mark uncertain values with `(estimated)` in the output.

## Step 3: Match Strategy

Use `AskUserQuestion` to let the user choose a generation strategy:

| Option | Description |
|--------|-------------|
| Match closest built-in style, customize on top | Identify the closest built-in style (vercel/shadcn/tailwind-ui/stripe/apple), override differences with extracted actual tokens |
| Fully custom from web app extraction | Generate an independent DESIGN.md entirely from analysis results, no built-in style reference |

**If "match built-in" is chosen:**

Based on Step 2 analysis, match against these characteristics to identify the closest built-in style:

| Built-in Style | Identifying Characteristics |
|---------------|----------------------------|
| Vercel | Black background, Geist font, no shadows, border depth |
| Shadcn | Zinc neutrals, CSS variables, dark mode support, Tailwind spacing |
| Tailwind UI | Indigo primary, white background, shadow-sm system, Inter font |
| Stripe | Purple gradient buttons, light gray background (#f6f9fc), weight-300 display |
| Apple | Pure white background, generous whitespace, SF Pro, rounded capsule buttons |

Read the corresponding built-in style file: `plugins/forge/skills/ui-design/templates/styles/<name>.md`

## Step 4: Build Design Tokens

**If "match built-in" is chosen:**

Use the built-in style file as a base, override differences with actual values from Step 2:
- Replace color values with extracted hex/rgb
- Replace font families with actually used fonts
- Replace border-radius, spacing, and other specific values
- Preserve the built-in style's Do's/Don'ts and Signature Patterns structure (can add site-specific patterns)
- Add a note at the top: `Based on: <built-in style name> (customized from <URL>)`

**If "fully custom" is chosen:**

Build all sections entirely from Step 2 analysis results, following the structure conventions of built-in style files.

## Step 5: Write DESIGN.md

Write the design system to `DESIGN.md` in the project root:

```markdown
# Design System: {{App Name or Domain}}

> Extracted from: {{URL}}
> Date: {{YYYY-MM-DD}}
> Based on: {{Built-in style name or "Custom"}}

## Visual Theme & Atmosphere

{{2-3 sentences describing overall visual style, atmosphere, and design philosophy}}

## Color Palette

| Role | Value | Usage |
|------|-------|-------|
| Background | #... | Page main background |
| Surface | #... | Card, panel background |
| Border | #... | Dividers, input borders |
| Text Primary | #... | Primary text |
| Text Secondary | #... | Secondary text, descriptions |
| Text Tertiary | #... | Placeholders, disabled text |
| Accent | #... | Primary interaction color, CTA buttons |
| Success | #... | Success state |
| Warning | #... | Warning state |
| Error | #... | Error state |

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

- **Primary**: {{background color, text color, border radius, padding, hover effect}}
- **Secondary**: {{style description}}
- **Ghost**: {{style description}}
- **Sizes**: sm / md / lg corresponding padding and font size

### Cards

- Background: {{value}}
- Border: {{value}}
- Border Radius: {{value}}
- Padding: {{value}}
- Shadow: {{value}}
- Hover: {{hover effect}}

### Inputs

- Background: {{value}}
- Border: {{value}}
- Border Radius: {{value}}
- Height: {{value}}
- Padding: {{value}}
- Focus: {{focus state description}}

### Navigation

- Layout: {{top nav / sidebar / other}}
- Background: {{value}}
- Active state: {{active item style}}
- Responsive: {{mobile behavior}}

## Layout

- Max content width: {{value}}
- Grid: {{columns}} columns, {{gap}} gap
- Section padding: {{desktop}} desktop / {{mobile}} mobile
- Component spacing: {{spacing system description}}

## Depth & Elevation

| Level | Shadow | Usage |
|-------|--------|-------|
| 0 | none | Flat elements |
| 1 | {{value}} | Cards, dropdowns |
| 2 | {{value}} | Modals, popovers |

## Do's and Don'ts

| Do | Don't |
|----|-------|
| {{correct practice}} | {{incorrect practice}} |

## Responsive Behavior

| Breakpoint | Behavior |
|-----------|---------|
| Mobile (<768px) | {{description}} |
| Tablet (768-1024px) | {{description}} |
| Desktop (>1024px) | {{description}} |

## Signature Patterns

{{2-5 signature visual patterns of this design system}}
```

## Step 6: Confirm & Next Step

Output completion message:

```
DESIGN.md written to project root
```

<EXTREMELY-IMPORTANT>
- Only write to `DESIGN.md` in the project root. Do not modify any other project files.
- If `DESIGN.md` already exists and the user declines overwrite, abort immediately without writing.
- All extracted values must come from actual CSS/HTML analysis. Mark estimated values with `(estimated)`.
</EXTREMELY-IMPORTANT>
