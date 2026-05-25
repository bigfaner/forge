# Freeform Expert Review

**提案**: init-justfile Surface 感知 + 测试编排简化
**评审视角**: Config Schema & Surface Detection Engineer
**日期**: 2026-05-24

---

## Background Assessment

这份提案试图解决 Forge 测试管线中两个相互耦合的问题。第一个问题是 `init-justfile` skill 在生成 justfile 配方时不感知项目的 surface 类型，导致所有 surface 类型获得相同的配方结构——但 Web/API 项目测试时需要先启动服务、等待就绪，而 CLI/TUI 项目只需编译运行，两者的编排序列本质不同。第二个问题是 `.forge/config.yaml` 中的 `test.execution` 节点形成了一层冗余委托：`run-tests` 从 config 中读取命令模板（如 `"just test {slug}"`），模板又解析为 just 命令，而 justfile 本身已经是命令抽象层。

提案的核心技术路径是：让 init-justfile 在生成配方时读取 surface 类型，为 web/api/cli/tui/mobile 五种 surface 生成差异化的 `dev`/`run`/`probe`/`test-setup` 配方组合；同时去掉 `test.execution` 委托层，让 `run-tests` 直接根据 surface 编排模式调用 just 配方。非命令类配置（`timeout`、`results-dir`）则保留在 config.yaml 中。

提案依赖的关键假设包括：(1) surface 信息可以从 `.forge/config.yaml` 的 `surfaces` 字段或 `forge surfaces` CLI 获取；(2) `test.execution` 在 Go 结构体中从未实现，因此移除没有存量兼容性负担；(3) justfile 可以作为唯一的命令抽象层覆盖所有编排需求。

---

## Key Risks

以下按风险严重程度排列，从最关键的问题开始。

问题：提案声称"去掉 `test.execution` 委托层"，但 run-tests SKILL.md 当前的工作流高度依赖 `test.execution` 中的多个字段——不仅是 `run`，还有 `setup`、`pre-check`、`teardown`、`results-dir`、`timeout`。提案用固定的 just 配方名（`just test-setup`、`just probe`、`just test-teardown`）替代了这些动态配置字段，但这意味着编排序列的每一步都被硬编码为固定的配方名。

> "去掉 `test.execution` 后，`run-tests` 的编排变为：run-tests → just test [journey] → 实际运行器"

这看起来简化了，但实际上将"命令可配置性"从 config.yaml 转移到了 justfile 配方体中。当一个项目的 `test-setup` 不是简单的 `just test-setup`（例如 Go 项目用 `go test` 直接跑单元测试不需要 test-setup），run-tests 仍然会无条件调用 `just test-setup`，只是在 justfile 中配方体为空。这不是"简化"，而是将配置从显式声明变成了隐式约定。更关键的是，run-tests SKILL.md 第 4 步的"Environment Readiness Check"已经通过 `rules/env-check.md` 和 `surface-<type>.md` 做 surface 感知的环境检查了——提案与现有机制之间的重叠没有被讨论。

风险：提案引入了 `test.timeout` 和 `test.results-dir` 两个新的 config 顶层字段，但 Go `Config` 结构体（`forge-cli/pkg/forgeconfig/config.go`）当前没有 `Test` 字段，也没有 `test` 键的 `GetConfigValue` 处理分支。

> "`GetConfigValue` 支持 `test.*` 键读取（`forge config get test.timeout`、`forge config get test.results-dir`）"

查看 `GetConfigValue` 的实现，它当前只处理 `auto.*`、`worktree.*`、`coverage.*` 和顶级 `test-framework` 键。要支持 `test.*`，需要新增 `TestConfig` 结构体、注册新的 `getTestKeyValue` 函数、在 `SetConfigValue` 中增加 `test.*` 分支。提案在"范围内"轻描淡写地说"Config 结构体新增 `TestConfig` 节点"，但这个改动的范围涉及：Go 结构体定义、YAML 序列化/反序列化、`GetConfigValue`/`SetConfigValue` 扩展、`config` 子命令测试、以及 config-schema.md 的更新。这不是"常规 Go 开发"——它改变了 Forge CLI 的配置接口契约。如果其他 skill 也需要读取 `test.timeout`（比如 `quality_gate.go` 或 `testrunner`），影响面会超出提案预估的 8-12 个任务。

