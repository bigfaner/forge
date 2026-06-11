---
created: 2026-05-27
author: "faner"
status: Draft
---

# Proposal: Forge Init 支持 Auto-Eval 配置

## Problem

两个问题：

1. **`forge init` 缺少 auto.eval 配置入口**：交互式配置流程覆盖 6 类自动行为（test/consolidateSpecs/cleanCode/validation/runTasks/knowledgeSave/gitPush），但缺少 `auto.eval.*`。用户无法通过 init 引导设置 eval 策略，只能事后 `forge config set`——大多数用户甚至不知道 auto.eval 存在。

2. **EvalConfig 使用 ModeToggle 但 eval 不依赖 mode**：brainstorm 执行时没有 feature 目录，`forge config get mode` 返回 `"none"`，导致 `auto.eval.proposal.none` key 不存在 → 始终走 FALLBACK_ASK。write-prd/ui-design/tech-design 仅存在于 full 流水线，`$MODE` 永远是 `"full"`。**四类 eval 配置的 ModeToggle 实质上是无效抽象**——quick 字段全是死代码，proposal 的 full/quick 区分在 brainstorm 阶段无意义。

### Evidence

- `autoBehaviorPrompts()`（init.go:257-337）包含 13 个提示，无 eval 相关
- `detectModeFromPath()`（config.go:127-171）依赖 `docs/features/<slug>/` 路径，brainstorm 在项目根目录执行 → mode = "none"
- brainstorm SKILL.md 的 eval 检查用 `auto.eval.proposal.$MODE`，MODE="none" 时 key 不存在 → 永远 FALLBACK_ASK → auto.eval.proposal 配置从未生效

### Urgency

auto.eval.proposal 已配置默认值 `quick:true, full:true`（暗示自动运行），但因 mode 检测失败从未生效。这是一个已上线但从未工作的功能。

## Proposed Solution

### Part 1: EvalConfig 全面扁平化为 bool

```go
// Before
type EvalConfig struct {
    Proposal   ModeToggle `yaml:"proposal"`
    Prd        ModeToggle `yaml:"prd"`
    UiDesign   ModeToggle `yaml:"uiDesign"`
    TechDesign ModeToggle `yaml:"techDesign"`
}

// After
type EvalConfig struct {
    Proposal   bool `yaml:"proposal"`
    Prd        bool `yaml:"prd"`
    UiDesign   bool `yaml:"uiDesign"`
    TechDesign bool `yaml:"techDesign"`
}
```

**理由**：

- **brainstorm**：执行时无 mode 上下文（mode="none"），proposal eval 不可能区分 quick/full
- **PRD/UI 设计/技术设计**：仅存在于 full 流水线，`$MODE` 永远是 "full"，区分 quick/full 无意义
- **结论**：四类 eval 都不需要 mode 维度，`bool` 是最小充分类型

### Part 2: Init 流程新增 4 个 Eval 提示

在 `autoBehaviorPrompts()` 中新增 4 个提示（gitPush 之前）：

| # | 提示文本 | 配置键 | 默认值 |
|---|---------|--------|--------|
| 1 | "Auto-evaluate proposals?" | `auto.eval.proposal` | true |
| 2 | "Auto-evaluate PRD documents?" | `auto.eval.prd` | false |
| 3 | "Auto-evaluate UI designs?" | `auto.eval.uiDesign` | true |
| 4 | "Auto-evaluate tech designs?" | `auto.eval.techDesign` | false |

与 `gitPush` 的单开关风格一致，无 quick/full 前缀。

### Part 3: 全部 4 个 SKILL.md Eval 检查统一适配

从：

```bash
MODE=$(forge config get mode 2>/dev/null)
if [ $? -ne 0 ]; then
  echo "FALLBACK_ASK"
else
  EVAL_ENABLED=$(forge config get auto.eval.proposal.$MODE 2>/dev/null)
  ...
fi
```

改为：

