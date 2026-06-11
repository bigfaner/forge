---
feature: "per-task-surface-scoped-gate"
journey: "non-testable-task-skip"
risk_level: "Low"
surface_types: ["cli"]
surface_keys: ["cli"]
sources:
  - docs/proposals/per-task-surface-scoped-gate/proposal.md
generated: "2026-06-08"
---

# Journey: non-testable-task-skip

**Risk Level**: Low

<!-- Risk Classification Criteria:
  High   = Workflow involves state mutation, data loss risk, or irreversible operations
  Medium = Workflow involves multi-step interaction without irreversible side effects
  Low    = Workflow is read-only or purely observational
-->

## Overview

非 coding 类型任务（doc/gate/summary/eval）提交时，`IsTestableType()` 返回 false，quality gate 被跳过，本改动不影响此行为。

<!-- One-sentence description of the user workflow and its goal -->

## Setup

<!-- Preconditions that must be established before the Journey starts.
     These are environment states, not user actions. -->

- 项目可以是任意 surface 配置（单 surface 或多 surface）
- 任务类型为 `doc`、`gate`、`summary` 或 `eval` 等非 coding 类型
- 任务已完成执行，准备提交

## Happy Path

<!-- The primary success scenario: steps the user takes to accomplish the goal.
     Each step describes a user action and the expected outcome.
     High-risk Journeys MUST have edge case count >= happy path step count. -->

### Step 1: Submit a doc-type task

**User Action**: 执行 `forge task submit` 提交一个类型为 `doc` 的已完成任务

**Expected Result**: `IsTestableType()` 对 `doc` 类型返回 false，`validateQualityGate()` 跳过 gate 执行，`RunGate()` 不被调用，任务直接提交成功

### Step 2: Submit a gate-type task

**User Action**: 执行 `forge task submit` 提交一个类型为 `gate` 的已完成任务

**Expected Result**: `IsTestableType()` 对 `gate` 类型返回 false，gate 执行被跳过，任务直接提交成功

## Edge Cases

<!-- Alternative scenarios where things go wrong or take an unexpected path.
     Each edge case references a happy path step (variant) and describes the
     divergent precondition and expected outcome.
     High-risk Journeys: number of edge cases MUST be >= number of happy path steps. -->

### Step 1b: Task type changes from doc to coding after reclassification

**Precondition**: 任务最初创建为 `doc` 类型，但被重新分类为 `coding` 类型

<!-- The precondition that differs from the happy path, causing this outcome -->

**User Action**: 提交该重新分类后的任务

**Expected Result**: `IsTestableType()` 对 `coding` 类型返回 true，gate 正常执行。如果任务有 `surface-key`，则执行 scoped gate；否则执行全量 gate

### Step 2b: Verify RunGate() is never invoked for non-testable types

**Precondition**: 任务类型为非 coding 类型

**User Action**: 提交任务并检查 gate 日志

**Expected Result**: 日志中无任何 `RunGate()` 调用记录，`resolvePrefixedRecipe()` 也不被调用。改动前的行为与改动后完全一致

## Journey Invariants

<!-- Cross-step properties that must hold throughout the entire Journey.
     At least one invariant is required per Journey.
     These are verified across all steps, not within a single step. -->

- 非 coding 类型任务的 gate 跳过行为在改动前后完全一致——本改动不应影响 `IsTestableType()` 的判断逻辑
- `RunGate()` 和 `resolvePrefixedRecipe()` 在非 testable 类型提交路径中根本不被调用
