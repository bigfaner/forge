---
created: 2026-05-19
updated: 2026-05-19
author: "faner + Claude"
status: Draft
supersedes: test-profile-system
---

# Proposal: Test Knowledge Convention-Driven Model

## Problem

contract-journey-test-model 提案已完整实现：Journey-Driven 测试组织、Contract 六维度验证、gen-journeys/gen-contracts/gen-test-scripts 四步 pipeline、forge test verify/promote/run-journey 命令——全部就绪。

但实现时保留了旧的 Profile 系统（6 个语言 profile + 语言检测 + embed 机制），导致两套模型共存：

**概念矛盾**：Contract 层用语义描述符"推迟精确匹配"——gen-contracts 阶段不猜测输出格式，精确正则延迟到 gen-test-scripts 阶段从 Fact Table 推导。但 Profile 层用 generate.md"提前硬编码技术决策"——在 LLM 看不到项目代码时就决定了断言库、import 模式、测试风格。

**职责混淆**：Profile 系统同时承担了三个本应独立的关注点：
1. 项目技术栈检测（语言、框架、接口类型）
2. 框架特定的代码生成知识（断言库、import、tag、helper）
3. 测试组织轴（language × interfaces）

Journey-Contract 模型已经解决了关注点 3（测试组织轴改为 Journey），但关注点 1 和 2 仍以 Profile 形式残留在系统中。

**config.yaml 职责越界**：当前 `.forge/config.yaml` 包含 `languages`、`interfaces`、`test-framework` 等可通过智能检测获得的信息，违反"config 只管 Forge 行为控制"的原则。

**Profile 渗透深度远超预期**：审计发现 Profile 系统被 19+ 个文件深度消费，不仅限于 gen-test-scripts 和 run-e2e-tests。pkg/journey/、init-justfile、forge config init、forge init、forge task index/add、pkg/e2e/、pkg/just/ 等都直接依赖 Profile 类型和函数。局部删除不可行，必须全面移除。

### Evidence

**Profile 的"零配置便利"是伪优势**：

| 场景 | Profile 行为 | 实际结果 |
|------|-------------|---------|
| Go + go-testing（默认） | detect→go→default→go-testing | 正确 |
| Go + ginkgo | detect→go→default→go-testing | **错误** |
| TypeScript + vitest | detect 不识别 vitest（只识别 @playwright/test） | **检测不到** |
| 框架升级 mocha→vitest | 仍用 mocha 默认映射 | **错误** |

Profile 只对恰好使用默认框架的项目有效。非默认框架反而因自动检测得到错误结果。

**generate.md 是硬编码假设**：它告诉 LLM "Go 项目用 testify"，但新项目可能还没决定用什么断言库。当假设与实际不符时，生成的代码比 LLM 自然生成（用 stdlib）更差。

**Profile 扩展成本高**：新增框架 = 新增 profile 目录 + 修改 embed.go + 修改 framework.go 映射 + 修改 KnownLanguages 常量 + 重新编译二进制。用户无法自行扩展。

**延迟成本**：每新增一个用户请求的非默认框架支持（ginkgo、vitest、pytest 等已在用户反馈中出现），需要修改 forge-cli Go 代码 + 重新编译 + 发布新版本——当前有 2 个待处理的框架支持请求被阻塞。Profile 系统还阻止了 consolidate-specs 对测试知识的完整管理——Convention 文件无法纳入 specs 管理，因为测试知识散落在 Go embed 文件中。继续保留 Profile 意味着每次框架支持请求都走"修改 Go 代码→编译→发布"的重量级流程。

**CLI 不需要框架信息**：`forge task index`、`forge task add` 等 CLI 命令当前调用 `profile.ReadLanguages()` 和 `profile.GetStrategy()`，但这些命令只需要知道"要不要做"（config.yaml `auto.*`），不需要知道"怎么做"。测试执行统一走 `just`（`just e2e-test`、`just e2e-compile`），justfile 自身屏蔽语言差异。

## Industry Benchmarking

### Existing Patterns

The problem of injecting framework-specific knowledge into code generation is well-explored. Three relevant external approaches:

1. **Hygen (template-driven generation)**: Hygen uses EJS templates organized by generator type. Framework knowledge lives in template files with programmatic variable substitution. Strength: deterministic output. Weakness: template maintenance scales linearly with framework count; non-developers cannot edit templates; template syntax constrains output flexibility. Hygen works when output is predictable and framework count is small, but breaks when generation needs contextual judgment (e.g., deciding assertion style based on existing test patterns).

2. **Plop.js (config-driven scaffolding)**: Plop uses a `plopfile.js` with inquirer prompts + Handlebars templates. Users select framework via config prompts, templates produce code. Strength: interactive, user-driven. Weakness: same template rigidity as Hygen. Adding a new framework requires writing new templates — the exact problem we face with Profile's generate.md files. Plop's pattern proves that config-driven framework selection works, but its template approach does not scale for LLM-based generation where output needs to adapt to project context.

3. **Cursor/Windsurf rules (LLM instruction files)**: These tools use `.cursorrules` or `.windsurfrules` markdown files as persistent LLM instructions — convention-like documents that shape generation behavior. Strength: user-editable, framework-agnostic, LLM-native. Weakness: no structured sections, no validation, no fallback when LLM ignores rules. Our Convention approach inherits this pattern's strength (user-editable markdown as LLM guidance) while adding fixed section structure for reliable parsing and a compile gate as the validation layer that pure rules files lack.

### Cross-Domain Parallels

