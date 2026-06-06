---
created: "2026-06-06"
author: "faner"
status: Draft
intent: "enhancement"  # valid values: new-feature | enhancement | refactor | cleanup | fix | doc
---

# Proposal: Surface-Key Test Pipeline Alignment

## Problem

Forge 的测试 pipeline 中，`gen-test-scripts` 任务的命名和输出目录使用 `surface-type`（如 `api`），而 `run-tests` 使用 `surface-key`（如 `backend`）。这导致多 surface 项目中，生成的测试文件被放到 `tests/<journey>/` 根目录而非 `tests/<surfaceKey>/<journey>/`，与项目 convention 不一致。

### Evidence

pm-work-tracker 项目（`backend=api`, `frontend=web`）实际遇到此问题：
- gen-test-scripts 生成 `tests/item-deletion/step1_delete_main_item_api.spec.ts`
- 项目 convention 要求 `tests/api/item-deletion/`
- 需要手动 `mv` 到正确位置，已记录为 lesson (`gotcha-gen-test-scripts-output-dir.md`)

代码层面不一致：
- `pipeline.go` 第 741 行：gen-test-scripts 用 `per-surface-type` expansion，文件名 `gen-test-scripts-api.md`
- `pipeline.go` 第 749 行：run-tests 用 `per-surface-key` expansion，文件名 `run-test-backend.md`

### Urgency

已造成实际困扰。每个多 surface 项目都会遇到，当前 workaround 是手动移动文件。随着多 surface 项目增多，问题会持续出现。

## Proposed Solution

统一整个测试 pipeline 使用 `surface-key` 作为任务命名和输出目录的依据：

1. **pipeline registry**：gen-test-scripts 从 `per-surface-type` 改为 `per-surface-key`
2. **输出目录**：多 surface → `tests/<surfaceKey>/<journey>/`，单 surface → `tests/<journey>/`
3. **全量同步**：所有引用 `tests/<journey>/` 的 skill 文件、模板、规则文件统一更新

### Innovation Highlights

这是一个对齐性修复，将分散的两个维度（type vs key）统一到一个维度（key）。Forge 的 surface 模型中，key 是用户定义的标识符（如 `backend`、`frontend`），type 是技术分类（如 `api`、`web`）。测试目录和任务命名应跟随用户可见的 key，因为：
- justfile recipe 已经用 key（`backend-test`，而非 `api-test`）
- 项目的 `tests/` 目录按功能分区（`tests/api/`、`tests/web/`），分区名对应 key
- type 可以重复（两个 surface 都是 `api`），key 天然唯一

## Requirements Analysis

### Key Scenarios

- **单 surface（scalar）**：`surfaces: tui` → 输出 `tests/<journey>/`，任务名 `gen-test-scripts.md`（无后缀）
- **单 surface（named）**：`surfaces: [{key: app, type: tui}]` → 输出 `tests/<journey>/`（单 surface 不加目录层），任务名 `gen-test-scripts.md`（无后缀）
- **多 surface**：`surfaces: [{key: backend, type: api}, {key: frontend, type: web}]` → 输出 `tests/backend/<journey>/` + `tests/frontend/<journey>/`，任务名 `gen-test-scripts-backend.md` + `gen-test-scripts-frontend.md`
- **多 surface 同 type**：`surfaces: [{key: admin, type: api}, {key: public, type: api}]` → 两个独立任务，输出到 `tests/admin/<journey>/` 和 `tests/public/<journey>/`

### Non-Functional Requirements

- **向后兼容**：单 surface 项目行为不变（输出仍是 `tests/<journey>/`，任务名无后缀）
- **全量同步**：变更必须覆盖所有引用测试目录的文件，不允许部分更新导致新的不一致

### Constraints & Dependencies

- Forge plugin 分发模型：skill 文件修改后需要验证在用户项目中的行为
- init-justfile 模板中的 just recipe 路径需要与新的目录结构匹配
- `forge surfaces` 文本输出解析规则不变（key=type 格式）

## Alternatives & Industry Benchmarking

### Industry Solutions

多模块项目的测试目录组织通常按模块分区（Maven 的 `module/src/test/`，Nx 的 `apps/<name>/tests/`），而非按技术类型分区。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 每个多 surface 项目都需手动修复 | Rejected: 已造成实际困扰 |
| 仅改 SKILL.md 输出目录 | — | 最小改动 | pipeline registry 仍用 type 命名，不一致持续存在 | Rejected: 治标不治本 |
| 仅改 pipeline expansion | — | 任务名对齐 | 输出目录仍错误，gen-test-scripts 的指令不变 | Rejected: 不完整 |
| **全量 surface-key 对齐** | — | 一致性彻底，所有文件同步 | 改动面大（~30 文件） | **Selected: 一次性解决，避免后续反复** |

## Feasibility Assessment

### Technical Feasibility

