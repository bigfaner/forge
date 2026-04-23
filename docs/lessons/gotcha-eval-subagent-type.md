# 评估任务应使用 doc-scorer / doc-reviser，而非 error-fixer

## Problem

在执行 `/eval-prd` 和 `/eval-design` 时，将 Agent 的 `subagent_type` 设为 `zcode:error-fixer`。评估虽然能运行，但语义错误。

## Root Cause

`error-fixer` 是为修复编译错误和测试失败设计的，不是文档评估。选错 subagent 是因为没有检查可用的 subagent 类型列表，凭印象猜测。

## Solution

文档评估任务（eval-prd、eval-design、eval-proposal 等）应使用：
- `doc-scorer` — 评分
- `doc-reviser` — 修订

## Key Takeaway

调用 Agent 工具前，先确认 subagent_type 是否与任务语义匹配。评估/评分类任务 → doc-scorer；修订类任务 → doc-reviser；错误修复类任务 → error-fixer。
