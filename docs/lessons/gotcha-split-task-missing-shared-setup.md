---
created: "2026-05-13"
tags: [testing, architecture]
---

# 拆分含全局 Setup 的任务导致 Subagent 跳过共享基础设施

## Problem

将 `/gen-test-scripts`（T-test-2）按输出大小拆分为 3 个并行 sub-task（fix-2 / disc-1 / disc-2）后，所有 3 个 subagent 都在 `beforeEach` 中调用 `login(page)` 重新登录，而非使用 Playwright 的 `storageState` 复用认证。

结果：
- 70 个 test × 每次登录 ~3-5s = 大量时间浪费在重复登录上
- 缺少 `auth-setup.ts` 和 `playwright.config.ts` 的 `projects` 配置
- 明确违反 skill 的 HARD-RULE："Auth-required tests MUST NOT call loginViaUI in beforeEach"

## Root Cause

**症状**：每个 spec 文件都 `beforeEach → login()`，无 `storageState` 复用

**直接原因**：3 个 subagent 都跳过了 gen-test-scripts 的 Step 1（Auth Classification），直接进入 spec 生成，默认用最简单的 `login()` helper

**根本原因**：任务拆分时只按"输出文件大小"拆分，没有识别出 gen-test-scripts 有一个**全局 setup phase**（auth classification → 生成 `auth-setup.ts` → 配置 `playwright.config.ts`）。每个 subagent 只看到自己的 TC 子集，不知道需要全局 auth 基础设施，导致：
1. 没有任何 sub-task 负责 `auth-setup.ts` 的生成
2. 没有 sub-task 修改 `playwright.config.ts` 添加 `projects` + `storageState`
3. 每个 subagent 独立决定用 `login()` 作为 fallback

## Solution

1. 生成 `auth-setup.ts`（从 skill template）
2. 在 `playwright.config.ts` 配置 `projects`：setup project 运行 `auth-setup.ts`，authenticated project 用 `storageState`
3. 移除所有 spec 中的 `beforeEach → login()`，改用 Playwright 自动注入的认证状态

## Reusable Pattern

**当拆分一个包含全局 setup phase 的任务时，必须：**

1. **识别全局 setup** — 分析 skill/任务的步骤，找出哪些是"一次性全局配置"（如 auth setup、DB migration、shared helpers）
2. **拆出 pre-task** — 将全局 setup 作为独立的 pre-task，所有并行子任务依赖它
3. **在子任务描述中注明** — "Auth infrastructure 已由 pre-task 完成，使用 storageState 而非 login()"

**判断标准**：如果 skill 的某个 step 产生的是**所有后续步骤共享的产物**（而非当前 step 专用的中间结果），它就是全局 setup，拆分时必须保留为独立任务。

**反模式**：按输出大小拆分任务时，只看最终产物大小，忽略中间共享产物的依赖关系。

## Example

```
# 错误拆分
fix-2: 生成 api.spec.ts (TC-053~TC-070)
disc-1: 生成 milestones-page.spec.ts (TC-001~TC-030)  ← 各自独立，跳过全局 auth setup
disc-2: 生成 existing-pages.spec.ts (TC-031~TC-052)

# 正确拆分
pre-e2e-auth: 生成 auth-setup.ts + 配置 playwright.config.ts  ← 全局 setup
fix-2: 生成 api.spec.ts (依赖 pre-e2e-auth)
disc-1: 生成 milestones-page.spec.ts (依赖 pre-e2e-auth)
disc-2: 生成 existing-pages.spec.ts (依赖 pre-e2e-auth)
```

## Related Files

- `tests/e2e/playwright.config.ts`
- `tests/e2e/auth-setup.ts` (缺失)
- `tests/e2e/features/milestone-map/milestones-page.spec.ts`
- `tests/e2e/features/milestone-map/existing-pages.spec.ts`
- `docs/lessons/gotcha-large-output-stall-subagent.md` (本次拆分的起因)
