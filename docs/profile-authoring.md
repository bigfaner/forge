# Profile Authoring Guide

How to create a new test profile for Forge's pluggable test strategy system.

## Directory Structure

```
task-cli/pkg/profile/profiles/<name>/
  manifest.yaml          # Metadata + capabilities + command declarations
  generate.md            # gen-test-scripts profile strategy
  run.md                 # run-e2e-tests profile strategy
  graduate.md            # graduate-tests profile strategy
  justfile-recipes       # Justfile recipe bodies for init-justfile
  templates/             # Code templates (spec files, helpers, config)
```

## Required Files

| File | Purpose |
|------|---------|
| `manifest.yaml` | Declares profile name, language, capabilities, templates, run/graduate commands |
| `generate.md` | AI-readable prompt for test script generation (framework-specific rules) |
| `run.md` | AI-readable prompt for test execution and result parsing |
| `graduate.md` | AI-readable prompt for test migration from staging to regression |
| `justfile-recipes` | Justfile recipe bodies for `test-e2e`, `e2e-setup`, `e2e-verify` |

## Manifest Schema

```yaml
name: <profile-name>          # lowercase, hyphenated
display: "Human-readable name"
language: <language>           # go, typescript, java, rust, python, yaml
file-extension: <ext>          # _test.go, .spec.ts, .java, .rs, .py, .yaml
test-directory: tests/e2e/     # always tests/e2e/

capabilities: [<cap>, ...]     # closed enum: web-ui, tui, mobile-ui, api, cli

templates:                     # paths relative to profile directory
  test-file: templates/<file>
  helpers: templates/<file>    # optional
  config-file: templates/<file> # optional
  additional: []               # extra template files

run:
  command: "<test-command>"
  compile: "<compile-command or null>"
  result-format: <format-name>

graduate:
  target-directory: tests/e2e/
  merge-strategy: <package|class|module|file>
  import-rewrite: <rule or null>
  compile-check: "<command or null>"
  list-tests: "<command>"
```

### Capabilities (Closed Enum)

| Capability | Meaning |
|-----------|---------|
| `web-ui` | Browser UI (DOM interaction) |
| `tui` | Terminal UI (text rendering, keyboard) |
| `mobile-ui` | Mobile UI (touch, gestures) |
| `api` | HTTP/network interface |
| `cli` | Command-line interface |

Adding a new capability requires changes to Forge core (gen-test-cases, eval-test-cases rubric).

## Strategy File Conventions

Each strategy file (`generate.md`, `run.md`, `graduate.md`) is an AI-readable prompt document. Follow these conventions:

- Use tables for structured data (commands, formats, classifications)
- Include code examples for common patterns
- List anti-patterns (forbidden behaviors) explicitly
- Keep sections concise — no fluff, no preamble
- Reference template files by their manifest key

### generate.md Must Cover

- Test runner and assertion library
- Spec template mapping (which template → which output file)
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

1. **test-e2e** — Run tests, support `--feature <slug>` for single-feature runs
2. **e2e-setup** — Install dependencies (idempotent)
3. **e2e-verify** — Check for unresolved `// VERIFY:` markers (use profile's file extension)

## Auto-Detection Registration

Add detection rules to `plugins/forge/references/shared/profile-detection.md`:

| Signal | Profile |
|--------|---------|
| `<marker-file>` exists | `<profile-name>` |

Detection signals should be unambiguous — don't overlap with existing profiles.

## Testing a New Profile

1. Verify `manifest.yaml` parses as valid YAML
2. Verify strategy files are non-empty and cover required sections
3. Run the full pipeline on a real project of the target type
4. Verify gen-test-scripts → run-e2e-tests → graduate-tests chain works
