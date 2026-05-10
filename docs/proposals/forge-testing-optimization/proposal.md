---
created: 2026-05-10
author: faner
status: Draft
---

# Proposal: Forge 测试能力优化

## Problem

Forge 的 E2E 测试生成 pipeline（gen-test-cases → gen-test-scripts → run-e2e-tests → graduate-tests）在真实项目中暴露了三类系统性问题：**生成的脚本质量低**、**缺乏生成后校验**、**缺乏运行时诊断**。

### Evidence

来源：pm-work-tracker 项目（350+ 测试用例，50 个 spec 文件）和 train-recorder 项目的生产事故。

#### A. 生成的脚本质量低

| 反模式 | 出现次数 | 后果 |
|--------|---------|------|
| `waitForTimeout` 滥用 | 100+ 处 | 每轮多等 60+ 秒，且 flaky（快了不够，慢了浪费） |
| `beforeEach` 重复登录 | 10+ 文件 | 30-test 套件多等 30-90 秒 |
| 无 `afterAll` 清理 | 14 文件 | 测试数据跨 run 污染 |
| 超大 serial 套件 (32+/27+) | 6 文件 | 1 个失败 → 全部级联失败 |
| CSS class 选择器 (`.ant-*`) | 多处 | UI 库升级即崩 |
| 三种 API 调用模式混用 | 全局 | 不一致的错误处理和超时 |

#### B. 缺乏生成后校验

gen-test-scripts 定义了 beforeAll Safety、Traceability 等规则，但唯一的后置检查是 TypeScript 编译（`tsc --noEmit`）。实际生成的代码可以违反所有规则而不被检测：

- gen-test-scripts 的 SKILL.md 中有 6 条 HARD-RULE，但无结构性验证
- TC ID 覆盖率无检查：test-cases.md 中有 10 个 TC，生成的 spec 可能只包含 8 个
- eval-test-cases 的 Step Actionability < 20 阻塞阈值是建议性的，用户可以绕过
- Element 字段是可选的，agent 可能匹配到错误的 locator

#### C. 缺乏运行时诊断

- run-e2e-tests 在 >30% UI 测试同时失败时，agent 直接逐个修复测试，不检查 app 是否崩溃
- fix task agent 启动 dev server 试图"亲眼确认"，陷入 npm install 循环（20+ 次重试，0 次文件编辑）
- 服务器生命周期管理不当：Playwright webServer 只检查 TCP 端口可达，不验证应用层就绪

### Urgency

每个生成低质量脚本的 feature 平均消耗 3-5 轮 agent 交互修复。pm-work-tracker 的 13 个 feature test suite 中只有 1 个通过毕业。校验脚本可在生成阶段立即拦截问题，避免下游多轮修复成本。

## Proposed Solution

### Phase 1：运行时诊断 + 模板改进（已完成）

#### 1.1 App Health First Gate

run-e2e-tests 增加 >30% 失败率时的诊断门控：

| 失败比例 | 诊断方向 | 首要操作 |
|----------|---------|---------|
| >30% UI 测试同时失败 | App 健康问题 | 检查截图 + 依赖兼容性 |
| 10-30% | 可能是测试问题 | Spot check 2-3 个失败截图 |
| <10% | 测试/选择器问题 | 逐个修复 |

诊断流程：截图检查 → 依赖兼容性 → 手动渲染验证 → 确认健康后再修测试。

#### 1.2 Fix Task 边界规则

error-fixer agent + fix-task 模板增加边界约束：

- 禁止启动 dev server
- npm install 最多重试 3 次
- 禁止运行 e2e 测试（由 dispatcher 统一执行）
- 正确流程：读测试 + 对比 DOM → 修改 → `just test` → record

#### 1.3 模板 Anti-Pattern 改进

gen-test-scripts SKILL.md + 模板文件增加 4 条 HARD-RULE：

| 规则 | 替代方案 |
|------|---------|
| 禁止 `waitForTimeout` | `waitForApiAction` / `expect().toBeVisible()` |
| Serial 套件上限 15，必须有 afterAll 清理 | 超过 15 拆分为多个 serial 块 |
| UI 测试登录用 `beforeAll`，禁止 `beforeEach` | Playwright storageState 或 beforeAll 共享 |
| 禁止 CSS class 选择器和 DOM 遍历 | role-based locator 或 data-testid |

新增 `waitForApiAction` helper：封装 `page.waitForResponse()` + action 模式。

### Phase 2：程序化校验（本轮实现）

#### 2.1 ts-morph 校验脚本

新建 `validate-specs.mjs`，使用 ts-morph AST 解析生成的 spec 文件，执行 8 条校验规则：

**ERROR 级（阻塞后续流程）：**

