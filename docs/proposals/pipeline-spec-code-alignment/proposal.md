---
created: "2026-05-26"
author: "faner"
status: Draft
---

# Proposal: 流水线全链路审计修复 (2026-05-26)

## Problem

对 Forge 任务流水线（quick-tasks / breakdown-tasks / run-tasks / task-executor / quality-gate / test skills / submit-task / fix-bug / eval / gen-contracts / gen-test-scripts / ARCHITECTURE / conventions / prompt templates / build/submit Go 代码）进行全链路审计，发现 **73 个问题**。原始 6 个 lesson 仅暴露了冰山一角。深层问题分为六类：**幽灵配置字段与过时引用**（scope→surface 迁移遗留、interfaces/surfaces 术语分裂、不存在的文件/字段/类型）、**跨文件逻辑矛盾**（契约不匹配、status 处理分歧）、**模板变量缺失**（fix-task 渲染失败）、**Go 代码双源真相与类型缺口**（index.json 与 .md 不一致、ValidTypes 缺失条目）、**测试管道缺口**（SKIP_EVAL_GATE、surface type 不一致）、**架构文档虚构**（不存在的 agent、不匹配的算法描述）。

### Evidence

#### A. 幽灵配置字段与过时引用（Critical/High）

- **A1**: `breakdown-tasks/rules/db-schema.md:36` 仍使用 `scope: "backend"`，scope 字段已被 surface-key/surface-type 替代。
- **A2**: `breakdown-tasks/rules/existing-code-split.md:24,32` 引用已废弃的 "scope assignment algorithm"，示例值 `"backend"`/`"all"` 是旧 scope 体系。
- **A3**: `breakdown-tasks/rules/scope-assignment.md` 已标记废弃但仍存在。SKILL.md 未明确指示忽略该文件。
- **A4**: `breakdown-tasks/rules/ui-placement.md:9` 使用从未定义的 `<HAS_UI>` / `<UI_ONLY>` 条件宏。
- **A5**: `breakdown-tasks/rules/existing-code-split.md:45` 维护笔记引用不存在的 "Step 5: Task Dependencies"。
- **A6**: `scope-to-surface-key.md:54` 错误恢复引用 `forge version` 命令（已验证存在，但为 hidden 命令）。
- **A7** [G1]: `quick-tasks/SKILL.md:54,169` 和 `breakdown-tasks/SKILL.md:173` 引用 `.forge/config.yaml` 的 `interfaces` 字段。Config 结构体中**不存在此字段**，实际字段是 `surfaces`。
- **A8** [G2]: `gen-test-scripts/SKILL.md:283` 引用 `.forge/config.yaml` 的 `test-template-dir` 字段。Config 结构体中**不存在此字段**。
- **A9** [G3]: `run-tasks.md:54` / `execute-task.md:25` 从 `forge task claim` 输出提取 `SCOPE` 字段。Go 代码已将 `SCOPE` 替换为 `SURFACE_KEY`/`SURFACE_TYPE`，输出合约测试显式验证 `SCOPE` 不在输出中。
- **A10** [G6]: `fix-bug.md:233`、`write-prd/rules/knowledge-extraction.md:22`、`tech-design/rules/knowledge-extraction.md:22` 引用 `decision-logging.md` 模板文件。该文件**不存在**于任何位置。实际文件是 `learn/templates/decision-entry.md`。

#### B. Surface Resolution 与类型一致性（High）

