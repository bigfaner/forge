---
id: "1"
title: "Embed just binaries for 6 platforms with build tags"
priority: "P1"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 1: Embed just binaries for 6 platforms with build tags

## Description

Create the embedded binary infrastructure for `just` (casey/just) across 6 target platforms. This is the foundational task — subsequent tasks will use these embedded binaries as a fallback installation method.

The approach uses Go build tags to include only the current platform's binary at compile time, keeping the final forge CLI binary increase to ~4 MB instead of ~24 MB.

## Reference Files
- `docs/proposals/forge-init-install-just/proposal.md` — Source proposal
- `forge-cli/internal/embedded/claudemd.go` — Existing embed pattern (reference for go:embed usage)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `forge-cli/internal/embedded/just/just_binary.go` | Interface file declaring `JustBinary() []byte` (unconditionally compiled) |
| `forge-cli/internal/embedded/just/just_binary_linux_amd64.go` | Build-tagged file embedding linux/amd64 binary |
| `forge-cli/internal/embedded/just/just_binary_linux_arm64.go` | Build-tagged file embedding linux/arm64 binary |
| `forge-cli/internal/embedded/just/just_binary_darwin_amd64.go` | Build-tagged file embedding darwin/amd64 binary |
| `forge-cli/internal/embedded/just/just_binary_darwin_arm64.go` | Build-tagged file embedding darwin/arm64 binary |
| `forge-cli/internal/embedded/just/just_binary_windows_amd64.go` | Build-tagged file embedding windows/amd64 binary |
| `forge-cli/internal/embedded/just/just_binary_windows_arm64.go` | Build-tagged file embedding windows/arm64 binary |
| `forge-cli/internal/embedded/just/binaries/.gitkeep` | Placeholder directory for platform binaries |
| `forge-cli/scripts/download-just.sh` | Script to download just release binaries for all 6 platforms |

### Modify
| File | Changes |
|------|---------|
| `forge-cli/.gitignore` | Add `internal/embedded/just/binaries/*.gz` to ignore downloaded archives |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] `JustBinary() []byte` function returns non-nil, non-empty byte slice on the current platform
- [ ] `go build ./...` succeeds on all platforms (build-tagged files compile correctly)
- [ ] Each platform file has correct build tags: `//go:build linux && amd64`, etc.
- [ ] Binary size increase for forge CLI is <= 5 MB per platform
- [ ] `scripts/download-just.sh` downloads just v1.40.0 binaries from GitHub releases for all 6 platforms

## Hard Rules

- Use `//go:build` directives (not `// +build`) for build tags
- The interface file (`just_binary.go`) MUST NOT embed anything — it only declares the function signature
- Each platform file MUST embed exactly one binary, tagged with both OS and architecture
- Do NOT commit the actual binary files to git — only the `.gitkeep` placeholder

## Implementation Notes

- Follow the pattern in `forge-cli/internal/embedded/claudemd.go` for go:embed usage, but adapt for binary files
- Download just v1.40.0 (latest stable as of 2025-05) from `https://github.com/casey/just/releases`
- Platform binary naming convention: `just-{version}-{os}-{arch}` (e.g., `just-1.40.0-x86_64-unknown-linux-gnu`)
- The download script should verify checksums if available
- Consider using `compress/gzip` if binaries need decompression at extraction time
