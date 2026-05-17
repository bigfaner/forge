---
created: 2026-05-17
author: "faner + Claude"
status: Draft
---

# Proposal: Task Lifecycle & Pipeline Orchestration Root-Cause Hardening

## Problem

Forge v3.0.0-beta-7 的 task state machine 和 pipeline orchestration 存在 **6 个系统性架构缺陷**，导致 **25+ 个具体逻辑缺口**。这些缺口不是孤立的 bug，而是由共通的根因模式驱动的系统性问题——分散的状态验证、混合的原子性模型、脆弱的 eval 协议、不完整的模板分派、以及缺失的 pipeline 契约。

### Root Cause Patterns

深入审计揭示了 6 个跨文件的根因模式（非孤立的缺口列表），每个模式衍生多个具体缺陷：

---

#### Pattern 1: Scattered State Machine — 状态验证分散且不一致

**核心问题**：状态转换验证逻辑分布在 `status.go`、`submit.go`、`claim.go`、`add.go` 四个文件中，各有不同的验证规则和豁免路径。没有单一的状态机守卫。

**具体缺陷**：

| ID | 缺陷 | 严重性 | 位置 |
|----|------|--------|------|
| SM-1 | `submit.go` 不检查当前状态——`completed`/`rejected`/`blocked` 任务可被任意提交 | **高** | `submit.go:105-195` |
| SM-2 | `status.go` 允许 `blocked → completed` 转换——绕过 submit 工作流（无 record、无 quality gate） | **高** | `status.go:118-129` |
| SM-3 | `add.go --block-source` 可将 `completed` 任务改为 `blocked`——违反终态不变量 | **高** | `add.go:173` |
| SM-4 | `claim.go` 自动解锁忽略 `BlockedReason`——被 `--block-source` 阻塞的任务可能在依赖满足时被错误解锁 | **中** | `claim.go:188-199` |
| SM-5 | `BlockedReason` 是 write-only 字段——被写入但从未在解锁决策中读取 | **中** | `claim.go:188` vs `types.go:118` |
| SM-6 | `--force` 标志允许离开终态但文档未说明此行为——用户可能意外使用 | **低** | `status.go:124` |
| SM-7 | `validateRecordData` auto-downgrade 设置 `blocked` 但不设置 `BlockedReason`——任务无原因地被阻塞 | **低** | `submit.go:306-309` |

**失败级联**：SM-1 + SM-3 形成最危险路径——一个已完成的任务被 `--block-source` 改为 blocked（SM-3），然后被 submit 覆盖为 completed（SM-1），绕过了所有质量保证。

---

#### Pattern 2: Mixed Atomicity — 混合的原子性和锁模型

**核心问题**：`submit.go` 使用原子写入 + 文件锁，其余所有命令使用非原子写入且无锁。这创造了并发竞态条件窗口。

**具体缺陷**：

| ID | 缺陷 | 严重性 | 位置 |
|----|------|--------|------|
| MA-1 | `claim.go` 无文件锁——两个 agent 可同时 claim 同一任务（双写导致双执行） | **高** | `claim.go:57-144` |
| MA-2 | `claim.go` 使用 `SaveIndex`（非原子）——与 `submit.go` 的 `SaveIndexAtomic` 混合使用 | **中** | `claim.go:103` |
| MA-3 | `status.go`、`add.go`、`build.go` 均使用非原子 `SaveIndex`——无锁保护 | **中** | 多文件 |
| MA-4 | `state.go` 的 `SaveState` 使用 `os.WriteFile`——崩溃时可截断，`LoadState` 解析失败导致 claim 丢失活跃任务 | **中** | `state.go:28` |
| MA-5 | `.forge/state.json` 无锁保护——`claim`(写 false)、`submit`(写 true)、`quality-gate`(删除) 三个写入者可交叉 | **中** | `claim.go:110`、`submit.go:235`、`quality_gate.go:88` |
| MA-6 | Windows 上 `os.Rename` 对已打开文件可能失败——非锁保护的写入可阻止原子 rename | **低** | `atomic.go` |

**失败级联**：MA-1 是最高风险的并发 bug——两个 `run-tasks` 实例同时 claim，两个 agent 执行同一任务，第一个 submit 创建 record，第二个 submit 必须 `--force` 覆盖。后果：第一个 agent 的工作记录被覆盖丢失。

---

#### Pattern 3: Broken Template Dispatch — 模板分派系统缺陷

**核心问题**：`TaskTypeRegistry`（19 类型）→ `typeToTemplate`（18 映射）→ 磁盘模板文件（20 文件）三者不一致。Prompt 合成链存在断裂。

**具体缺陷**：

| ID | 缺陷 | 严重性 | 位置 |
|----|------|--------|------|
| TD-1 | `code-quality.simplify` 注册了类型但无 `typeToTemplate` 映射——`Synthesize()` 返回 "unknown type"，任务**永久无法执行** | **高** | `prompt.go:22-41` vs `types.go:29` |
| TD-2 | `fix.md` 模板不注入 `SourceTaskID`——agent 执行 fix task 时无法追溯到源任务，完全依赖 fix task 文件的手写内容 | **中** | `prompt/data/fix.md` |
| TD-3 | `code-quality-simplify.md` 使用 `Skill(skill="simplify")` 不带 `forge:` 前缀——与其他所有模板的命名空间约定不一致 | **低** | `prompt/data/code-quality-simplify.md:23` |
| TD-4 | Scope 值不验证——`scope: "frontend"` 在 Go-only 项目中产生 `just compile frontend` 直接失败，无提前验证或友好错误 | **低** | 所有使用 `{{SCOPE}}` 的模板 |
| TD-5 | `verify-regression` 模板直接调用 `just test-e2e` 而非使用 skill——profile 不感知，与 run-e2e-tests skill 不一致 | **低** | `prompt/data/test-pipeline-verify-regression.md:31` |

