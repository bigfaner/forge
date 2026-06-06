---
type: tui
conventions:
  - testing-tui.md
---

# TUI Test Script Generation Instructions

Type-specific Steps for **TUI** (terminal full-screen application) test script generation. Loaded by the dispatcher when interface detection identifies TUI-type test cases.

**Test type**: 终端功能测试 (Terminal Functional Test). Generated test code MUST use `@tui-functional` tags, NOT `@e2e`. The "e2e" label does NOT apply to TUI tests — TUI tests verify terminal output rendering and key-sequence interactions via non-interactive stdin pipe execution.

## Classification Indicators

Classify test cases as **TUI** when they involve any of:

- Terminal screen rendering (full-screen redraw)
- Keyboard navigation and key bindings
- Text output assertions against screen content (exact strings, regex, snapshots)
- Screen transitions and cursor movement
- Terminal state changes (alternate screen buffer, raw mode)
- Interactive terminal UI elements (panels, dialogs, lists, status bars)

**TUI vs CLI disambiguation**: TUI clears the terminal and redraws (full-screen rendering, e.g., `vim`, `htop`, `lazygit`). CLI produces line-oriented sequential output. Interactive prompts (line-by-line Q&A) are CLI, not TUI. The distinguishing test: if the application takes over the entire terminal and renders its own screen boundaries, it is TUI; if it prints lines to stdout and exits, it is CLI.

## Golden Rules

Framework-agnostic constraints for TUI test generation. These rules define WHAT to enforce; Convention files define HOW to implement them in a specific framework.

### Shared Principles (per _shared.md)

All TUI tests must satisfy the five cross-type principles defined in `_shared.md`:

- **Isolation**: Each test creates its own working directory, environment, and resource scope
- **Determinism**: Tests must not depend on non-reproducible values or external services
- **Timeout Protection**: Every blocking operation has an explicit upper-bound timeout
- **Idempotency**: Running a test multiple times must produce the same result
- **Resource Cleanup**: Every acquired resource must be released when the test completes

The shared antipattern guards (Sleep-Based Waits, Hardcoded Configuration, Vacuous Assertions, Source-Code-Level Testing) are defined in `_shared.md` and apply to TUI tests without restatement.

### Non-Interactive Execution Model

TUI test scripts MUST use non-interactive execution (stdin pipe, not real terminal). No interactive test modes.

**Rationale**: A real terminal requires human interaction. Tests that depend on real terminal behavior hang in CI. Piped stdin with captured stdout enables deterministic, repeatable execution.

### Terminal Size Contract

Tests MUST control the terminal size environment to eliminate rendering variance. Set `TERM=dumb` or fixed `LINES` and `COLUMNS` environment variables before launching the TUI binary.

**Rationale**: TUI applications render differently based on terminal dimensions. Without controlling the size, the same test produces different output on different terminals or CI runners, causing false failures or hiding real bugs.

### ANSI Sanitization

Assertions MUST strip ANSI escape sequences from captured output before matching, or use dedicated terminal output parsing. Never assert on raw terminal output that contains control sequences.

**Rationale**: Terminal output contains ANSI escape codes for color, cursor positioning, and screen clearing. These sequences vary by terminal emulator and framework version. Assertions that include ANSI codes in expected strings break when the rendering implementation changes, even if the visible output is identical.

### Stable State Detection

Tests MUST define observable signals for "screen rendering complete" rather than relying on time-based assumptions. Acceptable signals include:

- stdout output has stabilized for a configurable duration (no new output for N milliseconds)
- The subprocess has exited (exit code returned)
- A specific marker string appeared in the output (e.g., a prompt, a status indicator)

**Rationale**: TUI rendering timing varies across systems and load conditions. Time-based assumptions (fixed delays) are either too short (test flakes on slow CI) or too long (wastes time on fast machines). Observable signals provide deterministic wait conditions.

### Key Sequence Encoding

Encode keyboard inputs as stdin characters:

| Key | stdin Encoding | Notes |
|-----|---------------|-------|
| Enter | `\n` or `\r` | Confirm/submit |
| Escape | `\x1b` | Cancel/back |
| Tab | `\t` | Next field |
| Arrow keys | `\x1b[A`, `\x1b[B`, `\x1b[C`, `\x1b[D` | Up/Down/Right/Left |
| Regular characters | Literal character | Letters, digits, symbols |
| Ctrl+C | `\x03` | Interrupt |

### TUI-Specific Antipattern Guards

#### 1. Sleep-Based Waits for Screen Transitions

**Pattern**: Using fixed-duration delays to wait for a TUI screen to render or transition.

