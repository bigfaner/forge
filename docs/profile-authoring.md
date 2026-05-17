# Language Strategy Authoring Guide

How to create a new test language strategy for Forge's pluggable test strategy system.

## Directory Structure

```
forge-cli/pkg/testing/languages/<key>/
  generate.md            # gen-test-scripts language strategy
  run.md                 # run-e2e-tests language strategy
  graduate.md            # graduate-tests language strategy
  justfile-recipes       # Justfile recipe bodies for init-justfile
  templates/             # Code templates (spec files, helpers, config)
```

## Required Files

| File | Purpose |
|------|---------|
| `generate.md` | AI-readable prompt for test script generation (framework-specific rules) |
| `run.md` | AI-readable prompt for test execution and result parsing |
| `graduate.md` | AI-readable prompt for test migration from staging to regression |
| `justfile-recipes` | Justfile recipe bodies for `test-e2e`, `e2e-setup`, `e2e-verify` |

## Language Key Convention

- Lowercase, no hyphens for single-word languages: `go`, `rust`, `python`, `java`
- Framework-specific keys for multi-framework languages: `javascript` (Playwright, the only supported JS framework in v3.0)
- Platform keys for non-language targets: `mobile`
- The language key is both the directory name and the internal identifier used by `forge testing` commands

## Supported Interfaces (Closed Enum)

Each language declares which interface types it supports. This metadata is hardcoded in the Go `languageCapabilities` map.

| Interface | Meaning |
|-----------|---------|
| `web-ui` | Browser UI (DOM interaction) |
| `tui` | Terminal UI (text rendering, keyboard) |
| `mobile-ui` | Mobile UI (touch, gestures) |
| `api` | HTTP/network interface |
| `cli` | Command-line interface |

Adding a new interface type requires changes to Forge core (gen-test-cases, eval-test-cases rubric).

## Strategy File Conventions

Each strategy file (`generate.md`, `run.md`, `graduate.md`) is an AI-readable prompt document. Follow these conventions:

- Use tables for structured data (commands, formats, classifications)
- Include code examples for common patterns
- List anti-patterns (forbidden behaviors) explicitly
- Keep sections concise -- no fluff, no preamble
- Reference template files by their relative path within the `templates/` directory

### generate.md Must Cover

- Test runner and assertion library
- Spec template mapping (which template -> which output file)
- CLI/API/TUI testing patterns with code examples
- Auth mechanism
- Import conventions
- Anti-patterns (forbidden in generated code)
- Compilation check command
- Traceability format

### run.md Must Cover

- Execution command
- Result format (with JSON/text examples)
- Result parsing rules (field mapping table)
- TC ID extraction pattern
- Test type classification rules
- Setup/teardown steps
- Timeout configuration
- Error handling table

### graduate.md Must Cover

- File extension and naming pattern
- Import rewrite rules (or "none needed")
- Validation commands (compilation, test discovery)
- Merge procedure for existing target files
- Shared infrastructure policy

## Justfile Recipes

Must define three recipes:

1. **test-e2e** -- Run tests, support `--feature <slug>` for single-feature runs
2. **e2e-setup** -- Install dependencies (idempotent)
3. **e2e-verify** -- Check for unresolved `// VERIFY:` markers (use language's file extension)

## Detection Registration

Add detection rules to the `DetectLanguages()` function in `forge-cli/pkg/testing/detect.go`:

| Signal | Language Key |
|--------|-------------|
| `<marker-file>` exists | `<key>` |

Detection signals should be unambiguous -- don't overlap with existing languages.

Also add the language's supported interfaces to the `languageCapabilities` map in `embed.go`.

## Steps to Add a New Language

1. Add detection case in `detect.go` (file existence check at project root)
2. Create `languages/<key>/` directory with all required strategy files
3. Add entry to `languageCapabilities` map in `embed.go`

No schema migration, no manifest files, no configuration changes needed.

## Testing a New Language

1. Verify strategy files are non-empty and cover required sections
2. Run the full pipeline on a real project of the target type
3. Verify gen-test-scripts -> run-e2e-tests -> graduate-tests chain works
