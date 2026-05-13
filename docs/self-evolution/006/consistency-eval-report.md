# Forge vs Task-CLI 一致性评估报告

**日期**: 2026-05-13
**范围**: forge plugin (skills/agents/commands/hooks) vs task-cli (Go 源码)
**方法**: 3 个并行子代理分别审查 CLI 源码、Skill 引用、模板格式，然后交叉比对

---

## 总结

| 类别 | 发现数 |
|------|--------|
| P0 (功能损坏) | 0 |
| P1 (文档不一致) | 2 |
| P2 (维护风险/小瑕疵) | 4 |

**整体评价**: forge 与 task-cli 的一致性**非常高**。所有 22 个 CLI 命令在 skill 中的引用均有效，状态枚举、数据结构、文件路径完全对齐。发现的问题集中在文档描述与 CLI 实际行为的细微差异。

---

## P1: 文档不一致

### P1-1: Fix-task ID 格式描述错误

**位置**: `plugins/forge/skills/breakdown-tasks/SKILL.md:407`

**文档说**: "Auto-generated fix-task IDs follow the `disc-N` format (e.g., `disc-1`, `disc-2`)"

**CLI 实际行为**: 使用 `--template fix-task` 时，ID 前缀是 `fix`（定义于 `pkg/template/template.go:31`），生成 `fix-1`, `fix-2`。`disc-N` 是 `task add` 不带 `--template` 时的通用前缀。

**证据**:
- `pkg/template/template.go:27-31`: `IDPrefix: "fix"`
- `internal/cmd/add.go:43`: `"id", "", "Custom task ID (auto-generated as disc-N if omitted)"`

**影响**: Agent 看到 `disc-N` 文档后可能对 `task query` 返回的 `fix-N` ID 感到困惑。不影响功能（`task add --template fix-task` 正确生成 `fix-N`），但破坏了可理解性。

**修复**: 将 breakdown-tasks/SKILL.md 第 407 行的 `disc-N` 改为 `fix-N`。

---

### P1-2: Quality Gate 中 fmt 步骤的行为描述不准确

**位置**: `plugins/forge/hooks/guide.md` — Quality Gate Protocol

**文档说**: "Quality gate sequence: `just compile → just fmt → just lint → just test`. On failure: ... fmt → blocked (toolchain issue)"

**CLI 实际行为**: `fmt` 被标记为 `Optional: true, Blocking: false`（`pkg/just/just.go:24`）。fmt 失败时 CLI 输出 WARNING 但**继续执行后续步骤**，不会阻塞。

**证据**:
- `pkg/just/just.go:24`: `{Name: "fmt", Optional: true, Blocking: false}`
- `pkg/just/just.go:120`: `"WARNING: non-blocking gate step %q failed\n"`

**影响**: Agent 可能误认为 fmt 失败会导致任务阻塞，执行不必要的等待或重试。

**修复**: guide.md 的 Quality Gate 部分应标注 fmt 为 "(non-blocking, warning only)"。

---

## P2: 维护风险 / 小瑕疵

### P2-1: InferType() 双重维护依赖

**位置**:
- `task-cli/pkg/prompt/prompt.go` — InferType() 源码
- `plugins/forge/skills/breakdown-tasks/SKILL.md` — Type assignment rules
- `plugins/forge/skills/quick-tasks/SKILL.md` — Type assignment rules

**问题**: 两个 skill 文件中写明 "These rules mirror InferType() in task-cli/pkg/prompt/prompt.go — both must stay in sync"。如果 CLI 修改了类型推断逻辑（如新增 ID pattern），skill 文档不会自动更新。

**风险**: 中等。当 task-cli 发布新类型或修改 pattern 时，skill 生成的任务可能缺少正确的 `type` 字段，导致 `task prompt` 选择错误模板。

**建议**: 在 task-cli 的 CLAUDE.md 中添加 doc-sync 规则：修改 InferType() 时必须同步更新 skill 文档，或考虑让 `task index` 命令自动推断 type（正在 task-index-command 分支实现）。

---

