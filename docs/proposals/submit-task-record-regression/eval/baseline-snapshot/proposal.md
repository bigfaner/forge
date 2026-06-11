---
created: 2026-05-24
author: "faner"
status: Draft
---

# Proposal: Submit-Task 按任务类型构建 Record 的回归验证

## Problem

submit-task skill 根据任务类型（task type）选择不同的 record-format 模板构建 record.json，经过 Go CLI 校验和渲染后生成 records/*.md。但当前缺乏系统性验证手段——无法确认 submit-task 为每种 task type 选择的模板、填充的字段、渲染的 markdown 是否都正确。

### Evidence

- forge-cli-clean-code 完成了 21 条任务记录（coding.cleanup/coding.refactor/coding.fix），但从未系统性验证 record-format 模板的示例 JSON 是否与 Go 端 RecordData schema 一致
- v3.0.0 重构期间产生了丰富的任务记录，覆盖了 12 种 task type，这些数据可用于验证 submit-task 的类型分发逻辑
- record-format 模板是纯文本指导文件，其中的示例 JSON 可能与 Go 端实际接受的 schema 存在偏差，但目前没有机制检测这种偏差
- 新增的 task type（doc.eval、doc.summary、doc.drift、gate）没有经过 golden dataset 回归验证

### Urgency

随着 task type 持续增加（当前 12 种活跃类型），submit-task 的类型分发逻辑（选择模板 → 填充字段 → 渲染）出错的风险逐步累积。现在用已有的丰富历史数据建立回归测试，成本最低——数据已存在，只需提取为 fixture。

## Proposed Solution

从已完成的 feature 提取历史 record 作为 golden dataset，按 task type 分组建立 Go CLI 层的回归测试：validateRecordData 校验 + markdown 模板渲染对比。重点关注 submit-task 是否为每种 task type 正确选择了模板并填充了正确的字段。

### Innovation Highlights

Golden dataset 回归测试是标准的软件工程实践，无特别创新。亮点在于验证目标：不仅验证 schema 一致性，还验证 submit-task 的类型分发逻辑（task type → template → fields → rendered markdown）的端到端正确性。

## Requirements Analysis

### Key Scenarios

- **Type dispatch correctness**: submit-task 为 `coding.feature`、`doc.review`、`gate` 等不同类型选择正确的 record-format 模板
- **Happy path**: 历史记录通过 validateRecordData 校验，渲染出的 markdown 与已提交记录一致
- **Schema 偏差**: record-format 模板中的示例字段名或类型与 Go RecordData 不匹配
- **模板渲染偏差**: Go markdown 模板渲染结果与已提交的 records/*.md 格式不一致
- **缺失字段**: record-format 模板未提及 Go 端 required 的字段
- **Legacy type 兼容**: 裸 `fix` 类型（非 `coding.fix`）的 record 构建是否正确

### Non-Functional Requirements

- 测试需在 CI 中可运行，执行时间 < 30s
- Golden dataset 易于扩展：新增 feature 时可增量添加 fixture
- 按 task type 分组组织 fixture，便于定位类型分发问题

### Constraints & Dependencies

- 依赖 forge-cli 已有的 RecordData struct、validateRecordData、markdown 模板
- 不修改 Go 端校验逻辑（只测试，不改行为）

## Alternatives & Industry Benchmarking

### Industry Solutions

Golden / snapshot testing 是验证数据管道输出一致性的标准方法（Go 的 `testing` 包内置 `-update` flag 支持此模式）。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 类型分发出错无检测手段，发现成本随时间增长 | Rejected: 已有足够历史数据 |
| JSON Schema 验证 | jsonschema | 标准化、工具多 | 需维护两套 schema（JSON Schema + Go struct），偏离 golden dataset 对比的目标 | Rejected: 间接验证，不测渲染 |
| **Golden dataset 对比（按 task type 分组）** | Go testing `-update` | 直接验证端到端输出、数据真实、按类型定位问题 | fixture 需随格式演进更新 | **Selected: 与目标最匹配** |

## Feasibility Assessment

### Technical Feasibility

完全可行。Go CLI 已有 `validateRecordData()` 和 `RenderRecord()` 两个核心函数，可直接在测试中调用。历史 record 数据已存在于各 feature 的 `docs/features/*/tasks/records/` 目录。

### Resource & Timeline

scope 小：提取 fixture（机械操作）+ 写测试（table-driven by task type）+ 分析修复模板。预计 5-8 个 coding task。

### Dependency Readiness

无外部依赖。所有数据在本地仓库。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 历史记录都是"正确"的 | Stress Test | Confirmed: 记录由 `forge task submit` 生成，通过了 Go 端校验，可作为 golden dataset |
| record-format 模板示例与 Go schema 一致 | Assumption Flip | 待验证：这是本次 feature 的核心目标 |
| 需要 test/validation category 的 golden data | Occam's Razor | Overturned: test/validation 类型由系统自动生成，不走 submit-task 的 record-format 模板，先跳过 |
| 原提案 4 个 feature 足够覆盖 | Need Gate | Refined: 扩展为 6+ 个 feature，确保每种 task type 至少 2 条代表性记录 |

## Golden Dataset: Feature Sources & Task Type Coverage

### Feature Sources

| Feature | Records | Key Task Types | 日期 |
|---------|---------|---------------|------|
| forge-cli-clean-code | 21 | coding.cleanup (3), coding.refactor (11), coding.fix (6), fix (6) | 05-24 |
| unify-surfaces | 9 | coding.refactor (1), coding.feature (3), coding.enhancement (2), doc (2), doc.review (1) | 05-24 |
| eval-freeform-pre-revision | 3 | doc (2), doc.review (1) | 05-24 |
| spec-authority-enforcement | 5 | doc (3), coding.enhancement (1), doc (1) | 05-24 |
| cli-query-knowledge-in-guide | 2 | doc (1), doc.review (1) | 05-24 |
| task-list-slug-worktree | 1 | coding.enhancement (1) | 05-24 |
| auto-gen-journeys-contracts | 8 | coding.feature (4), coding.enhancement (1), doc (3) | 05-24 |
| test-capability-v2 | 22 | coding.feature (5), coding.enhancement (5), coding.refactor (1), coding.cleanup (3), gate (3), doc.summary (3) | 05-23 |
| review-doc-pipeline | 4 | coding.refactor (1), coding.enhancement (1), doc (2) | 05-23 |
| forge-architecture-simplification | 35 | coding.feature (7), coding.refactor (6), coding.enhancement (6), coding.fix (9), coding.cleanup (1), gate (3), doc.summary (3) | 05-22 |

### Task Type Coverage Matrix

| Task Type | Records Available | Feature Sources | Coverage |
|-----------|------------------|-----------------|----------|
| coding.feature | 19 | unify-surfaces, auto-gen-journeys, test-capability-v2, forge-arch-simplification | Excellent |
| coding.enhancement | 16 | unify-surfaces, spec-authority, task-list-slug, review-doc, test-capability-v2, forge-arch | Excellent |
| coding.refactor | 20 | forge-cli-clean-code, unify-surfaces, review-doc, test-capability-v2, forge-arch | Excellent |
| coding.cleanup | 7 | forge-cli-clean-code, test-capability-v2, forge-arch | Good |
| coding.fix | 15 | forge-cli-clean-code, forge-arch | Good |
| fix (legacy bare) | 8 | forge-cli-clean-code, task-executor-prompt, forge-arch | Good |
| doc | 35+ | unify-surfaces, eval-freeform, spec-authority, cli-query, auto-gen-journeys, review-doc | Excellent |
| doc.review | 3 | unify-surfaces, eval-freeform, cli-query | Adequate |
| doc.eval | 5 | eval-freeform-expert, run-tasks-git-status, enforce-forge-task-add, run-tests-decouple | Adequate |
| doc.summary | 6 | test-capability-v2, forge-arch | Adequate |
| doc.drift | 6 | worktree-unpushed, auto-task-main, forge-research, list-tasks, cli-restructure, refactor-impact | Adequate |
| gate | 6 | test-capability-v2, forge-arch | Adequate |

## Scope

### In Scope

- 从 10 个已完成 feature 提取代表性 record 作为 golden dataset fixture（每种 task type 选 2-3 条）
- Go CLI 回归测试：validateRecordData 校验 + markdown 渲染对比，按 task type 分组 table-driven
- 分析历史记录，发现 record-format 模板中的问题并修复
- 覆盖 12 种 task type：coding（5 子类型）、doc（5 子类型）、gate（1 种）、fix（legacy bare type）

### Out of Scope

- test.* 和 validation.* category 的 golden dataset（无历史记录）
- Go 端 validateRecordData 校验规则的修改
- 新增 CLI 命令
- LLM agent 输出确定性的测试

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 历史 record 格式与当前 Go 模板不一致（因模板已迭代） | M | M | 测试用 `-update` flag 模式，确认是有意变更则更新 fixture |
| record.json（agent 输入）与 records/*.md（CLI 输出）不完整对应 | M | L | 从已提交的 markdown records 反向验证，不依赖 process/record.json |
| fixture 数量大导致测试维护成本高 | L | L | 每个 task type 选 2-3 条代表性记录，不追求全量 |
| legacy `fix` type 的 record 格式可能与 `coding.fix` 混淆 | M | L | 明确区分两种 type 的 fixture，验证各自走正确的模板 |

## Success Criteria

- [ ] ≥12 种 task type 的 golden dataset fixture 建立完成
- [ ] Go CLI 回归测试通过 `go test` 运行，覆盖 validateRecordData + markdown 渲染
- [ ] record-format 模板中发现的问题全部修复，示例 JSON 与 Go RecordData schema 一致
- [ ] 新增 fixture 时只需复制文件 + 加一行测试用例
- [ ] submit-task 为每种 task type 选择的模板经过 golden dataset 验证

## Next Steps

- Proceed to `/quick-tasks` to generate tasks
