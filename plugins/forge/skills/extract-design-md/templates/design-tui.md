# Design System: {{Terminal App Name}}

> Extracted from: screenshot analysis
> Date: {{YYYY-MM-DD}}
> Based on: {{modern-dark-tui / minimal-ascii-tui / Custom}}
> Note: All values are (estimated) — reverse-engineered from terminal screenshot via AI vision

## Visual Theme & Atmosphere

{{2-3 sentences describing overall TUI visual style, character density, and terminal aesthetic}}

## Color Space

{{256-color (xterm-256) / 16-color (standard ANSI) / monochrome}} (estimated)

## Character Set

{{Box-drawing + block elements / Pure ASCII / Mixed}} (estimated). Leverages {{Unicode capability / ASCII only}} for visual structure.

### Character Palette Reference

| Element | Character | Unicode/ASCII | Usage |
|---------|-----------|---------------|-------|
| Border corner TL | {{char}} | {{code}} | Panel top-left corner |
| Border corner TR | {{char}} | {{code}} | Panel top-right corner |
| Border corner BL | {{char}} | {{code}} | Panel bottom-left corner |
| Border corner BR | {{char}} | {{code}} | Panel bottom-right corner |
| Border horizontal | {{char}} | {{code}} | Panel top/bottom edges |
| Border vertical | {{char}} | {{code}} | Panel left/right edges |
| Divider | {{char}} | {{code}} | Section divider |
| Bar fill | {{char}} | {{code}} | Bar chart fill |
| Bar empty | {{char}} | {{code}} | Bar chart empty |
| Block full | {{char}} | {{code}} | Progress bar fill |
| Bullet | {{char}} | {{code}} | List item marker |
| Arrow | {{char}} | {{code}} | Navigation indicator |

All characters (estimated).

## Color Palette

{{Dark background / Default background}} with {{high-contrast / minimal}} semantic colors (estimated).

| Role | Color # | Preview | Usage |
|------|---------|---------|-------|
| Background | {{0-255}} | {{description}} | Primary surface |
| Background Alt | {{0-255}} | {{description}} | Alternating rows, inactive panels |
| Surface | {{0-255}} | {{description}} | Cards, focused panel bg |
| Border | {{0-255}} | {{description}} | Panel borders, dividers |
| Border Focus | {{0-255}} | {{description}} | Focused panel border |
| Text Primary | {{0-255}} | {{description}} | Headings, primary text |
| Text Secondary | {{0-255}} | {{description}} | Body text, descriptions |
| Text Tertiary | {{0-255}} | {{description}} | Captions, placeholders |
| Success | {{0-255}} | {{description}} | Positive values, success states |
| Error | {{0-255}} | {{description}} | Errors, destructive actions |
| Warning | {{0-255}} | {{description}} | Warnings, caution states |
| Info | {{0-255}} | {{description}} | Information, links |
| Accent | {{0-255}} | {{description}} | Highlights, selections |

All color values (estimated).

## Typography

Monospaced font only. No proportional fonts in TUI.

| Role | Style | Usage |
|------|-------|-------|
| Title | {{Bold, foreground #}} | Panel headers, app title |
| Heading | {{Bold, foreground #}} | Section headings |
| Body | {{Normal, foreground #}} | Content text |
| Emphasis | {{Bold, foreground #}} | Important values |
| Dim | {{Normal, foreground #}} | Labels, hints, metadata |
| Highlight | {{Bold + reverse video}} | Search matches, cursor |

All styles (estimated).

## Panel Layout

{{Compact / Loose}} density. {{description of overall layout}} (estimated).

- Terminal size: {{rows}} rows x {{columns}} columns (estimated)
- Visible panels: {{count}} (estimated)
- Panel dimensions:
  - {{Panel name}}: {{width}} cols x {{height}} rows (estimated)
- Vertical spacing: {{0-1 / 1-2}} lines between items (estimated)
- Horizontal padding: {{1-2 / 2-4}} characters (estimated)
- Status bar: {{always visible / not present}}, {{1 row / none}} (estimated)

## Key Bindings

| Key | Action |
|-----|--------|
| {{key}} | {{action}} (estimated) |

Key bindings extracted from visible status bar or help panel (estimated).

## Do's and Don'ts

| Do | Don't |
|----|-------|
| Use {{box-drawing / ASCII}} chars for all borders | Use arbitrary characters for borders |
| Use {{256-color / 16-color}} palette values | Hard-code hex colors or RGB |
| Keep {{compact / loose}} density | Inconsistent spacing |
| Specify Unicode codepoint or ASCII code for every char | Leave character choices as "TBD" |
| Use semantic colors (Success=green, Error=red) | Use arbitrary colors for status |

## Applicable Scenarios

- {{scenario 1}}
- {{scenario 2}}
