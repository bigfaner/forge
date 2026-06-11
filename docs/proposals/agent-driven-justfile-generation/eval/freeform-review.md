# Freeform Expert Review

**Reviewer**: Surface-Aware Dispatcher & Test Orchestration Architect
**Document**: `docs/proposals/agent-driven-justfile-generation/proposal.md`
**Date**: 2026-06-08

---

## Background Assessment

This proposal targets the `init-justfile` skill within the Forge plugin ecosystem. The skill currently generates a project's justfile through a three-layer pipeline: load a language-specific template (`templates/<lang>.just`), apply Convention overrides, then let the LLM customize remaining details. The core complaint is that the first and third layers frequently conflict -- templates provide defaults that mismatch the actual project, forcing the LLM to replace most template content anyway.

The proposed solution removes all six language templates plus the `rules/project-detection.md` classification rule. Instead, the agent would read surface rule files (`rules/surfaces/<type>.md`) to enumerate which recipes must exist (the "contract"), then generate each recipe's concrete commands from its own knowledge of the detected language/framework plus Convention files. A new `rules/server-lifecycle.md` would extract the PID-tracking / idempotent-startup / health-check bash pattern that is currently duplicated across all templates into a single reusable reference.

The fundamental architectural shift is from **template-driven generation** (static code templates with placeholder substitution) to **contract-driven LLM generation** (structural rules define what must exist; the LLM decides how to fill it). The proposal argues this aligns with Forge's surface-first philosophy: surface type determines the orchestration sequence and recipe structure, while language/framework only influences the specific shell commands within each recipe.

A critical architectural observation that the proposal does not address explicitly: the surface rule files (`rules/surfaces/<type>.md`) are **shared dependencies** consumed by two distinct skills. `run-tests` reads the "Orchestration Sequence" table as a **dispatch contract** (which steps to execute, what exit codes mean), while `init-justfile` would now read the same files as **generation contracts** (which recipes to generate, what their signatures and exit codes must be). The proposal's "Recipe Generation Requirements" section would be appended to files that already serve a different consumer. This dual-consumer constraint shapes many of the risks identified below.

---

## Key Risks

The proposal presents a clean conceptual argument but leaves several critical failure modes unaddressed. The risks are organized from highest to lowest systemic impact.

**风险：Server lifecycle bash 提取为 rule 文件后，agent 在没有可执行参考实现的情况下生成的 PID 追踪代码可能丢失边界条件处理。**

> 提案声明："server lifecycle bash 复杂度高，agent 可能遗漏边界情况" -- 通过 "提取为独立 rule 文件，提供完整参考模式" 缓解。

当前模板中的 server lifecycle 代码（见 `templates/go.just` 第 27-65 行）是一个经过实战验证的三层幂等模式：Layer 1 检查 PID 文件中进程是否存活，Layer 2 检查是否已有服务响应（手动启动的场景），Layer 3 才真正启动。这涉及 `kill -0` 信号检测、`tr -d '\r'` 跨平台 PID 文件清洗、`trap _cleanup EXIT INT TERM` 信号处理器注册、早期崩溃检测（启动后 1 秒检查进程是否仍在）、10 次重试的 30 秒健康检查窗口。

提案将此提取为 `rules/server-lifecycle.md` 作为"参考模式"——但 rule 文件是**描述性文本**，不是可执行代码。LLM 从文本描述生成 50+ 行 bash 时，丢失 `tr -d '\r'` 这类 Windows 边界处理（PID 文件在 Windows 上可能携带 `\r`）的概率很高。当前的模板是**已验证的可执行代码**，LLM 只需要替换启动命令本身；新方案要求 LLM 从描述重新生成整个 bash 块。这不是"轻微不确定性"——这是将确定性代码替换为概率性生成的本质降级。

**问题：Surface rule 文件同时服务 `init-justfile`（生成方）和 `run-tests`（调度方）两个消费者，提案对新增的 "Recipe Generation Requirements" section 如何与现有 Orchestration Sequence / Recipe Invocation Contract 保持一致没有定义验证机制。**

