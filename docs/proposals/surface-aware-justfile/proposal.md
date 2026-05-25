---
created: 2026-05-24
author: "faner"
status: Approved
---

# 提案：init-justfile Surface 感知 + 测试编排简化

## 问题

两个相互关联的问题：

### 问题 1：init-justfile 不感知 Surface

init-justfile 仅根据项目语言生成 just 配方，忽略了 surface 类型。但不同 surface 的**测试编排流程根本不同**：

- **Web/API**：须先启动应用 → 等待就绪 → 运行测试 → 关闭
- **CLI/TUI**：直接构建并测试，无需启动服务
- **Mobile**：启动模拟器 → 运行测试 → 清理

### 问题 2：test.execution 委托层冗余

当前 `run-tests` 通过 `.forge/config.yaml` 的 `test.execution.run` 读取命令模板（如 `"just test {slug}"`），再执行。这形成了层层委托：

```
run-tests → config.yaml test.execution.run → "just test {slug}" → just test → 实际测试运行器
```

justfile 本身已经是抽象层，config.yaml 再包一层增加了复杂度但没有增加灵活性 — 绝大多数模板变量最终都解析为 just 命令。少数非 just 路径（如 `go test`、`npx vitest`）可以通过封装到 justfile 配方中来保留。

### 证据

- Web UI 的 e2e 测试必须先启动应用，但当前配方没有 surface 特定的启动逻辑
- `run-tests` 的编排序列对 surface 不感知，只能依赖用户手动配置 `test.execution`
- `test.execution` 的多数示例指向 just 命令（`just test {slug}`、`just test-setup`、`just probe`），但 config-schema.md 也记录了非 just 示例（`go test -json -v ./...`、`npx vitest run --reporter=json`、`make test FEATURE={slug}`）。非 just 路径在简化方案中被牺牲——对于直接使用 `go test` 或 `npx vitest` 的项目，用户需要将这些命令封装到 justfile 的 `test` 配方中。这是可接受的 trade-off，因为 justfile 本身就是命令抽象层

### 紧迫性

随着 v3.0.0 test profile 的引入，surface 成为测试流程的核心维度。init-justfile 和 run-tests 需要协同工作 — init-justfile 生成正确的配方，run-tests 编排正确的执行序列。

## 提议方案

### 方案 A：Surface 感知 + 委托层简化

1. init-justfile 添加 surface 感知层，为不同 surface 生成差异化的配方
2. 废弃 `test.execution` 委托层（保留废弃检测，但不再使用其值），`run-tests` 直接调用 just 配方
3. `timeout`、`results-dir` 等非命令配置保留在 config.yaml 的简化字段中

### Surface 测试编排模式

| Surface | 测试编排序列 | 关键配方 |
|---------|------------|---------|
| **web** | `just dev`（后台）→ `just probe`（等待就绪）→ `just test` → teardown | `dev` 启动 dev server；`probe` 检查 HTTP 端点 |
| **api** | `just dev`（后台）→ `just probe`（等待就绪）→ `just test` → teardown | `dev` 启动 API 服务；`probe` 检查 /healthz |
| **cli** | `just build` → `just dev`（验证二进制可运行）→ `just test` | 无需 `run`/`probe`；`test` 通过子进程测试 |
| **tui** | `just build` → `just dev`（验证二进制可运行）→ `just test` | 无需 `run`/`probe`；`test` 通过 stdin 管道测试 |
| **mobile** | `just dev`（启动模拟器+应用）→ `just test` → teardown | `test-setup` 准备模拟器；`test` 用 maestro |

> **api/web 统一性说明**：api 和 web 的编排序列完全相同（dev → probe → test → teardown），唯一差异在于 `probe` 检查的目标端点（api 检查 `/healthz`，web 检查页面根路径）。当前选择保持两者为独立规则，原因是：(1) probe 端点差异可能在后续迭代中扩展为更大的行为差异（如 web 可能需要检查静态资源、api 可能需要检查数据库连接）；(2) 用户体验上，开发者通常以 surface 类型作为心智模型，合并为 `service` 会增加认知负担。若后续验证两者确实无实质性差异，可合并为 `service` 规则并共享编排模板。

**dev/run 分工**：
- **web/api**：保留 `dev`（热重载 dev server）和 `run`（生产模式启动），测试编排使用 `dev`
- **cli/tui**：只生成 `dev`（编译+运行二进制），不生成 `run`
- **mobile**：只生成 `dev`（启动模拟器+应用），不生成 `run`
- **混合项目（web+api）**：`dev` 配方接受 scope 参数 — `just dev frontend`、`just dev backend`、`just dev`（无 scope 时按依赖顺序启动所有 scope）。**scope 值必须是 `config.yaml` 中 `surfaces` map 的 key**（即用户定义的目录路径或逻辑名），而非 surface 类型枚举。例如，若 `surfaces: {admin-panel: web, payment-service: api}`，则 scope 值为 `admin-panel` / `payment-service`，而非 `frontend` / `backend`