**失败级联**：TD-1 是已知的**功能性断裂**——任何 `type: "code-quality.simplify"` 的任务在被 claim 时立即失败。`auto.cleanCode` 配置默认 `false` 掩盖了此 bug，但一旦启用就会触发。

---

#### Pattern 4: Fragile Eval Protocol — 脆弱的评估通信协议

**核心问题**：Eval scorer→orchestrator→reviser 的通信依赖自由文本解析和 prompt 级约束，无结构化协议、无编程式守卫、无回滚机制。

**具体缺陷**：

| ID | 缺陷 | 严重性 | 位置 |
|----|------|--------|------|
| EP-1 | Scorer 输出解析无错误恢复——`SCORE:` 行缺失时 orchestrator 行为未定义（可能将 nil 视为通过） | **高** | `eval/SKILL.md` Step 2 |
| EP-2 | Reviser 无编程式作用域强制——仅靠 prompt 指令限制编辑范围，单次幻觉 Edit 可损坏无关文档 | **高** | `doc-reviser.md` |
| EP-3 | Eval 循环无回滚——达到最大迭代但未达标时，文档停留在最后修订状态（可能比初始状态更差） | **中** | `eval/SKILL.md` Step 5 |
| EP-4 | Reviser 缺少项目上下文——scorer 看到 `CONTEXT_CONTENT`（项目约定），reviser 看不到，修订可能违反约定 | **中** | `eval/SKILL.md` Step 4 |
| EP-5 | 无并发控制——两个 eval 循环可同时修改同一文档，无文件锁或互斥 | **中** | 架构层面 |
| EP-6 | Reviser 30% 字数增长检查是 self-policed——orchestrator 不独立验证 | **低** | `doc-reviser.md` |
| EP-7 | Scorer 报告模板路径是两层间接（rubric → report template）——模板文件不存在时无验证 | **低** | `doc-scorer.md` |

**失败级联**：EP-1 + EP-2 + EP-3 形成最差路径——scorer 格式偏差导致 orchestrator 误判分数（EP-1），触发 reviser 在无上下文的情况下修订（EP-4），reviser 越界编辑（EP-2），最终状态比初始更差且无法恢复（EP-3）。

---

#### Pattern 5: Quality Gate Defects — 质量门禁系统性缺陷

**核心问题**：Quality gate 的 fix-task 创建、cap 限制、和 auto-restore 机制之间存在逻辑不一致。

**具体缺陷**：

| ID | 缺陷 | 严重性 | 位置 |
|----|------|--------|------|
| QG-1 | Fix-task `SourceTaskID: "quality-gate:" + step` 不是真实任务 ID——`autoRestoreSourceTask` 的 `FindTask` 静默失败，auto-restore 机制对 quality-gate fix-task **完全失效** | **高** | `quality_gate.go:429` vs `submit.go:244` |
| QG-2 | `countFixTasks` 统计所有状态（含 completed/skipped）——cap 是生命周期计数（3 次），不是活跃计数。一旦 3 个 fix-task 曾被创建（即使全部已完成），该 step 永久无法再创建 fix-task | **高** | `quality_gate.go:353-363` |
| QG-3 | Quality gate 在无 feature 时静默通过（exit 0）——无测试、无验证、用户无感知 | **中** | `quality_gate.go:81-83` |
| QG-4 | `ClearForgeState` 在 fix-task 添加前执行——如果 gate 在 step 1 后崩溃，`.forge/state.json` 已被删除，中间结果丢失 | **中** | `quality_gate.go:88` |

**失败级联**：QG-1 + QG-2 形成功能失效——quality-gate 发现 compile 失败，创建 fix-task（SourceTaskID="quality-gate:compile"），fix-task 完成后 auto-restore 失败（QG-1）。如果此过程重复 3 次（如间歇性编译错误），cap 锁死（QG-2），quality gate 永久无法创建新的 fix-task，用户必须手动干预。

---

#### Pattern 6: Pipeline Contract Gaps — Pipeline 契约系统性缺口

**核心问题**：Skill 间的交接契约不完整——缺少 slug 传播、ID 格式一致、状态查询统一、以及 quick/full 模式兼容性验证。

**具体缺陷**：

| ID | 缺陷 | 严重性 | 位置 |
|----|------|--------|------|
| PC-1 | TC ID 格式不匹配——`run-e2e-tests` 用 `TC-\d+` 提取，`gen-test-cases` 生成 `<target>/<title-slug>` 格式 | **高** | `run-e2e-tests/SKILL.md:142` vs `gen-test-cases/SKILL.md:106` |
| PC-2 | Slug 不传播——brainstorm 和 write-prd 独立确定 slug，feature 目录可能与 proposal 脱节 | **中** | `brainstorm/SKILL.md` → `write-prd/SKILL.md` |
| PC-3 | 测试执行两次——agent 收集 metrics 运行 `just test`，`forge task submit` 再次运行完整 quality gate | **中** | `submit-task/SKILL.md` |
| PC-4 | 状态验证命令不统一——`run-tasks` 用 `forge task status`，`execute-task` 用 `forge task query` | **中** | `run-tasks.md` vs `execute-task.md` |
| PC-5 | Quick mode 下游兼容性未文档化——`consolidate-specs` 总是 drift-only 模式（无 PRD/design），但 skill 未说明此为预期行为 | **中** | `quick-tasks/SKILL.md` → `consolidate-specs/SKILL.md` |
| PC-6 | Quick mode 状态更新非原子——两次 Edit 分别更新 proposal 和 manifest，第一次成功第二次失败导致不一致 | **低** | `quick.md` Step 4 |
| PC-7 | Manifest 文件名歧义——`testing/manifest.md` 和 feature `manifest.md` 同名不同路径，下游 skill "read manifest.md" 可能读错 | **低** | `gen-test-cases/SKILL.md` |
| PC-8 | `quick-tasks` 循环引用——prerequisite 表将 `/quick`（调用 quick-tasks）列为修复措施 | **低** | `quick-tasks/SKILL.md` |
| PC-9 | `execute-task` 提取 KEY 但不使用——死代码 | **低** | `execute-task.md` Step 1 |
| PC-10 | `consolidate-specs` early exit 锁定——`.integrated` marker 阻止重运行，错误分类（LOCAL 应为 CROSS）不可修正 | **低** | `consolidate-specs/SKILL.md` |

