---
feature: "justfile-e2e-integration"
---

# User Stories: Justfile E2E Integration

## Story 1: Skill 维护者使用统一命令接口

**As a** Skill 维护者
**I want to** 在 SKILL.md 中只写 `just e2e-setup` 和 `just test-e2e --feature <slug>`，而不是具体的 `npx` 或 `npm` 命令
**So that** 当底层工具链变化时，只需修改 justfile recipe，所有 skill 自动受益，无需逐一更新文档

**Acceptance Criteria:**

- Given `init-justfile` 已在应用项目中运行，justfile 包含 `e2e-setup` 和 `e2e-verify` 目标
- When Skill 维护者查看 `run-e2e-tests` SKILL.md 的 Step 1
- Then 文档中只出现 `just e2e-setup`，不出现 `cd tests/e2e && npm install` 或 `npx playwright install chromium`

---

## Story 2: AI Agent 执行 e2e 任务时获得明确指令

**As a** AI Agent（执行 skill 指令的自动化代理）
**I want to** 在 skill 文档中看到明确的 just 命令（如 `just test-e2e --feature <slug>`），而不是需要自行推断的占位符（如 `<project-test-command>`）
**So that** 我能直接执行命令，不需要判断项目语言或推断工具链，减少执行错误

**Acceptance Criteria:**

- Given Agent 正在执行 `task-executor` Step 3 Full Verification
- When Agent 读取 Step 3 的指令
- Then 指令明确写 `just build && just test`，不出现 `go test ./...`、`npm test`、`pytest` 等语言特定命令

---

## Story 3: AI Agent 在 gen-test-scripts 后验证 VERIFY 标记

**As a** AI Agent
**I want to** 在生成 spec 文件后运行 `just e2e-verify --feature <slug>` 并根据 exit code 决定是否继续
**So that** 含有未解析 `// VERIFY:` 标记的脚本不会进入 `run-e2e-tests`，避免浪费一轮执行

**Acceptance Criteria:**

- Given Agent 已通过 `gen-test-scripts` 生成 spec 文件
- When Agent 运行 `just e2e-verify --feature <slug>`，且 spec 文件中存在残留 `// VERIFY:` 标记
- Then 命令 exit 1，Agent 将 skill 标记为 incomplete，不执行 `run-e2e-tests`

- Given Agent 已通过 `gen-test-scripts` 生成 spec 文件
- When Agent 运行 `just e2e-verify --feature <slug>`，且无残留标记
- Then 命令 exit 0，Agent 继续执行 `run-e2e-tests`

---

## Story 4: AI Agent 在 fix-e2e 任务后验证修复结果

**As a** AI Agent
**I want to** 在修复 e2e 失败后，通过 `just test-e2e --feature <slug>` 验证修复是否有效
**So that** 修复结果有明确的验证步骤，不依赖 agent 自行推断如何重跑测试

**Acceptance Criteria:**

- Given Agent 已完成 fix-e2e 任务中的代码修复
- When Agent 读取 `fix-e2e` task 模板的 Implementation Notes
- Then 模板明确写 "运行 `just test-e2e --feature <slug>` 验证修复"，不出现 `npx tsx` 或其他原始命令

---

## Story 5: AI Agent 执行构建/测试任务时获得统一命令

**As a** AI Agent（执行 fix-bug、run-tasks、task-executor、error-fixer、execute-task、record-task、improve-harness 等 skill/command 的自动化代理）
**I want to** 在上述所有 skill 和 command 文档中看到 `just test` 或 `just build && just test`，而不是 `go test ./...`、`npm test`、`pytest --cov` 等语言特定命令
**So that** 我在跨语言项目中执行构建和测试时，无需判断项目语言或推断工具链，直接调用统一命令即可

**Acceptance Criteria:**

- Given Agent 正在执行 `fix-bug` command，项目为任意语言
- When Agent 读取 fix-bug 中的测试验证步骤
- Then 文档中只出现 `just test`，不出现 `<project-test-command>` 占位符或任何语言特定测试命令

- Given Agent 正在执行 `run-tasks` command 的 Breaking Gate 检查
- When Agent 读取 Breaking Gate 指令
- Then 指令明确写 `just test`，不出现 `npm test`、`go test` 等原始命令

- Given Agent 正在执行 `record-task` SKILL.md 的 Metrics Collection 步骤
- When Agent 读取语言示例中的测试命令
- Then 示例统一为 `just test`，不出现 `go test -cover ./...`、`npm test -- --coverage`、`pytest --cov=...`
