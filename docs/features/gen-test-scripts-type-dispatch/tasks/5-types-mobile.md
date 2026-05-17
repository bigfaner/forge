---
id: "5"
title: "Create types/mobile.md instruction file"
priority: "P2"
estimated_time: "30min-1h"
dependencies: []
type: "documentation"
mainSession: false
---

# 5: Create types/mobile.md instruction file

## Description

Extract Mobile-specific generation logic into a dedicated `types/mobile.md` type instruction file. Mobile tests verify touch interactions, gestures, screen transitions, and app lifecycle events. This is a minimal/template-based file since no active forge project currently uses mobile testing.

## Reference Files
- `docs/proposals/gen-test-scripts-type-dispatch/proposal.md` — Source proposal
- `plugins/forge/skills/gen-test-cases/types/mobile.md` — Reference structure (gen-test-cases Mobile type file)
- `forge-cli/pkg/profile/profiles/maestro/` — Maestro profile (mobile test framework reference)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/gen-test-scripts/types/mobile.md` | Mobile type instruction file with conventions frontmatter |

## Acceptance Criteria

- [ ] `plugins/forge/skills/gen-test-scripts/types/mobile.md` exists
- [ ] Frontmatter declares `type: mobile` and `conventions: [testing-mobile.md]`
- [ ] Contains a **Reconnaissance Strategy** section with mobile-specific search patterns (app manifest files, maestro flow configs, mobile entry points)
- [ ] Contains a **Fact Table Required Keys** section listing minimum keys for mobile type (app bundle identifier or entry point)
- [ ] Contains a **Verification Method** section describing how to confirm the project is a mobile app (detect mobile framework, app manifests)
- [ ] Contains a **Generation Patterns** section describing how mobile test cases translate to executable scripts (touch/gesture simulation, screen transition assertions, app lifecycle events, element location via accessibility labels)
- [ ] Contains a **Mobile Antipattern Guards** section (device-dependent tests without fixture isolation, hardcoded device dimensions, tests requiring physical device)
- [ ] At least 3 section headings are unique to this file

## Hard Rules

- Content should be minimal but structurally complete — mobile testing is not actively used but the file must exist for architectural completeness
- Ground patterns in the maestro profile's YAML template format

## Implementation Notes

- The proposal notes: "Mobile type file is speculative (no mobile profile usage yet)" — keep this file concise and template-based
- Reference the maestro profile at `forge-cli/pkg/profile/profiles/maestro/` for actual template structure
- Focus on structural completeness (all required sections present) rather than deep domain expertise
