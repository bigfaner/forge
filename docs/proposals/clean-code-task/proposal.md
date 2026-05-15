---
created: 2026-05-14
author: "faner"
status: Superseded
superseded-by: reject-clean-code-task
---

# Proposal: Add Clean Code Task to Forge Task Pipeline

## Problem

Forge 的任务流水线在实现（business tasks）和测试（test pipeline tasks）完成后，直接进入 consolidate-specs（从代码中提取知识）。缺少一个自动化代码清理环节，导致沉淀到 `docs/` 的 specs 基于未清理的代码提取，且最终代码质量依赖每个 task executor 的个人判断。

### Evidence

- Quick 模式（1-10 tasks）和 Breakdown 模式均无 post-implementation cleanup 步骤
- 多 task 并行实现后，常见累积问题：dead imports、commented-out code、重复逻辑、命名不一致
- **实测数据**：在 `migrate-e2e-to-go` feature（9 个 business tasks，subagent 并行执行）完成后，`git diff --stat` 显示 37 个变更文件中有 5 个包含 unused imports、2 个有 commented-out debug code、1 个存在跨文件重复的 helper function。这些均在后续手动 `/simplify` 时被发现，但若自动清理则可省去人工审查
- consolidate-specs 从代码中提取业务规则和技术规范——如果代码未清理，提取的质量也受影响

### Urgency

当前无清理机制意味着代码质量完全依赖实现阶段的自律。随着 forge pipeline 自动化程度提高，缺乏自动清理环节成为质量闭环的缺口。

## Proposed Solution

在 test pipeline 的 verify-regression 和 consolidate-specs 之间，插入一个新的 `clean-code` 任务类型。该任务由 `forge task index` 自动生成，审查本 feature 所有已变更文件并执行代码清理。

### User-Facing Behavior

1. **Task appearance**: After all business and test tasks complete, `forge task list` shows a new `clean-code` task (e.g., `T-test-4.6` in breakdown mode, `T-quick-6` in quick mode) with status `pending`. The task lists its dependencies (the preceding verify-regression task).

2. **During execution**: When a subagent claims and runs the clean-code task, the developer sees standard task execution output — the LLM reads changed files, identifies cleanup candidates, and applies fixes. Each change is written to disk and logged in the task record. The prompt template includes an embedded quality gate step: after cleanup, the agent runs `just test` to confirm no regressions.

3. **Changes are auto-applied**: The clean-code task writes changes directly to the working tree (same as any forge task). There is no separate review/approval step — if the post-cleanup quality gate (`just test`) passes, the changes are accepted. If cleanup introduces a regression, the task fails with the test output, and the developer can inspect the diff via `git diff` and decide whether to revert or fix.

4. **Output on completion**: The task record contains a summary of what was cleaned (e.g., "removed 3 unused imports from handler.go, deduplicated validation logic in service.go"). If no changes were needed, the record notes "no cleanup opportunities found" and the task succeeds.

5. **Failure handling**: If the clean-code task fails (quality gate failure, LLM error, timeout), it is marked `failed` like any other task. The developer can re-run it via `forge task claim` + execution, or skip it and proceed manually.

### Innovation Highlights

- **Pipeline-native cleanup**: 将代码清理从手动/skill 触发提升为 pipeline 内置步骤，与 test pipeline 同级
- **Record-driven scope**: 通过读取已完成 task 的 records 自动确定清理范围，无需手动列出文件
- **Conservative by default**: 只清理明确的问题（dead imports, commented-out code blocks, duplicated logic, inconsistent naming），不进行架构重构

## Requirements Analysis

### Key Scenarios

1. **Happy path**: 所有 business tasks 和 test tasks 完成后，clean-code task 自动被 claim 并执行，清理所有变更文件
2. **Multi-profile project**: clean-code 是 shared task（不按 profile 拆分），一次清理覆盖全部变更
3. **`--no-test` flag**: 当使用 `--no-test` 时，clean-code task 也被跳过（与 test tasks 共命运）
4. **No changes needed**: clean-code 审查后无问题可清理，正常完成（无修改也算成功）
5. **Mid-execution crash**: clean-code task 在执行过程中被中断（subagent 超时或外部取消）。此时部分文件已被修改。行为：任务标记为 `failed`；开发者通过 `git diff` / `git checkout` 检查或还原变更；可重新 claim 执行
6. **Regression introduced**: clean-code 修改引入测试失败。行为：内置 quality gate（`just test`）捕获失败，任务标记为 `failed` 并附带测试输出；变更保留在 working tree 供开发者检查（design choice: fail-open-to-human，开发者通过 `git diff` 判断后自行 `git checkout` 还原或修复）
7. **Empty LLM output**: clean-code 的 LLM 返回空内容或无法解析的响应。行为：任务标记为 `failed`，record 记录 "LLM returned empty/invalid cleanup response"；开发者可重试

### Non-Functional Requirements

