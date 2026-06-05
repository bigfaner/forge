---
created: 2026-06-04
author: "faner"
status: Draft
intent: "new-feature"
---

# Proposal: Eval Target & Iterations Configuration

## Problem

Eval 系统（7 种类型）的 target score 和 iterations 默认值硬编码在 rubric 文件 frontmatter 中，用户无法通过项目配置覆盖。需要自定义时只能每次手动传 CLI 参数（`--target 850 --iterations 5`）或直接编辑插件内部的 rubric 文件——前者需要记忆和重复操作，后者影响插件升级。

### Evidence

- rubric frontmatter 是唯一默认值来源：proposal 900/3, prd 900/3, ui 950/3, journey 850/3, contract 850/3, consistency 900/3（来自 `rubric-reference.md`）
- 已有 `auto.eval.*` 配置控制 eval **是否自动运行**（bool），但无配置控制 **运行参数**（target/iterations）
- 用户每次手动调用 `/eval-proposal --target 850 --iterations 5` 需要记住类型对应的参数值
- `auto-eval-config` 提案（已实现）解决了 on/off 问题，但 target/iterations 仍是盲区

### Urgency

随着 eval 自动运行（`auto.eval.*`）的启用，eval 在无用户干预下自动执行时使用 rubric 默认值。用户无法为自动运行场景设定不同的 target/iterations——配置化控制运行参数是自动化的必要补全。

## Proposed Solution

在 `.forge/config.yaml` 中新增顶层 `eval` 配置块，包含 7 种 eval 类型各一个 `{target, iterations}` 子结构。解析优先级：CLI args > config.yaml > rubric frontmatter。

**配置结构**：

```yaml
eval:
  proposal:
    target: 900
    iterations: 3
  prd:
    target: 900
    iterations: 3
  design:
    target: 900
    iterations: 3
  ui:
    target: 950
    iterations: 3
  journey:
    target: 850
    iterations: 3
  contract:
    target: 850
    iterations: 3
  consistency:
    target: 900
    iterations: 3
```

**解析层**：eval-* 命令（非 eval skill）负责通过 `forge config get eval.<type>.target|iterations` 读取配置，按优先级链解析后传递给 eval skill。eval skill 本身不感知 config，保持单一职责。

**init 行为**：`forge config init` 不询问 target/iterations，但自动从 rubric frontmatter 读取默认值填充到 config.yaml 输出中，供用户直接编辑。

### Innovation Highlights

无显著创新——本方案是对现有 config 系统的增量扩展，复用已由 `auto-eval-config` 建立的泛化反射路由机制。设计决策点在于"由命令层解析 config 并传递给 skill"而非"让 skill 自行读取 config"，这保持了 eval skill 的独立性和可测试性。

## Requirements Analysis

### Key Scenarios

- **自动 eval 使用项目默认值**: brainstorm 自动触发 eval-proposal 时，eval-proposal 命令从 config 读取 target=900, iterations=3，传递给 eval skill
- **手动调用覆盖配置**: 用户执行 `/eval-proposal --target 850` 时，CLI 参数覆盖 config 中的 900
- **config 未配置时回退 rubric**: config.yaml 无 `eval` 块时，eval-* 命令回退到 rubric frontmatter 默认值
- **config 部分覆盖**: 用户只设置 `eval.proposal.target: 850`，iterations 仍从 rubric 读取默认值 3
- **`forge config init` 生成完整 eval 块**: init 时从 rubric 文件读取各类型默认值，填充到 config.yaml 的 `eval` 块中
- **`forge config get eval.proposal.target`** 返回 `900`
- **`forge config set eval.proposal.target 850`** 正确写入配置

### Non-Functional Requirements

- 向后兼容：未配置 `eval` 块时行为与当前完全一致（rubric 默认值生效）
- 配置热生效：修改 config.yaml 后下次 eval 调用立即使用新值
- 不影响 auto.eval 开关逻辑（`auto.eval.*` 控制是否运行，`eval.*` 控制运行参数，职责分离）

### Constraints & Dependencies

