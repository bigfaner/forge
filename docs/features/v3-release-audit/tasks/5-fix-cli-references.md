---
id: "5"
title: "Fix broken CLI cross-references (5 locations)"
priority: "P0"
estimated_time: "30m"
dependencies: ["1", "2"]
type: "doc"
mainSession: false
---

# 5: Fix broken CLI cross-references (5 locations)

## Description
4 处 `forge config get surface` 应为 `forge surfaces`（已重命名命令）。1 处 `test.execution` 引用过时命令。这些断裂引用导致 agent 运行时调用不存在的 CLI 命令。

## Reference Files
- `proposal.md#Problem` — Evidence table: CLI Reference 2 Major + Skill-CLI cross-ref 2 Critical
- `proposal.md#Scope` — P0.3 defines CLI reference fix scope
- `proposal.md#Success-Criteria` — "零断裂 CLI 引用" acceptance criterion

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/run-tests/rules/env-check.md` | `forge config get surface` → `forge surfaces` |
| `plugins/forge/skills/run-tests/SKILL.md` | `forge config get surface` → `forge surfaces`; `test.execution` → correct command |
| `plugins/forge/skills/gen-test-scripts/rules/run-to-learn.md` | `forge config get surface` → `forge surfaces` |
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | `forge config get surface` → `forge surfaces` |
| `plugins/forge/skills/init-justfile/SKILL.md` | `test.execution` → correct command |

## Acceptance Criteria
- [ ] `grep -r "forge config get surface" plugins/forge/` 返回 0 结果
- [ ] `grep -r "test\.execution" plugins/forge/` 返回 0 结果（或仅 config-schema 定义性引用）
- [ ] 所有替换命令与 `forge --help` 输出一致

## Hard Rules
- 先运行 `forge surfaces` 和 `forge --help` 确认正确命令名
- 使用 `replace_all` 确保不遗漏同一文件中的多处引用

## Implementation Notes
这是运行时阻断级问题：agent 按文档调用不存在的命令会直接失败。需逐一替换并验证。
