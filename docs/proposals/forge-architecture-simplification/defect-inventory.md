---
created: 2026-05-17
author: "faner + Claude"
status: Draft
parent: proposal.md
---

# Defect Inventory: Forge v3.0.0-beta-7

本文档是 `proposal.md` 的附录，记录审计发现的所有具体缺陷。Proposal 本体只保留 redesign 方向。

---

## 已修复缺陷

| ID | 缺陷 | 修复方式 |
|----|------|----------|
| TD-1 | `code-quality.simplify` 无模板映射 → agent 执行时找不到 prompt | 重命名为 `clean-code` + `prompt.go:41` 添加映射 |

---

## Pattern 1: State Machine

| ID | 缺陷 | 位置 |
|----|------|------|
| SM-1 | `submit.go` 不检查当前状态——`completed`/`rejected` 任务可被任意提交 | `submit.go:105-195` |
| SM-2 | `status.go` 允许 `blocked → completed`——绕过 submit 工作流 | `status.go:118-129` |
| SM-3 | `add.go --block-source` 可将 `completed` 改为 `blocked`——违反终态不变量 | `add.go:173` |
| SM-4 | `claim.go` auto-unblock 忽略 `BlockedReason`——`--block-source` 阻塞的任务可能被错误解锁 | `claim.go:188-199` |
| SM-5 | `BlockedReason` write-only——被写入但从未在解锁决策中读取 | `claim.go:188` vs `types.go:118` |
| SM-6 | `submit.go` auto-downgrade 设置 `blocked` 但不设置 `BlockedReason` | `submit.go:306-309` |
| SM-7 | `checkDependenciesMet`（claim）比 `checkUnmetDeps`（status）更严格——fix-task awareness 只在 claim 中 | `claim.go:284-301` vs `status.go:157-184` |
| SM-8 | Wildcard `.x` 不排除 fix/auto-gen 任务——`1.x` 会匹配 `fix-1`、`T-test-1` | `claim.go:262` |

---

## Pattern 2: Write Path Atomicity

| ID | 缺陷 | 位置 |
|----|------|------|
| MA-1 | `claim.go` 无文件锁——两个 agent 可同时 claim 同一任务 | `claim.go:57-144` |
| MA-2 | `claim.go`/`status.go`/`add.go`/`build.go` 用非原子 `SaveIndex` | 多文件 |
| MA-3 | `state.go` 的 `SaveState` 用 `os.WriteFile`——崩溃时截断 | `state.go:28` |
| MA-4 | `.forge/state.json` 有 3 个写入者无锁 | `claim.go`/`submit.go`/`quality_gate.go` |

---

## Pattern 3: Template Dispatch

| ID | 缺陷 | 位置 |
|----|------|------|
| TD-2 | Fix task 模板不注入 `SourceTaskID`——agent 无法追溯到源任务 | `prompt/data/fix.md` |
| TD-3 | `verify-regression` 直接调用 `just test-e2e` 而非使用 skill | `prompt/data/test-pipeline-verify-regression.md` |
| TD-4 | Scope 值不验证——`scope: "frontend"` 在 Go-only 项目直接失败 | 所有 `{{SCOPE}}` 模板 |

---

## Pattern 4: Eval Protocol

| ID | 缺陷 | 位置 |
|----|------|------|
| EP-1 | Scorer 输出解析无错误恢复 | `eval/SKILL.md` Step 2 |
| EP-2 | Reviser 无编程式 scope 强制 | `doc-reviser.md` |
| EP-3 | Eval 无回滚——达到最大迭代时文档停留在最后修订状态 | `eval/SKILL.md` Step 5 |
| EP-4 | Reviser 缺少项目上下文——scorer 看到 context，reviser 看不到 | `eval/SKILL.md` Step 4 |

---

## Pattern 5: Quality Gate

| ID | 缺陷 | 位置 |
|----|------|------|
| QG-1 | SourceTaskID 使用 `"quality-gate:" + step`——`FindTask` 静默失败，auto-restore 完全失效 | `quality_gate.go:429` |
| QG-2 | `countFixTasks` 统计所有状态（含 completed）——永久 cap 锁死 | `quality_gate.go:353-363` |
| QG-3 | 无 feature 时静默通过 exit 0 | `quality_gate.go:81-83` |

---

## Pattern 6: Pipeline Contracts

