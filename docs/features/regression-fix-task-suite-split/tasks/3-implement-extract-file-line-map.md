---
id: "3"
title: "实现 extractFileLineMap 函数"
priority: "P0"
estimated_time: "2h"
complexity: "high"
dependencies: []
surface-key: "."
surface-type: "cli"
breaking: false
type: "coding.feature"
mainSession: false
---

# 3: 实现 extractFileLineMap 函数

## Description

新建 `extractFileLineMap` 函数，从 regression test output 中提取文件路径到相关输出行的映射。与现有 `extractSourceFiles`（返回扁平逗号分隔字符串）不同，新函数保留文件-输出行对应关系，用于按测试文件分组创建 fix task。

函数内部完成匹配行提取、上下文扩展（±2 行）和重叠去重，返回值已是可直接写入 task description 的内容。

## Reference Files
- `forge-cli/internal/cmd/quality_gate.go`: sourceFileRe（L574）和 extractSourceFiles（L586-608）作为基线参考 (source: proposal.md#In-Scope)
- `docs/lessons/gotcha-fix-task-broad-scope.md`: lesson 中记录的 agent 卡死根因分析 (source: proposal.md#Evidence)

## Acceptance Criteria
- [ ] 函数签名：`func extractFileLineMap(output string) map[string][]string`，输入原始 output，返回文件路径到提取行的映射
- [ ] 第一步扫描 output 收集所有包含直接 `--- FAIL:` 条目的测试文件路径作为"主文件"集合（仅这些文件生成独立映射条目）
- [ ] 第二步遍历 output 每行，对"主文件"集合中的文件路径匹配行（使用 sourceFileRe 提取 `file_test.go:line`），取前后各 2 行作为上下文窗口（共 5 行），同一文件多处匹配合并重叠的上下文窗口（避免重复行）
- [ ] 一行匹配多个"主文件"时，该行及其上下文归入所有匹配的主文件
- [ ] 仅作为栈 trace 引用出现的测试文件（无直接 `--- FAIL:` 条目）不生成独立映射条目，其输出行通过行匹配自然归入引用它的主文件上下文
- [ ] 单元测试覆盖：多文件多失败、单文件失败、无 FAIL 行返回空 map、重叠上下文窗口去重、行匹配多文件归入所有匹配主文件

## Implementation Notes
- MVP 仅实现 Go 的提取模式（`sourceFileRe` 提取 `file_test.go:line`，叠加 `--- FAIL:` 块缩进行解析）
- 非测试文件路径的输出行（如 `handler.go:42`）不归入任何映射条目，由调用方（Task 4）处理 fallback
- 性能：两次线性扫描，O(L × F)，对 10000 行 output 的典型场景 < 100ms
- sub-test（`--- FAIL: TestFoo/SubTest`）的嵌套缩进按父 test 归属处理；parallel test 交错输出不在 MVP 范围内