- **B1**: quick-tasks 和 breakdown-tasks 均对每个 task 的每个 file 调用 `forge surfaces --json`，单 surface 项目造成 N×M 次冗余调用。
- **B2**: 单 surface 项目 agent 将 `key: "."` 视为占位符，surface-type 留空。
- **B3**: `task-doc.md` 模板没有 `surface-key`/`surface-type` 字段，与 `task.md` 不对称。SKILL.md 没有解释这一差异。
- **B4**: `quality_gate.go:352` inferSurface 仅使用第一个提取的 source file，多 surface 场景可能归错 surface。
- **B5** [G4]: `gen-test-scripts/rules/step-0.5-validation.md:20` 使用 `webui` 作为 surface type，但规范名称是 `web`（Go 代码 `KnownSurfaceTypes`、gen-journeys 均使用 `web`）。
- **B6** [G7]: `submit-task/data/record-format-doc.md:3` 列出不存在的 `doc.eval` 类型，遗漏有效的 `doc.review` 类型。
- **B7** [G8]: `forgeconfig/config.go:115` 覆盖率默认值引用 `"coding.clean"` 类型，实际类型是 `"code-quality.simplify"`。该条目永远不会匹配任何 task。
- **B8** [G9]: `task/infer.go` 中 `InferType` 函数缺少 `T-eval-journey` 和 `T-eval-contract` 的处理，返回空字符串。
- **B9** [G16]: `web` surface type 映射到 `gen-test-scripts/types/ui.md`，命名不对称未解释。

#### C. Dispatcher 输出与管道逻辑（Critical/High）

- **C1**: `run-tasks.md:117` 循环后消息硬编码 `T-test-run`/`T-test-verify-regression`，quick mode 中不存在这些任务。
- **C2**: `run-tasks.md:37,109` "print summary" 无定义格式。
- **C3**: `quick.md:151,165` 声称 run-tasks 有 "knowledge extraction"，但 `run-tasks.md` 中没有此步骤。契约不匹配。
- **C4**: `execute-task.md:104` "Output your final summary" 无格式定义。
- **C5**: `execute-task.md:66` 将所有 STATUS≠completed 合并处理，缺少 `in_progress`→record-missing recovery 路径。
- **C6**: `execute-task.md:113` Agent 调用遗漏 `subagent_type="forge:task-executor"`。
- **C7**: `execute-task.md:35` 和 `run-tasks.md:62-65` MAIN_SESSION 缺失 instructions 处理不一致。
- **C8**: `run-tasks.md:110` / `execute-task.md:112` Agent timeout "Mark blocked" 未指定机制。
- **C9**: `task-executor.md:62,68` DONE 格式字段位置不一致，commit-hash 和 status 在同一位置无法区分。
- **C10** [G5]: `gen-contracts/SKILL.md:46-52` 有 SKIP_EVAL_GATE 模式，但 `gen-test-scripts/SKILL.md:29` 没有。Quick mode 在 gen-test-scripts 会被 eval gate 阻塞。
- **C11** [G10]: `run-tests/SKILL.md:66` 将 task 文件路径传给 `forge surfaces --json`，但 CLI 期望源码目录路径（使用 segment prefix matching），task 文件路径几乎必然匹配失败。
- **C12** [G11]: `run-tests/SKILL.md:187` 引用 "Convention loaded in Step 0"，但 Step 0 是 Stale State Recovery，从未加载 Convention。

#### D. Fix-Task 模板变量缺失与质量（High）

- **D1**: 所有 fix-task 创建点（共 10+ 处）均未传递 `--var SOURCE_FILES/TEST_SCRIPT/TEST_RESULTS`。ApplyVars 报错。
- **D2**: Fix-task 按问题类型分组而非按测试套件分组。
- **D3**: cleanup-task 始终 `Breaking: true`，lint/format 修复不应触发完整测试门禁。
- **D4**: `quality_gate.go` EstimatedTime 与 `coding.cleanup.md` frontmatter 不一致 — **双源真相**。
- **D5**: `addFixTask` 跳过 `BuildIndex`，.md 和 index.json 不一致时无法调和。
- **D6**: Breaking task 缺少集成测试影响评估。
- **D7**: `submit-task/SKILL.md:45` record format 文件路径无法从 subagent 工作目录解析 — **Critical**。

#### E. 模板与指引不完整（Medium）

