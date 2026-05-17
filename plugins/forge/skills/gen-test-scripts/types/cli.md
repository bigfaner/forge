---
type: cli
conventions:
  - testing-cli.md
---

# CLI Test Script Generation Instructions

Type-specific Steps for **CLI** (command-line binary) test script generation. Loaded by the dispatcher when interface detection identifies CLI-type test cases.

## Classification Indicators

Classify test cases as **CLI** when they involve any of:

- Commands and subcommands
- Flags and options
- Output format assertions (text, JSON, table)
- Exit code checks
- Positional arguments
- stdin/stdout/stderr content
- Error messages printed to terminal
- Configuration via command-line interface

**CLI vs TUI disambiguation**: CLI produces line-oriented sequential output (e.g., `git`, `docker`, `npm`). TUI clears the terminal and redraws (full-screen rendering). Interactive prompts (line-by-line Q&A using inquirer, cobra) are CLI, not TUI.

**Not CLI**: Build commands (`go build`, `npm run build`), lint/test tools (`grep`, `eslint`), CI scripts -- these are developer tooling, not product interfaces.

## Reconnaissance Strategy

CLI reconnaissance discovers the project's command structure, flag definitions, and entry points from source code.

### Search Commands

Run these searches to discover CLI interface details. Adapt file extensions to the project's language.

| Target | Grep Command | What It Finds |
|--------|-------------|---------------|
| Go CLI framework | `grep -rn "cobra.Command" --include='*.go' .` | Cobra command registration patterns, command names and descriptions |
| Go CLI framework | `grep -rn "\.Flags()\|\.PersistentFlags()" --include='*.go' .` | Flag definitions (required and optional) |
| Go entry points | `grep -rn "func main()" --include='*.go' .` | Binary entry points |
| Node.js CLI framework | `grep -rn "Command\\|program\\|option\\|command(" --include='*.js' --include='*.ts' .` | Commander/Yargs command registration |
| Python CLI framework | `grep -rn "@click\\|@app.command\\|argparse" --include='*.py' .` | Click/argparse/Typer command definitions |
| Binary declarations | `grep '"bin"' package.json` | Node.js binary entry points (name, path) |
| Command directory | `ls cmd/` | Go-style command directory structure (one file per subcommand) |
| Configuration flags | `grep -rn "flag\\.\\|pflag\\." --include='*.go' .` | Go standard/pflag flag parsing |

### Reconnaissance Procedure

1. **Detect CLI framework**: Run the grep commands above. Identify which CLI framework the project uses (cobra, commander, click, argparse, etc.).
2. **Map command tree**: Extract top-level commands and subcommands. Record each command's name, description, and flag set.
3. **Identify flag definitions**: For each command, collect flag names, types (string, bool, int), required/optional status, and default values.
4. **Locate binary entry point**: Find the main function or bin declaration that wires commands together. Record the binary name.
5. **Discover output formats**: Search for output formatting logic (JSON, table, text) to understand how to assert on command output.

## Fact Table Required Keys

After reconnaissance, the Fact Table must contain at least these CLI-specific entries for the completeness gate to pass:

| Key Pattern | Description | Example |
|-------------|-------------|---------|
| `CLI_BINARY` | Name of the executable binary | `CLI_BINARY` = `myapp` |
| `CLI_COMMAND_*` | At least one command name entry | `CLI_COMMAND_DEPLOY` = `deploy` |
| `CLI_FLAG_*` | Flag names used in test cases | `CLI_FLAG_ENV` = `--env` |

**Minimum requirement**: At least one `CLI_COMMAND_*` entry must be non-UNKNOWN. If all CLI Fact Table keys are UNKNOWN, skip CLI test generation and emit a WARNING.

**Completeness gate rule** (from SKILL.md Step 1.5): If all required keys for CLI are UNKNOWN, do NOT generate CLI tests. Individual UNKNOWN keys are acceptable -- only skip when every CLI key is UNKNOWN.

## Verification Method

Before generating CLI test scripts, confirm the project actually exposes a CLI interface. A project that only has HTTP handlers or a frontend does not need CLI test scripts.

Run these checks in order -- first success is sufficient:

