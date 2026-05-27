Detect spec drift between existing project specs and current code for the {{.FeatureSlug}} feature.

## Discovery Strategy
1. Run `git diff --name-only main...HEAD` to identify files changed by this feature
2. List all spec files in docs/business-rules/ and docs/conventions/
3. For each spec file, read its `domains` frontmatter
4. Only verify specs whose domains overlap with the changed files
5. Skip specs with no overlap — they are unaffected by this feature

Do NOT scan all spec files blindly. Use git diff to narrow scope first.
If git diff returns no changes, skip — nothing to drift against.

Auto-fix drifted specs and commit with [auto-specs] tag.
