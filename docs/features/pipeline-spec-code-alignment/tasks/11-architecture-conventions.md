---
id: "11"
title: "Fix architecture and conventions docs"
priority: "P1"
estimated_time: "1.5h"
dependencies: [10]
type: "doc"
mainSession: false
---

# 11: Fix architecture and conventions docs

## Description

Fix architecture documentation, convention files, and skill references that describe non-existent systems or use incorrect terminology. Covers Cluster 7 doc subset (issues G2-G7 from Evidence section):

1. **ARCHITECTURE.md** (~lines 119, 176-212): Describes doc-scorer/doc-reviser as independent agents, but they're protocol files within the eval skill. Only `agents/task-executor.md` exists. Fix: (a) remove fictitious agent descriptions, (b) describe eval skill's internal protocol correctly.

2. **ARCHITECTURE.md** (~lines 244-256): Scope resolution algorithm describes `just project-type` detection flow, but actual `ResolveScope` uses `just --dry-run compile <scope>` probing. Fix: align algorithm description with actual code.

3. **ARCHITECTURE.md** (~lines 148, 450): Duplicate spelling `forge forge task claim`. Fix to single `forge`.

4. **dispatcher-quality.md** (~lines 12, 34): References `go build ./...` and `go test ./...` (Go-specific), but actual code uses `just compile`/`just test` abstraction layer. Also only mentions `coding.fix`, missing `coding.cleanup` (for fmt/lint failures). Fix: use just abstractions, add `coding.cleanup`.

5. **gen-contracts/SKILL.md** (lines 61, 64, 70) + **rules/validation.md** (line 36): Uses "interfaces" terminology and `interfaces` config field. Actual config field is `surfaces`. Fix: replace "interfaces" → "surfaces".

6. **clean-code/SKILL.md** (~line 162): References `just test` but standard recipe is `just unit-test`. Fix.

7. **fix-bug.md** (lines 141, 143, 146, 179): References `just test <slug>` which doesn't exist as a target. Fix to correct target.

8. **execute-task.md** (frontmatter line 3): Description says "focused TDD workflow" but actual behavior is claim/dispatch/verify orchestrator. TDD logic is in task-executor. Fix description.

## Reference Files
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Problem` — Evidence G2 (gen-contracts interfaces), G3 (fictitious agents), G4 (algorithm mismatch), G5 (go build vs just compile), G6 (surface-rules scope), G7 (dispatcher missing coding.cleanup)
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Proposed-Solution` — Cluster 7 doc subset
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Success-Criteria` — SC for ARCHITECTURE.md, dispatcher-quality.md, just references

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `docs/ARCHITECTURE.md` | Remove fictitious agents, fix algorithm, fix `forge forge` typo |
| `docs/conventions/dispatcher-quality.md` | `go build`/`go test` → `just compile`/`just test`, add `coding.cleanup` |
| `plugins/forge/skills/gen-contracts/SKILL.md` | "interfaces" → "surfaces" |
| `plugins/forge/skills/gen-contracts/rules/validation.md` | "interfaces" → "surfaces" |
| `plugins/forge/skills/clean-code/SKILL.md` | `just test` → `just unit-test` |
| `plugins/forge/commands/fix-bug.md` | `just test <slug>` → correct target |
| `plugins/forge/commands/execute-task.md` | Frontmatter description: "TDD workflow" → "claim/dispatch/verify" |

## Acceptance Criteria
- [ ] ARCHITECTURE.md does not describe doc-scorer/doc-reviser as independent agents
- [ ] ARCHITECTURE.md scope resolution algorithm matches actual Go code
- [ ] No `forge forge` duplicate in ARCHITECTURE.md
- [ ] dispatcher-quality.md uses `just` abstractions and mentions `coding.cleanup`
- [ ] gen-contracts docs use "surfaces" terminology (not "interfaces")
- [ ] clean-code/SKILL.md references `just unit-test`
- [ ] fix-bug.md references correct just target for running tests
- [ ] execute-task.md frontmatter description is accurate

## Hard Rules
- Do not add new architecture descriptions — only fix existing incorrect ones
- Preserve document structure and heading hierarchy

## Implementation Notes
- For ARCHITECTURE.md: read the actual Go code for ResolveScope before rewriting the algorithm description
- For fix-bug.md just targets: check the project's justfile for available test targets
