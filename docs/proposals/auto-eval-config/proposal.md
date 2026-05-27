---
created: 2026-05-26
author: "faner"
status: Approved
---

# Proposal: Auto-Eval Configuration with Generic Config Key Resolution

## Problem

4 个文档评估 skill（eval-proposal、eval-prd、eval-ui、eval-design）在对应文档生成后，通过 `AskUserQuestion` 手动询问用户是否运行评估。每次交互都需要用户手动选择，增加了流水线的交互成本。项目已有 `auto.runTasks` 等 ModeToggle 配置驱动自动化的成功模式，eval 阶段应采用相同机制——让用户通过配置控制 eval 自动运行行为，而非统一消除所有手动确认（prd 和 techDesign 默认保持手动确认）。

同时，当前 `forge config get/set` 的 key 路由系统（`getAutoKeyValue`、`getWorktreeKeyValue`、`setAutoConfigValue` 等）基于硬编码的分发器和 `autoModeField` switch/case，只支持两层深度（`auto.{field}.{subfield}`）。任何新的嵌套配置都需要修改路由核心代码。本提案将路由系统泛化为任意深度 key 遍历，从根本上消除这个扩展瓶颈。

### Evidence

- brainstorm、write-prd、tech-design 三个 skill 均在文档提交后通过 `AskUserQuestion` 询问是否运行 eval
- ui-design 是唯一例外——无条件自动运行 eval-ui，与其他 skill 行为不一致
- `auto.runTasks` 已证明 ModeToggle 配置模式可以成功消除冗余交互
- 当前 `getAutoKeyValue` 使用 `strings.SplitN(rest, ".", 2)` 只拆一层，`autoModeField` 通过 switch/case 映射，`parseAutoRaw` 硬编码 `modeFields` 列表——三者共同限制了 key 深度和新增配置的灵活性

### Urgency

随着 quick 流水线不断优化自动化流程（auto.runTasks、auto.consolidateSpecs 等），eval 的手动确认成为剩余的主要交互摩擦点。每个 `AskUserQuestion` 需要 AI 生成提问 + 用户阅读选择 + AI 解析回答，约 15-20s 延迟（基于对 `AskUserQuestion` 函数的代码路径分析：LLM 生成提问 ~5-8s + 用户阅读选择 ~5-10s + LLM 解析回答 ~3-5s）。4 个 eval 点累计增加 ~60-80s 交互时间/每条流水线运行。统一配置模式可直接消除这部分成本。

泛化 config key resolution 的紧迫性来自 eval 配置的三层嵌套需求（`auto.eval.proposal.quick`）：当前 `getAutoKeyValue` 使用 `SplitN(rest, ".", 2)` 只拆一层，`autoModeField` 只支持二层分发。三层嵌套意味着硬编码路由需要新引入 eval 子分发器（~30 行/字段），且 `parseAutoRaw` 的硬编码 `modeFields` 列表也需要扩展。泛化路由在 eval 配置场景下直接触发——不是"未来可能需要"，而是"当前需求已到达现有路由的天花板"。

## Proposed Solution

### Part 1: Generic Config Key Resolution

重写 `GetConfigValue` / `SetConfigValue` 为基于任意深度 key 的通用遍历器：

**Get 路径**（反射）：
- 将 `key` 按 `.` 拆分为路径段，沿 Go Config struct 树递归走 reflect.Value
- 字段匹配规则：优先按 YAML tag 匹配（如 `source-branch` 匹配 `yaml:"source-branch"`），YAML tag 不存在时按 Go field name 匹配；遇到 `yaml:",inline"` 标记的字段时，将该字段的子项展开到父级（如 `CoverageConfig.ByType` 的 inline map 条目直接作为 `coverage` 的子 key 访问）
- 指针字段处理：遇到 `nil` 指针字段时返回 `errKeyNotFound`（get 不自动初始化）
- 叶子节点的类型决定输出格式：bool → "true"/"false"，string → 原值，ModeToggle → "quick:true full:true"，[]string → 换行连接，map → 按 key 查找后递归格式化 value
- 非叶子节点（struct/嵌套 struct）的 get 行为：遍历所有导出字段，按字段类型格式化后逐行输出。格式规则：ModeToggle → `<fieldName>: quick:<bool> full:<bool>`；bool → `<fieldName>: true/false`；嵌套 struct → 递归展开（缩进 2 空格）；map → 遍历 key 后逐行输出。例如 `forge config get auto` 输出：
  ```
  runTasks: quick:true full:true
  consolidateSpecs: quick:true full:true
  gitPush: false
  eval:
    proposal: quick:true full:true
    prd: quick:false full:false
    uiDesign: quick:true full:true
    techDesign: quick:false full:false
  ```
