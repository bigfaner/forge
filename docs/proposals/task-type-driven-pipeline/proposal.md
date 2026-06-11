---
created: 2026-05-14
author: faner
status: Draft
---

# Proposal: Task Type 驱动的自适应流水线

## Problem

当前 `forge task index` 对所有 feature 统一生成相同的自动化流水线（stage gates + summaries + test tasks），不区分任务类型。这导致：

1. **文档类 feature 生成无意义的测试任务**：如 `tui-ui-design` 这类只修改文档的 proposal，生成 `T-quick-1~5` 测试任务完全没有产出
2. **`type` 字段非必填**：`InferType` 对未知 ID fallback为 `implementation`，掩盖了信息缺失
3. **`--no-test` 手动 flag 维护成本高**：依赖 agent 或用户记住加 flag，不可靠

### Evidence

`tui-ui-design` proposal 产出的全是文档任务（修改 skill 文件、模板、rubric），不涉及任何代码编译或测试执行，但 pipeline 仍会生成 5 个测试任务。

## Proposed Solution

### 核心机制

**`isDocsOnlyFeature()` 自动检测**：在 `BuildIndex` 扫描完业务任务后，检查所有任务的 `type` 字段。如果全部不是 `implementation` 也不是 `fix`，则判定为 docs-only feature。

| Feature 类型 | 自动生成 |
|-------------|---------|
| 有 `implementation`/`fix` 任务 | gates + summaries + tests（不变） |
| 纯 `documentation` 任务 | 仅 `T-eval-doc`（文档质量评估） |

### 设计决策

#### D1: `type` 字段必填

`type` 从可选变为必填。业务任务 `.md` 文件的 frontmatter 缺少 `type` 时，`BuildIndex` 返回 hard error，中断 index 生成。

**理由**：type 是自适应流水线的决策依据。缺失 type 意味着无法正确判断 feature 性质，宁可中断也不能猜测。

**`InferType` 变更**：移除 fallback 到 `TypeImplementation` 的 default 分支。`InferType` 仅用于从 ID 模式推断自动生成任务（gates、summaries、test-pipeline 任务）的 type，不再为业务任务兜底。

#### D2: 新增 `documentation` 类型

| 类型 | 语义 | 质量门禁 | noTest |
|------|------|---------|--------|
| `implementation` | 实现代码 | `compile → fmt → lint → test` | false |
| `fix` | 修复 bug | `compile → fmt → lint → test` | false |
| `documentation` | 编写/修改文档 | 无（轻量自检替代） | true |
| `doc-evaluation` | 评估文档质量 | 无 | true |

`documentation` 类型的 prompt 模板结构：

```
Step 1: 理解任务（读取任务描述、参考文件）
Step 2: 执行文档工作（创建/修改文档）
Step 3: 自检（检查格式、交叉引用、术语一致性）
Step 4: 提交记录（forge task submit）
```

#### D3: `T-eval-doc` 自动评估任务

Docs-only feature 自动生成 `T-eval-doc` 任务，替代测试流水线的验证职能。

**参数**：

| 属性 | 值 |
|------|------|
| ID | `T-eval-doc` |
| Key | `eval-doc` |
| Type | `doc-evaluation` |
| noTest | true |
| 依赖 | 动态解析：最后一个业务任务 |

**评分**：1000 分制（8 维度 × 125 分）

| # | 维度 | 分值 |
|---|------|------|
| 1 | 结构完整性 | 125 |
| 2 | 逻辑一致性 | 125 |
| 3 | 可追溯性 | 125 |
| 4 | 准确性 | 125 |
| 5 | 完整性 | 125 |
| 6 | 术语一致性 | 125 |
| 7 | 格式规范 | 125 |
| 8 | 语言质量 | 125 |

**迭代机制**：
- 达标线：900/1000
- 最大迭代：3 轮（score → revise → re-score）
- 未达标处理：以最终分数完成任务，报告中标注"未达标"
- 评估逻辑直接写在 `doc-evaluation.md` prompt 模板中（不创建新 skill）

#### D4: 废弃 `--no-test` flag

硬删除 `--no-test` flag。自适应检测替代手动控制。

#### D5: 检测时机

在 `BuildIndex` 中，Step 5（扫描业务任务）之后、Step 6.5（生成 gates）之前，新增 `isDocsOnlyFeature()` 检查。判断逻辑只写一次，Step 6.5 和 Step 7 共享结果。

```
Step 5:   扫描业务 .md → 索引入 index（type 缺失则 error）
Step 5.5: isDocsOnlyFeature() 检查
Step 6:   检测孤儿
Step 6.5: 生成 stage gates + summaries（docs-only 时跳过）
Step 7:   生成 test tasks 或 T-eval-doc（docs-only 时生成 T-eval-doc）
```

### Developer Workflow — Before & After

#### Before（当前行为）

1. 所有 feature 统一生成 gates + summaries + tests
2. 文档 feature 的测试任务无意义但仍然生成
3. 需要手动 `--no-test` 跳过
4. `type` 缺失时静默 fallback 为 `implementation`

#### After（本提案实施后）

1. 文档 feature 自动检测，只生成 `T-eval-doc`
2. 代码 feature 行为不变
3. `type` 缺失直接报错，强制补全
4. 无需手动 flag

## Requirements Analysis

### Key Scenarios

