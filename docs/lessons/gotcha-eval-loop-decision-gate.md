# Eval Loop Decision Gate Must Fire After Every Scorer Run

## Problem

`/eval-proposal` 达到目标分数后，用户说"再迭代"，AI 继续运行 scorer + reviser，
在已超目标（93/100）后仍然调用 reviser，违反了 skill 的决策门规则。

## Root Cause

AI 将"用户要求继续迭代"解读为"跳过决策门，直接运行完整的 scorer + reviser 循环"。
实际上决策门（score >= target → stop）必须在**每次 scorer 返回后**执行，无论用户是否要求继续。

## Solution

每次 scorer 返回后，先检查决策门：

```
scorer 返回分数
  ↓
score >= target?
  → YES: 输出最终报告，停止。不运行 reviser。
  → NO + 还有迭代次数: 运行 reviser，然后回到 scorer
  → NO + 无迭代次数: 输出最终报告，停止
```

用户说"再迭代"时，正确行为是：
1. 运行一次 scorer
2. 检查决策门
3. 如果 score >= target → 报告并停止，不运行 reviser
4. 如果 score < target → 运行 reviser，继续循环

## Key Takeaway

**决策门是硬性规则，不受用户"继续"指令影响。**
"继续迭代"的语义是"再跑一次 scorer"，而不是"绕过决策门跑完整循环"。
每次 scorer 返回后必须先判断是否达标，达标则立即停止，不再调用 reviser。
