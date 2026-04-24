# eval-prd 应作为调度者，用 scorer/reviser subagent 实现对抗迭代

## Problem

执行 `/eval-prd --target 95 --iterations 6` 时：
1. 主会话把整个"评估 + 迭代"逻辑委托给单个 `general-purpose` agent
2. agent 跑偏，评估完成后额外生成了技术设计文档
3. `--target` / `--iterations` 参数被忽略，只执行了一次评估，没有迭代循环

## Root Cause

**设计缺陷**：eval-prd 把调度责任外包给了 agent，而不是由主会话充当调度者。

`general-purpose` agent 没有"完成即停止"的约束，也没有结构化输出格式，主会话无法可靠解析其分数并控制循环。

对比 eval-proposal 的正确架构：

```
主会话（调度者）
  │
  ├─ 第 N 轮：
  │   ├── Agent(zcode:proposal-scorer)  → SCORE: X/100 + ATTACKS
  │   ├── 主会话解析分数 → score >= target? → 停止
  │   └── Agent(zcode:proposal-reviser) → 修订文件
  │
  └─ 最终报告
```

eval-proposal 的循环控制在主会话，scorer/reviser 各自单一职责，输出格式固定可解析。

## Solution

执行带 `--target` / `--iterations` 参数的 eval 命令时，主会话必须自己控制循环：

1. **识别迭代信号**：看到 `--target` + `--iterations` 参数，立即切换到 eval-proposal 的调度模式
2. **评分**：`Agent(subagent_type: zcode:proposal-scorer)`，prompt 中指定 PRD 路径和 PRD 专用评分维度
3. **主会话解析**：从 scorer 输出中提取 `SCORE: X/100`，判断是否达到目标
4. **修订**：`Agent(subagent_type: zcode:proposal-reviser)`，传入 ATTACKS
5. **循环**：回到步骤 2，直到达到目标分数或耗尽迭代次数

若 zcode 未注册专用 PRD scorer/reviser，复用 `zcode:proposal-scorer` / `zcode:proposal-reviser`，在 prompt 中明确指定评估对象为 PRD 文档及其评分维度。

## Key Takeaway

**调度逻辑必须在主会话，不能委托给 agent。**

- `--target` / `--iterations` 参数只有在主会话控制循环时才有意义
- 专用 subagent（scorer/reviser）的价值在于：职责单一、输出格式固定、不会"顺手"做额外的事
- 遇到任何带迭代参数的 eval 命令，先对标 eval-proposal 的架构，而不是按技能描述字面执行