- **E1**: `manifest-quick.md` 双 slug 占位符 `{{FEATURE_SLUG}}` vs `{{SLUG}}`。
- **E2**: `manifest-quick.md:17` 引用不存在的 `testing/test-cases.md`。
- **E3**: 6 个模板占位符在 SKILL.md 中无映射说明。
- **E4**: quick-tasks Output Checklist 提及 stage-gate 文件（永远不会为 quick mode 生成）。
- **E5**: quick-tasks Step 5 声称自动生成 stage-gate 文件（对 quick mode 不成立）。
- **E6**: breakdown-tasks phase-inventory.json 写入有条件但检查无条件。
- **E7**: `coding.fix` 在 quality-gate 表但不在 Type Assignment 表。
- **E8**: breakdown-tasks 无 Commit 步骤。
- **E9**: `{{REFERENCE_FILES}}` 被称为 "template placeholder" 但无模板包含此占位符。
- **E10** [G13]: `run-tests/SKILL.md:122,165` 包含中文文本（"编排序"、"测试环境异常"），其余为英文。
- **E11** [G14]: `run-tests/SKILL.md:84,228` 引用 `BIZ-error-reporting-001`，未提供解析路径。
- **E12** [G15]: `gen-journeys/SKILL.md:59` CLI 错误消息示例与实际输出不匹配。
- **E13** [G17]: `submit-task/data/record-format-gate.md` 遗漏 `acceptanceCriteria`，与 SKILL.md "recommended" 声明不一致。
- **E14** [G12]: `journey-contract-model.md` 在 gen-journeys 和 gen-contracts 各有独立副本，可能分叉。

#### F. 低优先级问题（Low）

- **F1**: `execute-task.md` 疑似死代码（自动化流水线未调用）。
- **F2**: `task-executor.md:19` "Maximum 3 subagent calls" 无强制机制。
- **F3**: `task-executor.md:21` "STEP N DONE" 格式无消费者。
- **F4**: record-missing failure counter 递增/重置语义模糊。
- **F5**: `execute-task.md:115` fix-task 遗漏 `--description`。
- **F6**: `execute-task.md` allowed-tools 含 `TaskOutput` 但未使用，`run-tasks.md` 禁止。
- **F7**: Docs-only fast path 和 Type Assignment 表跨 skill 重复。
- **F8**: `doc.consolidate`/`doc.drift` 描述为 "user manually creates" 但出现在 Type Assignment 表。
- **F9** [H6]: `clean-code/SKILL.md:162` 引用 `just test` 但标准 recipe 是 `unit-test`。
- **F10** [H7]: `fix-bug.md:141,143,146,179` 引用 `just test <slug>` 不存在的 target。
- **F11** [H10]: `ARCHITECTURE.md:148,450` 重复拼写 `forge forge task claim`。
- **F12** [H11]: `forge-distribution.md:105-106` 引用已吸收的 `/record-decision` 和 `/learn-lesson`。
- **F13** [H12]: `prompt-template-hierarchy.md` 声称 "三级" 标签体系，但 `<HARD-RULE>` 是第四级未记录的标签。
- **F14** [H14]: `gen-test-scripts` 混用 "interface"/"surface" 术语。
- **F15** [H15]: `prompt.go:136` TASK_CATEGORY 注入使用脆弱的精确字符串匹配，失败时无错误提示。
- **F16** [H16]: `prompt.go:232` 向后兼容 InferType 包装器传 nil surfaces，禁用 surface-key 前缀匹配。

#### G. Go 代码类型缺口与架构文档虚构（High）