**混合项目多服务启动管理**：

当 `just dev`（无 scope）需要启动多个服务时：

- **端口冲突预防**：`just dev` 配方体在启动前检查端口是否已被占用（`lsof -i :$PORT` / `netstat -ano | findstr $PORT`），若占用则输出明确错误并退出
- **顺序启动策略**：`just dev` 配方体按顺序启动各 scope 的 dev server（先启动后端，再启动前端），每个启动后将 PID 写入 `.forge/dev-server.<scope>.pid`。顺序启动（而非并行）是为了避免多进程同时争抢系统资源导致启动失败
- **probe 顺序**：run-tests 先 probe 后端（api），再 probe 前端（web）。后端就绪是前端可用的前提条件
- **teardown 逆序清理**：先 teardown 前端，再 teardown 后端。逆序清理模拟生产环境的依赖关系

### 委托层简化

废弃 `test.execution` 委托层（不再读取和使用其命令值），`run-tests` 的编排变为：

```
当前（4 层委托）：
run-tests → config test.execution.run → "just test {slug}" → just test → 实际运行器

简化后（2 层）：
run-tests → just test [journey] → 实际运行器
```

`run-tests` 根据 surface 编排模式直接调用 just 配方，不再读 config 命令模板：

| 编排步骤 | 当前 | 简化后 |
|---------|------|--------|
| setup | `test.execution.setup` 配置 | `just test-setup` |
| pre-check | `test.execution.pre-check` 配置 | `just probe`（web/api）或跳过（cli/tui） |
| run | `test.execution.run` 配置 | `just test [journey]` |
| teardown | `test.execution.teardown` 配置 | `just test-teardown`（可选配方） |

保留在 config.yaml 中的配置（非命令类）：

```yaml
test:
  timeout: 300
  results-dir: "tests/{journey}/results"
```

**`test.execution` 废弃行为**：

`test.execution` 节点从**执行路径**中移除（不再读取其命令值），但保留**检测路径**（检查节点是否存在以输出迁移警告）。这不是矛盾——"不使用"不等于"不检测"。具体规则：

- **检测时机**：`run-tests` 启动时，在读取 config 后立即检查 `test.execution` 节点是否存在
- **警告格式**：若 `test.execution` 存在，输出标准警告（`⚠ DEPRECATED: test.execution is no longer used. Surface-aware test orchestration is now automatic. Remove test.execution from your config.yaml.`）
- **行为**：警告后继续执行（不中断），使用 surface 感知编排模式替代。`test.execution` 中的命令值被完全忽略，不会被读取或执行
- **文档同步**：在 config-schema.md 中标记 `test.execution` 为 `@deprecated (v3.0.0)`，附迁移说明

### Skill 流程

init-justfile 的生成流程：
```
检测语言 → 加载语言模板 → 生成 compile/build/lint/fmt
检测 surface → 加载 surface 规则 → 生成 test/dev/run/probe/test-setup
组装 → 验证
```

**test 配方生成 fallback 链**（替代被移除的 `test.execution.run`）：

当 init-justfile 为某个目录生成 `test` 配方时，按以下优先级查找指导规则，命中即停止：

1. **Surface 规则**：若该目录有 surface 类型（web/api/cli/tui/mobile），加载对应 `rules/surfaces/<type>.md`，按其中定义的编排模式生成 `test` 及配套配方（dev/probe/test-setup 等）
2. **Convention 框架**：若该目录无 surface 类型但有 convention 框架（如 `rules/test-python.md`），使用 convention 驱动的测试生成（当前行为）
3. **语言模板 cold start**：若以上均缺失，从语言模板（`templates/<lang>.just`）中提取 `test` 配方骨架，生成最小可运行配方
4. **报错提示**：若语言模板也无 `test` 配方定义，输出明确的提示信息，告知用户需手动编写 `test` 配方或配置 surface 类型

**语言模板 vs Surface 规则仲裁规则**：

init-justfile 生成流程中，语言模板和 surface 规则可能对同一配方的生成逻辑产生冲突。仲裁规则如下：

- **原则：Surface 规则优先覆盖 `test`/`dev`/`run`/`probe` 配方**
- **语言模板负责**：`compile`、`build`、`lint`、`fmt`、`unit-test`（语言级配方）
- **Surface 规则负责**：`test`、`dev`、`run`、`probe`、`test-setup`、`test-teardown`（编排级配方）
- **冲突处理**：当两者试图生成同名配方时，Surface 规则覆盖语言模板。例如，语言模板生成的 `test` 配方会被 surface 规则定义的编排模式替换
- **组装策略**：init-justfile 先加载语言模板生成基础配方，再加载 surface 规则覆盖编排配方，最终组装为完整的 justfile

