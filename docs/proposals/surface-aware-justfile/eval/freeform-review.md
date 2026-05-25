# Freeform Expert Review

**提案**: init-justfile Surface 感知 + 测试编排简化
**评审视角**: Surface-Aware Dispatcher & Test Orchestration Architect（调度器模式 / 规则文件驱动编排 / 跨平台进程生命周期管理）
**日期**: 2026-05-25

---

## Background Assessment

这份提案试图在 Forge v3.0.0 中解决两个耦合问题：其一，`init-justfile` 不感知 surface 类型，导致 Web/API 和 CLI/TUI 这两类本质不同的测试编排流程获得相同的 justfile 配方结构；其二，`test.execution` 配置节点在 run-tests 与实际测试运行器之间形成了一层冗余间接委托。提案选择将两者捆绑解决——先移除 `test.execution` 委托层，然后通过规则文件（`rules/surfaces/<type>.md`）同时指导 init-justfile 的配方生成和 run-tests 的编排执行。

核心技术路径是一个"检测 surface → 加载策略规则 → 按策略执行"的调度器模式。init-justfile 在检测到 surface 类型后加载 `rules/surfaces/<type>.md`，按规则生成 `dev`/`probe`/`test`/`test-setup`/`test-teardown` 等编排级配方；run-tests 在运行时同样加载 `rules/surfaces/<type>.md`，按规则中定义的编排序列和 just 配方调用契约执行。提案声称这一模式与现有的 `gen-test-cases` skill 同构——都是"检测类型 → 加载规则 → 按策略执行"。

提案的第二条主线是 scope 值域的统一迁移。当前混合项目的 scope 使用固定枚举 `frontend`/`backend`/`all`，提案将其迁移为 `config.yaml` 中 `surfaces` map 的用户自定义 key（如 `admin-panel`/`payment-service`）。这个迁移涉及 7 个以上组件，提案要求在同一提交中完成所有变更以保证原子性。

提案依赖的关键假设包括：(1) 规则文件可以同时服务两个消费者（init-justfile 的配方生成和 run-tests 的编排执行）而不产生职责膨胀；(2) LLM agent 通过参数化模板 + 退出码约束 + HARD-GATE 规则 + 状态机驱动可以实现编排序列的确定性下限；(3) just 原生的 `[linux]`/`[windows]` recipe attribute 可以在三个目标平台上可靠地实现后台启动、PID 追踪和 teardown 清理；(4) v3.0.0 未发布，因此 `test.execution` 的移除不涉及存量用户迁移。

提案对技术细节的覆盖相当深入——probe 重试循环中的 PID 存活检查、teardown 的命令行匹配防误杀、Windows 上 PowerShell 的 Get-CimInstance 替代已弃用的 wmic、scope 迁移的四阶段原子性保证——这些都是实现层面的关键决策。但正是这种深入的细节描述，暴露了提案在架构层面的一些值得关注的问题。

---

## Key Risks

**调度器模式同构性验证**。提案反复强调 run-tests 与 gen-test-cases 在调度模式上"同构"，但仔细审查后，两者的规则文件加载时机和参数传递路径存在微妙但重要的差异。

问题：gen-test-cases 加载规则文件时，它的输入是确定的（测试契约文件），输出也是确定的（测试代码）。规则文件指导的是生成逻辑，整个流程是一次性的"读取 → 生成 → 完成"。而 run-tests 加载规则文件后，它面对的是一个有状态的执行序列——需要在前一步的退出码基础上决定下一步动作，需要处理 probe 超时后的 teardown 路径，需要追踪 PID 文件的生命周期。

> "run-tests 不自行推断编排模式，也不依赖跨 skill 文件传递编排参数。它遵循与 gen-test-cases 相同的调度模式：1. 检测 surface type 2. 加载执行策略 3. 按策略执行"

"同构"这个术语暗示两者在架构上可以互换或共享基础设施，但实际上它们处理的是根本不同的问题域：一个是代码生成（无副作用），另一个是进程编排（有副作用、有状态、有失败恢复）。将两者标记为"同构"可能导致后续实现者在共享抽象层时过度泛化——例如试图为两者设计统一的"规则加载器"接口，而这个接口需要同时处理生成时的一次性读取和运行时的有状态编排，增加了不必要的复杂度。

