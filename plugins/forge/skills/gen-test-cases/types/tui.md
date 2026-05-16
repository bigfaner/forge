---
type: tui
conventions:
  - testing-tui.md
---

# TUI Test Case Generation Instructions

Type-specific Steps 3-4 for **TUI** (terminal text rendering, keyboard-driven) test cases. Loaded by the dispatcher after Step 2.5 interface detection.

## Classification Indicators

Classify a PRD criterion as **TUI** when it involves any of:

- Terminal screen rendering (full-screen redraw)
- Keyboard navigation and key bindings
- Text output assertions (exact strings, regex patterns)
- Screen transitions and cursor movement
- Terminal state changes
- Interactive terminal UI elements (panels, dialogs, lists)

**TUI vs CLI disambiguation**: TUI clears the terminal and redraws (full-screen rendering, e.g., `vim`, `htop`, `lazygit`). CLI produces line-oriented sequential output. Interactive prompts (line-by-line Q&A) are CLI, not TUI.

## Target Derivation

- **Target format**: `tui/<screen-name>`
- Derive `<screen-name>` from the screen/view name (e.g., `tui/dashboard`, `tui/help-panel`, `tui/file-browser`)

## Test ID Format

- **Test ID**: `<target>/<title-slug>`
- `title-slug` = lowercase title, spaces to hyphens, remove punctuation
- Example: `tui/dashboard/navigate-to-help-panel`

## Priority Assignment

1. Criterion tied to a core/critical Given/When/Then in the PRD → **P0**
2. Criterion tied to a secondary story, or an explicit error/boundary case for a core story → **P1**
3. Nice-to-have verifications, minor edge cases → **P2**

If the PRD has no explicit priority marking, default P0 for the first story's ACs and P1 for all others.

## TC Format

```markdown
## TC-{NNN}: {Title}
- **Source**: {Story N / AC-N} or {Spec Section X.Y}
- **Type**: TUI
- **Target**: tui/<screen-name>
- **Test ID**: tui/<screen-name>/<title-slug>
- **Pre-conditions**: {What must be true before testing}
- **Route**: {Screen identifier or path — required for TUI tests}
- **Steps**:
  1. {Step 1}
  2. {Step 2}
- **Expected**: {What the correct result looks like}
- **Priority**: P0 | P1 | P2
```

- `Route` field is required for TUI test cases — must contain a concrete screen identifier.
- Expected results for TUI must specify exact text, snapshot comparison points, or regex patterns for terminal output. Not "displays correctly".
- Keyboard inputs and key sequences must be explicitly described in Steps.

## Output Assertion Specificity

TUI test cases require concrete output assertions. Each expected result must include at least one of:

- **Exact text**: "Output contains `Error: file not found`"
- **Regex pattern**: "Output matches `/^\d+ files processed$/`"
- **Snapshot reference**: "Screen matches golden snapshot `dashboard-loaded.txt`"
- **Terminal state**: "Cursor is at row 5, column 1"

Vague assertions like "shows the result" or "displays output" are not acceptable.

## Route Validation

Cross-reference each TUI test case's `Route` field against actual screen/route definitions.

**Discovery patterns** (framework-specific):
- Go (tview, bubbletea, tcell): screen component definitions, `tview.NewPages`, tea.Model transitions
- Rust (ratatui, cursive): `App::run`, screen/state enums, event loop handlers
- Python (textual, urwid): `APP.push_screen`, screen class names

**Validation**: For each test case with a `Route` field:
- Match against discovered screen definitions → annotate `Matched (source:line)`
- No match → annotate `Route not found -- verify path`

If no route/screen definitions can be discovered, skip this step entirely. Do not fabricate validation results.

## Quality Rules

Apply the 6 Antipattern Prevention rules from the dispatcher's shared rules to every TUI test case. Key TUI-specific reminders:

- **Pre-conditions must be concrete and creatable**: Specify terminal state requirements (e.g., "terminal size 80x24", "application started with `--config test.yaml`").
- **Expected results must be specific and verifiable**: Use exact text, regex, or snapshot references. Not "output looks right".
- **Steps describe runtime behavior**: Interact with the running TUI application (press key, type input, navigate menu), not read source files.

## Output

Write to `docs/features/<slug>/testing/tui-test-cases.md`. Number test cases from TC-001 sequential. End the file with a traceability table:

```markdown
## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Story 1 / AC-1 | TUI | tui/dashboard | P0 |
```
