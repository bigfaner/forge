---
created: 2026-05-24
author: "faner"
status: Draft
---

# Proposal: Pipeline Integration Stitch — 修复 auto-gen-journeys-contracts 提案遗留的集成缝隙

## Problem

`auto-gen-journeys-contracts` 提案引入了 staged test pipeline（gen-journeys → gen-contracts → gen-scripts → run → verify-regression）和 eval 质量门控（eval-journey、eval-contract），但遗漏了执行阶段的配套工作：**4 个 prompt 模板文件缺失**、**eval 类型分类错误**、**依赖注入逻辑引用已废弃类型**、**gen-and-run 废弃代码残留**。

### Evidence

#### P0 — Pipeline 执行必定失败

1. **`prompt/data/` 缺少 4 个执行阶段模板文件**
   - 缺失：`test-gen-journeys.md`、`test-gen-contracts.md`、`eval-journey.md`、`eval-contract.md`
   - `task/data/` 中已有这 4 个文件（autogen 规划阶段模板），但 `prompt/data/`（执行阶段模板）未创建
   - `Synthesize()` 在渲染这些类型的任务 prompt 时 `ReadFile` 失败
   - 影响：所有使用 Forge 执行 test pipeline 或 eval gate 的 feature 必定失败

#### P1 — 特定场景失败

2. **`eval.*` 类型落入 CategoryCoding**（`category.go:16-33`）
   - `CategoryForType()` 无 `eval.` 前缀分支，`eval.journey`/`eval.contract` 落入 default → `CategoryCoding`
   - 影响：eval 任务被要求提供测试证据（testsPassed/coverage），但 eval 是 review 类任务
   - 影响面：`validateRecordData`（submit.go:296）、`RenderRecord`（record.go:235）、prompt 注入（prompt.go `renderTemplate`）

3. **`findFirstTestTaskIdx` quick-mode 分支匹配废弃类型**（`build.go:492-494`）
   - 查找 `T-quick-gen-and-run*`，但 Quick 模式已不再生成该类型（`GetQuickTestTasks` 使用 staged tasks）
   - 当前靠 `return 0` fallback 意外正确（gen-journeys 恰好排首位），但脆弱

4. **T-review-doc prepend 与 ResolveFirstTestDep 顺序耦合**（`build.go:329-347`）
   - 两步操作顺序硬编码：先 ResolveFirstTestDep（设置基础 deps），再 prepend T-review-doc
   - 虽然当前幂等（ResolveFirstTestDep 总是全新覆写），但重排顺序会导致 T-review-doc 丢失
   - 应合并为单步操作消除耦合

#### P2 — 维护风险

5. **`test.gen-and-run` 废弃代码残留**（生产代码 + 测试 + 活跃文档）
   - `types.go`: `TypeTestGenAndRun` 常量、ValidTypes、isTestTaskID 条目
   - `infer.go:32-33`: gen-and-run 推断分支
   - `prompt.go:297,304`: `genScriptBases` 中 `T-quick-gen-and-run` 条目
   - `validate_index.go:224`: `T-quick-gen-and-run-` 前缀检查
   - `build.go:484,492,494`: findFirstTestTaskIdx quick-mode 注释和匹配
   - `prompt/data/test-gen-and-run.md`、`task/data/test-gen-and-run.md`: 废弃模板文件
   - 14 个测试文件中 ~95 处引用
   - 活跃文档引用：OVERVIEW.md、task-lifecycle.md

6. **record-format 参考文档过期**
   - `record-format-test.md` 列出已废弃类型（`test.gen-cases`/`test.eval-cases`/`test.gen-and-run`），缺少新类型
   - 缺少 `record-format-eval.md`：agent 执行 eval 任务时无 JSON 字段参考

### Urgency

P0 意味着 `forge prompt get-by-task-id` 对 `test.gen-journeys`、`test.gen-contracts`、`eval.journey`、`eval.contract` 返回错误，任何执行 test pipeline 或 eval gate 的 feature 必定失败。

## Proposed Solution

三管齐下：

