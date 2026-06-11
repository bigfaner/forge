---
feature: "per-task-surface-scoped-gate"
journey: "scoped-gate-execution"
risk_level: "High"
surface_types: ["cli"]
surface_keys: ["cli"]
sources:
  - docs/proposals/per-task-surface-scoped-gate/proposal.md
generated: "2026-06-08"
---

# Journey: scoped-gate-execution

**Risk Level**: High

<!-- Risk Classification Criteria:
  High   = Workflow involves state mutation, data loss risk, or irreversible operations
  Medium = Workflow involves multi-step interaction without irreversible side effects
  Low    = Workflow is read-only or purely observational
-->

## Overview

多 surface 项目中，带有 surface-key 的任务提交时，per-task quality gate 只验证该 surface-key 对应的 prefixed recipe，跳过其他 surface 的验证。

<!-- One-sentence description of the user workflow and its goal -->

## Setup

<!-- Preconditions that must be established before the Journey starts.
     These are environment states, not user actions. -->

- 项目配置为多 surface（如 `backend=api` + `frontend=web`），justfile 包含 prefixed recipe（如 `backend-compile`、`frontend-lint`）
- 任务已创建并完成执行，状态为可提交，任务 frontmatter 中包含 `surface-key: backend` 或 `surface-key: frontend`
- Forge 已初始化，`forge surfaces` 输出多个 surface key

## Happy Path

<!-- The primary success scenario: steps the user takes to accomplish the goal.
     Each step describes a user action and the expected outcome.
     High-risk Journeys MUST have edge case count >= happy path step count. -->

### Step 1: Submit backend task with surface-key

**User Action**: 执行 `forge task submit` 提交一个 `surface-key: backend` 的已完成任务

**Expected Result**: `validateQualityGate()` 调用 `RunGate()` 时传入 `scope="backend"`，`RunGate()` 依次探测并执行 `just backend-compile` → `just backend-fmt` → `just backend-lint` → `just backend-unit-test`，全部通过后任务提交成功

### Step 2: Submit frontend task with surface-key

**User Action**: 执行 `forge task submit` 提交一个 `surface-key: frontend` 的已完成任务

**Expected Result**: `validateQualityGate()` 调用 `RunGate()` 时传入 `scope="frontend"`，`RunGate()` 依次探测并执行 `just frontend-compile` → `just frontend-fmt` → `just frontend-lint` → `just frontend-unit-test`，全部通过后任务提交成功，不执行任何 backend recipe

### Step 3: Verify scoped gate skips other surface

**User Action**: 在 backend 任务提交过程中，检查 gate 执行日志确认没有调用任何 frontend recipe（如 `frontend-compile`、`frontend-lint`）

**Expected Result**: gate 日志中只有 `backend-compile`/`backend-fmt`/`backend-lint`/`backend-unit-test` 四步，无任何 frontend 相关 recipe 调用

### Step 4: Verify prefixed recipe resolution mechanism

**User Action**: 检查 `RunGate()` 内部 `resolvePrefixedRecipe()` 的行为——当 `scope="backend"` 时，对每个 recipe 名拼接前缀并调用 `just.HasRecipe()` 探测

**Expected Result**: 每个 recipe（compile/fmt/lint/unit-test）均成功解析为 `backend-<recipe>` 形式，`HasRecipe()` 返回 true，不回退到 generic recipe

## Edge Cases

<!-- Alternative scenarios where things go wrong or take an unexpected path.
     Each edge case references a happy path step (variant) and describes the
     divergent precondition and expected outcome.
     High-risk Journeys: number of edge cases MUST be >= number of happy path steps. -->

### Step 1b: Backend task with missing surface-key in task frontmatter

**Precondition**: 任务的 frontmatter 中未设置 `surface-key` 字段，或 `surface-key` 为空字符串

<!-- The precondition that differs from the happy path, causing this outcome -->

**User Action**: 执行 `forge task submit` 提交该任务

**Expected Result**: `RunGate()` 收到 `scope=""`，跳过 prefixed 解析分支，回退到通用 recipe（`just compile` → `just fmt` → `just lint` → `just unit-test`），执行全量验证

### Step 2b: Surface-key references a non-existent prefixed recipe

**Precondition**: 任务的 `surface-key: backend`，但 justfile 中不存在 `backend-compile` 等 prefixed recipe（如 justfile 只生成了通用 recipe）

**User Action**: 执行 `forge task submit` 提交该任务

**Expected Result**: `resolvePrefixedRecipe()` 探测 `backend-compile` 时 `HasRecipe()` 返回 false，回退到通用 `compile` recipe，gate 执行全量验证但不报错

### Step 3b: Prefixed recipe partially exists (some steps prefixed, others not)

**Precondition**: justfile 中存在 `backend-compile` 和 `backend-lint`，但不存在 `backend-fmt` 和 `backend-unit-test`

**User Action**: 执行 `forge task submit` 提交 `surface-key: backend` 的任务

**Expected Result**: `compile` 和 `lint` 步骤使用 prefixed recipe（`backend-compile`、`backend-lint`），`fmt` 和 `unit-test` 步骤回退到通用 recipe（`fmt`、`unit-test`），混合执行成功

### Step 4b: Multiple surface-keys share same recipe prefix pattern

**Precondition**: 项目配置了 `backend-api=api` 和 `backend-worker=api` 两个 surface，justfile 中有 `backend-api-compile` 和 `backend-worker-compile`

**User Action**: 分别提交 `surface-key: backend-api` 和 `surface-key: backend-worker` 的任务

**Expected Result**: 每个任务只调用自己 surface-key 对应的 prefixed recipe，互不干扰。`backend-api` 任务不调用 `backend-worker-compile`，反之亦然

### Step 5b: NormalizeSurfaceKey character set constraint

**Precondition**: surface-key 包含大写字母、下划线或空格等非 `[a-z][a-z0-9-]*` 字符

**User Action**: 提交 `surface-key` 值为非规范格式的任务

**Expected Result**: `NormalizeSurfaceKey()` 将其规范化为合法字符集后再拼接 recipe 前缀，prefixed recipe 名称在 justfile 中合法且无歧义

## Journey Invariants

<!-- Cross-step properties that must hold throughout the entire Journey.
     At least one invariant is required per Journey.
     These are verified across all steps, not within a single step. -->

- 当 `scope`（surfaceKey）非空时，`RunGate()` 的每一步都必须先尝试 prefixed recipe，仅在 prefixed recipe 不存在时才回退到 generic recipe
- 回退到 generic recipe 不应视为错误——它是向后兼容的设计行为
- `scope=""` 时 `resolvePrefixedRecipe()` 不执行任何 prefixed 探测，直接返回原始 recipe 名
- 所有 gate recipe 步骤（compile/fmt/lint/unit-test）使用一致的 resolution 策略——不会出现某些步骤用 prefixed 而其他步骤因为 bug 使用 generic（除非 prefixed 不存在）
