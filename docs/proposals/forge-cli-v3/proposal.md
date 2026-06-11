---
created: 2026-05-13
author: "HuiFeng Fan + Claude"
status: Draft
---

# Proposal: Forge CLI v3 — Task CLI 扩展为 Forge CLI

## Problem

Forge 插件的品牌是 "forge"，但 CLI 二进制名是 `task`。19 个子命令平铺在根级别，缺乏组织，AI agent 和用户都难以发现所需命令。命令命名存在歧义（`check` check 什么？`record` 名词还是动词？），justfile 中 5 个 e2e 命令重复了相同的 profile 检测逻辑。

### Evidence

- 二进制名 `task` vs 插件名 `forge` — 品牌不一致
- 19 个平铺子命令中 `check`/`validate`/`verify-completion` 语义重叠
- justfile 中 `test-e2e`、`e2e-setup`、`e2e-verify`、`e2e-compile`、`e2e-discover` 各自重复 ~6 行 profile 检测 bash 代码
- AI agent 看到命令名无法推断用途（`record`、`prompt`、`all-completed`）

### Urgency

v3.0.0 是主版本升级，是做破坏性重命名的最佳时机。推迟意味着更多的代码和文档引用 `task` 命令，迁移成本持续增长。

## Proposed Solution

将 `task` CLI 重构为 `forge` CLI，按功能分组子命令，优化命名使其自解释（AI-friendly），并将 justfile 中的 e2e 命令迁移到 CLI 内。

### 最终命令结构

**5 个命令组：**

```
forge task claim              # 认领下一个可用任务
forge task submit             # 提交执行结果 + 更新任务状态
forge task status             # 查看/更新任务状态
forge task query              # 查询任务详情
forge task check-deps         # 检查任务依赖完整性
forge task validate-index     # 校验 index.json 结构和语义
forge task verify-task-done   # pre-commit 校验任务完成
forge task add                # 动态添加新任务
forge task index              # 构建/重建 index.json
forge task migrate            # 迁移 index.json schema
forge task list-types         # 列出所有支持的任务类型

forge e2e run                 # 运行 e2e 测试（profile-aware）
forge e2e setup               # 安装 e2e 依赖
forge e2e verify              # 检查 VERIFY 标记
forge e2e compile             # 编译检查 e2e 测试
forge e2e discover            # 列出所有 e2e 测试用例

forge forensic search         # 搜索历史会话
forge forensic extract        # 提取会话证据
forge forensic subagents      # 列出子代理记录

forge profile set             # 设置测试 profile
forge profile detect          # 检测可用 profile
forge profile get             # 获取 profile 配置

forge prompt get-by-task-id   # 获取任务的 agent prompt
```

**5 个顶层命令：**

```
forge feature                 # 设置/显示当前 feature 上下文
forge probe                   # HTTP 健康检查
forge cleanup                 # hook: 清理已完成任务状态
forge quality-gate            # hook: 质量门禁（编译+lint+测试+e2e）
forge verify-task-done        # hook: pre-commit 校验
forge version                 # 版本信息
```

### 命名变更对照表

| 旧命令 | 新命令 | 变更理由 |
|--------|--------|----------|
| `task` | `forge` | 品牌统一 |
| `task record` | `forge task submit` | 消除名词/动词歧义，"提交"更贴切 |
| `task check` | `forge task check-deps` | 明确检查对象是依赖 |
| `task validate` | `forge task validate-index` | 明确校验对象是 index.json |
| `task verify-completion` | `forge task verify-task-done` / 顶层 `forge verify-task-done` | 明确校验对象是任务完成状态 |
| `task all-completed` | `forge quality-gate` | 实际动作是跑质量门禁，不是查询 |
| `task prompt` | `forge prompt get-by-task-id` | AI 语境下 "prompt" 太泛 |
| `task template` | **删除** | 无直接使用者 |
| (justfile) `e2e-*` | `forge e2e *` | 消除 profile 检测代码重复 |
| (justfile) `probe` | `forge probe` | 迁移到 CLI |
| (新增) | `forge task list-types` | 列出所有任务类型，辅助 agent 理解 |

### Innovation Highlights