#### Pattern 7: Traceability Chain Fragility — 需求追溯链脆弱

**核心问题**：PRD acceptance criteria → test cases → test scripts → e2e results → graduation 的完整追溯链存在多处断裂点。追溯依赖人工生成的字符串引用和 agent 的"自觉性"，无编程式验证。

**具体缺陷**：

| ID | 缺陷 | 严重性 | 位置 |
|----|------|--------|------|
| TC-1 | 无自动化 PRD AC → test case 覆盖率验证——gen-test-cases 可能遗漏部分 AC，无代码检查覆盖率 | **高** | `gen-test-cases/SKILL.md` |
| TC-2 | TC ID 与 PRD section 通过自由文本 `Source` 字段关联——PRD 重构后引用静默失效 | **高** | `gen-test-cases/SKILL.md` → `run-e2e-tests/SKILL.md` |
| TC-3 | TC ID 格式不匹配——`run-e2e-tests` 用 `TC-\d+` 正则提取，`gen-test-cases` 生成 `<target>/<title-slug>` 格式 | **高** | `run-e2e-tests/SKILL.md:142` vs `gen-test-cases/SKILL.md:106` |
| TC-4 | graduate-tests 不验证 TC ID 保持——按函数名去重，不检查 TC ID 是否与原始 test cases 匹配 | **中** | `graduate-tests/SKILL.md` |
| TC-5 | PRD AC 可增删而不触发 test case 失效——无 checksum 或版本链接 | **中** | PRD → test-cases 链路 |
| TC-6 | `phase-inventory.json` 是死产物——breakdown-tasks 写入但无代码消费，与实际 phase 状态可发散 | **低** | `breakdown-tasks/SKILL.md` |

**失败级联**：TC-1 + TC-2 形成追溯黑洞——agent 遗漏了 PRD 中 2 个边界条件的 AC（TC-1），这些 AC 的 TC ID 不存在（TC-2 中 Source 引用失效），导致这些场景永远不被测试，直到生产环境暴露。

---

#### Pattern 8: Task Generation Integrity — 任务生成完整性缺陷

**核心问题**：`forge task index`（BuildIndex）作为任务的权威生成器，存在 orphan 管理、状态保持、类型推断 fallback 不一致等系统性缺陷。

**具体缺陷**：

| ID | 缺陷 | 严重性 | 位置 |
|----|------|--------|------|
| GI-1 | 已删除 .md 文件的任务成为永久 orphan——BuildIndex 只警告不清理，orphan 无限期累积 | **中** | `build.go:160-166` |
| GI-2 | .md 文件存在时 test task 不重新生成——但 index 条目被更新 → 下次 re-index 读回 stale .md frontmatter，**回滚**依赖更新 | **高** | `build.go:293-314` |
| GI-3 | fix-tasks 触发虚假 orphan 警告——`forge task add` 创建的 fix-task 无 .md 文件，每次 index 运行都警告 | **中** | `build.go:161` vs `isAutoGenTaskID` |
| GI-4 | 空类型 fallback 不一致——BuildIndex 硬失败（abort），migrate 命令静默默认为 `TypeFeature` | **中** | `build.go:150` vs `migrate.go:61` |
| GI-5 | 仅 3 个运行时字段被 preserve（Status, SourceTaskID, BlockedReason）——新增运行时字段必须手动添加到 4 个 preserve 点 | **低** | `build.go:130-132` 及 3 处类似 |
| GI-6 | BuildIndex 不自动注入跨 phase 依赖——生成 gate 但不将 Phase N+1 的任务依赖连接到 Phase N 的 gate | **中** | `stage_gates.go` vs `validate_index.go:285-349` |
| GI-7 | 空任务列表生成 T-eval-doc——`needsDocEval` 对零任务返回 true，评估"无" | **低** | `build.go:438-450` |
| GI-8 | `checkDependenciesMet`（claim.go）比 `checkUnmetDeps`（status.go）更严格——fix-task awareness 和 self-source blocking 只在 claim 中存在 | **中** | `claim.go:284-301` vs `status.go:157-184` |
| GI-9 | Wildcard dependency `.x` 不排除 fix/auto-gen 任务——`isBusinessTask` 只排除 `.gate/.summary`，`1.x` 会匹配 `fix-1`、`T-test-1` | **中** | `claim.go:262` |

**失败级联**：GI-2 是最隐蔽的——用户通过修改 task .md 更新了依赖关系，然后运行 `forge task index`。BuildIndex 更新了 index.json 中的条目。但如果 test task 的 .md 已存在，BuildIndex 不重新生成它，而是用 `TestTaskDef` 的新元数据更新 index 条目。下次 re-index 时，BuildIndex 从 stale .md 解析 frontmatter，**覆盖**了之前更新的 index 条目。依赖更新被静默回滚。

---

#### Pattern 9: Divergent Status Stores — 三层状态存储发散

**核心问题**：任务状态存储在三个独立的文件中（`index.json`、`process/state.json`、`.forge/state.json`），由不同的命令写入和读取，无统一的一致性保证。

**具体缺陷**：

| ID | 缺陷 | 严重性 | 位置 |
|----|------|--------|------|
| DS-1 | index.json（权威状态）、process/state.json（工作区标记）、.forge/state.json（会话标记）可在 submit 和 cleanup 之间发散 | **高** | 三文件交互 |
| DS-2 | Feature context 可被静默丢失——`ClearForgeState`（quality-gate）删除 .forge/state.json，fallback 从 git 分支解析可能指向不同 feature | **高** | `quality_gate.go:88` → `feature.go:34-72` |
| DS-3 | `.forge/state.json` 有 3 个写入者（claim=false, submit=true, quality-gate=delete），无锁保护 | **中** | `claim.go:110`、`submit.go:235`、`quality_gate.go:88` |
| DS-4 | verify_task_done 在 state.json 被清理后允许 commit——返回 nil（允许提交）而非阻塞 | **中** | `verify_task_done.go:44-45` |
| DS-5 | submit 后 cleanup 前——index.json 显示 completed 但 process/state.json 仍存在（stale），后续 claim 可看到不一致 | **低** | `submit.go` → `cleanup.go` 时间窗口 |

