---
feature: "per-task-surface-scoped-gate"
journey: "prefixed-recipe-failure"
risk_level: "High"
surface_types: ["cli"]
surface_keys: ["cli"]
sources:
  - docs/proposals/per-task-surface-scoped-gate/proposal.md
generated: "2026-06-08"
---

# Journey: prefixed-recipe-failure

**Risk Level**: High

<!-- Risk Classification Criteria:
  High   = Workflow involves state mutation, data loss risk, or irreversible operations
  Medium = Workflow involves multi-step interaction without irreversible side effects
  Low    = Workflow is read-only or purely observational
-->

## Overview

Prefixed recipe 存在但执行失败时，错误信息中包含 surface 上下文（step name 为 `<key>-<recipe>` 形式），可直接映射到对应 surface 的错误处理路径。

<!-- One-sentence description of the user workflow and its goal -->

## Setup

<!-- Preconditions that must be established before the Journey starts.
     These are environment states, not user actions. -->

- 多 surface 项目，justfile 包含 prefixed recipe（如 `backend-compile`、`backend-lint`）
- 任务的 `surface-key: backend`，代码中存在故意制造的编译错误或 lint 违规
- 任务已完成执行，准备提交

## Happy Path

<!-- The primary success scenario: steps the user takes to accomplish the goal.
     Each step describes a user action and the expected outcome.
     High-risk Journeys MUST have edge case count >= happy path step count. -->

### Step 1: Submit task and prefixed compile recipe fails

**User Action**: 执行 `forge task submit` 提交 `surface-key: backend` 的任务，backend 代码中存在编译错误

**Expected Result**: `just backend-compile` 执行失败，`onFail` 回调收到的 step name 为 `backend-compile`（含 surface 上下文），gate 报告失败，任务不被提交

### Step 2: Submit task and prefixed lint recipe fails

**User Action**: 执行 `forge task submit` 提交 `surface-key: backend` 的任务，backend 代码编译通过但存在 lint 违规

**Expected Result**: `backend-compile` 通过，`backend-fmt` 通过，`backend-lint` 执行失败，`onFail` 回调的 step name 为 `backend-lint`，错误信息可区分来自 backend surface 而非 frontend

### Step 3: Verify error message distinguishability between prefixed and generic failures

**User Action**: 对比 prefixed recipe 失败（`backend-lint`）和 generic recipe 失败（`lint`）的错误信息

**Expected Result**: prefixed 失败的 step name 为 `backend-lint`，generic 失败的 step name 为 `lint`，两者可通过 step name 区分错误来源的 surface

### Step 4: Fix error and resubmit successfully

**User Action**: 修复 backend 代码中的编译错误和 lint 违规，再次执行 `forge task submit`

**Expected Result**: `backend-compile`/`backend-fmt`/`backend-lint`/`backend-unit-test` 全部通过，任务提交成功

## Edge Cases

<!-- Alternative scenarios where things go wrong or take an unexpected path.
     Each edge case references a happy path step (variant) and describes the
     divergent precondition and expected outcome.
     High-risk Journeys: number of edge cases MUST be >= number of happy path steps. -->

### Step 1b: Prefixed recipe fails with non-zero exit code but no stderr output

**Precondition**: `just backend-compile` 退出码非零，但 stderr 为空（如 just recipe 本身吞掉了错误输出）

<!-- The precondition that differs from the happy path, causing this outcome -->

**User Action**: 提交任务

**Expected Result**: gate 仍然报告失败（基于退出码），`onFail` 回调的 step name 仍为 `backend-compile`，即使无 stderr 也能通过 step name 定位 surface

### Step 2b: First prefixed recipe passes but subsequent one fails

**Precondition**: `backend-compile` 通过，`backend-fmt` 通过，`backend-lint` 失败

**User Action**: 提交任务

**Expected Result**: gate 在 `backend-lint` 步骤停止，不执行 `backend-unit-test`，`onFail` 报告的失败 step 为 `backend-lint`。前面已通过的步骤不影响失败判定

### Step 3b: Generic recipe failure produces different step name format than prefixed

**Precondition**: 任务无 surface-key，回退到 generic `lint` recipe，`lint` 执行失败

**User Action**: 提交任务

**Expected Result**: `onFail` 的 step name 为 `lint`（无 surface 前缀），与 prefixed 失败的 `backend-lint` 在格式上可区分。错误处理路径可根据 step name 是否包含连字符前缀来区分 scoped vs 全量失败

### Step 4b: onFail callback receives correct surface context for error routing

**Precondition**: `backend-lint` 失败，`onFail` 回调被触发

**User Action**: 检查 `onFail` 回调的参数内容

**Expected Result**: step name 为 `backend-lint`，可从中提取 surface key（`backend`）和 recipe 名（`lint`），无需额外解析。错误处理路径可直接映射到对应 surface

### Step 5b: Failure in first gate step prevents subsequent step execution

**Precondition**: `backend-compile` 失败（编译错误）

**User Action**: 提交任务并检查 gate 日志

**Expected Result**: gate 只执行了 `backend-compile` 一步就停止，`backend-fmt`/`backend-lint`/`backend-unit-test` 均未执行。不存在部分执行后继续的情况——编译失败即终止

## Journey Invariants

<!-- Cross-step properties that must hold throughout the entire Journey.
     At least one invariant is required per Journey.
     These are verified across all steps, not within a single step. -->

- prefixed recipe 失败时，`onFail` 回调的 step name 始终为 `<key>-<recipe>` 形式（如 `backend-compile`），包含完整的 surface 上下文
- generic recipe 失败时，step name 始终为原始 recipe 名（如 `compile`），两者可区分
- gate 步骤按顺序执行，某步失败后立即停止，不会继续执行后续步骤
- 错误信息中 step name 的命名规则必须与 `resolvePrefixedRecipe()` 的输出一致——不会出现 step name 和实际执行的 recipe 名称不匹配的情况
