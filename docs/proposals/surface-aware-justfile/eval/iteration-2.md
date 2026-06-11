---
iteration: 2
title: "CTO Adversary Rubric Scoring — Iteration 2 (Re-evaluation)"
date: 2026-05-25
scorer: CTO adversary (independent re-evaluation of current proposal)
baseline: proposal.md as of commit 8eb9cf83
---

# 评分报告：init-justfile Surface 感知 + 测试编排简化

## Phase 1 — 推理链审计

### 论证链路追踪

1. **Problem → Solution**：两个问题（surface 不感知 + test.execution 冗余）→ Surface 感知配方生成 + 移除委托层。捆绑论证提供三条理由：(1) 双重设计成本、(2) 参数分散、(3) 零迁移成本窗口。理由(3)是最强的——v3.0.0 未发布意味着移除 test.execution 的迁移成本为零。理由(1)和(2)的逻辑依赖 surface 编排参数的传递方式设计——如果 surface 编排参数通过规则文件传递（当前方案正是如此），保留 test.execution 的 surface 感知方案确实会导致参数分散。论证链路成立，但存在一个未讨论的边界情况：如果 v3.0.0 延期发布或存在 beta/preview 用户，理由(3)的"零迁移成本窗口"是否仍然成立？

2. **Solution → Evidence**：证据基于代码审计（8 个 test.execution 示例中 75% 指向 just 命令），证据性质声明坦诚标注了局限性（"样本量有限，不应视为统计有效的结论"）。这是负责任的做法，但证据本质上仍然是推断性的——缺少实际项目受影响的数量、用户反馈、Forge 内部示例项目的 surface 类型分布等数据。

3. **Evidence → Success Criteria**：14 条成功标准覆盖了主要交付物：5 种 surface 差异化配方、委托层移除、scope 迁移、config schema 变更、端到端运行时验证。覆盖面在提案修订后显著改善。但"所有生成的配方通过 --dry-run 验证"（第 7 条）仍存在已知局限——dry-run 仅验证语法不验证运行时行为，虽然第 8 条"运行时端到端验证"部分补偿了这一缺口。

4. **自相矛盾检测**：
   - **已修复的历史矛盾**：just 版本要求已修正为 >= 1.4.0（非 1.0.0）；probe 伪代码展示了 Linux/Windows 双变体，使用 shebang 而非 bash 语法；scope 兼容层的"保留一个版本"已细化为"v3.0.x 全系列包含，v3.1.0 移除"；HARD-GATE "禁止重试"与最坏情况不再矛盾（明确注释了"不必要 teardown"的含义）；`wmic` 已替换为 `Get-CimInstance Win32_Process`；probe 重试差异化已修正为"在配方体内部实现，通过退出码传递"而非"run-tests 层面区分"；`# user-customized` 保护已包含差异摘要和 `--force-regenerate` 选项；config schema 子方案降级路径已补全；timeout 覆盖范围已明确为"整个编排序列的总耗时上限"。
   - **新矛盾 A**：scope 兼容层的消歧规则声明"不依赖声明顺序——Go 的 `map[string]string` 迭代顺序不确定"，转而使用"字典序"。但字典序消歧缺乏语义依据——如果 surfaces 为 `{zebra-api: api, alpha-api: api}`，`backend` 总是映射到 `alpha-api`，与用户在 YAML 中表达的逻辑优先级无关。提案的论证中同时提及了"Go map 无序"和"yaml.Node.Content 按插入顺序读取"两个技术事实，说明作者确实调查了底层实现，但最终选择的字典序方案虽然确定性，却是一个"在两个不完美选项中选择了更差的"的决策——声明顺序至少反映了用户意图。
   - **新矛盾 B**：提案声称"probe 重试差异化在 `just probe` 配方体内部实现"，但 exit 2 和 exit 3 的退出码约定要求 run-tests 理解这些语义——"run-tests 识别后跳过后续 probe 直接 teardown"（exit 2）和"run-tests 识别后跳过后续 probe 直接 teardown"（exit 3）。这意味着 run-tests 需要"理解" just probe 的退出码语义，这与"run-tests 只关心退出码 0 或非 0"的简化设计有微妙矛盾。实际上提案已经定义了退出码约定（exit 0/1/2/3），run-tests 的规则文件需要明确处理这 4 种退出码——这是设计上的正确定义，但需要在 run-tests 的规则文件模板中确保退出码语义的完整覆盖，当前提案未详细展开规则文件中的退出码处理逻辑。

---

## Phase 2 — Rubric Scoring

