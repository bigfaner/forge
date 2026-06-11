---
status: "completed"
started: "2026-06-01 21:27"
completed: "2026-06-01 21:29"
time_spent: "~2m"
---

# Task Record: 3 Create GitHub Actions release workflow

## Summary
Created .github/workflows/release-cli.yml — GitHub Actions workflow that builds forge CLI binaries for 6 platforms (darwin/linux/windows x amd64/arm64) and publishes a GitHub Release when a forge-cli/v* tag is pushed

## Changes

### Files Created
- .github/workflows/release-cli.yml

### Files Modified
无

### Key Decisions
无

## Document Metrics
~65 lines, single YAML workflow file covering 2 jobs (build + release)

## Referenced Documents
- docs/proposals/forge-cli-binary-distribution/proposal.md
- forge-cli/scripts/version.txt
- forge-cli/scripts/install.sh

## Review Status
final

## Acceptance Criteria
- [x] Workflow triggers on push tag pattern forge-cli/v*
- [x] Build matrix: 6 platforms, CGO_ENABLED=0, LDFLAGS inject version from version.txt
- [x] Windows binaries have .exe extension, others no extension; naming: forge-{version}-{os}-{arch}[.exe]
- [x] Release includes all binaries, checksums.txt (SHA256), and install.sh as assets
- [x] Uses softprops/action-gh-release@v2 with generate_release_notes: true

## Notes
Workflow follows proposal.md Implementation-3 exactly. Build job uses ubuntu-latest with cross-compilation via GOOS/GOARCH. Release job uses upload-artifact + download-artifact with merge-multiple to collect all binaries.