run-tests 的执行流程：
```
读取 .forge/surface-orchestration.yaml → 获取编排模式
→ just test-setup
→ just dev（web/api，后台启动，按 startup_order 顺序）
→ just probe（web/api，等待就绪，按 startup_order 顺序）
→ just test [journey]
→ just test-teardown（如有，逆序清理）
```

**与现有 surface 感知环境检查的关系**：

当前 `run-tests` SKILL.md 第 4 步的"Environment Readiness Check"已通过 `rules/env-check.md` 和 `surface-<type>.md` 做 surface 感知的环境检查（如检查 playwright 是否安装、端口是否可用）。提案与现有机制**互补而非重叠**：

- **现有机制（Environment Readiness Check）**：关注**前置条件检查** — 测试工具是否安装、依赖是否就绪、端口是否被占用。在编排执行之前运行
- **提案新增（Surface 感知编排）**：关注**执行序列编排** — 启动/停止服务的顺序、何时 probe、何时 teardown。在环境检查通过后运行

两者协作关系：`env-check` 确认环境可用 → 提案的编排模式按 surface 类型执行正确的序列。现有 `surface-<type>.md` 规则文件可以继续用于环境检查，无需修改。

### 创新亮点

将 init-justfile（配方生产者）和 run-tests（配方消费者）的 surface 感知统一设计。init-justfile 生成正确的配方组合，run-tests 按 surface 编排执行序列，废弃中间的 config 委托层。justfile 成为唯一的命令抽象层。

**统一的具体机制**：init-justfile 在生成 justfile 的同时，写入 `.forge/surface-orchestration.yaml` 文件，声明当前项目的编排模式。run-tests 读取该文件驱动执行序列，而非自行推断编排模式。这确保两个 skill 通过共享文件而非隐式约定保持一致。

```yaml
# .forge/surface-orchestration.yaml（由 init-justfile 生成，run-tests 消费）
version: 1
surfaces:
  admin-panel:
    type: web
    dev_recipe: "dev admin-panel"
    probe_target: "http://localhost:3000"
    teardown_order: 2
  payment-service:
    type: api
    dev_recipe: "dev payment-service"
    probe_target: "http://localhost:8080/healthz"
    teardown_order: 1  # 逆序：先清理后端
orchestration:
  startup_order: [payment-service, admin-panel]  # 按依赖顺序
  probe_after: dev
  teardown_reverse: true
```

**文件职责**：
- **init-justfile**（写入方）：根据 surface 检测结果和用户 config 生成此文件，包含每个 scope 的编排参数
- **run-tests**（读取方）：解析此文件获取编排序列，按 `startup_order` 顺序启动，按 `teardown_order` 逆序清理
- **用户**（可选修改）：高级用户可直接编辑此文件调整编排行为（如更改启动顺序、添加 probe 超时参数），无需修改 SKILL 或 justfile

### justfile 作为唯一抽象层的 trade-off 分析

提案将命令可配置性从 `config.yaml` 转移到 justfile 配方体中。这意味着 `run-tests` 硬编码调用 `just dev`/`just probe`/`just test` 等固定配方名，而非从配置中读取命令模板。

**选择此路径的理由**：
- **显式性优于隐式性**：justfile 配方是可见、可审计的；config.yaml 中的命令模板是隐式间接层，开发者难以追溯执行链路
- **justfile 已经是抽象层**：`test.execution` 的多数实际示例指向 just 命令，config 再包一层主要是转发。少数非 just 路径（`go test`、`npx vitest`、`make test`）可封装到 justfile 配方中保留，间接层增加了复杂度但未增加表达能力
- **可定制性保留在配方体中**：用户仍可通过编辑 justfile 配方体来定制行为（如添加环境变量、修改启动参数），只是不再通过 config.yaml 间接定制

**已知局限**：
- 若需要在不修改 justfile 的情况下切换编排策略（如 CI 环境用不同启动命令），当前方案无法支持。缓解措施：justfile 可通过环境变量参数化（`env_var := env("FORGE_DEV_CMD", "npm run dev")`）
- 若未来新增 surface 类型，`run-tests` 需要更新以识别新的编排模式。这是一个低频操作（每增加一种 surface 类型才需更新），且通过 surface 规则文件的定义可提前约束

## 需求分析

### 关键场景

- **Web+Node 项目**：`dev` 启动 next dev（端口 3000）→ `probe` curl localhost:3000 → `test` playwright e2e → teardown
- **API+Go 项目**：`dev` 启动 API server → `probe` curl /healthz → `test` HTTP 契约测试 → teardown
- **CLI+Rust 项目**：`build` 编译二进制 → `dev` 运行验证 → `test` 子进程集成测试
- **TUI+Go 项目**：`build` 编译二进制 → `dev` 运行验证 → `test` stdin 管道测试
- **Mobile+TS 项目**：`test-setup` 准备模拟器 → `dev` 启动应用 → `test` maestro YAML
- **无 surface 配置**：回退到当前行为（纯语言配方，run-tests 保持原有逻辑）
- **混合项目（web+api）**：`just dev` 无 scope 时按依赖顺序启动前端和后端（例如 `surfaces: {admin-panel: web, payment-service: api}` 时，先 `just dev payment-service` 再 `just dev admin-panel`），probe 分别检查两个服务

