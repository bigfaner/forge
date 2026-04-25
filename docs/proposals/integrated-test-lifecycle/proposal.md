---
created: 2026-04-24
author: fanhuifeng
status: Draft
---

# Proposal: Integrated Test Lifecycle

## Problem

测试生成（`gen-test-cases`、`gen-test-scripts`）与任务生命周期完全脱节：三个 skill 形成流水线，但没有任何机制保证它们在任务完成前被执行。`task all-completed` 执行 e2e 测试时，无法感知脚本是否已生成；测试失败后没有恢复路径，开发者只能手动介入。

### Evidence

- `gen-test-cases` → `gen-test-scripts` → `run-e2e-tests` 是三个独立 skill，没有任何 hook 或任务将它们串联到 `task claim` / `task record` 流程中。
- `task all-completed` 检测到 `testing/scripts/package.json` 后直接运行 `npm run test:all`，但没有检查脚本是否由 `gen-test-scripts` 生成过（可能是空目录或过期脚本）。
- 测试失败时 `task all-completed` 打印错误后退出，没有任何自动恢复机制，开发者需要手动分析、手动修复、手动重新触发。
- 项目现有 7 个 feature 目录，无一包含 `testing/` 子目录，测试生成步骤跳过率为 100%（7/7）。其中 `feat-auto-test-after-run-all-tasks` 是引入 `task all-completed` Stop hook 的 feature 本身——该 feature 完成时 `testing/scripts/` 不存在，`task all-completed` 跳过了 e2e 测试门控，自动化闭环在第一次实际使用时即已失效。

### Urgency

`task all-completed` 作为 Stop hook 是 zcode 自动化闭环的最后一道门。如果 e2e 测试脚本缺失或失败后无法自动恢复，这道门形同虚设——要么永远不触发（脚本未生成），要么触发后卡死（失败无路径）。当前 7 个 feature 的跳过率已达 100%；每新增一个 feature 就多一次无测试门控的合并风险。`feat-auto-test-after-run-all-tasks` 修改了 Stop hook 本身的逻辑——若该 feature 引入了 hook 的回归缺陷，缺陷将在零 e2e 覆盖的情况下合并，此后所有 feature 的自动化闭环都将静默失效。

## Proposed Solution

将测试生成变为**可见的、可追踪的任务**，纳入标准任务生命周期；同时将测试脚本按**页面/组件**归档，而非按 feature 存放，使其成为长期可维护的回归套件。

### 核心机制

**1. `/breakdown-tasks` 追加标准测试任务**

在所有业务任务之后，自动追加两个固定任务：

```
T-test-1: 生成 e2e 测试用例（调用 gen-test-cases skill）
T-test-2: 生成 e2e 测试脚本（调用 gen-test-scripts skill，依赖 T-test-1）
```

- 依赖关系：T-test-1 依赖最后一个业务任务；T-test-2 依赖 T-test-1
- AI agent 像执行普通任务一样认领并完成它们（`task claim` → 调用 skill → `task record`）
- 测试任务完成后，`testing/scripts/` 目录必然存在，`task all-completed` 可以安全执行

**2. 测试脚本按页面/组件归档（毕业模型）**

测试脚本的组织方式应与代码一致——按"测什么"而非"为什么写"。Feature 是开发历史，不是系统结构。

```
开发期（feature 隔离）：
  docs/features/<slug>/testing/scripts/   ← gen-test-scripts 生成到这里

task all-completed 首次成功（毕业）：
  tests/e2e/ui/<page>/                    ← 按页面归档
  tests/e2e/api/<resource>/              ← 按资源归档
  tests/e2e/cli/<command>/               ← 按命令归档
```

`gen-test-cases` 为每个测试用例增加 `target` 字段（如 `ui/login`、`api/auth`、`cli/deploy`），毕业时按此字段将脚本分发到 `tests/e2e/` 对应目录。

每个测试用例携带一个 **test ID**，格式为 `<type>/<target>/<slug>`，其中 `slug` 由测试用例标题规范化生成（小写、空格转连字符、去除标点）。例如标题 "Login with valid credentials" 在 `ui/login` 下的 ID 为 `ui/login/login-with-valid-credentials`。ID 由 `gen-test-cases` 自动生成，不需要人工编写；标题不变则 ID 稳定，标题变更则视为新用例。

同一页面被多个 feature 覆盖时，相同 test ID 用新脚本替换旧脚本，不同 test ID 追加——测试自然聚合，不重复。

