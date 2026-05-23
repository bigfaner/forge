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
- Then 生成的测试代码可执行且覆盖 PRD 的关键验收标准，且以下 gen-test-cases 相关内容已完全删除：
  - 技能目录：`skills/gen-test-cases/`
  - 评测命令：`commands/eval-test-cases/`
  - Rubric 文件：`skills/eval/rubrics/test-cases.md` 及类型子 rubric（`cli-test-cases.md`、`api-test-cases.md`、`ui-test-cases.md`、`tui-test-cases.md`、`mobile-test-cases.md`）
  - 任务类型：`test.gen-cases`（从 `forge task` 系统类型注册中移除）
  - 任务文件：`tasks/` 目录下所有 `type: test.gen-cases` 的 `.md` 文件
- And `forge task list --type test.gen-cases` 返回空列表
- And 全局搜索 `gen-test-cases` 关键词（除 PRD/历史文档外）无匹配

---

## Story 2: 深度测试 — 风险驱动的边界/异常覆盖

**As a** Forge 用户（项目开发者）
**I want to** 管线根据功能的风险等级自动生成不同密度的测试矩阵（高风险功能自动产出更多边界/异常测试）
**So that** 安全相关功能获得充分的测试覆盖，而低风险功能不浪费测试资源

**Acceptance Criteria:**
- Given 一个标记为 `risk_level: high` 的 Journey
- When gen-contracts 生成合约规范
- Then 每个 Step 包含 3-5 个 Outcome（含必须衍生的边界 Outcome），同一功能的 Journey 比较高风险 vs 低风险 variant 的总 Outcome 数量，高风险 ≥ 低风险 × 1.5；若同一功能仅有高风险 Journey（无低风险 variant 可比较），退化为绝对值验证：高风险 Journey 总 Outcome 数 ≥ 13
- And API 场景的每个认证端点自动衍生 `unauthorized` Outcome；TUI 场景的每个异步 Cmd 自动衍生 `timeout` Outcome；CLI 场景衍生 `not-found` + `already-exists` Outcome；WebUI 场景衍生 `validation-error` + `session-expired` Outcome

---

## Story 3: 场景差异化 — 按场景类型定制测试策略

**As a** Forge 用户（项目开发者）
**I want to** 管线根据项目的场景类型（CLI/TUI/WebUI/Mobile/API）自动采用不同的测试策略和层级组合
**So that** CLI 项目生成以 subprocess 断言为主的测试，WebUI 项目生成均衡的浏览器自动化测试，Mobile 项目生成 Maestro YAML 骨架，每种场景都能获得最适合的测试形态

**Acceptance Criteria:**
- Given 一个 CLI 类型的项目
- When 管线生成测试
- Then Contract 测试占比 ≥ 80%（分母：生成的测试函数总数，按 `Contract 测试函数 / (Contract 测试函数 + Journey 烟测试函数)` 计算），使用 subprocess 执行模型，自动检查二进制文件和环境变量隔离
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
- Then 评分 ≥ 850/1000 且每维度不低于最低阈值（完整性 ≥ 120、语义纯度 ≥ 120、前置条件互斥性 ≥ 90、事实依据 ≥ 90、场景适配 ≥ 90、一致性 ≥ 90）则通过；否则自动迭代修正，最多 3 轮
- And 3 轮后仍未达阈值（总分或任一维度），暂停管线并输出当前评分 + 未通过项明细，由用户决定后续操作
- Given gen-contracts 生成的 Contract 文档
- When eval-contract 自动评分
- Then 同样的门禁逻辑适用（≥ 850/1000 且每维度不低于最低阈值，最多 3 轮迭代）

---

## Story 5: 通用性 — 新项目快速接入测试

**As a** Forge 用户（项目开发者）
**I want to** 在一个全新项目（无 Convention 文件）中首次运行测试生成时，管线自动检测测试框架并生成 Convention 草稿
**So that** 我不需要手动编写 Convention 文件，只需审核微调即可开始使用测试管线

