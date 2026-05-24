---
created: "2026-05-24"
tags: [architecture, testing]
---

# Breaking Change Task Stuck at Quality Gate

## Problem

Task executor executes a "breaking change" task (e.g. remove `Interfaces []string` from Config struct, replace with `Surfaces SurfacesMap`). The code change compiles in isolation, but `forge task submit` runs the quality gate which includes `go build ./...` — and compilation fails because other files still call the deleted function (`ReadInterfaces()`). The executor gets stuck: fixing the callers feels like scope creep into a downstream task, but not fixing them means the quality gate can never pass.

## Root Cause

1. **Surface**: `go build` fails during quality gate — callers of the deleted API don't compile.
2. **Deeper**: Task boundaries were drawn along conceptual lines (Task 1: "struct change", Task 4: "migrate callers") rather than along **compilation boundaries**.
3. **Deepest**: The quality gate assumes each task leaves the codebase in a compilable state. Breaking change tasks violate this assumption by design — they remove an API surface that callers still depend on. The gate has no mechanism to distinguish "expected compilation failure from intentional removal" from "accidental breakage".

## Solution

Include the **minimal caller updates** in the breaking change task itself. The rule is: if deleting a public function/field causes `go build` to fail, updating those callers is part of the same task — not scope creep.

In the unify-surfaces case: Task 1 originally scoped only the Config struct change. The fix was to also update `build.go` (2 lines: `ReadSurfaces` + `SurfaceTypes` instead of `ReadInterfaces`). This is the minimal change needed for compilation. Deeper migrations (renaming `uiInterfaces` map, `BodyContext.Interfaces` field) can remain in Task 4 because they don't affect compilation.

## Reusable Pattern

**When planning breaking change tasks**: scope the task to include all changes required for `go build` to pass. Draw the boundary at "compiles + existing tests pass", not at "only touched the target struct".

Checklist for breaking change task planning:
1. Identify the API surface being removed (field, function, type)
2. `grep` for all references to it across the codebase
3. Classify references into:
   - **Compilation-blocking** (direct function call, field access) → must be in this task
   - **Non-blocking** (variable naming, comments, internal struct fields) → can defer to downstream task
4. Include step 3's "blocking" references in the task scope

## Example

```
Task 1: Remove Interfaces []string, add Surfaces SurfacesMap

Blocking references (must fix in Task 1):
  - build.go: forgeconfig.ReadInterfaces() → ReadSurfaces() + SurfaceTypes()

Non-blocking references (can defer to Task 4):
  - autogen.go: BodyContext.Interfaces field → rename later
  - autogen.go: uiInterfaces map → rename later
  - extract.go: Interfaces field in struct literal → rename later
```

## Related Files

- `forge-cli/pkg/forgeconfig/config.go` — Config struct definition
- `forge-cli/pkg/forgeconfig/detect.go` — ReadSurfaces, migration logic
- `forge-cli/pkg/task/build.go` — Minimal caller update (2 lines)
