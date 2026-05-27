---
created: "2026-05-27"
tags: [architecture]
---

# Reference Files 指向 proposal 全文导致 task-executor 越界

## Problem

Task 3（quality_gate 路径替换）的 executor 越界修改了 task 4 的范围文件（infer.go、types.go、infer_test.go），清理了 `TypeTestVerifyRegression` 相关代码。Task 3 的 6 条 AC 本身是精确的，全部指向 quality_gate 相关变更。

## Root Cause

因果链（3 层）：

1. **表面现象**：executor 在执行 task 3 时主动 grep 了 `TypeTestVerifyRegression`、`T-quick-verify-regression` 等 task 4 范围的关键词，发现"代码还存在"后直接修复
2. **直接原因**：task 3 的 Reference Files 指向 `proposal.md#Layer-1-Go-代码层术语路径统一`，该 section 描述了 tasks 1-6 的所有工作。executor 读到 proposal 中"清理 TypeTestVerifyRegression"的需求后，在 Step 1.5 spec-code conflict scan 中将其判定为"spec 与 code 不一致"并主动修复
3. **根因**：prompt 模板存在指令矛盾——CODING_PRINCIPLES 说 "Surgical Changes: Modify only the code directly relevant to the task"，但 Step 1.5 `<CRITICAL>` 级别要求"scan existing code against spec requirements across five dimensions"并修复差异。CRITICAL 优先级高于一般原则，executor 遵循了冲突扫描结果而非 scope 边界。而 Reference Files 指向 proposal 全文使得 executor 看到了超出本任务 scope 的 spec 要求

## Solution

将 Reference Files 从"指向 proposal section"改为"内联精确定位信息"：

```markdown
# 当前（有越界风险）
## Reference Files
- `proposal.md#Layer-1-Go-代码层术语路径统一` — 第 4 项定义了 quality_gate 的具体修复点

# 改为（self-contained）
## Reference Files
- quality_gate.go 中 `tests/e2e/results/raw-output.txt` 路径需替换为 `GetTestResultsDir()`
- quality_gate.go 中 "promoted scripts in tests/e2e/" 注释需更新
- `runSurfaceLifecycle()` 需在 mobile surface type 时插入 mobile-test-setup 调用
```

这样 executor 不需要读 proposal，只看 task doc 内的精确定位信息。Step 1.5 的 spec-code scan 也只能扫描 task doc 内列出的范围。

取舍：quick-tasks 生成阶段需要从 proposal 提取并内联（更多工作），且 proposal 变更时 task doc 不会自动同步。但对于 quick mode（≤15 任务），这个代价可接受。

## Reusable Pattern

- **Reference Files 是 scope 边界的泄漏点**：当 Reference Files 指向包含多个任务 scope 的文档时，executor 会看到超出本任务的需求并越界修复。预防方法：Reference Files 只包含本任务直接相关的精确信息，不指向包含全局 scope 的文档
- **模板指令矛盾时 executor 遵循高优先级指令**：`<CRITICAL>` > CODING_PRINCIPLES。如果两个指令方向矛盾（"只改你的" vs "找出所有不一致"），executor 选择 CRITICAL 级别的，导致 scope 膨胀
- **Self-contained task doc > 指针型 Reference Files**：将 spec 内容内联到 task doc 比"指向外部文档"更安全，因为 executor 的视野被限制在 task doc 内

## Related Files

- `forge-cli/pkg/prompt/data/coding-enhancement.md` — Step 1.5 CRITICAL 级别的 spec-code scan 指令
- [[gotcha-prompt-template-complexity-agnostic]] — prompt 模板不区分复杂度导致过度探索
- [[gotcha-task-executor-invisible-thinking-time]] — 过度探索 → 更多 tool call → 更多不可见 thinking 时间
- [[gotcha-quick-tasks-merge-threshold]] — 任务拆分粒度问题
