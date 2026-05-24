---
created: 2026-05-24
author: "faner"
status: Draft
---

# Proposal: Pipeline Integration Stitch — 修复 review-doc-pipeline 与 auto-gen-journeys-contracts 的集成缝隙

## Problem

两个已完成的提案（review-doc-pipeline、auto-gen-journeys-contracts）各自独立实现后，在交叉点留下 **14 处集成缝隙**，其中 2 处 P0 级（pipeline 必定失败）、4 处 P1 级（特定场景失败）、8 处 P2 级（维护风险）。

### Evidence

#### P0 — 必定失败

1. **`test.gen-journeys`/`test.gen-contracts` 缺少 prompt 模板映射**（`prompt.go` typeToTemplate）
   - autogen.go 有模板用于 .md 文件生成，但 prompt.go 无运行时执行映射
   - `forge prompt get-by-task-id` 对这两个类型返回 `"unknown type"` 错误
   - 影响：所有含 coding 任务的 feature 执行 test pipeline 时必定失败

2. **`eval.journey`/`eval.contract` 缺少 prompt 模板映射**
   - 同上，breakdown 模式的 eval 质量门控阶段必定失败

#### P1 — 特定场景失败

3. **`eval.*` 类型落入 CategoryCoding**（`category.go`）
   - `eval.` 前缀不匹配任何规则，落入 default → CategoryCoding
   - 影响：eval 任务提交时被要求测试证据（testsPassed/coverage），但 eval 不运行测试

4. **`record-format-doc.md` 仍列 `doc.eval`**（`submit-task/data/`）
   - agent 执行 T-review-doc 时可能使用错误的类型名称

5. **`record-format-test.md` 列出已移除的类型**（`submit-task/data/`）
   - 列出 `test.gen-cases`/`test.eval-cases`/`test.gen-and-run`，缺少 `test.gen-journeys`/`test.gen-contracts`

6. **mixed feature 依赖注入顺序脆弱**（`build.go`）
   - ResolveFirstTestDep 设置依赖后，才 prepend T-review-doc，re-index 幂等性风险

#### P2 — 维护风险

7. **`test.gen-and-run` 废弃类型仍完整存在**（types.go、infer.go、prompt.go、autogen.go、模板文件）
8. **`prompt.go` genScriptBases 包含死代码 `T-quick-gen-and-run`**
9. **`validate_index.go` 仍检查 `T-quick-gen-and-run-` 前缀**
10. **`isTestTaskID` 与 `IsAutoGenTaskID` 覆盖范围不一致**
11. **`category_test.go` 缺少新类型测试用例**
12. **测试文件大量引用废弃的 `T-quick-gen-and-run`**（8+ 文件）
13. **历史 feature 文件使用 `doc.eval` 类型**（10+ 已完成 feature 目录）
14. **README.md / ARCHITECTURE.md 引用过时的 `T-eval-doc`**

### Root Cause

P0 的根本原因是 `prompt.go` 和 `autogen.go` 使用**手写 map** 维护类型→模板映射。新增类型时，开发者更新了 `types.go` + `autogen.go`（生成阶段）但遗漏了 `prompt.go`（执行阶段）。**两个文件的命名约定完全一致**（`.` → `-`），完全可以用自动发现替代手写 map。

### Urgency

P0 意味着当前 pipeline 无法实际执行 gen-journeys、gen-contracts、eval-journey、eval-contract 任务。任何使用 Forge 执行 test pipeline 的 feature 都会失败。

## Proposed Solution

三管齐下：

1. **根治 P0**：将 `prompt.go` 和 `autogen.go` 的手写 map 替换为**基于命名约定的自动发现**，消除"忘记更新 map"的整个 bug 类别
2. **修复 P1**：纠正类型分类、更新过期引用、加固依赖注入
3. **清理 P2**：彻底移除 `test.gen-and-run` 废弃代码、更新文档和测试

### Innovation Highlights

**自动发现模板映射**是关键创新。当前手写 map 的维护成本和遗漏风险随类型数量线性增长。自动发现将映射逻辑从 O(N) 条目降为 O(1) 约定：

```
约定：类型名 "." → 文件名 "-" + ".md"
示例：test.gen-journeys → data/test-gen-journeys.md
```

Scope 中已包含将 `clean-code.md` 重命名为 `code-quality-simplify.md`，重命名后所有类型均遵循约定，零 override。

## Requirements Analysis

### Key Scenarios

