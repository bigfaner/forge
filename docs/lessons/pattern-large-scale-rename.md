---
created: "2026-05-19"
tags: [architecture, testing]
---

# Large-scale constant rename requires multi-pass, multi-layer approach

## Problem
Renaming ~20 type constants across a Go codebase (19 source files + 10 test files) required 4 distinct passes before all tests passed. Each pass discovered a new layer of references that the previous pass missed.

## Root Cause
1. **Layer 1: Go identifiers** — `TypeFeature` → `TypeCodingFeature` in source code (sed on constant names)
2. **Layer 2: Non-test string literals** — `other.Type == "fix"` in business logic, `"doc-generation.summary"` in generators
3. **Layer 3: Test struct fields** — `Type: "feature"` in test data structs (different formatting: `Type: "feature"` vs `Type:     "feature"`)
4. **Layer 4: Test expectations** — `want "feature"`, `Contains(content, "type:\"feature\"")`, JSON strings in test data
5. **Layer 5: Unused code** — functions that became dead after the rename (`isDocsOnlyType` after `needsDocEval` was simplified)
6. Root cause: constant renames propagate through 5+ distinct syntactic layers, each invisible to the others

## Solution
Execute rename in explicit passes, verifying compilation between each:
1. **Go identifiers** (constants in source) — `sed` with longest-first ordering
2. **String values in source** — grep for `"feature"`, `"fix"`, etc. in non-test files
3. **Test data + expectations** — sed for `Type: "x"`, `want "x"`, `"type":"x"` patterns
4. **Compile + test** — let the compiler and test runner find remaining references
5. **Lint** — catch unused dead code

## Reusable Pattern
For any rename touching >5 files with >30 references:
1. Pre-compute ALL layers: identifiers, string values, JSON, YAML, test expectations, comments
2. Use `sed` bulk replace ordered longest-first (avoids partial matches)
3. Run `go build` ONCE after ALL sed passes — never compile between individual edits
4. Run tests to discover remaining test-layer references
5. Run lint to catch dead code
6. Total expected time: ~30-45min for a 20-constant, 30-file rename (vs subagent's 30min timeout failure)

## Example
```bash
# Pass 1: Identifiers (longest first to avoid partial matches)
sed -i '' \
  -e 's/TypeTestPipelineVerifyRegression/TypeTestVerifyRegression/g' \
  -e 's/TypeFeature/TypeCodingFeature/g' \
  ... $(grep -rl OLD_NAMES --include="*.go")

# Pass 2: String values in source
sed -i '' 's/"fix"/"coding.fix"/g' quality_gate.go claim.go ...

# Pass 3: Test data (multiple patterns)
find . -name "*_test.go" -exec sed -i '' \
  -e 's/Type: "feature"/Type: "coding.feature"/g' \
  -e 's/Type:     "feature"/Type:     "coding.feature"/g' \
  {} +

# Pass 4: Compile + test + fix remaining
go build ./... && go test ./...
```

## Related Files
- `docs/lessons/arch-constant-rename-whack-a-mole.md` — related lesson on subagent failure
- `forge-cli/pkg/task/types.go` — renamed constants
- `forge-cli/pkg/task/build.go` — updated functions
