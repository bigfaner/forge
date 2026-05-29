# Report 04: Skills Deep Audit - Batch C

**Baseline commit**: `1542c8cc`
**Date**: 2026-05-30
**Scope**: 7 skills (extract-design-md, gen-contracts, gen-journeys, gen-sitemap, init-justfile, submit-task, test-guide)
**Layers**: Layer 2 (Instruction Consistency) + Layer 3 (Timing & Flow)

---

## 1. extract-design-md

**Files audited**: SKILL.md + 3 rules + 3 templates = 7 files

### Layer 2: Instruction Consistency

| # | File | Layer | Category | Severity | Description | Fix Suggestion | Confidence |
|---|------|-------|----------|----------|-------------|----------------|------------|
| C-01 | SKILL.md vs rules/match-strategy.md | 2 | REDUNDANT | P2 | SKILL.md Step 3 的匹配策略表格与 rules/match-strategy.md 重复。两者都列出了相同的两个选项（"Match closest built-in style" 和 "Fully custom"）以及相同的 5 个 built-in style 特征表。rules/match-strategy.md 未添加 SKILL.md 中缺少的信息增量。 | rules/match-strategy.md 应包含匹配算法的详细步骤（如何计算相似度、权重分配），而非仅重复选项列表。或 SKILL.md 移除详细特征表，仅引用 rules。 | high |
| C-02 | SKILL.md vs rules/extraction-layers.md | 2 | REDUNDANT | P2 | SKILL.md Step 2 的"Extraction Strategy"列表（Layer 1-5）与 rules/extraction-layers.md 内容高度重复。SKILL.md 列出了 5 个 layer 的简述，rules 文件展开描述同样的 5 个 layer。rules 文件增加了具体命令示例（agent-browser eval 命令），这构成有效信息增量。 | 部分冗余。rules 中的 agent-browser 命令示例是有效增量，但 layer 名称和基本描述重复。建议 SKILL.md 仅保留 layer 名称列表，将详细描述和命令完全委托给 rules 文件。 | medium |
| C-03 | SKILL.md vs rules/platform-routing.md | 2 | REDUNDANT | P2 | SKILL.md 的 "Mobile Extraction Summary" 和 "TUI Extraction Summary" 两节与 rules/platform-routing.md 高度重复。两者都描述了相同的 mobile breakpoint 分析、touch target 估算、safe area 检查、TUI ANSI color 分析等。rules 文件增加了具体实现细节（如 agent-browser Read tool 用法），构成部分信息增量。 | SKILL.md 的 "Mobile Extraction Summary" 和 "TUI Extraction Summary" 应缩减为一行引用，如"Follow mobile extraction rules in `rules/platform-routing.md`"。当前 SKILL.md 包含了过多实现细节。 | medium |
| C-04 | templates/design-web.md vs templates/design-mobile.md | 2 | INCOMPLETE | P2 | SKILL.md Step 5 指定 mobile 平台使用 `templates/design-web.md` 并追加 `templates/design-mobile.md`。但 design-web.md 模板缺少 SKILL.md "Target dimensions" 表中列出的 "Design philosophy" 维度的对应 section。design-web.md 有 "Visual Theme & Atmosphere" 节但不等同于 "Design philosophy"（后者要求 specific style keywords）。 | 在 design-web.md 模板中添加 `## Design Philosophy` section，或在 "Visual Theme & Atmosphere" 中增加 style keywords 占位符。 | low |
| C-05 | rules/platform-routing.md vs SKILL.md | 2 | CONFLICT | P1 | rules/platform-routing.md 第 2 节"TUI Extraction"第 4 步描述了独立的 match strategy（让用户选择 built-in TUI theme vs fully custom），但 SKILL.md Step 3 的 match strategy 仅提到 web built-in styles（vercel/shadcn/tailwind-ui/stripe/apple），未提到 TUI themes（modern-dark-tui/minimal-ascii-tui）。SKILL.md Step 3 说"Use AskUserQuestion to let the user choose"，而 rules/platform-routing.md 为 TUI 平台提供了独立的 match strategy。 | SKILL.md Step 3 应明确区分：web/mobile 平台使用 SKILL.md 定义的 match strategy（5 个 web built-in styles），TUI 平台使用 rules/platform-routing.md 定义的 match strategy（2 个 TUI themes）。当前 TUI match strategy 仅隐藏在 rules 文件中，SKILL.md 未提及。 | high |

### Layer 3: Timing & Flow

| # | File | Layer | Category | Severity | Description | Fix Suggestion | Confidence |
|---|------|-------|----------|----------|-------------|----------------|------------|
| T-01 | SKILL.md vs rules/platform-routing.md | 3 | TIMING | P2 | SKILL.md Process Flow 为 "1. Parse platform -> 2. Validate -> 3. Extract -> 4. Match -> 5. Build -> 6. Write -> 7. Confirm"，但 rules/platform-routing.md 为 TUI 平台在 extraction（step 3）之后、match（step 4）之前增加了 screenshot quality check 和 AI vision analysis 两个子步骤。这些子步骤的时序在 SKILL.md 的 "TUI Extraction Summary" 中被隐含提及（"Built-in TUI themes: Match against modern-dark-tui / minimal-ascii-tui"），但 "Screenshot quality check: Reject blurry" 出现在 SKILL.md Error Handling 表中而非流程步骤中。 | 将 screenshot quality check 明确作为 TUI extraction 的前置子步骤编入 SKILL.md 流程，而非仅放在 Error Handling。 | low |

