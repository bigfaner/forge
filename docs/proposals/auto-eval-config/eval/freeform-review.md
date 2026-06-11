# Freeform Expert Review

**Reviewer**: Developer Tooling & Configuration Architect
**Document**: `docs/proposals/auto-eval-config/proposal.md`
**Date**: 2026-05-26

---

## Background Assessment

This proposal tackles two distinct but coupled problems in the Forge CLI toolchain. The first is a developer experience friction: four document evaluation skills (`eval-proposal`, `eval-prd`, `eval-ui`, `eval-design`) each require manual `AskUserQuestion` confirmation before running, adding interaction cost to what should be an automated pipeline. The second is an architectural bottleneck: the `forge config get/set` routing system relies on hardcoded dispatchers (`autoModeField` switch/case, `getAutoKeyValue`, `setAutoConfigValue`, etc.) that only support two levels of key depth, making any nested configuration addition require routing-core changes.

The proposal's insight is to combine both problems into a single change: generalize the config key resolution to arbitrary depth (via Go reflection for get and YAML Node traversal for set), and then use that generalized infrastructure to add `auto.eval.*` as a nested struct with four ModeToggle fields. This is a sound architectural move -- it eliminates a class of extensibility problems rather than working around one instance.

After reading the current `config.go` implementation, I confirm the diagnosis is accurate. The `autoModeField` function (lines 515-531) uses a literal switch/case for six fields. The `getAutoKeyValue` function (lines 686-732) uses `strings.Index` to split at most one dot. The `parseAutoRaw` function (lines 357-407) hardcodes `modeFields := []string{"test", "consolidateSpecs", ...}`. The `setAutoConfigValue` function (lines 580-633) uses `strings.SplitN(rest, ".", 2)`. Each of these is a genuine extensibility bottleneck that the proposal correctly identifies and proposes to eliminate.

The proposal also correctly notes that `ui-design` is the only skill that unconditionally runs its eval, while the other three use `AskUserQuestion`. I verified this against the actual skill files: `brainstorm/SKILL.md` line 117 uses `AskUserQuestion` for eval-proposal; `write-prd/SKILL.md` line 215 uses `AskUserQuestion` for eval-prd; `tech-design/SKILL.md` line 170 uses `AskUserQuestion` for eval-design; while `ui-design/SKILL.md` line 136 unconditionally invokes `/eval-ui`. This behavioral asymmetry is real.

The default values proposed are backward-compatible: `proposal` defaults to `quick:true, full:true` (always auto-run, catches issues early); `uiDesign` defaults to `quick:true, full:true` (preserves current unconditional auto-run behavior); `prd` and `techDesign` default to `quick:false, full:false` (preserves current ask-first behavior). This is the correct strategy for zero-regression deployment.

## Key Risks

风险：Part 1 的泛化重写与 Part 2 的 eval 配置存在耦合交付风险
> "预计 3-4 小时：泛化路由重写（2h）+ eval 配置 + 4 个 skill（1h）+ 测试（1h）" — Part 1（泛化 key resolution）是 Part 2（eval 配置）的前置依赖。如果 Part 1 的反射遍历或 YAML Node set 路径存在 bug，Part 2 的 eval 配置将无法工作。这意味着两个 part 必须在同一 PR 中交付且同时通过测试。然而 Part 1 是对 config 系统核心路由的重写，影响所有现有 `forge config get/set` 路径（`auto.*`、`worktree.*`、`coverage.*`、`test-framework`），回归风险不限于 eval 功能。如果 Part 1 出现问题，eval 配置无法单独回滚。提案没有讨论增量交付策略（例如先合并 Part 1 并通过回归测试，再在下一个 PR 中添加 Part 2）。

风险：反射遍历对 `CoverageConfig.ByType`（`map[string]CoverageStrategy`）的处理需要特殊分支
> "反射遍历对 map 类型（SurfacesMap、CoverageConfig.ByType）的处理边界" — 在 Key Risks 表格中被列为 M/M，但实际复杂度被低估。当前 `getCoverageKeyValue` 的语义是 `coverage.coding.feature` 返回策略值（"80" 或 "maintain"），这不是简单的 map key 查找 -- 它需要理解 `CoverageStrategy` struct 的内部结构（`Type` + `Percentage`）。泛化的反射遍历器在遇到 `map[string]CoverageStrategy` 时，需要知道 `coding.feature` 是 map key 而非 struct field，然后还需要知道如何将 `CoverageStrategy` 序列化为字符串。提案提到"为 map 类型写专门的处理分支"，但这意味着反射遍历器内部需要一个类型特化注册表，否则 map 的 value 类型序列化逻辑无法通用化。这不是一个简单的 switch 分支能解决的。

