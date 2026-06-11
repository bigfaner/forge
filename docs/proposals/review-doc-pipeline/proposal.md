---
created: 2026-05-23
author: "faner"
status: Approved
---

# Proposal: Review-Doc Pipeline & Conditional User Stories

## Problem

当前系统中，doc 任务的验证机制存在两个结构性缺陷：

1. **eval-doc 是空壳**：`T-eval-doc` 任务模板仅 10 行，引用不存在的"8维度评分标准"，无 rubric 文件、无 eval 命令、无实际验证逻辑。纯文档特性生成的 eval-doc 任务无法有效核对交付物。
2. **混合特性无 doc 验证**：`build.go` 的路由逻辑是互斥的——有 coding 任务走测试流水线，纯 doc 走 eval-doc。混合特性中 doc 任务执行后无任何质量检查。

### Evidence

- `forge-cli/pkg/task/data/doc-eval.md` 仅含 10 行通用指令，引用不存在的 rubric
- `plugins/forge/skills/eval/rubrics/` 下无 `doc-eval.md`
- `build.go` 中 `needsDocEval()` 仅在所有任务都是 `TypeDoc` 时返回 true
- `needsTestPipeline()` 和 `needsDocEval()` 是互斥分支，不支持同时触发

### Urgency

test-capability-v2 刚合并，编码任务已有完整的 eval-journey → eval-contract → gen-test-scripts 质量链。doc 任务缺乏对等验证机制，文档质量完全依赖执行者自觉。

## Proposed Solution

三部分改动：

1. **eval-doc → review-doc**：重命名并重新定义任务。review-doc 采用 AC 核对清单模型（非 rubric 评分），单次核对 + 直接修复——task-executor 读取每个 doc 任务的交付物，核对 Acceptance Criteria，不符合的直接修改文档。
2. **混合特性组合流水线**：`build.go` 支持同时生成 review-doc 和测试流水线任务。执行顺序：review-doc → gen-journeys → gen-contracts → gen-test-scripts。
3. **write-prd 条件性 user stories**：write-prd 的 Step 7 仅在功能涉及代码时生成 `prd-user-stories.md`，doc-only 功能跳过（user stories 仅服务于 gen-journeys → 测试脚本生成）。

### Innovation Highlights

无特殊创新。核心设计决策是 **AC 核对而非 rubric 评分**——doc 任务通常较小，1000 分评分制过重。review-doc 作为第一道关卡确保文档交付物正确后再进入测试生成，避免测试基于错误文档生成。

## Requirements Analysis

### Key Scenarios

- **纯文档特性**：所有任务为 `doc` 类型 → 生成 T-review-doc，核对每个 doc 任务的 AC 并直接修复
- **纯代码特性**：所有任务为 `coding.*` → 走测试流水线，不生成 review-doc
- **混合特性**：同时有 `coding.*` 和 `doc` 任务 → 生成 review-doc + 测试流水线任务，review-doc 先执行
- **write-prd 代码功能**：Step 7 正常生成 user stories → 下游 gen-journeys 使用
- **write-prd 文档功能**：Step 7 跳过 user stories → 无 gen-journeys 输入

### Non-Functional Requirements

- review-doc 任务执行时间应 < 5 分钟（轻量核对，非完整评审）
- 重命名 `doc.eval` → `doc.review` 需向后兼容已存在的 index.json

### Constraints & Dependencies

- 依赖现有 task-executor agent 执行 review-doc 任务
- 依赖 `forge task index` 的 auto-generation 逻辑
- write-prd 条件判断需要可靠识别功能是否涉及代码

## Alternatives & Industry Benchmarking

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | doc 任务无验证，混合特性无 doc 质量检查 | Rejected: test pipeline v2 已为编码建立质量链，doc 缺口明显 |
| 完整 eval-doc skill + rubric | eval-journey/eval-contract 模式 | 与现有 eval 系统一致 | 1000 分评分对 doc 任务过重，需新建 skill/command/rubric | Rejected: 过度工程 |
| **AC 核对 + 直接修复** | 代码审查中的 checklist 模式 | 轻量、直接、改动小 | 无量化评分，难以跨特性比较 doc 质量 | **Selected: doc 任务核心目标是"按要求修改"，AC 核对是精确匹配** |

