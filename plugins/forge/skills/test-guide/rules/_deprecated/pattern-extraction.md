# Test File Pattern Extraction

## Test File Location Patterns

| Language       | Test file patterns                              |
| -------------- | ----------------------------------------------- |
| Go             | `**/*_test.go`                                  |
| JavaScript/TS  | `**/*.test.{ts,js,tsx,jsx}`, `**/*.spec.{ts,js,tsx,jsx}` |
| Rust           | `**/tests/*.rs` (integration), `**/*_test.rs`   |
| Python         | `**/test_*.py`, `**/*_test.py`                  |
| Java           | `**/*Test.java`, `**/*Tests.java`               |
| C#             | `**/*Tests.cs`, `**/*Test.cs`                   |

Focus on `tests/`, `tests/<surfaceKey>/<journey>/` (multi-surface), and `tests/<journey>/` (single-surface) directories first (forge convention), then project-wide.

## Import Detection (Assertion Library)

For each test file found, read and extract imports:

- Go: `"github.com/stretchr/testify/assert"`, `"github.com/stretchr/testify/require"`, `"github.com/onsi/gomega"`, `. "github.com/onsi/ginkgo/v2"`
- JS/TS: `from 'vitest'`, `from '@jest/globals'`, `from 'mocha'`, `from 'chai'`
- Python: `import pytest`, `import unittest`, `from assertpy import assert_that`
- Java: `import static org.junit.Assert.*`, `import static org.assertj.core.api.*`, `import org.testng.Assert`

## Tag / Marker Detection

- Go: `//go:build <surface>-<type>`, `//go:build feature`, `// +build <surface>-<type>`
- JS/TS: `describe('@feature', ...)`, `describe('@<surface>-<type>', ...)`, `{ tags: ['@feature'] }`
- Python: `@pytest.mark.<surface>_<type>`, `@pytest.mark.feature`
- Java: `@Tag("<surface>-<type>")`, `@Tag("feature")`

## Test Function Naming

- Go: `TestTC_NNN_Description`, `TestFeatureName`
- JS/TS: `describe('Feature: ...', () => { it('should ...', ...) })`
- Python: `test_tc_nnn_description`, `test_feature_name`
- Java: `testTcNnnDescription`, `shouldDoSomethingWhenCondition`

## Compiled Finding Format

```
Detected patterns:
  Framework: <framework name from imports>
  Assertion library: <library and style from imports + function calls>
  Test tags: <tag format from test files>
  Test naming: <naming pattern from function names>
  File pattern: <test file extension/pattern>
  Result format: <inferred from framework -- e.g., go test -json, vitest --reporter=json>
```