**Why harmful**: Terminal rendering timing varies across systems and load conditions. Sleep-based waits are either too short (test flakes on slow CI) or too long (wastes time on fast machines). Masks real timing issues.

**Instead**: Use polling with timeout: repeatedly check the captured output for the expected content within a timeout window. If the TUI framework supports event-driven feedback, wait for the specific event. The non-interactive model (stdin pipe) typically produces output synchronously, so polling is fast.

#### 2. Testing via Source Code Inspection Instead of Runtime

**Pattern**: Reading source code files (model definitions, view functions) and asserting on code text rather than running the TUI binary and capturing output.

**Why harmful**: Tests the source code structure, not runtime behavior. A refactoring that changes internal implementation without changing visible output would break the test for no valid reason. Zero verification of actual terminal rendering.

**Instead**: Always execute the TUI binary with piped stdin, capture the output, and assert on the rendered content. The Fact Table provides ground-truth values for constructing the command -- assertions must verify runtime output, not source code text.

#### 3. Interactive Prompts Without Stdin Pipe

**Pattern**: Launching the TUI binary without piping key sequences to stdin, expecting the test framework to interact with a real terminal.

**Why harmful**: The test hangs indefinitely waiting for user input, causing CI timeouts. The failure mode is a timeout, not an assertion failure -- providing no diagnostic value.

**Instead**: Always pipe the complete key sequence to stdin before launching the binary. If the test case describes a sequence of interactions, encode the full sequence in the stdin pipe. If the sequence cannot be encoded (mouse, resize), mark the test as manual-only.

### Tests Requiring Real Terminal Interaction

Some TUI behaviors cannot be tested via stdin piping (e.g., mouse interactions, resize events, terminal color rendering). These test cases must be explicitly handled:

TUI tests that require real terminal interaction must be explicitly marked as "manual only" with a skip rationale in the generated test file.

- **Detection**: If a test case's Steps involve mouse clicks, window resize, or color-specific assertions, mark it as manual.
- **Generation**: Generate the test function with a skip annotation and a comment explaining why.
- **Do not generate** interactive test modes. The skip rationale must state the specific interaction that requires a real terminal.

## Fact Table Required Keys

After reconnaissance, the Fact Table must contain at least these TUI-specific entries for the completeness gate to pass:

| Key Pattern | Description | Example |
|-------------|-------------|---------|
| `TUI_BINARY` | Name of the executable binary that launches the TUI | `TUI_BINARY` = `myapp-tui` |
| `TUI_ENTRY_POINT` | Source file and function where the TUI initializes | `TUI_ENTRY_POINT` = `cmd/tui/main.go:12` |
| `TUI_KEYBIND_*` | At least one key binding definition | `TUI_KEYBIND_QUIT` = `q`, `TUI_KEYBIND_NAV_DOWN` = `j` |

**Minimum requirement**: `TUI_BINARY` must be non-UNKNOWN, and at least one `TUI_KEYBIND_*` entry must be non-UNKNOWN. If all TUI Fact Table keys are UNKNOWN, skip TUI test generation and emit a WARNING.

**Completeness gate rule** (per SKILL.md Step 1.3 Fact Table build): If all required keys for TUI are UNKNOWN, do NOT generate TUI tests. Individual UNKNOWN keys are acceptable -- only skip when every TUI key is UNKNOWN.

## Verification Method

Before generating TUI test scripts, confirm the project actually exposes a TUI interface. A project that only has HTTP handlers or line-oriented CLI commands does not need TUI test scripts.

Run these checks in order -- first success is sufficient:

| Check | Command | Pass Condition |
|-------|---------|----------------|
| Go Bubble Tea | `grep -rn "bubbletea\|tea.Model\|tea.Program" --include='*.go' .` | At least one match found |
| Go tview/tcell | `grep -rn "tview\.New\|tcell\.NewScreen" --include='*.go' .` | At least one match found |
| Rust ratatui | `grep -rn "ratatui\|Terminal::new" --include='*.rs' .` | At least one match found |
| Python Textual | `grep -rn "from textual\|import textual" --include='*.py' .` | At least one match found |
| Alternate screen | `grep -rn "AlternateScreen\|EnterAltScreen\|enter_alt" --include='*.go' --include='*.rs' --include='*.py' .` | At least one match found (confirms full-screen rendering) |

**If all checks fail**: The project does not expose a TUI product interface. Skip TUI test generation and emit a WARNING suggesting the user verify source structure.

## Generation Patterns

TUI test cases translate to executable scripts using non-interactive process execution with stdin piping and terminal output capture. Follow the active strategy's `generate.md` for framework-specific syntax (imports, test runner, assertion library).

