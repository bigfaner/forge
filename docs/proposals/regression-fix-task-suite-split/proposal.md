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

**前置验证**（已完成）：经核实 `quality_gate.go:614-628`，当前 `countFixTasks` 按 title prefix `"fix " + step + ":"` 匹配并排除终态任务（Completed/Rejected/Skipped），**不依赖** `SourceTaskID` 字段过滤。`gotcha-quality-gate-fix-task-loop.md` 中记录的 `SourceTaskID` bug（`opts.SourceTaskID` 未赋值）不影响当前 cap 机制——cap 通过 title prefix 匹配正常工作。`addRegressionFixTasks` 的 regression 专用软上限可直接复用 `countFixTasks` 的 title prefix 匹配逻辑，无需额外修复。

两阶段无条件交付。Phase 0 独立交付且具有自身价值（改进 description 信息完整性），不作为 Phase 1 的"验证门"——Phase 1 无论 Phase 0 效果如何都会实施。理由：Phase 0 解决信息可见性问题（agent 能看到完整失败列表），Phase 1 解决 scope 过大问题（单 task 覆盖所有失败）。两者解决不同层面的问题，不互斥。

**Phase 0**：改进 `addSingleFixTask` 的 description 生成逻辑，从取 output tail 改为提取所有 `--- FAIL` 行列表。半天的低风险改动，使 agent 无需读 raw-output.txt。

**Phase 1**：在 `quality_gate.go` 中新建 `addRegressionFixTasks` 函数，新建 `extractFileLineMap` 函数从 output 中提取文件路径到相关输出行的映射（含上下文窗口，区别于现有 `extractSourceFiles` 的扁平字符串输出），通过命名约定识别测试文件，按测试文件分组创建独立 fix task。每个 task 只包含该测试文件相关的输出行。仅为包含直接 `--- FAIL:` 条目的测试文件创建独立 task；仅作为栈 trace 引用出现的文件（如 `utils_test.go:42`）不生成独立 task，其输出行归入拥有 `--- FAIL:` 条目的主文件。非测试文件路径（如 `handler.go:42`）的输出行归入 fallback task。

当无法识别测试文件时（`isTestFile` 返回零匹配），fallback 到现有按目录分组的 `addFixTask` 行为，产生与改动前完全一致的单个 monolithic task——零改进、零回归。fallback 触发时输出结构化日志警告（`WARNING: isTestFile returned zero matches, falling back to directory-grouped fix task`），便于事后排查识别覆盖率。

`addRegressionFixTasks` 使用 regression 专用软上限（10 个），不受 `maxFixTasksPerStep` 硬上限约束——拆分后每个 fix task scope 已收窄到单文件，regression 路径以独立软上限替代通用 cap。`addFixTask`（compile/fmt/lint/unit-test 步骤）保留 cap 作为 loop-breaker 安全阀。超出软上限时的合并策略：第 11 至 N 个文件的所有输出行合并为 1 个综合 overflow task（title 为 `"fix test: regression overflow (N-10 files)"`），不使用 `groupFilesByDir`（该函数对单目录返回 nil），直接将剩余文件的行拼接为一个 description。总 task 数 ≤ 11（10 个独立 task + 1 个 overflow task）。

### Innovation Highlights

无创新。原提案使用 Go 专属的 `FAIL <package>` 解析，改进为新建 `extractFileLineMap` 保留文件-输出行映射（现有 `extractSourceFiles` 返回扁平字符串且无行位置信息），结合测试文件命名约定识别。与 CI 系统中按 failing module 分组报告的常规做法一致。MVP 收窄到 Go 单语言，降低首次实现复杂度。

## Requirements Analysis

### Key Scenarios

- **Happy path**：4 个测试文件各有 2-5 个失败 → 创建 4 个 fix task，每个包含对应文件的输出行
- **单文件失败**：1 个测试文件有 3 个失败 → 创建 1 个 fix task，行为与当前一致
- **非标准命名**：测试文件名不匹配任何命名约定 → fallback 到按目录分组
- **栈 trace 引用辅助测试文件**：`handler_test.go` 失败但栈 trace 包含 `utils_test.go:42`（辅助函数）→ `handler_test.go` 拥有直接 `--- FAIL:` 条目，创建 fix task；`utils_test.go` 仅作为栈 trace 引用，不创建独立 task，其相关输出行归入 `handler_test.go` 的 task。避免 `utils_test.go` 被多个失败文件同时引用时的输出污染
- **全部通过**：无失败 → 不创建任何 fix task（现有行为不变）

