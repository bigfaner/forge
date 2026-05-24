# Freeform Expert Review

**Proposal**: Unify interfaces and surface into surfaces
**Reviewer**: Config-Schema & Surface-Detection Engineer
**Date**: 2026-05-24

---

## Background Assessment

This proposal aims to consolidate two overlapping concepts in Forge's test pipeline: `interfaces` (a project-level string array in config) and `surface` (a journey-level scalar, used in gen-journeys SKILL.md). Both accept the same value domain -- web, api, cli, tui, mobile -- but are set by different mechanisms (manual config vs. LLM-driven detection) and operate at different granularities. The proposal documents four concrete bugs stemming from this overlap, including a known silent failure mode (empty `interfaces` silently skips all test task generation) and a naming inconsistency between config schema values (`web-ui`/`mobile-ui`) and Go code expectations (`web`/`mobile`).

The core technical approach is to replace both fields with a single `surfaces` map in `.forge/config.yaml`, where keys are filesystem paths relative to project root and values are surface type strings. During `forge init`, Go code performs file-pattern-based detection and populates this map. A new independent CLI command `forge surfaces` provides query access, using longest-prefix-match to resolve arbitrary paths to their surface type. The naming convention normalizes to the short form (`web`, `mobile`) used by the existing Go code, dropping the `-ui` suffix from config schema values.

The proposal rests on several key assumptions: (1) that path-based mapping is more precise than type enumeration alone, (2) that Go-based file-pattern detection at init time is sufficient for most real-world projects, (3) that downstream skills can safely be adapted in later iterations without immediate breaking changes, and (4) that migration from old `interfaces` can be handled implicitly by re-running init rather than explicit schema migration.

Having traced the actual code paths, I can confirm the problem is real. `ReadInterfaces` in `detect.go` is a pure config read with no detection fallback. `BuildIndex` in `build.go` calls `ReadInterfaces` and feeds the result to `extractBodyContext` and `generateTestTasks`. The `uiInterfaces` map in `autogen.go` indeed uses short names (`web`, `mobile`) while config historically allowed `web-ui`/`mobile-ui`. The `surface` field referenced in gen-journeys SKILL.md has no corresponding Go struct field. The proposal accurately describes the current state.

---

## Key Risks

The proposal is well-structured for a problem statement, but as a config-schema engineer I see several areas where the data flow and migration story need sharpening before this moves to PRD.

### Detection signal table completeness and conflict resolution

The signal table lists 14 detection rules across 4 package managers plus mobile manifests. This is a reasonable starting set, but several real-world scenarios are unaddressed.

风险：信号冲突时缺少优先级规则
> "检测输出直接就是 path → surface 的 map，无需额外转换" — 提案暗示检测结果是无歧义的。但同一个 `package.json` 可能同时包含 `react`（web 信号）和 `express`（api 信号）。gen-journeys 的 surface-webui.md 已经识别了这种场景并给出了消歧规则（"如果用户交互主要通过浏览器，分类为 WebUI"），但提案的 Go 检测表没有定义类似的消歧逻辑。Go 代码无法执行"判断主要用户交互模式"这种语义推理，因此需要一个确定性的优先级规则表（例如：前端框架优先于后端框架），否则 init 检测会在 monorepo 的根 `package.json` 处产生不一致的结果。

问题：monorepo 根目录的 `package.json` 覆盖了子目录
> "路径信号 | `package.json` | react/vue/svelte + DOM entry | `web`" — 检测表按 package.json 位置定义路径，但在 pnpm workspace 或 yarn workspaces 的 monorepo 中，根目录的 `package.json` 只声明 workspace 配置，而实际的前端/后端依赖在子目录的 `package.json` 中。如果检测算法只扫描一层目录或固定深度，它会遗漏子目录中的信号。提案没有说明检测的目录遍历策略（深度限制？排除规则？`node_modules` 如何处理？）。

### Longest-prefix-match semantics and edge cases

The `forge surfaces <path>` command uses longest-prefix-match, which is the right general approach but has subtle failure modes.

风险：前缀匹配在路径规范化后的不一致行为
> "最长前缀匹配，无匹配则报错提示手动指定" — 配置中键的格式是"无前导 `./`，无尾随 `/`"，但 gen-journeys 调用 `forge surfaces <path>` 时传入的路径可能来自 journey 文件中的相对路径、LLM 推断的路径、或用户输入。如果传入 `./frontend` 或 `frontend/` 或 `frontend/src`，规范化逻辑必须在 Go 侧实现并测试。提案没有定义路径规范化的精确规则（是否解析 `..`？是否解析符号链接？Windows 路径分隔符如何处理？）。

问题：路径匹配的测试覆盖面没有声明
> "`forge surfaces <path>` 用最长前缀匹配" — 最长前缀匹配是经典的 trie/路由问题，边界情况很多。配置中有 `frontend: web` 和 `frontend/api: api`，查询 `frontend/api/routes` 应匹配 `frontend/api`，但查询 `frontend-new` 不应匹配 `frontend`（前缀匹配不是子串匹配）。提案没有声明匹配是按路径段（path segment）匹配还是按字符前缀匹配。按字符前缀匹配会导致 `frontend` 错误匹配 `frontend-new`。

