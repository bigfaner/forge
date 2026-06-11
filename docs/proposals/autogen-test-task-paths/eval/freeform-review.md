# Freeform Expert Review

**Expert**: Multi-Layer Context Propagation & Rendering Completeness Auditor
**Document**: `docs/proposals/autogen-test-task-paths/proposal.md`
**Date**: 2026-05-28

## Background Assessment

This proposal addresses a genuine "declared but not rendered" defect in the Forge task pipeline system. The core observation is sound: `FeatureSlug` is available as a field on both `promptTemplateData` and `autogenTemplateData`, and it flows into the rendering engine, but the six test-pipeline prompt templates (`test-gen-journeys`, `eval-journey`, `test-gen-contracts`, `eval-contract`, `test-gen-scripts`, `test-run`) never emit the slug into agent-visible output. The agent is forced to reverse-engineer the slug from `TASK_FILE` path components — a fragile, multi-step process that the proposal correctly identifies as wasteful.

The proposed fix is a two-pronged template change: add `FEATURE_SLUG: {{.FeatureSlug}}` to all six prompt templates, and add a `## Feature Paths` section with `ls` discovery commands to all six embed (task) templates. The skill layer is explicitly left unchanged as a fallback for standalone invocation outside the pipeline.

The proposal rests on several assumptions: that `FeatureSlug` is already populated in both data structs (confirmed by reading `promptTemplateData.FeatureSlug` and `autogenTemplateData.FeatureSlug`), that the embed template's `{{.FeatureSlug}}` is filled by the CLI at index time from the directory path, and that the prompt template's `{{.FeatureSlug}}` is filled by the dispatcher from `index.json`'s `feature` field. It further assumes that the six listed templates constitute the complete set of test-pipeline templates, and that no Go code changes are needed.

My reading aligns with the proposal's intent, though I note a structural asymmetry in the template fleet that the proposal does not acknowledge: the embed templates (task .md files) are architecturally distinct from the prompt templates in both their data model and rendering engine. The embed templates use `autogenTemplateData` rendered during `forge task index`, while the prompt templates use `promptTemplateData` rendered at dispatch time. This distinction matters more than the proposal suggests.

## Key Risks

The proposal's risk assessment is overly optimistic and misses several structural concerns. I will address the most material ones.

风险：Embed 模板和 Prompt 模板的 `FeatureSlug` 来源不同，但提案未分析两者不一致的场景。

> "CLI 从 `docs/features/<slug>/` 目录路径填充" vs. "Dispatcher 从 `index.json` 的 `feature` 字段传入"

The embed template's `{{.FeatureSlug}}` is resolved at `forge task index` time by reading the directory path. The prompt template's `{{.FeatureSlug}}` is resolved at dispatch time from `index.json`. Both originate from the same CLI command, but they are read from different stores at different moments in time. If a user runs `forge task index`, then renames the feature directory (or moves the feature), the embed templates will contain the old slug baked into the generated .md files, while the prompt template will use whatever is in `index.json` at dispatch time. The proposal claims "确定性强" (high certainty) for both paths but does not acknowledge this temporal coupling gap. In practice, directory renames during active development are a real scenario — the slug becomes stale in embed templates until the next `forge task index` re-run. If the agent sees `FEATURE_SLUG: my-feature` in the prompt but the embed template says `ls docs/features/my-renamed-feature/testing/`, the mismatch creates confusion rather than clarity.

风险：Proposal 将 "FeatureSlug 渲染为空" 的 Impact 评为 Low，但未分析空值在 embed 模板中产生的具体后果。

> "FeatureSlug 渲染为空 | L | L | prompt 模板中 FeatureSlug 已在 context 声明，dispatcher 始终传入"

If `FeatureSlug` renders empty in the prompt template, the proposal correctly notes the agent can fall back to parsing `TASK_FILE`. But the embed template scenario is worse: `ls docs/features//testing/` is an invalid command that either fails silently or lists the wrong directory. The proposal does not distinguish between the two failure modes — prompt-empty (recoverable via path parsing) versus embed-empty (produces broken discovery commands). Treating them as a single row in the risk table conflates two different severity levels. The embed template empty-slug scenario should be rated Medium impact at minimum, since the proposed `## Feature Paths` section becomes actively misleading rather than merely absent.

问题：Proposal 声称改动仅涉及 6 个 embed 模板和 6 个 prompt 模板，但实际 embed 模板舰队中存在结构性差异被忽略。

> "6 个测试流水线 embed 模板统一添加两类发现命令"