> 提案声明："简化 5 个 surface rule 文件：保留编排序列/recipe 契约/journey 策略，替换 TODO stub 模板为 'Recipe Generation Requirements' section"。

查看现有 surface rule 文件（如 `rules/surfaces/api.md`），"Recipe Invocation Contract" 表格定义了每个 recipe 的 `just` 签名和退出码语义。`run-tests` 依赖这些退出码语义来决定下一步行为（例如 dev 失败时直接跳到 teardown，exit code 2 表示可重试）。如果新增的 "Recipe Generation Requirements" section 与 Recipe Invocation Contract 之间存在不一致——例如 Generation Requirements 允许省略某个 recipe，但 Invocation Contract 声明它是必需的——`run-tests` 会在运行时遇到缺失 recipe 的失败，而不是在生成时被捕获。

当前的 TODO stub 模板（surface rule 文件末尾的 `just` 代码块）实际上充当了两个消费者之间的**结构一致性锚点**：它同时展示了 recipe 结构（给 init-justfile）和预期行为（给 run-tests 参考退出码语义）。替换为纯描述性的 "Recipe Generation Requirements" 打破了这个隐式一致性保证。

**风险：移除 `rules/project-detection.md` 后，混合语言多 surface 项目的语言检测信号源消失，agent 需要自行推断每个 surface 的语言/框架，但没有定义检测策略。**

> 提案声明删除 `rules/project-detection.md`，并且 "agent 根据 surface 配置独立处理每个 surface，无需知道整体'类型'"。

当前的 `project-detection.md` 定义了清晰的检测信号映射：`package.json` -> frontend、`go.mod` -> backend、以及混合项目的子目录发现（`find . -name package.json -maxdepth 3`）。更重要的是，mixed.just 模板中的 `frontend_dir` / `backend_dir` 变量和 per-service PID 追踪（`.pid-frontend`、`.pid-backend`）处理了**多服务并行启动**的场景。

提案移除了项目类型分类和检测算法，但没有说明 agent 在多 surface 项目中如何确定：每个 surface 使用什么语言？backend surface 的入口点在哪里？frontend surface 的运行脚本是 `start`、`preview` 还是 `dev`？这些信息在当前方案中来自 `project-detection.md` 的确定性行为，新方案中它们完全依赖 LLM 的推断。对于一个 `backend=api + frontend=web` 的 Go+Node 混合项目，agent 需要同时知道 Go 的入口点检测（`ls cmd/*/main.go`）和 Node 的运行脚本检测（解析 `package.json` 的 `scripts` 字段）——这些都是**确定性的文件系统检测**，不是需要 LLM 推断的模糊问题。

**风险：提案对 "一致性" 的定义——"recipe 名称、分组、边界标记不变；具体命令可能因 LLM 变化而有细微差异"——低估了命令差异对下游消费者的连锁影响。**

> 提案声明："一致性：相同项目多次运行生成的 justfile 结构一致（recipe 名称、分组、边界标记不变；具体命令可能因 LLM 变化而有细微差异）"。

`run-tests` 执行 `just <recipe-prefix>dev` 然后依赖 recipe 的**退出码语义**来驱动状态机。如果 LLM 在第一次运行时生成的 dev recipe 使用 `go run cmd/server/main.go & echo $! > .pid-file`，而第二次运行时生成 `go run ./... 2>&1 &`，这两种写法的后台进程行为、错误传播、PID 文件管理完全不同。`run-tests` 的 HARD-GATE 状态机（dev exit 0 -> probe -> test -> teardown）依赖的是 recipe 的**行为契约**，不仅仅是**名称契约**。命令的"细微差异"可能导致：进程未正确后台化、PID 文件未写入、`trap` 处理器未注册、teardown 时无法找到进程。这不是理论风险——这是从确定性模板切换到概率性生成后的必然结果。

**问题：混合项目的端口冲突处理在提案中完全缺失。**

