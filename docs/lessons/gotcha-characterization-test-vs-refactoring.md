---
created: "2026-05-22"
tags: [testing, architecture]
---

# 定标测试与重构行为变化冲突导致执行器卡死

## Problem

任务 `fix.7` 的执行器被分派后卡住无输出，被用户手动中断。任务内容是更新 task 0.1 的定标测试以适配 task 2.8 重构后的行为变化。

## Root Cause

因果链（4 层）：

1. **表层**：`fix.7` 的执行器卡住不返回
2. **第 1 层**：执行器被要求更新 4 个定标测试，但这些测试记录了**旧行为的全部细节**（包括已知缺陷），执行器无法判断是应该（a）更新测试以匹配新行为，还是（b）修复代码以保持旧行为
3. **第 2 层**：定标测试（characterization tests）的设计目的是"记录当前实际行为，无论对错"。但当上游重构任务（2.8）有意图地改变了行为，这些测试就从"保护层"变成了"阻碍层"
4. **第 3 层**：改变行为的重构任务和更新定标测试的修复任务是**紧密耦合的**——不更新测试，重构无法通过质量门；不先理解重构意图，又无从更新测试。执行器在这两个依赖之间反复横跳

## Solution

**分派层面解决**：改变行为的重构任务（breaking=true）应当在其任务描述中**显式声明哪些定标测试需要更新，以及新行为的预期**，而不是让后续的 fix 任务去反向推断。

例如，task 2.8 的描述应包含：
> 注：此重构会修改 `--block-source` 的行为——完成后，`TestAddCmd_BlockSource` 中的 `source 1.1 should be blocked` 断言应更新为 `source 1.1 is NOT blocked under new behavior`。

## Reusable Pattern

**定标测试的双刃剑**：定标测试记录当前行为以捕获非预期变化，但在意图性重构中会成为卡点。任何改变行为的重构任务必须同时承担更新相关定标测试的责任，不能让后续 fix 任务去推断"这到底是不是预期变化"。

## Example

fix.7 的任务描述：
```
Update characterization tests after SourceTaskID sentinel elimination
```

这四个测试的具体失败信息：
```
TestAddCmd_BlockSource:        source 1.1 should be blocked, got ""
TestAdd_BlockSource:           source 1.1 should have fix-task as dependency, got []
```

执行器无法判断：是 2.8 中去掉了 `--block-source` 的 blocked 效果，还是 2.8 改了实现方式导致 bug。这需要 2.8 的任务创建者说明意图。

## References

- task 0.1: 定标测试任务
- task 2.8: SourceTaskID sentinel 消除
- task fix.7: 更新定标测试（被卡）
- `docs/features/forge-architecture-simplification/tasks/records/2.8-quality-gate-fixes.md`