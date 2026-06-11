---
created: 2026-05-24
author: "faner"
status: Draft
---

# Proposal: Go CLI Record 渲染管线的回归验证（按任务类型）

## Problem

Go CLI 的 record 渲染管线（validateRecordData + markdown 模板渲染）接受 RecordData JSON 输入并输出 records/*.md。但当前缺乏系统性验证手段——无法确认每种 task type 的历史记录能否通过当前 Go CLI 的校验与渲染，输出与已提交 markdown 一致的结果。

### Evidence

- forge-cli-clean-code 完成了 21 条任务记录（coding.cleanup/coding.refactor/coding.fix），但从未系统性验证这些历史记录是否与当前 Go 模板渲染管线兼容
- v3.0.0 重构期间产生了丰富的任务记录，覆盖了 12 种 task type，这些数据可直接用于验证 Go CLI 的渲染一致性
- Go 端模板和 RecordData schema 持续迭代，历史记录格式已与当前管线产生偏差——例如 forge-cli-clean-code 的 `coding.fix` 记录使用旧的 `summary` 字段格式，而当前 Go struct 已重命名为 `resolution`（前向不兼容：历史记录包含当前 struct 已移除的字段），但目前没有机制系统性检测此类偏差
- 新增的 task type（doc.eval、doc.summary、doc.drift、gate）没有经过 golden dataset 回归验证

### Urgency

随着 task type 持续增加（当前 12 种活跃类型），Go CLI 渲染管线对各类 record 的兼容性风险逐步累积。现在用已有的丰富历史数据建立回归测试，成本最低——数据已存在，只需提取为 fixture。

## Proposed Solution

从已完成的 feature 提取历史 record 作为 golden dataset，按 task type 分组建立 Go CLI 层的回归测试：validateRecordData 校验 + markdown 模板渲染对比。测试文件位于 `internal/record/golden_test.go`，fixture 以 `testdata/{taskType}/{featureName}_{index}.json` + `testdata/{taskType}/{featureName}_{index}.golden.md` 的命名约定组织。table-driven 测试使用 `testCase` struct，字段包括 `name string`、`taskType string`、`inputJSON string`（fixture 路径）、`goldenMD string`（golden 路径）。

### Innovation Highlights

Golden dataset 回归测试是标准的软件工程实践，无特别创新。亮点在于按 task type 分组组织 golden dataset，直接验证 Go CLI 渲染管线一致性（RecordData JSON → validateRecordData → markdown 渲染 → golden 对比）。

## Requirements Analysis

### Key Scenarios

- **Type dispatch correctness**: Go CLI 为 `coding.feature`、`doc.review`、`gate` 等不同类型选择正确的 markdown 模板
- **Happy path**: 历史记录通过 validateRecordData 校验，渲染出的 markdown 与已提交记录一致
- **Schema 偏差（前向不兼容）**: 历史 RecordData JSON 包含当前 Go RecordData struct 已删除或重命名的字段，导致校验或渲染失败
- **模板渲染偏差**: Go markdown 模板渲染结果与已提交的 records/*.md 格式不一致
- **缺失字段（后向遗漏）**: 当前 Go 端新增的 required 字段在历史 record 中不存在——golden dataset 无法自动检测此类偏差，需通过 schema 对比人工补充
- **数据质量边界**: 历史记录中字段值极长（如 description 含大段文本）、特殊字符（markdown reserved characters）、空值（optional 字段未填充）——验证模板渲染对这些输入的处理是否稳健
- **Legacy type 兼容**: 裸 `fix` 类型（非 `coding.fix`）的 record 构建是否正确

### Non-Functional Requirements

- 测试需在 CI 中可运行，执行时间 < 30s（预估分解：~30 fixtures × 2 次函数调用/fixture ≈ ~60 次函数调用，无网络 I/O，文件 I/O 极小（仅读取 fixture JSON 和 golden MD），预计 <10s，预留 buffer 至 30s）
- Golden dataset 易于扩展：新增 feature 时可增量添加 fixture
- 按 task type 分组组织 fixture，便于定位类型分发问题

### Constraints & Dependencies

- 依赖 forge-cli 已有的 RecordData struct、validateRecordData、markdown 模板
- 不修改 Go 端校验逻辑（只测试，不改行为）

## Alternatives & Industry Benchmarking

### Industry Solutions

Golden / snapshot testing 是验证数据管道输出一致性的标准方法（Go 的 `testing` 包内置 `-update` flag 支持此模式）。业界成熟项目广泛采用此模式：

- **Hugo** (`hugolib/golden_test.go`): 按 content type 分组的 golden tests，验证 template rendering 输出与 `.golden` 文件一致。其 fixture 组织方式（`testdata/` 按 type 分目录）直接启发本方案的 `testdata/{taskType}/` 命名约定
- **Terraform** (`command/testdata/`): 大量 golden fixture 验证 CLI 子命令输出格式稳定性，每个子命令独立目录、每条 fixture 包含 input + expected output pair
- **Protoc** (`golden_test.go`): Google Protocol Buffers 编译器使用 golden tests 验证多语言代码生成输出，覆盖 schema evolution 场景（旧 `.proto` 文件在新编译器上的输出一致性）

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 类型分发出错无检测手段，发现成本随时间增长 | Rejected: 已有足够历史数据 |
| JSON Schema 验证 | jsonschema | 标准化、工具多、schema drift 检测精确；能覆盖 golden dataset 的盲区（后向遗漏：当前 Go 端新增 required 字段在历史 record 中不存在） | 需维护两套 schema（JSON Schema + Go struct），无法验证渲染输出正确性 | Rejected: 与 golden dataset 互补但本方案优先覆盖渲染管线一致性；schema drift 检测可作后续增强 |
| **Golden dataset 对比（按 task type 分组）** | Go testing `-update`；借鉴 Hugo `hugolib/golden_test.go` 的按 type 分组 fixture 和 Terraform `command/testdata/` 的 input+expected pair 模式 | 直接验证渲染管线输出、数据真实、按类型定位问题 | fixture 需随格式演进更新；无法检测 record-format 模板文档与 Go struct 的偏差（只能测已有记录的回归，不能发现文档描述与实现的不一致）；历史 markdown 与当前渲染不一致时需人为判断以哪方为准 | **Selected: Golden dataset 覆盖 validateRecordData + RenderRecord 的渲染管线一致性，JSON Schema 覆盖 schema drift 检测——两者互补，本方案优先覆盖渲染管线，schema drift 检测可作后续增强** |

## Feasibility Assessment

### Technical Feasibility

完全可行。Go CLI 已有 `validateRecordData()` 和 `RenderRecord()` 两个核心函数，可直接在测试中以 Go 函数调用方式（非 subprocess）调用。测试代码直接 import internal/record 包，构造 RecordData 输入，调用两个函数并对比输出。历史 record 数据已存在于各 feature 的 `docs/features/*/tasks/records/` 目录。

> **Scope note**: 本方案测试 Go CLI 渲染管线（RecordData JSON → 校验 → markdown 渲染 → golden 对比）。submit-task skill 层（模板选择逻辑、字段填充逻辑）属于 LLM agent 行为，不在自动化测试范围内——其正确性通过人工 code review 和 record-format 模板约定保证。

### Resource & Timeline

scope 小：提取 fixture（机械操作）+ 写测试（table-driven by task type）+ 分析修复模板。预计 5-8 个 coding task。

### Dependency Readiness

无外部依赖。所有数据在本地仓库。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 历史记录都是"正确"的 | Stress Test | Confirmed: 记录由 `forge task submit` 生成，通过了 Go 端校验，可作为 golden dataset |
| Go 端 RecordData schema 与历史记录格式兼容 | Assumption Flip | 部分验证：golden dataset 仅覆盖历史数据中出现的字段组合，不保证与所有可能的合法输入兼容 |
| 需要 test/validation category 的 golden data | Occam's Razor | Overturned: test/validation 类型由系统自动生成，不经过标准渲染管线，先跳过 |
| 原提案 4 个 feature 足够覆盖 | Need Gate | Refined: 扩展为 6+ 个 feature，确保每种 task type 至少 2 条代表性记录 |

## Golden Dataset: Feature Sources & Task Type Coverage

### Feature Sources

| Feature | Records | Key Task Types | 日期 |
|---------|---------|---------------|------|
| forge-cli-clean-code | 21 | coding.cleanup (3), coding.refactor (11), coding.fix (6), fix (6)† | 05-24 |
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
† `fix` 是 `coding.fix` 的 legacy alias，共享同一模板。forge-cli-clean-code 中 `coding.fix (6)` 和 `fix (6)` 为同一批记录的两种计数方式（task type 字段值为 `fix` 或 `coding.fix`），不重复计算——Coverage Matrix 中 coding.fix 的 15 条仅统计 task type 为 `coding.fix` 的记录（forge-arch 的 9 条 + forge-cli-clean-code 的 6 条），`fix (6)` 作为 legacy alias 单独列出。
| fix (legacy bare) | 8 | forge-cli-clean-code, task-executor-prompt, forge-arch | Good (alias of coding.fix) |
| doc | 35+ | unify-surfaces, eval-freeform, spec-authority, cli-query, auto-gen-journeys, review-doc | Excellent |
| doc.review | 3 | unify-surfaces, eval-freeform, cli-query | Adequate |
| doc.eval | 5 | eval-freeform-pre-revision | Adequate |
| doc.summary | 6 | test-capability-v2, forge-arch | Adequate |
| doc.drift | 6 | 待确认来源 feature（历史数据散布于多个已完成 feature，提取时按实际记录定位） | Adequate |
| gate | 6 | test-capability-v2, forge-arch | Adequate |

## Scope

### In Scope

- 从 10 个已完成 feature 提取代表性 record 作为 golden dataset fixture（每种 task type 选 2-3 条）
- Go CLI 回归测试：validateRecordData 校验 + markdown 渲染对比，按 task type 分组 table-driven
- 分析历史记录，发现 Go 端 RecordData schema 与历史格式的偏差并更新 golden fixture 以匹配当前 Go 行为（若发现 Go struct 需兼容调整则单独提 issue，不在本 feature scope 内）
- 覆盖 11 种独立 task type + 1 个 alias（`fix` 为 `coding.fix` 的 legacy alias，共享同一模板路径）：coding（5 子类型）、doc（5 子类型）、gate（1 种）、fix（legacy alias of coding.fix）

### Out of Scope

- test.* 和 validation.* category 的 golden dataset（无历史记录）
- Go 端 validateRecordData 校验规则的修改
- 新增 CLI 命令
- LLM agent 输出确定性的测试

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 历史 record 格式与当前 Go 模板不一致（因模板已迭代） | M | M | 测试用 `-update` flag 模式；CI 中 diff gating 通过 Go test 的 `-update` flag 实现：测试失败时输出 diff，开发者传入 `-update` 重新生成 golden 文件，CI 检测 `.golden` 文件 git diff 非空即标记失败 |
| record.json（agent 输入）与 records/*.md（CLI 输出）不完整对应 | M | L | 从已提交的 markdown records 反向验证，不依赖 process/record.json |
| fixture 数量大导致测试维护成本高 | L | L | 每个 task type 选 2-3 条代表性记录，不追求全量 |
| fixture 代表性不足，遗漏渲染边界情况 | M | L | 优先选取不同 feature 来源的记录以增加多样性；发现的边界 case 可增量补充 fixture |
| 历史 record 数据质量问题（JSON 格式不完整、字段缺失等非 schema 偏差） | M | M | 提取 fixture 时逐一验证 JSON 完整性；不完整的 record 标注为 known issue 并跳过，不强行补充 |
| legacy `fix` type 的 record 格式可能与 `coding.fix` 混淆 | L | L | `fix` 是 `coding.fix` 的 alias，共享同一模板路径；fixture 中标注 alias 关系即可 |

## Success Criteria

- [ ] ≥11 种独立 task type + fix alias 的 golden dataset fixture 建立完成
- [ ] Go CLI 回归测试通过 `go test` 运行，覆盖 validateRecordData + markdown 渲染
- [ ] Go 端 RecordData schema 与历史记录的偏差逐一识别并记录（发现的偏差记录在测试输出的诊断报告中），golden fixture 更新至与当前 Go 行为一致
- [ ] 新增 fixture 流程可操作：从已通过校验的历史 record 提取 RecordData JSON + 对应 markdown 为 golden pair，在 table-driven 测试中新增一条用例
- [ ] 每种 task type 的历史 record 经过 Go CLI golden dataset 验证（validateRecordData + markdown 渲染一致性）

## Next Steps

- Proceed to `/quick-tasks` to generate tasks