- **新类型零配置**: 在 `types.go` 添加新常量 + 在 `data/` 放入模板文件，prompt.go 自动识别
- **Mixed feature pipeline**: doc + coding 任务共存时，T-review-doc 正确插入为 test pipeline 前置依赖
- **Re-index 幂等性**: 对同一 index.json 多次执行 ResolveFirstTestDep 不丢失 T-review-doc 依赖
- **Type system 一致性**: 所有 auto-gen 任务类型在 category.go（含新增 CategoryEval）、isTestTaskID、IsAutoGenTaskID 中覆盖一致；eval 任务归入 CategoryEval 而非 CategoryTest，保持语义正确性
- **Clean removal**: test.gen-and-run 从所有代码路径中彻底消失

### Non-Functional Requirements

- **向后兼容**: 已有 index.json 引用 `test.gen-and-run` 时给出明确错误而非静默失败
- **部署时安全**: 自动发现应在 CLI 入口（`main()` 函数启动路径）验证所有已注册类型的模板文件存在和映射唯一性，缺失模板或碰撞导致启动失败。选择 CLI 入口而非 `init()` 的理由：(1) `init()` 在 `go test` 时也执行，会在单元测试中引入不必要的 embed.FS 依赖和启动延迟；(2) CLI 入口仅在用户实际使用 forge 时触发，测试可以独立运行；(3) CI 环境中 `go test` 不依赖模板文件存在，仅 `forge` 二进制运行时才需要
- **提交验证语义正确性**: eval 任务的 submit-task 验证字段必须与 eval 语义匹配（review 结果，非测试证据）
- **CategoryEval 测试覆盖率**: CategoryEval 的验证逻辑必须有独立的单元测试覆盖，包括：正向用例（接受 review 字段 summary/findings/severity）、负向用例（拒绝纯测试字段 testsPassed/coverage）、边界用例（混合提交 review + test 字段时的行为）

### Constraints & Dependencies

- 两个上游提案的实现代码必须已完成（已满足）
- 修改涉及 forge-cli Go 源码、plugin skill 数据文件、文档

## Alternatives & Industry Benchmarking

### Industry Patterns

此问题是 **Convention over Configuration**（CoC）的经典应用场景。该模式最早由 Rails 广泛推广（路由、数据库表名均由命名约定自动推导），后被 Spring Boot（自动配置）、ASP.NET Core（基于约定的控制器发现）等框架普遍采用。其核心主张：当命名遵循稳定约定时，显式注册表是冗余的。

在任务注册场景中，GitHub Actions 的 reusable workflow 发现机制、Airflow 的 DAG 自动发现（扫描指定目录的 Python 文件）都遵循相同原则——用约定替代手写注册表。

### Auto-Discovery 的等价故障模式

自动发现将故障模式从"忘记添加 map 条目"替换为"忘记创建模板文件"。两者在语义上等价（遗漏类型→模板映射），但存在关键差异：
- **map 遗漏**：编译通过，运行时静默返回错误类型
- **文件遗漏**：编译通过，运行时返回明确 "template not found" 错误

为弥补运行时才发现的不足，方案加入 **init-time 验证**：在程序启动时遍历所有已注册类型常量，逐一检查对应模板文件存在于 embed.FS。任何缺失立即报错退出，将故障左移到部署时而非请求时。init-time 校验还包含**映射唯一性检查**：遍历所有类型常量，对每个类型名执行 `strings.ReplaceAll(typeName, ".", "-") + ".md"`，将结果收集到 map 中，若发现重复文件名则 fatal exit（例如 `test.gen-jour-neys` 和 `test.gen.jour-neys` 均映射到 `test-gen-jour-neys.md`）。此检查在启动时一次性执行，O(N) 复杂度（N = 已注册类型数，当前 < 20）。

**为何选择自动发现而非 "init-time 校验 + 手写 map"**：两者在检测时机上等价（均能在启动时发现遗漏），但安全模型不同。手写 map 方案中，init-time 校验将遗漏从"运行时静默失败"升级为"启动时失败"——遗漏仍会发生，只是更早暴露。自动发现方案中，遗漏的根因被消除：开发者只需在 `types.go` 定义常量 + 在 `data/` 放置模板文件，无需在第三个文件同步注册。自动发现的实际风险是约定碰撞（attack point 7 已分析），这是确定性的、可通过 init-time 唯一性校验一次性解决的，而非每次新增类型都可能触发的人为遗漏。

### Comparison Table

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| 手动补全 map 条目 | 变更最小 | 不治本，下次新增类型还会遗漏；无启动期校验 | Rejected |
| Init-time 校验 + 手写 map | 编译后立即发现遗漏 | map 仍需手动维护；O(N) 条目随类型增长；init-time 校验虽能检测遗漏但无法消除遗漏的根因——开发者仍需在两个文件中同步注册类型，遗漏频率取决于人为纪律而非机制保障；维护负担本身是安全风险：O(N) map 的增长降低开发者添加条目的意愿，init-time 校验仅将遗漏的后果从"运行时静默失败"变为"启动时失败"，但遗漏本身仍会阻断 pipeline | Rejected |
| Code-gen from types.go | 编译时安全 | 引入代码生成管线；增加构建复杂度 | Rejected |
| **自动发现 + init-time 校验 + 全量修复** | 根治 P0；消除手写 map 维护负担；启动时校验提供早期错误检测；一次性清理所有缝隙 | 变更量较大；运行时文件查找（init 阶段一次性开销） | **Selected** |

