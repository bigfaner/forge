---
feature: "per-task-surface-scoped-gate"
journey: "generic-recipe-fallback"
risk_level: "Medium"
surface_types: ["cli"]
surface_keys: ["cli"]
sources:
  - docs/proposals/per-task-surface-scoped-gate/proposal.md
generated: "2026-06-08"
---

# Journey: generic-recipe-fallback

**Risk Level**: Medium

<!-- Risk Classification Criteria:
  High   = Workflow involves state mutation, data loss risk, or irreversible operations
  Medium = Workflow involves multi-step interaction without irreversible side effects
  Low    = Workflow is read-only or purely observational
-->

## Overview

单 surface 项目或无 prefixed recipe 的多 surface 项目中，per-task quality gate 回退到通用 recipe（`compile`/`fmt`/`lint`/`unit-test`），行为与改动前完全一致。

<!-- One-sentence description of the user workflow and its goal -->

## Setup

<!-- Preconditions that must be established before the Journey starts.
     These are environment states, not user actions. -->

- 项目配置为单 surface（如 `cli=cli`），justfile 只有通用 recipe（`compile`、`lint` 等），无 prefixed recipe
- 或者：多 surface 项目但 justfile 未生成 surface-specific recipe
- 任务已创建并完成执行，准备提交

## Happy Path

<!-- The primary success scenario: steps the user takes to accomplish the goal.
     Each step describes a user action and the expected outcome.
     High-risk Journeys MUST have edge case count >= happy path step count. -->

### Step 1: Submit task in single-surface project

**User Action**: 在单 surface 项目中执行 `forge task submit` 提交一个已完成任务

**Expected Result**: `RunGate()` 收到空 scope（单 surface 项目任务无 surface-key），跳过 prefixed 分支，执行通用 `just compile` → `just fmt` → `just lint` → `just unit-test`，退出码和 stdout 输出与改动前一致

### Step 2: Submit task in multi-surface project without prefixed recipes

**User Action**: 在一个多 surface 项目中提交任务，该项目的 justfile 只包含通用 recipe，不包含 `<key>-<recipe>` 形式的 prefixed recipe

**Expected Result**: `resolvePrefixedRecipe()` 对每个 recipe 探测 prefixed 版本均返回 false，回退到通用 recipe，gate 执行全量验证，行为与改动前一致

### Step 3: Verify feature-level gate remains unchanged

**User Action**: 执行 `forge quality-gate`（feature-level gate），传入空 scope

**Expected Result**: `RunGate()` 的 `scope=""` 导致跳过整个 prefixed 解析分支，执行全量 compile/fmt/lint/unit-test，行为与改动前完全一致

## Edge Cases

<!-- Alternative scenarios where things go wrong or take an unexpected path.
     Each edge case references a happy path step (variant) and describes the
     divergent precondition and expected outcome.
     High-risk Journeys: number of edge cases MUST be >= number of happy path steps. -->

### Step 1b: Single-surface project with task that has surface-key set

**Precondition**: 项目是单 surface，但任务的 frontmatter 中手动设置了 `surface-key: backend`

<!-- The precondition that differs from the happy path, causing this outcome -->

**User Action**: 提交该任务

**Expected Result**: `resolvePrefixedRecipe()` 尝试探测 `backend-compile` 等 prefixed recipe，justfile 中不存在，回退到通用 `compile`，最终执行全量验证。不会因 surface-key 存在就报错

### Step 2b: Mixed justfile with partial generic recipe coverage

**Precondition**: justfile 中通用 `compile` recipe 存在，但 `unit-test` 通用 recipe 不存在（项目只定义了 `backend-unit-test` prefixed recipe）

**User Action**: 在该项目中提交无 surface-key 的任务

**Expected Result**: `compile`、`fmt`、`lint` 步骤通过通用 recipe 执行；`unit-test` 步骤因通用 recipe 不存在而失败，gate 报告失败。这反映了 justfile 配置不完整，不是 Forge 的 bug

### Step 3b: Verify no performance regression in generic path

**Precondition**: 单 surface 项目，justfile 只有通用 recipe

**User Action**: 提交多个任务，测量 gate 执行时间

**Expected Result**: 由于 `scope=""` 直接跳过 prefixed 探测（0 次 `HasRecipe()` 调用），gate 执行时间与改动前无显著差异（探测开销为 0）

## Journey Invariants

<!-- Cross-step properties that must hold throughout the entire Journey.
     At least one invariant is required per Journey.
     These are verified across all steps, not within a single step. -->

- 回退到 generic recipe 时，recipe 名称必须与改动前完全一致（`compile`/`fmt`/`lint`/`unit-test`），不能因回退逻辑引入名称变更
- feature-level gate（`scope=""`）路径下，`resolvePrefixedRecipe()` 不执行任何探测操作，性能零开销
- 通用 recipe 回退不应产生额外的日志或警告——这是正常行为，不是降级
