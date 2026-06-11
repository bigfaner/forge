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

Every TUI panel design must include all mandatory structural requirements defined in `rules/tui-panel-requirements.md`. This includes the 5 mandatory items (ASCII Layout Mockup, Dimensions, Character Palette, Color Mapping, Edge Cases) plus additional per-panel specs (States, Key Bindings, Data Binding).

Missing any mandatory section is a spec defect. Refer to `rules/tui-panel-requirements.md` for the authoritative definitions, edge case scenarios, and enforcement rules.
