# L2 Conventions Batch 2 Audit Report

## Audit Baseline

- **Base commit**: `6bbb2593979f6858cb101b216ed89d5ba8a29f67`
- **Audit date**: 2026-06-03
- **Audit scope**: 10 files in `docs/conventions/`:
  1. naming.md
  2. package-organization.md
  3. prompt-template-hierarchy.md
  4. skill-self-containment.md
  5. skill-structure.md
  6. surface-cli.md
  7. surface-rules.md
  8. testing/index.md
  9. testing/cli/index.md
  10. testing/cli/core.md

## Issue Summary

- P0: 0 | P1: 7 | P2: 6 | P3: 3

## Issue Details

### [P1] naming.md: Non-existent Package `pkg/version`

- **File**: `docs/conventions/naming.md:191`
- **Declaration**: Section 4.3 "pkg/ 各包命名说明" lists `pkg/version` with meaning "版本常量" and justification "标准 Go 命名".
- **Actual**: `forge-cli/pkg/version/` directory does not exist. Version information is defined in `forge-cli/pkg/types/version.go` as `var Version = "dev"` and `func GetVersion() string`. The `version` concept is embedded within `pkg/types`, not a separate package.
- **Suggested action**: Remove `version` row from the package table in Section 4.3. The `pkg/types` entry already exists (line 183) and its description should note that it also contains version information (which it does not currently -- it says "共享类型定义" only). Alternatively, update `pkg/types` row to say "共享类型定义 + 版本信息".

### [P1] naming.md: Non-existent Packages `pkg/lesson` and `pkg/research`

- **File**: `docs/conventions/naming.md:196-197`
- **Declaration**: Section 4.3 lists `pkg/lesson` (meaning "经验教训记录") and `pkg/research` (meaning "深度研究管理") as pkg/ packages.
- **Actual**: Neither `forge-cli/pkg/lesson/` nor `forge-cli/pkg/research/` directories exist. Lesson and research functionality is implemented as top-level command files in `internal/cmd/lesson.go` and `internal/cmd/research.go`, without dedicated `pkg/` packages.
- **Suggested action**: Remove `lesson` and `research` rows from the package table in Section 4.3. These are command-only features without extracted business logic packages.

### [P1] package-organization.md: Stale Subpackage Count and Composition

- **File**: `docs/conventions/package-organization.md:40`
- **Declaration**: "当前状态：15 个顶层命令文件 + 8 个子包（`base`, `docs`, `fact`, `feature`, `forensic`, `prompt`, `task`, `worktree`）"
- **Actual**: There are 18 top-level non-test command files (not 15). There are 8 subpackages but the composition differs: `base`, `fact`, `feature`, `forensic`, `prompt`, `qualitygate`, `task`, `worktree`. The listed `docs` subpackage does not exist at all (D7 deviation describes it as an empty subpackage, but it has been fully removed). The `qualitygate` subpackage exists but is entirely missing from the document.
- **Suggested action**: (1) Update count from "15 个顶层命令文件" to "18 个顶层命令文件". (2) Replace `docs` with `qualitygate` in the subpackage list. (3) Add `qualitygate/` to the subpackage table in Section 2.1. (4) Remove or update D7 deviation entry since `cmd/docs/` no longer exists.

### [P1] package-organization.md: Missing `qualitygate` Subpackage Documentation

- **File**: `docs/conventions/package-organization.md:37-39`
- **Declaration**: Section 2.1 table lists subpackages as `task/`, `worktree/`, `feature/`, `fact/`, `prompt/`, `forensic/`. No mention of `qualitygate/`.
- **Actual**: `forge-cli/internal/cmd/qualitygate/` exists with 5 Go files (quality_gate.go, quality_gate_extract.go, quality_gate_fix_task.go, quality_gate_lifecycle.go, constants.go) plus a test file. It imports from `pkg/feature`, `pkg/forgeconfig`, `pkg/just`, `pkg/project`, `pkg/task`, `pkg/testrunner`, `pkg/types` and `internal/cmd/base`. It does NOT have a `register.go` file, violating the subpackage convention stated in Section 2.2.
- **Suggested action**: (1) Add `qualitygate/` to the subpackage table in Section 2.1. (2) Add a deviation entry for `qualitygate/` missing `register.go`. (3) Document its purpose (quality gate execution and lifecycle management).

### [P1] skill-structure.md: Constraint Violated -- Auxiliary Directories Beyond rules/ and templates/

- **File**: `docs/conventions/skill-structure.md:38`
- **Declaration**: "辅助文件只在 skill 目录内的 rules/ 或 templates/ 子目录中"
- **Actual**: Multiple skills use additional subdirectories beyond `rules/` and `templates/`:
  - `eval/`: has `experts/` and `rubrics/` subdirectories
  - `submit-task/`: has `data/` subdirectory
  - `tech-design/`: has `examples/` subdirectory
  - `write-prd/`: has `examples/` subdirectory
  - `gen-test-scripts/`: has `types/` subdirectory
  - `test-guide/`: has `references/` subdirectory
