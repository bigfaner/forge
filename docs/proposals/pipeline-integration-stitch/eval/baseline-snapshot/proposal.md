---
created: 2026-05-24
author: "faner"
status: Draft
---

# Proposal: Pipeline Integration Stitch — 修复 auto-gen-journeys-contracts 提案遗留的集成缝隙

## Problem

`auto-gen-journeys-contracts` 提案引入了 staged test pipeline（gen-journeys → gen-contracts → gen-scripts → run → verify-regression）和 eval 质量门控（eval-journey、eval-contract），但遗漏了执行阶段的配套工作：**4 个 prompt 模板文件缺失**、**eval 类型分类错误**、**依赖注入逻辑引用已废弃类型**、**gen-and-run 废弃代码残留**。

### Evidence

#### P0 — Pipeline 执行必定失败

1. **`prompt/data/` 缺少 4 个执行阶段模板文件**
   - 缺失：`test-gen-journeys.md`、`test-gen-contracts.md`、`eval-journey.md`、`eval-contract.md`
   - `task/data/` 中已有这 4 个文件（autogen 规划阶段模板），但 `prompt/data/`（执行阶段模板）未创建
   - `Synthesize()` 在渲染这些类型的任务 prompt 时 `ReadFile` 失败
   - 影响：所有使用 Forge 执行 test pipeline 或 eval gate 的 feature 必定失败

#### P1 — 特定场景失败

2. **`eval.*` 类型落入 CategoryCoding**（`category.go:16-33`）
   - `CategoryForType()` 无 `eval.` 前缀分支，`eval.journey`/`eval.contract` 落入 default → `CategoryCoding`
   - 影响：eval 任务被要求提供测试证据（testsPassed/coverage），但 eval 是 review 类任务
   - 影响面：`validateRecordData`（submit.go:296）、`RenderRecord`（record.go:235）、prompt 注入（prompt.go `renderTemplate`）

3. **`findFirstTestTaskIdx` quick-mode 分支匹配废弃类型**（`build.go:492-494`）
   - 查找 `T-quick-gen-and-run*`，但 Quick 模式已不再生成该类型（`GetQuickTestTasks` 使用 staged tasks）
   - 当前靠 `return 0` fallback 意外正确（gen-journeys 恰好排首位），但脆弱

4. **T-review-doc prepend 与 ResolveFirstTestDep 顺序耦合**（`build.go:329-347`）
   - 两步操作顺序硬编码：先 ResolveFirstTestDep（设置基础 deps），再 prepend T-review-doc
   - 虽然当前幂等（ResolveFirstTestDep 总是全新覆写），但重排顺序会导致 T-review-doc 丢失
   - 应合并为单步操作消除耦合

#### P2 — 维护风险

5. **`test.gen-and-run` 废弃代码残留**（生产代码 + 测试 + 活跃文档）
   - `types.go`: `TypeTestGenAndRun` 常量、ValidTypes、isTestTaskID 条目
   - `infer.go:32-33`: gen-and-run 推断分支
   - `prompt.go:297,304`: `genScriptBases` 中 `T-quick-gen-and-run` 条目
   - `validate_index.go:224`: `T-quick-gen-and-run-` 前缀检查
   - `build.go:484,492,494`: findFirstTestTaskIdx quick-mode 注释和匹配
   - `prompt/data/test-gen-and-run.md`、`task/data/test-gen-and-run.md`: 废弃模板文件
   - 14 个测试文件中 ~95 处引用
   - 活跃文档引用：OVERVIEW.md、task-lifecycle.md

6. **record-format 参考文档过期**
   - `record-format-test.md` 列出已废弃类型（`test.gen-cases`/`test.eval-cases`/`test.gen-and-run`），缺少新类型
   - 缺少 `record-format-eval.md`：agent 执行 eval 任务时无 JSON 字段参考

### Urgency

P0 意味着 `forge prompt get-by-task-id` 对 `test.gen-journeys`、`test.gen-contracts`、`eval.journey`、`eval.contract` 返回错误，任何执行 test pipeline 或 eval gate 的 feature 必定失败。

## Proposed Solution

三管齐下：

1. **根治 P0**：创建 4 个执行阶段 prompt 模板文件
2. **修复 P1**：新增 CategoryEval + 完整验证/记录分支 + 加固依赖注入为单步操作
3. **清理 P2**：从生产代码、测试文件和活跃文档中移除 gen-and-run 废弃代码