- Go 反射路由（`GetConfigValue`/`SetConfigValue`）已由 `auto-eval-config` 实现泛化，支持任意深度 key 遍历
- eval-* 命令为 markdown 文件，通过 Bash 调用 `forge config get` 读取配置
- rubric frontmatter 中 `target` 和 `iterations` 字段为 YAML 整数类型
- `forge config init` 需要读取 rubric 文件提取默认值——rubric 文件位于插件目录 `plugins/forge/skills/eval/rubrics/<type>.md`

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 每次手动传参；自动 eval 无法自定义参数 | Rejected: 自动化场景缺失 |
| 编辑 rubric 文件 | — | 直接修改默认值 | 侵入插件内部；升级时被覆盖；影响所有项目 | Rejected: 不属于用户可配置层 |
| **config.yaml 中间层** | auto.eval 模式延伸 | 项目级覆盖；向后兼容；init 可生成默认值；与 auto.eval 职责分离 | 需改 Go 结构体 + init 逻辑 + 7 个 eval 命令 | **Selected: 最符合现有架构** |
| eval skill 内部读取 config | — | 命令层无需改动 | eval skill 感知 config 违反单一职责；skill 为 markdown 无法调用 Go API | Rejected: 职责倒置 |

## Feasibility Assessment

### Technical Feasibility

完全可行。反射路由已支持任意深度 key，新增 `EvalSettings` 结构体后 `forge config get eval.proposal.target` 自动可用。

**Go 结构体变更**：

```go
type EvalTypeSettings struct {
    Target     *int `yaml:"target,omitempty"`
    Iterations *int `yaml:"iterations,omitempty"`
}

type EvalSettings struct {
    Proposal    EvalTypeSettings `yaml:"proposal"`
    Prd         EvalTypeSettings `yaml:"prd"`
    Design      EvalTypeSettings `yaml:"design"`
    Ui          EvalTypeSettings `yaml:"ui"`
    Journey     EvalTypeSettings `yaml:"journey"`
    Contract    EvalTypeSettings `yaml:"contract"`
    Consistency EvalTypeSettings `yaml:"consistency"`
}
```

`Target` 和 `Iterations` 使用 `*int` 指针类型：`nil` 表示未配置（回退 rubric），有值则覆盖。反射路由遇到 `nil` 指针时返回 `errKeyNotFound`，eval-* 命令据此判断是否回退。

**eval-* 命令变更**（以 eval-proposal 为例）：

```
# 解析 target
TARGET=$(forge config get eval.proposal.target 2>/dev/null)
if [ $? -eq 0 ] && [ -n "$TARGET" ]; then
  TARGET_ARG="--target $TARGET"
fi

# 解析 iterations
ITERATIONS=$(forge config get eval.proposal.iterations 2>/dev/null)
if [ $? -eq 0 ] && [ -n "$ITERATIONS" ]; then
  ITERATIONS_ARG="--iterations $ITERATIONS"
fi

Skill(skill="forge:eval", args="--type proposal $TARGET_ARG $ITERATIONS_ARG")
```

CLI `--target`/`--iterations` 参数优先于 config 值——命令层在拼接 args 时，如果用户显式传了 CLI 参数，则不使用 config 值。

### Resource & Timeline

预计 3-4 小时：Go 结构体 + config get/set 支持（1h）+ init 默认值生成（1h）+ 7 个 eval-* 命令更新（1h）+ 测试（0.5-1h）。

### Dependency Readiness

所有依赖已就绪。反射路由已实现。rubric 文件 frontmatter 格式稳定。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| eval skill 需要感知 config | XY Detection | Refined: eval-* 命令层负责解析 config 并传递参数给 skill，skill 保持无 config 依赖。用户需要的是"配置生效"，不是"skill 读取配置" |
| init 应该询问 target/iterations | Occam's Razor | Refuted: 交互式询问 7 个类型 × 2 个参数 = 14 个问题过于冗长。直接输出 rubric 默认值到 config.yaml 更简洁 |
| ui 类型需要按平台分别配置 | 5 Whys | Refuted: 三个平台（web/mobile/tui）默认值相同（950/3），无差异化需求。统一配置更简洁 |

