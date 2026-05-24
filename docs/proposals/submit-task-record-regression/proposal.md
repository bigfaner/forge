---
created: 2026-05-24
author: "faner"
status: Draft
---

# Proposal: Submit-Task Record 回归测试与模板优化

## Problem

submit-task skill 的 record-format 模板缺乏系统性验证手段——无法确认 agent 生成的 record.json 经过 Go CLI 校验和渲染后是否与已提交的 records/*.md 一致。

### Evidence

- forge-cli-clean-code 完成了 21 条任务记录（coding.cleanup/coding.refactor/coding.fix），但从未系统性验证 record-format 模板的示例 JSON 是否与 Go 端 RecordData schema 一致
- 其他 feature（test-capability-v2 22 条、auto-gen-journeys-contracts 8 条、review-doc-pipeline 4 条）积累了大量历史记录，这些数据未被用于回归验证
- record-format 模板是纯文本指导文件，其中的示例 JSON 可能与 Go 端实际接受的 schema 存在偏差，但目前没有机制检测这种偏差

### Urgency

随着 task type 持续增加（当前 22 种类型、6 个 category），record-format 模板与 Go 端 schema 的偏差会逐步累积。现在用已完成的丰富历史数据建立回归测试，成本最低。

## Proposed Solution

从 4 个已完成 feature 提取历史 record 作为 golden dataset，建立 Go CLI 层的回归测试：validateRecordData 校验 + markdown 模板渲染对比。基于测试发现的问题同步完善 record-format 模板。

### Innovation Highlights

Golden dataset 回归测试是标准的软件工程实践，无特别创新。亮点在于数据来源：利用 Forge 自身工作流产生的真实任务记录作为测试 fixture，而非手工构造测试数据。

## Requirements Analysis

### Key Scenarios

- **Happy path**: 历史记录通过 validateRecordData 校验，渲染出的 markdown 与已提交记录一致
- **Schema 偏差**: record-format 模板中的示例字段名或类型与 Go RecordData 不匹配
- **模板渲染偏差**: Go markdown 模板渲染结果与已提交的 records/*.md 格式不一致
- **缺失字段**: record-format 模板未提及 Go 端 required 的字段

### Non-Functional Requirements

- 测试需在 CI 中可运行，执行时间 < 30s
- Golden dataset 易于扩展：新增 feature 时可增量添加 fixture

### Constraints & Dependencies

- 依赖 forge-cli 已有的 RecordData struct、validateRecordData、markdown 模板
- 不修改 Go 端校验逻辑（只测试，不改行为）

## Alternatives & Industry Benchmarking

### Industry Solutions

Golden / snapshot testing 是验证数据管道输出一致性的标准方法（Go 的 `testing` 包内置 `-update` flag 支持此模式）。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 偏差持续累积，发现成本随时间增长 | Rejected: 已有足够历史数据 |
| JSON Schema 验证 | jsonschema | 标准化、工具多 | 需维护两套 schema（JSON Schema + Go struct），偏离 golden dataset 对比的目标 | Rejected: 间接验证，不测渲染 |
| **Golden dataset 对比** | Go testing `-update` | 直接验证端到端输出、数据真实 | fixture 需随格式演进更新 | **Selected: 与目标最匹配** |

## Feasibility Assessment

### Technical Feasibility

完全可行。Go CLI 已有 `validateRecordData()` 和 `RenderRecord()` 两个核心函数，可直接在测试中调用。历史 record 数据已存在于 `docs/features/*/tasks/records/` 和 `process/record.json`。

### Resource & Timeline

scope 小：提取 fixture（机械操作）+ 写测试（table-driven）+ 分析修复模板。预计 5-8 个 coding task。

### Dependency Readiness

无外部依赖。所有数据在本地仓库。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 历史记录都是"正确"的 | Stress Test | Confirmed: 记录由 `forge task submit` 生成，通过了 Go 端校验，可作为 golden dataset |
| record-format 模板示例与 Go schema 一致 | Assumption Flip | 待验证：这是本次 feature 的核心目标 |
| 需要 test/validation category 的 golden data | Occam's Razor | Overturned: test/validation 类型由系统自动生成，不走 submit-task 的 record-format 模板，先跳过 |

## Scope

### In Scope

- 从 4 个 feature 提取已完成任务的 record.json 作为 golden dataset fixture
- Go CLI 回归测试：validateRecordData 校验 + markdown 渲染对比
- 分析历史记录，发现 record-format 模板中的问题并修复
- 覆盖的 category：coding（5 种子类型）、doc（2 种子类型）、gate（1 种）

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
| fixture 数量大导致测试维护成本高 | L | L | 每个 category 选 2-3 条代表性记录，不追求全量 |

## Success Criteria

- [ ] ≥3 个 category（coding/doc/gate）的 golden dataset fixture 建立完成
- [ ] Go CLI 回归测试通过 `go test` 运行，覆盖 validateRecordData + markdown 渲染
- [ ] record-format 模板中发现的问题全部修复，示例 JSON 与 Go RecordData schema 一致
- [ ] 新增 fixture 时只需复制文件 + 加一行测试用例

## Next Steps

- Proceed to `/quick-tasks` to generate tasks