### Screen State Assertions

TUI tests must include concrete assertions against the captured terminal output:

| Assertion Type | Pattern | Example |
|---------------|---------|---------|
| Exact text | Output contains exact string | Assert captured output contains expected screen text |
| Regex match | Output matches pattern | Assert captured output matches expected pattern |
| Snapshot | Output matches golden file | Compare captured output against reference file |
| Absence | Output does not contain text | Assert captured output does not contain error text |

When a test case specifies screen state assertions (cursor position, panel visibility, highlight state), translate them to output content checks. Full cursor position assertion requires terminal escape sequence parsing -- if the test framework does not support this, skip cursor position assertions and note it in the traceability comment.

### Exit Code and Output Combination

Each TUI test must assert both exit code and output content:

1. Assert the exit code matches the test case's Expected field (0 for clean exit, non-zero for error).
2. Assert the captured output contains expected text or matches the expected pattern.
3. For error scenarios, assert stderr contains the expected error message if the test case specifies one.

## Test Ratio Constraint

TUI surface targets **Contract test ratio ≥ 80%**.

- **Formula**: `Contract test functions / (Contract test functions + Journey smoke test functions) × 100%`
- **Minimum**: For a Journey with N steps and M total Outcomes across all steps, generate M Contract test functions and exactly 1 Journey smoke test (happy path). This ensures ratio ≥ 80% for any Journey with ≥ 5 Outcomes.
- **Small Journey adjustment**: If a Journey has fewer than 5 total Outcomes, the 1 smoke test still counts as only 1 function, so the ratio naturally stays high. Do NOT skip the smoke test to inflate the ratio — every Journey MUST have at least 1 smoke test.

**Ratio enforcement rule**: If the generated test plan would result in Contract ratio < 80%, reduce the number of Journey smoke tests (minimum 1) before reducing Contract tests.

## Output

TUI test scripts are written to the adaptive output directory: multi-surface projects use `tests/<surfaceKey>/<journey>/`, single-surface projects use `tests/<journey>/`. Each test function includes a traceability comment linking back to the source test case ID.

## Reconnaissance Hints

<!-- Discovery hints — convert findings to Fact Table values, do not use for generation instructions -->

TUI reconnaissance discovers the project's terminal framework, screen definitions, key bindings, and entry points from source code.

### Search Commands

Run these searches to discover TUI interface details. Adapt file extensions to the project's language.

| Target | Grep Command | What It Finds |
|--------|-------------|---------------|
| Go Bubble Tea | `grep -rn "tea\|bubbletea" --include='*.go' .` | Bubble Tea model definitions, tea.Model implementations |
| Go tview | `grep -rn "tview\|tcell\|termbox" --include='*.go' .` | tview widget definitions, tcell screen setup |
| Rust TUI | `grep -rn "ratatui\|cursive\|termion" --include='*.rs' .` | Ratatui/cursive framework imports, terminal backend setup |
| Python Textual | `grep -rn "textual\|urwid\|rich" --include='*.py' .` | Textual app/screen definitions, urwid widget setup |
| Key bindings | `grep -rn "KeyBinding\|keybind\|key_map\|BindKey\|keybinding" --include='*.go' --include='*.rs' --include='*.py' .` | Key binding registrations and mappings |
| Screen rendering | `grep -rn "View\|Render\|Draw\|Screen\|Page" --include='*.go' --include='*.rs' --include='*.py' .` | Screen/view rendering functions, page transitions |
| Entry point | `grep -rn "func main()" --include='*.go' .` or `grep -rn "if __name__" --include='*.py' .` | Binary entry points that initialize the TUI |
| Alternate screen | `grep -rn "AlternateScreen\|raw.mode\|term.Raw\|EnterAltScreen" --include='*.go' --include='*.rs' --include='*.py' .` | Terminal alternate screen buffer usage (confirms full-screen TUI) |

### Reconnaissance Procedure

1. **Detect terminal framework**: Run the grep commands above. Identify which TUI framework the project uses (bubbletea, tview, ratatui, textual, urwid, etc.).
2. **Map screen definitions**: Extract screen/view/model definitions. Record each screen's name, key bindings, and transition triggers.
3. **Identify key bindings**: For each screen or global scope, collect key-to-action mappings (e.g., `q` -> quit, `j/k` -> navigate, `Enter` -> select).
4. **Locate binary entry point**: Find the main function that initializes the TUI application. Record the binary name.
5. **Discover rendering patterns**: Identify how the application renders output -- character-level drawing, widget-based layouts, or declarative views. This determines how test assertions should capture screen state.
