# Lesson: guide.md 缺少 proposals/ 目录导致 agent 无法遵循规范

## 问题

tui-ui-design 提案被写到项目根目录 `proposal-tui-ui-design.md`，格式也不符合 `docs/proposals/{slug}/proposal.md` 的项目规范。项目中有 30 个现有提案遵循该规范，但 agent 没有遵循。

### 根因

`guide.md`（通过 session-start hook 注入每个 Claude 会话）的 Project-Level Documents 部分没有收录 `docs/proposals/` 目录。agent 的上下文中完全没有这个规范的信息。

**Why:** guide.md 作为 forge 注入的唯一全局规范文档，是 agent 了解项目目录约定的唯一可靠来源。当某个目录未在 guide.md 中记录时，agent 只能通过手动扫描 `docs/` 来发现规范——这不是可靠的工作模式。

## 规则

`guide.md` 的 Project-Level Documents 必须收录所有有规范约束的 `docs/` 子目录。新增受规范约束的目录时，必须同步更新 `guide.md`。

**How to apply:**
1. 当创建新的 `docs/` 子目录并赋予规范约束时，同步更新 `plugins/forge/hooks/guide.md` 的 Project-Level Documents 代码块
2. 定期（如每次 eval-harness）检查 `docs/` 下的实际目录与 guide.md 记录是否一致

## 修复

已在 guide.md 中补充 `proposals/` 条目。

## Tags

`guide`, `directory-conventions`, `process`, `forge-improvement`