| ID | 规则 | 检测方式 |
|----|------|---------|
| E1 | 禁止 `waitForTimeout` / `setTimeout` | AST: CallExpression 中包含这些名称 |
| E2 | TC ID 全覆盖 | grep: test-cases.md 中所有 `TC-\d+` 必须在 spec 中出现 |
| E3 | 每个 test() 有 Traceability 注释 | AST: 检查 test() 调用上方或内部有无 `// Traceability:` |
| E4 | 禁止 DOM 父级遍历 `locator('..')` | AST: 字符串参数包含 `..` 的 locator 调用 |

**WARNING 级（报告但不阻塞）：**

| ID | 规则 | 检测方式 |
|----|------|---------|
| W1 | serial suite > 15 个 test() | AST: 统计 serial describe 内的 test() 数量 |
| W2 | serial suite 无 afterAll | AST: 检查 serial describe 内有无 afterAll 调用 |
| W3 | beforeEach 中有 login 调用 | AST: 找 beforeEach 回调中的 login/loginViaUI 调用 |
| W4 | CSS class 选择器 | AST: 找 locator 参数以 `.` 开头的字符串 |

#### 2.2 集成方式

- ts-morph 作为 devDependency 加入 `gen-test-scripts/templates/package.json`
- 校验脚本随 helpers.ts 一起生成到 `tests/e2e/`
- task-cli 新增 `validate-specs` 命令，spawn Node 脚本执行校验
- gen-test-scripts SKILL.md 在 Step 4 后增加 Step 4.5 结构校验，调用 `task validate-specs`
- 校验结果：ERROR → 阻塞（标记 T-test-2 为 blocked），WARNING → 报告继续

#### 2.3 eval-test-cases Step Actionability 强制阻塞

gen-test-scripts 的 Prerequisites 增加检查：

```
如果 eval-test-cases 报告存在且 Step Actionability 得分 < 20：
  中止 gen-test-scripts
  提示用户：先修复 test-cases.md 的 Step Actionability 到 20 分以上
```

#### 2.4 gen-test-cases Element 字段改为必填

当前 Element 字段是可选的，agent 省略后 gen-test-scripts 用模糊匹配选择 locator，导致错误匹配。

改为：Element 字段必填。如果 sitemap.json 中对应页面没有足够的元素数据，gen-test-cases 应在 Route Validation 阶段报告缺失，引导用户先运行 `/gen-sitemap` 补全元素数据。

### Phase 3：TC-ID 全局唯一（后续独立 PR）

#### 3.1 TC-{module}-{seq} 格式

当前格式 `TC-001` 在多 feature 毕业到同一回归套件时会 ID 冲突。

新格式：`TC-{module}-{seq}`（如 `TC-auth-001`、`TC-items-001`）。

#### 3.2 模块定义来源

模块定义放在 `docs/sitemap/sitemap.json` 顶层。gen-test-cases 根据测试用例引用的路由/实体映射到模块。如果模块归属模糊，提示用户确认。

#### 3.3 影响范围

TC-ID 格式变更涉及 5 个 SKILL.md + task-cli + 所有模板文件 + 现有项目的所有 TC ID。作为独立 PR 单独处理。

### Phase 3 记录待做（不计入本轮）

| # | 改进 | 说明 |
|---|------|------|
| 11 | graduate-tests 毕业后重跑测试 | 当前毕业只做 TypeScript 编译 + Playwright test discovery，不执行测试。import path 改写后可能运行时失败 |

## Quality Gate Summary

```
gen-test-cases
  ├─ Step 3.5: Route Validation（已有）
  └─ Element 字段必填（新增）
       ↓
eval-test-cases
  └─ Step Actionability < 20 → 阻塞 gen-test-scripts（新增强制检查）
       ↓
gen-test-scripts
  ├─ Step 1.5: Code Reconnaissance（已有）
  ├─ Step 4: 生成 spec 文件
  ├─ Step 4.5: 结构校验 ← ts-morph AST 校验（新增）
  └─ Anti-pattern HARD-RULEs（新增）
       ↓
run-e2e-tests
  ├─ App Health First Gate（新增）
  └─ Failure Diagnosis（增强）
       ↓
graduate-tests
  └─ 毕业后重跑测试（待做）
```

## Alternatives Considered

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| 只加 SKILL.md 规则，不加程序化校验 | 零代码改动 | pm-work-tracker 证明 LLM 不遵守规则 | Rejected: 约束力不足 |
| 纯 grep 校验（不用 AST） | 零依赖，简单 | 无法检查 serial suite 大小、afterAll 存在性、beforeEach 位置 | Rejected: 覆盖不完整 |
| TypeScript Compiler API 直用 | 零新增依赖 | 代码量多 150 行，API 不友好 | Rejected: 维护成本高 |
| ts-morph 作为 e2e devDependency | 项目自带校验能力，CI 可用 | 每个项目多 ~5-8MB | **Selected** |
| 校验嵌入 gen-test-scripts（Route A） | 生成后立即校验 | skill 膨胀，校验和生成耦合 | Rejected: 职责不清 |
| 校验脚本 + task-cli 命令（Route B） | 单一职责，可独立调用，dispatcher 可触发 | 多一个 skill 维护 | **Selected** |