**失败级联**：DS-2 + DS-3 形成最危险路径——quality-gate 清除 .forge/state.json（DS-2），如果在 cap 耗尽后退出（不创建 fix-task），.forge/state.json 被删除且不被重建。下一个命令的 `GetCurrentFeature` 回退到 git 分支解析，可能解析到 "main" 分支 → 无 feature → 错误。用户必须手动 `forge feature set` 恢复上下文。

---

#### Pattern 10: Agent Behavioral Containment — Agent 行为约束全部依赖 prompt

**核心问题**：task-executor 的所有铁律（一次一个任务、禁止 claim 下一个、最多 3 次 subagent 调用、submit-task 强制执行）全部是 prompt 级指令，零编程式强制。模型可以选择忽略任何约束而无任何后果。

**具体缺陷**：

| ID | 缺陷 | 严重性 | 位置 |
|----|------|--------|------|
| AB-1 | task-executor 可调用 `forge task claim` 领取下一个任务——无编程式拦截，仅靠 prompt 指令"禁止" | **高** | `task-executor.md:16` |
| AB-2 | `forge task claim` 无 cross-feature guard——task-executor 理论上可 claim 不同 feature 的任务 | **高** | `claim.go` 无 feature 验证 |
| AB-3 | "最多 3 次 subagent 调用" 无计数器——Claude Code 不暴露 subagent 调用计数 hook，完全依赖模型自觉 | **中** | `task-executor.md:14` |
| AB-4 | doc-reviser 可编辑 DOC_DIR 外文件——Edit 工具无路径限制，prompt 级约束可被忽略 | **中** | `doc-reviser.md:49-50` |
| AB-5 | doc-scorer 输出格式可带 markdown 格式——`**SCORE: 850/1000**` 会破坏正则解析 | **中** | `doc-scorer.md:57-72` |
| AB-6 | Forge CLI 未安装时 hook 静默失败——`forge quality-gate` 命令不存在 → exit 127 → 无质量门禁、无警告 | **中** | `hooks.json:61` + Claude Code hook 行为 |
| AB-7 | Manifest `in-progress`/`completed` 状态转换无 skill 管理——guide.md 记录了这些状态但无 skill/template 实现 | **低** | 多个 manifest template |

**失败级联**：AB-1 + AB-2 形成最差路径——task-executor 忽略"禁止 claim"指令，claim 了不同 feature 的任务（AB-2），在错误的 feature 中执行代码变更，然后 submit 更新了错误的 index.json。后续 run-tasks dispatcher 发现当前 feature 的任务被跳过，创建不必要的 fix-task。

---

#### Pattern 11: Config & Documentation Drift — 配置与文档漂移

**核心问题**：`.forge/config.yaml` 的 JSON schema、CLI 代码、guide.md 三者对同一配置项的描述不一致。PostToolUse hook 因 CLI 迁移后未更新而完全失效。

**具体缺陷**：

| ID | 缺陷 | 严重性 | 位置 |
|----|------|--------|------|
| CD-1 | `validate-index.sh` 调用 `task validate`（旧命令名）——应为 `forge task validate-index`，PostToolUse hook **完全失效** | **高** | `scripts/validate-index.sh:18,24` |
| CD-2 | JSON schema `project-type` enum 为 `["frontend","backend","fullstack","mobile","library"]`，但 CLI 只识别 `frontend/backend/mixed`——schema 有 3 个 CLI 不认识的值，CLI 有 1 个 schema 不认识的值 | **中** | `forge-config.schema.json:9` vs `just.go:76-83` |
| CD-3 | guide.md 记录 `auto.e2eTest.quick: true`（默认值），但 `AutoConfigDefaults()` 设置 `false`——文档与代码矛盾 | **中** | `guide.md:157` vs `config.go:46` |
| CD-4 | `run-tasks` dispatcher 的 agent 失败处理假设 "timeout"，但 `subagent_type not found`（插件未安装）是不同错误类型，无专门处理 | **低** | `run-tasks.md:107-113` |

**失败级联**：CD-1 是直接的功能失效——用户每次 Edit/Write index.json 时，PostToolUse hook 应该验证格式完整性。但由于命令名过期，验证永远不执行。index.json 格式错误（如依赖引用不存在的任务、循环依赖）在提交时不会被检测，直到 `forge task claim` 运行时才报错。

---

### Urgency

67 个 lessons-learned 文件反复暴露这些问题的表面症状。但症状驱动的修复（每个 gotcha 一个 patch）无法根除问题——根因在于 9 个架构层面的系统性缺陷：

1. **静默数据损坏**（SM-1/2/3）：系统可在无报错的情况下进入不一致状态
2. **并发数据丢失**（MA-1）：两个 agent 同时执行同一任务，记录被覆盖
3. **功能性断裂**（TD-1, QG-1）：`code-quality.simplify` 任务无法执行，quality-gate auto-restore 失效
4. **AI 行为不受控**（EP-2/4）：reviser 可在无守卫的情况下编辑任意文件
5. **追溯链断裂**（TC-1/2/3）：PRD 到测试的追溯依赖人工自觉，无编程式验证
6. **任务生成状态回滚**（GI-2）：BuildIndex 可静默回滚依赖更新
7. **状态存储发散**（DS-1/2）：三层状态文件可在命令间失去一致性

延迟修复的代价：每次使用 Forge 都在累积不可观测的状态不一致，而 lessons 的数量（67个）表明问题发现速度远快于修复速度。