> 当前 `mixed.just` 模板使用 `probe` recipe 检查 `tests/config.yaml` 中的 `baseUrl` 和 `apiBaseUrl` 来判断哪些服务需要启动，并且使用 per-service PID 文件（`.pid-frontend`、`.pid-backend`）避免重复启动已运行的服务。

提案没有说明 `server-lifecycle.md` 是否会包含这种 per-service 的端口感知启动模式，还是仅提供单一服务的 PID 追踪参考。对于一个 `backend=api`（Go，端口 8080）+ `frontend=web`（Node，端口 3000）的混合项目，两个 surface 的 dev recipe 各自启动一个服务。如果 `server-lifecycle.md` 只描述单服务模式，agent 生成的两个 recipe 可能会各自独立管理 PID 文件但缺乏协调——probe recipe 需要同时检查两个端口健康。

**风险：Verification step（dry-run + actual execution）的覆盖范围不足以捕获 server lifecycle 的错误。**

> 提案在 Key Risks 中声明 "verification step (dry-run + actual) 捕获错误"。

SKILL.md 的 Step 4b 定义了 verification 策略：long-running recipes（`<prefix>dev`）执行 `timeout 10 just <prefix>dev 2>&1 || true`。这个 10 秒 timeout 验证的是"进程没有在 10 秒内崩溃"，但它**无法验证**：PID 文件是否正确写入（因为 timeout 会 kill 进程，PID 文件可能被 cleanup handler 删除）、健康检查重试逻辑是否正确（需要服务真正 ready 才有意义）、teardown 是否能正确清理（因为 dev 被 timeout kill 了，不是正常 teardown）。

实际上当前模板中最复杂的部分——server lifecycle 的三层幂等逻辑——在 verification 中几乎不可测试。dry-run 只检查语法；actual execution 的 10 秒 timeout 只验证了启动不崩溃。真正的边界情况（PID 文件残留、进程被外部 kill 后 PID 被回收、Windows 上的 `\r` 污染）在 verification 中完全无法覆盖。

**风险：提案声称 "可在一次 session 内完成"，但实际涉及 5 个 surface rule 文件的结构性重写 + 1 个新的 server-lifecycle rule + SKILL.md 流程重写，scope 被低估。**

> 提案声明："改动集中在 `init-justfile` skill 目录内...预计单个 skill 改动，可在一次 session 内完成"。

表面上看是删除模板和重写 SKILL.md，但实际影响链是：删除模板 -> SKILL.md 的 HARD-RULE 中引用模板的三层生成流程需要重写 -> 每个 surface rule 文件的 TODO stub 模板需要替换为 "Recipe Generation Requirements" -> 新的 server-lifecycle.md 需要从零编写并覆盖所有现有模板中的 lifecycle 模式（单服务 + 多服务 + Windows 边界处理）-> verification step 可能需要增强以覆盖 LLM 生成代码的不确定性。

更关键的是，这是一次**不可逆的架构迁移**：一旦模板被删除，如果 agent-driven 方案在实际项目中表现不佳（比如冷启动项目的生成质量显著低于模板），回退路径是重新编写模板——而不是 git revert，因为 SKILL.md 的新流程与模板驱动流程不兼容。

---

## Improvement Suggestions

**建议：将 server-lifecycle.md 设计为可执行代码模板（带插槽），而不是描述性文档。**

Addresses: Server lifecycle bash 可靠性风险

提案应将 `server-lifecycle.md` 定位为**代码模板**而非参考文档。具体做法：在 rule 文件中保留当前 go.just / mixed.just 中经过验证的完整 bash 代码块，用 `{{START_COMMAND}}`、`{{HEALTH_CHECK_URL}}`、`{{PID_FILE_NAME}}` 等占位符标记需要替换的位置。Agent 的职责从"从描述生成 50 行 bash"降级为"识别正确的占位符值并替换"——这与当前模板的使用模式同构，但不再需要维护 6 份几乎相同的副本。

