# Freeform Expert Review

**文档**: `docs/proposals/forge-workspace/proposal.md`
**评审角色**: Multi-Project CLI Workspace Architect
**日期**: 2026-06-07

---

## Background Assessment

本提案旨在解决一个真实且日益突出的痛点：当独立开发者同时推进 4-8 个 Forge 项目时，过程文档（proposals、features、tasks、PRDs）散落在各个项目目录内，缺乏统一的管理入口和全局视图。提案的核心思路是在项目父目录引入一个"Workspace"叠加层——通过 `.forge-workspace.yaml` 注册表聚合多个项目，提供跨项目状态总览、feature 聚合和 workspace 级 proposals。

技术路径上，提案选择了纯文件系统方案：扫描直接子目录检测 `.forge/config.yaml` 实现项目发现，读取各项目已有的 manifest 和 task 文件实现状态聚合，workspace 级 proposals 存放在 workspace 自己的 `docs/proposals/` 目录。整个方案坚守"项目不变、workspace 是纯叠加层"的原则。

提案的一个核心设计决策是将多项目管理拆解为三个正交模块：Workspace（过程文档）、Dashboard（可视化）、Wiki（知识管理）。三个模块通过 `.forge-workspace.yaml` 共享注册信息，但各模块配置独立、schema 独立演进。v1 只实现 Workspace 模块。

提案的假设基础可以概括为：（1）所有项目平铺在同一个父目录下，没有嵌套；（2）每个项目的 manifest 和 task 格式基本一致或可通过 schema 版本兼容；（3）4-8 个项目的规模下，纯文件系统扫描的性能可接受；（4）CLI 优先是正确的交互层级，Dashboard 可以推迟。

## Key Risks

### 项目发现策略对非平铺目录布局的脆弱性

提案反复强调"仅扫描直接子目录（一层深度）"的发现规则，并且明确将嵌套项目列入 out-of-scope 的拒绝条件。这本身是一个合理的简化，但问题在于提案没有讨论"非平铺布局"出现时的应对策略——不是功能支持，而是**优雅拒绝**时的用户体验。

风险：`init` 零发现时缺少诊断性反馈，导致用户误判为工具故障。
> "仅扫描直接子目录（一层深度），跳过符号链接、`.gitignore` 标记的目录、隐藏目录" — 当父目录下存在嵌套分组目录时，扫描结果为零，用户无法区分"工具坏了"和"目录结构不符合预期"。后果是首次使用体验即产生挫败感，可能直接放弃。

### `.forge-workspace.yaml` 作为共享配置文件的耦合风险

三模块架构的协作设计中，`.forge-workspace.yaml` 被定位为三个模块的共享注册表。这是一个关键的架构决策点。

风险：共享注册表文件缺乏字段增长约束，后续模块演进中将退化为全局配置 dump，破坏模块独立性。
> "三者通过 `.forge-workspace.yaml` 共享项目注册信息" — 虽然提案限定该文件"仅含项目路径列表和 schema 版本号"，但没有机制保证这个约束在 Dashboard 和 Wiki 模块开发时不被打破。一旦任何模块向共享文件添加专属字段，三个模块的 schema 变更开始相互影响，与"模块代码隔离、schema 独立演进"的设计原则矛盾。

### 状态聚合缓存的一致性语义不明确

提案引入了缓存机制来优化 `status` 命令的响应时间，但缓存失效策略的描述存在模糊地带。

问题：mtime 指纹的追踪粒度未定义，缓存策略缺乏可实现的规格说明。
> "首次扫描后缓存各项目状态快照至 `.forge-workspace/cache.json`（含 mtime 指纹）；后续调用仅重新扫描有文件变动的项目，全量刷新阈值 5 分钟" — "mtime 指纹"是追踪项目根目录？`docs/` 下所有文件？还是仅 manifest 文件？粒度过粗触发无效扫描，粒度过细 mtime 收集本身接近全量扫描开销。5 分钟窗口内连续创建多个 tasks 时，status 输出与实际状态不一致，且用户不知晓数据是过期的。

### Proposal 到 Feature 的 assign 流程存在上下文损失风险

这是整个提案中最关键的 seam——workspace 级 proposal 到 project 级 feature 的转换。

风险：assign 流程的"核心上下文"继承定义模糊，brainstorm 产出的关键决策可能在转换中丢失。
> "assign 将 proposal 的 title、intent、status、核心上下文自动写入目标项目的 feature manifest" — 一个经过 brainstorm 和 eval 的 proposal 包含问题陈述、方案选项、取舍分析、评审反馈等结构化信息。"核心上下文"是一个无法实现的模糊术语：继承太少则丢失决策依据，继承全部则因 workspace proposal 与 project feature 的 schema 差异导致字段映射错误。这需要显式的字段映射表，而非留给实现阶段。

### workspace 级与项目级 proposal 的认知模型边界不够清晰

提案试图通过命令命名空间来区分两个层级的 proposals。

问题：workspace 级 proposal 被 assign 后的生命周期缺少闭环，`Assigned` 状态可能成为僵尸记录。
> "CLI 命令明确区分：workspace 级用 `forge workspace propose`，项目级用 `/brainstorm`（在项目内运行）" — 命名空间本身清晰，但 assign 后的追踪路径断裂。用户应在 workspace 级还是项目级追踪后续状态？提案提到 proposal 状态更新为 `Assigned`，但未定义 `Done` 的判定时机。缺少从 feature 完成回写 workspace proposal 的机制（无论自动还是手动），proposal 列表将积累大量 `Assigned` 状态的僵尸记录。