- **G1** [H1]: autogen.go 生成 surface-suffixed 类型（如 `test.gen-scripts.cli`），但 `ValidTypes` 不含这些变体。`prompt.Synthesize` 在 line 89 会拒绝它们。
- **G2** [H2]: `gen-contracts/SKILL.md:61,64,70` 和 `rules/validation.md:36` 使用 "interfaces" 术语和 `interfaces` config 字段引用，实际应为 "surfaces"。
- **G3** [H3]: `ARCHITECTURE.md:119,176-212` 描述 doc-scorer/doc-reviser 为独立 agent，但它们是 eval skill 内的 protocol 文件。agents/ 目录仅有 task-executor.md。文档说 "1 个专用 Agent" 却描述了 3 个。
- **G4** [H4]: `ARCHITECTURE.md:244-256` scope 解析算法描述 `just project-type` 检测流程，实际 `ResolveScope` 用 `just --dry-run compile <scope>` 探测。文档描述的算法不存在于代码中。
- **G5** [H5]: `docs/conventions/dispatcher-quality.md:12,34` 引用 `go build ./...` 和 `go test ./...`（Go 特定），实际代码使用 `just compile`/`just test` 抽象层。
- **G6** [H8]: `docs/conventions/surface-rules.md` 描述 surface-key 传播链，但 `submit.go` 仍用 legacy `scope` 字段（有 TODO 注释）。
- **G7** [H9]: `docs/conventions/dispatcher-quality.md:29` 只提 `coding.fix` 类型，遗漏 `coding.cleanup`（用于 fmt/lint 失败）。
- **G8** [H17]: `forge-cli/pkg/prompt/data/test-run.md:11,31` 引用 `forge:run-e2e-tests` skill，该 skill **不存在**。实际 skill 是 `forge:run-tests`。**test.run 任务执行时直接失败。** Critical。
- **G9** [H18]: `submit.go:157` `validateQualityGate` 硬编码 `scope=""` 忽略 `t.SurfaceKey`。多 surface 项目中 scoped recipe（如 `just compile backend`）从不在 task 提交时执行。有 TODO 注释确认是已知缺陷。
- **G10** [H19]: 17 个 prompt 模板（`prompt/data/*.md`）使用 `SCOPE: {{SURFACE_KEY}}` 标签。标签名 `SCOPE` 误导开发者以为使用 deprecated scope 字段。应改为 `SURFACE:` 或 `SURFACE_KEY:`。
- **G11** [H20]: `build.go:134,307` 仍从 frontmatter 填充 legacy `Scope` 字段到 index.json。导致 `CheckLegacyScope` 迁移检测自我维持——只要 .md 文件含 `scope:` key 就永远无法清除。
- **G12** [H21]: `execute-task.md:3` frontmatter description 声称 "focused TDD workflow"，但实际是 claim/dispatch/verify 编排器，TDD 逻辑在 task-executor 内。

### Urgency

多个 Critical 级别问题直接影响流水线执行：
- A7 (`interfaces` 幽灵字段) → agent 查找不存在的 config 字段
- A9 (`SCOPE` 过时) → claim 输出解析失败
- A10 (`decision-logging.md`) → 知识提取引用不存在的文件
- G2 (gen-contracts `interfaces`) → agent 查找错误的 config 字段
- D7 (submit-task 路径) → subagent 无法找到 record format 文件
- D1 (模板变量) → fix-task 创建在渲染时失败
- G1 (ValidTypes 缺口) → surface-suffixed 测试任务的 prompt 合成失败
- C5 (execute-task 状态处理) → in_progress 被错误创建 fix-task
- C10 (SKIP_EVAL_GATE 缺口) → Quick mode 在 gen-test-scripts 被阻塞
- G3 (虚构 agent) → 架构文档误导开发者
- **H17** (test-run 幽灵 skill) → **test.run 任务执行时调用不存在的 `forge:run-e2e-tests`，直接失败**
- **H18** (submit.go scope 硬编码) → 多 surface 项目 task 提交时 scoped recipe 从不执行
- G1 (ValidTypes 缺口) → surface-suffixed 测试任务的 prompt 合成失败
- C5 (execute-task 状态处理) → in_progress 被错误创建 fix-task
- C10 (SKIP_EVAL_GATE 缺口) → Quick mode 在 gen-test-scripts 被阻塞
- G3 (虚构 agent) → 架构文档误导开发者

## Proposed Solution

按 8 个集群修复，覆盖全部 73 个问题：

