---
id: "15"
title: "Unify gen-journeys test type format + fix gen-contracts stale paths"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 15: Unify gen-journeys test type format + fix gen-contracts stale paths

## Description
gen-journeys 5 个 surface rule文件的第 5 行 test type 描述存在双格式混写："CLI 功能测试 (CLI Functional Test). Test type: CLI 功能测试，通过子进程调用验证..."，同一概念重复描述两次。需统一为单行简洁格式。同时修复 gen-contracts 和 gen-journeys 中残留的过时路径引用。

## Reference Files
- `docs/proposals/surface-first-testing/proposal.md` — Proposed Solution
- `plugins/forge/skills/gen-journeys/rules/surface-cli.md:5`: 双格式混写
- `plugins/forge/skills/gen-journeys/rules/surface-web.md:5`: 双格式混写
- `plugins/forge/skills/gen-contracts/rules/journey-contract-model.md:155`: 旧路径 `testing-<scope>.md`
- `plugins/forge/skills/gen-journeys/SKILL.md:305`: 引用不存在的 model-and-directory-spec.md

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-journeys/rules/surface-cli.md` | 第 5 行 test type 统一为 "**Test type**: CLI Functional Test. 通过子进程调用验证进程退出码、stdout/stderr 输出和参数校验行为。" |
| `plugins/forge/skills/gen-journeys/rules/surface-api.md` | 第 5 行统一格式 |
| `plugins/forge/skills/gen-journeys/rules/surface-tui.md` | 第 5 行统一格式 |
| `plugins/forge/skills/gen-journeys/rules/surface-web.md` | 第 5 行统一格式 |
| `plugins/forge/skills/gen-journeys/rules/surface-mobile.md` | 第 5 行统一格式 |
| `plugins/forge/skills/gen-contracts/rules/journey-contract-model.md` | 旧路径 `testing-<scope>.md` → `testing/<surface>/core.md` |
| `plugins/forge/skills/gen-journeys/SKILL.md` | HARD-RULE 移除对不存在文件 `model-and-directory-spec.md` 的引用 |

## Acceptance Criteria
- [ ] gen-journeys 5 个 surface rule 文件 test type 描述统一为: `**Test type**: {English name}. {一句话英文描述}.`
- [ ] 移除重复的中文 test type 描述（首段中文描述保留，但第 5 行不再重复中文）
- [ ] journey-contract-model.md 中 `testing-<scope>.md` 改为 `testing/<surface>/core.md`
- [ ] gen-journeys SKILL.md HARD-RULE 不再引用不存在的 `model-and-directory-spec.md`

## Hard Rules
- 不引用其它 skill 的内部文件
- test type 英文名必须与 test-type-model.md 一致

## Implementation Notes
- 统一格式示例: `**Test type**: CLI Functional Test. Verifies process exit codes, stdout/stderr output, and argument validation via subprocess invocation. Generated test code MUST use \`@cli-functional\` tags, NOT \`@e2e\`.`
- 每个文件首段（第 3 行）的中文描述保留不变
