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

**前置修复**：修复 `gotcha-quality-gate-fix-task-loop.md` 中记录的 `SourceTaskID` 赋值 bug。当前 `addSingleFixTask` 设置了 `Vars["SOURCE_TASK_ID"]`（模板渲染用）但从未设置 `opts.SourceTaskID`（结构体字段），导致 `countFixTasks` 的过滤条件 `SourceTaskID != ""` 永远为 false，cap 形同虚设。不修复此 bug，regression 软上限同样无法生效。

在 `quality_gate.go` 中新建 `addRegressionFixTasks` 函数，新建 `extractFileLineMap` 函数从 output 中提取文件路径到相关输出行的映射（含上下文窗口，区别于现有 `extractSourceFiles` 的扁平字符串输出），通过命名约定识别测试文件，按测试文件分组创建独立 fix task。每个 task 只包含该测试文件相关的输出行。

当无法识别测试文件时，fallback 到现有按目录分组的 `addFixTask` 行为。

`addRegressionFixTasks` 使用 regression 专用软上限（10 个），不受 `maxFixTasksPerStep` 硬上限约束——拆分后每个 fix task scope 已收窄到单文件，regression 路径以独立软上限替代通用 cap。`addFixTask`（compile/fmt/lint/unit-test 步骤）保留 cap 作为 loop-breaker 安全阀。

### Innovation Highlights

无创新。原提案使用 Go 专属的 `FAIL <package>` 解析，改进为新建 `extractFileLineMap` 保留文件-输出行映射（现有 `extractSourceFiles` 返回扁平字符串且无行位置信息），结合测试文件命名约定识别。与 CI 系统中按 failing module 分组报告的常规做法一致。MVP 收窄到 Go 单语言，降低首次实现复杂度。

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

- 依赖测试文件命名约定（MVP 仅 `*_test.go`），无法识别的测试文件 fallback 到按目录分组
- 新建 `extractFileLineMap` 函数替代 `extractSourceFiles` 用于 regression 路径（保留文件路径到实际输出行的映射）。MVP 仅实现 Go 的提取模式：
  - `sourceFileRe` 提取 `file_test.go:line`
  - 叠加 `--- FAIL:` 块的缩进行解析（处理多行栈 trace 中文件引用仅出现在缩进行的情况）
- 后续迭代按需扩展其他语言（Python/JS-TS/Java/Ruby）的提取模式，扩展时所有模式同时执行，结果合并去重，不尝试识别输出使用的语言

## Alternatives & Industry Benchmarking

### Industry Solutions

CI 系统普遍按 test module/suite 分组报告失败（GitHub Actions test grouping、JUnit XML testsuite 元素）。错误指纹（error fingerprinting）策略也广泛用于错误聚合：Sentry 按栈 trace 指纹分组错误、JUnit Report XML 的 `testsuite` 元素天然按 suite 隔离失败。本方案借鉴了 fingerprinting 思路——以测试文件路径作为"指纹键"将输出行分配到独立 bucket。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | agent 卡死问题持续 | Rejected: 用户体验差 |
| **改进 description 信息呈现**（不拆分 task，只把完整 `--- FAIL` 行列表放进 task description 替代 tail） | lesson 直接原因分析 | 改动极小（~20 行），直接解决"agent 看不到完整失败列表"的问题 | scope 仍为单 task 覆盖所有失败，未拆分 | **Phase 0: 先验证信息不足是否是卡死的真正原因，如解决则无需拆分** |
| 按测试文件分组（本提案） | 通用分组策略 | scope 最窄、新代码量可控 | 依赖命名约定识别；同一根因 bug 创建多任务；加剧 claim priority 和 cross-feature pollution | **Selected (Phase 1): Phase 0 验证后如仍需拆分则实施** |
| LLM 分析失败输出自动分组 | AI 辅助 | 可识别语义关联的失败（同一根因） | 非确定性、增加 token 开销、增加外部依赖 | Rejected: 引入非确定性与额外开销 |
| JUnit XML 解析 + suite 级拆分 | JUnit/Ant | 结构化输入、suite 边界清晰 | 需要测试框架输出 JUnit XML、增加格式转换依赖 | Rejected: 强制约束测试框架输出格式 |
| Sentry 式栈 trace 指纹分组 | Sentry | 按根因聚合而非按文件 | 需要栈 trace 解析器、指纹算法复杂 | Rejected: 过度工程化，当前场景不需要去重 |
| 按目录分组（现有行为） | 当前实现 | 语言无关 | 同一目录多文件仍 scope 过大 | Rejected: 未改进 |

