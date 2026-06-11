---
id: "2"
title: "撰写 README.en.md（英文版）"
priority: "P1"
estimated_time: "1h"
dependencies: [1]
type: "doc"
mainSession: false
---

# 2: 撰写 README.en.md（英文版）

## Description
基于中文版 README.md 撰写英文版 README.en.md，内容对应但非直译，需符合英文母语表达习惯。采用相同的 7 节结构，确保海外开发者能快速理解 Forge 的差异化价值。

## Reference Files
- `docs/proposals/readme-showcase-rewrite/proposal.md` — Proposed Solution, Non-Functional Requirements, Success Criteria, Key Risks
- `README.md` — 中文版内容参考（Task 1 完成后）

## Affected Files

### Create
| File | Description |
|------|-------------|
| `README.en.md` | 英文版 README，7 节结构 |

### Modify
| File | Changes |
|------|---------|
| _(无修改文件)_ | |

### Delete
| File | Reason |
|------|--------|
| _(无删除文件)_ | |

## Acceptance Criteria
- [ ] 英文版内容与中文版对应（非直译），符合英文母语表达习惯
- [ ] 包含与中文版相同的 7 节结构
- [ ] 总长度 ≤ 200 行
- [ ] 竞品功能矩阵与中文版一致（3 个竞品、≥ 6 个维度）
- [ ] 4 大核心特性各有独立小节

## Implementation Notes
- 中英文版本内容对应但非直译，各自符合母语表达习惯
- 痛点叙述需用英文开发者熟悉的表达方式，避免翻译腔
- 安装步骤的技术命令保持不变（与中文版一致）
