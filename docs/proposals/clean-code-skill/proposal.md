---
created: 2026-05-17
author: "faner"
status: Completed
supersedes: reject-clean-code-task
---

# Proposal: Add Clean-Code Skill with Config-Gated Pipeline Integration

## Problem

Forge CLI 已注册 `TypeCleanCode` task类型和 `auto.cleanCode` config 选项，但 skill 层缺失：prompt template `code-quality-simplify.md` 直接委托给内置 `/simplify`，无 feature scope 限定、无 quality gate、无清理摘要。且 `typeToTemplate` 映射缺失，导致 `forge prompt get-by-task-id T-clean-code-1` 会报错。

### Evidence

- `forge-cli/pkg/prompt/prompt.go` 的 `typeToTemplate` 无 `TypeCleanCode` 条目 → `Synthesize()` 对该类型返回 `"unknown type"` 错误
- `code-quality-simplify.md` 仅调用 `Skill(skill="simplify")`，无 scope 控制，清理范围取决于 agent 的上下文窗口而非 feature diff
- 之前 `reject-clean-code-task` 的拒绝理由是 "pipeline 风险" — 但现在 `.forge/config.yaml` 的 `auto.cleanCode` 默认 `false`，用户可自行选择是否启用，风险已消除

### Urgency

`TypeCleanCode` 已注册但不可用（wiring 断裂），属于功能缺失。`auto.cleanCode` config 选项已暴露给用户但无对应实现。

## Proposed Solution

参照 Anthropic 官方 `code-simplifier` agent 模式，创建独立的 `clean-code` forge skill：

1. **Scope detection**: 通过 `git diff` 确定 feature 变更文件列表
2. **Independent cleanup**: 自有清理逻辑（非委托 `/simplify`），遵循五个原则：
   - **Preserve Functionality**: 只改写法不改行为
   - **Apply Project Standards**: 遵循 CLAUDE.md 和项目约定
   - **Enhance Clarity**: 减少复杂度、消除冗余、改善命名、移除无用注释
   - **Maintain Balance**: 不过度简化、不牺牲可读性换简洁
   - **Focus Scope**: 只处理 scope 内的变更文件
3. **Quality gate**: 清理后运行 `just test` 确认无回归
4. **Cleanup summary**: 输出修改了哪些文件、移除了哪些问题类型

删除旧的 `code-quality-simplify.md` 模板，新建 `code-quality-clean-code.md` 调用本 skill，修复 CLI wiring。

### User-Facing Behavior

**Pipeline 调用**（`auto.cleanCode: true`）:
1. `forge task index` 生成 `T-clean-code-1` task
2. Task executor 通过 prompt template 调用 `Skill(skill="forge:clean-code")`
3. Skill 自动 scope（`git diff main`），执行清理，运行 gate，输出摘要
4. 完成后通过 `forge:submit-task` 提交

**独立调用**（`/forge:clean-code`）:
1. 用户在 feature branch 上手动调用
2. 同样的 scope + cleanup + gate + summary 流程
3. 无 task 提交步骤

### Innovation Highlights

- **Config-gated pipeline integration**: 利用已有的 `.forge/config.yaml` 让用户选择是否启用，消除之前 rejection 的核心顾虑
- **Git-diff scope**: 不依赖 record 基础设施，用 `git diff` 自动确定清理范围，简单可靠
- **Independent cleanup logic**: 参照 `code-simplifier` 模式，不依赖内置 `/simplify`，可控性更强

## Requirements Analysis

### Key Scenarios

1. **Happy path**: 用户设置 `auto.cleanCode.full: true`，pipeline 生成 clean-code task，executor 执行后代码变干净、测试通过
2. **Standalone invocation**: 用户在 feature branch 手动执行 `/forge:clean-code`，清理当前变更
3. **No changes needed**: `git diff` scope 内无清理机会，skill 正常完成（无修改也算成功）
4. **Quality gate failure**: 清理引入回归，`just test` 失败，skill 报错并保留变更供用户检查
5. **No test infrastructure**: 项目无 `just test`，skill 跳过 gate 步骤（非硬性要求）

