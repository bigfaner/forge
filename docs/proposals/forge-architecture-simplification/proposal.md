---
created: 2026-05-17
updated: 2026-05-18
author: "faner + Claude"
status: Draft
---

# Proposal: Forge Architecture Simplification

> **附录**：所有具体缺陷的完整清单见 [defect-inventory.md](defect-inventory.md)（19 个模式、78 个存续缺陷 + 6 个新增缺陷、9 个死代码项）。

## Problem

Forge v3.0.0-rc.1 的架构存在 19 个系统性缺陷模式（78 个存续 + 6 个新增），根本原因不是"缺了某个验证"或"少了个锁"——而是**结构和逻辑层面的不清晰**：状态机验证分散在 4 个文件中各有不同规则，状态存储在 3 个独立文件中无一致性保证，eval 协议依赖自由文本解析，agent 行为约束全部依赖 prompt 而无编程式强制。

**目标不是修补 84 个缺口，而是重新设计使结构和逻辑清晰。**

### Urgency

**延迟成本**：每个新 PR 在不清晰的结构上叠加代码，引入新的同类缺陷。证据——#112 (skill-ecosystem-audit) 引入 3 个新缺陷（PR-2~PR-4），#113 (contract-journey-test-model) 引入 3 个新缺陷（TM-1~TM-3）。趋势是每个 PR 引入 3-5 个新缺陷，因为根本结构问题未解决。

**具体事故**：MA-1（claim.go 无锁）在并发 agent 场景下可导致两个 agent 同时执行同一任务，浪费 token 且产出冲突。SM-1（submit 不检查状态）允许对已完成的任务重复提交，覆盖原始 record。QG-1（sentinel SourceTaskID）导致 auto-restore 完全失效——fix-task 创建后无法自动追踪到源任务。

### Design Principles

1. **单一权威（Single Authority）**：每个决策只在一个地方做出——状态转换验证只在一个函数，index 写入只经过一个路径，eval 输出只有一个格式
2. **契约式约束（Design by Contract）**：关键约束不依赖 prompt 或文档——前置条件、后置条件、不变量在代码层面声明和强制。状态机守卫、文件锁、scope 验证、类型映射完整都必须在代码层面保证
3. **渐进式变更（Incremental Change）**：分阶段实施，每阶段独立可验证。不在同一变更窗口中混合逻辑变更与 cosmetic 变更
4. **保持行为等价性（Preserve Behavioral Equivalence）**：纯重构不改变外部行为；行为变更（如更严格的状态机验证）必须显式标注、有 characterization test 保护、有 `--force` 逃生舱口

### Changes Since v3.0.0-beta-7

**已修复**（13 项）：
- SM-2: `blocked → completed` 转换现已正确拒绝（`status.go:118-129`）
- SM-4/SM-5: `BlockedReason` 确认为审计日志字段，unblock 决策正确使用依赖状态
- SM-8: Wildcard `.x` 正确排除 `.gate`/`.summary`
- GI-4: BuildIndex 空 type 行为有意与 migrate 分化
- SI-1: forensic 硬编码路径已替换为通用模式
- SI-4/SI-5: breakdown-tasks/quick-tasks 和 run-tasks/execute-task 共享逻辑已提取
- SI-8: 全部 frontmatter 标准化（allowed-tools、argument-hint）
- CD-1: validate-index.sh 已删除
- 全部 `allowed_tools` → `allowed-tools` 和 `argument-hints` → `argument-hint`

**已移除**（4 项，因 contract-journey-test-model 替代旧测试模型）：
- PC-1: TC ID 格式不匹配（Journey 模型使用不同 ID 体系）
- TC-1~TC-3: Traceability 缺陷（Journey 模型通过 Contract 覆盖率替代）

**新增**（6 项，来自 #112 和 #113 合并）：
- MV-8~MV-10: 新测试模型命令的魔法值（`"tests"` 路径、`"exit code 0/1"` 匹配、`"CLAUDE_PROJECT_DIR="` 前缀）
- TM-1~TM-3: test promote 路径遍历、test verify 静默解析失败、unverifiable 标记为 OK

---

## Architectural Patterns to Redesign

### Pattern 1: State Machine — 散落的验证逻辑需集中

**现状**：`status.go` 有 `isTransitionAllowed`，`submit.go` 无状态验证，`claim.go` 有 auto-unblock 逻辑，`add.go` 有 `--block-source` 转换。四条路径、四套规则。`checkDependenciesMet`（claim）比 `checkUnmetDeps`（status）多 fix-task awareness。

**目标**：一个 `statemachine.go`，所有路径经过同一个函数。

```
statemachine.go
├── ValidateTransition(current, target, opts) error
│   ├── 终态保护（completed/rejected 不可离开，除非 Force）
│   ├── Submit-only 路径（in_progress → completed/blocked 必须经 submit）
│   ├── Blocked 解锁规则（必须检查依赖状态）
│   └── 依赖检查（统一 claim 和 status 的逻辑，含 fix-task awareness）
└── CanAutoUnblock(task, index) bool
    ├── 依赖是否满足（已有，正确实现）
    └── Fix-task 链是否活跃（claim 有，status 缺）
```

涉及缺陷：SM-1、SM-3、SM-6、SM-7（4 项）

---

### Pattern 2: Write Path — 混合的原子性需统一

**现状**：`submit.go` 用 `SaveIndexAtomic` + advisory lock，其余 5 个命令用 `SaveIndex` 无锁。3 个 `.forge/state.json` 写入者无锁无原子写入。

**目标**：所有 index 写入经过同一个函数 `SaveIndexLocked`（获取锁 + 原子写入）。state 写入改为原子。

