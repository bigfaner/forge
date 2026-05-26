---
created: "2026-05-26"
tags: [testing, architecture]
---

# Fix Task Scope Must Be Single Test Suite

## Problem
fix-6 覆盖 5 个测试套件的 config fixture 修复，task-executor 需要逐个诊断和修复，执行时间过长。用户等不到结果就中断了。

## Root Cause
1. fix-3 拆分策略是按"问题类型"而非按"文件/suite"分组——把所有"config fixture 问题"归为一个 fix task
2. 不同套件的 fixture 模式不同：有的完全没 config.yaml，有的有 config 但缺 surfaces，有的断言基于旧架构
3. task-executor 对 5 个套件逐一处理，每个套件需要诊断→修复→验证循环，累积耗时远超单个 suite

## Solution
Fix task 拆分粒度：一个 fix task = 一个测试套件（一个 `tests/<suite>/` 目录）。每个 fix task 的 scope 限定为该目录下的所有文件。

## Reusable Pattern
创建 fix task 时：
- SOURCE_FILES 限定在单一测试套件目录
- TEST_SCRIPT 只运行该套件的测试
- 如果多个套件有相同问题，创建多个独立 fix task 并行 dispatch，而非合并为一个