### Non-Functional Requirements

- **性能**：`extractFileLineMap` 对 output 做两次线性扫描（第一次收集 `--- FAIL:` 主文件集合，第二次按行匹配并扩展上下文窗口），时间复杂度 O(L × F)（L = output 行数，F = 主文件数量，受 soft cap 限制 F ≤ 10），对 10000 行 output 的典型场景，解析时间 < 100ms
- **大输出处理**：`extractFileLineMap` 处理 10000 行 output 的原始字符串约 1MB（每行 ~100 bytes）；`map[string][]string` 因步 5 跨文件复制产生的扩展倍率 ≤ F ≤ 10，典型场景总内存 < 15MB（原始 1MB × 10× 扩展 + map 开销）
- **并发执行预算**：regression 路径最多创建 10 个独立 fix task + 1 个综合 overflow task。建议（非约束）dispatcher 层面限制同一 feature 下同时被 claim 的 fix task 不超过 3 个，降低并发编辑冲突
- **兼容性**：不影响 compile/fmt/lint/unit-test 步骤的现有逻辑

### Constraints & Dependencies

- 依赖测试文件命名约定（MVP 仅 `*_test.go`），无法识别的测试文件 fallback 到按目录分组
- 新建 `extractFileLineMap` 函数替代 `extractSourceFiles` 用于 regression 路径（保留文件路径到实际输出行的映射）。MVP 仅实现 Go 的提取模式：
  - 先扫描 output 收集所有包含直接 `--- FAIL:` 条目的测试文件作为"主文件"集合
  - `sourceFileRe` 提取 `file_test.go:line`，仅对"主文件"集合中的文件生成映射条目
  - 叠加 `--- FAIL:` 块的缩进行解析（处理多行栈 trace 中文件引用仅出现在缩进行的情况）
  - 非测试文件路径的输出行（如 `handler.go:42`）和未归属行归入 fallback task
- 后续迭代按需扩展其他语言（Python/JS-TS/Java/Ruby）的提取模式，扩展时所有模式同时执行，结果合并去重，不尝试识别输出使用的语言

## Alternatives & Industry Benchmarking

### Industry Solutions

CI 系统普遍按 test module/suite 分组报告失败（GitHub Actions test grouping、JUnit XML testsuite 元素）。错误指纹（error fingerprinting）策略也广泛用于错误聚合：Sentry 按栈 trace 指纹分组错误、JUnit Report XML 的 `testsuite` 元素天然按 suite 隔离失败。本方案借鉴了 fingerprinting 思路——以测试文件路径作为"指纹键"将输出行分配到独立 bucket。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | agent 卡死问题持续 | Rejected: 用户体验差 |
| **改进 description 信息呈现**（不拆分 task，只把完整 `--- FAIL` 行列表放进 task description 替代 tail） | lesson 直接原因分析 | 改动极小（~20 行），直接解决"agent 看不到完整失败列表"的问题 | scope 仍为单 task 覆盖所有失败，未拆分 | **Selected (Phase 0): 无条件交付，解决信息可见性问题** |
| 按测试文件分组（本提案） | 通用分组策略 | scope 最窄、新代码量可控 | 依赖命名约定识别；同一根因 bug 创建多任务；加剧 claim priority 和 cross-feature pollution | **Selected (Phase 1): 无条件交付，解决 scope 过大问题** |
| LLM 分析失败输出自动分组 | AI 辅助 | 可识别语义关联的失败（同一根因） | 非确定性：相同输出可能产生不同分组；每次分组需 ~2K token；引入 LLM 调用延迟（~5s）和外部依赖 | Rejected: 非确定性 + 延迟 + token 开销无法接受 |
| JUnit XML 解析 + suite 级拆分 | JUnit/Ant | 结构化输入、suite 边界清晰、零歧义解析 | 需 `go test -json` 输出格式支持；forge 当前 `just test` 不产出 JUnit XML；需改造测试命令和 CI 流水线 | Rejected: 改造 `just test` 命令的 scope 超出本提案范围 |
| Sentry 式栈 trace 指纹分组 | Sentry | 按根因聚合而非按文件；去重能力 | 需栈 trace 规范化解析器（Go 栈 trace 格式有 3+ 变体）；指纹碰撞需人工调参；实现复杂度约 500 行 vs 按文件分组约 150 行 | Rejected: 复杂度 3x 于所选方案，当前场景不需要根因去重 |
| 按目录分组（现有行为） | 当前实现 | 语言无关 | 同一目录多文件仍 scope 过大 | Rejected: 未改进 |

