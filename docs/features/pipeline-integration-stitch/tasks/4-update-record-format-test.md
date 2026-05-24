---
id: "4"
title: "更新 record-format-test.md：替换废弃类型列表"
priority: "P1"
estimated_time: "15min"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 4: 更新 record-format-test.md

## Description

`record-format-test.md` 列出了已废弃/不存在的类型（`test.gen-cases`、`test.eval-cases`、`test.gen-and-run`），缺少新类型。需替换为当前有效类型列表。

## Reference Files
- `proposal.md#P1-—-Record-模板参考文档更新` — 精确的替换类型列表
- `proposal.md#Success-Criteria` — record-format-test.md 的验证标准

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/submit-task/data/record-format-test.md` | 替换类型列表 |

## Acceptance Criteria

- [ ] `record-format-test.md` 包含全部五个有效类型：`test.gen-journeys`、`test.gen-contracts`、`test.gen-scripts`、`test.run`、`test.verify-regression`
- [ ] `record-format-test.md` 不包含已废弃类型名：`test.gen-cases`、`test.eval-cases`、`test.gen-and-run`

## Hard Rules

- 仅修改类型列表，不改动文件其他内容

## Implementation Notes

当前 stale 列表：`test.gen-cases`, `test.eval-cases`, `test.gen-scripts`, `test.run`, `test.gen-and-run`, `test.verify-regression`

正确列表：`test.gen-journeys`, `test.gen-contracts`, `test.gen-scripts`, `test.run`, `test.verify-regression`
