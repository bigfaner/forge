---
created: "2026-06-07"
tags: [testing, local-dev-deployment]
---

# forge task validate is the correct command, not validate-index

## Problem
During breakdown-tasks, explored `forge task validate-index --help` to validate index.json. The command `validate-index` does not exist as a forge CLI subcommand — it returned the parent `forge task` help output instead, making it unclear how to validate the generated index.

## Root Cause
1. The Forge Guide (in CLAUDE.md system context) references `forge task validate-index <path>` as a CLI command
2. But the actual forge CLI binary only exposes `forge task validate [file]` as the validation subcommand
3. `validate-index` is either a guide documentation artifact or a deprecated alias — running it falls through to the parent `forge task` help with no error
4. This caused confusion about which command to use for index validation in Step 7 of breakdown-tasks

## Solution
Use `forge task validate [file]` (not `validate-index`). It accepts an optional file path argument:
```bash
forge task validate docs/features/<slug>/tasks/index.json
```
When no file is specified, it validates the current feature's index.json.

## Reusable Pattern
When a forge CLI command doesn't behave as documented in the guide:
1. Run `forge task --help` to see actual available subcommands
2. The guide may reference commands that are consolidated or renamed — trust `--help` output over guide text
3. `forge task validate` handles both structural validation and AC count validation (1-6 per task)

## Related Files
- Forge Guide (system prompt): references `validate-index` 
- CLI binary: only exposes `validate`
