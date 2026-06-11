# Freeform Expert Review

## Background Assessment

This proposal has evolved through four evaluation iterations (baseline 617 -> 666 -> 782 -> 766 -> now with substantive additions). The current revision adds four major features that were previously identified as blind spots: (1) re-run behavior when `config.yaml` already contains surfaces, (2) user override semantics (source annotation is TUI-only, not persisted), (3) a performance budget of <50ms, and (4) a `forge surfaces detect` subcommand with `--dry-run` and non-interactive support. These additions directly address the most persistent critiques from the rubric evaluations — the re-run behavior was flagged in iteration-0 through iteration-3 without resolution, and the user override flow was an undefined gap in every prior evaluation.

The proposal's core architecture remains a two-layer detection pipeline: dependency signals first, structural inference as fallback, manual entry as last resort. The new `Sources map[string]string` field on `DetectResult` provides per-path source annotation (e.g., `{"forge-cli/cli": "inference:cmd-dir"}`) — a design that correctly addresses the prior freeform review's concern about mixed sources in multi-surface projects. The `forge surfaces detect` command adds an explicit entry point for re-detection outside the init flow, which is architecturally clean: it reuses the detection + TUI infrastructure without coupling to the init pipeline.

Having re-examined the codebase, I note that `DetectResult` currently has no `Sources` field — only `Surfaces`, `Conflicts`, and `IsScalar`. The proposal's `Sources map[string]string` addition is backward-compatible because Go's zero-value for maps is `nil`, meaning existing construction sites do not need updates. The `runSurfaceConfig` function in `init.go` currently builds the summary detail as either `surfaces["."]` (scalar) or `"%d mappings"` (map form). The proposal correctly targets the map-form branch for display improvement, and the scalar form already shows the type name.

The new re-run prompt design (`"Surfaces already configured: cli. Re-detect?"` with Confirm/Re-detect options) replaces the previous unconstrained "preserves existing config" constraint with a concrete interaction. The user override design (source is TUI-only, discarded after confirmation) is a pragmatic choice that avoids config schema proliferation.

## Key Risks

The additions resolve several long-standing gaps but introduce new risks that the proposal does not fully acknowledge.

风险：`forge surfaces detect` 子命令与 `forge init` 中的 surface 配置流程存在功能重叠和语义分歧。

> "`forge surfaces detect` command: new subcommand that runs detection + inference and displays results, independent of `forge init`. Gives users an explicit entry point for re-detection without re-running the full init flow."

当前 `forge surfaces` 命令已经存在（`surfaces.go`），功能是查询已配置的 surfaces（读取 config，支持 list/query/types 三种模式）。新增的 `forge surfaces detect` 作为子命令挂载在 `surfaces` 下，但它的行为本质上是"重写 surfaces 配置"，这与现有 `surfaces` 命令的只读语义矛盾。`forge surfaces list` 是只读的，`forge surfaces query` 是只读的，而 `forge surfaces detect` 是写操作。这种混合语义在 CLI 设计中是反模式 — 用户预期 `get` 类命令不会产生副作用。

更关键的是，提案描述 `forge surfaces detect` 的交互行为为"runs detection + inference, shows results with source annotations, asks TUI confirmation (same flow as init), writes to config on confirm"。这意味着该命令在交互式终端中会触发 TUI 并写入配置。但提案同时声明非交互式终端"prints results to stdout and exit (no TUI, no config write)"。这就产生了一个语义不对称：同一个命令在不同终端模式下有完全不同的副作用行为。在 CI 脚本中运行 `forge surfaces detect` 会是只读的，但在开发者终端中运行会修改配置文件。这种隐式的行为切换比显式的 `--write` flag 更危险，因为它违反了最小惊讶原则。

建议：将 `forge surfaces detect` 的默认行为改为只读（显示结果，不写配置），添加 `--write` 或 `--apply` flag 来启用配置写入。`--dry-run` 则变为冗余（因为默认就是 dry-run），可以简化为 `--write` 和非交互模式两种控制维度。

问题：Re-run 行为设计中的 "Confirm" 选项会导致用户无法通过 `forge init` 修正之前错误的 surface 配置。

> "Re-run behavior: `config.yaml` with `surfaces: cli` present -> TUI shows `"Surfaces already configured: cli. Re-detect?"` with Confirm / Re-detect; Confirm returns `SKIPPED surfaces (already configured)`, Re-detect runs full detection + inference pipeline"

