# Forge 插件深度审计报告

> 审计日期: 2026-05-15
> 审计范围: Forge CLI + Plugin Skill System + Task Workflow + Architecture + Config
> 方法: 5 个并行 subagent 对不同维度进行对抗性评估

---

## 审计摘要

| 类别 | P0 | P1 | P2 | 合计 |
|------|-----|-----|-----|------|
| 并发/数据完整性 | 3 | 3 | 2 | 8 |
| Token/上下文效率 | 3 | 3 | 2 | 8 |
| 质量门禁可靠性 | 1 | 4 | 2 | 7 |
| Skill 架构/DRY | 0 | 3 | 1 | 4 |
| 错误处理 | 0 | 2 | 0 | 2 |
| 测试覆盖 | 0 | 1 | 0 | 1 |
| 配置管理 | 0 | 1 | 3 | 4 |
| 代码质量 | 0 | 0 | 5 | 5 |
| **合计** | **7** | **17** | **15** | **39** |

---

## P0 — 关键攻击点

### P0-1: Claim 无文件锁 — 竞态条件

**位置**: `forge-cli/internal/cmd/claim.go:69-103`

`executeClaim` 加载 index、修改状态、保存，全程无 advisory lock。`submit.go` 正确使用了 `LockFile`（`submit.go:89`），但 `claim`、`status`、`migrate`、`add` 均未加锁。

两个 agent 同时 `forge task claim` 会抢到同一任务：

```
Agent A: LoadIndex → find task-1 (pending)
Agent B: LoadIndex → find task-1 (pending)
Agent A: set status=in_progress → SaveIndex
Agent B: set status=in_progress → SaveIndex  ← 覆盖 A 的写入
```

**对比**: `submit.go:89` 有 `lock, err := indexPkg.LockFile(indexPath)`，`claim.go:31 runClaim` 无任何锁。

**受影响命令**: `claim`、`status`、`migrate`、`add`、`quality-gate`（调用 addFixTask）

**修复建议**: 所有写入 `index.json` 的路径统一使用 `LockFile` + `SaveIndexAtomic`。

---

### P0-2: 大多数命令使用非原子 SaveIndex

**位置**: `forge-cli/pkg/task/index.go:24-34`

```go
func SaveIndex(path string, index *TaskIndex) error {
    data, _ := json.MarshalIndent(index, "", "  ")
    return os.WriteFile(path, data, 0644)  // 直接写入，崩溃=数据丢失
}
```

对比 `forge-cli/pkg/index/atomic.go:12-47` 的 `SaveIndexAtomic`（先写 tmp 再 rename）：

- 只有 `submit.go` 使用 `SaveIndexAtomic`
- 其余 6+ 个命令使用非原子 `SaveIndex`
- 进程中断（kill、OOM、timeout）会留下零字节或残缺的 `index.json`

**修复建议**: 统一所有 `SaveIndex` 调用为 `SaveIndexAtomic`，或将 `SaveIndex` 本身改为原子操作。

---

### P0-3: Quality Gate 在验证前消费状态

**位置**: `forge-cli/internal/cmd/quality_gate.go:83`

```go
_ = feature.ClearForgeState(projectRoot)  // 先清除
// ... 然后才运行 compile → fmt → lint → test → e2e
```

`allCompleted` 状态采用 consume-once 模式。如果 quality gate 在 line 83 之后崩溃：
- 状态已消费，无法重触
- 所有任务已完成，但质量门禁永远不会再次执行
- 必须手动完成一个新任务才能重新触发

**修复建议**: 将 `ClearForgeState` 移到所有步骤成功完成之后（after line 197），或在步骤失败时重新写入状态。

---

### P0-4: Quality Gate 首步失败即退出

**位置**: `forge-cli/internal/cmd/quality_gate.go:245`

`handleGateFailure` 用 `os.Exit(0)` 退出，只修复第一个失败步骤：