### Constraints & Dependencies

- `just test` 需要项目有 justfile 和 test recipe（部分项目可能没有）
- Git diff scope 需要 feature branch 与 base branch 有明确分叉点
- 遵循 forge distribution 约束：skill 文件在 `plugins/forge/skills/` 下
- Prompt template 文件通过 `//go:embed` 嵌入，文件名变更需同步 Go 代码

## Alternatives & Industry Benchmarking

### Reference

Anthropic 官方 `code-simplifier` agent（`anthropics/claude-plugins-official`）— 一个 agent 定义，聚焦代码简化：保留功能、应用项目标准、增强清晰度、保持平衡、聚焦 scope。本 skill 借鉴其五原则框架，并增加 git-diff scope、quality gate、cleanup summary 三个 forge 特有能力。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | typeToTemplate wiring 断裂；auto.cleanCode config 无实际功能 | Rejected: 功能缺失 |
| 仅修 wiring，保持 /simplify 直接调用 | — | 最小改动 | 无 scope、无 gate、无 summary；依赖内置 skill 不可控 | Rejected: 与手动 /simplify 无差别 |
| **clean-code skill（独立逻辑）+ wiring fix** | This proposal + code-simplifier | 独立可控；feature-scoped；quality gate；cleanup summary；config-gated | 增加 skill 维护 | **Selected** |

## Feasibility Assessment

### Technical Feasibility

完全可行。所有基础设施已就绪：task type 已注册、config 已支持。需改动 5-6 个文件。

### Resource & Timeline

2-3 个 task，预计 1-2h 实现时间。

### Dependency Readiness

无外部依赖。`code-simplifier` 模式为纯 prompt 工程。

## Scope

### In Scope

- `plugins/forge/skills/clean-code/SKILL.md` — skill 定义（scope → 独立清理逻辑 → gate → summary）
- `plugins/forge/commands/clean-code.md` — slash command 入口
- 删除 `forge-cli/pkg/prompt/data/code-quality-simplify.md`
- 新建 `forge-cli/pkg/prompt/data/code-quality-clean-code.md`（调用 forge:clean-code）
- `forge-cli/pkg/prompt/prompt.go` — 添加 `TypeCleanCode` 到 `typeToTemplate`（指向新模板）
- `scripts/version.txt` — version bump

### Out of Scope

- 修改 `/simplify` 的行为或能力
- 跨 feature 全局清理
- Hook 集成（不在 all-completed hook 中自动触发）
- Record-driven scope（用 git diff 替代）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 清理引入回归 | M | M | 内嵌 quality gate（`just test`），失败则保留变更供用户检查 |
| 项目无 `just test` 导致 gate 跳过 | M | L | Gate 为可选步骤：有 just test 则运行，无则跳过并在输出中注明 |
| 大 diff（100+ 文件）导致 agent 上下文溢出 | L | M | Skill 指示 agent 按文件分批处理 |

## Success Criteria

- [ ] `/forge:clean-code` 可独立调用，对当前 feature branch 变更执行 scoped cleanup
- [ ] `auto.cleanCode.full: true` 时 `forge task index` 生成 `T-clean-code-1` task
- [ ] `forge prompt get-by-task-id T-clean-code-1` 返回有效 prompt（不再报 "unknown type"）
- [ ] `code-quality-simplify.md` 已删除，新模板 `code-quality-clean-code.md` 调用 `forge:clean-code`
- [ ] Skill 清理逻辑遵循 code-simplifier 五原则（preserve、apply standards、enhance clarity、balance、focus scope）
- [ ] Skill 输出包含 cleanup summary（哪些文件被修改、移除了哪些问题类型）
- [ ] 有 `just test` 的项目：清理后自动运行 quality gate
- [ ] `go test -race -cover ./...` 通过

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