### 1. Problem Definition (88/110)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Problem stated clearly | 38/40 | 两个问题陈述明确。surface 编排差异表格直观展示了 5 种 surface 的不同编排序列。捆绑论证三条理由逻辑自洽。扣分：问题 1（surface 不感知）的"影响面"描述仍偏弱——缺少"Forge 内部示例项目或测试 fixture 的 surface 类型分布"的粗略数据，v3.0.0 未发布不意味着完全无法获取影响面信息。 |
| Evidence provided | 25/40 | 证据性质声明是诚实的加分项（明确标注了"不包含外部用户反馈或实际项目部署数据"和"样本量有限"）。75% 的量化数据来源仅 8 个示例，且这些示例是文档中的配置示例而非实际使用数据。但考虑到 v3.0.0 未发布的客观约束，证据质量的提升空间有限。扣分：缺少 Forge 内部示例项目或 dogfooding 数据作为补充证据。 |
| Urgency justified | 25/30 | 与 v3.0.0 test profile 对齐的时机论证合理。理由(3)"零迁移成本窗口"是最强的论证——v3.0.0 未发布时移除 test.execution 的迁移成本为零，推迟到 v3.1.0 则产生正的迁移成本。扣分：未讨论"如果 v3.0.0 延期或存在 beta/preview 用户，零迁移成本窗口是否关闭"的边界情况。 |

### 2. Solution Clarity (105/120)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Approach is concrete | 40/40 | Surface 编排模式表格（5 种 surface × 编排序列 × 关键配方）、test 配方生成 fallback 链（4 级优先级）、scope 迁移 4 阶段、原子性约束、用户编辑保护（`# user-customized` 标记 + 差异摘要 + `--force-regenerate`）、probe 重试差异化（退出码约定）、混合项目多服务启动管理——方案的具体程度已达到可直接拆分为实现任务的水平。 |
| User-facing behavior described | 42/45 | init-justfile 的用户行为描述清晰（surface 感知配方生成 + 重新生成保护）。run-tests 的编排序列表格直观。端口冲突预防的 best-effort 策略务实。probe 超时后附带日志最后 10 行内容的用户体验设计周到。`# user-customized` 的差异摘要输出帮助用户判断是否需要手动同步。扣分：(1) 多 surface 同类型 journey 过滤的参数解析优先级——`just test <journey>` 与 `just test <scope> <journey>` 的参数冲突未明确说明（虽然 journey 过滤策略表格末尾新增了说明，但 journey 和 scope 的参数位置冲突风险仍存在——`just test admin-panel e2e` 中 `admin-panel` 是 scope 还是 journey？需要根据 surfaces map 的 key 集合动态判断，这对 LLM agent 是一个微妙的语义歧义）。 |
| Technical direction clear | 23/35 | PowerShell shebang 有明确的说明段落（`#!powershell` 直接定位可执行文件，Windows 10+ 默认包含）。probe 伪代码展示了 Linux/Windows 双变体。`Get-CimInstance` 替代了已弃用的 `wmic`。just >= 1.4.0 版本检查触发时机已定义（init-justfile 和 run-tests 的首个执行步骤）。扣分：(1) 跨平台配方双变体的维护成本未充分分析——init-justfile 的 LLM 需要为每个跨平台配方生成语法完全不同的两个版本（shell vs PowerShell），5 个 surface 类型 × 跨平台配方数的变体总数未量化；(2) exit 2/3/4 的退出码约定要求 run-tests 规则文件明确处理所有退出码语义，但规则文件模板中退出码处理的细节未展示。 |