1. **根治 P0**：创建 4 个执行阶段 prompt 模板文件
2. **修复 P1**：新增 CategoryEval + 完整验证/记录分支 + 加固依赖注入为单步操作
3. **清理 P2**：从生产代码、测试文件和活跃文档中移除 gen-and-run 废弃代码

### Innovation Highlights

Task 1（已完成）引入的**自动发现机制**已消除"忘记更新映射"的根因。本次提案聚焦于补全遗漏的**执行层面配套**（模板文件、类型分类、记录格式），使 staged test pipeline 和 eval gate 从"类型注册完成"推进到"端到端可执行"。本质上这是维护性工作——补全 Task 1 遗漏的 adapter 层实现，不引入新的架构模式或设计范式。

**防御性诊断层**：`CategoryForType` default 分支添加 `log.Printf` 警告是一个低成本的**防御性诊断辅助**——它不改变当前行为（仍返回 `CategoryCoding`），但使未来任何新类型的误分类立即产生可观测信号（日志中出现 `CategoryForType: unknown type` 即表示遗漏了分类分支）。这不是注册机制本身（实际的分类注册是 `eval.` 前缀分支），而是为分类遗漏提供运行时诊断能力。

## Requirements Analysis

### Key Scenarios

- **新类型零配置**: 在 `types.go` 添加常量 + 在 `data/` 放入模板文件，auto-discovery 自动识别（已由 Task 1 完成）
- **Eval 提交语义正确性**: `forge submit-task` 对 eval 任务接受 review 字段（summary/findings/severity），拒绝纯测试字段
- **Mixed feature 依赖注入**: T-review-doc 正确插入为 test pipeline 前置依赖，单步操作无顺序耦合
- **Quick-mode findFirstTestTaskIdx**: 正确匹配新 staged tasks 前缀（T-test-gen-journeys）
- **Gen-and-run 完全移除**: 生产代码零残留，测试和活跃文档零残留

### Non-Functional Requirements

- **向后兼容**: 旧 index.json 引用 `test.gen-and-run` 时给出明确迁移错误提示（覆盖 `validate_index.go` 验证路径和 `Synthesize()` 渲染路径）
- **CategoryEval 测试覆盖**: 正向用例（接受 review 字段）、负向用例（拒绝纯测试字段）、边界用例（混合提交）
- **eval record 模板**: 包含 score、findings、severity 等 eval 特有字段
- **性能**: CategoryEval 分支在 `RenderRecord`/`validateRecordData` 热路径上仅增加一次 string prefix 比较（`strings.HasPrefix(typ, "eval.")`），与现有 CategoryTest 分支开销一致；新增 eval record 模板为内存字符串查找，无 I/O 开销。预计影响 < 1μs，无需基准测试
- **版本范围**: 本次变更影响的 Forge 版本为引入 eval 类型的版本（即 Task 1 合并后的版本）。旧版本 index.json 不包含 eval 类型，因此 CategoryEval 分支不会被触发，性能影响为零。validate_index.go 的迁移错误仅针对旧版本中可能存在的 `test.gen-and-run` 引用
- **可观测性（显式排除）**: eval 任务的日志/指标/监控不在本次范围内。本次仅确保 eval 类型在类型系统中正确分类和渲染，不涉及 eval 执行结果的结构化日志或告警。eval task outcome observability 应作为独立提案处理，因为它需要设计日志格式、指标定义和告警策略，超出 adapter 层补全的范围

### Constraints & Dependencies

- Task 1（auto-discovery + init-time 校验 + clean-code.md 重命名）已完成，本次所有工作基于该基础

## Alternatives & Industry Benchmarking

### Industry Solutions

这是典型的**补全遗漏的 adapter 层**问题。在新类型系统中，注册（types.go）和发现（auto-discovery）已完成，但执行阶段的模板、分类和记录渲染未同步。

**行业参照**：

1. **Spring Boot auto-configuration**：Spring 通过 `@Configuration` + `@ConditionalOnClass` 实现约定优于配置——只需在 classpath 放入正确类，运行时自动发现并装配。Forge 的 auto-discovery 机制同理：`types.go` 注册类型 + `data/` 放入模板 = 运行时自动可用。本次 P0 模板缺失正是 Spring 所防范的"注册了 bean 但缺少自动装配配置"问题。