**规则文件双重职责的职责边界问题**。提案将 `rules/surfaces/<type>.md` 同时定位为 init-justfile 的配方生成指导和 run-tests 的编排序列定义。这两个消费者对规则文件的需求有本质差异。

> "每个 surface 类型的规则文件定义两部分内容：编排序列（该 surface 类型的测试执行流程）和 just 配方调用契约（序列中每个步骤调用的 just 配方名、参数、退出码语义）"

风险：当规则文件同时服务于"告诉 init-justfile 应该生成哪些配方"和"告诉 run-tests 应该按什么顺序调用这些配方"两个职责时，一个职责的变更可能意外影响另一个职责的消费者。例如，如果 run-tests 需要在编排序列中增加一个新步骤（如"在 probe 之前检查数据库连接"），这个变更会同时修改规则文件，而 init-justfile 可能需要为新步骤生成对应的配方——但如果 init-justfile 没有同步更新对规则文件的解析逻辑，它不会生成新步骤的配方，导致 run-tests 调用不存在的 just 配方。提案没有定义规则文件的两个职责之间是否存在变更同步约束，也没有定义当两个消费者的需求发生分歧时如何仲裁。

此外，提案描述的是 run-tests 的规则文件路径为 `rules/surfaces/<type>.md`（在 run-tests skill 目录下），而 init-justfile 的规则文件路径为 `skills/init-justfile/rules/surfaces/<type>.md`。如果是同一个文件（提案的"创新亮点"暗示这一点），那么两个 skill 共享一个规则文件意味着该文件必须同时满足两个 skill 的格式要求——这是一个隐式的耦合点。如果是不同文件（各 skill 目录下各一份），那么"同构"的说法就不成立，因为两者的规则文件可以独立演进。

> "将 init-justfile（配方生产者）和 run-tests（调度器）的 surface 感知通过规则文件统一设计"

提案需要明确澄清：规则文件是物理上共享的一个文件（两个 skill 指向同一个路径），还是逻辑上同构但物理上独立的两个文件。这个区分对实现和后续维护至关重要。

**scope 迁移原子性的实际可达性**。提案要求 4 个阶段（数据模型 → 规则引擎 → 模板层 → prompt 模板）在同一提交中完成，理由是分批提交会导致中间状态的不一致。

> "阶段 1-4 必须在同一提交中完成。如果分批提交，中间状态会导致 scope-assignment 输出 `admin-panel` 但 init-justfile 仍期望 `frontend` 的不一致"

问题：涉及 7 个以上组件的原子迁移在实际的 code review 和合并流程中面临严峻挑战。一个大提交修改了 `prompt.go`、`scope-assignment.md`、`db-schema.md`、`quick-tasks SKILL.md`、`init-justfile SKILL.md`、16 个 prompt 模板——这些变更跨越 Go 代码、Markdown 规则文件和 prompt 模板三种不同性质的文件。code review 需要同时理解 Go 逻辑变更、规则文件语义变更和模板语法变更，review 的认知负荷极高。如果 review 过程中发现某个阶段的实现有问题需要返工，整个提交都需要重新审查。

更关键的是，提案定义了过渡期兼容层（`resolveScope()` 中旧枚举到 surfaces map key 的映射），但这个兼容层本身引入了一个新的问题。

> "若 scope 值为旧枚举（`frontend`/`backend`），在 surfaces map 中查找对应的 surface key 并映射（`frontend` → 找到 type=web 的 key，`backend` → 找到 type=api 的 key）"

风险：这个映射假设了 `frontend` 总是映射到 type=web、`backend` 总是映射到 type=api，但 surfaces map 是用户自定义的。如果用户定义了 `surfaces: {admin-panel: web, payment-service: api, internal-tool: tui}`，那么 `frontend` 可以映射到 `admin-panel`（唯一 type=web），`backend` 可以映射到 `payment-service`（唯一 type=api）。但如果用户定义了 `surfaces: {admin-panel: api, payment-service: api}`（两个 api、零个 web），那么 `frontend` 就无法映射——兼容层在"找不到 type=web 的 surface"时的行为是什么？提案定义了多 surface 同类型冲突的消歧规则（优先声明顺序靠前的 key），但没有定义"类型完全不存在"时的行为。这可能导致兼容层在某些配置下静默失败。

**PID 命令行匹配的跨平台可靠性**。提案在 teardown 的命令行匹配中详细定义了三个平台的实现方式，但存在一些边界条件未被完全覆盖。