1. **纯文档 feature**：所有任务 type=documentation → 跳过 gates/summaries/tests → 生成 T-eval-doc → agent 执行评估
2. **纯代码 feature**：有 implementation 任务 → 行为完全不变
3. **混合 feature**：有 implementation + documentation → 视为代码 feature，生成完整流水线
4. **fix-only feature**：只有 fix 任务 → 视为代码 feature（fix 修改代码）
5. **type 缺失**：任何业务任务缺少 type → BuildIndex hard error
6. **T-eval-doc 评估不达标**：3 轮后 < 900 分 → 以最终分数完成，报告标注
7. **现有 docs-only feature 重新 index**：旧任务缺少 type → hard error，需补全

### Constraints & Dependencies

- 依赖 `breakdown-tasks` 和 `quick-tasks` skill 生成的任务包含 `type` 字段
- `documentation` 和 `doc-evaluation` 类型通过 `noTest: true` 跳过 `submit.go` 的质量门禁
- `T-eval-doc` 的 rubric 硬编码在 prompt 模板中，修改 rubric 需更新模板文件
- 现有已完成 feature 的任务文件如果缺少 `type`，重新 index 会报错

### Non-Functional Requirements

| NFR Category | Requirement | Verification Method |
|-------------|-------------|-------------------|
| **兼容性** | 有代码任务的 feature 行为完全不变 | 回归测试 |
| **安全性** | type 缺失时不猜测，宁可报错 | 单元测试 |
| **扩展性** | 未来新增类型只需更新 `isDocsOnlyFeature()` 的判断条件和 prompt 模板 | 代码审查 |

## Feasibility Assessment

### Technical Feasibility

所有改动在 `forge-cli`（Go 代码）和 prompt 模板层面，不涉及 skill 架构变更。`isDocsOnlyFeature()` 逻辑简单（遍历已索引任务检查 type），`T-eval-doc` 遵循现有自动生成任务的模式。

### Resource & Timeline

| Step | Task | 预计时间 |
|------|------|---------|
| 1 | `types.go` 新增常量 + 更新 registry/validTypes | 0.5h |
| 2 | `build.go` hard error + `isDocsOnlyFeature()` + 条件跳过 + T-eval-doc 生成 | 1.5h |
| 3 | `testgen.go` 新增 `GetDocEvalTask()` | 0.5h |
| 4 | `infer.go` 移除 fallback | 0.5h |
| 5 | 新增 `documentation.md` prompt 模板 | 1h |
| 6 | 新增 `doc-evaluation.md` prompt 模板（含 1000 分 rubric） | 1.5h |
| 7 | `prompt.go` 注册新类型 | 0.5h |
| 8 | 删除 `--no-test` flag 及相关文档 | 0.5h |
| 9 | skill 模板更新（breakdown-tasks 加入 type 字段） | 0.5h |
| 10 | 单元测试 + 集成测试 | 2h |
| **Total** | | **~9h** |

## Scope

### In Scope

- `pkg/task/types.go` — 新增 `TypeDocumentation`、`TypeDocEvaluation` 常量
- `pkg/task/build.go` — hard error、`isDocsOnlyFeature()`、条件跳过、T-eval-doc 生成
- `pkg/task/testgen.go` — `GetDocEvalTask()`
- `pkg/task/infer.go` — 移除 fallback
- `pkg/prompt/prompt.go` — typeToTemplate 注册新类型
- `pkg/prompt/data/documentation.md` — 新模板
- `pkg/prompt/data/doc-evaluation.md` — 新模板（含 rubric）
- `internal/cmd/index.go` — 删除 `--no-test` flag
- `plugins/forge/skills/breakdown-tasks/templates/task.md` — 加入 `type: "{{TYPE}}"`
- `plugins/forge/skills/quick-tasks/SKILL.md` — 删除 `--no-test` 文档
- `plugins/forge/commands/quick.md` — 删除 `--no-test` passthrough
- `scripts/version.txt` — 版本号 bump

### Out of Scope

- all-completed hook 改动（项目级质量门禁对 docs-only feature 无害）
- `submit.go` 改动（通过 `noTest: true` 复用现有机制）
- rubric 动态化（硬编码在模板中）
- 已完成 feature 的旧任务补全 type（手动处理）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 旧 feature 重新 index 报错 | Medium | Low — 已完成的 feature 通常不会再 index | error 信息明确指出哪个任务缺少 type |
| Skill agent 忘记填 type | Low | Medium — index 失败 | 模板强制包含 `type: "{{TYPE}}"` 占位符 |
| T-eval-doc rubric 不适配所有文档类型 | Medium | Low — 评估不够精确 | 8 维度是通用的，不针对特定文档类型 |
| docs-only 判断遗漏新的代码修改类型 | Low | High — 该生成测试的没生成 | 新增类型时必须更新 `isDocsOnlyFeature()` |

## Success Criteria

- [ ] 有 `implementation`/`fix` 任务的 feature，自动生成行为与改动前完全一致
- [ ] 纯 `documentation` 任务的 feature，不生成 gates/summaries/tests，只生成 `T-eval-doc`
- [ ] 业务任务缺少 `type` 字段时，`forge task index` 返回 hard error 并指出具体文件
- [ ] `T-eval-doc` 执行 1000 分评估，达标线 900 分，最多 3 轮迭代
- [ ] `documentation` 类型任务执行时不运行质量门禁
- [ ] `--no-test` flag 已移除，使用时报 unknown flag 错误

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