### 下游集成契约

配方签名不可变。`run-tests` 的调用方式简化但语义不变：

| 配方 | 签名（不可变） | 消费者 | 期望语义 |
|------|--------------|--------|---------|
| `unit-test` | `just unit-test` | forge task submit、clean-code、fix-bug、testrunner | 语言级单元测试；exit 0 = 通过 |
| `test` | `just test [journey]` | run-tests、forge quality-gate、fix-bug | Surface 级测试 |
| `probe` | `just probe` | run-tests（web/api 前置检查） | 服务健康检查；exit 0 = 健康；混合项目检查所有服务 |
| `test-setup` | `just test-setup` | run-tests | 安装测试依赖；幂等 |
| `test-teardown` | `just test-teardown` | run-tests（可选） | 测试后清理 |
| `dev` | `just dev [scope]` | run-tests（web/api 启动服务） | web/api: 后台启动并监听端口；cli/tui: 编译运行；`scope` 为可选参数，值为 `surfaces` map 的 key，仅混合项目需要 |
| `run` | `just run [scope]` | 仅 web/api 生成 | 生产模式启动 |

**关键约束**：
- `test` 必须始终接受 `journey=''` 参数
- web/api 的 `dev` 须能在后台运行
- cli/tui 不生成 `run` 配方
- `test-teardown` 为可选配方，不存在时 run-tests 跳过

### 约束与依赖

- Surface 信息来自 `.forge/config.yaml` 或 `forge surfaces` CLI
- **编排声明文件**：`.forge/surface-orchestration.yaml` 由 init-justfile 生成，run-tests 读取。该文件是 init-justfile 和 run-tests 的统一接口——init-justfile 写入编排参数，run-tests 按参数驱动执行
- **Surface 信息源优先级规则**（init-justfile 和 run-tests 统一遵循）：
  1. **`config.yaml` 的 `surfaces` 字段优先**：若 `surfaces` 字段存在且非空，以其中定义的 surface 类型和 scope 映射为准
  2. **`forge surfaces` CLI 回退**：若 `surfaces` 字段缺失或为空，通过 `forge surfaces [path]` CLI 基于文件信号检测获取 surface 类型
  3. **冲突处理**：当两者不一致时（config 说 web，文件信号检测出 api），以 `config.yaml` 为准（用户显式配置优先于自动检测）
- Surface 规则保持语言无关
- Standard Target Contract 配方名称和签名不变
- 遵循 forge-distribution.md 路径约定

### 非功能需求

| NFR | 要求 | 验证方式 |
|-----|------|---------|
| **跨平台兼容** | Windows/macOS/Linux 三平台均可运行编排序列。`just dev`/`probe`/`test-teardown` 配方体需包含平台分支 | 各平台手动验证；CI 矩阵（如果接入） |
| **向后兼容** | 无 surface 配置的项目行为与当前完全一致 | diff 输出对比 |
| **可观测性** | 编排过程中每个步骤（dev 启动、probe 轮询、test 执行、teardown 清理）输出结构化日志，用户可追踪执行进度 | 日志输出检查 |
| **性能** | surface 规则文件加载不增加 init-justfile 超过 1 秒的额外耗时 | 计时基准测试 |
| **可靠性** | dev server 崩溃时 probe 超时后执行 teardown，不遗留孤儿进程；会话中断后可通过 test-state.json 恢复清理 | 故障注入测试 |
| **just 版本** | just >= 1.0（支持 `[linux]`/`[windows]` recipe attribute 进行平台分支）；无需后台运行相关的特殊版本 | 版本检查 |

## 替代方案

| 方案 | 优势 | 劣势 | 结论 |
|------|------|------|------|
| 不做 | 零成本 | 下游无法 surface 感知编排；test.execution 委托层继续累积复杂度 | 拒绝 |
| 仅 surface 感知，保留 test.execution | 改动小 | 委托层冗余未解决；两处配置 surface 行为 | 拒绝：治标不治本 |
| **surface 感知 + 废弃 test.execution** | 统一抽象层；消除冗余委托 | 改动范围扩大到 run-tests；非 just 路径需封装到 justfile | **选定** |
| Go 代码直接管理进程生命周期 | 确定性最高：进程启动、PID 追踪、信号发送、超时控制全部由 Go 代码实现，不依赖 LLM 执行可靠性 | 需要 Go CLI 新增子命令（`forge test dev`/`probe`/`teardown`），开发成本高；justfile 退化为纯测试命令容器，丧失 just 生态的灵活性（用户无法通过编辑 justfile 定制启动行为）；与 Forge "justfile 作为唯一抽象层"的设计哲学冲突 | 拒绝（v3.0.0 范围过大），但采纳其核心思想作为兜底机制（见"LLM agent 执行确定性"中的分层防御策略） |

