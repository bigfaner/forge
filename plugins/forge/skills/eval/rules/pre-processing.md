# Pre-Processing by Type

| Type | Before Scoring |
|------|---------------|
| **All types** | If rubric has `context` frontmatter, load filtered context files: (1) for each string in `conventions`, glob `docs/conventions/<string>*.md` and read matching files; (2) if `business-rules: auto`, glob `docs/business-rules/*.md` and read all, else read listed filenames. Concatenate into `CONTEXT_CONTENT` for Step 2 injection. Skip missing files silently (no error, no abort). |
| `harness` | Gather project context, write snapshot. Scorer evaluates snapshot, not raw files. |
| `consistency` | Assemble document bundle -- copy relevant docs into flat directory for scorer. |
| `prd` | Detect mode: `prd-ui-functions.md` exists -> Mode A (with UI), else Mode B (no UI). |
| `validate-code` | 1) Read PRD -> extract user scenarios list (from prd-spec.md flow descriptions and prd-user-stories.md acceptance criteria). 2) Run `git diff <base-branch>...HEAD` to get changed files and diff hunks. 3) Compile changed file list. 4) Pass PRD scenarios + diff + file list to scorer as assembled input. |
| `validate-ux` | **Two-phase pre-processing** (must execute in git worktree or temp dir). Full sub-pipeline: `rubrics/validate-ux-pipeline.md`. |
| `journey` | Detect surface type from `.forge/config.yaml` `surface` field. Load the corresponding surface rule file (`skills/gen-journeys/rules/surface-<type>.md`) as additional context for the scorer. |
| `contract` | Detect surface type from `.forge/config.yaml` `surface` field. Load the corresponding surface rule file (`skills/gen-journeys/rules/surface-<type>.md`) as additional context for the scorer. |