- 默认值通过 `ReadConfig` 的 `applyDefaults` 已生效，反射读取的就是含默认值的最终状态

**Set 路径**（Go struct marshal，与当前 writeConfig 一致）：
- 将 `key` 按 `.` 拆分为路径段，通过反射沿 Go Config struct 树定位目标字段并设值
- 使用现有的 `readOrCreateConfig` + 反射 set + `writeConfig`（yaml.Marshal）写回文件
- 选择 struct marshal 而非直接操作 YAML Node 的理由：当前所有 set 路径均通过 Go struct → yaml.Marshal 写回，保持一致的序列化路径避免格式保真度差异
- 写回后下次 `ReadConfig` 会重新应用默认值
- 错误行为：key 路径中的不存在字段返回 `errKeyNotFound`；尝试 set 非 leaf 节点（如 `auto.eval`）返回错误 `"cannot set non-leaf key, use <key>.<field>"`；ModeToggle 字段（如 `auto.eval.proposal`）也视为非 leaf——set 时返回错误 `"cannot set ModeToggle directly, use <key>.quick or <key>.full"`；超过 struct 深度的 key 路径（如 `auto.eval.proposal.quick.extra`）返回 `errKeyNotFound`；value 类型不匹配（如对 bool 字段设 `"maybe"`）返回 `"invalid value \"maybe\" for bool field <path>: expected true or false"`

**消除的硬编码**：
- `autoModeField` switch/case → 删除（反射替代）
- `getAutoKeyValue` / `getWorktreeKeyValue` / `getCoverageKeyValue` → 合并为一个 `getByPath`
- `setAutoConfigValue` / `setWorktreeConfigValue` / `setCoverageConfigValue` → 合并为一个 `setByPath`
- `parseAutoRaw` 硬编码 `modeFields` 列表 → 递归扫描 auto 子树

### Part 2: Auto-Eval Configuration

在 `.forge/config.yaml` 的 `auto` 块中新增 `eval` 嵌套结构体，包含 4 个独立的 ModeToggle 字段：

- `auto.eval.proposal` — 控制 eval-proposal 是否自动运行
- `auto.eval.prd` — 控制 eval-prd 是否自动运行
- `auto.eval.uiDesign` — 控制 eval-ui 是否自动运行
- `auto.eval.techDesign` — 控制 eval-design 是否自动运行

每个 ModeToggle 支持 `quick`/`full` 子键，通过 CLI 级别的 mode 检测判断当前管道模式。

默认值：
- `proposal`: `quick: true, full: true` — 默认自动运行（proposal 是管道入口，尽早发现问题）
- `prd`: `quick: false, full: false` — 默认询问用户
- `uiDesign`: `quick: true, full: true` — 保持当前无条件自动行为，零回退风险
- `techDesign`: `quick: false, full: false` — 默认询问用户

### Innovation Highlights

1. **泛化 key resolution** 消除了 config 系统的扩展瓶颈——未来新增任何 `auto.*` 配置只需加 Go struct 字段 + 默认值，零路由代码改动
2. **eval 配置复用已有 ModeToggle 模式**，无新概念——这是本方案的保守之处，eval 部分不追求创新，而是追求与现有模式的一致性
3. **Go struct IS the schema**——结构体定义即配置 schema，反射路由将 struct 定义自动转换为 CLI get/set API，无需路由注册代码。这是 meta-programming 模式的核心：struct field 是声明式定义，反射是运行时解释器，两者组合消灭了整类路由 boilerplate
4. **mode 检测作为 CLI 级 API**（`forge config get mode`），不依赖 manifest 文件格式
5. **反射路由是平台思维而非功能思维**：每新增一个配置字段，硬编码路由需要修改 3 处（`autoModeField` switch、get 分发器、set 分发器），而反射路由的边际成本为零。需要承认：反射路由是标准 Go 技巧（Viper `Unmarshal`、koanf 都使用类似机制），本方案的创新不在于"用反射"本身，而在于将反射与 CLI get/set 路由结合——struct 定义即 CLI API，这种 meta-programming 模式在 Go 配置库中不常见（Viper 的 `Get()` 用 flat map，`Unmarshal()` 是单向的配置→struct，不支持反向的 struct→CLI get/set 路由）

