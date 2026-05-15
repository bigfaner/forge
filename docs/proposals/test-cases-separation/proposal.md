---
name: test-cases-separation
status: draft
created: 2026-05-15
---

# Proposal: test-cases.md 与 gen-test-scripts 的职责严格分离

## Problem

milestone-map 功能的 e2e 测试中，test-cases.md 包含了 44 个 provisional testid。gen-test-scripts agent 直接引用了这些 provisional testid，没有去前端 TSX 验证存在性，导致 19/24 个 UI 测试超时失败（修复成本 2h + 17h agent stall）。

当前 gen-test-cases 的 Element 字段虽然设计为引用 sitemap 语义 ID，但存在被误用为 provisional testid 容器的风险。gen-test-scripts 虽有 Code Reconnaissance 步骤，但没有硬性规则禁止引用 test-cases 中的 testid 类值。

## Solution

严格分离两个阶段的职责：
- **gen-test-cases** 输出只包含：测试场景、操作步骤（自然语言描述 UI 交互）、预期结果、前置条件
- **gen-test-scripts** 从源码提取所有技术实现细节（testid、定位策略、显示文本）

具体改动：移除 test-cases.md 中的 Element 字段，gen-test-scripts 强制源码优先。

## Alternatives

| 方案 | 优点 | 缺点 |
|------|------|------|
| **A: 严格分离（选定）** | 根治 provisional testid 问题，gen-test-scripts 被迫读源码 | gen-test-scripts 失去 Element 提示，需要更多源码读取 |
| B: 保留 Element + 添加 HARD-RULE | 最小变更 | Element 字段仍可能被误用；治标不治本 |
| C: 什么都不做 | 无变更成本 | 下次还会重蹈 milestone-map 覆辙 |

## Scope

### In Scope

1. **gen-test-cases SKILL.md**：移除 Element 字段相关规则，添加 HARD-RULE 禁止 provisional testid/selectors/实现细节
2. **gen-test-cases templates/test-cases.md**：移除 Element 字段
3. **gen-test-scripts SKILL.md**：移除 Element 字段处理，移除 `sitemap-missing` 相关逻辑，添加 HARD-RULE 强制源码读取
4. **eval-test-cases SKILL.md + rubric.md**：Dimension 3 web-ui 部分从 "Route & Element Accuracy" 改为 "Route Accuracy"，移除 "Elements are identifiable" 评估项
5. **web-playwright generate.md**：更新 Integration Tests 部分的 locator 引用

### Out of Scope

- breakdown-tasks、run-e2e-tests、graduate-tests 等下游技能（它们不直接引用 Element 字段）
- go-test、pytest、java-junit、rust-test 等 profile（generate.md 中无 Element 字段引用）
- 已生成的 test-cases.md 文件（历史数据不追溯）
- maestro generate.md（Element 引用为通用 UI 测试术语，非 sitemap Element 字段）

## Risks

| Risk | Likelihood | Mitigation |
|------|-----------|------------|
| gen-test-scripts 在没有 Element 提示时源码读取不完整 | Medium | 保留 Step 1.5 Code Reconnaissance 的 HARD-RULE，Fact Table Completeness Gate 已存在 |
| 已有 test-cases.md 包含 Element 字段，gen-test-scripts 需兼容 | Low | gen-test-scripts 移除 Element 处理逻辑，旧文件中的 Element 字段被忽略即可 |
| eval-test-cases Dimension 3 分值变化影响现有评分 | Low | 总分不变（200pts），重新分配到 Route Accuracy 维度即可 |

## Success Criteria

- [ ] gen-test-cases 输出的 test-cases.md 不包含任何 Element 字段、testid、CSS selector
- [ ] gen-test-scripts 在生成时强制读取源码构建 Fact Table，不依赖 test-cases 中的任何定位信息
- [ ] eval-test-cases 不再评估 Element 字段质量
- [ ] web-playwright generate.md 的 locator 策略引用更新为纯源码驱动
