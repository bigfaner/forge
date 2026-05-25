---
created: 2026-05-24
author: "faner"
status: Approved
---

# 提案：init-justfile Surface 感知 + 测试编排简化

## 术语定义

| 术语 | 定义 | 示例 |
|------|------|------|
| **surfaces map key** | `config.yaml` 中 `surfaces` 字段的 map key，原始值。仅作为配置源头的内部标识，用户不直接使用 | `admin-panel`、`payment-service`、`admin/service` |
| **surface-key** | surfaces map key 经 `/` → `-` 转换后的值。统一用于配方参数、文件路径、任务字段——是面向用户的唯一标识 | `admin-panel`、`payment-service`、`admin-service` |
| **surface type** | surface 的类型枚举值 | `web`、`api`、`cli`、`tui`、`mobile` |
| **编排模式** | surface 类型的抽象流程模板 | web = dev → probe → test → teardown |
| **执行策略** | 规则文件（`rules/surfaces/<type>.md`）中定义的步骤、退出码处理和调用契约。文中"执行策略规则文件"、"surface 规则文件"均指此文件 | `rules/surfaces/web.md` 的 `## 编排序列` 段 |

**关键映射**：`surface-key` 由 `surfaces map key` 经 `/` → `-` 转换而来。当 map key 不含 `/` 时两者相同（如 `admin-panel` → `admin-panel`）；当 map key 含 `/` 时（如 `admin/service`），surface-key 为 `admin-service`。配方体内部需要 surface-key → surfaces map key 的反向查找时，通过匹配转换规则实现（`admin-service` → 查 config 找到 key `admin/service`）。

## 问题

两个相互关联的问题（捆绑理由见下文论证）：

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

justfile 本身已经是抽象层，config.yaml 再包一层增加了复杂度但没有增加灵活性。对 config-schema.md 中记录的 8 个 `test.execution` 示例进行审计：6 个指向 just 命令（`just test {slug}`、`just test-setup`、`just probe` 等），2 个为非 just 路径（`go test -json -v ./...`、`npx vitest run --reporter=json`），即 **75% 的示例已通过 just 调用**，剩余 25% 可封装到 justfile 配方中保留。

### 证据

> **证据性质声明**：以下证据基于 Forge 仓库内的代码审计（config-schema.md 示例、SKILL.md 指令、Go 结构体定义）和逻辑推断，不包含外部用户反馈或实际项目部署数据。Forge v3.0.0 尚未发布，无存量用户和线上数据。量化数据（如"75% 的示例已通过 just 调用"）来源于 config-schema.md 中记录的 8 个示例，样本量有限，不应视为统计有效的结论。

**捆绑必要性论证**：为何不能在 v3.0.0 只做 surface 感知，委托层简化推迟到 v3.1.0？原因有三：(1) `run-tests` 的编排序列当前通过 `test.execution` 读取命令模板，若 v3.0.0 仅添加 surface 感知而保留 `test.execution`，则 surface 感知的编排参数（启动顺序、probe 目标、teardown 逆序）需要通过 `test.execution` 的子字段传递，导致 surface 编排与委托层深度耦合——v3.1.0 移除委托层时需要重新设计编排参数的传递路径，等于做两次设计。(2) 移除 `test.execution` 后 `run-tests` 的编排简化为直接调用 `just test`/`just dev`/`just probe`，surface 感知的编排差异通过 run-tests SKILL.md 中的 surface 执行策略规则文件传递——如果保留 `test.execution`，编排参数将分散在两个位置（config.yaml 命令模板 + SKILL 规则文件），增加维护负担。(3) v3.0.0 尚未发布，无存量用户，此时移除 `test.execution` 的迁移成本为零——推迟到 v3.1.0 意味着已有用户需要迁移。**零迁移成本窗口的边界条件**：此论证依赖"v3.0.0 未发布时无存量用户"的前提。如果 v3.0.0 延期发布但期间已有 beta/preview 用户（通过 GitHub pre-release 或早期测试计划），则这些用户的 `test.execution` 配置在升级后需要迁移。缓解措施：(a) 若 beta 用户数量 ≤ 10（Forge 团队可控范围），迁移通过人工沟通完成；(b) 若 beta 用户规模超出可控范围，在 v3.0.0 正式发布前提供一个迁移辅助脚本（读取旧 `test.execution` 配置并生成对应的 justfile 配方）。当前 beta 用户数量为零，此边界条件在可预见的未来不构成实际风险。

- Web UI 的 e2e 测试必须先启动应用，但当前配方没有 surface 特定的启动逻辑
- `run-tests` 的编排序列对 surface 不感知，只能依赖用户手动配置 `test.execution`
- `test.execution` 的多数示例指向 just 命令（`just test {slug}`、`just test-setup`、`just probe`），但 config-schema.md 也记录了非 just 示例（`go test -json -v ./...`、`npx vitest run --reporter=json`、`make test FEATURE={slug}`）。非 just 路径在简化方案中被牺牲——对于直接使用 `go test` 或 `npx vitest` 的项目，用户需要将这些命令封装到 justfile 的 `test` 配方中。这是可接受的 trade-off，因为 justfile 本身就是命令抽象层

### 紧迫性

随着 v3.0.0 test profile 的引入，surface 成为测试流程的核心维度。init-justfile 和 run-tests 需要协同工作 — init-justfile 生成正确的配方，run-tests 编排正确的执行序列。

## 提议方案

### 方案 A：Surface 感知 + 委托层简化

1. init-justfile 添加 surface 感知层，为不同 surface 生成差异化的配方
2. 移除 `test.execution` 委托层（不保留废弃检测），`run-tests` 直接调用 just 配方
3. `timeout`、`results-dir` 等非命令配置保留在 config.yaml 的简化字段中

### Surface 测试编排模式

| Surface | 测试编排序列 | 关键配方 |
|---------|------------|---------|
| **web** | `just dev`（后台）→ `just probe`（等待就绪）→ `just test` → teardown | `dev` 启动 dev server；`probe` 检查 HTTP 端点 |
| **api** | `just dev`（后台）→ `just probe`（等待就绪）→ `just test` → teardown | `dev` 启动 API 服务；`probe` 检查 /healthz |
| **cli** | `just build` → `just dev`（验证二进制可运行）→ `just test` | 无需 `run`/`probe`；`test` 通过子进程测试 |
| **tui** | `just build` → `just dev`（验证二进制可运行）→ `just test` | 无需 `run`/`probe`；`test` 通过 stdin 管道测试 |
| **mobile** | `just test-setup`（准备模拟器）→ `just dev`（启动应用）→ `just test` → teardown | `test-setup` 准备模拟器环境；`dev` 启动应用；`test` 用 maestro |

> **api/web 统一性说明**：api 和 web 的编排序列完全相同（dev → probe → test → teardown），唯一差异在于 `probe` 检查的目标端点（api 检查 `/healthz`，web 检查页面根路径）。当前选择保持两者为独立规则，原因是：(1) probe 端点差异可能在后续迭代中扩展为更大的行为差异（如 web 可能需要检查静态资源、api 可能需要检查数据库连接）；(2) 用户体验上，开发者通常以 surface 类型作为心智模型，合并为 `service` 会增加认知负担。若后续验证两者确实无实质性差异，可合并为 `service` 规则并共享编排模板。**合并收敛条件**：当以下所有条件同时满足时，触发合并评估：(1) 连续 3 个 Forge minor 版本中 web 和 api 的规则文件无实质性差异（仅 probe 端点路径不同，编排序列、退出码处理、teardown 逻辑完全相同）；(2) 无用户反馈要求区分两者（通过 GitHub issues、社区讨论等渠道确认）；(3) init-justfile 生成的 web/api 配方体在 5 个以上不同项目中验证无行为差异。合并时保留 probe 端点为参数化配置（`PROBE_TARGET` 环境变量），合并后的 `service` 规则通过 probe 端点差异覆盖 web/api 两种场景。

**dev/run 分工**：
- **web/api**：保留 `dev`（热重载 dev server）和 `run`（生产模式启动），测试编排使用 `dev`
- **cli/tui**：只生成 `dev`（编译+运行二进制），不生成 `run`
- **mobile**：生成 `dev`（启动应用）和 `test-setup`（准备模拟器环境），不生成 `run`
- **混合项目（web+api）**：`dev` 配方接受 surface-key 参数 — `just dev admin-panel`、`just dev payment-service`、`just dev`（无 surface-key 时按依赖顺序启动所有 surface）。**surface-key 值是 surfaces map key 经 `/` → `-` 转换后的结果**，而非 surface 类型枚举。例如，若 `surfaces: {admin-panel: web, payment-service: api}`，则 surface-key 值为 `admin-panel` / `payment-service`

**混合项目多服务启动管理**：

当 `just dev`（无 surface-key）需要启动多个服务时：


