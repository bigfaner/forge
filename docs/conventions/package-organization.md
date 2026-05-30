---
title: "Package Organization"
domains: [architecture, dependencies, pkg, cmd]
---

# Package Organization

本文档定义 forge-cli 的包组织规范。这是**规范性（normative）**文档，描述目标状态而非当前状态。所有新增代码必须符合此规范；已有代码的偏差记录在偏差分析表中，按计划逐步收敛。

## 1. 依赖方向规则

```
cmd/ → internal/ → pkg/
```

严格单向，禁止反向依赖：

- `internal/cmd/` 可以导入 `internal/` 下的其他包和 `pkg/` 的任何包
- `internal/` 下的非 cmd 包可以导入 `pkg/`，不能导入 `internal/cmd/`
- `pkg/` 下的包不能导入 `internal/` 或 `internal/cmd/` 的任何内容
- `pkg/` 内部遵循三层模型的方向约束（见第 3 节）

**违规示例**：`pkg/feature` 导入 `internal/cmd/task` 的函数 -- 这违反了依赖方向。

## 2. internal/cmd/ 目标状态

`internal/cmd/` 是 CLI 命令的入口层，职责是：
- 解析命令行参数
- 调用 `pkg/` 中的业务逻辑
- 格式化并输出结果

### 2.1 顶层文件与子包

| 层级 | 内容 | 示例 |
|------|------|------|
| `internal/cmd/*.go` | 顶层命令入口（简单命令） | `root.go`, `version.go`, `cleanup.go` |
| `internal/cmd/<command-group>/` | 复杂命令的子包（多文件拆分） | `task/`, `worktree/`, `feature/`, `fact/`, `prompt/`, `forensic/` |
| `internal/cmd/base/` | 命令共享的基础设施 | `output.go`, `errors.go`, `claude.go` |

当前状态：15 个顶层命令文件 + 8 个子包（`base`, `docs`, `fact`, `feature`, `forensic`, `prompt`, `task`, `worktree`）。

### 2.2 何时创建子包

- **单文件命令**（< 200 行，无子命令）：放在 `internal/cmd/<command>.go`
- **多子命令或超过 200 行**：创建 `internal/cmd/<command-group>/` 子包，至少包含：
  - `register.go` -- 注册 cobra 命令
  - `<subcommand>.go` -- 每个子命令一个文件
  - `<command_group>_test.go` / `testmain_test.go` -- 测试

### 2.3 cmd/ 内部禁止事项

- 禁止在 `cmd/` 中编写业务逻辑 -- 业务逻辑属于 `pkg/`
- 禁止 `cmd/` 的子包之间互相导入 -- 它们是独立的命令组
- 禁止直接操作文件系统的原始路径 -- 使用 `pkg/feature` 或 `pkg/project` 提供的路径工具

## 3. pkg/ 三层模型

`pkg/` 内部划分为三个层级，依赖方向严格自上而下：

```
domain（领域层）
  ↓ 可以导入
infrastructure（基础设施层）
  ↓ 可以导入
leaf（叶子层）
```

**关键约束**：
- leaf 包不能导入任何其他 `pkg/` 包
- infrastructure 包只能导入 leaf 包（主要是 `pkg/types`）
- domain 包可以导入 infrastructure 和 leaf 包
- domain 包之间的横向依赖（domain → domain）应尽量避免；如需引入，必须在代码审查中确认合理性

### 3.1 Leaf 层

零内部 `pkg/` 依赖的独立包。提供原子化的工具能力。

| 包 | 职责 |
|---|------|
| `pkg/types` | 纯类型定义（Status, Priority, Surface 等枚举与结构体） |
| `pkg/git` | Git CLI 操作封装 |
| `pkg/index` | 索引文件管理（原子写入、文件锁） |
| `pkg/infocmd` | 命令元数据描述 |
| `pkg/just` | Just 任务运行器集成 |
| `pkg/project` | 项目根目录检测与标记文件 |
| `pkg/facttable` | 事实表数据结构 |
| `pkg/version` | 版本常量 |

### 3.2 Infrastructure 层

仅依赖 leaf 包（主要是 `pkg/types`），提供跨领域共享的基础设施。

| 包 | 依赖 | 职责 |
|---|------|------|
| `pkg/forgeconfig` | `pkg/types` | Forge 配置文件解析、Surface 检测、执行顺序计算 |

### 3.3 Domain 层

实现核心业务逻辑，可依赖 infrastructure 和 leaf 包。

| 包 | 依赖 | 职责 |
|---|------|------|
| `pkg/feature` | `git`, `index`, `types` | Feature 生命周期管理（创建、状态追踪、路径计算） |
| `pkg/task` | `forgeconfig`, `index`, `infocmd`, `types` | 任务增删改查、状态机、依赖排序、模板生成 |
| `pkg/prompt` | `feature`, `forgeconfig`, `task` | Prompt 模板管理与渲染 |
| `pkg/proposal` | `feature`, `infocmd` | Proposal 文档管理 |
| `pkg/lesson` | `infocmd` | 经验教训记录 |
| `pkg/research` | `infocmd` | 深度研究任务管理 |
| `pkg/serverprobe` | `feature`, `just` | 服务器探测与 Just 集成 |
| `pkg/testrunner` | `just` | 测试运行器 |

