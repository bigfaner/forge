---
iteration: 1
title: "CTO Rubric Scoring — Iteration 1"
date: 2026-05-25
scorer: CTO adversary
---

# 评分报告：init-justfile Surface 感知 + 测试编排简化

## Phase 1 — 推理链审计

### 论证链路

1. **问题**：init-justfile 不感知 surface → 测试编排流程无法区分 web/api/cli/tui/mobile
2. **补充问题**：test.execution 委托层冗余 → 4 层间接转发
3. **方案**：surface 感知 + 去掉委托层 → justfile 成为唯一抽象层
4. **证据**：表面合理但缺乏量化数据支撑
5. **成功标准**：以 checklist 形式覆盖主要交付物

### 自相矛盾检测

- **矛盾 1**：提案声称"去掉 `test.execution` 委托层，`run-tests` 直接调用 just 配方"，但又在废弃行为中要求"检测到 `test.execution` 节点时输出废弃警告"——这意味着 `run-tests` 仍然需要**读取** config.yaml 来检测 `test.execution` 节点。这不是"去掉"，而是"读取但不使用"。委托逻辑被替换为检测逻辑，复杂度并未完全消除。
- **矛盾 2**：提案说 `test.execution` 在 Go 结构体中"从未实现"（属实，代码验证 `Config` 结构体无此字段），但 config-schema.md 文档中明确定义了该字段的 schema，且 run-tests SKILL.md 的 Step 1 要求读取 `test.execution`。所谓的"从未实现"只是 Go 层面未做结构化映射——LLM agent 实际上通过 `forge config get test.execution` 解析原始 YAML 来使用这些字段。提案低估了实际使用面。
- **矛盾 3**：混合项目的 scope 参数声明"值必须是 `surfaces` map 的 key"，但 run-tests 的编排模式表格只区分 surface **类型**（web/api/cli/tui/mobile），不区分 scope key。当项目有 `admin-panel: web` + `payment-service: api` 时，run-tests 需要同时执行 web 编排和 api 编排——但这两种编排都需要启动 dev server 并 probe，如何并发管理？提案未讨论。

---

## Phase 2 — 评分

### 1. 问题定义 (80/110)

| 子项 | 得分 | 理由 |
|------|------|------|
| 问题清晰 | 35/40 | 两个问题（surface 不感知 + 委托层冗余）陈述明确，容易理解。扣分点：问题 1 和问题 2 的因果关系未明确——它们是独立的还是必须捆绑解决？ |
| 证据 | 20/40 | **严重不足**。声称"Web UI 的 e2e 测试必须先启动应用，但当前配方没有 surface 特定的启动逻辑"——但没有提供任何一个实际项目的 justfile 作为证据。声称"`test.execution` 的所有示例都指向 just 命令"——但 config-schema.md 中有 `go test` 和 `npx vitest` 的示例（不以 just 开头），这个声称被代码事实反驳。证据缺乏量化：有多少用户受影响？有多少项目的 config 中使用了非 just 命令？ |
| 紧迫性 | 25/30 | 与 v3.0.0 test profile 对齐的时机论证合理。但"cost of delay"分析缺失——如果不做会怎样？现有 test.execution 方案能撑多久？ |

### 2. 方案清晰度 (92/120)

| 子项 | 得分 | 理由 |
|------|------|------|
| 方法具体 | 38/40 | Surface 编排模式表格是本文档的亮点——5 种 surface 的编排序列一目了然。读者可以准确复述方案内容。 |
| 用户行为描述 | 35/45 | init-justfile 的用户行为描述较好（生成什么配方）。但 run-tests 的用户行为描述不足——用户需要理解"不再配置 test.execution"意味着什么迁移步骤。废弃警告是写给谁的？AI agent 还是人类开发者？ |
| 技术方向 | 19/35 | **方向清晰但深度不足**。`test 配方生成 fallback 链`是好的设计，但没有讨论：(1) `just dev` 的后台运行机制（just 本身不原生支持后台运行，需要 `&` 或 `tmux`）；(2) `just probe` 的重试逻辑（端口可用前需要轮询）；(3) teardown 中如何可靠地杀掉后台 dev 进程。这些是"测试编排"的核心技术难点，完全未触及。 |

