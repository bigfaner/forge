---
status: "completed"
started: "2026-05-20 13:57"
completed: "2026-05-20 14:16"
time_spent: "~19m"
---

# Task Record: 4 Slim eval/quality domain (eval + gen-contracts + test-guide)

## Summary
对评测/质量域的 3 个 skill 进行精简拆分：eval (372→200行)、gen-contracts (365→186行)、test-guide (380→186行)，抽取规则细节到 rules/ 子目录

## Changes

### Files Created
- plugins/forge/skills/eval/rules/rubric-context.md
- plugins/forge/skills/eval/rules/pre-processing.md
- plugins/forge/skills/eval/rules/scorer-composition.md
- plugins/forge/skills/eval/rules/reviser-composition.md
- plugins/forge/skills/eval/rules/report-format.md
- plugins/forge/skills/eval/rules/rubric-reference.md
- plugins/forge/skills/gen-contracts/rules/code-reconnaissance.md
- plugins/forge/skills/gen-contracts/rules/dimension-rules.md
- plugins/forge/skills/gen-contracts/rules/tui-async.md
- plugins/forge/skills/gen-contracts/rules/validation.md
- plugins/forge/skills/test-guide/rules/signal-detection.md
- plugins/forge/skills/test-guide/rules/pattern-extraction.md
- plugins/forge/skills/test-guide/rules/convention-structure.md

### Files Modified
- plugins/forge/skills/eval/SKILL.md
- plugins/forge/skills/gen-contracts/SKILL.md
- plugins/forge/skills/test-guide/SKILL.md

### Key Decisions
- eval: 将重复的 context injection 模板（Step 2.1 和 Step 4.1 各一份）抽取到 scorer-composition.md 和 reviser-composition.md，消除重复
- eval: 将 Expert Dispatch Table、Scorer prompt composition、report path tables、multi-expert merge logic 合并到 scorer-composition.md
- eval: Pre-processing by Type 表格、Rubric Reference 表格、Final Report 模板分别抽取到独立 rules 文件
- gen-contracts: 将六维度声明规则、语义描述符规则、前置条件互斥性合并为 dimension-rules.md
- gen-contracts: 将 TUI Async、State Verification、Journey Invariants、Batch Processing 合并为 tui-async.md
- gen-contracts: 将 Code Reconnaissance 表格抽取为 code-reconnaissance.md
- gen-contracts: 将 Validation checklist 和 Error Handling 表格合并为 validation.md
- test-guide: 将语言检测/框架检测信号表抽取为 signal-detection.md
- test-guide: 将测试文件模式/断言库检测/命名模式抽取为 pattern-extraction.md
- test-guide: 将 Convention 文件结构和冷启动框架候选表抽取为 convention-structure.md
- 与 Tier 1 (consolidate-specs、tech-design、write-prd) 保持一致的拆分风格：rules/ 子目录 + ${CLAUDE_SKILL_DIR} 引用路径

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] 每个 SKILL.md 行数 ≤ 350 行
- [x] 所有步骤编号及描述保留
- [x] 引用的辅助文件路径均存在可读
- [x] 拆分风格与 Tier 1 保持一致

## Notes
eval 已有 28 个辅助文件（experts/、rubrics/），本次新增 6 个 rules/ 文件，现有结构未受影响。gen-contracts 有 2 个辅助文件（25 行），新增 4 个 rules/ 文件。test-guide 有 1 个辅助文件（39 行），新增 3 个 rules/ 文件。SKILL.md 总行数从 1117 降至 572（减少 49%）。
