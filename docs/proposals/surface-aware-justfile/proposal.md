---
created: 2026-05-24
author: "faner"
status: Draft
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

justfile 本身已经是抽象层，config.yaml 再包一层增加了复杂度但没有增加灵活性 — 所有模板变量最终都解析为 just 命令。

### 证据

- Web UI 的 e2e 测试必须先启动应用，但当前配方没有 surface 特定的启动逻辑
- `run-tests` 的编排序列对 surface 不感知，只能依赖用户手动配置 `test.execution`
- `test.execution` 的所有示例都指向 just 命令（`just test {slug}`、`just test-setup`、`just probe`），委托层只是转发

### 紧迫性

随着 v3.0.0 test profile 的引入，surface 成为测试流程的核心维度。init-justfile 和 run-tests 需要协同工作 — init-justfile 生成正确的配方，run-tests 编排正确的执行序列。

## 提议方案

### 方案 A：Surface 感知 + 委托层简化

1. init-justfile 添加 surface 感知层，为不同 surface 生成差异化的配方
2. 去掉 `test.execution` 委托层，`run-tests` 直接调用 just 配方
3. `timeout`、`results-dir` 等非命令配置保留在 config.yaml 的简化字段中

### Surface 测试编排模式

| Surface | 测试编排序列 | 关键配方 |
|---------|------------|---------|
| **web** | `just dev`（后台）→ `just probe`（等待就绪）→ `just test` → teardown | `dev` 启动 dev server；`probe` 检查 HTTP 端点 |
| **api** | `just dev`（后台）→ `just probe`（等待就绪）→ `just test` → teardown | `dev` 启动 API 服务；`probe` 检查 /healthz |
| **cli** | `just build` → `just dev`（验证二进制可运行）→ `just test` | 无需 `run`/`probe`；`test` 通过子进程测试 |
| **tui** | `just build` → `just dev`（验证二进制可运行）→ `just test` | 无需 `run`/`probe`；`test` 通过 stdin 管道测试 |
| **mobile** | `just dev`（启动模拟器+应用）→ `just test` → teardown | `test-setup` 准备模拟器；`test` 用 maestro |

**dev/run 分工**：
- **web/api**：保留 `dev`（热重载 dev server）和 `run`（生产模式启动），测试编排使用 `dev`
- **cli/tui**：只生成 `dev`（编译+运行二进制），不生成 `run`
- **mobile**：只生成 `dev`（启动模拟器+应用），不生成 `run`
- **混合项目（web+api）**：`dev` 配方接受 scope 参数 — `just dev frontend`、`just dev backend`、`just dev`（无 scope 时并发启动两者）

### 委托层简化

去掉 `test.execution` 后，`run-tests` 的编排变为：

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

### Skill 流程

init-justfile 的生成流程：
```
检测语言 → 加载语言模板 → 生成 compile/build/lint/fmt
检测 surface → 加载 surface 规则 → 生成 test/dev/run/probe/test-setup
组装 → 验证
```

run-tests 的执行流程：
```
检测 surface → 确定编排模式
→ just test-setup
→ just dev（web/api，后台启动）
→ just probe（web/api，等待就绪）
→ just test [journey]
→ just test-teardown（如有）
```

### 创新亮点

将 init-justfile（配方生产者）和 run-tests（配方消费者）的 surface 感知统一设计。init-justfile 生成正确的配方组合，run-tests 按 surface 编排执行序列，去掉中间的 config 委托层。justfile 成为唯一的命令抽象层。

## 需求分析

### 关键场景

- **Web+Node 项目**：`dev` 启动 next dev（端口 3000）→ `probe` curl localhost:3000 → `test` playwright e2e → teardown
- **API+Go 项目**：`dev` 启动 API server → `probe` curl /healthz → `test` HTTP 契约测试 → teardown
- **CLI+Rust 项目**：`build` 编译二进制 → `dev` 运行验证 → `test` 子进程集成测试
- **TUI+Go 项目**：`build` 编译二进制 → `dev` 运行验证 → `test` stdin 管道测试
- **Mobile+TS 项目**：`test-setup` 准备模拟器 → `dev` 启动应用 → `test` maestro YAML
- **无 surface 配置**：回退到当前行为（纯语言配方，run-tests 保持原有逻辑）
- **混合项目（web+api）**：`just dev` 无 scope 时须同时启动前端和后端（`just dev frontend` + `just dev backend` 并发），probe 分别检查两个服务

### 下游集成契约

配方签名不可变。`run-tests` 的调用方式简化但语义不变：

| 配方 | 签名（不可变） | 消费者 | 期望语义 |
|------|--------------|--------|---------|
| `unit-test` | `just unit-test` | forge task submit、clean-code、fix-bug、testrunner | 语言级单元测试；exit 0 = 通过 |
| `test` | `just test [journey]` | run-tests、forge quality-gate、fix-bug | Surface 级测试 |
| `probe` | `just probe` | run-tests（web/api 前置检查） | 服务健康检查；exit 0 = 健康；混合项目检查所有服务 |
| `test-setup` | `just test-setup` | run-tests | 安装测试依赖；幂等 |
| `test-teardown` | `just test-teardown` | run-tests（可选） | 测试后清理 |
| `dev` | `just dev [scope]` | run-tests（web/api 启动服务） | web/api: 后台启动并监听端口；cli/tui: 编译运行 |
| `run` | `just run [scope]` | 仅 web/api 生成 | 生产模式启动 |

**关键约束**：
- `test` 必须始终接受 `journey=''` 参数
- web/api 的 `dev` 须能在后台运行
- cli/tui 不生成 `run` 配方
- `test-teardown` 为可选配方，不存在时 run-tests 跳过

