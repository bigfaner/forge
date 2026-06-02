---
id: "16"
title: "Fix gen-contracts concept attribution + run-tests cross-skill reference"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 16: Fix gen-contracts concept attribution + run-tests cross-skill reference

## Description
gen-contracts SKILL.md Steps 3.7-3.10 将 State Verification Levels、Journey Invariants、Batch Processing 等通用概念全部归入 `rules/tui-async.md`，但这些概念对所有 surface 都适用。需要在 SKILL.md 中明确标注通用性。同时 run-tests `rules/env-check.md` 跨 skill 引用 gen-journeys 的 `rules/surface-<type>.md`，属于脆弱依赖，需改为自包含描述。

## Reference Files
- `docs/proposals/surface-first-testing/proposal.md` — Proposed Solution
- `plugins/forge/skills/gen-contracts/SKILL.md:217-228`: Steps 3.7-3.10 通用概念归入 TUI rule (ref: Proposed Solution)
- `plugins/forge/skills/run-tests/rules/env-check.md:10-20`: 跨 skill 引用 gen-journeys surface rule
- `plugins/forge/skills/gen-contracts/rules/tui-async.md`: 当前 TUI 专属 rule 文件

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-contracts/SKILL.md` | Steps 3.8-3.10 标注 State Verification / Journey Invariants / Batch Processing 为通用概念（适用所有 surface），仅 3.7 TUI Async 为 TUI 专属 |
| `plugins/forge/skills/run-tests/rules/env-check.md` | 跨 skill 引用改为自包含：内联各 surface 的环境检查清单或引用自身 surface rule 文件 |

## Acceptance Criteria
- [ ] gen-contracts SKILL.md Step 3.7 明确标注为 "TUI-specific"
- [ ] gen-contracts SKILL.md Steps 3.8/3.9/3.10 明确标注为 "applies to all surface types"
- [ ] run-tests env-check.md 不再引用 gen-journeys skill 的内部文件路径
- [ ] env-check.md 改为引用 run-tests 自身的 `rules/surfaces/<type>.md` 或自包含描述

## Hard Rules
- 不引用其它 skill 的内部文件
- 必须先加载 `docs/conventions/forge-distribution.md`

## Implementation Notes
- gen-contracts 的 tui-async.md 包含了通用验证逻辑和 TUI 专属 async 逻辑的混合，短期内通过标注区分即可，长期可考虑拆分
- env-check.md 的 Environment Readiness Checks 信息来源可从 gen-journeys 的 surface rule 文件中提取关键检查项，内联到 env-check.md 或各自 surface rule 文件中
