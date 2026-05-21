---
id: "1"
title: "Create types/_shared.md cross-type golden rules"
priority: "P0"
estimated_time: "1-2h"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Create types/_shared.md cross-type golden rules

## Description

Create `plugins/forge/skills/gen-test-scripts/types/_shared.md` defining five cross-type universal golden rules that all 5 type files reference instead of repeating. This is the foundation for the three-layer model: `_shared.md` (abstract principles) → type file Golden Rules (type-specific constraints) → Convention (framework implementation).

## Reference Files
- `docs/proposals/gen-test-scripts-golden-rules/proposal.md` — Source proposal
- `plugins/forge/skills/gen-test-scripts/types/cli.md` — Existing CLI type (for shared antipattern extraction)
- `plugins/forge/skills/gen-test-scripts/types/tui.md` — Existing TUI type
- `plugins/forge/skills/gen-test-scripts/types/api.md` — Existing API type

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/gen-test-scripts/types/_shared.md` | Cross-type universal golden rules |

### Modify
| File | Changes |
|------|---------|
| — | — |

### Delete
| File | Reason |
|------|--------|
| — | — |

## Acceptance Criteria

- [ ] `_shared.md` defines 5 universal principles: Isolation, Determinism, Timeout Protection, Idempotency, Resource Cleanup
- [ ] Each principle has: declarative constraint statement, rationale, and shared antipattern guard
- [ ] Determinism principle is expanded with sub-dimensions: (a) no random dependency (timestamps, UUIDs → fixed values), (b) no external service dependency (third-party APIs, email), (c) no order dependency (each test independently runnable)
- [ ] Timeout Protection principle covers: all I/O operations (subprocess, HTTP, wait conditions) must have timeout upper bound; default value + Convention override mechanism
- [ ] Resource Cleanup principle covers: tests must not leave behind temp files, background processes, database records, browser sessions
- [ ] Shared antipattern guards extracted from duplicates across 5 type files: Sleep-based waits, hardcoded config, vacuous assertions, source-code-level testing
- [ ] File is completely framework-agnostic — no language-specific code, commands, or grep patterns
- [ ] Frontmatter includes `conventions` field (empty, as _shared.md is universal)

## Hard Rules

- `_shared.md` defines ONLY abstract principles — no type-specific details (no CLI subprocess, no TUI terminal size, no UI selector strategy)
- 5 type files will REFERENCE these principles in their Golden Rules section, not duplicate them

## Implementation Notes

- Structure each principle as: `## Principle: {Name}` → constraint statement → rationale → shared antipattern guard
- The shared antipattern guards here replace overlapping guards in type files (type files keep only type-specific antipatterns)
- Idempotency applies primarily to API + CLI; for UI/TUI/Mobile, define as "repeated interaction should not break subsequent tests"
