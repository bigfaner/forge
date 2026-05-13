---
date: "2026-05-13"
doc_dir: "docs/features/forge-cli-v3/design/"
iteration: 1
target_score: 900
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration 1

**Score: 763/1000** (target: 900)

```
┌─────────────────────────────────────────────────────────────────┐
│                     DESIGN QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Architecture Clarity      │  170     │  200     │ ⚠️         │
│    Layer placement explicit  │  65/70   │          │            │
│    Component diagram present │  55/70   │          │            │
│    Dependencies listed       │  50/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Interface & Model Defs    │  133     │  200     │ ⚠️         │
│    Interface signatures typed│  48/70   │          │            │
│    Models concrete           │  50/70   │          │            │
│    Directly implementable    │  35/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Error Handling            │  115     │  150     │ ⚠️         │
│    Error types defined       │  35/50   │          │            │
│    Propagation strategy clear│  45/50   │          │            │
│    HTTP status codes mapped  │  35/50   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Testing Strategy          │  115     │  150     │ ⚠️         │
│    Per-layer test plan       │  35/50   │          │            │
│    Coverage target numeric   │  45/50   │          │            │
│    Test tooling named        │  35/50   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Breakdown-Readiness ★     │  130     │  200     │ ❌         │
│    Components enumerable     │  55/70   │          │            │
│    Tasks derivable           │  45/70   │          │            │
│    PRD AC coverage           │  30/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Security Considerations   │  100     │  100     │ N/A        │
│    Threat model present      │  N/A/50  │          │            │
│    Mitigations concrete      │  N/A/50  │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  763     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

★ Breakdown-Readiness < 180/200 blocks progression to `/breakdown-tasks`

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Section 2 (Interfaces) | forensic, profile, cleanup, feature commands have no interface signatures — described as "Unchanged" without showing current signature | -22 pts from Interface signatures typed |
| Section 2 (Interfaces) | `addFixTask` described as "Signature mirrors current all_completed.go" — prose reference instead of typed signature | included in above |
| Section 2 (Models) | Only 3 model items listed under "Key Go structs affected" — most renamed commands' affected structs not enumerated | -20 pts from Models concrete |
| Section 2 (Directly implementable) | Phase 4 reference updates (hooks, 23 skills, docs, tests) have no per-file change spec; developer must grep and guess | -25 pts from Directly implementable |
| Section 3 (Error Handling) | No Go error types defined for new error cases — only exit codes and stderr strings specified | -15 pts from Error types defined |
| Section 3 (Error Handling) | HTTP status codes criterion not applicable to CLI; exit codes mapped instead — partial credit | -15 pts from HTTP status codes mapped |
| Section 4 (Testing) | No test plan for Phase 4 reference updates; no specific mock framework named | -15 pts from Per-layer test plan |
| Section 4 (Testing) | Mock framework not named — "ExecRunner interface mock" describes pattern, not tooling | -15 pts from Test tooling named |
| Section 5 (Breakdown-Readiness) | PRD AC S3 "concurrent write conflict" flagged as "PRD amendment needed" but not addressed — unaddressed AC gap | -30 pts from PRD AC coverage |
| Section 5 (Breakdown-Readiness) | 23 skills and 6 command files not individually enumerated — cannot derive specific tasks | -25 pts from Tasks derivable |
| Section 5 (Breakdown-Readiness) | PRD Goal "AI agent command selection >= 9/10" has no corresponding design verification mechanism | -10 pts from PRD AC coverage |
| Section 1 (Architecture) | Component diagram ASCII alignment is messy; pkg layer boxes overflow and are hard to read | -15 pts from Component diagram |
| Section 1 (Architecture) | External dependency versions not pinned (cobra, yaml.v3) | -10 pts from Dependencies listed |

---

## Attack Points

### Attack 1: [Breakdown-Readiness — PRD AC coverage gap with unaddressed acceptance criterion]

**Where**: "S3: concurrent write conflict | **PRD amendment needed**: remove from scope | Pre-existing gap predating this feature (no locking mechanism exists in codebase). Design recommends PRD amendment to drop this AC, since adding file-locking is an orthogonal concern requiring its own design."
**Why it's weak**: The design does not have the authority to declare a PRD AC out of scope. If the PRD says "concurrent write conflict" is an acceptance criterion, the design must either address it or produce evidence that the PRD has been formally amended *before* this design is evaluated. Recommending an amendment is not the same as having one. This leaves a gap that blocks downstream task derivation.
**What must improve**: Either (a) add a concrete design for concurrent write conflict handling, or (b) attach proof that the PRD was amended (e.g., PRD commit hash or updated PRD file with the AC struck through). Until one of these exists, the AC is unaddressed.

### Attack 2: [Interface & Model Definitions — missing signatures for majority of existing commands]

**Where**: "forensic.go | Unchanged, already has subcommands" (line 574), "profile.go | Unchanged" (line 577), "cleanup.go | Unchanged" (line 564), "feature | 设置/显示当前 feature 上下文" (PRD line 228).
**Why it's weak**: The document shows detailed Cobra command signatures only for new commands (list-types, probe, e2e subcommands, quality-gate cap). For the renamed commands (submit, check-deps, validate-index, verify-task-done) and unchanged commands (forensic, profile, cleanup, feature), the design provides no interface signatures at all. A developer implementing `forge task submit` cannot determine the exact flag set, argument validation, or return behavior from this document alone — they must read the existing `record.go` source. The "Unchanged" label is insufficient for breakdown-readiness because tasks must specify what to change AND what to preserve.
**What must improve**: Add a "Renamed Commands Interface Table" that lists each renamed command's full Cobra struct (Use, Short, Args, Flags, Run function name). For unchanged commands, add a minimal "Existing Interface Reference" section with file paths and key struct fields so that tasks can reference them without ambiguity.

### Attack 3: [Breakdown-Readiness — Phase 4 reference updates lack per-file enumeration]

**Where**: "Update 23 skills" (line 169), "External consumers: hooks.json, 23 skill files, 2 agent files, 6 command files, 2 doc files" (line 22).
**Why it's weak**: Phase 4 is the largest bulk-work phase. The design lists categories of files to update but does not enumerate them. A task breakdown cannot derive "update skill X" tasks without knowing which 23 skills exist and what command references each one contains. The PRD Scope explicitly lists "更新 23 个 skills 中的命令引用" as an in-scope item, but the design provides no mapping from skill file to old/new command references. This means the `/breakdown-tasks` step would need to do its own codebase analysis, defeating the purpose of a complete design.
**What must improve**: Add a Phase 4 enumeration table listing every file to be updated with columns: File Path | Old Reference | New Reference. At minimum, enumerate the 23 skill file names and the specific old/new command strings to replace. For hooks.json and doc files, show the exact line-level changes.

---

## Previous Issues Check

<!-- Only for iteration > 1 — not applicable for iteration 1 -->

---

## Verdict

- **Score**: 763/1000
- **Target**: 900/1000
- **Gap**: 137 points
- **Breakdown-Readiness**: 130/200 — **cannot proceed to /breakdown-tasks** (below 180 threshold)
- **Action**: Continue to iteration 2. Priority fixes: (1) Address or formally amend PRD S3 concurrent write conflict AC, (2) Add interface signatures for ALL renamed/unchanged commands, (3) Enumerate Phase 4 files with per-file change specs.
