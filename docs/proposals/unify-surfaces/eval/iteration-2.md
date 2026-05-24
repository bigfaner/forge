# Eval-Proposal Complete
**Final Score**: 830/1000 (target: 900)
**Iterations Used**: 2/3

### Score Progression
| Iteration | Score | Delta |
|-----------|-------|-------|
| 1         | 645   | —     |
| 2         | 830   | +185  |

### Dimension Breakdown (final)

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 105 | 110 |
| 2. Solution Clarity | 110 | 120 |
| 3. Industry Benchmarking | 85 | 120 |
| 4. Requirements Completeness | 95 | 110 |
| 5. Solution Creativity | 75 | 100 |
| 6. Feasibility | 85 | 100 |
| 7. Scope Definition | 75 | 80 |
| 8. Risk Assessment | 80 | 90 |
| 9. Success Criteria | 70 | 80 |
| 10. Logical Consistency | 50 | 90 |

### Dimension Details

#### 1. Problem Definition: 105/110

- **Problem stated clearly**: 40/40 — 核心问题无歧义：`interfaces` 和 `surface` 功能重叠，取值相同但机制不同，存在同步风险。修订版保持了 iteration-1 的高水准，问题定义清晰完整。
- **Evidence provided**: 38/40 — 4 个具体 bug 引用含文件路径和代码行为描述。修订版无变化，仍缺少 issue tracker 链接或用户反馈作为外部证据，但内部证据足够充实。轻微扣分。
- **Urgency justified**: 27/30 — v3.0.0 窗口期论证合理："发布前统一比发布后迁移成本低"。仍缺少发布后迁移成本的量化（影响用户数、迁移脚本复杂度），紧迫性论证偏主观。

#### 2. Solution Clarity: 110/120

- **Approach is concrete**: 38/40 — `surfaces` map + `forge surfaces` CLI + 检测逻辑表方案可复述。修订版新增了 Go struct 声明（含 `yaml:"surfaces"` tag 和 omitempty 解释）、兼容读取过渡逻辑的伪代码注释，具体性大幅提升。轻微扣分：检测逻辑表中依赖检测的具体方式（解析 JSON 的 dependencies 字段？正则匹配？）仍未明确说明。
- **User-facing behavior described**: 40/45 — CLI 命令有退出码契约表、输入输出示例、gen-journeys 调用契约。修订版解决了 iteration-1 的主要扣分点。扣分点：init TUI 确认界面的交互流程仍缺少描述（冲突信号如何"标注"？用户如何"覆盖"？文本标签还是选择菜单？），影响可测试性。
- **Technical direction clear**: 32/35 — 修订版新增了路径规范化规则（5 条）、路径段前缀匹配算法（含示例）、检测遍历策略（4 条规则含示例）。大幅改进。轻微扣分：检测逻辑的 JSON 解析方式（全量解析 dependencies 对象 vs 字符串匹配）未声明，对实现的精确指导有微小模糊。

#### 3. Industry Benchmarking: 85/120

- **Industry solutions referenced**: 32/40 — 修订版从笼统一句扩展为 4 个工具的具体分析：Cypress（检测机制、配置方式、差异分析）、Turborepo（声明式 vs 检测模式的权衡）、ESLint flat config（路径映射的行业验证）、Jest projects（声明式数组对比）。每个工具都描述了具体检测机制、配置方式和与 Forge 的差异。扣分点：仍有"自动检测 Web 项目"的概括性表述出现在 Assumptions Challenged 表中（"init 中可以做完整检测"行），但总体已大幅改进。
- **At least 3 meaningful alternatives**: 22/30 — 5 个方案含 do-nothing。修订版在比较表中增加了 Source 列，标注每个方案的来源。但前 3 个替代方案仍是 self-invented（本 proposal v1/v2），只有 Selected 方案标注了"借鉴 ESLint + Turborepo"。Jest projects 作为行业方案出现但仅作为对比参考而非替代方案。Do-nothing 合理。
- **Honest trade-off comparison**: 16/25 — 修订版增加了 Trade-off 深度分析段落，量化了"同路径多 surface 需拆子路径"的实际影响（fullstack 框架占比 15-20%，纯 API/前端占 80%+）。扣分点：比较表中各方案的 Cons 仍偏简略（每方案只有一条 Con），且未分析 Selected 方案在 monorepo 深度嵌套（>3 层）场景下的限制。
- **Chosen approach justified against benchmarks**: 15/25 — 修订版明确说明"选择检测模式是因为 Forge 目标用户是开发者个人，减少手动配置是核心价值"和"选择路径段前缀匹配而非 ESLint 的 glob 匹配以降低实现成本"。改进明显。扣分点：未分析为什么不用 Turborepo 的纯声明式模式（其确定性和可调试性优势在 CI 环境中可能更有价值），权衡不够全面。