| ID | 缺陷 | 位置 |
|----|------|------|
| PC-1 | TC ID 格式不匹配——`run-e2e-tests` 用 `TC-\d+`，`gen-test-cases` 生成 slug 格式 | `run-e2e-tests` vs `gen-test-cases` |
| PC-2 | Slug 不传播——brainstorm 和 write-prd 独立确定 slug | `brainstorm` → `write-prd` |
| PC-3 | 测试执行两次——agent 收集 metrics + submit quality gate | `submit-task` |
| PC-4 | `run-tasks` 用 `forge task status`，`execute-task` 用 `forge task query` | 两命令 |
| PC-5 | `quick-tasks` prerequisite 循环引用 `/quick` | `quick-tasks/SKILL.md` |
| PC-6 | `execute-task` 提取 KEY 但不使用 | `execute-task.md` Step 1 |
| PC-7 | Config schema enum 与 CLI 不一致——schema 有 fullstack/mobile/library 但 CLI 不认识，CLI 有 mixed 但 schema 不认识 | `forge-config.schema.json` vs `just.go` |
| PC-8 | guide.md `auto.e2eTest.quick` 默认 true 但代码默认 false | `guide.md:157` vs `config.go:46` |

---

## Pattern 7: Traceability

| ID | 缺陷 | 位置 |
|----|------|------|
| TC-1 | 无自动化 PRD AC → test case 覆盖率验证 | `gen-test-cases` |
| TC-2 | TC ID 与 PRD section 通过自由文本 Source 字段关联 | `gen-test-cases` → `run-e2e-tests` |
| TC-3 | `phase-inventory.json` 是死产物——breakdown-tasks 写入但无代码消费 | `breakdown-tasks` |

---

## Pattern 8: BuildIndex

| ID | 缺陷 | 位置 |
|----|------|------|
| GI-1 | 已删除 .md 文件的任务成为永久 orphan——无清理机制 | `build.go:160-166` |
| GI-2 | .md 文件存在时 test task 不重新生成——但 index 条目更新 → re-index 从 stale .md 回滚 | `build.go:293-314` |
| GI-3 | fix-tasks 触发虚假 orphan 警告——`isAutoGenTaskID` 不含 `fix-`/`disc-` | `build.go:161` |
| GI-4 | BuildIndex 空类型硬失败 vs migrate 静默默认 TypeFeature——不一致 | `build.go:150` vs `migrate.go:61` |
| GI-5 | 仅 3 字段被 preserve——新运行时字段需手动添加 4 处 | `build.go:130-132` |
| GI-6 | BuildIndex 不自动注入跨 phase 依赖 | `stage_gates.go` |

---

## Pattern 9: Status Stores

| ID | 缺陷 | 位置 |
|----|------|------|
| DS-1 | index.json、process/state.json、.forge/state.json 可发散 | 三文件交互 |
| DS-2 | Feature context 可被静默丢失——ClearForgeState 后 fallback 可能指向不同 feature | `quality_gate.go:88` |
| DS-3 | verify_task_done 在 state.json 清理后允许 commit | `verify_task_done.go:44-45` |

---

## Pattern 10: Agent Constraints

| ID | 缺陷 | 位置 |
|----|------|------|
| AB-1 | task-executor 可 claim 下一个任务——无编程式拦截 | `task-executor.md:16` |
| AB-2 | `forge task claim` 无 cross-feature guard | `claim.go` |
| AB-3 | "最多 3 次 subagent 调用" 无计数器 | `task-executor.md:14` |
| AB-4 | doc-reviser 可编辑 DOC_DIR 外文件 | `doc-reviser.md:49-50` |

---

## Pattern 11: Config & Doc Drift

| ID | 缺陷 | 位置 |
|----|------|------|
| CD-1 | `validate-index.sh` 调用 `task validate`——PostToolUse hook 完全失效 | `scripts/validate-index.sh:24` |
| CD-2 | Config schema `project-type` enum 与 CLI 不一致 | `forge-config.schema.json` vs `just.go` |
| CD-3 | guide.md `auto.e2eTest.quick` 默认值与代码矛盾 | `guide.md` vs `config.go` |

---

## Pattern 12: Code Organization

