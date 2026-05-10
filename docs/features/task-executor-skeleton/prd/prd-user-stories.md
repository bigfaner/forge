---
feature: "task-executor-skeleton"
---

# User Stories: task-executor-skeleton

## Story 1: Template Author 编写 Execution Workflow

**As a** template author
**I want to** 在任务模板中声明 `## Execution Workflow` 段落指定执行步骤
**So that** task-executor 按我定义的流程执行任务，而非被硬编码为 TDD

**Acceptance Criteria:**

- Given 一个任务模板文件, When 它包含 `## Execution Workflow` 标题且正文非空, Then task-executor 将其正文作为 Step 2 执行指令（替换硬编码 TDD）
- Given 一个任务模板文件, When 它不包含 `## Execution Workflow` 标题, Then task-executor 回退到 TDD + Quality Gate（行为与当前一致）
- Given 一个任务模板文件, When `## Execution Workflow` 标题存在但正文为空, Then task-executor 记录配置错误警告并回退到 TDD

---

## Story 2: Task-executor Agent 跳过无效 TDD 循环

**As a** task-executor agent
**I want to** 读取任务中的 Execution Workflow 而非被强制进入 TDD 循环
**So that** 执行型任务（如运行 e2e 测试）不会在失败时无限重试，14 分钟浪费降为 < 5 分钟

**Acceptance Criteria:**

- Given T-test-3 模板包含"失败时创建 fix task，禁止重试"的 workflow, When e2e 测试失败, Then agent 创建 fix task 并停止（不进入 TDD 循环）
- Given 执行记录, When 检查 Step 2 输出, Then 包含 "Execution Workflow" 关键词，不包含 "TDD implementation" 或 "RED/GREEN/REFACTOR"
- Given 执行型任务模板, When workflow 执行完成, Then 直接进入 Step 3（记录 + 提交），不经过 Quality Gate

---

## Story 3: Forge Maintainer 清理 noTest 歧义

**As a** forge maintainer
**I want to** `noTest` / `NO_TEST` 从代码中完全消失
**So that** 新模板无需理解这个歧义字段，模板行为完全由 Execution Workflow 决定

**Acceptance Criteria:**

- Given 所有 task-cli / agent / command / skill 文件, When grep `noTest` 或 `NO_TEST`（大小写不敏感）, Then 零匹配
- Given task-cli 代码（types.go, record.go）, When 审查条件分支, Then 无基于 noTest 的残留逻辑
- Given 所有 16 个任务模板的 frontmatter, When 验证, Then 无 `noTest` 字段
- Given `index.schema.json`（breakdown + quick）, When 验证, Then 无 `noTest` 字段定义且 `ajv validate` 对所有模板通过
