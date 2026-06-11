# Existing Code Modification Split Rules

**Load condition**: load this file IF the tech-design references modifications to existing shared code (interfaces, models, API contracts, utility functions).

**Guard clause**: if the tech-design references only new code (new files, new interfaces, new endpoints), skip this rule. Purely additive changes do not need splitting.

## When to Apply

Inspect the tech-design for changes to artifacts that already exist in the codebase. Apply this split when:

- The change propagates to **>5 downstream files**, OR
- The change **spans multiple architectural layers** (e.g., repository -> service -> handler)

If neither threshold is met, create a single task as normal.

## Split Procedure

When the thresholds are met, split the task into two sub-tasks by dependency layers. Each sub-task must be independently compilable and testable.

### Sub-task A: Shared Artifact Update

- **Sub-ID**: `<seq>.<sub>a`
- **`breaking: true`** — always, because the shared artifact's contract changes
- **Surface**: resolve `surface-key` and `surface-type` via `forge surfaces --json` on affected files (see Step 4a Surface-Key/Type Assignment)
- **Content**: Apply changes to the shared artifact, reconcile ALL downstream consumers so existing code compiles and tests pass. No new business logic — only signature changes + stubs/adapters.
- **Dependencies**: same as the original unsplit task

### Sub-task B: Feature Implementation

- **Sub-ID**: `<seq>.<sub>b`
- **`breaking`**: determined by the feature's own changes (not the shared artifact update)
- **Surface**: resolve `surface-key` and `surface-type` via `forge surfaces --json` on new/modified feature files (see Step 4a Surface-Key/Type Assignment)
- **Content**: Implement the actual feature logic using the updated shared artifact. Standard acceptance criteria from the design.
- **Dependencies**: depends on `<seq>.<sub>a`

## Exclusion

Purely additive new code (new files, new interfaces that nothing yet implements) does not need splitting. This split only applies when an existing shared artifact with downstream consumers is being modified.

## Maintenance Note

This rule file depends on the following sections in the skeleton SKILL.md:

- **Step 4a: Create Task Files** — task file creation, sub-ID conventions, breaking classification
- **Surface-Key/Type Assignment** — surface resolution for each sub-task

If either of these sections changes in the skeleton, verify that the split procedure and sub-ID conventions in this file remain consistent.
