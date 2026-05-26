---
id: "5"
title: "更新 init-justfile skill 文件术语和 recipe 别名"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 5: 更新 init-justfile skill 文件术语和 recipe 别名

## Description

更新 init-justfile skill 的 SKILL.md 及 5 个 surface 规则文件。添加向后兼容的 justfile recipe alias（旧名 → 新名），更新 recipe 描述中的测试类型术语。

## Reference Files
- `docs/proposals/surface-test-type-model/proposal.md#Technical-Direction` — 定义了 justfile alias 兼容方案和过渡期机制
- `docs/proposals/surface-test-type-model/proposal.md#User-Facing-Experience` — justfile recipe 输出标签变更
- `docs/proposals/surface-test-type-model/proposal.md#Non-Functional-Requirements` — 向后兼容要求（2 版本过渡期）

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/init-justfile/SKILL.md` | 更新测试类型术语描述 |
| `plugins/forge/skills/init-justfile/rules/surfaces/cli.md` | 添加 `alias test-e2e := cli-test-functional`，更新 recipe 描述 |
| `plugins/forge/skills/init-justfile/rules/surfaces/tui.md` | 添加 alias，更新 recipe 描述 |
| `plugins/forge/skills/init-justfile/rules/surfaces/api.md` | 添加 alias，更新 recipe 描述 |
| `plugins/forge/skills/init-justfile/rules/surfaces/web.md` | 添加 alias，更新 recipe 描述 |
| `plugins/forge/skills/init-justfile/rules/surfaces/mobile.md` | 添加 alias，更新 recipe 描述 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria

- [ ] 每个 surface 的 justfile 规则文件包含向后兼容 alias（`alias test-e2e := <surface>-test-<type>`）
- [ ] alias 行带有 `# DEPRECATED: removed after v{current+2}` 注释
- [ ] recipe 描述使用 surface-specific 测试类型名称（如 "Run CLI functional tests" 而非 "Run e2e tests"）
- [ ] `just --list` 输出中 recipe 名称和描述清晰区分测试类型
- [ ] 聚合 recipe（`test`）的描述更新为 "Run all surface tests"
- [ ] 所有 rules 文件引用概念文档 `docs/reference/test-type-model.md`

## Hard Rules

- 现有 recipe 执行逻辑不变（命令、参数、环境变量），仅更新名称和描述
- alias 必须使用 just 原生 `alias` 语法
- 单 surface 项目保持无前缀 recipe（`test-cli-functional` 而非 `cli-test-functional`），多 surface 项目使用 `<surfaceKey>-test-<type>` 格式

## Implementation Notes

当前 justfile recipe 已使用 `<surfaceKey>-<step>` 命名（如 `cli-test`、`web-test`），不需要重命名 recipe 本身。需要做的是：(1) 在 recipe 描述中将 "e2e tests" 替换为 surface-specific 名称；(2) 添加 alias 供旧命令名使用。注意区分 surface key（用户自定义，如 "backend"）和 surface type（固定枚举，如 "api"），recipe 名称前缀用 surface key，测试类型标签用 surface type。