## Feasibility Assessment

### Technical Feasibility

自动发现实现简单：`strings.ReplaceAll(typeName, ".", "-") + ".md"`，加上 embed.FS 的 ReadFile 验证文件存在。重命名 `clean-code.md` 后无 override。init-time 校验在 CLI 入口（`main()` 函数）中执行：遍历所有类型常量、逐一对 embed.FS 执行 ReadFile 验证文件存在，同时检查映射唯一性（无重复文件名），缺失或碰撞立即 fatal exit。

### Resource & Timeline

预计 2 个 coding task（自动发现重构 + 废弃代码清理）+ 若干 doc task。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 手写 map 是唯一方式 | Occam's Razor | Overturned: 命名约定已完全确定，自动发现更简单 |
| `test.gen-and-run` 需保留以兼容旧 index.json | XY Detection | Overturned: 旧 index.json 应迁移到新 pipeline |
| `eval.*` 归入 CategoryCoding 无实际影响 | Assumption Flip | Overturned: category 影响提交验证和记录渲染。进一步分析发现 CategoryTest 也不适合 eval 任务（要求测试证据而非 review 结果），需新建 CategoryEval |

## Scope

### In Scope

**P0 — 自动发现 + 缺失模板**
- prompt.go: typeToTemplate → 自动发现（纯约定，零 override）
- autogen.go: autogenTypeToFile → 同样自动发现（纯约定，零 override）
- 重命名 `prompt/data/clean-code.md` → `code-quality-simplify.md`（统一命名约定）
- 创建 4 个缺失的 prompt 模板文件：test-gen-journeys.md、test-gen-contracts.md、eval-journey.md、eval-contract.md。这些是**执行阶段模板**（agent 运行时指令），非 autogen 规划阶段模板（生成 .md 文件）。每个模板须包含：任务上下文说明、输入格式描述（从 index.json 的 task configuration 中读取）、期望输出格式、质量门控标准。结构参考现有执行阶段模板（如 clean-code.md / code-quality-simplify.md），而非 autogen.go 的规划阶段模板

**P1 — 类型修复 + 引用更新**
- category.go: 新增 `CategoryEval`（提交验证要求 review 相关字段：summary/findings/severity），将 `eval.` 前缀映射到 CategoryEval 而非 CategoryTest。eval 任务是质量门控（review），不是测试生成器，两者提交时需要的证据字段完全不同。CategoryTest 要求 testsPassed/coverage 等测试证据，CategoryEval 要求 review 结论和发现
- submit-task 验证逻辑: 在 validateSubmission 中为 CategoryEval 添加专用验证分支，接受 review 类字段
- record-format-doc.md: `doc.eval` → `doc.review`
- record-format-test.md: 移除已废弃类型，添加新类型
- build.go: 加固 mixed feature 依赖注入（T-review-doc 与 ResolveFirstTestDep 合并处理）

**P2 — 废弃代码清理**
- 按 gen-and-run Removal Checklist（见 Key Risks 章节）逐一移除 `test.gen-and-run` 在 10 个文件中的引用
- validate_index.go: 移除 `T-quick-gen-and-run-` 前缀检查，替换为迁移感知错误提示（"task type 'test.gen-and-run' is deprecated, regenerate index.json"）
- build.go: 更新 findFirstTestTaskIdx 中 quick-mode fallback，将 T-quick-gen-and-run* 替换为新 test task 前缀匹配
- isTestTaskID: 扩展语义覆盖 T-review-doc，并更新函数文档注释说明包含的任务类型范围
- 补充 category_test.go 测试用例
- 更新 README.md / ARCHITECTURE.md 中 T-eval-doc 引用
- task-lifecycle.md: 更新系统类型列表
- 更新引用废弃类型的测试文件

### Out of Scope