**毕业触发机制**：`task all-completed` 通过检查 `tests/e2e/` 目录是否已包含当前 feature slug 的迁移标记文件（`tests/e2e/.graduated/<slug>`）来区分首次与非首次成功。首次成功时创建该标记文件并执行迁移；后续成功时跳过迁移直接 exit 0。标记文件内容为迁移时间戳，便于审计。

`docs/features/<slug>/testing/` 保留 `test-cases.md` 和 `results/` 作为可追溯性记录，随 feature 归档。

**长期收益**：`tests/e2e/ui/login/` 包含所有 feature 对登录页的覆盖，未来可作为全量回归套件独立运行。

**3. `task all-completed` 失败时追加修复任务**

e2e 测试失败后，每个失败的测试用例独立创建一个修复任务，ID 格式为 `fix-e2e-{round}-{index}`：

```json
{
  "id": "fix-e2e-1-1",
  "title": "修复 e2e 测试失败: Login with invalid credentials",
  "priority": "P0",
  "dependencies": [],
  "status": "pending",
  "file": "fix-e2e-1-1.md"
}
```

- `e2eRound` 字段持久化到 `index.json`，记录当前修复轮次（0 = 尚未失败）
- 每轮最多 3 次，超限后 agent 停止并提示人工介入
- 追加后通过 Stop hook JSON `{"decision": "block", "reason": "..."}` 驱动 agent 继续工作
- agent 完成修复后再次触发 `task all-completed`，形成闭环

**4. `run-e2e-tests` skill 保留，用于手动深度分析**

开发中随时可调用，提供截图、详细日志、逐步分析——这是 `task all-completed` 轻量执行不具备的能力。两者互补，不重叠。

### 数据流

```
/breakdown-tasks
    → 业务任务 1..N
    → T-test-1: gen-test-cases  (依赖任务 N，输出含 target 字段)
    → T-test-2: gen-test-scripts (依赖 T-test-1，输出到 testing/scripts/)

task all-completed (Stop hook)
    → 检查所有任务 completed/skipped
    → 运行 e2e 测试（testing/scripts/）
        ├── 成功（首次）→ 写 latest.md → 毕业：迁移脚本到 tests/e2e/<type>/<target>/ → exit 0
        ├── 成功（非首次）→ 写 latest.md → exit 0
        └── 失败 → 写 latest.md → 追加 fix-e2e-N 任务 → exit 1
                                        ↓
                                  agent 认领 fix-e2e-N
                                  读 latest.md → 修复 → task record
                                        ↓
                                  task all-completed 再次触发
```

## Alternatives Considered

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | 零改动 | 测试生成持续被跳过；失败无恢复路径；`task all-completed` 形同虚设 | Rejected: 核心问题未解决 |
| 方案 B：独立 `task test-phase` 命令 | 显式、可控 | 需要新增 1 个 CLI 子命令、1 套调用约定、1 份文档说明；"何时调用 test-phase vs all-completed"本身是新概念负担。与 `task all-completed` 的重叠在于：两者都需要检查所有任务已完成、都需要运行 e2e 测试——方案 B 实质上是把 `all-completed` 的测试逻辑拆成两步调用，但没有消除任何逻辑，只是增加了一个调用点。方案 A 复用现有 `task claim/record` 原语，零新 CLI 概念 | Rejected: 方案 B 引入 1 个新命令 + 1 套调用约定，而方案 A 在现有原语内解决同一问题 |
| 方案 C：约定优于配置（纯文档） | 零改动 | 没有强制力，容易被跳过，与现状无本质区别 | Rejected: 无法解决根本问题 |
| 方案 A（本提案） | 复用现有任务机制；零新概念；失败自动闭环 | 始终追加导致无 e2e 需求的 feature 产生两个 skipped 任务（可接受的噪音） | Accepted |

## Scope

### In Scope

- `/breakdown-tasks` skill：在任务列表末尾追加 `T-test-1`（gen-test-cases）和 `T-test-2`（gen-test-scripts）两个标准任务
- `gen-test-cases` skill：为每个测试用例增加 `target` 字段（如 `ui/login`、`api/auth`、`cli/deploy`）
- `task all-completed` CLI：
  - e2e 测试失败时，向 `index.json` 追加修复任务并 exit 1
  - e2e 测试首次成功时，按 `target` 字段将脚本迁移到 `tests/e2e/<type>/<target>/`（毕业）
- 定义修复任务的标准格式（id 命名规则、file 字段指向 latest.md）
- 定义毕业规则：相同 test ID 更新，新 ID 追加，`docs/features/<slug>/testing/` 保留为归档记录
- 更新 `task-cli/docs/OVERVIEW.md` 和 `docs/todos/test-capture-design.md` 反映新行为
- 更新 `gen-test-cases` 和 `gen-test-scripts` skill 的 Prerequisites 说明，明确它们也可以作为任务被 agent 调用

