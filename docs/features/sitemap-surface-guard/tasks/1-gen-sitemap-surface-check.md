---
id: "1"
title: "gen-web-sitemap: Add Step 0 surface type check"
priority: "P0"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
complexity: "medium"
---

# 1: gen-web-sitemap: Add Step 0 surface type check

## Description

`gen-web-sitemap` SKILL.md 缺少 surface 类型校验，可在任何项目中被直接调用。需在 Process Flow 的 Step 1 之前新增 Step 0，通过 `forge surfaces --json` 检测项目是否有 `web` surface。无 web surface 时中止执行并输出明确提示。

## Reference Files
- `plugins/forge/skills/gen-web-sitemap/SKILL.md`: Process Flow 当前从 Step 1 开始，需在 Step 1 前插入 Step 0 Surface Check (source: proposal.md#Proposed-Solution)
- `plugins/forge/skills/gen-test-scripts/types/ui.md`: Sitemap Resolution 段落有已有守卫模式可作参考: "Only execute when the project has `web-ui` interface AND UI-type test cases exist." (source: proposal.md#Proposed-Solution)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-web-sitemap/SKILL.md` | 在 Process Flow 的 Step 1 前新增 Step 0 Surface Check |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] Step 0 执行 `forge surfaces --json`，解析返回的 surface 列表
- [ ] 无 `web` surface 类型时，STOP 并输出明确提示（如 "No web surface detected. gen-web-sitemap is only applicable to web projects."）
- [ ] 有 `web` surface（含 monorepo 多 surface 场景）时，Step 0 通过，正常进入 Step 1
- [ ] `forge surfaces --json` 返回空或命令失败时，等同于无 web surface，中止执行

## Implementation Notes

守卫指令放在 Process Flow 的最早位置（Step 0），配合 STOP 关键词增强 LLM agent 可见性。参考 `gen-test-scripts/types/ui.md` 已有守卫的措辞风格保持一致。