1. **幽灵字段与过时引用清理** (A1-A10): 删除废弃文件、修复 `interfaces`/`SCOPE`/`decision-logging.md` 引用
2. **Surface 与类型一致性** (B1-B9): 两层 resolution、修复 surface type 命名、清理幽灵类型
3. **Dispatcher 与管道逻辑对齐** (C1-C12): 条件化消息、定义缺失格式、修复 eval gate 缺口
4. **Fix-Task 模板变量与质量** (D1-D2): 所有创建点传递模板变量、按套件拆分
5. **Quality-Gate Go 代码修复** (D3-D6): cleanup-task 参数、双源真相、inferSurface
6. **submit-task 路径修复** (D7, B6, E13): record format 路径解析、doc 类型修正
7. **Go 代码类型缺口与架构文档修复** (G1-G7): ValidTypes 补全、虚构 agent 描述、算法描述对齐
8. **模板与指引补全** (E1-E14, F1-F16): 占位符映射、误导性声明、中英混合、just target 引用

### Innovation Highlights

三个值得注意的设计改进：
- 两层 surface resolution: 项目级快捷路径消除冗余调用。
- Fix-task 按套件拆分: 可并行、范围有界。
- 双源真相消除: 统一 frontmatter 字段的权威来源。

## Requirements Analysis

### Key Scenarios

- 单 surface 项目: 零次 `forge surfaces` 调用，surface-type 自动填充
- 多 surface 项目: 路径前缀优先，仅歧义文件查询 `forge surfaces`
- quick-mode 完成: 循环后消息准确反映状态，无硬编码任务名
- Quick mode 测试管道: SKIP_EVAL_GATE 在 gen-contracts 和 gen-test-scripts 均生效
- Fix-task 创建: 模板变量全部填充，按测试套件拆分
- Lint 失败: cleanup-task 使用 `Breaking: false`、`EstimatedTime: "15min"`
- submit-task: record format 文件路径可正确解析
- Breaking task: 描述包含集成测试 fixture 影响评估
- Skill 文档中 config 字段引用均指向实际存在的字段

### Constraints & Dependencies

- 变更必须保持与现有 task frontmatter 格式的向后兼容
- Go 代码变更限于 `quality-gate.go`、`template.go`、`config.go`、`infer.go`
- 双源真相修复需决定 frontmatter 字段权威来源
- SKILL.md / command 模板变更必须保持 agent 可解析的散文格式

## Alternatives & Industry Benchmarking

| 方案 | 优点 | 缺点 | 结论 |
|------|------|------|------|
| 不做任何改动 | 零成本 | Critical 问题导致流水线执行失败 | **拒绝** |
| 仅修 Critical/High (~30 个) | 覆盖主要故障 | Medium/Low 问题持续影响 agent 行为 | 部分 |
| **全面修复（全部 73 个问题）** | 完整解决 | 工作量较大 | **选择**: 问题互相关联，部分修复留下不一致 |

## Scope

### In Scope

**集群 1 — 幽灵字段与过时引用清理（~10 个文件）：**
- 删除 `breakdown-tasks/rules/scope-assignment.md`
- 更新 `db-schema.md`: scope → surface-type
- 更新 `existing-code-split.md`: scope assignment → surface inference
- 修复 `ui-placement.md`: 定义或移除 `<HAS_UI>`/`<UI_ONLY>`
- 更新 `scope-to-surface-key.md`: 修正 `forge version` 引用
- 更新 `quick-tasks/SKILL.md` + `breakdown-tasks/SKILL.md`: `interfaces` → `surfaces`
- 更新 `gen-test-scripts/SKILL.md`: 移除 `test-template-dir` 引用
- 更新 `run-tasks.md` + `execute-task.md`: `SCOPE` → `SURFACE_KEY`/`SURFACE_TYPE`
- 更新 `fix-bug.md` + `write-prd/rules/knowledge-extraction.md` + `tech-design/rules/knowledge-extraction.md`: `decision-logging.md` → `decision-entry.md`

