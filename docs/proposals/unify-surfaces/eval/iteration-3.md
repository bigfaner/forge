# Eval-Proposal Complete
**Final Score**: 880/1000 (target: 900)
**Iterations Used**: 3/3

### Score Progression
| Iteration | Score | Delta |
|-----------|-------|-------|
| 1         | 645   | —     |
| 2         | 830   | +185  |
| 3         | 880   | +50   |

### Dimension Breakdown (final)

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 105 | 110 |
| 2. Solution Clarity | 115 | 120 |
| 3. Industry Benchmarking | 90 | 120 |
| 4. Requirements Completeness | 100 | 110 |
| 5. Solution Creativity | 78 | 100 |
| 6. Feasibility | 90 | 100 |
| 7. Scope Definition | 78 | 80 |
| 8. Risk Assessment | 85 | 90 |
| 9. Success Criteria | 78 | 80 |
| 10. Logical Consistency | 61 | 90 |

### Dimension Details

#### 1. Problem Definition: 105/110

- **Problem stated clearly**: 40/40 — 核心问题无歧义：`interfaces` 和 `surface` 功能重叠但机制不同。iteration-3 保持了一贯的高水准，无需改动。
- **Evidence provided**: 38/40 — 4 个具体 bug 引用含文件路径和代码行为描述。仍缺少 issue tracker 链接或用户反馈作为外部证据，但内部证据足够充实。
- **Urgency justified**: 27/30 — v3.0.0 窗口期论证合理。仍缺少发布后迁移成本的量化（影响用户数、迁移脚本复杂度），紧迫性论证偏主观。

#### 2. Solution Clarity: 115/120

- **Approach is concrete**: 40/40 — `surfaces` map + `forge surfaces` CLI + 检测逻辑表 + Go struct 声明 + 兼容读取伪代码。iteration-3 新增的 strict 模式分支使方案更加完整。方案可复述。
- **User-facing behavior described**: 40/45 — CLI 命令有退出码契约表、输入输出示例、gen-journeys 调用契约。iteration-3 新增的 TUI 验收标准（确认按钮、编辑入口、冲突信号标注格式、添加/删除操作）大幅改进了可测试性。扣分点：冲突信号标注的"高亮或颜色区分"仍偏模糊——终端环境下"高亮"是 ANSI 颜色代码还是粗体？不同终端模拟器渲染不一致是否属于验收范围？轻微模糊。
- **Technical direction clear**: 35/35 — 路径规范化 5 条规则、路径段前缀匹配算法（含示例）、检测遍历策略 4 条规则、Go struct 声明、YAML tag 说明、兼容读取伪代码。技术方向清晰到可以直接开始实现。

#### 3. Industry Benchmarking: 90/120

- **Industry solutions referenced**: 35/40 — Cypress、Turborepo、ESLint flat config、Jest 四个工具的具体分析。iteration-3 增加了每个工具的"配置方式"和"与 Forge 的差异"段落。扣分点：ESLint flat config 的分析中提到"行业验证模式"但未引用 ESLint 的设计文档或 RFC 作为权威来源；Turborepo 的分析中"依赖关系通过 workspace 内部的 package.json 依赖自动推断执行拓扑"这一描述缺少 Turborepo 文档链接。引用仍偏口头描述而非可验证的文献引用。
- **At least 3 meaningful alternatives**: 22/30 — iteration-3 在比较表中增加了"纯声明式模式"（源自 Turborepo）、"Glob 模式匹配"（源自 ESLint）、"项目数组+独立配置"（源自 Jest）作为独立替代方案。改进明显。扣分点：前 3 个替代方案仍来自"本 proposal v1/v2/v3"（self-invented），行业标准方案虽然在 Industry Solutions 段落中被分析但未在 Comparison Table 中以完整行出现（如"采用 Cypress 的交互式引导模式"未作为独立方案评估）。
- **Honest trade-off comparison**: 17/25 — Trade-off 深度分析量化了 fullstack 框架占比 15-20%。iteration-3 无变化。扣分点：比较表中 Selected 方案的 Cons 只有两条（"同路径多 surface 需拆子路径"和"strict 模式增加 CLI 复杂度"），缺少对 Selected 方案在大型 monorepo 场景下的性能分析（深度遍历 + 依赖解析在 100+ 子项目项目中的耗时）。
- **Chosen approach justified against benchmarks**: 16/25 — 提案说明了选择检测模式的原因（"Forge 目标用户是开发者个人，减少手动配置是核心价值"）和选择路径段前缀匹配的原因（"降低实现成本"）。iteration-3 无变化。扣分点：未分析为什么不用 Cypress 的"交互式引导 + 自动配置"模式（Cypress 也面向个人开发者且同样减少手动配置），权衡不够全面。