## Feasibility Assessment

### Technical Feasibility

完全可行。新建 `extractFileLineMap` 函数（返回文件路径到实际输出行的映射，每条记录为该文件匹配的原始输出行——函数内部完成匹配行提取、上下文扩展和重叠去重，返回值已是可直接写入 task description 的内容），新增 `isTestFile` 命名约定匹配函数。`extractFileLineMap` 以 `sourceFileRe` 为基础叠加框架专属模式，输入输出类型简单（string → map[string][]string），测试用例可从现有 regression 输出样本构造。

### Resource & Timeline

**Phase 0（改进 description 信息呈现）**：半天。改动集中在 `addSingleFixTask` 的 description 生成逻辑，从取 tail 改为提取所有 `--- FAIL` 行列表。

**Phase 1（按测试文件拆分）**：1-2 天实现 + 测试。改动集中在 `quality_gate.go`，新增 `extractFileLineMap`（MVP 仅 Go 模式）、`isTestFile`（仅 Go）、`createFixTask` helper（从 `addSingleFixTask` 提取）和 `addRegressionFixTasks`（含软上限逻辑），加上 `SourceTaskID` bug 修复。

### Dependency Readiness

无外部依赖。`extractFileLineMap` 以现有 `sourceFileRe` 正则为基线，叠加框架专属模式，依赖项成熟。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 需要按 Go suite 解析 FAIL 行 | XY Detection: 真正的需求是缩小 fix task scope，不是解析 suite | Overturned: 新建 extractFileLineMap + 测试文件命名约定即可，无需框架专属解析 |
| maxFixTasksPerStep cap 对 regression 必要 | 5 Whys: cap 的原始作用是防止 fix task 循环创建导致失控（loop-breaker），而非仅限制 scope | Overridden for regression path: `addRegressionFixTasks` 使用独立软上限（10 个），不受 `maxFixTasksPerStep` 硬上限约束（按测试文件拆分后 scope 已收窄，降低循环触发概率），但 `addFixTask`（compile/fmt/lint/unit-test 路径）保留 cap 作为 loop-breaker |
| 需要接口抽象层支持多框架 | Occam's Razor: 只有一个真实场景 | Overturned: 函数封装 + 命名约定已足够，无需接口 |
| 需要同时支持 5 种语言 | YAGNI: 当前 forge 项目自身是 Go，非 Go 语言的测试失败是假设场景 | Overturned: MVP 仅实现 Go，后续按需扩展。多语言支持的设计保留在提案中作为扩展路径 |
| agent 卡死一定因为 scope 过大 | First Principles: lesson 中直接原因是"concise error 只展示输出尾部"，agent 看不到完整失败列表 | Challenged: 增加 Phase 0（改进 description 信息呈现），先验证信息不足是否是卡死的真正原因 |

## Scope

### In Scope