## Scope

### In Scope

- **Go `EvalSettings` 结构体**: 新增到 `forgeconfig/config.go`，7 种 eval 类型各含 `target`（*int）和 `iterations`（*int）
- **`Config` 结构体扩展**: 新增 `Eval *EvalSettings` 字段
- **`forge config get eval.<type>.target|iterations`**: 通过反射路由自动支持
- **`forge config set eval.<type>.target|iterations`**: 通过反射路由自动支持
- **`forge config init` 输出 eval 块**: 从 rubric frontmatter 读取默认值填充
- **7 个 eval-* 命令更新**: 各命令通过 `forge config get` 读取 config 值，拼接为 skill args
- **`forge config get eval`**: 返回所有 eval 类型的汇总输出
- **单元测试**: config_test.go 新增 EvalSettings 相关测试

### Out of Scope

- `auto.eval` 开关逻辑变更（已有，不改）
- eval skill（eval/SKILL.md）内部逻辑变更
- rubric 文件内容变更
- ui 子类型（web/mobile/tui）分别配置
- journey/contract pipeline 硬编码逻辑
- eval 评分、scorer、expert 选择
- validate-code/validate-ux 类型的配置支持（这两种为单次运行，无迭代意义）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| `*int` 指针在反射路由中的序列化——`nil` 指针应返回 `errKeyNotFound` 而非空字符串 | L | M | `getByPath` 已有 nil 指针处理（返回 `errKeyNotFound`）；新增测试验证 `forge config get eval.proposal.target` 在未配置时正确返回 not found |
| init 读取 rubric 文件路径——rubric 文件位于插件目录，init 需要感知插件安装路径 | M | M | init 命令通过 `forge config get` 已能感知项目根目录；rubric 文件路径可通过项目目录结构推断（`plugins/forge/skills/eval/rubrics/<type>.md`）或嵌入 Go 代码作为常量 |
| 7 个 eval-* 命令的 config check 模板不一致——markdown 命令无编译时强制一致性 | M | L | 使用统一的 bash 模板（resolve → fallback → pass args），保持 7 个命令结构一致 |
| `forge config get eval.proposal.target` 返回空字符串 vs `errKeyNotFound`——未配置时的行为歧义 | L | M | `*int` 为 nil 时 `getByPath` 返回 `errKeyNotFound`，CLI 输出错误码 1；eval-* 命令据此判断回退。保持与现有 `auto.eval` 的 get 行为一致 |

## Success Criteria

- [ ] `forge config get eval.proposal.target` 在 config.yaml 已配置时返回正确数值（如 `900`）
- [ ] `forge config get eval.proposal.target` 在未配置时返回 exit code 1（`errKeyNotFound`）
- [ ] `forge config set eval.proposal.target 850` 正确写入 config.yaml
- [ ] `forge config set eval.journey.iterations 5` 正确写入 config.yaml
- [ ] `forge config get eval` 返回所有已配置 eval 类型的汇总输出
- [ ] `forge config init` 生成的 config.yaml 包含完整 `eval` 块，默认值与 rubric frontmatter 一致
- [ ] eval-proposal 命令在 config 已配置 target=850 时传递 `--target 850` 给 eval skill
- [ ] eval-proposal 命令在 config 未配置时不传递 `--target`，eval skill 使用 rubric 默认值
- [ ] CLI `--target` 参数优先于 config 值——用户传 `/eval-proposal --target 800` 时使用 800 而非 config 中的 850
- [ ] 所有 7 个 eval-* 命令（proposal/prd/design/ui/journey/contract/consistency）正确读取并传递 config 值
- [ ] 现有配置测试（config_test.go）全部通过，无回归
- [ ] 现有 eval 自动运行逻辑（auto.eval.*）行为不变

consistency_check_result:
  status: pass
  pairs_checked: 66
  conflicts_found: 0

## Next Steps

- Proceed to `/write-prd` to formalize requirements
