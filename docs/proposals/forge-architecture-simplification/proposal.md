---
created: 2026-05-17
author: "faner + Claude"
status: Draft
---

# Proposal: Forge Architecture Simplification

> **附录**：所有具体缺陷的完整清单见 [defect-inventory.md](defect-inventory.md)（18 个模式、~100 个缺陷、11 个死代码项）。

## Problem

Forge v3.0.0-beta-7 的架构存在 18 个系统性缺陷模式，但根本原因不是"缺了某个验证"或"少了个锁"——而是**结构和逻辑层面的不清晰**：状态机验证分散在 4 个文件中各有不同规则，状态存储在 3 个独立文件中无一致性保证，eval 协议依赖自由文本解析，agent 行为约束全部依赖 prompt 而无编程式强制。

**目标不是修补 100 个缺口，而是重新设计使结构和逻辑清晰。**

### Design Principles

1. **单一权威（Single Authority）**：每个决策只在一个地方做出——状态转换验证只在一个函数，index 写入只经过一个路径，eval 输出只有一个格式
2. **契约式约束（Design by Contract）**：关键约束不依赖 prompt 或文档——前置条件、后置条件、不变量在代码层面声明和强制。状态机守卫、文件锁、scope 验证、类型映射完整都必须在代码层面保证
3. **渐进式变更（Incremental Change）**：分阶段实施，每阶段独立可验证。不在同一变更窗口中混合逻辑变更与 cosmetic 变更
4. **保持行为等价性（Preserve Behavioral Equivalence）**：纯重构不改变外部行为；行为变更（如更严格的状态机验证）必须显式标注、有 characterization test 保护、有 `--force` 逃生舱口

### Current State After Recent Updates

最近合并的三个功能（forge-worktree, stop-hook-completion, clean-code-skill）已修复部分问题：

| 已修复 | 变更 |
|--------|------|
| TD-1 `code-quality.simplify` 无模板映射 | 重命名为 `clean-code` + 映射已添加 (`prompt.go:41`) |
| Manifest `completed` 状态转换无管理 | `forge feature complete --if-done` 已实现 |
| Clean-code 任务无 scope 控制 | `forge:clean-code` skill 已创建 |

但核心架构问题仍然存在，且新增了 5 个问题（NI-1 ~ NI-5，见附录）。

---

## Architectural Patterns to Redesign

### Pattern 1: State Machine — 散落的验证逻辑需集中

**现状**：`status.go` 有 `isTransitionAllowed`，`submit.go` 无验证，`claim.go` 有 auto-unblock 逻辑，`add.go` 有 `--block-source` 转换。四条路径、四套规则。

**目标**：一个 `statemachine.go`，所有路径经过同一个函数。

```
statemachine.go
├── ValidateTransition(current, target, opts) error
│   ├── 终态保护（completed/rejected 不可离开，除非 Force）
│   ├── Submit-only 路径（in_progress → completed/blocked 必须经 submit）
│   ├── Blocked 解锁规则（必须检查 BlockedReason）
│   └── 依赖检查（统一 claim 和 status 的逻辑）
└── CanAutoUnblock(task, index) bool
    ├── 依赖是否满足
    ├── BlockedReason 是否允许解锁
    └── Fix-task 链是否活跃
```

涉及缺陷：SM-1 ~ SM-8（8 项）

---

### Pattern 2: Write Path — 混合的原子性需统一

**现状**：`submit.go` 用 `SaveIndexAtomic` + advisory lock，其余所有命令用 `SaveIndex` 无锁。8 个写入者、2 种原子性、1 个锁。

**目标**：所有 index 写入经过同一个函数 `SaveIndexLocked`（获取锁 + 原子写入）。

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

**现状**：scorer 输出自由文本（`SCORE: X/Y`），orchestrator 用正则提取。无解析容错、无回滚、无 scope 强制。

**目标**：第一版只做"解析失败则中止"。JSON 格式化和双层解析留后续迭代。

```
Phase 1（本提案）:
  eval/SKILL.md → 解析失败则中止（非崩溃）
  eval/SKILL.md → Step 1 备份 → Step 5 失败时回滚

Phase 2（后续迭代）:
  doc-scorer.md → JSON 输出 + 文本 fallback
  eval/SKILL.md → 双层解析 + scope 验证 + context 对称
```

涉及缺陷：EP-1、EP-3（EP-2、EP-4 延后）

