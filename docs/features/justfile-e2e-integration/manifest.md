---
feature: "justfile-e2e-integration"
status: tasks
---

# Feature: justfile-e2e-integration

<!-- Status flow: prd → design → tasks → in-progress → done -->

## Documents

| Document | Path | Summary |
|----------|------|---------|
| PRD Spec | prd/prd-spec.md | 扩展 init-justfile 新增 e2e-setup/e2e-verify 目标，将 13 个 skill/agent 文件中的原始命令（代码块+文字描述）统一替换为 just 命令 |
| User Stories | prd/prd-user-stories.md | 5 个故事：Skill 维护者使用统一接口、Agent 获得明确指令、VERIFY 硬门控、fix-e2e 验证修复、构建/测试统一命令 |
| Tech Design | design/tech-design.md | 4 阶段实现计划：Phase 1 新增 just 目标、Phase 2 e2e skill 文件、Phase 3 构建/测试文件（7个）、Phase 4 breakdown-tasks 模板 |

## Traceability

| PRD Section | Design Section | UI Component | Tasks |
|-------------|----------------|--------------|-------|
| 5.1 新增 just 目标 | Phase 1: init-justfile | — | — |
| 5.2 E2E 命令替换 | Phase 2: gen-test-scripts, run-e2e-tests | — | — |
| 5.2 单元测试/构建命令替换 | Phase 3: fix-bug, run-tasks, task-executor, error-fixer, execute-task, record-task, improve-harness | — | — |
| 5.2 breakdown-tasks 模板 | Phase 4: run-e2e-tests.md, gen-test-scripts.md, fix-e2e.md | — | — |
