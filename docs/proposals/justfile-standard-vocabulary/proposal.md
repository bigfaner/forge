---
created: 2026-04-29
author: faner
status: Draft
---

# Proposal: Justfile 标准命令词汇表

## Problem

Forge 插件的 skill/agent/command 文件中散落着原始 shell 命令。当前 init-justfile 命令定义了 6 个标准目标（`test`, `test-e2e`, `build`, `lint`, `e2e-setup`, `e2e-verify`），但存在两个核心缺陷：

1. **词汇不完整**：缺少 `compile`（编译）、`run`（启动服务）、`dev`（热重载）、`fmt`（格式化）、`check`（lint+compile）、`clean`、`install`、`ci` 等常见操作。Skills 遇到这些需求时只能内联原始命令。
2. **不支持前后端分离项目**：混合项目（如 forge 自身有 Go 后端 + React 前端）无法选择性构建/运行某个 scope。Skills 无法感知项目结构，只能执行全量操作。

### Evidence

- `run-e2e-tests` skill 手动执行 `npx serve` 启动服务器，而非调用标准化的 `just run`
- `execute-task` 和 `error-fixer` agent 调用 `just build && just test`，但 `build` recipe 在本项目中不存在
- 多个项目（如 pm-work-tracker）前后端混合，无法用一条命令选择性操作
- 没有机制让 skill 知道项目是纯前端、纯后端还是混合项目

## Solution

### 标准命令词汇表（16 个命令）

| 命令 | 位置参数 | 用途 |
|------|---------|------|
| `compile [scope]` | `frontend`/`backend` | 类型检查 + 转译，快速反馈 |
| `build [scope]` | `frontend`/`backend` | 完整构建（打包、编译产物） |
| `run [scope]` | `frontend`/`backend` | 启动服务（dev server / API server） |
| `dev [scope]` | `frontend`/`backend` | 热重载开发模式 |
| `test [scope]` | `frontend`/`backend` | 单元 + 集成测试 |
| `test-e2e [--feature <slug>]` | — | E2E 测试 |
| `lint [scope]` | `frontend`/`backend` | 静态分析 |
| `fmt [scope]` | `frontend`/`backend` | 自动格式化代码 |
| `check [scope]` | `frontend`/`backend` | lint + compile（CI 门禁） |
| `clean [scope]` | `frontend`/`backend` | 清理构建产物 |
| `install [scope]` | `frontend`/`backend` | 安装依赖 |
| `ci` | — | 完整 CI 流水线 |
| `e2e-setup` | — | 安装 e2e 依赖（幂等） |
| `e2e-verify --feature <slug>` | — | 检查未解析的 `// VERIFY:` 标记 |
| `project-type` | — | 返回项目类型（单字输出） |

### `just project-type` — 项目自描述

返回单个单词，skill 据此决定是否传 scope 参数：

| 输出 | 含义 | skill 行为 |
|------|------|-----------|
| `frontend` | 纯前端项目 | `just build`（无 scope） |
| `backend` | 纯后端项目 | `just build`（无 scope） |
| `mixed` | 前后端混合 | `just build frontend` / `just build backend` / `just build` |

### 自适应 Justfile 生成

`init-justfile` 探测项目结构后，生成与项目类型匹配的 justfile：

| 项目类型 | Justfile 风格 | 示例 |
|---------|-------------|------|
| 纯后端 | 无 scope 参数 | `just build` → `go build ./...` |
| 纯前端 | 无 scope 参数 | `just build` → `npm run build` |
| 前后端混合 | 有 scope 参数 | `just build frontend` / `just build backend` / `just build`（全部） |

### 任务级 scope 标记

`breakdown-tasks` 在拆分任务时，根据任务涉及的代码路径为每个任务在 `index.json` 中标记 scope：

```json
{
  "id": "1.1",
  "title": "Implement login API endpoint",
  "scope": "backend",
  ...
}
```

| scope 值 | 含义 | skill 行为 |
|----------|------|-----------|
| `frontend` | 仅涉及前端代码 | `just build frontend` |
| `backend` | 仅涉及后端代码 | `just build backend` |
| `all`（默认） | 跨前后端或全栈 | `just build`（全部） |

### Skill 集成方式

Skills 通过两层信息决定 scope：

