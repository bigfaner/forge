---
created: 2026-05-14
author: faner
status: Draft
---

# Proposal: Interface-Type-Specific Verification Strategies

## Problem

Forge 的测试生成 pipeline 不区分 interface 类型的验证策略。gen-test-cases 将所有 UI 测试用例归为同一类型，生成相同的验证条件；gen-test-scripts 使用同一套模板逻辑生成测试代码。但不同 interface 类型（TUI、web-ui、mobile-ui、API、CLI）的"正确性"定义完全不同：

- TUI 的正确性是**视觉渲染正确性**（对齐、溢出、宽度），不是 DOM 状态
- web-ui 的正确性是**DOM 交互正确性**（元素可见、可点击、状态变化）
- mobile-ui 的正确性是**设备适配正确性**（触摸、方向、平台差异）
- API 的正确性是**契约正确性**（请求/响应符合 spec、错误路径完整）
- CLI 的正确性是**输出正确性**（退出码、输出格式、参数组合）

### Evidence

**TUI lesson**（docs/lessons/lesson-tui-visual-verify.md）：deep-drill-analytics feature 的 11 个 bug 全部通过了"编译 + 测试"verify gate，因为测试只检查了逻辑正确性（函数返回值），没有检查视觉渲染正确性（行数、宽度、对齐）。根因是 gen-test-cases 不为 TUI 类型生成 golden file 测试用例和维度检查用例。

**API/CLI 盲区**：当前 profile 声明了 `api` 和 `cli` capability，但 gen-test-cases 只生成功能性测试用例，不生成契约测试、边界值测试、参数组合测试等集成层面的测试。

### Urgency

每次 TUI feature 都需要手动在 task verify criteria 中追加 golden file 和维度检查条件，依赖人工记忆。lesson 已经证明这条路不可靠——漏掉的 bug 比发现的多。

## Proposed Solution

### 核心机制

1. **Profile 策略文件**：每个 profile 新增 `verification-strategies.md`，定义其各 capability 的验证策略（包含验证维度、边界场景、测试数据要求、测试级别标记）
2. **类型化用例生成**：gen-test-cases 读取 profile 策略，按 interface 类型生成不同的验证条件，并自动标记测试级别（e2e / integration）
3. **级别化脚本生成**：gen-test-scripts 根据测试级别和 interface 类型选择不同的代码生成策略

### 类型验证策略矩阵

| Interface Type | Test Level | 核心验证维度 | 边界场景 |
|---------------|------------|-------------|---------|
| **TUI** | e2e | Golden file 对比、维度检查（行数=高度/宽度<=终端宽度）、ANSI 色码一致性 | CJK 字符、长路径(>50)、多位数字(>9)、空字段、窄终端(80x24)、宽终端(140x40) |
| **web-ui** | e2e | DOM 交互、视觉回归（截图）、响应式布局、可访问性 | 空状态、加载态、错误态、边界数据量、不同视口尺寸 |
| **mobile-ui** | e2e | 设备渲染、触摸交互、屏幕方向、平台差异 | 横竖屏切换、低网速、推送通知中断、小屏设备 |
| **API** | integration | 契约验证（请求/响应对 spec）、错误路径、边界值、真实依赖 | 无效输入、认证失败、超限、并发、空响应、大数据量 |
| **CLI** | integration | 输出 golden file、退出码、参数组合、管道兼容 | 无效参数、--help、参数互斥、空输入、管道+重定向 |

### 测试级别定义

| Level | 含义 | 触发条件 | 验证方式 |
|-------|------|---------|---------|
| **e2e** | 端到端，验证用户可见的完整行为 | interface 类型为 visual/interactive（TUI、web-ui、mobile-ui） | 渲染输出、DOM 状态、截图、设备行为 |
| **integration** | 集成测试，验证组件间协作行为 | interface 类型为 non-visual（API、CLI） | HTTP 契约、退出码、输出格式、参数解析 |

### Developer Workflow — Before & After

#### Before（当前行为）

```markdown
## TC-001: 验证 analytics panel 显示正确数据
- 操作：启动 tui-app，选择 analytics panel
- 预期：panel 显示正确的统计数据
- Level: （无此字段）
```

验证条件是"显示正确数据"——无法检测行溢出、CJK 对齐、窄终端截断等渲染 bug。

#### After（本提案实施后）