## Feasibility Assessment

### Technical Feasibility

完全可行。新建 `extractFileLineMap` 函数（返回文件路径到实际输出行的映射，每条记录为该文件匹配的原始输出行——函数内部完成匹配行提取、上下文扩展和重叠去重，返回值已是可直接写入 task description 的内容），新增 `isTestFile` 命名约定匹配函数。`extractFileLineMap` 以 `sourceFileRe` 为基础叠加框架专属模式，输入输出类型简单（string → map[string][]string），测试用例可从现有 regression 输出样本构造。

### Resource & Timeline

**Phase 0（改进 description 信息呈现）**：半天。改动集中在 `addSingleFixTask` 的 description 生成逻辑，从取 tail 改为提取所有 `--- FAIL` 行列表。

**Phase 1（按测试文件拆分）**：1-2 天实现 + 测试（实现约 1 天，测试约 1 天）。改动集中在 `quality_gate.go`，新增 `extractFileLineMap`（MVP 仅 Go 模式）、`isTestFile`（仅 Go）、`createFixTask` helper（从 `addSingleFixTask` 提取，含独立单元测试）和 `addRegressionFixTasks`（含软上限逻辑）。`countFixTasks` cap 机制已正常工作，无需额外改动。

### Dependency Readiness

无外部依赖。`extractFileLineMap` 以现有 `sourceFileRe` 正则为基线，叠加框架专属模式，依赖项成熟。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| 需要按 Go suite 解析 FAIL 行 | XY Detection: 真正的需求是缩小 fix task scope，不是解析 suite | Overturned: 新建 extractFileLineMap + 测试文件命名约定即可，无需框架专属解析 |
| maxFixTasksPerStep cap 对 regression 必要 | 5 Whys: cap 的原始作用是防止 fix task 循环创建导致失控（loop-breaker），而非仅限制 scope | Overridden for regression path: `addRegressionFixTasks` 使用独立软上限（10 个），不受 `maxFixTasksPerStep` 硬上限约束（按测试文件拆分后 scope 已收窄，降低循环触发概率），但 `addFixTask`（compile/fmt/lint/unit-test 路径）保留 cap 作为 loop-breaker |
| 需要接口抽象层支持多框架 | Occam's Razor: 只有一个真实场景 | Overturned: 函数封装 + 命名约定已足够，无需接口 |
| 需要同时支持 5 种语言 | YAGNI: 当前 forge 项目自身是 Go，非 Go 语言的测试失败是假设场景 | Overturned: MVP 仅实现 Go，后续按需扩展。多语言支持的设计保留在提案中作为扩展路径 |
| agent 卡死一定因为 scope 过大 | First Principles: lesson 中直接原因是"concise error 只展示输出尾部"，agent 看不到完整失败列表 | Resolved: Phase 0 和 Phase 1 解决不同层面的问题（信息可见性 vs scope），两者无条件交付 |

## Scope

### In Scope

