---
id: "14"
title: "Consolidate-specs extraction, UTF-8 handling, and CLI test path updates"
priority: "P2"
estimated_time: "1h"
dependencies: ["6"]
type: "doc"
mainSession: false
---

# 14: Consolidate-specs extraction, UTF-8 handling, and CLI test path updates

## Description
三个小项：(1) consolidate-specs SKILL.md 提取 ≥50 行到 rules/ 以降低主文件体积；(2) guide.md UTF-8 字符处理评估；(3) CLI 测试中过时 skill 路径引用更新。

## Reference Files
- `proposal.md#Scope` — P2.18: consolidate-specs extraction; P2.19: guide.md UTF-8; P2.20: CLI test path updates

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/consolidate-specs/SKILL.md` | Extract ≥50 lines to rules/ |
| `docs/conventions/guide.md` or `docs/reference/guide.md` | UTF-8 character handling assessment |
| CLI test files | Update stale skill path references |

## Acceptance Criteria
- [ ] consolidate-specs SKILL.md 体积减小 ≥50 行
- [ ] UTF-8 字符处理问题已评估，影响已记录
- [ ] CLI 测试中无过时 skill 路径引用

## Hard Rules
- consolidate-specs 提取需保持 SKILL.md 流程完整
- CLI 测试路径更新需与实际 skill 目录匹配

## Implementation Notes
UTF-8 评估可能是无操作（评估后发现无需修改）。
