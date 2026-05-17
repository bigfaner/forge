# Lesson Entry Template

Used by `/learn` when writing to `docs/lessons/`.

## File Format

```markdown
---
created: "<DATE>"
tags: [<TAG1>, <TAG2>, ...]
---

# <Title>

## Problem
<!-- Symptom description: what happened? what error? -->

## Root Cause
<!-- Root cause: why did it happen? Trace causal chain at least 3 levels deep -->

## Solution
<!-- Solution: how was it fixed? -->

## Reusable Pattern
<!-- Reusable knowledge: what to do next time a similar issue occurs? -->

## Example
<!-- Code example or command (optional) -->

## Related Files
<!-- Related file paths (optional) -->

## References
<!-- Related documentation or links (optional) -->
```

## Tag Vocabulary

Tags must use values from the fixed 8-category vocabulary:

| Tag | Domain |
|-----|--------|
| `architecture` | System structure, layering |
| `interface` | API contracts, data shapes |
| `data-model` | Schema, indexing, soft-delete |
| `dependencies` | Library choices, version constraints |
| `error-handling` | Error types, status codes, propagation |
| `testing` | Test patterns, coverage, mocking |
| `security` | Auth, permissions, data protection |
| `local-dev-deployment` | Dev environment, tooling, deployment |

Select 1-4 tags per entry. If no exact match, pick the closest fit.

## File Naming

Pattern: `<category-prefix><slug>.md`

| Category | Prefix | Example |
|----------|--------|---------|
| Debugging | `debug-` | `debug-race-condition.md` |
| Architecture | `arch-` | `arch-dependency-direction.md` |
| Tooling | `tool-` | `tool-go-test-coverage.md` |
| Pattern | `pattern-` | `pattern-error-wrapping.md` |
| Gotcha | `gotcha-` | `gotcha-context-cancellation.md` |

## Core Principle

Record "what to do next time you encounter a similar problem", not "what I did".
