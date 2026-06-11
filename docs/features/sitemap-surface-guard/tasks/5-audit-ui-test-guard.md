---
id: "5"
title: "gen-test-scripts/types/ui: Audit existing surface guard for consistency"
priority: "P2"
estimated_time: "0.5h"
dependencies: [2, 3, 4]
type: "doc"
mainSession: false
complexity: "low"
---

# 5: gen-test-scripts/types/ui: Audit existing surface guard for consistency

## Description

gen-test-scripts/types/ui.md 已有守卫 "Only execute when the project has `web-ui` interface AND UI-type test cases exist"，作为本次 surface guard 加固的参考基准。需审查该守卫在新的 surface 检测体系下是否仍然充分，确认是否需要调整措辞或逻辑。

## Reference Files
- `plugins/forge/skills/gen-test-scripts/types/ui.md`: Sitemap Resolution 段落已有守卫 "Only execute when the project has `web-ui` interface AND UI-type test cases exist."，需审查是否与 Task 1-4 统一的 `forge surfaces --json` 守卫模式一致 (source: proposal.md#Scope-In-Scope)

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-scripts/types/ui.md` | 审查并可能更新 Sitemap Resolution 段落的守卫措辞 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] 审查已有守卫 "Only execute when the project has `web-ui` interface" 是否与 `forge surfaces --json` 检测方式一致
- [ ] 记录审查结论：守卫充分无需修改，或实施修改并记录变更内容
- [ ] 若修改，守卫措辞与 Task 1-4 中建立的模式保持一致

## Implementation Notes

已有守卫使用 "web-ui interface" 措辞，而新守卫使用 `forge surfaces --json` 返回的 `web` surface type。需确认 "web-ui interface" 与 surface type `web` 的对应关系是否清晰，是否需要统一措辞。
