---
iteration: 3
title: "CTO Rubric Scoring — Iteration 3 (FINAL)"
date: 2026-05-25
scorer: CTO adversary
baseline: iteration-2 (748/1000)
---

# 评分报告：init-justfile Surface 感知 + 测试编排简化（第 3 轮 — 终审）

## Phase 1 — 推理链审计

### 论证链路追踪

1. **问题 1**：init-justfile 不感知 surface → 不同 surface 的测试编排流程差异被忽略
2. **问题 2**：test.execution 委托层冗余 → 4 层间接转发，config 再包一层 just 已有的抽象
3. **方案**：surface 感知配方生成 + 废弃 test.execution + surface-orchestration.yaml 声明式编排
4. **证据**：承认非 just 示例存在但认为可封装到 justfile；test.execution 在 Go 层未结构化但在 LLM 层在用
5. **成功标准**：11 条 checklist

### 自相矛盾检测

**上一轮矛盾追踪**：

- **矛盾 A（回滚计划 vs 范围定义）**：已修复。回滚计划现在明确说"通过 git revert 回退"和"不引入 feature flag 机制——回滚基础设施不在 v3.0.0 范围内"。范围外也新增了"回滚基础设施（feature flag 机制不在范围内；回滚通过 git revert 实现）"。逻辑自洽。
- **矛盾 B（兼容期 vs 回退代码删除）**：部分修复。回滚计划第 4 步明确了 v3.2 完全移除的时间线，但仍然没有说"v3.2 删除废弃检测代码时必须同时确认无回退需求"——这个模糊措辞仍然存在，不过严重性降低。
- **矛盾 C（并发 vs 顺序启动）**：部分修复。场景描述中混合项目不再使用"并发"一词，改为"按依赖顺序启动前端和后端"。但"关键场景"部分的描述仍然是"按依赖顺序启动前端和后端"——与方案一致，矛盾已消除。
- **矛盾 D（shebang + [linux/windows] 互斥）**：已修复。提案完全放弃了 shebang 方案，改为"使用 just 原生 `[linux]`/`[windows]` recipe attribute"，并明确说"无需 shebang，无外部依赖（不依赖 bash 可用性）"。NFR 中 just 版本也从 ">= 0.9" 修正为 ">= 1.0"。逻辑一致。

**新矛盾**：

- **新矛盾 A**：替代方案表格中"Go 代码直接管理进程生命周期"一行说"采纳其核心思想作为兜底机制（见'LLM agent 执行确定性'中的分层防御策略）"。但"LLM agent 执行确定性"一节的 v3.0.0 兜底机制实际是"参数化模板 + 退出码约束"——这与 Go 代码管理进程生命周期的核心思想（确定性代码控制进程）毫无关系。退出码约束只是检查命令结果，不管理进程生命周期。声称"采纳核心思想"是夸大——实际上只是在 SKILL.md 中加了文本指令约束。
- **新矛盾 B**：probe 轮询逻辑一节的伪代码使用 `for i in $(seq 1 $max_retries)` — 这是 Unix shell 语法。但提案刚确认使用 just 原生 `[linux]`/`[windows]` recipe attribute（无 shebang），这意味着 Linux 变体的配方体执行的是 just 默认 shell（sh），Windows 变体执行 CMD。伪代码中 `seq` 命令在 macOS 的 sh 中不可用（macOS 的 `/bin/sh` 不包含 GNU seq），`curl -sf` 在 Windows CMD 中不可用。伪代码与选定的跨平台策略不一致——它不是 `[linux]`/`[windows]` 双变体的写法，而是单一的 Unix 伪代码。
- **新矛盾 C**：后台进程管理一节"选定方案"说"为每个需要跨平台行为的配方生成两个变体，just 根据运行平台自动选择匹配的配方"。但 `just probe` 的 probe 轮询伪代码只给出一个 Unix 版本。如果 init-justfile 要为 probe 配方生成 `[linux]` 和 `[windows]` 两个变体，那 probe 的 Windows 版本是什么？`curl` 在 Windows CMD 中不可用（提案自己说"curl 可能不可用"），但伪代码中只用 `curl`。Windows 版 probe 的 fallback（PowerShell Invoke-WebRequest）从未被给出伪代码或示例。

