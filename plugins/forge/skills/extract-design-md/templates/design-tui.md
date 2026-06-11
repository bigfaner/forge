# Design System: {{TERMINAL_APP_NAME}}

> Extracted from: screenshot analysis
> Date: {{DATE}}
> Based on: {{TUI_THEME_NAME_OR_CUSTOM}}
> Note: All values are (estimated) — reverse-engineered from terminal screenshot via AI vision

## Visual Theme & Atmosphere

{{TUI_VISUAL_STYLE_DESCRIPTION}}

## Color Space

{{COLOR_SPACE_TYPE}} (estimated)

## Character Set

{{CHARACTER_SET_TYPE}} (estimated). Leverages {{UNICODE_CAPABILITY}} for visual structure.

### Character Palette Reference

| Element | Character | Unicode/ASCII | Usage |
|---------|-----------|---------------|-------|
| Border corner TL | {{CHAR}} | {{CODE}} | Panel top-left corner |
| Border corner TR | {{CHAR}} | {{CODE}} | Panel top-right corner |
| Border corner BL | {{CHAR}} | {{CODE}} | Panel bottom-left corner |
| Border corner BR | {{CHAR}} | {{CODE}} | Panel bottom-right corner |
| Border horizontal | {{CHAR}} | {{CODE}} | Panel top/bottom edges |
| Border vertical | {{CHAR}} | {{CODE}} | Panel left/right edges |
| Divider | {{CHAR}} | {{CODE}} | Section divider |
| Bar fill | {{CHAR}} | {{CODE}} | Bar chart fill |
| Bar empty | {{CHAR}} | {{CODE}} | Bar chart empty |
| Block full | {{CHAR}} | {{CODE}} | Progress bar fill |
| Bullet | {{CHAR}} | {{CODE}} | List item marker |
| Arrow | {{CHAR}} | {{CODE}} | Navigation indicator |

All characters (estimated).

## Color Palette

{{TUI_BACKGROUND_TYPE}} with {{TUI_CONTRAST_LEVEL}} semantic colors (estimated).

| Role | Color # | Preview | Usage |
|------|---------|---------|-------|
| Background | {{COLOR_NUMBER}} | {{COLOR_DESCRIPTION}} | Primary surface |
| Background Alt | {{COLOR_NUMBER}} | {{COLOR_DESCRIPTION}} | Alternating rows, inactive panels |
| Surface | {{COLOR_NUMBER}} | {{COLOR_DESCRIPTION}} | Cards, focused panel bg |
| Border | {{COLOR_NUMBER}} | {{COLOR_DESCRIPTION}} | Panel borders, dividers |
| Border Focus | {{COLOR_NUMBER}} | {{COLOR_DESCRIPTION}} | Focused panel border |
| Text Primary | {{COLOR_NUMBER}} | {{COLOR_DESCRIPTION}} | Headings, primary text |
| Text Secondary | {{COLOR_NUMBER}} | {{COLOR_DESCRIPTION}} | Body text, descriptions |
| Text Tertiary | {{COLOR_NUMBER}} | {{COLOR_DESCRIPTION}} | Captions, placeholders |
| Success | {{COLOR_NUMBER}} | {{COLOR_DESCRIPTION}} | Positive values, success states |
| Error | {{COLOR_NUMBER}} | {{COLOR_DESCRIPTION}} | Errors, destructive actions |
| Warning | {{COLOR_NUMBER}} | {{COLOR_DESCRIPTION}} | Warnings, caution states |
| Info | {{COLOR_NUMBER}} | {{COLOR_DESCRIPTION}} | Information, links |
| Accent | {{COLOR_NUMBER}} | {{COLOR_DESCRIPTION}} | Highlights, selections |

All color values (estimated).

## Typography

Monospaced font only. No proportional fonts in TUI.

| Role | Style | Usage |
|------|-------|-------|
| Title | {{TITLE_STYLE}} | Panel headers, app title |
| Heading | {{HEADING_STYLE}} | Section headings |
| Body | {{BODY_STYLE}} | Content text |
| Emphasis | {{EMPHASIS_STYLE}} | Important values |
| Dim | {{DIM_STYLE}} | Labels, hints, metadata |
| Highlight | {{HIGHLIGHT_STYLE}} | Search matches, cursor |

All styles (estimated).

## Panel Layout

{{LAYOUT_DENSITY}} density. {{LAYOUT_DESCRIPTION}} (estimated).

- Terminal size: {{TERMINAL_ROWS}} rows x {{TERMINAL_COLUMNS}} columns (estimated)
- Visible panels: {{PANEL_COUNT}} (estimated)
- Panel dimensions:
  - {{PANEL_NAME}}: {{PANEL_WIDTH}} cols x {{PANEL_HEIGHT}} rows (estimated)
- Vertical spacing: {{VERTICAL_SPACING}} lines between items (estimated)
- Horizontal padding: {{HORIZONTAL_PADDING}} characters (estimated)
- Status bar: {{STATUS_BAR_VISIBILITY}}, {{STATUS_BAR_ROWS}} (estimated)

## Key Bindings

| Key | Action |
|-----|--------|
| {{KEY}} | {{KEY_ACTION}} (estimated) |

Key bindings extracted from visible status bar or help panel (estimated).

## Do's and Don'ts

| Do | Don't |
|----|-------|
| Use {{RECOMMENDED_CHAR_TYPE}} chars for all borders | Use arbitrary characters for borders |
| Use {{RECOMMENDED_COLOR_PALETTE}} palette values | Hard-code hex colors or RGB |
| Keep {{RECOMMENDED_DENSITY}} density | Inconsistent spacing |
| Specify Unicode codepoint or ASCII code for every char | Leave character choices as "TBD" |
| Use semantic colors (Success=green, Error=red) | Use arbitrary colors for status |

## Applicable Scenarios

- {{SCENARIO_1}}
- {{SCENARIO_2}}