```bash
EVAL_ENABLED=$(forge config get auto.eval.proposal 2>/dev/null)
if [ "$EVAL_ENABLED" = "true" ]; then
  echo "AUTO_RUN"
elif [ "$EVAL_ENABLED" = "false" ]; then
  echo "SKIP"
else
  echo "FALLBACK_ASK"
fi
```

去掉 MODE 查询和 `$MODE` 后缀。4 个 SKILL.md 统一模式。

### Innovation Highlights

无创新。修正一个从未工作的功能（auto.eval.proposal 的 mode 检测缺陷），补齐 init 配置入口，消除过度抽象（ModeToggle → bool）。

## Requirements Analysis

### Key Scenarios

- **用户首次 init**：TUI 引导设置 4 个 eval 开关
- **Brainstorm eval auto-run**：`auto.eval.proposal=true` → 提交后自动运行 eval-proposal（此前因 mode 检测失败从未生效）
- **PRD eval manual**：`auto.eval.prd=false`（默认）→ 提交后询问用户
- **已有 config.yaml 兼容**：旧格式 `auto.eval.prd: {quick: false, full: false}` 需兼容读取

### Constraints & Dependencies

- 依赖泛化 config key resolution（反射路由）支持扁平 bool 的 get/set
- `forge config init`（stdin 版本）与 `forge init`（TUI 版本）共享 `askAutoBehavior()` 函数，改动同时覆盖两个入口

## Alternatives & Industry Benchmarking

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | 零改动 | auto.eval.proposal 从未生效；init 缺少 eval 配置 | Rejected: 功能不工作 |
| 仅在 init 添加提示，保留 ModeToggle | 改动最小 | proposal eval 在 brainstorm 仍因 mode="none" 不工作 | Rejected: 未修核心缺陷 |
| **全部改为 bool + init 提示** | 修正功能缺陷 + 补齐配置入口 | 4 个 SKILL.md + config.go 需适配 | **Selected: 唯一让 auto.eval.proposal 真正工作的方案** |

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| eval 配置需要区分 quick/full 模式 | XY Detection: 真正的需求是"eval 是否自动运行"，不是"按模式区分 eval 行为" | Overturned: mode 维度是过度设计 |
| brainstorm 可以通过 `forge config get mode` 获取当前模式 | Evidence: brainstorm 在项目根目录执行，mode 始终为 "none" | Overturned: mode 检测在此场景无效 |
| 非 proposal 的 ModeToggle 有意义 | Assumption Flip: 如果 $MODE 永远是 "full"，ModeToggle 等价于 bool | Confirmed: 无意义的间接层 |

## Scope

### In Scope

- `EvalConfig` 全面改为 bool（4 个字段）
- `AutoConfigDefaults()` 简化 eval 默认值
- `autoBehaviorPrompts()` 新增 4 个 eval 提示
- 全部 4 个 SKILL.md eval 检查片段统一适配
- config_test.go 更新
- 旧格式 config.yaml 兼容读取（ModeToggle → bool 迁移）

### Out of Scope

- 泛化 config key resolution（已实现）
- UI/UX 设计（沿用现有 TUI 模式）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 已有 config.yaml 中 ModeToggle 格式不兼容 | M | M | ReadConfig 中对 bool 字段遇到 map 值时取 `full` 子键值 |
| SKILL.md 改动遗漏 | L | L | grep `auto.eval.*\$MODE` 验证无残留 |

## Success Criteria

- [ ] `forge init` 交互流程包含 4 个 eval 提示，默认值与 `AutoConfigDefaults()` 一致
- [ ] brainstorm 执行后 `auto.eval.proposal=true` 时自动运行 eval-proposal（此前从未工作）
- [ ] `forge config get auto.eval.proposal` 返回 `"true"` 或 `"false"`（非 ModeToggle 格式）
- [ ] `forge config set auto.eval.prd true` 正确写入并持久化
- [ ] 全部 4 个 SKILL.md 的 eval 检查不再依赖 `$MODE`
- [ ] 已有 config.yaml 中 ModeToggle 格式兼容读取

## Next Steps

- Proceed to task breakdown（范围较小，可直接 `/quick-tasks`）