## 4. 偏差分析表

以下记录当前代码与目标状态之间的偏差。数据来源：[pkg-dependency-graph.md](../../features/forge-cli-codebase-standards/pkg-dependency-graph.md)。

| 编号 | 偏差项 | 当前状态 | 目标状态 | 差距描述 |
|:----:|--------|----------|----------|----------|
| D1 | `pkg/infocmd` 高扇入 | 被 4 个 pkg/ 包导入（`lesson`, `proposal`, `research`, `task`），是最被依赖的 leaf 包 | 扇入可控，接口稳定 | 高扇入意味着变更影响面大；需确保其 API 稳定或拆分职责 |
| D2 | `pkg/feature` domain 中心枢纽 | 被 3 个 domain 包导入（`prompt`, `proposal`, `serverprobe`），同时自身依赖 3 个 leaf 包 | 减少被依赖面或明确为 infrastructure 层 | 既是最被依赖的 domain 包，又有多依赖；可考虑提取稳定接口降低耦合 |
| D3 | domain 横向依赖 | 4 条 domain → domain 边（`prompt→feature`, `prompt→task`, `proposal→feature`, `serverprobe→feature`） | 最小化横向依赖 | 横向依赖增加了测试复杂度和变更传播风险 |
| D4 | `pkg/task` 高耦合 | 导入 4 个内部 pkg/ 包（`forgeconfig`, `index`, `infocmd`, `types`） | 减少直接依赖数量 | 最宽的依赖面；考虑通过依赖注入或接口抽象降低直接导入 |
| D5 | `pkg/prompt` 深链 | 依赖深度 3 层（`prompt→task→forgeconfig→types`） | 依赖深度不超过 2 层 | 最深的依赖链增加了理解和测试的难度 |
| D6 | 孤立包未验证归属 | `pkg/facttable`, `pkg/project`, `pkg/version` 零内部导入 | 确认 pkg/ 层级归属是否合理 | 无内部消费者，需评估是否应保留在 pkg/ 或降级为 internal |
| D7 | `internal/cmd/docs` 空子包 | `cmd/docs/features/` 目录树存在但无文件 | 不创建空包目录 | 空目录违反最小化原则；应在有实际代码时再创建 |

## 5. 开发者工作流

### 5.1 新增命令

```
1. 确定命令复杂度：
   - 简单命令（< 200 行，无子命令）→ 创建 internal/cmd/<command>.go
   - 复杂命令（有子命令或预计超过 200 行）→ 创建 internal/cmd/<command-group>/

2. 如果创建子包，文件结构：
   internal/cmd/<command-group>/
   ├── register.go          # cobra 命令注册
   ├── <subcommand1>.go     # 子命令实现
   ├── <subcommand2>.go
   └── <command_group>_test.go

3. 业务逻辑放在 pkg/ 的对应包中，cmd/ 只做参数解析和调用
```

### 5.2 新增 pkg/ 包

```
1. 确定层级：
   - 纯数据/工具，无 pkg/ 依赖 → leaf
   - 依赖 pkg/types，提供共享能力 → infrastructure
   - 实现业务逻辑，依赖多个 pkg/ 包 → domain

2. 遵守层级约束：
   - leaf 包的 import 中不能出现 "forge-cli/pkg/" 的任何其他包
   - infrastructure 包只能导入 leaf 包
   - domain 包导入其他 domain 包需在 PR 中说明理由

3. 包命名使用小写单数形式：pkg/task, pkg/feature（非 tasks, features）
```

### 5.3 重构检查

当修改涉及包间依赖时：
1. 确认依赖方向：`cmd → internal → pkg`，不反向
2. 确认 pkg/ 内三层约束：`domain → infrastructure → leaf`，不反向
3. 确认不引入新的 domain 横向依赖
4. 运行 `go vet ./...` 确保无循环依赖

## 6. PR Review Checklist

复制以下清单到 PR 描述中，逐项确认：

```markdown
## Package Structure Review

- [ ] **依赖方向**：新增导入符合 `cmd → internal → pkg` 单向规则，无反向依赖
- [ ] **pkg/ 三层约束**：
  - [ ] leaf 包未导入其他 `pkg/` 包
  - [ ] infrastructure 包仅导入 leaf 包
  - [ ] domain 包的横向依赖（domain → domain）已在 PR 中说明合理性
- [ ] **包归属**：新代码放在正确的层级（cmd 做入口、pkg 做业务逻辑）
- [ ] **命名规范**：包名使用小写单数形式
- [ ] **子包决策**：如新增 `internal/cmd/` 下的子包，确认命令复杂度需要拆分（> 200 行或有子命令）
- [ ] **扇入影响**：修改 `pkg/infocmd`、`pkg/feature`、`pkg/types` 等高扇入包时，确认对所有消费者的影响
- [ ] **循环依赖**：运行 `go vet ./...` 通过，无循环导入
```

## 7. 参考

- [pkg-dependency-graph.md](../../features/forge-cli-codebase-standards/pkg-dependency-graph.md) -- pkg/ 依赖图的完整事实基线
- [code-structure.md](./code-structure.md) -- 代码结构规范（嵌套、缩进、控制流）
- [forge-cli-reference.md](./forge-cli-reference.md) -- forge-cli 开发参考
