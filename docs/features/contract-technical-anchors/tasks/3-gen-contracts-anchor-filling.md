---
id: "3"
title: "gen-contracts 从 handbook 填充锚点字段"
priority: "P0"
estimated_time: "2h"
dependencies: [1, 2]
type: "doc"
mainSession: false
---

# 3: gen-contracts 从 handbook 填充锚点字段

## Description
gen-contracts 在生成 Contract 时读取对应 surface 的 handbook，自动填充技术锚点字段。包含 handbook 新鲜度检查（比对 handbook 生成时间戳与 tech-design 最后修改时间），过期时提示用户重新生成。缺少 handbook 时跳过锚点填充并提示用户。

## Reference Files
- `docs/proposals/contract-technical-anchors/proposal.md` — Proposed Solution, Key Scenarios
- `plugins/forge/skills/gen-contracts/SKILL.md`: 增加 handbook 读取和锚点填充步骤 (ref: Proposed Solution)
- `plugins/forge/skills/gen-contracts/templates/contract.md`: 锚点字段定义 (ref: Anchor Field Schema)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-contracts/SKILL.md` | 增加 handbook 读取、锚点填充、新鲜度检查逻辑 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] gen-contracts 读取 api-handbook 自动填充 API Contract 的 endpoint/method 字段
- [ ] gen-contracts 读取 cli-handbook 自动填充 CLI/TUI Contract 的 command/subcommand 字段
- [ ] gen-contracts 读取 page-map/screen-map 自动填充 Web/Mobile Contract 的 page/screen 字段
- [ ] Handbook 新鲜度检查：handbook 生成时间早于 tech-design 最后修改时间时，提示用户"handbook 可能过期，建议重新生成"
- [ ] 缺少 handbook 时跳过锚点填充并提示用户"缺少 handbook，建议运行 tech-design 生成"
- [ ] `last_anchor_sync` 时间戳在填充时自动更新

## Implementation Notes
- handbook 读取在 gen-contracts 现有步骤之后、最终输出之前执行
- 向后兼容：缺少 handbook 或锚点字段时管道不中断，降级为现有行为
- 填充逻辑以 handbook 为权威源，不从代码中逆向提取
