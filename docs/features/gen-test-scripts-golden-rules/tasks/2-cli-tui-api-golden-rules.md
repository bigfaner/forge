---
id: "2"
title: "Restructure cli.md, tui.md, api.md with golden rules"
priority: "P0"
estimated_time: "1-2h"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 2: Restructure cli.md, tui.md, api.md with golden rules

## Description

Restructure 3 type files (cli.md, tui.md, api.md) into Golden Rules + Reconnaissance Hints dual-zone structure. Add missing golden rules identified by expert evaluation. Fix Output path and step number references.

## Reference Files
- `docs/proposals/gen-test-scripts-golden-rules/proposal.md` — Source proposal
- `plugins/forge/skills/gen-test-scripts/types/_shared.md` — Cross-type shared principles (created in task 1)
- `plugins/forge/skills/gen-test-scripts/types/cli.md` — Current CLI type
- `plugins/forge/skills/gen-test-scripts/types/tui.md` — Current TUI type
- `plugins/forge/skills/gen-test-scripts/types/api.md` — Current API type
- `plugins/forge/skills/gen-test-scripts/SKILL.md` — For step number reference alignment

## Affected Files

### Create
| File | Description |
|------|-------------|
| — | — |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-scripts/types/cli.md` | Restructure + add golden rules + fix references |
| `plugins/forge/skills/gen-test-scripts/types/tui.md` | Restructure + add golden rules + fix references |
| `plugins/forge/skills/gen-test-scripts/types/api.md` | Restructure + add golden rules + fix references |

### Delete
| File | Reason |
|------|--------|
| — | — |

## Acceptance Criteria

### Structural (all 3 files)
- [ ] Each file has `## Golden Rules` section (framework-agnostic constraints, no language-specific code or commands)
- [ ] Each file has `## Reconnaissance Hints` section with `<!-- Discovery hints — convert findings to Fact Table values, do not use for generation instructions -->` annotation
- [ ] Golden Rules reference `_shared.md` principles (e.g., "per _shared.md: Timeout Protection") instead of redefining them
- [ ] Shared antipattern guards (Sleep, hardcoded config, vacuous assertions) removed — replaced by `_shared.md` reference
- [ ] Type-specific antipattern guards retained in Golden Rules section
- [ ] Output path changed from `tests/e2e/features/<feature>/` to `tests/<journey>/`
- [ ] All "SKILL.md Step 1.5" references corrected to match actual SKILL.md step numbers

### CLI-specific additions
- [ ] Timeout Protection: two-level timeout — (1) test function-level timeout from `_shared.md`, (2) process-level timeout guard (subprocess must exit within N seconds, SIGKILL on timeout, clean up process tree)
- [ ] Binary Isolation: tests MUST compile a dedicated binary in TestMain/setup, never use `go run` or PATH resolution
- [ ] Environment Hermeticity: explicit env inheritance + override pattern (`cmd.Env = append(os.Environ(), ...)`) or `t.Setenv()`; never rely on host environment

### TUI-specific additions
- [ ] Terminal Size Contract: set `TERM=dumb` or fixed `LINES`/`COLUMNS` env vars, eliminate rendering variance
- [ ] ANSI Sanitization: assertions must strip ANSI escape sequences before matching, or use dedicated terminal output parsing
- [ ] Stable State Detection: define observable signals for "screen rendering complete" (stdout output stable for N ms, process exited, specific marker appeared) — not time-based assumptions

### API-specific additions
- [ ] Idempotency Check: for PUT/DELETE endpoints, verify repeated requests produce identical results
- [ ] Request Timeout: HTTP client must set connection + read timeouts, prevent test hangs
- [ ] Content-Type Verification: requests must declare Accept header, responses must verify Content-Type header

## Hard Rules

- Golden Rules section must contain ZERO language-specific code, import paths, or grep commands
- Reconnaissance Hints grep commands are acceptable but must be annotated as discovery-only
- Each type file frontmatter `conventions` field lists the Convention file(s) it expects (keep existing values)

## Implementation Notes

- cli.md, tui.md, api.md are structurally similar — restructure them consistently
- The existing Classification Indicators, Fact Table Required Keys, and Verification Method sections can remain as-is (they're already well-structured)
- Keep existing Generation Patterns but move framework-specific parts to Reconnaissance Hints or mark as "Convention-driven"
- The Key Sequence Encoding table in tui.md (escape codes like `\x1b[A`) should stay in Golden Rules as it defines WHAT to encode, not HOW
