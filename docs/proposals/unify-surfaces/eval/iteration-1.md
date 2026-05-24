# Eval-Proposal Complete
**Final Score**: 645/1000 (target: 900)
**Iterations Used**: 1/3

### Score Progression
| Iteration | Score | Delta |
|-----------|-------|-------|
| 1         | 645   | —     |

### Dimension Breakdown (final)

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Problem Definition | 100 | 110 |
| 2. Solution Clarity | 80 | 120 |
| 3. Industry Benchmarking | 50 | 120 |
| 4. Requirements Completeness | 60 | 110 |
| 5. Solution Creativity | 65 | 100 |
| 6. Feasibility | 60 | 100 |
| 7. Scope Definition | 45 | 80 |
| 8. Risk Assessment | 50 | 90 |
| 9. Success Criteria | 55 | 80 |
| 10. Logical Consistency | 80 | 90 |

### Dimension Details

#### 1. Problem Definition: 100/110
- **Problem stated clearly**: 38/40 — 核心问题明确：`interfaces` 和 `surface` 功能重叠，取值相同但机制不同。唯一模糊点是"同步风险"——"存在同步风险"的表述不够精确，具体是哪两种场景会不同步（config 与 journey 不一致？Go 代码与 schema 不一致？）需要一句补充。
- **Evidence provided**: 35/40 — 4个具体 bug 引用（含文件路径和代码行为描述），包括静默跳过、命名不一致、字段缺失、init 未配置。扣分点：缺少用户反馈或 issue tracker 链接作为外部证据，全部是作者自己发现的内部问题。
- **Urgency justified**: 27/30 — v3.0.0 窗口期论证合理："发布前统一比发布后迁移成本低"。扣分点：未量化"发布后迁移"的成本（影响的用户数、迁移脚本复杂度），使得紧迫性论证偏主观。

#### 2. Solution Clarity: 80/120
- **Approach is concrete**: 30/40 — `surfaces` map + `forge surfaces` CLI 的方案可以复述。但检测逻辑表只列了"依赖/内容检测"列，没有说明具体的检测算法（正则匹配 package.json 的 dependencies 字段？解析 JSON？），init 确认界面的交互流程也没有描述。
- **User-facing behavior described**: 25/45 — CLI 命令的输入输出有示例，但缺少关键的 CLI 契约：退出码未定义（"无匹配则报错"——stderr+非零？stdout+零？）、输出格式不够精确（`forge surfaces` 输出每行一个？空格分隔？）。gen-journeys 作为 LLM skill 需要解析这些输出，格式契约至关重要。
- **Technical direction clear**: 25/35 — Go 实现方向明确，但缺少 Go struct 定义（`Surfaces map[string]string` 的 YAML tag 是 `omitempty` 还是不是？这直接影响提案声称要修复的 bug）。路径匹配算法未声明（字符前缀 vs 路径段前缀），检测目录遍历策略未定义。

#### 3. Industry Benchmarking: 50/120
- **Industry solutions referenced**: 10/40 — 仅一句提及 Cypress 和 Postman，没有说明它们的具体做法（Cypress 如何自动检测？Postman 如何推断 API 测试？），也没有引用任何文档或源码。"例如 Cypress 自动检测 Web 项目"过于笼统。
- **At least 3 meaningful alternatives**: 20/30 — 4个方案（do nothing、保留两个概念、surfaces 数组、surfaces map + config 查询、surfaces map + 独立命令），但前3个都是 self-invented，无一个是行业验证方案。"保留两个概念，interfaces 自动检测"有可能是行业常见做法（如 Jest 的 projects 配置），但提案没有调研。Do-nothing 合理。
- **Honest trade-off comparison**: 10/25 — 对比表只列了 Pros/Cons 各一条，且 cons 描述简略。"同路径多 surface 需拆子路径"是 cons，但没有分析这个约束在实际项目中的影响面（多少项目会遇到？拆分子路径的用户体验如何？）。
- **Chosen approach justified against benchmarks**: 10/25 — 未说明为什么不采用行业常见的做法（如 Jest projects 的声明式配置、Turborepo 的 pipeline 配置），而是自创 surfaces map 模式。"路径级精度 + 关注点分离"作为理由太抽象。

#### 4. Requirements Completeness: 60/110
- **Scenario coverage**: 20/40 — 6个场景列出了 happy path（单接口、monorepo、Next.js）和部分 edge case（检测失败、检测多选、用户覆盖）。但以下场景遗漏：(1) 信号冲突场景（同一个 package.json 同时有 react 和 express）；(2) 空 surfaces 的行为（是否复现静默跳过 bug？）；(3) 已有项目不重新 init 的过渡场景；(4) CI 环境中 surfaces 不存在时的行为。
- **Non-functional requirements**: 20/40 — 检测速度（5秒）有量化。但"向后兼容"声称与实际方案矛盾——实际是"不兼容，但 init 会重新检测"，freeform review 已指出这点。缺少：路径规范化性能、Windows 兼容性、YAML 序列化一致性（omitempty 问题）。
- **Constraints & dependencies**: 20/30 — Go 实现、LLM skill 调用 CLI 的约束已声明。但遗漏：路径分隔符跨平台约束、`forge-init-config-sync` proposal 的状态（"已 Approved 但未实现"——未说明如果它先实现了怎么办）。