- compile 失败 → 只创建 compile 的 fix-task，lint/test/e2e 不执行
- 用户无法得知其他步骤是否正常
- 修复后需要重新走完整门禁才能发现下一个问题

**修复建议**: 改为收集所有失败步骤，一次性报告或为每个失败步骤创建 fix-task。

---

### P0-5: guide.md 注入所有 subagent

**位置**: `plugins/forge/hooks/session-start`、`plugins/forge/hooks/hooks.json:14-22`

`hooks.json` 的 `SubagentStart` hook 将 7.4KB guide.md 注入**每个** subagent：

- `doc-scorer` — 只需评分协议，不需要完整工作流图
- `doc-reviser` — 只需目录约定，不需要质量门禁协议
- `task-executor` — 只需质量门禁，不需要 eval 参数表

典型 `run-tasks` 工作流的 token 开销：

| 组件 | 大小 |
|------|------|
| 主会话: guide + run-tasks | ~16.8KB |
| 每个 subagent: guide + task-executor + skill prompt | ~20KB+ |
| 3 个 subagent 合计指令开销 | ~80,000 tokens |

**修复建议**: 按 subagent 类型提供精简版 guide，或在 session-start hook 中根据 agent 类型选择注入内容。

---

### P0-6: Skill 文件过大

| Skill | 大小 | 行数 | 问题 |
|-------|------|------|------|
| `gen-test-scripts/SKILL.md` | 27.4KB | 429 | 包含 5 个 profile 指令，单 profile 项目处理 80% 无关内容 |
| `breakdown-tasks/SKILL.md` | 25.6KB | 478 | 条件标签系统增加 LLM 解析开销 |
| `init-justfile/SKILL.md` | 18.7KB | 344 | 内联了 6 个完整 justfile 模板，模板文件已单独存在 |

加载 skill 时完整 SKILL.md 进入上下文窗口。`gen-test-scripts` 单个就消耗 ~7,000 tokens。

**修复建议**: 按 profile 拆分为片段，加载时按活跃 profile 选择注入；内联模板改为引用外部文件。

---

### P0-7: 无懒加载/选择性加载机制

profile 系统在 CLI 层感知活跃 profile（`forge profile`），但 SKILL.md 无条件包含所有 profile 的指令。没有机制只加载当前 profile 相关的部分。

**修复建议**: 引入 profile-aware 的 skill 片段加载：skill 定义骨架 + profile-specific 片段，加载时按 config 组合。

---

## P1 — 高优先级攻击点

### P1-1: validate-index.sh 调用 `task` 而非 `forge`

**位置**: `plugins/forge/scripts/validate-index.sh:18,24`

```bash
command -v task &> /dev/null    # line 18 — 检查旧名称
task validate "$FILE_PATH"      # line 24 — 调用旧命令
```

PostToolUse hook 在 Edit/Write 操作后触发，如果用户只安装了 `forge` 而非旧 `task` 二进制，hook 每次静默失败，index.json 变更不验证。

---

### P1-2: Claim 静默忽略损坏的 state.json

**位置**: `forge-cli/internal/cmd/claim.go:146-151`

```go
state, err := task.LoadState(statePath)
if err != nil {
    fmt.Fprintf(os.Stderr, "Warning: failed to load task state: %v\n", err)
    return false, false, nil  // 当作无状态继续
}
```

JSON 解析失败只打 warning，agent 重新 claim 新任务。可能丢弃一个正在执行的 in-progress 任务。

---

### P1-3: verify-task-done 缺状态时返回 nil

**位置**: `forge-cli/internal/cmd/verify_task_done.go:41-58`

当 project root 未找到、feature 未设置、或 state 缺失时，函数返回 `nil`（允许提交）。pre-commit hook 应拦截时放行。

---

### P1-4: Cleanup 删除 blocked 任务状态

**位置**: `forge-cli/internal/cmd/cleanup.go:65`