#### 4. Requirements Completeness: 95/110

- **Scenario coverage**: 36/40 — 修订版从 6 个场景扩展到 10 个，覆盖了 iteration-1 遗漏的所有关键场景：信号冲突（#7）、已有项目过渡（#8）、CI 环境（#9）、空 surfaces（#10）。扣分点：(1) `forge surfaces` 命令本身无参数时的行为虽已声明（"空 map 也返回 0"），但 surfaces map 中只有一个 `".": web` 时的 `forge surfaces frontend` 行为未作为独立场景列出；(2) 配置文件手动编辑损坏（如 `surfaces: "not-a-map"`）的错误恢复未覆盖。
- **Non-functional requirements**: 36/40 — 修订版 NFR 大幅改进：增加了向后兼容过渡期（含具体版本计划 v3.0.0 引入/v3.1.0 移除）、路径规范化性能要求、Windows 兼容性、YAML 序列化一致性（omitempty 问题）。扣分点："向后兼容（兼容读取过渡期）"描述详细，但未声明过渡期内 `interfaces` 和 `surfaces` 同时存在且 `surfaces` 优先时的合并/冲突行为（如果用户手动同时配置了两者，`interfaces` 的值是否被完全忽略？）。
- **Constraints & dependencies**: 23/30 — 修订版明确了 `forge-init-config-sync` 的处理方式（"标记为 Superseded-by 本 proposal"），解决了 iteration-1 的扣分点。路径分隔符跨平台约束已声明。扣分点：检测遍历策略中 3 层深度限制的理由未说明（为什么是 3 而不是 2 或 4？），`node_modules` 等排除目录列表是否可配置未声明。

#### 5. Solution Creativity: 75/100

- **Novelty over industry baseline**: 30/40 — surfaces map 的路径映射在 Forge 上下文中是创新性的。修订版通过 Assumptions Challenged 表展示了多个被推翻的假设，体现了迭代思考。与行业方案的关系在 Industry Benchmarking 中有分析。扣分点：创新主要在"组合已有模式"（路径映射 + 检测 + CLI 查询），而非全新概念。价值明确但非突破性。
- **Cross-domain inspiration**: 25/35 — 修订版通过 ESLint 按路径配置和 Turborepo workspace 隔离的引用，承认了跨领域灵感来源。longest-prefix-match 的路由匹配灵感也通过路径段前缀匹配的详细描述间接体现。扣分点：未明确提及 longest-prefix-match 来自路由匹配（HTTP router、CDN 路由）这一经典领域，跨领域引用仍偏隐含。
- **Simplicity of insight**: 20/25 — "检测本就是在特定路径发现信号，直接记录路径"的核心洞察依然优雅。Assumptions Challenged 表中的"字符前缀匹配足够精确" → "Overturned" 和 "不迁移旧字段是安全的" → "Overturned" 展示了清晰的推理链。

#### 6. Feasibility: 85/100

