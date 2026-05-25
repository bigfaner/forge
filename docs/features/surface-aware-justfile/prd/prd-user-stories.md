---
feature: "Surface-Aware Justfile"
---

# User Stories: Surface-Aware Justfile

## Story 1: Surface 感知配方生成

**As a** Forge 用户（项目开发者）
**I want to** 配置 surfaces 字段后运行 init-justfile 自动生成对应 surface 类型的 dev/test/probe 配方
**So that** web/api 项目获得正确的"启动→等待→测试→清理"编排，cli/tui 项目获得"构建→测试"编排

**Acceptance Criteria:**
- Given 项目 config.yaml 定义了 `surfaces: {admin-panel: web}`
- When 运行 init-justfile
- Then 生成的 justfile 包含 dev(后台启动)、probe(重试轮询)、test、test-teardown 配方
- And 生成的配方包含 `[linux]`/`[windows]` 双平台变体
- And CLI/TUI surface 不生成 `run` 配方

---

## Story 2: 自动化测试编排

**As a** Forge 用户（项目开发者）
**I want to** 运行 run-tests 时自动检测 surface 类型并执行正确的编排序列
**So that** 我不需要手动配置 test.execution，测试编排开箱即用

**Acceptance Criteria:**
- Given 任务 frontmatter 包含 `surface-type: web`
- When run-tests 执行
- Then 按序列执行 just dev → just probe → just test → just test-teardown
- And probe 失败时执行 teardown 后中止：exit 1（retryable）或 exit 2（blocking）区分语义
- And probe 失败后禁止在同一编排周期内重试 probe 或重启 dev（HARD-GATE）

---

## Story 3: Surface-key 统一迁移

**As a** Forge 插件开发者
**I want to** 将 surface-key 值域从固定枚举（frontend/backend）统一迁移为用户自定义 surface-key 名称（surface-type 保持 5 种固定类型）
**So that** 混合项目的所有配方和任务使用一致的标识符，消除硬编码约束

**Acceptance Criteria:**
- Given 项目 surfaces 定义为 `{admin-panel: web, payment-service: api}`
- When breakdown-tasks 生成任务
- Then 任务的 surface-key 为 `admin-panel` 或 `payment-service`（而非 frontend/backend）
- And 无 surfaces 配置的项目行为不变
- And prompt.go resolveScope() 基于 surfaces map 集合查询而非 projectType 硬编码

---

## Story 4: Task 数据模型扩展

**As a** Forge 插件开发者
**I want to** Task 数据模型新增 surface-key 和 surface-type 字段
**So that** 下游 skill（run-tests/execute-task/init-justfile）能直接从任务获取 surface 信息，无需额外查询

**Acceptance Criteria:**
- Given breakdown-tasks 生成一个涉及 web surface 文件的任务
- When 查看 index.json
- Then 任务包含 `surface-key: "admin-panel"` 和 `surface-type: "web"`
- And 旧任务文件含 `scope: frontend` 时，`forge task migrate` 可将其自动迁移为 `surface-key` + `surface-type`，迁移前所有 task 读取命令返回阻塞错误（exit 2）
- And forge task add 从源任务继承 surface-key 和 surface-type
- And quality-gate fix-task 从失败文件路径自动推断 surface-key/type
