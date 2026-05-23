# Surface: CLI (Command-Line Interface)

CLI surface 适用于以命令行方式交互的应用程序。测试重点是 exit code、stdout/stderr 输出、参数校验和 subprocess 隔离。

## Detection Signals

| Signal | File Pattern | Dependency Pattern | Exclusion |
|--------|-------------|-------------------|-----------|
| Go CLI | `main.go` exists at project root or in `cmd/` subdirectory | `cobra.Command` or `urfave/cli` in `go.mod` / import paths | No frontend framework entry (no `package.json` with React/Vue/Svelte) |
| Node.js CLI | `package.json` exists at project root | `commander`, `yargs`, `oclif`, or `inquirer` in `dependencies` or `devDependencies` | No frontend framework (`react`, `vue`, `svelte`) in dependencies |
| Python CLI | `pyproject.toml` or `setup.py` exists | `click`, `typer`, or `argparse` in dependencies | No frontend entry (no browser DOM references) |
| Rust CLI | `Cargo.toml` exists | `clap`, `structopt`, or `gum` in `[dependencies]` | No frontend entry (no browser DOM or web server handler) |
| Rust CLI (test-driven) | `Cargo.toml` exists | `#[cfg(test)]` or `cargo test` in project (no explicit CLI framework, no HTTP handler) | No frontend entry, no `ratatui`/`cursive` (TUI), no HTTP framework |

**Confidence Levels**:

- **High**: Primary language entry file (`main.go` / `package.json` / `pyproject.toml` / `Cargo.toml`) + CLI framework dependency + no conflicting signals
- **Medium**: Primary language entry file + CLI framework dependency + some ambiguous signals (e.g., also has HTTP handler for internal API)
- **Low**: Only partial signals detected (e.g., `main.go` without recognized CLI framework)

**Disambiguation Rules**:

1. If `package.json` contains both `commander`/`yargs` and `express`/`fastify`, check for frontend framework. If no frontend framework: prefer API (server-first) over CLI.
2. If `main.go` contains both `cobra.Command` and `http.Handler`, check if the HTTP handler serves an API vs. being an internal health endpoint. Internal endpoints don't disqualify CLI.
3. If both CLI and TUI signals are present (e.g., `cobra` + `tea.Program`), prefer TUI when terminal I/O dominates the user interaction model.

## General Testing Principles

1. **Exit code verification**: Every test must assert the subprocess exit code. Exit code 0 = success, non-zero = failure. Never assume a command succeeds without checking its exit code.
2. **stdout/stderr separation**: Assert stdout for expected output and stderr for error/diagnostic messages. Do not conflate the two streams.
3. **Argument validation**: Test invalid arguments, missing required arguments, and mutually exclusive flag combinations.
4. **Subprocess isolation**: The CLI under test must run as an isolated subprocess. This means:
   - Compile a dedicated test binary (or use `go test -c`, the compiled binary itself)
   - Isolate environment variables (do not inherit test runner's env; set only what the test needs)
   - Use temporary directories for all file I/O (`t.TempDir()` in Go, `tmpdir` fixture in Python)
   - Never share global state between test cases
5. **Idempotency**: CLI commands should be testable in any order when properly isolated. Avoid tests that depend on execution sequence.
6. **Concurrent safety**: If the CLI accesses shared resources (files, databases), test concurrent invocations to detect race conditions.

## Test Strategy Guidance

**Test Level Emphasis**: Contract 80% / Journey smoke 20%

CLI tests are inherently subprocess-based -- each invocation starts a new process with clean state. This makes Contract-level testing (individual command behavior) highly reliable and cost-effective. Journey smoke tests validate the end-to-end user workflow but are slower and less isolatable.

**Execution Model**: Subprocess

- Compile the CLI binary once before the test suite runs
- Each test case spawns a new subprocess with isolated arguments and environment
- Capture exit code, stdout, and stderr for assertions
- Set per-test timeouts (default 30s, configurable) to prevent hanging tests

**Environment Readiness Checks**:

| Check | How to Verify |
|-------|--------------|
| Binary compiles | `go build` or equivalent compilation succeeds |
| Binary is executable | File exists and has execute permission |
| Required external tools available | `which <tool>` returns 0 |
| Config files accessible | Config paths exist or can be created in temp dir |

**Why subprocess over in-process**: CLI applications often modify global state (signal handlers, working directory, environment variables). In-process testing risks polluting the test runner's state. Subprocess isolation guarantees clean state per test.

## Required Outcome Reference

These are common boundary/error Outcomes for CLI applications. Use as reference anchors -- the LLM combines these with actual project context to determine which Outcomes are relevant for each Journey.

**Mandatory derived Outcomes** (must be considered for every CLI Journey):

- **not-found**: The requested resource does not exist. Example: `forge task status 99` when task 99 doesn't exist. Assert: exit code != 0, stderr contains descriptive error message.
- **already-exists**: Attempting to create a resource that already exists. Example: `forge task add` with a task ID that's already taken. Assert: exit code != 0, stderr indicates conflict.

**Additional common CLI boundary Outcomes**:

- **invalid-argument**: Missing required flag, invalid flag value, or mutually exclusive flags combined. Assert: exit code != 0, stderr shows usage/help text.
- **permission-denied**: User lacks permissions for the requested operation (filesystem or API). Assert: exit code != 0, stderr contains permission error.
- **timeout**: Command exceeds expected execution time. Assert: process killed, exit code indicates signal (SIGKILL/SIGTERM).
- **output-format**: Verify output matches expected format (table, JSON, plain text) based on flags like `--output=json`.
