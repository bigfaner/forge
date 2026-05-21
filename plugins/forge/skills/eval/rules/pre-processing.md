# Pre-Processing by Type

| Type | Before Scoring |
|------|---------------|
| **All types** | If rubric has `context` frontmatter, load filtered context files: (1) for each string in `conventions`, glob `docs/conventions/<string>*.md` and read matching files; (2) if `business-rules: auto`, glob `docs/business-rules/*.md` and read all, else read listed filenames. Concatenate into `CONTEXT_CONTENT` for Step 2 injection. Skip missing files silently (no error, no abort). |
| `harness` | Gather project context, write snapshot. Scorer evaluates snapshot, not raw files. |
| `consistency` | Assemble document bundle -- copy relevant docs into flat directory for scorer. |
| `test-cases` | Load Convention files from `docs/conventions/` by `domains` frontmatter (match `testing`, `go`, `typescript`, etc.), extract language from `Framework` section. Fallback: scan source/test files (`go.mod`, `package.json`, `*_test.go`, etc.). Pass interfaces to scorer. |
| `ui-test-cases`, `tui-test-cases`, `mobile-test-cases`, `api-test-cases`, `cli-test-cases` | Load Convention files from `docs/conventions/` by `domains` frontmatter (match `testing`, `go`, `typescript`, etc.), extract language from `Framework` section. Fallback: scan source/test files (`go.mod`, `package.json`, `*_test.go`, etc.). Pass interfaces to scorer. |
| `prd` | Detect mode: `prd-ui-functions.md` exists -> Mode A (with UI), else Mode B (no UI). |
| `validate-code` | 1) Read PRD -> extract user scenarios list (from prd-spec.md flow descriptions and prd-user-stories.md acceptance criteria). 2) Run `git diff <base-branch>...HEAD` to get changed files and diff hunks. 3) Compile changed file list. 4) Pass PRD scenarios + diff + file list to scorer as assembled input. |
| `validate-ux` | **Two-phase pre-processing** (must execute in git worktree or temp dir). Full sub-pipeline: `rubrics/validate-ux-pipeline.md`. |