| Check | Command | Pass Condition |
|-------|---------|----------------|
| Node.js binary | `grep '"bin"' package.json` | Key exists with a path value |
| Go command directory | `ls cmd/` | Directory exists and contains `.go` files |
| Cobra framework | `grep -rn "cobra.Command" --include='*.go' .` | At least one match found |
| Commander/Yargs | `grep -rn "new Command\\|yargs\\|program" --include='*.js' --include='*.ts' .` | At least one match found |
| Click/argparse | `grep -rn "@click\\|argparse" --include='*.py' .` | At least one match found |

**If all checks fail**: The project does not expose a CLI product interface. Skip CLI test generation and emit a WARNING suggesting the user verify source structure. Build/lint/test commands are developer tooling, not CLI product interfaces.

## Generation Patterns

CLI test cases translate to executable scripts using process execution patterns. Follow the active profile's `generate.md` for framework-specific syntax (imports, test runner, assertion library).

### Process Execution

Each CLI test function invokes the project's binary as a subprocess:

1. **Build the binary** (if needed): Run the project's build command to ensure the binary exists. Use `t.TempDir()` or framework equivalent for isolation.
2. **Execute the command**: Spawn the binary with the command, flags, and arguments specified in the test case's Steps field.
3. **Capture output**: Collect stdout, stderr, and exit code from the subprocess.
4. **Assert results**: Compare captured values against the test case's Expected field.

### Assertion Patterns

CLI tests must include concrete assertions for each dimension:

| Dimension | Assertion Pattern | Example |
|-----------|-------------------|---------|
| Exit code | Assert exact exit code | Exit code 0 for success, 1 for error |
| stdout | Assert contains/exact/matches | `assert.Contains(stdout, "Deployed to staging")` |
| stderr | Assert contains (error cases) | `assert.Contains(stderr, "Error: config not found")` |
| Output format | Assert JSON structure or table headers | `assert.JSONEq(stdout, expectedJSON)` |

### Argument and Flag Testing

- **Required flags**: Test with and without required flags to verify both success and error paths.
- **Flag values**: Use concrete values from test cases. Do not invent flag values.
- **Positional arguments**: Test each argument position explicitly.
- **Flag combinations**: When test cases specify multiple flags, pass all of them in the command invocation.

### Environment Isolation

Each CLI test must create its own isolated environment:

- Use `t.TempDir()` (or framework equivalent) as working directory.
- Set environment variables explicitly (do not rely on host environment).
- Clean up created resources within the test scope.

## CLI Antipattern Guards

Beyond the generic 6 antipattern guards in the main SKILL.md, CLI-specific generation must avoid these additional patterns:

### 1. Recursive Test Invocation

**Pattern**: Calling `exec.Command("go", "test"` (or the equivalent test-runner invocation for any language) from within a test that belongs to the same package/module being tested.

**Why harmful**: Causes process explosion -- each recursion level spawns more processes. On Windows, orphaned children persist indefinitely, consuming 6GB+ RAM.

**Instead**: If a meta-test must verify "all tests pass", use a recursion guard: set an environment variable before spawning the subprocess and check it at the top of the calling test. Or exclude the meta-test via `-run` flag filtering.

### 2. Static File Text Grep

**Pattern**: Reading static source files (`.md`, `.go`, `.json`, `.yaml`) and asserting on text content via `assert.Contains(out, "some string")`.

**Why harmful**: Tests documentation text, not runtime behavior. A typo fix in a markdown file breaks the test without any functional regression. Zero verification value -- the test is coupled to prose, not behavior.

**Instead**: Only test runtime behavior: invoke the CLI binary and assert on outputs, exit codes, or rendered content. Never read source files as test input.

### 3. Interactive Prompts Without Automation

**Pattern**: Invoking a CLI command that triggers an interactive prompt (e.g., "Are you sure? [y/N]") without piping the expected input.

**Why harmful**: The test hangs indefinitely waiting for stdin input, causing CI timeouts. The failure mode is a timeout, not an assertion failure -- providing no diagnostic value.

**Instead**: Either pipe the expected input to stdin (`echo "y" | myapp deploy`), or pass the non-interactive flag if the CLI provides one (`--yes`, `--no-input`, `--force`). If the test case requires interactive behavior, explicitly set up the stdin pipe in the test.

## Output

CLI test scripts are written to `tests/e2e/features/<feature>/` following the profile's template naming convention. Each test function includes a traceability comment linking back to the source test case ID.