## Proposed Solution

9 个根因模式对应 7 个机制级修复工作流：

**W1 — 统一状态机守卫**（CLI）：提取 `ValidateTransition(current, target, opts)` 为单一权威验证函数，所有状态变更路径（submit、status、claim、add）必须经过。终结 Pattern 1（分散验证）。

**W2 — 原子性与锁统一**（CLI）：所有 index.json 写入统一使用 `SaveIndexAtomic` + advisory lock。修复 quality-gate fix-task 的 SourceTaskID 和 cap 语义。统一 status store 一致性。终结 Pattern 2（混合原子性）和 Pattern 9（三层状态发散）。

**W3 — 结构化 Eval 协议**（Plugin）：Scorer 输出结构化 JSON，orchestrator 双层解析，reviser scope 编程式守卫，上下文对称注入。终结 Pattern 4（脆弱 eval 协议）。

**W4 — Pipeline 契约硬化**（Plugin + CLI）：修复模板分派断裂、slug 传播、status query 统一、dead code 清理。终结 Pattern 3（模板分派断裂）和 Pattern 6（Pipeline 契约缺口）。

**W5 — 追溯链硬化**（Plugin）：添加 PRD AC → TC ID 的编程式链接和覆盖率验证。统一 TC ID 格式。终结 Pattern 7（追溯链脆弱）。

**W6 — BuildIndex 完整性**（CLI）：修复 orphan 管理、stale .md 回滚、fix-task orphan 噪音、类型推断 fallback 统一。终结 Pattern 8（任务生成完整性）。

**W7 — Agent 约束与配置对齐**（Plugin + CLI）：添加 claim 命令的 cross-feature guard 和 agent 执行边界检查。修复 validate-index.sh 命令名。统一 config schema 和文档。终结 Pattern 10（Agent 行为约束）和 Pattern 11（配置漂移）。

### Innovation Highlights

**单一守卫模式（Single Guard）**：借鉴数据库 CHECK 约束——所有状态变更通过一个 `ValidateTransition` 函数。这不是新概念，但应用于 AI agent 状态管理的关键洞察是：**AI agent（和人类一样）会绕过分散的规则，只有集中式守卫能可靠拦截**。当前 Forge 有 4 个独立验证点（status、submit、claim、add），每个有不同的豁免路径，agent 可选择最弱的路径绕过。

**结构化优先协议（Structured-First Protocol）**：借鉴编译器的 error recovery——scorer 输出 JSON 块（`~~~json` fence 包裹），orchestrator 先尝试 JSON 解析，失败则 fallback 到正则提取。这种"结构化优先、容错其次"的设计承认 LLM 输出的不稳定性，但不放弃结构化的好处。

**编程式 scope 守卫（Programmatic Scope Guard）**：当前 eval reviser 的文件编辑范围完全依赖 prompt 指令——相当于"用自然语言做权限控制"。提案在 orchestrator 层面添加路径验证：reviser 完成后，检查所有编辑的文件路径是否在 `DOC_DIR + 白名单` 内，越界编辑自动回滚。借鉴 Git pre-commit hook 的"不信任提交者"哲学。

**锁优先架构（Lock-First）**：当前只有 `submit` 持有 advisory lock，其他所有写入者（claim、status、add、index）无锁运行。提案将 advisory lock 提升为所有 index.json 写入的前置条件——未获取锁的操作必须等待或失败。这消除了整个"混合原子性"缺陷模式。

## Requirements Analysis

### Key Scenarios

**Happy path:**
- Agent claim 任务 → lock 获取 → 原子写入 → 释放 lock
- Agent submit 任务 → `ValidateTransition` 验证当前状态为 `in_progress` → quality gate → 原子写入状态
- Scorer 返回 JSON 分数 → orchestrator 解析 → 达标 → 完成
- Reviser 在 DOC_DIR 内修订 → orchestrator 验证范围 → 下轮 scoring

**Edge cases:**
- 并发 claim → 第二个 claim 等待 lock → 第一个完成 → 第二个 claim 不同任务
- Submit 已 completed 任务 → `ValidateTransition` 拒绝 → 错误信息含恢复建议
- Scorer JSON 输出格式偏差 → fallback 到正则 → 仍可提取分数 → 继续循环
- Reviser 编辑 DOC_DIR 外文件 → orchestrator 检测 → 回滚越界编辑 → 警告
- Fix-task 完成 → auto-restore 正确找到 source → 恢复为 pending
- `code-quality.simplify` 任务 → `typeToTemplate` 映射存在 → 正常执行

**Error scenarios:**
- Lock 获取超时 → 友好错误信息 + 提示检查其他 agent 进程
- Quality-gate 无 feature → 明确错误（非静默通过）
- Fix-task cap 达到上限 → 清晰提示"3 次 fix 已用尽" + 建议人工干预
- Eval 达到最大迭代 → 报告失败 + 文档保持初始状态（非最终修订状态）

### Non-Functional Requirements

- **向后兼容**：所有变更兼容现有 `index.json` 格式和 record 文件
- **零行为变更（正常路径）**：所有已正常工作的流程不产生可观测的行为差异
- **性能**：Lock 获取增加延迟 < 50ms（advisory lock 已存在，只是扩展使用范围）
- **可观测性**：新守卫的拦截输出遵循现有 `---` 分隔的 KEY: VALUE 错误格式

### Constraints & Dependencies

- Go CLI 的 Cobra 命令结构不变
- Eval 协议变更需 doc-scorer.md 和 doc-reviser.md 同步更新
- 不引入新的外部依赖
- `index.json` schema 不变（兼容性约束）

## Alternatives & Industry Benchmarking

### Industry Solutions

