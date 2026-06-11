---
id: "7"
title: "Write v2-to-v3 migration guide and update config docs"
priority: "P2"
estimated_time: "1.5h"
dependencies: ["6"]
type: "documentation"
mainSession: false
---

# 7: Write v2-to-v3 migration guide and update config docs

## Description
Update all config.yaml examples across documentation to use v3 schema (project-type + interfaces + optional languages). Write migration guide at `docs/proposals/simplify-testing-model/migration-v2-to-v3.md` covering field mapping, profile-to-language mapping, common override patterns, and troubleshooting detection failures.

## Reference Files
- `docs/proposals/simplify-testing-model/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|
| docs/proposals/simplify-testing-model/migration-v2-to-v3.md | v2-to-v3 config migration guide |

### Modify
| File | Changes |
|------|---------|
| Any config.yaml examples in docs/ | Replace test-profiles/capabilities with interfaces/languages |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- Migration guide exists at `docs/proposals/simplify-testing-model/migration-v2-to-v3.md`
- Guide documents v2→v3 field mapping table: test-profiles → (removed, auto-detect), capabilities → interfaces
- Guide documents profile-name-to-language mapping (6 entries): go-test→go, web-playwright→javascript, rust-test→rust, pytest→python, java-junit→java, maestro→mobile
- Guide documents common override patterns: multi-language false positive (`languages: [go]`), monorepo subdirectories
- Guide documents troubleshooting: detection failures, no language detected error
- All config.yaml examples in docs/ use v3 schema (no test-profiles or capabilities fields)
- No references to "profile" or "capability" in user-facing documentation (excluding this migration guide)

## Hard Rules
- Migration guide is the ONLY place where v2 field names should appear post-migration
- Do not include migration tooling — manual config edit is intentional (2-line change)

## Implementation Notes
- The migration guide serves dual purpose: (1) helps existing users upgrade, (2) documents the design rationale for future contributors
- Include concrete before/after config.yaml examples for the most common project types