风险：YAML Node set 路径的序列化保真度问题未充分探索
> "YAML Node set 操作的序列化保真度（注释丢失、格式变化）" — 在 Key Risks 表格中被评为 L/L，但这取决于具体实现。当前 `SetConfigValue` 的所有 set 路径（`setAutoConfigValue`、`setWorktreeConfigValue`、`setCoverageConfigValue`）都通过 Go struct marshal 再写回文件（`yaml.Marshal(cfg)` 写入），这会完全重新格式化 YAML。如果泛化的 set 路径改用 YAML Node 操作（直接在 Node 树上修改并序列化），格式保真度取决于 `yaml.Node` 的序列化行为。Go 的 `gopkg.in/yaml.v3` 在 Node 序列化时保留注释和格式的能力有限 -- 映射节点的键顺序可能改变，行间注释可能丢失。提案需要明确选择：set 路径是走 struct marshal（当前模式，格式会被重排但行为已知）还是 YAML Node 修改（格式更好保留但行为需要验证）。

问题：反射 get 路径的 YAML tag 匹配机制需要明确定义
> "将 `key` 按 `.` 拆分为路径段，沿 Go Config struct 树递归走 reflect.Value...struct 按 YAML tag 匹配字段" — 当前 `AutoConfig` 的字段使用 `yaml:"consolidateSpecs"` 等 tag，但 `WorktreeConfig` 使用 `yaml:"source-branch"`（含连字符）和 `yaml:"copy-files"`。这意味着反射遍历器在匹配 `worktree.source-branch` 时，需要将 `source-branch` 与 `SourceBranch` 字段的 YAML tag `"source-branch"` 匹配，而不是与 Go 字段名匹配。这本身不难实现，但 `Config` struct 的 `Auto` 字段是 `*AutoConfig`（指针），`Worktree` 是 `*WorktreeConfig`（指针），反射遍历器需要正确解引用指针。对于 nil 指针（用户未配置 worktree），`get` 应返回 `errKeyNotFound`，`set` 需要先初始化指针再继续遍历。提案的函数签名 `getStructValueByPath(v reflect.Value, segments []string)` 使用 `reflect.Value` 而非 pointer，暗示调用者负责解引用，但这个约定未显式说明。

问题：`parseAutoRaw` 泛化后的 `raw` 追踪粒度对嵌套结构体的影响
> "泛化 `parseAutoRaw`：递归扫描 auto 子树，不再硬编码 modeFields" — 当前 `parseAutoRaw` 返回 `map[string]map[string]bool`，其中外层 key 是字段名（如 `"test"`），内层 key 是子字段名（如 `"quick"`、`"full"`）。对于嵌套结构体 `auto.eval.proposal`，如果 `raw` 的 key 仍然是扁平路径（`"eval.proposal"` → `{"quick": true, "full": true}`），那么 `applyDefaults` 需要知道哪些 key 是嵌套的、哪些是扁平的。如果改为嵌套结构（`"eval"` → `{"proposal": {"quick": true}}`），`applyModeDefault` 的签名需要改变。提案说"递归扫描保持相同的叶子节点追踪粒度"，但没有定义 `raw` map 的 key 格式将如何变化。这是默认值填充逻辑的核心数据结构，未定义则无法实现。

风险：skill 中 mode 检测的实现依赖 manifest 文件格式假设
> "通过 manifest 文件检测（`docs/features/<slug>/manifest.md` 中的 `mode: quick`）判断当前管道模式" — 我检查了 `feature_complete.go` 的实现（第 105-109 行），它通过检查 `proposal.md` 是否存在于 feature 目录来判断 quick mode（`quickMode := proposalErr == nil`），而非读取 manifest.md 的 `mode` 字段。提案声称"manifest 文件格式需包含 `mode` 字段（已存在于 `feature_complete.go`）"，但代码中并没有读取 manifest 的 mode 字段。这意味着提案中的 mode 检测机制实际上不存在于当前代码中，需要新建。skill markdown 中的 config check 需要一个可靠的方式判断当前是 quick 还是 full 模式 -- 如果这个判断依赖一个尚未实现的 API，整个 eval 配置的 quick/full 区分就是不可工作的。

