---
feature: "surface-key-test-alignment"
journey: "same-type-multi-surface-gen"
risk_level: "High"
surface_types: ["cli"]
sources:
  - docs/proposals/surface-key-test-alignment/proposal.md
generated: "2026-06-06"
---

# Journey: same-type-multi-surface-gen

**Risk Level**: High

<!-- Risk Classification Criteria:
  High   = Workflow involves state mutation, data loss risk, or irreversible operations
  Medium = Workflow involves multi-step interaction without irreversible side effects
  Low    = Workflow is read-only or purely observational
-->

## Overview

验证多个 surface 共享同一 type（如两个 `api` surface：`admin` 和 `public`）时，gen-test-scripts 按 surface-key 各自生成独立任务和独立输出目录，确保 key 的唯一性在任务命名和目录结构中得到体现。

## Setup

- Forge 项目配置为多 surface 同 type（如 `surfaces: [{key: admin, type: api}, {key: public, type: api}]`）
- 项目中已有测试 Journey 文档和对应的 Contract 文件
- pipeline.go 中 gen-test-scripts 使用 `per-surface-key` expansion

## Happy Path

### Step 1: 同 type 多 surface 生成独立任务文件

**User Action**: 在同 type 多 surface 项目中运行 pipeline 生成任务

**Expected Result**: 为每个 surface-key 各生成一个独立任务文件：`gen-test-scripts-admin.md` 和 `gen-test-scripts-public.md`，即使两个 surface 的 type 都是 `api`

### Step 2: admin surface 测试文件输出到独立目录

**User Action**: 执行 `gen-test-scripts-admin.md` 任务

**Expected Result**: 测试文件输出到 `tests/admin/<journey>/` 目录

### Step 3: public surface 测试文件输出到独立目录

**User Action**: 执行 `gen-test-scripts-public.md` 任务

**Expected Result**: 测试文件输出到 `tests/public/<journey>/` 目录，与 admin 的测试文件完全隔离

### Step 4: 两个 surface 的测试互不干扰

**User Action**: 同时存在 `tests/admin/<journey>/` 和 `tests/public/<journey>/` 下的测试文件，运行 `forge run-tests`

**Expected Result**: 每个 surface 的测试独立运行，admin 的测试不会发现 public 的测试文件，反之亦然

## Edge Cases

### Step 1b: 同 type 多 surface 的 justfile recipe 独立性

**Precondition**: 两个 surface key 不同但 type 相同，init-justfile 需要为每个 key 生成独立 recipe

**User Action**: 在同 type 多 surface 项目中生成 justfile

**Expected Result**: 生成 `admin-test` 和 `public-test` 两个独立 recipe（使用 key 而非 type），各自引用 `tests/admin/` 和 `tests/public/` 路径

### Step 2b: 同 type 多 surface 的 convention 文件不受影响

**Precondition**: convention 文件按 type 组织（如 `docs/conventions/testing/api/core.md`）

**User Action**: 检查 convention 文件路径

**Expected Result**: convention 文件路径不变，仍按 type 组织。surface-key 只影响任务命名和测试输出目录，不影响 convention 路径

### Step 3b: 同 type 多 surface 的 tag 命名不变

**Precondition**: 测试 tag 使用 type 标注（如 `@api-functional`）

**User Action**: 检查生成的测试文件中的 tag

**Expected Result**: tag 仍为 `@api-functional`（按 type），不因使用 surface-key 命名而改变 tag 格式

### Step 4b: pipeline expansion 不混淆同 type 的 surface

**Precondition**: `expandPerSurfaceKey` 在处理多个同 type surface 时

**User Action**: 验证 pipeline.go 的 expansion 逻辑在两个 `api` type surface 下的行为

**Expected Result**: expansion 按 key 而非 type 区分，生成 `admin` 和 `public` 两个独立任务，不会因 type 相同而合并或覆盖

### Step 5b: surface rule 文件加载不重复

**Precondition**: 两个 surface 都是 `api` type，对应的 surface rule 文件是同一个 `rules/surface-api.md`

**User Action**: pipeline 或 skill 加载 surface rule

**Expected Result**: 只加载一次 `rules/surface-api.md`，不因两个 surface 同 type 而重复加载

## Journey Invariants

- 同 type 多 surface 项目中，任务命名始终使用 surface-key（`admin`/`public`），不使用 type（`api`）
- 同 type 多 surface 项目的测试目录通过 key 分区（`tests/admin/` vs `tests/public/`），type 的相同性不影响目录隔离
- convention 路径、surface rule 文件、test tag 等按 type 组织的部分不受 surface-key 对齐的影响
- pipeline expansion 的去重基于 key，不基于 type