## Requirements Analysis

### Key Scenarios

- **自动评估 proposal（默认）**: 用户执行 `/quick`，brainstorm 提交 proposal 后自动运行 eval-proposal，无需手动确认
- **手动确认 prd eval**: 用户执行 `/write-prd`，提交 PRD 后询问是否运行 eval-prd（默认行为）
- **配置驱动的灵活控制**: 用户通过 `forge config set auto.eval.prd.full true` 开启 PRD 自动评估
- **config 文件异常 YAML**: 用户手动编辑 config.yaml 写入 `auto.eval.proposal: "invalid"`（string 而非 mapping），`forge config get auto.eval.proposal` 遇到类型不匹配时返回错误 `"type mismatch at auto.eval.proposal: expected struct, got string"`
- **ui-design 行为统一**: ui-design 从无条件自动运行改为读取配置，默认值保持 ON 以避免回退
- **任意深度 config get**: `forge config get auto.eval.proposal.quick` → `true`（四层深度自动支持）
- **任意深度 config set**: `forge config set auto.eval.proposal.quick false` → 正确写入嵌套结构
- **mode 检测**: `forge config get mode` → `quick` 或 `full`，skill 据此判断当前管道模式
- **mode 检测（无 feature 上下文）**: 在 feature 目录外调用 `forge config get mode` 返回 `"none"`，skill 应跳过 eval 自动运行逻辑，回退到 AskUserQuestion
- **无效 key 深度**: `forge config get auto.eval.proposal.quick.extra`（在 bool leaf 之后继续访问）返回 `errKeyNotFound`
- **无效 key 路径**: `forge config get auto.nonexistent` 返回 `errKeyNotFound`

### Non-Functional Requirements

- 向后兼容：未配置时使用默认值，不影响现有行为（ui-design 默认保持 ON）
- 错误消息质量：所有 CLI 错误消息包含具体信息——无效 key 包含 key 路径（`"key \"auto.nonexistent\" not found"`），类型不匹配包含期望类型（`"invalid value \"maybe\" for bool field auto.eval.proposal.quick: expected true or false"`），非 leaf set 包含建议操作（`"cannot set non-leaf key auto.eval, use auto.eval.<field>.<subfield>"`）
- 配置热生效：修改配置后无需重启即可生效
- 性能：泛化路由的反射开销应有明确上限——`GetConfigValue` 端到端延迟 <1ms（当前硬编码分发器 <0.1ms），反射遍历开销增加一个数量级但仍在 CLI 可接受范围内。反射调用次数等于 key 深度（`auto.eval.proposal.quick` = 4 次 reflect 操作），每次 <0.1ms。测量条件：CLI 正常运行态（非首次冷启动，Go reflect 的首次类型查找较慢，但 CLI 每次调用是新进程，首次即唯一次。实测 CLI 进程启动本身约 20-50ms，反射遍历 <1ms 是相对于进程启动时间的可接受开销）

### Constraints & Dependencies