这个设计假设用户总是知道当前配置是否正确。但考虑一个常见场景：用户首次运行 `forge init` 时，依赖信号检测到了错误的 type（比如 `api` 实际应该是 `web`），用户当时没有注意到就 Confirm 了。之后用户想修正为 `web`。运行 `forge init` 时看到 "Surfaces already configured: api. Re-detect?" — 如果选择 Confirm，配置保持 `api`；如果选择 Re-detect，检测管道可能再次返回 `api`（因为依赖信号没变）。用户没有直接的 "Edit" 选项来修正已有配置。

现有代码中 `runConfigInitIfNeeded` 已经有一个类似的模式：当 config 存在时询问 "Config already exists. Reconfigure?" 并提供 Yes/No。但那个流程会重新走完整的 config 设置（auto-behavior, worktree）。而 surface 的 re-run 设计只有 Confirm/Re-detect 两个选项，缺少 Edit（手动修正）路径。这与 `askMapConfirmation` 的四选项设计（Confirm/Edit/Add/Delete）形成不一致。

风险：`Sources map[string]string` 的 value 格式 `"inference:cmd-dir"` 和 `"dependency:cobra"` 使用了冒号分隔的命名空间格式，但没有定义格式规范。

> "`Sources` map correctly populated: `{"forge-cli/cli": "inference:cmd-dir"}` for inferred paths, `{"forge-cli/cli": "dependency:cobra"}` for detected paths"

提案定义了两种 source 值格式：`"inference:<rule-id>"` 和 `"dependency:<signal-name>"`。但这个格式只在 success criteria 中通过示例出现，没有在 Schema 或 API 契约中正式定义。这会导致以下问题：

1. 下游消费者（TUI formatter、init summary builder）需要解析这个字符串来提取 `rule-id` 和 `signal-name`，但没有文档化的解析规则。如果未来添加新的 source 类型（如 `"hybrid:..."` 或 `"user:..."`），解析逻辑会变得脆弱。

2. success criteria 中的 key 格式不一致：有的 criterion 用 `"forge-cli/cli": "inference:cmd-dir"`，有的用 `"cli": "dependency:cobra"`。前者 key 是相对路径，后者 key 是 surface type。这暗示 `Sources` map 的 key 语义在不同场景下不同（per-path vs per-type），但提案没有明确 key 的含义是 path 还是 type。

查看现有代码，`DetectResult.Surfaces` 的 key 是路径（`.` 表示根，`forge-cli/cli` 表示子目录）。如果 `Sources` 的 key 应该与 `Surfaces` 的 key 一致（即路径），那么 success criteria 中 `"cli": "dependency:cobra"` 的 key 应该是 `"."` 而非 `"cli"`。这个不一致需要修正。

问题：性能预算 <50ms 缺乏测量基准和分摊方案。

> "Inference performance budget: all inference functions combined must complete in <50ms (filesystem stat + directory listing only, no file content reads beyond manifest parsing already done by dependency detection)"

50ms 是一个合理的目标，但提案没有分析现有检测管道的基础耗时。当前 `DetectSurfacesWithConflicts` 已经执行了文件系统遍历（`os.ReadDir`、`os.ReadFile` 读取 go.mod/package.json 等）。在大型 monorepo 中（如 forge 自身的目录结构），目录遍历可能已经消耗了 10-20ms。推断函数是在依赖信号返回空之后才调用的，但前提是依赖信号检测本身已经执行完毕（读取了所有 manifest 文件）。如果依赖检测耗时 30ms，推断函数只剩下 20ms 的预算。

更重要的是，50ms 预算是"所有推断函数组合"的预算，但提案声明推断函数只在依赖信号为空时才调用。这意味着推断函数是串行执行的（先检查 Go，再检查 Node.js，再检查 Python），还是并行？如果串行，每个函数的预算是多少？如果项目只有 `go.mod`，是否仍然调用 Node.js 和 Python 的推断函数？提案声称"no file content reads beyond manifest parsing already done by dependency detection"，但推断函数需要检查目录结构（如 `cmd/` 子目录），这是额外的 `os.ReadDir` 调用。

建议：将性能预算分拆为两部分：(1) 依赖信号检测（现有代码，已有基准），(2) 推断函数（新增代码，50ms 预算适用于此部分）。同时，推断函数应该有 short-circuit 逻辑：当项目的 manifest 文件类型已知时（如有 `go.mod` 就不需要检查 `package.json`），只调用对应生态系统的推断函数。

