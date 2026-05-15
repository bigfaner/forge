# TUI Platform Rules

Navigation patterns for Terminal User Interface applications (Bubble Tea, Textual, Ratatui, etc.).

## Navigation Architecture

TUI navigation is keyboard-driven. All interactions are expressed through keymaps, panel layouts, and modes.

### Keymap

| Key | Action | Context/Mode |
|-----|--------|--------------|
| `Tab` | Next panel focus | All modes |
| `Shift+Tab` | Previous panel focus | All modes |
| `1`-`9` | Jump to panel by index | All modes |
| `j` / `Down` | Next item / scroll down | List, Table |
| `k` / `Up` | Previous item / scroll up | List, Table |
| `h` / `Left` | Previous column / collapse | Table, Tree |
| `l` / `Right` | Next column / expand | Table, Tree |
| `Enter` | Select / confirm / drill-in | List, Table, Form |
| `Esc` | Back / cancel / close | All modes |
| `q` | Quit / exit mode | All modes |
| `/` | Search / filter | List, Table |
| `:` | Enter command mode | Normal |
| `g` | Go to top | List, Table |
| `G` | Go to bottom | List, Table |
| `?` | Show help / keybindings | All modes |

### Panel Layout

| Panel | View | Position | Size Hint |
|-------|------|----------|-----------|
| Header | All | Top, full width | 1 row, fixed |
| Navigation | All | Below header, left or full width | Variable height |
| Content | All | Center, main area | Fills remaining space |
| Status Bar | All | Bottom, full width | 1 row, fixed |
| Sidebar | Detail | Left side | 20-30% width |
| Detail | Detail | Right of sidebar | Remaining width |
| Modal/Overlay | Overlay | Centered on content | 40-60% width, 40-60% height |

### Modes

| Mode | Description | Default Keybindings |
|------|-------------|---------------------|
| Normal | Browse and navigate panels | `Tab`, `j/k`, `1-9`, `Enter`, `Esc` |
| Insert/Edit | Edit form fields | `Enter` confirm, `Esc` cancel, `Tab` next field |
| Command | Enter `:`-prefixed commands | `Enter` execute, `Esc` cancel |
| Search | Filter current view | `/` enter, `Enter` confirm, `Esc` cancel |
| Help | Show keybinding reference | `?` enter, `Esc`/`q` exit |

### Navigation Rules

- Follow the Platform-Agnostic Rules in prototype.md Navigation Contract
- Keyboard is the only input method; no mouse interaction in design specs
- All panels must be reachable by keyboard (Tab order or number keys)
- Modal/overlay panels capture focus until dismissed (Esc)
- Status bar always shows current mode and available keybindings
- Content panel scrolls independently; header and status bar remain fixed
- Panel focus is visually indicated (border highlight, cursor, or color inversion)

## Structural Requirements

Every TUI panel design must include these 5 mandatory sections. Missing any section is a spec defect.

### 1. ASCII Layout Mockup

Box-drawing characters showing the precise visual structure of each panel. Use the theme's character set for all visual elements.

### 2. Dimensions

Concrete numeric values for every size. No "appropriate" or "approximately".

```
Panel width: viewport - 2 (borders)
Content area: panel width - 4 (border + padding)
Column widths: explicit character counts
Bar chart max: explicit character count
Path truncation: explicit maxLen formula
```

### 3. Character Palette

Every visual element must specify its Unicode character and selection rationale.

| Element | Character | Unicode | Reason |
|---------|-----------|---------|--------|

### 4. Color Mapping

Foreground and background colors from the theme palette for every visual element.

| Element | Character | Foreground | Background |
|---------|-----------|------------|------------|

### 5. Edge Cases

Must cover these 5 mandatory scenarios:

| # | Scenario | Expected |
|---|----------|----------|
| 1 | Narrow terminal (80x24) | Layout does not overflow |
| 2 | Wide terminal (140+ col) | Layout does not distort |
| 3 | Mixed numeric widths (1 vs 100) | Right-pad to widest value, column-aligned |
| 4 | Long paths/strings (>50 chars) | Truncate with ".../" preserving trailing segment |
| 5 | No data | Centered "No data" placeholder |

For features involving CJK text, add scenario 6: CJK character width calculated correctly (2 columns per CJK character).
