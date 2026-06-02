---
id: "9"
title: "Fix stale test-type-model references in gen-journeys + run-tests + init-justfile"
priority: "P1"
estimated_time: "1h"
dependencies: []
type: "doc"
mainSession: false
---

# 9: Fix stale test-type-model references in gen-journeys + run-tests + init-justfile

## Description
修复 gen-journeys（5 个 rules/surface-*.md）、run-tests（5 个 rules/surfaces/*.md + result-parsing.md）、init-justfile（5 个 rules/surfaces/*.md）中共 16 处引用已删除的 `docs/reference/test-type-model.md` 的地方。全部改为自包含内联描述（路径 A）。每个 rule 文件已有完整的 surface-specific test type 信息，只需移除无效的外部引用。

同时修复 `run-tests/rules/result-parsing.md:62` 中 "UI tests" → "Web tests"。

## Reference Files
- `docs/proposals/surface-first-testing/proposal.md` — Proposed Solution, Success Criteria
- `plugins/forge/skills/gen-journeys/rules/surface-cli.md:5`: 旧路径引用 (ref: Proposed Solution)
- `plugins/forge/skills/gen-journeys/rules/surface-api.md:5`: 旧路径引用 (ref: Proposed Solution)
- `plugins/forge/skills/gen-journeys/rules/surface-tui.md:5`: 旧路径引用 (ref: Proposed Solution)
- `plugins/forge/skills/gen-journeys/rules/surface-web.md:5`: 旧路径引用 (ref: Proposed Solution)
- `plugins/forge/skills/gen-journeys/rules/surface-mobile.md:5`: 旧路径引用 (ref: Proposed Solution)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-journeys/rules/surface-cli.md` | 移除旧路径引用，内联 test type 定义 |
| `plugins/forge/skills/gen-journeys/rules/surface-api.md` | 移除旧路径引用，内联 test type 定义 |
| `plugins/forge/skills/gen-journeys/rules/surface-tui.md` | 移除旧路径引用，内联 test type 定义 |
| `plugins/forge/skills/gen-journeys/rules/surface-web.md` | 移除旧路径引用，内联 test type 定义 |
| `plugins/forge/skills/gen-journeys/rules/surface-mobile.md` | 移除旧路径引用，内联 test type 定义 |
| `plugins/forge/skills/run-tests/rules/surfaces/cli.md` | 移除旧路径引用 |
| `plugins/forge/skills/run-tests/rules/surfaces/api.md` | 移除旧路径引用 |
| `plugins/forge/skills/run-tests/rules/surfaces/tui.md` | 移除旧路径引用 |
| `plugins/forge/skills/run-tests/rules/surfaces/web.md` | 移除旧路径引用 |
| `plugins/forge/skills/run-tests/rules/surfaces/mobile.md` | 移除旧路径引用 |
| `plugins/forge/skills/run-tests/rules/result-parsing.md` | "UI tests" → "Web tests" |
| `plugins/forge/skills/init-justfile/rules/surfaces/cli.md` | 移除旧路径引用 |
| `plugins/forge/skills/init-justfile/rules/surfaces/api.md` | 移除旧路径引用 |
| `plugins/forge/skills/init-justfile/rules/surfaces/tui.md` | 移除旧路径引用 |
| `plugins/forge/skills/init-justfile/rules/surfaces/web.md` | 移除旧路径引用 |
| `plugins/forge/skills/init-justfile/rules/surfaces/mobile.md` | 移除旧路径引用 |

## Acceptance Criteria
- [ ] gen-journeys 5 个 surface rule 文件中 `docs/reference/test-type-model.md` 引用已移除，替换为自包含的 test type 声明
- [ ] run-tests 5 个 surface rule 文件中旧路径引用已移除
- [ ] run-tests `result-parsing.md` 中 "UI tests" 改为 "Web tests"
- [ ] init-justfile 5 个 surface rule 文件中旧路径引用已移除
- [ ] 所有修改后的文件不包含任何对其它 skill 内部文件的引用

## Hard Rules
- 不引用其它 skill 的内部文件，所有知识自包含
- 每个 rule 文件修改后必须仍能独立被 LLM agent 理解和使用

## Implementation Notes
- gen-journeys rules 中的引用格式为 "Test type definition and classification criteria: see `docs/reference/test-type-model.md`"，替换为自包含的一句话声明（如 "Test type: CLI 功能测试，通过子进程调用验证输入-输出行为。"）
- run-tests/init-justfile rules 中的引用格式为 "测试类型术语定义参见 `docs/reference/test-type-model.md`"，直接删除此引用即可（文件已有完整的 surface-specific 信息）