### manifest 格式不一致的处理策略可能不足

提案预见了跨项目 manifest 格式漂移的风险，并提出了 schema 版本号的方案。

问题：现有项目 manifest 可能没有 `schema_version` 字段，聚合器的迁移路径未定义。
> "引入 manifest schema 版本号（`schema_version` 字段）；聚合器按版本分发解析逻辑，未知版本降级为纯文本摘要" — 当前各项目是否已有该字段？无版本号时视为 v0 还是格式未知？"降级为纯文本摘要"意味着这些项目在 `forge workspace features` 输出中丢失结构化信息，用户可能困惑为什么某些项目没有状态字段。降级应在 CLI 输出中有明确的视觉区分。

### "项目不变"原则下 brainstorm 技能在 workspace 上下文运行的矛盾

提案声称 workspace 是"纯叠加层"，现有项目结构完全不变。但同时又要求 brainstorm 技能在 workspace 上下文运行。

问题：brainstorm 技能需要感知 workspace 上下文，这暗示对现有技能的侵入性修改，与"纯叠加层"的直觉承诺不符。
> "brainstorm/eval 技能在 workspace 上下文运行" — 这意味着 brainstorm 技能需要新增"上下文模式"参数，调整输出路径（写到 workspace 的 `docs/proposals/` 而非项目的）。这是对现有技能的侵入性改动。提案未显式声明哪些技能需要修改、修改范围是什么，"项目不变"的声明可能给 reviewer 造成"现有代码零改动"的错误印象。

## Improvement Suggestions

建议：为 `forge workspace init` 的零发现场景增加诊断性输出。
Addresses: `init` 零发现时缺少诊断性反馈的风险。
> What changes: 当 `init` 扫描完成但未发现任何有效项目时，不应仅输出空注册表，而应输出诊断信息——扫描了多少子目录、被跳过的目录及原因分类（缺少 `.forge/config.yaml` / 是隐藏目录 / 是符号链接），以及建议检查项（"你的项目是否在嵌套子目录中？当前仅扫描一层深度"）。这不改变发现逻辑，只在输出层增加诊断，是低成本的体验改进。

建议：为 `.forge-workspace.yaml` 定义字段增长的治理约束——采用"注册表与配置分离"的硬性规则。
Addresses: 共享注册表文件缺乏字段增长约束的耦合风险。
> What changes: 在提案的"模块间隔离原则"部分，显式声明 `.forge-workspace.yaml` 只允许包含 `{ schema_version, projects: [...] }` 两个字段，任何模块专属配置一律存入 `.forge-workspace/<module>.yaml`。在 `init` 时通过 JSON Schema 验证该文件的结构合规性。这条规则一旦写入提案，就为后续 code review 提供了明确的判断依据。

建议：明确定义缓存策略的 mtime 追踪粒度和缓存失效时的用户提示。
Addresses: mtime 指纹追踪粒度未定义的问题。
> What changes: 将 mtime 追踪粒度从"含 mtime 指纹"细化为具体策略：建议追踪每个项目 `docs/` 目录下所有 `.md` 文件的最新 mtime 作为该项目的脏标记——平衡精确性和成本。当 `status` 输出使用了缓存数据时，在表格底部增加提示："Data cached Xs ago. Use --fresh for live scan."。5 分钟阈值应作为 `.forge-workspace/config.yaml` 中的可配置项。

建议：为 assign 流程定义显式的字段映射表，替代模糊的"核心上下文"。
Addresses: assign 流程"核心上下文"继承定义模糊的风险。
> What changes: 在 v1 scope 部分增加字段映射表：`title -> title`、`intent -> intent`、`status -> status`（重置为 Draft 或保持 Approved）、`body（problem/solution） -> description`、`eval_results -> design_context`。未被映射的字段（alternatives、assumptions）通过 `source: workspace-proposal:<slug>` 链接保留可追溯性。这个映射表本身就是 assign 功能的核心设计文档，不应推迟到实现阶段。

建议：为 workspace 级 proposal 的生命周期定义完整的闭环规则。
Addresses: workspace 级 proposal `Assigned` 状态可能成为僵尸记录的问题。
> What changes: 补充 proposal 状态机完整定义，特别是 `Assigned -> Done` 的转换条件。建议规则：assign 产出的 feature 状态变为 `Done` 时，workspace 级 proposal 自动更新为 `Done`，记录 `completed_at`。需要 proposal frontmatter 维护 `assigned_to` 和 `feature_slug` 字段以支持状态同步。如果自动回写对 v1 过于复杂，至少提供 `forge workspace close <proposal>` 手动关闭命令。

建议：显式列出需要为 workspace 上下文修改的现有技能清单及修改范围。
Addresses: brainstorm 技能需要侵入性修改但未声明的问题。
> What changes: 在 Constraints & Dependencies 或 Scope 部分增加小节，列出 v1 需要修改的现有技能：至少包括 `brainstorm`（支持 workspace 上下文模式、输出路径可配置）和 `eval`（如果 workspace 级 proposal 也需评审）。对每个技能说明修改范围（新增参数 vs 内部路径重构）。这不会改变方案设计，但让实现阶段的范围估算更准确，避免"项目不变"声明造成零改动的错误预期。