- **Suggested action**: Update the constraint to acknowledge the actual auxiliary directory patterns. Allow `types/`, `references/`, `examples/`, `data/`, `experts/`, and `rubrics/` as additional valid auxiliary directories, or restructure the affected skills to consolidate content into `rules/` and `templates/`.

### [P1] skill-structure.md: 5 SKILL.md Files Exceed 350-Line Limit

- **File**: `docs/conventions/skill-structure.md:37`
- **Declaration**: "每个 SKILL.md 行数不超过 350 行"
- **Actual**: Five SKILL.md files exceed the 350-line limit:
  - `init-justfile/SKILL.md`: 504 lines (44% over)
  - `write-prd/SKILL.md`: 453 lines (29% over)
  - `tech-design/SKILL.md`: 445 lines (27% over)
  - `gen-journeys/SKILL.md`: 392 lines (12% over)
  - `gen-test-scripts/SKILL.md`: 370 lines (6% over)
- **Suggested action**: For each oversized SKILL.md, extract rule details to `rules/` files and output templates to `templates/` files per the splitting heuristic defined in the same document. Prioritize `init-justfile` (504 lines, most over limit).

### [P1] testing/index.md: Incomplete Surface Coverage -- Only Lists `cli`

- **File**: `docs/conventions/testing/index.md:10-12`
- **Declaration**: The test conventions index table only lists one surface type: `cli` with file location `tests/cli/`.
- **Actual**: The codebase defines 5 surface types in `pkg/types/surface.go`: `web`, `api`, `cli`, `tui`, `mobile`. Corresponding rule files exist in `plugins/forge/skills/init-justfile/rules/surfaces/` and `plugins/forge/skills/run-tests/rules/surfaces/` for all 5 types. The testing conventions should document all surface types with test strategies. Additionally, `tests/cli/` does not exist as a directory -- the actual test directories use journey names (e.g., `tests/surface-aware-recipe-generation/`, `tests/task-lifecycle/`) rather than surface-type-based names.
- **Suggested action**: (1) Add rows for `web`, `api`, `tui`, `mobile` surface types. (2) Update the `cli` file location from `tests/cli/` to reflect the actual journey-based directory pattern (`tests/<journey>/`). (3) Add a column or note explaining that test directories are organized by journey, not by surface type.

### [P2] package-organization.md: D7 Deviation Describes Non-Existent `cmd/docs` Subpackage

- **File**: `docs/conventions/package-organization.md:121`
- **Declaration**: Deviation D7 states "`internal/cmd/docs` 空子包" with current state "cmd/docs/features/ 目录树存在但无文件".
- **Actual**: `forge-cli/internal/cmd/docs/` does not exist at all. The empty directory has been cleaned up since the deviation was recorded.
- **Suggested action**: Remove D7 from the deviation table. The issue has been resolved.

### [P2] prompt-template-hierarchy.md: TASK-CONSTRAINTS Tag Described but Not Used

- **File**: `docs/conventions/prompt-template-hierarchy.md:30-33`
- **Declaration**: Defines `<TASK-CONSTRAINTS>` as a template-level workflow constraint tag, stating it is "用于 test.* 模板".
- **Actual**: A search of `plugins/forge/` for `TASK-CONSTRAINTS` returns zero results. The tag is defined in the convention document but not used in any actual template or skill file. This means either (a) the tag was planned but never implemented, or (b) it was removed from templates after the convention was written.
- **Suggested action**: Either implement `TASK-CONSTRAINTS` in test-related templates as described, or remove the tag definition from the convention document if it is no longer planned.

### [P2] testing/cli/core.md: Test File Location Claims `tests/cli/` but Directory Does Not Exist

- **File**: `docs/conventions/testing/cli/core.md:12`
- **Declaration**: "目录: `tests/cli/` 或 `tests/<journey>/`（当 Journey 仅包含 CLI 测试时）"
- **Actual**: `tests/cli/` directory does not exist. All test directories use journey-based names (e.g., `tests/surface-aware-recipe-generation/`, `tests/task-lifecycle/`, `tests/quality-gate/`). No tests are organized under a `tests/cli/` path.
- **Suggested action**: Update the directory guidance. The primary pattern is `tests/<journey>/`. The `tests/cli/` path appears to be a theoretical organization that was never adopted. Consider removing `tests/cli/` as a valid location or documenting when it should be used vs. journey-based naming.

### [P2] testing/cli/core.md: Build Tag `cli_functional` Used Across Non-CLI Journey Tests

- **File**: `docs/conventions/testing/cli/core.md:14`
- **Declaration**: "Build tag: `//go:build cli_functional`（Go）、`@cli-functional`（BDD tag）" -- implies this tag is CLI-specific.
- **Actual**: The `cli_functional` build tag is used across all journey tests, not just CLI surface tests. For example, `tests/surface-aware-recipe-generation/` tests use `cli_functional` build tags despite testing surface configuration and recipe generation (not CLI-specific behavior). This makes the tag name misleading -- it really means "functional tests that invoke the CLI binary" rather than "tests for the CLI surface type".
- **Suggested action**: Clarify the semantics of the `cli_functional` build tag in the documentation. It should be documented as a general "integration test that spawns CLI subprocesses" tag rather than being tied to the CLI surface type specifically.

