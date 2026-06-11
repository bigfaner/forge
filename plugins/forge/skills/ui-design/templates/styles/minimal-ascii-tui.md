# Design System: Minimal ASCII TUI

16-color terminal, pure ASCII characters. Default terminal background, loose density. Maximum compatibility.

## Visual Theme & Atmosphere

Minimalist and universally compatible. Pure ASCII characters ensure rendering on every terminal, including legacy systems, SSH sessions, and basic terminal emulators. Relies on character weight and spacing for visual distinction rather than color intensity. Clean and unobtrusive -- the content speaks for itself.

## Color Space

16-color (standard ANSI). Works on any terminal, including Windows CMD, basic SSH sessions, and minimal terminal emulators. No 256-color or true-color requirement.

## Character Set

Pure ASCII. No Unicode characters. Every visual element uses characters from the basic ASCII printable range (0x20-0x7E).

### Character Palette Reference

| Element | Character | ASCII Code | Usage |
|---------|-----------|------------|-------|
| Border corner | + | 0x2B | Panel corners |
| Border horizontal | - | 0x2D | Panel top/bottom edges |
| Border vertical | \| | 0x7C | Panel left/right edges |
| Divider horizontal | - | 0x2D | Section dividers |
| Divider junction | + | 0x2B | Junction points |
| Bar fill | # | 0x23 | Bar chart fill, progress fill |
| Bar empty | . | 0x2E | Bar chart empty |
| Bullet/indicator | * | 0x2A | List item marker, active indicator |
| Arrow right | -> | 0x3E | Navigation indicator |
| Check | [x] | - | Completion indicator |
| Cross | [ ] | - | Incomplete indicator |
| Block indicator | = | 0x3D | Heavy fill, emphasis bar |
| Separator | ~ | 0x7E | Wave separator between sections |
| Subtle fill | - | 0x2D | Dim bar, placeholder |

## Color Palette

Default terminal background. Colors used sparingly for semantic distinction only. Most visual differentiation comes from character weight and spacing.

| Role | Color | ANSI Code | Usage |
|------|-------|-----------|-------|
| Background | Default | - | Terminal default background |
| Text Primary | Default bold | 1 | Headings, emphasis |
| Text Secondary | Default | 0 | Body text |
| Text Dim | Default dim | 2 | Labels, hints, metadata |
| Success | Green | 32 | Positive values, success |
| Error | Red | 31 | Errors, destructive |
| Warning | Yellow | 33 | Warnings |
| Info | Blue | 34 | Information, links |
| Highlight | Reverse video | 7 | Cursor line, selection |
| Accent | Cyan | 36 | Active elements, focus |

Note: When colors are unavailable (monochrome terminal), all visual distinction falls back to character weight (bold vs normal) and spacing.

## Typography

Monospaced font only. No font specification.

| Role | Style | Usage |
|------|-------|-------|
| Title | Bold | Panel headers, app title |
| Heading | Bold | Section headings |
| Body | Normal | Content text |
| Emphasis | Bold | Important values |
| Dim | Dim/faint | Labels, hints, metadata |
| Highlight | Reverse video | Cursor line, search matches |

## Density

Loose. Prioritize readability and universal compatibility over information density.

- Vertical spacing: 1-2 lines between items
- Horizontal padding: 2-4 characters
- Panel borders: single-line ASCII (`+-|+`), 1 char width
- Section separators: blank line + optional `---` line
- Generous whitespace for visual clarity

## Components

### Panels
- ASCII border: `+--- Title ---+`
- No focus state color change (rely on cursor position)
- 2-char horizontal padding inside borders

### Lists
- Items prefixed with `*` or index number
- Cursor line: reverse video highlight
- Blank line between logical groups

### Tables
- Columns aligned with space padding
- Header row: text with `-` separator below
- No alternating row colors (distinguish by spacing)

### Bar Charts
- Fill: `#` in semantic color (green/red/blue)
- Empty: `.` in dim style
- Label right-aligned, value left-padded
- Maximum bar width: 20 chars

### Progress Bars
- Fill: `#` in semantic color
- Empty: `-` in dim style
- Percentage label right of bar
- Example: `[##########----------] 50%`

### Status Bar
- Full width, bottom of screen
- Enclosed in `[]` or `--` borders
- Left: current mode, Right: key hints
- Dim text, minimal color usage

### Scrollbar
- Not recommended for ASCII theme
- Use pagination indicators instead: `[1/5]`, `[More]`, `(n items)`
- If needed: `|` track with `#` thumb

## Do's and Don'ts

| Do | Don't |
|-----|-------|
| Use only ASCII printable characters | Use any Unicode (box-drawing, block elements) |
| Rely on spacing for visual distinction | Depend on color alone |
| Keep layout simple and wide | Crowd content into tight spaces |
| Use reverse video for cursor/selection | Assume color is always available |
| Test on basic terminals (CMD, xterm) | Assume modern terminal features |
| Use bold for emphasis sparingly | Use underline (renders poorly on many terminals) |
| Provide pagination for long lists | Attempt smooth scrolling indicators |

## Applicable Scenarios

- Minimal CLI tools prioritizing maximum compatibility
- SSH sessions to remote servers with unknown terminal capabilities
- CI/CD output and log viewers
- Boot/ Rescue environment tools
- Tools targeting Windows CMD or legacy terminals
- Any TUI application where universal rendering is more important than visual richness
- Quick utility scripts with simple TUI interfaces