- **数据库 CHECK 约束**：PostgreSQL 的统一验证层——所有写入经过同一套约束
- **LLM Structured Outputs**：Anthropic 的 tool_use JSON 返回、OpenAI 的 response_format——强制模型输出可解析格式
- **File Locking**：SQLite 的 WAL lock、Git 的 index.lock——写入者必须获取锁
- **Circuit Breaker**：质量门禁的 fix-task cap 借鉴断路器模式——防止无限重试

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 6 个系统性缺陷持续存在，25+ 缺口累积，67 个 lessons 继续增长 | Rejected: 静默数据损坏风险不可接受 |
| 逐缺口 patch | 经验主义 | 快速 | 缺口数量持续增长，不解决根因（分散验证、混合原子性） | Rejected: 已证明不够——67 个 lessons 表明 patch 速度 < 发现速度 |
| **四层根因机制** | DB CHECK + File Lock + Structured Output + Contract Testing | 从机制层面消除整类问题 | 变更面广，需仔细测试 | **Selected: 根因修复 > 症状 patch** |

## Feasibility Assessment

### Technical Feasibility

所有变更在当前技术栈内可行：
- Go CLI：advisory lock 机制已存在于 `index/lock.go`，仅需扩展使用范围
- `SaveIndexAtomic` 已存在于 `index/atomic.go`，仅需替换所有 `SaveIndex` 调用
- Plugin skills：SKILL.md/agent 定义是 Markdown 文档变更
- Eval 协议：JSON 输出格式对 LLM 是成熟实践

### Resource & Timeline

7 个 workstream，约 10-12 个任务，预估 5-6 天工作量。

### Dependency Readiness

- Advisory lock：已实现（`pkg/index/lock.go`），`submit.go` 已在使用
- Atomic save：已实现（`pkg/index/atomic.go`），`submit.go` 已在使用
- 无外部依赖

## Scope

### In Scope

**W1 — 统一状态机守卫（CLI）：**

新建 `pkg/task/statemachine.go`：
- `ValidateTransition(current, target, opts)` 函数，opts 包含 `Force bool`、`ViaSubmit bool`、`ViaAdd bool`
- 终态保护：`completed`/`rejected` 不可离开（除非 `Force`）
- Submit 专用：只有 `in_progress` 可通过 submit 转为 `completed`/`blocked`
- `blocked → completed` 必须经过 submit（不可通过 `status` 命令）
- `BlockedReason` 参与解锁决策：有 reason 的 blocked 任务不自动解锁

修改现有命令：
- `submit.go`：调用 `ValidateTransition` 替代当前无验证
- `status.go`：`isTransitionAllowed` 委托给共享函数，`checkUnmetDeps` 补齐 fix-task awareness（与 claim.go 对齐，修复 GI-8）
- `add.go`：`--block-source` 转换经过 `ValidateTransition`
- `claim.go`：auto-unblock 检查 `BlockedReason`；wildcard dependency 排除 fix/auto-gen 任务（修复 GI-9）
- `submit.go`：auto-downgrade 时设置 `BlockedReason`

**W2 — 原子性与锁统一 + Quality Gate 修复 + 状态一致性（CLI）：**

原子性与锁：
- `claim.go`：`SaveIndex` → `SaveIndexAtomic` + 获取 advisory lock
- `status.go`：`SaveIndex` → `SaveIndexAtomic` + 获取 advisory lock
- `build.go`：`SaveIndex` → `SaveIndexAtomic`（index 命令通过 build 调用）
- `state.go`：`SaveState` → atomic write（temp + rename）

Quality Gate 修复：
- `quality_gate.go`：fix-task `SourceTaskID` 改为被阻塞的实际任务 ID（非 sentinel）
- `quality_gate.go`：`countFixTasks` 仅统计活跃 fix-task（非终态），实现活跃 cap 而非生命周期 cap
- `quality_gate.go`：无 feature 时输出明确错误而非静默通过

状态一致性（Pattern 9）：
- `quality_gate.go`：`ClearForgeState` 后立即写入 `allCompleted=false`（而非删除），避免 feature context 丢失
- `verify_task_done.go`：state.json 不存在时返回 warning（而非静默允许），提示用户检查是否已 submit

**W3 — 结构化 Eval 协议（Plugin）：**

- `doc-scorer.md`：输出格式改为 `~~~json` 包裹的 JSON 块（`{"score": N, "total": N, "dimensions": [...], "attacks": [...]}`），保留文本 fallback
- `eval/SKILL.md` Step 2：添加 JSON 解析优先 + 正则 fallback 的两层解析，解析失败中止循环并报错
- `eval/SKILL.md` Step 4：添加 scope 验证——reviser 完成后 diff 检查编辑路径是否在 `DOC_DIR + 白名单`，越界编辑回滚
- `eval/SKILL.md` Step 4：注入 `CONTEXT_CONTENT` 到 reviser prompt（与 scorer 对称）
- `eval/SKILL.md` Step 5：达到最大迭代时回滚文档到 eval 前 state（需在 Step 1 备份）

**W4 — Pipeline 契约硬化（Plugin + CLI）：**

- `prompt.go`：添加 `code-quality.simplify` → `data/code-quality-simplify.md` 映射（修复 TD-1）
- `prompt/data/fix.md`：添加 `{{SOURCE_TASK_ID}}` placeholder + 源任务上下文注入
- `prompt/data/code-quality-simplify.md`：`Skill(skill="simplify")` → `Skill(skill="simplify")`（确认是否需要 forge: 前缀）
- `prompt/data/test-pipeline-verify-regression.md`：改用 skill 调用替代直接 `just test-e2e`
- `run-tasks.md` + `execute-task.md`：统一使用 `forge task query`
- `submit-task/SKILL.md`：metrics 从 quality gate 结果提取，避免二次测试
- `quick-tasks/SKILL.md`：prerequisite 表 `/quick` → `/brainstorm`
- `execute-task.md`：移除 KEY 字段提取
- `quick.md`：Step 4 改为单次原子状态更新
- `consolidate-specs/SKILL.md`：文档说明 quick mode 总是 drift-only 为预期行为

**W5 — 追溯链硬化（Plugin）：**