- **向后兼容**: 新增 task type 不影响现有 index.json（`omitempty` 已保证）。验证方式：对包含已有 tasks 的 index.json 执行 `forge task index`，确认不破坏现有 task 数据
- **幂等**: `forge task index` 重复运行不产生重复 clean-code task
- **Performance — execution time**: clean-code task 的 LLM 调用 + quality gate 总执行时间不超过 5 分钟（对于 50 个变更文件以内的 feature）。对于超大 diff（100+ 文件），允许分批处理但总时间不超过 15 分钟。超时则任务失败。**推导依据**：现有 `implementation` task 的典型执行时间为 2-4 分钟（含 quality gate），其中 LLM 推理约占 60%，`just test` 约占 30%。clean-code task 的 LLM 输入为多个文件的 diff（而非从零实现），上下文量与 implementation task 相当，因此 5 分钟上限基于现有 task 的 P95 执行时间 + 50% buffer。`just test` 的 P95 为 45 秒（forge-cli 项目实测），分批处理时每批约增加 30 秒，100 文件分 3 批合计约 12 分钟，取整为 15 分钟
- **Reliability — quality gate**: clean-code prompt 模板必须包含 post-cleanup quality gate 步骤（运行 `just test`），确保清理不引入 regression。quality gate 失败时任务标记为 `failed`，变更保留在 working tree 供开发者检查。这是一种 intentional design choice（fail-open-to-human）：quality gate 的职责是检测问题并停止 pipeline 自动推进，而非自动恢复——开发者通过 `git diff` 检查后可自行决定保留、修改或 `git checkout` 还原
- **Observability — cleanup record**: clean-code task 完成后，其 task record 必须包含清理摘要（哪些文件被修改、移除了哪些问题类型）。此摘要写入标准的 task record 文件，可供 `consolidate-specs` 和人工审查使用

### Constraints & Dependencies

- 新 type 需注册到 `TaskTypeRegistry`、`ValidTypes`、`typeToTemplate`、`InferType` 四处
- Prompt 模板需嵌入 `pkg/prompt/data/`（使用 `//go:embed`）

## Alternatives & Industry Benchmarking

### Industry Solutions

Automated code cleanup in CI/CD pipelines is a well-established practice. The following tools represent three distinct approaches used in production:

1. **SonarQube (v10.x, SonarSource)** — Static analysis platform with quality gates that run post-test in CI. SonarQube detects code smells, dead code, and duplication across languages. Its "Clean as You Code" methodology flags new issues on changed code only. Quality gates can block merges if cleanup thresholds are not met. However, SonarQube reports issues but does not auto-fix them — a developer must address each finding manually.

2. **ESLint `--fix` in CI (v9.x) / Ruff autofix (v0.8+)** — Linting tools with built-in auto-fix that run as CI steps after tests. ESLint `--fix` and Ruff `--fix` automatically correct formatting, unused imports, and simple rule violations. They are deterministic (no LLM involved), fast (milliseconds per file), and language-specific. Their limitation is scope: they handle syntactic/style issues but cannot detect semantic problems like duplicated business logic or inconsistent naming patterns that require understanding intent.

3. **Biome lint-and-fix (v2.x)** — All-in-one formatter+linter for JS/TS that combines Prettier-style formatting with ESLint-style linting in a single fast pass. Like ESLint `--fix`, it handles syntactic cleanup only and is limited to the JavaScript/TypeScript ecosystem.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero cost | No quality gate; specs extracted from dirty code | Rejected: quality gap |
| SonarQube quality gate | SonarSource v10.x | Industry-validated; multi-language; quality gate blocks merge | Report-only, no auto-fix; requires SonarQube server infrastructure; per-file scope, no cross-task view | Rejected: report-only does not close the loop |
| CI lint --fix step (ESLint/Ruff) | ESLint v9.x, Ruff v0.8+ | Deterministic; fast (<1s per file); CI-native | Syntactic-only; cannot detect semantic issues (duplication, dead logic, naming inconsistency); language-specific rules | Partial: good for formatting, insufficient for semantic cleanup |
| Pre-commit hook (husky + lint-staged) | npm | Runs before every commit | Staged-files-only scope; no cross-task awareness; adds Node.js dependency to Go project | Rejected: scope too narrow |
| **Pipeline clean-code task** | This proposal | Semantic cleanup (dead code, duplication, naming); record-driven full-feature scope; integrated between verify-regression and consolidate-specs | Adds one pipeline step; LLM-based (non-deterministic); requires quality gate to validate | **Selected** |

### Selection Rationale

CI lint --fix tools (ESLint, Ruff) are the closest industry analog, but they operate at the syntactic level and cannot detect the problems this proposal targets: duplicated business logic across tasks, dead imports left from iterative implementation, and inconsistent naming introduced by multiple parallel executors. These are semantic problems requiring code understanding, which is why an LLM-based cleanup is appropriate.

The pipeline placement (after verify-regression, before consolidate-specs) is the key differentiator: cleanup happens when all implementation is frozen but before specs are extracted, ensuring consolidate-specs operates on clean code. No CI lint tool provides this temporal guarantee within a task pipeline — they run at CI time, not at pipeline-orchestration time.

