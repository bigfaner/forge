---
name: extract-design-md
description: Extract visual style from a web, mobile, or TUI application and generate a DESIGN.md for use with ui-design skill. Supports --platform flag (web, mobile, tui).
allowed-tools: Bash Read Write WebFetch
argument-hint: "[url] [--platform web|mobile|tui]"
---

# Extract Design MD

Auto-extract visual style from an application and generate a forge-compatible `DESIGN.md` for direct consumption by the `ui-design` skill.

**Core principle**: Observe a real product's visual language and distill it into a reusable design system specification.

## Process Flow

```
1. Parse platform flag → 2. Validate input → 3. Platform-specific extraction → 4. Match strategy → 5. Build design tokens → 6. Write DESIGN.md → 7. Confirm
```

## Platform Routing

Extract the `--platform` flag from command arguments. If not provided, default to `web`.

**Valid values**: `web`, `mobile`, `tui`

**Validation**: If `--platform` is provided with any other value, stop immediately and output:

> ERROR: unsupported platform "<value>". Must be one of: web, mobile, tui

Then route to the appropriate extraction section:

| Platform | Input | Extraction Method |
|----------|-------|-------------------|
| `web` (default) | URL | CSS extraction from HTML (Steps below) |
| `mobile` | URL | Mobile-adapted CSS extraction with mobile User-Agent + responsive analysis |
| `tui` | Local screenshot path | AI vision analysis — ANSI colors, character set, panel layout |

**Mobile extraction**: When `--platform mobile`, reuse the web extraction pipeline (Layers 1-5) with a mobile User-Agent viewport context, then add mobile-specific analysis:

1. **Fetch with mobile context**: When using WebFetch or agent-browser, set mobile viewport headers (viewport width: 375px, User-Agent: mobile) to trigger responsive CSS. Reuse all web extraction layers (Layer 1-5) unchanged — the same CSS bundle parsing, custom property extraction, multi-page sampling, and visual inference apply.

2. **Responsive breakpoint analysis**: Scan CSS for `@media` queries. Extract common mobile breakpoints:
   - 320px (small phone / iPhone SE)
   - 375px (standard phone / iPhone 12/13/14)
   - 414px (large phone / iPhone Plus/Pro Max)
   - 768px (tablet / iPad)
   Record which breakpoints the target site uses and what layout changes occur at each.

3. **Touch target estimation**: Analyze interactive elements (buttons, links, inputs) from CSS for minimum size compliance. Check `width`, `height`, `min-width`, `min-height`, `padding` on interactive selectors. Flag elements below the 44x44pt minimum touch target guideline. Values extracted from computed CSS; if not directly specified, mark as `(estimated)`.

4. **Safe area handling**: Check CSS for `env(safe-area-inset-*)` usage (notch/home indicator on iOS). Check HTML `<meta name="viewport">` for `viewport-fit=cover`. If neither is present, note that safe area handling was not detected and values are `(estimated)`.

> **Limitation**: Mobile extraction depends on the target URL serving responsive CSS. Sites without responsive stylesheets will produce web-equivalent results with mobile-specific sections marked `(estimated)`.

**TUI extraction**: When `--platform tui`, the input must be a **local file path** to a terminal screenshot (not a URL — screenshots cannot be fetched remotely). AI vision analyzes the screenshot to reverse-engineer design tokens, since TUI has no CSS or structured source to parse.

1. **Validate screenshot input**: The argument must be a local file path (e.g. `./screenshot.png`, `/tmp/terminal.png`). If a URL is provided instead, stop and output:

   > ERROR: TUI platform requires a local screenshot file path, not a URL. Provide a path like ./screenshot.png

   Use the `Read` tool to load the image file. If the file does not exist or cannot be read, stop with a clear error.

2. **Screenshot quality check**: Before detailed analysis, assess the screenshot quality. If the screenshot is blurry, low-resolution, or unreadable, stop and output:

   > ERROR: Screenshot quality is too low for reliable analysis. Please provide a clear, high-resolution terminal screenshot. Tips: use native screenshot tool (not photo of screen), ensure text is legible, capture at 1x scale.

3. **AI vision analysis**: Use the `Read` tool on the screenshot file to perform visual analysis. Extract the following categories, marking **ALL values as `(estimated)`** since AI vision inference is inherently approximate:

   - **ANSI color palette**: Identify the xterm-256 color numbers used in the screenshot. Map observed colors to the closest xterm-256 palette entries. Record background, text primary/secondary/tertiary, border, and semantic colors (success/error/warning/info). All color values must be xterm-256 numbers (0-255).
   - **Character set**: Determine whether the TUI uses box-drawing characters (┌─┐│└┘), block elements (█▄░▪), pure ASCII (+-\|*#), or a mix. Identify the specific characters used for borders, dividers, indicators, bar charts, and progress bars.
   - **Panel layout dimensions**: Estimate the number of rows and columns visible in the terminal. Identify panel boundaries and their dimensions (width in columns, height in rows). Note the overall terminal grid size.
   - **Key bindings**: If a status bar, help panel, or key binding legend is visible in the screenshot, extract the key-to-action mappings.

4. **Match strategy for TUI**: After extraction, use `AskUserQuestion` to let the user choose:

   | Option | Description |
   |--------|-------------|
   | Match closest built-in TUI theme, customize on top | Identify the closest built-in TUI theme (modern-dark-tui or minimal-ascii-tui), override differences with extracted tokens |
   | Fully custom from screenshot analysis | Generate an independent TUI DESIGN.md entirely from analysis results |

   **If "match built-in" is chosen:**

   Match against these characteristics to identify the closest built-in TUI theme:

   | Built-in Theme | Identifying Characteristics |
   |---------------|----------------------------|
   | modern-dark-tui | Dark background, 256-color (xterm-256), box-drawing + block elements, compact density |
   | minimal-ascii-tui | Default terminal background, 16-color (standard ANSI), pure ASCII characters, loose density |

   Read the corresponding built-in style file: `${CLAUDE_SKILL_DIR}/../ui-design/templates/styles/<name>.md`

5. **Build TUI design tokens and write DESIGN.md**: Read the template at `${CLAUDE_SKILL_DIR}/templates/design-tui.md`. Fill in results from analysis. All extracted values must be marked `(estimated)`.

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

Read the corresponding built-in style file: `${CLAUDE_SKILL_DIR}/../ui-design/templates/styles/<name>.md`

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

Select template based on platform:

| Platform | Template | Additional sections |
|----------|----------|-------------------|
| `web` | `${CLAUDE_SKILL_DIR}/templates/design-web.md` | — |
| `mobile` | `${CLAUDE_SKILL_DIR}/templates/design-web.md` | Append `${CLAUDE_SKILL_DIR}/templates/design-mobile.md` after "Responsive Behavior", before "Signature Patterns" |
| `tui` | `${CLAUDE_SKILL_DIR}/templates/design-tui.md` | — |

Write the design system to `DESIGN.md` in the project root.

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
