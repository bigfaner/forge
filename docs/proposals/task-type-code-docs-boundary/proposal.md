---
created: 2026-05-18
author: faner
status: Approved
---

# Proposal: Task Type Code/Docs Boundary

## Problem

Task `type` field 同时承担**描述职责**（任务在做什么）和**控制流职责**（是否走 quality-gate），但 agent 不清楚这个边界。结果：只修改 `.md` 文件的任务被标记为 `type: "enhancement"`（语义上确实是在增强功能），quality-gate 照跑——compile、fmt、lint、test 全部执行，全部无意义。

### Evidence

- `eval-adversarial-scorer` feature 的 2 个任务全部修改 `.md` 文件，但 `type: "enhancement"` + `scope: "backend"` → 走了完整 quality-gate
- `tui-ui-design` feature 的 6 个任务全部修改 `plugins/forge/skills/` 下的 markdown，标注 `type: "implementation"` → 同样走了 quality-gate
- 现有 `type-assignment.md` 定义 `documentation` 为 "non-compilable, non-runnable"，但缺乏执行层面的强制机制——agent 按任务意图而非产出物分类

### Urgency

中——不会导致错误产出，但浪费时间（quality-gate 在无代码变更时全部空转），且 agent 在 `forge profile` 等步骤上浪费注意力。

## Proposed Solution

**一条规则**：`type: "documentation"` 的任务跳过整个 quality-gate（compile + fmt + lint + test）。不需要额外设置 `noTest: true`。

### 核心改动

**1. type-assignment.md 增加 code/docs 判定规则**

```
判定标准：看产出物，不看意图。

Code types (走 quality-gate): feature, enhancement, cleanup, refactor, fix
Doc type   (跳过 quality-gate): documentation
Meta type  (特殊处理): gate

如果任务只修改 .md / .yaml / .json（非编译产物），无论意图是"实现新功能"还是"写文档"，type 必须为 documentation。
```

**2. Go 代码：`testableTypes` 扩展**

`forge-cli/pkg/task/build.go` 的 `testableTypes` 当前只有 `{feature, enhancement, fix}`，缺少 `cleanup` 和 `refactor`。导致 cleanup/refactor 任务被当作 docs-only，不走 quality-gate 也不生成 test pipeline。

修正：将 `cleanup` 和 `refactor` 加入 `testableTypes`，使 `IsTestableType` 和 `needsTestPipeline` 正确识别它们。

**3. Go 代码：submit-task quality-gate 改用 type 驱动**

`forge-cli/internal/cmd/submit.go` 当前仅用 `t.NoTest` 决定是否跳过 quality-gate。修正为：`!t.NoTest && IsTestableType(t.Type)` 才触发 quality-gate。同时，coverage auto-set（`coverage = -1.0`）也应对非 testable type 生效。

**4. quality-gate 执行点（guide.md）增加 type 检查**

guide.md quality-gate 协议统一规则：

> 当 `type === "documentation"` 时，跳过 compile / fmt / lint / test。等同于 `noTest: true`。

**5. task 生成器（quick-tasks、breakdown-tasks）强化分类**

当检测到 docs-only feature 时，使用 `task-doc.md` 模板（已设 `type: "documentation"`），而非 `task.md` 模板（默认 `type: "feature"`）。

### 设计决策

#### D1: Type 驱动而非 noTest 驱动

**选择**: `type: "documentation"` 作为 quality-gate skip 的唯一触发条件

**理由**: Agent 已经必须设置 type。让 type 同时传达 "这是文档任务" + "跳过 quality-gate"，减少认知负担和遗漏风险。`noTest: true` 保留作为显式覆盖机制（edge case：code task 但不需要测试）。

**否决方案**: 扩展 `noTest` 语义（`notest-docs-only-detection` 提案）。需要 agent 设置两个字段（type + noTest），增加遗漏风险。

#### D2: 按产出物分类而非意图分类

**选择**: 只要产出物是非编译文件（.md, .yaml, .json under skills/docs），type 就是 `documentation`

**理由**: 意图模糊（"改进 agent" 是 implementation 还是 documentation？），产出物客观可判断。Quality-gate 关心的也是产出物——compile/fmt/lint/test 对 .md 文件没有意义。

#### D3: 与现有 `noTest` 共存

**选择**: `type: "documentation"` 隐式等于 `noTest: true`，但 `noTest` 字段保留不废弃

**理由**: `noTest` 对 edge case 仍有用（code task 但不需要测试）。向后兼容——旧 index.json 无 `noTest` 字段时，type 仍能正确驱动行为。

## Requirements Analysis

### Key Scenarios