- Go reflect 包用于 struct 遍历（stdlib，无外部依赖）
- yaml.v3 用于 YAML 序列化/反序列化（已有依赖），需处理 `yaml:",inline"` tag 语义（inline map 字段的 key 直接出现在父节点下）
- mode 检测需要新增 CLI 级 API：`forge config get mode` 返回 `quick`、`full` 或 `none`。当前 `feature_complete.go` 通过检查 proposal.md 是否存在于 feature 目录来推断 quick mode，该逻辑需提取为独立的 mode 检测函数。CLI 通过解析当前工作目录推断 feature slug：从 `pwd` 中提取匹配 `.forge/features/<slug>` 路径模式的最后一段作为 feature slug，然后检查对应 feature 目录下是否存在 `proposal.md`。若 pwd 不匹配 `.forge/features/` 模式，返回 `"none"`
- 4 个 skill 文件为 markdown，通过 Bash 调用 `forge config get` 实现 config check。典型的 skill 端 config check 模式如下：
  ```bash
  # Eval 自动运行检查（替换原 AskUserQuestion 步骤）
  MODE=$(forge config get mode 2>/dev/null)
  if [ $? -ne 0 ]; then
    # CLI 不可用，回退到手动确认
    AskUserQuestion "是否运行 eval-{type}?"
  else
    EVAL_ENABLED=$(forge config get auto.eval.{skillKey}.$MODE 2>/dev/null)
    if [ "$EVAL_ENABLED" = "true" ]; then
      # 自动运行 eval
      /{eval-skill}
    elif [ "$EVAL_ENABLED" = "false" ]; then
      echo "eval-{type} 已通过配置跳过"
    else
      # 配置不可读，回退到手动确认
      AskUserQuestion "是否运行 eval-{type}?"
    fi
  fi
  ```
  每次调用 spawn 一个 CLI 子进程，约 50-100ms 延迟。选择 CLI 子进程而非直接读取 config.yaml（via yq）的理由：CLI 封装了默认值应用逻辑（`applyDefaults`），直接读 YAML 文件无法获取默认值；且 CLI 是唯一经过测试的配置读取路径。quick 流水线 4 个 eval 点共 200-400ms subprocess 开销，相对于整个流水线运行时间（分钟级）可忽略

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 无开发成本 | 保留交互摩擦 + 路由系统扩展瓶颈 | Rejected: 双重问题 |
| 扁平命名空间 (auto.evalProposal) | — | 零路由改动；完全解决 eval 用户需求（配置化控制 eval 自动运行） | 不解决路由瓶颈，命名不直观 | Rejected: 仅解决用户面的 eval 问题，不解决开发者面的路由扩展问题——但需承认：如果仅关注 eval 需求，此方案性价比最高 |
| **泛化路由 + 嵌套 eval 配置** | auto.runTasks 模式 + 通用 key 遍历 | 根本性解决扩展瓶颈 + 灵活的 eval 控制；反射路由零边际成本 | 路由层重写的一次性投入（~2.5h）；影响所有现有 `forge config get/set` 调用路径，任何反射遍历 bug 会同时影响所有配置操作；反射代码调试难度高于显式分发器（stack trace 通过 reflect.Value 间接调用，不如函数调用直观） | **Selected: 最优长期方案** |
| 仅扩展路由支持三层 | — | 改动较小（约 50-80 行：`autoModeField` 新增 eval case + `getEvalKeyValue` + `setEvalConfigValue`，首字段约 50 行含 eval 子分发器 + 后续字段约 15 行/字段） | autoModeField/get/set 仍需修改，未来四层时又要改；200 行反射重写的盈亏平衡点约为 5 个新字段（50 行首字段 + 15 行后续 vs 200 行一次性投入），当前已规划 eval（4 字段）+ 未来知识管理和验证自动化预计 3-4 字段 | Rejected: N 个字段的边际成本为 O(N)×15 行，反射路由边际成本为 0 行 |

### Industry References

- **Viper** (Go): 支持 `viper.Get("a.b.c")` 任意深度 key 访问。Viper 内部使用 flat key map（`map[string]interface{}`）将嵌套结构展平为 dot-separated key，与反射遍历 Go struct 的方案不同。Viper 的 `Unmarshal()` 方法确实使用反射将配置映射到 struct，但这是单向的（配置 → struct），不支持反向的 struct → CLI get/set 路由。Viper 不处理 nil 指针（配置始终初始化为零值）、不支持自定义类型格式化（如 ModeToggle 的 `quick:true full:true` 格式）、不追踪原始写入 vs 默认值。本方案的反射方法保留 Go struct 的类型信息，能实现 ModeToggle 格式化和 `parseAutoRaw` 的 raw tracking——这些是 Viper flat map 无法做到的
- **koanf** (Go): Go-native 配置库，使用 `konf` 的 struct tag 驱动配置映射。koanf 与本方案最接近——都通过 struct 定义作为配置 schema。但 koanf 侧重多源合并（file/env/flag），不提供 CLI get/set 路由能力。本方案的差异化在于：struct 定义不仅是 schema，还是 CLI 的操作接口——每个字段自动成为 `forge config get/set` 的合法 key，无需额外的路由注册
- **最小可行替代方案**: 在 `AutoConfig` 中添加 4 个扁平 ModeToggle 字段（`auto.evalProposal`/`auto.evalPrd`/`auto.evalUiDesign`/`auto.evalTechDesign`），扩展 `autoModeField` switch 加 4 个 case。优点：20 分钟内完成，零路由架构改动，零回归风险（现有路径完全不变）。缺点：命名不直观（`auto.evalProposal` vs `auto.eval.proposal`）；不解决路由瓶颈，四层嵌套场景仍需硬编码扩展；`parseAutoRaw` 的 `modeFields` 列表也需要扩展。适用于：仅需解决 eval 问题且不关心路由架构的场景

