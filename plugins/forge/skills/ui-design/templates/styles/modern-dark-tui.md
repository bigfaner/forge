# Design System: Modern Dark TUI

256-color terminal, box-drawing + block elements. Dark background, high contrast, compact density.

## Visual Theme & Atmosphere

Rich, information-dense terminal interface. Box-drawing characters create clean panel boundaries. Block elements enable compact visualizations (bar charts, progress indicators). High contrast on dark background ensures readability in any terminal. Feels professional and modern, like a well-designed CLI dashboard.

## Color Space

256-color (xterm-256). Requires terminal with 256-color support. Most modern terminals (iTerm2, Windows Terminal, Alacritty, Kitty) support this by default.

## Character Set

Box-drawing + block elements. Leverages Unicode capability for precise visual structure.

### Character Palette Reference

| Element | Character | Unicode | Usage |
|---------|-----------|---------|-------|
| Border corner TL | ┌ | U+250C | Panel top-left corner |
| Border corner TR | ┐ | U+2510 | Panel top-right corner |
| Border corner BL | └ | U+2514 | Panel bottom-left corner |
| Border corner BR | ┘ | U+2518 | Panel bottom-right corner |
| Border horizontal | ─ | U+2500 | Panel top/bottom edges |
| Border vertical | │ | U+2502 | Panel left/right edges |
| Divider | ├ | U+251C | Left junction (section divider) |
| Divider | ┤ | U+2524 | Right junction (section divider) |
| Divider cross | ┼ | U+253C | Cross junction |
| Thick border vert | ┃ | U+2503 | Scrollbar thumb |
| Bar fill (half) | ▄ | U+2584 | Half-height bar chart fill |
| Bar empty (half) | _ | U+005F | Half-height bar chart empty |
| Block small | ▪ | U+25AA | Inline indicator, file bar |
| Block full | █ | U+2588 | Full block, progress bar fill |
| Block shade light | ░ | U+2591 | Light shade, progress bar empty |
| Arrow right | → | U+2192 | Navigation indicator |
| Bullet | • | U+2022 | List item marker |
| Check | ✓ | U+2713 | Success/completion indicator |
| Cross | ✗ | U+2717 | Error/failure indicator |

## Color Palette

Dark background with high-contrast semantic colors.

| Role | Color # | Preview | Usage |
|------|---------|---------|-------|
| Background | 235 | Dark gray | Primary surface |
| Background Alt | 236 | Slightly lighter | Alternating rows, inactive panels |
| Surface | 237 | Elevated surface | Cards, focused panel bg |
| Border | 239 | Subtle border | Panel borders, dividers |
| Border Focus | 75 | Bright blue | Focused panel border |
| Text Primary | 252 | Near white | Headings, primary text |
| Text Secondary | 246 | Gray | Body text, descriptions |
| Text Tertiary | 241 | Dim gray | Captions, placeholders |
| Success | 82 | Green | Positive values, success states |
| Error | 196 | Red | Errors, destructive actions |
| Warning | 220 | Yellow | Warnings, caution states |
| Info | 75 | Blue | Information, links |
| Accent | 213 | Pink/magenta | Highlights, selections |
| Accent Hover | 177 | Purple | Hover state on selections |
| Cursor Line BG | 55 | Dark purple | Active/cursor row background |
| Scrollbar Track | 236 | Subtle | Scrollbar background |
| Scrollbar Thumb | 246 | Gray | Scrollbar handle |

## Typography

Monospaced font only. No proportional fonts in TUI.

| Role | Style | Usage |
|------|-------|-------|
| Title | Bold, foreground 252 | Panel headers, app title |
| Heading | Bold, foreground 252 | Section headings |
| Body | Normal, foreground 246 | Content text |
| Emphasis | Bold, foreground 252 | Important values |
| Dim | Normal, foreground 241 | Labels, hints, metadata |
| Highlight | Bold + reverse video | Search matches, cursor |

Font stack: Terminal default monospace. No font specification needed -- the terminal controls the font.

## Density

Compact. Maximize information per screen.

- Vertical spacing: 0-1 lines between items
- Horizontal padding: 1-2 characters
- Panel borders: single-line box-drawing (1 char width)
- Status bar: always visible, 1 row
- Minimize vertical whitespace; use it only for logical section separation

## Components

### Panels
- Single-line box-drawing border (┌─┐│└┘├┤)
- Title in top border: `┌─ Title ──────┐`
- Focus state: border color changes to Border Focus (75)

### Lists
- Items prefixed with bullet (•) or index number
- Cursor line: bg=Cursor Line BG (55), text=bold
- Dim items: foreground 241

### Tables
- Columns aligned with padding
- Header row: bold, underline separator ─
- Alternating row background: toggle between Background (235) and Background Alt (236)
- Cursor row: bg=Cursor Line BG (55)

### Bar Charts
- Fill: ▄ (U+2584) in semantic color (Success/Error/Info)
- Empty: _ (U+005F) in dim color (241)
- Label right-aligned, value left-padded

### Progress Bars
- Fill: █ (U+2588) in semantic color
- Empty: ░ (U+2591) in border color (239)
- Percentage label right of bar

### Status Bar
- Full width, bottom of screen
- Left: current mode
- Center: contextual info
- Right: keybinding hints (dim text)

### Scrollbar
- Track: │ (U+2502) in scrollbar track color (236)
- Thumb: ┃ (U+2503) in scrollbar thumb color (246)
- Right edge of scrollable panel, 1 char wide

## Do's and Don'ts

| Do | Don't |
|----|-------|
| Use box-drawing chars for all borders | Use ASCII chars for borders (`+`, `\|`, `-`) |
| Use 256-color palette values | Hard-code hex colors or RGB |
| Keep compact density, minimize whitespace | Add decorative blank lines |
| Specify Unicode codepoint for every char | Leave character choices as "TBD" |
| Use semantic colors (Success=green, Error=red) | Use arbitrary colors for status |
| Test at 80x24 minimum terminal size | Assume wide terminals only |
| Pad numeric columns to widest value | Allow column misalignment |

## Applicable Scenarios

- CLI dashboards and monitoring tools
- Interactive terminal applications (file managers, system monitors)
- Developer productivity tools with data visualization
- Log viewers and analysis tools
- Git/TUI clients
- Any TUI application targeting modern terminal emulators with 256-color support
