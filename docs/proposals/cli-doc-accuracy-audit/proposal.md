---
created: "2026-06-07"
author: faner
status: Draft
intent: "doc"
---

# Proposal: CLI Documentation Accuracy Audit

## Problem

Forge 的两套 CLI 文档（guide.md 系统提示 + CLI `--help` 输出）与实际代码行为存在 19 处不一致，导致 AI agent 执行任务时使用错误命令或误解命令行为。

### Evidence

- **G1（关键）**：guide.md 引用 `forge task validate-index`，CLI 实际命令是 `forge task validate`。这已在 pm-work-tracker 项目中造成实际困惑（参见 `gotcha-forge-validate-command.md`）
- **C2**：`forge cleanup` help 说 "when task is completed"，但代码也清理 blocked/suspended/rejected 状态
- **C4**：`forge task validate` help 仅列 5 项验证，实际执行 12+ 项
- **C11**：`forge quality-gate` help 未提及自动创建 fix task 这一重要副作用
- **C7/C8**：`forensic search` 和 `forensic subagents` 完全没有 Long 描述

### Urgency

Guide.md 作为系统提示注入所有 agent 会话，其中的错误会影响每一个 forge 用户的每一次任务执行。修复成本极低（纯文本修改），收益即时可见。

## Proposed Solution

全面刷新 guide.md 的 CLI 参考部分，同步修复所有 CLI 命令的 cobra Long/Short 描述，使其准确反映代码行为。

### Innovation Highlights

 straightforward 文档修复，无创新成分。核心价值在于系统性审计方法——逐命令比对 help text 与 RunE 实现，而非仅修复已知报错。

## Requirements Analysis

### Key Scenarios

- Agent 读取 guide.md 后执行 CLI 命令，命令名、参数、行为描述均正确
- 开发者运行 `forge <command> --help`，输出完整反映该命令的全部功能和副作用
- 新命令或行为变更后，help text 是准确的行为参考

### Non-Functional Requirements

- CLI 二进制体积不变（仅修改字符串常量）
- 无性能影响

### Constraints & Dependencies

- 修改 guide.md 需遵循 forge-distribution.md 规范（guide.md 是 hook，作为系统提示分发）
- CLI help text 修改仅涉及 cobra Command 定义的 Long/Short 字段，不涉及逻辑变更

## Alternatives & Industry Benchmarking

### Industry Solutions

大多数 CLI 工具（如 kubectl, gh）通过自动化测试保证 help text 与代码一致。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 持续误导 agent | Rejected: 已有实际混淆案例 |
| 仅修 guide.md | 本次 lesson 的最小修复 | 解决最直接问题 | CLI help 仍不准确 | Rejected: 问题不止一处 |
| 自动化 help text freshness 测试 | kubectl 模式 | 一劳永逸 | 需要额外基础设施 | Deferred: 值得未来考虑 |
| **全面手动修复 guide + CLI help** | 本次审计结果 | 一次解决全部已知问题 | 未来新增命令仍需手动维护 | **Selected: 最高 ROI** |

## Feasibility Assessment

### Technical Feasibility

纯文本修改，100% 可行。无代码逻辑变更。

### Resource & Timeline

单人 1-2 小时。19 处修改，每处平均 5 分钟。

### Dependency Readiness

无外部依赖。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "guide.md 中的命令名是准确的" | Evidence-based: 实际运行 CLI 比对 | Overturned: `validate-index` 不存在 |
| "CLI help text 完整反映代码行为" | 逐函数审计 RunE vs Long | Overturned: 11 处不一致 |
| "低严重度问题可以暂缓" | Occam's Razor | Confirmed: 但全修成本极低，不如一次清完 |

## Scope

### In Scope

**Guide.md 修复（8 处）：**
- G1: `validate-index` → `validate [file]`
- G2: `quality-gate` 描述更新为准确反映实际行为
- G3: `cleanup` 描述从 "clean stale artifacts" 改为具体行为说明
- G4: `task submit` 补充 `--quiet` 标志
- G5: 新增 `task query <id-or-key>` 命令
- G6: 新增 `task check-deps` 命令
- G7: 新增 `feature list` 命令
- G8: `task list` 补充 `--tree` 标志

**CLI Help Text 修复（11 处）：**
- C1: `forge init` Long 补充 surface detection 步骤
- C2: `forge cleanup` Short/Long 补充 blocked/suspended/rejected 状态
- C3: `forge task claim` Long 补充 auto-unblock 行为
- C4: `forge task validate` Long 补充全部 12+ 项验证
- C5: `forge task add` Long 补充使用概述
- C6: `forge feature` Long 补充 `set` 子命令和行为差异说明
- C7: `forge forensic search` 新增 Long 描述
- C8: `forge forensic subagents` 新增 Long 描述
- C9: `forge worktree status` Long 补充 UNPUSHED 字段
- C10: `forge fact summary` Long 补充 COVERAGE 指标
- C11: `forge quality-gate` Long 补充 fix task 自动创建、retry-once、docs-only 跳过

### Out of Scope

- 自动化 help text freshness 测试（值得未来单独提案）
- justfile 中 `validate-index` 正则的更新（不影响 agent 行为）
- `forge-cli/docs/features/` 中历史测试用例文档的更新
- CLI 二进制行为变更（仅修改 help 文本）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 修改 Long 文本时引入新的不准确描述 | L | L | 逐条与 RunE 代码核对 |
| guide.md 修改引入 markdown 格式问题 | L | L | 审查渲染结果 |
| 未来 CLI 行为变更再次导致 help text 漂移 | M | M | 在 forge-distribution.md 中添加 "修改命令时同步 help text" 的规范 |

## Success Criteria

- [ ] 运行 `forge task validate-index docs/features/any/tasks/index.json` 返回 "unknown command" 错误（确认旧命令已从 guide 中移除）
- [ ] guide.md 中每一个 `forge ` 开头的命令名都与 `forge --help` 或 `forge <group> --help` 输出的 Available Commands 完全匹配
- [ ] 对 11 个 CLI 修改的命令分别运行 `--help`，Long 描述中提及的每一项功能都能在对应 RunE 函数中找到代码实现
- [ ] `forge task validate --help` 的 Validations 列表覆盖全部实际执行的验证步骤
- [ ] 修改后 `go build ./...` 和 `go test ./...` 全部通过（确保字符串修改未破坏编译）

## Next Steps

- 本提案 scope 为纯文档修复，无需 PRD 或 tech-design，可直接执行
- 执行后运行 `go build ./...` 和 `go test ./...` 验证
