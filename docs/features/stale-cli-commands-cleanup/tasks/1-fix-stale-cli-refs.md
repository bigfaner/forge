---
id: "1"
title: "替换所有 skill 文件中的过期 CLI 命令引用"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 1: 替换所有 skill 文件中的过期 CLI 命令引用

## Description

多个 skill 文件引用了已从 forge CLI 移除的命令。agent 按指令执行时会失败。需要将所有过期命令引用替换为"读取项目文件推断"的明确指令。

涉及两类主要替换：
- `forge test detect`（10 处）→ 通过检查项目文件推断测试语言
- `forge test interfaces`（6 处）→ 通过检查项目结构和配置推断接口类型

另有两处低风险引用需一并修复。

## Reference Files
- `docs/proposals/stale-cli-commands-cleanup/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-contracts/SKILL.md` | 替换 `forge test detect` (L46) + `forge test interfaces` (L48) |
| `plugins/forge/skills/gen-test-cases/SKILL.md` | 替换 `forge test detect` (L21, L26) + `forge test interfaces` (L23, L78) |
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | 替换 `forge test detect` (L34) |
| `plugins/forge/skills/tech-design/SKILL.md` | 替换 `forge test detect` (L53, L57) |
| `plugins/forge/skills/quick-tasks/SKILL.md` | 替换 `forge test detect` (L39, L43) |
| `plugins/forge/skills/eval/rubrics/validate-ux-pipeline.md` | 替换 `forge test interfaces` (L7) + `forge test detect` (L15) |
| `plugins/forge/skills/eval/rules/pre-processing.md` | 替换 `forge test detect` (L8-9) |
| `plugins/forge/skills/eval/rubrics/test-cases.md` | 替换 `forge test interfaces` (L49) |
| `plugins/forge/skills/gen-test-cases/types/cli.md` | 替换 `forge deploy` 示例 (L57) |
| `plugins/forge/skills/eval/rubrics/cli-test-cases.md` | 替换 `forge task list` 示例 (L41) |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] `plugins/forge/skills/` 下所有文件中零引用 `forge test detect`
- [ ] `plugins/forge/skills/` 下所有文件中零引用 `forge test interfaces`
- [ ] `forge task list` 引用替换为 `forge task query`
- [ ] `forge deploy` 示例替换为实际存在的命令示例
- [ ] 替换后的指令明确描述 agent 应检查哪些项目文件来推断信息

## Hard Rules

- 替换文本必须是 agent 可直接执行的明确指令，不能只是"自行推断"
- 测试语言推断需覆盖主流语言：JavaScript/TypeScript (package.json)、Go (go.mod)、Python (pyproject.toml/setup.py)、Rust (Cargo.toml)
- 接口类型推断需覆盖：UI、TUI、API、CLI、Mobile

## Implementation Notes

替换 `forge test detect` 时，建议使用以下推断逻辑：
1. 检查项目根目录的 `package.json`（JS/TS）、`go.mod`（Go）、`Cargo.toml`（Rust）、`pyproject.toml`/`setup.py`（Python）
2. 检查 `.forge/config.yaml` 中的 `languages` 字段作为备选
3. 如以上均无法确定，要求用户手动配置

替换 `forge test interfaces` 时：
1. 检查 `docs/conventions/` 中是否有接口类型相关配置
2. 检查项目目录结构（如 `pages/` → UI, `api/` → API）
3. 检查 `.forge/config.yaml` 中的 `project-type` 字段
