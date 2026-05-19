---
created: "2026-05-19"
tags: [testing, data-model, dependencies]
---

# Task Type "documentation" vs "doc" — Template-Validation Mismatch

## Problem

When running `forge task index`, three tasks failed validation with `invalid type 'documentation'`. The `forge task index` command rejected the type and refused to generate `index.json` until the type was corrected to `"doc"`.

## Root Cause

Three-level causal chain:

1. **Immediate cause**: The task files used `type: "documentation"` in YAML frontmatter, but the CLI's `ValidTypes` map in `types.go` only accepts `"doc"`.
2. **Proximate cause**: The quick-tasks skill template `task-doc.md` hardcodes `type: "documentation"` as its default frontmatter value — a human-readable placeholder that doesn't match the actual type constant.
3. **Systemic cause**: Skill templates and the CLI validation schema are maintained in separate packages (templates in `pkg/prompt/data/` and `pkg/template/data/`, validation in `pkg/task/types.go`) with no compile-time or test-time mechanism to detect drift between template defaults and valid type constants.

## Solution

Changed `type: "documentation"` to `type: "doc"` in all three task files and re-ran `forge task index`. This is a per-session fix — the template still contains the wrong default.

## Reusable Pattern

When writing task frontmatter from skill templates, always cross-reference the template's default `type` value against `ValidTypes` in `forge-cli/pkg/task/types.go`. The authoritative type constants are the `Type*` constants (e.g., `TypeDoc = "doc"`, not `"documentation"`). Never trust template defaults blindly.

## Example

```go
// types.go — the authoritative source
const (
    TypeDoc = "doc"              // correct
    // NOT "documentation"
)

// task-doc.md template — misleading default
// type: "documentation"  ← WRONG, should be "doc"
```

## Related Files

- `forge-cli/pkg/task/types.go` — `ValidTypes` map and `TypeDoc` constant
- `C:\Users\panda\.claude\plugins\cache\forge\forge\3.0.0-rc.5\skills\quick-tasks\templates\task-doc.md` — template with incorrect default
- `forge-cli/pkg/task/testgen.go` — auto-generated tasks use correct constants

## References

- Discovered during `deduplicate-quality-gate` quick-tasks planning session