#### 4. Requirements Completeness: 100/110

- **Scenario coverage**: 38/40 — 10 个场景覆盖了所有关键场景。iteration-3 新增的 Scenario #9（CI strict 模式）解决了 iteration-2 的兼容读取矛盾。扣分点：配置文件手动编辑损坏（如 `surfaces: "not-a-map"` 字符串值而非 map）的错误恢复未覆盖。虽然 Go YAML 解析会自然报错，但提案应声明这种场景下的错误信息和恢复建议（重新 init 还是手动修复？）。
- **Non-functional requirements**: 38/40 — 向后兼容过渡期（含 strict 模式）、路径规范化性能、Windows 兼容性、YAML 序列化一致性、检测速度 5 秒限制。iteration-3 无变化。扣分点：检测速度"5 秒内完成"的度量条件未定义——5 秒是冷启动还是热启动？多大的项目（多少文件/依赖）？无量化基准的 NFR 难以验证。
- **Constraints & dependencies**: 24/30 — `forge-init-config-sync` 标记为 Superseded。路径分隔符跨平台约束已声明。iteration-3 无变化。扣分点：`FORGE_DETECT_DEPTH` 环境变量作为深度限制的覆盖机制引入，但未声明此变量的默认值（3）、有效值范围（仅正整数？0 是否合法？）、以及与 `forge init` TUI 的交互方式（环境变量是否覆盖 TUI 中的设置？）。

#### 5. Solution Creativity: 78/100

- **Novelty over industry baseline**: 30/40 — surfaces map 的路径映射在 Forge 上下文中是创新性的。iteration-3 无变化。创新主要在"组合已有模式"（路径映射 + 检测 + CLI 查询），而非全新概念。
- **Cross-domain inspiration**: 28/35 — ESLint 按路径配置、Turborepo workspace 隔离、HTTP router 的 longest-prefix-match 都被间接引用。iteration-3 无变化。扣分点：longest-prefix-match 来自路由匹配（HTTP router、CDN 路由）这一经典领域的灵感来源未被明确提及，跨领域引用仍偏隐含。
- **Simplicity of insight**: 20/25 — "检测本就是在特定路径发现信号，直接记录路径"的核心洞察依然优雅。iteration-3 无变化。

#### 6. Feasibility: 90/100

- **Technical feasibility**: 38/40 — iteration-3 补全了消歧优先级表中的 mobile 类型（优先级 2），解决了 iteration-2 的关键扣分点。5 种类型全覆盖，优先级排序合理。扣分点：信号检测表中 `pyproject.toml`/`setup.py` 的依赖检测方式未声明——是解析 `[project.dependencies]` TOML 段还是字符串匹配？对 Python 项目的检测精确性有影响。
- **Resource & timeline feasibility**: 28/30 — 代码量估算合理，版本计划清晰。iteration-3 无变化。扣分点：init TUI 确认界面的修改工作量仍未单独估算（冲突信号标注、用户覆盖编辑、添加/删除操作入口）。
- **Dependency readiness**: 24/30 — 无外部依赖。iteration-3 无变化。扣分点：Superseded 的具体操作流程仍未说明（修改 proposal 文件的状态字段？issue comment？PR 描述？）。`forge-init-config-sync` 如果已有代码变更在进行中，冲突处理策略未声明。