#### 5. Solution Creativity: 65/100
- **Novelty over industry baseline**: 25/40 — surfaces map 的路径映射确实比简单的类型列表更精确，是一个有价值的创新。但与行业方案（如 Turborepo 的 pipeline 配置、pnpm 的 workspace packages）的关系未分析，无法判断是否真的优于行业实践。
- **Cross-domain inspiration**: 20/35 — "路径映射"灵感可能来自路由匹配（longest-prefix-match），但提案没有承认这种关联。也没有参考其他工具（如 ESLint 的 override 按路径配置）。
- **Simplicity of insight**: 20/25 — "检测本就是在特定路径发现信号，直接记录路径"确实是一个简洁的洞察。

#### 6. Feasibility: 60/100
- **Technical feasibility**: 20/40 — 信号表覆盖了 4 种包管理器 + 移动端，但信号冲突问题未解决（react+express 同时存在时的优先级规则）。路径规范化规则未定义（`..` 处理、符号链接、Windows 分隔符）。目录遍历策略未声明（monorepo 子目录深度）。这些是实现前必须解答的技术问题。
- **Resource & timeline feasibility**: 25/30 — 代码量估算（150-250 行检测 + 50-100 行命令）合理。skill 文档更新范围已识别。扣分点：未考虑下游 skill 适配的工作量（gen-journeys 至少需要同步更新）。
- **Dependency readiness**: 15/30 — `forge-init-config-sync` "已 Approved 但未实现"，且本 proposal "可替代它"——但未说明替代意味着什么（是否需要先撤销那个 proposal？是否有冲突的代码变更？）。无外部依赖是加分项。

#### 7. Scope Definition: 45/80
- **In-scope items are concrete**: 25/30 — 7个 in-scope 条目都是可交付的具体工作项。扣分点："gen-journeys skill：从独立 surface 检测改为调用 `forge surfaces <path>` 查询"——但没有包含 gen-journeys 的 surface rule 文件重命名（`surface-webui.md` -> `surface-web.md`），freeform review 指出这个遗漏会导致字段名不一致。
- **Out-of-scope explicitly listed**: 10/25 — 3 个 out-of-scope 条目已列出。但"下游 skill 全面适配"被推迟，而 gen-journeys 是数据流的核心消费者——推迟它意味着同一版本内新旧字段并存。freeform review 标记为高风险。"旧 config 自动迁移工具"被列为 out-of-scope，但 NFR 声称"向后兼容"，两者矛盾。
- **Scope is bounded**: 10/25 — 没有明确的时间线或版本绑定（"v3.0.0 之前完成"？）。"后续迭代"没有定义。缺少版本计划。

#### 8. Risk Assessment: 50/90
- **Risks identified**: 15/30 — 4 个风险已列出，但遗漏了以下重要风险：(1) 旧 `interfaces` 字段被忽略导致的静默信息丢失（freeform review 标为 high）；(2) surfaces map 空 map 被 omitempty 丢弃复现同一 bug；(3) gen-journeys 新旧字段并存导致同步问题（freeform review 标为 high）。
- **Likelihood + impact rated**: 15/30 — 评分基本合理，但"下游 skill 引用旧字段名"标为 M/H 是诚实的。扣分点："同路径多 surface 拆分不直观"标为 M/M，但实际上 Next.js 等 fullstack 框架非常常见，impact 可能是 H。
- **Mitigations are actionable**: 20/30 — "init 确认界面允许用户手动编辑"是可操作的。但"out-of-scope skill 暂时保留对旧字段名的兼容读取"——"暂时"没有定义期限，兼容读取的具体策略也未定义。

#### 9. Success Criteria: 55/80
- **Criteria are measurable and testable**: 40/55 — 7 个成功标准大部分可测试（检测 3 种类型、map 格式写入、CLI 命令行为）。扣分点：(1) "检测结果以 path → surface map 形式在 TUI 中展示"——"展示"如何测试？需要更具体的 UI 验收标准；(2) "旧 `interfaces` 字段的命名不一致 bug 被修复"——如何验证？需要具体的测试场景（如 `interfaces: ["web-ui"]` 时 `forge task index` 应正确生成任务）。
- **Coverage is complete**: 15/25 — 缺少以下 in-scope 条目的成功标准：(1) 命名规范统一（`web-ui` -> `web`）的具体验证；(2) gen-journeys rule 文件更新的验证；(3) 路径匹配边界情况的验证（`frontend` vs `frontend-new`）。

#### 10. Logical Consistency: 80/90
- **Solution addresses the stated problem**: 30/35 — surfaces map 统一了两个字段，解决了同步风险和命名不一致。扣分点：提案声称要修复"interfaces 为空时静默跳过"的 bug，但如果新字段用了 `omitempty`，空 surfaces 也会被静默丢弃，复现同一类 bug。
- **Scope ↔ Solution ↔ Success Criteria aligned**: 25/30 — 基本对齐。但 success criteria 缺少对 "命名规范统一"（in-scope 第2项）和 "gen-journeys rule 文件更新"（in-scope 第7项）的覆盖。
- **Requirements ↔ Solution coherent**: 25/25 — 需求和方案基本对应，无孤立需求。