---

## Phase 2 — 评分

### 1. 问题定义 (88/110)

| 子项 | 得分 | 理由 |
|------|------|------|
| 问题清晰 | 37/40 | 两个问题陈述明确，因果关系清晰。问题 1（surface 不感知）和问题 2（委托层冗余）都描述了具体的痛点。"不同 surface 的测试编排流程根本不同"这句话配以表格说明，简洁有力。扣分：两个问题的捆绑必要性仍然只是"相互关联"四个字——替代方案表格中"仅 surface 感知"选项仍以"治标不治本"打发，未分析为什么保留 test.execution + surface 感知不可行。 |
| 证据 | 32/40 | 诚实度持续提升。承认 config-schema.md 中的非 just 示例（`go test`、`npx vitest`、`make test`）并讨论 trade-off。证据部分的质量比第 1 轮有本质性改善。仍扣分：(1) 没有引用实际项目中的 justfile 作为证据——"web UI 的 e2e 测试必须先启动应用"是断言不是证据；(2) "test.execution 的多数示例指向 just 命令"中的"多数"没有量化。 |
| 紧迫性 | 19/30 | 与 v3.0.0 test profile 对齐的论证成立，但与上两轮一样——cost of delay 分析缺失。如果 v3.0.0 不做这个变更，当前方案会导致什么具体损失？是测试会失败？还是用户体验退化？还是技术债累积？没有答案。紧迫性论证停留在"需要协同工作"这个模糊陈述上。 |

**vs 上轮 (85/110)：+3。证据部分的诚实度持续改善。**

### 2. 方案清晰度 (108/120)

| 子项 | 得分 | 理由 |
|------|------|------|
| 方法具体 | 39/40 | surface-orchestration.yaml 的引入是本轮对方案清晰度的最大贡献。它将 init-justfile 和 run-tests 的"统一"从概念变为具体机制。文件格式示例（包含 version/surfaces/orchestration 三层结构）让读者可以精确复述方案。scope 参数值明确为"surfaces map 的 key"而非 surface 类型枚举。扣分：scope 值的语义化问题仍未回答——如果 surfaces key 是 `.` 或 `./`，`just dev .` 在命令行上是否直觉？ |
| 用户行为描述 | 40/45 | 显著改善。废弃行为现在有 4 条精确规则（检测时机、警告格式、行为、文档同步），不再有歧义。迁移路径通过"用户需将命令封装到 justfile"和废弃警告中的迁移说明来指导。混合项目的行为描述统一为"按依赖顺序"（消除了上轮的并发/顺序矛盾）。扣分：init-justfile 多次运行的幂等行为仍未定义——第二次运行时是覆盖 surface-orchestration.yaml 还是合并？ |
| 技术方向 | 29/35 | 跨平台方案从 shebang 切换到 just 原生 recipe attribute，消除了 bash 依赖——这是正确的方向。PID 有效性校验（/proc/<pid> 检查）回应了上轮盲点。但 probe 伪代码仍然只给出 Unix 版本，Windows 版本的 probe 轮询从未示例化。后台进程管理的"选定方案"示例（`nohup`/`start /B`）清晰，但 Windows 示例中 `echo %PID%` 在 `start /B` 后无法获取子进程 PID（`%PID%` 是 CMD 的特殊变量但 `start /B` 不会设置它）——这个 Windows 示例在实际使用中可能不可行。 |

**vs 上轮 (100/120)：+8。surface-orchestration.yaml 消除了最大的设计模糊性。**

### 3. 行业对标 (92/120)