---

### Pattern 4: Quality Gate — fix-task 机制需修正

**现状**：fix-task 的 SourceTaskID 使用 sentinel（非真实 ID），auto-restore 静默失败。cap 是生命周期计数而非活跃计数。

**目标**：fix-task 使用被阻塞的实际任务 ID，cap 改为活跃计数。

涉及缺陷：QG-1 ~ QG-3（3 项）

---

### Pattern 5: Pipeline Contracts — 隐式契约需显式化

**现状**：TC ID 格式不统一、slug 不传播、状态查询命令不一致、config schema 与 CLI enum 错位。

**目标**：统一 ID 格式、命令调用、config schema。

涉及缺陷：PC-1 ~ PC-8（8 项）

---

### Pattern 6: Traceability — 追溯链需编程式验证

**现状**：PRD AC 到 test case 无自动覆盖率验证，TC ID 与 PRD section 通过自由文本关联，`phase-inventory.json` 是死产物。

**目标**：AC 覆盖率自检、machine-parseable 关联字段。

涉及缺陷：TC-1 ~ TC-3（3 项）

---

### Pattern 7: BuildIndex — 任务生成需完整

**现状**：orphan 任务永不清理、stale .md 回滚依赖更新、fix-task 误报 orphan、preserve 字段需手动维护。

**目标**：orphan 清理（默认行为，非 flag）、diff 检测重生成、preserve 函数化。

涉及缺陷：GI-1 ~ GI-6（6 项）

---

### Pattern 8: Status Stores — 三层状态需一致性

**现状**：index.json、process/state.json、.forge/state.json 三个文件可发散。Feature context 可被静默丢失。

**目标**：State 清理改为写 false 而非删除，verify_task_done 检查 state 一致性。

涉及缺陷：DS-1 ~ DS-3（3 项）

---

### Pattern 9: Agent Constraints — 行为约束需编程式强制

**现状**：task-executor 可 claim 下一个任务、无 cross-feature guard、subagent 调用无计数器、reviser 可越界编辑。

**目标**：claim 添加 feature guard、reviser scope 编程式验证。subagent 计数器依赖 Claude Code API，暂不解决。

涉及缺陷：AB-1 ~ AB-4（4 项）

---

### Pattern 10: Config & Doc Drift — 配置与文档需对齐

**现状**：validate-index.sh 调用错误命令、config schema 与 CLI enum 不一致、guide.md 默认值与代码矛盾。

**目标**：修复命令名、统一 enum、对齐默认值。

涉及缺陷：CD-1 ~ CD-3（3 项）

---

### Pattern 11: Code Organization — 70 个文件平铺在单包内

**现状**：`forge-cli/internal/cmd/` 有 70 个文件、27k 行，所有命令共享一个 Go package。

**目标**：按职责拆分子包。此模式**独立为后续提案**（见 Out of Scope），本提案先在现有包内完成逻辑修正。

```
后续提案的目标结构：
internal/cmd/
├── cmd.go, errors.go, output.go, test_results.go
├── task/       # forge task *
├── e2e/        # forge e2e *
├── feature/    # forge feature *
├── worktree/   # forge worktree *
├── forensic/   # forge forensic *
├── testing/    # forge testing *
├── prompt/     # forge prompt *
└── 顶层命令（init, config, quality-gate, cleanup 等）
```

涉及缺陷：CO-1 ~ CO-5（5 项）

---

### Pattern 12: Naming Drift — 名称反映旧职责范围

**现状**：`testgen.go` / `TestTaskDef` 系列命名现在涵盖 spec consolidation、clean-code、doc-eval 等非测试任务。`build.go` 暗示编译实际是构建 index。`TaskFromFile` 无文件操作。

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

### Pattern 13: Error Handling — AIError taxonomy 未全面采用

**现状**：`errors.go` 定义了完整的 AIError 结构（Code/Message/Cause/Hint/Action）和 14 个工厂函数，但 `worktree.go` 完全不用、`submit.go` 绕过 `Exit()`、`quality_gate.go` 吞错。

**目标**：所有命令统一走 `Exit()` + AIError。禁止直接 `os.Exit`、`fmt.Fprintln(os.Stderr)`。

涉及缺陷：EH-1 ~ EH-7（7 项）

---

### Pattern 14: CLI UX 一致性 — 两条路径混用