2. **ASP.NET Core convention-based endpoint discovery**：ASP.NET Core 从 assembly 中按命名约定（`*Controller`）自动发现并注册 endpoint。Forge Task 1 的 auto-discovery 采用相同模式（文件名 → task type 映射）。本次遗留问题等同于"注册了 route template 但缺少对应的 controller action"。

3. **Temporal workflow/activity registration**：Temporal 要求 workflow type 和 activity type 都必须注册到 worker，否则在执行时返回 `NotFoundError`。Forge 的 `Synthesize()` 对缺失模板的 `ReadFile` 失败与此完全一致——类型存在但执行适配器缺失。

类比：Airflow 添加新 DAG 类型后需要配套的 executor plugin 和 UI renderer。

### Comparison Table

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| 手动补全模板 + 最小化修复 | 变更最小 | eval 分类仍错；gen-and-run 僵尸代码持续积累 | Rejected |
| 仅修 P0+P1，保留 gen-and-run | 减少变更量 | 废弃代码干扰开发和测试 | Rejected |
| **完整修复 P0+P1+P2** | 端到端可执行；零僵尸代码；类型系统一致 | 变更量较大——14 文件 ~95 处测试引用需逐文件编辑，预估 P2 单独耗时 2h+ | **Selected** |
| **Code-gen：从类型定义自动生成 adapter 代码** | 消除手动遗漏的根因；与 auto-discovery 形成闭环 | 引入代码生成工具链依赖；现有代码库未使用 code-gen 模式；ROI 不合理——仅模板层面值得自动化的部分为 4 个文件（不值得建生成器），而 CategoryEval、RecordData 字段、RenderRecord case 等 adapter 逻辑涉及 Go 类型系统和业务语义，无法通过模板化代码生成覆盖 | Rejected |
| **Schema-based init-time validation** | 启动时即发现缺失模板而非运行时崩溃 | 需设计 schema 格式；Task 1 已通过 init-time 校验覆盖部分场景；本次 P0 问题本质是遗漏创建文件而非验证缺失 | Rejected |

**选择理由**：完整修复与行业最佳实践一致（Spring/Temporal 的完整注册模型），同时避免引入 code-gen 的工具链复杂度（ROI 不合理）和 schema validation 的过度工程化。P2 的 ~95 处引用每处修改模式相同（删除引用 → 验证编译），风险可控但 2h+ 的耗时不可低估。

## Feasibility Assessment

### Technical Feasibility

- 4 个 prompt 模板：test-gen 模板参考 `test-gen-scripts.md`（skill 委托模式，~45 行），eval 模板不适用 `validation-code.md`（该模板为代码质量验证，含 compile/fmt/lint/test 流程，语义不同）——eval 模板需新建独立模式：quality evaluation 角色定义 + `forge:eval` skill 委托 + score/findings/severity/passed record 字段
- CategoryEval：参考 CategoryTest 的实现模式（分类常量 + CategoryForType 分支 + record 模板 + submit 验证）
- 依赖注入加固：将 ResolveFirstTestDep + T-review-doc prepend 合并为单函数调用
- Gen-and-run 移除：机械性操作，按文件清单逐一清理

### Resource & Timeline

预计 4 个 coding task（含集成测试）+ 2 个 doc task，总工作量 ~7h。集成测试覆盖 Quick mode 依赖链，预计增加 ~1h 工作量（需要在测试上下文中创建 Quick mode feature 并断言依赖链）。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| re-index 幂等性是实际 bug | Code Audit | Overturned: 当前代码已幂等（ResolveFirstTestDep 总是全新覆写），风险仅为代码耦合 |
| gen-and-run 清理应包含历史 feature/proposal 文档 | Occam's Razor | Refined: 历史文档不影响运行时，仅清理活跃文档 |
| eval 需要 CategoryTest 分类 | Assumption Flip | Overturned: eval 是 review 类任务（非测试生成器），需新建 CategoryEval |

## Scope