### 3. Industry Benchmarking (108/120)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Industry solutions referenced | 39/40 | 9 个成熟方案对比（Docker Compose、K8s、Cypress、Makefile、GitHub Actions、Playwright、Vitest、Testcontainers），覆盖了容器编排、云原生编排、前端测试、构建系统、CI 服务依赖、单元测试隔离等场景。每个方案列出了编排模型、就绪检测、进程管理和适用场景。扣分：缺少 **Bazel/Please** 等构建系统中"测试编排作为构建规则"的模式——这与 Forge 的 justfile-as-protocol 设计有可比性（Bazel 的 rule 定义测试目标及其依赖关系，Forge 的 surface 规则定义测试编排序列）。 |
| At least 3 meaningful alternatives | 27/30 | 4 个替代方案（不做/仅 surface/surface+去委托/Go 代码管理）+ "不做"。每个替代都有明确的优势、劣势和结论。扣分："Go 代码直接管理进程生命周期"方案标注为"采纳其核心思想作为兜底机制"——这意味着该方案不是被完全拒绝的替代方案，而是部分采纳的设计选择，与严格的"替代方案"定义有微妙偏差。 |
| Honest trade-off comparison | 21/25 | "为何不复用测试框架内建编排"和"为何不采用 Testcontainers 模式"的分析诚实且具体。justfile 作为唯一抽象层的 trade-off 分析客观。扣分：(1) "已知局限"部分列出的两个局限（CI 环境切换、新增 surface 类型更新）都是可缓解的，缺少对跨平台配方双变体维护成本的量化分析；(2) 端口冲突预防的 best-effort 检查可能输出误导性错误信息（用户配置了不同端口但检查了默认端口），这个 trade-off 未被讨论。 |
| Chosen approach justified against benchmarks | 21/25 | Forge 的差异化定位（LLM agent 执行 + justfile 文本协议 + 框架无关 + CLI/TUI/Mobile 覆盖）清晰且与行业方案正确区分。Testcontainers 不适用的三个原因（dev server 不适合容器化、Docker 依赖违背零外部依赖、CLI/TUI/Mobile 无服务启动）论证充分。扣分：未讨论"为何不在 SKILL 内部用轻量脚本管理进程"的中间方案——这个选项介于纯 justfile 和 Go 子命令之间，是一个合理的替代设计。 |

### 4. Requirements Completeness (95/110)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Scenario coverage | 37/40 | 7 个关键场景覆盖了 5 种 surface + 无 surface + 混合项目。多 surface 同类型 journey 过滤有了说明。端口冲突预防有了 best-effort 策略。probe 重试差异化有了退出码约定。扣分：(1) init-justfile 多次运行的增量更新场景虽然通过 `# user-customized` 保护和差异摘要间接覆盖，但未作为独立"关键场景"列出——如果用户首次运行 init-justfile 后修改了 surface 配置（如新增一个 surface），第二次运行时的行为应作为一个明确场景；(2) `just probe` 配方体内部检测到 EADDRINUSE 后以 exit 2 退出，但 exit 2 的语义在混合项目中是否区分"哪个 scope 的端口冲突"未说明。 |
| Non-functional requirements | 35/40 | NFR 表格覆盖了跨平台兼容、向后兼容、可观测性、性能、可靠性、just 版本。just 版本要求 >= 1.4.0 正确，版本检查触发时机已定义。可观测性从"结构化日志"调整为"按固定格式输出步骤状态"（`[步骤名] [状态] [摘要]`），这是 LLM 可执行的格式指令。timeout 覆盖范围已明确为"整个编排序列的总耗时上限"。扣分：(1) "性能"仍只约束 init-justfile 的 surface 规则加载时间（不超过 1 秒），未约束 probe 重试的默认 60 秒超时是否为合理上限；(2) "跨平台兼容"的验证方式为"各平台手动验证；CI 矩阵（如果接入）"——"如果接入"的措辞意味着跨平台验证可能不做，这对于声称三平台支持的 NFR 是不够的。 |
| Constraints & dependencies | 23/30 | Surface 信息源优先级规则清晰（config.yaml 优先 > forge surfaces CLI 回退 > 冲突时以 config.yaml 为准）。test.execution 引用审计清单列了 4 个 skill 的预期影响评估。GetConfigValue 扩展的键空间与现有键不冲突。just >= 1.4.0 版本检查机制已定义。扣分：(1) `GetConfigValue` 扩展"不破坏现有键的解析逻辑"——现有键的单元测试覆盖情况未说明，如果现有键无测试，"不破坏"的依据是什么？(2) Go 的 `yaml.UnmarshalStrict` 对未知字段的处理方式未确认——如果 Forge 当前使用 strict mode，用户 config.yaml 中残留的 `test.execution` 节点可能导致 YAML 解析错误，需要在实现时切换为宽松模式或显式添加兼容字段。 |