### 3. 行业对标 (65/120)

| 子项 | 得分 | 理由 |
|------|------|------|
| 行业方案引用 | 10/40 | **零引用**。没有提及任何行业工具的类似设计——Docker Compose 的 health check + depends_on、Kubernetes 的 init container + readinessProbe、GitHub Actions 的 service container、Makefile 的 target dependency graph。这些都是"启动服务 → 等待就绪 → 运行测试"的成熟解决方案。 |
| 替代方案 | 20/30 | 列了 3 个替代方案（不做 / 仅 surface / surface+去掉委托），但都属于"做 vs 不做"和"做多少"的增量选择。缺乏根本不同的架构替代——例如：为什么不让 run-tests 自己管理进程生命周期（像 Cypress 的 `start-server-and-test` 那样），而要委托给 just？ |
| 权衡对比 | 20/25 | justfile 作为唯一抽象层的 trade-off 分析是亮点——明确列出了局限性和缓解措施。 |
| 选定方案论证 | 15/25 | "去掉委托层"的理由（所有示例都指向 just 命令）被 config-schema.md 中的非 just 示例削弱。论证应诚实承认：某些项目可能直接用 go test / npx vitest，这些路径在简化方案中被牺牲了。 |

### 4. 需求完整性 (70/110)

| 子项 | 得分 | 理由 |
|------|------|------|
| 场景覆盖 | 30/40 | 7 个关键场景覆盖了 5 种 surface + 无 surface + 混合。遗漏：(1) `just dev` 启动失败（端口被占、依赖缺失）；(2) `just probe` 超时（服务启动慢）；(3) 多个 surface 的同一类型（如两个 web surface 并发启动）；(4) 已有 justfile 但需要覆盖更新的场景 |
| 非功能需求 | 15/40 | **严重缺失**。没有讨论：(1) 性能——5 个 surface 规则文件的加载耗时？(2) 安全——`just dev` 后台进程的权限和资源泄漏？(3) 兼容性——just 版本要求（1.50.0 够吗？后台运行需要什么版本？）(4) 可观测性——编排过程中用户看到什么输出？ |
| 约束与依赖 | 25/30 | Surface 信息源优先级规则（config.yaml > forge surfaces CLI > 冲突以 config 为准）是好的。但"GetConfigValue 扩展为 config schema 子方案的一部分，需独立评审"——这是一个关键依赖，却只用了半句话带过。 |

### 5. 方案创新性 (55/100)

| 子项 | 得分 | 理由 |
|------|------|------|
| 新颖度 | 20/40 | "surface 感知配方生成"是 forge 特有的设计，但不是新概念——Docker Compose 和 Kubernetes 都有类似的"类型感知编排"。真正的创新点（init-justfile 和 run-tests 的双向 surface 感知统一）被提出但未深入：统一体现在哪里？共享了什么数据结构或接口？ |
| 跨域灵感 | 15/35 | 没有证据表明借鉴了其他领域的设计。提案中的编排模式本质上是传统的"启动-等待-测试-清理"流水线。 |
| 洞察简洁度 | 20/25 | "justfile 已经是抽象层，config 再包一层只是转发"——这个洞察是简洁有力的，是整个提案的逻辑支点。 |

### 6. 可行性 (60/100)