- **Technical feasibility**: 35/40 — 修订版解决了 iteration-1 的三个主要技术问题：(1) 信号冲突增加了消歧优先级规则表；(2) 路径规范化增加了 5 条精确规则；(3) 检测遍历增加了 4 条策略规则（含 workspace 处理和排除目录）。大幅改进。扣分点：消歧优先级表中 `web` > `api` > `cli` > `tui` 的排序理由是"前端应用通常内含后端"，但未考虑 `mobile` 类型——如果 `package.json` 同时包含 `react-native` 和 `react`，优先级如何？`mobile` 未出现在优先级表中。
- **Resource & timeline feasibility**: 28/30 — 代码量估算已更新（增加了兼容读取过渡逻辑 30-50 行和 gen-journeys 适配工作量）。版本计划（v3.0.0/v3.1.0）提供了时间框架。扣分点：init TUI 确认界面的修改工作量未单独估算（冲突信号标注、用户覆盖编辑）。
- **Dependency readiness**: 22/30 — `forge-init-config-sync` 标记为 Superseded 的处理策略已声明。无外部依赖是加分项。扣分点：Superseded 的具体操作（需要修改 proposal 状态？需要 issue comment？）未说明，存在流程模糊。

#### 7. Scope Definition: 75/80

- **In-scope items are concrete**: 28/30 — 修订版将 gen-journeys 的 4 项适配工作（SKILL.md 更新、surface 字段引用更新、rule 文件重命名、rule 内容更新）全部纳入 In Scope，解决了 iteration-1 的核心扣分点。每个条目都是可交付的具体工作项。扣分点：兼容读取过渡逻辑的"移除"条件未声明——v3.1.0 还是 v4.0.0 移除？两个版本都提到了，但未给出决策条件。
- **Out-of-scope explicitly listed**: 22/25 — 3 个 out-of-scope 条目已列出，且每个都有推迟理由。修订版将"旧 config 自动迁移工具"的推迟理由改为"通过兼容读取过渡期解决，无需独立迁移工具"，解决了与 NFR 的矛盾。扣分点："下游 skill 全面适配"中列举了 gen-contracts 等具体 skill，但未说明它们在 v3.0.0 中遇到 `surfaces` 时的实际行为（报错？忽略？回退？）。
- **Scope is bounded**: 25/25 — 修订版增加了"版本计划"段落，明确了 v3.0.0 和 v3.1.0/v4.0.0 的交付边界。"后续迭代"有了具体版本锚定。

#### 8. Risk Assessment: 80/90

- **Risks identified**: 27/30 — 修订版从 4 个风险增加到 7 个，新增了 iteration-1 要求的 3 个高风险项：(1) 旧 interfaces 字段被忽略导致静默信息丢失；(2) surfaces 空 map 被 omitempty 丢弃复现 bug；(3) gen-journeys 新旧字段并存导致同步问题。扣分点：`mobile` 类型在消歧优先级表中缺失，这是一个潜在的风险遗漏——包含 `react-native` + `react` 的项目会产生未定义的检测行为。
- **Likelihood + impact rated**: 25/30 — 评分基本诚实。"同路径多 surface 拆分不直观"的 impact 已从 M 调整为 H（与 iteration-1 建议一致）。扣分点："下游 skill 引用旧字段名导致运行时错误"标为 M/H，但 out-of-scope skill 在 v3.0.0 中遇到 `surfaces` 时的行为完全未定义，impact 可能被低估。
- **Mitigations are actionable**: 28/30 — 修订版的缓解措施具体可操作：兼容读取过渡逻辑有伪代码和版本计划、omitempty 问题有明确的 struct tag 声明、gen-journeys 已纳入 In Scope。扣分点：out-of-scope skill 的兼容读取策略仍是"暂时保留对旧字段名的兼容读取"，"暂时"的具体版本和"兼容读取"的具体机制（哪些字段？什么行为？）仍不够精确。

#### 9. Success Criteria: 70/80