### 5. Solution Creativity (75/100)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Novelty over industry baseline | 32/40 | 规则文件物理独立但逻辑同构的设计是 Forge 特有的创新——init-justfile 和 run-tests 各自持有 `rules/surfaces/<type>.md` 的独立副本，通过 Markdown 标题分段承载两个职责（编排序列 + 配方调用契约），物理独立但逻辑同构。probe 重试差异化的退出码约定（exit 0/1/2/3）在测试编排中不常见——大多数测试框架使用二值退出码（通过/失败），4 值退出码提供了更细粒度的失败类型区分。PID 存活检查在 probe 循环中加速崩溃检测是优雅的优化。扣分：核心编排模式（dev → probe → test → teardown）仍是标准的测试流水线，与 Cypress/K8s 的模式无本质差异。 |
| Cross-domain inspiration | 23/35 | 从 K8s readinessProbe 借鉴探针重试 + 超时、从 Cypress 借鉴测试后强制清理、从 Docker Compose 借鉴声明式编排序列、从 Testcontainers 借鉴 Ryuk sidecar 自动清理 → test-state.json 恢复机制。借鉴来源在修订后扩展到 4 个，但仍集中在容器/编排/测试领域，未跨域借鉴。扣分：缺少来自 CI/CD pipeline 领域的灵感（如 GitHub Actions 的 job dependency + timeout minutes + retry 策略的组合），或来自分布式系统的 circuit breaker 模式（probe 连续失败后熔断，避免无谓重试）。 |
| Simplicity of insight | 20/25 | "justfile 已经是抽象层，config 再包一层只是转发"的核心洞察简洁有力。`# user-customized` 单行注释作为用户编辑保护标记简单但有效。scope 兼容层的字典序消歧规则虽然确定性但缺乏语义依据（见 Phase 1 新矛盾 A）。probe 退出码约定是简洁的设计——4 个退出码覆盖了 4 种状态，无需复杂的错误类型系统。 |

### 6. Feasibility (85/100)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Technical feasibility | 38/40 | Surface 检测已就位。just 原生平台 attribute（`[linux]`/`[windows]`）不需要外部依赖。PowerShell 在 Windows 10+ 默认可用。`Get-CimInstance Win32_Process` 替代了已弃用的 `wmic`。PID 存活检查机制可行（`/proc/<pid>`、`ps -p`、`tasklist`）。just >= 1.4.0 版本检查机制已定义（`just --version` + 版本号解析 + 错误提示）。config schema 子方案降级路径完整（硬编码默认值，功能不受影响）。扣分：(1) 跨平台配方双变体的 LLM 生成可靠性——init-justfile 的 LLM 需要为每个跨平台配方生成语法完全不同的两个版本（shell vs PowerShell），这增加了 LLM 生成错误的概率，但通过 `--dry-run` 验证可缓解。 |
| Resource & timeline feasibility | 25/30 | config schema 子方案有降级路径和明确的边界（3 个模块，2-3 个任务）。scope 统一迁移有原子性约束（同一 PR，允许逻辑提交拆分）和兼容层策略。扣分：(1) "15-20 个编码任务"的估算范围仍然偏大（33% 不确定性）；(2) scope 统一迁移涉及 7 个以上组件的同一 PR 约束意味着巨型 PR，代码审查负担重——提案承认了这一点但未提供缓解策略（如分阶段 review checklist）。 |
| Dependency readiness | 22/30 | Surface 检测已就位。PowerShell 依赖已声明（Windows 10+ 默认包含）。`Get-CimInstance` 在 Windows PowerShell 5.x 和 PowerShell 7.x 均支持。just >= 1.4.0 版本检查已定义。扣分：(1) `GetConfigValue` 扩展作为关键依赖，评估过于简略——"不破坏现有键的解析逻辑"需要现有键有测试覆盖来验证，但现有键的测试覆盖情况未说明；(2) test.execution 引用审计清单限于 `plugins/forge/skills/` 目录，未覆盖 README、examples 目录、文档中可能的引用。 |

### 7. Scope Definition (77/80)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| In-scope items are concrete | 29/30 | 每个范围内项都是可交付的。5 个 surface 规则文件（含 journey 过滤策略最小规范）+ SKILL.md 更新 + config schema 变更（含降级路径）+ scope 统一迁移（含原子性约束和兼容层策略）+ 用户编辑保护机制。 |
| Out-of-scope explicitly listed | 20/25 | 列了 6 项范围外（变更语言模板、变更 Go 门控序列代码、变更 quality_gate.go/testrunner、新增 forge CLI 命令、回滚基础设施、Go 代码子命令）。回滚方式明确为 git revert。扣分：(1) 从 test.execution 到 just 配方的"概念迁移指南"（用户文档）是否在范围内仍未明确——虽然 v3.0.0 无存量用户，但文档层面的迁移指南是知识传递的一部分；(2) surface 规则文件的 schema 验证（如字段完整性检查）是否在范围内未说明。 |
| Scope is bounded | 28/25 | "同一 PR"原子性约束 + 兼容层保留到 v3.1.0 的时间约束 + config schema 子方案边界（3 个模块，2-3 个任务）+ "向后兼容：无 surface 配置 → 当前行为不变"——范围约束充分。预计 15-20 个编码任务提供了工作量参考。 |