- **端口冲突预防（best-effort）**：`just dev` 配方体在启动前尝试检查端口是否已被占用。此检查为 **best-effort**——不保证跨平台一致性，且存在 TOCTOU 竞态（检查和启动之间端口可能被占用）。具体策略：Linux/macOS 使用 `lsof -i :$PORT`（注意：Linux 上可能需要 root 权限，失败时静默跳过）；Windows 使用 `netstat -ano | findstr :$PORT`。**端口号来源**：`$PORT` 的值优先从配方体的环境变量解析结果获取（即 `PORT := env("PORT", "3000")`，先检查环境变量覆盖值，无覆盖时使用默认值），确保端口预检与 dev server 实际监听的端口一致。若环境变量解析失败或端口值无法确定，使用默认端口执行检查并附带说明："端口预检使用默认端口 N，若已通过环境变量覆盖，此结果可能不准确"。若检查失败（权限不足或命令不可用），跳过检查直接启动，将端口冲突检测交给 dev server 本身的启动错误处理。**端口冲突不作为编排流程的硬性门控**——它是用户体验优化（提前报错而非等待 probe 超时），probe 超时是端口冲突的最终兜底检测机制。端口检查报告的错误信息包含限定词："端口预检（best-effort，结果仅供参考）：端口 N 可能已被占用，实际状态以 dev server 启动结果为准"
- **顺序启动策略**：`just dev` 配方体按顺序启动各 surface-key 的 dev server（先启动后端，再启动前端），每个启动后将 PID 写入 `.forge/dev-server.<surface-key>.pid`。顺序启动（而非并行）是为了避免多进程同时争抢系统资源导致启动失败
- **probe 重试差异化**：差异化重试在 `just probe` 配方体内部实现（而非 run-tests 层面），通过特定退出码向 run-tests 传递失败类型。配方体内区分三种失败模式：(1) **连接被拒绝**（ECONNREFUSED）— 端口无进程监听，继续重试（默认行为）；(2) **连接超时** — 端口有进程但不响应 HTTP，配方体内部增加计数权重（每次超时计为 2 次重试）加速退出，最终以 exit 3 退出（run-tests 识别后跳过后续 probe 直接 teardown）；(3) **端口冲突**（日志 EADDRINUSE）— 立即以 exit 2 退出（run-tests 识别后跳过所有重试直接 teardown）。退出码约定：exit 0 = 健康，exit 1 = 通用失败（默认重试行为），exit 2 = 端口冲突（立即中止），exit 3 = 连接超时（加速退出）。run-tests 根据退出码执行对应的后续动作，无需解析 probe 的 stderr 输出
- **probe 顺序**：run-tests 先 probe 后端（api），再 probe 前端（web）。后端 probe 完成后输出 `[probe] 后端就绪，开始探测前端...`，提供进度反馈
- **teardown 逆序清理**：先 teardown 前端，再 teardown 后端。逆序清理模拟生产环境的依赖关系

### 委托层简化

移除 `test.execution` 委托层（不再读取、不保留检测），`run-tests` 的编排变为：

```
当前（4 层委托）：
run-tests → config test.execution.run → "just test {slug}" → just test → 实际运行器

简化后（2 层）：
run-tests → just test [surface-key] [journey] → 实际运行器
```

`run-tests` 根据 surface 编排模式直接调用 just 配方，不再读 config 命令模板：

| 编排步骤 | 当前 | 简化后 |
|---------|------|--------|
| setup | `test.execution.setup` 配置 | `just test-setup`（仅 mobile；其他 surface 无此步骤） |
| pre-check | `test.execution.pre-check` 配置 | `just probe`（web/api）或跳过（cli/tui/mobile） |
| run | `test.execution.run` 配置 | `just test [surface-key] [journey]` |
| teardown | `test.execution.teardown` 配置 | `just test-teardown`（可选配方，web/api/mobile） |

`test.*` 配置项从 `.forge/config.yaml` 中移除。测试相关配置下沉到各 surface 的测试目录中：

- **results-dir**：由测试目录结构约定隐式决定（见下一节），无需配置
- **timeout / 运行时参数**：由 `tests/{surface-key}/config.yaml`（单 surface 项目为 `tests/config.yaml`）提供，testkit 解析后供测试脚本使用

### 测试目录结构约定

测试用例按 surface-key 组织目录：

- **多 surface 项目**：`tests/{surface-key}/`
- **单 surface 项目**：`tests/`（无子目录）

其中 `surface-key` 由 `config.yaml` 中 `surfaces` map 的 key 转换而来——将其中的 `/` 替换为 `-`。例如：

| surfaces map key | surface-key | 测试目录 |
|-----------------|-------------|---------|
| `admin-panel` | `admin-panel` | `tests/admin-panel/` |
| `admin/service` | `admin-service` | `tests/admin-service/` |
| `payment-service` | `payment-service` | `tests/payment-service/` |

**单 surface 判定规则**：`surfaces` map 仅有一个条目时，测试用例直接放在 `tests/` 下，不创建 surface-key 子目录。例如 `surfaces: {web-app: web}` → 测试目录为 `tests/`，而非 `tests/web-app/`。

**对下游的影响**：

| 组件 | 变更 |
|------|------|
| `results-dir` | 由目录结构约定隐式决定：多 surface `tests/{surface-key}/{journey}/results`，单 surface `tests/{journey}/results` |
| `just test` 配方签名 | 单 surface: `just test [journey]`，直接运行 `tests/` 下的测试；多 surface: `just test [surface-key] [journey]`，按 surface-key 枚举测试目录，配方体内含参数反转逻辑 |
| `run-tests` 编排 | 探测 surface 数量，决定测试目录路径；不配置任何 timeout 参数，编排流程由各 just 配方体自身的退出码和重试逻辑控制 |
| `gen-test-cases` | 生成测试用例时写入 `tests/{surface-key}/`（多 surface）或 `tests/`（单 surface）；生成 `config.yaml` 模板 |
| testkit / 辅助工具 | 解析 `tests/{surface-key}/config.yaml`，向测试用例暴露 surface 特定的测试参数 |

**per-surface 测试配置**：每个 surface 的测试目录下放置 `config.yaml`，由 testkit 解析，向测试脚本暴露配置项。**与 run-tests 编排层完全独立**——run-tests 不读取此文件。单 surface 项目配置文件为 `tests/config.yaml`，多 surface 项目为 `tests/{surface-key}/config.yaml`。不同 surface 类型的配置参数不同，由各自的 surface 规则文件定义。示例：

```yaml
# tests/admin-service/config.yaml — web surface E2E 测试配置
# E2E Test Configuration
baseUrl: http://localhost:5173
apiBaseUrl: http://localhost:8080

# Test account credentials (leave empty if no auth required)
username: 'admin'
password: 'admin123'

# Login page locators (matched to sitemap E-002, E-003, E-004)
loginLocators:
  usernameField: "账号"
  passwordField: "密码"
  submitButton: "登录"
```

**各 surface 类型的典型配置项**（具体字段由 surface 规则文件定义，以下为示例）：

| Surface | 典型配置项 |
|---------|-----------|
| **web** | `baseUrl`、`apiBaseUrl`、`username`/`password`、页面 locators |
| **api** | `baseUrl`、`authToken`、`defaultHeaders` |
| **cli** | `binaryPath`、`envVars` |
| **tui** | `binaryPath`、`terminalSize` |
| **mobile** | `appPath`、`deviceName`、`platformVersion` |

**配置生成**：`gen-test-cases` 在生成测试用例时，同时生成 `config.yaml` 模板（含默认值和注释说明）。用户可按需修改。

**超时与重试控制**：SKILL.md 不配置任何 timeout 参数——LLM/agent 无法精确控制工具的执行超时。超时和重试由各 just 配方体自身的代码逻辑控制：

- `just probe`：配方体内置重试循环（默认 30 次 / 2 秒间隔），重试耗尽后根据失败模式以对应退出码退出（exit 1 = 通用超时，exit 2 = 端口冲突，exit 3 = 连接超时）
- `just test-teardown`：配方体内置 SIGTERM → 等待 5 秒 → SIGKILL 的超时清理
- `just dev`：后台启动，失败由 probe 检测（PID 存活检查 + 端口监听检测）


### Skill 流程

init-justfile 的生成流程：
```
检测语言 → 加载语言模板 → 生成 compile/build/lint/fmt
检测 surface → 加载 surface 规则 → 生成 test/dev/run/probe/test-setup
组装 → 验证
```