### [P2] naming.md: Deviation N3 Describes `runWorktree*` Prefix as Redundant but Code Still Uses It

- **File**: `docs/conventions/naming.md:244-248`
- **Declaration**: Deviation N3 states `internal/cmd/worktree/` uses `runWorktreeRemove`, `runWorktreeList` which is "冗余" and suggests simplifying to `runRemove`, `runList`.
- **Actual**: All worktree run functions still use the `runWorktree*` prefix: `runWorktreePush`, `runWorktreeList`, `runWorktreeResume`, `runWorktreeRemove`, `runWorktreeStatus`, `runWorktreeStart`. The deviation was documented but never addressed. The conclusion says "新代码应在子包中使用 `run<Subcommand>` 模式" but the old code has not been updated.
- **Suggested action**: Either (a) update the worktree run functions to remove the redundant `Worktree` prefix (low-risk rename since they are unexported), or (b) update the deviation entry to mark this as a permanent exception rather than a temporary deviation.

### [P2] surface-cli.md: Minor Error Message Discrepancy in Quote Style

- **File**: `docs/conventions/surface-cli.md:19`
- **Declaration**: Error message: `"no surface configured; run 'forge init' to configure surfaces"` (using single quotes around `forge init`).
- **Actual**: The actual error message in `forge-cli/internal/cmd/surfaces.go` uses backticks: `"no surface configured; run \`forge init\` to configure surfaces"`. This applies to both the JSON error path (`jsonError`) and the plain-text error path (`write`).
- **Suggested action**: Update the convention to use backticks around `forge init` to match the actual code output. This ensures documentation accurately reflects what consumers see.

### [P3] testing/cli/index.md: Minimal Content -- Only Contains Link to core.md

- **File**: `docs/conventions/testing/cli/index.md:1-11`
- **Declaration**: The file serves as an index with a single link to `core.md`.
- **Actual**: The file content is minimal and auto-generated (`<!-- auto-generated by forge:test-guide -->`). It provides no summary of CLI testing conventions, making it a poor navigation entry point. Users must always click through to `core.md`.
- **Suggested action**: Consider adding a brief summary of CLI testing conventions (e.g., isolation model, assertion requirements) to the index page, or document that this is intentionally a redirect-only page.

### [P3] skill-self-containment.md: Very Brief -- Lacks Concrete Examples

- **File**: `docs/conventions/skill-self-containment.md:1-11`
- **Declaration**: States that skills must be logically self-contained and that cross-skill duplication is acceptable.
- **Actual**: The convention document is only 11 lines and provides no concrete examples of what constitutes a violation or how to judge the boundary between self-containment and excessive duplication. This makes it difficult to enforce consistently.
- **Suggested action**: Add concrete examples: (1) an example of acceptable duplication (e.g., two skills both defining how to run `go test` with similar parameters), (2) an example of a violation (e.g., a skill that delegates critical workflow steps to another skill without documenting them), and (3) guidance on when to use `references/` vs. duplicating content.

### [P3] package-organization.md: Top-Level Non-Command Files in cmd/

- **File**: `docs/conventions/package-organization.md:28-30`
- **Declaration**: `internal/cmd/*.go` contains "顶层命令入口（简单命令）". Section 2.3 prohibits business logic in cmd/.
- **Actual**: Some top-level files in `internal/cmd/` are not command entry points but utility/infrastructure files: `output.go` (output formatting -- its own doc comment redirects to `base` sub-package), `styles.go` (color constants for terminal display). These are shared infrastructure files placed at the top level rather than in `base/` or a dedicated utility location.
- **Suggested action**: Either move `output.go` and `styles.go` to `internal/cmd/base/` where they logically belong (the `output.go` doc comment already says to import `base`), or document this pattern as an acceptable exception in Section 2.

## Cross-Layer Influence Items (for L3 Reference)

The following findings from this L2 audit may affect L3 knowledge base entries:

| L2 Finding | Potential L3 Impact |
|-----------|-------------------|
| `pkg/version` does not exist (version is in `pkg/types`) | Any lesson/decision referencing `pkg/version` or version package structure is likely outdated |
| `pkg/lesson` and `pkg/research` do not exist | Any lesson/decision referencing these as `pkg/` packages is outdated; they are `internal/cmd/` files |
| `cmd/docs` subpackage removed (D7 resolved) | Any lesson/decision about docs command restructuring should be checked |
| `tests/cli/` does not exist as a test directory | Any lesson/decision referencing `tests/cli/` path organization is outdated |
| `qualitygate` subpackage undocumented | Any lesson/decision about quality gate architecture may reference incorrect paths |

## Audit Quality Review

- **Sampling ratio**: 10% | **Sampling result**: pass | **Missed items**: 0 | **Extended review**: no