**集群 2 — Surface 与类型一致性（~6 个文件 + 2 Go 文件）：**
- `quick-tasks/SKILL.md` + `breakdown-tasks/SKILL.md` + `scope-to-surface-key.md`: 两层 resolution
- `quick-tasks/templates/task-doc.md`: 添加 surface 字段或豁免说明
- `gen-test-scripts/rules/step-0.5-validation.md`: `webui` → `web`
- `submit-task/data/record-format-doc.md`: 移除 `doc.eval`，添加 `doc.review`
- `forgeconfig/config.go`: 移除 `coding.clean` 死代码条目
- `task/infer.go`: 添加 `T-eval-journey` 和 `T-eval-contract` 处理
- `quality_gate.go`: inferSurface 使用所有 source files

**集群 3 — Dispatcher 与管道逻辑对齐（~6 个文件）：**
- `run-tasks.md`: 条件化消息 + 定义摘要格式 + timeout 机制
- `quick.md`: 删除知识提取虚假声明
- `execute-task.md`: 摘要格式 + 状态区分 + subagent_type + MAIN_SESSION 统一
- `task-executor.md`: 统一 DONE 格式
- `gen-test-scripts/SKILL.md`: 添加 SKIP_EVAL_GATE 条件
- `run-tests/SKILL.md`: 修复 surface 检测路径 + Step 0 引用 + 移除中文文本

**集群 4 — Fix-Task 模板变量（~5 个文件）：**
- `run-tasks.md`: 4 个创建点补充 `--var`
- `task-executor.md`: fix-task 创建补充 `--var`
- `execute-task.md`: fix-task 创建补充 `--var` + `--description`
- `submit-task/SKILL.md`: recovery 补充 `--var`
- `quick-tasks/SKILL.md` + `breakdown-tasks/SKILL.md`: breaking task IT 影响评估指引

**集群 5 — Quality-Gate Go 代码修复（3 个文件）：**
- `quality_gate.go`: fix-task 按套件拆分 + cleanup-task 参数修正 + 双源真相修复
- `template.go`: cleanup-task `Breaking: false`
- `coding.cleanup.md`: frontmatter `breaking: false`

**集群 6 — submit-task 修复（1 个文件）：**
- `submit-task/SKILL.md`: record format 路径改为 `plugins/forge/skills/submit-task/data/record-format-{TASK_CATEGORY}.md` 或等效可解析路径

**集群 7 — Go 代码类型缺口与架构文档修复（~8 个文件）：**
- `types.go` ValidTypes: 添加 surface-suffixed 类型变体或在 Synthesize 中跳过校验
- `prompt/data/test-run.md`: `forge:run-e2e-tests` → `forge:run-tests`（**Critical 修复**）
- `submit.go:157`: 将 `scope=""` 改为 `scope=t.SurfaceKey`
- `prompt/data/*.md` (17 个模板): `SCOPE:` 标签 → `SURFACE_KEY:` 或 `SURFACE:`
- `gen-contracts/SKILL.md` + `rules/validation.md`: "interfaces" → "surfaces"
- `docs/ARCHITECTURE.md`: (a) 移除虚构的 doc-scorer/doc-reviser agent 描述 (b) 修正 scope 解析算法描述 (c) 修正 `forge forge` 重复拼写 (d) `fmt` failure 描述修正为 non-blocking
- `docs/conventions/dispatcher-quality.md`: `go build`/`go test` → `just compile`/`just test`，补充 `coding.cleanup`
- `clean-code/SKILL.md`: `just test` → `just unit-test`
- `fix-bug.md`: `just test <slug>` → 正确的 target
- `execute-task.md:3`: description 修正为 claim/dispatch/verify

**集群 8 — 模板与指引补全（~8 个文件）：**
- `manifest-quick.md`: 统一 slug 占位符 + 处理 test-cases.md 引用
- `quick-tasks/SKILL.md`: 占位符映射 + stage-gate 修正 + REFERENCE_FILES 修正
- `breakdown-tasks/SKILL.md`: Commit 步骤 + phase-inventory 条件化
- `gen-journeys/SKILL.md`: 错误消息示例对齐
- `run-tests/SKILL.md`: BIZ-error-reporting-001 添加解析路径 + 移除中文文本
- `forge-distribution.md`: 更新 `/record-decision`/`/learn-lesson` 引用为 `/learn`
- `prompt-template-hierarchy.md`: 记录 `<HARD-RULE>` 第四级标签
- 考虑 `journey-contract-model.md` 统一为单副本