```go
if t.Status == "completed" || t.Status == "blocked" || t.Status == "rejected" {
    _ = os.Remove(statePath)
```

blocked 任务可能有活跃的 fix-task。删除状态切断了 blocked 任务与 fix-task 的关联。

---

### P1-5: autoRestoreSourceTask 不更新 state.json

**位置**: `forge-cli/internal/cmd/submit.go:222-236`

将 blocked 任务恢复为 pending，但不更新 `.forge/state.json`。如果 `allCompleted=true`，状态保持 true 即使出现了新的 pending 任务。

---

### P1-6: Feature 命令为 git 推断的 slug 自动创建目录

**位置**: `forge-cli/pkg/feature/feature.go:27-29`

```go
// Feature doesn't exist but we inferred it from git
if err := EnsureFeatureDir(projectRoot, feature); err == nil {
    return feature, nil
```

分支名 `fix/typo` 会自动创建 `docs/features/typo/` 及所有子目录。读操作产生副作用。

---

### P1-7: 依赖检查逻辑重复 6 次

以下函数实现近乎相同的依赖检查：

1. `checkDependenciesMet` — `claim.go:238-269`
2. `checkUnmetDeps` — `status.go:157-184`
3. `validateDependencies` — `validate_index.go:141-166`
4. `GetUnmetDependencies` — `pkg/task/add.go:337-370`
5. `validateWildcardSelfDeps` — `validate_index.go:258-282`
6. `validateLiveness` — `validate_index.go:467-528`

所有函数处理 `.x` 通配符模式、遍历 `TasksMap()`、检查 completed/skipped 状态。修 bug 需改 6+ 处。

---

### P1-8: add.go 零测试覆盖

**位置**: `forge-cli/internal/cmd/add.go`

动态创建任务的关键写入路径，包括 fix-task 创建、SourceTaskID 解析、BlockSource 逻辑、auto-ID 生成，全部无单测。

以下命令也无测试：`query.go`、`task_parent.go`、`version.go`、`profile.go`、`prompt_*.go`、`e2e_*.go`、`test_results.go`。

---

### P1-9: 必需 recipe 缺失时质量门禁通过

**位置**: `forge-cli/pkg/just/just.go:105-107`

```go
// Missing required recipe → WARNING but skip (does NOT fail)
if !has {
    if !step.Optional {
        fmt.Fprintf(os.Stderr, "WARNING: required recipe %q not found\n", step.Recipe)
    }
    continue  // 跳过，不算失败
}
```

没有 `compile` 或 `test` recipe 的项目，质量门禁静默通过。

---

### P1-10: E2E setup/probe 失败静默跳过

**位置**: `forge-cli/internal/cmd/quality_gate.go:159-180`

- `e2e-setup` 失败 → `e2eReady = false` → 跳过全部 e2e，只打 warning
- `ProbeServers` 失败 → 同上
- 不创建 fix-task
- 用户必须手动注意到 stderr warning

---

### P1-11: Record-missing 恢复不验证代码正确性

**位置**: `plugins/forge/commands/run-tasks.md:107-118`

当 subagent 完成但 record 文件缺失时，spawn 恢复 agent 只创建最小 record，不验证代码改动是否正确或完整。

---

### P1-12: Fix-task 上限(3)可导致静默卡死

**位置**: `forge-cli/internal/cmd/quality_gate.go:26`

`maxFixTasksPerStep = 3`，超限后打印手动创建提示。如果用户未监控终端输出，feature 表现为"完成"但门禁永远不过。

---

### P1-13: Manifest drift 无自动同步

**位置**: `plugins/forge/skills/breakdown-tasks/templates/manifest-update-tasks.md`

Manifest 更新是手动模板，无自动化验证。PostToolUse hook 只验证 `index.json`，不验证 `manifest.md`。skill 创建/重命名文件后 Documents 表可能过时。

---

### P1-14: Eval skill 60% 结构相同却无共享基类

