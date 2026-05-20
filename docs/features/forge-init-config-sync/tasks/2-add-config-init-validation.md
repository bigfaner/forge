---
id: "2"
title: "Add auto.validation prompts to forge config init (stdin)"
priority: "P1"
estimated_time: "30m"
dependencies: []
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 2: Add auto.validation prompts to forge config init (stdin)

## Description

The `forge config init` command (stdin-based interactive) prompts for auto e2e, consolidate-specs, clean-code, and git-push — but NOT for `auto.validation`. Add the missing validation quick/full prompts in the auto-behavior section.

## Reference Files
- `docs/proposals/forge-init-config-sync/proposal.md` — Source proposal
- `forge-cli/internal/cmd/config.go` — Auto-behavior prompts at lines 106-127, insert between cleanCode and gitPush
- `forge-cli/pkg/forgeconfig/config.go` — `AutoConfig.Validation` field at line 35

## Acceptance Criteria
- [ ] `runConfigInit()` prompts for `auto.validation` quick (y/N) and full (y/N) using the same stdin prompt pattern as other auto fields
- [ ] Validation prompts appear between cleanCode and gitPush prompts
- [ ] Config struct construction includes the validation values
- [ ] Existing tests in `config_test.go` still pass (update if needed)

## Hard Rules
- Follow existing stdin prompt pattern exactly (fmt.Print + reader.ReadString)
- Do NOT modify `forge init` (that's Task 1)

## Implementation Notes
- See lines 120-126 for the cleanCode/gitPush prompt pattern. The validation prompts should follow the same structure.
- Worktree prompts already exist in `forge config init` — no change needed there.
