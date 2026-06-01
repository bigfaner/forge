---
id: "3"
title: "Create GitHub Actions release workflow"
priority: "P0"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 3: Create GitHub Actions release workflow

## Description

Create `.github/workflows/release-cli.yml` that automatically builds forge CLI binaries for 6 platforms and publishes a GitHub Release when a `forge-cli/v*` tag is pushed. This is the CI/CD backbone of the binary distribution model.

## Reference Files

- `forge-cli/scripts/version.txt`: version file read during build to set LDFLAGS (source: proposal.md#Implementation-3)
- `forge-cli/scripts/install.sh`: uploaded as release asset alongside binaries (source: proposal.md#Implementation-3)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `.github/workflows/release-cli.yml` | CI workflow: tag trigger → 6-platform build matrix → release with binaries + checksums |

### Modify
| File | Changes |
|------|---------|
| (none) |  |

### Delete
| File | Reason |
|------|--------|
| (none) |  |

## Acceptance Criteria

- [ ] Workflow triggers on push tag pattern `forge-cli/v*`
- [ ] Build matrix: 6 platforms (darwin/linux/windows × amd64/arm64), `CGO_ENABLED=0`, LDFLAGS inject version from `version.txt`
- [ ] Windows binaries have `.exe` extension, other platforms have no extension; binary naming: `forge-{version}-{os}-{arch}[.exe]`
- [ ] Release includes all platform binaries, `checksums.txt` (SHA256), and `install.sh` as assets
- [ ] Uses `softprops/action-gh-release@v2` for release creation with `generate_release_notes: true`

## Hard Rules

- Compile command must use `CGO_ENABLED=0` for pure static builds with no external C dependencies

## Implementation Notes

- Use `actions/upload-artifact@v4` + `actions/download-artifact@v4` with `merge-multiple: true` to collect all binaries in the release job
- The `build` job runs on `ubuntu-latest` with cross-compilation via `GOOS`/`GOARCH` — no macOS/Windows runners needed
- LDFLAGS: `-s -w -X forge-cli/pkg/types.Version=${VERSION}` strips debug symbols (~30% size reduction) and injects version