```
index/atomic.go
├── SaveIndexLocked(path, index, lockOpts) error
│   ├── LockFile() — advisory lock, 5s timeout
│   ├── SaveIndexAtomic() — temp + rename
│   └── UnlockFile()
└── SaveStateAtomic(path, state) error  — 同样 temp + rename
```

**前提**：验证当前锁机制在 Windows（WSL）上的跨平台行为。

涉及缺陷：MA-1 ~ MA-4（4 项）

---

### Pattern 3: Eval Protocol — 自由文本协议需结构化

**现状**：scorer 输出自由文本（`SCORE: X/Y`），orchestrator 用正则提取。无解析容错、无回滚、reviser 无项目上下文。

**目标**：第一版只做"解析失败则中止 + 回滚 + context 对称"。JSON 格式化留后续迭代。

```
Phase 1（本提案）:
  eval/SKILL.md → Step 1 备份 → 解析失败则中止（非崩溃）→ Step 5 失败时回滚
  eval/SKILL.md → Reviser 注入与 Scorer 相同的 CONTEXT_CONTENT

Phase 2（后续迭代）:
  doc-scorer.md → JSON 输出 + 文本 fallback
  eval/SKILL.md → 双层解析 + scope 验证 + context 对称
```

涉及缺陷：EP-1 ~ EP-4（4 项）

---

### Pattern 4: Quality Gate — fix-task 机制需修正

**现状**：fix-task 的 SourceTaskID 使用 sentinel（非真实 ID），auto-restore 静默失败。cap 是生命周期计数而非活跃计数。无 feature 时静默通过。

**目标**：fix-task 使用被阻塞的实际任务 ID，cap 改为活跃计数，无 feature 时返回错误。

涉及缺陷：QG-1 ~ QG-3（3 项）

---

### Pattern 5: Pipeline Contracts — 隐式契约需显式化

**现状**：slug 不传播、状态查询命令不一致、config schema 与 CLI enum 错位、测试执行两次、prerequisite 循环引用。

**目标**：统一 slug 传播、命令调用、config schema。

涉及缺陷：PC-2 ~ PC-8（7 项）

---

### Pattern 6: BuildIndex — 任务生成需完整

**现状**：orphan 任务永不清理（仅 emit warning）、fix-task 误报 orphan、preserve 仅 3 字段、无跨 phase 依赖注入。

**目标**：orphan 清理（默认行为，非 flag）、fix-task 排除、preserve 函数化。

涉及缺陷：GI-1 ~ GI-3、GI-5、GI-6（5 项）

---

### Pattern 7: Status Stores — 三层状态需一致性

**现状**：index.json、process/state.json、.forge/state.json 三个文件可发散。Feature context 可被静默丢失。

**目标**：State 清理改为写 false 而非删除，verify_task_done 检查 state 一致性。

涉及缺陷：DS-1 ~ DS-3（3 项）

---

### Pattern 8: Agent Constraints — 行为约束需编程式强制

**现状**：task-executor 可 claim 下一个任务、无 cross-feature guard、subagent 调用无计数器、reviser 可越界编辑。

**目标**：claim 添加 feature guard、reviser scope 编程式验证。subagent 计数器依赖 Claude Code API，暂不解决。

涉及缺陷：AB-1 ~ AB-4（4 项）

---

### Pattern 9: Code Organization — ~75 个文件平铺在单包内

**现状**：`forge-cli/internal/cmd/` 有 ~75 个文件，所有命令共享一个 Go package。`isBusinessTask` 重复定义，`validator` 业务逻辑放在 cmd 层。

**目标**：按职责拆分子包。此模式**独立为后续提案**（见 Out of Scope），本提案先在现有包内完成逻辑修正。

```
后续提案的目标结构：
internal/cmd/
├── cmd.go, errors.go, output.go, test_results.go
├── task/       # forge task *
├── e2e/        # forge e2e *  (现为 forge test *)
├── feature/    # forge feature *
├── worktree/   # forge worktree *
├── forensic/   # forge forensic *
├── testing/    # forge testing *
├── prompt/     # forge prompt *
└── 顶层命令（init, config, quality-gate, cleanup 等）
```

涉及缺陷：CO-1 ~ CO-5（5 项）

---

### Pattern 10: Naming Drift — 名称反映旧职责范围

**现状**：`testgen.go` / `TestTaskDef` 系列命名现在涵盖 spec consolidation、clean-code、doc-eval 等非测试任务。`build.go` 暗示编译实际是构建 index。

**目标**：名称准确反映当前职责。核心重命名簇：

```
testgen.go          → autogen.go
TestTaskDef         → AutoGenTaskDef
GenerateTestTaskMD  → GenerateAutoGenTaskMD
generateTestTasks   → generateAutoGenTasks
GetBreakdownTestTasks → GetBreakdownAutoTasks
GetQuickTestTasks   → GetQuickAutoTasks
ResolveFirstTestDep → ResolveFirstAutoGenDep
TaskFromFile        → ToTask
build.go            → build_index.go
TestInterfaces      → Interfaces
NewTestIndex        → NewTaskIndexForTest
```

涉及缺陷：NM-1 ~ NM-13（13 项）

---

### Pattern 11: Error Handling — AIError taxonomy 未全面采用

**现状**：`errors.go` 定义了完整的 AIError 结构和 14 个工厂函数，但 `worktree.go` 完全不用、`submit.go` 绕过 `Exit()`、`quality_gate.go` 吞错。

**目标**：所有命令统一走 `Exit()` + AIError。禁止直接 `os.Exit`、`fmt.Fprintln(os.Stderr)`。

涉及缺陷：EH-1 ~ EH-7（7 项）

---

### Pattern 12: CLI UX 一致性 — 两条路径混用

**现状**：35 个命令用 `Run`（自行 Exit），8 个用 `RunE`（cobra 打印 error）。两套 config init 向导产出不同。

