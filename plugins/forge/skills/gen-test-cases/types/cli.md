---
type: cli
conventions:
  - testing-cli.md
---

# CLI Test Case Generation Instructions

Type-specific Steps 3-4 for **CLI** (command-line binary) test cases. Loaded by the dispatcher after Step 2.5 interface detection.

## Classification Indicators

Classify a PRD criterion as **CLI** when it involves any of:

- Commands and subcommands
- Flags and options
- Output format (text, JSON, table)
- Exit codes
- Positional arguments
- stdin/stdout/stderr content
- Error messages printed to terminal
- Configuration via command-line interface

**CLI vs TUI disambiguation**: CLI produces line-oriented sequential output (e.g., `git`, `docker`, `npm`). TUI clears the terminal and redraws (full-screen rendering). Interactive prompts (line-by-line Q&A using inquirer, cobra) are CLI, not TUI.

**Not CLI**: Build commands (`go build`, `npm run build`), lint/test tools (`grep`, `eslint`), CI scripts — these are developer tooling, not product interfaces.

## Target Derivation

- **Target format**: `cli/<command>`
- Derive `<command>` from the command name or subcommand path (e.g., `cli/deploy`, `cli/config-set`, `cli/auth-login`)

## Test ID Format

- **Test ID**: `<target>/<title-slug>`
- `title-slug` = lowercase title, spaces to hyphens, remove punctuation
- Example: `cli/deploy/valid-config-deploys-successfully`

## Priority Assignment

1. Criterion tied to a core/critical Given/When/Then in the PRD → **P0**
2. Criterion tied to a secondary story, or an explicit error/boundary case for a core story → **P1**
3. Nice-to-have verifications, minor edge cases → **P2**

If the PRD has no explicit priority marking, default P0 for the first story's ACs and P1 for all others.

## TC Format

```markdown
## TC-{NNN}: {Title}
- **Source**: {Story N / AC-N} or {Spec Section X.Y}
- **Type**: CLI
- **Target**: cli/<command>
- **Test ID**: cli/<command>/<title-slug>
- **Pre-conditions**: {What must be true before testing}
- **Steps**:
  1. {Exact command invocation, e.g., forge task query --status pending}
  2. {Additional commands or verification steps if needed}
- **Expected**: {Exit code, stdout content, stderr content}
- **Priority**: P0 | P1 | P2
```

- CLI test cases do NOT include a `Route` field. They use `Target` and describe the command invocation in Steps.
- Steps must specify the exact command with all relevant flags and arguments.
- Expected results must include concrete exit codes and stdout/stderr content assertions.

## Command Coverage

CLI test cases must cover all commands, subcommands, and flag combinations mentioned in the PRD:

- **Commands**: Each top-level command gets at least one happy-path test case.
- **Flags**: Required and optional flags are explicitly tested with concrete values.
- **Arguments**: Positional argument combinations are covered.
- **Output format**: If the PRD specifies output formats (text, JSON, table), each format gets a test case.

## Output Assertion Specificity

CLI test cases require concrete output assertions. Each expected result must include:

- **Exit code**: e.g., "Exit code 0" or "Exit code 1"
- **stdout content**: Exact text, substring match, or pattern (e.g., "stdout contains `Deployed to staging`")
- **stderr content**: Error messages when applicable (e.g., "stderr contains `Error: config not found`")

Vague assertions like "command succeeds" or "shows output" are not acceptable.

## Quality Rules

Apply the 6 Antipattern Prevention rules from the dispatcher's shared rules to every CLI test case. Key CLI-specific reminders:

- **Pre-conditions must be concrete and creatable**: Specify how to set up required state (e.g., "a config file at `./test-config.yaml` with valid credentials" not "config exists").
- **Expected results must be specific and verifiable**: State exact exit code and output content. Not "command works" or "output is correct".
- **Steps describe runtime behavior**: Execute the actual CLI binary, not read source files or check documentation.
- **No meta-testing**: Do not generate test cases like "all commands succeed" or "help flag works for every command" unless the PRD explicitly requires it.

## Output

Write to `docs/features/<slug>/testing/cli-test-cases.md`. Number test cases from TC-001 sequential. End the file with a traceability table:

```markdown
## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Story 1 / AC-1 | CLI | cli/task-query | P0 |
```
