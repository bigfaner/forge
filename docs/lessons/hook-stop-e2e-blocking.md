# Stop Hook 的 e2e 测试阻塞机制

## Problem

对话结束时，Stop hook 反复输出 `e2e tests failed 3 times. Manual intervention required.`，阻塞 Claude 停止响应。即使当前工作与任务执行完全无关（如 `/extract-design-md`），也会触发。

## Root Cause

**因果链（3层）**:

1. **触发层**: forge 插件 `hooks.json` 的 `Stop` hook 配置了 `task all-completed`，在每次 Claude 停止响应时执行
2. **逻辑层**: `task all-completed` 检测到当前 feature 所有任务已完成 → 自动运行 `just test-e2e --feature <slug>` → e2e 测试失败 → 触发 auto-fix 重试（最多3次） → 全部失败 → 输出 `{"continue":false,"stopReason":"e2e tests failed 3 times."}` 阻止停止
3. **根因层**: e2e 测试本身已过时/损坏（API 返回值变化、页面结构变化、断言不匹配），不是当前工作的问题

**关键**: `Stop` hook 不区分"任务执行结束"和"其他工作结束"——只要 active feature 的任务全完成但 e2e 失败，就会阻塞每次对话。

## Solution

1. **即时修复**: 修复或跳过过时的 e2e 测试
   - `just test-e2e --feature <slug>` 查看具体失败原因
   - 更新测试预期值（如 API 返回 200 而非 201）
   - 或跳过不适用的测试用例

2. **长期方案**:
   - 使用 `task feature` 切换到没有失败测试的 feature
   - 或修复 e2e 测试使其通过
   - Stop hook 的阻塞只在 e2e 全部通过或重试耗尽后才停止

## Key Takeaway

当 Stop hook 持续阻塞时，排查路径是：
1. 检查 `hooks.json` 中 `Stop` hook 执行的命令
2. 手动运行该命令查看输出（本例中是 `task all-completed -v`）
3. 定位具体失败原因（通常是过时的测试而非代码 bug）
4. 修复测试或切换 active feature 来解除阻塞
