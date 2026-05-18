---
feature: "remove-references-dir"
status: tasks
mode: quick
---

# Feature (Quick): remove-references-dir

<!-- Status flow: tasks -> in-progress -> completed -->

Remove the non-standard `plugins/forge/references/` directory by inlining all shared reference content directly into consuming skills/commands and relocating CLI-specific files to the CLI repository. Six tasks cover the full migration, ending with directory deletion and documentation update.

## Documents

| Document | Path |
|----------|------|
| Proposal | ../../proposals/remove-references-dir/proposal.md |

## Tasks

| ID | Title | Status | File |
|----|-------|--------|------|
| 1 | Inline decision-logging.md into consuming skills | pending | tasks/1-inline-decision-logging.md |
| 2 | Inline knowledge-extraction.md into consuming skills and commands | pending | tasks/2-inline-knowledge-extraction.md |
| 3 | Inline step0-profile-resolution, type-assignment, and intent-propagation into consuming skills | pending | tasks/3-inline-task-type-refs.md |
| 4 | Inline config.yaml and sitemap.json examples into gen-sitemap command | pending | tasks/4-inline-gen-sitemap-examples.md |
| 5 | Move forge-config schema and example YAML to CLI, update test paths | pending | tasks/5-relocate-cli-schema.md |
| 6 | Update forge-distribution.md and remove references/ directory | pending | tasks/6-update-docs-remove-dir.md |
