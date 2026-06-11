# Freeform Review: init-justfile 精简 -- CLI scaffold 替代 prompt 层模板 (REVISED)

**Reviewer**: Skill-Prompt-to-CLI Scaffold Architect
**Date**: 2026-06-09
**Document**: `docs/proposals/init-justfile-slim/proposal.md`
**Basis**: Revised proposal (post round-1/2/3 reviews), compared against current codebase

---

## Section 1: Background Assessment

本提案解决的核心问题：`init-justfile` skill 作为 Forge 最大的 skill（1645 行，8 个文件），将大量机械性 bash 代码模板维护在 LLM prompt 中。修订版在初版基础上做了显著改进：(1) 补充了 Consumer Impact 分析表格；(2) 明确了 `ci` 聚合排除 surface test 的设计理由；(3) 增加了多服务编排模式的生成策略；(4) 统一了 `# user-customized` 标记策略（lifecycle + quality 均标记，aggregate 不标记）；(5) 补充了 "替代方案与行业基准" 章节；(6) 增加了 "成功标准" 和 "范围外" 章节；(7) 保留了 Convention Cold Start Fallback 摘要；(8) 增加了回滚机制；(9) 将 Go 代码估算从 ~500 行修正为 ~600-1000 行。

技术路径不变：将 recipe bash 骨架生成从 prompt 层下沉到 `forge justfile scaffold` CLI 命令，agent 职责简化为"调用命令 + 检测语言 + 填占位符"。核心假设是：(1) 5 种 surface type 枚举稳定；(2) CLI 是 trusted producer；(3) 占位符清单完备；(4) agent 填值逻辑可在精简的 SKILL.md 中可靠描述。

修订版在很大程度上回应了前一轮评审的核心关切。以下聚焦于修订后仍残留的问题、修订引入的新问题，以及代码库验证后发现的偏差。

---

## Section 2: Key Risk Identification

### `<<URL_KEY>>` 占位符语义与 server-lifecycle.md 不一致

问题：提案的占位符表中 `<<URL_KEY>>` 的说明是"服务标识键名（用于 PID 文件命名）"，但当前 `server-lifecycle.md` 中 `<URL_KEY>` 的实际用途完全不同。

引用提案原文：
> `<<URL_KEY>>` | 服务标识键名（用于 PID 文件命名） | Surface key（与 `--key` 参数一致）

引用当前 `server-lifecycle.md`（第 266 行、第 308 行）：
> `_url=$(sed -n 's/^<URL_KEY>:[[:space:]]*\(.*\)/\1/p' "$_config" | head -1)`
> "`<URL_KEY>` is the YAML key for this surface's URL (e.g., `baseUrl` for frontend, `apiBaseUrl` for backend)."

`<URL_KEY>` 在当前代码中是 config.yaml 里 URL 字段的 YAML key 名（如 `baseUrl`、`apiBaseUrl`），用于从配置文件中提取服务 URL。提案将其描述为"服务标识键名（用于 PID 文件命名）"并将其解析来源定义为"Surface key（与 `--key` 参数一致）"，这与实际语义完全不同。PID 文件命名使用的是 `<surfaceKey>`，而非 `<URL_KEY>`。

后果：实施时如果按提案描述填值（将 surface key 填入 `<<URL_KEY>>`），probe recipe 中的 `sed` 解析会失败（config.yaml 中不存在名为 surface key 的 URL 字段），导致健康检查在运行时静默失败或报错。

### `<<SERVICE_LIST>>` 占位符的解析来源缺乏可操作性

风险：提案声明了 `<<SERVICE_LIST>>` 占位符，但解析来源为"Convention multi-service 定义"，这是一个不存在的 Convention 字段。

引用提案原文：
> `<<SERVICE_LIST>>` | 多服务编排时的服务启动依赖列表 | Convention multi-service 定义

当前 Convention 系统中不存在 `multi-service` 相关的字段定义。`forge surfaces` 命令的输出是 `key=type` 格式（或 scalar 形式下直接输出 type），不含依赖顺序信息。多服务编排的启动顺序（"先 api 后 web"）在提案中由一句话描述，但没有说明 agent 如何从 `forge surfaces` 的平铺输出中推导出依赖拓扑。