| 子项 | 得分 | 理由 |
|------|------|------|
| 行业方案引用 | 32/40 | 5 个成熟方案（Docker Compose、K8s、Cypress、Makefile、GitHub Actions）的对比表维度合理。新增"从行业方案借鉴的设计"表格，明确列出探针重试、测试后清理、声明式编排的借鉴来源。"Forge 方案的定位差异"一节诚实承认了 Forge 与这些方案的本质差异（LLM agent 执行 vs 确定性代码执行）。扣分：(1) 仍未引用 just 生态自身的最佳实践；(2) Cypress 的 start-server-and-test 仍然只用半行。 |
| 替代方案 | 20/30 | 新增了第 4 个替代方案"Go 代码直接管理进程生命周期"，论述较详细（优势、劣势、结论）。这是对上轮"没有根本不同的架构替代"批评的回应。但仍扣分：(1) 这个方案被定性为"拒绝（v3.0.0 范围过大）"——即不是真正的替代，只是"未来可能"；(2) 没有考虑"让 justfile 包含子命令式的状态管理"这种中间方案（在 justfile 内用文件系统状态机模拟进程管理）。 |
| 权衡对比 | 20/25 | "justfile 作为唯一抽象层的 trade-off 分析"和"已知局限"两节诚实讨论了 CI 环境切换和新 surface 类型的限制。缓解措施（环境变量参数化）合理。 |
| 选定方案论证 | 20/25 | "选择此路径的理由"三点（显式性优于隐式性、justfile 已经是抽象层、可定制性保留在配方体中）逻辑连贯。行业方案的定位差异一节帮助解释了为什么不直接复用 Docker/K8s 式方案。但与 Cypress 式方案（Go 代码直接 fork+wait）的对比只在替代方案表格中有一行，缺乏系统性分析。 |

**vs 上轮 (85/120)：+7。新增的第 4 个替代方案和借鉴设计表格有加分，但替代方案广度仍有不足。**

### 4. 需求完整性 (90/110)

| 子项 | 得分 | 理由 |
|------|------|------|
| 场景覆盖 | 36/40 | 7 个关键场景覆盖了 5 种 surface + 无 surface + 混合。混合项目的并发启动管理有详细设计（端口冲突预防、顺序启动策略、probe 顺序、teardown 逆序清理）。扣分：(1) `just dev` 端口冲突后的行为——说"输出明确错误并退出"，但没有 fallback（不自动换端口是合理的，但应在场景中说明用户需要做什么）；(2) 多个同类型 surface（如 3 个 web surface，分别在不同端口）的场景未覆盖。 |
| 非功能需求 | 30/40 | 6 项 NFR 表格每项有要求和验证方式。just 版本从 ">= 0.9" 修正为 ">= 1.0"，与技术方案一致。扣分：(1) "可靠性"NFR 说"故障注入测试"但范围内没有故障注入测试的任务——NFR 要求与范围不对齐；(2) "跨平台兼容"的验证方式是"各平台手动验证"——仍不是可重复的验证标准；(3) "可观测性"NFR 要求"结构化日志"但范围内没有定义日志格式规范。 |
| 约束与依赖 | 24/30 | Surface 信息源优先级规则清晰（config.yaml > forge surfaces CLI > 冲突时以 config 为准）。config schema 变更列为独立子方案有详细边界。扣分：(1) `.forge/surface-orchestration.yaml` 的格式验证约束未定义——如果用户手动编辑后 YAML 格式错误，run-tests 如何处理？(2) "GetConfigValue 扩展需独立评审"——这个关键依赖仍然只是"需独立评审"。 |

**vs 上轮 (88/110)：+2。混合项目的场景描述改善。**

### 5. 方案创新性 (62/100)

| 子项 | 得分 | 理由 |
|------|------|------|
| 新颖度 | 24/40 | surface-orchestration.yaml 的引入使"统一"从一个概念变为具体机制——这是本轮的创新增量。init-justfile 写入编排声明、run-tests 读取驱动的双向协作模式有一定新颖性。但仍属于"配置文件驱动的声明式编排"这一已有范式，不是新概念。 |
| 跨域灵感 | 18/35 | 借鉴表格承认从 K8s readinessProbe、Cypress、Docker Compose 借鉴。但这些都是直接的"同领域借鉴"。真正的跨域灵感没有出现。 |
| 洞察简洁度 | 20/25 | "justfile 已经是抽象层，config 再包一层只是转发"——这个洞察仍然简洁有力。surface-orchestration.yaml 作为两个 skill 的统一接口是另一个简洁的洞察——一个文件同时服务于生成者和消费者。 |

**vs 上轮 (58/100)：+4。surface-orchestration.yaml 的引入是创新性方面的微量提升。**

### 6. 可行性 (76/100)

