---
id: "T-review-doc"
title: "Review Documentation Quality"
priority: "P1"
estimated_time: "30min"
dependencies: ["3", "4", "5", "6"]
type: "doc.review"
surface-key: ""
surface-type: ""
---

Review documentation quality for the forge-cli-binary-distribution feature (quick mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

### 3-create-release-workflow

- [ ] Workflow triggers on push tag pattern `forge-cli/v*`
- [ ] Build matrix: 6 platforms (darwin/linux/windows × amd64/arm64), `CGO_ENABLED=0`, LDFLAGS inject version from `version.txt`
- [ ] Windows binaries have `.exe` extension, other platforms have no extension; binary naming: `forge-{version}-{os}-{arch}[.exe]`
- [ ] Release includes all platform binaries, `checksums.txt` (SHA256), and `install.sh` as assets
- [ ] Uses `softprops/action-gh-release@v2` for release creation with `generate_release_notes: true`


### 4-create-release-cli-command

- [ ] Command reads current version from `forge-cli/scripts/version.txt`
- [ ] Prompts developer for new version number with semver bump suggestion (patch/minor/major)
- [ ] Updates `forge-cli/scripts/version.txt` with new version
- [ ] Executes: `git commit -m "chore(forge-cli): bump version to {version}"` → `git tag forge-cli/v{version}` → `git push origin HEAD forge-cli/v{version}`
- [ ] Documents that GitHub Actions will auto-trigger from the tag push to build and publish the release


### 5-remove-version-bump-rules

- [ ] Version Bump section (lines 28-33, containing "Code changes must bump the version..." and the 3 semver bullet points) is removed from `forge-cli/CLAUDE.md`
- [ ] No other content in `forge-cli/CLAUDE.md` is modified


### 6-remove-init-forge

- [ ] `plugins/forge/commands/init-forge.md` deleted
- [ ] No references to `/init-forge` remain in other plugin command files under `plugins/forge/commands/`


## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/forge-cli-binary-distribution/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/forge-cli-binary-distribution/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.

## Acceptance Criteria

- [ ] All doc task deliverables (tasks 3-6) reviewed against their AC
- [ ] Review covers only files under the allowlist directories, excluding tasks/ and build artifacts
- [ ] Issues found are reported with specific file and line references