**跨平台变体验证检查点**：init-justfile 在生成跨平台配方（`[linux]`/`[windows]` 双变体）后，在"验证"步骤中增加变体验证检查——对每个双变体配方执行以下检查：(1) 两个变体的配方签名完全一致（参数名、参数默认值相同）；(2) 两个变体的环境变量引用一致（引用的变量名集合相同）；(3) 两个变体的退出码语义一致（成功/失败的退出码约定相同）。此检查确保 LLM 生成两个语法完全不同的变体时不会引入语义偏差。

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
- **编排序列中的引用**：surface 规则的编排序列可能引用语言模板负责的配方（如 cli/tui 编排中的 `just build`），此时 surface 规则仅**引用**该配方作为编排步骤，不负责生成其配方体——配方体仍由语言模板生成
- **冲突处理**：当两者试图生成同名配方时，Surface 规则覆盖语言模板。例如，语言模板生成的 `test` 配方会被 surface 规则定义的编排模式替换。**用户编辑保护**：init-justfile 在覆盖编排级配方前，检查 justfile 中是否包含用户手动编辑标记（`# user-customized` 注释行）。若目标配方已有此标记，init-justfile 执行以下保护行为：(1) 跳过覆盖该配方；(2) 输出差异摘要——将当前配方体与新生成版本逐行对比，列出变更点（如"probe 默认重试次数从 20 变为 30"、"新增 PROBE_INTERVAL 环境变量支持"），用户据此判断是否需要手动同步；(3) 提示用户两种更新途径：移除 `# user-customized` 标记后重新运行 init-justfile（自动覆盖），或运行 `init-justfile --force-regenerate`（强制覆盖所有 `# user-customized` 配方，适用于确认所有自定义可丢弃的场景）
- **组装策略**：init-justfile 先加载语言模板生成基础配方，再加载 surface 规则覆盖编排配方，最终组装为完整的 justfile

run-tests 的执行流程（调度器模式，与 gen-test-cases 在流程骨架上同构）：
```
检测 surface type（forge surfaces / config.yaml）
→ 加载对应执行策略规则（rules/surfaces/<type>.md）
→ 按策略编排执行 just 配方序列

示例（web surface）：
just dev → just probe → just test [journey] → just test-teardown

示例（cli surface）：
just build → just dev → just test [journey]
```

**与 gen-test-cases 的同构范围与差异点**：

两者共享"检测 → 加载规则 → 按规则执行"三步流程骨架，但执行层有本质差异：
- **gen-test-cases**：按规则生成**文本产物**（测试用例代码），无进程管理、无退出码门控
- **run-tests**：按规则编排**进程生命周期**（启动/等待/测试/清理），需处理退出码门控、PID 追踪、状态机转移

同构仅覆盖流程骨架，不意味着可共享抽象基类。

**与现有 surface 感知环境检查的关系**：

当前 `run-tests` SKILL.md 第 4 步的"Environment Readiness Check"已通过 `rules/env-check.md` 和 `surface-<type>.md` 做 surface 感知的环境检查（如检查 playwright 是否安装、端口是否可用）。提案与现有机制**互补而非重叠**：

- **现有机制（Environment Readiness Check）**：关注**前置条件检查** — 测试工具是否安装、依赖是否就绪、端口是否被占用。在编排执行之前运行
- **提案新增（Surface 感知编排）**：关注**执行序列编排** — 启动/停止服务的顺序、何时 probe、何时 teardown。在环境检查通过后运行

两者协作关系：`env-check` 确认环境可用 → 提案的编排模式按 surface 类型执行正确的序列。现有 `surface-<type>.md` 规则文件可以继续用于环境检查，无需修改。

### 创新亮点

将 init-justfile（配方生产者）和 run-tests（调度器）的 surface 感知通过规则文件统一设计。**规则文件物理独立但逻辑同构**——每个 skill 各自持有 `rules/surfaces/<type>.md` 的独立副本，但同一 surface 类型的两个文件共享相同的编排序列和配方调用契约，确保生成与执行的一致性。物理独立是因为：(1) 遵循 forge-distribution.md 路径约定，各 skill 规则文件位于各自目录下；(2) 两个 skill 消费方式不同（init-justfile 读取配方签名生成配方体，run-tests 读取编排序列调度执行），各自可独立演进。逻辑同构通过新增 surface 类型时同步创建两个 skill 的规则文件保证——若仅更新其中一个，另一个运行时回退到默认行为。

**调度器模式**（流程骨架与 gen-test-cases 同构，见上文差异点说明）：

run-tests 不自行推断编排模式，也不依赖跨 skill 文件传递编排参数。它遵循与 gen-test-cases 相同的"检测 → 加载 → 执行"流程骨架：

1. **检测 surface type**：通过 `forge surfaces` CLI 或 `config.yaml` 的 `surfaces` 字段获取 surface 类型
2. **加载执行策略**：读取 `rules/surfaces/<type>.md` 规则文件，获取该 surface 类型的编排序列定义
3. **按策略执行**：策略规则中定义了编排序列（哪些 just 配方需要调用、调用顺序、退出码处理逻辑），run-tests 按策略执行

**执行策略规则文件**（run-tests `rules/surfaces/<type>.md`）的职责：

每个 surface 类型的规则文件物理上独立存在于各 skill 目录下（`init-justfile/rules/surfaces/<type>.md` 和 `run-tests/rules/surfaces/<type>.md`），通过 Markdown 标题分段承载两个职责：

- **编排序列**（`## 编排序列` 段，供 run-tests 消费）：该 surface 类型的测试执行流程（如 web = dev → probe → test → teardown），定义步骤顺序、退出码门控逻辑和**退出码处理表**。退出码处理表是规则文件的强制性结构元素，定义了每个编排步骤可能返回的退出码及其对应的后续动作。新增退出码时必须同步更新此表，否则 run-tests 会按默认行为（通用失败 → teardown）处理未定义的退出码。**注意**：probe 的退出码是配方体重试循环**耗尽后**的最终退出码（run-tests 不感知单次重试的结果，只感知 probe 命令的最终退出码）。示例如下：

  | 步骤 | 退出码 | 语义 | run-tests 后续动作 |
  |------|--------|------|-------------------|
  | probe | 0 | 服务健康 | 继续下一步（test） |
  | probe | 1 | 重试耗尽（服务未就绪） | 执行 teardown 后中止，报告"服务启动超时" |
  | probe | 2 | 端口冲突（EADDRINUSE，配方体检测到后立即退出） | 执行 teardown 后中止，报告"端口冲突" |
  | probe | 3 | 连接超时（端口有进程但不响应，加速退出） | 执行 teardown 后中止，报告"连接超时" |
  | dev | 0 | 启动成功 | 继续下一步（probe） |
  | dev | 1 | 启动失败 | 执行 teardown 后中止 |
  | test | 0 | 测试通过 | 继续下一步（teardown） |
  | test | 1 | 测试失败 | 执行 teardown 后中止 |

- **just 配方调用契约**（`## 配方调用契约` 段，供 run-tests + init-justfile 共同消费）：序列中每个步骤的 just 配方名、参数、退出码语义。init-justfile 读取此段确定需要生成哪些配方及签名；run-tests 读取此段确定调用参数和退出码处理

两个分段独立、无交叉引用，消费方按需读取对应分段。

`probe_target` 等项目特定参数内嵌在 just 配方体中（由 init-justfile 生成时写入），不通过中间文件传递。run-tests 只关心"调用 `just probe` 并检查退出码"，不关心 probe 的具体目标 URL。

**新增 surface 类型的扩展方式**：

新增 surface 类型只需两步：
1. 在 init-justfile 的 `rules/surfaces/` 下新增 `<type>.md`，定义配方生成指导
2. 在 run-tests 的 `rules/surfaces/` 下新增 `<type>.md`，定义编排序列和调用契约

**Go 代码校验验证**：经代码审计确认，Go 的 `surfaces` 字段类型为 `map[string]string`，不对 surface 类型值域做枚举校验——新增类型无需更新 Go 代码。但需注意：若 `forge surfaces` CLI 内部对返回值做了类型过滤（如只返回已知类型），需在实现时验证，必要时同步更新 CLI 类型白名单。

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
- **混合项目（web+api）**：`just dev` 无 surface-key 时按依赖顺序启动前端和后端（例如 `surfaces: {admin-panel: web, payment-service: api}` 时，先 `just dev payment-service` 再 `just dev admin-panel`），probe 分别检查两个服务

### 下游集成契约

配方签名不可变。`run-tests` 的调用方式简化但语义不变：

| 配方 | 签名（不可变） | 消费者 | 期望语义 |
|------|--------------|--------|---------|
| `unit-test` | `just unit-test` | forge task submit、clean-code、fix-bug、testrunner | 语言级单元测试；exit 0 = 通过 |
| `test`（单 surface） | `just test [journey]` | run-tests、forge quality-gate、fix-bug | Surface 级测试；`journey` 为可选参数，为空时运行全部。单 surface 项目不生成 surface-key 参数 |
| `test`（多 surface） | `just test [surface-key] [journey]` | run-tests、forge quality-gate、fix-bug | Surface 级测试；`surface-key` 为可选参数；`journey` 为可选参数，为空时运行全部。单参数时按优先级反转（见关键约束） |
| `probe` | `just probe` | run-tests（web/api 前置检查） | 服务健康检查；exit 0 = 健康；混合项目检查所有服务 |
| `test-setup` | `just test-setup` | run-tests | 安装测试依赖；幂等 |
| `test-teardown` | `just test-teardown` | run-tests（可选） | 测试后清理 |
| `dev` | `just dev [surface-key]` | run-tests（web/api 启动服务） | web/api: 后台启动并监听端口；cli/tui: 编译运行；`surface-key` 为可选参数，仅混合项目需要 |
| `run` | `just run [surface-key]` | 仅 web/api 生成 | 生产模式启动 |