- 重构 resolveBreakdownDeps/resolveQuickDeps 的重复逻辑
- eval rollback 改进
- 旧 index.json 迁移工具
- 历史 feature 文件中的 doc.eval 类型（不影响运行时）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 自动发现遗漏 edge case（如类型名含 `-`） | L | H | 约定已稳定，所有现有类型名仅含 `.` 和字母；重命名 clean-code.md 消除唯一例外；CLI 入口校验确保启动时发现遗漏和映射碰撞 |
| 移除 gen-and-run 后旧 index.json 报错 | L | M | validate_index.go 提供迁移感知错误提示："task type 'test.gen-and-run' is deprecated, regenerate index.json via forge quick-tasks" |
| category 修改影响质量门控现有行为 | M | M | 新增 CategoryEval 专项测试：(1) 正向用例验证 eval 任务提交 review 字段被接受；(2) 负向用例验证 eval 任务仅提交 test 字段被拒绝；(3) 遍历所有 switch-on-category 分支确认 CategoryEval 不落入 default/else 路径；现有 CategoryTest 测试套件作为回归基线 |
| 4 个新 prompt 模板内容不准确 | M | H | 新模板为执行阶段 agent 指令（非 autogen 规划阶段的 .md 生成模板）。每个模板须包含：任务上下文说明、输入格式描述、期望输出格式、质量标准。不可照搬 autogen 模板结构 |
| clean-code.md 重命名导致外部引用断裂 | M | M | 重命名后执行 `grep -r "clean-code" forge-cli/ plugins/` 验证零残留引用；更新 CLAUDE.md 或其他配置中可能的引用 |
| findFirstTestTaskIdx quick-mode 回退失效 | H | M | 移除 T-quick-gen-and-run 后更新 findFirstTestTaskIdx 中的 quick-mode 匹配逻辑，使用新的 test task 前缀（如 T-gen-journeys/T-gen-contracts） |
| gen-and-run 引用移除不完整导致编译失败 | M | H | 逐一枚举需修改的文件清单（见下方 Removal Checklist），按清单逐项验证 |
| eval/CategoryEval 提交验证字段与实际不匹配 | L | M | CategoryEval 验证分支参考现有 eval 任务的 submit-task 实际字段（summary/findings/severity），编写单元测试覆盖 |

### gen-and-run Removal Checklist

需移除 `T-quick-gen-and-run` / `test.gen-and-run` 引用的完整文件清单：
1. `types.go`: 移除 TypeGenAndRun 常量
2. `infer.go`: 移除 InferType 中 gen-and-run 分支（~line 32-33）
3. `prompt.go`: 移除 genScriptBases 中 T-quick-gen-and-run 条目（~line 294）
4. `prompt.go`: 移除 typeToTemplate 中 gen-and-run 映射（如仍有）
5. `autogen.go`: 移除 gen-and-run 相关逻辑（如仍有）
6. `validate_index.go`: 移除 T-quick-gen-and-run- 前缀检查（~line 224-226），替换为迁移错误提示
7. `build.go`: 更新 findFirstTestTaskIdx 中 quick-mode fallback 的匹配模式
8. `category.go`: 确认 gen-and-run 分类规则已无引用
9. 模板文件: 删除 data/ 下 gen-and-run 相关 .md 文件
10. 测试文件: 更新所有引用废弃类型的测试用例（确定性命令：`grep -rl "gen-and-run\|quick-gen-and-run\|T-quick-gen" forge-cli/ plugins/ --include="*_test.go"`，预期命中文件待执行该命令后逐一修改）

## Success Criteria

- [ ] `forge prompt get-by-task-id` 对 `test.gen-journeys`、`test.gen-contracts`、`eval.journey`、`eval.contract` 均返回有效 prompt（P0 修复）
- [ ] prompt.go 和 autogen.go 中不再有手写类型→文件名 map（自动发现生效）
- [ ] `grep -r "doc.eval" forge-cli/ plugins/` 返回零结果
- [ ] `eval.journey` 和 `eval.contract` 归入 CategoryEval（非 CategoryTest）
- [ ] `forge submit-task` 对 eval 任务接受 review 类字段（summary/findings），不要求测试证据（testsPassed/coverage）
- [ ] `forge submit-task` 对 eval 任务**拒绝**仅含测试字段（testsPassed/coverage 无 summary/findings）的提交，返回明确的验证错误
- [ ] `grep -r "gen-and-run" forge-cli/ plugins/` 返回零结果
- [ ] `grep -r "clean-code" forge-cli/ plugins/` 仅匹配重命名后的 code-quality-simplify 文件名，无残留 clean-code.md 引用
- [ ] validate_index.go 对引用 `test.gen-and-run` 的旧 index.json 返回包含 "deprecated" 和 "regenerate" 的迁移指引错误信息
- [ ] mixed feature re-index 幂等：T-review-doc 依赖不丢失
- [ ] findFirstTestTaskIdx 对 quick-mode pipeline 正确返回首个 test task 索引（非 -1）
- [ ] init-time 校验（CLI 入口）：缺失模板文件或映射碰撞时 CLI 启动失败并报告缺失/碰撞的类型名
- [ ] 所有现有测试通过

## Next Steps

- Proceed to `/quick-tasks` to generate tasks