风险：混合项目（web+api）的 `dev` 配方 scope 语义定义不完整。提案说 `dev` 接受 scope 参数——`just dev frontend`、`just dev backend`、`just dev`（无 scope 时并发启动两者），但 run-tests 如何知道应该传哪个 scope？

> "`dev` 配方接受 scope 参数 — `just dev frontend`、`just dev backend`、`just dev`（无 scope 时并发启动两者）"

当前 `surfaces` 在 config.yaml 中的表示是 map 形式（如 `{frontend: web, backend: api}`）或 scalar 形式（如 `api`）。当 map 有多个 surface 时，run-tests 需要将 map 的 key 作为 scope 传给 `just dev`。但这里有一个关键歧义：map 的 key 是用户自定义的目录路径，不是 `frontend`/`backend` 这样的固定枚举值。如果用户的 config 是 `{admin-panel: web, payment-service: api}`，run-tests 应该调用 `just dev admin-panel` 还是 `just dev frontend`？提案没有定义 scope 值到 config key 的映射规则。

问题：提案说 run-tests "检测 surface → 确定编排模式"，但没有定义 run-tests 如何获取 surface 信息。

> "run-tests 通过 `forge surfaces` 或 `.forge/config.yaml` 获取 surface 类型"

这两种获取方式的语义不同：`forge surfaces [path]` 基于文件信号检测，可能返回检测结果；`.forge/config.yaml` 的 `surfaces` 字段是用户显式配置。当两者不一致时（config 说 web，文件信号检测出 api），以哪个为准？init-justfile 用哪个？run-tests 用哪个？如果 init-justfile 用检测结果生成配方，但 run-tests 用 config 读编排模式，两者可能不匹配。提案需要一个统一的 surface 信息源优先级规则。

风险：提案将 `test.execution` 的移除标记为安全——"v3.0.0 未发布，无存量用户，直接移除"——但 `test.execution` 在 run-tests SKILL.md 和 config-schema.md 中都有完整的文档定义、模板变量系统（`{slug}`、`{journey}`、`{test-dir}`、`{results-dir}`）、错误处理逻辑和 escape 规则。

> "v3.0.0 未发布，无存量用户，直接移除"

虽然 Go 结构体层面 `test.execution` 确实未实现（Config 结构体中没有 `Execution` 字段），但 run-tests skill 的 SKILL.md 文档和 config-schema.md 参考文档已经完整描述了这个接口。如果有开发者或 AI agent 按照现有文档已经配置了 `test.execution`（在 SKILL 层面它是可用的，因为 run-tests 作为 skill 通过 `forge config get` 读取 YAML 而不依赖 Go 结构体），移除后会静默忽略这些配置而不报错。提案应该要求在 run-tests 启动时检测到 `test.execution` 节点时输出废弃警告，而不是静默忽略。

问题：提案为 5 种 surface 定义了测试编排模式表，但 api 和 web 的编排序列完全相同（`dev` → `probe` → `test` → teardown），区别只在 probe 的端点不同（HTTP 端点 vs /healthz）。如果 api 和 web 的配方生成逻辑可以合并，那 5 个 surface 规则文件是否过度拆分？

> | **web** | `just dev`（后台）→ `just probe`（等待就绪）→ `just test` → teardown |
> | **api** | `just dev`（后台）→ `just probe`（等待就绪）→ `just test` → teardown |

这两个 surface 在"测试编排序列"和"关键配方"上的唯一差异是 probe 检查的目标不同——这是 init-justfile 生成 `probe` 配方时根据语言/框架决定的细节，不是 surface 规则需要区分的。将 api 和 web 合并为一个"服务型 surface"（service surface）规则，在 probe 配方生成指导中记录差异，可能更简洁。

问题：提案的"下游集成契约"表列出了 `dev` 配方的签名为 `just dev [scope]`，但 init-justfile SKILL.md 当前的 Standard Target Contract 中 `dev` 没有定义 scope 参数。

> | `dev` | `just dev [scope]` | run-tests（web/api 启动服务） | web/api: 后台启动并监听端口；cli/tui: 编译运行 |