> "Linux/macOS：检查 `/proc/<pid>` 是否存在（Linux）或 `ps -p <pid>` 是否成功（macOS）；进一步校验命令行——读取 `/proc/<pid>/cmdline`（Linux）或 `ps -p <pid> -o command=`（macOS），确认命令行包含预期的 dev server 关键词"

风险：macOS 上 `ps -p <pid> -o command=` 的输出存在截断问题。macOS 的 `ps` 命令在 `-o command=` 模式下默认截断长命令行到终端宽度（通常 80 或 120 字符）。对于复杂的 dev server 启动命令（如 `node /Users/developer/project/node_modules/.bin/next dev --port 3000 --experimental-app`），截断后的输出可能不包含预期的关键词（如 `npm run dev`），因为 `ps` 显示的是实际可执行路径而非 shell 别名。如果 init-justfile 通过 `npm run dev` 启动 dev server，但 `ps` 显示的是 `node /path/to/next dev`，命令行匹配会失败，teardown 会误判为"PID 已被回收"而跳过终止步骤。提案要求命令行包含"预期的 dev server 关键词（如 `npm run dev`、`go run`）"，但没有讨论关键词匹配策略应该匹配的是用户命令（npm run dev）还是实际进程路径（node .../next dev）。

> "Windows：通过 PowerShell `Get-CimInstance Win32_Process -Filter "ProcessId=<pid>" | Select-Object CommandLine` 校验命令行"

Windows 的 `Get-CimInstance` 在某些受限环境下（如某些企业策略禁止 WMI 查询）会失败。虽然这是边缘场景，但提案没有定义 `Get-CimInstance` 失败时的 fallback 行为——是跳过命令行校验直接终止（风险：可能杀错进程），还是跳过终止（风险：可能遗留孤儿进程）？

**LLM 编排确定性的失败模式深度分析**。提案定义了四层防御（参数化模板 + 退出码门控 + HARD-GATE 规则 + 状态机驱动），并声明"确定性下限"足以保证最坏情况下的系统安全。让我逐一分析每层防御的实际失败模式和最坏后果。

> "HARD-GATE 规则：probe 失败后禁止重试 probe 或重试 dev——唯一允许的下一步是执行 teardown 后中止"

问题：HARD-GATE 规则对 LLM agent 而言是 SKILL.md 中的文本指令，不是运行时强制的程序约束。如果 LLM agent 忽略了 HARD-GATE 标记（这在 LLM 驱动的系统中是已知的失效模式——当指令过长或上下文窗口接近饱和时，末尾的 HARD-GATE 规则可能被"遗忘"），四层防御的最外层（HARD-GATE）就失效了。此时退化为三层防御：参数化模板 + 退出码门控 + 状态机驱动。

退出码门控是最关键的确定性层。但退出码的含义在 `just dev` 后台启动场景下有一个已知的盲区——提案自己也承认了这一点。

> "just dev 后台启动 dev server 后配方本身以 exit 0 返回（因为后台进程已分离），但 dev server 可能在启动后的几秒内崩溃"

这意味着退出码门控在 `just dev` 这一步是不可靠的。提案的缓解措施是通过 probe 重试循环中的 PID 存活检查来加速崩溃检测。

> "每次 probe 重试失败后，检查 `.forge/dev-server.pid` 中记录的进程是否仍然存活。如果进程已死，probe 立即退出并报告'dev server 崩溃'，不等待剩余重试"

风险：PID 存活检查本身存在一个微妙的竞态——PID 文件被写入的时刻和 dev server 实际 fork 出子进程的时刻之间有时间差。如果 dev server 启动脚本是 `npm run dev`，nohup 后台执行后 `echo $!` 写入的是 npm 进程的 PID，但 npm 可能会 fork 出 node 子进程然后自身退出（某些 npm 版本的行为）。此时 `.forge/dev-server.pid` 中记录的 PID 指向已退出的 npm 进程，而实际的 dev server（node 进程）是孤儿进程。probe 的 PID 存活检查会发现"PID 已死"，但 teardown 尝试终止的也是这个已死的 PID——真正的 node 进程仍在运行但无法被追踪。

这个场景的最坏后果是：dev server 以孤儿进程形式运行在后台，占用端口和资源，而 Forge 系统认为它已经被清理。提案的 `test-state.json` 异常恢复路径无法捕获这种情况，因为系统认为 teardown 已成功完成（PID 文件中的进程确实已被终止，只是不是正确的进程）。

