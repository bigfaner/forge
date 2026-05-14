---
feature: "typed-verification-strategies"
---

# User Stories: typed-verification-strategies

## Story 1: 类型化测试用例自动生成

**As a** Forge 用户（开发者）
**I want to** gen-test-cases 根据我项目的 interface 类型自动生成不同的验证条件（TUI 生成 golden file + 维度检查，API 生成契约测试，CLI 生成输出 golden file）
**So that** 测试能捕获类型特有的 bug（如 TUI 渲染溢出、API 契约偏差、CLI 输出格式错误），而不是只检查功能正确性

**Acceptance Criteria:**
- Given 一个使用 go-test profile（capabilities: tui, api, cli）的项目
- When gen-test-cases 执行
- Then 生成的 test-cases.md 中 TUI 用例包含 golden file 断言 + 维度检查 + ≥2 个边界场景，API 用例包含契约验证 + 错误路径 + ≥2 个边界值，CLI 用例包含输出 golden file + 退出码 + 参数组合

---

## Story 2: 测试级别自动标记

**As a** Forge 用户（开发者）
**I want to** gen-test-cases 自动为测试用例标记 e2e 或 integration 级别（TUI/web-ui/mobile-ui → e2e，API/CLI → integration）
**So that** gen-test-scripts 能根据级别生成不同结构的测试代码，且我能在 test-cases.md 中直观区分测试层次

**Acceptance Criteria:**
- Given gen-test-cases 已读取 profile 的 verification-strategies.md
- When 生成测试用例
- Then TUI/web-ui/mobile-ui 用例的 Level 字段为 "e2e"，API/CLI 用例的 Level 字段为 "integration"，Level 字段覆盖率 ≥ 95%

---

## Story 3: Profile 策略文件定义验证策略

**As a** Profile 作者
**I want to** 在 profile 目录中新增 verification-strategies.md 文件，为每个 capability 定义验证维度、边界场景和测试数据要求
**So that** forge 能按我的框架特性生成针对性的测试用例，而不是使用通用策略

**Acceptance Criteria:**
- Given 一个 profile 目录（如 forge-cli/pkg/profile/profiles/go-test/）
- When 我在该目录创建 verification-strategies.md
- Then 文件必须为每个 capability 包含 `## <capability-key>` section，每个 section 包含 `### 验证维度`（≥3 条）和 `### 边界场景`（≥2 条），否则 gen-test-cases 拒绝该 profile 并输出验证错误："Strategy file validation failed for profile: X. Section <capability-key> missing required subsections (验证维度 ≥3, 边界场景 ≥2)."

---

## Story 4: 策略缺失时优雅降级

**As a** Forge 用户（开发者）
**I want to** 在 profile 没有 verification-strategies.md 时，gen-test-cases 仍然能正常工作（回退到当前行为）
**So that** 升级 forge 不会因为缺少策略文件而中断我的工作流

**Acceptance Criteria:**
- Given 一个没有 verification-strategies.md 的 profile
- When gen-test-cases 执行
- Then 输出 warning（"No verification strategy found for profile: X"），生成无 Level 字段的通用测试用例，不中断执行

---

## Story 5: Capability Key 不一致时明确报错

**As a** Profile 作者
**I want to** gen-test-cases 在 verification-strategies.md 的 capability sections 与 manifest.yaml 声明的 capabilities 不一致时，明确告诉我哪些 key 缺失或多余
**So that** 我能快速定位并修复不一致，而不是生成错误的测试用例

**Acceptance Criteria:**
- Given 策略文件包含 `## tui` + `## api`，但 manifest.yaml 声明 `tui` + `cli`
- When gen-test-cases 执行
- Then 中止执行并输出错误："Strategy/manifest mismatch. Missing in strategy: cli. Extra in strategy: api."

---

## Story 6: gen-test-scripts 按测试级别生成不同代码结构

**As a** Forge 用户（开发者）
**I want to** gen-test-scripts 根据测试用例的 Level 字段（e2e/integration）生成结构不同的测试代码
**So that** e2e 测试包含渲染截获和 golden file 比对逻辑，integration 测试包含 HTTP 断言或子进程退出码检查逻辑，而不是统一模板

**Acceptance Criteria:**
- Given test-cases.md 包含 Level=e2e（TUI capability）和 Level=integration（API capability）的测试用例
- When gen-test-scripts 执行
- Then e2e 用例生成的代码包含 golden file 读取函数（`os/exec` + golden file path 构造），测试文件写入 `tests/e2e/` 目录；integration 用例生成的代码包含 HTTP 断言或退出码检查函数（`net/http` + assert 库），测试文件写入 `tests/integration/` 目录；两类代码在 import 列表、assertion 方式、目录结构三方面均存在至少一处差异

---

## Story 7: test-cases.md 模板包含 Level 和 Interface 字段

**As a** Forge 用户（开发者）
**I want to** gen-test-cases 输出的 test-cases.md 中每个用例包含 Level（e2e/integration）和 Interface（TUI/API/CLI/web-ui/mobile-ui）字段
**So that** 下游 gen-test-scripts 能直接读取这两个字段选择代码生成策略，而不需要二次推断 interface 类型

**Acceptance Criteria:**
- Given 一个使用 go-test profile（capabilities: tui, api, cli）的项目，且 profile 包含有效的 verification-strategies.md
- When gen-test-cases 执行完成
- Then 输出的 test-cases.md 中每个测试用例包含 `Level: e2e` 或 `Level: integration` 字段，值与 interface 类型映射一致（tui→e2e, api→integration, cli→integration）；每个用例包含 `Interface: <type>` 字段，值为 capability 对应的 interface 类型；Level 和 Interface 字段覆盖率 ≥ 95%

---

## Story 8: eval-test-cases 按类型化验证完整度评分

**As a** Forge 用户（开发者）
**I want to** eval-test-cases 在评估 test-cases.md 时，额外检查类型化验证的完整度（是否包含 interface 类型对应的验证维度、边界场景、测试数据）
**So that** 我能知道生成的测试用例是否充分利用了 profile 定义的类型化策略，而不是仅靠通用评分维度

**Acceptance Criteria:**
- Given test-cases.md 包含 TUI 用例（Interface=TUI, Level=e2e）且 profile 的 verification-strategies.md 为 TUI 定义了 golden file 对比 + 维度检查 + CJK 边界场景
- When eval-test-cases 执行评分
- Then 评分报告包含"类型化验证完整度"维度，检查项包括：(1) 该 interface 类型的验证维度是否全部覆盖（≥3/3），(2) 边界场景是否包含至少 2 个策略定义的场景，(3) 测试数据是否满足策略中的数据要求；完整度得分 = 已覆盖检查项数 / 总检查项数 × 权重分值；当无策略文件时该维度得分率 ≥ 0.8（不因缺少策略而惩罚）