## Feasibility Assessment

### Technical Feasibility

完全可行。Go reflect（stdlib）和 yaml.v3（已有依赖）均为成熟技术。

**Get 路径实现要点**：
```go
func getStructValueByPath(v reflect.Value, segments []string) (string, error)
```
按 segments 递归走 reflect.Value：struct → 按 YAML tag 匹配字段（tag 优先于 Go field name，遇到 `,inline` 标记的 map 字段时跳过字段名层级，直接在 map 中查找）；指针 → 解引用，nil 时返回 errKeyNotFound；map → 按 key 查找后递归；自定义 YAML 类型（实现 `yaml.Unmarshaler`） → 回退到硬编码路径（绕行 SurfacesMap）；到达叶子节点时按类型格式化。

**Set 路径实现要点**：
```go
func setStructValueByPath(v reflect.Value, segments []string, value string) error
```
按 segments 递归走 reflect.Value：struct → 按 YAML tag 匹配字段（含 `,inline` tag 处理）；指针 → nil 时自动初始化（`reflect.New` + `Set`）；map → 按 key 查找或创建；ModeToggle → 拒绝直接 set，要求指定 `.quick` 或 `.full`；到达叶子节点时解析 value 字符串并设值。通过 `writeConfig`（yaml.Marshal）写回文件。

### Resource & Timeline

预计 4-5 小时：泛化路由重写含 `,inline` tag 和 SurfacesMap 处理（2.5h）+ eval 配置 + 4 个 skill（1h）+ 测试含 parseAutoRaw 回归（1-1.5h）。额外 1 小时用于处理边缘情况：`yaml:",inline"` tag 的反射兼容（CoverageConfig.ByType）、`SurfacesMap` 自定义类型的反射路由绕行（config get 遇到 SurfacesMap 时走现有的 `getCoverageKeyValue` 路径，不尝试反射遍历自定义 marshaler）、`WorktreeConfig.CopyFiles`（[]string 类型）的格式化。

### Dependency Readiness

所有依赖均已就绪，除 mode 检测 API 需要新增实现（提取 `feature_complete.go` 中的 quick mode 判定逻辑）。`yaml:",inline"` tag 处理需要在 `getStructValueByPath` 中增加 tag 检测逻辑，不引入新依赖但增加实现复杂度。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| config get/set 需要硬编码路由 | Occam's Razor | Refuted: 反射 + YAML tag 遍历可实现通用路由，消除所有分发器。Go struct IS the schema——字段定义即 CLI API |
| 嵌套结构需要修改路由核心 | First Principles | Refuted: 泛化后嵌套和扁平实现成本相同 |
| 每次评估都需要用户确认 | Occam's Razor | Challenged: 配置化控制比无条件自动化更安全——让用户按需选择，而非强制全自动化 |
| ui-design 的无条件自动评估是错误的 | Assumption Flip | Refined: 保持默认 ON，由配置控制而非硬编码 |
| proposal 默认应询问用户 | 5 Whys | Challenged: proposal 是流水线入口，默认自动评估有助于尽早发现问题（设计选择，非逻辑推论） |

## Scope

### In Scope