**混合项目多服务端口检查的时序影响**。

> "顺序启动策略：just dev 配方体按顺序启动各 scope 的 dev server（先启动后端，再启动前端），每个启动后将 PID 写入 `.forge/dev-server.<scope>.pid`"

> "probe 顺序：run-tests 先 probe 后端（api），再 probe 前端（web）。后端就绪是前端可用的前提条件"

风险：在多 scope 场景下，probe 的总等待时间是线性叠加的。如果后端 probe 需要 30 次重试（最多 60 秒），前端也需要 30 次重试（最多 60 秒），最坏情况下的总 probe 等待时间为 120 秒。提案没有讨论这个等待时间对开发者体验的影响，也没有提供并行 probe 的可能性。虽然后端是前端的依赖前提（因此后端 probe 必须先完成），但后端 probe 完成后、前端 probe 期间，系统处于"等待前端启动"的状态，这段等待中没有任何有意义的进度反馈——如果同时还有端口冲突导致的 probe 反复失败，用户体验会非常差。

> "probe 超时作为端口冲突的最终兜底检测机制"

这意味着在端口冲突场景下，用户需要等待完整的 probe 超时时间（每个 scope 最多 60 秒）才能看到有意义的错误信息。在 3 个 scope 的混合项目中，这可能意味着 3 分钟的无效等待。提案承认端口检查是 best-effort，但 probe 超时兜底的实际等待时间对开发者体验的影响未被充分评估。

**新增 surface 类型扩展方式的隐含依赖**。提案声称新增 surface 类型只需两步，不涉及 Go 代码变更。

> "新增 surface 类型只需两步：1. 在 init-justfile 的 `rules/surfaces/` 下新增 `<type>.md` 2. 在 run-tests 的 `rules/surfaces/` 下新增 `<type>.md`。无需修改 config schema、无需新增中间文件、无需更新 Go 代码"

问题：`config.yaml` 中的 `surfaces` 字段值域是否由 Go 代码的枚举校验？如果 `surfaces: {my-service: custom-type}` 中的 `custom-type` 不在 Go 代码的已知 surface 类型列表中，Forge CLI 是否会拒绝这个配置？如果 CLI 会拒绝，那么新增 surface 类型就需要修改 Go 代码（更新枚举），"无需更新 Go 代码"的声明就不成立。如果 CLI 不拒绝未知类型，那么 typo（如 `surfaces: {my-service: wbe}` 而非 `web`）会导致 init-justfile 和 run-tests 静默生成无效的编排序列——它们会尝试加载 `rules/surfaces/wbe.md`，找不到后回退到什么行为？提案定义了 init-justfile 的 fallback 链（surface 规则 → convention 框架 → 语言模板 cold start → 报错提示），但没有定义 run-tests 找不到对应 surface 规则文件时的行为。

---

## Improvement Suggestions

建议：明确规则文件的物理归属和职责分离机制。要么将规则文件物理拆分为"生成指导"和"编排定义"两个独立文件（如 `rules/surfaces/<type>-generation.md` 和 `rules/surfaces/<type>-orchestration.md`），要么在单一文件中用显式的节标记（如 `## Generation Guidance` 和 `## Orchestration Definition`）分隔两个职责，并定义消费者只读取自己职责对应的节。同时定义规则文件的变更同步约束——当编排定义节变更时，是否强制要求同步更新生成指导节。
Addresses: 规则文件双重职责的职责边界不清晰
> What changes: 提案"创新亮点"部分和"执行策略规则文件"部分需要澄清：(1) `rules/surfaces/<type>.md` 是一个物理文件被两个 skill 共享，还是两个 skill 各有独立的副本；(2) 如果是共享文件，定义两个职责节的分隔格式和消费者只读约束；(3) 定义当一个消费者的需求变更时，是否需要同步验证另一个消费者不受影响。