后果：agent 在实际填值时无法确定 `<<SERVICE_LIST>>` 的正确值。如果是 CLI 自动从 `forge surfaces` 推导，则不需要作为占位符暴露给 agent；如果需要 agent 填值，则必须定义依赖顺序的推断规则或 Convention schema。

### Consumer Impact 不完整 -- `forge quality-gate` Go 代码的 fallback 机制未纳入行动项

风险：提案的行动项第 4 条"更新 quality gate"和第 5 条"更新 `forge quality-gate` Go binary"描述了移除 fallback 链的意图，但代码库验证表明 Go 端的 fallback 机制比提案描述的更深层。

引用提案原文：
> "更新 quality gate：移除 fallback 链（`<key>-compile` 不存在时 fallback 到 `compile`），直接调 `<key>-compile`"
> "更新 `forge quality-gate` Go binary：将硬编码的 `just compile` / `just unit-test` 改为按 surface key 拼接 recipe 名"

引用当前 `forge-cli/pkg/just/just.go` 第 98-119 行（`ResolvePrefixedRecipe` 函数）：
> "When scope is non-empty, it probes for a prefixed recipe (e.g., 'backend-compile'). Falls back to the generic recipe name (e.g., 'compile') when the prefixed recipe does not exist or scope is empty."

`ResolvePrefixedRecipe` 是一个通用的 recipe 名解析函数，被 `RunGate` 调用。它同时处理"有 scope 时尝试前缀版本"和"无前缀版本时 fallback 到通用名"。提案要"移除 fallback 链"，但当前行为实际上正是提案的 recipe 命名模型所依赖的（单 surface scalar 项目用无前缀名，多 surface 项目用前缀名）。完全移除 fallback 意味着单 surface scalar 项目也无法工作（因为 `compile` 本身就是"通用名"）。

后果：行动项"移除 fallback 链"如果被字面执行，会破坏单 surface scalar 项目的 quality gate。需要区分"移除多 surface 项目中不必要的前缀 fallback"和"保留单 surface 项目的通用名解析"。

### `--aggregate` 模式读取 `forge surfaces` 但无法获取依赖拓扑

风险：提案声明 `--aggregate` 模式"读取 `forge surfaces` 获取全部 surface key"，但 `forge surfaces` 的输出格式是扁平的 `key=type` 列表（或 scalar 形式下仅输出 type），不包含 surface 之间的依赖关系。

引用提案原文：
> "CLI 读取 `forge surfaces` 获取全部 surface key"
> "多服务编排模式：...额外生成 `test-setup` 聚合 recipe，按依赖顺序编排各 surface 的启动（先 api 后 web）和 teardown（逆序）"

代码库验证：`forge surfaces` 的输出格式（`internal/cmd/surfaces.go` 第 100-119 行）是 `key=type` 行输出或 JSON `{key, type}` 数组。没有任何字段描述 surface 间的依赖关系。启动顺序（"先 api 后 web"）需要外部信息源，但提案未说明这个信息从何而来。

后果：CLI 的 `--aggregate` 模式无法独立确定多服务启动顺序。要么需要引入新的配置字段（如 Convention 中的 `depends_on`），要么硬编码 api-before-web 的隐式规则，要么将依赖编排推迟到 agent 层处理。当前提案对此沉默。

### Recipe 命名模型中 mobile 的 aggregate recipe 含 test-setup 与其他 service type 不一致

问题：提案的"每个 surface type 生成的 recipes"表格中，mobile 的 aggregate recipe 包含 `test-setup` 步骤，而 api 和 web 不包含。

引用提案原文：
> `api` | ... | `<key>` (dev->probe->test->teardown)
> `web` | ... | `<key>` (dev->probe->test->teardown)
> `mobile` | `test-setup`, `dev`, `probe`, `test`, `teardown`, `<key>` | ... | `<key>` (test-setup->dev->probe->test->teardown)