| 子项 | 得分 | 理由 |
|------|------|------|
| 技术可行性 | 25/40 | 声称"直接可行"但回避了关键技术难点：(1) `just dev` 后台运行——just 本身不原生支持后台运行模式，需要 bash `&` + PID 追踪，这在跨平台（Windows）上有严重问题；(2) `just probe` 轮询——需要超时+重试逻辑，放在 just 配方体中会很丑陋；(3) teardown 杀后台进程——跨平台信号处理差异大。 |
| 资源与时间 | 15/30 | "10-15 个编码任务"的估算偏乐观。config schema 变更被单独列为"独立子方案"，但如果这个子方案受阻（GetConfigValue 扩展需要兼容旧配置），整个提案的时间线会延后。 |
| 依赖就绪度 | 20/30 | Surface 检测就位是事实。但 `test.execution` "从未实现"的说法过于简化——LLM agent 通过原始 YAML 读取使用它，移除后需要确保所有 agent 路径都更新。 |

### 7. 范围定义 (60/80)

| 子项 | 得分 | 理由 |
|------|------|------|
| 范围内 | 25/30 | 每个范围内项都是可交付的。5 个 surface 规则文件 + SKILL.md 更新 + config schema 变更，粒度合理。 |
| 范围外 | 15/25 | 范围外列了 4 项，但缺少关键排除：(1) 现有项目的迁移指南（如何从 test.execution 过渡？）；(2) CI/CD 环境的适配。 |
| 范围边界 | 20/25 | "向后兼容：无 surface 配置 → 当前行为不变"是好的约束。但 config schema 子方案的边界模糊——它需要独立评审，但又在此提案的范围内。 |

### 8. 风险评估 (60/90)

| 子项 | 得分 | 理由 |
|------|------|------|
| 风险识别 | 20/30 | 列了 6 个风险，覆盖了主要场景。遗漏：(1) Windows 平台兼容性（后台进程、信号处理）；(2) 多 surface 并发启动的资源竞争；(3) just 配方体的跨平台可移植性。 |
| 可能性+影响 | 20/30 | `test [journey]` 过滤与原生运行器不兼容的风险标记为"中/高"是合理的。但"run-tests 无法感知 surface"标记为"低/高"——如果 run-tests 是通过 LLM agent 执行的 SKILL，它感知 surface 的能力取决于 LLM 是否正确读取 config，这不是"低"可能性。 |
| 缓解措施 | 20/30 | "LLM 组合语言模板 + surface 规则"不是一个可操作的缓解措施——它是设计本身。真正的缓解应该是"当 LLM 生成的配方不正确时，有自动化验证机制捕获"。`--dry-run` 验证算一个，但 dry-run 不覆盖运行时行为。 |

### 9. 成功标准 (55/80)

| 子项 | 得分 | 理由 |
|------|------|------|
| 可测量可测试 | 35/55 | 10 条标准中有 6 条是明确的（生成差异化的配方、不生成 run、不再依赖 test.execution.run 等）。但"所有生成的配方通过 `--dry-run` 验证"——dry-run 只验证语法，不验证运行时行为（如 probe 是否真的能等到服务就绪）。缺少端到端的集成测试标准。 |
| 覆盖完整性 | 20/25 | 成功标准覆盖了范围内的主要交付物。但缺少：(1) config schema 变更的成功标准（GetConfigValue 扩展是否通过测试？）；(2) 性能基准（surface 规则加载是否影响 init-justfile 的响应时间？） |

### 10. 逻辑一致性 (65/90)

| 子项 | 得分 | 理由 |
|------|------|------|
| 方案解决问题 | 28/35 | Surface 感知解决"测试编排流程不同"是成立的。去掉 test.execution 委托层解决"冗余"也是成立的。但两者捆绑的必要性论证不够——为什么不只做 surface 感知而保留 test.execution？提案说"治标不治本"但没有充分论证为什么保留委托层就是"治标"。 |
| 对齐 | 20/30 | 范围 ↔ 方案基本对齐，但成功标准 ↔ 方案有缝隙：方案中详细描述了"混合项目 scope 并发启动"，但成功标准中只有"无 surface 配置的项目输出与当前一致"——混合项目的成功标准缺失。 |
| 需求↔方案一致性 | 17/25 | 下游集成契约表格与方案设计一致。但"语言模板与 surface 规则的配方职责边界"——这个仲裁规则在方案中定义了（surface 优先覆盖 test/dev/run/probe），但需求分析中未列出"当语言模板已生成 test 配方时，surface 规则覆盖它的用户影响"作为一个场景。 |

