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
| `web` (default) | URL | CSS extraction from HTML (Layers 1-5 in Step 2) |
| `mobile` | URL | Mobile-adapted CSS extraction (see `rules/platform-routing.md`) |
| `tui` | Local screenshot path | AI vision analysis (see `rules/platform-routing.md`) |

For **mobile** and **tui** platform-specific extraction details, follow the rules in `rules/platform-routing.md`.

### Mobile Extraction Summary

When `--platform mobile`, reuse the web extraction pipeline (Layers 1-5) with a mobile User-Agent and viewport context (viewport width: 375px), then add mobile-specific analysis:

- **Responsive breakpoint analysis**: Scan CSS for `@media` queries targeting common mobile breakpoints: 320px (small phone / iPhone SE), 375px (standard phone), 414px (large phone), 768px (tablet). Record which breakpoints the target site uses.
- **Touch target estimation**: Analyze interactive elements for minimum size compliance against the 44x44pt touch target guideline. Flag elements below minimum. Mark values as `(estimated)` if not directly specified.
- **Safe area handling**: Check CSS for `env(safe-area-inset-*)` usage and HTML `<meta name="viewport">` for `viewport-fit=cover`. Note safe-area results.
- **Responsive CSS limitation**: Mobile extraction depends on the target URL serving responsive CSS (responsive stylesheet with media queries). Sites without responsive stylesheets produce web-equivalent results with mobile sections marked `(estimated)`.

### TUI Extraction Summary

When `--platform tui`, input must be a **local file path** to a terminal screenshot (not a URL). AI vision analyzes the screenshot to reverse-engineer design tokens:

- **ANSI color palette**: Identify xterm-256 color numbers used. Map observed colors to closest xterm-256 palette entries (0-255). All values `(estimated)`.
- **Character set**: Determine whether the TUI uses box-drawing characters, block elements, pure ASCII, or a mix. Identify specific characters for borders, dividers, indicators, bar charts, and progress bars.
- **Panel layout dimensions**: Estimate terminal rows and columns. Identify panel boundaries and their dimensions.
- **Key bindings**: Extract key-to-action mappings from visible status bar or help panel.
- **Built-in TUI themes**: Match against `modern-dark-tui` (dark background, 256-color, box-drawing + block elements) or `minimal-ascii-tui` (default background, 16-color, pure ASCII).
- **Screenshot quality check**: Reject blurry, low-resolution, or unreadable screenshots.

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

SPAs (React/Vue, etc.) have styles not in HTML source. Extract layer by layer, stopping at the first successful layer. Follow the 5-layer strategy in `rules/extraction-layers.md`:

1. **Layer 1**: Trace CSS bundle from HTML source
2. **Layer 2**: Extract CSS custom properties (design tokens)
3. **Layer 3**: Multi-page sampling for missing component styles
4. **Layer 4**: agent-browser runtime extraction (local apps)
5. **Layer 5**: Visual inference (fallback, mark as `(estimated)`)

## Step 3: Match Strategy

### Web / Mobile

Use `AskUserQuestion` to let the user choose a generation strategy. Follow the match options and built-in style identification rules in `rules/match-strategy.md`:

| Option | Description |
|--------|-------------|
| Match closest built-in style, customize on top | Identify the closest built-in style, override differences with extracted tokens |
| Fully custom from web app extraction | Generate an independent DESIGN.md entirely from analysis results |

If "match built-in" is chosen, match against the 5 web built-in style characteristics per `rules/style-matching.md` (Vercel, Shadcn, Tailwind UI, Stripe, Apple) and read the corresponding style file.

### TUI

For `--platform tui`, the match strategy differs from web/mobile. Follow the complete TUI match strategy defined in `rules/platform-routing.md` section 4, which uses 2 built-in TUI themes instead of the 5 web styles:

| Option | Description |
|--------|-------------|
| Match closest built-in TUI theme, customize on top | Identify the closest built-in TUI theme (modern-dark-tui or minimal-ascii-tui), override differences with extracted tokens |
| Fully custom from screenshot analysis | Generate an independent TUI DESIGN.md entirely from analysis results |

If "match built-in" is chosen, match against TUI theme characteristics per `rules/style-matching.md` and read the corresponding style file from `ui-design/templates/styles/<name>.md`.

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
| `web` | `templates/design-web.md` | — |
| `mobile` | `templates/design-web.md` | Append `templates/design-mobile.md` after "Responsive Behavior", before "Signature Patterns" |
| `tui` | `templates/design-tui.md` | — |

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

## Error Handling

| Scenario | Action |
|----------|--------|
| URL unreachable (network error, timeout) | Report error with URL. Suggest checking URL and network. Abort. |
| URL returns 4xx/5xx | Report status code. Suggest verifying URL. Abort. |
| CSS extraction returns empty (Layers 1-3 all fail) | Fall back to Layer 4 (agent-browser) if local. Otherwise fall back to Layer 5 (visual inference, mark all as `(estimated)`), warn user about accuracy. |
| agent-browser not available (Layer 4) | Skip Layer 4, use Layer 5. Warn user: "agent-browser unavailable, using visual inference (lower accuracy)". |
| TUI screenshot not found or unreadable | Report error with path. Suggest providing a valid local file path. Abort. |
| TUI screenshot quality too low | Report specific quality issue (blurry, low-res, unreadable). Suggest: use native screenshot tool, ensure text legibility, capture at 1x scale. Abort. |
