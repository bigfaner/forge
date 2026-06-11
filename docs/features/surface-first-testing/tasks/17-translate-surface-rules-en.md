---
id: "17"
title: "Translate init-justfile + run-tests + gen-contracts surface rules to English"
priority: "P2"
estimated_time: "2h"
dependencies: []
type: "doc"
mainSession: false
---

# 17: Translate init-justfile + run-tests + gen-contracts surface rules to English

## Description
init-justfile 和 run-tests 各有 5 个 surface rule 文件（rules/surfaces/*.md）使用中文编写，gen-contracts 的 journey-contract-model.md 也是中文。而 gen-journeys 的 surface rule 文件以英文为主，形成不一致。所有面向 LLM agent 的规则文件应统一为英文，提高跨语言项目的通用性。

## Reference Files
- `docs/proposals/surface-first-testing/proposal.md` — Non-Functional Requirements
- `plugins/forge/skills/init-justfile/rules/surfaces/cli.md`: 全文中文
- `plugins/forge/skills/run-tests/rules/surfaces/cli.md`: 全文中文
- `plugins/forge/skills/gen-contracts/rules/journey-contract-model.md`: 全文中文

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/init-justfile/rules/surfaces/cli.md` | 中文 → 英文 |
| `plugins/forge/skills/init-justfile/rules/surfaces/api.md` | 中文 → 英文 |
| `plugins/forge/skills/init-justfile/rules/surfaces/web.md` | 中文 → 英文 |
| `plugins/forge/skills/init-justfile/rules/surfaces/tui.md` | 中文 → 英文 |
| `plugins/forge/skills/init-justfile/rules/surfaces/mobile.md` | 中文 → 英文 |
| `plugins/forge/skills/run-tests/rules/surfaces/cli.md` | 中文 → 英文 |
| `plugins/forge/skills/run-tests/rules/surfaces/api.md` | 中文 → 英文 |
| `plugins/forge/skills/run-tests/rules/surfaces/web.md` | 中文 → 英文 |
| `plugins/forge/skills/run-tests/rules/surfaces/tui.md` | 中文 → 英文 |
| `plugins/forge/skills/run-tests/rules/surfaces/mobile.md` | 中文 → 英文 |
| `plugins/forge/skills/gen-contracts/rules/journey-contract-model.md` | 中文 → 英文 |

## Acceptance Criteria
- [ ] init-justfile 5 个 surface rule 文件翻译为英文，逻辑内容不变
- [ ] run-tests 5 个 surface rule 文件翻译为英文，逻辑内容不变
- [ ] gen-contracts journey-contract-model.md 翻译为英文，逻辑内容不变
- [ ] 翻译后文件不包含中文（技术术语如 `tests/<journey>/` 路径除外）
- [ ] surface type 术语保持小写：web/api/cli/tui/mobile

## Hard Rules
- 翻译不改变任何逻辑、规则、约束
- 必须先加载 `docs/conventions/forge-distribution.md`

## Implementation Notes
- 优先保持技术精确性，避免意译导致歧义
- 英文格式参考 gen-journeys 的 surface rule 文件作为模板
