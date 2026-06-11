---
created: "2026-05-26"
tags: [testing, architecture]
---

# Single Surface Project Task Frontmatter Should Not Leave Surface Fields Empty

## Problem
quick-tasks 生成的所有 task 文件中 surface-key 和 surface-type 均为空字符串，丢失了项目的 surface 信息。

## Root Cause
1. `forge surfaces --json` 对 single surface 项目（scalar 形式 `surfaces: cli`）返回 `[{"key":".","type":"cli"}]`
2. Agent 看到 key 为 `"."` 认为是占位符无意义，自行决定将两个字段都留空
3. 模板指令原文是 "single surface → use its key+type; mixed or no match → leave both empty"——Agent 误读了条件分支，将 single surface 当作 "无需填充" 处理

## Solution
Single surface 项目的 task frontmatter：
- `surface-type`：填入 forge surfaces 返回的 type 值（如 `"cli"`）
- `surface-key`：scalar 形式无显式 key，留空字符串

判断规则：读 `.forge/config.yaml` 的 `surfaces` 字段，scalar 形式（`surfaces: api`）→ type=该值、key 空；map 形式且只有一个 key → type=该 value、key=该 key。

## Reusable Pattern
当模板指令包含条件分支时（"A → 做 X；B → 做 Y"），逐条核对自己属于哪个分支，不要凭直觉跳过。surface-key/type 字段为空意味着"多 surface 项目且无法确定归属"——与 single surface 项目的语义不同。

## References
- `skills/quick-tasks/SKILL.md` — Step 2 Surface-Key/Type Inference
- `docs/lessons/pattern-surface-resolution-shortcut.md` — surface resolution 分层短路策略
