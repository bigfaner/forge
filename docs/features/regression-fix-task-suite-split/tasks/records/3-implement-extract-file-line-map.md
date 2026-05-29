---
status: "completed"
started: "2026-05-29 14:16"
completed: "2026-05-29 14:26"
time_spent: "~10m"
---

# Task Record: 3 实现 extractFileLineMap 函数

## Summary
实现 extractFileLineMap 和 isTestFile 函数。两步扫描：收集 --- FAIL: 块主文件集合，按行匹配并扩展 ±2 上下文窗口。8 个单元测试全部通过。

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/internal/cmd/quality_gate_test.go

### Key Decisions
- 主文件识别：每个 --- FAIL: 块的首个 test file ref 为主文件
- 栈 trace 引用文件不生成独立条目，其输出行通过行匹配归入主文件上下文

## Test Results
- **Tests Executed**: Yes
- **Passed**: 8
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] 函数签名正确
- [x] 收集主文件集合
- [x] 上下文窗口 ±2 行
- [x] 多文件归入所有匹配主文件
- [x] 栈 trace 引用不生成独立条目
- [x] 单元测试覆盖

## Notes
无