## Feasibility Assessment

### Technical Feasibility

完全可行。改动涉及：
- Go CLI：`types.go` 类型常量重命名、`autogen.go` 任务定义、`build.go` 路由逻辑
- Skill 文件：`write-prd/SKILL.md` 条件判断、`quick-tasks/SKILL.md` 类型表更新
- 任务模板：`doc-eval.md` → `doc-review.md` 重写

### Resource & Timeline

- 预估 3-5 个任务，1-2 小时/任务
- 全部为 coding 任务（Go CLI 修改）+ 少量 doc 任务（skill 文件修改）

### Dependency Readiness

- task-executor agent 无需修改（它读取任务模板指令执行）
- `forge task index` 的 auto-generation 已有完整框架

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| eval-doc 需要 rubric 评分 | Assumption Flip | Overturned: doc 任务的核心目标是"按要求修改"，AC 核对更精确。rubric 评分适用于开放式质量评估（如 journey 的"语义纯度"），不适用于对照清单的符合性检查 |
| write-prd 应始终生成 user stories | 5 Whys | Refined: user stories 唯一消费者是 gen-journeys，后者为测试脚本生成服务。doc-only 功能无代码可测试，user stories 无下游消费者 |
| 混合特性只能走一条流水线 | XY Detection | Overturned: 互斥路由是历史实现细节，不是设计约束。review-doc 和 test pipeline 可独立运行 |

## Scope

### In Scope

- 重命名 `doc.eval` → `doc.review`：types.go、autogen.go、build.go、category.go
- 重写任务模板：`doc-eval.md` → `doc-review.md`，含 AC 核对指令 + 直接修复流程
- 修改 `build.go` 路由：支持混合特性同时生成 review-doc 和测试流水线任务
- 修改 `build.go` 依赖：review-doc 在测试流水线之前执行
- 修改 `write-prd/SKILL.md`：Step 7 条件性生成 user stories
- 更新 `quick-tasks/SKILL.md` 类型表：反映 review-doc 类型
- 更新 `breakdown-tasks/SKILL.md` 类型表：同步变更

### Out of Scope

- 新建 `/review-doc` 命令或 skill
- 创建 rubric 文件
- 修改 gen-journeys、gen-contracts、gen-test-scripts
- 修改 eval skill 或其 scorer composition
- 已存在 index.json 的迁移策略（现有特性完成或手动更新）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| `doc.eval` → `doc.review` 重命名导致已存在的 index.json 任务类型不匹配 | M | L | 直接替换，不保留旧值。现有特性若包含 `doc.eval` 任务需手动更新 index.json |
| review-doc 任务执行者（LLM）过度修改文档，引入新错误 | M | M | 任务模板明确限制：仅修复 AC 不符合项，不做风格/结构优化 |
| write-prd 条件判断误判（将混合功能误判为 doc-only） | L | M | 判断依据为 In Scope 中是否包含任何可编译/可运行文件路径，与 type assignment 规则一致 |

## Success Criteria

- [ ] `forge task index` 对纯 doc 特性生成 T-review-doc 任务，任务模板包含 AC 核对 + 直接修复指令
- [ ] `forge task index` 对混合特性同时生成 T-review-doc 和测试流水线任务，review-doc 依赖先于 gen-journeys
- [ ] `forge task index` 对纯代码特性不生成 review-doc
- [ ] write-prd 对 doc-only 功能跳过 Step 7 (user stories)，对含代码功能正常生成
- [ ] `doc.eval` 类型值从 ValidTypes 和 SystemTypes 中移除，完全替换为 `doc.review`
- [ ] review-doc 任务执行能核对 AC 并直接修复不符合项

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