### In Scope

**P0 — 执行阶段模板**
- 创建 `prompt/data/test-gen-journeys.md`（agent 执行 test.gen-journeys 任务的指令模板）
- 创建 `prompt/data/test-gen-contracts.md`（agent 执行 test.gen-contracts 任务的指令模板）
- 创建 `prompt/data/eval-journey.md`（agent 执行 eval.journey 任务的指令模板）
- 创建 `prompt/data/eval-contract.md`（agent 执行 eval.contract 任务的指令模板）

模板内容模式：

**test-gen 类型**（test-gen-journeys.md、test-gen-contracts.md）遵循现有 test-gen-scripts.md 模式：
- Header: `TASK_ID`/`TASK_FILE`/`SCOPE`/`PHASE_SUMMARY` 占位符
- Role: "focused task executor running a [journey|contract] generation task"
- Task Constraints: MUST invoke 对应 skill（`forge:gen-journeys` / `forge:gen-contracts`），MUST NOT 手动写文件
- Workflow 2-step: Read Task Definition → Invoke Skill
- Record Fields: scriptsCreated, casesGenerated

**eval 类型**（eval-journey.md、eval-contract.md）为全新模式，与 test-gen 根本不同：
- Header: 同样使用 `TASK_ID`/`TASK_FILE`/`SCOPE`/`PHASE_SUMMARY` 占位符
- Role: "focused task executor running a quality evaluation task"
- Task Constraints: MUST invoke `forge:eval` skill，MUST NOT 修改被评估的文件
- Workflow 2-step: Read Task Definition → Invoke Eval Skill
- Eval 执行委托：`Skill(skill="forge:eval", args="--type [journey|contract] --target 850")`。委托链说明：`forge:eval` 是统一的 eval 入口 skill，根据 `--type` 参数内部分发到 `forge:eval-journey`（journey 评估）或 `forge:eval-contract`（contract 评估）的具体评估逻辑。`--target 850` 来源：Forge eval 系统使用 1000-point 评分制（见 `plugins/forge/skills/eval/rubrics/`），850/1000 为 PRD 和 proposal 类别的默认通过阈值。当前硬编码为 850 是合理的——eval-journey/eval-contract 均使用同一评分制和通过标准。若未来需要不同阈值，应在 eval skill 的命令定义中参数化，而非在 prompt 模板中暴露
- Record Fields: score, findings, severity, passed（eval 特有字段）

**P1 — CategoryEval + 依赖加固**
- `category.go`: 新增 `CategoryEval = "eval"`，`CategoryForType` 添加 `eval.` 前缀分支。default 分支改为返回 `CategoryCoding` 并通过 `log.Printf("CategoryForType: unknown type %q, defaulting to coding", typ)` 记录警告，使未来新类型的误分类可被日志追踪而非静默
- `submit.go`: `validateRecordData` 为 CategoryEval 添加验证分支（接受 review 字段 summary/findings/severity，拒绝纯测试字段）
- `types.go`: RecordData 结构添加 eval 特有字段，遵循现有无前缀命名惯例和 JSON tag 规范：
  - `Score float64 \`json:"score,omitempty"\`` — eval 评分（0-1000）
  - `Findings []string \`json:"findings,omitempty"\`` — eval 发现问题列表
  - `Severity string \`json:"severity,omitempty"\`` — eval 问题严重程度（critical/major/minor）
  - `Passed bool \`json:"passed,omitempty"\`` — eval 是否通过质量门控
  - 命名遵循现有惯例：Doc 字段用 `ReferencedDocs`/`ReviewStatus`（无 doc 前缀），Test 字段用 `CasesGenerated`（无 test 前缀），Validation 字段用 `ValidationPassed`（因 Validation 本身已是分类名），Gate 字段用 `GatePassed`（同理）
