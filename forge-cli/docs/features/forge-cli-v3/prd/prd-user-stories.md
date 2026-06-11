---
feature: "forge-cli-v3"
---

# User Stories: Forge CLI v3

## Story 1: AI Agent 发现并执行正确命令

**As a** AI agent（任务执行者）
**I want to** 通过 `forge --help` 快速定位到正确的任务管理命令
**So that** 无需阅读文档即可完成 claim → 执行 → submit 的完整任务生命周期

**Acceptance Criteria:**
- Given AI agent 收到任务执行指令
- When agent 调用 `forge --help`
- Then 输出显示 5 个命令组（task/e2e/forensic/profile/prompt）+ 5 个顶层命令，共 10 个入口
- And agent 能通过 `forge task --help` 看到 11 个子命令，每个子命令描述包含"命令名+动词+宾语"三要素（如 "submit task execution result"），描述长度 <= 80 字符
- Given AI agent 调用 `forge <不存在的命令>`
- When 输入拼写错误的命令如 `forge taks`
- Then 退出码为 1，stderr 输出包含 "unknown command" 提示及最接近的命令建议
- Given AI agent 调用 `forge task <不存在的子命令>`
- When 输入无效子命令
- Then 退出码为 1，stderr 输出 "unknown subcommand" 并列出有效子命令列表

---

## Story 2: AI Agent 通过 prompt 获取执行指令

**As a** AI agent（任务执行者）
**I want to** 通过任务 ID 获取对应的 agent prompt
**So that** 知道该执行什么类型的任务（implementation/fix/gate/test-pipeline 等）

**Acceptance Criteria:**
- Given AI agent 需要执行任务 T-impl-1（type: implementation）
- When agent 调用 `forge prompt get-by-task-id T-impl-1`
- Then 命令输出该任务的 implementation 类型 prompt，退出码 0
- And prompt 包含 `{{TASK_ID}}`、`{{TASK_FILE}}`、`{{SCOPE}}` 等变量的替换结果
- Given AI agent 使用不存在的任务 ID
- When agent 调用 `forge prompt get-by-task-id NONEXISTENT-999`
- Then 退出码为 1，stderr 输出 "task not found" 错误信息，stdout 为空
- Given 任务 ID 对应的任务缺少 type 字段或 type 值不在已知类型列表中
- When agent 调用 `forge prompt get-by-task-id <id>`
- Then 退出码为 1，stderr 输出 "unknown task type" 或 "missing task type" 错误信息

---

## Story 3: AI Agent 提交任务结果

**As a** AI agent（任务执行者）
**I want to** 提交任务执行结果并自动更新任务状态
**So that** 任务生命周期正确流转（pending → completed/blocked/rejected）

**Acceptance Criteria:**
- Given AI agent 完成任务 T-impl-1 的执行
- When agent 调用 `forge task submit T-impl-1 --result success --summary "..."`
- Then index.json 中 T-impl-1 状态更新为 completed，records/ 目录下生成对应记录文件
- And 退出码为 0
- Given 任务 T-impl-1 已处于 completed 状态
- When agent 再次调用 `forge task submit T-impl-1 --result success --summary "..."`
- Then 退出码为 1，stderr 输出 "task already in terminal state" 错误信息，index.json 不变更
- Given agent 提交时缺少必需参数（如省略 --result）
- When agent 调用 `forge task submit T-impl-1 --summary "..."`
- Then 退出码为 1，stderr 输出 "required flag(s) not set: result" 错误信息
- Given 两个 agent 同时对同一任务 T-impl-1 调用 `forge task submit`
- When 并发写 index.json 发生 lock 竞争
- Then 恰好一个 agent 退出码为 0（提交成功），另一个退出码为 1 且 stderr 包含 "concurrent write conflict, retry"；index.json 可被 `jq .` 正常解析且 JSON 语法合法

---

## Story 4: Hook 自动触发生命周期管理

**As a** CI/Hook 自动化流程
**I want to** 在正确时机自动调用 cleanup、quality-gate、verify-task-done
**So that** 任务状态一致性和代码质量得到保障

**Acceptance Criteria:**
- Given SessionEnd hook 触发
- When hook 调用 `forge cleanup`
- Then 已完成/已阻塞/已拒绝任务的状态文件被清理，.forge/state.json 保持不变
- And 退出码为 0

- Given Stop hook 触发（所有任务完成后）
- When hook 调用 `forge quality-gate`
- Then 依次执行 compile → fmt → lint → test，任何步骤失败则创建 P0 fix-task
- And 退出码反映质量门禁结果（0=通过，1=失败）

- Given SessionEnd hook 触发但无任何任务处于终态（completed/blocked/rejected）
- When hook 调用 `forge cleanup`
- Then 退出码为 0，无文件被删除，stdout 输出 "no tasks to clean up"

- Given quality-gate 创建了 fix-task 但该 fix-task 对应的修复也失败
- When 下一次 Stop hook 触发时
- Then `forge quality-gate` 再次检测到失败步骤，创建新的 P0 fix-task（不覆盖已有 fix-task）
- And 新 fix-task 的 title 包含失败步骤名称和序号（如 "fix-compile-3"），grep `fix-compile-` 在 index.json 中恰好返回 N 条（N = 该步骤的 fix-task 总数）