**目标**：统一 `RunE` + `Exit()`。config init 合并为一条路径。补全 `.Args` 验证。

涉及缺陷：UX-1 ~ UX-7（7 项）

---

### Pattern 13: Skill 接口一致性 — 结构与内容漂移

**现状**：3 种"流程描述"命名。`task-executor` agent 缺 `inputs` frontmatter。Skill step 编号和 Prerequisites 排序不统一。

**目标**：Skill 模板标准化。Agent inputs frontmatter 统一。共享结构规范。

涉及缺陷：SI-2、SI-3、SI-6、SI-7（4 项）

---

### Pattern 14: Config & Schema 管理 — 零散的配置路径

**现状**：无 `config set`。`config get` 只暴露 1/7 auto 字段。`e2eprobe` 用手写 YAML 解析器。常量在两处重复定义。无 schema 版本。

**目标**：统一 config CRUD。消灭手写解析器。常量统一管理。建立 schema 版本机制。

涉及缺陷：CF-1 ~ CF-6（6 项）

---

### Pattern 15: Magic Values & Constants — 魔法值散落

**现状**：`maxFixTasksPerStep=3`、`probeTimeout=5s`、`defaultLockTimeout=5s` 等影响运行时行为的值硬编码。`submit.go` 硬编码中文 `"无"`。`.forge` 目录名散落为字符串字面量。新测试模型命令引入 `"tests"` 路径、`"exit code 0/1"` 匹配等新魔法值。

**目标**：所有影响行为的常量统一到 `pkg/constants/` 或 config。中文默认值改为英文 `"None"`。目录/文件名常量化。

涉及缺陷：MV-1 ~ MV-10（10 项）

---

### Pattern 16: Config & Doc Drift — 配置与文档需对齐

**现状**：config schema 与 CLI enum 不一致、guide.md 默认值与代码矛盾。

**目标**：修复 enum、对齐默认值。

涉及缺陷：CD-2、CD-3（2 项）

---

### Pattern 17: Template Dispatch — 模板分派缺陷

**现状**：Fix task 模板不注入 SourceTaskID、verify-regression 直接调用 just、Scope 值不验证。

**目标**：修复模板注入、调用路径、scope 验证。

涉及缺陷：TD-2、TD-3、TD-4（3 项）

---

### Pattern 18: Newly Introduced Issues — 合并 PR 遗留

**现状**：`feature_complete.go` 状态大小写不一致、`worktree.go remove` 提示不存在的 flag、`clean-code` 用 grep 解析 JSON、template 与 SKILL.md 字段数不匹配、worktree 重复 pre-flight check。

**目标**：逐项修复。

涉及缺陷：NI-1 ~ NI-5（5 项）

---

### Pattern 19: Test Model Command Quality — 新测试命令缺陷

**现状**：test promote 无路径遍历校验、test verify 解析失败静默返回零值、Fact Table 无条目时误标 OK。

**目标**：添加输入校验、错误传播、正确标记 unverifiable。

涉及缺陷：TM-1 ~ TM-3（3 项）

---

## Innovation Highlights

### 1. Single Authority Principle — Systematic Application

将 "每个决策只在一个地方做出" 系统性应用到 19 个模式：状态机验证集中到一个函数、index 写入集中到一个路径、eval 输出集中到一个格式。这不是逐个补丁，而是结构性重新设计——使得新的缺陷无法在分散路径中隐藏。

### 2. Characterization Tests as Safety Net

在重构前为当前行为（含"不合法但被允许"的行为）编写集成测试。这借鉴了 Michael Feathers 在《Working Effectively with Legacy Code》中的方法：先锁定行为，再改变行为。Phase 0 的 characterization tests 确保 Phase 2 的每个行为变更都是**有意识的**，而非意外引入的回归。

### 3. 4-Phase Incremental Approach

Phase 1（重命名/清理）零风险不改行为 → Phase 2（逻辑修正）有 characterization tests 保护 → Phase 3（结构优化）在稳定基础上执行。每个 Phase 独立可验证、可回滚。Phase 3 可部分延后而不影响 Phase 1/2 的核心修复。

### 4. `--force` Escape Hatch Pattern

所有行为收紧（状态机更严格验证、终态保护、blocked 解锁条件）都提供 `--force` 覆盖机制。这使得部署不会因新的更严格验证而破坏已有工作流——用户可以在过渡期内使用 `--force` 维持旧行为，逐步适应新约束。

### 5. Defect-Driven Architecture Redesign

从 84 个具体缺陷中归纳出 19 个系统性模式，再从模式中提炼出 4 个设计原则。这种 bottom-up 方法确保每个设计决策都有实证基础，而非架构审美偏好。

---

## Alternatives & Industry Benchmarking

### Industry Solutions

| 领域 | 行业方案 | 与 Forge 的关系 |
|------|----------|----------------|
| **状态机** | Go `looplab/state`、`qlang/semver` 有限状态机库 | Forge 的状态机极简（6 个状态、~10 个转换），库引入过重。参考其 Guard/Action 模式，但手写 `ValidateTransition` 更合适 |
| **原子文件写入** | SQLite WAL 模式、PostgreSQL `fsync` + rename | Forge 采用相同模式：temp file + `os.Rename`。已有参考实现 `SaveIndexAtomic`，需推广到所有写入者 |
| **错误处理标准化** | Go `fmt.Errorf` + `%w`、Cobra `RunE` 最佳实践 | Forge 已有 `AIError` 结构（Code/Message/Cause/Hint/Action），但未全面采用。统一到 `RunE` + `Exit()` 是 Cobra 推荐做法 |
| **遗留代码重构** | Michael Feathers《Working Effectively with Legacy Code》——Characterization Tests + Sprout Class | Forge 采用 Characterization Tests（Phase 0）+ 新建 `statemachine.go`（Sprout Class 模式），与书中方法一致 |
| **Go 项目布局** | `golang-standards/project-layout`、Kubernetes `cmd/` 子包模式 | Forge 的 W12 包拆分参照 Kubernetes 的 `cmd/` 子包组织方式 |

