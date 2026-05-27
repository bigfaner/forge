---
title: "Go Ginkgo Testing Convention"
---

# Go Ginkgo Testing Convention

Convention for generating Go test code using the Ginkgo BDD framework with Gomega matchers.

## framework

- **name**: Ginkgo v2 + Gomega
- **version**: ginkgo v2+
- **language**: Go
- **runner_command**: `ginkgo -v --json-report=report.json -tags=cli_functional`

## discovery

- **test_dir**: `tests/e2e/`
- **file_pattern**: `*_test.go`
- **exclude_pattern**: `vendor/`, `node_modules/`

## structure

- **suite_pattern**: `var _ = Describe("...", func() { ... })` — top-level BDD container
- **case_pattern**: `It("should ...", func() { ... })` — individual spec within Describe/Context
- **hook_pattern**: `BeforeEach` / `AfterEach` / `BeforeSuite` / `AfterSuite`

### Spec Structure

Use `Describe` / `Context` / `It` BDD nesting:

```go
var _ = Describe("Feature: Task Lifecycle", func() {
    var (
        projectDir string
    )

    BeforeEach(func() {
        projectDir = setupTestProject()
    })

    Describe("Task claiming", func() {
        Context("when tasks are available", func() {
            It("should claim a task successfully", func() {
                cmd := exec.Command("forge", "task", "claim")
                cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectDir)
                out, err := cmd.CombinedOutput()
                Expect(err).ShouldNot(HaveOccurred())
                Expect(string(out)).Should(ContainSubstring("claimed task"))
            })
        })
    })
})
```

### Test Entry Point

Every Ginkgo test suite needs a bootstrap file:

```go
//go:build e2e

package e2e

import (
    "testing"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

func TestE2E(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "E2E Test Suite")
}
```

### Table-Driven Specs

Use Ginkgo's `DescribeTable` / `Entry`:

```go
DescribeTable("validation",
    func(input, expected string) {
        result := process(input)
        Expect(result).Should(Equal(expected))
    },
    Entry("valid input", "hello", "HELLO"),
    Entry("empty input", "", ""),
)
```

### Traceability

Each `It` block should include a traceability comment:

```go
It("should login successfully", Label("TC-001"), func() {
    // Traceability: TC-001 -> PRD User Auth section
})
```

## assertions

- **style**: should
- **library**: `github.com/onsi/gomega`
- **custom_matchers**: none

### Key Functions

- `Expect(actual).Should(Equal(expected))` — equality
- `Expect(actual).Should(BeNil())` — nil check
- `Expect(err).ShouldNot(HaveOccurred())` — no error
- `Expect(str).Should(ContainSubstring(substr))` — substring
- `Expect(actual).Should(BeTrue())` — boolean true
- `Expect(actual).Should(BeFalse())` — boolean false
- `Expect(actual).Should(BeEmpty())` — empty collection
- `Expect(actual).ShouldNot(BeNil())` — not nil
- `Expect(err).Should(HaveOccurred())` — error expected

**Rule**: Always use `Expect` with `Should`/`ShouldNot` matcher chains.

## Tags

- **Build tag**: `//go:build e2e` must be the first line of every test file
- **Format**: Pure Go build tag syntax

```go
//go:build e2e

package e2e
```

## Result Format

- **Output flags**: `-json -v`
- **Format type**: `json-stream` (same as Go testing — one JSON object per line)

## Import Patterns

Standard imports for Ginkgo e2e tests:

```go
import (
    "os/exec"
    "os"
    "time"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)
```

- HTTP tests add: `"net/http"`, `"net/http/httptest"`
- File tests add: `"path/filepath"`

## Anti-patterns (Forbidden)

| Forbidden | Replacement |
|-----------|-------------|
| `time.Sleep` for synchronization | `Eventually()` with matcher |
| `assert`/`require` from testify | `Expect` with Gomega matchers |
| Hardcoded ports | `net.Listen("tcp", ":0")` or environment variables |
| Real secrets/tokens in code | `os.Getenv("E2E_API_TOKEN")` |
| `//go:build` tag missing | Always include `//go:build e2e` as first line |
| Deep nesting (>3 levels of Describe/Context) | Flatten with extracted helper functions |
| Mixed assertion libraries | Use only Gomega, never testify |

## Helpers

### runCLI helper

```go
func runCLI(args ...string) string {
    cmd := exec.Command("forge", args...)
    out, err := cmd.CombinedOutput()
    Expect(err).ShouldNot(HaveOccurred(), "forge %s failed: %s",
        strings.Join(args, " "), string(out))
    return string(out)
}
```

### Eventually helper

```go
Eventually(func() string {
    out, _ := exec.Command("forge", "task", "status").CombinedOutput()
    return string(out)
}, 5*time.Second, 100*time.Millisecond).Should(ContainSubstring("completed"))
```