## 行业对标

"启动服务 → 等待就绪 → 运行测试 → 清理"是通用的编排需求。以下成熟方案的对比帮助定位 Forge 方案的设计空间：

### 成熟方案对比

| 方案 | 编排模型 | 就绪检测 | 进程管理 | 适用场景 |
|------|---------|---------|---------|---------|
| **Docker Compose** `healthcheck` + `depends_on` | 声明式（YAML） | 容器内 HTTP/TCP 探针，支持 `interval`/`timeout`/`retries` | 容器生命周期（Docker daemon 托管） | 多容器微服务 |
| **Kubernetes** `readinessProbe` + `initContainers` | 声明式（YAML） | HTTP GET / TCP / Exec 三种探针类型，`periodSeconds`/`failureThreshold` | Pod 生命周期（kubelet 托管） | 云原生编排 |
| **Cypress** `start-server-and-test` | 命令式（CLI） | 轮询 URL 直到 HTTP 200，可配置超时 | 自动管理 dev server PID，测试后清理 | 前端 E2E 测试 |
| **Makefile** target dependency | 声明式（依赖图） | 无内建探针；依赖 `make -j` 并行 + `.WAIT` 串行控制 | 无进程管理 | 通用构建编排 |
| **GitHub Actions** `service containers` | 声明式（YAML） | 自动健康检查，基于 Docker health status | Runner 自动启停容器 | CI 服务依赖 |

### Forge 方案的定位差异

以上方案的共同特点：**编排逻辑由确定性代码执行**（Docker daemon、kubelet、Node.js 进程）。Forge 的关键差异在于：

- **run-tests 是 LLM agent 执行的 SKILL**，编排序列由 LLM 按步骤执行，而非确定性代码
- **justfile 是文本协议**，不是运行时——just 配方启动的进程在 just 退出后可能成为孤儿进程

这意味着 Forge 不能直接复用上述方案的进程管理能力，需要在 SKILL 指导层面建立可靠性保证（见"LLM agent 执行确定性"和"后台进程管理"两节）。

### 从行业方案借鉴的设计

| 借鉴点 | 来源 | Forge 实现 |
|--------|------|-----------|
| 探针重试 + 超时 | K8s readinessProbe | `just probe` 配方内实现重试循环（见"技术方向深化"） |
| 测试后强制清理 | Cypress start-server-and-test | `just test-teardown` + `.forge/test-state.json` 状态恢复 |
| 声明式编排序列 | Docker Compose depends_on | Surface 规则文件声明编排步骤序列 |

## 可行性评估

### 后台进程管理

just 本身不原生支持后台运行模式。web/api 编排中的 `just dev`（后台启动 dev server）需要解决三个问题：

**1. 后台启动**

just 配方体内的 `&` 后台操作符在 Unix 上可用，但 Windows 上不可用。方案：

- **选定方案**：使用 just 原生 `[linux]`/`[windows]` recipe attribute 实现平台分支。为每个需要跨平台行为的配方生成两个变体，just 根据运行平台自动选择匹配的配方。无需 shebang，无外部依赖（不依赖 bash 可用性）
- **替代**：使用 just 的 `python3 -c '...'` 内联脚本调用 `subprocess.Popen`，跨平台兼容但依赖 Python

**just 原生平台 attribute 示例**：

```just
[linux]
dev:
    nohup npm run dev > /dev/null 2>&1 & echo $! > .forge/dev-server.pid

[windows]
dev:
    start /B npm run dev > NUL 2>&1 & echo %PID% > .forge\\dev-server.pid
```

**优势**：just 原生属性由 just 解释器在配方选择阶段处理，不依赖 bash、Python 等外部运行时。Windows 上直接使用 CMD 语法，无兼容性问题

**2. PID 追踪**

后台进程的 PID 需要持久化以供 teardown 使用：

- `just dev` 将 PID 写入 `.forge/dev-server.pid`（约定路径）
- `just probe` 不需要 PID——它只检查 HTTP 端点可达性
- `just test-teardown` 读取 `.forge/dev-server.pid`，发送 SIGTERM（Unix）或 taskkill（Windows）
- run-tests 在 `.forge/test-state.json` 中记录 teardown 命令（现有机制），同时记录 PID 文件路径用于异常恢复

**3. Windows 兼容性**

Windows 上存在以下差异：
- 无 `SIGTERM`/`SIGKILL` 信号——需要使用 `taskkill /PID <pid> /T`（/T 杀进程树）
- 无 `nohup`——需要使用 `start /B` 或 PowerShell `Start-Process`
- `curl` 可能不可用——`just probe` 需要考虑使用 PowerShell `Invoke-WebRequest` 作为 fallback

