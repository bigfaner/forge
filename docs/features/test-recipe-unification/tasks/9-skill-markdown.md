---
id: "9"
title: "Update skill and command markdown for new recipe model"
priority: "P1"
estimated_time: "1h"
dependencies: ["1", "7"]
type: "doc"
mainSession: false
---

# 9: Update skill and command markdown for new recipe model

## Description

更新 5+ 个 skill/command markdown 文件中的 recipe 引用和配置示例，使其反映新的两层测试模型。重点更新 init-justfile（Standard Target Contract）和 run-tests（config-schema 示例）。

## Reference Files
- `proposal.md#Impact-Analysis` — Tier 2 lists skill/command markdown files with specific change descriptions
- `proposal.md#Proposed-Solution` — defines recipe naming conventions and retires e2e-test/e2e-setup/e2e-verify
- `proposal.md#Requirements-Analysis` — Recipe 参数签名约定 for init-justfile generation rules
- `proposal.md#Key-Risks` — discoverability concern for `auto.test` naming and new recipe semantics

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `commands/fix-bug.md` | `just test` → `just unit-test`; `just e2e-test` → `just test` |
| `skills/clean-code/SKILL.md` | `just test` → `just unit-test` |
| `skills/gen-test-scripts/rules/run-to-learn.md` | `just test` → `just unit-test` |
| `skills/init-justfile/SKILL.md` | Standard Target Contract 全面更新：unit-test, test, test-setup, probe |
| `skills/run-tests/SKILL.md` | config-schema 示例中 recipe 名更新 |
| `skills/run-tests/references/config-schema.md` | `e2eTest` → `test` in schema examples |

## Acceptance Criteria
- All skill/command markdown references `unit-test`, `test`, `test-setup`, `probe` recipe names
- No residual `e2e-test`, `e2e-setup`, `e2e-verify` references
- `init-justfile/SKILL.md` Standard Target Contract reflects new recipe model with per-language/per-surface generation
- `run-tests` config schema examples use `test` key (not `e2eTest`)

## Implementation Notes
- init-justfile SKILL.md 的 Standard Target Contract 是最重要的更新——它定义了生成的 justfile 结构
- 注意 SKILL.md 是分发到用户环境的，修改前须加载 `docs/conventions/forge-distribution.md` 了解分发模型