- 确认 `countFixTasks` cap 机制正常工作：当前按 title prefix `"fix " + step + ":"` 匹配（不依赖 `SourceTaskID` 字段），cap 已正常工作（对应 `gotcha-quality-gate-fix-task-loop.md`）。本提案直接复用此逻辑
- 新建 `isTestFile` 函数，MVP 仅支持 Go 测试文件命名约定（`*_test.go`），后续迭代按需扩展 Python/JS-TS/Java/Ruby。Go 是 forge 项目自身的语言，先解决眼前实际问题并验证效果
- 新建 `addRegressionFixTasks` 函数，按测试文件分组创建 fix task。仅为包含直接 `--- FAIL:` 条目的测试文件创建独立 task；仅作为栈 trace 引用的文件不生成独立 task，其输出行归入拥有 `--- FAIL:` 条目的主文件。非测试文件路径（如 `handler.go:42`）的输出行归入 fallback task。每个拆分 task 的 title 包含文件名以确保唯一性和可导航性，格式为 `"fix test: <filename> failure in quality gate"`（如 `"fix test: handler_test.go failure in quality gate"`），使 `countFixTasks` 的 title prefix 匹配仍能正确计数
- 新建 `extractFileLineMap` 函数，签名：`func extractFileLineMap(output string) map[string][]string`——输入原始 output 字符串，返回文件路径到提取行的映射（每条记录为该文件相关的原始输出行含上下文窗口，由函数内部完成匹配、上下文扩展和重叠去重）。仅为拥有直接 `--- FAIL:` 条目的文件生成映射条目。MVP 仅实现 Go 的提取模式（`sourceFileRe` 提取 `file_test.go:line`，叠加 `--- FAIL:` 块解析），后续迭代叠加其他语言模式
- Phase 0：改进 `addSingleFixTask` 的 description 生成逻辑，从取 output tail 改为提取所有 `--- FAIL` 行列表，使 agent 能看到完整失败信息而不必读 raw-output.txt
- 每个 fix task 只包含该测试文件相关的输出行，description 格式为：测试文件路径 + 筛选后的相关输出行（含上下文窗口），关联算法如下：
  1. 遍历 output，收集所有包含直接 `--- FAIL:` 条目的测试文件作为"主文件"集合
  2. 遍历 output 每行，检查该行是否包含"主文件"集合中的测试文件路径
  3. 匹配行取 **前后各 2 行** 作为上下文窗口（共 5 行）
  4. 同一测试文件的多处匹配，合并重叠的上下文窗口（避免重复行）
  5. 一行匹配多个"主文件"时，该行及其上下文归入所有匹配的主文件
  6. 仅作为栈 trace 引用出现的测试文件（无直接 `--- FAIL:` 条目）不生成独立映射，其输出行通过步骤 2 自然归入引用它的主文件上下文
  7. 示例 description 内容：
     ```
     handler_test.go
     --- FAIL: TestGetUser (0.00s)
         handler_test.go:42: expected status 200, got 500
         handler_test.go:43: response body: {"error": "unauthorized"}
     ```
- 保留 `maxFixTasksPerStep` 用于非 regression 步骤（compile/fmt/lint/unit-test 的 `addFixTask` 调用路径），将 `addSingleFixTask` 的 task 创建逻辑（surface inference、template defaults、opts 构造、task 创建、markdown 生成、state 更新）提取为共享 helper（如 `createFixTask`），`addRegressionFixTasks` 和 `addSingleFixTask` 均调用此 helper，仅 `addSingleFixTask` 执行 cap 检查。共享 helper 须有独立单元测试覆盖，确保两条调用路径行为一致（task 字段填充、markdown 生成、state 更新）
- `runTestRegression` 失败时调用新函数替代 `addFixTask`
- 无法识别测试文件时 fallback 到现有按目录分组的 `addFixTask`，并输出结构化日志警告（`WARNING: isTestFile returned zero matches for output, falling back to directory-grouped fix task`），便于事后排查
- 更新或新增单元测试覆盖新逻辑

### Out of Scope