- `gen-test-cases/SKILL.md`：添加 AC 覆盖率自检步骤——生成 test cases 后，反向验证每个 PRD Given/When/Then 块是否被至少一个 TC 引用，输出覆盖率报告
- `gen-test-cases/SKILL.md`：TC ID 格式统一为 `TC-{NNN}` 数字格式，Source 字段包含 PRD section 的 machine-parseable 标识（如 `story:US-001/ac:AC-3`）
- `run-e2e-tests/SKILL.md`：TC ID 提取正则兼容 `TC-NNN` 和 `TC-NNN-*`（子步骤扩展），更新 `TC-\d+` 为 `TC-\d+(-\w+)?`
- `gen-test-scripts/SKILL.md`：测试函数名必须包含 TC ID，添加验证步骤
- `graduate-tests/SKILL.md`：迁移后验证 TC ID 保持一致

**W6 — BuildIndex 完整性（CLI）：**

- `build.go`：添加 orphan 清理选项（`--clean-orphans` flag），删除无 .md 文件的非 auto-gen orphan 任务
- `build.go`：test task .md 文件存在时，与 `TestTaskDef` 的依赖/元数据做 diff——如有差异则重新生成（修复 GI-2）
- `build.go`/`isAutoGenTaskID`：添加 `fix-` 和 `disc-` 前缀到 auto-gen 检测，消除 fix-task 的虚假 orphan 警告（修复 GI-3）
- `build.go` + `migrate.go`：统一空类型 fallback 策略——BuildIndex 的 hard failure 保留，但添加 `--default-type` flag 允许用户指定 fallback 类型（修复 GI-4）
- `build.go`：将 preserve 逻辑提取为 `PreserveRuntimeFields(existing, new)` 函数，集中管理运行时字段列表，新增字段只需修改一处（修复 GI-5）
- `build.go`/`needsDocEval`：零任务时返回 false，不生成 T-eval-doc（修复 GI-7）

**W7 — Agent 约束与配置对齐（Plugin + CLI）：**

Agent 行为约束：
- `claim.go`：添加 feature scope guard——claim 时验证 `process/state.json`（如果存在）的 feature 与当前 feature 一致，不一致时警告或拒绝（修复 AB-2）
- `claim.go`：添加 `--agent-mode` flag，在 agent 模式下限制每次只能 claim 一个已有 `process/state.json` 的任务（不可重复 claim）（缓解 AB-1）
- `run-tasks.md`：添加 dispatcher 层面的 claim 后验证——claim 返回后检查 claim 的任务是否属于当前 feature（修复 AB-2 的 dispatcher 端）

Hook 修复：
- `scripts/validate-index.sh`：`task validate` → `forge task validate-index`（修复 CD-1）

配置对齐：
- `forge-config.schema.json`：`project-type` enum 统一为 `["frontend", "backend", "mixed", "fullstack", "mobile", "library"]`（CLI 实际支持的 + schema 原有的）（修复 CD-2）
- `guide.md`：`auto.e2eTest.quick` 默认值从 `true` 改为 `false`，与 `AutoConfigDefaults()` 对齐（修复 CD-3）
- Manifest lifecycle：`stop-hook-completion` 提案实现后，`in-progress`/`completed` 状态由 `forge feature complete` 管理（AB-7 依赖独立提案）

### Out of Scope

- `ARCHITECTURE.md` 过时内容更新（属于 `skill-ecosystem-audit` 范围）
- 42 个 `plugins/forge/` 硬编码路径修复（属于 `skill-ecosystem-audit` 范围）
- `eval-forge-runtime-audit` 提案的 6 维度重组（独立的 eval-forge 重构）
- `stop-hook-completion` 提案的 `forge feature complete --if-done` 命令
- `hard-rules-compliance` 提案的 Hard Rules 机制
- 质量门禁跨功能污染问题（需 feature isolation 机制）
- guide.md 与 ARCHITECTURE.md 的内容重复治理
- gen-test-scripts 的 Fact Table 持久化
- graduate-tests 的 marker 校验和
- `forensic` 的 JSONL 路径硬编码
- `validate-ux` 的 teardown 机制
- Manifest `in-progress`/`completed` 状态转换的实现（依赖 `stop-hook-completion` 提案）
- Forge CLI 未安装时的 hook 恢复机制（需要 Claude Code hook runner 改进）
- Agent subagent 调用计数的编程式强制（需要 Claude Code API 支持调用计数 hook）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 统一状态机引入回归——更严格的验证可能拒绝之前允许的转换 | M | M | 先运行现有 e2e 测试建立基线，修改后全量回归。任何被拒绝的转换必须可通过 `--force` 绕过 |
| Claim 加锁增加延迟或死锁——如果 lock 获取失败，整个 claim 阻塞 | L | H | Lock 超时机制（5s）+ 超时时友好错误信息。参考 submit.go 现有的 lock 使用模式 |
| Scorer JSON 格式对 LLM 不稳定 | M | H | 两层解析（JSON 优先 + 正则 fallback）+ eval-test-cases 6 轮验证 |
| Scope 守卫回滚误判——reviser 合法编辑 DOC_DIR 外文件被错误拦截 | M | M | 白名单：`DOC_DIR` + `manifest.md` + `.forge/` 始终允许 |
| Quality-gate fix-task SourceTaskID 改为真实任务 ID 可能影响现有 e2e 测试 | M | M | 修改前检查 `countFixTasks` 和 `autoRestoreSourceTask` 的 e2e 覆盖 |
| Eval 备份/回滚增加复杂度——备份大文档的磁盘开销 | L | L | 仅备份被修改的文件（非全目录），eval 完成后清理备份 |
| BuildIndex stale .md 修复（GI-2）可能触发大规模 .md 重写 | M | L | 仅在依赖/元数据实际变化时重写，不变的 .md 保持不动 |
| AC 覆盖率自检依赖 PRD 格式一致性——agent 生成的 PRD 可能格式不统一 | M | M | 自检步骤使用宽松的 Given/When/Then 匹配（正则），非严格解析 |
| `.forge/state.json` 从"删除"改为"写 false"可能影响依赖 `IsNotExist` 的逻辑 | L | M | 全面 grep `ClearForgeState` 和 `os.IsNotExist` 的使用点，确保兼容 |