## Scope

### In Scope — Phase 1（已完成）

- run-e2e-tests SKILL.md：App Health First Gate
- error-fixer.md：E2E Fix Boundary Rules
- fix-task.md 模板（task-cli）：边界规则
- breakdown-tasks/templates/run-e2e-tests.md：App Health First 检查
- gen-test-scripts SKILL.md：Anti-Pattern HARD-RULEs
- gen-test-scripts/templates/helpers.ts：waitForApiAction helper
- gen-test-scripts/templates/playwright-ui.spec.ts：更新模式
- gen-test-scripts/templates/api.spec.ts：cleanup 模式

### In Scope — Phase 2（本轮实现）

- 新建 `gen-test-scripts/templates/validate-specs.mjs`：ts-morph 校验脚本（8 条规则）
- 修改 `gen-test-scripts/templates/package.json`：加入 ts-morph devDependency
- task-cli 新增 `validate-specs` 命令：spawn Node 脚本
- gen-test-scripts SKILL.md：Step 4.5 结构校验
- gen-test-scripts SKILL.md Prerequisites：Step Actionability < 20 强制阻塞
- gen-test-cases SKILL.md + 模板：Element 字段改为必填

### Out of Scope

- TC-{module}-{seq} ID 格式变更（Phase 3，独立 PR）
- sitemap.json 模块定义（随 TC-ID 变更一起）
- graduate-tests 毕业后重跑测试（记录待做）
- 现有项目的测试脚本迁移（手动处理）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| ts-morph 与项目 TypeScript 版本不兼容 | Low | High — 校验脚本无法运行 | package.json 中 ts-morph 版本与 typescript 版本对齐；校验失败时 fallback 到 WARNING 而非 ERROR |
| 校验脚本误报（合法代码被标记为 ERROR） | Medium | Medium — 阻塞正常流程 | Phase 2 初期可设为全 WARNING 模式，收集误报数据后再升级为 ERROR |
| Element 必填导致 gen-test-cases 在无 sitemap 时无法生成 | Medium | High — 阻断测试生成 | 如果 sitemap.json 不存在，Element 标记为 "sitemap-missing" 并在 test-cases.md 中添加 WARNING；gen-test-scripts 看到 "sitemap-missing" 时使用 Fact Table 中的实际 DOM 结构推断 |
| task-cli spawn Node 脚本在 Windows 上路径问题 | Medium | Medium — `task validate-specs` 失败 | 使用 `node` 命令而非绝对路径；task-cli 中处理路径分隔符 |
| 校验规则维护负担（Playwright API 变化） | Low | Low — 规则简单稳定 | 规则基于 Playwright 稳定 API（test, describe, locator），不依赖实验性特性 |

## Success Criteria

- [ ] validate-specs.mjs 脚本能检测 E1-E4 四种 ERROR 和 W1-W4 四种 WARNING
- [ ] ts-morph 在 tests/e2e/package.json 中作为 devDependency 存在
- [ ] `task validate-specs` 命令能执行校验并返回结构化输出
- [ ] gen-test-scripts SKILL.md 包含 Step 4.5 结构校验步骤
- [ ] gen-test-scripts 在 eval-test-cases Step Actionability < 20 时中止
- [ ] gen-test-cases SKILL.md 和模板中 Element 字段标记为必填
- [ ] 对 pm-work-tracker 的现有 spec 文件运行校验，能检测出 waitForTimeout、缺失 cleanup、缺失 Traceability 等已知问题

## Implementation Plan

### Phase 1 — 已完成

9 个文件已修改，涵盖运行时诊断、fix task 边界、anti-pattern 规则和模板改进。

### Phase 2 — 本轮实现

| Step | Task | 预计时间 |
|------|------|---------|
| 1 | 创建 validate-specs.mjs 校验脚本（ts-morph） | 2h |
| 2 | package.json 模板加入 ts-morph | 10min |
| 3 | task-cli `validate-specs` 命令 | 1h |
| 4 | gen-test-scripts SKILL.md Step 4.5 | 30min |
| 5 | gen-test-scripts Prerequisites 加 Step Actionability 检查 | 20min |
| 6 | gen-test-cases Element 字段必填 | 30min |
| 7 | 集成测试：对 pm-work-tracker 运行校验 | 30min |

### Phase 3 — 后续独立 PR

| Step | Task | 依赖 |
|------|------|------|
| 1 | sitemap.json 加 module 定义 | 需要设计 module schema |
| 2 | TC-{module}-{seq} 格式变更（5 个 SKILL.md + task-cli + 模板） | Phase 2 Step 1 |
| 3 | graduate-tests 毕业后重跑测试 | Phase 2 稳定后 |
