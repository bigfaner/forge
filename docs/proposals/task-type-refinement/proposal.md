---
created: "2026-05-16"
author: faner
status: Draft
---

# Proposal: Task Type Refinement — 细化业务任务类型

## Problem

当前 forge 的业务任务类型只有两个：`implementation` 和 `fix`。`implementation` 覆盖了从"加新 API 端点"到"删几个测试文件"的所有代码变更，语义过于宽泛。

这导致两个结构性问题：

### 问题一：cleanup/refactor feature 生成无意义的测试 pipeline

`isDocsOnlyFeature()` 的判断逻辑只区分"有 `implementation`/`fix`"和"纯 `documentation`"两种情况。但现实中存在第三种：**有代码变更但无可测运行时行为的 feature**（cleanup、refactor）。

典型案例：`e2e-test-quality-cleanup` feature 删除低质量测试文件、修改 SKILL.md。它的 4 个业务任务都是 `type: implementation`，所以 `isDocsOnlyFeature()` 返回 false，系统生成了完整的 T-quick-1~6 测试 pipeline。但这些测试要么验证的是静态文本（antipattern #6），要么根本没有可测的运行时行为——**一个清理低质量测试的 feature 自己生成了低质量测试**。

### 问题二：所有代码任务共享同一套执行策略

`implementation` 类型在 task-executor prompt 里走 TDD 流程。但 cleanup 任务（删除死代码、修改 markdown）不需要先写失败测试再实现。类型信息没有渗透到执行层，导致：

- cleanup 任务被要求走 TDD（无意义的开销）
- fix 任务的 TDD 纪律完全依赖 skill prompt，type 系统没有参与

### Evidence

- `e2e-test-quality-cleanup` 生成了 6 个 T-quick 测试任务，全部需要 skip
- `isDocsOnlyFeature()` 只检查 `TypeImplementation` 和 `TypeFix`，无法识别 cleanup/refactor
- `implementation.md` prompt 模板对所有代码变更使用相同的 TDD 步骤

## Proposed Solution

### 核心变更：细化 `implementation` 为 4 个具体类型

| 新类型 | 替代 | 语义 | 可测运行时行为 | 执行策略 | Quality Gate |
|--------|------|------|---------------|---------|-------------|
| `feature` | `implementation` | 新增运行时行为 | 是 | 实现 → gate | 正常执行 |
| `enhancement` | `implementation` | 增强已有行为 | 是 | 实现 → gate | 正常执行 |
| `cleanup` | `implementation` | 删除死代码/文件、修复测试 | 否 | 改善 → gate | 正常执行 |
| `refactor` | `implementation` | 重构不改行为 | 否 | 重构 → gate | 正常执行 |

`implementation` 废弃。`fix`、`documentation` 保持不变。

### 设计决策

#### D1: Type 信号来源 — proposal intent + task 级覆盖

Proposal frontmatter 新增 `intent` 字段，声明 feature 的主导意图：

```yaml
---
status: Draft
intent: fix   # fix | feature | enhancement | cleanup | refactor
---
```

`/quick-tasks` 和 `/breakdown-tasks` 生成任务时：
1. 读取 proposal intent，作为默认 type
2. 个别 task 可在 frontmatter 中显式覆盖 `type`

这避免了 AI 推断 type 的不可靠性，同时保留了混合 feature 的灵活性。

#### D2: 测试 pipeline 生成逻辑

`isDocsOnlyFeature()` 扩展为 `needsTestPipeline()`，逻辑从"是否纯文档"变为"是否存在可测类型"：

```go
testableTypes = map[string]bool{
    TypeFeature:      true,
    TypeEnhancement:  true,
    TypeFix:          true,
}

func needsTestPipeline(tasks map[string]Task) bool {
    for _, t := range tasks {
        if isAutoGenTaskID(t.ID) {
            continue
        }
        if testableTypes[t.Type] {
            return true
        }
    }
    return false
}
```

三档决策：

