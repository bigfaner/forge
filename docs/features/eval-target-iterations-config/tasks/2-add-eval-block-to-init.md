---
id: "2"
title: "Add eval block to forge config init output"
priority: "P1"
estimated_time: "1.5h"
complexity: "medium"
dependencies: [1]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 2: Add eval block to forge config init output

## Description

When `forge config init` (or `forge init`) generates config.yaml, automatically populate the `eval` block with default values read from rubric frontmatter. Do NOT add interactive prompts for target/iterations — the values are seeded from rubrics for user convenience (edit-in-place), not interactive questions.

Rubric default values: proposal 900/3, prd 900/3, design 900/3, ui 950/3, journey 850/3, contract 850/3, consistency 900/3.

## Reference Files
- `docs/proposals/eval-target-iterations-config/proposal.md` — Proposed Solution (init behavior), Scope > In Scope (item 5)
- `forge-cli/internal/cmd/init_config.go` — runConfigInitIfNeeded (line 233), Config construction at line 274
- `plugins/forge/skills/eval/rubrics/<type>.md` — rubric frontmatter with target/iterations defaults

## Acceptance Criteria
- [ ] `forge config init` generates config.yaml with complete `eval` block containing all 7 types (proposal, prd, design, ui, journey, contract, consistency)
- [ ] Default values match rubric frontmatter: proposal 900/3, prd 900/3, design 900/3, ui 950/3, journey 850/3, contract 850/3, consistency 900/3
- [ ] Generated config.yaml is valid YAML parseable by the `Config` struct (no serialization errors)

## Implementation Notes
- `runConfigInitIfNeeded` constructs `Config{Auto: auto, Worktree: worktree}` at line 274. Add `Eval` field populated with rubric defaults.
- Rubric file paths: `plugins/forge/skills/eval/rubrics/<type>.md`. Init needs to parse YAML frontmatter to extract `target` and `iterations`.
- Alternative: embed default values as Go constants (avoids file I/O during init). This is simpler and the values are stable.
- The `forge config init` TUI flow (`forge config init` command) also calls `runConfigInitIfNeeded` — both code paths must include the eval block.