**现状**：35 个命令用 `Run`（自行 Exit），4 个用 `RunE`（cobra 打印 error）。两套 config init 向导产出不同。

**目标**：统一 `RunE` + `Exit()`。config init 合并为一条路径。补全 `.Args` 验证。

涉及缺陷：UX-1 ~ UX-7（7 项）

---

### Pattern 15: Skill 接口一致性 — 结构与内容漂移

**现状**：forensic skill 有硬编码用户路径。3 种"流程描述"命名。`task-executor` agent 缺 `inputs` frontmatter。`run-tasks`/`execute-task` 和 `quick-tasks`/`breakdown-tasks` 有大段复制粘贴。

**目标**：Skill 模板标准化。Agent inputs frontmatter 统一。共享内容提取。

涉及缺陷：SI-1 ~ SI-8（8 项）

---

### Pattern 16: Config & Schema 管理 — 零散的配置路径

**现状**：无 `config set`。`config get` 只暴露 1/7 auto 字段。`e2eprobe` 用手写 YAML 解析器。常量在两处重复定义。无 schema 版本。

**目标**：统一 config CRUD。消灭手写解析器。常量统一管理。建立 schema 版本机制。

涉及缺陷：CF-1 ~ CF-6（6 项）

---

### Pattern 17: Magic Values & Constants — 魔法值散落

**现状**：`maxFixTasksPerStep=3`、`probeTimeout=5s`、`defaultLockTimeout=5s` 等影响运行时行为的值硬编码。`submit.go` 硬编码中文 `"无"`。`.forge` 目录名散落为字符串字面量。

**目标**：所有影响行为的常量统一到 `pkg/constants/` 或 config。中文默认值改为英文 `"None"`。目录/文件名常量化。

涉及缺陷：MV-1 ~ MV-7（7 项）

---

## Proposed Solution: 4 Phases, 12 Workstreams

### Phase 0 — 前提：Characterization Tests

**在开始任何工作流之前**，为以下模块的当前行为编写集成测试，锁定当前行为（含"不合法但被允许"的行为）：

| 模块 | 覆盖场景 | 对应工作流 |
|------|---------|-----------|
| `claim.go`、`submit.go`、`status.go`、`add.go` | SM-1 ~ SM-8（状态机所有转换路径） | W4 |
| `build.go`（BuildIndex） | orphan 处理、preserve 行为、test task 生成 | W6 |
| `quality_gate.go` | fix-task 创建、cap 计数、SourceTaskID | W4/W9 |
| `worktree.go` | 错误处理路径、flag 行为 | W7 |

这些测试在 Phase 2 实施后需更新为新的预期行为——但确保了**有意识地**改变每一个行为，而非意外引入回归。

### Phase 1 — 卫生清理（2-3 天）

不改变任何运行时行为语义。仅允许：删除死代码、重命名、常量提取、文档修正。所有变更前后现有测试必须同等通过。

**W1 — Dead Code & Naming Cleanup（CLI + Plugin）**

- 删除 11 个 TRULY DEAD 函数/文件
- 13 个命名漂移修复（`TestTaskDef` → `AutoGenTaskDef`、`testgen.go` → `autogen.go` 等，完整列表见 Pattern 12）
- `validate-index.sh` 修复命令名 → `forge task validate-index`
- `indexCmd.Long` 描述更新
- `gen-test-cases/SKILL.md` 引用修正
- `feature_complete.go` 状态大小写统一
- `clean-code` template 字段与 SKILL.md 对齐

**W2 — Constants & Magic Values（CLI + Plugin）**

- 常量 `ForgeDir`/`ForgeConfigFileName` 统一到 `pkg/constants/`
- `.forge`/`config.yaml`/`index.json` 路径字面量统一为常量引用
- `submit.go` 中文 `"无"` → 英文 `"None"`
- 运行时魔法值（`maxFixTasksPerStep`、`probeTimeout`、`lockTimeout`）提取到 constants 或 config
- `isBusinessTask` 统一到 `pkg/task/` 导出

### Phase 2 — 核心修复（4-6 天）

**修复真实缺陷**，每个工作流完成后通过完整测试（含 Phase 0 characterization tests）。

**W4 — State Machine Centralization（CLI）** ← 行为变更