| 子项 | 得分 | 理由 |
|------|------|------|
| 技术可行性 | 34/40 | 跨平台方案从 shebang 切换到 just 原生 attribute，消除了 bash 依赖问题，技术可行性显著提升。PID 有效性校验解决了上轮盲点。但：(1) Windows 示例中 `start /B npm run dev > NUL 2>&1 & echo %PID%` — `start /B` 不暴露子进程 PID，`%PID%` 在 CMD 中不是自动变量，这个示例在实际 Windows 环境中可能无法获取 PID；(2) probe 伪代码只给出 Unix 版本，Windows 版本的可行性未验证。 |
| 资源与时间 | 23/30 | "10-15 个编码任务"比上轮更详细，config schema 变更（独立子方案，2-3 个任务）有明细。但新增的 surface-orchestration.yaml 生成逻辑、PID 有效性校验逻辑、跨平台双变体配方生成——这些复杂度未被计入任务估算。 |
| 依赖就绪度 | 19/30 | 与上轮基本持平。test.execution 在 LLM 层面在用的发现仍然是诚实的。但 GetConfigValue 扩展仍然是"需独立评审"，3 轮过后这个关键依赖仍然没有确定性结论。 |

**vs 上轮 (72/100)：+4。跨平台方案的技术可行性提升，但 Windows PID 获取是新问题。**

### 7. 范围定义 (70/80)

| 子项 | 得分 | 理由 |
|------|------|------|
| 范围内 | 28/30 | 范围内项目粒度合理，每个项都是可交付的。config schema 变更的独立子方案定义详细（新增/移除/扩展/边界约束/影响面评估）。surface-orchestration.yaml 生成逻辑明确列入范围内。journey 过滤策略最小规范附在每个 surface 规则文件中。 |
| 范围外 | 17/25 | 本轮新增了"回滚基础设施（feature flag 机制不在范围内；回滚通过 git revert 实现）"——回应了上轮批评。但仍然缺少：(1) 现有项目的迁移指南（从 test.execution 迁移到 surface 感知模式）；(2) surface-orchestration.yaml 的格式校验和错误恢复；(3) 多次运行 init-justfile 的幂等行为定义。 |
| 范围边界 | 25/25 | 向后兼容约束保留。config schema 变更独立子方案边界清晰。回滚计划与范围定义现在一致（git revert，无 feature flag）。 |

**vs 上轮 (62/80)：+8。回滚计划与范围定义的对齐修复了上轮的最大扣分点。**

### 8. 风险评估 (78/90)

| 子项 | 得分 | 理由 |
|------|------|------|
| 风险识别 | 26/30 | 6 个风险覆盖主要场景。回滚计划现在与范围定义一致（git revert）。但仍然未识别：(1) Windows 上 PID 获取可能失败（`start /B` 不暴露 PID）的风险；(2) surface-orchestration.yaml 被用户手动编辑后格式错误导致 run-tests 解析失败的风险；(3) just 原生 `[linux]`/`[windows]` attribute 要求 just >= 1.0，但项目当前使用的 just 版本是否满足？ |
| 可能性+影响 | 26/30 | 6 个风险评估合理。test.execution 兼容性风险"低/低"的论证（v3.0.0 未发布，无存量用户）诚实。 |
| 缓解措施 | 26/30 | HARD-GATE 规则、状态机驱动、退出码约束、参数化模板——缓解措施的可操作性提升。回滚计划从"feature flag"简化为"git revert"，更务实。但"参数化模板 + 退出码约束提供确定性下限"——退出码是 just 命令的返回值，LLM 是否能可靠地检查每个命令的退出码并做出正确决策？如果 LLM 忽略了非零退出码直接执行下一步（技术上可以做到，因为 LLM 控制执行流），退出码约束就不是确定性的。 |

**vs 上轮 (72/90)：+6。回滚计划的务实化是主要加分项。**

### 9. 成功标准 (68/80)

