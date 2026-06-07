---
feature: "cli-doc-accuracy-audit"
journey: "guide-accuracy"
risk_level: "High"
surface_types: ["cli"]
surface_keys: ["cli"]
sources:
  - docs/proposals/cli-doc-accuracy-audit/proposal.md
generated: "2026-06-07"
---

# Journey: guide-accuracy

**Risk Level**: High

<!-- Risk Classification Criteria:
  High = guide.md 修改影响所有 agent 会话，错误命令名会导致任务执行失败
-->

## Overview

Agent 读取 guide.md 系统提示后执行 Forge CLI 命令，验证 guide 中引用的命令名、参数、标志和行为描述均与实际 CLI 行为一致。

## Setup

<!-- Preconditions that must be established before the Journey starts. -->

- Forge CLI 已编译且可执行（`go build` 成功）
- guide.md 文件存在且包含 CLI 参考部分
- 测试环境有一个合法的 Forge 项目目录（含 `.forge/config.yaml`）

## Happy Path

### Step 1: Agent 验证 guide.md 中命令名存在

**User Action**: 提取 guide.md 中所有以 `forge ` 开头的命令名，逐一运行 `forge <command> --help` 确认命令存在且 Available Commands 列表匹配。

**Expected Result**: guide.md 中每个 `forge ` 命令名都出现在 `forge --help` 或 `forge <group> --help` 的 Available Commands 中。特别是 `forge task validate [file]` 存在（而非已删除的 `validate-index`）。

### Step 2: Agent 按 guide.md 描述使用 task validate 命令

**User Action**: 按 guide.md 描述运行 `forge task validate docs/features/any/tasks/index.json`，验证命令接受文件路径参数。

**Expected Result**: 命令成功执行（exit 0）或返回有意义的验证错误。不再出现 "unknown command validate-index" 错误。

### Step 3: Agent 按 guide.md 描述使用 quality-gate 命令

**User Action**: 运行 `forge quality-gate` 并观察输出，确认 guide.md 中描述的行为（包括自动创建 fix task 的副作用）与实际一致。

**Expected Result**: quality-gate 输出包含 pass/fail 状态，失败时自动创建 fix task。guide.md 中对该命令的描述涵盖这一副作用。

### Step 4: Agent 按 guide.md 描述使用 cleanup 命令

**User Action**: 运行 `forge cleanup` 并确认清理范围包括 completed、blocked、suspended、rejected 状态的任务（而非仅 completed）。

**Expected Result**: cleanup 命令清理 guide.md 中列出的所有任务状态。输出显示被清理的工件列表。

### Step 5: Agent 按 guide.md 使用新增命令和标志

**User Action**: 逐一运行 guide.md 中新增的命令和标志：`forge task query <id-or-key>`、`forge task check-deps`、`forge feature list`、`forge task list --tree`、`forge task submit --quiet`。

**Expected Result**: 每个命令和标志都存在且行为与 guide.md 描述一致。无 "unknown command" 或 "unknown flag" 错误。

## Edge Cases

### Step 1b: guide.md 引用不存在的命令

**Precondition**: 修改前 guide.md 中仍包含 `forge task validate-index`（旧命令名）

**User Action**: 尝试运行 `forge task validate-index docs/features/any/tasks/index.json`

**Expected Result**: 命令返回 "unknown command" 错误（exit code != 0），stderr 包含错误信息。这证明 guide.md 修复前的引用是错误的。

### Step 2b: validate 命令无参数调用

**Precondition**: 工作目录不是 Forge 项目根目录，且未提供文件路径参数

**User Action**: 运行 `forge task validate`（不带文件路径）

**Expected Result**: 命令使用默认的 `index.json` 路径查找，若不存在则报告验证失败。不会 crash。

### Step 3b: quality-gate 在无任务项目上运行

**Precondition**: Forge 项目中没有任何任务记录

**User Action**: 在空项目目录中运行 `forge quality-gate`

**Expected Result**: 命令返回有意义的结果（可能是 "no tasks to validate" 或类似信息），不会 panic 或产生误导性输出。

### Step 4b: cleanup 在干净项目上运行

**Precondition**: 项目中没有 stale 状态的任务或工件

**User Action**: 运行 `forge cleanup`

**Expected Result**: 命令正常完成，输出指示无内容需要清理（exit 0），不会删除有效文件。

### Step 5b: task query 使用无效 ID

**Precondition**: 提供一个不存在的任务 ID 或 key

**User Action**: 运行 `forge task query non-existent-task-id`

**Expected Result**: 命令返回 "not found" 错误（exit code != 0），stderr 包含描述性错误信息。

### Step 5c: task list --tree 在无依赖项目中运行

**Precondition**: 项目中的任务没有任何依赖关系

**User Action**: 运行 `forge task list --tree`

**Expected Result**: 命令正常输出任务列表（扁平结构），不会因缺少依赖树而 crash。

## Journey Invariants

- guide.md 中引用的每个 CLI 命令名必须在 `forge --help` 或对应子命令的 help 输出中存在
- 修改后的 guide.md 不引入任何新的 markdown 格式错误（代码块闭合、链接有效）
- 所有 CLI 行为变更仅限于字符串常量修改，`go build ./...` 和 `go test ./...` 始终通过