### Alternatives Comparison

| 方案 | 实现成本 | 风险 | 完整性 | Verdict |
|------|---------|------|--------|---------|
| **Do nothing**——逐个修 bug 不重构 | 每个 bug 1-2h | 每修一个引入新回归（#112/#113 各引入 3-5 个新缺陷） | 永远修不完——根本原因是结构不清晰 | **Rejected**: 84 个缺陷的根因是结构问题，逐个修补无法根治 |
| **Incremental fixes**——只修 SM/MA/QG 核心缺陷（~20 项），不改命名/结构 | ~5 天 | 中——核心缺陷修复，但命名混乱和结构问题持续积累 | 20/84 项——仅解决数据安全关键路径 | **Rejected**: XY 问题——修了结果不改原因，新 PR 继续引入同类缺陷 |
| **Full rewrite**——从零重写 CLI | ~30 天 | 高——重写期间无法合入其他 PR，功能回归风险大 | 理论 100%——但重写可能引入新的设计缺陷 | **Rejected**: 投入产出比差——现有代码 90%+ 可复用，问题是局部结构而非全局腐烂 |
| **4-Phase incremental redesign**（本提案） | 14-20 天 | 中——分阶段可控，Phase 3 可延后 | 84/84 项——全覆盖，含新增缺陷 | **Selected**: 在现有代码基础上结构化修复，每阶段独立可验证可回滚 |

### Chosen Approach Justification

选择 4-Phase incremental redesign 的理由：

1. **与 Feathers 方法一致**：Characterization Tests (Phase 0) → Sprout (Phase 1-2) → Optimize (Phase 3)，经过验证的遗留代码重构方法论
2. **与 Cobra 最佳实践对齐**：`RunE` 统一迁移是 Cobra 官方推荐的错误处理方式
3. **与 SQLite/PostgreSQL 原子写入模式一致**：temp + rename 是文件系统原子写入的事实标准
4. **可延后性**：Phase 3 全部是改善性工作，可根据团队带宽灵活排期

---

## Requirements Analysis

### Key Scenarios

**Scenario 1: 开发者提交完成的任务 [Risk: High]**
- Happy path: `forge task submit --result success` → 任务状态变为 completed，record 文件创建
- Edge case: 任务已是 completed → 拒绝提交，提示使用 `--force`
- Edge case: 测试失败但 result=success → auto-downgrade 到 blocked，设置 BlockedReason
- Error scenario: index.json 并发写入 → 锁冲突，重试或报错

**Scenario 2: 并发 claim 同一任务 [Risk: High]**
- Happy path: Agent A claim → Agent B claim → 各自获得不同任务
- Race condition: Agent A 和 B 同时读取 index → 选中同一任务 → 锁机制保证只有一个成功
- Error scenario: 锁超时 → 返回错误，agent 稍后重试

**Scenario 3: Quality gate 触发 fix-task [Risk: Medium]**
- Happy path: 测试失败 → fix-task 创建，SourceTaskID 指向实际任务 → auto-restore 可追踪
- Edge case: 已有 3 个活跃 fix-task → 拒绝创建新 fix-task（cap 达到）
- Edge case: 已有 3 个 completed fix-task → 允许创建新 fix-task（活跃计数，非生命周期计数）

**Scenario 4: BuildIndex 重建 [Risk: Medium]**
- Happy path: 新增/修改 .md 文件 → BuildIndex 正确更新 index.json
- Edge case: 删除 .md 文件 → orphan 清理，打印警告
- Edge case: fix-task 的 .md 被删除 → 不触发 orphan 警告（豁免）
- Error scenario: .md 文件损坏/缺 frontmatter → 跳过并警告

**Scenario 5: Eval pipeline 失败 [Risk: Low]**
- Happy path: Scorer 评分 → Reviser 修订 → 达到目标分数
- Edge case: Scorer 输出格式不匹配 → 中止并报错（非崩溃）
- Edge case: 达到最大迭代仍未达标 → 回滚到 Step 1 备份
- Error scenario: Reviser 引入违反项目 convention 的修改 → Reviser 可看到 CONTEXT_CONTENT 避免此类修改

**Scenario 6: 开发者配置项目 [Risk: Low]**
- Happy path: `forge config init` → TUI 向导引导配置所有字段
- Edge case: `forge config set auto.gitPush true` → 直接修改配置文件
- Edge case: `forge config get auto.e2eTest` → 显示当前值

### Non-Functional Requirements

| NFR | 要求 | 验证方式 |
|-----|------|----------|
| **Performance** | `forge task index` 在 50+ 任务项目上延迟 ≤ 基线 120% | 基准测试（W5 前后对比） |
| **Performance** | `SaveIndexLocked` 锁获取超时 5s，重试间隔 50ms | 单元测试 |
| **Security** | `forge test promote` 拒绝 `../` 路径遍历 | 单元测试 |
| **Compatibility** | Windows/WSL 上 advisory lock 行为一致 | Go/No-Go checkpoint（W5） |
| **Compatibility** | 已有项目 BuildIndex orphan 清理时打印警告不删除 .md | 集成测试 |
| **Observability** | 所有错误通过 AIError 结构化输出（Code/Message/Cause/Hint/Action） | grep 验证 |
| **Backward compatibility** | `--force` 逃生舱口覆盖所有行为收紧 | characterization tests |

### Constraints & Dependencies