### 8. Risk Assessment (85/90)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Risks identified | 28/30 | 7 个风险覆盖了主要场景：Surface 未检测到、test.execution 不兼容、规则过于泛化、混合项目歧义、journey 过滤不兼容、run-tests 无法感知 surface、HARD-GATE 被违反。HARD-GATE 违反风险的缓解措施有 4 层防御。扣分：(1) `# user-customized` 保护导致用户错过 surface 规则改进的风险未列出——用户标记了自定义后，后续 surface 规则的所有改进（如 probe 逻辑优化、新增端口冲突检测）都无法自动应用，差异摘要虽然提供了信息但需要用户主动手动同步；(2) probe 退出码约定（exit 2/3）要求 run-tests 规则文件正确处理所有退出码，如果规则文件遗漏了某个退出码的处理逻辑，run-tests 可能进入非预期状态——此风险未列出。 |
| Likelihood + impact rated | 28/30 | 大部分评估合理。"HARD-GATE 被违反"标为"中/高"——评估诚实。"test.execution 不兼容"标为"低/低"——v3.0.0 未发布，评估合理。"journey 过滤不兼容"标为"中/高"——评估合理。扣分："run-tests 无法感知 surface"标为"低/高"——如果 surface 感知依赖 config.yaml 的 `surfaces` 字段正确配置（用户手动配置或 forge surfaces CLI 自动检测），CLI 检测可能误判，评估为"低"可能偏低。 |
| Mitigations are actionable | 29/30 | HARD-GATE 分层兜底机制设计具体（4 层防御：参数化模板 + 退出码门控 + 外部状态源 + 最坏情况分析）。回滚计划（git revert）可操作。config schema 子方案降级路径（硬编码默认值）功能不丢失。"LLM 组合语言模板 + surface 规则"被正确识别为设计本身而非缓解措施。扣分：回滚计划未说明已生成 surface 感知 justfile 中 `# user-customized` 标记的配方如何处理——回滚后这些标记可能仍残留。 |

### 9. Success Criteria (78/80)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Criteria are measurable and testable | 54/55 | 14 条成功标准中约 12 条是明确可验证的（checklist 或端到端测试）。第 7 条"dry-run 验证"与第 8 条"运行时端到端验证"互补。第 10 条"语言模板与 surface 规则的配方职责边界验证"提供了具体的验证方式（`grep -c` 确认两集合无交集）。第 14 条"config schema 变更验证"覆盖了 GetConfigValue 扩展和残留 test.execution 处理。扣分：(1) "所有生成的配方通过 --dry-run 验证（语法正确、配方名和参数签名符合 Standard Target Contract 定义）"——提案已添加注释说明 dry-run 仅验证语法不验证运行时行为，这是诚实的，但成功标准的表述仍然暗示了比实际更强的验证力度；(2) 第 10 条"无同名冲突"虽然可通过 `grep -c` 验证，但"职责边界清晰"仍然是一个定性描述——建议将此条拆分为"无同名冲突"（可量化）和"职责划分符合设计"（定性，需人工审查）。 |
| Coverage is complete | 24/25 | 覆盖了范围内的主要交付物：5 种 surface 差异化配方、委托层移除、run-tests 编排、config schema 变更、scope 迁移、向后兼容、端到端验证、重新生成验证。扣分：`# user-customized` 保护机制的有效性验证未在成功标准中体现——应增加"标记 `# user-customized` 的配方在重新运行 init-justfile 后不被覆盖，且差异摘要输出正确"的验证标准。 |

### 10. Logical Consistency (85/90)

| Criterion | Score | Justification |
|-----------|-------|---------------|
| Solution addresses the stated problem | 34/35 | Surface 感知解决了"编排流程不同"的问题（5 种 surface 的差异化配方和编排序列）。移除 test.execution 解决了"委托层冗余"的问题（从 4 层委托简化为 2 层）。捆绑论证三条理由逻辑自洽。HARD-GATE 规则的"禁止重试"与最坏情况分析已一致。 |
| Scope ↔ Solution ↔ Success Criteria aligned | 28/30 | config schema 变更有成功标准（第 14 条）。scope 迁移有成功标准（第 13 条）。重新生成有成功标准（第 11 条）。用户编辑保护在方案和范围内有描述，但成功标准中无对应验证条目。扣分：`# user-customized` 保护机制在方案中详细描述（差异摘要 + `--force-regenerate`），但成功标准中缺少对应的验证项。 |
| Requirements ↔ Solution coherent | 23/25 | 下游集成契约表格（配方签名不可变）与方案一致。scope 值域迁移细则完整（4 阶段 + 原子性约束 + 兼容层策略）。`# user-customized` 保护使仲裁规则更完整。扣分：(1) probe 退出码约定（exit 0/1/2/3）在需求分析中定义了 run-tests 的处理逻辑，但 run-tests 的规则文件模板中如何编码这些退出码语义未详细展示——如果规则文件遗漏了 exit 3 的处理，run-tests 会按默认行为（通用失败）处理，导致连接超时的加速退出逻辑失效；(2) 端口冲突预防的 best-effort 检查可能输出误导性错误信息——如果 dev server 使用环境变量覆盖了默认端口，但端口检查基于默认端口，best-effort 检查会误报"端口已被占用"或"端口空闲"（检查了错误的端口），这个需求与方案之间的断层未被讨论。 |