**缓解方案**：`just dev`/`just probe`/`just test-teardown` 配方体通过 just 原生 `[linux]`/`[windows]` recipe attribute 分别定义平台分支，无需 shebang 或平台检测脚本。init-justfile 生成的配方模板为每个跨平台配方生成两个变体（`[linux]` 和 `[windows]`），just 自动选择。

### 技术方向深化

#### Probe 轮询逻辑

`just probe` 不是单次 HTTP 请求——dev server 启动需要时间。probe 配方必须实现重试循环：

```
# just probe 配方骨架（伪代码）
max_retries=30        # 最多重试 30 次
interval=2            # 每次间隔 2 秒
total_timeout=60      # 总超时 60 秒

for i in $(seq 1 $max_retries); do
  if curl -sf http://localhost:$PORT/healthz; then
    exit 0  # 就绪
  fi
  sleep $interval
done
exit 1  # 超时未就绪
```

**参数配置来源**：
- `max_retries` 和 `interval` 通过环境变量覆盖（`PROBE_RETRIES`、`PROBE_INTERVAL`），默认值在 surface 规则文件中定义
- web surface 默认：30 次重试 / 2 秒间隔 / 检查根路径
- api surface 默认：30 次重试 / 2 秒间隔 / 检查 `/healthz`

**超时失败行为**：`just probe` 以 exit 1 退出，run-tests 捕获后：(1) 执行 `just test-teardown` 清理已启动的 dev server；(2) 报告"服务启动超时"错误；(3) 不执行 `just test`。

#### Teardown 进程回收机制

teardown 必须可靠地回收所有测试启动的进程，即使测试中途崩溃：

1. **正常路径**：`just test-teardown` 读取 `.forge/dev-server.pid`，**先校验 PID 有效性**，再向对应进程发送 SIGTERM/taskkill
2. **PID 有效性校验**（防止陈旧 PID 导致杀错进程）：
   - Linux/macOS：检查 `/proc/<pid>` 是否存在（Linux）或 `ps -p <pid>` 是否成功（macOS）
   - Windows：`tasklist /FI "PID eq <pid>" /NH` 检查进程是否存在
   - 若 PID 无效（进程已不存在），跳过终止步骤，仅清理 PID 文件并输出警告
3. **异常恢复路径**：如果 run-tests 会话被中断（如用户 Ctrl+C），`.forge/test-state.json` 中的 teardown 命令在下次 run-tests 启动时执行（现有机制）
4. **孤儿进程兜底**：`just test-teardown` 在发送 SIGTERM 后等待 5 秒，如果进程仍存活则发送 SIGKILL/taskkill /F
5. **PID 文件清理**：teardown 完成后（无论 PID 有效与否）删除 `.forge/dev-server.pid`

**跨平台信号映射**：

| 操作 | Unix | Windows |
|------|------|---------|
| 优雅终止 | `kill -TERM $PID` | `taskkill /PID $PID` |
| 强制终止 | `kill -KILL $PID` | `taskkill /PID $PID /F /T` |
| 等待退出 | `wait $PID` 或 `timeout` | `timeout /t 5` |

### LLM agent 执行确定性

run-tests 是由 LLM agent 执行的 SKILL，不是确定性代码。编排序列的可靠性面临以下挑战：

1. **Surface 类型误判**：LLM 可能错误识别 surface 类型，导致使用错误的编排模式
2. **步骤跳过**：LLM 可能在 probe 未通过时直接执行 test
3. **时序错误**：LLM 可能在 dev 未完全启动时就执行 probe

**缓解措施（分层防御）**：

- **v3.0.0 兜底机制（确定性下限）**：run-tests 从 `.forge/surface-orchestration.yaml` 读取完整编排参数（启动顺序、probe 目标、teardown 逆序），SKILL.md 中将编排序列定义为参数化的固定模板，LLM 填充模板参数而非自由编排。关键约束通过 just 命令退出码强制执行——每个步骤的退出码决定下一步动作（0=继续，非0=触发 teardown 后中止）。即使 LLM 误读指令，just 命令的退出码是确定性的，不会因 LLM 行为改变
- **HARD-GATE 规则**：在 run-tests SKILL.md 中用 `<HARD-GATE>` 标记编排序列的强制顺序。违反顺序的行为列为"禁止"，与现有"禁止跳过失败测试"同等级别
- **每步退出码检查**：每个编排步骤（dev/probe/test/teardown）必须检查前一步的退出码。非零退出码触发 abort 或 teardown，不继续后续步骤
- **状态机驱动**：run-tests 的编排步骤本质上是状态机（init → dev → probe → test → teardown）。SKILL.md 中显式声明当前步骤和下一步骤的映射关系，LLM 按状态机转移而非自由选择
- **长期方向**：将编排逻辑迁移到 Go CLI 子命令（如 `forge test dev`/`probe`/`teardown`），Go 代码直接管理进程生命周期（后台启动、PID 校验、信号发送、超时控制），SKILL 只负责调用 Go 命令和读结果。但这不是 v3.0.0 的范围——v3.0.0 通过参数化模板 + 退出码约束提供确定性下限

