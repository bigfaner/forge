---
created: 2026-05-26
author: "faner"
status: Draft
---

# Proposal: Auto-Eval Configuration with Generic Config Key Resolution

## Problem

4 个文档评估 skill（eval-proposal、eval-prd、eval-ui、eval-design）在对应文档生成后，通过 `AskUserQuestion` 手动询问用户是否运行评估。每次交互都需要用户手动选择，增加了流水线的交互成本。项目已有 `auto.runTasks` 等 ModeToggle 配置驱动自动化的成功模式，eval 阶段应采用相同机制。

同时，当前 `forge config get/set` 的 key 路由系统（`getAutoKeyValue`、`getWorktreeKeyValue`、`setAutoConfigValue` 等）基于硬编码的分发器和 `autoModeField` switch/case，只支持两层深度（`auto.{field}.{subfield}`）。任何新的嵌套配置都需要修改路由核心代码。本提案将路由系统泛化为任意深度 key 遍历，从根本上消除这个扩展瓶颈。

### Evidence

- brainstorm、write-prd、tech-design 三个 skill 均在文档提交后通过 `AskUserQuestion` 询问是否运行 eval
- ui-design 是唯一例外——无条件自动运行 eval-ui，与其他 skill 行为不一致
- `auto.runTasks` 已证明 ModeToggle 配置模式可以成功消除冗余交互
- 当前 `getAutoKeyValue` 使用 `strings.SplitN(rest, ".", 2)` 只拆一层，`autoModeField` 通过 switch/case 映射，`parseAutoRaw` 硬编码 `modeFields` 列表——三者共同限制了 key 深度和新增配置的灵活性

### Urgency

随着 quick 流水线不断优化自动化流程（auto.runTasks、auto.consolidateSpecs 等），eval 的手动确认成为剩余的主要交互摩擦点。统一配置模式可降低流水线使用成本。同时，泛化 config key resolution 是一次投资长期收益的基础设施改进。

## Proposed Solution

### Part 1: Generic Config Key Resolution

重写 `GetConfigValue` / `SetConfigValue` 为基于任意深度 key 的通用遍历器：

**Get 路径**（反射）：
- 将 `key` 按 `.` 拆分为路径段，沿 Go Config struct 树递归走 reflect.Value
- 叶子节点的类型决定输出格式：bool → "true"/"false"，string → 原值，ModeToggle → "quick:true full:true"，[]string → 换行连接
- 默认值通过 `ReadConfig` 的 `applyDefaults` 已生效，反射读取的就是含默认值的最终状态

**Set 路径**（YAML Node）：
- 将 `key` 按 `.` 拆分为路径段，沿 YAML Node 树递归走 `findMappingKey`
- 叶子节点直接设值，中间节点不存在则创建 MappingNode
- 写回后下次 `ReadConfig` 会重新应用默认值

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

每个 ModeToggle 支持 `quick`/`full` 子键，通过 manifest 文件检测（`docs/features/<slug>/manifest.md` 中的 `mode: quick`）判断当前管道模式。

默认值：
- `proposal`: `quick: true, full: true` — 默认自动运行（proposal 是管道入口，尽早发现问题）
- `prd`: `quick: false, full: false` — 默认询问用户
- `uiDesign`: `quick: true, full: true` — 保持当前无条件自动行为，零回退风险
- `techDesign`: `quick: false, full: false` — 默认询问用户

### Innovation Highlights

1. **泛化 key resolution** 消除了 config 系统的扩展瓶颈——未来新增任何 `auto.*` 配置只需加 Go struct 字段 + 默认值，零路由代码改动
2. **eval 配置复用已有 ModeToggle 模式**，无新概念
3. **manifest 检测复用现有机制**（`feature_complete.go` 已有 manifest detection）

## Requirements Analysis

### Key Scenarios

- **自动评估 proposal（默认）**: 用户执行 `/quick`，brainstorm 提交 proposal 后自动运行 eval-proposal，无需手动确认
- **手动确认 prd eval**: 用户执行 `/write-prd`，提交 PRD 后询问是否运行 eval-prd（默认行为）
- **配置驱动的灵活控制**: 用户通过 `forge config set auto.eval.prd.full true` 开启 PRD 自动评估
- **ui-design 行为统一**: ui-design 从无条件自动运行改为读取配置，默认值保持 ON 以避免回退
- **任意深度 config get**: `forge config get auto.eval.proposal.quick` → `true`（三层深度自动支持）
- **任意深度 config set**: `forge config set auto.eval.proposal.quick false` → 正确写入嵌套结构

### Non-Functional Requirements

- 向后兼容：未配置时使用默认值，不影响现有行为（ui-design 默认保持 ON）
- 配置热生效：修改配置后无需重启即可生效
- 性能：泛化路由不应增加显著的延迟开销（反射开销在配置读取场景可忽略）

### Constraints & Dependencies

- Go reflect 包用于 struct 遍历（stdlib，无外部依赖）
- yaml.Node API 用于 YAML 树操作（已有依赖 gopkg.in/yaml.v3）
- manifest 文件格式需包含 `mode` 字段（已存在于 `feature_complete.go`）
- 4 个 skill 文件为 markdown，通过 Bash 调用 `forge config get` 实现 config check

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 无开发成本 | 保留交互摩擦 + 路由系统扩展瓶颈 | Rejected: 双重问题 |
| 扁平命名空间 (auto.evalProposal) | — | 零路由改动 | 不解决路由瓶颈，命名不直观 | Rejected: 治标不治本 |
| **泛化路由 + 嵌套 eval 配置** | auto.runTasks 模式 + 通用 key 遍历 | 根本性解决扩展瓶颈 + 灵活的 eval 控制 | 路由层重写的一次性投入 | **Selected: 最优长期方案** |
| 仅扩展路由支持三层 | — | 改动较小 | autoModeField/get/set 仍需修改，未来四层时又要改 | Rejected: 渐进式技术债 |

