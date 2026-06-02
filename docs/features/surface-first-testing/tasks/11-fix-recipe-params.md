---
id: "11"
title: "Fix recipe parameter mismatch + run-tests surface-key usage"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 11: Fix recipe parameter mismatch + run-tests surface-key usage

## Description
修复 init-justfile 和 run-tests 之间的 recipe 签名不匹配问题。init-justfile 生成的 `just <key>-test` recipe 不接受参数，但 run-tests 调用 `just <surface>-test <journey>` 传入 journey 参数。需要：(1) init-justfile 的 recipe 签名改为接受可选 journey 参数；(2) run-tests 在多 surface 项目中使用 surface-key（非 surface-type）构造 recipe 名称，与 init-justfile 的 `<key>-<verb>` 命名对齐。

## Reference Files
- `docs/proposals/surface-first-testing/proposal.md` — Proposed Solution (下游 skill 适配)
- `plugins/forge/skills/init-justfile/SKILL.md:75,250-251`: recipe 签名声明为无参数 (ref: Proposed Solution)
- `plugins/forge/skills/init-justfile/rules/surfaces/cli.md:22,39-66`: recipe 模板无参数 (ref: Proposed Solution)
- `plugins/forge/skills/run-tests/SKILL.md:191`: 调用 `just <surface>-test <journey>` (ref: Proposed Solution)
- `plugins/forge/skills/run-tests/rules/surfaces/cli.md:46`: `just cli-test <journey>` (ref: Proposed Solution)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/init-justfile/SKILL.md` | recipe 签名从 `just <key>-test` 改为 `just <key>-test [journey]`，支持可选 journey 参数 |
| `plugins/forge/skills/init-justfile/rules/surfaces/cli.md` | recipe 模板增加可选 journey 参数 |
| `plugins/forge/skills/init-justfile/rules/surfaces/api.md` | recipe 模板增加可选 journey 参数 |
| `plugins/forge/skills/init-justfile/rules/surfaces/web.md` | recipe 模板增加可选 journey 参数 |
| `plugins/forge/skills/init-justfile/rules/surfaces/tui.md` | recipe 模板增加可选 journey 参数 |
| `plugins/forge/skills/init-justfile/rules/surfaces/mobile.md` | recipe 模板增加可选 journey 参数 |
| `plugins/forge/skills/run-tests/SKILL.md` | 多 surface 项目使用 surface-key 构造 recipe 名称；单 surface 项目使用 surface-type |
| `plugins/forge/skills/run-tests/rules/surfaces/cli.md` | recipe 调用支持 journey 参数 |
| `plugins/forge/skills/run-tests/rules/surfaces/api.md` | recipe 调用支持 journey 参数 |
| `plugins/forge/skills/run-tests/rules/surfaces/web.md` | recipe 调用支持 journey 参数 |
| `plugins/forge/skills/run-tests/rules/surfaces/tui.md` | recipe 调用支持 journey 参数 |
| `plugins/forge/skills/run-tests/rules/surfaces/mobile.md` | recipe 调用支持 journey 参数 |

## Acceptance Criteria
- [ ] init-justfile SKILL.md 声明 recipe 签名为 `just <key>-test [journey]`（可选 journey 参数）
- [ ] init-justfile 5 个 surface rule 文件的 recipe 模板接受可选 journey 参数
- [ ] run-tests SKILL.md 在多 surface 项目中使用 `forge surfaces --json` 的 `key` 字段构造 recipe 名称
- [ ] run-tests 5 个 surface rule 文件的 per-journey 执行伪代码与 init-justfile recipe 签名对齐
- [ ] 单 surface 项目中 run-tests 仍使用 surface-type 作为 recipe 前缀（向后兼容）

## Hard Rules
- 必须先加载 `docs/conventions/forge-distribution.md`
- recipe 参数签名使用 just 的可选参数语法 `[journey]`，不是必选参数

## Implementation Notes
- just 可选参数语法：`recipe-name [param]` — 方括号表示可选。当不传 journey 时，recipe 运行全部测试；传入时只运行指定 journey
- run-tests 获取 surface-key 的方式：`forge surfaces --json` 输出的 `key` 字段。单 surface 项目 key 为 "."，此时使用 surface-type
