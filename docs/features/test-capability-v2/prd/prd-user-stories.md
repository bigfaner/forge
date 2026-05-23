---
feature: "test-capability-v2"
---

# User Stories: test-capability-v2

## Story 1: 管线统一 — 使用唯一测试生成路径

**As a** Forge 用户（项目开发者）
**I want to** 只有一条清晰的测试生成路径（Journey-Contract），不需要在 gen-test-cases 和 Journey-Contract 之间选择
**So that** 我能专注于编写 PRD 而不是理解管线内部差异，降低使用门槛

**Acceptance Criteria:**
- Given 一个已有 PRD 的功能
- When 我运行 `/gen-journeys` → `/gen-contracts` → `/gen-test-scripts` 完整流程
- Then 生成的测试代码可执行且覆盖 PRD 的关键验收标准，且没有任何 gen-test-cases 相关的技能、命令、或 rubric 文件残留

---

## Story 2: 深度测试 — 风险驱动的边界/异常覆盖

**As a** Forge 用户（项目开发者）
**I want to** 管线根据功能的风险等级自动生成不同密度的测试矩阵（高风险功能自动产出更多边界/异常测试）
**So that** 安全相关功能获得充分的测试覆盖，而低风险功能不浪费测试资源

**Acceptance Criteria:**
- Given 一个标记为 `risk_level: high` 的 Journey
- When gen-contracts 生成合约规范
- Then 每个 Step 包含 3-5 个 Outcome（含必须衍生的边界 Outcome），高风险 Journey 的测试总数比低风险 Journey 多 ≥ 50%
- And API 场景的每个认证端点自动衍生 `unauthorized` Outcome；TUI 场景的每个异步 Cmd 自动衍生 `timeout` Outcome；CLI 场景衍生 `not-found` + `already-exists` Outcome；WebUI 场景衍生 `validation-error` + `session-expired` Outcome

---

## Story 3: 场景差异化 — 按场景类型定制测试策略

**As a** Forge 用户（项目开发者）
**I want to** 管线根据项目的场景类型（CLI/TUI/WebUI/Mobile/API）自动采用不同的测试策略和层级组合
**So that** CLI 项目生成以 subprocess 断言为主的测试，WebUI 项目生成均衡的浏览器自动化测试，Mobile 项目生成 Maestro YAML 骨架，每种场景都能获得最适合的测试形态

**Acceptance Criteria:**
- Given 一个 CLI 类型的项目
- When 管线生成测试
- Then Contract 测试占比 ≥ 80%，使用 subprocess 执行模型，自动检查二进制文件和环境变量隔离
- Given 一个 Mobile 类型的项目
- When 管线生成测试
- Then 输出 Maestro YAML 骨架（app lifecycle + navigation）和 deep link 测试，复杂场景标记 `manual-only`

---

## Story 4: 评测门禁 — Journey 和 Contract 质量自动验证

**As a** Forge 维护者
**I want to** Journey 和 Contract 文档在生成后自动通过评测技能验证质量，未达阈值自动迭代修正
**So that** 下游的 gen-test-scripts 收到的输入质量有保障，减少因上游文档质量问题导致的测试生成失败

**Acceptance Criteria:**
- Given gen-journeys 生成的 Journey 文档
- When eval-journey 自动评分
- Then 评分 ≥ 850/1000 则通过；否则自动迭代修正，最多 3 轮
- And 3 轮后仍未达阈值，暂停管线并输出当前评分 + 未通过项明细，由用户决定后续操作
- Given gen-contracts 生成的 Contract 文档
- When eval-contract 自动评分
- Then 同样的门禁逻辑适用（≥ 850/1000 通过，最多 3 轮迭代）

---

## Story 5: 通用性 — 新项目快速接入测试

**As a** Forge 用户（项目开发者）
**I want to** 在一个全新项目（无 Convention 文件）中首次运行测试生成时，管线自动检测测试框架并生成 Convention 草稿
**So that** 我不需要手动编写 Convention 文件，只需审核微调即可开始使用测试管线

**Acceptance Criteria:**
- Given 一个包含 pytest 测试文件但无 Convention 文件的 Python 项目
- When 运行 test-guide
- Then 自动检测到 pytest 框架，生成 Convention 草稿，草稿包含全部 4 个必需 section（framework、discovery、structure、assertions），通过 Convention schema 验证
- And 用户审核后修改量 ≤ 草稿总内容的 20% 即可作为正式 Convention 使用
- Given 内置 Convention 库
- Then 包含 pytest、JUnit、Rust/cargo test 共 ≥ 3 个新增 Convention 文件

---

## Story 6: 信息增强 — Run-to-Learn 迭代提升测试精度

**As a** Forge 用户（项目开发者）
**I want to** 通过 Run-to-Learn 机制自动捕获被测系统的实际输出，用运行时信息丰富 Fact Table
**So that** 静态侦察无法获取的信息（动态输出、i18n、渲染结果）通过实际运行补充，生成的测试断言更精确

**Acceptance Criteria:**
- Given 一个初始 Fact Table 覆盖率 < 60% 的项目
- When 启用 Run-to-Learn 机制
- Then 经过 ≤ 3 轮迭代后，Fact Table 覆盖率提升 ≥ 20 个百分点
- And 重新生成的测试中边界/异常 Outcome 占比 ≥ 30%（与不使用 Run-to-Learn 的基线对比）
- And 每轮迭代的骨架测试执行有超时保护，不会无限挂起