| Feature 任务类型 | 测试 Pipeline | T-eval-doc |
|-----------------|-------------|------------|
| 全部 `documentation` | 不生成 | 生成 T-eval-doc |
| 有可测类型（feature/enhancement/fix） | 生成 | 不生成 |
| 只有 cleanup/refactor | 不生成 | 不生成 |

`docsOnly` 概念拆为两个独立判断：`needsTestPipeline()` 和 `needsDocEval()`。

#### D3: 执行策略 — 每个 type 独立 prompt

每个业务类型有独立的 prompt 模板：

| Type | Prompt 模板 | 关键步骤 |
|------|------------|---------|
| `feature` | `data/feature.md` | 实现功能 → quality gate |
| `enhancement` | `data/enhancement.md` | 增强功能 → quality gate |
| `fix` | `data/fix.md` | TDD：写失败测试复现 bug → 修代码 → 验证 |
| `cleanup` | `data/cleanup.md` | 改善技术债务 → quality gate |
| `refactor` | `data/refactor.md` | 重构 → quality gate（行为不变验证） |
| `documentation` | `data/documentation.md` | 文档工作 → 自检（无 quality gate） |

`implementation.md` 废弃，替换为 `feature.md` 和 `enhancement.md`。

"修测试"任务归 `cleanup`——修复测试本身不是用户可感知的行为 bug，不适用 TDD。

#### D4: 动态任务 type 由 quality gate 失败步骤决定

质量门禁失败时动态创建的 fix/cleanup 任务，type 由失败步骤确定性地决定：

| 失败步骤 | 动态任务 type |
|----------|-------------|
| `compile` | `fix` |
| `fmt` | `cleanup` |
| `lint` | `cleanup` |
| `test`（单元） | `fix` |
| `test`（e2e，all-completed） | `fix` |

不需要推断。哪一步断了，type 就定了。

#### D5: Type 偏移记录

Executor 在执行中发现 type 不匹配时（如 `fix` 任务调查后发现是 flaky test），直接按实际 type 的流程处理，并在 `forge task submit` 的 record 中记录偏移：

```markdown
## Type Reclassification
- Original: fix (quality gate: test failure)
- Actual: cleanup (investigation: flaky test, not introduced by this feature)
- Reason: e2e test TestTC_003_Login has race condition in assertion timing
```

#### D6: `noTest` 行为

`cleanup`/`refactor` 不设 `noTest`——它们正常走 quality gate。区别只在 quick-tasks 不为它们生成 test pipeline。

这跟 `documentation` type 不同：documentation 才设 `noTest: true`，因为它不改代码。

### 完备性验证

6 个业务 type 对常见开发场景的覆盖情况：

| 场景 | Type | 验证方式 | 判定 |
|------|------|---------|------|
| 新 API 端点 | `feature` | test pipeline + quality gate | 覆盖 |
| 给 API 加分页 | `enhancement` | test pipeline + quality gate | 覆盖 |
| 修登录 bug | `fix` (TDD) | test pipeline + quality gate | 覆盖 |
| 删死代码 | `cleanup` | quality gate 回归 | 覆盖 |
| 重构命名 | `refactor` | quality gate（行为不变） | 覆盖 |
| 修 flaky test | `cleanup` | quality gate | 覆盖 |
| 写文档 | `documentation` | T-eval-doc | 覆盖 |
| 加测试覆盖 | `cleanup` | quality gate | 覆盖 |
| 升级依赖 | `cleanup` | quality gate 回归，不生成新测试 | 覆盖 |
| 移除废弃 API | `cleanup` | 删 API + 删对应测试，quality gate 验证剩余测试通过 | 覆盖 |
| 数据迁移脚本 | `refactor` | 重构数据存储形式，行为不变，quality gate 验证 | 覆盖 |
| 配置变更（CI/justfile） | `cleanup` | quality gate 回归 | 覆盖 |
| 性能优化 | `enhancement` | 可测的运行时行为改善，test pipeline | 覆盖 |

关键推理：