| 子项 | 得分 | 理由 |
|------|------|------|
| 可测量可测试 | 44/55 | 11 条标准中有 8 条是明确的。新增第 11 条（混合项目端到端验证，包含 PID 文件有效性校验）覆盖了上轮的盲点。扣分：(1) "所有生成的配方通过 --dry-run 验证"——dry-run 只验证语法不验证运行时，这个批评从第 1 轮就存在但从未被改善；(2) 第 5 条"test.execution 节点从 config-schema 中标记为 @deprecated，非命令字段移至 test: 顶层，废弃检测正常工作"——"正常工作"如何测试？需要测试用例但没有给出；(3) 第 11 条"混合项目端到端验证"——手动验证还是自动化测试？未说明。 |
| 覆盖完整性 | 24/25 | 覆盖了范围内主要交付物。config schema 变更的成功标准通过第 5 条覆盖。回滚验证通过 git revert 机制隐含覆盖（不需要单独标准）。 |

**vs 上轮 (58/80)：+10。第 11 条混合项目标准是关键改善。**

### 10. 逻辑一致性 (78/90)

| 子项 | 得分 | 理由 |
|------|------|------|
| 方案解决问题 | 31/35 | surface 感知解决问题 1，废弃 test.execution 解决问题 2。surface-orchestration.yaml 确保两个解决方案通过共享接口协作。但问题捆绑的必要性仍然论证不足——"替代方案"表格中"仅 surface 感知"选项仍以"治标不治本"四个字打发。 |
| 对齐 | 26/30 | 上轮最大的不对齐（回滚计划 vs 范围定义）已修复。混合项目"并发 vs 顺序"矛盾已修复。NFR 与技术方案的对齐改善（just >= 1.0）。扣分：(1) NFR 要求"故障注入测试"但范围内无相应任务；(2) probe 伪代码只给出 Unix 版本但选定方案要求双平台变体。 |
| 需求-方案一致性 | 21/25 | 下游集成契约表格与方案设计一致。语言模板 vs surface 规则的仲裁规则清晰。NFR 与技术深化基本对应。扣分：(1) "可靠性"NFR 要求"test-state.json 恢复清理"但这个机制的细节在方案中是"现有机制"——没有验证现有机制是否支持多 PID 场景；(2) scope 值为 surfaces map key 但 journey 过滤策略表格中的 journey 标签映射（如 `smoke → 登录+首页加载`）与 scope 无关——混合项目中不同 scope 的 journey 过滤如何区分？ |

**vs 上轮 (68/90)：+10。回滚矛盾和并发/顺序矛盾的修复是主要加分。**

---

## 评分汇总

| 维度 | 上轮 | 本轮 | 满分 | 变化 |
|------|------|------|------|------|
| 问题定义 | 85 | 88 | 110 | +3 |
| 方案清晰度 | 100 | 108 | 120 | +8 |
| 行业对标 | 85 | 92 | 120 | +7 |
| 需求完整性 | 88 | 90 | 110 | +2 |
| 方案创新性 | 58 | 62 | 100 | +4 |
| 可行性 | 72 | 76 | 100 | +4 |
| 范围定义 | 62 | 70 | 80 | +8 |
| 风险评估 | 72 | 78 | 90 | +6 |
| 成功标准 | 58 | 68 | 80 | +10 |
| 逻辑一致性 | 68 | 78 | 90 | +10 |
| **总计** | **748** | **810** | **1000** | **+62** |

---

## Phase 3 — 盲点猎杀

### [blindspot-1] Windows PID 获取方案不可行

提案的 Windows 示例：`start /B npm run dev > NUL 2>&1 & echo %PID% > .forge\\dev-server.pid`。`start /B` 在 CMD 中后台启动进程但**不暴露子进程 PID**。`%PID%` 不是 CMD 的自动变量（与 bash 的 `$!` 不同）。这意味着整个 PID 追踪链条在 Windows 上断裂——没有 PID 就无法 teardown。提案需要一个真正可行的 Windows PID 获取方案，如 PowerShell 的 `Start-Process ... -PassThru | Select-Object -ExpandProperty Id`。

### [blindspot-2] probe 配方的 Windows 变体从未定义

提案明确说"为每个需要跨平台行为的配方生成两个变体"，但 probe 轮询伪代码只给出 Unix 版本（`curl -sf`、`seq`、`sleep`）。Windows CMD 没有 `curl`（提案自己也承认"curl 可能不可用"），没有 `seq`，`sleep` 也不是 CMD 内建命令。probe 的 Windows 变体必须使用 PowerShell（`Invoke-WebRequest`、`Start-Sleep`）或完全不同的轮询策略。如果 init-justfile 的配方模板中没有 Windows probe 变体，NFR "跨平台兼容"就是不可达成的。