## Success Criteria

- [ ] `forge task submit` 对 `completed`/`rejected` 状态返回错误（不使用 `--force`），错误信息包含当前状态、目标状态、恢复建议
- [ ] `forge task status <id> completed` 对 `blocked` 状态返回错误（必须通过 submit）
- [ ] `forge task add --block-source` 对 `completed` 源任务返回错误（不使用 `--force` 时）
- [ ] 所有 index.json 写入使用 `SaveIndexAtomic` + advisory lock（grep 确认）
- [ ] `code-quality.simplify` 任务的 `forge prompt get-by-task-id` 成功返回模板内容
- [ ] Quality-gate fix-task 的 SourceTaskID 是真实任务 ID，auto-restore 不再静默失败
- [ ] `countFixTasks` 仅统计活跃 fix-task（completed/skipped 不计数）
- [ ] Eval scorer 输出不包含 SCORE 时，orchestrator 不崩溃，而是输出解析错误并中止
- [ ] Eval reviser 编辑 DOC_DIR 外文件时，orchestrator 检测到并回滚
- [ ] `run-e2e-tests` 的 TC ID 提取兼容 slug 格式（非仅数字）
- [ ] `run-tasks` 和 `execute-task` 使用同一 CLI 命令验证状态
- [ ] 现有 e2e 测试全部通过（零回归）
- [ ] `quick-tasks` prerequisite 表不引用 `/quick`
- [ ] `gen-test-cases` 输出包含 AC 覆盖率报告（覆盖率 < 100% 时警告）
- [ ] BuildIndex 对 stale test task .md 文件检测到依赖差异时重新生成（非静默保留旧值）
- [ ] fix-task 不再触发 orphan 警告（isAutoGenTaskID 包含 fix-/disc- 前缀）
- [ ] 零任务 feature 不生成 T-eval-doc
- [ ] `checkUnmetDeps`（status.go）包含 fix-task awareness（与 claim.go 对齐）
- [ ] `.forge/state.json` 在 quality-gate 后保持 feature context（不丢失）
- [ ] `validate-index.sh` 调用 `forge task validate-index`（非 `task validate`）
- [ ] `forge-config.schema.json` 的 `project-type` enum 包含 `mixed`（CLI 实际使用的值）
- [ ] `guide.md` 中 `auto.e2eTest.quick` 默认值与 `AutoConfigDefaults()` 一致
- [ ] `forge task claim` 验证 claim 的任务属于当前 feature scope（cross-feature guard）

## Next Steps

- Proceed to `/write-prd` to formalize requirements, or
- Use `/quick` for streamlined task generation

---

## Appendix A: Dead Code & Unused Artifacts Inventory

以下构件经审计确认为 TRULY DEAD（从未被任何非测试代码或 skill/agent/command 引用）或 RARELY USED（注册但从未被 skill 管道调用）：

### TRULY DEAD — Go CLI 导出函数

| 函数 | 位置 | 说明 |
|------|------|------|
| `profile.WriteLanguages` | `pkg/profile/config.go:395` | 定义但从未从非测试代码调用。`config init` 使用本地 `writeConfigFile` |
| `git.GetCurrentBranch` | `pkg/git/git.go:15` | 从未被非测试代码调用 |
| `git.ExtractFeatureFromBranch` | `pkg/git/git.go:91` | 从未被非测试代码调用 |
| `git.IsGitRepository` | `pkg/git/git.go:107` | 从未被非测试代码调用 |
| `git.Run` | `pkg/git/git.go:131` | 从未被非测试代码调用 |
| `e2eprobe.ProbeEndpoint` | `pkg/e2eprobe/e2eprobe.go:16` | 内部使用但不必要导出 |
| `e2eprobe.ExtractYAMLStringField` | `pkg/e2eprobe/e2eprobe.go:74` | 内部使用但不必要导出 |
| `cmd/forge.GetName` | `cmd/forge/run.go:21` | 从未被非测试代码调用 |
| `cmd/forge.IsTestMode` | `cmd/forge/run.go:26` | 从未被非测试代码调用 |

### TRULY DEAD — Plugin 文件

| 文件 | 位置 | 说明 |
|------|------|------|
| `api-test-cases.md` | `skills/gen-test-cases/templates/` | Legacy 模板，SKILL.md 引用 `types/` 而非 `templates/` |
| `cli-test-cases.md` | `skills/gen-test-cases/templates/` | 同上 |
| `mobile-test-cases.md` | `skills/gen-test-cases/templates/` | 同上 |
| `test-cases.md` | `skills/gen-test-cases/templates/` | 同上 |
| `tui-test-cases.md` | `skills/gen-test-cases/templates/` | 同上 |
| `ui-test-cases.md` | `skills/gen-test-cases/templates/` | 同上 |
| `simplify-skill.md` | `commands/simplify-skill.md` | 一次性工具，无任何 plugin 文件引用 |

### RARELY USED — CLI 命令（注册但从未被 skill/agent 调用）

| 命令 | 位置 | 说明 |
|------|------|------|
| `forge task list-types` | `internal/cmd/list_types.go` | 仅交互式调试用 |
| `forge task migrate` | `internal/cmd/migrate.go` | 一次性迁移工具 |
| `forge probe` | `internal/cmd/probe.go` | 仅 justfile 集成项目用 |
| `forge lesson` | `internal/cmd/lesson.go` | 仅 info 命令 |
| `forge proposal` | `internal/cmd/proposal.go` | 仅 info 命令 |
| `forge claude` | `internal/cmd/claude.go` | 便捷 wrapper |
| `forge config` | `internal/cmd/config.go` | 配置管理，未被 skill 引用 |