| Constraint | Impact | Workstream |
|------------|--------|------------|
| Claude Code API 不支持 subagent 调用计数 | AB-3 暂不解决，依赖 prompt 约束 | — |
| Go import cycle 阻止 `profile` ↔ `feature` 互相引用 | CF-4 常量需引入 `pkg/constants/` 第三方包 | W2 |
| Windows `flock` 等价机制（`LockFileEx`）行为不确定 | W5 有 Go/No-Go checkpoint | W5 |
| Git blame 干扰——大规模重命名后 git blame 显示重命名 commit | Phase 1 结束后设 tag，Phase 2 从 tag 后开始 | W1 |
| `gen-test-cases/templates/` 中的旧模板仍被部分 skill 引用 | 需确认引用清理完毕再删除 | W1 |

---

## Feasibility Assessment

### Technical Feasibility — 高

| 技术点 | 可行性 | 理由 |
|--------|--------|------|
| State machine 集中化 | **高** | Go 接口 + struct 即可实现。状态空间极简（6 状态、~10 转换），不需要状态机库 |
| 原子文件写入 | **高** | 已有参考实现 `SaveIndexAtomic`（temp + rename），推广到其他写入者是机械性工作 |
| Advisory file locking | **高** | 已有 `index/lock.go` 实现。Windows 兼容性需验证（Go/No-Go checkpoint） |
| RunE 迁移 | **高** | Cobra 原生支持。每个命令独立迁移，渐进式 |
| AIError 统一 | **高** | 已有 14 个工厂函数。`worktree.go` 是最大的迁移对象（~300 行） |
| 包拆分（W12） | **中** | ~75 文件，需逐一决定 unexported helper 的导出策略。爆炸半径大，但可在独立分支安全执行 |

**无外部依赖**：所有变更使用 Go stdlib（`os`、`sync`、`fmt`）+ 已有依赖（cobra、yaml）。不需要引入新的第三方库。

### Resource & Timeline Feasibility

| Phase | 估算 | 风险缓冲 | 可延后性 |
|-------|------|----------|----------|
| Phase 0 | 2 天 | 1 天 | 否——Phase 2 的前提 |
| Phase 1 | 2-3 天 | 0 天 | 否——Phase 2 的前提；但风险极低（纯重命名/删除） |
| Phase 2 | 4-6 天 | 2 天 | 部分——W7 可延后到 Phase 3 |
| Phase 3 | 4-6 天 | 2 天 | **全部可延后**——改善性工作 |
| **总计** | **14-20 天** | — | Phase 1+2 必须（~8-12 天），Phase 3 可延后 |

**Phase 1 风险最低**：纯重命名和删除，现有测试保证行为不变。适合作为第一个 PR。

**Phase 2 核心路径**：W4（状态机）是最高价值工作流，解决 8 个缺陷（SM + QG）。W5（写入路径）有 Go/No-Go checkpoint，验证失败不影响其他工作流。

### Dependency Readiness

| 依赖 | 状态 | 阻塞风险 |
|------|------|----------|
| Phase 0 characterization tests | 需新建 | 无——纯测试编写 |
| Go stdlib `os.Rename` 原子性 | 已验证（POSIX） | 低——Windows NTFS 上也保证原子 |
| Advisory lock 跨平台 | 需验证 | 中——W5 Go/No-Go checkpoint |
| 现有测试套件 | 126+ 测试通过 | 低——提供回归安全网 |
| `AIError` 错误体系 | 已实现（14 个工厂函数） | 无——只需推广使用 |
| `SaveIndexAtomic` 参考实现 | 已存在 | 无——作为 W5 的模板 |

---

## Proposed Solution: 4 Phases, 12 Workstreams

### Phase 0 — 前提：Characterization Tests

**在开始任何工作流之前**，为以下模块的当前行为编写集成测试，锁定当前行为（含"不合法但被允许"的行为）：

| 模块 | 覆盖场景 | 对应工作流 |
|------|---------|-----------|
| `claim.go`、`submit.go`、`status.go`、`add.go` | SM-1、SM-3、SM-6、SM-7（状态机转换路径） | W4 |
| `build.go`（BuildIndex） | orphan 处理、preserve 行为、test task 生成 | W6 |
| `quality_gate.go` | fix-task 创建、cap 计数、SourceTaskID | W4/W9 |
| `worktree.go` | 错误处理路径、flag 行为 | W7 |

这些测试在 Phase 2 实施后需更新为新的预期行为——但确保了**有意识地**改变每一个行为，而非意外引入回归。

### Phase 1 — 卫生清理（2-3 天）

不改变任何运行时行为语义。仅允许：删除死代码、重命名、常量提取、文档修正。所有变更前后现有测试必须同等通过。

**W1 — Dead Code & Naming Cleanup（CLI + Plugin）**

- 删除 9 个 TRULY DEAD 函数/文件
- 13 个命名漂移修复（`TestTaskDef` → `AutoGenTaskDef`、`testgen.go` → `autogen.go` 等，完整列表见 Pattern 10）
- `indexCmd.Long` 描述更新
- `gen-test-cases/SKILL.md` 引用修正
- `feature_complete.go` 状态大小写统一（NI-1）
- `clean-code` template 字段与 SKILL.md 对齐（NI-4）

**W2 — Constants & Magic Values（CLI + Plugin）**

