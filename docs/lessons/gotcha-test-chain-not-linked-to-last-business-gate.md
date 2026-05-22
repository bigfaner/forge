---
created: "2026-05-22"
tags: [architecture]
---

# resolveBreakdownDeps 未将测试链依赖到业务链尾部导致执行顺序错误

## Problem

`forge task index` 自动生成的测试链任务在执行时有依赖缺陷：测试链的依赖是内部串联的，但缺少 Test-1→LastBusinessGate 这条关键边。导致 `forge task claim` 在 Phase 2 完成后直接领取测试链任务，跳过了 Phase 3 的业务任务。

## Root Cause

因果链（4 层）：

1. **表层**：Phase 2 完成后 dispatcher 跳过 Phase 3 直接执行测试任务
2. **第 1 层**：测试链头部 `T-test-gen-cases` 没有任何前序依赖（deps: []），测试链与业务链是两条独立的 DAG 分支，都只依赖 2.gate

3. **第 2 层**：`resolveBreakdownDeps` 函数（`autogen.go:253-289`）只负责测试链**内部**的依赖编排：

   ```
   gen-cases → eval-cases → gen-scripts → run-e2e → graduate → verify
   ```

   但从不把测试链头部连接到业务链的尾部（`3.gate`）。

4. **第 3 层**：模板中设计了 `{{LAST_BUSINESS_TASK_ID}}` 占位符，但：
   - 无代码自动将其替换为实际的最后一个 gate ID
   - 验证器 `validateFirstTestTaskTemplate` 只检查占位符**是否已被替换**（发现未替换则报错），但不负责替换
   - 于是占位符未被填充，依赖边从未被建立

## Solution

**修复方案**：在 `BuildIndex` 的 Step 7.5（`GenerateTestTasks`）中，在调完 `resolveBreakdownDeps` 之后，检测所有业务任务的最后一个 gate/summary ID，将其添加到 `T-test-gen-cases` 的 `Dependencies` 中。

代码位置：`forge-cli/pkg/task/build.go` 的 Step 7.5 部分（~line 278-313），在 `GenerateTestTasks` 调用之后追加：

```go
// 将测试链挂到最后一个业务 gate 之后
lastGate := findLastBusinessGate(index)
if lastGate != "" {
    for i, td := range testTasks {
        if td.ID == "T-test-gen-cases" {
            testTasks[i].Dependencies = append(
                testTasks[i].Dependencies, lastGate)
            break
        }
    }
}
```

## Reusable Pattern

**自动生成的链必须接到业务链上**：任何自动生成的后续任务（测试、验证、文档评估等）必须以其前驱业务链的最后一个同步点为依赖。依赖图的"入口节点"不能是浮动的。

## References

- `forge-cli/pkg/task/autogen.go:253-289` — `resolveBreakdownDeps`
- `forge-cli/pkg/task/build.go:278-313` — Step 7.5 测试链生成
- `forge-cli/internal/cmd/validate_index.go:231-233` — 占位符验证逻辑
- task index 中的 T-test-* 任务依赖配置