---

## 2. gen-contracts

**Files audited**: SKILL.md + 6 rules + 2 templates = 9 files

### Layer 2: Instruction Consistency

| # | File | Layer | Category | Severity | Description | Fix Suggestion | Confidence |
|---|------|-------|----------|----------|-------------|----------------|------------|
| C-06 | SKILL.md Step 4 vs rules/validation.md | 2 | REDUNDANT | P2 | SKILL.md Step 4 "Validate Contracts" 的 schema validation checks 表格与 rules/validation.md 的 validation checks 表格几乎完全相同。两者都列出了 6 个相同的检查项（Mandatory dimensions, Semantic descriptor purity, Outcome name uniqueness, Preconditions mutual exclusivity, Journey Invariants, Side-effect default）。rules/validation.md 增加了 2 项 SKILL.md 未列出的检查："Outcome count checkpoint" 和 "Unclassified validation points"。 | SKILL.md Step 4 应仅引用 rules/validation.md 并列出差异（如 retry logic），而非重复整个 validation checks 表格。当前重复增加了维护风险——更新一处时可能遗漏另一处。 | high |
| C-07 | SKILL.md vs rules/dimension-rules.md | 2 | REDUNDANT | P2 | SKILL.md Section 3.2 "Six-Dimension Declaration" 和 3.3 "Semantic Descriptors" 与 rules/dimension-rules.md 的前两个 section 高度重复。两者都描述了 4 mandatory + 2 optional 维度、semantic descriptor 规则、MUST NOT regex 语法。rules/dimension-rules.md 增加了 good/bad examples，构成有效增量。 | 同 C-06。SKILL.md 应缩减为引用 + 差异描述。 | medium |
| C-08 | SKILL.md Section 3.5 vs rules/risk-density.md | 2 | REDUNDANT | P2 | SKILL.md Section 3.5 的 risk density targets 表格（High: 3-5/13-20, Medium: 2-3/8-12, Low: 1-2/4-7）与 rules/risk-density.md 的 Density Targets 表格完全相同。rules/risk-density.md 还包含 density 应用步骤和 surface-required outcome derivation 规则，构成有效增量。 | SKILL.md Section 3.5 应仅引用 `rules/risk-density.md` 并保留 Outcome generation priority 列表（这是 rules 文件中没有的）。 | high |
| C-09 | SKILL.md vs rules/tui-async.md | 2 | INCOMPLETE | P2 | SKILL.md Sections 3.7-3.10 将 TUI async await semantics、state verification levels、journey-level invariants、batch processing 全部委托给 `rules/tui-async.md`。但 `rules/tui-async.md` 同时包含了 State Verification Levels 和 Batch Processing 的规则。SKILL.md 仅各有一行引用，但用词是"per rules/tui-async.md"，信息流向清晰。这是一个**合理的设计模式**（SKILL.md 概述 + rules 详细规则），标记为 INCOMPLETE 因为 SKILL.md 未提及 rules/tui-async.md 中的 `tea.Batch(cmd1, cmd2)` 并发语义。 | 在 SKILL.md Section 3.7 中添加一行说明：TUI Steps with concurrent Cmds (tea.Batch) must declare individual timeout for each Cmd。当前此规则仅存在于 rules 文件中。 | low |
| C-10 | SKILL.md vs rules/code-reconnaissance.md | 2 | CONFLICT | P1 | SKILL.md Step 2 描述 Fact Table 输出格式使用 JSON（`"source": "static"`, `"confidence": "inferred"`, `"updated_at"` 等字段），但 rules/code-reconnaissance.md 描述 Fact Table 输出格式使用 Markdown 表格（`| Key | Value | Source |`）。两种格式不一致。SKILL.md 还提到 Fact Table canonical schema 定义在 `forge-cli/pkg/facttable/facttable.go`，这暗示 JSON 是权威格式，但 rules 文件的 Markdown 格式可能导致混淆。 | 更新 rules/code-reconnaissance.md 使用与 SKILL.md 一致的 JSON Fact Table 格式，或在 rules 文件中明确说明 Markdown 是中间展示格式，最终写入 `.forge/fact-table.json` 时转换为 JSON。 | high |
| C-11 | SKILL.md Section 3.9 vs rules/tui-async.md "Journey-Level Invariants" | 2 | REDUNDANT | P3 | SKILL.md Section 3.9 "Journey-Level Invariants" 说 "Every Contract file MUST end with a `## Journey Invariants` section per `rules/tui-async.md`"。但 Journey Invariants 并非仅与 TUI 相关——它是所有 surface type 的通用要求。将其委托给 `rules/tui-async.md`（文件名暗示 TUI-only）会让非 TUI 项目忽略此规则。 | 将 Journey-Level Invariants 规则移至 rules/dimension-rules.md 或 rules/validation.md，使其不被误认为 TUI-only。 | medium |
| C-12 | templates/contract.md vs SKILL.md | 2 | INCOMPLETE | P3 | SKILL.md 提到 Contract 文件 "MUST include a `## Journey Invariants` section"，templates/contract.md 也包含 `## Journey Invariants` 节。但 SKILL.md 提到的 `skip_eval: true` frontmatter 字段（SKIP_EVAL_GATE 模式下）未在 contract.md 模板中体现。 | 在 templates/contract.md 的 frontmatter 中添加可选的 `skip_eval` 字段占位符和说明注释。 | low |
| C-13 | templates/outcome-block.md vs SKILL.md | 2 | INCOMPLETE | P3 | SKILL.md 提到 Outcome 可以包含 `Invariants` 维度（step-level），且 rules/dimension-rules.md 定义其为可选维度。templates/outcome-block.md 包含 `{{STEP_INVARIANTS}}` 占位符，但当省略时应默认无约束。模板中 `Side-effect` 使用硬编码的 `{{SIDE_EFFECT}}`（无默认值），但 rules/dimension-rules.md 说 "When Side-effect is omitted or empty, defaults to `none`"。模板未体现此默认值行为。 | 在 outcome-block.md 模板中为 Side-effect 添加默认值说明：`Side-effect: {{SIDE_EFFECT}} <!-- defaults to "none" when omitted -->`。 | low |