**Acceptance Criteria:**
- Given 一个包含 pytest 测试文件但无 Convention 文件的 Python 项目
- When 运行 test-guide
- Then 自动检测到 pytest 框架，生成 Convention 草稿，草稿包含全部 4 个必需 section（framework、discovery、structure、assertions），通过 Convention schema 验证
- And 草稿通过 Convention schema 验证后，用户审核修改量评估方式为：`diff --stat` 统计用户修改行数占草稿总行数的比例；目标 ≤ 20%（此指标为人工审核参考，标记为 **human-verified**）
- Given 用户拒绝 Convention 草稿
- When 用户指出草稿中的错误部分
- Then test-guide 基于用户反馈重新生成草稿（保留用户认可的部分，仅修正被指出的错误），最多重试 2 次；2 次重试耗尽后输出最终草稿供用户手动编辑，管线暂停等待用户确认
- Given 内置 Convention 库
- When 用户查看 Forge 插件的 `conventions/` 目录
- Then 包含 pytest、JUnit、Rust/cargo test 共 ≥ 3 个新增 Convention 文件

---

## Story 6: 信息增强 — Run-to-Learn 迭代提升测试精度

**As a** Forge 用户（项目开发者）
**I want to** 通过 Run-to-Learn 机制自动捕获被测系统的实际输出，用运行时信息丰富 Fact Table
**So that** 静态侦察无法获取的信息（动态输出、i18n、渲染结果）通过实际运行补充，生成的测试断言更精确

**Acceptance Criteria:**
- Given 一个初始 Fact Table 覆盖率 < 60% 的项目
- When 启用 Run-to-Learn 机制
- Then 经过 ≤ 3 轮迭代后，Fact Table 覆盖率从初始值（gen-contracts 静态侦察结果）提升 ≥ 20 个百分点
- And 重新生成的测试中边界/异常 Outcome 占比 ≥ 30%（基线：初始静态侦察时的占比，标记为 **human-verified**）
- And 每轮迭代的骨架测试执行有超时保护，不会无限挂起
- And 生成的测试附带置信度评级：`HIGH`（运行验证事实 ≥ 80%）、`MEDIUM`（运行验证事实 40-80%）、`LOW`（运行验证事实 < 40%），评级基于 Fact Table 中 confirmed 事实占比计算；LOW 置信度的测试仍会执行但在报告中标记为需人工审核后方可采信结果

---

## Story 7: 可扩展场景类型系统 — 维护者新增场景类型

**As a** Forge 维护者
**I want to** 通过添加一个场景类型配置文件（定义检测规则、测试策略、环境检测项、必须 Outcome），让新场景类型无缝接入测试管线
**So that** 社区可以扩展 Forge 支持更多场景（如 Desktop/Electron、嵌入式、gRPC），而不需要修改管线核心代码

**Acceptance Criteria:**
- Given 一个包含 `detect_rules`（项目文件信号模式）、`strategy`（AI 侧重比例）、`env_check`（环境就绪检测命令）、`required_outcomes`（必须衍生 Outcome 列表）的场景类型配置文件
- When 将配置文件放入 `scenarios/` 目录
- Then 管线自动识别新场景类型，`gen-journeys`/`gen-contracts`/`gen-test-scripts` 按配置的策略和规则执行，无需修改任何管线技能代码
- And 新增场景类型后，已有场景类型（CLI/TUI/WebUI/Mobile/API）的测试生成结果不受影响（回归验证）：可接受的差异范围为——测试函数数量变化 ≤ 5%、无新增编译错误、eval 评分偏差 ≤ 30 分（同一 Journey/Contract 输入下）
- And 新场景类型的 `required_outcomes` 配置自动反映到 eval rubric 的"场景适配"维度评分中：eval 评分时按该场景类型的 `required_outcomes` 列表检查 Outcome 覆盖率，缺少任何一个必须 Outcome 扣 30 分/个；新增场景类型需提交至少 1 个人工标注的 gold standard 文档对用于校准
