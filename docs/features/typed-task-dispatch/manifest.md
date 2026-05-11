---
feature: "typed-task-dispatch"
status: tasks
---

# Feature: typed-task-dispatch

<!-- Status flow: prd → design → tasks → in-progress → completed -->

## Documents

| Document | Path | Summary |
|----------|------|---------|
| PRD Spec | prd/prd-spec.md | 将策略逻辑从 agent markdown 下沉到 CLI，通过 task prompt 命令合成类型专属 prompt，task-executor 变为薄执行器；引入显式 type 字段替代 noTest 补丁字段；提供 task migrate 迁移命令 |
| User Stories | prd/prd-user-stories.md | 7 个用户故事：非编码任务获得正确流程、新增类型只需一个模板文件、独立检查合成结果、迁移旧 index.json、breakdown-tasks 自动设置 type、execute-task 路由一致性、error-fixer 废弃等价覆盖 |
| Tech Design | design/tech-design.md | 两层模型（约束层 + 策略层）；新增 pkg/prompt 包（11 种类型 Go embed 模板）；task prompt / task migrate 命令；task-executor.md 精简为 ~40 行；run-tasks / execute-task 路由更新；breakdown-tasks / quick-tasks Type Assignment 规则 |

## Traceability

| PRD Section | Design Section | UI Component | Placement | Tasks |
|-------------|----------------|--------------|-----------|-------|
| Flow: task prompt 内部流程 | New Package: pkg/prompt | — | — | 1.2 |
| Functional Specs: task prompt 命令 | New Commands: task prompt | — | — | 1.3 |
| Functional Specs: task migrate 命令 | New Commands: task migrate | — | — | 1.4 |
| Related Changes: task validate 扩展 | New Commands: task validate extension | — | — | 1.5 |
| Related Changes: task claim TYPE 输出 | Agent Changes: run-tasks.md | — | — | 1.6 |
| Related Changes: index.schema.json | Data Models: index.schema.json | — | — | 2.1 |
| Related Changes: 所有任务模板 frontmatter | Skill Layer Changes: 任务模板 | — | — | 2.2 |
| Related Changes: breakdown-tasks / quick-tasks | Skill Layer Changes | — | — | 2.3 |
| Related Changes: task-executor.md 精简 | Agent Changes: task-executor.md | — | — | 3.1 |
| Related Changes: run-tasks 路由变更 | Agent Changes: run-tasks.md | — | — | 3.2 |
| Related Changes: execute-task 路由变更 | Agent Changes: execute-task.md | — | — | 3.3 |
| Related Changes: error-fixer 废弃 | Agent Changes: error-fixer deprecation | — | — | 4.1 |