```markdown
## TC-001: 验证 analytics panel 渲染正确性
- Interface: tui
- Level: e2e
- 操作：启动 tui-app --width 80 --height 24，选择 analytics panel
- 预期（渲染维度）：
  - 输出行数 <= 24 行（无溢出）
  - 每行宽度 <= 80 字符（无截断）
  - Golden file 与 testdata/analytics-80x24.golden 完全匹配
- 预期（边界场景）：
  - CJK 路径不导致列错位
  - 数字 > 9999 时格式化不溢出
```

#### gen-test-scripts 行为变化

Before：所有用例生成相同结构的测试函数，无级别区分。

After：
- e2e 用例生成渲染+golden file 对比函数（读取 testdata/*.golden，截获 ANSI 输出逐行比对）
- integration 用例生成 HTTP 断言或子进程退出码检查函数（不涉及渲染）

开发者可观察到：test-cases.md 多出 `Level` 和 `Interface` 字段；测试脚本按级别分目录（如 `e2e/tui_test.go` vs `integration/api_test.go`）。

### 创新点

**策略与 profile 解耦但就近定义**：策略文件在 profile 目录内（每个 profile 独立定义），但 gen-test-cases 通过统一的 capability key 查询策略，不需要硬编码类型映射。新增 profile 只需加策略文件，不需要改 skill 逻辑。

## Requirements Analysis

### Key Scenarios

1. **TUI feature 测试生成**：gen-test-cases 检测到 `tui` capability → 自动生成 golden file 测试用例 + 维度检查用例 + 边界场景用例，标记为 e2e 级别
2. **web-ui feature 测试生成**：gen-test-cases 检测到 `web-ui` capability → 生成 DOM 交互用例 + 视觉回归用例 + 响应式用例，标记为 e2e 级别
3. **API 集成测试生成**：gen-test-cases 检测到 `api` capability → 生成契约测试用例 + 错误路径用例 + 边界值用例，标记为 integration 级别
4. **CLI 集成测试生成**：gen-test-cases 检测到 `cli` capability → 生成输出 golden file 用例 + 退出码用例 + 参数组合用例，标记为 integration 级别
5. **Mixed profile**：go-test profile 同时有 `tui` + `api` + `cli` → 生成 3 类测试用例，TUI 标记 e2e，API/CLI 标记 integration
6. **Unknown capability type**：manifest.yaml 声明了 capability（如 `grpc`），但 verification-strategies.md 中无对应 `## grpc` section → gen-test-cases 输出 warning（`No strategy found for capability: grpc`），为该 capability 生成通用测试用例（无 Level 字段），继续处理其他已知 capability，不中断执行
7. **Capability key mismatch**：策略文件包含 `## tui` + `## api`，但 manifest.yaml 声明 `tui` + `cli` → gen-test-cases 中止并输出错误：`Strategy/manifest mismatch. Missing in strategy: cli. Extra in strategy: api. Please align verification-strategies.md sections with manifest.yaml capabilities.`
8. **Golden file staleness**：run-e2e-tests 执行 TUI golden file 比对时检测到输出不匹配 → 测试失败并输出 diff。agent review diff 后：(a) 若为 intentional behavior change → 运行 `gen-test-scripts --update-golden` 更新 golden file 并重新执行测试；(b) 若为 regression → 以 diff 作为证据 file bug 并修复代码。策略文件无需改动
9. **Strategy file parse failure**：verification-strategies.md 存在但格式不符合 section 结构要求（如缺少 `### 验证维度` 子标题）→ gen-test-cases 输出 warning（`Strategy file for profile {name} is malformed: missing ### 验证维度 in section ## tui`），回退到当前行为（无 Level 字段），不中断执行。开发者根据 warning 修复策略文件结构

### Constraints & Dependencies

- 依赖现有 profile 的 capability 声明（manifest.yaml）
- 策略文件是 profile 目录的一部分，随 profile 版本更新
- 不引入新的外部依赖
- **策略文件解析机制**：verification-strategies.md 是 Markdown 格式，gen-test-cases 通过 LLM prompt 上下文读取（非结构化 YAML 解析）。策略文件必须包含以下 section 结构以供 gen-test-cases 正确解读：每个 capability 以 `## <capability-key>` 为标题（如 `## tui`、`## api`），标题下按 `### 验证维度`、`### 边界场景`、`### 测试数据要求` 子标题组织内容。gen-test-cases 的 SKILL.md prompt 模板将包含策略文件的读取指令和 section 结构说明
- **Capability key 一致性**：gen-test-cases 读取策略文件后，比对策略文件中的 capability 标题（`## <key>`）与 manifest.yaml 中声明的 capability 列表。不一致时（策略文件缺少已声明 capability、或包含未声明 capability）gen-test-cases 中止并输出明确的 mismatch 错误信息（列出缺失/多余的 key）
- **策略文件有效性**：一个策略文件视为"有效"当且仅当：(1) 包含 ≥1 个 capability section（`## <capability-key>`），(2) 每个 capability section 包含 `### 验证维度` 子标题且其下 ≥3 个维度条目，(3) 每个 capability section 包含 `### 边界场景` 子标题且其下 ≥2 个场景条目。CI 中通过 lint 脚本自动检查，不满足条件的策略文件导致 profile lint 失败

### Non-Functional Requirements

| NFR Category | Requirement | Verification Method |
|-------------|-------------|-------------------|
| **Performance** | 策略文件读取增加延迟 < 2s（gen-test-cases 总耗时 < 5%）。策略文件上限 200 行，CI 检查 | 对比有/无策略文件时 gen-test-cases 执行耗时 |
| **Compatibility — 缺失策略文件** | 无 verification-strategies.md 时回退到当前行为（无 Level 字段），输出 warning | 验证：warning 出现、用例无 Level、不中断 |
| **Compatibility — 格式错误** | 解析失败时同"缺失"：warning + fallback | 构造错误格式文件验证 graceful degradation |
| **Security — 策略注入** | 策略内容仅作 gen-test-cases 参考输入，不拼入 shell 命令或 eval | 审计 gen-test-cases 确认无命令注入路径 |
| **Consistency — 策略漂移** | 跨 profile lint：同一 capability 验证维度差异 > 50% 时 warning | eval-harness 检查，不阻塞生成 |
| **Observability** | test-cases.md 头部注释：`Applied strategy: {profile}/{cap-count} capabilities, {dim-count} dimensions` | 检查头部注释包含策略元数据 |

## Alternatives & Industry Benchmarking

### Industry Patterns & Tools

**Testing Quadrant（Brian Marick / Martin Fowler）**：Q1-Q4 四象限模型。本提案的 e2e/integration 二级划分是其简化版（Q4+Q2 vs Q1+Q3）。forge 的约束是测试由 AI agent 生成，分级粒度不宜过细——agent 无法判断"Q2 vs Q4"，但能判断"UI 渲染 vs API 契约"。

**Pact（Contract Testing）**：consumer-driven contract 模式，需要 provider/consumer 双方运行验证服务器。具体对比：(1) 断言风格——Pact 使用 consumer-driven provider states（consumer 定义期望的请求/响应对，provider 验证这些状态是否满足）；forge 的 API 测试是单侧的，gen-test-cases 根据 verification-strategies.md 中的维度生成请求/响应匹配断言，不需要双方协调。(2) 输出格式——Pact 产出 JSON contract 文件（Pact specification），通过 Pact broker 分享和版本化；forge 产出的是可执行的测试脚本（如 `integration/api_test.go`），assertion 内联在测试代码中。(3) 生命周期——Pact 需要 broker 基础设施来管理 contract 版本和 provider verification；forge 的 `gen-test-scripts → run-e2e-tests` 链路不需要外部服务，策略文件即 contract。微服务场景下可引入 Pact 作为子策略，但当前 forge 的 API 测试场景（单项目验证）不需要 broker 的协作能力。

**Bats / Cram（CLI Golden Testing）**：`$ command` + `expected output` 模式。本提案直接借鉴此模式，差异化在于 golden 内容由 AI 自动生成而非人工编写。具体对比：(1) 输入格式——Bats 使用 `.t` 文件以 `$ command` / `expected output` 行对定义测试，Cram 使用同类 `.t` 格式；forge 生成的测试是 Go/Rust/Python 源文件，golden file 存放在 `testdata/*.golden` 中，内容为 ANSI-stripped 的实际输出。(2) 断言粒度——Bats/Cram 逐行比对 stdout 输出；forge 的 golden test 对整个 stdout 做 diff（full-output diff），但额外在策略中定义维度断言（如行数、宽度），这些维度断言在 Bats/Cram 中需要手动编写 `[[ ${#lines[@]} -le 24 ]]` 之类的检查。(3) 生命周期——Bats/Cram 测试由开发者手动编写和维护；forge 的 golden file 由 gen-test-scripts 自动生成、run-e2e-tests 执行比对，开发者仅在 mismatch 时介入。

**charmbracelet/vhs（TUI Testing）**：固定终端尺寸+输出比对。具体对比：(1) 测试定义方式——VHS 使用 `.tape` 脚本描述终端交互序列（`Output` / `Type` / `Wait` 命令），输出为 GIF 或文本截图；forge 的 TUI 测试不使用录制脚本，而是以固定 `--width`/`--height` 启动被测程序，通过子进程 stdout 截获 ANSI 输出与 golden file 比对。(2) 依赖——VHS 依赖 `ffmpeg` 和终端模拟器（`vhs` 本身启动 pty）；forge 的测试不需要终端模拟器，只依赖子进程 stdout 捕获，CI 环境兼容性更好。(3) 维度断言——VHS 主要做视觉输出比对（截图 diff）；forge 除了 golden file 全量 diff，还通过策略文件生成维度断言（行数 <= 高度、每行宽度 <= 终端宽度），这些是 VHS 不提供的结构性验证。

**Appium / Espresso（Mobile Testing）**：Appium 跨平台但慢，Espresso 仅 Android。forge 的 mobile-ui 策略不绑定工具——profile 层决定框架（如 Maestro），verification-strategies.md 只定义验证维度，与 Appium 互补而非竞争。

### Alternatives Analysis

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| **Do nothing** | 当前状态 | 零改动，无迁移成本 | TUI lesson 已证明不可靠：11 个渲染 bug 全部通过现有 verify gate（bug 捕获率 0/11） | Rejected：无法接受已验证的失败模式 |
| **Centralized config registry** | 类似 eslint-config-shared | 一致性高；更新一处全局生效 | 6 个 profile 框架差异大（Go 用 `testdata/`，Rust 用 `include_str!`，Python 用 `conftest.py` fixtures），集中基础策略无法表达框架特异的测试数据加载方式。混合方案（base strategy + per-profile override patches）看似折中，但在 forge 的 AI-generation 上下文中增加不必要的复杂度：base 层需要定义 override 机制（JSON patch？YAML merge key？），gen-test-cases 的 prompt 需要同时读取 base 和 override 并合并理解，token 开销和 prompt 复杂度双重增加——而 AI 直接读取完整的 profile-local 策略文件无需合并逻辑。一致性漂移通过跨 profile lint（见 NFR Consistency）检测即可，不需要中央注册来强制。 | Rejected：一致性收益被框架差异抵消，混合方案引入的 override 机制在 AI-generation 场景中无价值 |
| **Strategy hardcoded in SKILL.md** | 类似 Playwright 内置 annotation | 简单，无额外文件 | 新增 profile 需改 skill 代码；策略与 skill 版本强耦合 | Rejected：违背 profile 可插拔原则 |
| **Profile 策略文件** | 本提案 | 按 profile 框架定制验证手段；新增 profile 只需加文件，零改动 skill；策略可随 profile 版本独立演进 | 6 个 profile 各维护一份策略文件，存在一致性漂移风险（如 TUI 维度在 go-test 和 rust-test 中不同步） | **Selected**：一致性漂移可通过 lint 规则检测（见 NFR），且框架差异使得强制一致性不可取 |

## Feasibility Assessment

### Technical Feasibility

所有改动都在 skill 文件层面（SKILL.md、模板文件、profile 目录内的策略文件），不涉及编译代码。profile 已有 capability 声明，策略文件只是补充验证维度信息。

### Resource & Timeline

| Step | Task | 预计时间 |
|------|------|---------|
| 1 | 6 个 profile 各写 verification-strategies.md | 2h |
| 2 | test-cases.md 模板更新（Level 字段 + 类型 section） | 1h |
| 3 | gen-test-cases SKILL.md 增强策略读取和类型化生成 | 2h |
| 4 | gen-test-scripts SKILL.md 增强级别化代码生成 | 1.5h |
| 5 | eval-test-cases rubric 更新（类型化验证完整度维度） | 1h |
| 6 | 端到端验证 | 1h |
| **Total** | | **~8.5h** |

## Scope

### In Scope

- 6 个 profile 各新增 `verification-strategies.md`（go-test、web-playwright、maestro、pytest、rust-test、java-junit）
- gen-test-cases SKILL.md 增强：读取 profile 策略 → 按 interface 类型生成不同验证条件 → 自动标记 e2e/integration 级别
- gen-test-scripts SKILL.md 增强：按测试级别和 interface 类型选择代码生成策略
- test-cases.md 模板更新：新增 Level 字段、类型专属验证 section
- eval-test-cases rubric 更新：新增"类型化验证完整度"评分维度

### Out of Scope

- breakdown-tasks 改动（verify 模板注入）
- execute-task quality gate 改动
- eval-design 改动
- 新 profile 或新测试框架
- 视觉回归基础设施（截图对比服务）
- run-e2e-tests / graduate-tests 改动
- CI/CD 管道改动

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 策略文件定义的验证维度与实际 PRD 不匹配 | Medium | Medium — 生成的测试用例覆盖不全 | **Owner**: eval-test-cases rubric 自动检查。**Trigger**: gen-test-cases 生成后，eval-test-cases 对比 PRD 中 interface 行为描述与策略文件维度，维度覆盖率 < 70% 时标记为 incomplete 并输出缺失维度清单。**Corrective action**: 开发者根据清单补充策略文件维度后重新生成 |
| 6 个 profile 策略文件维护成本 | Medium | Low — 策略相对稳定 | **Owner**: Profile maintainer。**Trigger**: 每次 profile 发布时运行跨 profile lint（见 NFR Consistency 行）。**Corrective action**: 同一 capability 验证维度差异 > 50% 触发人工 review，maintainer 决定是同步维度还是保留框架差异并记录理由 |
| gen-test-cases 读取策略后 token 开销增加 | Medium | Low — 策略文件精简 | **Owner**: gen-test-cases skill。**Trigger**: gen-test-cases 执行完毕时输出 token 计数日志。**Corrective action**: 策略文件超过 200 行时 CI 报错阻塞；token 开销 > gen-test-cases 总消耗 15% 时 warning 并建议精简策略文件 |
| TUI golden file 测试在 CI 环境中的 terminal 模拟差异 | Medium | High — CI 中 golden test 不稳定 | **Owner**: gen-test-scripts skill。**Trigger**: TUI golden file 测试失败时。**Corrective action**: 策略文件强制要求 golden test 使用固定 terminal 尺寸（80x24）启动参数，gen-test-scripts 生成的测试代码硬编码 `--width 80 --height 24`，不读取环境变量 |

## Success Criteria

- [ ] 每个 verification-strategies.md 为每个 capability 定义 ≥3 个验证维度、≥2 个边界场景。自动验证：脚本扫描策略文件，计数规则为——每个 `## <capability-key>` 标题下，`### 验证维度` 子标题之后、下一个 `###` 或 `##` 标题之前的所有无序列表项（`- ` 开头）计为一个 dimension 条目；`### 边界场景` 子标题之后、下一个 `###` 或 `##` 标题之前的所有无序列表项计为一个 boundary 条目。低于阈值则报错。示例：`## tui` 下 `### 验证维度` 含 4 个 `- ` 列表项 = 4 dimensions（通过）；`## api` 下 `### 边界场景` 仅 1 个 `- ` 列表项 = 1 boundary（不通过）
- [ ] gen-test-cases 生成的 test-cases.md 中 Level 字段覆盖率 ≥ 95%（允许未命中策略文件的用例不带 Level，但命中策略的用例必须标记）。interface 类型与 manifest.yaml capability 声明一致率 100%（不一致则 gen-test-cases 中止并报错）
- [ ] TUI 类型用例必须包含 ≥1 个 golden file 断言 + ≥1 个维度检查（行数/宽度/对齐）+ ≥2 个边界场景用例（如 CJK、窄终端）。覆盖检查通过 eval-test-cases rubric 自动评分
- [ ] API 类型用例必须包含 ≥1 个契约验证断言（请求/响应匹配 spec）+ ≥1 个错误路径用例 + ≥2 个边界值用例（无效输入、超限）
- [ ] gen-test-scripts 为 e2e 和 integration 级别生成结构不同的测试代码：e2e 生成渲染截获 + golden file 比对函数；integration 生成 HTTP 断言或子进程退出码检查函数。差异体现在 import 列表、assertion 库选择、测试目录结构三方面，且三方面均必须至少存在一处差异（例如 import 中 e2e 引用 `os/exec` + `golden` 而 integration 引用 `net/http` + `assert`）。自动验证：逐项对比 e2e 与 integration 生成输出的 import/structure/assertion 差异，任一方面无差异则失败
- [ ] eval-test-cases rubric 新增"类型化验证完整度"维度，权重 ≥ 15%（总 rubric 1000 分中 ≥ 150 分），通过阈值为该维度得分 ≥ 70%

## Next Steps

- Proceed to `/write-prd` to formalize requirements