---

### Attacks

1. **[3. Industry Benchmarking]**: 行业调研严重不足——仅笼统提及 Cypress 和 Postman，无具体做法描述，无文档/源码引用——"例如 Cypress 自动检测 Web 项目、Postman 自动推断 API 测试"是 hand-waving，需要在每个行业方案中说明其具体检测机制、配置方式、与 Forge 的差异

2. **[4. Requirements Completeness]**: 信号冲突场景完全遗漏——"检测输出直接就是 path → surface 的 map，无需额外转换"暗示检测结果无歧义，但同一个 package.json 可同时包含 react（web 信号）和 express（api 信号），提案未定义消歧优先级规则——需要增加冲突消歧规则表

3. **[4. Requirements Completeness]**: 向后兼容声称不实——NFR 声称"旧的 `interfaces` 字段应被识别并迁移"，但 Out-of-Scope 明确排除"旧 config 自动迁移工具"，实际方案是"init 重新检测会覆盖，旧字段可忽略"——如果不重新 init，旧字段被完全忽略，forge task index 静默停止生成测试任务——需要在 NFR 中改为"通过兼容读取过渡期实现向后兼容"或将迁移工具纳入 In Scope

4. **[2. Solution Clarity]**: 路径匹配算法未声明——"最长前缀匹配"未定义是字符前缀还是路径段前缀——按字符前缀匹配会导致 `frontend` 错误匹配 `frontend-new`——需要声明按路径段（path segment）匹配

5. **[2. Solution Clarity]**: 路径规范化规则未定义——配置键格式声明了"无前导 `./`，无尾随 `/`"，但查询时传入的路径格式未约束——`./frontend`、`frontend/`、`frontend/src`、`../other` 的规范化逻辑未声明——需要增加路径规范化规则段落

6. **[2. Solution Clarity]**: CLI 退出码契约缺失——"无匹配则报错提示手动指定"未定义退出码和输出流——gen-journeys 作为 LLM skill 需要区分"成功"和"未找到"——需要声明：未匹配时退出码 1 + stderr 输出错误信息

7. **[7. Scope Definition]**: gen-journeys 核心适配被推迟会导致字段并存——gen-journeys 是数据流的核心消费者（"gen-journeys 通过 `forge surfaces <path>` 查询 surface"），但 surface rule 文件重命名和旧 `surface` 字段引用更新未纳入 In Scope——如果 v3.0.0 发布时 gen-journeys 仍写 `surface: webui` 而 Go 代码只读 `surfaces` map，会形成新的同步问题——需要将 gen-journeys 的字段和 rule 文件更新纳入 In Scope

8. **[8. Risk Assessment]**: 遗漏 3 个高风险项——(1) 旧 interfaces 字段被忽略导致静默信息丢失（high）；(2) surfaces 空 map 被 omitempty 丢弃复现静默跳过 bug（high）；(3) gen-journeys 新旧字段并存导致同步问题（high）——需要在风险表中增加这 3 项

9. **[6. Feasibility]**: 目录遍历策略未定义——在 pnpm workspace 或 yarn workspaces 的 monorepo 中，根目录的 package.json 只声明 workspace 配置，实际依赖在子目录——提案未说明检测的遍历深度、排除目录列表、workspace manifest 的特殊处理——需要增加 "Detection Traversal Strategy" 段落

10. **[6. Feasibility]**: Config struct YAML tag 未声明——当前 `Interfaces` 用 `omitempty`（推测），如果新 `Surfaces map[string]string` 也用 `yaml:"surfaces,omitempty"`，空 map 会被序列化时省略，Go 解析为 nil map，复现提案声称要修复的"静默跳过"bug——需要在提案中显式声明 YAML tag 不使用 omitempty

11. **[beyond-rubric]**: `[blindspot]` 提案未声明 surfaces map 中存在未知类型值时的行为——`forge task index` 从 surfaces 提取去重类型列表时，如果遇到 `"frontend": "unknown-type"`，是忽略、报错还是传透？当前 `hasUIInterface` 对未知类型返回 false，但新逻辑的行为未声明——需要在 Requirements 中定义未知类型的处理策略

12. **[beyond-rubric]**: `[blindspot]` interfaces 到 surfaces 的语义映射非一对一——旧的 `interfaces` 是去重类型列表 `["api", "cli"]`，无路径信息。迁移时无法自动生成 path → surface 映射。如果选择不迁移而是要求 re-init，对已有 worktree 和 CI 环境缺少过渡方案——需要定义至少一个版本的兼容读取过渡期

### Outcome
Target NOT reached — 645/1000 (target: 900). Significant gaps in Industry Benchmarking, Requirements Completeness, and Feasibility. Priority revisions: (1) add signal conflict resolution rules, (2) define path normalization and matching semantics, (3) add backward compatibility read-through, (4) move gen-journeys core adaptation into In Scope, (5) conduct real industry benchmarking.