### Backward compatibility and migration

The proposal handles migration optimistically, relying on init re-detection rather than explicit schema migration.

风险：旧 `interfaces` 字段的静默忽略可能导致信息丢失
> "旧 config 自动迁移工具 — init 重新检测会覆盖，旧字段可忽略" — 如果用户已经手动配置了 `interfaces: ["api", "cli"]`，并且不重新运行 `forge init`（例如在 CI 环境或已有的项目中），旧的 `interfaces` 字段会被 Go 代码完全忽略（因为 `Config` struct 将从 `interfaces` 切换到 `surfaces` map）。这意味着 `forge task index` 会静默停止生成测试任务，复现提案中提到的已知 bug。提案声称"向后兼容"但实际方案是"不兼容，但 init 会重新检测"，这两者不是同一回事。

问题：`interfaces` 到 `surfaces` 的语义映射不是一对一的
> "config schema 变更：删除 `interfaces` 和 `surface`，新增 `surfaces` map 字段（path → surface）" — 旧的 `interfaces` 是去重的类型列表 `["api", "cli"]`，没有路径信息。迁移时无法自动生成路径映射。如果提案选择不迁移而是要求重新 init，那么对于已有的、不处于 init 状态的项目（CI、已有 worktree），需要至少一个兼容读取的过渡期：当 `surfaces` 不存在时，回退读 `interfaces` 作为扁平类型列表（映射为 `".": <type>` 或者直接提取去重类型）。提案没有描述这个过渡期策略。

### Schema structure and data flow integrity

问题：`Config` struct 变更缺少字段类型声明
> "键：相对于项目根目录的路径（无前导 `./`，无尾随 `/`）" — 提案展示了 YAML 示例但没有给出 Go struct 定义。当前的 `Config` struct 用 `Interfaces []string`，新字段应该是 `Surfaces map[string]string`。提案应该声明精确的 struct 变更，包括 YAML tag（`yaml:"surfaces,omitempty"` 还是 `yaml:"surfaces"`？`omitempty` 会导致空 map 被丢弃，`forge task index` 又回到静默跳过的 bug）。这是提案中声称要修复的 bug 的同一类问题，如果新字段也用了 `omitempty`，空 `surfaces: {}` 会被序列化为 `surfaces:` 然后 Go 解析为 nil map。

问题：`forge task index` 的去重类型提取逻辑未定义
> "`forge task index` 从读 `interfaces` 改为读 `surfaces`（提取去重 surface 类型列表）" — 当前代码中 `capabilities` 变量（`build.go:66`）是 `[]string` 类型，直接传给 `extractBodyContext` 和后续的 `GetBreakdownTestTasks`。从 `surfaces` map 提取去重类型列表的逻辑虽然简单（遍历 map values，去重），但提案没有说明当 `surfaces` map 中存在未知类型值时的行为（是忽略、报错、还是传透？）。当前 `hasUIInterface` 函数对未知类型返回 false，这可能是正确的行为，但应该被显式声明。

### Downstream skill adaptation scope boundary

风险：延迟下游 skill 适配可能导致运行时字段名不一致
> "下游 skill 全面适配（gen-contracts, gen-test-scripts, eval-journey, eval-contract, run-tests）— 可在后续迭代中更新引用" — 当前的 gen-journeys SKILL.md 引用了 `surface` 字段（"Write the surface type to `.forge/config.yaml` in the `surface` field"），并且 surface rule 文件命名为 `surface-webui.md`（包含 `web-ui` 前缀）。提案将字段改为 `surfaces` map，命名改为 `web`。如果 gen-journeys 在 v3.0.0 发布时仍然引用旧字段名 `surface` 和旧命名 `webui`，它写入的 `surface` 字段和 `surfaces` map 会同时存在于 config 中，形成新的同步问题。提案承认 "out-of-scope skill 暂时保留对旧字段名的兼容读取"，但没有定义这个"暂时"的期限和兼容读取的具体策略。

### CLI command design

建议关注 `forge surfaces` 的退出码和输出格式契约。

问题：`forge surfaces <path>` 无匹配时的行为不够明确
> "无匹配则报错提示手动指定" — "报错"是输出到 stderr 并返回非零退出码，还是输出提示信息到 stdout 并返回零退出码？gen-journeys 作为 LLM 驱动的 skill 调用此命令时，需要区分"成功返回 surface 类型"和"未找到"两种情况。如果通过退出码区分，skill 的 bash 调用需要检查退出码；如果通过 stdout 内容区分，需要定义明确的输出格式。这应该作为 CLI 契约的一部分声明。

---

## Improvement Suggestions

