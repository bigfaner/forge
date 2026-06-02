# Surface: TUI (Terminal User Interface)

TUI surface 适用于在终端中提供交互式界面的应用程序（基于 Bubble Tea、tview、Textual 等框架）。测试重点是键盘输入、终端渲染输出、异步 Cmd 超时处理。

**Test type**: 终端功能测试 (Terminal Functional Test). Test type: 终端功能测试，通过子进程 + stdin pipe 验证终端渲染输出、键盘输入响应和异步 Cmd 超时处理。Generated test code MUST use `@tui-functional` tags, NOT `@e2e`.

## Detection Signals

| Signal | File Pattern | Dependency Pattern | Exclusion |
|--------|-------------|-------------------|-----------|
| Go TUI | `main.go` exists at project root or in `cmd/` subdirectory | `tea.Program` or `tview.Application` in imports | No frontend framework in `package.json` |
| Node.js TUI | `package.json` exists at project root | `blessed`, `ink`, or `neo-blessed` in `dependencies` | No frontend framework (`react`, `vue`, `svelte`) in dependencies |
| Python TUI | `pyproject.toml` or `setup.py` exists | `rich`, `textual`, or `prompt_toolkit` in dependencies | No frontend entry (no browser DOM references) |
| Rust TUI | `Cargo.toml` exists | `ratatui` or `cursive` in `[dependencies]` | No frontend entry (no browser DOM or web server handler) |

**Confidence Levels**:

- **High**: Primary language entry file + TUI framework dependency + no conflicting signals
- **Medium**: Primary language entry file + TUI framework dependency + some overlapping signals (e.g., also has CLI commands)
- **Low**: Only partial signals (e.g., TUI framework in dependencies but no clear terminal I/O entry point)

**Disambiguation Rules**:

1. If both CLI and TUI signals are present (e.g., `cobra` + `tea.Program`), prefer TUI when the application's primary interaction is terminal-based interactive UI (keyboard navigation, screen rendering) rather than command-line argument processing.
2. If `rich` is detected in Python without `textual`, it may be CLI with rich output rather than TUI. Check for interactive prompts or screen management to confirm TUI.
3. `ink` (React-based terminal UI) in Node.js is a TUI signal even though it uses React concepts -- the rendering target is the terminal, not a browser.

## General Testing Principles

1. **Terminal I/O model**: TUI applications read from stdin (keyboard events) and write to stdout (terminal rendering). Tests must simulate keyboard input sequences and verify rendered output.
2. **Async Cmd timeout**: Every asynchronous Cmd in a TUI application must have a reasonable timeout. Test that:
   - Timed-out Cmds produce a timeout message rather than hanging
   - The UI remains responsive during async operations
   - User can cancel long-running operations
3. **Deterministic rendering**: When possible, test the final rendered state rather than intermediate animations. Use "snapshot" assertions on the model state rather than character-by-character output comparison.
4. **Input simulation**: Test complete input sequences (keystrokes, key combinations, mouse events where applicable) rather than testing individual key presses in isolation.
5. **State machine verification**: TUI applications are state machines. Test state transitions explicitly:
   - Valid transitions produce expected UI updates
   - Invalid inputs are gracefully handled (error message, no crash)
   - Navigation (back/forward/quit) works from every state

## Test Strategy Guidance

**Test Level Emphasis**: Contract 80% / Journey smoke 20%

TUI testing is most effective at the Contract level (individual component/interaction behavior). Journey smoke tests validate complete user workflows through the TUI but are slower due to terminal I/O simulation.

**Execution Model**: Subprocess with stdin pipe

- Compile the TUI binary once before the test suite
- Each test case spawns a subprocess with stdin piped for input simulation
- Capture stdout for output verification
- Set per-test timeouts (default 30s for sync operations, 60s for async Cmd sequences)

**Environment Readiness Checks**:

| Check | How to Verify |
|-------|--------------|
| Binary compiles | `go build` or equivalent succeeds |
| Terminal capability | Verify pseudo-terminal (pty) support is available |
| Stdin pipe works | Test process accepts input via stdin |
| No GUI dependency | TUI should not require X11/Wayland/display server |

**Why subprocess with stdin pipe**: TUI applications typically take over the terminal (raw mode, alternate screen buffer). Testing them in-process would conflict with the test runner's own terminal I/O. Subprocess isolation via pty or stdin pipe provides clean terminal state per test.

## Required Outcome Reference

**Mandatory derived Outcome** (must be considered for every TUI Journey):

- **timeout**: Every asynchronous Cmd must handle timeout gracefully. Example: a network request Cmd that takes too long should display "operation timed out" and return the UI to a responsive state. Assert: UI does not freeze, timeout message is displayed, user can continue interacting.

**Additional common TUI boundary Outcomes**:

- **invalid-input**: Key press or input sequence that the current state does not accept. Assert: error message displayed, state unchanged or gracefully recovered.
- **screen-transition**: Navigation between screens/views. Assert: correct screen rendered, previous state preserved or cleaned up as designed.
- **concurrent-update**: UI receives update while user is typing. Assert: no data corruption, user input is preserved.
- **resize**: Terminal window resize event. Assert: UI reflows correctly, no crash or garbled output.
- **quit-confirmation**: Unsaved state when user attempts to quit. Assert: confirmation dialog shown, user choice respected.
