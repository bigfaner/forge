---
created: "2026-05-28"
author: fanhuifeng
status: Draft
---

# Proposal: Regression Fix Task 按 Test File 拆分（语言无关）

## Problem

Quality-gate hook 检测到 regression test 失败后创建单个 fix task 覆盖所有失败，当失败数量多（20+）且跨多个不相关文件时，agent 执行时卡住。

### Evidence

实际发生：quality-gate 创建的 fix-1 包含 4 个 test suite 共 20+ 失败，agent 执行后长时间无响应被用户手动中断。Concise error 只展示输出尾部，agent 需要读 raw-output.txt 才能看到完整失败列表。详见 `docs/lessons/gotcha-fix-task-broad-scope.md`。

### Urgency

高。每次 regression 测试出现多文件失败都会触发此问题，agent 卡死浪费时间和 token。

## Proposed Solution

在 `quality_gate.go` 中新建 `addRegressionFixTasks` 函数，新建 `extractFileLineMap` 函数从 output 中提取文件路径到相关输出行的映射（含上下文窗口，区别于现有 `extractSourceFiles` 的扁平字符串输出），通过命名约定识别测试文件，按测试文件分组创建独立 fix task。每个 task 只包含该测试文件相关的输出行。

当无法识别测试文件时，fallback 到现有按目录分组的 `addFixTask` 行为。

`addRegressionFixTasks` 使用 regression 专用软上限（10 个），不受 `maxFixTasksPerStep` 硬上限约束——拆分后每个 fix task scope 已收窄到单文件，regression 路径以独立软上限替代通用 cap。`addFixTask`（compile/fmt/lint/unit-test 步骤）保留 cap 作为 loop-breaker 安全阀。

### Innovation Highlights

无创新。原提案使用 Go 专属的 `FAIL <package>` 解析，改进为新建 `extractFileLineMap` 保留文件-输出行映射（现有 `extractSourceFiles` 返回扁平字符串且无行位置信息），结合测试文件命名约定识别。与 CI 系统中按 failing module 分组报告的常规做法一致。

## Requirements Analysis

### Key Scenarios

- **Happy path**：4 个测试文件各有 2-5 个失败 → 创建 4 个 fix task，每个包含对应文件的输出行
- **单文件失败**：1 个测试文件有 3 个失败 → 创建 1 个 fix task，行为与当前一致
- **非标准命名**：测试文件名不匹配任何命名约定 → fallback 到按目录分组
- **栈 trace 引用辅助测试文件**：`handler_test.go` 失败但栈 trace 包含 `utils_test.go:42`（辅助函数）→ 两个文件各创建 fix task，`utils_test.go` 的 task 为误报 → 可接受：误报 task 执行后会发现无需修改，成本可控
- **全部通过**：无失败 → 不创建任何 fix task（现有行为不变）

### Non-Functional Requirements

- **性能**：解析 test output 的时间可忽略（现有 `extractSourceFiles` 已优化）
- **大输出处理**：`extractFileLineMap` 处理 10000 行 output 的内存占用 < 50MB（纯字符串操作，无复杂对象）
- **兼容性**：不影响 compile/fmt/lint/unit-test 步骤的现有逻辑

### Constraints & Dependencies

- 依赖测试文件命名约定（`*_test.go`、`test_*.py`、`*.test.ts` 等），无法识别的测试文件 fallback 到按目录分组
- 新建 `extractFileLineMap` 函数替代 `extractSourceFiles` 用于 regression 路径（保留文件路径到实际输出行的映射），内部先用 `sourceFileRe` 提取 `file:line` 模式，再叠加框架专属模式覆盖 `sourceFileRe` 无法匹配的格式。所有模式的结果合并去重。按语言的提取模式：
  - Go：`sourceFileRe` 提取 `file_test.go:line`，叠加 `--- FAIL:` 块的缩进行解析（处理多行栈 trace 中文件引用仅出现在缩进行的情况）
  - Python (pytest)：`sourceFileRe` 提取 `file.py:line`，叠加 `FAILED path/to/test_file.py::Class::method` 截取 `::` 前路径
  - JS/TS (Jest)：`sourceFileRe` 提取 `file.test.ts:line`，叠加 `FAIL path/to/test.test.ts` 标题行匹配
  - Java：`sourceFileRe` 提取 `File.java:line`，叠加标准栈 trace `at com.example.TestClass.testMethod(TestClass.java:42)` 解析
  - Ruby：`sourceFileRe` 提取 `file.rb:line`，叠加 Minitest `Failure:\n test_method [path/to/test_file.rb:line]` 和 RSpec `rspec ./spec/test_file.rb:line`
  - **多语言混合输出**：所有模式同时执行，每个提取结果独立映射到对应文件路径，最后按测试文件分组合并。不尝试识别输出使用的语言——让所有模式同时匹配，命中的结果自然归属于对应文件