风险：用户覆盖行为的"source is TUI-only, not persisted"设计在 `forge surfaces detect` 子命令中会产生信息丢失。

> "User override: after editing an inferred `cli` to `api` in TUI, config contains `surfaces: api` with no source field; source annotation does not appear in serialized config"
> "Source information is display-only and discarded after TUI confirmation completes."

这个设计在 `forge init` 的首次运行中是合理的：用户确认后，来源不再重要，配置值就是最终值。但在 `forge surfaces detect` 的 re-detection 场景中，这个设计有问题：用户运行 `forge surfaces detect`，看到推断结果是 `cli (inferred)`，决定覆盖为 `api`，确认后配置变成 `surfaces: api`。下次运行 `forge surfaces detect` 时，系统无法区分这个 `api` 是来自"用户手动覆盖"还是"之前的依赖检测"还是"之前的推断"。如果此时项目依赖发生了变化（比如新增了 Flask），检测管道会返回 `api (detected from flask)`。TUI 显示 "Surfaces already configured: api. Re-detect?"，用户选择 Re-detect，看到新的检测结果是 `api (detected from flask)` — 这看起来是正确的，但用户无法知道之前的 `api` 是否就是自己手动设置的。

更严重的场景：用户首次推断得到 `cli`，覆盖为 `web`，配置写入 `surfaces: web`。之后项目结构变化，推断规则现在会返回 `api`。用户运行 `forge surfaces detect` 看到推断结果是 `api (inferred)`，但现有配置是 `web`。TUI 提示 "Surfaces already configured: web. Re-detect?" — 如果用户不记得当初为什么设置 `web`，可能盲目选择 Re-detect 覆盖掉自己之前的审慎决定。如果 config 中保留了来源信息（如注释 `# source: user override, was: cli (inferred)`），用户就能做出更明智的决定。

建议：在 config 中以 YAML 注释形式保留来源元数据。注释不会被程序化读取，不增加 schema 复杂度，但为人类读者提供了上下文。例如：`surfaces: api  # was: cli (inferred:cmd-dir), overridden by user on 2026-05-24`。

问题：Node.js 推断规则中 `index.html` 检测为 `web` 的规则没有限定条件。

> "Node.js minimal project: `package.json` exists but no framework deps -> `bin` field -> `cli`; `index.html` at root -> `web`"

`index.html` 存在于几乎所有 Node.js 项目的某个子目录中（如 `public/index.html`、`dist/index.html`、`node_modules/*/index.html`）。提案说"at root"，但 `detectSurfaceAtDirWithConflicts` 的调用链是通过 `scanSubdirsWithConflicts` 递归扫描的。当递归到子目录时，如果子目录包含 `index.html`，也会触发 `web` 推断。例如，一个 Express 项目有 `public/index.html`，在深度扫描时 `public/` 子目录会被检测为 `web` surface — 这显然是错误的。

提案应该在 Node.js 推断规则中明确 `index.html` 检测只在项目根目录（与 `package.json` 同级）生效，不应在子目录扫描中触发。这需要在 `inferNodeSurface` 中加入路径深度检查。

风险：非交互模式下 `forge surfaces detect` 的 exit code 语义与 `forge init` 不一致。

> "`forge surfaces detect` in non-interactive terminal: prints results to stdout, no TUI, no config write; exit code 0 if detection succeeds, 1 if no surfaces found"

`forge init` 在非交互模式下跳过 surface 配置并返回 `SKIPPED surfaces (non-interactive terminal)`，exit code 始终为 0。但 `forge surfaces detect` 在检测不到 surfaces 时返回 exit code 1。这种不一致意味着同一个检测逻辑在两个入口点有不同的退出码语义，对脚本消费者来说不可预测。

## Improvement Suggestions

建议：将 `forge surfaces detect` 重新设计为默认只读命令，通过显式 flag 控制写入行为。

> "runs detection + inference and displays results, independent of `forge init`. Gives users an explicit entry point for re-detection without re-running the full init flow."