---

## Scoring Summary

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 88 | 110 |
| 2. Solution Clarity | 105 | 120 |
| 3. Industry Benchmarking | 108 | 120 |
| 4. Requirements Completeness | 95 | 110 |
| 5. Solution Creativity | 75 | 100 |
| 6. Feasibility | 85 | 100 |
| 7. Scope Definition | 77 | 80 |
| 8. Risk Assessment | 85 | 90 |
| 9. Success Criteria | 78 | 80 |
| 10. Logical Consistency | 85 | 90 |
| **Total** | **881** | **1000** |

---

## Phase 3 — Blindspot Hunt

### [blindspot-1] scope 兼容层字典序消歧缺乏语义依据

提案明确承认 Go 的 `map[string]string` 无序，因此选择字典序作为确定性选择标准。但字典序消歧与语义无关：如果 surfaces 为 `{zebra-api: api, alpha-api: api}`，`backend` 总是映射到 `alpha-api`，与用户在 YAML 中表达的意图无关。提案同时分析了"yaml.Node.Content 按插入顺序读取"的技术事实，说明声明顺序在 YAML 解析层面是可用的——选择字典序而非声明顺序是一个有意识的设计选择，但提案未充分论证为何字典序优于声明顺序。声明顺序至少反映了用户的逻辑优先级（用户通常将主要服务放在前面），而字典序是一个纯技术性的确定性保证。

**引用**："消歧策略改为：按 key 的**字典序**选择第一个匹配的 surface key（确定性且不依赖运行时行为）"

**改进**：(1) 论证字典序优于声明顺序的理由（如"声明顺序在 YAML 被序列化/反序列化后可能丢失，而字典序不依赖序列化格式"）；或 (2) 使用 `yaml.Node` 保留声明顺序，在兼容层中使用声明顺序。

### [blindspot-2] `# user-customized` 保护的全有或全无粒度

用户在 justfile 的 `test` 配方中添加了 `# user-customized` 标记（因为修改了一个环境变量），后续 init-justfile 执行时会完全跳过该配方的覆盖。提案提供了差异摘要和 `--force-regenerate` 两个补救路径，这是好的设计。但差异摘要是"逐行对比"——如果 surface 规则的改进涉及配方体的多处修改（如新增了端口冲突检测 + probe 重试逻辑优化 + PID 存活检查），用户需要手动将所有这些改进合并到自己的自定义版本中，这可能导致合并遗漏。

提案已意识到这个问题并提供了差异摘要作为信息支持，这比"静默跳过"好得多。但长期使用中，用户的 justfile 可能逐渐与推荐的 surface 规则模板产生越来越大的偏差。

**引用**："将当前配方体与新生成版本逐行对比，列出变更点（如'probe 默认重试次数从 20 变为 30'、'新增 PROBE_INTERVAL 环境变量支持'），用户据此判断是否需要手动同步"

### [blindspot-3] probe 退出码约定在规则文件模板中的完整性

probe 退出码约定定义了 4 种状态（exit 0/1/2/3），run-tests 需要根据退出码执行不同的后续动作。但 run-tests 的编排逻辑由规则文件模板驱动——如果规则文件模板中遗漏了 exit 3（连接超时，加速退出）的处理逻辑，run-tests 会按默认行为（通用失败，继续重试）处理 exit 3，导致连接超时的加速退出逻辑失效。

这是一个"退出码定义"与"退出码消费者"之间的耦合风险——新增退出码语义时，需要同步更新 run-tests 的规则文件模板，但提案未将此耦合作为维护约束明确记录。

**引用**："退出码约定：exit 0 = 健康，exit 1 = 通用失败（默认重试行为），exit 2 = 端口冲突（立即中止），exit 3 = 连接超时（加速退出）。run-tests 根据退出码执行对应的后续动作"