**关键约束**：
- **签名按项目类型分化**：单 surface 项目生成 `test journey=''`（无 surface-key 参数），多 surface 项目生成 `test surface-key='' journey=''`。init-justfile 根据项目 surface 数量决定生成哪个签名
- **多 surface 项目的参数解析优先级**：当只传入一个位置参数时（如 `just test e2e`），just 将其绑定为 surface-key——配方体需做反转判断：先与所有 surface-key 集合匹配，匹配成功视为 surface-key，否则重新解释为 journey。此反转逻辑仅在多 surface 项目的配方体中存在
- **单 surface 项目无歧义**：`just test e2e` 直接绑定为 journey，无需反转。`just test` 运行全部测试
- **journey 与 surface-key 职责分离**（仅多 surface 项目）：journey 限定测试类型（如 `e2e`、`smoke`），surface-key 限定测试的 surface 范围。`just test admin-panel e2e` 仅运行 admin-panel 的 e2e 测试；`just test e2e` 运行所有 surface 的 e2e 测试；`just test admin-panel` 运行 admin-panel 的全部测试
- web/api 的 `dev` 须能在后台运行
- cli/tui 不生成 `run` 配方
- `test-teardown` 为可选配方，不存在时 run-tests 跳过

### Surface-key 解析变更

提案涉及两个独立的参数维度，当前与提案的差异如下：

#### 维度 1：test 命令参数——从 feature slug 到 journey

| 对比项 | 当前 | 提案 |
|--------|------|------|
| 命令模板 | `just test {slug}` | `just test [journey]` |
| 参数语义 | feature slug（测试哪个功能） | journey 分类（测试哪类场景） |
| 参数来源 | `forge feature` → state.json / git 分支 / 目录扫描 | run-tests 传入；为空时运行全部测试 |
| 空值行为 | 无 slug → abort | journey 为空 → 运行全部 |

`{slug}` 参数从 test 命令中完全移除。测试范围由 journey 分类（如 `e2e`、`contract`、`smoke`）限定，不再按 feature 过滤。

#### 维度 2：统一参数值域——从固定枚举到 surface-key

当前混合项目的所有配方共享 `frontend`/`backend`/`all` 枚举。提案统一迁移为 surface-key（surfaces map key 经 `/` → `-` 转换后的值）。

**迁移前后对比**：

```just
# 当前（frontend/backend 枚举）
compile surface-key="":
    case "{{ surface-key }}" in
        backend)  go vet ./... ;;
        frontend) npx tsc --noEmit ;;
        *)        go vet ./... && npx tsc --noEmit ;;
    esac

# 迁移后（surface-key）
compile surface-key="":
    case "{{ surface-key }}" in
        payment-service)  go vet ./... ;;
        admin-panel)      npx tsc --noEmit ;;
        *)                go vet ./... && npx tsc --noEmit ;;
    esac
```

| 对比项 | 当前 | 提案 |
|--------|------|------|
| 参数值域 | 固定枚举：`frontend` / `backend` / `all` | 用户自定义：surface-key（如 `admin-panel` / `payment-service`）；跨 surface 仍用 `all` |
| 参数来源 | `project-type` 字段 + 文件路径模式匹配 | surface-key 由 surfaces map key 转换而来 + 路径归属映射 |
| 适用范围 | 混合项目所有配方 | 有 surfaces 配置的混合项目所有配方 |

**迁移影响面**：

| 组件 | 当前行为 | 迁移变更 |
|------|---------|---------|
| `breakdown-tasks` `rules/surface-key-assignment.md` | 文件路径模式：`ui/`/`components/` → frontend，`cmd/`/`internal/` → backend | 按 surfaces 路径归属：文件属于哪个 surface 的目录 → 该 surface-key |
| `quick-tasks` SKILL.md | 内联推断：UI → frontend，API → backend | 按 surface 路径归属推断 surface-key |
| `breakdown-tasks` `rules/db-schema.md` | 硬编码 `surface-key: 'backend'` | 找到 type=api 的 surface-key，若无则回退 `all` |
| `init-justfile` SKILL.md | 混合项目配方 `case surface-key in frontend/backend/all` | `case surface-key in <surface-keys>/all` |
| `prompt.go` `resolveSurfaceKey()` | 按 `project-type` 枚举校验（backend 项目清 frontend） | 按 surface-key 集合校验：值不在集合中且非 `all`/空则清除 |
| `just.go` `ResolveSurfaceKey()` | 干运行探测 justfile 是否接受 surface-key 参数 | **无需变更**（值无关的探测逻辑） |
| 16 个 prompt 模板 | `just compile {{SURFACE_KEY}}`，值为 frontend/backend | 同一模板语法，值为 surface-key |
| Task 数据模型 | `surface-key` 字段值为 frontend/backend/all，无 surface-type | `surface-key` 字段值为 surface-key 或 `all`；新增 `surface-type` 字段 |
| 任务模板（breakdown-tasks / quick-tasks） | frontmatter 仅含 `scope` | `scope` 重命名为 `surface-key`，新增 `surface-type` 字段 |

**Task 数据模型扩展**：

Go struct 新增字段（`forge-cli/pkg/task/types.go`）：

```go
type Task struct {
    // ... 现有字段 ...
    SurfaceKey   string `json:"surface-key,omitempty"`   // 替代原 Scope 字段，值为 surface-key 或 "all"
    SurfaceType  string `json:"surface-type,omitempty"`  // surface 类型枚举：web/api/cli/tui/mobile
    // Scope 字段保留一个版本用于向后兼容（v3.1.0 移除），与 SurfaceKey 同步写入
}
```

任务模板 frontmatter 变更（`plugins/forge/skills/*/templates/task.md`）：

```markdown
---
id: "001"
title: "..."
surface-key: "admin-panel"    # 替代原 scope
surface-type: "web"           # 新增
---
```

**字段填充规则**：

| 字段 | 填充时机 | 填充方式 |
|------|---------|---------|
| `surface-key` | `breakdown-tasks` / `quick-tasks` 生成任务时 | 按 surface-key-assignment 规则推断（同现有 scope 推断逻辑，值域变更） |
| `surface-type` | `breakdown-tasks` / `quick-tasks` 生成任务时 | 根据 surface-key 查 config.yaml 的 `surfaces` map 获取对应类型 |
| `surface-key`（兼容层） | `forge task claim` 读取旧任务时 | 若存在 `scope` 字段但无 `surface-key`，按过渡期兼容层映射 |

**全入口填充策略**：

| 创建入口 | surface-key | surface-type | 填充方式 |
|---------|-------------|-------------|---------|
| `breakdown-tasks` / `quick-tasks` 模板 | surface-key 或 `all` | web/api/cli/tui/mobile | 按 surface-key-assignment 规则推断；surface-type 从 surfaces map 查类型 |
| `forge task add` CLI | 继承源任务值 | 继承源任务值 | 从 `SourceTaskID` 指向的源任务复制 `surface-key` 和 `surface-type`。无源任务时：单 surface 项目 surface-key 为空、surface-type 为唯一 surface 的类型；多 surface 项目均默认空 |
| quality-gate fix-task | surface-key 或 `all` | web/api/cli/tui/mobile | 复用 `resolveSurfaceKey()` 从失败文件路径推断；surface-type 从 surfaces map 查类型。无法推断时默认空 |
| auto-gen test tasks | 硬编码 `"all"` | 从任务类型提取 | 现有 `scope: "all"` 改为 `surface-key: "all"`；`surface-type` 从 per-type 任务定义提取（如 `T-test-gen-journeys-api` → `api`） |
| stage gates | 不填充 | 不填充 | 项目级门控，非 surface 特定 |
| `forge task reopen` | 保留原值 | 保留原值 | 不涉及 |

**消费方**：

| 消费方 | 如何使用 `surface-type` |
|--------|----------------------|
| `execute-task` / `task-executor` | 加载对应的 surface 执行策略规则文件（`rules/surfaces/<type>.md`） |
| `init-justfile` | 知道任务涉及哪个 surface 类型，生成对应配方 |
| `run-tests` | 根据任务的 surface-type 选择编排模式 |
| `fix-bug` | 验证修复后选择正确的测试命令 |

**surface-key-assignment 迁移细则**：

当前规则按目录名模式推断（`ui/` → frontend，`cmd/` → backend）。迁移后按 surfaces 路径归属推断：

1. 读取 `config.yaml` 的 `surfaces` 字段，获取每个 surface 的目录路径
2. 对任务的每个受影响文件，检查其路径属于哪个 surface 的目录
3. 所有文件属于同一 surface → surface-key 为该 surface 的转换值
4. 文件跨越多个 surface → surface-key 为 `all`
5. 文件不属于任何 surface 目录 → surface-key 为 `all`

`all` 语义不变：跨 surface 或无法归类。配方体中 `all` 等同于空值（编译所有 surface-key）。


**迁移顺序约束与原子性保证**：

