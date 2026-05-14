<!-- Snippet to update manifest.md after /ui-design completes -->

## Documents (updated)

Add row(s) based on platform count:

<!-- Single platform (web or mobile only): -->
| UI Design | ui/ui-design.md | {{UI_DESIGN_SUMMARY}} |
| Prototype | prototype/ | {{PROTOTYPE_SUMMARY}} |

<!-- Single platform (TUI only): -->
<!-- | UI Design | ui/ui-design-tui.md | {{UI_DESIGN_SUMMARY}} | -->
<!-- | Prototype | prototype/ | {{PROTOTYPE_SUMMARY}} | -->

<!-- Multi-platform (e.g. web + tui): -->
<!-- | UI Design (web) | ui/ui-design-web.md | {{UI_DESIGN_SUMMARY_WEB}} | -->
<!-- | Prototype (web) | prototype/web/ | {{PROTOTYPE_SUMMARY_WEB}} | -->
<!-- | UI Design (tui) | ui/ui-design-tui.md | {{UI_DESIGN_SUMMARY_TUI}} | -->
<!-- | Prototype (tui) | prototype/tui/ | {{PROTOTYPE_SUMMARY_TUI}} | -->

## Traceability (updated)

Add entries linking PRD UI functions to UI design sections.
For multi-platform features, repeat the table for each platform's design file:

<!-- Per platform: -->
| PRD Section | Design Section | UI Component | Tasks |
|-------------|----------------|--------------|-------|
| "UI Functions > {{Function Name}}" | "UI Design > {{Component Name}}" | "{{Component Name}}" (ui-design §N) | <!-- task IDs added by /breakdown-tasks --> |

## Frontmatter

Update `status` to `design`.