## Alternatives & Industry Benchmarking

### Industry Solutions

CI 系统普遍按 test module/suite 分组报告失败（GitHub Actions test grouping、JUnit XML testsuite 元素）。错误指纹（error fingerprinting）策略也广泛用于错误聚合：Sentry 按栈 trace 指纹分组错误、JUnit Report XML 的 `testsuite` 元素天然按 suite 隔离失败。本方案借鉴了 fingerprinting 思路——以测试文件路径作为"指纹键"将输出行分配到独立 bucket。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | agent 卡死问题持续 | Rejected: 用户体验差 |
| LLM 分析失败输出自动分组 | AI 辅助 | 可识别语义关联的失败（同一根因） | 非确定性、增加 token 开销、增加外部依赖 | Rejected: 引入非确定性与额外开销 |
| JUnit XML 解析 + suite 级拆分 | JUnit/Ant | 结构化输入、suite 边界清晰 | 需要测试框架输出 JUnit XML、增加格式转换依赖 | Rejected: 强制约束测试框架输出格式 |
| Sentry 式栈 trace 指纹分组 | Sentry | 按根因聚合而非按文件 | 需要栈 trace 解析器、指纹算法复杂 | Rejected: 过度工程化，当前场景不需要去重 |
| 按目录分组（现有行为） | 当前实现 | 语言无关 | 同一目录多文件仍 scope 过大 | Rejected: 未改进 |
| **按测试文件分组** | 通用分组策略 | 语言无关、scope 最窄、新代码量可控 | 依赖命名约定识别；同一根因 bug 创建多任务 | **Selected: 最小改动，最大通用性** |

## Feasibility Assessment

### Technical Feasibility

完全可行。新建 `extractFileLineMap` 函数（返回文件路径到实际输出行的映射，每条记录为该文件匹配的原始输出行——函数内部完成匹配行提取、上下文扩展和重叠去重，返回值已是可直接写入 task description 的内容），新增 `isTestFile` 命名约定匹配函数。`extractFileLineMap` 以 `sourceFileRe` 为基础叠加框架专属模式，输入输出类型简单（string → map[string][]string），测试用例可从现有 regression 输出样本构造。

### Resource & Timeline

预计 1-2 天实现 + 测试。改动集中在 `quality_gate.go`，新增 `extractFileLineMap`（含 5 种语言的提取模式）、`isTestFile`、`createFixTask` helper（从 `addSingleFixTask` 提取）和 `addRegressionFixTasks`（含软上限逻辑）。

### Dependency Readiness

无外部依赖。`extractFileLineMap` 以现有 `sourceFileRe` 正则为基线，叠加框架专属模式，依赖项成熟。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 需要按 Go suite 解析 FAIL 行 | XY Detection: 真正的需求是缩小 fix task scope，不是解析 suite | Overturned: 新建 extractFileLineMap + 测试文件命名约定即可，无需框架专属解析 |
| maxFixTasksPerStep cap 对 regression 必要 | 5 Whys: cap 的原始作用是防止 fix task 循环创建导致失控（loop-breaker），而非仅限制 scope | Overridden for regression path: `addRegressionFixTasks` 使用独立软上限（10 个），不受 `maxFixTasksPerStep` 硬上限约束（按测试文件拆分后 scope 已收窄，降低循环触发概率），但 `addFixTask`（compile/fmt/lint/unit-test 路径）保留 cap 作为 loop-breaker |
| 需要接口抽象层支持多框架 | Occam's Razor: 只有一个真实场景 | Overturned: 函数封装 + 命名约定已足够，无需接口 |

## Scope

### In Scope

