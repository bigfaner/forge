---
id: "3"
title: "Update forge init to stop generating claude justfile recipes"
priority: "P1"
estimated_time: "30min"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 3: Update forge init to stop generating claude justfile recipes

## Description

`forge init` currently appends `claude` and `claude-c` recipes to the project justfile. Since `forge claude` now provides this functionality, remove the recipe generation from `forge init`.

## Reference Files
- `docs/proposals/migrate-claude-commands/proposal.md` — Source proposal
- `forge-cli/internal/cmd/init.go` — Init command with justfile recipe generation (lines 48-54, 186-266)

## Acceptance Criteria

- [ ] `justfileRecipes` slice in init.go is empty or removed
- [ ] `forge init` no longer appends claude-related recipes to justfile
- [ ] Init command description updated to remove claude/claude-c recipe mention
- [ ] Existing tests for `forge init` updated (no longer expect claude recipes in justfile)
- [ ] `buildJustfileAppend` and related helper functions cleaned up if no longer needed

## Hard Rules

- Remove all claude-related recipe generation code, not just the entries
- If `justfileRecipes` becomes empty, remove the entire justfile update step from init flow

## Implementation Notes

- The `justfileRecipes` var (lines 48-54) contains the two recipes
- The `updateJustfile`, `buildJustfileAppend`, `recipeExists`, `collectAddedRecipeNames` functions (lines 186-266) handle recipe appending
- If no recipes remain, these functions and the entire "Step 4: Update justfile" block can be removed
- Update `initCmd.Long` description (line 22-26) to remove "appends claude/claude-c recipes to justfile"
- Check `init_test.go` for tests that assert claude recipes in generated justfile content