- `record.go`: 新增 `record-eval.md` Go 模板 + `RenderEvalRecord` 函数 + `RenderRecord` switch 添加 CategoryEval case + `RecordTemplateData` 添加 eval 格式化字段（ScoreFormatted、FindingsFormatted、SeverityFormatted、PassedFormatted）+ `NewRecordTemplateData` 填充这些字段
- `plugins/forge/skills/submit-task/data/record-format-eval.md`: eval 任务 JSON 字段参考文档
- `category_test.go` + `submit_test.go`: CategoryEval 专项测试
- 集成测试（`pkg/task/integration_test.go` 或扩展现有 `build_test.go`）：覆盖 Quick mode 完整依赖链——创建 Quick mode feature → 验证任务列表包含 T-review-doc（当 needsEval=true）→ 验证 gen-journeys 依赖指向 T-review-doc → 验证 needsEval=false 时不包含 T-review-doc 且与旧行为一致
- `build.go`: 将 ResolveFirstTestDep + T-review-doc prepend 合并为单步操作。合并函数签名：`func resolveTestDepsAndInjectReviewDoc(testTasks []AutoGenTaskDef, index *TaskIndex, mode string, needsEval bool)`。当 `needsEval=true` 时，T-review-doc 注入与 first-test-task 依赖解析在同一函数内完成；当 `needsEval=false` 时，仅执行 ResolveFirstTestDep 逻辑
- `build.go`: `findFirstTestTaskIdx` quick-mode 分支替换为 `findTaskIndexByPrefix(tasks, "T-test-gen-journeys")`。**关于 auto-discovery 与前缀匹配的区分**：Task 1 的 auto-discovery 解决的是**模板发现**（文件名 → task type 的自动映射），而此处 `findTaskIndexByPrefix` 解决的是**任务定位**（从已有任务列表中找到特定功能的任务）。两者是不同层次的问题：模板发现应完全自动化（已实现），任务定位需要明确的语义标识符。**为何使用前缀而非类型常量**：此处定位的是 task ID（如 `T-test-gen-journeys-featureX`），而非 task type（如 `test.gen-journeys`）。task ID 由 `GetQuickTestTasks` 按命名约定生成，前缀匹配是 task ID 层面唯一可用的定位手段。与旧方案 `T-quick-gen-and-run` 相比，新前缀 `T-test-gen-journeys` 直接对应当前 staged pipeline 的首个任务类型，语义稳定——除非 staged pipeline 本身被重新设计（此时此处逻辑必然需要同步修改），否则前缀不会失效。这是同一失败类别下的最优权衡：接受前缀耦合但将其锚定在当前架构的稳定语义上

**注意**：`findFirstTestTaskIdx` 修改同时出现在 P1（功能修复）和 P2（gen-and-run 清理，`build.go:484,492-494`）。这是同一处代码的两个不同修改意图：P1 将 quick-mode 匹配从废弃前缀更新为新前缀，P2 清理同一位置的注释和残留代码。实现时应合并为单次编辑，避免冲突。

**P1 — Record 模板参考文档更新**
- `record-format-test.md`: 将类型列表替换为 `test.gen-journeys`、`test.gen-contracts`、`test.gen-scripts`、`test.run`、`test.verify-regression`（移除不存在的 `test.gen-cases`、`test.eval-cases` 和已废弃的 `test.gen-and-run`）

**P2 — gen-and-run 废弃代码移除**

生产代码移除（按此顺序确保增量编译通过）：

1. **`infer.go:32-33`** — 移除 gen-and-run 推断分支（`case id == "T-quick-gen-and-run", typeSuffixedID(id, "T-quick-gen-and-run"): return TypeTestGenAndRun`）。此为引用 TypeTestGenAndRun 常量的消费者，必须先于常量定义移除。
2. **`prompt.go:293-305`** — 移除 `genScriptBases` 中 `"T-quick-gen-and-run"` 条目（line 304）及注释中 line 297 的 gen-and-run 说明。
3. **`validate_index.go:224-226`** — 移除 `T-quick-gen-and-run-` 前缀检查分支，替换为迁移感知错误提示（`"test.gen-and-run is deprecated; use staged test pipeline types (test.gen-journeys, test.gen-contracts, test.gen-scripts)"`）。
4. **`build.go:484,492-494`** — 更新 findFirstTestTaskIdx 注释（line 484）和 quick-mode 匹配逻辑（line 492-494），替换为 `findTaskIndexByPrefix(tasks, "T-test-gen-journeys")`（与 ResolveFirstTestDep 使用相同发现机制）。
5. **`types.go:55,88,114,134`** — 移除 TypeTestGenAndRun 常量定义（line 55）、TaskTypeRegistry 条目（line 88）、ValidTypes 条目（line 114）、SystemTypes 条目（line 134）、isTestTaskID 条目。所有消费者已在前 4 步移除，此步安全删除常量定义。