### Innovation Highlights

Task 1（已完成）引入的**自动发现机制**已消除"忘记更新映射"的根因。本次提案聚焦于补全遗漏的**执行层面配套**（模板文件、类型分类、记录格式），使 staged test pipeline 和 eval gate 从"类型注册完成"推进到"端到端可执行"。

## Requirements Analysis

### Key Scenarios

- **新类型零配置**: 在 `types.go` 添加常量 + 在 `data/` 放入模板文件，auto-discovery 自动识别（已由 Task 1 完成）
- **Eval 提交语义正确性**: `forge submit-task` 对 eval 任务接受 review 字段（summary/findings/severity），拒绝纯测试字段
- **Mixed feature 依赖注入**: T-review-doc 正确插入为 test pipeline 前置依赖，单步操作无顺序耦合
- **Quick-mode findFirstTestTaskIdx**: 正确匹配新 staged tasks 前缀（T-test-gen-journeys）
- **Gen-and-run 完全移除**: 生产代码零残留，测试和活跃文档零残留

### Non-Functional Requirements

- **向后兼容**: 旧 index.json 引用 `test.gen-and-run` 时给出明确迁移错误提示
- **CategoryEval 测试覆盖**: 正向用例（接受 review 字段）、负向用例（拒绝纯测试字段）、边界用例（混合提交）
- **eval record 模板**: 包含 score、findings、severity 等 eval 特有字段

### Constraints & Dependencies

- Task 1（auto-discovery + init-time 校验 + clean-code.md 重命名）已完成，本次所有工作基于该基础

## Alternatives & Industry Benchmarking

### Industry Solutions

这是典型的**补全遗漏的 adapter 层**问题。在新类型系统中，注册（types.go）和发现（auto-discovery）已完成，但执行阶段的模板、分类和记录渲染未同步。类比：Airflow 添加新 DAG 类型后需要配套的 executor plugin 和 UI renderer。

### Comparison Table

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| 手动补全模板 + 最小化修复 | 变更最小 | eval 分类仍错；gen-and-run 僵尸代码持续积累 | Rejected |
| 仅修 P0+P1，保留 gen-and-run | 减少变更量 | 废弃代码干扰开发和测试 | Rejected |
| **完整修复 P0+P1+P2** | 端到端可执行；零僵尸代码；类型系统一致 | 变更量较大（但大部分是机械性清理） | **Selected** |

## Feasibility Assessment

### Technical Feasibility

- 4 个 prompt 模板：参考现有模板（如 `code-quality-simplify.md`）结构编写，每模板 ~30-50 行
- CategoryEval：参考 CategoryTest 的实现模式（分类常量 + CategoryForType 分支 + record 模板 + submit 验证）
- 依赖注入加固：将 ResolveFirstTestDep + T-review-doc prepend 合并为单函数调用
- Gen-and-run 移除：机械性操作，按文件清单逐一清理

### Resource & Timeline

预计 3 个 coding task + 2 个 doc task，总工作量 ~6h。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| re-index 幂等性是实际 bug | Code Audit | Overturned: 当前代码已幂等（ResolveFirstTestDep 总是全新覆写），风险仅为代码耦合 |
| gen-and-run 清理应包含历史 feature/proposal 文档 | Occam's Razor | Refined: 历史文档不影响运行时，仅清理活跃文档 |
| eval 需要 CategoryTest 分类 | Assumption Flip | Overturned: eval 是 review 类任务（非测试生成器），需新建 CategoryEval |

## Scope

### In Scope

**P0 — 执行阶段模板**
- 创建 `prompt/data/test-gen-journeys.md`（agent 执行 test.gen-journeys 任务的指令模板）
- 创建 `prompt/data/test-gen-contracts.md`（agent 执行 test.gen-contracts 任务的指令模板）
- 创建 `prompt/data/eval-journey.md`（agent 执行 eval.journey 任务的指令模板）
- 创建 `prompt/data/eval-contract.md`（agent 执行 eval.contract 任务的指令模板）
- 每个模板包含：任务上下文说明、输入格式（从 index.json task configuration 读取）、期望输出格式、质量标准