| ID | 缺陷 | 位置 |
|----|------|------|
| CO-1 | 70 个文件平铺在 `internal/cmd/`——无包边界，无法按层级理解代码 | `internal/cmd/` |
| CO-2 | `isBusinessTask` 在 `cmd/validate_index.go` 和 `pkg/task/add.go`（名 `isBusinessTaskID`）重复定义 | 两处 |
| CO-3 | 依赖检查 `checkDependenciesMet`（claim）和 `checkUnmetDeps`（status）比 `pkg/task.GetUnmetDependencies` 更完整（fix-task awareness）但留在 cmd 层 | `claim.go`、`status.go` |
| CO-4 | `addFixTask`（quality_gate.go）重复调用 `task.AddTask()` + `task.CreateTaskMarkdown()` + `feature.EnsureForgeState()`，与 `add.go` 的 `executeAdd` 逻辑重叠 | `quality_gate.go` vs `add.go` |
| CO-5 | `validator` struct（validate_index.go）是纯业务逻辑（无 CLI 依赖）但放在 cmd 层 | `validate_index.go` |

---

## Pattern 13: Naming Drift

函数/类型/文件名仍反映旧职责范围，与当前行为不符。核心问题：`testgen.go` / `TestTaskDef` 系列命名现在涵盖 spec consolidation、clean-code、doc-eval 等非测试任务。

| ID | 当前名称 | 问题 | 位置 | 建议名称 |
|----|----------|------|------|----------|
| NM-1 | `TestTaskDef` | 定义所有自动生成任务（含 spec、clean-code、doc-eval），不仅是 test | `task/testgen.go:12` | `AutoGenTaskDef` |
| NM-2 | `testgen.go` | 文件生成所有自动生成任务类型，~40% 非测试 | `task/testgen.go` | `autogen.go` |
| NM-3 | `isTestTaskID` | 对 `T-clean-code-1` 也返回 true | `task/build.go:404` | `isPipelineTaskID`（将 clean-code 移入 `isAutoGenTaskID`） |
| NM-4 | `GenerateTestTaskMD` | 为所有自动生成任务生成 MD，不仅是 test | `task/testgen.go:185` | `GenerateAutoGenTaskMD` |
| NM-5 | `generateTestTasks` | 生成 spec/clean-code 任务 | `task/build.go:377` | `generateAutoGenTasks` |
| NM-6 | `GetBreakdownTestTasks` / `GetQuickTestTasks` | 返回包含 spec、clean-code 的完整任务集 | `task/testgen.go:35,118` | `GetBreakdownAutoTasks` / `GetQuickAutoTasks` |
| NM-7 | `TestInterfaces` | 这些是应用界面类型（cli/api/tui），不是"测试界面" | `task/build.go:25` | `Interfaces` 或 `AppInterfaces` |
| NM-8 | `ResolveFirstTestDep` | 为 clean-code 任务也解析依赖 | `task/testgen.go:406` | `ResolveFirstAutoGenDep` |
| NM-9 | `TaskFromFile` | 无文件操作——纯 struct 转换 | `task/testgen.go:538` | `ToTask` |
| NM-10 | `build.go` | "build" 暗示编译，实际是构建 task index | `task/build.go` | `build_index.go` |
| NM-11 | `NewTestIndex` | 歧义——是"测试任务 index"还是"用于测试的 index"？ | `task/index.go:94` | `NewTaskIndexForTest` |
| NM-12 | `indexCmd.Long` | 描述说"Test tasks auto-generated"，但实际生成 gates、specs、clean-code 等 | `cmd/index.go:28-29` | 更新描述 |
| NM-13 | `gen-test-cases` SKILL.md | 引用 `/breakdown-tasks` 但测试任务来自 `forge task index` | `gen-test-cases/SKILL.md:39` | 更新为 `forge task index` |

---

## Newly Introduced Issues

| ID | 缺陷 | 位置 |
|----|------|------|
| NI-1 | `feature_complete.go` 状态大小写不一致——manifest `"completed"` vs proposal `"Completed"` | `feature_complete.go:121,129` |
| NI-2 | `worktree.go remove` 提示 `--force` 但无此 flag | `worktree.go:104` |
| NI-3 | `clean-code` SKILL.md scope detection 用 `grep` 解析 JSON | `SKILL.md` |
| NI-4 | `clean-code` template summary.md 有 11 个字段但 SKILL.md 只定义 6 个 | 模板 vs SKILL.md |
| NI-5 | `worktree.go` 三个子命令重复相同的 pre-flight check | `runWorktreeStart/Resume/Remove` |

