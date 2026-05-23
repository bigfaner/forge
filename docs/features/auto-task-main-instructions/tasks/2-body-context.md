---
id: "2"
title: "Add BodyContext struct and renderBody placeholder substitution"
priority: "P1"
estimated_time: "1h"
dependencies: []
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 2: Add BodyContext struct and renderBody placeholder substitution

## Description

Add the `BodyContext` struct to `autogen.go` that carries planning-time data from `BuildIndex()` to template rendering. Add the `renderBody()` function that substitutes `{{PLACEHOLDER}}` tokens in template content with BodyContext fields.

This is the foundation for all subsequent tasks — templates and wiring both depend on this struct and function.

## Reference Files
- `docs/proposals/auto-task-main-instructions/proposal.md` — Source proposal (BodyContext struct, Placeholder Substitution table)
- `forge-cli/pkg/task/autogen.go` — Target file

## Acceptance Criteria

- [ ] `BodyContext` struct added to `autogen.go` with fields: FeatureSlug, Mode, Scope, SuccessCriteria, AcceptanceCriteria, ProjectType, Interfaces
- [ ] `renderBody(templateContent string, def AutoGenTaskDef, ctx BodyContext) string` function added
- [ ] Placeholder substitution handles all 6 tokens: `{{FEATURE_SLUG}}`, `{{MODE}}`, `{{SCOPE}}`, `{{INTERFACES}}`, `{{TEST_TYPE}}`, `{{ACCEPTANCE_CRITERIA}}`
- [ ] Empty fields are handled per spec: `{{MODE}}` → omit mode line, `{{SCOPE}}` → omit scope section, `{{INTERFACES}}` → "See .forge/config.yaml", `{{ACCEPTANCE_CRITERIA}}` → "- [ ] All acceptance criteria met"
- [ ] `GenerateTestTaskMD()` signature updated to `(def AutoGenTaskDef, ctx BodyContext) ([]byte, error)` — second param changes from `_ string` to `ctx BodyContext`
- [ ] Existing callers of `GenerateTestTaskMD()` in `build.go` pass empty `BodyContext{}` to maintain backward compatibility
- [ ] Existing tests pass without modification (backward compatible)

## Hard Rules

- MUST NOT change frontmatter generation logic in `GenerateTestTaskMD()`
- MUST handle empty BodyContext gracefully (all fields zero-valued) — this is the backward-compatibility path
- `renderBody()` MUST use `strings.ReplaceAll` for placeholder substitution (no regex needed)

## Implementation Notes

### BodyContext struct (from proposal):

```go
type BodyContext struct {
    FeatureSlug        string
    Mode               string        // "quick" or "breakdown"
    Scope              []string      // In-scope items from proposal/PRD
    SuccessCriteria    []string      // Success criteria from proposal/PRD
    AcceptanceCriteria []string      // PRD acceptance criteria (breakdown mode)
    ProjectType        string        // From .forge/config.yaml
    Interfaces         []string      // Test interfaces from config
}
```

### Placeholder handling rules:

| Placeholder | Source | When Empty |
|-------------|--------|-----------|
| `{{FEATURE_SLUG}}` | BodyContext.FeatureSlug | Never empty (required) |
| `{{MODE}}` | BodyContext.Mode | Omit mode line |
| `{{SCOPE}}` | BodyContext.Scope | Omit scope section |
| `{{INTERFACES}}` | BodyContext.Interfaces | "See .forge/config.yaml" |
| `{{TEST_TYPE}}` | AutoGenTaskDef.TestType | Omit type line |
| `{{ACCEPTANCE_CRITERIA}}` | BodyContext.AcceptanceCriteria | "- [ ] All acceptance criteria met" |

### GenerateTestTaskMD changes:

In the template-loading branch, replace `buf.Write(data)` with:
1. Call `renderBody(string(data), def, ctx)` to substitute placeholders
2. Write the rendered result