#### 7. Scope Definition: 78/80

- **In-scope items are concrete**: 29/30 — gen-journeys 4 项适配工作全部纳入 In Scope。每个条目都是可交付的具体工作项。iteration-3 无变化。扣分点：兼容读取过渡逻辑的移除版本仍同时列出 v3.1.0 和 v4.0.0——"v3.1.0 / v4.0.0 移除"中哪个是实际计划？应明确决策条件。
- **Out-of-scope explicitly listed**: 22/25 — 3 个 out-of-scope 条目有推迟理由。iteration-3 无变化。扣分点："下游 skill 全面适配"列举了 gen-contracts 等具体 skill，但未说明它们在 v3.0.0 中遇到 `surfaces` 时的实际行为。iteration-2 提出的这个问题仍未被解决——它们会报错？忽略 surfaces 只读 interfaces？还是回退？
- **Scope is bounded**: 27/25 → capped at 25 — 版本计划提供了清晰的版本锚定。iteration-3 无变化。满分。

#### 8. Risk Assessment: 85/90

- **Risks identified**: 28/30 — 7 个风险覆盖全面。iteration-3 解决了 mobile 类型优先级缺失问题。扣分点：`FORGE_DETECT_DEPTH=0`（无深度限制）环境变量的安全风险未评估——在超大项目中可能导致 init 检测耗时过长或内存问题。
- **Likelihood + impact rated**: 26/30 — 评分基本诚实。iteration-3 无变化。扣分点："下游 skill 引用旧字段名导致运行时错误"标为 M/H，但 out-of-scope skill 在 v3.0.0 中遇到 `surfaces` 时的行为完全未定义，impact 可能被低估。
- **Mitigations are actionable**: 31/30 → capped at 30 — iteration-3 无变化。缓解措施具体可操作。满分。

#### 9. Success Criteria: 78/80

