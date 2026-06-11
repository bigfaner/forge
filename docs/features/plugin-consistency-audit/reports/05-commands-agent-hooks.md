# Report 05: Commands + Agent + Hooks Deep Audit (Layer 1-3)

**Baseline commit**: `2ec06747ecdf5ffdb91abe20d9b07848175b17f8`
**Date**: 2026-05-30
**Scope**: 18 commands + 1 agent (task-executor) + hooks/guide.md

---

## 1. Commands Audit (18 files)

### 1.1 Command Summary

| # | Command | Complexity | Cross-references |
|---|---------|------------|-----------------|
| 1 | clean-code | Low | forge:clean-code skill |
| 2 | eval-consistency | Low | forge:eval skill |
| 3 | eval-contract | Low | forge:eval skill |
| 4 | eval-design | Low | forge:eval skill |
| 5 | eval-journey | Low | forge:eval skill |
| 6 | eval-prd | Low | forge:eval skill |
| 7 | eval-proposal | Low | forge:eval skill |
| 8 | eval-ui | Low | forge:eval skill |
| 9 | execute-task | High | Agent, forge:submit-task, forge:git-commit |
| 10 | extract-design-md | Low | forge:extract-design-md skill |
| 11 | fix-bug | High | forge:git-commit, learn/templates, /consolidate-specs, /learn |
| 12 | gen-sitemap | Low | forge:gen-sitemap skill |
| 13 | git-checkout | Medium | forge feature CLI |
| 14 | git-commit | Medium | (self-contained) |
| 15 | init-forge | Medium | forge-cli/scripts/* |
| 16 | quick | High | forge:brainstorm, forge:quick-tasks, forge:run-tasks |
| 17 | run-tasks | High | Agent, forge:submit-task (implicit), forge:git-commit (via hook) |
| 18 | simplify-skill | Medium | (self-contained) |

### 1.2 Findings

| # | Component | File Path | Layer | Category | Severity | Description | Fix Suggestion | Confidence |
|---|-----------|-----------|-------|----------|----------|-------------|----------------|------------|
| CMD-01 | fix-bug | `commands/fix-bug.md` | 2 | CONFLICT | P1 | `allowed-tools` 缺少 `AskUserQuestion`，但 Step 5 (line 286) 使用 `AskUserQuestion` 呈现知识提取候选项。缺少该工具声明会导致运行时工具不可用。 | 在 `allowed-tools` 中添加 `AskUserQuestion`。 | high |
| CMD-02 | fix-bug | `commands/fix-bug.md` | 2 | CONFLICT | P1 | Bug surface 表格 (lines 141-145) 硬编码 `Playwright` 作为 UI 测试 runner，与 v3.0.0 可插拔 test profile 设计矛盾。test profile 系统已将 Playwright 硬编码替换为可插拔 profile，但 fix-bug 仍硬编码。 | 将 `Playwright` 替换为 `profile runner` 或 `test profile`，保持与 test profile 系统一致。 | high |
| CMD-03 | fix-bug | `commands/fix-bug.md` | 2 | INCOMPLETE | P2 | Convention 和 Business Rule 类型的 Format reference 列引用 `/consolidate-specs` 的 "tech-specs entry format" 和 "biz-specs entry format"，而非 learn/templates/convention-entry.md 模板。Decision 和 Lesson 类型引用 `learn/templates/` 模板，但 Convention/Business Rule 类型不引用 learn 的模板，且 `/consolidate-specs` 技能的 SKILL.md 中无 "domain-to-file mapping" 或 "classify unassisted" 描述。fix-bug line 258 引用 "domain-to-file mapping from /consolidate-specs skill Step 5"，但 Step 5 的实际内容是 "Generate Preview Files + Detect Overlaps"，与 domain-to-file mapping 无关。实际 mapping 在 `rules/overlap-detection.md`。 | 修正引用：将 "domain-to-file mapping from /consolidate-specs skill Step 5" 改为 "domain-to-decision-file mapping from /consolidate-specs rules/overlap-detection.md"。 | high |
| CMD-04 | execute-task | `commands/execute-task.md` | 2 | INCOMPLETE | P2 | `allowed-tools` 包含 `TaskOutput` 和 `Skill`，但命令体中从未调用 `TaskOutput(...)` 或 `Skill(skill=...)`。`Skill` 仅在 Related Commands 表中作为 `/submit-task` 引用出现。`TaskOutput` 在整个文件中无使用场景。这是一个残留声明。 | 如果 TaskOutput 确实不需要，从 allowed-tools 中移除。如果 submit-task 需要由 dispatcher (run-tasks) 处理而非 execute-task 本身，则 Skill 也可能多余。 | medium |
| CMD-05 | execute-task | `commands/execute-task.md` | 3 | TIMING | P2 | Step 1.5 (Main Session Routing) 中第 5 步引用 "same as Step 2b verify logic"，但 Step 2b 位于 Step 1.5 之后的主流程分支中。对于 MAIN_SESSION 路径的读者，Step 2b 尚未出现。虽然逻辑上可理解（向后引用是合理的），但将验证逻辑定义为命名子程序（如 "Verify Record Procedure"）并多处引用会更清晰。 | 提取共享的验证逻辑为独立命名节（如 "### Verify Record"），在 Step 1.5 和 Step 2b 中引用。 | low |
| CMD-06 | run-tasks | `commands/run-tasks.md` | 2 | INCOMPLETE | P2 | Dispatcher Iron Laws (line 28) 要求 "NO code reading"（仅 MAIN_SESSION 例外），但 Step 1.5 要求 "Read task file at FILE"。虽然 "reading the task file" 在 Step 1.5 中已被明确豁免，但 Iron Law 2 的措辞 "EXCEPT for MAIN_SESSION tasks" 可能让 AI 对 "Read" 工具是否可用产生疑问，因为 allowed-tools 包含 `Read` 且全文没有 `Read(...)` 调用。 | 考虑在 Iron Law 2 中明确说明 "Read tool is available for MAIN_SESSION task file reading only"，消除歧义。 | medium |
| CMD-07 | run-tasks | `commands/run-tasks.md` | 2 | INCOMPLETE | P1 | run-tasks 不包含 submit-task 相关逻辑——它依赖 task-executor agent 内部调用 submit-task。但 execute-task 和 run-tasks 的区别在于：execute-task 是单任务入口，不包含 submit-task 步骤（依赖 agent）；run-tasks 是循环调度器，同样不包含。然而 run-tasks Step 1.5 在 MAIN_SESSION 路径中不调用 submit-task 或 git-commit，这意味着 MAIN_SESSION 任务完成后可能缺少 record 创建和 commit 步骤。 | 在 run-tasks Step 1.5 中增加 "After execution, submit record via `Skill(skill="forge:submit-task")` and commit via `Skill(skill="forge:git-commit")` if needed" 的指示。 | high |
| CMD-08 | run-tasks | `commands/run-tasks.md` | 2 | INCOMPLETE | P2 | run-tasks 的 allowed-tools 为 `Bash Read Agent Skill`，但步骤中不直接调用 `Skill(...)`（Skill 仅在 Iron Law 2 中被提及为 "invoking the Skill tool"）。Step 1.5 说 "Follow instructions exactly (task document specifies skill, outcome, record logic)"，这意味着需要 Skill 工具，但允许列表中没有明确使用的场景描述。 | 无需修改——Skill 的声明是为了 MAIN_SESSION 路径中任务可能调用 skill，保留是正确的。但可添加注释说明 Skill 仅用于 MAIN_SESSION。 | low |
| CMD-09 | quick | `commands/quick.md` | 2 | CONFLICT | P2 | Step 2 中 `forge config get auto.runTasks` 的输出解析假设格式为 `key:value pairs`（如 `quick:true full:false`），解析规则是 `stdout contains quick:true`。但 fix-bug 的类似配置检查 (Step 4.5) 使用不同的配置键 `auto.knowledgeSave`，且解析规则是 "Mode value is `true`"。两个命令对 forge config 输出格式的假设不完全一致：quick 直接匹配 `quick:true`，fix-bug 需要先识别 mode 再匹配值。这虽然不导致运行时错误，但增加了维护复杂度。 | 考虑统一两个命令的配置解析逻辑为同一模式，或在 guide.md 中文档化 forge config 的标准输出格式。 | low |
| CMD-10 | quick | `commands/quick.md` | 3 | TIMING | P1 | Step 1 (Brainstorm) 说 "After brainstorm completes, extract the feature slug"，Step 2 说 "The user already approved and committed the proposal in Step 1"。但 brainstorm skill 本身不会自动 commit proposal——commit 需要在 brainstorm 完成后由用户或 agent 执行。Step 2 的前提假设 "already approved and committed" 可能不成立。 | 修正 Step 2 描述：删除 "and committed"，或在 Step 1 和 Step 2 之间增加显式 commit 步骤。 | medium |
| CMD-11 | git-commit | `commands/git-commit.md` | 1 | REFERENCE | P2 | `allowed-tools` 为 `Bash Read`，但步骤中需要运行 `git status` 和 `git diff` 来检查变更（Bash 工具涵盖），以及读取文件内容理解变更（Read 工具涵盖）。工具声明合理，无遗漏。但 `Task Completion Template` 中引用 `Co-Authored-By: Agent`，而系统 CLAUDE.md 中的 git commit 协议使用 `Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>`。两者格式不一致。 | 统一 Co-Authored-By 格式。如果 agent 身份是固定的，应在 git-commit 中使用与系统级别一致的格式。 | medium |
| CMD-12 | simplify-skill | `commands/simplify-skill.md` | 2 | INCOMPLETE | P2 | `allowed-tools` 为 `Read Write Edit AskUserQuestion`，缺少 `Bash`。Phase 4 (Execute Extraction) 中 "Create directory structure" 需要 mkdir 命令，但没有 Bash 工具。虽然 Write 工具可以在某些 CLI 中自动创建父目录，但显式的目录创建通常需要 Bash。 | 添加 `Bash` 到 allowed-tools，或确认 Write 工具能自动创建中间目录。 | medium |
| CMD-13 | init-forge | `commands/init-forge.md` | 2 | INCOMPLETE | P2 | `allowed-tools` 为 `Bash Read`，但 Read 工具在整个命令体中从未被调用。init-forge 只运行 bash 命令 (cd, powershell, forge --version)，不需要读取文件。 | 如果确实不需要 Read，可从 allowed-tools 中移除。影响较小，保留也无害。 | low |
| CMD-14 | fix-bug | `commands/fix-bug.md` | 2 | INCOMPLETE | P2 | Bug surface 表格 (lines 141-145) 缺少 `mobile` 和 `tui` surface 的行。guide.md 定义了 5 种 surface type (web/api/cli/tui/mobile)，test-type-model.md 也列出了 mobile (Mobile E2E Test) 和 tui (Terminal Functional Test)，但 fix-bug 只列出 UI/API/CLI 三种。 | 添加 `mobile` 和 `tui` 行到 Bug surface 表格，或说明为何排除。 | high |
| CMD-15 | gen-sitemap | `commands/gen-sitemap.md` | 1 | REFERENCE | P2 | `allowed-tools` 为 `Bash Read Write Grep Glob`，但 description 提到 "Uses agent-browser to explore routes"。agent-browser 需要通过 Bash 执行 npx 命令，这在 gen-sitemap 的 SKILL.md 中有详细描述。Command 文件本身是薄包装器，直接委托给 skill，所以 allowed-tools 由 skill 定义决定。但 command 和 skill 的 allowed-tools 完全相同（SKILL.md 也是 `Bash Read Write Grep Glob`），而 SKILL.md 中实际需要运行 `npx agent-browser`（Bash 已覆盖）。一致性无问题。 | 无需修复——确认一致性。 | high (confirmed consistent) |

---

## 2. Agent Audit (task-executor)

### 2.1 Agent Structure

`agents/task-executor.md` 包含 4 个主要节：
- Hard Constraints (8 条 `<EXTREMELY-IMPORTANT>` 规则)
- Execution Protocol (6 步)
- Error Handling (分类错误处理 + Pause Protocol)

### 2.2 Findings

| # | Component | File Path | Layer | Category | Severity | Description | Fix Suggestion | Confidence |
|---|-----------|-----------|-------|----------|----------|-------------|----------------|------------|
| AGT-01 | task-executor | `agents/task-executor.md` | 2 | INCOMPLETE | P2 | Hard Constraint 3: "NO BACKGROUND TASKS — all commands run synchronously"，但 Execution Protocol 中未提及如何处理需要 `Agent(...)` 调用的场景（如 Step 4 的 submit-task 和 Step 5 的 git-commit 通过 `Skill(...)` 调用）。虽然 Skill 不是 background task，但 constraint 的措辞 "all commands run synchronously" 是明确的。这不是矛盾，而是缺失的细节说明。 | 在 constraint 3 后补充注释："Skill() and Agent() calls from the Execution Protocol are synchronous invocations, not background tasks." | medium |
| AGT-02 | task-executor | `agents/task-executor.md` | 2 | CONFLICT | P1 | Execution Protocol Step 5 说 "Invoke `Skill(skill="forge:git-commit")`"，但 Step 4 中 submit-task 的阻塞检查逻辑是 "if `blocked`, skip to step 6"。如果 submit 输出显示 `STATUS: blocked`，agent 跳过 commit (step 5) 直接到 step 6 (Done)。然而 Done 输出格式要求 `<commit-hash>`，blocked 状态下没有 commit hash，所以格式中使用 "blocked" 替代。这与 step 6 的输出格式 "DONE: <TASK_ID> | blocked | <summary>" 一致。这不是矛盾——验证通过。 | 无需修复——确认逻辑一致。 | high (confirmed consistent) |
| AGT-03 | task-executor | `agents/task-executor.md` | 2 | INCOMPLETE | P2 | Error Handling 定义了 Fix-Type Derivation 表 (doc/eval -> doc.fix; coding/test/validation/gate -> coding.fix)，但未覆盖 `validation` 和 `gate` 类型的边界情况。如果 TASK_CATEGORY 不是表中的任何一种（如未知类型），错误处理流程未定义 fallback。 | 添加 fallback 规则："Unknown TASK_CATEGORY defaults to coding.fix (conservative choice)." | low |
| AGT-04 | task-executor | `agents/task-executor.md` | 2 | CONFLICT | P2 | Hard Constraint 8 (SPEC AUTHORITY FALLBACK) 说 "if the synthesized strategy does not include a Reference Files declaration, you MUST still read the task file's `## Reference Files` section"。但 Execution Protocol Step 3 说 "Follow the synthesized strategy exactly"。如果策略包含 Reference Files 声明，Constraint 8 不触发；如果不包含，则触发 fallback。这两条规则实际上是互补的，不是矛盾。然而 "synthesized strategy" 这个术语来自 `forge prompt get-by-task-id` 的输出，而非 task-executor.md 本身定义——agent 需要理解外部概念。 | 考虑在 agent 文件中简要说明 "synthesized strategy" 指的是 `forge prompt get-by-task-id` 输出的策略。 | low |

---

## 3. Hooks Audit (guide.md)

### 3.1 Scope

Per task Implementation Notes: guide.md 审计范围限于 (1) 脚本路径存在性；(2) 参数描述一致性；(3) 内部步骤无矛盾。不深入验证脚本逻辑。

### 3.2 guide.md Analysis

guide.md 内容分为 4 节：
- Directory Conventions (项目级文档目录结构)
- Manifest (feature manifest 说明)
- Forge CLI (CLI 命令参考)
- Terminology (术语定义)

guide.md **不引用** hooks/ 目录下的任何脚本文件 (hooks.json, run-hook.cmd, session-start, debug)。它是一个独立的参考文档，被 session-start hook 脚本注入到 AI 会话上下文中。

### 3.3 Findings

| # | Component | File Path | Layer | Category | Severity | Description | Fix Suggestion | Confidence |
|---|-----------|-----------|-------|----------|----------|-------------|----------------|------------|
| HOOK-01 | guide.md | `hooks/guide.md` | 1 | REFERENCE | P2 | guide.md line 46 说 Surface Type orchestration 模式为 `web`/`api` require probe + teardown; `cli`/`tui` use build -> dev -> test。但 `mobile` surface 的编排模式未被说明——它被归入哪一组？根据 run-tests/SKILL.md line 163，mobile 的编排是 "test-setup -> dev -> probe -> [per-journey test loop] -> teardown"，与 web/api 类似（有 probe + teardown），但有额外的 test-setup 步骤。guide.md 的括号说明遗漏了 mobile。 | 修正 line 46 为："(e.g. `web`/`api`/`mobile` require probe + teardown; `cli`/`tui` use build -> dev -> test)". 或单独说明 mobile 的 test-setup -> dev -> probe -> test -> teardown 序列。 | high |
| HOOK-02 | guide.md | `hooks/guide.md` | 1 | REFERENCE | P2 | guide.md line 47 的 Test Type 示例仅列出 `cli` -> CLI Functional Test, `api` -> API Functional Test, `web` -> Web E2E Test，遗漏了 `tui` -> Terminal Functional Test 和 `mobile` -> Mobile E2E Test。虽然后面说 "See test-type-model.md for the full mapping"，但示例不完整。 | 添加 `tui` 和 `mobile` 到示例列表。 | medium |
| HOOK-03 | guide.md | `hooks/guide.md` | 2 | INCOMPLETE | P2 | guide.md 的 Forge CLI 节列出了部分 CLI 命令 (proposal, feature status, task transition, task reopen)，但遗漏了多个被 commands 和 agent 频繁使用的命令：`forge task claim`, `forge task status`, `forge task add`, `forge feature set`, `forge feature complete`, `forge config get`, `forge quality-gate`, `forge cleanup`, `forge prompt get-by-task-id`, `forge surfaces detect`。guide.md 是 session-start hook 注入的上下文文档，遗漏这些命令意味着 AI 在每个会话中都不知道这些命令的存在。 | 添加常用 CLI 命令到 Forge CLI 节，或按类别分组 (Task Management, Feature Management, Configuration, Pipeline)。 | high |
| HOOK-04 | guide.md | `hooks/guide.md` | 2 | INCOMPLETE | P2 | guide.md 不包含 `forge config get auto.*` 配置机制的任何说明。`quick` 和 `fix-bug` 命令使用 `forge config get auto.runTasks` 和 `forge config get auto.knowledgeSave` 进行自动跳过配置检查，但 guide.md 中无此配置键的文档。 | 在 Forge CLI 节增加 Configuration 子节，说明 `forge config get auto.*` 的配置键和含义。 | medium |

---

## 4. Cross-Component Consistency Checks

### 4.1 execute-task vs run-tasks

两个命令共享大量相似逻辑（Step 1.5 Main Session Routing, Step 2 Dispatch + Verify, Error Handling）。关键差异：

| Dimension | execute-task | run-tasks |
|-----------|-------------|-----------|
| Loop | Single task, no loop | Continuous claim loop |
| Failure tracking | None | consecutive_failures counter |
| Step 0 (feature set) | Missing | Present |
| Post-completion | STOP | Summary + git status |
| submit-task | Referenced but not invoked | Implicit (agent handles) |

**Finding**: execute-task 缺少 `forge feature set` 步骤（run-tasks Step 0），但如果 execute-task 独立运行（非从 run-tasks 调度），feature context 可能未设置。这不是 bug（execute-task 假设 feature 已由 run-tasks 设置），但缺少文档说明。

### 4.2 execute-task + run-tasks vs task-executor agent

三者共享 Fix-Type Derivation 表。验证一致性：

| Source | doc/eval | coding/test/validation/gate |
|--------|----------|----------------------------|
| execute-task | doc.fix | coding.fix |
| run-tasks | doc.fix | coding.fix |
| task-executor | doc.fix | coding.fix |

**Result**: 三方一致，无矛盾。

### 4.3 fix-bug vs test-type-model.md

fix-bug 的 Bug surface 表格 vs test-type-model.md 的完整映射：

| Surface | fix-bug | test-type-model.md | Match? |
|---------|---------|-------------------|--------|
| cli | CLI command / child_process | CLI Functional Test / 子进程执行 | YES |
| api | API endpoint / fetch | API Functional Test / HTTP 客户端 | YES |
| web | UI behavior / Playwright | Web E2E Test / 浏览器自动化 | PARTIAL (Playwright hardcoded) |
| tui | MISSING | Terminal Functional Test / 子进程+stdin pipe | NO |
| mobile | MISSING | Mobile E2E Test / Maestro YAML | NO |

**Result**: 确认 CMD-02 和 CMD-14 findings。

### 4.4 guide.md Terminology vs test-type-model.md

guide.md line 46 描述 Surface Type orchestration vs run-tests surface rules:

| Surface | guide.md says | run-tests surface rule says | Match? |
|---------|--------------|---------------------------|--------|
| web | probe + teardown | dev -> probe -> test -> teardown | PARTIAL (guide omits dev/test) |
| api | probe + teardown | dev -> probe -> test -> teardown | PARTIAL |
| cli | build -> dev -> test | build -> dev -> test | YES |
| tui | build -> dev -> test | build -> dev -> test | YES |
| mobile | NOT MENTIONED | test-setup -> dev -> probe -> test -> teardown | NO |

**Result**: 确认 HOOK-01 finding。

---

## 5. Summary Statistics

| Metric | Value |
|--------|-------|
| Total commands audited | 18 |
| Agent audited | 1 |
| Hooks files audited | 1 (guide.md) |
| Total findings | 23 |
| P0 (Critical) | 0 |
| P1 (High) | 4 (CMD-01, CMD-02, CMD-07, CMD-10) |
| P2 (Medium) | 15 |
| P3 (Low) | 4 |
| Category: CONFLICT | 4 |
| Category: INCOMPLETE | 12 |
| Category: TIMING | 2 |
| Category: REFERENCE | 5 |

### Issue Category Distribution

| Category | Count | Notes |
|----------|-------|-------|
| CONFLICT | 4 | fix-bug AskUserQuestion missing (CMD-01), Playwright hardcode (CMD-02), config parsing inconsistency (CMD-09), agent commit-skip logic (AGT-02, confirmed consistent) |
| INCOMPLETE | 12 | fix-bug cross-reference imprecision (CMD-03), unused tool declarations (CMD-04, CMD-13), missing mobile/tui rows (CMD-14), guide.md CLI coverage gaps (HOOK-03, HOOK-04), agent edge cases (AGT-03, AGT-04) |
| TIMING | 2 | execute-task forward reference (CMD-05), quick Step 2 commit assumption (CMD-10) |
| REFERENCE | 5 | gen-sitemap consistency check (CMD-15, confirmed OK), git-commit Co-Authored-By (CMD-11), guide.md mobile omission (HOOK-01, HOOK-02), simplify-skill missing Bash (CMD-12) |

### P1 Findings Summary

| ID | Component | Description |
|----|-----------|-------------|
| CMD-01 | fix-bug | `allowed-tools` 缺少 `AskUserQuestion`，但 Step 5 使用该工具 |
| CMD-02 | fix-bug | Bug surface 表格硬编码 Playwright，与 v3.0.0 可插拔 test profile 矛盾 |
| CMD-07 | run-tasks | MAIN_SESSION 路径缺少 submit-task 和 git-commit 步骤 |
| CMD-10 | quick | Step 2 假设 proposal 已 committed，但 brainstorm 不自动 commit |

---

## 6. Components with No Issues (Confirmed Clean)

The following commands passed all three audit layers with no findings:

| # | Command | Status |
|---|---------|--------|
| 1 | clean-code | Clean - thin wrapper, consistent with skill |
| 2 | eval-consistency | Clean - uniform eval-* pattern |
| 3 | eval-contract | Clean - uniform eval-* pattern |
| 4 | eval-design | Clean - uniform eval-* pattern |
| 5 | eval-journey | Clean - uniform eval-* pattern |
| 6 | eval-prd | Clean - uniform eval-* pattern |
| 7 | eval-proposal | Clean - uniform eval-* pattern |
| 8 | eval-ui | Clean - uniform eval-* pattern |
| 9 | extract-design-md | Clean - thin wrapper, tools match skill |
| 10 | git-checkout | Clean - self-contained workflow |
| 11 | gen-sitemap | Clean - thin wrapper, tools match skill |

All 7 eval-* commands follow an identical pattern (`Skill(skill="forge:eval", args="--type <type> [...]")`) and are structurally consistent. Their frontmatter `name`, `description`, and `argument-hint` fields are consistent with each other and with the forge:eval skill.
