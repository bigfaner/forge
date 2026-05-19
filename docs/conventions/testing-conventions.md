---
title: Testing Convention Files
domains:
  - testing
  - e2e
  - code-generation
---

# Testing Convention Files

Convention files (`testing-{framework}.md`) define the test generation rules for each language/framework combination. They replace the previous hardcoded profile package approach.

## File Structure

Convention files live in `docs/conventions/` and follow the naming pattern `testing-{framework}.md`:

- `testing-go.md` — Go testing + testify/assert conventions
- `testing-ginkgo.md` — Ginkgo v2 + Gomega conventions
- `testing-vitest.md` — TypeScript Vitest conventions

## Content Format

Each convention file includes:

1. **File patterns** — naming conventions for test files (e.g., `*_test.go`)
2. **Test function patterns** — how tests are declared (e.g., `func Test*`)
3. **Assertion style** — preferred assertion library and patterns
4. **Project structure** — where test files go relative to source
5. **Build/run commands** — how to compile and execute tests
6. **Result parsing** — output format for test result extraction

## Usage by Skills

- **gen-test-scripts**: Loads convention files matching the project's detected languages via domains frontmatter
- **run-e2e-tests**: Uses convention rules for result parsing and output format detection
- **init-justfile**: Generates justfile recipes from convention file build/run commands

## Adding New Frameworks

To add support for a new framework:

1. Create `docs/conventions/testing-{framework}.md` following the content format above
2. Add the language to `forgeconfig.KnownLanguages` in the CLI
3. Add interface capabilities to `forgeconfig.languageCapabilities`
