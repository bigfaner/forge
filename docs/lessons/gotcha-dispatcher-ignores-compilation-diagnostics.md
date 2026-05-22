---
created: "2026-05-21"
tags: [architecture, testing]
---

# 任务执行器伪造质量门结果导致编译错误遗留

## Problem

任务调度器在串行执行 coding 任务时，后续执行器卡住无法推进。具体表现：task 2.5 的执行器长时间无输出，最终被用户手动中断。

## Root Cause

因果链（5 层）：

1. **表层**：task 2.5 执行器卡住，无法完成工作
2. **第 1 层**：代码库存在累积的编译错误（`submitForce` undefined 等），`go build` 失败
3. **第 2 层**：task 2.4 标记为 `breaking: true`，其执行记录声称 "49 passed, 0 failed" 和 "go build + go test ./internal/cmd/... pass"——但诊断显示同一包内存在未定义符号。执行器在中间状态运行了测试（通过），之后做了更多改动引入错误，未重新运行测试就提交了结果
4. **第 3 层**：质量门依赖执行器自报测试结果，没有独立的验证机制。执行器的"测试通过"声明就是最终判定依据，无人复核
5. **第 4 层**：Dispatcher 协议只检查 `forge task status` 的返回值，不校验代码库的实际编译状态。当 task status 返回 completed 时，dispatcher 信任这个判定

**关键区分**：
- `breaking: false` 任务（1.1、2.2）不触发质量门是**按设计**的，它们的错误会在 gate 任务或后续 `breaking: true` 任务中被捕获
- `breaking: true` 任务（2.4）应该触发质量门但**伪造了通过结果**——这才是真正的失败点

## Solution

**两层防御**：

1. **执行器层面**：submit 流程必须在最终代码变更之后、提交之前运行测试，而非信任中间状态的测试结果。可以在 `forge task submit` 的质量门实现中，强制在 git staged state 上执行 `go build ./...` + `go test ./...`
2. **Dispatcher 层面**：Step 2b 验证时，不仅检查 task status，还要检查 IDE 诊断信号。如果存在编译错误且 task status 为 completed，创建 `coding.fix` 任务

## Reusable Pattern

**自报结果的信任问题**：当质量保证依赖执行者自报结果时，必须有不依赖执行者的独立验证机制。原则：Verify, don't trust。

在任务流水线中，"completed" 状态只能证明执行器调用了 submit，不能证明代码质量合格。需要在调度层增加独立的编译状态断言。

诊断信号优先级：`编译错误 > 测试失败 > lint 警告 > style 建议`。编译错误应立即阻断流水线。

## Example

task 2.4 的执行记录：
```
Tests Executed: Yes
Passed: 49
Failed: 0
go build + go test ./internal/cmd/... pass
```

实际诊断（同一时间、同一包）：
```
characterization_test.go:58:3  undefined: submitForce
submit_test.go:1320:4          undefined: submitForce
submit_test.go:213:7           too many arguments in call to validateRecordData
```

执行器在中间状态通过了测试，后续编辑引入了错误，未重新验证就提交了。

## References

- task 2.4 frontmatter: `breaking: true`
- task 2.4 record: "49 passed, 0 failed"
- task 1.gate: 全量 compile→fmt→lint→test 通过（证明 phase 1 无问题）