当前 `server-lifecycle.md` 中 test-setup 是 mobile surface 的特有步骤（模拟器启动），不适用于 api/web。但提案将 test-setup 列为 mobile 的 lifecycle recipe 而非 api/web 的，这一点本身是正确的。问题在于：多服务编排模式中也提到"额外生成 `test-setup` 聚合 recipe"，这里的 test-setup 含义是什么？是仅编排 mobile 的 test-setup，还是也包含其他 surface 的前置步骤？

后果：多服务编排的 `test-setup` 聚合 recipe 的生成条件和内容边界模糊。如果一个项目包含 api + mobile，`test-setup` 聚合应该只编排 mobile 的模拟器启动，还是也包含 api 的启动？

### 占位符语法 `<<...>>` 与 justfile 的 heredoc/字符串语义可能冲突

风险：提案选择 `<<...>>` 作为占位符语法，避免与 Go template `{{...}}` 和 justfile 变量 `{{var}}` 冲突。但 `<<` 在 bash 和 justfile 中有 heredoc 语义。

引用提案原文：
> "占位符语法：使用 `<<...>>` 而非 `{{...}}` 避免与 Go template `{{...}}` 和 justfile 变量语法 `{{var}}` 冲突。"

在 bash 中 `<<` 是 heredoc 操作符，在 justfile recipe body 中如果占位符出现在行首，`just` 解析器可能将其解释为 heredoc 开始标记。虽然占位符通常出现在赋值右侧（如 `_port=<<PORT>>`），但在字符串上下文中使用时（如 `echo "Waiting for <<PORT>>"`），双尖括号不会造成语法冲突。但如果 CLI 生成的模板中有裸占位符出现在行首（如某些 heredoc-like 上下文），可能导致 just 解析错误。

后果：低可能性但高影响的语法陷阱。建议在提案中明确约束：CLI 生成的模板中占位符不得出现在行首作为独立 token，或确保所有占位符都嵌入在赋值或字符串上下文中。

### 行动项 2 的 SKILL.md 目标行数 ~250 行与成功标准中的 <300 行不一致

问题：提案成功标准第 4 条声明 "prompt 层总行数 < 300 行"，但行动项第 2 条说 "从 548 行精简到 ~250 行"。

引用提案原文：
> "prompt 层（SKILL.md + 保留的 rules）总行数 < 300 行"（成功标准第 4 条）
> "重写 SKILL.md：从 548 行精简到 ~250 行"（行动项第 2 条）

保留的 `self-correction.md` 有 34 行。如果 SKILL.md 是 ~250 行，总计 ~284 行，刚好在 300 行以内。但成功标准是 "< 300 行"的硬约束，行动项的"~250 行"是估算值。如果实际精简结果为 270 行 + 34 行 = 304 行，将不满足成功标准。

后果：实施者可能以满足"~250 行"估算为目标，忽略 300 行的硬约束。建议成功标准改为"<= 280 行"留出安全边际，或将行动项估算改为"<= 250 行"作为硬目标。

### `forge quality-gate` 的 `NonBreakingGateSequence` 硬编码 recipe 名未被提案的行动项完全覆盖

风险：提案行动项第 5 条说"将硬编码的 `just compile` / `just unit-test` 改为按 surface key 拼接 recipe 名"，但代码库中的硬编码不仅是命令调用，而是 gate sequence 定义本身。

引用当前 `forge-cli/pkg/just/just.go` 第 48-53 行：
```go
func NonBreakingGateSequence() []GateRecipe {
    return []GateRecipe{
        {Name: "compile", Optional: false, Blocking: true},
        {Name: "fmt", Optional: true, Blocking: true},
        {Name: "lint", Optional: true, Blocking: true},
    }
}
```

这些 GateRecipe 的 Name 字段是固定的字符串 `"compile"`、`"fmt"`、`"lint"`，不包含任何 surface key 信息。运行时通过 `ResolvePrefixedRecipe` 动态解析为 `backend-compile` 或 `compile`。提案要"移除 fallback 链"，但移除 `ResolvePrefixedRecipe` 的 fallback 行为后，`NonBreakingGateSequence` 返回的固定名称 `"compile"` 对于多 surface 项目将无法匹配到任何 recipe（因为没有无前缀的 `compile`）。

