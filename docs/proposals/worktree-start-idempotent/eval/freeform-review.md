# Freeform Expert Review

**Expert**: Developer Tooling & Configuration Architect
**Document**: `proposal.md` — Make `forge worktree start` Idempotent & Rename `copy-files` → `includes`
**Date**: 2026-06-09

---

## Background Assessment

这份提案旨在解决两个开发者体验问题。第一个是 `forge worktree start <slug>` 命令在 worktree 目录已存在时直接报错退出，无法用同一条命令「进入已有 worktree 并启动全新 Claude 会话」。提案建议将 `start` 改为幂等操作——worktree 不存在则创建并启动新会话，已存在则跳过创建直接启动新会话。这个设计的核心洞察是 `start` 的用户语义应该是「开始工作」而非「创建 worktree」，与 `mkdir -p` 的幂等模式对齐。

第二个问题是将配置项 `worktree.copy-files` 重命名为 `worktree.includes`，使其描述意图（"包含什么"）而非实现行为（"复制文件"），并与 Claude Code 的 `.worktreeinclude` 命名风格保持一致。提案明确声明不做旧字段兼容——直接替换。

提案还详细列出了边界场景的处理策略：`--source-branch` 在 worktree 已存在时被忽略并输出 warning；`includes` 文件复制在已存在时被跳过；`--no-launch` 仅验证路径。

在技术实现层面，提案预估改动范围为 2-3 个文件、约 40 行变更，依赖 `cmd_resume.go` 中已有的 worktree 验证逻辑，无外部依赖。

提案的基本假设是：用户语义的 `start` 等同于「开始工作」，而非「创建资源」；幂等行为符合最小惊讶原则；配置重命名不需要迁移路径。

---

## Key Risks

提案在整体方向上是合理的——将 CLI 命令幂等化是成熟的工具设计模式。但存在若干值得深入分析的风险和问题。

风险：配置项重命名不保留旧字段兼容逻辑，将导致现有用户配置静默失效。
> "直接替换，不保留任何旧字段兼容逻辑" — 这是整个提案中最大的风险点。`worktree.copy-files` 一旦被重命名为 `worktree.includes`，所有已使用 `copy-files` 的用户配置文件在升级后会被静默忽略（YAML 解析器通常对未知字段不报错）。用户可能不会立即发现问题——因为 worktree 创建不会报错，只是文件不再被复制。这种「静默失败」模式在配置系统设计中是尤其危险的，因为它违背了开发者工具最基本的原则：当行为发生变更时，必须让用户知道。正确的做法至少应该是在检测到旧字段时输出 deprecation warning，或者提供一个版本的兼容期。即使是 "breaking change" 也可以接受，但必须是 "loud breaking change" 而非 "silent breaking change"。

问题：提案将配置重命名和命令行为变更捆绑在同一个变更中，但没有论证为什么要同时做这两件事。
> 提案的标题同时包含"幂等化 start"和"重命名 copy-files"两个独立变更。在 Section "Constraints & Dependencies" 中列出的修改文件同时涉及命令逻辑和配置结构体。从配置架构的角度看，这两个变更是正交的——幂等化只涉及控制流逻辑，重命名只涉及配置 schema。捆绑在一起会增加变更的风险表面积：如果重命名导致配置读取问题，用户可能误以为是幂等逻辑的 bug。提案没有解释为什么不能分两步发布（先幂等化，再重命名），也没有讨论这种捆绑的风险缓解策略。

风险：`--source-branch` 被忽略的场景可能导致用户在不知情的情况下使用了错误的分支。
> "已有 worktree 但用了 `--source-branch` → 忽略该 flag（worktree 已存在，分支已确定）" — 这个场景有一个微妙的假设：用户记得这个 worktree 是基于哪个分支创建的。在实际使用中，用户可能在很久之前创建了 worktree，现在运行 `forge worktree start feature-foo --source-branch hotfix-bar` 时，期望进入的是 hotfix-bar 分支上的工作环境，但实际上进入了 feature-foo 原来的分支。仅输出 warning 可能不够——用户可能不会注意到终端中混在大量输出里的一行 warning。提案应该考虑：是否应该在 `--source-branch` 指定的分支与 worktree 实际分支不同时，输出更显眼的提示（甚至可以考虑非零退出码配合 `--no-launch` 使用）？

风险：`includes` 文件跳过策略假设「首次创建时已复制」始终成立，但没有考虑 worktree 内容可能被损坏或过时的情况。
> "首次创建时已复制，后续不应再覆盖（会丢失 worktree 内的修改）" — 这个假设在理想情况下成立，但现实场景中 worktree 可能因为 git 操作（如 force push 后的 rebase 冲突残留）或手动修改而处于不一致状态。如果用户希望「重新同步 includes 文件到最新状态」，当前提案没有提供任何路径。建议考虑增加一个 `--reinclude` 或类似的 flag，允许用户在已有 worktree 中重新应用 includes 逻辑。

问题：提案声称改动范围是 "2-3 个文件，约 40 行变更"，但列出的约束与依赖暗示实际影响范围可能更大。
> "修改 `forge-cli/internal/cmd/worktree/cmd_start.go`（幂等逻辑）" 和 "修改 `forge-cli/pkg/forgeconfig/` 中的配置结构体（`CopyFiles` → `Includes` + 兼容读取）" — 注意这里出现了"兼容读取"的字样，但在 In Scope 部分又明确说"直接替换，不保留旧字段"。这两处描述是矛盾的。如果真的做"直接替换"，那么所有引用 `CopyFiles` 的代码（不仅是配置结构体，还包括读取逻辑、测试代码、可能的其他命令）都需要同步修改。实际文件数可能远超 2-3 个。此外，Success Criteria 中要求"代码中不存在任何 `copy-files` / `CopyFiles` 兼容逻辑"，这意味着需要全量搜索替换，改动范围的可控性需要更仔细的评估。