6 个 eval skill 各自定义相同的：
- 架构图（相同 Mermaid 流程图）
- Orchestrator Iron Laws（仅文档类型名词不同）
- Step 2: Invoke Scorer（相同模式）
- Step 3: Decision Gate（逐字相同）
- Step 4: Invoke Reviser（仅路径替换不同）

~60 行 x 6 文件 = 360 行纯重复。

---

### P1-15: Report 模板 36 行 x 7 文件重复

**位置**: `skills/eval-*/templates/report.md`

7 个 eval report 模板共享：
- Frontmatter（8 行）
- Deductions 表格（4 行）
- 3x Attack Point 块（15 行）
- Previous Issues Check 表格（4 行）
- Verdict 部分（5 行）

合计 ~36 行 x 7 文件 = 252 行仅 scorecard 头和维度不同的重复。

---

### P1-16: Profile Resolution Step 0 重复 6 处

以下 SKILL.md 包含近乎相同的 "Step 0: Resolve Profile" 块（~15 行 x 6）：

- `quick-tasks/SKILL.md:20`
- `gen-test-scripts/SKILL.md:28`
- `eval-test-cases/SKILL.md:8`
- `breakdown-tasks/SKILL.md:25`
- `run-e2e-tests/SKILL.md:19`
- `graduate-tests/SKILL.md:18`

---

### P1-17: 版本三重分裂

| 组件 | 版本 | 位置 |
|------|------|------|
| Plugin manifest | `3.0.0-beta-6` | `plugins/forge/.claude-plugin/plugin.json` |
| CLI binary | `3.10.0` | `forge-cli/scripts/version.txt` |
| README badge | `2.16.1` | `README.md:5` |

`/upgrade-forge` 只更新 plugin manifest 和 marketplace JSON，不更新 CLI version.txt。

---

## P2 — 中等优先级攻击点

### P2-1: `rejected` 任务无逃逸路径

**位置**: `quality_gate.go:94-95`

rejected 任务阻止 all-completed 触发，但无机制将 rejected 状态转为其他状态。

---

### P2-2: Scope resolution 每次 gate step 重读配置

**位置**: `forge-cli/pkg/just/just.go:67-87`

`ResolveScope` 每次调用 `profile.ReadConfig()`，质量门禁序列中读取 4 次同一配置文件。配置中途变更可导致不同步骤使用不同 scope。

---

### P2-3: `checkDependenciesMet` O(n*m) 嵌套扫描

**位置**: `forge-cli/internal/cmd/claim.go:238-269`

对每个 eligible 任务的每个依赖，扫描整个 TasksMap。50 任务 x 3 依赖 = 7,500 次迭代。应预建 ID→status 索引。

---

### P2-4: `feature list` 每特性 6 次 I/O

**位置**: `forge-cli/internal/cmd/feature.go:217-268`

`discoverFeatures` 读取 manifest.md、index.json、prd-spec.md、tech-design.md、ui-design.md、results.json。无缓存无并行。

---

### P2-5: `FindProjectRoot` 每次遍历到文件系统根

**位置**: `forge-cli/pkg/project/root.go:41-114`

每个命令调用都从 cwd 遍历到 `/`。深层目录 ~10 次 stat() 调用。结果应缓存。

---

### P2-6: `add` 双写 index.json

**位置**: `forge-cli/internal/cmd/add.go:179-211`

`AddTask` 写入 index，然后 `BuildIndex` 读取所有 markdown 再次写入。两次非原子写入。

---

### P2-7: `fmt` 非阻塞可掩盖真实问题

**位置**: `forge-cli/pkg/just/just.go:23-29`

格式化失败只打 WARNING。如果 CI 环境中 fmt 是阻塞检查，本地通过但 CI 失败。

---

### P2-8: `type` 字段 omitempty 可丢值

**位置**: `forge-cli/pkg/task/types.go:94`