---

## 评分汇总

| 维度 | 得分 | 满分 |
|------|------|------|
| 问题定义 | 80 | 110 |
| 方案清晰度 | 92 | 120 |
| 行业对标 | 65 | 120 |
| 需求完整性 | 70 | 110 |
| 方案创新性 | 55 | 100 |
| 可行性 | 60 | 100 |
| 范围定义 | 60 | 80 |
| 风险评估 | 60 | 90 |
| 成功标准 | 55 | 80 |
| 逻辑一致性 | 65 | 90 |
| **总计** | **662** | **1000** |

---

## Phase 3 — 盲点猎杀

### [blindspot-1] Windows 平台兼容性被完全忽略

提案中的编排模式依赖 Unix shell 语义——`just dev` 后台运行（`&`）、进程信号（`SIGTERM`/`SIGKILL`）、`curl` 命令。当前项目在 Windows 11 上开发（环境信息明确标注 `Platform: win32`），但提案完全未讨论跨平台兼容性。这是一个**阻塞性盲点**。

### [blindspot-2] `just dev` 后台进程的生命周期管理

提案将"启动 dev server"委托给 `just dev`，但未讨论：(1) 后台进程的 PID 如何追踪；(2) 如果 `just dev` 在 `just test` 之前崩溃，如何检测和恢复；(3) `just test-teardown` 如何知道要杀哪个进程。run-tests 现有的 `.forge/test-state.json` 只存储 teardown 命令字符串，不存储 PID。这个设计空白会导致测试后大量僵尸进程。

### [blindspot-3] LLM agent 执行的不确定性

提案假设 run-tests（一个 SKILL）会"根据 surface 编排模式决定执行序列"——但 SKILL 是由 LLM agent 执行的，不是由确定性代码执行。LLM 可能：(1) 错误识别 surface 类型；(2) 跳过某个编排步骤；(3) 在错误的时间点执行 teardown。提案未讨论如何将编排序列的**确定性**保证到足以可靠执行——应该考虑是否需要将编排逻辑从 SKILL（LLM 执行）迁移到 Go 代码（确定性执行）。

### [blindspot-4] config-schema.md 的权威性矛盾

提案声称要"移除" `test.execution`，但 config-schema.md 是 run-tests SKILL 的 `references/` 文件——它指导 LLM agent 如何生成和读取配置。移除这个文档意味着所有依赖它的 agent 行为都需要更新。提案将此列为"文档更新"但未评估影响面：还有哪些 agent/SKILL 引用了 config-schema.md？

### [blindspot-5] 回滚计划缺失

如果 surface 感知方案上线后发现严重问题（如 web surface 的编排模式在某种框架下不工作），如何回滚？test.execution 已经被移除，config.yaml 结构已变更。提案没有定义回滚路径或 feature flag 机制。

### [blindspot-6] init-justfile 的幂等性未讨论

如果用户多次运行 init-justfile（surface 配置变化后），justfile 是被完全覆盖还是增量更新？现有的 boundary marker 机制可以处理，但 surface 规则的变化可能导致配方冲突（旧的 web 配方 vs 新的 api 配方）。

---

## 评级

**662/1000 — 及格但存在结构性缺陷**

主要问题：
1. 证据链薄弱（config-schema.md 中的非 just 示例直接反驳了"所有示例都指向 just 命令"的核心论据）
2. 行业对标完全缺失（没有引用任何成熟的项目编排方案）
3. 后台进程生命周期管理是核心技术难点但被回避
4. Windows 平台兼容性被忽略
5. 缺少回滚计划