---

## Dead Code Inventory

### TRULY DEAD — 可立即删除

| 文件/函数 | 位置 | 说明 |
|-----------|------|------|
| `profile.WriteLanguages` | `pkg/profile/config.go:395` | 定义但从未调用 |
| `git.GetCurrentBranch` | `pkg/git/git.go:15` | 从未调用 |
| `git.ExtractFeatureFromBranch` | `pkg/git/git.go:91` | 从未调用 |
| `git.IsGitRepository` | `pkg/git/git.go:107` | 从未调用 |
| `git.Run` | `pkg/git/git.go:131` | 从未调用 |
| `e2eprobe.ProbeEndpoint` | `pkg/e2eprobe/e2eprobe.go:16` | 不必要导出 |
| `e2eprobe.ExtractYAMLStringField` | `pkg/e2eprobe/e2eprobe.go:74` | 不必要导出 |
| `cmd/forge.GetName` | `cmd/forge/run.go:21` | 从未调用 |
| `cmd/forge.IsTestMode` | `cmd/forge/run.go:26` | 从未调用 |
| `gen-test-cases/templates/` (6 files) | `skills/gen-test-cases/templates/` | Legacy，SKILL.md 已改用 `types/` |
| `commands/simplify-skill.md` | `commands/` | 一次性工具，无引用 |

### RARELY USED — 保留但标记为调试工具

`forge task list-types`、`forge task migrate`、`forge probe`、`forge lesson`、`forge proposal`、`forge claude`、`forge config`——均为交互式调试/管理命令，skill 管道不调用。

---

## Pattern 14: Error Handling — AIError taxonomy 未全面采用

`errors.go` 定义了完整的 `AIError` 结构（Code/Message/Cause/Hint/Action）和 14 个工厂函数，但多个命令绕过该体系。

| ID | 缺陷 | 位置 |
|----|------|------|
| EH-1 | `worktree.go` 4 个子命令完全不用 AIError——全部 `fmt.Errorf`，无 CAUSE/HINT/ACTION | `worktree.go` 全文 |
| EH-2 | `worktree.go` 手动写 stderr + 返回 error——用户看到双重错误输出 | `worktree.go:67-68, 145-146, 214-215` |
| EH-3 | `submit.go` 锁冲突用 `fmt.Fprintln + os.Exit(1)` 绕过 `Exit()`——无结构化输出 | `submit.go:92-96` |
| EH-4 | `quality_gate.go` 从不用 AIError——所有错误路径返回 nil，基础设施问题（JSON 损坏、目录缺失）静默吞掉 | `quality_gate.go:64-111` |
| EH-5 | `add.go` 丢弃 3 个 config/profile 读取错误（`_` 赋值）——损坏的配置文件导致空 capabilities 无警告 | `add.go:206-216` |
| EH-6 | `claim.go:221` "no task available with met dependencies" 是 plain `fmt.Errorf`——对比 `claim.go:209` 同场景用 `ErrNoPendingTasks()` AIError，体验不一致 | `claim.go:221` vs `claim.go:209` |
| EH-7 | `claim.go:83-84` `LoadState` 和 `index.ByID` 错误被丢弃——数据损坏时下游打印空字段 | `claim.go:83-84` |

---

## Pattern 15: CLI UX 一致性 — 两条路径混用

| ID | 缺陷 | 位置 |
|----|------|------|
| UX-1 | `Run` vs `RunE` 混用——35 个命令用 `Run`（自行 Exit），4 个用 `RunE`（cobra 打印 error），两种错误输出路径 | 全局 |
| UX-2 | 两套 config init——`init.go` 用 huh TUI（4 步含 auto），`config.go` 用 bufio（3 步无 auto），产出不同；`config.go` 在非交互终端挂起 | `init.go:174-295` vs `config.go:87-169` |
| UX-3 | `configGetCmd` 是唯一设 `SilenceErrors: true` + `SilenceUsage: true` 的命令——错误行为与其他命令不一致 | `config.go:35-36` |
| UX-4 | 12 个命令缺 `.Args` 验证——多余位置参数被静默忽略 | `check_deps.go`、`cleanup.go`、`quality_gate.go`、`version.go` 等 |
| UX-5 | `verify_task_done.go` 是唯一用 exit code 2 的命令——无其他命令用非 0/1 区分错误类型 | `verify_task_done.go:33` |
| UX-6 | 输出格式不一致——`list_types.go`、`forensic.go`、`migrate.go`、`prompt_get.go` 无 `PrintBlock` 包装，直接 `fmt.Printf` | 多文件 |
| UX-7 | `e2e_*` 子命令全部用 raw `os.Exit(1)` 而非 `Exit()`——绕过结构化 AI 错误格式 | 5 个 e2e 文件 |

