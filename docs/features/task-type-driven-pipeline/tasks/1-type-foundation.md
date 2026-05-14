---
id: "1"
title: "Add documentation/doc-evaluation type constants and remove InferType fallback"
priority: "P1"
estimated_time: "1h"
dependencies: []
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 1: Add documentation/doc-evaluation type constants and remove InferType fallback

## Description
Foundation task for the task-type-driven pipeline. Add two new task type constants (`TypeDocumentation`, `TypeDocEvaluation`) to the type system, update the registry and validation map, and remove the `TypeImplementation` fallback from `InferType` so that business tasks without an explicit type will surface as empty (enabling hard-error detection in task 2).

## Reference Files
- `docs/proposals/task-type-driven-pipeline/proposal.md` — Source proposal (D1, D2)

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/task/types.go` | Add `TypeDocumentation` and `TypeDocEvaluation` constants; add entries to `TaskTypeRegistry` and `ValidTypes` |
| `forge-cli/pkg/task/infer.go` | Remove `default: return TypeImplementation` branch; add `T-eval-doc` ID pattern |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] `TypeDocumentation = "documentation"` and `TypeDocEvaluation = "doc-evaluation"` constants exist in `types.go`
- [ ] `TaskTypeRegistry` includes both new types with descriptions
- [ ] `ValidTypes` map includes both new types
- [ ] `InferType` returns `""` for unknown IDs (no `TypeImplementation` fallback)
- [ ] `InferType` recognizes `T-eval-doc` pattern and returns `TypeDocEvaluation`
- [ ] Existing tests in `types_test.go` and `infer_test.go` pass after changes
- [ ] New table-driven tests cover: `InferType` returning `""` for business IDs, `T-eval-doc` pattern, both new constants, `ValidTypes` entries

## Implementation Notes
- The fallback removal is safe because `build.go` will handle empty types with a hard error (task 2). Until task 2 lands, `InferType` returning `""` for business tasks is non-breaking — the current `build.go` still sets type to the `InferType` result, which would be empty but won't error until task 2 adds the check.
- `T-eval-doc` is a single task ID (not profile-suffixed), so the pattern check is a simple string comparison.