废弃模板文件删除：
- `prompt/data/test-gen-and-run.md`
- `task/data/test-gen-and-run.md`

测试文件更新（按文件）：
- `pkg/task/types_test.go:576,602,633,693` — 移除 TypeTestGenAndRun 引用
- `pkg/task/category_test.go:31` — 移除 test.gen-and-run 测试用例
- `pkg/task/infer_test.go:22,46,47,104,105` — 移除 gen-and-run 推断测试用例
- `pkg/task/build_test.go:708,710,748,1451` — 移除 gen-and-run 构建测试用例
- `pkg/task/autogen_test.go:106-107,666,1104,1298,1596-1597` — 移除 gen-and-run 自动生成测试用例
- `pkg/task/record_test.go:562,570` — 移除 gen-and-run record 测试用例
- `pkg/task/stage_gates_test.go:32` — 更新 stage gate 测试中的 task ID
- `pkg/prompt/prompt_test.go:637-644,693-694,1676-1680` — 移除 gen-and-run prompt 测试用例
- `internal/cmd/quality_gate_test.go:1129-1130` — 移除 gen-and-run quality gate 测试
- `internal/cmd/task/validate_index_test.go:1480,1483,1632-1676,1703-1708,1768-1770` — 更新验证测试用例
- `internal/cmd/task/list_test.go:271` — 更新列表测试标题
- `internal/cmd/base/output_test.go:12` — 更新输出测试标题
- `tests/task-lifecycle/task_stage_gates_test.go:263` — 更新生命周期测试

活跃文档更新：
- `docs/OVERVIEW.md` — 移除 gen-and-run 引用
- `docs/WORKFLOW.md`（task-lifecycle 等效文档）— 移除 gen-and-run 引用

### Out of Scope

- 历史 feature/proposal 文档中的 gen-and-run 引用（~35 文件 ~130 处，不影响运行时）
- 重构 resolveBreakdownDeps/resolveQuickDeps 的重复逻辑
- eval rollback 改进
- 旧 index.json 自动迁移工具
- record-format-doc.md 中 `doc.eval` → `doc.review`（`doc.eval` 不在运行时使用）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 4 个新 prompt 模板内容不准确 | M | H | test-gen 模板复用 test-gen-scripts.md 的已验证模式（skill 委托 + 2-step workflow + record fields），eval 模板使用全新模式但通过以下方式验证正确性：(1) 模板创建后立即运行 `forge prompt get-by-task-id` 验证渲染输出；(2) 集成测试验证 eval 类型完整流程：模板渲染 → skill 委托 → record 提交。两者均使用 header/constraints/workflow/record-fields 四段式结构 |
| CategoryEval 提交验证字段与实际不匹配 | L | M | 验证分支接受 review 字段（summary/findings/severity），编写单元测试覆盖正向/负向/边界用例 |
| gen-and-run 引用移除不完整导致编译失败 | M | H | 按指定顺序移除（消费者先于常量定义），确保每步增量编译通过：infer.go → prompt.go → validate_index.go → build.go → types.go，每步后执行 `go build ./...` 验证 |
| 删除 `prompt/data/test-gen-and-run.md` 后 Synthesize() 对旧 feature 的 file-not-found 错误无迁移感知 | M | M | 两路径防护：(1) `validate_index.go` 验证路径已覆盖迁移错误提示；(2) 在 `Synthesize()` 的 `ReadFile` 失败分支中，检查文件名包含 `gen-and-run` 时输出迁移指引而非通用 file-not-found 错误。确保两个入口（验证时和渲染时）都给出一致的迁移信息 |
| findFirstTestTaskIdx 修改影响现有 dependency wiring | M | M | 修改后运行 `forge task list` 验证 Quick mode feature 的任务依赖链完整（T-review-doc → T-test-gen-journeys → 后续 staged tasks）；集成测试断言 `needsEval=true` 时 T-review-doc 出现在依赖链中且 gen-journeys 依赖指向 T-review-doc |
| CategoryForType default 分支静默误分类未来新类型 | L | M | default 分支添加 log.Printf 警告，使未知类型的误分类可被日志追踪 |
| eval record 模板字段设计不合理 | L | M | eval 任务为本次新建（无既有实现可参考），字段设计基于 eval skill 的输出契约：`plugins/forge/skills/eval/` 下的 rubric 文件定义了 score（0-1000 数值）、findings（问题列表）、severity（critical/major/minor）、passed（布尔门控结果）四个标准输出维度。record 模板字段直接映射这些维度。若 eval skill 后续调整输出格式，需同步更新 record 模板 |