当前设计让 `forge surfaces detect` 在交互式终端中自动写入配置，这违反了 CLI 工具中查询/检测命令应该是只读的惯例。更好的设计是：默认行为等同于 `--dry-run`（只显示结果），添加 `--apply` flag 启用配置写入。`--dry-run` 变为默认行为的冗余别名，可以保留作为文档化用途。这样 `forge surfaces detect` 在任何终端模式下都是只读的，写入行为需要显式 opt-in。这与 `forge surfaces` 的其他子命令（list、query、types）的只读语义保持一致。

建议：为 re-run 行为增加 "Edit" 选项，与现有的 `askMapConfirmation` 四选项设计保持一致。

> "TUI shows `"Surfaces already configured: cli. Re-detect?"` with Confirm (keep existing) / Re-detect options"

将二选项改为三选项：Confirm（保持现有）/ Re-detect（重新检测）/ Edit（手动输入）。Edit 选项直接进入 `manualSurfaceEntry` 流程，让用户可以绕过检测管道直接修正配置。这解决了"检测管道返回相同错误结果而用户无法手动修正"的场景。三选项的设计也可以通过 `huh.NewSelect` 实现，与 `askMapAction` 的实现模式一致。

建议：为 `Sources` map 的 value 格式建立明确的枚举类型或格式规范，而非使用字符串约定。

> "`Sources` map correctly populated: `{"forge-cli/cli": "inference:cmd-dir"}` for inferred paths, `{"forge-cli/cli": "dependency:cobra"}` for detected paths"

当前的冒号分隔字符串格式缺乏形式化定义。建议在 `detect_surface.go` 中定义常量：`const SourceDependency = "dependency"` 和 `const SourceInference = "inference"`，value 格式定义为 `fmt.Sprintf("%s:%s", SourceInference, ruleID)`。这样下游消费者可以用 `strings.SplitN(source, ":", 2)` 安全解析，而格式变更只需要修改一处。同时，`Sources` map 的 key 应该明确定义为与 `Surfaces` map 相同的路径 key（而非 surface type），修正 success criteria 中的不一致。

建议：在 config 中以 YAML 注释形式保留 surface 来源元数据。

> "Source information is display-only and discarded after TUI confirmation completes."

当前设计在确认后丢弃来源信息，这对 `forge surfaces detect` 的重复调用场景不利。建议在 `writeConfigFile` 写入 surfaces 时，如果 `Sources` map 非空，在同一行追加 YAML 注释：`surfaces: api  # source: inference:api-dir, confirmed by user`。这不需要修改 `SurfacesMap` 的序列化逻辑（注释在反序列化时被忽略），也不增加 schema 复杂度，但为后续 re-detection 提供了人类可读的上下文。实现方式：在 `MarshalYAML` 之后对输出字符串做后处理，插入注释行。

建议：明确 `forge surfaces detect` 与 `forge init` 的 surface 步骤之间的代码复用边界。

> "Still small scope — task 5 reuses the detection + TUI infrastructure from tasks 1-4."

提案声明 task 5（`forge surfaces detect`）复用 task 1-4 的基础设施，但没有分析复用的边界。具体而言：`askSurfaceConfirmation` 当前被 `runSurfaceConfig`（init 流程）调用，它内部调用 `DetectSurfacesWithConflicts` 并返回 `(SurfacesMap, bool)`。`forge surfaces detect` 需要：1) 调用同一个检测函数，2) 使用同一个 TUI 确认流程，3) 但有不同的 re-run 逻辑（`forge surfaces detect` 总是执行检测，而 init 有 "already configured" 检查）。这意味着 `askSurfaceConfirmation` 需要重构为两个部分：检测（返回 `DetectResult` + `Sources`）和确认（接收 `DetectResult` + `Sources`，执行 TUI）。提案应该在技术方向中明确这个重构，否则 task 5 的实现会导致 `askSurfaceConfirmation` 出现条件分支（init 路径 vs detect 路径），增加复杂度。

建议：为 Node.js `index.html` 推断规则增加路径深度守卫。

> "`index.html` at root -> `web`"

在 `inferNodeSurface` 中，`index.html` 检测应该只在 `dir == projectRoot` 时生效，不应在子目录扫描中触发。具体实现：`inferNodeSurface` 接收一个 `isRoot bool` 参数（或比较 `dir` 与 project root），只有当 `isRoot == true` 时才检查 `index.html`。`bin` 字段检测不受此限制，因为 `bin` 只存在于 `package.json` 中，而子目录的 `package.json` 通常代表 workspace 成员（此时 workspace 模式会分别处理）。