当前 init-justfile 的混合项目 scope 参数只定义在 `compile` 和 `run` 的示例中（Step 3b），`dev` 没有明确的 scope 参数规范。提案需要将 `[scope]` 加入 Standard Target Contract 表中的 `dev` 行，并更新 init-justfile SKILL.md。

风险：`test [journey]` 过滤策略未定义。

> "`test` 必须始终接受 `journey=''` 参数"

提案提到每个 surface 规则文件应记录"journey 过滤策略"，但没有给出任何示例。不同 surface 类型的 journey 过滤方式可能根本不同：Web 用 Playwright 的 `--grep` 或配置文件过滤；API 用 Go 的 build tag 或测试函数名前缀；CLI 用环境变量或子命令选择。如果 surface 规则只是说"用 journey 参数过滤"，但没有指导如何在 just 配方中将 journey 映射到具体的测试运行器过滤机制，LLM 在生成配方时会缺乏确定性。风险表中也标记了这个问题的可能性为中、影响为高，但缓解措施只是"Surface 规则记录映射关系"——这个"记录"本身需要具体设计。

问题：提案说 init-justfile 的生成流程是"检测语言 → 加载语言模板 → 生成 compile/build/lint/fmt"然后"检测 surface → 加载 surface 规则 → 生成 test/dev/run/probe/test-setup"，这暗示了语言和 surface 是两个独立的正交维度。

> "检测语言 → 加载语言模板 → 生成 compile/build/lint/fmt" → "检测 surface → 加载 surface 规则 → 生成 test/dev/run/probe/test-setup"

但 `dev` 和 `run` 配方的内容实际上同时依赖语言和 surface：一个 Go+Web 项目的 `dev` 是 `air` 或 `go run`，而一个 Node+Web 项目的 `dev` 是 `next dev` 或 `vite dev`。提案将这两个步骤分开，但没有讨论当语言模板生成的 `dev` 与 surface 规则指导的 `dev` 发生冲突时如何仲裁。例如，init-justfile 当前 Step 3b 的 `dev` 配方由语言决定，Step 3a/3c 可能被 surface 规则覆盖——哪个优先？

风险：提案在范围内列出了"去掉 Step 3a 中对 `test.execution.run` 的依赖"，但 init-justfile SKILL.md 的 Step 3a 当前是测试配方生成的核心步骤，它的整个逻辑链是"Read Config → 读 `test.execution.run` → 用作 `test` 配方体"。

> "去掉 Step 3a 中对 `test.execution.run` 的依赖"

移除这个依赖后，`test` 配方体的来源变为"Convention + Surface 规则"。但 Step 3a 的当前逻辑有完整的 fallback 链：Config → Convention → cold start → prompt。提案只说去掉 `test.execution.run` 依赖，没有描述新的 fallback 链是什么。如果 Surface 规则也缺失（无 surface 配置），`test` 配方的生成逻辑是什么？回退到当前行为？那当前行为依赖的就是 `test.execution.run`——但这个字段已经被移除了。这里有一个循环依赖。

---

## Improvement Suggestions

建议：为 `run-tests` 的 surface 信息获取定义明确的优先级规则，并在提案中记录。
Addresses: "run-tests 如何获取 surface 信息"的问题

> What changes: 在提案的"约束与依赖"部分新增一条规则："Surface 信息优先从 `.forge/config.yaml` 的 `surfaces` 字段读取。若该字段缺失，通过 `forge surfaces` CLI 检测并回退。`init-justfile` 和 `run-tests` 使用相同的优先级规则，确保配方生成与编排消费的 surface 视图一致。" 这消除了 init-justfile 和 run-tests 对同一项目看到不同 surface 类型的风险。

建议：将 `test.timeout` 和 `test.results-dir` 的 config schema 变更从"范围内"的简略描述升级为独立的子方案，明确定义 Go 结构体变更、`GetConfigValue`/`SetConfigValue` 扩展点和迁移策略。
Addresses: `GetConfigValue` 不支持 `test.*` 键的风险