后果：移除 fallback 链需要重新设计 `GateSequence` 的生成方式（从固定名称变为动态名称），这比行动项描述的"硬编码替换"复杂得多。当前行动项低估了改动范围。

---

## Section 3: Improvement Suggestions

建议：修正 `<<URL_KEY>>` 占位符的语义描述，使其与 `server-lifecycle.md` 中的实际用途一致。将描述改为"config.yaml 中该 surface URL 的 YAML key 名（如 `baseUrl`、`apiBaseUrl`）"，解析来源改为"Convention URL key 映射或 agent 从 config.yaml 结构推断"。如果决定在 CLI scaffold 中简化 probe 逻辑（不再从 config.yaml 读取 URL，而是直接使用 `<<HEALTH_URL>>`），则应删除 `<<URL_KEY>>` 占位符并用 `<<HEALTH_URL>>` 替代其功能。这一修正直接消除占位符语义不一致导致的运行时错误风险。

建议：将 `<<SERVICE_LIST>>` 的解析逻辑从 agent 侧移入 CLI 侧。CLI 可以直接从 `forge surfaces --json` 的输出获取所有 surface key 列表，并在 `--aggregate` 模式中自动构建服务列表，无需暴露为占位符。如果依赖顺序需要额外信息，应在提案中新增 Convention 字段定义（如 `dependencies` 数组），或声明"按 surface type 的字母序/固定规则排序"作为默认策略。这消除了 agent 无法解析该占位符的风险。

建议：细化行动项第 4 条和第 5 条，明确 fallback 机制的保留策略。建议措辞改为："保留 `ResolvePrefixedRecipe` 的 fallback 行为（前缀 -> 无前缀），但移除 `run-tests` skill 中额外的 fallback 链逻辑。Go 端的 recipe 名解析逻辑不变，因为当前实现已正确支持提案的命名模型。"这避免误删对单 surface scalar 项目至关重要的 fallback 路径。

建议：在多服务编排模式的描述中明确依赖顺序的信息来源。三个选项：(a) 声明"按 api -> web -> mobile 固定顺序"，适用于 Forge 典型场景；(b) 新增 Convention `dependencies` 字段；(c) 将依赖编排完全交给 agent（CLI 只生成独立的单 surface recipe，跨 surface 编排在 SKILL.md 流程中由 agent 按需组装）。推荐选项 (a)，因为 surface type 枚举固定且依赖关系在大多数场景下可预测。

建议：在提案中增加对占位符在模板中位置的约束说明。声明"所有占位符必须嵌入在赋值语句右侧（如 `_port=<<PORT>>`）或字符串字面量中，不得出现在行首作为独立 token"，以避免 `<<` 的 heredoc 歧义。这是一个实现约束而非设计约束，但在提案中声明可以防止实施时遗漏。

建议：统一成功标准的行数约束。将成功标准第 4 条改为"prompt 层（SKILL.md + 保留的 rules）总行数 <= 280 行"，与行动项第 2 条的 ~250 行 SKILL.md + 34 行 self-correction.md = ~284 行的估算一致，留出 4 行安全边际。

建议：在行动项第 5 条中增加具体的 Go 代码改动范围描述。当前措辞"将硬编码的 `just compile` / `just unit-test` 改为按 surface key 拼接 recipe 名"可能被理解为简单的字符串替换。建议改为"审查 `pkg/just/just.go` 中的 `GateSequence` 函数和 `ResolvePrefixedRecipe` 函数，确认其对 `<key>-<recipe>` / `<recipe>` 双路径的解析行为与提案的命名模型一致。当前 `ResolvePrefixedRecipe` 的 fallback 行为（前缀 -> 通用名）已覆盖提案的命名模型，无需修改。仅需确认 `run-tests` skill 不再依赖额外的 fallback 链。"
