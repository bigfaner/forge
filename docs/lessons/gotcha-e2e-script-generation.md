# E2E 测试脚本生成：SKILL.md 流程结构性盲区导致测试与代码脱节

## Problem

生成的 e2e 测试脚本存在大量基础错误，导致需要多轮修复才能运行：
- 相对路径计算错误（多了一层 `../`）
- API 路径前缀假设错误（`/api/v1/` vs 实际的 `/v1/`）
- 端口号与实际运行服务不符
- 路由名称错误（`/items` vs `/main-items`，`/members/invite` vs `/members`）
- 测试用例引用了不存在的路由
- 测试依赖硬编码的预置数据，但数据库中没有这些数据
- UI testid 猜测错误（`map-card` vs 实际的 `map-card-${bizKey}`），24 个测试全部超时

## Root Cause

**不是 agent 不遵守指令，而是 SKILL.md 流程存在结构性盲区。**

### 第一层：原始版本没有 Code Reconnaissance 步骤

SKILL.md 最初的流程是从 "读 test cases" 直接到 "生成 spec files"，中间没有读源码的步骤。agent 完全按照指令执行——但指令本身缺少关键环节。

**Fix**：2026-04-29 加入 Step 1.5 Code Reconnaissance，要求读 router/config/handler/auth/CLI 源码并建 Fact Table。

### 第二层：Step 1.5 只覆盖后端，遗漏了 Frontend UI

Step 1.5 的 Required reads 表有 5 类源文件（Router、Config、API handlers、Auth、CLI），但没有 "Frontend UI components" 类别。UI 选择器验证依赖两条路径：
1. **sitemap.json**（Step 2）——但 sitemap 不含 `data-testid`，且新页面可能不在 sitemap 中
2. **Step 1.5 Fact Table** ——但没有 frontend 源码读取类别

当两条路径都失效时（新页面 + Step 1.5 不读前端），agent 标注 "provisional" 后继续生成——因为 SKILL.md 没有要求它在无法验证 testid 时停止。

**Fix**：在 Step 1.5 Required reads 表加入 "Frontend UI components" 行，要求 grep `data-testid` 并记录到 Fact Table。

### 第三层：sitemap 缺失路由时的 fallback 是空话

SKILL.md 写了 "use Fact Table DOM structure from Step 1.5"，但 Step 1.5 没有 frontend 读取——这个 fallback 引用了一个不存在的数据来源。

**Fix**：sitemap 缺失路由时 emit WARNING 并建议 re-run `/gen-sitemap`，同时用 Step 1.5 的 Fact Table 继续推断（现在 Fact Table 有 frontend 数据了）。

### 第四层：策略和语法在 SKILL.md 与 generate.md 之间混在一起

Locator 策略（用什么选择器、优先级）和 Auth 策略（如何分类、如何缓存）同时出现在 SKILL.md 和 web-playwright/generate.md 中。两份不同步就有冲突——web-playwright 的 generate.md 完整复制了 Auth Classification 表，如果 SKILL.md 增加新类别，generate.md 不会自动更新。

**Fix**：明确分工原则——SKILL.md 负责策略决策（选什么、优先级），generate.md 只负责框架语法（怎么写）。

具体改动：
- Step 3 Map Locators：策略收归此步骤，含 integration test + testid HARD-RULE
- web-playwright generate.md：Locator Mapping 改为 Locator Syntax（去掉策略决策），Auth Classification 改为 Auth Implementation（去掉分类表）

### 第五层：UNKNOWN 值没有完整性门槛

Step 1.5 HARD-RULE 说 "note as UNKNOWN, do not fabricate"，但没说 UNKNOWN 之后该怎么办。agent 标了 UNKNOWN 然后继续生成——和原来标 "provisional" 继续生成是同一个模式。

**Fix**：加入 Fact Table Completeness Gate——如果某个 test type 的全部关键 Fact Table 值都是 UNKNOWN，跳过该类型并 WARNING。

### 第六层：UI 探测命令硬编码路径

Step 4 的 UI probe 用 `grep -r ... src/`，但前端代码可能在 `frontend/src/`、`web/src/` 等。到 Step 4 执行时 Step 1.5 已经做过前端 grep，重复验证没有意义。

**Fix**：UI probe 改为查 Fact Table 里有没有 Frontend 行。

## Open Question：PRD→代码路由翻译

gen-test-cases 忠于 PRD（Route 字段按 PRD 原文写），gen-test-scripts 负责翻译为实际代码路由（通过 Fact Table）。但当前 Fact Table 没有明确承担"PRD Route → Code Route 翻译"的职责——Step 1.5 只说 "Use corrected routes where available"，没定义纠正规则。

**待决定**：是否在 Step 1.5 HARD-RULE 中加入"当 Fact Table 路由和 test case Route 不一致时，用 Fact Table 覆盖"？

## Architecture Principle

**SKILL.md 负责策略决策，generate.md 负责框架语法。**

- 策略：选什么 locator、auth 如何分类、完整性门槛
- 语法：Playwright 怎么写 locator、Go 怎么写 auth header

混合两者导致重复定义和不同步风险。

## Timeline

| 日期 | 事件 |
|------|------|
| 2026-04-28 | 原始问题发生（pm-work-tracker），API/路径/端口/路由/数据全部猜错 |
| 2026-04-29 | Step 1.5 Code Reconnaissance 加入 SKILL.md，覆盖后端源码 |
| 2026-05-13 | UI testid 问题发生（pm-work-tracker milestones 页面），Step 1.5 不覆盖 frontend |
| 2026-05-14 | 修复：前端 testid 被改为静态值以适配测试（方向有误，应反过来） |
| 2026-05-15 | 全面结构性修复：Step 1.5 补 Frontend、策略/语法分离、Fact Table 完整性门槛、UI probe 去 grep 化 |

## Anti-Pattern

1. **根据 PRD/设计文档推断实现细节** → 测试与实际代码脱节
2. **SKILL.md 设计 fallback 路径但不提供 fallback 所需的数据来源** → agent 按"指令"执行却产出错误结果
3. **策略和语法混在两处定义** → 修改一处时另一处不同步
4. **"provisional"/"UNKNOWN" 作为逃生舱但不设门槛** → agent 永远不会停下来
