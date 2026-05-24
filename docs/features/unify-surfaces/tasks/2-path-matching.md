---
id: "2"
title: "Path normalization and segment prefix matching"
priority: "P0"
estimated_time: "1.5h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 2: Path normalization and segment prefix matching

## Description

Implement the path normalization rules and segment prefix matching algorithm. These are used by `forge surfaces <path>` queries and internal surface resolution. Scalar form bypasses matching entirely (any path returns the value).

## Reference Files
- `proposal.md#路径规范化与匹配算法` — 5 normalization rules, segment prefix matching algorithm, why-segments-not-chars rationale
- `proposal.md#CLI-命令与退出码契约` — exit code contract for path queries (exit 0 vs exit 1)
- `proposal.md#Success-Criteria` — path boundary validations (.. rejection, Windows backslash, symlink, segment vs char prefix)

## Acceptance Criteria

- [ ] Path normalization: strip leading `./`, trailing `/`, convert `\` to `/`
- [ ] Paths containing `..` return error ("path contains '..'")
- [ ] Symlinks NOT resolved — literal path matching only
- [ ] Scalar form: any path query returns the value directly, no matching
- [ ] Map form — segment prefix matching: `frontend/api/routes` matches `frontend/api` (2 segments) over `frontend` (1 segment)
- [ ] Map form — no partial match: `frontend-new` does NOT match `frontend`
- [ ] Map form — unmatched path returns error with manual config hint

## Hard Rules

- Matching is by path SEGMENTS (split by `/`), NOT character prefix
- No path resolution or symlink resolution — purely string-based segment comparison

## Implementation Notes

- Suggested location: new file `forge-cli/pkg/forgeconfig/surfaces.go` or `forge-cli/pkg/surfaces/match.go`
- Export two functions: `NormalizePath(string) (string, error)` and `MatchSurface(surfaces map[string]string, query string) (string, error)`
- `MatchSurface` returns `("", error)` when no match found (for exit code 1 propagation)
- For scalar form: when `Surfaces` has single key `"."`, skip matching and return value directly
