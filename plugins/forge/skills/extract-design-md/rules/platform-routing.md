# Platform-Specific Extraction

## Mobile Extraction

When `--platform mobile`, reuse the web extraction pipeline (Layers 1-5) with a mobile User-Agent viewport context, then add mobile-specific analysis:

1. **Fetch with mobile context**: When using WebFetch or agent-browser, set mobile viewport headers (viewport width: 375px, User-Agent: mobile) to trigger responsive CSS. Reuse all web extraction layers (Layer 1-5) unchanged -- the same CSS bundle parsing, custom property extraction, multi-page sampling, and visual inference apply.

2. **Responsive breakpoint analysis**: Scan CSS for `@media` queries. Extract common mobile breakpoints:
   - 320px (small phone / iPhone SE)
   - 375px (standard phone / iPhone 12/13/14)
   - 414px (large phone / iPhone Plus/Pro Max)
   - 768px (tablet / iPad)
   Record which breakpoints the target site uses and what layout changes occur at each.

3. **Touch target estimation**: Analyze interactive elements (buttons, links, inputs) from CSS for minimum size compliance. Check `width`, `height`, `min-width`, `min-height`, `padding` on interactive selectors. Flag elements below the 44x44pt minimum touch target guideline. Values extracted from computed CSS; if not directly specified, mark as `(estimated)`.

4. **Safe area handling**: Check CSS for `env(safe-area-inset-*)` usage (notch/home indicator on iOS). Check HTML `<meta name="viewport">` for `viewport-fit=cover`. If neither is present, note that safe area handling was not detected and values are `(estimated)`.

> **Limitation**: Mobile extraction depends on the target URL serving responsive CSS. Sites without responsive stylesheets will produce web-equivalent results with mobile-specific sections marked `(estimated)`.

## TUI Extraction

When `--platform tui`, the input must be a **local file path** to a terminal screenshot (not a URL -- screenshots cannot be fetched remotely). AI vision analyzes the screenshot to reverse-engineer design tokens, since TUI has no CSS or structured source to parse.

1. **Validate screenshot input**: The argument must be a local file path (e.g. `./screenshot.png`, `/tmp/terminal.png`). If a URL is provided instead, stop and output:

   > ERROR: TUI platform requires a local screenshot file path, not a URL. Provide a path like ./screenshot.png

   Use the `Read` tool to load the image file. If the file does not exist or cannot be read, stop with a clear error.

2. **Screenshot quality check**: Before detailed analysis, assess the screenshot quality. If the screenshot is blurry, low-resolution, or unreadable, stop and output:

   > ERROR: Screenshot quality is too low for reliable analysis. Please provide a clear, high-resolution terminal screenshot. Tips: use native screenshot tool (not photo of screen), ensure text is legible, capture at 1x scale.

3. **AI vision analysis**: Use the `Read` tool on the screenshot file to perform visual analysis. Extract the following categories, marking **ALL values as `(estimated)`** since AI vision inference is inherently approximate:

   - **ANSI color palette**: Identify the xterm-256 color numbers used in the screenshot. Map observed colors to the closest xterm-256 palette entries. Record background, text primary/secondary/tertiary, border, and semantic colors (success/error/warning/info). All color values must be xterm-256 numbers (0-255).
   - **Character set**: Determine whether the TUI uses box-drawing characters, block elements, pure ASCII (+-\|*#), or a mix. Identify the specific characters used for borders, dividers, indicators, bar charts, and progress bars.
   - **Panel layout dimensions**: Estimate the number of rows and columns visible in the terminal. Identify panel boundaries and their dimensions (width in columns, height in rows). Note the overall terminal grid size.
   - **Key bindings**: If a status bar, help panel, or key binding legend is visible in the screenshot, extract the key-to-action mappings.

4. **Match strategy for TUI**: After extraction, use `AskUserQuestion` to let the user choose:

   | Option | Description |
   |--------|-------------|
   | Match closest built-in TUI theme, customize on top | Identify the closest built-in TUI theme (modern-dark-tui or minimal-ascii-tui), override differences with extracted tokens |
   | Fully custom from screenshot analysis | Generate an independent TUI DESIGN.md entirely from analysis results |

   **If "match built-in" is chosen:**

   Match against TUI theme characteristics defined in `rules/style-matching.md` to identify the closest built-in TUI theme, then read the corresponding style file from the ui-design skill: `ui-design/templates/styles/<name>.md` (resolve relative to the skills parent directory)

5. **Build TUI design tokens and write DESIGN.md**: Read the template at `templates/design-tui.md`. Fill in results from analysis. All extracted values must be marked `(estimated)`.
