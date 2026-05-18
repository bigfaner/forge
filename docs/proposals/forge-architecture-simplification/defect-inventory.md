---
created: 2026-05-17
updated: 2026-05-18
author: "faner + Claude"
status: Draft
parent: proposal.md
---

# Defect Inventory: Forge v3.0.0-rc.1

本文档是 `proposal.md` 的附录，记录审计发现的所有具体缺陷。Proposal 本体只保留 redesign 方向。

基于 v3.0.0-rc.1 代码库全面重审（2026-05-18），含 skill-ecosystem-audit (#112) 和 contract-journey-test-model (#113) 合并后的状态。

---

## 已修复缺陷

| ID | 缺陷 | 修复方式 | 确认 |
|----|------|----------|------|
| TD-1 | `code-quality.simplify` 无模板映射 | 重命名为 `clean-code` + `prompt.go:41` 添加映射 | 早期修复 |
| SM-2 | `status.go` 允许 `blocked → completed` 绕过 submit | `isTransitionAllowed` 现在拒绝所有到 `completed` 的转换 (`status.go:118-129`) | #code |
| SM-4 | `claim.go` auto-unblock 忽略 `BlockedReason` | 设计修正：unblock 基于 `checkDependenciesMet()` 依赖状态检查，`BlockedReason` 是审计日志 | #design |
| SM-5 | `BlockedReason` write-only | 设计修正：intentionally write-only audit metadata，unblock 决策正确使用依赖状态 | #design |
| SM-8 | Wildcard `.x` 不排除 `.gate`/`.summary` | `isBusinessTask()` 过滤 `.gate` 和 `.summary` 后缀 (`validate_index.go:252-255`) | #code |
| GI-4 | BuildIndex 空类型与 migrate 不一致 | BuildIndex 现在对空类型 hard error（`build.go:143-153`），migrate 用 fallback（`migrate.go:60`），行为有意分化 | #code |
| SI-1 | forensic 硬编码用户路径 | 替换为通用模式 (`70d5054`) | #112 |
| SI-4 | run-tasks 与 execute-task 复制粘贴 | 提取到 `references/shared/` (`b15b190`) | #112 |
| SI-5 | quick-tasks 与 breakdown-tasks 复制粘贴 | 提取到 `references/shared/` (`b15b190`) | #112 |
| SI-8 | commands/ 格式不一致 | 全部 frontmatter 标准化（allowed-tools、argument-hint） | #112 |
| CD-1 | validate-index.sh 调用错误命令 | 脚本已完全删除 | #removed |
| — | `allowed_tools` → `allowed-tools` | 全部 18 个 command/skill 文件已修正 | #112 |
| — | `argument-hints` → `argument-hint` | 全部 command 文件已修正 | #112 |

**注**：SM-4 和 SM-5 不是传统意义上的"修复"——经过审计确认，`BlockedReason` 作为审计日志字段是正确设计，unblock 决策应基于依赖状态而非阻塞原因。

---

## 已移除缺陷（模型替代）

以下缺陷因 contract-journey-test-model (#113) 完全替代旧测试模型而过时：

| ID | 缺陷 | 移除原因 |
|----|------|----------|
| PC-1 | TC ID 格式不匹配（`run-e2e-tests` 用 `TC-\d+` vs `gen-test-cases` slug 格式） | 旧模型已替代，Journey-Driven 模型使用不同 ID 体系 |
| TC-1 | 无自动化 PRD AC → test case 覆盖率验证 | Journey 模型通过 Contract 覆盖率替代 |
| TC-2 | TC ID 与 PRD section 自由文本关联 | Journey 步骤有结构化关联 |
| TC-3 | `phase-inventory.json` 是死产物 | Journey 模型不使用此文件 |

---

## Pattern 1: State Machine

| ID | 缺陷 | 状态 | 位置 |
|----|------|------|------|
| SM-1 | `submit.go` 不检查当前状态——`completed`/`rejected` 任务可被任意提交（仅 record-file write-once 保护） | **存在** | `submit.go:105-195` |
| SM-3 | `add.go --block-source` 可将 `completed` 改为 `blocked`——违反终态不变量 | **存在** | `pkg/task/add.go:169-173` |
| SM-6 | `submit.go` auto-downgrade 设置 `blocked` 但不设置 `BlockedReason` | **存在** | `submit.go:306-309` |
| SM-7 | `checkDependenciesMet`（claim）比 `checkUnmetDeps`（status）更严格——fix-task awareness 只在 claim 中 | **部分修复** | `claim.go:251-303` vs `status.go:157-184` |

---

## Pattern 2: Write Path Atomicity

| ID | 缺陷 | 状态 | 位置 |
|----|------|------|------|
| MA-1 | `claim.go` 无文件锁——两个 agent 可同时 claim 同一任务 | **存在** | `claim.go:97-104` |
| MA-2 | `claim.go`/`status.go`/`add.go`/`build.go`/`migrate.go` 用非原子 `SaveIndex`（仅 `submit.go` 用 `SaveIndexAtomic`） | **存在** | 多文件 |
| MA-3 | `state.go` 的 `SaveState` 用 `os.WriteFile`——崩溃时截断 | **存在** | `pkg/task/state.go:37` |
| MA-4 | `.forge/state.json` 有 3 个写入者无锁（`WriteForgeState`、`EnsureForgeState`、`ClearForgeState`） | **存在** | `feature/forge_state.go` |

---

## Pattern 3: Template Dispatch

| ID | 缺陷 | 状态 | 位置 |
|----|------|------|------|
| TD-2 | Fix task 模板不注入 `SourceTaskID`——agent 无法追溯到源任务 | **存在** | `prompt/data/fix.md` |
| TD-3 | `verify-regression` 直接调用 `just test-e2e` 而非使用 skill | **存在** | `prompt/data/test-pipeline-verify-regression.md` |
| TD-4 | Scope 值不验证——`scope: "frontend"` 在 Go-only 项目直接失败 | **存在** | 所有 `{{SCOPE}}` 模板 |

---

## Pattern 4: Eval Protocol

| ID | 缺陷 | 状态 | 位置 |
|----|------|------|------|
| EP-1 | Scorer 输出解析无错误恢复——`SCORE: X/Y` 格式不匹配时无 fallback | **存在** | `eval/SKILL.md:172-175` |
| EP-2 | Reviser 无编程式 scope 强制 | **存在** | `doc-reviser.md` |
| EP-3 | Eval 无回滚——达到最大迭代时文档停留在最后修订状态，原始版本丢失 | **存在** | `eval/SKILL.md:185-189` |
| EP-4 | Reviser 缺少项目上下文——scorer 看到 context，reviser 看不到 | **存在** | `eval/SKILL.md:197-199` |

---

## Pattern 5: Quality Gate

| ID | 缺陷 | 状态 | 位置 |
|----|------|------|------|
| QG-1 | SourceTaskID 使用 `"quality-gate:" + step` sentinel——`FindTask` 静默失败 | **存在** | `quality_gate.go:429` |
| QG-2 | `countFixTasks` 统计所有状态（含 completed）——永久 cap 锁死 | **存在** | `quality_gate.go:350-363` |
| QG-3 | 无 feature 时静默通过 exit 0——掩盖配置错误 | **存在** | `quality_gate.go:127-130` |

---

## Pattern 6: Pipeline Contracts

| ID | 缺陷 | 状态 | 位置 |
|----|------|------|------|
| PC-2 | Slug 不传播——brainstorm 和 write-prd 独立确定 slug | **存在** | `brainstorm` → `write-prd` |
| PC-3 | 测试执行两次——agent 收集 metrics + submit quality gate | **存在** | `submit-task` |
| PC-4 | `run-tasks` 用 `forge task status`，`execute-task` 用 `forge task query` | **存在** | 两命令 |
| PC-5 | `quick-tasks` prerequisite 循环引用 `/quick` | **存在** | `quick-tasks/SKILL.md` |
| PC-6 | `execute-task` 提取 KEY 但不使用 | **存在** | `execute-task.md` Step 1 |
| PC-7 | Config schema enum 与 CLI 不一致——schema 有 fullstack/mobile/library 但 CLI 不认识，CLI 有 mixed 但 schema 不认识 | **存在** | `forge-config.schema.json` vs `just.go` |
| PC-8 | guide.md `auto.e2eTest.quick` 默认 true 但代码默认 false | **存在** | `guide.md:157` vs `config.go:46` |

---

## Pattern 7: BuildIndex

| ID | 缺陷 | 状态 | 位置 |
|----|------|------|------|
| GI-1 | 已删除 .md 文件的任务成为永久 orphan——无清理机制（仅 emit warning） | **存在** | `build.go:159-166` |
| GI-2 | .md 文件损坏/缺 frontmatter 时被静默跳过——旧任务条目作为 orphan 残留 | **部分修复** | `build.go:93-94` |
| GI-3 | fix-tasks 触发虚假 orphan 警告——`isAutoGenTaskID` 不含 `fix-`/`disc-` 前缀 | **存在** | `build.go:161` |
| GI-5 | 仅 3 字段被 preserve——新运行时字段需手动添加 4 处 | **存在** | `build.go:129-132` |
| GI-6 | BuildIndex 不自动注入跨 phase 依赖 | **存在** | `stage_gates.go` |

---

## Pattern 8: Status Stores

| ID | 缺陷 | 状态 | 位置 |
|----|------|------|------|
| DS-1 | index.json、process/state.json、.forge/state.json 可发散 | **存在** | 三文件交互 |
| DS-2 | Feature context 可被静默丢失——ClearForgeState 后 fallback 可能指向不同 feature | **存在** | `quality_gate.go:88` |
| DS-3 | verify_task_done 在 state.json 清理后允许 commit | **存在** | `verify_task_done.go:44-45` |

---

## Pattern 9: Agent Constraints

| ID | 缺陷 | 状态 | 位置 |
|----|------|------|------|
| AB-1 | task-executor 可 claim 下一个任务——无编程式拦截 | **存在** | `task-executor.md:16` |
| AB-2 | `forge task claim` 无 cross-feature guard | **存在** | `claim.go` |
| AB-3 | "最多 3 次 subagent 调用" 无计数器 | **存在** | `task-executor.md:14` |
| AB-4 | doc-reviser 可编辑 DOC_DIR 外文件 | **存在** | `doc-reviser.md:49-50` |

---

## Pattern 10: Code Organization

| ID | 缺陷 | 状态 | 位置 |
|----|------|------|------|
| CO-1 | ~75 个文件平铺在 `internal/cmd/`——无包边界 | **存在** | `internal/cmd/` |
| CO-2 | `isBusinessTask` 在 `cmd/validate_index.go` 和 `pkg/task/add.go`（名 `isBusinessTaskID`）重复定义 | **存在** | 两处 |
| CO-3 | 依赖检查 `checkDependenciesMet`（claim）和 `checkUnmetDeps`（status）比 `pkg/task.GetUnmetDependencies` 更完整但留在 cmd 层 | **存在** | `claim.go`、`status.go` |
| CO-4 | `addFixTask`（quality_gate.go）重复调用 task.AddTask + CreateTaskMarkdown + EnsureForgeState，与 add.go 重叠 | **存在** | `quality_gate.go` vs `add.go` |
| CO-5 | `validator` struct（validate_index.go）是纯业务逻辑但放在 cmd 层 | **存在** | `validate_index.go` |

---

## Pattern 11: Naming Drift

| ID | 当前名称 | 问题 | 状态 | 建议名称 |
|----|----------|------|------|----------|
| NM-1 | `TestTaskDef` | 定义所有自动生成任务（含 spec、clean-code、doc-eval），不仅是 test | **存在** | `AutoGenTaskDef` |
| NM-2 | `testgen.go` | 文件生成所有自动生成任务类型，~40% 非测试 | **存在** | `autogen.go` |
| NM-3 | `isTestTaskID` | 对 `T-clean-code-1` 也返回 true | **存在** | `isPipelineTaskID` |
| NM-4 | `GenerateTestTaskMD` | 为所有自动生成任务生成 MD | **存在** | `GenerateAutoGenTaskMD` |
| NM-5 | `generateTestTasks` | 生成 spec/clean-code 任务 | **存在** | `generateAutoGenTasks` |
| NM-6 | `GetBreakdownTestTasks` / `GetQuickTestTasks` | 返回包含非测试任务的完整任务集 | **存在** | `GetBreakdownAutoTasks` / `GetQuickAutoTasks` |
| NM-7 | `TestInterfaces` | 应用界面类型，不是"测试界面" | **存在** | `Interfaces` |
| NM-8 | `ResolveFirstTestDep` | 为 clean-code 任务也解析依赖 | **存在** | `ResolveFirstAutoGenDep` |
| NM-9 | `TaskFromFile` | 无文件操作——纯 struct 转换 | **存在** | `ToTask` |
| NM-10 | `build.go` | "build" 暗示编译，实际是构建 task index | **存在** | `build_index.go` |
| NM-11 | `NewTestIndex` | 歧义——"测试任务 index"还是"用于测试的 index"？ | **存在** | `NewTaskIndexForTest` |
| NM-12 | `indexCmd.Long` | 描述说 "Test tasks auto-generated"，但实际生成 gates、specs 等 | **存在** | 更新描述 |
| NM-13 | `gen-test-cases` SKILL.md | 引用 `/breakdown-tasks` 但测试任务来自 `forge task index` | **存在** | 更新为 `forge task index` |

---

## Pattern 12: Error Handling

| ID | 缺陷 | 状态 | 位置 |
|----|------|------|------|
| EH-1 | `worktree.go` 4 个子命令完全不用 AIError——全部 `fmt.Errorf` | **存在** | `worktree.go` 全文 |
| EH-2 | `worktree.go` 手动写 stderr + 返回 error——用户看到双重错误输出 | **存在** | `worktree.go:67-68, 145-146, 214-215` |
| EH-3 | `submit.go` 锁冲突用 `fmt.Fprintln + os.Exit(1)` 绕过 `Exit()` | **存在** | `submit.go:91-97` |
| EH-4 | `quality_gate.go` 基础设施错误返回 nil + 自定义 JSON 输出 | **部分修复** | `quality_gate.go:64-111, 232` |
| EH-5 | `add.go` 丢弃 3 个 config/profile 读取错误（`_` 赋值） | **存在** | `add.go:206-220` |
| EH-6 | `claim.go:221` plain `fmt.Errorf` vs `claim.go:209` 同场景用 AIError | **存在** | `claim.go:209 vs 221` |
| EH-7 | `claim.go:83-84` `LoadState` 和 `index.ByID` 错误被丢弃 | **存在** | `claim.go:83-84` |

---

## Pattern 13: CLI UX

| ID | 缺陷 | 状态 | 位置 |
|----|------|------|------|
| UX-1 | `Run` vs `RunE` 混用——35 个 `Run` + 8 个 `RunE` | **存在** | 全局 |
| UX-2 | 两套 config init——`init.go` 用 huh TUI（含 auto），`config.go` 用 bufio（无 auto），产出不同 | **存在** | `init.go:174-295` vs `config.go:87-169` |
| UX-3 | `configGetCmd` 是唯一设 `SilenceErrors: true` + `SilenceUsage: true` 的命令 | **存在** | `config.go:35-36` |
| UX-4 | 12 个命令缺 `.Args` 验证——多余位置参数被静默忽略 | **存在** | `check_deps.go`、`cleanup.go`、`quality_gate.go`、`version.go` 等 |
| UX-5 | `verify_task_done.go` 是唯一用 exit code 2 的命令 | **存在** | `verify_task_done.go:33` |
| UX-6 | 输出格式不一致——部分命令无 `PrintBlock` 包装 | **存在** | 多文件 |
| UX-7 | `e2e_*` 子命令（现为 `test_*`）全部用 raw `os.Exit(1)` 而非 `Exit()` | **存在** | 5+ 个 test 文件 |

---

## Pattern 14: Skill Interface

| ID | 缺陷 | 状态 | 位置 |
|----|------|------|------|
| SI-2 | "流程描述" section 有 3 种命名：Process Flow、Architecture、Workflow | **存在** | 多个 skill |
| SI-3 | `task-executor` agent 缺 `inputs` frontmatter | **存在** | `agents/task-executor.md` |
| SI-6 | Skill Prerequisites/When-to-Use 排序不统一 | **存在** | 所有 skill |
| SI-7 | Skill step 编号起点不一致——有的从 Step 0，有的从 Step 1 | **存在** | 多个 skill |

---

## Pattern 15: Config & Schema

| ID | 缺陷 | 状态 | 位置 |
|----|------|------|------|
| CF-1 | 无 `config set` 命令 | **存在** | `cmd/config.go` |
| CF-2 | `config get` 只暴露 `auto.gitPush`——其余 auto 字段不可查询 | **存在** | `config.go:391-401` |
| CF-3 | `e2eprobe` 用手写 YAML 解析器（`ExtractYAMLStringField`） | **存在** | `e2eprobe/e2eprobe.go:74-84` |
| CF-4 | 常量 `ForgeDir`/`ForgeConfigFileName` 在 `profile/config.go` 和 `feature/constants.go` 重复定义 | **存在** | 两处 |
| CF-5 | 无 config schema 版本字段 | **存在** | `ForgeConfig` struct |
| CF-6 | `config init` 不暴露 `auto` block | **存在** | `config.go:157-168` |

---

## Pattern 16: Magic Values & Constants

| ID | 缺陷 | 状态 | 位置 |
|----|------|------|------|
| MV-1 | `maxFixTasksPerStep = 3` 硬编码 | **存在** | `quality_gate.go:30` |
| MV-2 | `probeTimeout = 5s` 硬编码 | **存在** | `e2eprobe/e2eprobe.go:61` |
| MV-3 | `defaultLockTimeout = 5s` 硬编码 | **存在** | `index/lock.go:16` |
| MV-4 | `time.Sleep(50ms)` 锁重试间隔硬编码 | **存在** | `index/lock.go:55` |
| MV-5 | `forensicLast` 默认 20 | **存在** | `forensic.go:68` |
| MV-6 | `submit.go` 硬编码中文字符串 `"无"` | **存在** | `submit.go:371,451,475` |
| MV-7 | `.forge` 目录名散落为字符串字面量（`worktree.go` 3处、`markers.go`、`root.go`、`ensure.go`） | **存在** | 多文件 |
| MV-8 | 新测试模型命令硬编码 `"tests"` 路径 | **新增** | `test_promote.go:27`、`verify.go:441` |
| MV-9 | 新测试模型命令硬编码 `"exit code 0/1"` 字符串匹配 | **新增** | `verify.go:263` |
| MV-10 | Journey 隔离硬编码 `"CLAUDE_PROJECT_DIR="` 和 `"forge-journey-"` | **新增** | `journey_isolation.go:123,188` |

---

## Pattern 17: Config & Doc Drift

| ID | 缺陷 | 状态 | 位置 |
|----|------|------|------|
| CD-2 | Config schema `project-type` enum 与 CLI 不一致 | **存在** | `forge-config.schema.json` vs `just.go` |
| CD-3 | guide.md `auto.e2eTest.quick` 默认值与代码矛盾 | **存在** | `guide.md` vs `config.go` |

---

## Pattern 18: Newly Introduced Issues

| ID | 缺陷 | 状态 | 位置 |
|----|------|------|------|
| NI-1 | `feature_complete.go` 状态大小写不一致——manifest `"completed"` vs proposal `"Completed"` | **存在** | `feature_complete.go:121,129` |
| NI-2 | `worktree.go remove` 提示 `--force` 但无此 flag | **存在** | `worktree.go:105` |
| NI-3 | `clean-code` SKILL.md scope detection 用 `grep` 解析 JSON | **存在** | `SKILL.md:70` |
| NI-4 | `clean-code` template summary.md 有 10 个字段但 SKILL.md 只定义 6 个 | **存在** | 模板 vs SKILL.md |
| NI-5 | `worktree.go` 三个子命令重复相同的 pre-flight check | **存在** | `worktree.go` |

---

## Pattern 19: Test Model Command Quality（PR #113 引入）

| ID | 缺陷 | 状态 | 位置 |
|----|------|------|------|
| TM-1 | `test promote` 无路径遍历校验——`journeyName` 可包含 `../` | **新增** | `test_promote.go:27` |
| TM-2 | `test verify` 解析失败静默返回零值 Contract——不报错 | **新增** | `verify.go:485` |
| TM-3 | `test verify` Fact Table 无条目时标记 Contract 为 OK 而非 unverifiable | **新增** | `verify.go:237` |

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

### RARELY USED — 保留但标记为调试工具

`forge task list-types`、`forge task migrate`、`forge probe`、`forge lesson`、`forge proposal`、`forge claude`、`forge config`——均为交互式调试/管理命令，skill 管道不调用。

---

## 统计

| 类别 | 数量 |
|------|------|
| 已修复 | 13 |
| 已移除（模型替代） | 4 |
| 仍存在 | 78 |
| 新增（PR 引入） | 6 |
| 死代码 | 9 |