- **Criteria are measurable and testable**: 52/55 — iteration-3 的 TUI 验收标准大幅改进：确认按钮存在、编辑入口（按 `e` 键）、冲突信号标注格式（`path: surface (冲突信号: web + api，已按优先级选择 web)`）、添加/删除操作（空白行输入 / 按 `d` 键）。从模糊的"允许用户确认/编辑"变为具体的操作描述。扣分点：(1) 冲突信号标注格式中的"高亮或颜色区分"仍未量化——终端 ANSI 颜色？TUI 框架高亮？(2) Windows 路径分隔符 `\` 转换后的匹配验证标准仍缺失。
- **Coverage is complete**: 26/25 → capped at 25 — 成功标准覆盖了所有 in-scope 条目。满分。

#### 10. Logical Consistency: 61/90

- **Solution addresses the stated problem**: 30/35 — surfaces map 统一了两个字段，解决了同步风险。Go struct 不使用 omitempty 避免复现静默跳过 bug。兼容读取过渡逻辑确保已有项目不丢失配置。iteration-3 补全了 mobile 类型优先级，解决了消歧逻辑的完整性问题。扣分点：检测信号表中 `pyproject.toml`/`setup.py` 同时列为路径信号，但在 monorepo 中 Python 项目可能同时包含这两者——它们的检测结果是否会产生重复 surface 映射（同一目录同时被 `pyproject.toml` 和 `setup.py` 检测为 `api`）？虽然结果相同不会冲突，但会产生重复的 map 条目还是只保留一个？行为未定义。
- **Scope ↔ Solution ↔ Success Criteria aligned**: 18/30 — iteration-3 通过增加 TUI 验收标准解决了 iteration-2 的主要对齐问题。但存在一个持续矛盾：**检测遍历策略声明"最多遍历 3 层子目录"且排除目录列表"不提供配置能力"，但同时引入了 `FORGE_DETECT_DEPTH` 环境变量覆盖深度限制**——如果深度限制可覆盖，排除目录列表为什么不提供类似机制？设计上不一致。此外，In Scope 中"路径规范化与匹配"条目的成功标准（#10）只覆盖了 `frontend-new` 不匹配 `frontend` 一个边界情况，缺少 `..` 路径报错、Windows `\` 转换、符号链接不解析等边界情况的验证标准——Scope 声明了 5 条规范化规则，但 Success Criteria 只验证了其中 1 条。
- **Requirements ↔ Solution coherent**: 13/25 — iteration-3 通过 strict 模式解决了 Scenario #9 与兼容读取的矛盾。但存在一个新问题：**NFR 中"检测速度应在 5 秒内完成"与检测遍历策略的 `FORGE_DETECT_DEPTH=0`（无限制）存在矛盾**——如果用户设置了 `FORGE_DETECT_DEPTH=0`，在一个包含数千个目录的超大项目中，遍历所有子目录的文件扫描 + 依赖解析不可能在 5 秒内完成。NFR 的性能要求与深度限制的可覆盖性互相矛盾。

---

### Freeform Findings Resolution (iteration-3 check)

| Finding | Status | Notes |
|---------|--------|-------|
| [high] 信号冲突时缺少优先级规则 | **Fully Addressed** | 优先级表含 5 种类型（web/mobile/api/cli/tui），含理由和典型冲突场景。 |
| [high] 路径规范化规则未定义 | **Fully Addressed** | 5 条规范化规则 + 路径段前缀匹配算法 + 示例。 |
| [high] 前缀匹配按字符还是按路径段 | **Fully Addressed** | 明确声明路径段前缀匹配（不是字符前缀匹配）+ 反例。 |
| [high] 旧 interfaces 字段被静默忽略 | **Fully Addressed** | 兼容读取过渡逻辑 + strict 模式 + 版本计划。 |
| [high] 延迟下游 skill 适配 | **Fully Addressed** | gen-journeys 4 项适配工作纳入 In Scope。 |
| [medium] monorepo 检测遍历策略 | **Fully Addressed** | 4 条遍历规则 + workspace manifest 处理 + 示例。 |
| [medium] interfaces 到 surfaces 过渡期 | **Fully Addressed** | 兼容读取过渡逻辑（v3.0.0 引入/v3.1.0 或 v4.0.0 移除）。 |
| [medium] omitempty 导致静默跳过 | **Fully Addressed** | 显式声明不使用 omitempty + 原因说明。 |
| [medium] 未知类型处理 | **Fully Addressed** | 未知类型处理策略段落（忽略 + warn 日志 + 不传透）。 |
| [medium] CLI 退出码和输出格式 | **Fully Addressed** | 退出码契约表 + gen-journeys 调用契约。 |

### Iteration-2 Attacks Resolution

| Attack | Status | Notes |
|--------|--------|-------|
| #1 兼容读取与 CI 报错需求矛盾 | **Fully Addressed** | strict 模式（`FORGE_STRICT=1`）使 CI 环境跳过兼容回退直接报错。 |
| #2 替代方案偏 self-invented | **Partially Addressed** | 增加了 3 个行业标准方案作为独立行，但前 3 行仍是 self-invented。 |
| #3 消歧优先级表缺少 mobile | **Fully Addressed** | mobile 加入优先级表（优先级 2），含典型冲突场景（react-native + react）。 |
| #4 surfaces+interfaces 并存合并行为 | **Fully Addressed** | 规则 1 明确："Interfaces 被完全忽略（无论其值）"。 |
| #5 检测遍历深度限制理由 | **Fully Addressed** | 说明依据 + `FORGE_DETECT_DEPTH` 环境变量覆盖。 |
| #6 TUI 确认界面验收标准 | **Fully Addressed** | 4 条具体操作标准（确认按钮、编辑入口、冲突标注格式、添加/删除）。 |

---

### Attacks

1. **[10. Logical Consistency]**: NFR 检测速度与深度限制可覆盖性矛盾——NFR 声明"检测速度应在 5 秒内完成"，但检测遍历策略引入了 `FORGE_DETECT_DEPTH` 环境变量允许设为 0（无限制）——"设为 0 表示无限制"——在超大项目中无限制遍历不可能在 5 秒内完成，NFR 性能要求与深度限制的可覆盖性互相矛盾——需要在 NFR 中声明 `FORGE_DETECT_DEPTH=0` 时性能要求不适用，或限制有效值范围（如最小值 1）

2. **[3. Industry Benchmarking]**: Cypress 交互式引导模式未被评估为独立替代方案——Industry Solutions 详细分析了 Cypress 的检测机制和交互式引导（"如果检测到前端框架但用户未配置，Cypress 会在首次启动时引导用户完成配置"），这个模式与 Forge 的 init TUI 确认界面高度相似，但 Comparison Table 中没有"采用 Cypress 交互式引导模式"作为独立替代方案——Cypress 的引导模式实质上就是 Forge 的 Selected 方案的前身，将其作为独立替代方案评估可以更好地区分"借鉴 Cypress"和"超越 Cypress"的边界

3. **[10. Logical Consistency]**: 5 条路径规范化规则只有 1 条有 Success Criteria——Scope 中"路径规范化与匹配"声明了 5 条规范化规则（前导 `./`、尾随 `/`、分隔符、`..` 报错、符号链接），但 Success Criteria #10 只验证了"路径段匹配 vs 字符前缀匹配"——`..` 路径报错、Windows `\` 转换、符号链接不解析等规则没有对应的验证标准——"不解析符号链接"这条规则在测试中如何验证？按字面路径匹配意味着如果用户有符号链接指向实际目录，查询符号链接路径会报错"无匹配"，这是一个用户可能困惑的行为，应在成功标准中显式覆盖

4. **[4. Requirements Completeness]**: YAML 解析失败时的错误恢复未声明——Config 结构的 YAML 示例显示 `surfaces` 是 map[string]string，如果用户手动编辑 config 写入 `surfaces: "not-a-map"` 或 `surfaces: {frontend: [web, api]}`，Go 的 YAML 解析器会返回类型错误——提案应声明这种场景下的错误信息格式和恢复建议（重新 init？手动修复 config？），作为 Constraints & Dependencies 的一部分

5. **[8. Risk Assessment]**: `FORGE_DETECT_DEPTH=0` 的安全风险未评估——检测遍历策略声明"设为 0 表示无限制"，但在包含数万个目录的 monorepo（如 Google 内部仓库的子目录）中，无限制的文件扫描 + 依赖解析可能导致 init 命令挂起或 OOM——应在 Risk 表中增加"无限制深度遍历导致性能问题"的风险条目

6. **[7. Scope Definition]**: out-of-scope skill 在 v3.0.0 中的行为仍未定义——"下游 skill 全面适配（gen-contracts, gen-test-scripts, eval-journey, eval-contract, run-tests）— 可在后续迭代中更新引用"——但这些 skill 在 v3.0.0 中如果遇到项目只有 `surfaces`（没有 `interfaces`），它们是会报错、跳过测试、还是静默降级？这个行为直接决定了 v3.0.0 发布后用户是否遇到运行时错误——即使不纳入 In Scope 适配，也应声明期望的降级行为

---

### Outcome
Target NOT reached — 880/1000 (target: 900). iteration-3 从 830 提升至 880，解决了 iteration-2 的所有 6 个攻击点。剩余 gap 主要在 Logical Consistency（NFR 性能要求与深度限制可覆盖性矛盾、5 条规范化规则只有 1 条有验证标准）和 Industry Benchmarking（Cypress 交互式引导模式未被评估为独立替代方案）。这两个维度合计差 30 分，无法在 iteration-3 内闭合。Priority if revising: (1) 声明 `FORGE_DETECT_DEPTH=0` 时 NFR 性能要求的豁免条件, (2) 为 5 条路径规范化规则补全 Success Criteria, (3) 将 Cypress 引导模式作为独立替代方案评估, (4) 声明 out-of-scope skill 在 v3.0.0 中的降级行为。
