---
id: "4"
title: "Create /release-cli slash command"
priority: "P1"
estimated_time: "1h"
dependencies: [3]
type: "doc"
mainSession: false
---

# 4: Create /release-cli slash command

## Description

Create `.claude/commands/release-cli.md` — a local slash command that automates the CLI release workflow: version bump → commit → tag → push. Developers run `/release-cli` to trigger the full release pipeline.

## Reference Files

- `.claude/commands/upgrade-forge.md`: reference for existing version bump command pattern — reuse prompt structure for version input and git operations (source: proposal.md#Implementation-4)
- `forge-cli/scripts/version.txt`: version file to bump (source: proposal.md#Implementation-4)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `.claude/commands/release-cli.md` | Slash command: read version → prompt bump → update version.txt → commit → tag → push |

### Modify
| File | Changes |
|------|---------|
| (none) |  |

### Delete
| File | Reason |
|------|--------|
| (none) |  |

## Acceptance Criteria

- [ ] Command reads current version from `forge-cli/scripts/version.txt`
- [ ] Prompts developer for new version number with semver bump suggestion (patch/minor/major)
- [ ] Updates `forge-cli/scripts/version.txt` with new version
- [ ] Executes: `git commit -m "chore(forge-cli): bump version to {version}"` → `git tag forge-cli/v{version}` → `git push origin HEAD forge-cli/v{version}`
- [ ] Documents that GitHub Actions will auto-trigger from the tag push to build and publish the release

## Hard Rules

- Tag format: `forge-cli/v{version}` (with `v` prefix). CLI version numbers are independent from Plugin version numbers.

## Implementation Notes

- This command is independent from `/upgrade-forge` (which handles Plugin version bumps). CLI and Plugin versions are managed separately.
- CLI version starts at 5.x.x (historical accumulation), Plugin version is 3.0.0-rc.x — no correspondence between them.