- **升级依赖**：维护性质，不改业务行为。需要的是回归测试（quality gate 跑现有测试），不是生成新 e2e 测试。`cleanup` 正确。
- **移除废弃 API**：删代码和删对应测试是同一个 `cleanup` task 的一体两面。不需要 test pipeline 来"反向生成"测试——被删的测试本身就是要清理的对象。quality gate 验证剩余测试通过。
- **数据迁移脚本**：改变数据表现形式，不改变业务语义。这是 `refactor` 的定义——重构数据层，行为不变。

结论：**6 个业务 type 覆盖所有场景，不需要新增 type。** 边界场景的风险不在 type 分类，而在现有测试覆盖率——那是 code review 和人工判断的范畴，不是 type 系统的职责。

### 与已有 Proposal 的关系

| Proposal | 关注点 | 与本提案的关系 |
|----------|--------|--------------|
| `task-type-driven-pipeline` | docs-only 检测、T-eval-doc、`--no-test` 废弃 | 本提案将 `isDocsOnlyFeature()` 扩展为 `needsTestPipeline()`，覆盖更多场景 |
| `typed-task-dispatch` | CLI 驱动 prompt 合成、`task prompt` 命令 | 本提案新增 4 个业务 type 和对应 prompt 模板，`task-type-refinement` 是 `typed-task-dispatch` 的前置：type taxonomy 必须先确定，prompt 合成才能按 type 分派 |

实施顺序：`task-type-refinement` → `task-type-driven-pipeline` → `typed-task-dispatch`

## Scope

### In Scope

- `pkg/task/types.go` — 新增 `TypeFeature`、`TypeEnhancement`、`TypeCleanup`、`TypeRefactor` 常量，废弃 `TypeImplementation`
- `pkg/task/build.go` — `isDocsOnlyFeature()` → `needsTestPipeline()` + `needsDocEval()`
- `pkg/prompt/` — 新增 `feature.md`、`enhancement.md`、`cleanup.md`、`refactor.md` 模板，废弃 `implementation.md`
- `internal/cmd/` — 动态任务创建时根据失败步骤决定 type
- `internal/cmd/migrate.go` — 新增迁移规则：旧 `implementation` → `feature`（保守默认）
- proposal frontmatter — 新增 `intent` 字段
- `/quick-tasks` skill — 读取 proposal intent，传播到 task type
- `/breakdown-tasks` skill — 同上
- record 格式 — 新增 Type Reclassification 块

### Out of Scope

- `task prompt` 命令实现（属于 `typed-task-dispatch`）
- prompt 模板的详细内容（属于 `typed-task-dispatch`）
- `--no-test` flag 废弃（属于 `task-type-driven-pipeline`）
- `error-fixer` 废弃（属于 `typed-task-dispatch`）
- all-completed hook 改动

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 旧 index.json 含 `implementation` type | 高 | 中 | `task migrate` 将 `implementation` 映射到 `feature`（保守默认） |
| Agent 生成任务时选错 type | 中 | 中 | proposal intent 作为默认值减少错误；task type 枚举有明确语义 |
| cleanup prompt 模板对"修测试"场景不够具体 | 低 | 低 | "修测试"本质上就是改代码后验证，cleanup prompt 足以覆盖 |
| prompt 模板冷启动质量不足 | 中 | 中 | 参照现有 `implementation.md` 逐步拆分，每类型迁移后立即验证 |

## Success Criteria

- [ ] `TypeImplementation` 从 `types.go` 中移除，所有引用替换为新 type
- [ ] `needsTestPipeline()` 对 cleanup-only feature 返回 false，对 feature/enhancement/fix feature 返回 true
- [ ] `needsDocEval()` 对 documentation-only feature 返回 true，对其余返回 false
- [ ] proposal frontmatter 含 `intent` 字段时，`/quick-tasks` 生成的任务继承对应 type
- [ ] 个别 task frontmatter 显式设置 `type` 时覆盖 proposal intent
- [ ] 动态 fix 任务：compile/test 失败 → `fix`，fmt/lint 失败 → `cleanup`
- [ ] record 中 Type Reclassification 块在 type 偏移时存在，无偏移时不存在
- [ ] 现有 `implementation` type 的 index.json 通过 `task migrate` 正确迁移到 `feature`
- [ ] 所有 forge-cli 单元测试通过，coverage ≥ 80%

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