### Out of Scope

- 新增 CLI 命令或标志
- 修改 `forge surfaces` 命令本身
- 重构 tech-design/write-prd/fix-bug 中重复的知识提取规则
- 删除 execute-task.md（保留作为手动入口）
- enforce "Maximum 3 subagent calls"（需要运行时计数器）

## Key Risks

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|----------|
| 两层 resolution 遗漏边界情况 | 低 | 低 | 路径前缀回退到 `forge surfaces` |
| Fix-task 拆分过度碎片化 | 中 | 中 | 底线规则：同目录保持单个 |
| 双源真相修复引入新不一致 | 中 | 高 | opts 为权威来源，frontmatter 用 `{{...}}` |
| cleanup-task Breaking:false 遗漏回归 | 低 | 中 | compile+lint 仍运行 |
| gen-test-scripts SKIP_EVAL_GATE 过度放宽 | 低 | 中 | 仅在 Quick mode 生效，full mode 不受影响 |
| 大范围 SKILL.md 修改被误解 | 中 | 中 | 显式条件块 + 具体示例 |

## Success Criteria

- [ ] 所有 skill 文档中的 config 字段引用指向实际存在的字段（无 `interfaces`、`test-template-dir`）
- [ ] gen-contracts 和 gen-test-scripts 使用 "surfaces" 术语
- [ ] `scope-assignment.md` 已删除，`db-schema.md` 和 `existing-code-split.md` 使用 surface 字段
- [ ] `forge task claim` 输出字段引用正确（`SURFACE_KEY`/`SURFACE_TYPE`，无 `SCOPE`）
- [ ] `decision-logging.md` 引用已替换为 `decision-entry.md`
- [ ] 单 surface 项目零次 `forge surfaces` 调用，非空 `surface-type`
- [ ] surface type 一致使用规范名称（`web`，非 `webui`）
- [ ] `record-format-doc.md` 不含幽灵类型
- [ ] `coding.clean` 死代码已清理
- [ ] `InferType` 处理所有自动生成的 task ID（含 `T-eval-*`）
- [ ] ValidTypes 包含 surface-suffixed 变体或 Synthesize 跳过校验
- [ ] test-run.md 引用 `forge:run-tests`（非 `forge:run-e2e-tests`）
- [ ] submit.go 使用 `t.SurfaceKey` 传递 scope 参数
- [ ] prompt 模板使用 `SURFACE_KEY:` 标签（非 `SCOPE:`）
- [ ] 循环后消息基于实际状态，摘要格式已定义
- [ ] quick.md 不声明 run-tasks 未执行的知识提取
- [ ] execute-task.md 正确区分 blocked/in_progress
- [ ] Quick mode 测试管道 SKIP_EVAL_GATE 在 gen-contracts 和 gen-test-scripts 均生效
- [ ] 所有 fix-task 创建点传递 `--var SOURCE_FILES/TEST_SCRIPT/TEST_RESULTS`
- [ ] cleanup-task 使用 `Breaking: false` 和 `EstimatedTime: "15min"`
- [ ] quality gate fix-task index.json 和 .md 字段一致
- [ ] submit-task record format 文件路径可正确解析
- [ ] manifest-quick.md 使用统一 slug 占位符
- [ ] quick-tasks SKILL.md 无 stage-gate 误导性声明
- [ ] ARCHITECTURE.md 不描述虚构的 agent，scope 解析算法与代码一致
- [ ] dispatcher-quality.md 使用 `just` 抽象层而非 Go 特定命令，提及 `coding.cleanup`
- [ ] clean-code 和 fix-bug 中 `just test` 引用修正为 `just unit-test`
- [ ] 所有现有集成测试通过
- [ ] 所有现有集成测试通过