surface-key 值域迁移涉及 7 个以上组件，变更必须在特定顺序下执行以保证系统一致性。迁移顺序约束如下：

1. **阶段 1——数据模型（最先变更）**：`prompt.go resolveSurfaceKey()` 校验逻辑从固定枚举切换为 surface-key 集合查询。此变更不影响无 surfaces 配置的项目（枚举值在集合中无匹配时回退 `all`），但为后续组件提供正确的 surface-key 值域基础
2. **阶段 2——规则引擎**：`breakdown-tasks rules/surface-key-assignment.md` 和 `rules/db-schema.md`、`quick-tasks SKILL.md` 同步更新，surface-key 推断逻辑切换到 surfaces 路径归属
3. **阶段 3——模板层**：`init-justfile SKILL.md` 混合项目配方生成更新（case 分支从 `frontend/backend` 改为 surface-key）
4. **阶段 4——16 个 prompt 模板**：`SURFACE_KEY` 变量值域同步更新。模板语法不变，仅运行时传入值变化

**原子性保证**：阶段 1-4 必须在**同一 PR** 中完成，允许 PR 内按阶段拆分为多个逻辑提交（便于 code review 分阶段审查）。PR 合并前必须通过集成测试验证所有阶段联合行为一致。PR 不允许在中间状态合入，因为部分合入会导致 surface-key-assignment 输出 `admin-panel` 但 init-justfile 仍期望 `frontend` 的不一致。

**过渡期兼容层**：阶段 1 的 `resolveSurfaceKey()` 保留一个版本的向后兼容逻辑——若值为旧枚举（`frontend`/`backend`），转换为对应的 surface-key。具体映射：`frontend` → 查找 type=web 的 surface-key，`backend` → 查找 type=api 的 surface-key。若 surfaces 中不存在对应类型，输出警告并回退为 `all`，确保任何配置下都有确定性行为。

**多 surface 同类型冲突解决规则**：当多个 surface 同类型时（如两个 type=web 的 surface `admin-panel` + `marketing-site`），兼容层按以下规则消歧：(1) 若只有一个匹配类型的 surface，直接映射；(2) 若有多个匹配，**不依赖声明顺序**——Go 的 `map[string]string` 迭代顺序不确定。消歧策略改为：按 surface-key 的**字典序**选择第一个匹配的 surface-key（确定性且不依赖运行时行为）；(3) 同时输出警告："'frontend' 映射到 'X'，但存在多个 type=web surface（X, Y），建议显式使用 surface-key"。此兼容层在 **v3.1.0** 中移除（即保留一个 minor version：v3.0.x 全系列包含兼容层，v3.1.0 起移除）。这确保了如果某个组件在迁移窗口中遗漏，旧值仍能被正确映射而非静默丢弃。移除时需同步清理 `resolveSurfaceKey()` 中的旧枚举映射代码，并在 CHANGELOG 中标注 breaking change。

**单 surface 项目和无 surfaces 配置的项目**：

- 单 surface 项目（surfaces 只有一个条目）：surface-key 始终为空（无需区分），行为与当前单语言项目一致
- 无 surfaces 配置：surface-key 逻辑不生效，`compile`/`lint`/`fmt` 不生成 surface-key 参数，回退到当前行为

### 约束与依赖

- Surface 信息来自 `.forge/config.yaml` 或 `forge surfaces` CLI（**待验证**：需确认 `forge surfaces` CLI 是否已实现及返回格式）
- **Surface 信息源优先级规则**（init-justfile 和 run-tests 统一遵循）：
  1. **`config.yaml` 的 `surfaces` 字段优先**：若 `surfaces` 字段存在且非空，以其中定义的 surface 类型和 surface-key 映射为准
  2. **`forge surfaces` CLI 回退**：若 `surfaces` 字段缺失或为空，通过 `forge surfaces [path]` CLI 基于文件信号检测获取 surface 类型（**待验证**：需在实现时确认 CLI 返回格式）
  3. **冲突处理**：当两者不一致时（config 说 web，文件信号检测出 api），以 `config.yaml` 为准（用户显式配置优先于自动检测）
- Surface 规则保持语言无关
- Standard Target Contract 配方名称和签名不变
- 遵循 forge-distribution.md 路径约定

### 非功能需求

| NFR | 要求 | 验证方式 |
|-----|------|---------|
| **跨平台兼容** | Windows/macOS/Linux 三平台均可运行编排序列。`just dev`/`probe`/`test-teardown` 配方体需包含平台分支 | 各平台手动验证；CI 矩阵（如果接入） |
| **向后兼容** | 无 surface 配置的项目行为与当前完全一致 | diff 输出对比 |
| **可观测性** | 编排过程中每个步骤（dev 启动、probe 轮询、test 执行、teardown 清理）按 SKILL.md 定义的固定格式输出步骤状态（`[步骤名] [状态] [摘要]`，如 `[probe] [retry 3/30] http://localhost:3000 — 连接被拒绝`），用户可追踪执行进度 | SKILL.md 输出格式与实际输出对比 |
| **性能** | surface 规则文件加载不增加 init-justfile 超过 1 秒的额外耗时 | 计时基准测试 |
| **可靠性** | dev server 崩溃时 probe 超时后执行 teardown，不遗留孤儿进程；会话中断后可通过 test-state.json 恢复清理 | 故障注入测试 |
| **just 版本** | just >= 1.4.0（`[linux]`/`[windows]` recipe attribute 在 just 1.4.0 引入，1.0-1.3.x 不支持此功能，会报 "unknown attribute" 错误）；无需后台运行相关的特殊版本。**版本检查触发时机**：在 init-justfile 和 run-tests 的 SKILL.md 中，作为首个执行步骤（在检测 surface 类型之前）运行 `just --version` 并解析版本号。若版本 < 1.4.0，输出明确的错误信息（"just 版本不满足要求：当前 X.Y.Z，需要 >= 1.4.0。升级方式：<平台特定安装指令>"）并中止执行。run-tests 同样在编排序列执行前检查，防止在不兼容的 just 版本下启动 dev server 后因配方语法错误导致进程无法清理 | SKILL.md 步骤定义；手动安装 just 1.3.x 后运行 init-justfile 验证错误提示 |

## 替代方案

| 方案 | 优势 | 劣势 | 结论 |
|------|------|------|------|
| 不做 | 零成本 | 下游无法 surface 感知编排；test.execution 委托层继续累积复杂度 | 拒绝 |
| 仅 surface 感知，保留 test.execution | 改动小 | 委托层冗余未解决；两处配置 surface 行为 | 拒绝：治标不治本 |
| **surface 感知 + 移除 test.execution** | 统一抽象层；消除冗余委托 | 改动范围扩大到 run-tests；非 just 路径需封装到 justfile | **选定** |
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
| **Playwright** `webServer` config | 声明式（配置文件） | 自动轮询 URL 直到 HTTP 响应，`timeout` 可配 | 框架托管 dev server 进程，测试结束自动清理 | 前端 E2E 测试 |
| **Vitest** `pool` + `setupFiles` | 配置驱动 | 无 HTTP 探针（单元测试无需服务就绪） | 进程池管理测试隔离 | 单元/集成测试 |
| **Testcontainers** | 声明式（代码 API） | 容器内 HTTP/TCP 探针，支持 `waitingFor` 策略 | Docker 容器生命周期（Ryuk sidecar 自动清理） | 集成测试中需要真实服务依赖（数据库、消息队列） |

### Forge 方案的定位差异

以上方案的共同特点：**编排逻辑由确定性代码执行**（Docker daemon、kubelet、Node.js 进程）。Forge 的关键差异在于：

- **run-tests 是 LLM agent 执行的 SKILL**，编排序列由 LLM 按步骤执行，而非确定性代码
- **justfile 是文本协议**，不是运行时——just 配方启动的进程在 just 退出后可能成为孤儿进程

这意味着 Forge 不能直接复用上述方案的进程管理能力，需要在 SKILL 指导层面建立可靠性保证（见"LLM agent 执行确定性"和"后台进程管理"两节）。

**为何不复用测试框架内建的编排能力**：Playwright 的 `webServer` 配置和 Cypress 的 `start-server-and-test` 在各自框架内提供了完善的"启动 → 等待 → 测试 → 清理"编排。但 Forge 的 run-tests 需要编排**任意语言和框架**的测试流程，不能绑定特定测试框架。此外，Forge 的 surface 感知覆盖了 CLI/TUI/Mobile 等无 HTTP 服务的场景，这些场景不在 Playwright/Vitest 的编排范围内。因此 Forge 必须在 justfile 层面建立框架无关的编排协议，测试框架的编排能力作为其各自 just 配方体内的实现细节存在。