### P2-2: CLI 输出解析的脆弱性

**位置**: 多个 skill 文件

**问题**: Skills 通过 grep/sed 解析 CLI 的 stdout 输出（如 `grep '^FEATURE:'`, `grep 'PROFILE: (none)'`）。如果 task-cli 修改输出格式，这些解析会静默失败。

| Skill | 解析模式 |
|-------|---------|
| run-e2e-tests | `task feature 2>/dev/null \| grep '^FEATURE:' \| sed 's/^FEATURE:[[:space:]]*//'` |
| 多个 skill | `task profile` 输出检查 `PROFILE: (none)` |
| breakdown-tasks/quick-tasks | `task add` 输出检查 `ACTION: ADDED` / `ACTION: SKIPPED` |

**风险**: 低。CLI 使用结构化 block 输出（`---` 分隔），格式相对稳定。但缺少 contract test 保证。

**建议**: 考虑为 task-cli 的输出格式添加 `--json` flag（部分命令已支持），或至少在测试中覆盖输出格式的稳定性。

---

### P2-3: record.json 必填字段语义差异

**位置**:
- `plugins/forge/skills/record-task/SKILL.md` — 文档
- `task-cli/internal/cmd/record.go:249-298` — 源码

**差异**:

| 字段 | Skill 文档语义 | CLI 实际行为 |
|------|---------------|-------------|
| `summary` | 必填 | **必填** (hard-required) ✅ |
| `testsPassed/testsFailed/coverage` | "must come from actual output" (语义上是强制的) | 仅在 `completed + coverage >= 0 + 两者都为 0` 时拒绝 (无测试证据) |
| `acceptanceCriteria` | 文档列出但未明确标注必填 | `recommended`（缺失仅 warning），仅当存在 unmet AC 时拒绝 |
| `keyDecisions` | 未明确 | `recommended`（缺失仅 warning） |

**风险**: 低。Skill 文档比 CLI 更严格，不会导致功能问题。但如果 Agent 遵循文档的严格语义，会在 CLI 不要求的情况下花费额外精力收集指标。

**建议**: record-task/SKILL.md 应区分 "hard-required" 和 "recommended" 字段，与 CLI 对齐。

---

### P2-4: `task index` 命令尚未集成到 skill 工作流

**位置**: `task-cli/internal/cmd/index.go` (task-index-command 分支)

**问题**: `task index` 命令已实现（扫描 .md 文件生成 index.json），但 `breakdown-tasks` 和 `quick-tasks` skill 仍通过 skill 内部逻辑手动构建 index.json。未来需要将 index 构建委托给 CLI。

**风险**: 低。当前不影响功能。属于已知的待整合项（见 task-index-command 分支的 plan）。

---

## 一致性确认项（无问题）

以下维度经交叉验证完全一致：

| 维度 | 状态 |
|------|------|
| 22 个 CLI 命令及其 flags | 全部正确引用 ✅ |
| 状态枚举 (6 值) | skills/hooks/CLI 完全对齐 ✅ |
| 类型枚举 (11 值) | skills/index.schema/CLI 完全对齐 ✅ |
| 优先级枚举 (P0/P1/P2) | 完全对齐 ✅ |
| 文件路径约定 | tasks/, records/, process/, testing/ 全部对齐 ✅ |
| 状态转换规则 | terminal/completed guards 对齐 ✅ |
| Auto-downgrade (completed→blocked) | 文档与 CLI 一致 ✅ |
| --force 无法覆盖 auto-downgrade | 文档与 CLI 一致 ✅ |
| Auto-restore (fix完成→source恢复) | 文档与 CLI 一致 ✅ |
| Wildcard dependencies (.x) | validate/check/claim 全部正确处理 ✅ |
| Scope resolution (frontend/backend/mixed) | guide.md 与 just.go 一致 ✅ |
| Fix-task 创建 (--block-source) | execute-task/run-tasks/breakdown-tasks 全部正确 ✅ |
| Hook 集成 (cleanup/all-completed) | hooks.json 与 CLI 行为一致 ✅ |