- 新建 `isTestFile` 函数，识别 Go/Python/JS-TS/Java/Ruby 的测试文件命名约定（功能显式限定为这 5 类语言）
- 新建 `addRegressionFixTasks` 函数，按测试文件分组创建 fix task
- 新建 `extractFileLineMap` 函数，签名：`func extractFileLineMap(output string) map[string][]string`——输入原始 output 字符串，返回文件路径到提取行的映射（每条记录为该文件相关的原始输出行含上下文窗口，由函数内部完成匹配、上下文扩展和重叠去重）
- 每个 fix task 只包含该测试文件相关的输出行，description 格式为：测试文件路径 + 筛选后的相关输出行（含上下文窗口），关联算法如下：
  1. 遍历 output 每行，检查该行是否包含当前测试文件路径
  2. 匹配行取 **前后各 2 行** 作为上下文窗口（共 5 行）
  3. 同一测试文件的多处匹配，合并重叠的上下文窗口（避免重复行）
  4. 一行匹配多个测试文件时，该行及其上下文归入所有匹配的测试文件
  5. 示例 description 内容：
     ```
     handler_test.go
     --- FAIL: TestGetUser (0.00s)
         handler_test.go:42: expected status 200, got 500
         handler_test.go:43: response body: {"error": "unauthorized"}
     ```
- 保留 `maxFixTasksPerStep` 用于非 regression 步骤（compile/fmt/lint/unit-test 的 `addFixTask` 调用路径），将 `addSingleFixTask` 的 task 创建逻辑（surface inference、template defaults、opts 构造、task 创建、markdown 生成、state 更新）提取为共享 helper（如 `createFixTask`），`addRegressionFixTasks` 和 `addSingleFixTask` 均调用此 helper，仅 `addSingleFixTask` 执行 cap 检查
- `runTestRegression` 失败时调用新函数替代 `addFixTask`
- 无法识别测试文件时 fallback 到现有按目录分组的 `addFixTask`
- 更新或新增单元测试覆盖新逻辑

### Out of Scope

- 基线对比过滤预存在失败（lesson 第二层改进）：本提案仅实现第一层（按测试文件拆分），基线过滤作为后续迭代独立实现。lesson 文档指出两层可独立生效，第一层已能解决当前 agent 卡死问题
- unit-test / compile / lint 步骤的改动
- Surface inference 改进
- `coding.fix.md` 模板变更

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 测试文件命名不规范导致识别失败 | M | L | fallback 到按目录分组，零功能损失 |
| 同一根因 bug 导致多个测试文件失败时创建冲突修复任务 | M | M | 接受此 trade-off：按测试文件分组是独立修复的最小 scope 单元。主动缓解：在 task description 中列出其他可能编辑同一生产代码的相关任务 ID（RELATED_TASKS 字段），供 agent 避免并发冲突 |
| `addRegressionFixTasks` 超出软上限后大量 fix task 并发 | M | M | 每个 fix task scope 已收窄到单文件，并发执行风险可控；引入软上限：regression 路径最多创建 10 个 fix task，超出部分按目录合并为综合 task（此时 scope 大但数量可控，不会同时面临 scope 大 + 数量多两个问题） |
| 输出行关联测试文件的准确性（过度包含导致相邻测试文件输出互相污染） | M | L | 上下文窗口固定为前后各 2 行，重叠窗口合并去重；一行匹配多个测试文件时归入所有匹配文件（宁可多包含） |
| Rust 等不在显式支持列表中的语言（top-10 in sourceExts） | M | L | 显式限定支持范围为 Go/Python/JS-TS/Java/Ruby，其他语言 fallback 到按目录分组（接受此局限，后续迭代按需扩展 isTestFile 规则） |
| 回滚复杂度 | L | M | 回滚方案：将 `runTestRegression` 中的 `addRegressionFixTasks` 调用替换回 `addFixTask`，删除新增函数即可恢复原有行为 |

## Success Criteria

- [ ] 4 个测试文件各有失败时，创建 4 个独立 fix task（而非 1 个综合 task）
- [ ] 每个 fix task 的 description 包含该测试文件相关的输出行（包含该文件路径的行及上下文）
- [ ] `maxFixTasksPerStep` 保留用于非 regression 步骤的 `addFixTask` 调用路径，`addRegressionFixTasks` 受 regression 专用软上限（10 个）约束，不受 `maxFixTasksPerStep` 硬上限限制
- [ ] regression 路径最多创建 10 个 fix task，超出部分 fallback 到按目录分组
- [ ] 无法识别测试文件时 fallback 创建按目录分组的 fix task，行为与改动前一致
- [ ] 现有 compile/fmt/lint/unit-test 步骤的 fix task 创建不受影响
- [ ] 支持至少 5 种语言的测试文件命名约定（Go、Python、JS/TS、Java、Ruby）

## Next Steps

- Proceed to `/write-prd` to formalize requirements