Examining the actual embed templates reveals a critical asymmetry. The templates `test-gen-journeys.md` and `test-gen-contracts.md` in `forge-cli/pkg/task/templates/` are rich documents — they contain multi-step process flows, discovery strategies, input source specifications, acceptance criteria, and already reference `{{.FeatureSlug}}` extensively in path patterns like `docs/features/{{.FeatureSlug}}/testing/`. In contrast, `test-run.md` and `test-gen-scripts.md` are thin wrappers that delegate entirely to skills — they contain only a two-step workflow ("Read Task Definition" then "Invoke Skill"). Adding `## Feature Paths` with `ls` commands to the thin-wrapper templates is largely redundant because the agent immediately invokes the skill, which has its own discovery logic. The proposal applies a uniform change to a non-uniform template fleet without acknowledging that the value of the change varies dramatically by template type.

问题：Proposal 未验证 prompt 模板中 `FEATURE_SLUG` 行的放置位置是否与现有模板结构兼容。

> "6 个测试流水线 prompt 模板在 `TASK_FILE` 行之后添加：`FEATURE_SLUG: {{.FeatureSlug}}`"

Looking at the existing prompt template `doc-summary.md`, which already renders `FEATURE_SLUG: {{.FeatureSlug}}`, the placement is between `TASK_FILE` and `SURFACE_KEY` lines. However, the six test-pipeline prompt templates have a different structure: `SURFACE_KEY` is conditionally rendered via `{{if .SurfaceKey}}SURFACE_KEY: {{.SurfaceKey}}{{end}}`, and is followed by an optional `## PhaseSummary` block. The proposal does not specify whether `FEATURE_SLUG` should be conditional (only rendered when non-empty) or always rendered. If `FeatureSlug` is always populated in practice, unconditional rendering is fine — but the proposal should state this explicitly rather than leaving it as an implicit assumption. Furthermore, the prompt templates do not currently use `{{.FeatureSlug}}` in their body at all; only the embed templates do. Adding `FEATURE_SLUG` to the header block of prompt templates changes the information density of the agent's initial context window without demonstrating that the agent actually needs the slug before it reads the task file in Step 1.

风险：Embed 模板的 discovery 命令路径约定 (`docs/features/<slug>/testing/`) 与 skill 中的路径可能不同步。

> "发现 journeys（与 `run-tests` skill Step 1.5 一致）"

The proposal claims the discovery commands are "consistent with" the skill's Step 1.5, but this is a snapshot assertion. The skill's discovery logic lives in plugin YAML that is independently versioned from the CLI's embed templates. If the skill's path convention changes (e.g., from `docs/features/<slug>/testing/` to `tests/<slug>/` or any other convention shift), the embed templates become out of date with no validation mechanism to detect the drift. The proposal creates a new coupling point between two independently evolving artifacts without establishing a synchronization contract or validation check. This is the "declared but stale" variant of the "declared but not rendered" defect — the variable reaches the agent, but the value is wrong.

问题：Proposal 的 Scope 章节列出了 6 个 embed 模板和 6 个 prompt 模板，但实际文件系统中的 embed 模板集合因 worktree 而异，且主分支缺少部分模板。

> "修改 6 个测试流水线 embed 模板（`forge-cli/pkg/task/templates/`）"

On the current `v3.0.0` branch, the `forge-cli/pkg/task/templates/` directory contains exactly 6 test-pipeline embed templates: `test-gen-journeys.md`, `eval-journey.md`, `test-gen-contracts.md`, `eval-contract.md`, `test-gen-scripts.md`, and `test-run.md`. This matches the proposal's list. However, the worktree branches show additional templates that may be merged before this proposal is implemented. The proposal does not address whether the template list should be validated against `task.ValidTypes` at build time, or whether the count of 6 is an unstable number that could grow. If new test-pipeline task types are added, the proposal's "6 templates" scope becomes stale — and there is no mechanism in the proposal to detect this.

风险：Embed 模板添加 `## Feature Paths` 后，agent 可能在 Step 1（读 task file）时执行 discovery 命令，但此时 skill 尚未被调用，agent 缺少执行上下文。

> "Subagent 读 task file，看到 discovery 命令，可预先定位 journeys 和 contracts"

The proposal frames the discovery commands in the embed template as a "pre-exploration" capability. But the thin-wrapper templates (`test-run.md`, `test-gen-scripts.md`) instruct the agent to read the task file in Step 1 and immediately invoke the skill in Step 2. If the agent starts running `ls` commands during Step 1 based on the new `## Feature Paths` section, it is performing work that the skill will redo. For the rich templates (`test-gen-journeys.md`, `test-gen-contracts.md`), the discovery commands duplicate information already present in the template body — the `gen-journeys` template already contains "Scan `docs/features/{{.FeatureSlug}}/testing/` for all Journey files" and similar path references. The proposal claims "不重复但互相补充" (no duplication, mutually complementary) but does not audit the actual template content to verify this claim.