- 常量 `ForgeDir`/`ForgeConfigFileName` 统一到 `pkg/constants/`
- `.forge`/`config.yaml`/`index.json` 路径字面量统一为常量引用
- `submit.go` 中文 `"无"` → 英文 `"None"`（MV-6）
- 运行时魔法值（`maxFixTasksPerStep`、`probeTimeout`、`lockTimeout`）提取到 constants 或 config（MV-1~MV-5）
- 新测试模型魔法值提取：`"tests"` 路径、`"exit code 0/1"`、`"CLAUDE_PROJECT_DIR="`、`"forge-journey-"`（MV-8~MV-10）
- `isBusinessTask` 统一到 `pkg/task/` 导出（CO-2）
- `worktree.go` 移除不存在的 `--force` 提示或添加该 flag（NI-2）
- `worktree.go` 提取共享 pre-flight check 函数（NI-5）
- `clean-code` SKILL.md scope detection 改用 `forge task` 命令（NI-3）

### Phase 2 — 核心修复（4-6 天）

**修复真实缺陷**，每个工作流完成后通过完整测试（含 Phase 0 characterization tests）。

**W4 — State Machine Centralization（CLI）** ← 行为变更

新建 `pkg/task/statemachine.go`：
- `ValidateTransition(current, target, opts)` — 所有路径唯一的验证入口
- `CanAutoUnblock(task, index)` — 统一 auto-unblock 逻辑（已有，保持）
- 修改 `submit.go`、`status.go`、`add.go`、`claim.go` 各自调用此函数
- `checkDependenciesMet` 统一到 statemachine，`checkUnmetDeps` 改为调用同一逻辑（SM-7）
- `--force` 逃生舱口始终可用
- Quality-gate fix-task SourceTaskID 使用真实任务 ID（QG-1）
- `countFixTasks` 改为活跃计数（QG-2）
- 无 feature 时返回错误而非 exit 0（QG-3）
- `addFixTask` 复用 task 包逻辑（CO-4）
- `add.go --block-source` 添加终态保护（SM-3）
- `submit.go` 添加当前状态检查（SM-1）
- auto-downgrade 时设置 `BlockedReason`（SM-6）

**W5 — Write Path Unification（CLI）**

- **Go/No-Go checkpoint**：验证当前锁机制在 Windows 上的跨平台行为。验证失败则重新设计，不影响其他工作流
- 所有 `SaveIndex` 调用替换为 `SaveIndexLocked`（锁 + 原子写入）
- 删除 `SaveIndex`（非原子版本）
- `SaveState` 改为 `SaveStateAtomic`（temp + rename）（MA-3）
- `.forge/state.json` 从"删除"改为"写 false"（DS-1 部分）
- 3 个 state 写入者统一使用 `SaveStateAtomic`（MA-4）

**W6 — BuildIndex Integrity（CLI）**

- orphan 清理改为默认行为（非 `--clean-orphans` flag）（GI-1）
- fix-task 添加 `fix-` 前缀到 orphan 豁免（GI-3）
- Test task .md diff 检测，变化时重新生成
- `isAutoGenTaskID` 添加 `fix-`/`disc-` 前缀
- Preserve 逻辑提取为 `PreserveRuntimeFields` 函数（GI-5）
- BuildIndex 和 migrate 统一空类型 fallback
- `validator` struct 从 cmd 层移到 `pkg/task/`（CO-5）
- `checkDependenciesMet` 和 `checkUnmetDeps` 依赖检查逻辑移到 `pkg/task/`（CO-3）

**W7 — Error Handling & Eval Safety（CLI + Plugin）**

- `worktree.go` 全部改用 AIError 工厂函数 + `Exit()`（EH-1、EH-2）
- `submit.go` 锁冲突改用 `Exit()`（EH-3）
- `quality_gate.go` 基础设施错误返回 error 而非 nil（EH-4）
- `add.go` config/profile 读取错误改为警告日志（EH-5）
- `claim.go` 统一 AIError 使用（EH-6、EH-7）
- test_* 子命令改用 `Exit()`（UX-7）
- Eval scorer 解析失败时中止（非崩溃）+ Step 1 备份 + 失败时回滚（EP-1、EP-3）
- Eval reviser 注入与 scorer 相同的 CONTEXT_CONTENT（EP-4）
- Test promote 添加路径遍历校验（TM-1）
- Test verify 解析失败返回错误而非零值（TM-2）
- Test verify Fact Table 无条目标记为 unverifiable 而非 OK（TM-3）

### Phase 3 — 结构优化（4-6 天）

**改善性工作**，在 Phase 2 的行为修复稳定后执行。

**W8 — CLI UX Consistency（CLI）**

- 全部命令统一迁移到 `RunE`，移除 `Run` 模式（UX-1）
- `config init` 合并到 `init.go` 的 huh TUI 向导（UX-2）
- `configGetCmd` 移除 `SilenceErrors`/`SilenceUsage`（UX-3）
- 12 个命令补全 `.Args` 验证（UX-4）
- 无 block 包装的输出改用 `PrintBlock`（UX-6）
- `worktree.go remove` 添加 `--force` flag 或修正提示
- `worktree.go` pre-flight 提取共享函数
- `verify_task_done.go` exit code 2 统一到 1（UX-5）

**W9 — Pipeline Contracts（Plugin + CLI）**

- `fix.md` 模板添加 `{{SOURCE_TASK_ID}}`（TD-2）
- `verify-regression` 改用 skill 调用（TD-3）
- Scope 值添加验证（TD-4）
- `run-tasks` + `execute-task` 统一使用 `forge task query`（PC-4）
- `quick-tasks` prerequisite `/quick` → `/brainstorm`（PC-5）
- 移除 `execute-task` KEY 提取（PC-6）
- `submit-task` metrics 从 quality gate 提取避免二次测试（PC-3）
- `forge task claim` 添加 cross-feature guard（AB-2）

**W10 — Config & Schema（CLI + Plugin）**