- **Generic config key resolution**: 重写 `GetConfigValue`/`SetConfigValue` 为任意深度 key 遍历（反射 get + 反射 set + struct marshal 写回），含 `yaml:",inline"` tag 检测和自定义 YAML 类型的 fallback 处理
- **消除硬编码路由**: 删除 `autoModeField`、`getAutoKeyValue`/`getWorktreeKeyValue`/`getCoverageKeyValue`、`setAutoConfigValue`/`setWorktreeConfigValue`/`setCoverageConfigValue`（保留函数体作为 SurfacesMap 等自定义类型的 fallback）
- **泛化 `parseAutoRaw`**: 递归扫描 auto 子树，不再硬编码 modeFields。泛化后的 raw 数据结构使用扁平路径 key：`map[string]map[string]bool`，其中外层 key 为字段路径（如 `"eval.proposal"`），内层 map 记录该字段下哪些子键（`"quick"`/`"full"`）在 YAML 中显式出现。`applyDefaults` 相应改为 `applyModeDefault(&a.Eval.Proposal, a.raw, "eval.proposal", d.Eval.Proposal)` 的 flat-path 调用方式
- **Go CLI AutoConfig 新增 Eval 嵌套结构体**: proposal/prd/uiDesign/techDesign 各为 ModeToggle
- **AutoConfigDefaults 更新**: proposal {true,true}, prd {false,false}, uiDesign {true,true}, techDesign {false,false}
- **forge config get/set auto.eval.* 支持**: 任意深度自动支持
- **4 个 skill 增加 config check**: 用 bash config check 模式替换各 skill 中的 AskUserQuestion 步骤（具体修改点：brainstorm 步骤 6-7 之间、write-prd 步骤 6-7 之间、ui-design 步骤 6-7 之间、tech-design 步骤 5-6 之间）
- **mode 检测 API**: 新增 `forge config get mode`，提取 `feature_complete.go` 中的 quick mode 判定逻辑为独立函数，供 skill 通过 CLI 调用。返回值：feature 目录内 + proposal.md 存在 → `"quick"`，feature 目录内无 proposal.md → `"full"`，非 feature 目录 → `"none"`。CLI 通过解析当前工作目录路径（匹配 `.forge/features/<slug>` 模式）推断 feature slug 上下文
- **单元测试更新**: config_test.go、config_schema_test.go，包含 `TestGetByPath_InlineMap`（coverage inline tag）、`TestParseAutoRaw_EvalConfig`（eval raw tracking）、`TestParseAutoRaw_ExistingFields_Regression`（现有字段回归）

### Out of Scope

- eval skill 本身（eval/SKILL.md）的行为修改
- 新的 eval 类型添加
- config 文件异常 YAML 结构的容错恢复（检测到类型不匹配时返回错误，不尝试自动修复）
- quick/full 管道流程的其他变更
- forge guide 文档更新（后续迭代）
- `forge config init` TUI 的更新（可后续迭代）
- 并发配置修改保护（两个 `forge config set` 同时运行的 last-write-wins 问题——当前系统已有此行为，非新增风险）
- 性能基准测试（反射 vs 硬编码的微基准对比）
- config 文件校验（无效 key pattern 检测）
- PR-1 rollback 方案的自动化测试（rollback 机制为 git revert，不引入代码级 rollback 路径）
- 外部工具输出格式兼容性迁移指南
- 反射路由 vs 显式分发器的 debuggability 对比文档（反射代码 stack trace 可读性较差，但 CLI 工具单次调用、无长期运行状态，debug 影响有限）

## Delivery Strategy

**推荐拆分为 2 个独立 PR 交付**：

1. **PR-1: Generic Config Key Resolution**（Part 1）— 仅重写 GetConfigValue/SetConfigValue，不新增 eval 配置。交付后运行完整回归测试（config_test.go、config_schema_test.go），确保 `worktree.source-branch`、`coverage.coding.feature` 等现有路径行为不变。PR-1 验收标准：现有所有 config get/set 路径行为不变 + `auto.eval.*` 路径可正确解析（先合并 EvalConfig 结构体定义，不含 skill 集成）。PR-1 回滚方案：revert commit，恢复旧分发器代码。
2. **PR-2: Auto-Eval Configuration**（Part 2）— 依赖 PR-1 的泛化路由，新增 4 个 skill 的 config check（用 bash 模板替换 AskUserQuestion 步骤）。PR-2 回滚方案：revert commit，恢复 AskUserQuestion 步骤。

