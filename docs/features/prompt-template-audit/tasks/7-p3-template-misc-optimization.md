---
id: "7"
title: "P3: coding-refactor fmt处理简化 + coding-enhancement同步注释"
priority: "P2"
estimated_time: "20m"
dependencies: ["6"]
type: "doc"
mainSession: false
---

# 7: P3: coding-refactor fmt处理简化 + coding-enhancement同步注释

## Description
两个低优先级优化：(1) coding-refactor.md Step 4 的 fmt 失败处理使用了复杂的 git stash 流程，简化为更直观的描述；(2) coding-enhancement.md 与 coding-feature.md 有 90% 内容重复，添加同步维护注释提醒修改时需同步检查。

## Reference Files
- `docs/proposals/prompt-template-audit/proposal.md` — Source proposal (P3 #14/#15)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `forge-cli/pkg/prompt/data/coding-refactor.md` | 简化 Step 4 fmt 失败处理：移除 git stash 复杂流程，改为简单的"如果 fmt 只影响 refactor 涉及的文件则修复，否则继续" |
| `forge-cli/pkg/prompt/data/coding-enhancement.md` | 在文件头部添加同步维护注释：修改此文件时同步检查 coding-feature.md |

## Acceptance Criteria
- [ ] coding-refactor.md 不再包含 `git stash && just fmt && git diff --name-only && git stash pop` 等复杂 git 操作
- [ ] coding-enhancement.md 头部包含同步维护注释
- [ ] 简化后的 fmt 处理仍然覆盖关键场景（fmt 失败时的处理策略）

## Implementation Notes
- coding-refactor.md 的 batch sizing 硬编码数值（≤10/15-20/3-5）维持现状——模板无法获取项目规模信息