The Convention file pattern is not novel in isolation — it follows a well-established pattern from developer tooling:

- **`.editorconfig`**: User declares formatting rules (indent style, charset, line endings). Editors read and respect these rules. Projects without `.editorconfig` fall back to editor defaults. The pattern: user-editable config + tool respects it + sensible fallback.

- **`.prettierrc`**: User declares code formatting preferences (print width, semicolons, trailing commas). Prettier reads the config and enforces it. Missing config triggers sensible defaults. The pattern: user-editable config + tool respects it + sensible fallback.

- **`tsconfig.json`**: User declares type system behavior (strict mode, module resolution, target). The TypeScript compiler reads and enforces these settings. Missing config falls back to compiler defaults. The pattern: user-editable config + tool respects it + sensible fallback.

Convention files follow the same pattern: user-editable config (markdown with fixed sections) + tool respects it (LLM reads Convention during generation) + sensible fallback (LLM defaults when Convention is missing). The difference is the enforcement mechanism — `.editorconfig` and `tsconfig.json` use programmatic enforcement, while Convention uses the compile gate as its enforcement layer. This is a deliberate trade-off: weaker enforcement than a compiler, but applicable to a domain (code generation style) where no compiler exists.

### Alternatives Considered

| Alternative | Description | Pros | Cons | Verdict |
|-------------|-------------|------|------|---------|
| **A. Config-driven profiles** | User selects framework in config.yaml (e.g., `test-framework: ginkgo`); system loads corresponding template set | Deterministic output with no LLM interpretation risk; proven pattern used by most code generators; easy to reason about correctness; simple mental model for users | Template per framework does not adapt to project-specific patterns (e.g., a team's custom assertion helper); framework additions require Forge code changes and binary release; user selection can be wrong just as easily as auto-detection | Not selected — while deterministic and well-understood, the per-framework template model creates the same maintenance bottleneck as Profile (Forge release cycle for framework support). Acceptable for projects that use default frameworks exclusively. |
| **B. AST-based detection + generation** | Parse existing test files with language-specific AST parsers to extract framework, assertion style, and patterns programmatically | Fully deterministic detection with no LLM ambiguity; precise extraction of actual project patterns from code; mature parser libraries available (Go/AST, TypeScript compiler API, Python/lib2tope) | Requires per-language parser integration and maintenance; cannot handle cold starts (no existing files to parse); parser rules must be updated when framework conventions change; multi-language support multiplies maintenance cost | Not selected — strongest option for detection accuracy on existing projects, but cold start gap and per-language parser maintenance make it a complement (Code Reconnaissance) rather than a primary mechanism. |
| **C. LLM-only (no Convention files)** | Remove Profile; rely entirely on LLM training data + Code Reconnaissance with no user-editable knowledge layer | Simplest implementation with fewest new components; zero user-facing changes or new files to manage; LLM training data already covers mainstream frameworks well | No user control over generation behavior when LLM defaults are wrong; no persistence of discovered patterns across sessions; generation style may vary between runs; debugging generation failures requires reading LLM context, not editing a file | Not selected — most elegant for default-framework projects, but user agency loss is significant: when generation is wrong, the only recourse is re-running and hoping for different output. Acceptable as the fallback behavior (which Convention uses for missing files). |
| **D. Convention files (proposed)** | User-editable markdown with fixed structure; LLM reads as guidance; compile gate validates output | User-controllable and editable; persistent across sessions; minimal infrastructure; graceful degradation (missing Convention falls back to LLM defaults) | LLM may ignore Convention content; relies on LLM interpretation, not programmatic parsing; user must maintain Convention file when framework changes | Selected — the compile gate mitigates the LLM interpretation risk, and user-editability addresses the maintenance bottleneck of A and B. See POC validation below |

### Why Convention over Alternatives

Each alternative has genuine strengths: A is deterministic, B is precise, C is simple. The Convention approach trades away deterministic correctness for user-editability — the same trade-off that `.editorconfig` and `.prettierrc` make successfully. The key insight is that the compile gate converts "LLM might misinterpret Convention" from a silent failure into a caught-and-retried failure, shifting the risk profile from "wrong code silently accepted" to "wrong code caught at generation time." See "LLM Convention Compliance Validation" in Feasibility for the POC that validates this trade-off.

## Proposed Solution

将框架知识注入从"二进制内嵌的 Profile 文件"改为"用户可编辑的 Convention 文件"，全面移除 Profile 系统。

### 核心原则

1. **config.yaml 只保留 Forge 行为控制项**（`auto.*`、`test-command`、`worktree`），不存放可智能检测的信息。
2. **测试执行统一走 just**——skill 和 CLI 通过 `just e2e-test` / `just e2e-compile` 等命令操作，justfile 屏蔽语言差异。
3. **框架知识由 Convention 提供**——用户可编辑、slash command 可引导生成、按 domains 按需加载。
4. **全面移除 Profile**——不保留内部实现，所有消费者重写。

### 三层分离

```
┌─────────────────────────────────────────────────┐
│  方法论层（Skill prompt）                         │
│  Forge 发明的规则，不可扩展也不需要扩展             │
│  输出规则、标签生命周期、质量门禁、Pipeline 结构    │
│  不含任何框架特定代码或分支逻辑                     │
├─────────────────────────────────────────────────┤
│  框架知识层（docs/conventions/）                   │
│  用户可编辑，/forge:test-guide 可引导生成          │
│  固定结构，skill 按结构解读                        │
│  替代 generate.md + templates/                   │
├─────────────────────────────────────────────────┤
│  项目实际层（Code Reconnaissance）                 │
│  运行时自动，不需要用户参与                        │
│  从已有测试文件提取代码模式                        │
│  补充 conventions 未覆盖的细节                    │
└─────────────────────────────────────────────────┘
```

### 方法论层：Skill Prompt 保留什么

gen-test-scripts 的 SKILL.md 只包含 Forge 方法论：

**保留**：
- Pipeline 结构（Step 1-4 的输入输出）
- 输出组织规则（`tests/<journey>/`，每步一个文件，每 Journey 一个烟测试）
- 标签生命周期规则（`@feature` → `@regression`，具体语法从 Convention 获取）
- 可追溯性注释要求（格式从 Convention 获取）
- 接口类型侦察指引（告诉 Code Reconnaissance 该收集什么）
- 质量门禁（编译检查通过 `just e2e-compile`、VERIFY 标记检查、重复名称检查）
- Pipeline 硬规则（单 Journey 调用、批次拆分阈值、Outcome 互斥校验）

**移除**：
- 框架特定 import 列表 → Convention 提供
- 框架特定断言语法 → Convention 提供
- 框架特定代码示例 → Convention 提供
- 框架特定反模式规则 → Convention 提供
- 模板文件引用 → Convention + Code Reconnaissance 替代
- 所有 `forge test` CLI 调用（detect / get generate / framework / interfaces）

**关键设计**：Skill prompt 中没有 `if framework == "go" then ...` 的分支逻辑。框架差异完全由 Convention 内容描述，Skill 不感知具体框架。这一原则同样适用于 run-e2e-tests：结果解析策略从 Convention 的 Result Format section 读取，skill prompt 不包含按框架切换的解析逻辑。

### 框架知识层：docs/conventions/

Convention 文件是用户可编辑的项目级文档，替代 generate.md + templates/ 的角色。

#### 固定结构

Convention 文件遵循固定结构，skill 按 section 标题解读（LLM 读取 markdown，非程序化解析）。最小集包含以下必需 section：

```markdown
---
domains: [testing, go]
---

# Test Conventions

## Framework
- <framework name + version>
- File pattern: <file naming convention>
- Package/module: <package name>

## Assertion
- <assertion library>
- <key assertion functions>

## Tags
- <build tag / marker syntax and value>

## Result Format
- Output command flags: <e.g., "-json", "-v", "--reporter=json">
- Format type: <json-stream | json-report | text-verbose>
```

The Result Format section is part of the minimum set. It declares how the test runner produces output, enabling run-e2e-tests to parse results without embedding framework-specific parsing knowledge. Common values: Go (`go test -json`, json-stream), Python (`pytest -v`, text-verbose), JavaScript (`npx playwright test --reporter=json`, json-report). When Result Format is missing from a Convention file, run-e2e-tests falls back to text-based output parsing.

Convention 文件可以包含更多 section（Helpers、Import Patterns、Code Style、Anti-patterns 等），随项目演进逐步添加。test-guide 首次生成只写最小集（Framework + Assertion + Tags + Result Format）。

缺失 section 的处理：skill 在 Code Reconnaissance 步骤中尝试补充缺失信息（如 Convention 有 Framework 但无 Assertion → 从已有测试文件 import 分析断言库）。无法补充时用 LLM 默认值。不阻塞生成。

#### 文件组织

用户可按项目需要组织多个测试 Convention 文件，通过 `domains` frontmatter 按需加载：

```
docs/conventions/
  testing-go.md              domains: [testing, go]
  testing-javascript.md      domains: [testing, javascript, web-ui]
  testing-antipatterns.md    domains: [testing]           ← 共享反模式（用户后续添加）
```

#### Convention 加载机制

Convention 加载是 LLM 行为指引，不是程序化过滤。Skill prompt 指示 LLM：
1. 列出 `docs/conventions/` 目录中所有文件
2. 读取每个文件的 `domains` frontmatter
3. 仅加载 domains 与当前任务匹配的文件内容
4. 按 Convention 文件的固定结构解读内容

这是现有机制的自然使用——guide.md 已指示 agent 按 domains 加载 convention 文件。不新增基础设施。

#### Convention 文件内容示例（forge-cli 项目）

最小集：

```markdown
---
domains: [testing, go, cli]
---

# Test Conventions

## Framework
- Go testing package + testify/assert
- File pattern: *_test.go
- Package: e2e

## Assertion
- assert (not require)
- assert.NoError / assert.Contains / assert.Equal

## Tags
- //go:build e2e

## Result Format
- Output command flags: -json
- Format type: json-stream
```

随着项目演进，Convention 文件逐步充实，最终可能包含：

```markdown
## Helpers
- runForge(args ...string) (string, int, error)
- readJSON(t, path string) map[string]interface{}

## Import Patterns
- "os/exec" for CLI invocation
- "github.com/stretchr/testify/assert" for assertions

## Anti-patterns
- See testing-antipatterns.md
```

#### 增长路径

Convention 不是一次写完的，随项目演进逐步充实：

```
/forge:test-guide 首次运行 → 最小集（Framework + Assertion + Tags + Result Format）
                    ↓
项目写了几轮测试后，再次 test-guide → 补充 Helpers、Import Patterns
                    ↓
用户手动编辑 → 加反模式引用、风格偏好
                    ↓
consolidate-specs → 纳入管理、drift 检测、去重
```

#### 用户维护成本说明

Convention 文件将框架知识维护从 Forge 系统转移到用户项目。这一转移是有意的：
- **Profile 时代**：框架升级时用户无法自行更新，必须等 Forge 发版（当前有 2 个框架支持请求被阻塞）
- **Convention 时代**：框架升级时用户直接编辑 Convention 文件，无需等 Forge 发版

成本对比：Profile 时代用户无法解决框架不匹配问题（阻塞）；Convention 时代用户需要编辑一个 markdown 文件（3-5 分钟）。`/forge:test-guide` 可重新运行以检测框架变化并建议更新。consolidate-specs 的 drift 检测会在 Convention 内容与实际测试模式不一致时发出提醒。

#### 用户发现机制

- 用户手册说明 Convention 文件的用途、格式和位置
- `/forge:test-guide` slash command 引导用户生成
- gen-test-scripts 执行时若未找到相关 Convention，输出提示信息，不阻塞生成流程

### Slash Command：/forge:test-guide

引导用户生成测试 Convention 文件的 slash command。统一流程（首次和更新走同一路径）。

#### 交互模型

test-guide 是多轮对话式 skill。Skill 在单次响应中完成检测和展示，用户通过后续消息确认或调整。这是 Claude Code skill 的自然交互模式——skill 输出结果和建议，用户在对话中跟进。

```
Step 1: 检测项目上下文
        扫描项目文件系统：
        - 文件信号（go.mod / package.json / Cargo.toml 等）
        - 已有测试文件（*_test.go / test_*.py / *.spec.ts 等）
        - 已有测试文件的 import 和模式分析
        │
        ├── 有测试文件 → Step 2a
        └── 无测试文件 → Step 2b

Step 2a: 提取模式（有测试文件）
        分析已有测试文件：
        - import → 断言库、helper 包
        - 函数签名 → 命名模式
        - build tag → tag 约定
        - 包声明 → 包名
        展示提取结果，询问用户确认。

Step 2b: 询问用户（冷启动）
        基于文件信号检测到的语言，给出候选框架列表：
        "检测到 Go 项目，选择测试框架：
         1. testing + stdlib
         2. testing + testify
         3. ginkgo
         4. 其他..."
        只问框架，其他细节用 LLM 默认值。

Step 3: 生成 Convention 文件
        用户确认后，写入 docs/conventions/testing-<scope>.md
        每框架一个文件
        固定结构（最小集：Framework + Assertion + Tags）
        带 "<!-- auto-generated by forge:test-guide -->" 标记
        consolidate-specs 将自动纳入管理
```

#### 检测策略

test-guide 的检测由 LLM 执行（非确定性），与当前 detect.go 的确定性检测不同。这是可接受的权衡——test-guide 的检测结果由用户确认后才写入文件，错误会被用户捕获。而编译门禁在 gen-test-scripts 阶段提供最终兜底。

### 执行层：just 统一抽象

run-e2e-tests、init-justfile 及其他需要执行测试命令的组件，统一通过 `just` 调用：

| 操作 | 调用方式 | justfile 屏蔽的差异 |
|------|---------|-------------------|
| 编译检查 | `just e2e-compile` | Go: `go test -c` / Python: `pytest --collect-only` / JS: `tsc --noEmit` |
| 运行测试 | `just e2e-test` | Go: `go test -tags=e2e` / Python: `pytest` / JS: `npx playwright test` |
| 环境准备 | `just e2e-setup` | Go: `go build` / Python: `pip install` / JS: `npm install` |

#### 结果输出契约

justfile 的 `e2e-test` recipe 必须以统一格式输出测试结果。init-justfile 生成的 recipe 包含结果格式化逻辑：

- Go：`go test -json` 输出 JSON event stream
- Python：`pytest --json-report` 输出 JSON report（或 `-v` 文本 + run-e2e-tests skill 解析）
- JavaScript：`npx playwright test --reporter=json` 输出 JSON report

run-e2e-tests skill reads the Convention file's Result Format section to determine parsing strategy — not hardcoded framework knowledge in the skill prompt. The Convention declares the output format type (json-stream, json-report, text-verbose) and the command flags that produce it. run-e2e-tests uses this declaration to select the appropriate parsing logic. This is a deliberate improvement over the current approach: framework-specific knowledge moves from an embedded Go binary (invisible to users, requiring Forge release to change) to a user-editable Convention file (visible, modifiable without Forge release). When Convention has no Result Format section, run-e2e-tests falls back to generic text-based parsing.

### 项目实际层：Code Reconnaissance

Code Reconnaissance（Fact Table）是 skill prompt 中的一个思考步骤——LLM 执行时读源码，在上下文里构建 markdown 表格作为中间笔记，不落盘。

**定位**：它是运行时的补充侦察，不是持久化的知识库。它补充 Convention 未覆盖的项目特定细节。

**侦察范围扩展**：

当前 Fact Table 只收集被测代码信息（CLI 入口、API handler、TUI Model）。扩展后同时收集测试代码信息：

```
Code Reconnaissance
├── 被测代码侦察（现有，不变）
│   ├── CLI entry points
│   ├── API handlers
│   └── TUI components
└── 测试框架侦察（新增）
    ├── 测试文件模式 → test.file-pattern
    ├── import 分析 → test.assertion-lib, test.helpers
    ├── build tag 分析 → test.build-tag-value
    └── 函数签名分析 → test.naming-pattern
```

侦察结果作为 LLM 的中间笔记，补充 Convention 文件中可能缺失的项目特定模式。不替代 Convention——Convention 是稳定的、用户审核过的知识；Reconnaissance 是即时的、自动提取的补充。

### 质量保证：Convention + Reconnaissance + 编译门禁

去掉 generate.md 后，质量保证链为：

```
预生成：Convention 提供框架知识 + Code Reconnaissance 补充项目模式
    ↓
生成：LLM 按 Convention + Reconnaissance + Skill 方法论生成代码
    ↓
后验证：just e2e-compile
    ↓
不通过 → 编译错误反馈给 LLM → 重新生成（最多 2 次重试）
```

**冷启动处理**：新项目没有测试文件、没有 Convention 文件时：
- Skill 使用 LLM 训练数据生成代码
- 编译门禁兜底（`just e2e-compile`）
- 首次成功生成后，用户可通过 `/forge:test-guide` 固化模式
- 冷启动初始化顺序：`init-justfile`（生成 justfile，内含 e2e-compile recipe）→ gen-test-scripts（LLM 默认 + 编译门禁）→ test-guide（固化模式）

**justfile 自身验证**：冷启动时 justfile 是 LLM 生成的，存在循环依赖风险（justfile 错误导致编译门禁基于错误 recipe 执行）。应对措施：
1. init-justfile 生成后展示 recipe 内容给用户确认（`just --dry-run e2e-compile` 显示实际命令）
2. init-justfile 使用保守策略：只生成 `just --list` 可验证的 recipe，不生成复杂 shell 逻辑
3. 冷启动的 e2e-compile recipe 使用最常见模式（Go: `go test -c ./...`，Python: `python -m pytest --collect-only`，JS: `npx tsc --noEmit`），LLM 对这些 recipe 的训练数据充足，错误率低
4. 如果 justfile recipe 编译失败，错误信息包含实际执行的命令，用户可直接修正 recipe

### config.yaml 清理

**只保留 Forge 行为控制**：

```yaml
# Forge behavioral controls only
test-command: "just e2e-test"     # run-journey 使用的测试执行命令
auto:
  e2eTest:
    quick: false
  consolidateSpecs:
    quick: true
    full: true
  cleanCode:
    quick: false
    full: true
  gitPush: true
worktree:                          # worktree 配置
  source-branch: main
  copy-files: []
```

`test-command`、`auto.*`、`worktree` 是 Forge 行为控制，不是可检测信息，保留在 config.yaml 中。

**移除可检测字段**：

| 移除字段 | 理由 |
|---------|------|
| `languages` | CLI 不需要框架信息；Skill 从 Code Reconnaissance 获取 |
| `interfaces` | Skill 从 Code Reconnaissance 获取 |
| `test-framework` | Convention 声明 |
| `project-type` | pkg/just 通过其他方式解析 scope |

### 全面移除 Profile：所有消费者处理

以下为 Profile 系统所有消费者的完整清单及处理方式：

#### 删除的组件

| 组件 | 原位置 | 删除理由 |
|------|--------|---------|
| `pkg/profile/` 整个目录 | forge-cli/pkg/profile/ | 全面移除 Profile 系统 |
| `pkg/profile/languages/`（6 个语言子目录） | pkg/profile/languages/ | 框架知识迁移到 Convention |
| `generate.md`（每语言） | pkg/profile/languages/*/generate.md | 被 Convention 替代 |
| `run.md`（每语言） | pkg/profile/languages/*/run.md | 测试执行统一走 just，结果解析策略从 Convention 的 Result Format section 读取 |
| `graduate.md`（每语言） | pkg/profile/languages/*/graduate.md | 已被 forge test promote 替代 |
| `templates/`（每语言） | pkg/profile/languages/*/templates/ | 被 Convention + Code Reconnaissance 替代 |
| `justfile-recipes`（每语言） | pkg/profile/languages/*/justfile-recipes | init-justfile 重写，从 Convention + LLM 知识生成 recipe |
| `embed.go` | pkg/profile/embed.go | 不再需要 embed |
| `detect.go` | pkg/profile/detect.go | 检测由 LLM 在 test-guide 和 Code Reconnaissance 中执行 |
| `framework.go` | pkg/profile/framework.go | 框架信息从 Convention 获取 |
| `config.go` 的 Languages/Interfaces 字段 | pkg/profile/config.go | config.yaml 不再存放可检测信息 |
| `forge test detect` 命令 | internal/cmd/test.go | 不再需要 |
| `forge test get` 命令组（generate/run/justfile/template） | internal/cmd/test.go | 不再需要 |
| `forge test interfaces` 命令 | internal/cmd/test.go | 不再需要 |
| `forge test framework` 命令 | internal/cmd/test.go | 不再需要 |

#### 重写的组件

| 组件 | 原位置 | 重写内容 |
|------|--------|---------|
| gen-test-scripts skill | plugins/forge/skills/gen-test-scripts/ | 移除 Profile 依赖，Convention 加载 + Code Reconnaissance + just 编译门禁 |
| run-e2e-tests skill | plugins/forge/skills/run-e2e-tests/ | 移除 Profile 依赖，just 执行 + Convention Result Format 驱动的结果解析（无硬编码框架知识） |
| init-justfile skill | plugins/forge/skills/init-justfile/ | 移除 `forge test detect` / `forge test get justfile` 调用，从 Convention + LLM 知识生成 e2e recipe |
| pkg/journey/journey.go | forge-cli/pkg/journey/ | 移除 `profile.FrameworkInfo` 依赖，`FeatureTag`/`RegressionTag` 从 Convention 的 Tags section 获取，`GenerateDispatched()` 重写为 Convention 驱动 |
| internal/cmd/config.go（forge config init） | forge-cli/internal/cmd/ | 移除 languages/interfaces 交互收集，只收集 auto.* + worktree |
| internal/cmd/init.go（forge init） | forge-cli/internal/cmd/ | 同上，移除 languages/interfaces 收集 |
| internal/cmd/index.go（forge task index） | forge-cli/internal/cmd/ | 移除 Profile 调用，e2e 测试任务生成基于 config.yaml auto.e2eTest 而非 profile |
| internal/cmd/add.go（forge task add） | forge-cli/internal/cmd/ | 同上 |
| pkg/task/testgen.go | forge-cli/pkg/task/ | 移除 `profile.Language`、`profile.AutoConfig` 类型依赖，重写测试任务生成模型 |
| pkg/task/build.go | forge-cli/pkg/task/ | 移除 `StrategyResolver`、`Languages`、`TestInterfaces` 字段 |
| pkg/e2e/e2e.go | forge-cli/pkg/e2e/ | 移除 `profile.ReadLanguages()`、`profile.IsKnownLanguage()` 调用，测试执行通过 just |
| pkg/just/just.go | forge-cli/pkg/just/ | 移除 `profile.ReadConfig()` 中对 project-type 的依赖，用其他方式解析 scope |

#### 保留的组件

| 组件 | 位置 | 保留理由 |
|------|------|---------|
| `forge test verify` | internal/cmd/test_verify.go | 契约断裂检测，与 Profile 无关 |
| `forge test promote` | internal/cmd/test_promote.go | 标签晋升，与 Profile 无关 |
| `forge test run-journey` | internal/cmd/test.go | Journey 隔离执行，与 Profile 无关（使用 config.yaml test-command） |
| `pkg/contract/` | forge-cli/pkg/contract/ | Contract 解析和验证，与 Profile 无关 |
| gen-journeys skill | plugins/forge/skills/gen-journeys/ | 纯叙述性提取，不依赖 Profile |
| gen-contracts skill | plugins/forge/skills/gen-contracts/ | 代码侦察 + Contract 生成，不依赖 Profile |
| Journey isolation 基础设施 | internal/cmd/journey_isolation*.go | 与 Profile 无关 |
| `/forge:test-guide` | plugins/forge/skills/test-guide/ | 新建 |

### forge task index/add 重设计

当前 `forge task index` 使用 `profile.ReadLanguages()` 和 `profile.GetStrategy()` 决定是否生成 e2e 测试任务，并嵌入 `StrategyContent`（generate.md 内容）到任务文件中。

移除 Profile 后：

- **是否生成 e2e 任务**：由 config.yaml `auto.e2eTest.full` / `auto.e2eTest.quick` 决定，不再依赖 language detection
- **任务内容中的测试策略**：不再嵌入 generate.md 内容。测试策略由 gen-test-scripts skill 在执行时从 Convention + Code Reconnaissance 获取
- **task index 生成的测试任务**：简化为引用 Convention 文件路径 + Journey 名称，不包含框架特定策略

### pkg/journey/ 重设计

当前 `journey.go` 的 `FeatureTag()`/`RegressionTag()` 对 `profile.FrameworkInfo.Name` 做 switch 匹配，`GenerateDispatched()` 按框架名分派到不同生成器。

移除 Profile 后：

- **标签生成**：从 Convention 文件的 `Tags` section 读取标签语法和值。LLM 在 gen-test-scripts 执行时根据 Convention 内容生成正确的标签代码
- **框架分派**：移除 Go 代码中的 switch 分派。所有框架的代码生成统一由 LLM 完成（通过 Convention + Code Reconnaissance）
- **journey.go 简化**：移除 `profile.FrameworkInfo` 依赖，保留 Journey 数据模型和文件操作逻辑

### forge config init / forge init 重设计

当前 `forge config init` 和 `forge init` 交互收集 project-type、languages、interfaces。

移除后：

- **交互流程简化**：只收集 Forge 行为控制项
  - auto.e2eTest.quick / full
  - auto.consolidateSpecs.quick / full
  - auto.cleanCode.quick / full
  - auto.gitPush
  - worktree.source-branch / copy-files
  - test-command（默认 `just e2e-test`）
- **不再收集**：project-type、languages、interfaces、test-framework
- **生成的 config.yaml**：只包含 auto.* + worktree + test-command

### init-justfile 重设计

当前 init-justfile Step 0 调用 `forge test detect` 和 `forge test get justfile` 获取语言和 recipe 内容。

移除后：

- **Step 0**：读取 Convention 文件获取框架信息。无 Convention 时，LLM 通过 Code Reconnaissance 分析项目文件信号推断
- **e2e recipe 生成**：LLM 根据 Convention 的 Framework section + LLM 训练数据生成正确的 e2e-compile、e2e-test、e2e-setup recipe
- **recipe 内容不再来自 embed 文件**，而是由 LLM 动态生成

## Requirements Analysis

### Key Scenarios

- **Scenario 1: 常见框架零配置**：Go + go-testing 项目，无 Convention 文件。gen-test-scripts 通过 Code Reconnaissance 发现 go.mod 和已有 *_test.go 文件，LLM 用训练数据生成标准 Go testing 代码，`just e2e-compile` 验证通过。
- **Scenario 2: Convention 自定义**：项目使用 Go + ginkgo。用户创建 Convention 文件声明 ginkgo。Skill 加载 Convention，生成 ginkgo 风格代码。`just e2e-compile` 验证通过。
- **Scenario 3: 多框架项目**：Go 后端 + JS 前端。test-guide 生成 `testing-go.md` 和 `testing-javascript.md`。gen-test-scripts 根据 Journey 涉及的接口类型按需加载。
- **Scenario 4: 冷启动**：全新项目。init-justfile 生成 justfile（含 e2e-compile recipe）→ gen-test-scripts 用 LLM 默认生成代码 → 编译门禁兜底 → test-guide 固化模式。
- **Scenario 5: Convention 缺失提示**：gen-test-scripts 未找到 Convention，输出提示。不阻塞。
- **Scenario 6: forge task index 正常**：基于 config.yaml auto.e2eTest 决定是否生成 e2e 测试任务。不依赖 profile。
- **Scenario 7: forge config init 简化**：只收集 auto.* + worktree + test-command。不收集 languages/interfaces。
- **Scenario 8: init-justfile 正常**：从 Convention + Code Reconnaissance 获取框架信息，生成正确的 e2e recipe。
- **Scenario 9: 向后兼容**：现有 126+ 测试通过。gen-test-scripts 通过 Reconnaissance 提取模式，代码风格一致（diff 为空）。

### Non-Functional Requirements

- **向后兼容**：所有现有 e2e 测试通过（`just e2e-compile` + 测试输出 diff 为空）
- **渐进迁移**：用户按节奏创建 Convention 文件，不存在时不阻塞生成
- **零新增基础设施**：Convention 加载使用现有 `docs/conventions/` + `domains` 指引机制
- **编译门禁兜底**：生成后必须通过 `just e2e-compile`，不通过则反馈重试（最多 2 次）
- **Convention 按需加载**：LLM 按 domains frontmatter 只加载相关文件

### Constraints & Dependencies

- 依赖现有 `docs/conventions/` 目录和 domains 指引加载机制
- 依赖 `just` 作为测试执行的统一抽象层
- 不依赖外部服务或 API
- Convention 文件为 markdown 格式，固定 section 结构，用户可编辑

### Phased Delivery Timeline

| Phase | Scope | Estimated Duration | Dependencies |
|-------|-------|--------------------|--------------|
| **Phase 0: POC** | Validate LLM Convention compliance — generate test code for 3 frameworks (Go testing, ginkgo, vitest) using Convention files only; measure first-pass compile rate | 2-3 days | None — can run on current codebase |
| **Phase 1: Profile removal + config cleanup** | Delete `pkg/profile/`, clean config.yaml, rewrite CLI commands (forge config init, forge init, forge task index/add, forge test subcommand removal), rewrite pkg/journey/, pkg/e2e/, pkg/just/, pkg/task/ | 5-7 days | Phase 0 validates core assumption |
| **Phase 2: Skill rewrites** | gen-test-scripts (Convention + Reconnaissance), run-e2e-tests (just + inline parsing), init-justfile (Convention-driven recipe generation) | 3-5 days | Phase 1 complete (skills depend on Profile-free CLI) |
| **Phase 3: test-guide + Convention bootstrap** | New `/forge:test-guide` slash command, Convention fixed structure definition, user manual update, consolidate-specs integration | 2-3 days | Phase 2 complete |
| **Validation** | Run against forge-cli's 126+ e2e tests; measure first-pass compile rate and diff equivalence against Profile-generated output | 1-2 days | Phase 3 complete |

**Total estimate: 13-20 working days, single developer.** Phase 0 is the critical gate — if LLM Convention compliance POC fails (first-pass compile rate < 70%), the approach needs rethinking before committing to Phase 1.

### LLM Convention Compliance Validation (Phase 0 POC)

The foundational assumption — "LLM reads Convention markdown and correctly applies it during code generation" — must be validated before committing to the rewrite.

**POC design**:
1. Write Convention files for 3 frameworks: Go testing (forge-cli's actual setup), Go ginkgo, TypeScript vitest
2. Run gen-test-scripts with Convention files loaded, without any Profile content
3. Measure: (a) first-pass `just e2e-compile` success rate, (b) Convention-specified import/assertion/tag accuracy in generated code
4. Compare against Profile-generated output for the same inputs

**Fallback if POC fails**: If LLM ignores Convention content in > 30% of generation runs:
- Add structured extraction: skill prompt explicitly extracts Convention values into a Fact Table entry before code generation (programmatic read step before creative generation step)
- Add Convention content repetition: key sections (Framework, Assertion, Tags) are included in the generation prompt's system instructions, not just in the context
- Worst case: Convention files include YAML frontmatter for programmatic extraction of critical fields (framework name, assertion library, tag value), with markdown body for LLM-consumable guidance

## Scope

### In Scope

**全面移除 Profile**：
- 删除 `pkg/profile/` 整个目录
- config.yaml 清理（移除 languages / interfaces / test-framework / project-type）

**Skill 重写**：
- gen-test-scripts skill（Convention + Reconnaissance + just）
- run-e2e-tests skill（just 执行 + 内联结果解析）
- init-justfile skill（Convention + Reconnaissance 驱动 recipe 生成）

**Go 代码重写**：
- pkg/journey/（移除 FrameworkInfo 依赖，Convention 驱动标签和生成）
- internal/cmd/config.go（forge config init 简化）
- internal/cmd/init.go（forge init 简化）
- internal/cmd/index.go（forge task index 移除 Profile）
- internal/cmd/add.go（forge task add 移除 Profile）
- pkg/task/testgen.go + build.go（重写测试任务生成模型）
- pkg/e2e/e2e.go（移除 Profile 依赖）
- pkg/just/just.go（移除 project-type 依赖）
- internal/cmd/test.go（移除 detect / get / interfaces / framework 命令）

**新建**：
- `/forge:test-guide` slash command
- Convention 文件固定结构定义
- 用户手册更新

**集成**：
- Convention 文件纳入 consolidate-specs 管理

### Out of Scope

- gen-journeys skill（不依赖 Profile）
- gen-contracts skill（不依赖 Profile）
- forge test verify / promote（不依赖 Profile）
- 现有测试迁移
- Unit 测试
- 反模式文档生成
- Convention 文件自动同步

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 全面重写引入回归 | H | H | 现有 126+ e2e 测试作为回归安全网；逐模块重写，每模块先通过回归验证 |
| 无 Convention 时 LLM 生成质量低于 generate.md | M | M | Code Reconnaissance 从已有测试提取 import 和断言模式；LLM 对 Go testing / pytest / Jest 等主流框架训练数据充足，非 Convention 场景不劣于 Profile |
| 冷启动生成质量不可控 | M | M | init-justfile 先行生成 justfile（含 e2e-compile recipe）；编译门禁验证；首次生成后立即通过 test-guide 固化模式 |
| journey.go 重写后标签生成错误 | M | M | Convention Tags section 提供明确标签语法（如 `//go:build e2e`）；编译门禁验证标签语法 |
| run-e2e-tests 结果解析不稳定 | M | H | justfile recipe 使用结构化输出格式（`go test -json` / `pytest --json-report` / `playwright --reporter=json`）；run-e2e-tests 从 Convention 的 Result Format section 读取格式类型和解析策略，不硬编码框架知识 |
| forge task index 重写后测试任务生成缺失 | M | M | auto.e2eTest 开关决定是否生成；测试策略由 skill 运行时获取而非预嵌入；task index 的 e2e 任务生成逻辑简化为 config 开关检查 |
| 工期超预期（19+ 文件重写） | M | M | 按依赖链分阶段交付；Phase 0 POC 作为早期中止点 |

### Rollback Strategy

Profile 系统的全面移除是高风险操作。如果 Convention 方案在 Phase 2（Skill 重写）后验证失败（首过编译率持续低于 70%），回滚路径如下：

1. **分支策略**：在 `v3.0.0` 分支上开发，`main` 分支保持 Profile 系统不变。合并前所有 Phase 通过验证。
2. **Phase 1 回滚**（Profile 删除 + config 清理）：Phase 1 只删除代码，不改变 skill 行为。回滚 = revert Phase 1 commit，恢复 `pkg/profile/` 目录和 config 字段。`pkg/profile/` 在删除前完整保留在 git 历史中。
3. **Phase 2 回滚**（Skill 重写）：Skill 文件在 `plugins/forge/skills/` 中，每个 skill 独立。如果某个 skill 的 Convention 版本不达标，可单独回滚该 skill 到 Profile 版本，其他已验证的 skill 保持 Convention 版本。混合模式可行——Profile 版 skill 和 Convention 版 skill 可以共存，因为它们通过不同机制获取框架知识。
4. **不可回滚点**：Phase 2 全部完成且 Phase 3 开始后，回滚成本显著上升（test-guide 依赖无 Profile 的环境）。此点之前是回滚窗口。

## Success Criteria

- [ ] `pkg/profile/` 目录不存在
- [ ] config.yaml 只包含 `auto.*` + `test-command` + `worktree`，不包含 `languages` / `interfaces` / `test-framework` / `project-type`
- [ ] 无任何代码导入 `forge-cli/pkg/profile`
- [ ] gen-test-scripts 不调用任何 `forge test` 子命令，从 Convention + Code Reconnaissance 获取框架知识
- [ ] run-e2e-tests 通过 `just e2e-test` 执行测试，结果解析策略从 Convention 的 Result Format section 读取，不包含硬编码的框架特定解析逻辑
- [ ] init-justfile 从 Convention + Code Reconnaissance 获取框架信息，不调用已删除的 CLI 命令
- [ ] pkg/journey/ 不导入 pkg/profile，标签生成基于 Convention 或 Reconnaissance
- [ ] forge task index 基于 config.yaml auto.e2eTest 决定是否生成 e2e 任务，不依赖 Profile
- [ ] forge config init / forge init 只收集 auto.* + worktree + test-command
- [ ] 所有现有 e2e 测试通过（`just e2e-compile` + 测试输出 diff 为空）
- [ ] **首过编译率 >= 85%**：gen-test-scripts 生成的代码在首次 `just e2e-compile` 中通过率不低于 85%（在 forge-cli 项目自身的 126+ 测试上测量，作为质量基准线）
- [ ] **生成代码 diff 等效性**：对现有 Journey 的测试文件，Convention 模式重新生成的代码与 Profile 模式生成代码的 diff 中，框架核心模式（import、断言函数、tag 语法）差异为 0；允许风格差异（变量命名、注释措辞）
- [ ] `/forge:test-guide` 引导式生成 Convention 文件，每框架一个文件
- [ ] Convention 文件固定结构（Framework / Assertion / Tags / Result Format），skill 在 3 个不同框架项目（Go testing、Go ginkgo、TypeScript vitest）上正确解读 Convention 内容，生成的 import/断言/tag 与 Convention 声明一致
- [ ] Convention 文件纳入 consolidate-specs 管理
- [ ] gen-test-scripts 未找到 Convention 时输出提示，不阻塞生成

## Next Steps

- Proceed to `/write-prd` to formalize requirements