**为何不采用 Testcontainers 模式**：Testcontainers 是"启动服务依赖 → 等待就绪 → 测试 → 清理"模式在 Java/Go/Node 生态中成熟的实现，通过 Docker 容器管理服务生命周期（启动、探针、清理由 Ryuk sidecar 保证）。Forge 不采用此模式的原因：(1) Forge 编排的被测服务是用户项目本身的 dev server（如 `npm run dev`、`go run`），不是外部依赖服务（如 PostgreSQL、Redis）——dev server 不适合容器化，因为它需要访问宿主机的文件系统和热重载能力；(2) Testcontainers 依赖 Docker daemon 运行，对 Forge 用户的开发环境引入了额外的基础设施要求，违背 Forge"零外部依赖"的设计原则；(3) Forge 的 CLI/TUI/Mobile surface 类型不涉及服务启动，Testcontainers 的容器化编排对这些场景无意义。如果未来 Forge 需要编排外部服务依赖（如测试数据库），可以考虑在 justfile 配方体中使用 Testcontainers 作为实现细节，但这不在 v3.0.0 的范围内。

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
    mkdir -p .forge && nohup npm run dev > .forge/dev-server.log 2>&1 & echo $! > .forge/dev-server.pid

[windows]
dev:
    powershell -Command "if (!(Test-Path '.forge')) { New-Item -ItemType Directory -Path '.forge' | Out-Null }; $p = Start-Process -FilePath 'npm' -ArgumentList 'run','dev' -RedirectStandardOutput '.forge\\dev-server.log' -RedirectStandardError '.forge\\dev-server-err.log' -PassThru -WindowStyle Hidden; $p.Id | Out-File -FilePath '.forge\\dev-server.pid' -Encoding utf8"
```

> **日志输出保留**：dev server 的 stdout/stderr 重定向到 `.forge/dev-server.log`（非 /dev/null），启动失败时（编译错误、依赖缺失）用户和 run-tests 可查看日志定位原因。probe 超时后 run-tests 在错误消息中附带日志最后 10 行内容，避免用户在"服务启动超时"后还需手动查找失败原因。


> **注意**：Windows CMD 的 `start /B` 不暴露子进程 PID——`%PID%` 变量不会被自动赋值。必须使用 PowerShell `Start-Process -PassThru` 获取子进程 PID。Windows 配方变体因此采用 PowerShell 内联命令。

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

```just
# just probe 配方骨架 — Linux/macOS 变体
# 退出码：0=健康，1=重试耗尽，2=端口冲突，3=连接超时
[linux]
probe:
    #!/usr/bin/env sh
    retries=0; max=30; interval=2; effective=0
    while [ $retries -lt $max ]; do
        status=$(curl -s -o /dev/null -w "%{http_code}" --connect-timeout 5 http://localhost:$PORT/healthz 2>/dev/null)
        case $? in
            0)  if [ "$status" -ge 200 ] 2>/dev/null && [ "$status" -lt 500 ] 2>/dev/null; then exit 0; fi ;;
            6|7)  ;; # ECONNREFUSED / 连接被拒绝 — 继续重试 ;;
            28) # curl 超时 — 连接超时计为 2 次重试
                retries=$((retries + 2)); effective=$((effective + 1))
                echo "[probe] [retry $((retries/2))/$((max/2))] :$PORT — 连接超时" >&2
                continue ;;
        esac
        # 检查 dev server 日志是否有端口冲突
        if grep -q "EADDRINUSE" .forge/dev-server.log 2>/dev/null; then
            echo "[probe] [abort] 端口冲突 — EADDRINUSE" >&2; exit 2
        fi
        retries=$((retries + 1))
        echo "[probe] [retry $retries/$max] :$PORT — 连接被拒绝" >&2
        sleep $interval
    done
    # 重试耗尽：根据超时模式选择退出码
    if [ $effective -gt 0 ]; then exit 3; fi
    exit 1

# just probe 配方骨架 — Windows 变体
[windows]
probe:
    #!powershell
    $retries = 0; $max = 30; $interval = 2; $timeouts = 0
    while ($retries -lt $max) {
        try {
            $resp = Invoke-WebRequest -Uri "http://localhost:$PORT/healthz" -UseBasicParsing -TimeoutSec 5
            if ($resp.StatusCode -lt 500) { exit 0 }
        } catch {
            $err = $_.Exception.Message
            if ($err -match "timeout|timed out") {
                $retries += 2; $timeouts++
                Write-Host "[probe] [retry $($retries/2)/$($max/2)] :$PORT — 连接超时"
                continue
            }
        }
        if (Test-Path ".forge\dev-server.log" -and (Select-String -Path ".forge\dev-server.log" -Pattern "EADDRINUSE" -ErrorAction SilentlyContinue)) {
            Write-Host "[probe] [abort] 端口冲突 — EADDRINUSE"; exit 2
        }
        $retries++
        Write-Host "[probe] [retry $retries/$max] :$PORT — 连接被拒绝"
        Start-Sleep -Seconds $interval
    }
    if ($timeouts -gt 0) { exit 3 }
    exit 1
```

> **Windows shebang 说明**：just 在 Windows 上处理 `#!` shebang 时，会查找 PATH 中的对应解释器。`#!powershell` 直接定位 PowerShell 可执行文件，不依赖 `/usr/bin/env`。Windows 10+ 和 Windows Server 2019+ 默认包含 PowerShell（`powershell.exe` for 5.x 或 `pwsh.exe` for 7.x），因此 `#!powershell` 在现代 Windows 上可用。若 PowerShell 不可用（极罕见的 Windows Nano Server 场景），just 会报错并跳过该配方——此时用户需安装 PowerShell 或使用 `[linux]` 变体在 WSL 中运行。此依赖已在 NFR"跨平台兼容"中隐含覆盖。

**参数配置来源**：
- `max_retries` 和 `interval` 通过环境变量覆盖（`PROBE_RETRIES`、`PROBE_INTERVAL`），默认值在 surface 规则文件中定义
- web surface 默认：30 次重试 / 2 秒间隔 / 检查根路径
- api surface 默认：30 次重试 / 2 秒间隔 / 检查 `/healthz`

**超时失败行为**：`just probe` 重试耗尽后以退出码退出（具体退出码取决于失败模式，见退出码处理表），run-tests 捕获后：(1) 执行 `just test-teardown` 清理已启动的 dev server；(2) 报告对应错误（超时/端口冲突/连接超时）；(3) 不执行 `just test`。

**probe 进度输出**：probe 配方体每次重试输出 `[probe] [retry N/M] <target> — <result>` 到 stderr，使用户在等待过程中看到实时进度。检测到端口冲突时立即输出 `[probe] [abort] 端口冲突` 并以 exit 2 退出（见退出码处理表），提供早期退出条件。


**已知局限——后台启动退出码不可靠**：`just dev` 后台启动后以 exit 0 返回，但 dev server 可能数秒内崩溃。提案采用**三层检测机制**：

1. **PID 存活检查**：probe 重试失败后检查 PID 存活（Linux `/proc/<pid>`，macOS `ps -p <pid> -o pid=`，Windows `tasklist`）。进程已死则立即退出报告崩溃
2. **端口监听检测**：PID 检查失败后连接目标端口（`curl -sf`/`nc -z`），端口可达说明 npm wrapper 退出但子进程存活，继续 probe
3. **日志尾部检查**：检查 `.forge/dev-server.log` 最后 5 行是否含 `EADDRINUSE` 等错误模式，匹配则立即退出报告端口冲突

dev server 在 probe 成功后崩溃由测试错误报告捕获。

#### Teardown 进程回收机制

teardown 必须可靠地回收所有测试启动的进程，即使测试中途崩溃：

1. **正常路径**：`just test-teardown` 读取 `.forge/dev-server.pid`，**先校验 PID 有效性和命令行匹配**，再向对应进程发送 SIGTERM/taskkill

2. **PID 有效性校验 + 命令行匹配 + 端口反查**（防止 PID 回收杀错进程，覆盖 npm wrapper 退出场景）：
   - **PID 校验**：Linux `/proc/<pid>`；macOS `ps -p <pid>`；Windows `tasklist`
   - **命令行校验**：Linux `/proc/<pid>/cmdline`；macOS `ps -p <pid> -o command= -w`（`-w` 禁用截断）；Windows `Get-CimInstance Win32_Process`（失败时回退 `tasklist /V`）
   - **端口反查**（PID 校验失败时触发）：通过监听端口定位实际进程 PID——Linux `lsof -i :$PORT -t`/`ss -tlnp`；macOS `lsof -i :$PORT -t`；Windows `netstat -ano | findstr :$PORT`。反查 PID 替代原 PID 作为终止目标
   - PID 无效且端口反查无结果时跳过终止，仅清理 PID 文件
   - **混合项目**：每 surface-key 独立 PID 文件（`.forge/dev-server.<surface-key>.pid`），teardown 逆序处理
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

