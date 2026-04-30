---
feature: "justfile-standard-vocabulary"
---

# User Stories: Justfile 标准命令词汇表

## Story 1: 标准化命令执行

**As a** Forge skill 用户（开发者）
**I want to** 所有 forge skill 通过标准化的 `just <verb>` 命令执行构建、测试、运行等操作
**So that** 我不需要关心项目使用什么语言工具链，一套命令适用于所有 forge 项目

**Acceptance Criteria:**
- Given 一个使用 forge 插件的项目，When 任意 skill 执行构建/测试/编译操作时，Then 通过 `just <verb>` 调用，不直接调用 `go`、`npm`、`cargo` 等语言工具链
- Given 一个纯后端 Go 项目，When 开发者执行 `just test` 时，Then 实际执行 `go test -race ./...`
- Given 一个混合项目，When 开发者执行 `just build frontend` 时，Then 仅执行前端构建

---

## Story 2: 自适应项目初始化

**As a** 项目维护者
**I want to** 运行 `/init-justfile` 时自动识别项目类型并生成匹配的 justfile
**So that** 我不需要手动编写或调整 justfile 中的 scope 参数配置

**Acceptance Criteria:**
- Given 一个包含 `package.json` 但无 `go.mod` 的项目，When 运行 `/init-justfile`，Then 生成的 justfile 无 scope 参数，`just project-type` 输出 `frontend`
- Given 一个包含 `go.mod` 但无 `package.json` 的项目，When 运行 `/init-justfile`，Then 生成的 justfile 无 scope 参数，`just project-type` 输出 `backend`
- Given 一个同时包含 `package.json`（前端目录）和 `go.mod`（后端目录）的项目，When 运行 `/init-justfile`，Then 生成的 justfile 支持 `frontend`/`backend` scope 参数，`just project-type` 输出 `mixed`

---

## Story 3: 智能任务 scope 标记

**As a** Forge skill 作者
**I want to** breakdown-tasks 自动为每个任务标记 scope（frontend/backend/all）
**So that** 执行 skill 时能根据任务 scope 选择性构建/测试，避免不必要的全量操作

**Acceptance Criteria:**
- Given 一个混合项目的 tech design，When `/breakdown-tasks` 拆分任务时，Then `index.json` 中每个任务包含 `scope` 字段
- Given 一个任务仅涉及前端目录文件（如 `web/src/components/Button.tsx`），When `/breakdown-tasks` 处理该任务时，Then scope 标记为 `frontend`
- Given 一个任务涉及前后端目录文件，When `/breakdown-tasks` 处理该任务时，Then scope 标记为 `all`
- Given 一个非混合项目，When `/breakdown-tasks` 拆分任务时，Then 所有任务的 scope 标记为 `all`（无前后端区分）

---

## Story 4: Agent 友好的命令输出

**As a** Forge agent
**I want to** 所有 `just <verb>` 命令返回可预测的退出码、无交互式提示、错误信息输出到 stderr
**So that** 我能可靠地判断命令执行结果并自动处理错误，无需人工干预

**Acceptance Criteria:**
- Given 任意 `just <verb>` 命令，When 执行成功时，Then 退出码为 0，stdout 包含正常输出
- Given 任意 `just <verb>` 命令，When 执行失败时，Then 退出码非 0，stderr 包含错误信息
- Given agent 执行 `just compile`，When 代码有类型错误时，Then 退出码非 0 且错误详情输出到 stderr，agent 可解析错误信息并尝试修复
- Given agent 连续执行 `just install`、`just compile`、`just test`，When 全部成功时，Then 每步退出码均为 0，agent 无需人工介入

---

## Story 5: scope 不匹配防御

**As a** Forge skill 用户（开发者）
**I want to** 当任务的 scope 与项目类型不匹配时收到警告
**So that** 我能发现 breakdown-tasks 可能分配了错误的 scope

**Acceptance Criteria:**
- Given 一个纯后端项目和 scope=frontend 的任务，When skill 执行 `just build` 时，Then 显示警告信息 `[forge] scope=frontend but project-type=backend; falling back to just build` 并回退执行无 scope 的命令
- Given 一个混合项目和 scope=frontend 的任务，When skill 执行时，Then 正常执行 `just build frontend`，无警告
