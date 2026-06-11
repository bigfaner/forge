---
status: "completed"
started: "2026-05-20 17:07"
completed: "2026-05-20 17:10"
time_spent: "~3m"
---

# Task Record: 1 替换所有 skill 文件中的过期 CLI 命令引用

## Summary
替换所有 skill 文件中已移除的 CLI 命令引用（forge test detect, forge test interfaces, forge task list, forge deploy），替换为通过检查项目文件推断语言/接口类型的明确指令

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-contracts/SKILL.md
- plugins/forge/skills/gen-test-cases/SKILL.md
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/tech-design/SKILL.md
- plugins/forge/skills/quick-tasks/SKILL.md
- plugins/forge/skills/eval/rubrics/validate-ux-pipeline.md
- plugins/forge/skills/eval/rules/pre-processing.md
- plugins/forge/skills/eval/rubrics/test-cases.md
- plugins/forge/skills/gen-test-cases/types/cli.md
- plugins/forge/skills/eval/rubrics/cli-test-cases.md

### Key Decisions
- 语言检测替换为检查 package.json/go.mod/Cargo.toml/pyproject.toml/setup.py 文件，覆盖 JS/TS/Go/Rust/Python 四种主流语言
- 接口类型检测替换为检查 docs/conventions/、项目目录结构、.forge/config.yaml、package.json 依赖，覆盖 UI/TUI/API/CLI/Mobile 五种类型
- forge deploy 示例替换为 forge task query --status staging（实际存在的命令）
- forge task list 替换为 forge task query

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] plugins/forge/skills/ 下所有文件中零引用 forge test detect
- [x] plugins/forge/skills/ 下所有文件中零引用 forge test interfaces
- [x] forge task list 引用替换为 forge task query
- [x] forge deploy 示例替换为实际存在的命令示例
- [x] 替换后的指令明确描述 agent 应检查哪些项目文件来推断信息

## Notes
共修改 10 个文件，替换了 16 处过期命令引用。所有替换文本包含具体的文件检查清单和推断逻辑，agent 可直接执行。