新建 `pkg/task/statemachine.go`：
- `ValidateTransition(current, target, opts)` — 所有路径唯一的验证入口
- `CanAutoUnblock(task, index)` — 统一 auto-unblock 逻辑（含 BlockedReason 检查）
- 修改 `submit.go`、`status.go`、`add.go`、`claim.go` 各自调用此函数
- `checkDependenciesMet` 统一 claim 和 status 的 fix-task awareness
- Wildcard `.x` 排除 fix/auto-gen 任务
- `--force` 逃生舱口始终可用
- Quality-gate fix-task SourceTaskID 使用真实任务 ID（QG-1）
- `countFixTasks` 改为活跃计数（QG-2）
- `addFixTask` 复用 task 包逻辑

**W5 — Write Path Unification（CLI）**

- **Go/No-Go checkpoint**：验证当前锁机制在 Windows 上的跨平台行为。验证失败则重新设计，不影响其他工作流
- 所有 `SaveIndex` 调用替换为 `SaveIndexLocked`（锁 + 原子写入）
- 删除 `SaveIndex`（非原子版本）
- `SaveState` 改为 `SaveStateAtomic`（temp + rename）
- `.forge/state.json` 从"删除"改为"写 false"

**W6 — BuildIndex Integrity（CLI）**

- orphan 清理改为默认行为（非 `--clean-orphans` flag）
- Test task .md diff 检测，变化时重新生成
- `isAutoGenTaskID` 添加 `fix-`/`disc-` 前缀
- Preserve 逻辑提取为 `PreserveRuntimeFields` 函数
- 零任务不生成 T-eval-doc
- BuildIndex 和 migrate 统一空类型 fallback

**W7 — Error Handling & Eval Safety（CLI + Plugin）**

- `worktree.go` 全部改用 AIError 工厂函数 + `Exit()`
- `submit.go` 锁冲突改用 `Exit()`
- `quality_gate.go` 基础设施错误返回 error 而非 nil（QG-3）
- `add.go` config/profile 读取错误改为警告日志
- `e2e_*` 子命令改用 `Exit()`
- 消除 `worktree.go` 双重输出
- 补充缺失的 AIError 工厂函数
- Eval scorer 解析失败时中止（非崩溃）+ 失败时回滚（EP-1、EP-3）

### Phase 3 — 结构优化（4-6 天）

**改善性工作**，在 Phase 2 的行为修复稳定后执行。

**W8 — CLI UX Consistency（CLI）**

- 全部命令统一迁移到 `RunE`，移除 `Run` 模式
- `config init` 合并到 `init.go` 的 huh TUI 向导
- `configGetCmd` 移除 `SilenceErrors`/`SilenceUsage`
- 12 个命令补全 `.Args` 验证
- 无 block 包装的输出改用 `PrintBlock`
- `worktree.go remove` 添加 `--force` flag 或修正提示
- `worktree.go` pre-flight 提取共享函数
- `clean-code` SKILL.md scope detection 改用 `forge task` 命令

**W9 — Pipeline Contracts（Plugin + CLI）**

- `fix.md` 模板添加 `{{SOURCE_TASK_ID}}`
- `verify-regression` 改用 skill 调用
- `run-tasks` + `execute-task` 统一使用 `forge task query`
- `quick-tasks` prerequisite `/quick` → `/brainstorm`
- 移除 `execute-task` KEY 提取
- `submit-task` metrics 从 quality gate 提取避免二次测试
- `forge task claim` 添加 cross-feature guard

**W10 — Config & Schema（CLI + Plugin）**

- 添加 `config set <key> <value>` 命令
- `config get` 扩展支持所有 auto 字段
- Config schema 添加版本字段 + enum 添加 `mixed`
- guide.md `auto.e2eTest.quick` 默认值对齐代码
- `e2eprobe` 改用 typed YAML unmarshal

**W11 — Skill Interface Polish（Plugin）**

- forensic skill 硬编码路径改为动态解析
- `task-executor.md` 添加 `inputs` frontmatter
- 所有 skill "流程描述" section 统一为 `## Process Flow`
- `run-tasks` 与 `execute-task` 共享内容提取
- `quick-tasks` 与 `breakdown-tasks` 共享内容提取

**W12 — Code Organization（CLI）**

按 cobra 命令树将 `internal/cmd/` 拆分为子包。这是 Phase 3 最后一个工作流——在所有逻辑修正（Phase 2）和 CLI 一致性（W8）完成后，在稳定的代码基础上执行。