这样 server-lifecycle.md 同时服务两种消费者：init-justfile 读取它作为代码模板进行占位符替换；run-tests 不直接使用它（仍然依赖 Recipe Invocation Contract 的退出码语义）。

**建议：保留 `project-detection.md` 的检测信号映射，重构为每个 surface 的语言检测规则，而非整体项目类型分类。**

Addresses: 混合语言检测信号源消失

提案正确地指出项目类型分类（frontend/backend/mixed）可以移除——因为 surface 配置已经提供了这个信息。但检测信号映射（哪个 marker file 对应什么语言/工具链）和入口点检测（Go 的 `cmd/*/main.go`、Node 的 `scripts.start`）是**确定性的机械操作**，不应交给 LLM 推断。

建议将 `project-detection.md` 重构为 `rules/language-detection.md`，保留 marker file -> 语言的映射和 per-language 的入口点检测命令，但移除 frontend/backend/mixed 分类算法。Agent 读取 surface 配置后，对每个 surface 所在目录执行对应语言的确定性检测，然后将结果注入 recipe 生成。

**建议：为 surface rule 文件引入结构化验证约束，确保 "Recipe Generation Requirements" section 与 Recipe Invocation Contract 的一致性。**

Addresses: 双消费者一致性风险

在每个 surface rule 文件的 "Recipe Generation Requirements" section 中，增加一个明确的 recipe 清单矩阵：recipe name | required/optional | expected exit codes。init-justfile 生成后，执行一个简单的 post-generation check：对 Recipe Invocation Contract 中声明的每个 required recipe，验证生成的 justfile 中是否存在同名 recipe 且具有正确的参数签名。这个 check 可以通过 `just --list` 的输出解析完成，成本极低但能捕获结构不一致。

**建议：在混合项目的 surface rule 文件中增加 multi-service lifecycle 的专门指导。**

Addresses: 混合项目端口冲突

当前的 api.md 和 web.md 各自只描述单 surface 的 lifecycle。对于多 surface 项目（例如 `backend=api + frontend=web`），需要在至少一个 surface rule 文件中（或 server-lifecycle.md 中）明确说明：当多个 surface 各自需要启动服务时，PID 文件必须 per-service 命名（`.pid-backend`、`.pid-frontend`），probe recipe 必须检查所有相关服务的健康状态，teardown 必须清理所有已启动的服务。mixed.just 模板中的 multi-service 模式（第 42-88 行）已经解决了这个问题，新方案不应丢失这个能力。

**建议：增加"黄金回归测试"作为 success criteria，验证 agent 生成的 justfile 在结构关键点上与当前模板输出匹配。**

Addresses: 一致性定义不足

当前提案的 success criteria 包括"对 Go/Node/Python/Rust 项目，agent 生成的 recipe 通过 verification step"，但 verification step 的覆盖范围有限。建议增加一个明确的回归标准：选取 3-5 个代表性项目配置（单 surface Go、单 surface Node、多 surface mixed），对每个配置执行 agent-driven 生成，然后对比生成的 justfile 与当前模板输出，验证以下关键点完全一致：boundary markers 位置、recipe 分组结构、`[linux]/[windows]` 双平台属性、`# user-customized` 标记位置。具体命令内容允许差异，但结构骨架必须匹配。这为不可逆迁移提供了安全网。

**建议：增加回退策略的显式文档，定义何时认为 agent-driven 方案失败以及如何恢复。**

Addresses: Scope 被低估、不可逆迁移风险

提案应该在 Key Risks 表中增加一行：如果 agent 在 N 个连续项目上的生成质量（以 verification pass rate 衡量）低于某个阈值（例如 80% 的 recipe 通过 actual execution），则触发回退评估。回退方案可以是：保留 server-lifecycle.md 作为代码模板（见第一个建议），但同时保留一个精简版的 `templates/fallback.just`，在 agent generation 失败时作为 safety net。这不是回归到当前方案，而是在新方案中内置降级路径。
