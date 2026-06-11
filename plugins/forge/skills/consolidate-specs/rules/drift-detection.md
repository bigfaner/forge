# Drift Detection and Auto-Fix Rules

## Step 9: Detect Drift in Project-Level Specs

Read all project-level spec files and validate each rule against the current codebase:

1. **Read all spec files**:
   - `docs/business-rules/*.md` -- all business rule files
   - `docs/conventions/*.md` -- all technical convention files

2. **Validate each rule against code**: For each rule in every spec file, search the codebase for the keywords and patterns described in the rule. Compare the rule's stated behavior against the actual code implementation:
   - Extract key domain terms, function names, file paths, or behavior descriptions from each rule
   - Search the relevant source files for those terms
   - Determine whether the rule's description still matches the code's current behavior

3. **Classify each rule**:
   - `current` -- rule description matches current code behavior
   - `drifted` -- rule description is partially or fully inconsistent with current code (e.g., renamed function, changed threshold, modified behavior)
   - `orphaned` -- the code the rule describes no longer exists (e.g., deleted module, removed feature)

4. **Output drift report**: Write a summary of all classifications. If no drift is found (all `current`), skip Steps 10-11 and proceed to Step 12 (vocabulary generation).

### Drift-Only Mode Entry

If running in drift-only mode (no PRD/design files exist), start here at Step 9. Skip Steps 1-8 entirely.

## Step 10: Auto-Fix Drifted Specs

For each rule classified as `drifted` or `orphaned` in Step 9:

1. **Drifted rules**: Update the rule's description/behavior text in-place to match the current code. Preserve the project-global ID (e.g., `BIZ-auth-001`) -- only update the descriptive text, not the ID or structural format.

2. **Orphaned rules**: Remove the rule entry from the spec file. Record the deletion for the commit message in Step 11:
   - Rule ID (e.g., `TECH-api-002`)
   - Reason for deletion (e.g., "corresponding module `X` removed in commit abc1234")

3. **Detect implicit new rules**: While scanning the code for drift, if you discover new patterns, conventions, or business logic that should be documented at the project level but are not in any spec file:
   - Extract the candidate rule with `[CROSS]` classification
   - **Interactive mode**: Present to user for confirmation before appending
   - **Non-interactive mode**: Auto-append with `[auto-specs]` tag -- include in commit message
   - Append confirmed rules to the appropriate spec file with a new project-global ID

4. **Re-derive `domains` frontmatter**: When a file's content changes substantially (rules updated, added, or removed), re-derive the `domains` field per `rules/domain-frontmatter.md`. Compare the new domain set against the existing one:
   - If domains have changed, update the frontmatter in-place
   - If the updated `domains` cause a new >50% overlap with another file's domains, flag in the commit message and notify the user

### Preservation Rules

- Project-global IDs must never change during auto-fix -- only description and behavior text updates
- File structure and formatting must remain consistent with the existing spec file conventions
- Deleted rules must be recorded with ID and reason for traceability in git history
