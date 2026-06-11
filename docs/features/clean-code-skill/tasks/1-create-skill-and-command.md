---
id: "1"
title: "Create clean-code skill and command"
priority: "P0"
estimated_time: "30m"
dependencies: []
scope: "backend"
breaking: false
type: "feature"
mainSession: false
---

# 1: Create clean-code skill and command

## Description

创建 `clean-code` forge skill 和对应的 slash command 入口。Skill 参照 Anthropic 官方 `code-simplifier` agent 模式，实现独立的代码清理逻辑（不依赖内置 `/simplify`），包含 feature-scoped cleanup、quality gate、cleanup summary 三个核心能力。

## Reference Files

- `docs/proposals/clean-code-skill/proposal.md` — Source proposal
- `plugins/forge/skills/` — Existing skill directory for pattern reference
- `plugins/forge/commands/` — Existing command directory for pattern reference

## Acceptance Criteria

- [ ] `plugins/forge/skills/clean-code/SKILL.md` exists with complete skill definition
- [ ] Skill workflow: scope detection (git diff) → code cleanup (5 principles) → quality gate (just test, optional) → cleanup summary
- [ ] Cleanup logic follows code-simplifier 5 principles: Preserve Functionality, Apply Project Standards, Enhance Clarity, Maintain Balance, Focus Scope
- [ ] Skill can be invoked standalone via `/forge:clean-code`
- [ ] `plugins/forge/commands/clean-code.md` exists as slash command entry point

## Hard Rules

- 遵循 forge distribution 约束（见 `docs/conventions/forge-distribution.md`）
- Skill 只修改 `git diff` scope 内的文件
- Quality gate 为可选：有 `just test` 则运行，无则跳过

## Implementation Notes

- 参考 Anthropic `code-simplifier` agent: https://github.com/anthropics/claude-plugins-official/blob/main/plugins/code-simplifier/agents/code-simplifier.md
- Scope detection 用 `git diff --name-only main`（或用户指定的 base branch）
- 参考现有 skill 结构（如 `learn-lesson`、`brainstorm`）了解 SKILL.md frontmatter 格式
- Command 文件参考 `plugins/forge/commands/` 下已有的 command 格式
