---
id: "1"
title: "Update existing command descriptions in guide.md"
priority: "P0"
estimated_time: "30m"
dependencies: []
type: "doc"
mainSession: false
---

# 1: Update existing command descriptions in guide.md

## Description
修复 guide.md 中 4 处 CLI 命令引用错误（G1-G4）。guide.md 作为系统提示注入所有 agent 会话，其中 `validate-index` 命令名错误已在 pm-work-tracker 项目中造成实际混淆，需立即修复。

## Reference Files
- `docs/proposals/cli-doc-accuracy-audit/proposal.md` — In Scope (G1-G4), Success Criteria, Evidence
- `plugins/forge/hooks/guide.md`: CLI 参考部分命令名和描述修正 (ref: In Scope)
- `docs/conventions/forge-distribution.md`: guide.md 分发规范，修改前必读 (ref: Constraints & Dependencies)

## Affected Files

### Create
| File | Description |
|------|-------------|
| _(none)_ |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/hooks/guide.md` | G1: validate-index → validate [file]; G2: quality-gate 描述; G3: cleanup 描述; G4: task submit --quiet |

### Delete
| File | Reason |
|------|--------|
| _(none)_ |

## Acceptance Criteria
- [ ] guide.md 中 `forge task validate-index` 替换为 `forge task validate [file]`，且 `forge task validate-index` 在 CLI 中返回 "unknown command" 错误
- [ ] guide.md 中 `forge quality-gate` 描述准确反映实际行为（含 fix task 自动创建、retry-once、docs-only 跳过）
- [ ] guide.md 中 `forge cleanup` 描述从 "clean stale artifacts" 改为具体行为说明（包含 blocked/suspended/rejected 状态的清理）
- [ ] guide.md 中 `forge task submit` 描述补充 `--quiet` 标志

## Hard Rules
- 修改 guide.md 前必须先阅读 `docs/conventions/forge-distribution.md`，遵循其分发规范

## Implementation Notes
- G1 是关键修复：`validate-index` 命令不存在于 CLI，必须替换为 `validate [file]`
- 所有描述修改需与对应命令的 `forge <command> --help` 输出核对一致
- Key Risk: 修改时避免引入新的不准确描述，逐条与 RunE 代码核对
