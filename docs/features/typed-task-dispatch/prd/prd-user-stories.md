---
feature: "typed-task-dispatch"
---

# User Stories: typed-task-dispatch

## Story 1: 执行非编码类任务时获得正确的执行流程

**As a** forge 用户（使用 run-tasks 执行任务链路的开发者）
**I want to** 运行 doc-generation、test-pipeline、fix、gate 等非编码类任务时，agent 自动采用对应的执行流程
**So that** 不再出现文档生成任务走 TDD 流程、测试流水线任务被 noTest 字段绕过等流程错位问题

**Acceptance Criteria:**
- Given 一个 type 为 `doc-generation.summary` 的任务已在 index.json 中
- When run-tasks 执行该任务
- Then task-executor 不出现任何 TDD 相关步骤（无 RED/GREEN/REFACTOR，无 `just test` 调用），直接执行文档生成流程

- Given 一个 type 为 `fix` 的任务已在 index.json 中
- When run-tasks 执行该任务
- Then task-executor 执行诊断 → 定位 → 修复 → 验证 → 提交的五步流程，`go build ./...` 或 `go test ./...` 在修复后退出码为 0

---

## Story 2: 新增任务类型只需添加一个 prompt 模板文件

**As a** forge 维护者（维护和扩展 forge 插件的开发者）
**I want to** 新增一种任务类型时，只需在 task-cli 中添加一个 Go embed markdown 模板文件
**So that** 不再需要同时修改 task-executor.md、index.schema.json 和任务模板三处文件，将新增任务类型的改动成本从 ~2 小时降至可预期的单文件改动

**Acceptance Criteria:**
- Given 需要新增一种任务类型
- When 在 task-cli 的 prompt 模板目录中添加对应的 markdown 模板文件，并在 type 枚举中注册
- Then `task prompt <id>` stdout 包含任务 ID、scope、新 type 的执行步骤（来自新增模板文件），且不包含其他 type 的执行步骤；Go 单元测试可覆盖该类型，无需修改 task-executor.md 或任务模板

- Given 新增模板文件存在语法错误（如占位符格式错误）或 type 枚举未注册
- When 执行 `task prompt <id>` 指向该类型的任务
- Then 命令退出码非零，stderr 输出包含模板文件路径和具体错误原因，stdout 无输出

---

## Story 3: 通过 task prompt 命令独立检查合成结果

**As a** forge 维护者
**I want to** 在不运行 agent 的情况下，通过 `task prompt <id>` 直接查看某个任务的合成 prompt
**So that** 调试任务执行问题时，可以快速定位是 CLI 合成层（模板错误）还是 agent 执行层（约束层问题）的问题，将排查时间从 ~1 小时缩短到分钟级

**Acceptance Criteria:**
- Given 当前 feature 有一个 in_progress 任务（任意 type）
- When 执行 `task prompt <id>`
- Then stdout 输出该任务的完整合成 prompt，包含以下字段：任务 ID（`id`）、scope、任务 type、执行步骤（来自对应 type 模板）、phase summary 路径（若为新 phase 第一个任务）；命令在 500ms 内完成

- Given `task prompt <id>` 执行时 type 字段缺失或模板不存在
- When 命令执行
- Then 错误信息输出到 stderr，退出码非零，stdout 无输出

---

## Story 4: 迁移旧 feature 的 index.json 到新 type 字段

**As a** forge 用户
**I want to** 对已有 feature 的 index.json 运行 `task migrate`，自动填充 type 字段
**So that** 不需要手动判断每个任务的类型，可以快速完成迁移并继续使用新版 forge

**Acceptance Criteria:**
- Given 一个不含 type 字段的旧 index.json，且所有任务均为非 in_progress 状态
- When 执行 `task migrate`
- Then index.json 中所有任务均填充了正确的 type 字段（按推断规则），任务状态保持不变，`task validate` 对迁移后的 index.json 无报错

- Given index.json 中存在 in_progress 状态的任务
- When 执行 `task migrate`
- Then 命令报错并提示需先完成在途任务，不修改 index.json

---

## Story 5: breakdown-tasks 生成的任务自动包含 type 字段

**As a** forge 用户
**I want to** 运行 breakdown-tasks 或 quick-tasks 生成任务列表时，生成的 index.json 自动包含正确的 type 字段
**So that** 新建 feature 无需手动运行 task migrate，可以直接使用新版 forge 执行任务

**Acceptance Criteria:**
- Given 一个完整的 tech-design 文档
- When 执行 breakdown-tasks
- Then 生成的 index.json 中所有任务均包含 type 字段，`task validate` 无报错，且 type 值与任务的实际类型（implementation、doc-generation.summary、gate 等）一致

- Given tech-design 文档中某个任务的描述无法匹配任何已知 type 推断规则
- When 执行 breakdown-tasks
- Then 该任务的 type 字段回退为 `implementation`，并在 stderr 输出警告（包含任务 ID 和无法推断的原因），不中断整体生成流程

---

## Story 6: execute-task 与 run-tasks 路由行为保持一致

**As a** forge 用户（直接调用 execute-task 执行单个任务的开发者）
**I want to** execute-task 采用与 run-tasks 相同的 task prompt 路由逻辑
**So that** 无论通过 run-tasks 还是 execute-task 执行任务，行为一致，不出现路由差异导致的执行结果不同

**Acceptance Criteria:**
- Given index.json 中存在一个 type 为 `implementation` 的任务
- When 通过 execute-task 执行该任务
- Then execute-task 调用 `task prompt <id>` 合成 prompt 后 dispatch 给 forge:task-executor，不使用旧的 TASK_FILE + NO_TEST 参数组合；task-executor 收到的 prompt 内容与 run-tasks 路径下对同一任务生成的 prompt 内容相同

- Given `task prompt <id>` 执行失败（退出码非零）
- When execute-task 调用 task prompt
- Then execute-task 将该任务标记为 blocked，stderr 输出包含 task prompt 的错误信息，不静默失败

---

## Story 7: error-fixer agent 废弃后 fix 流程通过新路径等价覆盖

**As a** forge 维护者
**I want to** 确认 fix 类任务和 record 缺失恢复均已由新的 prompt 路径处理
**So that** error-fixer agent 可以安全移除，不遗留任何对 error-fixer 的 dispatch 调用

**Acceptance Criteria:**
- Given index.json 中存在一个 type 为 `fix` 的任务
- When run-tasks 执行该任务
- Then task-executor 收到的 prompt 来自 `task prompt <id>`，包含诊断 → 定位 → 修复 → 验证 → 提交五步流程；run-tasks.md 中不存在对 forge:error-fixer 的 dispatch 调用

- Given 任务执行完成后 record 文件缺失
- When run-tasks 检测到 record 缺失
- Then run-tasks 调用 `task prompt <id> --fix-record-missed` 并将结果 dispatch 给 forge:task-executor；error-fixer agent 不被调用；error-fixer.md 中的 dispatch 入口已从 run-tasks.md 移除