拆分的好处：PR-1 的回归风险（影响所有现有 config get/set 路径）与 eval 功能解耦，可独立回滚。PR-2 是纯增量变更，失败时只需移除 eval 相关代码。PR-1 合并前需通过 acceptance test：泛化路由支持任意嵌套 struct 的 get/set，包含 inline tag 和自定义类型的 fallback 处理。

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Part 1 路由重写的回归风险影响所有现有 config 路径 | M | H | 拆分为独立 PR（见 Delivery Strategy）；PR-1 交付后全量回归测试通过再启动 PR-2；保留 `getCoverageKeyValue`/`getWorktreeKeyValue` 等函数体作为 fallback，在新路由返回 `errUnsupportedType` 时回退到旧路径 |
| 反射遍历对 `yaml:",inline"` 标记字段（`CoverageConfig.ByType`）的解析复杂度 — inline map 的 key 在 YAML 中直接出现在父节点下，反射路由需要识别 `,inline` tag 并将 map 条目展开为父级子 key | M | H | 在 `getStructValueByPath` 中检测 `yaml:",inline"` tag：遇到 inline 标记的 map 字段时，跳过 map 字段名层级，直接在 map 中按 segment 查找 key；新增测试用例 `TestGetByPath_InlineMap` 验证 `coverage.coding.feature` 正确解析 |
| 反射遍历对 SurfacesMap 自定义 YAML 类型（`UnmarshalYAML`/`MarshalYAML`）的兼容性 — SurfacesMap 在 YAML 中可序列化为 string 或 mapping，反射看到的是 `map[string]string` 但实际 YAML 表示可能不一致 | L | M | 为自定义 YAML 类型注册反射路由绕行：config get 遇到实现 `yaml.Unmarshaler` 接口的字段时，回退到现有的硬编码路径（如 `getCoverageKeyValue` 中对 SurfacesMap 的处理），而非尝试反射遍历 |
| 反射遍历对 map 类型（CoverageConfig.ByType）的序列化复杂度 — CoverageConfig.ByType 的 value 是 struct（Type+Percentage），非简单标量 | M | H | 为 map[string]struct 类型注册类型感知格式化器：Percentage → 输出数值，maintain → 输出 "maintain"；getByPath 遇到 map 时按 key 查找后递归格式化 value struct |
| YAML 注释和格式在 set 操作后的保真度 | L | L | Set 路径完全通过 Go struct → yaml.Marshal 写回，与当前 writeConfig 行为一致。YAML 注释丢失是 yaml.Marshal 的已知行为，但当前所有 set 路径已有此特性，不是新增风险。用户首次使用新路由系统 set 配置时可能注意到注释消失，但行为与旧系统一致 |
| parseAutoRaw 泛化后 raw tracking 精度变化 | M | M | 递归扫描保持相同的叶子节点追踪粒度；新增测试用例 `TestParseAutoRaw_EvalConfig` 验证三项：(a) raw map 包含正确的 flat-path key（如 `map["eval.proposal"]["quick"]=true`），(b) `applyDefaults` 仅补充 YAML 中未显式出现的子键，(c) `applyDefaults` 不覆盖用户显式设置的值；新增测试用例 `TestParseAutoRaw_ExistingFields_Regression` 验证现有 auto 字段（test、consolidateSpecs、gitPush）在泛化后仍产生正确的 raw 数据 |
| 4 个 skill 的 config check 逻辑漂移 — markdown skill 文件无编译时/运行时强制一致性检查 | M | M | 使用统一的 bash config check 模板（见 Constraints & Dependencies 中的示例代码），在 skill 中用 EXTREMELY-IMPORTANT 标注；将模板固化为 skill 的共享 snippet（skill markdown 可 include 的公共代码块）以减少复制漂移风险 |
| skill → CLI 子进程依赖：skill 调用 `forge config get` 时 CLI binary 未构建或版本过旧 | M | M | skill 的 config check 模板包含 fallback 逻辑：CLI 退出码非零时回退到 AskUserQuestion；CI 中在 skill 测试前强制 `go install` |

## Success Criteria

**PR-1: Generic Config Key Resolution**