### 约束与依赖

- Surface 信息来自 `.forge/config.yaml` 或 `forge surfaces` CLI
- Surface 规则保持语言无关
- Standard Target Contract 配方名称和签名不变
- 遵循 forge-distribution.md 路径约定

## 替代方案

| 方案 | 优势 | 劣势 | 结论 |
|------|------|------|------|
| 不做 | 零成本 | 下游无法 surface 感知编排；test.execution 委托层继续累积复杂度 | 拒绝 |
| 仅 surface 感知，保留 test.execution | 改动小 | 委托层冗余未解决；两处配置 surface 行为 | 拒绝：治标不治本 |
| **surface 感知 + 去掉 test.execution** | 统一抽象层；消除冗余委托 | 改动范围扩大到 run-tests | **选定** |

## 可行性评估

### 技术可行性

直接可行。init-justfile 已有 convention 驱动的测试生成。Surface 检测增加一次 config 读取。`run-tests` 简化为直接调用 just 配方。

### 资源与时间

中等范围：
- 5 个 surface 规则文件
- init-justfile SKILL.md 更新
- run-tests SKILL.md 简化（去掉 test.execution 读取）
- config-schema.md 更新
- init-justfile Step 3a 去掉 test.execution 依赖

预计 8-12 个编码任务。

### 依赖就绪度

Surface 检测已就位。`test.execution` 在 Go 结构体中从未实现，移除无影响。`GetConfigValue` 扩展为常规 Go 开发。

## 假设挑战

| 假设 | 挑战工具 | 发现 |
|------|---------|------|
| `test.execution` 提供灵活性 | XY 检测 | 推翻：实际所有示例都指向 just 命令。真正的灵活性在 justfile 配方中，不在 config 委托层 |
| Surface 规则应包含语言特定指导 | 假设翻转 | 推翻：用户确认保持语言无关 |
| 只需优化单个配方 | 压力测试 | 推翻：核心是**编排序列**不同 |
| `dev`/`run` 只用于人工开发 | XY 检测 | 推翻：run-tests 编排测试时需要 `dev` 启动被测服务 |
| CLI/TUI 需要 `run` 配方 | Occam's Razor | 推翻：CLI/TUI 无服务启动概念，`dev`（编译+运行）足够 |

## 范围

### 范围内

**init-justfile**：
- 5 个 surface 规则文件：`skills/init-justfile/rules/surfaces/{web,api,cli,tui,mobile}.md`
  - 每个包含：测试编排模式、配方生成指导、journey 过滤策略
- SKILL.md 更新：新增 surface 检测步骤，surface 感知配方生成
- CLI/TUI 只生成 `dev`，不生成 `run`
- 去掉 Step 3a 中对 `test.execution.run` 的依赖

**run-tests**：
- SKILL.md 简化：直接调用 just 配方，去掉 `test.execution` 读取
- 根据 surface 编排模式决定执行序列（是否启动 dev、是否 probe）

**配置（forge-cli）**：
- Config 结构体新增 `TestConfig` 节点（`timeout`、`results-dir`）
- `GetConfigValue` 支持 `test.*` 键读取（`forge config get test.timeout`、`forge config get test.results-dir`）
- Surface 信息通过现有 `forge surfaces [path]` 获取，无需扩展 `forge config get`
- 移除 run-tests/config-schema.md 中的 `test.execution` 引用
- init-justfile SKILL.md 的 `test` 配方生成不再参考 `test.execution.run`

**通用**：
- 向后兼容：无 surface 配置 → 当前行为不变
- 混合项目：surface 规则按 scope 应用

### 范围外

- 变更语言模板（`templates/*.just`）
- 变更 `forge-cli/pkg/just/` 门控序列
- 变更 `forge-cli/internal/cmd/quality_gate.go` 或 `testrunner` 的 Go 代码
- 新增 forge CLI 命令（仅扩展 `config get` 键支持）

## 主要风险

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|---------|
| Surface 未检测到 | 中 | 低 | 回退到当前行为 |
| run-tests 简化导致已配置 `test.execution` 的项目不兼容 | — | — | v3.0.0 未发布，无存量用户，直接移除 |
| Surface 规则过于泛化 | 中 | 中 | LLM 组合语言模板 + surface 规则 |
| 混合项目 surface 歧义 | 中 | 中 | `forge surfaces` 基于路径检测；按 scope 应用 |
| `test [journey]` 过滤与原生运行器不兼容 | 中 | 高 | Surface 规则记录映射关系 |
| run-tests 无法感知 surface（skill 内无 config 读取） | 低 | 高 | run-tests 通过 `forge surfaces` 或 `.forge/config.yaml` 获取 surface 类型 |

## 成功标准

- [ ] init-justfile 为 web/api/cli/tui/mobile 5 种 surface 生成差异化的配方组合
- [ ] CLI/TUI 不生成 `run` 配方，统一使用 `dev`
- [ ] `run-tests` 不再依赖 `test.execution.run`，直接调用 `just test [journey]`
- [ ] `run-tests` 根据 surface 编排模式决定是否启动 dev/probe
- [ ] `test.execution` 节点从 config-schema 中完全移除，非命令字段移至 `test:` 顶层
- [ ] 无 surface 配置的项目输出与当前一致
- [ ] 所有生成的配方通过 `--dry-run` 验证
- [ ] 每个 surface 规则文件记录了测试编排模式和 journey 过滤策略

## 下一步

- 继续执行 `/quick-tasks` 生成实现任务