问题：提案没有讨论 `start` 在幂等模式下与 `resume` 命令的边界清晰度问题。
> "`start` 成为「全新会话」的统一入口，`resume` 保持为「恢复旧会话」的专用命令" — 这个语义划分在理论上清晰，但在实际使用中可能造成困惑。当前 `start` 报错时消息是 `worktree already exists, use "resume" instead`，这个报错实际上起到了「引导用户了解 resume 命令」的作用。改为幂等后，用户可能永远不需要了解 `resume` 命令的存在（因为 `start` 总是能用），导致 `resume` 成为一个「隐藏命令」。这不是一个严重问题，但提案应该明确这个 UX 后果，并考虑是否需要在 `start` 进入已有 worktree 时提示「要恢复上次会话，请使用 resume 命令」。

风险：提案中的 Non-Functional Requirements 部分过于简略，遗漏了对开发者工具来说至关重要的可观测性需求。
> "向后兼容：`start` 在 worktree 不存在时的行为完全不变" 和 "性能：无影响（只是跳过了创建步骤）" — 提案只列了这两条非功能需求，但缺少了「日志/输出的可观测性」这一关键需求。在配置系统设计中，当一个命令的行为从「报错」变为「正常执行但走了不同路径」时，用户必须有明确的方式知道发生了什么。虽然 In Scope 中提到了"输出区分性日志"，但这应该在非功能需求中被正式化，明确日志的格式、级别和内容要求。

---

## Improvement Suggestions

建议：为配置重命名增加「loud deprecation」策略，而非静默替换。
Addresses: 配置项静默失效风险
> 具体变更：在第一个版本中，同时支持 `copy-files` 和 `includes`，但当检测到 `copy-files` 时输出醒目的 deprecation warning（如 `WARNING: 'worktree.copy-files' is deprecated, use 'worktree.includes' instead. Support for 'copy-files' will be removed in v4.0`）。在下一个大版本（如 v4.0）中再移除 `copy-files` 支持。如果确实决定在一个版本内 breaking，则应该在检测到旧字段时输出错误信息并退出，而非静默忽略。这样做的好处是：用户不会在不知情的情况下丢失功能，且迁移路径清晰明确。

建议：将配置重命名和幂等化拆分为两个独立的变更步骤。
Addresses: 捆绑变更的风险表面积
> 具体变更：先发布幂等化 `start` 的变更（只改 `cmd_start.go`），验证稳定后再发布配置重命名变更。这两个变更没有技术依赖关系——幂等逻辑不关心配置项叫什么名字。拆分后，如果出现问题，可以精确定位是哪个变更引起的，降低回滚成本。在提案中明确标注这两个变更的发布顺序和各自的版本号。

建议：当 `--source-branch` 被忽略且指定的分支与 worktree 实际分支不同时，输出更显眼的提示。
Addresses: `--source-branch` 被静默忽略的风险
> 具体变更：在检测到 worktree 已存在时，比较 `--source-branch` 指定的分支与 worktree 的实际分支。如果不同，输出类似 `WARNING: worktree 'feature-foo' is on branch 'main', but --source-branch 'hotfix-bar' was specified. The existing worktree will be used as-is.` 的提示，使用醒目的颜色（如黄色或红色前缀）。如果相同，可以只输出 info 级别的提示。这确保用户不会在不知情的情况下进入了错误的分支环境。

建议：增加 `--reinclude` flag 允许用户在已有 worktree 中重新应用 includes 文件。
Addresses: includes 文件跳过策略的局限性
> 具体变更：在 `start` 命令中增加 `--reinclude` flag。当 worktree 已存在且指定了 `--reinclude` 时，重新执行文件复制逻辑（覆盖 includes 中列出的文件）。这个 flag 的默认行为仍然是跳过，不会破坏幂等语义。这为用户提供了一个「重置」路径，当 worktree 的 includes 文件过时或被意外修改时可以恢复。

建议：修正文档中关于"兼容读取"的矛盾描述，并做更准确的改动范围评估。
Addresses: 文档中"兼容读取"与"直接替换"的矛盾
> 具体变更：在 Constraints & Dependencies 部分，将 `（CopyFiles → Includes + 兼容读取）` 改为 `（CopyFiles → Includes，全量替换）`，或者根据实际采用的策略（如果有兼容期则为兼容读取，如果直接 breaking 则为全量替换）保持一致。同时，重新评估实际涉及的文件数——使用 `grep -r CopyFiles` 或 `grep -r copy-files` 搜索整个代码库，列出所有需要修改的文件，而非笼统地说"2-3 个文件"。

建议：在非功能需求中增加「可观测性」条目，明确日志输出要求。
Addresses: 非功能需求遗漏可观测性
> 具体变更：在 Non-Functional Requirements 中增加：`可观测性：start 命令在两种路径（新建 vs 进入已有）下都必须输出明确的区分性信息。日志格式应包含动作关键词（如 "Created new worktree" vs "Entering existing worktree"），便于用户确认执行路径和脚本解析。` 这将一个散落在 In Scope 中的要求正式化为非功能约束，确保实现时不会被忽略。