- **v3.0.0 兜底机制（确定性下限）**：run-tests 从 surface 执行策略规则文件（`rules/surfaces/<type>.md`）读取编排序列，该规则在 SKILL.md 中定义为参数化的固定模板，LLM 按模板执行而非自由编排。关键约束通过 just 命令退出码强制执行——每个步骤的退出码决定下一步动作（0=继续，非0=触发 teardown 后中止）。即使 LLM 误读指令，just 命令的退出码是确定性的，不会因 LLM 行为改变
- **HARD-GATE 规则**：在 run-tests SKILL.md 中用 `<HARD-GATE>` 标记编排序列的强制顺序。违反顺序的行为列为"禁止"，与现有"禁止跳过失败测试"同等级别。新增以下 HARD-GATE 规则：**probe 失败后禁止重试 probe 或重试 dev——唯一允许的下一步是执行 teardown 后中止**。这防止 LLM 在 probe 失败后进入"重试 probe → 再失败 → 再重试"的无限循环，或试图重新执行 `just dev` 而不先清理已启动的进程
- **运行时强制保障**：除文本指令外，增加两层保障：(1) **test-state.json 状态门控**——run-tests 每步将 `{phase, status}` 写入状态文件，下一步执行前检查，若上一步 failed 且当前不是 teardown 则拒绝执行。(2) **teardown 幂等性**——PID 不存在时跳过，违反 HARD-GATE 执行不必要的 teardown 也不产生副作用
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
- **config schema 变更**（独立子方案，1-2 个任务）：移除 `test.execution` 文档、未知字段处理、per-surface `tests/{surface-key}/config.yaml` 格式定义
- **surface-key 统一迁移**（见"维度 2"迁移影响面）：surface-key-assignment 规则更新、quick-tasks surface-key 推断更新、db-schema 规则更新、prompt.go resolveSurfaceKey 更新、init-justfile 混合项目配方生成更新

预计 15-20 个编码任务。

### 依赖就绪度

Surface 检测已就位。`test.execution` 在 Go `Config` 结构体中未被映射为独立字段（代码验证：结构体无 `Execution` 子字段），但 run-tests SKILL.md 通过 `forge config get test.execution` 让 LLM agent 读取原始 YAML 并使用这些字段——即 `test.execution` 在 Go 层面未结构化实现，但在 LLM agent 层面实际在用。移除后需要确保 run-tests 的所有编排路径都更新为 surface 感知模式，config-schema.md 同步删除 `test.execution` 文档。`timeout` 和 `probe` 参数不再通过 `.forge/config.yaml` 传递，改由 per-surface 测试配置（`tests/{surface-key}/config.yaml`）提供，testkit 等辅助工具负责解析。


**test.execution 引用审计**（移除前的必要前置）：

移除 `test.execution` 前需审计所有 SKILL.md 对 `test.execution` 的引用，确保无其他 skill 静默依赖此配置。审计清单：

| Skill | 可能引用点 | 预期影响 |
|-------|-----------|---------|
| `fix-bug` | 验证修复后运行测试 | 已通过 `just test` 调用，不依赖 `test.execution` |
| `clean-code` | 重构后验证测试 | 已通过 `just test` 调用，不依赖 `test.execution` |
| `run-tests` | 编排序列 | **直接移除目标**，需全面更新 |
| `quality-gate` | 门控检查 | Go 代码直接调用 just，不依赖 `test.execution` |

审计通过 `grep -r "test.execution" plugins/forge/skills/` 执行。若发现除 `run-tests` 外的引用，需在对应 SKILL.md 中同步更新为 `just test` 调用路径，作为 surface-key 统一迁移的附加任务。

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

  **多 surface 同类型的 journey 过滤规则**：当混合项目包含多个同类型 surface（如 `admin-panel: web` + `marketing-site: web`）时，`just test e2e` 运行所有 web surface 的 e2e 测试（按 surface 维度聚合而非全局去重）。如果用户只需测试特定 surface 的某个 journey，需通过 surface-key 参数指定：`just test admin-panel e2e`。journey 标签映射在各 surface 的规则文件中独立定义，`just test <journey>` 等价于对所有 surface 执行该 journey 的测试。

  **参数解析优先级**（仅多 surface 项目适用）：多 surface 项目的 `test` 配方签名为 `just test [surface-key] [journey]`，just 按位置绑定参数。当只传入一个位置参数时（如 `just test admin-panel`），just 将其绑定为 surface-key——配方体需做反转判断：(1) 先与所有 surface-key 集合匹配——若匹配成功，视为 surface-key 参数，journey 为空（运行该 surface-key 的全部测试）；(2) 若不匹配任何 surface-key，视为 journey 参数，surface-key 为空（运行所有 surface 的该类型测试）。当传入两个位置参数时（如 `just test admin-panel e2e`），第一个为 surface-key，第二个为 journey。**单 surface 项目**签名为 `just test [journey]`，无反转逻辑——参数直接绑定为 journey。

  journey 过滤通过 `test` 配方的参数实现：多 surface 项目 `just test [surface-key] [journey]` 中 journey 为空时运行全部，否则按 surface 规则中定义的标签映射过滤；surface-key 为空时覆盖所有 surface。单 surface 项目 `just test [journey]` 中 journey 为空时运行全部。具体过滤逻辑由 justfile 配方体中的 `if/else` 或 `test runner` 的标签机制实现。
- SKILL.md 更新：新增 surface 检测步骤，surface 感知配方生成
- CLI/TUI 只生成 `dev`，不生成 `run`
- 去掉 Step 3a 中对 `test.execution.run` 的依赖

**run-tests**：
- SKILL.md 简化：改为调度器模式，检测 surface type 后加载对应执行策略规则
- 根据 surface 编排模式决定执行序列（是否启动 dev、是否 probe）
- 执行策略规则文件（`rules/surfaces/<type>.md`）定义编排序列和 just 配方调用契约

**配置 schema 变更（独立子方案）**：

此变更改变了 Forge CLI 的配置接口契约，影响面超出常规 Go 开发。单独列出以确保边界明确：

> **子方案被拒时的降级路径**：如果 config schema 变更子方案未通过评审，主提案仍可执行，但需要以下调整：(1) timeout 在 run-tests SKILL.md 中硬编码为默认值 300 秒，不通过任何配置文件读取——功能性不受影响，但用户无法自定义 timeout；(2) `test.execution` 的移除仍然执行（Go 结构体未映射此字段，移除仅涉及文档和 SKILL.md），`GetConfigValue` 对 `test.execution.*` 的查询继续返回"未找到"（当前行为）；(3) 降级路径的功能性损失仅限于"用户无法自定义 timeout 和 probe 参数"，这不影响编排序列的正确性。config schema 变更可在后续版本独立完成，不阻塞 v3.0.0 的 surface 感知和编排简化。

- **移除 `test.execution` 节点**：config-schema.md 中删除 `test.execution` 相关文档（Go 结构体本身未映射此字段，无需代码变更）。用户 config.yaml 中残留的 `test.execution` 节点通过以下策略处理：
  - Go Config 结构体使用 `yaml:"-"` 标签或 `squelch` 机制忽略未知字段，确保残留的 `test.execution` 节点不会导致 YAML 反序列化报错
  - 若 Go YAML 库默认不允许未知字段（如 `yaml.UnmarshalStrict`），则需切换为宽松模式或显式添加 `Execution` 字段标记为 `yaml:"execution" json:"-"`（仅用于兼容性解析，不暴露给业务逻辑）
  - `GetConfigValue` 对 `test.execution.*` 键的查询返回明确提示："test.execution 已移除，请使用 just 配方替代"
- **不新增 `.forge/config.yaml` 的 `test.*` 配置**：`timeout` 和 `probe` 参数由 per-surface 测试配置（`tests/{surface-key}/config.yaml`）提供，testkit 等辅助工具负责解析。`.forge/config.yaml` 不再承载测试相关配置
- **边界约束**：`surfaces` 字段保持现有 schema 不变；`GetConfigValue` 无需扩展 `test.*` 键支持
- **影响面评估**：涉及 `forge-cli/internal/config/`（未知字段处理）、`plugins/forge/references/run-tests/config-schema.md`（文档更新）2 个模块，预计 1-2 个独立任务

**surface-key 统一迁移**：

- `breakdown-tasks` `rules/surface-key-assignment.md`：文件路径分类改为按 surfaces 路径归属推断（见"维度 2"迁移细则）
- `quick-tasks` SKILL.md：内联 surface-key 推断同步更新
- `breakdown-tasks` `rules/db-schema.md`：`surface-key: 'backend'` 改为动态判定 type=api 的 surface-key
- `prompt.go` `resolveSurfaceKey()`：校验值域从 `frontend`/`backend` 改为 surface-key 集合
- `init-justfile` SKILL.md：混合项目所有配方的 surface-key case 分支改为 surface-key
- **Task Go struct**：`Scope` 字段重命名为 `SurfaceKey`，新增 `SurfaceType` 字段；旧 `scope` JSON 键保留一个版本用于兼容（v3.1.0 移除）
- **任务模板 frontmatter**：`scope` 重命名为 `surface-key`，新增 `surface-type` 字段（`breakdown-tasks/templates/task.md`、`quick-tasks/templates/task.md`）
- **任务生成逻辑**：`breakdown-tasks` 和 `quick-tasks` 生成任务时，同时填充 `surface-key`（从 surface-key-assignment 规则推断）和 `surface-type`（从 config.yaml 的 surfaces map 按 surface-key 查类型）
- **`forge task add` CLI**：从 `SourceTaskID` 指向的源任务继承 `surface-key` 和 `surface-type`；无源任务时单 surface 项目 surface-key 为空（surface-type 填充唯一 surface 的类型），多 surface 项目均默认空。**不新增 CLI flag**——继承优先于手动指定
- **quality-gate fix-task**：`addFixTask()` 复用 `resolveSurfaceKey()` 从失败文件路径推断 surface-key，查 surfaces map 填充 surface-type；无法推断时默认空
- **auto-gen test tasks**：`autogen.go` 中 `AutoGenTaskDef.Scope` 重命名为 `SurfaceKey`（保持 `"all"`），新增 `SurfaceType` 从 per-type 任务定义提取（如 `T-test-gen-journeys-api` → surface-type = `api`）