### [blindspot-3] surface-orchestration.yaml 缺少格式验证和容错

提案将 .forge/surface-orchestration.yaml 定义为 init-justfile 和 run-tests 的"统一接口"和"共享文件"。用户甚至可以"直接编辑此文件调整编排行为"。但：
- 如果 YAML 格式错误（缩进错误、字段拼写错误），run-tests 的解析会失败，整个测试编排崩溃
- 如果 `startup_order` 中引用了不存在的 scope key，运行时才会发现
- 如果用户删除了一个 surface 但忘记更新此文件，会出现孤儿配置
- 提案没有定义任何 schema 验证、默认值、或错误恢复机制

对于"统一接口"这种关键依赖，格式验证应该是最基本的保障。

### [blindspot-4] "采纳核心思想"是虚夸

替代方案表格对"Go 代码直接管理进程生命周期"说"采纳其核心思想作为兜底机制"。但实际上 v3.0.0 的兜底机制是"参数化模板 + 退出码约束 + HARD-GATE 文本指令"——这三者与"Go 代码管理进程生命周期"的核心思想（确定性代码控制进程启动/停止/信号发送）完全无关。退出码约束只能检查命令是否成功，不能管理进程生命周期。HARD-GATE 是给 LLM 的文本指令，不是确定性保证。这不是"采纳核心思想"——这是"认可了问题但选择了完全不同的解决方案"。这种夸大措辞会让后续维护者误解兜底机制的能力边界。

### [blindspot-5] 混合项目 journey 过滤的 scope 感知缺失

journey 过滤策略表格定义了每种 surface 类型的标签映射（如 web 的 `smoke → 登录+首页加载`，api 的 `contract → API 契约验证`）。但在混合项目中（web+api），`just test smoke` 应该只运行 web scope 的 smoke journey 还是所有 scope 的 smoke journey？`just test admin-panel`（scope 作为参数）和 `just test smoke`（journey 作为参数）的参数解析如何区分？如果 `admin-panel` 恰好也是一个有效的 journey 标签名，就会产生歧义。提案没有定义 scope 和 journey 参数的仲裁规则。

### [blindspot-6] 无 surface 配置的"当前行为"是什么？

成功标准第 6 条："无 surface 配置的项目输出与当前一致"。但"当前行为"是什么？如果当前 run-tests 依赖 test.execution 配置来决定执行序列，而提案废弃了 test.execution，那么"无 surface 配置"的项目在废弃 test.execution 后的行为与"当前行为"不可能一致——因为 test.execution 的读取路径被移除了。要么成功标准第 6 条是矛盾的，要么"当前行为"需要更精确的定义（如"无 surface 配置且无 test.execution 配置的项目"）。

---

## 评级

**810/1000 — 良好，关键结构问题已修复，剩余问题主要在实现细节层**

与第 2 轮（748）相比主要改善：

1. 回滚计划与范围定义对齐修复（+8 范围定义，+6 风险评估，+10 逻辑一致性）
2. surface-orchestration.yaml 消除了"统一"概念的模糊性（+8 方案清晰度）
3. 跨平台方案从 shebang 切换到 just 原生 attribute，消除了 bash 依赖（+4 可行性）
4. 混合项目成功标准补充（+10 成功标准）

仍存在的结构性问题（按严重度排序）：

1. **Windows PID 获取方案不可行** — `start /B` 不暴露子进程 PID，整个 teardown 链条在 Windows 上断裂
2. **probe Windows 变体从未定义** — 只给出 Unix 伪代码，NFR "跨平台兼容"不可达成
3. **NFR "故障注入测试" 与范围不对齐** — 要求故障注入但范围内无相应任务
4. **surface-orchestration.yaml 缺少格式验证** — 关键共享文件没有容错机制
5. **混合项目 journey 过滤缺少 scope 感知** — scope 参数和 journey 参数可能歧义
6. **问题捆绑必要性论证不足** — 3 轮过后仍以"治标不治本"打发