问题：默认值选择中 `proposal` 默认 `true` 的语义合理性需要更多论证
> "proposal: `quick: true, full: true` — 默认自动运行（proposal 是管道入口，尽早发现问题）" — proposal 评估的质量取决于 proposal 本身的完成度。在 `/quick` 流水线中，brainstorm 产出的 proposal 通常比较粗糙（快速探索），立即运行一个 900 分阈值的对抗性评估可能产生大量迭代，反而增加而非减少交互成本。提案的 Assumption Flip 表格说"proposal 默认应询问用户"被 5 Whys 推翻了，但 5 Whys 的结论是"自动评估可尽早发现问题"，这假设评估结果是建设性的。如果 proposal 质量较低导致评估分数远低于阈值，用户需要反复修改 proposal 并重新评估，这比一次手动确认的开销更大。

建议关注：`forge config get auto.eval`（中间节点）的行为未定义
> 成功标准列出了 `auto.eval.proposal`（三段）和 `auto.eval.proposal.quick`（四段）的行为，但没有定义 `auto.eval` 本身的行为。如果用户输入 `forge config get auto.eval`，期望看到什么？如果返回错误，用户可能困惑（为什么不支持？）。如果返回整个 eval 块的汇总，格式是什么？同样，`forge config set auto.eval true` 是否等同于设置所有四个子字段？这些边界情况影响 CLI 的用户体验一致性，应在提案中明确。

## Improvement Suggestions

建议：将 Part 1 和 Part 2 拆分为独立可交付的 PR
> 对应 "泛化路由重写（2h）+ eval 配置 + 4 个 skill（1h）" — Part 1 的泛化路由重写是一个独立有价值的基础设施改进，可以在没有 eval 配置的情况下独立合并。建议先交付 Part 1 并确保所有现有 config_test.go 和 config_schema_test.go 的回归测试通过，再在第二个 PR 中添加 eval 配置和 skill 更新。这降低了单个 PR 的回归风险，也使得问题定位更容易。拆分后的 Part 1 PR 可以包含一个简单的验证：使用一个测试用的嵌套 struct 证明三层路径工作正常，而不需要立即引入 eval 的领域逻辑。

建议：为 `parseAutoRaw` 的泛化定义新的 `raw` 数据结构
> 对应 "泛化 `parseAutoRaw`：递归扫描 auto 子树" — 建议将 `raw` 的类型从 `map[string]map[string]bool` 改为 `map[string]any`（或等效的递归结构），使得 `raw["eval"]` 可以是一个嵌套 map `{"proposal": {"quick": true, "full": true}}`。同时在提案中给出这个数据结构的具体定义和示例，这样 `applyDefaults` 的实现者可以准确理解每个字段的追踪粒度。没有这个定义，"递归扫描保持相同的叶子节点追踪粒度"是一个无法验证的实现承诺。

建议：为 mode 检测提供 CLI 级别的 API
> 对应 "通过 manifest 文件检测判断当前管道模式" — 当前代码中没有统一的 mode 检测 API。建议在提案的 Scope 中新增一个 `forge config get mode` 或 `forge mode` 命令，返回当前工作目录对应的管道模式（quick/full/none）。这样 skill 中的 config check 逻辑可以简化为两步：（1）`forge config get mode` 获取当前模式；（2）`forge config get auto.eval.{skillName}.{mode}` 获取对应配置。这比让每个 skill 独立实现 manifest 解析更可靠，也避免了 skill markdown 中嵌入文件路径假设的脆弱性。

建议：重新考虑 `proposal` 的默认值
> 对应 "proposal: `quick: true, full: true`" — 建议将 proposal 的默认值改为 `quick: false, full: true`。在 quick 模式下，proposal 通常是快速探索产物，强制自动评估可能适得其反。在 full 模式下，proposal 经过更充分的讨论，自动评估的投入产出比更高。这与 `auto.test` 的默认值（`quick: false, full: true`）保持一致的哲学：quick 模式追求速度，跳过验证步骤；full 模式追求质量，启用验证步骤。

建议：在提案中补充 `getStructValueByPath` 和 `setYAMLValueByPath` 的完整签名和边界行为
> 对应 "Get 路径（反射）" 和 "Set 路径（YAML Node）" — 当前提案只给出了函数名和一行注释。建议补充：指针 nil 的处理方式（get 返回 errKeyNotFound？set 自动初始化？）、map 类型的 key 匹配策略、YAML tag 与 Go field name 的优先级、叶子节点类型不支持时的错误消息格式、以及中间节点查询的行为定义。这些是实现者在编写代码时需要做出的一连串决策，如果不在提案中锁定，实现结果可能与提案意图产生偏差。