建议：为 scope 迁移原子性提供更务实的实施策略。将"同一提交"约束弱化为"同一 PR，但允许逻辑分组的多个提交"，配合过渡期兼容层确保每个独立提交都不破坏系统一致性。同时补充兼容层在"类型完全不存在"时的行为定义（如输出警告并回退 `all`）。
Addresses: scope 迁移 7+ 组件原子提交的实际执行难度和兼容层的边缘场景
> What changes: 在"迁移顺序约束与原子性保证"部分：(1) 将"同一提交"约束改为"同一 PR"，允许按阶段分提交但 PR 合并时原子入主干；(2) 兼容层增加"类型不存在"的处理——若旧值 `frontend` 在 surfaces map 中找不到 type=web 的 surface，输出警告"scope 'frontend' 无法映射到任何 surface，回退为 'all'"并使用 `all`；(3) 在 PR 模板中增加 checklist 确认所有 7 个组件已同步更新。

建议：将 PID 命令行匹配策略从"匹配用户命令"改为"匹配进程可执行文件路径的关键部分"，并为 macOS 的 `ps` 截断提供明确的处理策略。
Addresses: macOS `ps` 截断和命令行匹配目标不一致导致 teardown 误判
> What changes: 在"Teardown 进程回收机制"部分：(1) 命令行匹配的关键词应同时包含用户命令（`npm run dev`）和可能的实际进程路径（`node.*next`、`go run`），用 OR 逻辑匹配而非精确匹配；(2) 为 macOS 显式使用 `ps -p <pid> -o command= -w`（`-w` 标志禁用截断）或 `ps -p <pid> -o args=` 获取完整命令行；(3) 定义命令行匹配失败的 fallback 行为：输出完整的不匹配信息（预期关键词 + 实际命令行），由用户判断是否安全终止，而非静默跳过。

建议：为 probe 超时场景提供更快的失败反馈，特别是在多 scope 混合项目中。考虑在 probe 重试循环中增加早期退出条件——连续 N 次收到"连接被拒绝"（而非超时）时，降低剩余重试次数或缩短重试间隔，因为"连接被拒绝"通常意味着服务根本不会启动。
Addresses: 多 scope 场景下 probe 线性叠加的等待时间过长
> What changes: 在"Probe 轮询逻辑"部分：(1) 区分"连接被拒绝"（connection refused）和"连接超时"（connection timeout）两种失败模式——前者意味着服务未监听端口，后者可能只是启动慢；(2) 连续 5 次收到"连接被拒绝"时，在输出中提示"服务可能未能启动，建议检查 .forge/dev-server.log"，但继续 probe 直到 PID 存活检查或超时；(3) 考虑在后端 probe 完成后、前端 probe 期间，提供进度输出（如 `[probe] [frontend] [retry 5/30]`），避免用户在长时间等待中失去耐心而手动中断。

建议：定义 run-tests 在加载不到对应 surface 规则文件时的明确行为，并验证 config.yaml 的 surfaces 值域是否由 Go 代码校验。
Addresses: 新增 surface 类型扩展方式的隐含依赖和 typo 风险
> What changes: 在"新增 surface 类型的扩展方式"部分：(1) 验证 `forge surfaces` CLI 和 Go 代码是否校验 surface 类型值域——如果是，补充说明"新增 surface 类型需同步更新 Go 代码中的枚举定义"；如果不是，增加 typo 防护——init-justfile 在加载规则文件前校验 surface 类型值是否在已知列表中，未知值输出警告并回退到 fallback 链；(2) 定义 run-tests 加载不到规则文件时的行为：输出明确错误"surface 类型 '<type>' 没有对应的编排规则文件，请检查 config.yaml 中的 surfaces 配置或添加 rules/surfaces/<type>.md"——而非静默回退到默认行为。

建议：为 teardown 的 PID 追踪增加 npm/node 进程树识别的说明，确保 PID 文件记录的是实际的服务进程而非 npm wrapper 进程。
Addresses: npm fork 子进程导致 PID 指向已退出的 wrapper 而非实际 dev server
> What changes: 在"后台进程管理"的 PID 追踪部分：(1) 明确说明 `echo $!` 记录的是 nohup 后台命令的直接子进程 PID（在 `npm run dev` 场景下是 npm 进程）；(2) 如果 npm 进程退出而 node 子进程继续运行（孤儿进程），teardown 的命令行匹配应该能捕获这种情况——当 PID 无效但端口仍被占用时，通过端口反查（`lsof -i :$PORT`）定位实际的进程 PID 并终止；(3) 考虑在 dev 配方中使用 `exec` 替换 npm 进程（如 `exec npm run dev`），使 PID 文件直接指向 node 进程而非 npm wrapper——但这需要验证 `exec` 在 nohup 后台场景下的行为是否一致。
