# fix-e2e Task "unknown" Placeholder Issue

## Problem

When `task all-completed` auto-generates a fix-e2e task after detecting e2e test failures, the task file contains vague "unknown" placeholders:

- Title: `"修复 e2e 测试失败: unknown"`
- Reference files: `failure-unknown.md` (doesn't exist)
- Descriptions contain no actionable failure details

This makes the task nearly useless for the error-fixer subagent, which needs concrete failure information to diagnose and fix issues.

## Root Cause

**Causal chain (3 levels):**

1. **Symptom** → Task generated with "unknown" placeholders
2. **Direct cause** → `task all-completed` couldn't parse specific failure details from test output
3. **Root cause** → e2e tests produced **zero discovered tests** (infrastructure-level failure, not individual test failures). The test runner exited non-zero, but there were no parseable test results — only a `latest.md` showing "0 tests, 0 passed, 0 failed". The fix-e2e generator has no fallback for this scenario and defaults to "unknown".

The generator assumes failures are at the individual-test level (where it can extract test names and error messages). It doesn't handle infrastructure failures where no tests run at all (missing dependencies, script syntax errors, environment misconfiguration, etc.).

## Solution

Two-part fix:

1. **Improve fix-e2e task generation**: When the test runner produces no parseable failures, the generator should:
   - Capture the raw stderr/stdout from the test runner as the failure detail
   - Use a more descriptive title like "e2e infrastructure failure" instead of "unknown"
   - Reference the raw output file instead of a non-existent `failure-unknown.md`

2. **Manual workaround**: When encountering a fix-e2e task with "unknown", first run the e2e tests directly to see the actual error output, then fix based on real diagnostics rather than the placeholder task content.

## Key Takeaway

Auto-generated fix tasks must handle **infrastructure-level failures** (0 tests discovered) differently from **test-level failures** (specific assertions failed). The generator should always capture and include raw test runner output as context, even when it can't parse structured failure data.


## fix-e2e-1-1 原文
```
# fix-e2e-1-1: 修复 e2e 测试失败

## Description

这是第 1 轮修复尝试。修复步骤：

1. 读取 `testing/results/latest.md` 查看失败概览
2. 读取 `testing/results/failures/failure-unknown.md` 了解具体失败详情
3. 定位根本原因（代码逻辑 / 测试脚本 / 环境配置）
4. 修复并验证

## Reference Files

- `testing/results/latest.md` — 测试结果概览
- `testing/results/failures/failure-unknown.md` — unknown

- `testing/test-cases.md` — 测试用例文档
- `testing/scripts/` — 测试脚本目录

## Acceptance Criteria

- [ ] 已定位失败的根本原因
- [ ] 已修复代码或测试脚本
- [ ] 单元测试全部通过
```