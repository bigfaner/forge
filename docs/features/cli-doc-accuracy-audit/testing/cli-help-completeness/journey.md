---
feature: "cli-doc-accuracy-audit"
journey: "cli-help-completeness"
risk_level: "Low"
surface_types: ["cli"]
surface_keys: ["cli"]
sources:
  - docs/proposals/cli-doc-accuracy-audit/proposal.md
generated: "2026-06-07"
---

# Journey: cli-help-completeness

**Risk Level**: Low

<!-- Risk Classification Criteria:
  Low = 只读验证操作，运行 --help 查看输出，不涉及状态变更
-->

## Overview

开发者运行 `forge <command> --help` 验证每个 CLI 命令的 Long 描述完整反映该命令的全部功能和副作用，确保 help text 与 RunE 实现一致。

## Setup

<!-- Preconditions that must be established before the Journey starts. -->

- Forge CLI 已编译且可执行（`go build` 成功）
- 已知需要验证的 11 个 CLI 命令列表（来自提案 C1-C11）
- 可访问对应 RunE 函数的源代码用于交叉核对

## Happy Path

### Step 1: 验证 forge init 的 Long 描述

**User Action**: 运行 `forge init --help`，检查 Long 描述是否包含 surface detection 步骤说明。

**Expected Result**: Long 描述提及 surface detection 流程（自动检测项目 surface 类型如 cli/web/api 等）。输出中包含完整的初始化流程描述。

### Step 2: 验证 forge cleanup 的 Short/Long 描述

**User Action**: 运行 `forge cleanup --help`，检查描述是否包含 blocked、suspended、rejected 状态（而非仅 "completed"）。

**Expected Result**: Short 或 Long 描述明确列出 cleanup 覆盖的任务状态包括 completed、blocked、suspended 和 rejected。

### Step 3: 验证 forge task claim 的 Long 描述

**User Action**: 运行 `forge task claim --help`，检查 Long 描述是否提及 auto-unblock 行为。

**Expected Result**: Long 描述包含 claim 操作会自动解除被当前任务阻塞的任务（auto-unblock）这一行为说明。

### Step 4: 验证 forge task validate 的 Long 描述

**User Action**: 运行 `forge task validate --help`，检查 Long 描述中列出的验证项是否覆盖全部 12+ 项实际验证步骤。

**Expected Result**: Validations 列表完整覆盖所有实际执行的验证步骤（12+ 项），而非仅列出原来的 5 项。

### Step 5: 验证 forge task add 的 Long 描述

**User Action**: 运行 `forge task add --help`，检查 Long 描述是否包含使用概述。

**Expected Result**: Long 描述提供清晰的使用概述，说明如何添加不同类型的任务。

### Step 6: 验证 forge feature 的 Long 描述

**User Action**: 运行 `forge feature --help`，检查 Long 描述是否包含 `set` 子命令和行为差异说明。

**Expected Result**: Long 描述提及 `set` 子命令，并说明 feature 相关子命令之间的行为差异。

### Step 7: 验证 forge forensic search 和 forensic subagents 的 Long 描述

**User Action**: 分别运行 `forge forensic search --help` 和 `forge forensic subagents --help`。

**Expected Result**: 两个命令都有非空的 Long 描述，说明各自的功能、参数和输出。不再为空白或缺失。

### Step 8: 验证 forge worktree status 的 Long 描述

**User Action**: 运行 `forge worktree status --help`，检查 Long 描述是否提及 UNPUSHED 字段。

**Expected Result**: Long 描述提及输出包含 UNPUSHED 字段，用于指示未推送的提交。

### Step 9: 验证 forge fact summary 的 Long 描述

**User Action**: 运行 `forge fact summary --help`，检查 Long 描述是否提及 COVERAGE 指标。

**Expected Result**: Long 描述提及输出包含 COVERAGE 指标，说明其含义和计算方式。

### Step 10: 验证 forge quality-gate 的 Long 描述

**User Action**: 运行 `forge quality-gate --help`，检查 Long 描述是否涵盖 fix task 自动创建、retry-once 行为和 docs-only 跳过逻辑。

**Expected Result**: Long 描述完整描述：(1) 质量门禁失败时自动创建 fix task；(2) 只重试一次；(3) 跳过仅文档类型的任务。

## Edge Cases

### Step 1b: 未知子命令的 help 输出

**Precondition**: CLI 安装正常，但运行一个不存在的子命令

**User Action**: 运行 `forge nonexistent-command --help`

**Expected Result**: 返回 "unknown command" 错误，exit code != 0，stderr 包含可用命令列表提示。

### Step 2b: 未修改的命令 help 输出一致性

**Precondition**: 提案范围外的命令（不在 C1-C11 列表中）

**User Action**: 运行未修改命令（如 `forge task status --help`）确认其 help 输出未被意外影响。

**Expected Result**: 未修改命令的 help 输出保持不变，Long/Short 描述与修改前一致。

## Journey Invariants

- 每个 CLI 命令的 `--help` 输出都必须是非空的（至少有 Short 描述）
- `--help` 的 exit code 必须是 0（help 输出不是错误）
- 修改后的 Long 描述中提及的每一项功能都能在对应 RunE 函数中找到代码实现
- 所有修改仅涉及 cobra Command 的 Long/Short 字段，不涉及任何逻辑变更