### Layer 3: Timing & Flow

| # | File | Layer | Category | Severity | Description | Fix Suggestion | Confidence |
|---|------|-------|----------|----------|-------------|----------------|------------|
| T-02 | SKILL.md Process Flow | 3 | TIMING | P2 | SKILL.md Process Flow 定义为 "0. Resolve language -> 1. Read Journeys -> 2. Code Reconnaissance -> 3. Generate Contracts -> 4. Validate -> 5. Write"。Step 0 "Resolve language and surfaces" 包含 "Detect surfaces: Check .forge/config.yaml surfaces field"。但 Step 1 第 4 子步要求 "Load the project's surface type from .forge/config.yaml and read the corresponding surface rule from gen-journeys skill's rules/surface-<type>.md"。这意味着 surface 检测在 Step 0 和 Step 1 各执行一次——Step 0 检测 surface 用于 language resolution，Step 1 重新加载 surface rules 用于 required_outcomes。这不是时序错误但存在重复执行。 | 在 Step 0 完成后将检测到的 surface types 缓存，Step 1 直接使用缓存结果而非重新检测。或明确说明 Step 0 的 surface 检测结果是 Step 1 的前置输入。 | low |

---

## 3. gen-journeys

**Files audited**: SKILL.md + 1 cross-skill reference + 5 surface rules + 1 template = 8 files

### Layer 2: Instruction Consistency

| # | File | Layer | Category | Severity | Description | Fix Suggestion | Confidence |
|---|------|-------|----------|----------|-------------|----------------|------------|
| C-14 | SKILL.md "Surface Detection" vs surface rules | 2 | REDUNDANT | P2 | SKILL.md 的 "Per-Surface Rule Application" 节为每个 surface type 列出了 Mandatory outcomes、Test level emphasis、Edge case focus。这些信息是对 surface rules 文件内容的摘要。例如 SKILL.md 说 API surface 的 "Mandatory error outcomes: HTTP status code boundaries (4xx, 5xx)"，而 rules/surface-api.md 的 "Required Outcome Reference" 节详细定义了 `unauthorized` outcome。SKILL.md 的摘要与 rules 文件不完全对齐：SKILL.md 说 API "integration-heavy ratio"，但 rules/surface-api.md 说 "Balanced 50/50 (Contract 50% / Journey smoke 50%)"。 | 这是 CONFLICT 而非 REDUNDANT。SKILL.md "Per-Surface Rule Application" 节中 API surface 的 test level emphasis 描述（"integration-heavy ratio"）与 rules/surface-api.md 的 "Balanced 50/50" 矛盾。WILL FOLLOW SPEC: 以 rules/surface-api.md 为准。 | high |
| C-15 | SKILL.md "Surface Detection" vs surface rules (upgrade from C-14) | 2 | CONFLICT | P1 | SKILL.md "Per-Surface Rule Application" 节中：API surface 说 "integration-heavy ratio"；Web surface 说 "e2e-heavy ratio"；TUI surface 说 "integration-heavy ratio"。但对应 rules 文件中：API (rules/surface-api.md) 说 "Balanced 50/50"；Web (rules/surface-web.md) 说 "Balanced 50/50"；TUI (rules/surface-tui.md) 说 "Contract 80% / Journey smoke 20%"。三处全部不一致。WILL FOLLOW SPEC: 以 rules 文件为准。 | 更新 SKILL.md "Per-Surface Rule Application" 节中所有 surface type 的 test level emphasis 描述，与对应 rules 文件保持一致。或移除 SKILL.md 中的摘要，仅引用 rules 文件。 | high |
| C-16 | SKILL.md vs templates/journey.md | 2 | INCOMPLETE | P3 | SKILL.md Step 4 要求 Journey 文件 frontmatter 包含 `quality: low` 字段（Proposal Mode smoke-level Journeys），但 templates/journey.md 的 frontmatter 不包含 `quality` 字段占位符。 | 在 templates/journey.md frontmatter 中添加可选的 `quality` 字段：`quality: "{{QUALITY_LEVEL}}" # Optional: "low" for smoke-level Journeys`。 | medium |
| C-17 | SKILL.md "Proposal Mode" section | 2 | INCOMPLETE | P2 | SKILL.md 的 Proposal Mode 描述了一个 `quality: low` frontmatter 注释和 Quality Notice 警告文本。这些内容是 SKILL.md 内联定义的，没有对应的 rules 文件或 template 来标准化格式。Quality Notice 的格式是在 SKILL.md 中用 markdown blockquote 给出的，但 templates/journey.md 不包含此 blockquote 占位符。 | 在 templates/journey.md 中添加 Quality Notice blockquote 占位符，或创建一个 rules/proposal-mode.md 来集中定义 Proposal Mode 的特殊处理规则。 | low |