**通用**：
- 向后兼容：无 surface 配置 → 当前行为不变
- 混合项目：surface 规则按 surface-key 应用

### 范围外

- 变更语言模板（`templates/*.just`）
- 变更 `forge-cli/pkg/just/` 门控序列
- 变更 `forge-cli/internal/cmd/quality_gate.go` 或 `testrunner` 的 Go 代码
- 新增 forge CLI 命令（`config get` 不新增 `test.*` 键支持，timeout 由 per-surface 测试配置提供）
- **回滚基础设施**（feature flag 机制不在范围内；回滚通过 git revert 实现）
- Go 代码子命令直接管理进程（长期方向，非 v3.0.0 范围）

## 主要风险

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|---------|
| Surface 未检测到 | 中 | 低 | 回退到当前行为 |
| run-tests 简化导致已配置 `test.execution` 的项目不兼容 | 低 | 低 | v3.0.0 未发布，无存量用户 |
| Surface 规则过于泛化 | 中 | 中 | LLM 组合语言模板 + surface 规则 |
| 混合项目 surface 歧义 | 中 | 中 | `forge surfaces` 基于路径检测；按 surface-key 应用 |
| `test [journey]` 过滤与原生运行器不兼容 | 中 | 高 | Surface 规则记录映射关系 |
| run-tests 无法感知 surface（skill 内无 config 读取） | 低 | 高 | run-tests 通过 `forge surfaces` 或 `.forge/config.yaml` 获取 surface 类型 |
| HARD-GATE 规则被 LLM 违反（如 probe 失败后重试而非 teardown） | 中 | 高 | **分层兜底机制**：(1) 结构化约束——编排序列由执行策略规则文件定义（LLM 按规则执行而非自由编排），降低自由度；(2) 退出码强制门控——just 命令退出码是确定性的（非 LLM 生成），每个步骤的退出码决定下一步动作，即使 LLM 尝试跳过步骤，非零退出码仍会触发 abort；(3) 执行策略规则文件作为外部约束——编排序列从规则文件读取而非由 LLM 推断，减少 LLM 误判空间；(4) 如果 LLM 仍违反 HARD-GATE（如忽略退出码直接继续），最坏情况为"执行不必要的 teardown 后中止"——不会导致数据损坏或不可恢复状态，因为 teardown 本身是幂等操作。**注意**：HARD-GATE 规则严格禁止 probe 失败后的任何重试行为（包括 teardown 后重试 dev/probe），最坏情况的"不必要 teardown"仅指 LLM 可能先执行了完整的 teardown 再中止（而非跳过 teardown 直接重试） |
| `# user-customized` 导致用户错过 surface 规则改进 | 中 | 中 | **缓解措施**：(1) init-justfile 每次运行时输出差异摘要（逐行对比当前自定义配方与最新生成版本的变更点），即使跳过覆盖也告知用户有哪些改进被跳过；(2) 差异摘要按类型分类（"优化类"、"新功能类"、"bug 修复类"），帮助用户判断同步优先级；(3) `--force-regenerate` 提供一键更新路径；(4) 长期风险：用户可能持续忽略差异摘要，导致 justfile 逐渐与推荐模板产生较大偏差——此风险通过差异摘要的显式输出（非静默跳过）控制为"知情选择"而非"意外偏离" |

### 回滚计划

如果 surface 感知方案上线后发现严重问题（如某种 surface 类型的编排在特定框架下不可用），按以下策略回滚：

1. **回滚方式**：通过 `git revert` 回退提案的所有变更（SKILL.md、surface 规则文件、config-schema.md、init-justfile 模板）。不引入 feature flag 机制——回滚基础设施不在 v3.0.0 范围内
2. **回滚影响面**：回滚后 `run-tests` 恢复到当前行为（读取 `test.execution` 委托路径），已生成的 surface 感知 justfile 中 dev/probe/test 配方仍然可用（它们只是 just 配方，不依赖 run-tests 的编排逻辑）

## 成功标准

- [ ] init-justfile 为 web/api/cli/tui/mobile 5 种 surface 生成差异化的配方组合
- [ ] CLI/TUI 不生成 `run` 配方，统一使用 `dev`
- [ ] `run-tests` 不再依赖 `test.execution.run`，直接调用 `just test`：单 surface 项目调用 `just test [journey]`，多 surface 项目调用 `just test [surface-key] [journey]`
- [ ] `run-tests` 根据 surface 编排模式决定是否启动 dev/probe
- [ ] `test.execution` 节点从 config-schema 中完全删除，`.forge/config.yaml` 不再承载任何 `test.*` 配置；timeout 和 probe 参数由 `tests/{surface-key}/config.yaml`（单 surface 为 `tests/config.yaml`）提供
- [ ] 无 surface 配置的项目输出与当前一致
- [ ] 所有生成的配方通过 `--dry-run` 验证（语法正确、配方名和参数签名符合 Standard Target Contract 定义）。注意：`--dry-run` 仅验证配方语法和依赖解析，不验证运行时行为（如 dev server 是否实际启动、probe 目标是否可达）。运行时行为通过端到端验证标准（下一项）覆盖
- [ ] 运行时端到端验证（至少一个 web/api 项目）：`just dev` 后台启动 dev server 成功 → `just probe` 检测到服务就绪 → `just test` 执行测试并通过 → `just test-teardown` 清理进程，无孤儿进程残留
- [ ] 每个 surface 规则文件记录了测试编排模式和 journey 过滤策略
- [ ] 语言模板与 surface 规则的配方职责边界验证：生成的 justfile 中，`compile`/`build`/`lint`/`fmt`/`unit-test` 配方来源于语言模板，`test`/`dev`/`run`/`probe`/`test-setup`/`test-teardown` 配方来源于 surface 规则。验证方式：对生成的 justfile 执行 `grep -c` 确认两集合无交集（无同名配方同时出现在两个来源中）
- [ ] 混合项目（web+api）端到端验证：`just dev` 无 surface-key 时按依赖顺序启动所有 surface-key 的 dev server，各 surface-key 的 probe 依次通过，测试执行完成后 teardown 逆序清理所有进程。PID 文件有效性校验覆盖混合项目的多 PID 场景
- [ ] surface-key 值域统一迁移：混合项目所有配方的 surface-key 参数值从 `frontend`/`backend` 迁移为 surface-key，surface-key-assignment、prompt.go resolveSurfaceKey、init-justfile 配方生成同步更新
- [ ] Task 数据模型扩展验证：Go struct 新增 `SurfaceType` 字段，`Scope` 重命名为 `SurfaceKey`；任务模板 frontmatter 包含 `surface-key` 和 `surface-type`；`breakdown-tasks` / `quick-tasks` 生成的任务中两个字段正确填充；`forge task add` 从源任务继承 `surface-key` 和 `surface-type`（无源任务时单 surface 项目 surface-key 为空，不新增 CLI flag）；quality-gate fix-task 从失败文件路径自动推断 surface-key/type；auto-gen test tasks 的 `SurfaceType` 从任务定义提取；旧任务的 `scope` 字段通过兼容层正确映射
- [ ] 测试目录结构验证：多 surface 项目测试用例放入 `tests/{surface-key}/`（surface-key = surfaces map key 中 `/` 替换为 `-`）；单 surface 项目测试用例放入 `tests/`（无子目录）；`results-dir` 路径与目录结构一致
- [ ] config schema 变更验证：`.forge/config.yaml` 中无 `test.*` 配置项；用户 config.yaml 中残留的 `test.execution` 节点不导致 YAML 解析错误；`GetConfigValue` 对 `test.execution.*` 的查询返回明确的移除提示；SKILL.md 中无任何 timeout 参数配置；testkit 能正确解析 per-surface config 并向测试脚本暴露配置项（如 web 的 baseUrl/credentials、api 的 authToken 等）
- [ ] `# user-customized` 保护机制验证：在生成的 justfile 中标记 `# user-customized` 的配方重新运行 init-justfile 后不被覆盖；差异摘要正确输出当前配方与新生成版本的变更点；`--force-regenerate` 参数正确覆盖所有 `# user-customized` 配方

## 下一步

- 继续执行 `/quick-tasks` 生成实现任务