- **Criteria are measurable and testable**: 48/55 — 修订版从 7 个成功标准扩展到 10 个。新增了兼容读取过渡验证（#6）、命名规范统一验证（#7）、gen-journeys rule 文件重命名验证（#9）、路径匹配边界验证（#10）。改进显著。扣分点：(1) #2 "检测结果以 path → surface map 形式在 TUI 中展示，含冲突信号标注（如'检测到 web + api 信号'），并允许用户确认/编辑"——"允许用户确认/编辑"如何客观验证？需要更具体的 UI 验收标准；(2) #10 "路径匹配边界验证"只覆盖了 `frontend-new` 不匹配 `frontend` 的情况，缺少 Windows 路径分隔符 `\` 转换后的匹配验证标准。
- **Coverage is complete**: 22/25 — 修订版成功标准覆盖了所有 in-scope 条目。扣分点：检测遍历策略（In Scope 中的"路径规范化与匹配"条目）缺少对应的遍历深度和排除目录的验证标准。

#### 10. Logical Consistency: 50/90

- **Solution addresses the stated problem**: 25/35 — surfaces map 统一了两个字段，解决了同步风险。Go struct 不使用 omitempty 避免复现静默跳过 bug。兼容读取过渡逻辑确保已有项目不丢失配置。扣分点：消歧优先级表不完整——缺少 `mobile` 类型。如果项目包含 `react-native`（mobile 信号）和 `react`（web 信号），优先级未定义，解决方案无法正确处理此类冲突，与声称解决"命名不一致 bug"的目标存在微小间隙。
- **Scope ↔ Solution ↔ Success Criteria aligned**: 10/30 — 修订版大幅改进了对齐度：gen-journeys 适配有了对应的成功标准（#8、#9），路径匹配有了验证标准（#10）。但存在一个关键矛盾：**检测遍历策略中声明"最多遍历 3 层子目录"，但检测信号表中的 `apps/web/client` 示例正好是 3 层，而信号表中的路径模式（如 `pyproject.toml`/`setup.py`）暗示可能在任意深度出现。如果实际项目结构是 4 层（`packages/team-a/frontend`），检测会遗漏，但成功标准 #1 只要求"至少 3 种 surface 类型"不要求覆盖深度嵌套项目。Scope 中"检测逻辑基于文件模式匹配"的声明与遍历深度限制存在隐含矛盾——真正的文件模式匹配不应受深度限制。**
- **Requirements ↔ Solution coherent**: 15/25 — 需求和方案基本对应。但存在一个逻辑间隙：NFR 声称"兼容读取过渡期"中"当 `surfaces` 为空但 `interfaces` 不为空时"回退读 interfaces，但 Scenario #9 "CI 环境：CI 中 surfaces 不存在或为空时"期望"明确错误而非静默跳过"。如果 CI 环境中有旧的 `interfaces` 配置但无 `surfaces`，兼容读取逻辑会回退读 `interfaces` 而非报错，与 Scenario #9 的期望矛盾——应该报错还是应该兼容读取？两个需求互相矛盾。

---

### Freeform Findings Resolution

| Finding | Status | Resolution |
|---------|--------|------------|
| [high] 信号冲突时缺少优先级规则 | **Addressed** | 增加了消歧优先级规则表（web > api > cli > tui）+ 冲突处理流程 + TUI 标注。但缺少 `mobile` 类型。 |
| [high] 路径规范化规则未定义 | **Addressed** | 增加了 5 条路径规范化规则（前导./、尾随/、分隔符、..、符号链接）。 |
| [high] 前缀匹配按字符还是按路径段未声明 | **Addressed** | 明确声明"路径段前缀匹配（不是字符前缀匹配）"+ 含示例。 |
| [high] 旧 interfaces 字段被静默忽略可能导致信息丢失 | **Addressed** | 增加了兼容读取过渡逻辑（surfaces 为空时回退读 interfaces + deprecation 警告）+ 版本计划。 |
| [high] 延迟下游 skill 适配可能导致 surface 和 surfaces 字段同时存在 | **Addressed** | gen-journeys 适配已纳入 In Scope（4 项具体工作），风险表中增加了对应条目。 |
| [medium] monorepo 根目录 package.json 覆盖子目录 | **Addressed** | 增加了检测遍历策略（深度限制、排除目录、workspace manifest 处理）。 |
| [medium] interfaces 到 surfaces 的语义映射非一对一 | **Addressed** | 兼容读取过渡逻辑处理了无路径信息的旧 `interfaces` 回退。 |
| [medium] Config struct 变更缺少 Go 字段类型声明 | **Addressed** | 显式声明了 `Surfaces map[string]string \`yaml:"surfaces"\``（不带 omitempty）+ 原因说明。 |
| [medium] surfaces map 中存在未知类型值时的行为未定义 | **Addressed** | 增加了"未知类型处理策略"段落（忽略 + warn 日志 + 不传透）。 |
| [medium] forge surfaces 无匹配时的退出码和输出格式不够明确 | **Addressed** | 增加了退出码契约表（exit 0/1 + stdout/stderr）+ gen-journeys 调用契约。 |