### Layer 3: Timing & Flow

No TIMING issues found. The process flow (Read Input -> Identify Workflows -> Classify Risk -> Generate Files -> Validate -> Review & Commit) has clear sequential steps with no circular dependencies. Step 5 (Validate) correctly depends on Step 4 (Generate), and Step 6 (Review) correctly depends on Step 5.

---

## 4. gen-sitemap

**Files audited**: SKILL.md + 3 rules + 1 template = 5 files

### Layer 2: Instruction Consistency

| # | File | Layer | Category | Severity | Description | Fix Suggestion | Confidence |
|---|------|-------|----------|----------|-------------|----------------|------------|
| C-18 | SKILL.md vs rules/merge-validation.md | 2 | REDUNDANT | P2 | SKILL.md Step 5 "Merge, Dedup & Validate" 仅引用 rules/merge-validation.md（"See rules/merge-validation.md for the complete rules on: 5a. Element merge, 5b. Post-collection dedup, 5c. Stale route detection, 5d. Validation"），未在 SKILL.md 中重复具体规则。这是一个良好的设计模式——SKILL.md 仅列出子步骤名称，完整规则委托给 rules 文件。无需修改。 | 无需修改（positive example）。 | high |
| C-19 | SKILL.md vs rules/schema.md | 2 | REDUNDANT | P3 | SKILL.md Schema 节说 "See rules/schema.md for the complete field reference"，SKILL.md 未重复 schema 细节。这是良好的委托模式。但 SKILL.md Config Resolution 节中内联定义了 `baseUrl` priority 和 `apiBaseUrl` 路径前缀规则，这些规则也可以提取到 rules 文件中。当前的分散是可接受的，因为这些是 SKILL.md 层面的流程逻辑而非 schema 约束。 | 无需修改。 | high |
| C-20 | SKILL.md Step 4 vs rules/page-exploration.md | 2 | INCOMPLETE | P2 | SKILL.md Step 4 的 "Layout Element Filtering" 说使用 `role + name` 匹配来过滤 layout elements。rules/page-exploration.md 的 Element Extraction 也说 "Exclude elements already in layout.elements (matched by role + name)"。两者一致。但 SKILL.md Step 4 的 HARD-RULE 说 "MUST skip layout elements"，而 rules/page-exploration.md 中的 Dynamic State Exploration 未提到 layout filtering——动态状态中的新元素是否也应排除 layout elements 未明确说明。 | 在 rules/page-exploration.md 的 Dynamic State Exploration 节中添加说明：动态状态中出现的 layout-level 元素（如 modal overlay 中包含 nav elements）也应被过滤。 | low |
| C-21 | SKILL.md "Prerequisites" vs templates/test-config.yaml | 2 | INCOMPLETE | P3 | SKILL.md Prerequisites 节提到 "agent-browser is installed (optional tool)"，但 templates/test-config.yaml 不包含 agent-browser 相关配置。Config Resolution 节定义了 `baseUrl`、`apiBaseUrl`、`username`、`password` 字段，templates/test-config.yaml 也包含这些字段和 `loginLocators`、`timeout` 字段。SKILL.md 未引用 `loginLocators` 字段但 Config Resolution 节使用了 `username`/`password` 的登录流程。template 的 `loginLocators` 字段在 SKILL.md 中无对应使用说明。 | 在 SKILL.md Config Resolution 节中说明 `loginLocators` 的用途（覆盖默认的 regex-based locator detection），或删除 template 中的该字段。 | low |

### Layer 3: Timing & Flow

