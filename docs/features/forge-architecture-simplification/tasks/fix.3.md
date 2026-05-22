---
id: "fix.3"
title: "Remove unused indexPkg imports from build.go and index.go"
priority: "P0"
dependencies: []
status: 
type: "coding.fix"
breaking: true
---

# fix.3: Remove unused indexPkg imports from build.go and index.go

## Problem

Task 2.6 replaced `task.SaveIndex` calls with `index.WithLock` / `SaveIndexAtomic` but left unused `forge-cli/pkg/index` imports in two files.

## Scope
- `forge-cli/pkg/task/build.go:12` — unused import `"forge-cli/pkg/index"` as `indexPkg`
- `forge-cli/internal/cmd/index.go:10` — unused import `"forge-cli/pkg/index"` as `indexPkg`

## Acceptance Criteria
- [ ] `go build ./...` passes (0 errors)
