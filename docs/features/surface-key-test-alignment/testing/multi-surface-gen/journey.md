---
feature: "surface-key-test-alignment"
journey: "multi-surface-gen"
risk_level: "High"
surface_types: ["cli"]
sources:
  - docs/proposals/surface-key-test-alignment/proposal.md
generated: "2026-06-06"
---

# Journey: multi-surface-gen

**Risk Level**: High

<!-- Risk Classification Criteria:
  High   = Workflow involves state mutation, data loss risk, or irreversible operations
  Medium = Workflow involves multi-step interaction without irreversible side effects
  Low    = Workflow is read-only or purely observational
-->

## Overview

验证多 surface 项目（如 `backend=api`, `frontend=web`）中，gen-test-scripts 任务按 surface-key 命名、测试文件输出到 `tests/<surfaceKey>/<journey>/` 目录结构，解决原有 per-surface-type 导致的目录不一致问题。

## Setup

- Forge 项目配置为多 surface（如 `surfaces: [{key: backend, type: api}, {key: frontend, type: web}]`）
- 项目中已有测试 Journey 文档（如 `docs/features/<slug>/testing/<journey>/journey.md`）
- 项目中有对应的 Contract 文件
- pipeline.go 中 gen-test-scripts 的 expansion 模式已从 `per-surface-type` 改为 `per-surface-key`

## Happy Path

### Step 1: 多 surface 项目生成任务文件命名

**User Action**: 用户在多 surface 项目中运行 pipeline 生成任务

**Expected Result**: 为每个 surface-key 各生成一个任务文件：`gen-test-scripts-backend.md` 和 `gen-test-scripts-frontend.md`（使用 surface-key `backend`/`frontend` 而非 surface-type `api`/`web`）

### Step 2: backend surface 的测试文件输出目录

**User Action**: 执行 `gen-test-scripts-backend.md` 任务，为 backend (api) surface 生成测试脚本

**Expected Result**: 测试文件输出到 `tests/backend/<journey>/` 目录（如 `tests/backend/item-deletion/step1_delete_main_item.spec.ts`）

### Step 3: frontend surface 的测试文件输出目录

**User Action**: 执行 `gen-test-scripts-frontend.md` 任务，为 frontend (web) surface 生成测试脚本

**Expected Result**: 测试文件输出到 `tests/frontend/<journey>/` 目录（如 `tests/frontend/item-deletion/step1_delete_item.spec.ts`）

### Step 4: run-tests 兼容性验证

**User Action**: 运行 `forge run-tests` 或对应的 justfile recipe（如 `just backend-test <journey>`）

**Expected Result**: run-tests 能正确发现 `tests/backend/<journey>/` 和 `tests/frontend/<journey>/` 下的测试文件并执行；justfile recipe 中的路径与新目录结构一致

### Step 5: pipeline.go 测试通过

**User Action**: 运行 `go test ./pkg/task/...` 验证 pipeline 的 expansion 逻辑

**Expected Result**: 所有测试通过，gen-test-scripts 使用 `per-surface-key` expansion 的行为符合预期

## Edge Cases

### Step 1b: 多 surface 项目中 surface-key 包含特殊字符

**Precondition**: surface 配置中 key 值包含连字符或下划线（如 `key: admin-api`）

**User Action**: 运行 pipeline 生成任务

**Expected Result**: 任务文件名使用完整 key 值（如 `gen-test-scripts-admin-api.md`），输出目录为 `tests/admin-api/<journey>/`，不截断或转换 key

### Step 2b: 从单 surface 升级为多 surface 后的目录迁移

**Precondition**: 项目原为单 surface，测试文件在 `tests/<journey>/`；后新增第二个 surface 变为多 surface 配置

**User Action**: 修改 `.forge/config.yaml` 添加第二个 surface 后重新运行 pipeline

**Expected Result**: 新生成的任务按 surface-key 命名；已有测试文件需要迁移到 `tests/<surfaceKey>/<journey>/`；gen-test-scripts 任务指令中应体现新目录结构

### Step 3b: gen-test-scripts SKILL.md 指令中的路径模板正确性

**Precondition**: gen-test-scripts 的 SKILL.md 中引用了输出目录规则

**User Action**: 在多 surface 项目中读取 gen-test-scripts 任务文件，检查其输出目录指令

**Expected Result**: 任务文件中的路径模板为 `tests/<surfaceKey>/<journey>/`，而非旧的 `tests/<journey>/`；模板变量使用 surface-key 而非 surface-type

### Step 4b: justfile recipe 路径适配

**Precondition**: init-justfile 模板中 just recipe 引用了测试目录路径

**User Action**: 在多 surface 项目中运行 `forge init-justfile` 或检查生成的 justfile

**Expected Result**: just recipe 中的测试路径为 `tests/<surfaceKey>/<journey>/` 或通过参数传递正确组装，与 gen-test-scripts 的输出目录匹配

### Step 5b: init-justfile 模板中多 surface 的 recipe 前缀

**Precondition**: 多 surface 项目的 justfile recipe 使用 surface-key 作为前缀（如 `backend-test`）

**User Action**: 检查 init-justfile 模板生成的 recipe

**Expected Result**: recipe 名为 `backend-test` 和 `frontend-test`（使用 key 而非 type），recipe 内部的 `cd tests/<surfaceKey>/<journey>` 路径正确

### Step 6b: 全量 grep 确认无遗漏引用

**Precondition**: 所有 skill 文件、模板、规则文件已更新

**User Action**: 对 Forge plugin 目录执行 `grep -r "tests/<journey>"` 和 `grep -r "tests/{{journey}}"`

**Expected Result**: 所有剩余的 `tests/<journey>/` 引用仅在单 surface 上下文中出现，多 surface 上下文中均已更新为 `tests/<surfaceKey>/<journey>/` 或自适应规则

## Journey Invariants

- 多 surface 项目中，每个 surface 的任务文件名后缀必须是 surface-key（如 `backend`），不能是 surface-type（如 `api`）
- 多 surface 项目中，测试输出目录结构为 `tests/<surfaceKey>/<journey>/`，surface-key 层级始终存在
- run-tests 的 surface-key expansion 和 gen-test-scripts 的 surface-key expansion 必须使用同一个 key 值
- 所有 skill 文件中引用测试目录的描述必须与 pipeline.go 的 expansion 行为一致
