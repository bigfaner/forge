---
id: "3"
title: "更新 gen-test-scripts 和 gen-journeys skill 文件术语"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 3: 更新 gen-test-scripts 和 gen-journeys skill 文件术语

## Description

更新 gen-test-scripts 和 gen-journeys 两个 skill 的 SKILL.md 及规则文件中的测试类型术语。这两个 skill 负责测试生成流程的前半段（Journey 提取 → 测试脚本生成），需要确保生成的测试代码使用 surface-specific 的测试类型名称和标签。

## Reference Files
- `docs/proposals/surface-test-type-model/proposal.md#Proposed-Solution` — Test Type Mapping 表定义了每种 surface 的测试类型名称和验证维度
- `docs/proposals/surface-test-type-model/proposal.md#Technical-Direction` — task type 三段式命名规则
- `docs/proposals/surface-test-type-model/proposal.md#Success-Criteria` — SC4（skill 文件使用 surface-specific 名称）和 SC5（rules 文件引用概念文档）

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | 更新 "e2e" / "E2E" 引用为 surface-specific 测试类型名称 |
| `plugins/forge/skills/gen-test-scripts/types/cli.md` | 引用概念文档，使用 "CLI 功能测试" 术语 |
| `plugins/forge/skills/gen-test-scripts/types/tui.md` | 引用概念文档，使用 "终端功能测试" 术语 |
| `plugins/forge/skills/gen-test-scripts/types/api.md` | 引用概念文档，使用 "API 功能测试" 术语 |
| `plugins/forge/skills/gen-test-scripts/types/ui.md` | 引用概念文档，使用 "Web 端到端测试" 术语 |
| `plugins/forge/skills/gen-test-scripts/types/mobile.md` | 引用概念文档，使用 "移动端端到端测试" 术语 |
| `plugins/forge/skills/gen-test-scripts/types/_shared.md` | 更新共享术语定义 |
| `plugins/forge/skills/gen-journeys/rules/surface-cli.md` | 使用 "CLI 功能测试" 术语 |
| `plugins/forge/skills/gen-journeys/rules/surface-tui.md` | 使用 "终端功能测试" 术语 |
| `plugins/forge/skills/gen-journeys/rules/surface-api.md` | 使用 "API 功能测试" 术语 |
| `plugins/forge/skills/gen-journeys/rules/surface-web.md` | 使用 "Web 端到端测试" 术语 |
| `plugins/forge/skills/gen-journeys/rules/surface-mobile.md` | 使用 "移动端端到端测试" 术语 |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria

- [ ] gen-test-scripts/SKILL.md 中不再使用 "e2e" 作为统称，改为引用 `docs/reference/test-type-model.md`
- [ ] gen-test-scripts/types/ 5 个文件各自使用对应 surface 的测试类型名称和语义定义
- [ ] gen-journeys/rules/surface-*.md 5 个文件使用对应的测试类型名称
- [ ] 所有 rules 文件引用概念文档 `docs/reference/test-type-model.md` 中的测试类型定义
- [ ] 生成的测试代码注释/标签中使用 surface-specific 测试类型名称（如 `@cli-functional` 而非 `@e2e`）

## Hard Rules

- types/ 文件中的生成策略逻辑不变，仅更新术语和标签
- 保持每个 types/ 文件的 frontmatter `domains` 字段与 `testing` 关键词关联

## Implementation Notes

gen-test-scripts/SKILL.md 第 210 行有 `task_lifecycle_smoke_test.go <- Journey smoke test (happy path E2E)` 注释，需更新。第 214 行的硬规则 "NOT to tests/e2e/features/" 保留（这是旧路径限制，与术语无关）。gen-journeys 的 surface rules 文件主要负责 Journey 生成的风险分类和标签规则，需要确保 @feature 标签与新的测试类型名称一致。