---

## Pattern 16: Skill 接口一致性 — 结构与内容漂移

| ID | 缺陷 | 位置 |
|----|------|------|
| SI-1 | forensic skill 有硬编码用户路径（`~/.claude/projects/-Users-fanhuifeng-...`）——其他用户无法使用 | `forensic/SKILL.md:37,89,93,100` |
| SI-2 | "流程描述" section 有 3 种命名：Process Flow、Architecture、Workflow——同一概念三种叫法 | 多个 skill |
| SI-3 | `task-executor` agent 缺 `inputs` frontmatter——doc-scorer 和 doc-reviser 都有，三个 agent 接口规范不统一 | `agents/task-executor.md` |
| SI-4 | `run-tasks` 与 `execute-task` dispatch/verify/breaking-gate 逻辑复制粘贴 | `commands/run-tasks.md` vs `commands/execute-task.md` |
| SI-5 | `quick-tasks` 与 `breakdown-tasks` Type Assignment、Intent Propagation、Scope Assignment、Template Selection 大段 copy-paste | `quick-tasks/SKILL.md` vs `breakdown-tasks/SKILL.md` |
| SI-6 | Skill Prerequisites/When-to-Use 排序不统一——有的先 Prereq 有的先 When-to-Use | 所有 skill |
| SI-7 | Skill step 编号起点不一致——有的从 Step 0 开始，有的从 Step 1 | 多个 skill |
| SI-8 | `commands/` 下两种极端格式——1 行 stub（eval-*）vs 500 行完整文档（gen-sitemap），无结构契约 | `commands/*.md` |

---

## Pattern 17: Config & Schema 管理 — 零散的配置路径

| ID | 缺陷 | 位置 |
|----|------|------|
| CF-1 | 无 `config set` 命令——用户必须手改 YAML 或重跑 `config init`（覆盖整个文件） | `cmd/config.go` |
| CF-2 | `config get` 只暴露 `auto.gitPush`——其余 6 个 auto 字段不可查询 | `config.go:368-379` |
| CF-3 | `e2eprobe` 用手写 YAML 解析器（`ExtractYAMLStringField`）而非统一 typed unmarshal——不支持嵌套、引号、注释 | `e2eprobe/e2eprobe.go:74-83` |
| CF-4 | 常量 `ForgeDir`/`ForgeConfigFileName` 在 `profile/config.go:16-17` 和 `feature/constants.go:59,61` 重复定义（为避免 import cycle）——任一侧 drift 导致配置加载静默失败 | 两处 |
| CF-5 | 无 config schema 版本字段——字段重命名/删除时用户配置静默丢失数据 | `.forge/config.yaml` |
| CF-6 | `config init` 不暴露 `auto` block——7 个 auto 字段用户只能手动编辑 | `config.go:157-168` |

---

## Pattern 18: Magic Values & Constants — 魔法值散落

| ID | 缺陷 | 位置 |
|----|------|------|
| MV-1 | `maxFixTasksPerStep = 3` 硬编码——影响运行时行为，不可配置 | `quality_gate.go:30` |
| MV-2 | `probeTimeout = 5s` 硬编码 | `e2eprobe/e2eprobe.go:61` |
| MV-3 | `defaultLockTimeout = 5s` 硬编码 | `index/lock.go:16` |
| MV-4 | `time.Sleep(50ms)` 锁重试间隔硬编码 | `index/lock.go:55` |
| MV-5 | `forensicLast` 默认 20——仅 CLI flag 可调 | `forensic.go:68` |
| MV-6 | `submit.go` 硬编码中文字符串 `"无"` 作空值默认——在英文 CLI 中不一致 | `submit.go:371,451,475` |
| MV-7 | `.forge` 目录名和 `config.yaml` 文件名分散在多处字符串字面量中，未统一为常量 | 多文件 |
