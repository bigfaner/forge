---
created: "2026-05-28"
tags: [testing, architecture]
---

# Quality-gate fix 任务 scope 过大导致 agent 卡住

## Problem

Quality-gate hook 检测到 `just test` 失败后自动创建 fix-1 任务。该任务包含 4 个 test suite 共 20+ 个失败，agent 执行时卡住被用户中断。

## Root Cause

因果链（3 层）：

1. **表面现象**：fix-1 任务执行后 agent 长时间无响应，用户手动中断
2. **直接原因**：
   - 任务 scope 包含 20+ 失败，跨越 4 个 test suite（automated-test-orchestration、quality-gate、test-generation、test-suite-health）
   - Concise error 只展示输出尾部（test-suite-health），agent 需要读 raw-output.txt 才能看到完整失败列表
   - 部分失败与当前改动无关（feature-management 测试），agent 无法区分"我引入的"和"预存在的"
3. **根因**：
   - **quality-gate hook 创建 fix 任务的粒度为"整次 test 运行"**：一个 fix 任务覆盖所有 test suite 的所有失败。当失败来自多个不相关的原因时，单个 fix 任务 scope 过大
   - **fix 任务缺少失败分组信息**：任务描述只有尾部输出，没有按 suite 或 root cause 分组的失败摘要，agent 需要自行分析全部失败才能确定修复策略

## Solution

### 第一层：按 test suite 目录机械拆分（不需要 LLM）

Hook 从失败输出中按 `FAIL forge-tests/<suite>` 行分组，每个 suite 创建独立 fix task：

```bash
# 从 raw-output.txt 提取 failing suites
grep "^FAIL\s" tests/results/raw-output.txt | awk '{print $2}' | sort -u
# → forge-tests/automated-test-orchestration
# → forge-tests/quality-gate
# → forge-tests/test-generation
# → forge-tests/test-suite-health
```

每个 fix task 的 Concise error 包含该 suite 的完整失败详情（`grep` 该 suite 的全部 `--- FAIL` 行），而非当前的"只展示输出尾部"。

效果：20+ 失败自动拆为 4 个 fix task，每个 scope 收窄到单个 suite。

### 第二层：基线对比过滤预存在失败

解决更根本的问题——区分"我引入的"和"本来就坏的"：

1. 任务开始前（或 Task 1 baseline 阶段）跑一次 `just test`，记录已知失败到 `tests/results/baseline.txt`
2. Hook 创建 fix task 时，将当前失败与基线做 diff
3. 只有**增量失败**（当前有但基线没有的）才创建 fix task

即使第一层不拆分，基线过滤也能将 scope 自然收窄到与当前改动相关的失败。

### 两层组合

| 场景 | 第一层（suite 拆分） | 第二层（基线过滤） | 效果 |
|------|---------------------|-------------------|------|
| 20 失败，全由当前改动引入 | 4 个 fix task（按 suite） | 无基线可过滤 | 每个 fix task 5 个失败 |
| 20 失败，5 个预存在 + 15 个新增 | 4 个 fix task | 过滤掉 5 个预存在 | 每个 fix task 3-4 个新增失败 |
| 3 失败，全由同一 root cause | 1-2 个 fix task | 无基线可过滤 | 单一 fix task 即可 |

两层都不需要 LLM 分析——纯 shell 脚本（grep + awk + diff）即可实现。

## Reusable Pattern

- **fix 任务的 scope 应与 root cause 的 scope 匹配**：如果 20 个失败来自 3 个不同的 root cause，应该有 3 个 fix 任务，不是 1 个
- **预存在失败的隔离**：fix 任务描述中应标注基线失败（已知的、与当前变更无关的），agent 只修复增量失败
- **与 [[gotcha-split-rules-operational-blindness]] 同源**：任务粒度过大的问题不仅出现在业务任务拆分中，也出现在 fix 任务自动生成中

## Related Files

- `plugins/forge/hooks/` — quality-gate hook 的 fix 任务生成逻辑
- `tests/results/raw-output.txt` — 完整测试输出
- [[gotcha-split-rules-operational-blindness]] — 业务任务操作粒度盲视