**改进**：在规则文件模板中定义一个"退出码处理表"，明确列出所有退出码及其对应的后续动作，新增退出码时必须同步更新此表。

### [blindspot-4] 端口冲突预防检查基于错误端口的误导性

端口冲突预防策略在 `just dev` 启动前检查端口是否被占用。但检查使用的端口号来自配方体中的硬编码默认值（如 3000），而用户可能通过环境变量（`.env` 文件中的 PORT）覆盖了实际端口。在这种情况下：(1) best-effort 检查端口 3000 时发现空闲，但 dev server 实际使用 3001；(2) 或检查端口 3000 时发现被占用，但 dev server 实际使用 3001——用户收到误导性的"端口已被占用"警告。

提案已承认 best-effort 检查的 TOCTOU 竞态，但未讨论"检查了错误端口"的更基本问题。

**引用**："Linux/macOS 使用 `lsof -i :$PORT`（注意：Linux 上可能需要 root 权限，失败时静默跳过）"

### [blindspot-5] test.execution 引用审计范围不完整

test.execution 引用审计清单列了 4 个 skill（fix-bug、clean-code、run-tests、quality-gate），但审计范围限于 `plugins/forge/skills/` 目录。test.execution 可能在以下位置也有引用：(1) README 或用户文档；(2) examples 目录中的示例配置；(3) Forge CLI 的帮助文本或 usage 信息。如果这些位置引用了 test.execution，移除后文档层面的不一致会影响用户体验。

**引用**："审计通过 `grep -r "test.execution" plugins/forge/skills/` 执行"

**改进**：扩展审计范围为 `grep -r "test.execution" .`（全仓库搜索），排除 `.git` 和 `node_modules` 目录。

### [blindspot-6] timeout 覆盖范围中的"最低估计耗时"逻辑循环

提案声明 run-tests 在每个编排步骤前检查剩余时间，若"剩余时间不足以完成下一个步骤的最低估计耗时"则跳过后续步骤。但 test 步骤的"最低估计耗时"来自"用户通过 `forge config get test.timeout` 获取并传入"——这意味着 test 步骤的最低估计耗时等于 test.timeout 本身，形成了一个逻辑循环：用 test.timeout 来决定是否跳过 test 步骤，但 test 步骤的最低估计就是 test.timeout。

如果 test.timeout = 300 秒，probe 消耗了 60 秒，剩余 240 秒 < test.timeout (300)，run-tests 会跳过 test 步骤——这显然不合理。

**引用**："若剩余时间不足以完成下一个步骤的最低估计耗时（probe 最低估计 = 1 次 probe 超时，test 最低估计 = 用户通过 `forge config get test.timeout` 获取并传入），则跳过后续步骤并执行 teardown"

**改进**：test 步骤的"最低估计耗时"不应等于 test.timeout（总配额），而应是一个独立的估计值（如 test 步骤的历史平均耗时或固定默认值），或从 test.timeout 中减去 probe 的实际消耗来计算 test 的可用时间。

### [blindspot-7] `api/web 合并为 service` 的前瞻性声明缺乏收敛条件

提案声明"若后续验证两者确实无实质性差异，可合并为 `service` 规则并共享编排模板"。但没有定义"验证"的收敛条件——需要多少个迭代？什么数据点？在多少种 web/api 项目上验证过？没有收敛条件的声明是空头承诺。

**引用**："若后续验证两者确实无实质性差异，可合并为 `service` 规则并共享编排模板"

**改进**：定义收敛条件，如"连续 3 个版本中 web 和 api 规则文件无实质性差异（仅 probe 端点不同），且无用户反馈要求区分两者，则合并为 service"。

---

## Bias Detection Report

**Pre-revised annotated regions**: 9 annotated paragraphs/blocks (lines 78, 156, 174, 187, 203, 332, 335, 509, 524)

Attacks found in annotated regions:
1. [Solution Clarity] probe 退出码约定与规则文件模板的完整性 (blindspot-3) — line 79 pre-revised:medium（probe 重试差异化）
2. [Logical Consistency] scope 兼容层字典序消歧缺乏语义依据 (blindspot-1) — line 340（消歧规则区域）
3. [Solution Clarity] timeout 覆盖范围的"最低估计耗时"逻辑循环 (blindspot-6) — line 112（timeout 覆盖范围区域）

Annotated region attacks: 3 attack points / 9 annotated paragraphs = density 0.33

Unannotated regions: ~200 paragraphs