- [ ] `forge config get auto.eval.proposal` 返回 `quick:true full:true`（三层深度）
- [ ] `forge config get auto.eval.proposal.quick` 返回 `true`（四层深度）
- [ ] `forge config get auto.eval` 返回所有 eval 子字段的汇总输出——每行一个字段，格式为 `<fieldName>: quick:<bool> full:<bool>`（如 `proposal: quick:true full:true`），字段按 Go struct 定义顺序排列（SC 验证此中间节点格式）
- [ ] `forge config get auto` 返回所有 auto 子字段的汇总输出，包含混合类型：ModeToggle 字段（runTasks、consolidateSpecs）格式为 `<name>: quick:<bool> full:<bool>`，bool 字段（gitPush）格式为 `<name>: <bool>`，嵌套 struct（eval）缩进 2 空格后递归展开
- [ ] `forge config set auto.eval` 被拒绝（返回错误："cannot set non-leaf key, use auto.eval.<field>.<subfield>"）
- [ ] `forge config set auto.eval.proposal true` 被拒绝（返回错误："cannot set ModeToggle directly, use auto.eval.proposal.quick or auto.eval.proposal.full"）
- [ ] `forge config set auto.eval.prd.full true` 正确写入嵌套 config
- [ ] `forge config get auto.eval.proposal.quick.extra`（超过 bool leaf 继续访问）返回 `errKeyNotFound`
- [ ] `forge config get auto.nonexistent` 返回 `errKeyNotFound`
- [ ] `forge config get coverage.coding.feature` 对 `yaml:",inline"` map 字段正确解析并保持现有行为不变（回归测试 + inline tag 兼容验证）
- [ ] `forge config get worktree.source-branch` 保持现有行为不变（回归测试）
- [ ] `parseAutoRaw` 对 `auto.eval.*` 字段生成正确的 flat-path raw map（如 `map["eval.proposal"]["quick"]=true`），`applyDefaults` 仅补充 YAML 中未显式出现的子键
- [ ] `parseAutoRaw` 对现有 auto 字段（test、consolidateSpecs、gitPush）在泛化后仍产生正确的 raw 数据（回归测试）
- [ ] 现有配置测试（config_test.go、config_schema_test.go）全部通过

**PR-2: Auto-Eval Configuration**

- [ ] `forge config get mode` 在 feature 目录内 + proposal.md 存在时返回 `"quick"`
- [ ] `forge config get mode` 在 feature 目录内无 proposal.md 时返回 `"full"`
- [ ] `forge config get mode` 在非 feature 目录时返回 `"none"`
- [ ] brainstorm 在 `auto.eval.proposal` 对应模式为 true 时跳过 AskUserQuestion，直接运行 eval-proposal
- [ ] brainstorm 在 `auto.eval.proposal` 对应模式为 false 时保持现有 AskUserQuestion 行为
- [ ] write-prd 在 `auto.eval.prd` 对应模式为 true 时跳过 AskUserQuestion，直接运行 eval-prd
- [ ] write-prd 在 `auto.eval.prd` 对应模式为 false 时保持现有 AskUserQuestion 行为
- [ ] ui-design 在 `auto.eval.uiDesign` 对应模式为 true 时跳过 AskUserQuestion，直接运行 eval-ui
- [ ] ui-design 在 `auto.eval.uiDesign` 对应模式为 false 时保持现有 AskUserQuestion 行为
- [ ] tech-design 在 `auto.eval.techDesign` 对应模式为 true 时跳过 AskUserQuestion，直接运行 eval-design
- [ ] tech-design 在 `auto.eval.techDesign` 对应模式为 false 时保持现有 AskUserQuestion 行为
- [ ] 未配置时（config missing）：proposal 和 uiDesign 默认自动运行，prd 和 techDesign 默认询问
- [ ] Skill 行为 SC 代理验证：手动验证至少 1 个 skill（brainstorm）的 config check 集成后，通过 bash 模板代码审查确认其余 3 个 skill 使用相同模板（因为 skill 行为由 AI agent 执行 markdown 指令，无法自动化测试，使用代码审查作为代理验证方法）

## Next Steps

- Proceed to `/quick-tasks` for task generation (quick mode)