### Industry References

- **Viper** (Go): 支持 `viper.Get("a.b.c")` 任意深度 key 访问，证明反射/YAML 树遍历是成熟的配置访问模式
- **dotnet Microsoft.Extensions.Configuration**: 基于 `:` 分隔的路径遍历，支持任意深度

## Feasibility Assessment

### Technical Feasibility

完全可行。Go reflect 和 yaml.Node API 均为成熟技术。

**Get 路径实现要点**：
```go
func getStructValueByPath(v reflect.Value, segments []string) (string, error)
```
按 segments 递归走 reflect.Value：struct → 按 YAML tag 匹配字段；map → 按 key 查找；到达叶子节点时按类型格式化。

**Set 路径实现要点**：
```go
func setYAMLValueByPath(root *yaml.Node, segments []string, value string) error
```
按 segments 递归走 yaml.Node.MappingNode：找到对应 key 的 value node；不存在则创建。叶子节点设值后序列化回文件。

### Resource & Timeline

预计 3-4 小时：泛化路由重写（2h）+ eval 配置 + 4 个 skill（1h）+ 测试（1h）。

### Dependency Readiness

所有依赖（Go reflect、yaml.Node、manifest detection）均已就绪。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| config get/set 需要硬编码路由 | Occam's Razor | Refuted: 反射 + YAML 树遍历可实现通用路由，消除所有分发器 |
| 嵌套结构需要修改路由核心 | First Principles | Refuted: 泛化后嵌套和扁平实现成本相同 |
| 每次评估都需要用户确认 | Occam's Razor | Refuted: 熟练用户知道何时需要评估，手动确认是摩擦 |
| ui-design 的无条件自动评估是错误的 | Assumption Flip | Refined: 保持默认 ON，由配置控制而非硬编码 |
| proposal 默认应询问用户 | 5 Whys | Refuted: proposal 是流水线入口，自动评估可尽早发现问题 |

## Scope

### In Scope

- **Generic config key resolution**: 重写 `GetConfigValue`/`SetConfigValue` 为任意深度 key 遍历（反射 get + YAML Node set）
- **消除硬编码路由**: 删除 `autoModeField`、`getAutoKeyValue`/`getWorktreeKeyValue`/`getCoverageKeyValue`、`setAutoConfigValue`/`setWorktreeConfigValue`/`setCoverageConfigValue`
- **泛化 `parseAutoRaw`**: 递归扫描 auto 子树，不再硬编码 modeFields
- **Go CLI AutoConfig 新增 Eval 嵌套结构体**: proposal/prd/uiDesign/techDesign 各为 ModeToggle
- **AutoConfigDefaults 更新**: proposal {true,true}, prd {false,false}, uiDesign {true,true}, techDesign {false,false}
- **forge config get/set auto.eval.* 支持**: 任意深度自动支持
- **4 个 skill 增加 config check**: forge config get + manifest 文件模式推断
- **单元测试更新**: config_test.go、config_schema_test.go

### Out of Scope

- eval skill 本身（eval/SKILL.md）的行为修改
- 新的 eval 类型添加
- quick/full 管道流程的其他变更
- forge guide 文档更新（文档变更随代码一起完成）
- `forge config init` TUI 的更新（可后续迭代）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 反射遍历对 map 类型（SurfacesMap、CoverageConfig.ByType）的处理边界 | M | M | 为 map 类型写专门的处理分支，保持与当前行为一致 |
| YAML Node set 操作的序列化保真度（注释丢失、格式变化） | L | L | YAML v3 Node API 保留格式；当前实现已有 MarshalYAML 自定义 |
| parseAutoRaw 泛化后 raw tracking 精度变化 | L | M | 递归扫描保持相同的叶子节点追踪粒度，测试验证默认值行为不变 |
| 4 个 skill 的 config check 逻辑漂移 | L | M | 使用统一的配置检查模板，在 skill 中用 EXTREMELY-IMPORTANT 标注 |

## Success Criteria

- [ ] `forge config get auto.eval.proposal` 返回 `quick:true full:true`（三层深度）
- [ ] `forge config get auto.eval.proposal.quick` 返回 `true`（四层深度）
- [ ] `forge config set auto.eval.prd.full true` 正确写入嵌套 config
- [ ] `forge config get worktree.source-branch` 保持现有行为不变（回归测试）
- [ ] `forge config get coverage.coding.feature` 保持现有行为不变（回归测试）
- [ ] brainstorm 在 `auto.eval.proposal` 对应模式为 true 时跳过 AskUserQuestion，直接运行 eval-proposal
- [ ] brainstorm 在 `auto.eval.proposal` 对应模式为 false 时保持现有 AskUserQuestion 行为
- [ ] write-prd/tech-design/ui-design 同理遵循配置驱动
- [ ] 未配置时（config missing）：proposal 和 uiDesign 默认自动运行，prd 和 techDesign 默认询问
- [ ] 现有配置测试（config_test.go、config_schema_test.go）全部通过

## Next Steps

- Proceed to `/quick-tasks` for task generation (quick mode)