1. **纯 markdown feature**: 所有任务 `type: "documentation"` → 自动跳过 quality-gate → 生成 `T-eval-doc`
2. **纯 code feature**: 任务 `type: "feature"/"enhancement"` → 走 quality-gate → 生成 `T-quick-1~5`（行为不变）
3. **Mixed feature**: 部分任务 `type: "documentation"`，部分 `type: "feature"` → 各自按 type 决定是否 quality-gate → 整体非 docs-only，生成测试管线
4. **向后兼容**: 旧 index.json 无 type 字段 → 默认 `feature` → 走 quality-gate（行为不变）

### Constraints & Dependencies

- 修改 `plugins/forge/` 下的文件前必须先加载 `docs/conventions/forge-distribution.md`
- `type-assignment.md` 是 shared reference，被 quick-tasks 和 breakdown-tasks 引用
- `guide.md` 是全局 hook，所有 task-executing 工作流读取
- Go 代码变更涉及 `forge-cli/`，需遵循 TDD（RED → GREEN → REFACTOR），版本号需 bump
- `testableTypes` 扩展后，`isDocsOnly()`（quality_gate.go）和 `needsTestPipeline()`（build.go）自动受益——它们已使用 `IsTestableType`

## Alternatives & Industry Benchmarking

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| **Type 驱动（本提案）** | 单字段控制，agent 已必须设置 type | 修改 type 语义从"意图"到"产出物" | **Selected** |
| noTest 扩展（`notest-docs-only-detection`） | 不改变 type 语义 | Agent 需设两个字段 | Rejected |
| 自动检测文件扩展名 | 全自动 | 解析 markdown 脆弱，实现复杂 | Rejected |
| Do nothing | 零成本 | 空转 quality-gate，agent 困惑 | Rejected |

## Feasibility Assessment

### Technical Feasibility

8 个文件变更：5 个 markdown（skill/reference 文档）+ 2 个 Go 文件 + 1 个测试文件。Go 变更范围小（修改 map + 条件判断），测试用 table-driven 覆盖。

### Resource & Timeline

| Step | Task | 预计时间 |
|------|------|---------|
| 1 | 更新 `type-assignment.md` 增加判定规则 | 15min |
| 2 | 更新 `guide.md` quality-gate 协议增加 type 检查 | 15min |
| 3 | 更新 `submit-task` 检查 type=documentation | 15min |
| 4 | 更新 quick-tasks / breakdown-tasks 的 docs-only 检测引导 | 15min |
| 5 | Go: 扩展 `testableTypes`，修改 submit.go quality-gate skip 逻辑 | 30min |
| 6 | Go: 单元测试覆盖 | 30min |
| **Total** | | **~2h** |

## Scope

### In Scope

- `plugins/forge/references/shared/type-assignment.md` — 增加 code/docs 判定规则
- `plugins/forge/hooks/guide.md` — quality-gate 协议增加 `type: "documentation"` 跳过规则
- `plugins/forge/skills/submit-task/SKILL.md` — type 检查跳过 quality-gate
- `plugins/forge/skills/quick-tasks/SKILL.md` — 强化 docs-only 分类引导
- `plugins/forge/skills/breakdown-tasks/SKILL.md` — 强化 docs-only 分类引导
- `forge-cli/pkg/task/build.go` — `testableTypes` 增加 cleanup 和 refactor
- `forge-cli/internal/cmd/submit.go` — quality-gate skip 改用 `IsTestableType` 驱动
- Go 单元测试覆盖上述变更

### Out of Scope

- `noTest` 字段废弃
- `docs-only-fast-path` 提案的范围（技能文档中的 skip 行为记录）
- `notest-docs-only-detection` 提案（不同方案，保留但不实施）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Agent 仍按意图分类而非产出物 | M | M — type 错误导致空转 quality-gate | type-assignment.md 明确规则 + 示例 |
| Mixed feature 中文档任务被标 documentation 但整体不应跳过测试管线 | L | L — isDocsOnlyFeature 用 ALL 语义，不会误判 | 无需额外处理 |
| 与 `notest-docs-only-detection` 方向冲突 | L | L — 两个提案共存但不实施 notest 版本 | 明确标注 notest 版本为 deprecated |

## Success Criteria

- [ ] `type-assignment.md` 明确 "按产出物分类" 规则，包含 code types / doc type 分类表
- [ ] `guide.md` quality-gate 协议明确 `type: "documentation"` 跳过 quality-gate
- [ ] `submit-task` skill 检查 type 并跳过 quality-gate
- [ ] quick-tasks / breakdown-tasks 的 docs-only 分类引导更新
- [ ] 纯 markdown feature（如 eval-adversarial-scorer 风格）的任务自动标 `type: "documentation"`
- [ ] `testableTypes` 包含 cleanup 和 refactor，`IsTestableType` 正确返回 true
- [ ] submit.go quality-gate skip 同时检查 `noTest` 和 `IsTestableType`
- [ ] Go 单元测试覆盖 `testableTypes` 扩展和 submit type-based skip
