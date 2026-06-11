---
name: step-1-contract-loading
description: Code reconnaissance (Fact Table building), semantic descriptor resolution, and framework/domain pattern extraction
---

# Step 1: Code Reconnaissance (Build Fact Table)

Read source code to extract ground-truth values for semantic descriptor resolution AND test framework patterns.

## 1.1 Framework Reconnaissance

Scan existing test files for framework-specific patterns to supplement or validate Convention:

| Source | What to extract | Purpose |
|--------|-----------------|---------|
| Test file names | File pattern (`*_test.go`, `*.test.ts`) | Confirm Convention's `file-pattern` |
| Test file imports | Assertion library, test runner imports | Confirm Convention's `Assertion` section |
| Build tags / markers | Tag syntax (`//go:build e2e`, `@pytest.mark.e2e`) | Confirm Convention's `Tags` section |
| Test function signatures | Naming pattern, parameter types | Infer naming conventions |
| Test helper files | Utility patterns, fixture setup | Infer project-specific test patterns |

If no test files exist and no file signals are recognizable, Reconnaissance produces an empty Fact Table for framework columns. This is expected cold-start behavior -- proceed with Convention alone, or LLM defaults if no Convention.

## 1.2 Domain Reconnaissance

Extract ground-truth values from application source code for semantic descriptor resolution:

| Source | What to extract |
|--------|-----------------|
| CLI entry points | Command names, flag names, output format strings |
| API handlers | Request/response schemas, status codes |
| TUI components | Model fields, View output patterns |
| Config files | Ports, base paths, auth mechanisms |
| Auth implementation | Login endpoint, token field, header format |

## 1.3 Build Fact Table

Combine all reconnaissance into a single Fact Table with source citations:

```markdown
## Fact Table
| Key | Value | Source |
|-----|-------|--------|
| CLI_TASK_CLAIM_OUTPUT | claimed task <task_id> | internal/cmd/claim.go:42 |
| CLI_FEATURE_CREATE_OUTPUT | Feature <slug> created successfully | internal/cmd/feature.go:45 |
| TEST_FRAMEWORK | go-testing | tests/<surfaceKey>/<journey>/step1_test.go or tests/<journey>/step1_test.go (import analysis) |
| TEST_ASSERTION_LIB | testify/assert | tests/<surfaceKey>/<journey>/step1_test.go:3 or tests/<journey>/step1_test.go:3 (import) |
| TEST_BUILD_TAG | //go:build e2e | tests/<surfaceKey>/<journey>/step1_test.go:1 or tests/<journey>/step1_test.go:1 |
```

<HARD-RULE>
- Every Fact Table value must cite source file and line number. Unknown sources -> `UNKNOWN`. Do not fabricate.
- Fact Table values drive semantic descriptor to regex conversion. All `// VERIFY:` markers must be resolved using Fact Table values.
- When Reconnaissance finds signals that conflict with Convention -> Convention wins, but log the conflict for user awareness.
</HARD-RULE>

## 1.4 Semantic Descriptor to Regex Conversion

Contract Output dimensions use semantic descriptors (natural language), not regex. This step converts them to precise regex patterns:

1. For each Outcome's Output dimension, look up matching Fact Table entries.
2. Convert the Fact Table value to a regex pattern:
   - Placeholder tokens like `<task_id>` become named capture groups: `(?P<task_id>[\w-]+)`
   - Literal text is regex-escaped.
3. If no Fact Table match is found, keep the original descriptor as a `// VERIFY:` marker.

Example pipeline:
```
Semantic descriptor: "success confirmation containing feature-slug"
  -> Fact Table lookup: CLI_FEATURE_CREATE_OUTPUT = "Feature my-feature created successfully"
  -> Generated regex: Feature\s+([\w-]+)\s+created\ successfully
```
