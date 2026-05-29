---
id: "1"
title: "Register doc.fix type in Go type system"
priority: "P0"
estimated_time: "1h"
complexity: "medium"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.feature"
mainSession: false
---

# 1: Register doc.fix type in Go type system

## Description

The task dispatcher hardcodes `coding.fix` for all fix tasks regardless of source task category. This causes doc/eval task failures to spawn code-level fix tasks with irrelevant quality gates (golangci-lint, go test), leading to cascading failures and infinite fix-task chains.

Register `doc.fix` as a first-class type in the Go type system so that doc-category failures can spawn properly scoped fix tasks. This covers: type constant, validation registration, InferType pattern matching, and template defaults.

## Reference Files

- `forge-cli/pkg/task/types.go`: Add `TypeDocFix = "doc.fix"` constant, register in ValidTypes, SystemTypes, and TaskTypeRegistry (source: proposal.md#In-Scope)
- `forge-cli/pkg/task/infer.go`: InferType() must recognize `doc-fix-` prefixed IDs as `doc.fix` type — add case before existing `fix-` match to avoid collision (source: proposal.md#In-Scope)
- `forge-cli/pkg/task/category.go`: CategoryForType() uses `HasPrefix("doc")` which auto-maps `doc.fix` → CategoryDoc — verify no change needed (source: proposal.md#Constraints-&-Dependencies)
- `forge-cli/pkg/task/tasktemplate.go`: Add `"doc.fix"` entry to taskTemplateDefaults map with IDPrefix `"doc-fix"` (source: proposal.md#Constraints-&-Dependencies)

## Acceptance Criteria

- [ ] `TypeDocFix` constant with value `"doc.fix"` registered in ValidTypes, SystemTypes, and TaskTypeRegistry
- [ ] `InferType()` recognizes `doc-fix-` prefixed IDs (e.g. `doc-fix-1`) as `doc.fix` type
- [ ] `taskTemplateDefaults` includes `"doc.fix"` entry with `IDPrefix: "doc-fix"`, `Priority: P0`, `Breaking: false`
- [ ] `forge task add --type doc.fix --title "Fix: test" --source-task-id <any-id>` succeeds without validation error

## Implementation Notes

### Category auto-mapping

`CategoryForType()` already uses `strings.HasPrefix(typ, "doc")` which returns `CategoryDoc` — no change needed in category.go. Similarly, `IsTestableType()` checks `strings.HasPrefix(typ, "coding.")` so `doc.fix` correctly returns false (no test pipeline).

### InferType ordering

The `doc-fix-` prefix check MUST appear before the existing `fix-` check in InferType's switch statement, because `"fix-"` is a prefix of `"fix-doc-"` but not vice versa. The correct order: check `doc-fix-` first, then `fix-`.

### Task Impact

- Affected test suite(s): `forge-cli/pkg/task/` (types_test.go, infer_test.go, tasktemplate_test.go)
- Expected fixture changes: None
- Risk level: low