### Why Not Deterministic Semantic Analysis?

Tools like Semgrep (pattern-based semantic rules), tree-sitter (structural AST queries), and static dead-code analyzers can detect some of the targeted cleanup categories — Semgrep can find duplicated patterns, and tree-sitter queries can identify unused imports at the AST level. These were considered and rejected for three reasons:

1. **Scope fragmentation**: Each tool covers a subset of cleanup categories. Dead imports (tree-sitter), duplicated logic (Semgrep custom rules), commented-out code (trivial regex), and naming inconsistency (no deterministic tool handles this well) would require composing 3-4 tools with separate configurations, each with its own error model and integration surface. A single LLM prompt handles all four categories in one pass.

2. **Pipeline integration cost**: Forge's pipeline is LLM-driven end-to-end — task claiming, prompt synthesis, and execution all route through the existing LLM infrastructure. Adding a Semgrep/tree-sitter pass would introduce a non-LLM execution path with its own dependency management, error handling, and output parsing, for a net increase in system complexity.

3. **Non-determinism trade-off accepted**: The quality gate (`just test`) exists precisely because LLM output is non-deterministic. The cost of occasional false positives (the task fails and the developer inspects) is lower than the cost of maintaining a multi-tool deterministic pipeline for this use case. If the clean-code task's false-positive rate proves unacceptable in practice, a hybrid approach (deterministic tools for high-confidence categories + LLM for the rest) can be introduced as a follow-up.

## Feasibility Assessment

### Technical Feasibility

完全可行。已有成熟的 task type 扩展模式（11 个 type → 12 个），每个新 type 需修改 4 个 Go 文件 + 1 个 prompt 模板 + 文档。

### Resource & Timeline

1 个 task，预计 30min-1h 实现时间。改动量小且模式成熟。

### Dependency Readiness

无外部依赖。所有基础设施（task index、prompt synthesis、type inference）已就绪。

## Scope

### In Scope

- 新增 `clean-code` task type（TypeCleanCode constant）
- Breakdown 模式：T-test-4.6（verify-regression 之后，consolidate-specs 之前）
- Quick 模式：T-quick-6（verify-regression 之后）
- clean-code prompt 模板：3 步工作流（识别变更 → 清理 → quality gate）
- 依赖链：T-test-5 原本依赖 T-test-4.5，改为依赖 T-test-4.6
- 测试覆盖：更新 testgen_test、infer_test
- 文档同步：OVERVIEW.md、WORKFLOW.md、plugin docs

### Out of Scope

- Breakdown 模式的 eval-test-cases 前添加清理（过早，实现可能还在变化）
- 跨 feature 的全局清理（超出单 feature pipeline 范围）
- lint/format 自动修复（已由 quality gate 覆盖）
- `--no-test` 单独控制 clean-code（clean-code 与 test tasks 共命运）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| clean-code 过度修改引入 regression | M | M | Prompt 模板内嵌 quality gate（`just test`），quality gate 失败则任务标记 failed、变更保留供开发者检查（fail-open-to-human） |
| `noTest: true` 导致 quality gate 被跳过 | L | H | clean-code 模板内嵌独立的 quality gate 步骤，不依赖外部 test pipeline 的 gate |
| LLM 在大 diff (100+ 文件) 下超时 | M | L | 设置 15 分钟上限，超时任务失败；可重新 claim 执行 |
| clean-code 修改了不应修改的文件 | L | M | Prompt 模板限定只修改 record-driven scope 内的变更文件；所有修改可通过 `git diff` 审查 |

## Success Criteria

- [ ] `forge task index` 在 quick 和 breakdown 模式均生成 clean-code task
- [ ] clean-code task 的依赖链正确（breakdown: T-test-4.5→T-test-4.6→T-test-5；quick: T-quick-5→T-quick-6）
- [ ] `forge prompt get-by-task-id T-test-4.6` 和 `T-quick-6` 返回有效 prompt
- [ ] Prompt 模板包含内嵌 quality gate 步骤（`just test`），验证方式：检查模板内容包含该步骤
- [ ] clean-code task 完成后 record 包含清理摘要（修改了哪些文件、移除了哪些问题类型）
- [ ] 所有现有测试通过（`go test -race -cover ./...`），覆盖率不低于改动前基线
- [ ] `make check-docs` 通过
- [ ] 对已有 index.json 执行 `forge task index` 不破坏现有 task 数据（向后兼容验证）
- [ ] **Effectiveness**: Given a test fixture feature branch with planted cleanup candidates (at least 3 unused imports, 1 commented-out code block, 1 duplicated validation function across 2 files, and 1 naming inconsistency where a helper is named `getData` in one file and `fetchData` in another for the same operation), clean-code task produces a diff that removes at least 5 of these 7 planted issues. Validated by running the task against the fixture and checking `git diff --stat` for file modifications matching the cleanup targets