这是标准的 CLI 重构模式，核心创新在于：
- **AI-first 命名（可度量）**：命令名自解释，AI agent 无需读文档即可理解用途。度量方法：将旧命令列表（`task check`, `task record`, `task all-completed` 等）和新命令列表（`forge task check-deps`, `forge task submit`, `forge quality-gate` 等）分别提供给 LLM，要求其为 10 个任务场景选择正确命令。验收标准：新命名正确率 >= 9/10，旧命名正确率 <= 7/10。命名灵感来自通用 CLI 设计原则（gh/kubectl 的 noun-first 分组）而非特定 AI 工具调用规范——MCP 和 function calling 的工具命名建议（如 Anthropic MCP 文档中的 `domain_action` 格式、OpenAI 的 `verb_object` 推荐）与我们的设计方向一致，但这是事后验证而非设计约束。核心洞察是：agent 的命令发现依赖 `--help` 输出的分组结构，而非逐个猜测命令名；noun-first 分组将候选集从 19 个缩减到每组 3-11 个，显著降低选择错误率。代价：`forge task claim` 比 `task claim` 多一层前缀，但通过前缀过滤的净收益为正。
- **从 justfile 到 CLI 的 e2e 迁移**：将 profile-aware 的 e2e 编排逻辑从 bash 脚本提升为 Go 代码，消除重复，提高可测试性

## Requirements Analysis

### Key Scenarios

- AI agent 调用 `forge task claim` 获取任务 → `forge task submit` 提交结果
- AI agent 调用 `forge prompt get-by-task-id` 获取 agent prompt
- Hook 自动调用 `forge cleanup`、`forge quality-gate`、`forge verify-task-done`
- 开发者调用 `forge e2e run --feature <slug>` 运行 feature 级 e2e 测试
- 开发者调用 `forge task list-types` 查看支持的任务类型
- CI 调用 `forge e2e compile && forge e2e run` 执行测试

### Non-Functional Requirements

- **启动延迟**：`forge --help` 响应时间 <= 当前 `task --help` 基线 + 50ms（基线值：迁移前测量 `task --help` 3 次取中位数，记录为 T_baseline；验收时 `forge --help` 3 次中位数 <= T_baseline + 50ms）
- **行为等价**：所有已迁移命令的退出码与原命令一致（0 成功、1 业务错误、2 参数错误）；stdout 格式不变（仅 `--help` 输出因分组结构变化除外）
- **二进制体积**：编译产物大小增量 <= 500KB（基线：迁移前 `task` 二进制大小；验收时 `forge` 二进制 <= 基线 + 500KB）
- **构建时间**：`go build` 增量不超过 10 秒（基线：迁移前测量值）

### Error & Edge-Case Scenarios

- **旧命令别名（old binary aliasing）**：用户系统中可能残留 `task` 二进制或 shell alias。处理策略：v3 不提供向后兼容 shim；在 `task` 二进制最后版本的 `--help` 输出中添加迁移提示（deprecation notice），指向 `forge` 命令。验收标准：`task --help` 输出包含 "请使用 forge 替代" 提示。
- **Profile 检测失败**：`forge e2e run` 依赖 `.forge/config.yaml` 中的 profile 配置。若文件缺失或 profile 字段为空，命令应退出码 1 并输出明确的错误信息（"未检测到测试 profile，请运行 forge profile detect"），而非静默使用默认值。验收标准：删除 config.yaml 后 `forge e2e run` 退出码为 1，stderr 包含 "profile" 关键字。
- **Flag 冲突（重命名后）**：命令重命名可能导致与 Cobra 全局 flag 的冲突（如 `forge task submit --help` vs `forge --help`）。验收标准：`go vet ./...` 无错误，每个子命令的 `--help` 正确显示该子命令的 flags 而非全局 flags。
- **并发执行**：多个 AI agent 同时调用 `forge task claim` 时，应依赖已有的文件锁机制（index.json.lock），不应出现任务被重复认领。验收标准：并发 2 个 `forge task claim` 进程，两者不会认领同一任务（行为等价，无代码变更——仅验证重命名未破坏已有锁机制）。

### Constraints & Dependencies

- Go 1.25 + Cobra 框架
- 无向后兼容要求（v3.0.0 clean break）
- 需要同步更新 hooks.json、justfile、所有文档
- 需要更新 23 个 skills 目录中的命令引用（完整清单：`grep -rl 'task ' plugins/forge/skills/ --include='*.md'`，当前匹配 23 个文件）

## Alternatives & Industry Benchmarking

### Industry Solutions

**kubectl（verb-resource 模式）**

kubectl 采用 `kubectl <verb> <resource>` 二级结构（如 `kubectl get pods`、`kubectl describe node`）。早期版本（v1.0 前）命令为扁平的 `kubectl getpods`、`kubectl stop` 等，v1 引入分组后命令数从 ~30 扁平命令增长到 100+ 而仍保持可发现性。其核心教训：verb-first 让用户按意图（get/create/delete）发现命令，resource 作为第二参数缩小范围。kubectl 的 `--help` 按资源类型分组展示，与 Cobra 的 `CommandGroups` 机制直接对应。Forge CLI 的 task 组采用类似的 verb-first 子命令（claim/submit/query），但不需要 resource 层级——因为 forge 的操作对象固定为 task。