### Effort Estimate & Phasing

8 个变更点分布在 3 个子系统，按复杂度分两个阶段交付：

**Phase 1（~3 engineer-days）：任务注入 + 失败恢复**
- `/breakdown-tasks` 追加 T-test-1/T-test-2（prompt 修改，低风险）
- `gen-test-cases` 增加 `target` 字段（prompt 修改）
- `task all-completed` 失败时追加 fix-e2e 任务（Go 代码，写 index.json）
- 定义 fix-task 格式与追加上限规则
- 更新 Prerequisites 说明

**Phase 2（~3 engineer-days）：毕业模型**
- `task all-completed` 首次成功时按 `target` 迁移脚本（文件迁移 + 冲突解决逻辑，最高复杂度）
- 定义毕业规则（相同 test ID 更新，新 ID 追加）
- 更新 `OVERVIEW.md` 和 `test-capture-design.md`

Phase 1 可独立上线并立即解决跳过率问题；Phase 2 依赖 Phase 1 的 `target` 字段。

### Out of Scope

- 修改 `gen-test-scripts`、`run-e2e-tests` skill 的核心逻辑
- 单元测试生成（仅 e2e）
- CI/CD 集成（`tests/e2e/` 目录的 CI 接入留给用户配置）
- 修复任务的自动执行（agent 认领后手动执行，不自动）
- 测试用例的版本管理或 diff
- `tests/e2e/` 的全量回归运行机制（本提案只建立目录结构，运行留给后续）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 修复任务无限追加（每次失败都加一个） | 中 | 中：index.json 膨胀，agent 困惑 | 追加前检查是否已有 pending 的 fix-e2e 任务，有则跳过 |
| `/breakdown-tasks` 对无 e2e 需求的 feature 也追加测试任务 | 中 | 低：多余任务被 skipped | 始终追加（与核心机制 1 一致）；agent 执行 T-test-1 时若 PRD 无 UI/API/CLI 需求，将任务标记为 skipped 并说明原因，T-test-2 同步 skipped |
| agent 修复任务后测试仍失败，陷入循环 | 低 | 高：无限循环 | 限制 fix-e2e 任务最多追加 3 次，超出后打印警告并 exit 0 |
| 毕业时 `tests/e2e/` 中同名文件冲突（两个 feature 覆盖同一页面同一 test ID） | 中 | 中：旧测试被覆盖 | 相同 test ID 视为更新（新版本替换旧版本）；不同 test ID 追加；迁移前打印 diff |
| `task all-completed` 追加任务需要写 index.json，可能与并发操作冲突 | 低 | 中 | 单进程写入，原子替换（写临时文件后 rename） |

## Success Criteria

- [ ] 运行 `/breakdown-tasks` 后，index.json 末尾包含 `T-test-1` 和 `T-test-2` 两个任务，依赖关系正确
- [ ] agent 执行 `T-test-1` 后，`testing/test-cases.md` 中每个测试用例含 `target` 字段和 test ID；执行 `T-test-2` 后，`testing/scripts/package.json` 存在
- [ ] `task all-completed` e2e 首次成功后，`tests/e2e/<type>/<target>/` 目录下出现对应脚本文件，`tests/e2e/.graduated/<slug>` 标记文件存在，`docs/features/<slug>/testing/scripts/` 保留不删除
- [ ] 两个不同 feature 覆盖同一页面时，`tests/e2e/ui/<page>/` 中不存在两个 test ID 相同的脚本文件；若两个 feature 的测试用例标题不同，则各自的 test ID 不同，两个文件均保留
- [ ] `task all-completed` 在 e2e 测试失败时，index.json 中出现新的 `fix-e2e-N` 任务（pending 状态），且 `file` 字段指向 `testing/results/latest.md`
- [ ] 同一 feature 连续失败时，不会重复追加 fix-e2e 任务（已有 pending 任务则跳过）
- [ ] fix-e2e 任务追加上限为 3 次，超出后 `task all-completed` 打印警告并 exit 0
- [ ] `run-e2e-tests` skill 行为不变，可独立调用
- [ ] `task-cli/docs/OVERVIEW.md` 包含对 `tests/e2e/<type>/<target>/` 目录结构和毕业机制的说明
- [ ] `gen-test-cases` 和 `gen-test-scripts` skill 的 Prerequisites 说明中明确列出"可作为 `/breakdown-tasks` 追加的标准任务被 agent 调用"

## Next Steps

- Proceed to `/write-prd` to formalize requirements