> What changes: 提案新增一个"Config Schema 变更"子节：定义 `TestConfig` 结构体（`timeout int`、`resultsDir string`），在 `Config` 中增加 `Test *TestConfig` 字段，在 `GetConfigValue` 中新增 `test.*` 处理分支，在 `SetConfigValue` 中新增 `test.*` 写入支持。同时定义当配置文件中存在旧的 `test.execution` 节点时的迁移行为：`run-tests` 在启动时检测到 `test.execution` 应输出 `DEPRECATED` 警告并忽略，而不是静默跳过。这使得改动范围从"常规 Go 开发"变为可审计的、有明确边界的接口变更。

建议：定义混合项目 scope 的映射规则，将 `surfaces` map 的 key 作为 scope 值直接传递。
Addresses: 混合项目 scope 歧义风险

> What changes: 在"下游集成契约"部分明确约束：当 `surfaces` 为 map 形式时（如 `{frontend: web, backend: api}`），`just dev` 的 scope 参数值必须是 config 中 surfaces map 的 key（如 `frontend`、`backend`），不是固定的枚举值。`init-justfile` 在生成混合项目的配方时，从 config 的 surfaces map 提取 key 列表，生成对应的 scope 分支。`run-tests` 在编排时遍历所有 surfaces，为每个需要启动的 surface 调用 `just dev <key>`。这样 scope 值与 config 保持一对一映射，没有歧义。

建议：合并 api 和 web 的 surface 规则为一个"service"规则，或者至少让它们共享编排序列模板，只在 probe 目标上差异化。
Addresses: api/web 编排序列完全相同的问题

> What changes: 将 5 个 surface 规则文件改为 4 个：`service.md`（覆盖 web+api）、`cli.md`、`tui.md`、`mobile.md`。`service.md` 中定义通用的服务编排序列，在 probe 配方生成指导中列出不同框架的默认健康检查端点（`/healthz`、`localhost:3000` 等）。如果将来确实需要 web 和 api 有不同的编排行为（比如 web 需要额外的浏览器环境检查），可以通过 service 规则内的子类型区分，而不是从一开始就维护两个完全相同的规则文件。

建议：在提案中补充 `test` 配方生成的新 fallback 链，替代被移除的 `test.execution.run` 依赖。
Addresses: Step 3a `test` 配方体 fallback 循环依赖

> What changes: 定义新的 `test` 配方生成 fallback 链：Surface 规则 → Convention 框架 → 语言模板 cold start → 报错提示。具体来说：(1) 如果检测到 surface 且对应规则文件存在，从 surface 规则中提取 `test` 配方生成指导（包括 journey 过滤策略、底层运行器调用方式）。(2) 如果无 surface 配置但 Convention 存在，从 Convention 的 Framework 和 Result Format 构建配方（当前 Step 3a 的 fallback 逻辑）。(3) 如果两者都缺失，使用 cold start 默认值。(4) 不再 prompt 用户去配置 `test.execution.run`——因为这个字段已经被移除了，改为 prompt 运行 `/init-justfile` 重新生成。

建议：在"下游集成契约"表中为 `dev` 配方增加 `[scope]` 参数的正式定义，并同步更新 init-justfile SKILL.md 的 Standard Target Contract。
Addresses: `dev` 配方 scope 参数未在 Standard Target Contract 中定义的问题

> What changes: 将 init-justfile SKILL.md 的 Standard Target Contract 表中 `dev` 行的描述更新为："Start the service / dev mode; mixed projects accept optional scope parameter (`just dev [scope]`)". Recipe Parameter Signatures 表中增加 `dev` 行：`just dev [scope]`，描述为"Optional `scope` parameter: start all services when omitted; start specific service when provided (mixed projects only)".

建议：在提案中为每个 surface 规则文件定义 journey 过滤策略的最小规范，至少包含一个完整示例。
Addresses: `test [journey]` 过滤策略未定义的风险

> What changes: 在提案的"Surface 测试编排模式"表中增加一列"Journey 过滤机制"，为每个 surface 给出 just 配方中如何将 journey 参数映射到具体运行器过滤的示例。例如：Web surface 的 `test` 配方中 `journey` 映射为 `npx playwright test --grep "{{journey}}"`；Go API surface 映射为 `go test ./tests/{{journey}}/... -tags={{journey}}`。这为 surface 规则文件的编写提供了确定性模板。