- Given 同一失败步骤已产生 3 个未完成的 fix-task
- When `forge quality-gate` 再次检测到该步骤失败
- Then 不再创建新 fix-task，改为 stderr 输出 "max fix-tasks reached for <step>, manual intervention required"，退出码为 1

---

## Story 5: 开发者运行 e2e 测试

**As a** 开发者
**I want to** 使用 `forge e2e run` 运行 profile-aware 的 e2e 测试
**So that** 无需关心 profile 检测逻辑即可执行正确的测试套件

**Acceptance Criteria:**
- Given 项目配置了 web-playwright profile（.forge/config.yaml 中 profile 字段为 web-playwright）
- When 开发者调用 `forge e2e run --feature my-feature`
- Then 命令读取 config.yaml 中 profile 字段，选择对应的 e2e 测试套件执行，stdout 输出 "profile: web-playwright"
- And 退出码与原 justfile `test-e2e` 一致

- Given 项目 .forge/config.yaml 中未配置 profile 字段或 profile 值为空
- When 开发者调用 `forge e2e run --feature my-feature`
- Then 退出码为 1，stderr 输出 "no e2e profile configured" 错误信息

- Given .forge/config.yaml 中 profile 值为不在支持列表中的字符串（如 "unknown-profile"）
- When 开发者调用 `forge e2e run --feature my-feature`
- Then 退出码为 1，stderr 输出 "unknown profile: unknown-profile" 并列出所有支持的有效 profile

- Given 开发者调用 `forge e2e run --feature nonexistent-feature`
- When 指定的 feature 名称在项目中不存在对应目录或测试文件
- Then 退出码为 1，stderr 输出 "feature not found: nonexistent-feature"

---

## Story 6: 开发者查看任务类型

**As a** 开发者
**I want to** 列出所有支持的任务类型及其含义
**So that** 在编写 task markdown 时使用正确的 type 字段

**Acceptance Criteria:**
- Given 开发者不确定可用的任务类型
- When 开发者调用 `forge task list-types`
- Then 输出包含所有 11 种任务类型（implementation, fix, gate, doc-generation.*, test-pipeline.*）
- And 每种类型附带描述，描述格式为"动词+宾语"结构且长度 <= 60 字符，退出码 0

- Given 任务类型注册表为空（无已定义类型）
- When 开发者调用 `forge task list-types`
- Then stdout 输出空列表（0 行类型记录），退出码为 0

---

## Story 7: 开发者排查 Agent 会话问题

**As a** 开发者
**I want to** 搜索和分析 Claude Code 会话记录，定位 agent 偏差根因
**So that** 在任务执行出现异常时能快速定位问题发生的时间点和原因

**Acceptance Criteria:**
- Given 开发者需要排查某次任务执行中的 agent 行为偏差
- When 开发者调用 `forge forensic search --project-path .`
- Then 命令扫描 history.jsonl，输出匹配的会话列表（包含 session ID、时间戳、skill 名称），退出码 0

- Given 开发者找到目标会话 ID
- When 开发者调用 `forge forensic extract <session-jsonl-path>`
- Then 命令输出该会话的紧凑证据摘要（包含关键决策点、工具调用序列、偏离预期的节点），退出码 0

- Given 开发者需要查看目标会话中的子代理行为
- When 开发者调用 `forge forensic subagents <session-dir-path>`
- Then 命令列出该会话下所有子代理的 transcript 文件路径和摘要信息，退出码 0

- Given 开发者指定的 session-jsonl-path 不存在
- When 开发者调用 `forge forensic extract /nonexistent/path.jsonl`
- Then 退出码为 1，stderr 输出 "file not found: /nonexistent/path.jsonl"

---

## Story 8: 开发者管理测试 Profile

**As a** 开发者
**I want to** 设置、检测和查看 e2e 测试 profile 配置
**So that** 在不同项目环境中使用正确的测试工具链（Playwright、Go test、Maestro 等）

**Acceptance Criteria:**
- Given 项目目录中存在多种测试框架的配置文件（playwright.config.ts、*_test.go 等）
- When 开发者调用 `forge profile detect`
- Then 命令扫描项目结构，输出检测到的所有 profile 及其依据（如 "web-playwright: found playwright.config.ts"），退出码 0

- Given 开发者确认使用 web-playwright profile
- When 开发者调用 `forge profile set web-playwright`
- Then .forge/config.yaml 中 profile 字段更新为 "web-playwright"，退出码为 0

- Given 开发者需要确认当前 profile 的详细配置
- When 开发者调用 `forge profile get web-playwright`
- Then 命令输出该 profile 的 strategy 文件内容（generate.md、run.md 等），退出码为 0

- Given 开发者设置不存在的 profile
- When 开发者调用 `forge profile set nonexistent-profile`
- Then 退出码为 1，stderr 输出 "unknown profile: nonexistent-profile" 并列出所有支持的有效 profile