**P1 — CategoryEval + 依赖加固**
- `category.go`: 新增 `CategoryEval = "eval"`，`CategoryForType` 添加 `eval.` 前缀分支
- `submit.go`: `validateRecordData` 为 CategoryEval 添加验证分支（接受 review 字段 summary/findings/severity，拒绝纯测试字段）
- `types.go`: RecordData 结构添加 eval 特有字段（evalScore、evalFindings、evalSeverity、evalPassed）
- `record.go`: 新增 `record-eval.md` Go 模板 + `RenderEvalRecord` 函数 + `RenderRecord` switch 添加 CategoryEval case
- `plugins/forge/skills/submit-task/data/record-format-eval.md`: eval 任务 JSON 字段参考文档
- `category_test.go` + `submit_test.go`: CategoryEval 专项测试
- `build.go`: 将 ResolveFirstTestDep + T-review-doc prepend 合并为单步操作
- `build.go`: `findFirstTestTaskIdx` quick-mode 分支更新为匹配 `T-test-gen-journeys` 前缀

**P1 — Record 模板参考文档更新**
- `record-format-test.md`: 移除废弃类型，添加 `test.gen-journeys`/`test.gen-contracts`

**P2 — gen-and-run 废弃代码移除**
- 生产代码（5 文件 ~15 处引用）：
  - `types.go`: 移除 TypeTestGenAndRun 常量、ValidTypes、isTestTaskID 条目
  - `infer.go`: 移除 gen-and-run 推断分支（line 32-33）
  - `prompt.go`: 移除 genScriptBases 中 `T-quick-gen-and-run` 条目（line 304）
  - `validate_index.go`: 移除 `T-quick-gen-and-run-` 前缀检查（line 224），替换为迁移感知错误提示
  - `build.go`: 清理 findFirstTestTaskIdx 中 gen-and-run 相关注释和匹配
- 废弃模板文件：删除 `prompt/data/test-gen-and-run.md`、`task/data/test-gen-and-run.md`
- 测试文件（14 文件 ~95 处引用）：更新所有引用 gen-and-run 的测试用例
- 活跃文档：更新 OVERVIEW.md、task-lifecycle.md 中的 gen-and-run 引用

### Out of Scope

- 历史 feature/proposal 文档中的 gen-and-run 引用（~35 文件 ~130 处，不影响运行时）
- 重构 resolveBreakdownDeps/resolveQuickDeps 的重复逻辑
- eval rollback 改进
- 旧 index.json 自动迁移工具
- record-format-doc.md 中 `doc.eval` → `doc.review`（`doc.eval` 不在运行时使用）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 4 个新 prompt 模板内容不准确 | M | H | 参考现有执行阶段模板结构（code-quality-simplify.md），包含上下文/输入/输出/质量标准四段式 |
| CategoryEval 提交验证字段与实际不匹配 | L | M | 验证分支接受 review 字段（summary/findings/severity），编写单元测试覆盖正向/负向/边界用例 |
| gen-and-run 引用移除不完整导致编译失败 | M | H | 按文件清单逐项清理，每项后执行 `go build` 验证 |
| findFirstTestTaskIdx 修改影响现有 dependency wiring | M | M | 更新后添加集成测试验证 Quick mode 依赖链正确性 |
| eval record 模板字段设计不合理 | L | L | 参考现有 eval 任务的 submit-task 实际字段设计 |

## Success Criteria

- [ ] `prompt/data/` 包含 test-gen-journeys.md、test-gen-contracts.md、eval-journey.md、eval-contract.md
- [ ] `Synthesize()` 对 `test.gen-journeys`、`test.gen-contracts`、`eval.journey`、`eval.contract` 返回有效 prompt（P0 修复）
- [ ] `CategoryForType("eval.journey")` 返回 `CategoryEval`（非 `CategoryCoding`）
- [ ] `forge submit-task` 对 eval 任务接受含 summary/findings 的提交，拒绝仅含 testsPassed/coverage 的提交
- [ ] `RenderRecord` 对 CategoryEval 使用 eval 专用 record 模板
- [ ] `findFirstTestTaskIdx` 对 Quick mode 正确返回 gen-journeys 任务索引
- [ ] ResolveFirstTestDep + T-review-doc prepend 为单步操作，无顺序耦合
- [ ] `grep -r "gen-and-run\|quick-gen-and-run\|T-quick-gen" forge-cli/` 返回零结果
- [ ] `grep -r "gen-and-run\|quick-gen-and-run\|T-quick-gen" plugins/forge/` 返回零结果
- [ ] validate_index.go 对引用 `test.gen-and-run` 的旧 index.json 返回迁移指引错误信息
- [ ] 所有现有测试通过

## Next Steps

- Proceed to `/quick-tasks` to generate tasks