| # | File | Layer | Category | Severity | Description | Fix Suggestion | Confidence |
|---|------|-------|----------|----------|-------------|----------------|------------|
| T-03 | SKILL.md Process Flow | 3 | TIMING | P1 | SKILL.md Process Flow 定义为 "1. Load existing -> 2. Analyze layout & build route registry -> 3. Discover routes -> 4. Explore pages -> 5. Merge & dedup -> 6. Write"。Step 2 包含 2a (route registry from code) 和 2b (identify shared elements via multi-page comparison)。但 Step 2b 使用 agent-browser 访问 3-5 个 wrapped routes，这意味着 Step 2b 实际上已经开始了页面探索（与 Step 4 重叠）。Step 4 "Explore Pages" 对每条路由逐个探索。如果 Step 2b 探索了 5 条路由，Step 4 是否应跳过这 5 条？SKILL.md 未说明 Step 2b 探索过的页面在 Step 4 中的处理方式。 | 在 Step 4 中明确说明：Step 2b 中已探索的页面，Step 4 仅需提取 page-specific elements（跳过 layout comparison），或直接复用 Step 2b 的 snapshot 结果。 | high |
| T-04 | SKILL.md Step 3 vs Step 2a | 3 | TIMING | P2 | Step 3 "Discover Routes" 使用 agent-browser crawl links 来发现路由。但 Step 2a 已经通过分析代码构建了完整的 route registry。Step 3 应以 Step 2a 的结果作为 BASE，然后通过 crawling 补充。SKILL.md 确实说"Start with route registry from Step 2a as the BASE"，时序正确。但 Step 3 的 crawling 依赖 agent-browser，而 Step 2b 也依赖 agent-browser。如果 agent-browser 不可用（SKILL.md Prerequisites 中说明的降级场景），Step 2b 和 Step 3 都应跳过。SKILL.md 说 "skip Steps 3-4 link crawling"，但 Step 2b 也需要 agent-browser，未明确说明 Step 2b 的降级行为。 | 在 SKILL.md Prerequisites 的 agent-browser 降级说明中明确：无 agent-browser 时跳过 Step 2b（layout identification）和 Step 3-4（route discovery/page exploration），仅依赖 Step 2a 的静态分析结果。 | medium |

---

## 5. init-justfile

**Files audited**: SKILL.md + 2 rules + 5 surface rules + 6 templates = 14 files

### Layer 2: Instruction Consistency

