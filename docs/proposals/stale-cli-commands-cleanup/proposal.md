---
created: 2026-05-20
author: "fanhuifeng"
status: Draft
---

# Proposal: 清理 Skill 中引用的已移除 CLI 命令

## Problem

多个 skill 文件引用了已从 forge CLI 中移除的命令（`forge test detect`、`forge test interfaces`），agent 按指令执行时会失败。

### Evidence

CLI 测试文件 `forge-cli/internal/cmd/test_test.go:54` 明确列出已移除的子命令：`detect`、`get`、`interfaces`、`framework`。当前 `forge test` 仅注册了 `promote`、`run-journey`、`verify`。

但在 skill 文件中仍有 **16 处** 引用这些不存在的命令作为实际执行指令：

- `forge test detect`：7 个文件，10 处引用
- `forge test interfaces`：5 个文件，6 处引用

### Urgency

agent 遵循 skill 指令时直接运行这些命令会失败，导致工作流中断。这是已确认的功能缺陷。

## Proposed Solution

1. 将所有 `forge test detect` 引用替换为"读取项目文件推断测试语言"的指令（检查 package.json/go.mod/Cargo.toml/pyproject.toml 等）
2. 将所有 `forge test interfaces` 引用替换为"读取项目结构和配置推断接口类型"的指令
3. 修复低风险引用（rubric 中的 `forge task list` → `forge task query`，示例中的 `forge deploy`）
4. 新增 `docs/conventions/forge-cli-reference.md` 记录所有有效 CLI 命令，防止未来再次出现过期引用

### Innovation Highlights

无特殊创新——这是标准的文档-代码同步维护工作。通过新增 CLI 命令参考文档降低 recurrence 风险。

## Requirements Analysis

### Key Scenarios

- Agent 执行 gen-contracts/gen-test-cases/quick-tasks 等 skill 时，需要确定项目的测试语言 → 直接读项目文件推断
- Agent 需要确定项目的接口类型（UI/TUI/API/CLI） → 读项目结构和 docs/conventions/ 推断
- Skill 作者编写新 skill 时需要知道哪些 CLI 命令可用 → 查阅 forge-cli-reference.md

### Non-Functional Requirements

- 修改后的 skill 指令必须足够明确，agent 无需额外 CLI 命令即可完成语言/接口检测

### Constraints & Dependencies

- 无外部依赖
- 仅修改 skill 文档文件，不涉及 CLI 代码变更

## Alternatives & Industry Benchmarking

### Industry Solutions

标准的代码-文档同步问题，通常通过文档化 + lint 检查解决。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | agent 执行失败，用户体验差 | Rejected: 已确认的功能缺陷 |
| 恢复 CLI 命令 | — | skill 不需要改 | CLI 维护成本高，命令功能可通过文件推断完成 | Rejected: 过度工程化 |
| 自动化 lint 检查 | — | 长期防护 | 需要额外开发和维护 lint 脚本 | Deferred: 投入产出比不高 |
| **修复引用 + CLI 参考文档** | 行业标准 | 解决当前问题 + 降低 recurrence 风险 | 无自动化保障 | **Selected: 最小有效方案** |

## Feasibility Assessment

### Technical Feasibility

完全可行。所有修改都是文本替换，目标文件路径已全部定位。

### Resource & Timeline

1 人 30 分钟内可完成。

### Dependency Readiness

无外部依赖。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 需要恢复 CLI 命令来解决此问题 | Occam's Razor | Overturned: agent 直接读项目文件即可获取相同信息，无需 CLI 中间层 |
| 低风险引用可以忽略 | Stress Test | Refined: `forge task list` 在 rubric 中作为示例可能导致用户混淆，一并修复成本极低 |

## Scope

### In Scope

- 替换 7 个 skill 文件中的 `forge test detect` 引用（10 处）
- 替换 5 个 skill 文件中的 `forge test interfaces` 引用（6 处）
- 修复 `forge task list` → `forge task query`（1 处）
- 修复 `forge deploy` 示例（1 处）
- 新增 `docs/conventions/forge-cli-reference.md`

### Out of Scope

- 新增或恢复任何 CLI 命令
- 自动化 lint 检查脚本
- skill 逻辑重构（仅替换命令引用）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 替换后的"读文件推断"指令不够具体，agent 行为不一致 | L | M | 在替换文本中提供具体的文件检查清单和推断逻辑 |
| 未来新增/移除 CLI 命令时文档未同步更新 | M | L | forge-cli-reference.md 添加到 docs/conventions/，由 consolidate-specs 管理 |

## Success Criteria

- [ ] 所有 skill 文件中零引用 `forge test detect` 和 `forge test interfaces`
- [ ] `docs/conventions/forge-cli-reference.md` 包含完整的有效 CLI 命令清单
- [ ] 替换后的指令明确描述了 agent 应如何通过项目文件推断测试语言和接口类型

## Next Steps

- Proceed to `/quick-tasks` to generate fix tasks
