---
feature: "surface-key-test-alignment"
journey: "single-surface-gen"
risk_level: "Medium"
surface_types: ["cli"]
sources:
  - docs/proposals/surface-key-test-alignment/proposal.md
generated: "2026-06-06"
---

# Journey: single-surface-gen

**Risk Level**: Medium

<!-- Risk Classification Criteria:
  High   = Workflow involves state mutation, data loss risk, or irreversible operations
  Medium = Workflow involves multi-step interaction without irreversible side effects
  Low    = Workflow is read-only or purely observational
-->

## Overview

验证单 surface 项目（scalar 和 named 两种形式）中，`forge` pipeline 生成的 gen-test-scripts 任务文件命名无 surface-key 后缀，测试文件输出到 `tests/<journey>/` 目录而非 `tests/<surfaceKey>/<journey>/`。

## Setup

- Forge 项目配置为单 surface（scalar 形式：`surfaces: tui`，或 named 形式：`surfaces: [{key: app, type: tui}]`）
- 项目中已有测试 Journey 文档（如 `docs/features/<slug>/testing/<journey>/journey.md`）
- 项目中有对应的 Contract 文件

## Happy Path

### Step 1: Scalar 单 surface 项目生成任务文件

**User Action**: 用户在 scalar 单 surface 项目（`surfaces: tui`）中运行 `forge task index`（或 pipeline 自动触发生成任务）

**Expected Result**: 生成的 gen-test-scripts 任务文件名为 `gen-test-scripts.md`（无 surface-key 后缀，无 surface-type 后缀）

### Step 2: Scalar 单 surface 项目测试输出目录

**User Action**: 在 scalar 单 surface 项目中，执行 gen-test-scripts 生成的测试脚本任务

**Expected Result**: 测试文件输出到 `tests/<journey>/` 目录（如 `tests/task-lifecycle/step1_claim_task.spec.ts`），不含 surface-key 层级目录

### Step 3: Named 单 surface 项目生成任务文件

**User Action**: 用户在 named 单 surface 项目（`surfaces: [{key: app, type: tui}]`）中运行 pipeline 生成任务

**Expected Result**: 生成的 gen-test-scripts 任务文件名同样为 `gen-test-scripts.md`（单 surface 无论 scalar 还是 named，都不加后缀）

### Step 4: Named 单 surface 项目测试输出目录

**User Action**: 在 named 单 surface 项目中，执行 gen-test-scripts 生成的测试脚本任务

**Expected Result**: 测试文件输出到 `tests/<journey>/` 目录，与 scalar 形式行为一致，不含 key 层级目录

## Edge Cases

### Step 1b: 从多 surface 降级为单 surface 后的任务命名

**Precondition**: 项目原来配置了多 surface（`surfaces: [{key: backend, type: api}, {key: frontend, type: web}]`），后来改为单 surface（`surfaces: tui`）

**User Action**: 修改 `.forge/config.yaml` 为单 surface 后重新运行 pipeline 生成任务

**Expected Result**: 新生成的任务文件名从 `gen-test-scripts-backend.md` 变为 `gen-test-scripts.md`，旧的多 surface 任务文件应被清理或不再生效

### Step 2b: 单 surface 项目中 gen-test-scripts 任务引用了不存在的 surface-key

**Precondition**: gen-test-scripts SKILL.md 中的输出目录规则错误地引用了 surface-key 路径模板

**User Action**: 在单 surface 项目中执行 gen-test-scripts 任务

**Expected Result**: 不应出现 surface-key 层级目录；如果 skill 指令中仍引用旧路径，应生成到正确的 `tests/<journey>/` 而非 `tests/<surfaceKey>/<journey>/`

### Step 3b: 单 surface 项目 pipeline.go 的 expansion 模式验证

**Precondition**: pipeline.go 中 gen-test-scripts 的 expansion 模式已从 `per-surface-type` 改为 `per-surface-key`

**User Action**: 检查 `expandPerSurfaceKey` 的 `isSingleSurface` 分支在单 surface 下的行为

**Expected Result**: 单 surface 项目触发 `isSingleSurface` 去后缀逻辑，任务名不含 surface-key 后缀

### Step 4b: 向后兼容性 -- 现有单 surface 项目升级 Forge 后行为不变

**Precondition**: 现有单 surface 项目已有 `tests/<journey>/` 目录下的测试文件，升级 Forge 到新版本

**User Action**: 升级 Forge plugin 后重新运行 gen-test-scripts

**Expected Result**: 测试文件仍然输出到 `tests/<journey>/`，不新增 surface-key 层级；已有测试文件路径不被破坏

## Journey Invariants

- 单 surface 项目（无论 scalar 还是 named 形式）的任务文件名始终为 `gen-test-scripts.md`（无后缀）
- 单 surface 项目的测试输出目录始终为 `tests/<journey>/`，不含 surface-key 层级
- pipeline.go 中 `expandPerSurfaceKey` 的 `isSingleSurface` 分支是单 surface 去后缀的唯一机制