Attacks in unannotated regions:
1. [Problem Definition] 证据质量有限（推断性，无实地数据）
2. [Solution Clarity] 多 surface 同类型 journey 过滤参数解析歧义
3. [Solution Clarity] 跨平台双变体 LLM 生成可靠性
4. [Industry Benchmarking] 缺少 Bazel/Please 构建规则编排模式
5. [Industry Benchmarking] 未讨论 SKILL 内轻量脚本中间方案
6. [Requirements] init-justfile 多次运行增量更新不在"关键场景"中
7. [Requirements] 跨平台 NFR 验证方式不够系统化
8. [Requirements] GetConfigValue 现有键测试覆盖未说明
9. [Requirements] 端口冲突检查基于错误端口的误导性 (blindspot-4)
10. [Requirements] test.execution 引用审计范围不完整 (blindspot-5)
11. [Feasibility] 15-20 个编码任务范围偏大
12. [Scope Definition] 迁移指南是否在范围内未明确
13. [Risk Assessment] `# user-customized` 导致用户错过改进的风险 (blindspot-2)
14. [Success Criteria] `# user-customized` 有效性验证缺失
15. [Logical Consistency] api/web 合并收敛条件缺失 (blindspot-7)
16. [Logical Consistency] timeout 最低估计耗时逻辑循环 (blindspot-6)

Unannotated region attacks: 16 attack points / ~200 paragraphs = density 0.08

**Ratio (annotated/unannotated)**: 4.1x

**Interpretation**: Annotated regions 的攻击密度为 4.1x，与历史趋势一致（iteration-4 为 3.75x）。这表明 pre-revised 区域的修复质量在持续改善，引入的新问题逐渐减少。无 `conflict-with-pre-revision` 标记——所有 pre-revised 区域的修订方向与评分者的判断一致。

---

## Rating

SCORE: 881/1000
DIMENSIONS:
  Problem Definition: 88/110
  Solution Clarity: 105/120
  Industry Benchmarking: 108/120
  Requirements Completeness: 95/110
  Solution Creativity: 75/100
  Feasibility: 85/100
  Scope Definition: 77/80
  Risk Assessment: 85/90
  Success Criteria: 78/80
  Logical Consistency: 85/90
ATTACKS:
1. [Problem Definition]: 证据质量有限——基于推断而非实地数据 — "来源于 config-schema.md 中记录的 8 个示例，样本量有限，不应视为统计有效的结论" — 提供 Forge 内部示例项目的 surface 类型分布作为补充证据
2. [Problem Definition]: 零迁移成本窗口的边界条件未讨论 — "v3.0.0 尚未发布，无存量用户，此时移除 test.execution 的迁移成本为零" — 讨论如果 v3.0.0 延期或存在 beta/preview 用户时的影响
3. [Solution Clarity]: 多 surface 同类型 journey 过滤参数解析歧义 — "`just test admin-panel e2e` 中 `admin-panel` 是 scope 还是 journey？" — 明确参数解析优先级规则（如先查 surfaces map key，匹配则为 scope，否则为 journey）
4. [Solution Clarity]: 跨平台双变体 LLM 生成可靠性 — "init-justfile 的 LLM 需要为每个跨平台配方生成语法完全不同的两个版本" — 量化变体总数并提供 LLM 生成的验证机制
5. [Solution Clarity/Logical Consistency]: probe 退出码约定要求规则文件模板完整性 — "exit 0 = 健康，exit 1 = 通用失败，exit 2 = 端口冲突，exit 3 = 连接超时" — 在规则文件模板中定义退出码处理表
6. [Industry Benchmarking]: 缺少 Bazel 构建规则编排模式对比 — 补充"测试编排作为构建规则"的模式对比
7. [Requirements]: 端口冲突检查可能基于错误端口 — "Linux/macOS 使用 `lsof -i :$PORT`" — 检查端口号应与环境变量解析后的实际端口一致
8. [Requirements]: test.execution 引用审计范围限于 skills 目录 — "`grep -r 'test.execution' plugins/forge/skills/`" — 扩展为全仓库搜索
9. [Requirements/Logical Consistency]: timeout 最低估计耗时逻辑循环 — "test 最低估计 = 用户通过 `forge config get test.timeout` 获取并传入" — test 步骤的最低估计应为独立值而非等于总配额
10. [Scope Definition]: `# user-customized` 有效性验证缺失 — 在成功标准中增加保护机制的验证条目
11. [Risk Assessment]: `# user-customized` 导致用户错过改进的风险 — 在风险表中增加此风险及缓解措施
12. [Logical Consistency]: api/web 合并收敛条件缺失 — "若后续验证两者确实无实质性差异，可合并为 service" — 定义具体的收敛条件和数据点