| # | File | Layer | Category | Severity | Description | Fix Suggestion | Confidence |
|---|------|-------|----------|----------|-------------|----------------|------------|
| C-22 | SKILL.md Step 0 vs templates/*.just | 2 | CONFLICT | P1 | SKILL.md Step 0 HARD-RULE 明确说 "Do NOT use framework-specific recipe templates. Generate unit-test and surface-level test recipes from Convention content and LLM knowledge of the framework." 但 templates/ 目录下存在 6 个 .just 文件（generic.just, go.just, node.just, python.just, rust.just, mixed.just），这些是按语言/项目类型分类的框架特定模板。SKILL.md 的 "generate from LLM knowledge" 指令与 .just 模板的存在直接矛盾。此外，go.just 模板中的 `test` recipe 硬编码了 `npx playwright test`（Node.js/Playwright 命令），这是 Go 项目模板中的遗留错误。 | (1) 确认 .just 模板是否仍为活跃组件。如果是：更新 SKILL.md Step 0 和 Step 3a 以引用模板。如果否：删除 .just 模板。(2) go.just 中的 `test` recipe 不应包含 `npx playwright test`——Go 项目应使用 `go test`。 | high |
| C-23 | SKILL.md Step 3a vs templates/go.just | 2 | CONFLICT | P1 | SKILL.md Step 3a 说 "Generate unit-test, compile, build, lint, fmt, check, clean, install, ci recipes from Convention knowledge and LLM understanding"。但 templates/go.just 提供了硬编码的 Go 特定命令（如 `go vet ./...`、`go test ./...`、`golangci-lint run ./...`）。如果 LLM 被指示不使用模板，那这些硬编码命令不应被使用；如果使用模板，则违反 SKILL.md 的 HARD-RULE。 | 同 C-22。需要决定模板的角色：是活跃模板还是遗留文件。 | high |
| C-24 | templates/go.just vs templates/node.just | 2 | CONFLICT | P1 | go.just 和 node.just 的 `test` recipe 都硬编码了 `npx playwright test`（第 65-69 行和整个 test recipe body）。Go 项目的 `test` recipe 不应使用 Node.js 的 Playwright 命令。这表明 go.just 的 test recipe 是从 node.just 复制而来，未正确适配 Go 项目。这是 go.just 模板的一个 P0 级 bug（会导致运行时错误）。 | go.just 的 `test` recipe 应使用 Convention-defined 的 Go test 命令（如 `go test ./tests/... -v -tags=web-e2e`），而非 `npx playwright test`。 | high |
| C-25 | SKILL.md "Standard Target Contract" vs templates/go.just | 2 | INCOMPLETE | P2 | SKILL.md "Standard Target Contract" 定义了 `unit-test`（required）、`compile`、`build`、`lint`、`fmt`、`check`、`clean`、`install`、`ci` 作为 language-level targets。templates/go.just 包含这些加上额外的 `run`、`dev`、`test`、`test-setup`、`probe` targets。SKILL.md 未在 Standard Target Contract 中列出 `run`、`dev`、`test`、`test-setup`、`probe` 作为 language-level targets，但模板中它们存在。 | 在 SKILL.md Standard Target Contract 中明确说明 `run`、`dev`、`test`、`test-setup`、`probe` 是 additional convenience targets（不属于 standard contract 但存在于模板中），或将它们从模板中移除。 | medium |
| C-26 | SKILL.md "Surface-Level Targets" vs surface rules | 2 | CONFLICT | P1 | SKILL.md "Surface-Level Targets" 表说 mobile surface 有 `<key>-dev` 和 `<key>-probe` targets，但 init-justfile 的 surface rules/surfaces/mobile.md 定义了 mobile surface 的编排序列为 `test-setup -> dev -> probe -> test -> teardown`，多了一个 `test-setup` 步骤。SKILL.md 的 Surface-Level Targets 表未列出 `test-setup` target。此外，SKILL.md 说 CLI/TUI surfaces "do NOT generate dev, probe, or aggregate recipes"，但 CLI/TUI 不在 Surface-Level Targets 表中列出——一致性上没问题。 | 在 SKILL.md Surface-Level Targets 表中添加 mobile-specific 的 `<key>-test-setup` target 行，注明 "Mobile only: prepare emulator and test environment"。 | high |
| C-27 | SKILL.md Step 3b vs surface rules/surfaces/*.md | 2 | REDUNDANT | P3 | SKILL.md Step 3b "Generate surface-level recipes" 描述了 recipe 命名规则（single surface 用 type prefix，mixed 用 key prefix）和 dual-platform variants。Surface rules 文件也包含了完整的 recipe templates（含 `[linux]`/`[windows]` 属性和 `# user-customized` 标记）。SKILL.md 的描述与 surface rules 的模板一致，但 SKILL.md 没有提到 surface rules 中的 recipe templates——它说 LLM 应从 Convention knowledge 和 surface rule file 的 orchestration sequence 生成命令。 | 这是一个设计冲突（类似 C-22）：SKILL.md 说 LLM 生成，但 surface rules 提供了完整模板。应统一：要么 SKILL.md 引用 surface rules 中的模板作为起点，要么 surface rules 不提供具体模板代码。 | medium |
| C-28 | SKILL.md Step 3b vs rules/surfaces/api.md, web.md | 2 | INCOMPLETE | P2 | SKILL.md Step 3b 说 "For each surface, read the surface rule file... extract Orchestration sequence, Recipe contracts, Journey filter strategy"。但 rules/surfaces/api.md 和 web.md 的 recipe templates 中使用的是硬编码的 surface type 作为 prefix（如 `api-dev`、`web-test`），而 SKILL.md Step 3b 说 mixed project 应使用 `<surface-key>-<verb>` 命名（如 `admin-panel-dev`）。Surface rules 文件未提及 `<surface-key>` 命名规则——它们只展示了 `<type>-<verb>` 模式。 | 在 surface rules 文件的 LLM 指令中添加命名规则说明：single surface project 使用 type 作为 prefix，mixed project 使用 key 作为 prefix。 | medium |

### Layer 3: Timing & Flow

| # | File | Layer | Category | Severity | Description | Fix Suggestion | Confidence |
|---|------|-------|----------|----------|-------------|----------------|------------|
| T-05 | SKILL.md Process Flow | 3 | TIMING | P2 | SKILL.md Process Flow 为 "0. Load Convention -> 1. Detect project type -> 1s. Detect surfaces -> 2. Check existing justfile -> 3. Generate recipes -> 4. Verify -> 5. Output"。Step 0 (Load Convention) 依赖 `docs/conventions/` 目录中的文件。但如果 Convention 文件不存在（cold start），Step 0 降级为 "proceed to Step 1 for file signal detection"。Step 1 检测语言/框架，这提供的信息本应反馈到 Convention 加载——但 Step 0 已在 Step 1 之前执行。这意味着 cold start 路径中，Convention 加载永远失败（因为 Step 1 尚未检测语言），导致 LLM 使用通用默认值。 | 这是一个轻微的时序问题。在 cold start 场景中，Step 1 的检测结果应在 Step 0 的 Convention fallback 中被考虑。建议：当 Step 0 未找到 Convention 时，先执行 Step 1 的语言检测，再尝试基于检测到的语言重新查找 Convention。 | low |

---

## 6. submit-task

**Files audited**: SKILL.md + 6 data files = 7 files

### Layer 2: Instruction Consistency

| # | File | Layer | Category | Severity | Description | Fix Suggestion | Confidence |
|---|------|-------|----------|----------|-------------|----------------|------------|
| C-29 | SKILL.md "Common Fields" vs data/record-format-coding.md | 2 | REDUNDANT | P3 | SKILL.md "Common Fields" 表定义了 `taskId`, `status`, `summary`, `filesCreated`, `filesModified`, `acceptanceCriteria`, `notes`。data/record-format-coding.md 也列出了 `taskId`, `status`, `summary`, `filesCreated`, `filesModified`, `acceptanceCriteria`, `notes`（作为 JSON example 中的字段）。两者一致，无矛盾。但 SKILL.md 重复了每个 data 文件中的 Common Fields 部分——这是一个轻微的维护风险。 | 可接受的设计：SKILL.md 定义 Common Fields 作为规范，data 文件在 JSON example 中体现它们。保持同步是维护负担但增加了每个文件的独立性。 | high |
| C-30 | SKILL.md "Common Rules" vs data/record-format-coding.md | 2 | CONFLICT | P1 | SKILL.md "Common Rules" 说 "`acceptanceCriteria` with any `met: false` entry is rejected for `completed` status — use `blocked` instead"。但 data/record-format-coding.md 的 Rules 节说 "`testsFailed > 0` with `completed` triggers auto-downgrade to `blocked`"。这两条规则功能类似但触发条件不同：SKILL.md 的规则基于 `acceptanceCriteria`，data 文件的规则基于 `testsFailed`。对于 coding tasks，如果 `acceptanceCriteria` 全部 `met: true` 但 `testsFailed > 0`，应该用哪条规则？SKILL.md 说 `completed` 可以（因为 AC 全 pass），data 文件说必须 `blocked`（因为 testsFailed > 0）。 | 在 SKILL.md Common Rules 中明确说明：coding tasks 还必须满足 category-specific rules（如 data/record-format-coding.md 中的 testsFailed 检查），category-specific rules 与 common rules 冲突时以 category-specific 为准。 | high |
| C-31 | SKILL.md "Common Rules" vs data/record-format-coding.md | 2 | CONFLICT | P2 | SKILL.md "Common Rules" 说 "`testsPassed`, `testsFailed`, `coverage` are coding-only. Do NOT include them in doc, test, validation, or gate records." 但 data/record-format-coding.md 在 Required 列中将 `testsPassed`、`testsFailed`、`coverage` 标记为 "conditional"（required when status: completed），而 SKILL.md 说它们是 coding-only 但未说明它们在 coding records 中是 conditional。 | 在 SKILL.md Common Rules 中将 "`testsPassed`, `testsFailed`, `coverage` are coding-only" 补充为 "coding-only, conditional on status: completed"。 | medium |
| C-32 | data/record-format-gate.md vs SKILL.md | 2 | INCOMPLETE | P3 | data/record-format-gate.md Rules 节说 "The gate MUST verify all project code compiles before passing" 和 "Do NOT include `acceptanceCriteria` unless the gate explicitly checks against criteria"。但 SKILL.md Common Fields 表将 `acceptanceCriteria` 标记为 "recommended"（对全部 category），未说明 gate category 的例外。 | 在 SKILL.md Common Fields 表的 `acceptanceCriteria` 行添加注释："For gate tasks: only include if the gate explicitly checks against criteria"。 | low |

### Layer 3: Timing & Flow

No TIMING issues found. The submit-task flow (Determine Record Format -> Write record.json -> Submit via CLI) is simple and linear with no timing dependencies between steps.

---

## 7. test-guide

**Files audited**: SKILL.md + 4 rules + 1 template = 6 files

### Layer 2: Instruction Consistency

| # | File | Layer | Category | Severity | Description | Fix Suggestion | Confidence |
|---|------|-------|----------|----------|-------------|----------------|------------|
| C-33 | SKILL.md Step 3 vs rules/draft-generation.md | 2 | REDUNDANT | P2 | SKILL.md Step 3 "Generate Convention Draft" 的子步骤（3a Select template source、3b Customize with extracted patterns、3c Validate draft completeness）与 rules/draft-generation.md 的内容高度重复。两者都列出了：built-in Convention template 映射表（相同的 6 个 framework）、Extracted Pattern -> Convention Field Override 映射表（相同的 5 个映射）、validation rules（相同的 6 项检查）。 | 同 C-06 模式。SKILL.md 应仅引用 rules/draft-generation.md 并保留 SKILL.md 独有的内容（如 Step 4 的 user review 交互流程）。 | high |
| C-34 | SKILL.md Step 3c vs rules/convention-structure.md | 2 | CONFLICT | P2 | SKILL.md Step 3c "Validate draft completeness" 列出的 validation checks 为：framework (name, language, runner_command)、discovery (file_pattern)、structure (suite_pattern)、assertions (style, library)、Tags (tag format)、Result Format (execution command)。但 rules/convention-structure.md 的 Required sections 定义为：Framework (name, file pattern, package, test runner, build tag)、Assertion (library, key functions, rule)、Tags (build tag syntax)、Result Format (output flags, format type, execution command)。两者 section 名称不同：SKILL.md 用 "discovery" 和 "structure"，convention-structure.md 用 "Framework" 和无 discovery/structure 分节。templates/convention-template.md 也使用 convention-structure.md 的 section 名称（Framework, Assertion, Tags, Result Format），不包含 "discovery" 或 "structure" section。 | 更新 SKILL.md Step 3c 的 validation checks 使用与 convention-structure.md 一致的 section 名称。或更新 rules/draft-generation.md 的 schema validation 表使用统一的 section 名称。当前的不一致可能导致生成的 Convention 文件与 SKILL.md 的 validation 不匹配。 | medium |
| C-35 | SKILL.md Step 1 vs rules/signal-detection.md | 2 | REDUNDANT | P3 | SKILL.md Step 1 的子步骤（1a-1d）描述了完整的检测流程。rules/signal-detection.md 也描述了相同的检测流程（Language Detection -> Framework Detection -> Cross-validation -> Confidence Levels）。SKILL.md 的子步骤与 rules 文件的节对应，但 SKILL.md 增加了 Step 1 的 "Handle detection results" 子步骤（1d: high/medium/low confidence handling），这在 rules 文件中也有对应（Detection Result Format 和 Cold Start Handling）。 | 可接受的冗余——SKILL.md 提供流程概览，rules 文件提供实现细节。但 high/medium/low confidence 的处理逻辑在两处都有描述，可考虑仅保留 rules 文件中的版本。 | medium |
| C-36 | templates/convention-template.md vs rules/convention-structure.md | 2 | INCOMPLETE | P3 | rules/convention-structure.md 说 "Optional sections: Helpers, Import Patterns, Code Style, Anti-patterns"。但 templates/convention-template.md 仅包含 4 个 required sections（Framework, Assertion, Tags, Result Format），不包含 optional sections 的占位符。生成的 Convention 文件如果要包含 optional sections，LLM 需要自行推断格式。 | 在 templates/convention-template.md 底部添加 optional sections 的注释占位符，如 `<!-- Optional sections: uncomment as needed -->` + 各 optional section 的模板。 | low |

### Layer 3: Timing & Flow

No TIMING issues found. The process flow (Check existing -> Detect framework -> Scan files -> Generate draft -> User review -> Write) is clear and sequential. The retry feedback loop (Step 4b -> 4b with counter) is well-defined with max_retries = 2.

---

## Summary Statistics

| Metric | Value |
|--------|-------|
| Skills audited | 7 |
| Total files read | 56 |
| CONFLICT issues | 8 (C-05, C-10, C-15, C-22, C-23, C-24, C-26, C-30) |
| REDUNDANT issues | 10 (C-01, C-02, C-03, C-06, C-07, C-08, C-11, C-27, C-29, C-33) |
| INCOMPLETE issues | 12 (C-04, C-09, C-12, C-13, C-16, C-17, C-20, C-21, C-25, C-28, C-32, C-36) |
| TIMING issues | 4 (T-01, T-02, T-03, T-04, T-05) |
| Total issues | 34 (some merged/escalated) |
| Baseline commit | `1542c8cc` |

## Issue Priority Summary

### P0 (Critical)

| ID | Component | Description |
|----|-----------|-------------|
| C-24 | init-justfile | go.just template's `test` recipe uses `npx playwright test` (Node.js command) -- Go project runtime failure |

### P1 (High)

| ID | Component | Description |
|----|-----------|-------------|
| C-05 | extract-design-md | TUI match strategy defined only in rules/platform-routing.md, not in SKILL.md Step 3 |
| C-10 | gen-contracts | Fact Table format inconsistency: SKILL.md says JSON, rules says Markdown |
| C-15 | gen-journeys | Test level emphasis mismatch between SKILL.md and surface rules (3 out of 5 surfaces wrong) |
| C-22 | init-justfile | SKILL.md HARD-RULE says "Do NOT use templates" but 6 .just template files exist |
| C-23 | init-justfile | LLM-generation instruction conflicts with hardcoded template contents |
| C-26 | init-justfile | Mobile `test-setup` target missing from SKILL.md Surface-Level Targets table |
| C-30 | submit-task | `acceptanceCriteria` rule vs `testsFailed` rule conflict for coding tasks |
| T-03 | gen-sitemap | Step 2b page exploration overlap with Step 4 not handled |

### P2 (Medium)

| ID | Component | Description |
|----|-----------|-------------|
| C-01 | extract-design-md | Match strategy table duplicated in SKILL.md and rules |
| C-02 | extract-design-md | Extraction layers duplicated in SKILL.md and rules |
| C-03 | extract-design-md | Mobile/TUI extraction summaries duplicated |
| C-04 | extract-design-md | design-web.md template missing Design Philosophy section |
| C-06 | gen-contracts | Schema validation checks duplicated in SKILL.md and rules/validation.md |
| C-07 | gen-contracts | Six-dimension rules duplicated in SKILL.md and rules |
| C-08 | gen-contracts | Risk density targets duplicated |
| C-17 | gen-journeys | Quality Notice format not in journey.md template |
| C-20 | gen-sitemap | Dynamic state exploration lacks layout element filtering |
| C-25 | init-justfile | Extra targets in templates not documented in Standard Target Contract |
| C-28 | init-justfile | Surface rules don't mention `<surface-key>` naming convention |
| C-31 | submit-task | testsFailed/testsPassed conditional requirement unclear |
| C-34 | test-guide | Validation section names mismatch between SKILL.md and convention-structure.md |
| T-02 | gen-contracts | Surface detection executed twice (Step 0 and Step 1) |
| T-04 | gen-sitemap | Step 2b degradation when agent-browser unavailable not documented |
| T-05 | init-justfile | Cold start Convention loading happens before language detection |

### P3 (Low)

| ID | Component | Description |
|----|-----------|-------------|
| C-09 | gen-contracts | tea.Batch concurrent semantics not in SKILL.md |
| C-11 | gen-contracts | Journey Invariants rule delegated to TUI-only rules file |
| C-12 | gen-contracts | SKIP_EVAL_GATE frontmatter not in contract.md template |
| C-13 | gen-contracts | Side-effect default "none" not documented in outcome-block.md template |
| C-16 | gen-journeys | quality: low frontmatter field not in journey.md template |
| C-21 | gen-sitemap | loginLocators field in template not referenced in SKILL.md |
| C-27 | init-justfile | Surface rule templates vs LLM-generation design tension |
| C-29 | submit-task | Common Fields duplicated across SKILL.md and all data files |
| C-32 | submit-task | Gate category exception for acceptanceCriteria not in Common Fields |
| C-35 | test-guide | Detection flow duplicated in SKILL.md and rules |
| C-36 | test-guide | Convention template lacks optional section placeholders |
| T-01 | extract-design-md | Screenshot quality check in Error Handling vs process flow |