- `cmd/task/` — claim, submit, status, query, add, index, migrate, check-deps, validate-index, list-types（11 文件，~1,635 行）
- `cmd/e2e/` — validate-specs, run, setup, verify, compile, discover（7 文件，~452 行）
- `cmd/feature/` — feature, feature-complete（2 文件，~694 行）
- `cmd/worktree/` — worktree（1 文件，~321 行）
- `cmd/forensic/` — forensic（1 文件，~997 行）
- `cmd/testing/` — testing（1 文件，~292 行）
- `cmd/prompt/` — get-by-task-id（2 文件，~69 行）
- `cmd/` 根 — root, errors, output, test_results, quality-gate, cleanup, probe, verify-task-done, init, config, claude, lesson, proposal, version

共享基础设施（`errors.go`、`output.go`、`test_results.go`）留在 `cmd/` 根。

**风险控制**：在独立分支上执行，完整 CI 验证后合并。unexported helper 需决定导出或提取到 `internal/`。forensic.go（997 行）应同步拆分为多个文件。

---

## Workstream Dependency Matrix

### 工作流级别 DAG

```
Phase 0: Characterization Tests
    │
    ▼
Phase 1: W1(Dead Code+Naming) ──→ W2(Constants)
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
- W6 依赖 W1（autogen.go 重命名先于 BuildIndex 逻辑变更）
- W4 包含 QG-1/QG-2 修复（依赖 `addFixTask` 复用，来自 W2 的 `isBusinessTask` 统一）
- W12 必须在 W8（RunE 迁移）和 W11 之后执行——包拆分在稳定的 CLI 接口和 error 处理基础上进行

### 文件级冲突矩阵

| 文件 | 修改方 | 顺序约束 |
|------|--------|----------|
| `claim.go` | W4(状态机) → W5(写入) → W7(错误) → W8(RunE) | 严格顺序 |
| `submit.go` | W1(命名) → W5(写入) → W4(状态机) → W7(错误) → W8(RunE) | 严格顺序 |
| `status.go` | W4(状态机) → W5(写入) → W7(错误) → W8(RunE) | 严格顺序 |
| `add.go` | W3(isBusinessTask整合) → W4(状态机) → W5(写入) → W7(错误) → W8(RunE) | 严格顺序 |
| `quality_gate.go` | W5(写入) → W7(错误) → W8(RunE) | 严格顺序 |
| `worktree.go` | W3(pre-flight) → W7(错误) → W8(RunE) | 严格顺序 |
| `build.go` → `build_index.go` | W1(重命名) → W6(orphan) | 先改名再改逻辑 |
| `testgen.go` → `autogen.go` | W1(重命名) → W6(preserve) | 先改名再改逻辑 |

**核心规则**：同一文件上，Phase 1（改名/清理）先于 Phase 2（逻辑修正）先于 Phase 3（结构优化）。

---

## Scope

### In Scope

Phase 0 ~ Phase 3 全部工作流。估算 **16-19 tasks，14-20 days**。

### Out of Scope

- 42 个 `plugins/forge/` 硬编码路径修复（属 `skill-ecosystem-audit`）
- `eval-forge-runtime-audit` 6 维度重组（独立提案）
- Agent subagent 调用计数编程式强制（需 Claude Code API 支持）
- Forge CLI 未安装时的 hook 恢复（需 Claude Code hook runner 改进）
- Manifest `in-progress` 状态转换（已由 `stop-hook-completion` 部分解决）
- 跨功能质量门禁污染（需 feature isolation 机制）
- Eval JSON 格式化 + 双层解析（本提案只做"解析失败中止+回滚"，JSON 留后续迭代）
- Traceability AC 覆盖率自检（优先级低，延后）

---

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 状态机更严格验证引入回归——拒绝之前允许的转换 | M | M | `--force` 逃生舱口；Phase 0 characterization tests 锁定当前行为 |
| W5 文件锁在 Windows/WSL 上的跨平台行为不确定 | M | H | W5 前提条件：先验证锁机制在 Windows 上的行为，必要时重新设计 |
| Phase 1 命名重命名与 Phase 2 逻辑变更的 git blame 干扰 | L | M | 严格分阶段提交；Phase 1 结束后设 tag；每个 Phase 独立分支 |
| Run→RunE 迁移遗漏 error 处理——error 传播路径变化 | M | M | 逐命令迁移，每个命令迁移后跑测试验证 |
| 范围仍偏大——10-15 天估算可能不足 | M | M | Phase 1 和 Phase 2 是必须的；Phase 3 可部分延后 |
| 删除死代码可能破坏测试 | L | M | 死代码函数只在 `_test.go` 中使用，一并清理 |
| BuildIndex orphan 清理改为默认行为——已有项目可能有 orphan | L | L | 清理时打印警告，不删除 .md 文件 |
| W1 大规模重命名影响并行分支——正在进行的 feature branch 无法 rebase | M | M | W1 在所有分支合并后、新分支创建前执行；Phase 1 完成后设 git tag |
| W5 加锁增加 index 写入延迟 | L | M | 基准测试：`forge task index` 在 50+ 任务项目上延迟不超过基线 120% |
| W12 包拆分爆炸半径覆盖 74 文件——无法增量回滚 | H | M | 独立分支执行；完整 CI 验证后合并；unexported helper 逐一决定导出策略 |

---

## Success Criteria

### Phase 0

- [ ] Characterization tests 覆盖 claim/submit/status/add 的所有 SM-1 ~ SM-8 场景（含当前"不合法但被允许"的行为）
- [ ] Characterization tests 覆盖 BuildIndex 的 orphan 处理、preserve 行为
- [ ] Characterization tests 覆盖 quality_gate 的 fix-task 创建、cap 计数行为

### Phase 1

- [ ] 死代码函数全部删除（grep 确认无引用）
- [ ] `TestTaskDef` 重命名为 `AutoGenTaskDef`（grep 确认无旧名称引用）
- [ ] `testgen.go` → `autogen.go`，`build.go` → `build_index.go`
- [ ] `ForgeDir`/`ForgeConfigFileName` 只定义一次
- [ ] `submit.go` 无中文字符串
- [ ] `validate-index.sh` 调用 `forge task validate-index`
- [ ] `isBusinessTask` 只在 `pkg/task/` 定义一次

### Phase 2

- [ ] `forge task submit` 对 `completed`/`rejected` 状态返回错误（非 `--force`）
- [ ] `forge task status <id> completed` 对 `blocked` 状态返回错误
- [ ] `forge task add --block-source` 对 `completed` 源任务返回错误（非 `--force`）
- [ ] 所有 index.json 写入使用 `SaveIndexLocked`（grep 确认无直接 `SaveIndex` 调用）
- [ ] `worktree.go` 全部使用 AIError + Exit()（grep 确认无 raw `os.Exit`）
- [ ] Orphan 任务在 BuildIndex 时默认清理
- [ ] fix-tasks 不触发 orphan 警告
- [ ] Quality-gate fix-task SourceTaskID 是真实任务 ID
- [ ] `countFixTasks` 仅统计活跃 fix-task
- [ ] Eval scorer 解析失败时中止并输出错误（非 panic/crash）
- [ ] `forge task index` 在 50+ 任务项目上延迟不超过基线 120%
- [ ] Phase 0 characterization tests 全部通过（更新预期后的版本）

### Phase 3

- [ ] 所有命令使用 `RunE`（grep 确认无 `Run:` 赋值）
- [ ] 只有一套 config init 流程
- [ ] `config get` 支持所有 7 个 auto 字段
- [ ] `config set <key> <value>` 可用
- [ ] `e2eprobe` 使用 typed YAML unmarshal（`ExtractYAMLStringField` 已删除）
- [ ] `forge-config.schema.json` 包含 `mixed` + 版本字段
- [ ] `task-executor.md` 有 `inputs` frontmatter
- [ ] forensic skill 无硬编码用户路径（`grep -r '~/.claude/projects/-Users-' forensic/SKILL.md` 返回空）
- [ ] 所有 skill "流程描述" section 统一为 `## Process Flow`（grep 确认无 `## Workflow` 或 `## Architecture` 作为流程描述标题）
- [ ] `cmd/` 拆分为子包（task/、e2e/、feature/、worktree/、forensic/、testing/、prompt/）
- [ ] `forensic.go`（997 行）拆分为多个文件（每个 < 300 行）
- [ ] 每个 Phase 结束后现有 e2e 测试全部通过

---

## Next Steps

- Proceed to `/write-prd` to formalize requirements, or
- Use `/quick` for streamlined task generation
