---
id: "2"
title: "Clean documentation references to forge test commands"
priority: "P1"
estimated_time: "45m"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 2: Clean documentation references to forge test commands

## Description
Remove or update all documentation references to `forge test` commands (promote, run-journey, verify) across README, CLI reference, architecture docs, conventions, skill docs, and command docs. This ensures agents and users reading these files won't attempt to invoke nonexistent commands.

## Reference Files
- `proposal.md#Scope` — lists every doc file containing `forge test` references
- `proposal.md#Success-Criteria` — zero-residue search verification
- `proposal.md#Key-Risks` — risk of agents reading stale docs and invoking removed commands

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `README.md` | Remove `forge test` pipeline references (lines 85, 110, 314) |
| `forge-cli/docs/OVERVIEW.md` | Remove `forge test` section (lines 83-105) and promote/CI references |
| `forge-cli/docs/OVERVIEW.zh.md` | Remove `forge test` section (lines 90-92) |
| `docs/ARCHITECTURE.md` | Remove `forge test promote` references (lines 286, 304, 531, 542) |
| `docs/conventions/forge-cli-reference.md` | Remove `forge test` command table (lines 51-59) and `forge e2e → forge test` mappings (lines 149, 151) |
| `docs/conventions/forge-distribution.md` | Remove `forge test promote` references (lines 112, 119, 179) |
| `docs/profile-authoring.md` | Remove stale `forge testing` reference (line 30) |
| `plugins/forge/skills/run-tests/SKILL.md` | Remove `forge test promote` references (lines 265, 267) |
| `plugins/forge/skills/consolidate-specs/SKILL.md` | Remove `forge test promote` reference (line 267) |
| `plugins/forge/skills/gen-contracts/rules/journey-contract-model.md` | Update promote/stage/graduate references (lines 138, 193, 194) |
| `plugins/forge/skills/gen-journeys/rules/journey-contract-model.md` | Update promote/stage/graduate references (lines 138, 193, 194) |
| `plugins/forge/commands/run-tasks.md` | Remove `forge test promote` suggestion (line 119) |

## Acceptance Criteria
- [ ] Full-text search for `forge test promote`, `forge test run-journey`, `forge test verify` returns zero results (excluding `docs/features/` history docs)
- [ ] No documentation file instructs users or agents to run `forge test` subcommands

## Hard Rules
- **DO NOT modify** files under `docs/features/` — they are historical design records.
- When removing `forge test` references from docs, ensure surrounding text remains coherent. Remove entire paragraphs/sections if they only describe `forge test` commands; otherwise update only the specific references.

## Implementation Notes
- Some docs describe the test promotion model conceptually (e.g., `@feature → @regression` tag lifecycle). The model itself remains valid — only remove references to the CLI commands that implement it. Skill-layer promotion (via `run-tests`) is unaffected.
- `journey-contract-model.md` exists in both `gen-contracts/rules/` and `gen-journeys/rules/` — both copies need updating.