- 修复 `SourceTaskID` 赋值 bug：`addSingleFixTask` 中将 `opts.SourceTaskID` 设为固定标识（如 `"quality-gate"`），使 `countFixTasks` 的 cap 检查生效（对应 `gotcha-quality-gate-fix-task-loop.md`）
- 新建 `isTestFile` 函数，MVP 仅支持 Go 测试文件命名约定（`*_test.go`），后续迭代按需扩展 Python/JS-TS/Java/Ruby。Go 是 forge 项目自身的语言，先解决眼前实际问题并验证效果
- 新建 `addRegressionFixTasks` 函数，按测试文件分组创建 fix task
- 新建 `extractFileLineMap` 函数，签名：`func extractFileLineMap(output string) map[string][]string`——输入原始 output 字符串，返回文件路径到提取行的映射（每条记录为该文件相关的原始输出行含上下文窗口，由函数内部完成匹配、上下文扩展和重叠去重）。MVP 仅实现 Go 的提取模式（`sourceFileRe` 提取 `file_test.go:line`，叠加 `--- FAIL:` 块解析），后续迭代叠加其他语言模式
- Phase 0：改进 `addSingleFixTask` 的 description 生成逻辑，从取 output tail 改为提取所有 `--- FAIL` 行列表，使 agent 能看到完整失败信息而不必读 raw-output.txt
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
- fix task claim priority 修复（`gotcha-fix-task-claim-priority.md`）：短期沿袭现有 dispatcher 行为，长期方案（fix task ID 改为数字前缀或 dispatcher 优先 claim fix task）作为独立提案
- cross-feature pollution 修复（`gotcha-quality-gate-cross-feature-pollution.md`）：短期沿袭"mark as skipped"工作流，长期方案（quality-gate 按 feature task 类型过滤测试范围）作为独立提案
- unit-test / compile / lint 步骤的改动
- Surface inference 改进
- `coding.fix.md` 模板变更
- 非 Go 语言的测试文件识别和输出解析（Python/JS-TS/Java/Ruby 作为后续迭代）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 测试文件命名不规范导致识别失败 | M | L | fallback 到按目录分组，零功能损失 |
| 同一根因 bug 导致多个测试文件失败时创建冲突修复任务 | M | M | 接受此 trade-off：按测试文件分组是独立修复的最小 scope 单元。主动缓解：在 task description 中列出其他可能编辑同一生产代码的相关任务 ID（RELATED_TASKS 字段），供 agent 避免并发冲突 |
| `addRegressionFixTasks` 超出软上限后大量 fix task 并发 | M | M | 每个 fix task scope 已收窄到单文件，并发执行风险可控；引入软上限：regression 路径最多创建 10 个 fix task，超出部分按目录合并为综合 task（此时 scope 大但数量可控，不会同时面临 scope 大 + 数量多两个问题） |
| 输出行关联测试文件的准确性（过度包含导致相邻测试文件输出互相污染） | M | L | 上下文窗口固定为前后各 2 行，重叠窗口合并去重；一行匹配多个测试文件时归入所有匹配文件（宁可多包含） |
| 拆分后 fix task 数量增加，加剧 claim priority 问题（`gotcha-fix-task-claim-priority.md`：字母 ID `fix-1` 排序输给数字 ID `4`） | H | M | 短期：dispatcher 层面对同一 feature 下 fix task 添加优先 claim 逻辑（修复 task 优先于业务 task）。长期：fix task ID 改为数字前缀（如 `0-fix-1`） |
| 拆分后 cross-feature pollution 加剧（`gotcha-quality-gate-cross-feature-pollution.md`：project-wide 测试的失败挂在当前 feature 下） | H | M | 拆分前已存在，拆分后一个 docs-only feature 可能被塞入更多 fix task。短期：沿袭现有"mark as skipped"工作流。长期：quality-gate 按 feature task 类型过滤测试范围 |
| Rust 等不在显式支持列表中的语言（top-10 in sourceExts） | M | L | MVP 仅支持 Go，其他语言 fallback 到按目录分组（接受此局限，后续迭代按需扩展 isTestFile 规则） |
| 回滚复杂度 | L | M | 回滚方案：将 `runTestRegression` 中的 `addRegressionFixTasks` 调用替换回 `addFixTask`，删除新增函数即可恢复原有行为 |

## Success Criteria

- [ ] 修复 `SourceTaskID` 赋值 bug：`addSingleFixTask` 设置 `opts.SourceTaskID = "quality-gate"`，`countFixTasks` 的 cap 检查生效
- [ ] 4 个测试文件各有失败时，创建 4 个独立 fix task（而非 1 个综合 task）
- [ ] 每个 fix task 的 description 包含该测试文件相关的输出行（包含该文件路径的行及上下文）
- [ ] `maxFixTasksPerStep` 保留用于非 regression 步骤的 `addFixTask` 调用路径，`addRegressionFixTasks` 受 regression 专用软上限（10 个）约束，不受 `maxFixTasksPerStep` 硬上限限制
- [ ] regression 路径最多创建 10 个 fix task，超出部分 fallback 到按目录分组
- [ ] 无法识别测试文件时 fallback 创建按目录分组的 fix task，行为与改动前一致
- [ ] 现有 compile/fmt/lint/unit-test 步骤的 fix task 创建不受影响
- [ ] MVP 支持 Go 测试文件命名约定（`*_test.go`）

## Next Steps

- Proceed to `/write-prd` to formalize requirements