- 添加 `config set <key> <value>` 命令（CF-1）
- `config get` 扩展支持所有 auto 字段（CF-2）
- Config schema 添加版本字段 + enum 添加 `mixed`（CF-5、PC-7）
- guide.md `auto.e2eTest.quick` 默认值对齐代码（PC-8、CD-3）
- `e2eprobe` 改用 typed YAML unmarshal（CF-3）
- Config schema enum 与 CLI 对齐（CD-2）
- `config init` 暴露 auto block（CF-6）

**W11 — Skill Interface Polish（Plugin）**

- `task-executor.md` 添加 `inputs` frontmatter（SI-3）
- 所有 skill "流程描述" section 统一为 `## Process Flow`（SI-2）
- Skill Prerequisites/When-to-Use 排序统一（SI-6）
- Skill step 编号起点统一为 Step 1（SI-7）
- Doc-reviser scope 编程式验证（AB-4、EP-2）

**W12 — Code Organization（CLI）**

按 cobra 命令树将 `internal/cmd/` 拆分为子包。这是 Phase 3 最后一个工作流——在所有逻辑修正（Phase 2）和 CLI 一致性（W8）完成后，在稳定的代码基础上执行。

- `cmd/task/` — claim, submit, status, query, add, index, migrate, check-deps, validate-index, list-types
- `cmd/test/` — test detect, interfaces, promote, verify, run-journey, framework, get-* (原 e2e 命令)
- `cmd/feature/` — feature, feature-complete
- `cmd/worktree/` — worktree
- `cmd/forensic/` — forensic
- `cmd/` 根 — root, errors, output, test_results, quality-gate, cleanup, probe, verify-task-done, init, config, claude, lesson, proposal, version

共享基础设施（`errors.go`、`output.go`、`test_results.go`）留在 `cmd/` 根。

**风险控制**：在独立分支上执行，完整 CI 验证后合并。unexported helper 需决定导出或提取到 `internal/`。forensic.go 应同步拆分为多个文件。

---

## Workstream Dependency Matrix

### 工作流级别 DAG

```
Phase 0: Characterization Tests
    │
    ▼
Phase 1: W1(Dead Code+Naming) ──→ W2(Constants+MagicValues+SmallFixes)
                                          │
                                          ▼
Phase 2: W4(State Machine + QG) ──→ W5(Write Path*)
                                          │
                                     W6(BuildIndex)
                                          │
                                          ▼
                                     W7(Error + Eval)
                                          │
                                          ▼
Phase 3: W8(CLI UX) ──→ W9(Pipeline) ──→ W10(Config) ──→ W11(Skill) ──→ W12(Package Split)
```

- W5 有 go/no-go checkpoint（Windows 锁验证）。验证失败不影响 W4、W6、W7
- W9 依赖 W4（claim cross-feature guard 需要 W4 的状态机就位）
- W6 依赖 W1（autogen.go 重命名先于 BuildIndex 逻辑变更）和 W2（isBusinessTask 统一先于 BuildIndex 逻辑调整）
- W4 包含 QG-1/QG-2 修复（依赖 `addFixTask` 复用，来自 W2 的 `isBusinessTask` 统一）
- W12 必须在 W8（RunE 迁移）和 W11 之后执行——包拆分在稳定的 CLI 接口和 error 处理基础上进行

### 文件级冲突矩阵

| 文件 | 修改方 | 顺序约束 |
|------|--------|----------|
| `claim.go` | W4(状态机) → W5(写入) → W7(错误) → W8(RunE) | 严格顺序 |
| `submit.go` | W1(命名) → W5(写入) → W4(状态机) → W7(错误) → W8(RunE) | 严格顺序 |
| `status.go` | W4(状态机) → W5(写入) → W7(错误) → W8(RunE) | 严格顺序 |
| `add.go` | W2(isBusinessTask整合) → W4(状态机) → W5(写入) → W7(错误) → W8(RunE) | 严格顺序 |
| `quality_gate.go` | W4(QG修复) → W5(写入) → W7(错误) → W8(RunE) | 严格顺序 |
| `worktree.go` | W2(pre-flight提取) → W7(错误) → W8(RunE) | 严格顺序 |
| `build.go` → `build_index.go` | W1(重命名) → W6(orphan) | 先改名再改逻辑 |
| `testgen.go` → `autogen.go` | W1(重命名) → W6(preserve) | 先改名再改逻辑 |
| `test_promote.go` | W2(常量提取) → W7(输入校验) | 先提取再校验 |
| `verify.go` | W2(常量提取) → W7(错误处理) | 先提取再修 |

**核心规则**：同一文件上，Phase 1（改名/清理）先于 Phase 2（逻辑修正）先于 Phase 3（结构优化）。

---

## Scope

### In Scope

Phase 0 ~ Phase 3 全部工作流。估算 **14-17 tasks，14-20 days**。

### Out of Scope

- `eval-forge-runtime-audit` 6 维度重组（独立提案）
- Agent subagent 调用计数编程式强制（需 Claude Code API 支持）（AB-3）
- Forge CLI 未安装时的 hook 恢复（需 Claude Code hook runner 改进）
- 跨功能质量门禁污染（需 feature isolation 机制）
- Eval JSON 格式化 + 双层解析（本提案只做"解析失败中止+回滚+context对称"，JSON 留后续迭代）
- W12 包拆分的详细方案（本提案只定义目标和约束，详细设计由独立提案完成）
- Slug 传播机制（PC-2）的具体方案（需设计 brainstorm ↔ write-prd 的 slug 传递协议）