完全可行。核心代码改动仅 `pipeline.go` 一处（expansion 模式切换），其余均为 skill 文档和模板的文本替换。`per-surface-key` expansion 机制已存在并被 run-tests 使用。

### Resource & Timeline

改动量大但机械性高，预计 2-3 小时完成。

### Dependency Readiness

`per-surface-key` expansion 函数、`isSingleSurface` 判断、`expandPerSurfaceKey` 逻辑均已存在于 `pipeline.go` 中。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| surface-type 是测试分区的正确维度 | Assumption Flip | Overturned: 用户按功能分区（backend/frontend），不按技术类型（api/web）。type 可重复，key 天然唯一 |
| 单 surface named 项目也需要 surfaceKey 目录层 | Stress Test | Overturned: 单 surface 无论 scalar 还是 named，都应保持 `tests/<journey>/` 的简洁结构。只有多 surface 才需要 key 层级 |
| init-justfile 的 recipe 路径不受影响 | Codebase Check | Confirmed: justfile recipe 通过 `<prefix>test <journey>` 参数传递 journey 名，recipe 内部的 `tests/` 路径需要适配新目录 |

## Scope

### In Scope

1. **pipeline.go**：gen-test-scripts expansion 从 `per-surface-type` 改为 `per-surface-key`
2. **gen-test-scripts SKILL.md**：输出目录从 `tests/<journey>/` 改为自适应规则
3. **gen-test-scripts/types/**（5 文件）：更新输出目录描述
4. **gen-test-scripts/rules/**（step-0.5-validation, step-1-contract-loading, quality-gates, run-to-learn）：更新目录引用
5. **run-tests SKILL.md**：验证兼容性，确认 journey discovery 和 recipe 调用路径正确
6. **init-justfile SKILL.md**：确认 recipe 前缀逻辑与新目录对齐
7. **init-justfile/templates/**（6 文件）：justfile 模板中 `tests/` 路径适配多 surface
8. **init-justfile/rules/surfaces/**（5 文件）：确认 recipe 定义与新目录兼容
9. **gen-journeys SKILL.md + templates/journey.md**：Journey frontmatter 的 `surface_types` 是否需要同步 surface-key 信息
10. **gen-contracts SKILL.md + rules/journey-contract-model.md**：更新 `tests/<journey>/` 引用
11. **test-guide SKILL.md + templates/surfaces/**（5 文件）+ rules/（3 文件）+ references/：全量更新测试目录约定
12. **hooks/guide.md**：第 77 行测试目录约定声明
13. **submit-task/data/record-format-test.md**：示例路径更新
14. **forge-cli 测试文件**：引用 `gen-scripts-{type}` 的测试用例更新为 `{key}` 命名

### Out of Scope

- Convention 文件目录结构（`docs/conventions/testing/<type>/core.md`）不变 — convention 按 type 组织是正确的
- Surface rule 文件命名（`rules/surfaces/<type>.md`）不变 — rule 按 type 分类
- Test tag 命名（`@api-functional`）不变 — tag 按 type 标注
- `eval` skill 的 rubric 路径不变 — eval 引用的是 `docs/features/<slug>/testing/` 下的文件，不受影响
- `breakdown-tasks` 和 `quick-tasks` 的 Surface-Key/Type Inference 逻辑不变 — 已正确处理

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 单 surface 项目行为退化 | L | H | `expandPerSurfaceKey` 的 `isSingleSurface` 分支已处理单 surface 去后缀逻辑 |
| init-justfile 模板路径错误导致 run-tests 失败 | M | H | 逐个模板验证 journey filter 路径在单/多 surface 下的正确性 |
| 遗漏某个引用 `tests/<journey>/` 的文件导致新的不一致 | M | M | 全量 grep `tests/<journey>` 和 `tests/{{journey}}` 确认覆盖完整 |
| 已有多 surface 项目的 justfile 需要手动更新 | H | L | 变更仅影响 Forge plugin 分发的模板和 skill 指令，用户项目已有的 justfile 不受影响（recipe 不变，变的是测试文件位置） |

## Success Criteria

- [ ] `pipeline.go` 中 gen-test-scripts 使用 `per-surface-key` expansion，通过 `go test ./pkg/task/...` 全部测试
- [ ] 多 surface 项目（如 pm-work-tracker）中，gen-test-scripts 任务文件命名为 `gen-test-scripts-backend.md` 而非 `gen-test-scripts-api.md`
- [ ] 多 surface 项目中，生成的测试文件输出到 `tests/<surfaceKey>/<journey>/` 而非 `tests/<journey>/`
- [ ] 单 surface 项目行为不变：任务名无后缀，输出到 `tests/<journey>/`
- [ ] 所有 skill 文件中引用 `tests/<journey>/` 的描述已统一为自适应规则，grep 确认无遗漏

## Next Steps

- Proceed to `/write-prd` to formalize requirements