```
1. 优先读取当前任务的 scope 字段（来自 index.json）
2. 若任务无 scope 或 scope=all：
     执行 `just build`（无 scope 参数，justfile 自行决定构建范围）
3. 若任务有明确 scope（frontend/backend）：
     执行 `just project-type`
     if output == "mixed":
       → `just build frontend` 或 `just build backend`
     else:
       ⚠ scope 与项目类型不匹配（如 scope=frontend 但项目为纯 backend）
       记录警告：可能是 breakdown-tasks scope 分配错误
       → 回退执行 `just build`（无 scope 参数）
```

### Developer Experience

迁移前后，开发者在终端中看到的命令调用方式发生如下变化：

| 场景 | 迁移前（当前） | 迁移后（本提案） |
|------|---------------|-----------------|
| 纯后端项目运行测试 | `go test ./...` | `just test` |
| 混合项目运行前端测试 | `cd frontend && npm test` | `just test frontend` |
| 混合项目启动开发服务器 | `npx serve -l 3000`（硬编码端口） | `just dev frontend`（justfile 配置端口） |
| 纯前端项目构建 | `npm run build` | `just build` |
| scope 不匹配 | 无检测，执行失败或静默跳过 | `[forge] scope=frontend but project-type=backend; falling back to just build`（显式警告后回退） |

## Alternatives

### A. 本提案：完整词汇表 + 全量迁移

更新 init-justfile 契约（16 命令）、添加 `project-type`、更新所有 skill/agent/command 文件、更新 forge 项目 justfile 作为参考实现。

- **优点**：统一接口，所有 skill 立即受益，消除散落的原始命令
- **缺点**：工作量大（8+ 文件），需同时更新 skill 文档和 e2e 测试

### B. 仅词汇表 + 选择性迁移

只更新 init-justfile 契约和 2 个高价值 skill（run-e2e-tests、execute-task），其余逐步采纳。

- **优点**：工作量小，降低风险
- **缺点**：接口不一致，部分 skill 继续使用原始命令

### C. 不标准化，仅本地添加

只在 forge 项目 justfile 中添加缺失的 recipe，不改契约。

- **优点**：最小改动
- **缺点**：其他项目无法复用，skill 无法通用化

### Decision: 选择方案 A

方案 B 看似务实的折中，但会导致同一项目中部分 skill 走 `just` 路径、其余 skill 继续内联原始 shell 命令。具体失败模式：`execute-task` 调用 `just compile`（方案 B 已更新），而同一次任务流中 `fix-bug` 仍直接调用 `go test`（方案 B 未更新），两者对 "如何运行测试" 的认知不一致——justfile 中 `just test` 可能包含覆盖率标志或环境变量设置，而内联的 `go test` 不包含，导致任务执行与修复验证在非等同条件下运行。

方案 A 需更新 10 个 skill/agent/command 文件 + 4 个 task 模板 + 1 个 breakdown-tasks skill + 1 个 init-justfile 契约 = 共 16 处修改。这些修改按复杂度分为三类：

- **Category A — 机械替换（14 处）**：skill/agent/command 文件和 task 模板中的命令字符串替换（如 `go test ./...` → `just test`）。风险低，每处约 5 分钟，总计约 1 小时。
- **Category B — 设计工作（1 处）**：init-justfile 自适应生成逻辑，需为三种项目类型设计探测规则和 justfile 模板。中等复杂度，约 1-2 天。
- **Category C — Prompt 工程（1 处）**：breakdown-tasks scope 标注，需修改 AI agent 的指令使其根据代码路径推断 scope。结果不完全确定，需迭代验证，约 1 天。

一次性迁移消除了方案 B 中持续维护两套调用约定的过渡期，而过渡期的长度不确定，且每次部分迁移都需回归测试两套路径。

## Scope

### In Scope

1. 更新 `init-justfile` 命令：扩展标准目标契约至 16 个命令，含自适应生成逻辑
2. 添加 `project-type` recipe 到契约
3. 更新以下 skill 文件以使用新词汇：
   - `run-e2e-tests/SKILL.md`（`just run` 替代 `npx serve`）
   - `gen-test-scripts/SKILL.md`（`just e2e-setup`、`just e2e-verify`）
   - `record-task/SKILL.md`（`just test`）
   - `execute-task.md`（`just compile`、`just test`）
   - `fix-bug.md`（`just test`、`just test-e2e`）
   - `run-tasks.md`（`just test`）
   - `improve-harness.md`（`just test`）
   - `graduate-tests/SKILL.md`（`just test-e2e`）
4. 更新 agent 文件：
   - `task-executor.md`（`just compile && just test`）
   - `error-fixer.md`（`just compile && just test`）