---

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 状态机更严格验证引入回归——拒绝之前允许的转换 | M | M | `--force` 逃生舱口；Phase 0 characterization tests 锁定当前行为 |
| W5 文件锁在 Windows/WSL 上的跨平台行为不确定 | M | H | W5 前提条件：先验证锁机制在 Windows 上的行为，必要时重新设计 |
| Phase 1 命名重命名与 Phase 2 逻辑变更的 git blame 干扰 | L | M | 严格分阶段提交；Phase 1 结束后设 tag；每个 Phase 独立分支 |
| Run→RunE 迁移遗漏 error 处理——error 传播路径变化 | M | M | 逐命令迁移，每个命令迁移后跑测试验证 |
| 范围仍偏大——14-17 天估算可能不足 | M | M | Phase 1 和 Phase 2 是必须的；Phase 3 可部分延后 |
| 删除死代码可能破坏测试 | L | M | 死代码函数只在 `_test.go` 中使用，一并清理 |
| BuildIndex orphan 清理改为默认行为——已有项目可能有 orphan | L | L | 清理时打印警告，不删除 .md 文件 |
| W1 大规模重命名影响并行分支——正在进行的 feature branch 无法 rebase | M | M | W1 在所有分支合并后、新分支创建前执行；Phase 1 完成后设 git tag |
| W5 加锁增加 index 写入延迟 | L | M | 基准测试：`forge task index` 在 50+ 任务项目上延迟不超过基线 120% |
| W12 包拆分爆炸半径覆盖 ~75 文件——无法增量回滚 | H | M | 独立分支执行；完整 CI 验证后合并；unexported helper 逐一决定导出策略 |
| Test model 新命令质量缺陷影响已有 e2e 测试 | L | M | TM-1~TM-3 修复不影响 gen-test-scripts 生成逻辑，仅改善 CLI 命令健壮性 |

---

## Success Criteria

### Phase 0

- [ ] Characterization tests 覆盖 claim/submit/status/add 的所有 SM-1、SM-3、SM-6、SM-7 场景（含当前"不合法但被允许"的行为）
- [ ] Characterization tests 覆盖 BuildIndex 的 orphan 处理、preserve 行为
- [ ] Characterization tests 覆盖 quality_gate 的 fix-task 创建、cap 计数行为

### Phase 1

- [ ] 死代码函数全部删除（grep 确认无引用）
- [ ] `TestTaskDef` 重命名为 `AutoGenTaskDef`（grep 确认无旧名称引用）
- [ ] `testgen.go` → `autogen.go`，`build.go` → `build_index.go`
- [ ] `ForgeDir`/`ForgeConfigFileName` 只定义一次
- [ ] `submit.go` 无中文字符串
- [ ] `isBusinessTask` 只在 `pkg/task/` 定义一次
- [ ] `feature_complete.go` 状态大小写统一
- [ ] `clean-code` template 与 SKILL.md 字段数一致
- [ ] 新测试模型魔法值全部提取为常量

### Phase 2

- [ ] `forge task submit` 对 `completed`/`rejected` 状态返回错误（非 `--force`）
- [ ] `forge task add --block-source` 对 `completed` 源任务返回错误（非 `--force`）
- [ ] 所有 index.json 写入使用 `SaveIndexLocked`（grep 确认无直接 `SaveIndex` 调用）
- [ ] `worktree.go` 全部使用 AIError + Exit()（grep 确认无 raw `os.Exit`）
- [ ] Orphan 任务在 BuildIndex 时默认清理
- [ ] fix-tasks 不触发 orphan 警告
- [ ] Quality-gate fix-task SourceTaskID 是真实任务 ID
- [ ] `countFixTasks` 仅统计活跃 fix-task
- [ ] Quality-gate 无 feature 时返回非零 exit code
- [ ] Eval scorer 解析失败时中止并输出错误（非 panic/crash）
- [ ] Eval 有回滚机制——达到最大迭代时恢复 Step 1 备份
- [ ] Eval reviser 接收与 scorer 相同的 CONTEXT_CONTENT
- [ ] `forge task index` 在 50+ 任务项目上延迟不超过基线 120%
- [ ] Phase 0 characterization tests 全部通过（更新预期后的版本）
- [ ] `test promote` 拒绝包含 `../` 的 journeyName
- [ ] `test verify` 解析失败时返回错误
- [ ] `test verify` 对无 Fact Table 条目的 Contract 标记为 `unverifiable`
- [ ] auto-downgrade 时 `BlockedReason` 被设置
- [ ] `checkDependenciesMet` 和 `checkUnmetDeps` 调用同一函数

### Phase 3

- [ ] 所有命令使用 `RunE`（grep 确认无 `Run:` 赋值）
- [ ] 只有一套 config init 流程
- [ ] `config get` 支持所有 7 个 auto 字段
- [ ] `config set <key> <value>` 可用
- [ ] `e2eprobe` 使用 typed YAML unmarshal（`ExtractYAMLStringField` 已删除）
- [ ] `forge-config.schema.json` 包含 `mixed` + 版本字段
- [ ] Config schema `project-type` enum 与 CLI 完全一致（CD-2）
- [ ] guide.md `auto.e2eTest.quick` 默认值与代码一致（CD-3）
- [ ] `fix.md` 模板包含 `{{SOURCE_TASK_ID}}`（TD-2）
- [ ] `verify-regression` 使用 skill 调用而非 `just test-e2e`（TD-3）
- [ ] Scope 值对无效输入返回错误（TD-4）
- [ ] `task-executor.md` 有 `inputs` frontmatter
- [ ] 所有 skill "流程描述" section 统一为 `## Process Flow`（grep 确认无 `## Workflow` 或 `## Architecture` 作为流程描述标题）
- [ ] Doc-reviser scope 验证阻止 DOC_DIR 外文件编辑（AB-4）
- [ ] `cmd/` 拆分为子包（task/、test/、feature/、worktree/、forensic/）
- [ ] 每个 Phase 结束后现有 e2e 测试全部通过

---

## Next Steps

- Proceed to `/write-prd` to formalize requirements, or
- Use `/quick` for streamlined task generation