`json:"type,omitempty"` 意味着手动编辑 index.json 移除 type 后，Load/Save 往返会静默丢弃该字段。

---

### P2-9: Test graduation merge 累积冗余测试

**位置**: `plugins/forge/skills/graduate-tests/SKILL.md:110-135`

去重仅基于"完整标题字符串匹配"。不同 feature 的不同 TC ID 测试相同行为时，两套都保留，随时间膨胀回归测试套件。

---

### P2-10: Mermaid 图在 guide.md 中浪费 ~1.5KB token

**位置**: `plugins/forge/hooks/guide.md:35-62,72-79`

两个 mermaid 代码块 ~1.5KB，LLM 无法渲染。应替换为简洁文本描述或移除。

---

### P2-11: 硬编码 profile 名列表在 6+ 文件中

`web-playwright`、`go-test`、`maestro`、`java-junit`、`rust-test`、`pytest` 列表在 guide.md 和至少 6 个 SKILL.md 中重复。新增 profile 需改 6+ 文件。

---

### P2-12: `tmp_api_server.mjs` 含其他开发者硬编码路径

**位置**: 项目根目录

18KB Node.js 文件，line 9 有 `/Users/nasuki/zcode` 硬编码路径。应移除或 gitignore。

---

### P2-13: `settings.local.json` 含无关项目路径

**位置**: `.claude/settings.local.json:10`

```json
"Bash(grep -r \"test\" /Users/fanhuifeng/.../zcode/task-cli/...)"
```

引用完全不同的项目 `zcode/task-cli`，可能是复制粘贴遗留。

---

### P2-14: Config schema 与 Go struct 不一致

**位置**: `plugins/forge/references/shared/forge-config.schema.json` vs `forge-cli/pkg/profile/config.go:16-20`

JSON schema 声明 `additionalProperties: false` 且只定义 `test-profiles`，但 Go struct 接受 `project-type` 和 `capabilities`。schema 从未在运行时执行。

---

### P2-15: YAML 解析不拒绝未知字段

**位置**: `forge-cli/pkg/profile/config.go:55,75`

`gopkg.in/yaml.v3` 解析无 `DisallowUnknownFields()`。配置键拼写错误静默忽略。

---

## 修复优先级建议

### 第一阶段：数据安全（解决 P0-1, P0-2）

1. 为 `claim`/`status`/`migrate`/`add`/`quality-gate` 统一添加 `LockFile`
2. 将所有 `SaveIndex` 替换为 `SaveIndexAtomic`
3. 影响文件：`claim.go`、`status.go`、`migrate.go`、`add.go`、`quality_gate.go`、`pkg/task/index.go`

### 第二阶段：质量门禁可靠性（解决 P0-3, P0-4, P1-9, P1-10, P1-12）

1. 将 `ClearForgeState` 移到所有步骤成功后
2. 收集所有失败步骤而非首步退出
3. 缺失必需 recipe 应阻塞而非 warning
4. E2E setup/probe 失败创建 fix-task
5. fix-task 上限时主动通知用户

### 第三阶段：Token 效率（解决 P0-5, P0-6, P0-7）

1. 按 subagent 类型提供精简版 guide
2. 大 skill 按 profile 拆分为片段
3. 引导 profile-aware 的选择性加载
4. 移除 guide.md 中的 Mermaid 图

### 第四阶段：DRY 重构（解决 P1-7, P1-14, P1-15, P1-16）

1. 提取共享依赖检查函数
2. 创建 eval-base 共享模板
3. 提取 report template 公共部分
4. 创建 profile-resolution 共享片段

### 第五阶段：清理与加固（解决剩余 P1, P2）

1. 修复 `validate-index.sh` 的 `task` → `forge`
2. 同步版本号
3. 清理无关文件（`tmp_api_server.mjs`、`settings.local.json`）
4. 为 `add.go` 补充单测
5. 添加配置验证（schema enforcement、DisallowUnknownFields）
