---
id: "2"
title: "tech-design 增加全 surface handbook 自动生成"
priority: "P0"
estimated_time: "2h"
dependencies: []
type: "doc"
mainSession: false
---

# 2: tech-design 增加全 surface handbook 自动生成

## Description
扩展 tech-design skill，支持为 CLI/TUI surface 生成 cli-handbook、Web surface 生成 page-map、Mobile surface 生成 screen-map。以现有 api-handbook 为参考模式，确保格式一致。这些 handbook 是后续 gen-contracts 填充锚点字段的权威数据源。

## Reference Files
- `docs/proposals/contract-technical-anchors/proposal.md` — Proposed Solution, Scope, Phased Implementation Roadmap
- `plugins/forge/skills/tech-design/SKILL.md`: 增加 handbook 生成分支逻辑 (ref: Proposed Solution)
- `plugins/forge/skills/tech-design/templates/api-handbook.md`: 参考格式模式 (ref: Feasibility Assessment)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/tech-design/templates/cli-handbook.md` | CLI/TUI handbook 模板 |
| `plugins/forge/skills/tech-design/templates/page-map.md` | Web page-map 模板 |
| `plugins/forge/skills/tech-design/templates/screen-map.md` | Mobile screen-map 模板 |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/tech-design/SKILL.md` | 增加 handbook 生成分支，根据 surface 类型选择模板 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] tech-design skill 根据项目 surface 配置（`forge surfaces`）自动选择并生成对应 handbook
- [ ] cli-handbook 模板覆盖 CLI/TUI surface 的命令、子命令、参数、别名等信息
- [ ] page-map 模板覆盖 Web surface 的页面名称、路由、布局、认证要求
- [ ] screen-map 模板覆盖 Mobile surface 的屏幕名称、导航路径、deeplink、平台
- [ ] 新 handbook 模板格式与 api-handbook 保持一致的 frontmatter 和 section 结构

## Implementation Notes
- 参照 api-handbook 的成熟格式模式，新增类型尽量对齐
- 分批实现建议：Phase 1 CLI（与 API 同步）→ Phase 2 Web → Phase 3 Mobile
- handbook 文件需包含生成时间戳，供后续新鲜度检查使用