建议：为检测信号表增加消歧优先级规则
Addresses: 信号冲突时缺少优先级规则
> What changes: 在 Feasibility Assessment 的信号表之后增加一个"Disambiguation Priority"小节。定义规则：当同一个 manifest 文件匹配多个信号时，前端框架 > 后端框架 > CLI 框架 > TUI 框架（理由：前端应用通常内含后端，但用户面向的是前端 surface）。同时声明：当检测到冲突信号时，init TUI 应展示所有候选并让用户选择，而不是自动决定。这样即使优先级规则不够完美，用户仍有最终决定权。

建议：定义路径规范化的精确规则和匹配算法
Addresses: 前缀匹配在路径规范化后的不一致行为、路径匹配的测试覆盖面没有声明
> What changes: 在 Config 结构小节增加"Path Normalization & Matching"段落，声明：(1) 所有路径按 `/` 分隔符分割为 segment 数组后按 segment 前缀匹配（不是字符前缀匹配，避免 `frontend` 匹配 `frontend-new`）；(2) 输入路径规范化规则：去除前导 `./`，去除尾随 `/`，不解析 `..`（如果包含 `..` 则报错），统一使用 `/` 分隔符；(3) 最长前缀匹配定义为：匹配的 segment 数最多者胜出。这些规则应在 PRD 中作为 CLI 契约的一部分声明，并在 Go 实现中有对应的表驱动测试。

建议：增加 `interfaces` 到 `surfaces` 的兼容读取过渡期
Addresses: 旧 `interfaces` 字段的静默忽略可能导致信息丢失、`interfaces` 到 `surfaces` 的语义映射不是一对一的
> What changes: 在 `ReadInterfaces` 函数（或其替代函数）中增加过渡逻辑：当 `surfaces` map 不为空时，直接从 `surfaces` 提取去重类型列表；当 `surfaces` 为空但 `interfaces` 不为空时，回退到读取 `interfaces` 字段（同时将旧值映射到新命名：`web-ui` -> `web`, `mobile-ui` -> `mobile`）。在控制台输出一条 deprecation 警告。这个过渡逻辑应在 v3.0.0 中引入，并在 v3.1.0 或 v4.0.0 中移除。这样既有真正的向后兼容，又不需要用户必须重新运行 init。

建议：声明 `surfaces` 字段的 YAML tag 为非 omitempty
Addresses: `Config` struct 变更缺少字段类型声明、`omitempty` 导致空 map 被丢弃
> What changes: 在提案中显式声明 Go struct 定义为 `Surfaces map[string]string \`yaml:"surfaces"\``（不带 `omitempty`）。空 map 应序列化为 `surfaces: {}` 而非被省略。这与 `interfaces` 的问题根源相同：空值被 omitempty 省略后，下游读取时无法区分"未配置"和"配置为空"。如果确实需要区分"未配置"和"配置为空"，可以用指针 `*map[string]string`，但这增加了复杂度，建议直接不使用 omitempty。

建议：将 gen-journeys 的 surface rule 文件重命名纳入 In Scope
Addresses: 延迟下游 skill 适配可能导致运行时字段名不一致
> What changes: 至少将 gen-journeys skill 的 `surface` 字段引用更新和 rule 文件重命名（`surface-webui.md` -> `surface-web.md`）纳入 In Scope。gen-journeys 是提案数据流的核心消费者（"gen-journeys 通过 `forge surfaces <path>` 查询 surface"），如果它的 surface 检测逻辑不随提案一起更新，就会出现旧检测逻辑写入 `surface: webui` 而新系统期望 `surfaces: {".": web}` 的不一致。下游的 gen-contracts、gen-test-scripts 等可以安全推迟，但 gen-journeys 必须同步更新。

建议：为检测逻辑定义目录遍历策略
Addresses: monorepo 根目录的 `package.json` 覆盖了子目录
> What changes: 在 Feasibility Assessment 中增加"Detection Traversal Strategy"段落，声明：(1) 遍历深度限制（建议 3 层，覆盖 `packages/frontend` 这类结构）；(2) 排除目录列表（`node_modules`, `.git`, `vendor`, `dist`, `build`）；(3) 遇到 workspace manifest（`pnpm-workspace.yaml` 或 `package.json` 中的 `workspaces` 字段）时，跳过根目录的依赖检测，只检测子目录。这个策略应该在 PRD 中声明并在实现中有对应的测试用例。

建议：声明 `forge surfaces` 的退出码和输出格式契约
Addresses: `forge surfaces <path>` 无匹配时的行为不够明确
> What changes: 在 CLI 命令小节增加退出码定义：(1) `forge surfaces`（全量查看）：成功时退出码 0，输出每行一个 `path=surface`；(2) `forge surfaces <path>`（路径查询）：匹配成功时退出码 0，输出单个 surface 类型字符串（无额外格式化）；未匹配时退出码 1，stderr 输出错误信息（包含匹配建议或手动配置提示）；(3) `forge surfaces --types`（类型列表）：成功时退出码 0，输出空格分隔的去重类型列表。这个契约是 gen-journeys skill 解析命令输出的基础，必须在 PRD 中精确定义。