## Improvement Suggestions

建议：区分 embed 模板的改动策略 —— 仅对内容稀薄的模板添加 `## Feature Paths`，对已有丰富路径引用的模板跳过。

Addresses: "Embed 模板改动对非统一模板舰队价值不一" 问题

> What changes: 将 6 个 embed 模板分为两类处理。对 `test-run.md` 和 `test-gen-scripts.md` 等薄包装模板，添加 `## Feature Paths` discovery 命令，因为这些模板当前完全不含路径信息。对 `test-gen-journeys.md` 和 `test-gen-contracts.md` 等已有丰富路径引用的模板，验证现有 `{{.FeatureSlug}}` 引用是否已覆盖所有必要路径，仅补充缺失部分而非机械地添加统一的 `## Feature Paths` 区域。对 `eval-journey.md` 和 `eval-contract.md`，评估其 discovery strategy 是否已充分。这样做避免了在已有路径信息的模板中制造冗余，同时确保薄模板获得必要的上下文。

建议：为 embed 模板中的 FeatureSlug 空值添加防御性渲染逻辑。

Addresses: "FeatureSlug 渲染为空在 embed 模板中产生无效命令" 风险

> What changes: 在 embed 模板的 `## Feature Paths` 区域添加条件渲染：仅当 `FeatureSlug` 非空时才输出 discovery 命令。例如，在 embed 模板中使用 `{{if .FeatureSlug}}## Feature Paths\n\n...{{end}}` 模式，或在 CLI 填充阶段添加断言：如果 `FeatureSlug` 为空则终止索引生成并报错。后一种方案更强——在编译期而非运行期捕获缺陷，符合 "快速失败" 原则。具体做法是在 `renderBody()` 调用前校验 `autogenTemplateData.FeatureSlug` 非空，为空时返回明确的错误信息而非生成含无效路径的 .md 文件。这需要少量 Go 代码改动，但考虑到风险场景的严重性（生成损坏的 task 文件并推入 agent pipeline），这是合理的代价。

建议：在 Prompt 模板中将 `FEATURE_SLUG` 设为条件渲染或明确标注其始终非空的保证。

Addresses: "Prompt 模板中 FEATURE_SLUG 放置位置和条件性不明确" 问题

> What changes: 要么在提案中明确声明 `FeatureSlug` 在所有测试流水线调度场景下均为非空值（并验证此断言——检查 dispatcher 是否存在 FeatureSlug 为空的合法路径），要么在模板中使用 `{{if .FeatureSlug}}FEATURE_SLUG: {{.FeatureSlug}}{{end}}` 条件渲染。同时，参考 `doc-summary.md` 的现有实现，将 `FEATURE_SLUG` 放在 `TASK_FILE` 行之后、`SURFACE_KEY` 行之前，保持模板间的结构一致性。

建议：建立 embed 模板路径与 skill discovery 路径的同步验证机制。

Addresses: "Embed 模板 discovery 命令与 skill 路径不同步" 风险

> What changes: 在 `ValidatePromptTemplates()` 或类似的启动时校验中，添加对 embed 模板 `## Feature Paths` 区域路径格式的约定检查——验证路径符合 `docs/features/<placeholder>/testing/` 模式。更进一步，可以在 CI 中添加一个轻量级测试：解析 embed 模板中的路径字符串，与 skill 定义中的路径模式做文本比对，检测不一致时报警。不需要复杂的契约系统，一个简单的字符串常量提取 + 比对测试即可大幅降低漂移风险。Proposal 应在 Scope 的 Out of Scope 章节中明确提及此同步风险，并建议后续迭代建立验证。

建议：在提案中明确 agent 使用 embed 模板 discovery 命令的时机——仅限 Step 1 之前或 Step 1 之中，且不应替代 skill 的 discovery 逻辑。

Addresses: "Agent 可能在不当时机执行 discovery 命令" 风险

> What changes: 在 embed 模板的 `## Feature Paths` 区域添加使用说明，明确这些命令是供 agent 在 Step 1（Read Task Definition）阶段理解任务范围时使用的参考路径，不应作为独立执行命令。或者，更直接地，将 discovery 命令改为路径列表而非可执行命令格式——用 `Journeys directory: docs/features/{{.FeatureSlug}}/testing/` 替代 `ls docs/features/{{.FeatureSlug}}/testing/`。这样做既提供了路径信息（agent 可在需要时自行决定是否探索），又避免了诱导 agent 在错误时机执行 shell 命令。