### 技术可行性

直接可行。init-justfile 已有 convention 驱动的测试生成。Surface 检测增加一次 config 读取。`run-tests` 简化为直接调用 just 配方。

### 资源与时间

中等范围：
- 5 个 surface 规则文件
- init-justfile SKILL.md 更新
- run-tests SKILL.md 简化（去掉 test.execution 读取）
- config-schema.md 更新
- init-justfile Step 3a 去掉 test.execution 依赖
- **config schema 变更**（独立子方案，2-3 个任务）：结构体变更、GetConfigValue 扩展、文档更新

预计 10-15 个编码任务。

### 依赖就绪度

Surface 检测已就位。`test.execution` 在 Go `Config` 结构体中未被映射为独立字段（代码验证：结构体无 `Execution` 子字段），但 run-tests SKILL.md 通过 `forge config get test.execution` 让 LLM agent 读取原始 YAML 并使用这些字段——即 `test.execution` 在 Go 层面未结构化实现，但在 LLM agent 层面实际在用。移除后需要确保 run-tests 的所有编排路径都更新为 surface 感知模式，同时废弃检测覆盖 config-schema.md 的文档更新。`GetConfigValue` 扩展为 config schema 子方案的一部分，需独立评审。

## 假设挑战

| 假设 | 挑战工具 | 发现 |
|------|---------|------|
| `test.execution` 提供灵活性 | XY 检测 | 部分推翻：多数示例指向 just 命令，但 config-schema.md 中确实存在 `go test`、`npx vitest`、`make test` 等非 just 示例。对于这些路径，简化方案要求用户将其封装到 justfile 中——灵活性从 config.yaml 转移到 justfile 配方体，不丢失但位置变了 |
| Surface 规则应包含语言特定指导 | 假设翻转 | 推翻：用户确认保持语言无关 |
| 只需优化单个配方 | 压力测试 | 推翻：核心是**编排序列**不同 |
| `dev`/`run` 只用于人工开发 | XY 检测 | 推翻：run-tests 编排测试时需要 `dev` 启动被测服务 |
| CLI/TUI 需要 `run` 配方 | Occam's Razor | 推翻：CLI/TUI 无服务启动概念，`dev`（编译+运行）足够 |

## 范围

### 范围内

**init-justfile**：
- 5 个 surface 规则文件：`skills/init-justfile/rules/surfaces/{web,api,cli,tui,mobile}.md`
  - 每个包含：测试编排模式、配方生成指导、journey 过滤策略
- **journey 过滤策略最小规范**（每个 surface 规则文件必须定义）：

  | Surface | 默认 journey 范围 | 过滤规则 | 示例 |
  |---------|-----------------|---------|------|
  | **web** | 所有 journeys | `just test` 运行全部；`just test smoke` 运行冒烟测试 | journey 标签映射：`smoke` → 登录+首页加载；`e2e` → 全流程 |
  | **api** | 所有 journeys | `just test` 运行全部；`just test contract` 运行契约测试 | journey 标签映射：`contract` → API 契约验证；`integration` → 数据库集成 |
  | **cli** | 所有 journeys | `just test` 运行全部；`just test unit` 运行子进程测试 | journey 标签映射：`unit` → 命令解析测试；`integration` → 文件系统交互 |
  | **tui** | 所有 journeys | `just test` 运行全部；`just test snapshot` 运行快照测试 | journey 标签映射：`snapshot` → UI 快照对比；`interaction` → 键盘输入测试 |
  | **mobile** | 所有 journeys | `just test` 运行全部；`just test e2e` 运行 maestro flow | journey 标签映射：`e2e` → maestro YAML 流程；`snapshot` → 屏幕截图对比 |

  journey 过滤通过 `test` 配方的 `[journey]` 参数实现：`just test [journey]` 中 `journey` 为空时运行全部，否则按 surface 规则中定义的标签映射过滤。具体过滤逻辑由 justfile 配方体中的 `if/else` 或 `test runner` 的标签机制实现。
- SKILL.md 更新：新增 surface 检测步骤，surface 感知配方生成
- CLI/TUI 只生成 `dev`，不生成 `run`
- 去掉 Step 3a 中对 `test.execution.run` 的依赖

**run-tests**：
- SKILL.md 简化：直接调用 just 配方，去掉 `test.execution` 读取
- 根据 surface 编排模式决定执行序列（是否启动 dev、是否 probe）
- 读取 `.forge/surface-orchestration.yaml` 获取编排序列（startup_order、probe 目标、teardown 逆序）

**配置 schema 变更（独立子方案）**：

此变更改变了 Forge CLI 的配置接口契约，影响面超出常规 Go 开发。单独列出以确保边界明确：