## Success Criteria

- [ ] `prompt/data/` 包含 test-gen-journeys.md、test-gen-contracts.md、eval-journey.md、eval-contract.md
- [ ] `Synthesize()` 对 `test.gen-journeys`、`test.gen-contracts`、`eval.journey`、`eval.contract` 返回有效 prompt（P0 修复）
- [ ] `CategoryForType("eval.journey")` 返回 `CategoryEval`（非 `CategoryCoding`）
- [ ] `forge submit-task` 对 eval 任务接受含 summary/findings 的提交，拒绝仅含 testsPassed/coverage 的提交
- [ ] `RenderRecord` 对 CategoryEval 使用 eval 专用 record 模板，且渲染输出包含 `ScoreFormatted`、`FindingsFormatted`、`SeverityFormatted`、`PassedFormatted` 字段值。单元测试验证：构造包含 eval 字段的 `RecordData{Score: 850, Findings: []string{"finding1"}, Severity: "major", Passed: true}` 的 CategoryEval record，调用 `RenderRecord` 后输出包含格式化后的 score 字符串、findings 列表、"major" 和 passed 标识
- [ ] `findFirstTestTaskIdx` 对 Quick mode 正确返回 gen-journeys 任务索引
- [ ] `resolveTestDepsAndInjectReviewDoc(testTasks, idx, "quick", true)` 返回的依赖列表包含 T-review-doc；`resolveTestDepsAndInjectReviewDoc(testTasks, idx, "quick", false)` 返回的依赖列表不包含 T-review-doc 且与旧 ResolveFirstTestDep 输出一致
- [ ] `grep -r "gen-and-run\|quick-gen-and-run\|T-quick-gen" forge-cli/ --exclude-dir=docs/proposals` 返回零结果
- [ ] `grep -r "gen-and-run\|quick-gen-and-run\|T-quick-gen" plugins/forge/ --exclude-dir=docs/proposals` 返回零结果
- [ ] validate_index.go 对引用 `test.gen-and-run` 的旧 index.json 返回迁移指引错误信息
- [ ] 所有现有测试通过（`go test ./pkg/task/... ./pkg/prompt/... ./internal/cmd/... ./tests/...`）
- [ ] `plugins/forge/skills/submit-task/data/record-format-eval.md` 存在且包含 score、findings、severity、passed 字段定义
- [ ] 单元测试验证：对 `CategoryForType("unknown.type")` 调用后，日志输出包含 `"CategoryForType: unknown type"` 警告字符串
- [ ] 集成测试覆盖 Quick mode 依赖链：创建 Quick mode feature → 验证生成的任务列表包含 T-review-doc（当 needsEval=true）且 gen-journeys 依赖正确指向 T-review-doc
- [ ] `record-format-test.md` 包含全部五个新类型（`test.gen-journeys`、`test.gen-contracts`、`test.gen-scripts`、`test.run`、`test.verify-regression`）且不包含已废弃类型名（`test.gen-cases`、`test.eval-cases`、`test.gen-and-run`）

## Next Steps

- Proceed to `/quick-tasks` to generate tasks
