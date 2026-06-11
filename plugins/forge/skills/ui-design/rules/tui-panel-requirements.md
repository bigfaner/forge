# TUI Panel Design Requirements

When platform=tui, each panel MUST include all 5 mandatory structural requirements:

1. **ASCII Layout Mockup** -- box-drawing (or ASCII) illustration showing the exact visual structure of the panel
2. **Dimensions** -- concrete numeric values for every size (e.g., "panel width: 60 chars"). No "approximately" or "appropriate".
3. **Character Palette** -- every visual element must specify its Unicode character (with code point) and selection reason
4. **Color Mapping** -- foreground/background color codes from the theme palette for every visual element
5. **Edge Cases** -- must cover 5 mandatory scenarios: (1) narrow terminal 80x24, (2) wide terminal 140+col, (3) mixed numeric widths, (4) long strings/paths >50 chars, (5) no data

<HARD-RULE>
These 5 structural requirements are MANDATORY for every TUI panel. Skipping any item is a spec defect. This is the key lesson from the deep-drill-analytics feature -- without enforcement, agents skip visual specs and produce iterative trial-and-error fix/style commits.
</HARD-RULE>

In addition to the 5 structural requirements, each TUI panel must also specify:
- **States** (loading, empty, error, populated)
- **Key Bindings** (what keys interact with this panel, from `platforms/tui.md` keymap)
- **Data Binding** (UI element -> data field)