- 基线对比过滤预存在失败（lesson 第二层改进）：本提案仅实现第一层（按测试文件拆分），基线过滤作为后续迭代独立实现。lesson 文档指出两层可独立生效，第一层已能解决当前 agent 卡死问题
- fix task claim priority 修复（`gotcha-fix-task-claim-priority.md`）：短期沿袭现有 dispatcher 行为，长期方案（fix task ID 改为数字前缀或 dispatcher 优先 claim fix task）作为独立提案
- cross-feature pollution 修复（`gotcha-quality-gate-cross-feature-pollution.md`）：短期沿袭"mark as skipped"工作流，长期方案（quality-gate 按 feature task 类型过滤测试范围）作为独立提案
- dispatcher 层并发限制（NFR 中的"建议"不约束本提案实现范围）
- unit-test / compile / lint 步骤的改动
- Surface inference 改进
- `coding.fix.md` 模板变更
- 非 Go 语言的测试文件识别和输出解析（Python/JS-TS/Java/Ruby 作为后续迭代）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 测试文件命名不规范导致识别失败 | M | L | fallback 到按目录分组，产生与改动前完全一致的 monolithic task——零回归、零改进。触发时输出结构化日志警告便于排查 |
| 同一根因 bug 导致多个测试文件失败时创建冲突修复任务 | M | M | 接受此 trade-off：按测试文件分组是独立修复的最小 scope 单元。建议 dispatcher 限制同 feature 同时 claim 的 fix task ≤ 3 个，降低并发冲突概率 |
| `addRegressionFixTasks` 超出软上限后大量 fix task 并发 | M | M | 软上限 10 个 + 1 个综合 overflow task（直接拼接剩余文件输出行，不依赖 `groupFilesByDir`）。不会同时面临 scope 大 + 数量多两个问题 |
| 输出行关联测试文件的准确性（栈 trace 引用文件误创建独立 task） | L | M | 仅为有直接 `--- FAIL:` 条目的文件创建独立 task，栈 trace 引用文件不生成 task |
| `createFixTask` helper 重构引入回归影响现有 fix task 路径 | M | M | 独立单元测试覆盖两条调用路径行为一致性；重构先于新功能合并，确保现有 compile/fmt/lint/unit-test 路径不受影响 |
| 拆分后 fix task 数量增加，加剧 claim priority 问题（`gotcha-fix-task-claim-priority.md`：字母 ID `fix-1` 排序输给数字 ID `4`） | H | M | 短期：dispatcher 层面对同一 feature 下 fix task 添加优先 claim 逻辑（修复 task 优先于业务 task）。长期：fix task ID 改为数字前缀（如 `0-fix-1`） |
| 拆分后 cross-feature pollution 加剧（`gotcha-quality-gate-cross-feature-pollution.md`：project-wide 测试的失败挂在当前 feature 下） | H | M | 拆分前已存在，拆分后一个 docs-only feature 可能被塞入更多 fix task。短期：沿袭现有"mark as skipped"工作流。长期：quality-gate 按 feature task 类型过滤测试范围 |
| Go `--- FAIL:` 块解析覆盖不全（sub-test、parallel test 输出格式差异） | M | M | MVP 明确限定支持的 Go test 输出格式：标准 `--- FAIL: TestName` + 缩进栈 trace；sub-test（`--- FAIL: TestFoo/SubTest`）的嵌套缩进按父 test 归属处理；parallel test 交错输出不在 MVP 范围内，归入 fallback |
| Rust 等不在显式支持列表中的语言（top-10 in sourceExts） | M | L | MVP 仅支持 Go，其他语言 fallback 到按目录分组（接受此局限，后续迭代按需扩展 isTestFile 规则） |
| 回滚复杂度 | L | L | 回滚方案：将 `runTestRegression` 中的 `addRegressionFixTasks` 调用替换回 `addFixTask`，删除新增函数即可恢复原有行为 |

## Success Criteria

- [ ] 确认 `countFixTasks` 按 title prefix 匹配正常工作（已验证，无需额外改动）
- [ ] **Phase 0**：`addSingleFixTask` 的 description 包含所有 `--- FAIL` 行条目（从 output 中提取完整列表，而非 tail 截断），agent 无需读取 `raw-output.txt` 即可看到全部失败信息
- [ ] N（N ≤ 10）个测试文件各有直接 `--- FAIL:` 条目时，创建 N 个独立 fix task（而非 1 个综合 task）
- [ ] 每个 fix task 的 description 包含该测试文件相关的输出行（包含该文件路径的行及 ±2 行上下文），仅作为栈 trace 引用的文件不生成独立 task
- [ ] `maxFixTasksPerStep` 保留用于非 regression 步骤的 `addFixTask` 调用路径，`addRegressionFixTasks` 受 regression 专用软上限（10 个）约束，不受 `maxFixTasksPerStep` 硬上限限制
- [ ] 超过 10 个测试文件有失败时，第 11 至 N 个文件的所有输出行合并为 1 个综合 overflow task，总 task 数 ≤ 11
- [ ] 无法识别测试文件时 fallback 创建按目录分组的 fix task（与改动前行为一致），并输出结构化日志警告
- [ ] 现有 compile/fmt/lint/unit-test 步骤的 fix task 创建不受影响
- [ ] MVP 支持 Go 测试文件命名约定（`*_test.go`）
- [ ] **端到端验证**：给定包含 4+ 个测试文件各有失败的 regression output，拆分后的每个 fix task 被 agent 单独 claim 后在 10 分钟内产出有效修复或正确判定无需修改（不出现"长时间无响应"）

## Next Steps

- Proceed to `/write-prd` to formalize requirements