- **新增 `TestConfig` 节点**：Config 结构体新增 `test` 顶层节点，包含 `timeout`（int，秒）和 `results-dir`（string，模板路径）字段
- **移除 `test.execution` 节点**：从 Config 结构体中移除 `test.execution`，config-schema.md 中标记为 `@deprecated`（保留一个版本的废弃检测兼容期）
- **扩展 `GetConfigValue`**：支持 `test.*` 键读取（`forge config get test.timeout`、`forge config get test.results-dir`）
- **废弃检测**：run-tests 启动时检测到 `test.execution` 节点存在时，输出废弃警告（见下方"废弃行为"定义）
- **边界约束**：不修改其他 skill 对 `test.timeout` 等非命令字段的读取方式；`surfaces` 字段保持现有 schema 不变
- **影响面评估**：涉及 `forge-cli/internal/config/`（结构体定义）、`forge-cli/internal/cmd/config.go`（GetConfigValue 扩展）、`plugins/forge/references/run-tests/config-schema.md`（文档更新）3 个模块，预计 2-3 个独立任务

**通用**：
- 向后兼容：无 surface 配置 → 当前行为不变
- 混合项目：surface 规则按 scope 应用
- init-justfile 生成 `.forge/surface-orchestration.yaml` 编排声明文件（见"创新亮点"中的具体机制定义）

### 范围外

- 变更语言模板（`templates/*.just`）
- 变更 `forge-cli/pkg/just/` 门控序列
- 变更 `forge-cli/internal/cmd/quality_gate.go` 或 `testrunner` 的 Go 代码
- 新增 forge CLI 命令（仅扩展 `config get` 键支持）
- **回滚基础设施**（feature flag 机制不在范围内；回滚通过 git revert 实现）
- Go 代码子命令直接管理进程（长期方向，非 v3.0.0 范围）

## 主要风险

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|---------|
| Surface 未检测到 | 中 | 低 | 回退到当前行为 |
| run-tests 简化导致已配置 `test.execution` 的项目不兼容 | 低 | 低 | v3.0.0 未发布，无存量用户；且运行时检测到 `test.execution` 节点时输出废弃警告，引导用户移除 |
| Surface 规则过于泛化 | 中 | 中 | LLM 组合语言模板 + surface 规则 |
| 混合项目 surface 歧义 | 中 | 中 | `forge surfaces` 基于路径检测；按 scope 应用 |
| `test [journey]` 过滤与原生运行器不兼容 | 中 | 高 | Surface 规则记录映射关系 |
| run-tests 无法感知 surface（skill 内无 config 读取） | 低 | 高 | run-tests 通过 `forge surfaces` 或 `.forge/config.yaml` 获取 surface 类型 |

### 回滚计划

如果 surface 感知方案上线后发现严重问题（如某种 surface 类型的编排在特定框架下不可用），按以下策略回滚：

1. **回滚方式**：通过 `git revert` 回退提案的所有变更（SKILL.md、surface 规则文件、config-schema.md、init-justfile 模板）。不引入 feature flag 机制——回滚基础设施不在 v3.0.0 范围内
2. **配置兼容期**：`test.execution` 节点在 config-schema.md 中标记为 `@deprecated` 而非立即删除，保留一个版本（v3.1）的兼容期。废弃检测逻辑在兼容期结束后才移除
3. **回滚影响面**：回滚后 `run-tests` 恢复到当前行为（读取 `test.execution` 委托路径），已生成的 surface 感知 justfile 中 dev/probe/test 配方仍然可用（它们只是 just 配方，不依赖 run-tests 的编排逻辑）
4. **完全移除**：在兼容期（预计 v3.2）后，确认无回退需求时，移除 `test.execution` 的废弃检测代码

## 成功标准

- [ ] init-justfile 为 web/api/cli/tui/mobile 5 种 surface 生成差异化的配方组合
- [ ] CLI/TUI 不生成 `run` 配方，统一使用 `dev`
- [ ] `run-tests` 不再依赖 `test.execution.run`，直接调用 `just test [journey]`
- [ ] `run-tests` 根据 surface 编排模式决定是否启动 dev/probe
- [ ] `test.execution` 节点从 config-schema 中标记为 `@deprecated`，非命令字段移至 `test:` 顶层，废弃检测正常工作
- [ ] 无 surface 配置的项目输出与当前一致
- [ ] 所有生成的配方通过 `--dry-run` 验证
- [ ] 每个 surface 规则文件记录了测试编排模式和 journey 过滤策略
- [ ] `run-tests` 检测到 `test.execution` 配置时输出废弃警告（不中断执行）
- [ ] 语言模板与 surface 规则的配方职责边界清晰（语言级 vs 编排级），无同名冲突
- [ ] 混合项目（web+api）端到端验证：`just dev` 无 scope 时按依赖顺序启动所有 scope 的 dev server，各 scope 的 probe 依次通过，测试执行完成后 teardown 逆序清理所有进程。PID 文件有效性校验覆盖混合项目的多 PID 场景

## 下一步

- 继续执行 `/quick-tasks` 生成实现任务