5. 更新 task 模板：
   - `breakdown-tasks/templates/run-e2e-tests.md`
   - `breakdown-tasks/templates/gen-test-scripts.md`
   - `breakdown-tasks/templates/fix-e2e.md`
   - `breakdown-tasks/templates/graduate-tests.md`
6. 更新 `breakdown-tasks` skill：任务拆分时在 `index.json` 中为每个任务添加 `scope` 字段
7. 更新 forge 项目 justfile 作为参考实现（mixed 项目）

### Out of Scope

- 其他项目（pm-work-tracker 等）的 justfile 更新（用户自行运行 `/init-justfile`）
- `just project-type` 在 CI 环境中的集成
- 自动检测 feature 变更文件归属 frontend/backend 的算法（通过 breakdown-tasks 的 scope 标记解决）
- justfile recipe 的并行执行优化（just 原生支持）

## Risks

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|---------|
| 自适应生成逻辑复杂，难以覆盖所有项目结构 | Medium | High — init-justfile 生成不正确的 justfile，导致所有依赖 just 的 skill 执行失败 | 为三种项目类型（frontend/backend/mixed）各提供经过验证的 fallback 模板；生成后运行 `just --list` 验证所有 recipe 可达 |
| Skill 行为变更可能导致现有工作流中断 | High | High — e2e 测试或任务执行失败阻塞开发流程 | 分两阶段部署：(1) 先更新 justfile 确保 recipe 存在，(2) 再更新 skill 引用。当前 e2e 覆盖：`justfile-e2e-integration` 测试套件验证 skill 文件内容匹配（20/20 通过），但本次修改的 8 个 skill 文件和 4 个 task 模板均无独立的 per-skill e2e 测试。对无 e2e 覆盖的 skill，逐个手动验证：执行更新后的 skill 文件中引用的每个 `just <verb>` 命令，确认 justfile 中对应 recipe 存在且返回码为 0；若执行失败，回滚该 skill 文件至上一版本并记录失败 recipe 名称 |
| `just project-type` 输出不准确导致 skill 传错 scope | Medium | Medium — 在非 mixed 项目上多传 scope 参数，recipe 报参数错误 | init-justfile 自动生成该 recipe（基于项目探测，非人工填写）；skill 层做防御：若 `project-type` 返回非 `mixed` 则忽略 scope 参数 |
| 已有自定义 justfile 的项目重新运行 init-justfile 时丢失自定义 recipe | Medium | High — 用户手动添加的 recipe 和变量被覆盖，项目构建中断 | init-justfile 检测已有 justfile 时提示用户确认，而非静默覆盖；标准 recipe 区域用 `# --- forge standard recipes ---` 标记边界，保留标记外的用户自定义内容 |

## Success Criteria

1. `init-justfile` 根据项目结构生成正确的 justfile 变体，具体验证场景：
   - 给定一个包含 `package.json` 但无 `go.mod` 的项目，init-justfile 生成无 scope 参数的 justfile，且 `just project-type` 输出 `frontend`
   - 给定一个包含 `go.mod` 但无 `package.json` 的项目，init-justfile 生成无 scope 参数的 justfile，且 `just project-type` 输出 `backend`
   - 给定一个同时包含 `package.json`（前端目录）和 `go.mod`（后端目录）的项目，init-justfile 生成带 scope 参数的 justfile（`just build frontend` / `just build backend` / `just build`），且 `just project-type` 输出 `mixed`
2. `just project-type` 在三种项目类型下都返回正确的单个单词
3. 所有 skill/agent/command 文件中，任何可执行指令（Bash tool call 或代码块）若执行 build、test、lint、format、compile、run、dev、clean、install 操作，必须通过 `just <verb>` 调用。直接调用语言工具链（`go`、`npm`、`cargo`、`pytest`、`npx` 等）为违规。文档注释、示例说明、非可执行文本中的命令不受此规则约束
4. `run-e2e-tests` skill 通过 `just run` 启动服务器，不再硬编码 `npx serve`
5. forge 项目自身的 justfile 包含全部 16 个标准命令，且可通过 `just project-type` → `mixed` 验证
6. 现有 e2e 测试（`justfile-e2e-integration`）全部通过
7. `breakdown-tasks` 生成的 `index.json` 中每个任务包含 `scope` 字段，值为 `frontend`、`backend` 或 `all`。给定一个包含前端专用任务（如 "Update login form component"）和后端专用任务（如 "Implement login API endpoint"）的混合项目，breakdown-tasks 正确分配 `scope=frontend` 和 `scope=backend`；跨端任务默认 `scope=all`