**gh（extension 可扩展模型）**

GitHub CLI 采用 `gh <noun> <verb>` 结构（如 `gh pr create`、`gh repo clone`），并通过 `gh extension` 机制允许社区扩展命令树。早期 gh 仅支持 `pr` 和 `repo` 两组，后续按用户需求逐步增加 `issue`、`actions`、`run` 等组。其教训：(1) noun-first 分组让 `--help` 输出按业务领域组织，比 verb-first 更适合命令数 > 20 的场景；(2) 顶层命令数控制在 ~10 以内是 `--help` 可读性的关键阈值。Forge CLI 的 5 组 + 5 顶层命令（共 10 个 `--help` 入口）直接参照这一阈值。

**设计选择：noun-first（gh 模式）vs verb-first（kubectl 模式）**

Forge CLI 选择 noun-first 分组（`forge task claim` 而非 `forge claim task`），理由：(1) 命令数 19 个，按业务域分组后每组 3-11 个子命令，符合 gh 的分组规模；(2) AI agent 通过前缀 `forge task` 即可过滤出所有任务相关命令，无需遍历所有 verb；(3) Cobra 的 `CommandGroups` 原生支持 noun-first 分组的 help 输出。

参考：kubectl commands (<https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands>)，gh CLI reference (<https://cli.github.com/manual/>)

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 品牌不一致、命令混乱持续恶化 | **Rejected:** 问题会随命令增多而加剧 |
| 仅重命名二进制 | — | 最小改动 | 不解决命令组织和命名问题 | **Rejected:** 治标不治本 |
| 保留 task 名 + 分组 | — | 减少改动 | 仍然与 forge 品牌不一致 | **Rejected:** 错过主版本升级的窗口 |
| **全量重构为 forge CLI** | gh noun-first 分组模式（见上方分析） | 品牌统一、分组清晰、命名自解释、e2e 代码去重 | 改动面大（CLI + hooks + skills + docs） | **Selected:** noun-first 分组适配 19 命令规模，v3.0.0 是做破坏性变更的时机 |

## Feasibility Assessment

### Technical Feasibility

完全可行。Cobra 原生支持 command groups。Go module rename 是标准操作。justfile → CLI 的 e2e 迁移是将 bash 逻辑转写为 Go。

### Resource & Timeline

预计 4-8 小时工作量，主要在：
- Cobra 命令重组和重命名（1-2h）
- Go module + 目录重命名（0.5h）
- e2e 命令从 bash 迁移到 Go（1-2h）
- hooks + justfile 更新（0.5h）
- skills 引用更新（1h）
- 文档更新（1h）

### Dependency Readiness

无外部依赖。所有变更在项目内部完成。

## Scope

### In Scope

- 二进制 `task` → `forge` 重命名
- Go module `task-cli` → `forge-cli` 重命名
- 目录 `task-cli/` → `forge-cli/` 重命名
- Cobra 命令分组（task/e2e/forensic/profile/prompt）
- 命令重命名（record→submit, check→check-deps, validate→validate-index, verify-completion→verify-task-done, all-completed→quality-gate, prompt→prompt get-by-task-id）
- 删除 `template` 命令
- 新增 `forge task list-types` 命令
- e2e 命令从 justfile 迁移到 CLI（run/setup/verify/compile/discover）
- `probe` 从 justfile 迁移到 CLI
- 更新 hooks.json 中的命令引用
- 更新 justfile（移除迁移的 recipe）
- 更新 23 个 skills 中的 `task` 命令引用为 `forge` 命令（识别命令：`grep -rl 'task ' plugins/forge/skills/ --include='*.md'`）
- 更新文档（OVERVIEW.md, WORKFLOW.md, 中文版）
- 更新 pkg/version/version.go 中的 Name 常量
- 更新所有 Go 测试中的命令引用

### Out of Scope

- 新增业务功能
- 改变命令的执行逻辑或输出格式
- prompt 组的扩展命令（如 list、validate）
- `forge feature` 的子命令拆分（保持当前的 get/set 双行为）
- justfile 中 build/dev/lint/test 等 Go 语言 recipe 的迁移

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Skills/hooks/tests 中遗漏 `task` 引用导致运行时错误 | M | H | 新增 `just check-stale-refs` CI target：grep `exec.Command("task"` 和 `"task "` 模式，扫描 skills/、hooks/、`*_test.go`，匹配则退出码 1；集成到 `just lint` 流水线 |
| e2e 命令迁移时 profile 检测逻辑不一致 | M | M | 保留 justfile recipe 作为 fallback，并行验证一个迭代周期（sprint）；fallback recipe 在首次成功运行 `forge e2e run` 后的下一个 sprint 起点移除 |
| Go module rename 导致 go.sum 冲突 | L | L | 清理 go.sum 后重新生成 |
| 文档不同步 | M | L | `just check-stale-refs` 同时扫描 `docs/` 和 `*.md`，阻塞合并 |

## Success Criteria

**覆盖全部 In Scope 条目，每条可客观验证。**

### 命令结构与分组（scope #4）

- [ ] `forge --help` 显示 5 个命令组（task/e2e/forensic/profile/prompt）+ 5 个顶层命令（feature/probe/cleanup/quality-gate/verify-task-done/version），共 10 个入口
- [ ] `forge task --help` 显示 11 个子命令，名称与命名变更对照表一致
- [ ] `forge e2e --help` 显示 5 个子命令（run/setup/verify/compile/discover）
- [ ] `forge forensic --help` 显示 3 个子命令（search/extract/subagents）
- [ ] `forge profile --help` 显示 3 个子命令（set/detect/get）

### 重命名（scope #1-3, #5-6）

- [ ] `which forge` 返回有效路径，`which task` 不再指向 forge-cli 二进制
- [ ] Go module 名为 `forge-cli`，目录为 `forge-cli/`，`go mod tidy` 无错误
- [ ] 全局搜索 `task ` 在 skills/ 和 hooks/ 中零匹配（排除非 CLI 语义的 "task" 出现，如 "task type" 等自然语言）

### 新增与删除（scope #7, #8）

- [ ] `forge task list-types` 输出至少包含当前 index.json 中定义的所有 task type，退出码 0
- [ ] `forge task template` 返回 "unknown command" 错误（退出码非 0），确认已删除

### e2e 与 probe 迁移（scope #9, #10）

- [ ] `forge e2e run --feature <slug>` 等价于原 justfile `test-e2e`，定义：对全部 5 个 profile（default/web/api/plugin/full），(a) 退出码一致（justfile `test-e2e` 退出码 = `forge e2e run` 退出码），(b) stdout 包含相同的测试名称集合（忽略 ANSI 转义和行序差异），(c) `forge e2e run -v`（verbose 模式）输出包含 "profile: <name>" 行，确认使用了 config.yaml 中的 profile 而非硬编码默认值
- [ ] `forge probe` 成功执行 HTTP 健康检查，退出码与原 justfile `probe` recipe 一致

### hooks 与 justfile 更新（scope #11, #12）

- [ ] hooks.json 中所有命令引用为 `forge` 前缀（grep `hooks.json` 中 `task ` 返回零行）
- [ ] justfile 中 `e2e-*` 和 `probe` recipe 已移除或替换为 `forge` 命令调用

### skills 更新（scope #13）

- [ ] 23 个 skill 文件中 `task claim`/`task submit`/`task record` 等旧命令引用已全部替换为 `forge` 对应命令（验证：`grep -rl 'task claim\|task record\|task status\|task query\|task check\b\|task validate\|task verify\|task cleanup\|task all-completed\|task prompt\b\|task feature\|task add\b\|task index\b\|task migrate\|task profile\|task template\|task version' plugins/forge/skills/ --include='*.md'` 返回零匹配）

### 文档与代码更新（scope #14, #15, #16）

- [ ] OVERVIEW.md、WORKFLOW.md 及其中文版中 `task` 命令引用已替换为 `forge` 命令
- [ ] `pkg/version/version.go` 中 `Name` 常量值为 `"forge"`
- [ ] 所有 Go 测试通过（`go test ./...` 退出码 0）

### Go 测试命令引用更新（scope #17）

- [ ] grep `exec.Command("task"` 在所有 `*_test.go` 文件中返回零匹配（确认测试代码已引用 `forge` 而非 `task`）
- [ ] grep `"task "` 模式在 `*_test.go` 的命令构造参数中返回零匹配（排除注释和字符串中的自然语言 "task"）

### 集成验证

- [ ] `just compile && just fmt && just lint && just test` 全部通过
- [ ] `go vet ./...` 无错误，无 flag 冲突警告
- [ ] `just check-stale-refs` 通过（零残留 `task` 命令引用）

### NFR 验证

- [ ] `forge --help` 响应时间 <= T_baseline + 50ms（T_baseline 为迁移前 `task --help` 三次中位数）
- [ ] `forge` 二进制大小 <= 原始 `task` 二进制 + 500KB

## Next Steps

- Proceed to `/write-prd` to formalize requirements
