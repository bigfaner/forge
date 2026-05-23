---
id: "5"
title: "Update existing tests and add body content verification tests"
priority: "P1"
estimated_time: "1-2h"
dependencies: ["2", "4"]
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 5: Update existing tests and add body content verification tests

## Description

Update existing tests in `autogen_test.go` to work with the new `GenerateTestTaskMD(def AutoGenTaskDef, ctx BodyContext)` signature, and add new tests that verify body content is correctly generated per strategy type (A/B/C) with proper placeholder substitution.

## Reference Files
- `docs/proposals/auto-task-main-instructions/proposal.md` — Success Criteria for test verification
- `forge-cli/pkg/task/autogen_test.go` — Existing tests to update
- `forge-cli/pkg/task/autogen.go` — BodyContext, renderBody

## Acceptance Criteria

- [ ] All existing `TestGenerateTestTaskMD*` tests updated to pass `BodyContext{}` (backward compat)
- [ ] New test: `TestRenderBody_FeatureSlug` — verifies `{{FEATURE_SLUG}}` substitution
- [ ] New test: `TestRenderBody_ScopeAndInterfaces` — verifies `{{SCOPE}}` and `{{INTERFACES}}` substitution
- [ ] New test: `TestRenderBody_AcceptanceCriteria` — verifies Strategy B criteria filling
- [ ] New test: `TestRenderBody_EmptyFields` — verifies graceful handling (omit sections, use fallbacks)
- [ ] New test: `TestBodyContentPerStrategy` — verifies each of the 13 types gets correct body content with populated BodyContext
- [ ] Strategy A body tests: feature slug + scope + interfaces injected
- [ ] Strategy B body tests: acceptance criteria pre-filled as validation checklist
- [ ] Strategy C body tests: discovery strategy steps present (git diff, directory scan)
- [ ] `go test -race -cover ./pkg/task/...` passes with 80%+ coverage

## Hard Rules

- MUST use table-driven tests for per-type verification
- MUST NOT mock `embed.FS` — use the real embedded templates
- Test coverage target: 80%+ for `autogen.go`

## Implementation Notes

### Test structure:

```go
func TestRenderBody_FeatureSlug(t *testing.T) {
    ctx := BodyContext{FeatureSlug: "my-feature", Mode: "quick"}
    result := renderBody("...{{FEATURE_SLUG}}...({{MODE}} mode)...", def, ctx)
    assert.Contains(t, result, "my-feature")
    assert.Contains(t, result, "quick mode")
}

func TestBodyContentPerStrategy(t *testing.T) {
    tests := []struct{
        name string
        typ string
        ctx BodyContext
        wantContains []string
    }{
        // Strategy A: context injected
        {"gen-cases", TypeTestGenCases, BodyContext{FeatureSlug: "feat", Mode: "quick", Scope: []string{"item1"}, Interfaces: []string{"api"}}, []string{"feat", "item1", "api"}},
        // Strategy B: criteria filled
        {"validation-code", TypeValidationCode, BodyContext{FeatureSlug: "feat", AcceptanceCriteria: []string{"AC1", "AC2"}}, []string{"AC1", "AC2"}},
        // Strategy C: discovery steps
        {"doc-drift", TypeDocDrift, BodyContext{FeatureSlug: "feat"}, []string{"git diff", "domains"}},
    }
    // ...
}
```

### Backward compatibility test:

All existing `GenerateTestTaskMD(def, "feature-slug")` calls need updating to `GenerateTestTaskMD(def, BodyContext{})`. Verify that with empty BodyContext, the output still contains proper frontmatter and a usable body (no broken placeholders left in output).