---

### Attacks

1. **[10. Logical Consistency]**: 兼容读取与 CI 报错需求矛盾——NFR 声称"当 surfaces 为空但 interfaces 不为空时回退读 interfaces"，但 Scenario #9 期望"CI 中 surfaces 不存在或为空时应明确错误而非静默跳过"——如果 CI 环境有旧 `interfaces` 但无 `surfaces`，兼容读取会静默回退而非报错，与 Scenario #9 矛盾——需要明确：CI 环境中是否应该跳过兼容读取逻辑？是否需要 `--strict` 模式？

2. **[3. Industry Benchmarking]**: 替代方案仍偏 self-invented——比较表中 5 个方案有 3 个标注为"本 proposal v1/v2/v3"，虽然引用了 ESLint 和 Turborepo 作为灵感来源，但没有将"采用 Turborepo 纯声明式模式"或"采用 Jest projects 数组模式"作为独立的替代方案评估——这些行业验证的模式应作为独立选项出现在比较表中而非仅作为参考脚注

3. **[10. Logical Consistency]**: 消歧优先级表缺少 `mobile` 类型——"当同一个 manifest 文件同时匹配多个 surface 类型信号时，按以下优先级自动消歧"列出了 web/api/cli/tui 四种，但 `mobile`（react-native + flutter 信号）完全缺失——如果 `package.json` 同时包含 `react-native`（mobile）和 `react`（web），检测结果未定义——需要将 `mobile` 加入优先级表

4. **[4. Requirements Completeness]**: surfaces 和 interfaces 同时配置时的合并行为未定义——NFR 的兼容读取逻辑说"当 surfaces 不为空时直接使用"，但如果用户同时配置了 `surfaces: {frontend: web}` 和 `interfaces: [api, cli]`，`interfaces` 中的 `api` 和 `cli` 是否被忽略？如果用户在过渡期内手动添加了 `interfaces` 的新值，这些值会丢失——需要声明：surfaces 非空时 interfaces 是否被完全忽略，还是合并（取 surfaces 中不存在的类型）

5. **[6. Feasibility]**: 检测遍历深度限制的理由未说明——"最多遍历 3 层子目录"但未给出理由，且示例中 `apps/web/client` 已达 3 层边界——实际项目中 `packages/team-a/web-app` 这种 3 层结构很常见，而 `packages/team-a/frontend/web-app` 这种 4 层结构在大型 monorepo 中也不罕见——需要说明深度限制的选择依据或声明为可配置

6. **[9. Success Criteria]**: TUI 确认界面的验收标准不够客观——"含冲突信号标注（如'检测到 web + api 信号'），并允许用户确认/编辑"——"允许用户确认/编辑"如何验证？需要具体到：确认按钮存在、编辑入口可见、冲突信号以特定格式展示

---

### Outcome
Target NOT reached — 830/1000 (target: 900). 修订版从 645 提升至 830，大幅改进。剩余 gap 集中在 Logical Consistency（兼容读取与 CI 报错需求矛盾、mobile 类型优先级缺失）和 Industry Benchmarking（替代方案仍偏 self-invented）。Priority revisions: (1) 解决兼容读取与 CI 报错的逻辑矛盾, (2) 补全 mobile 类型的消歧优先级, (3) 声明 surfaces+interfaces 并存时的合并策略, (4) 将行业验证模式作为独立替代方案评估。
