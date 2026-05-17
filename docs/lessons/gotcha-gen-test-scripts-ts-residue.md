---
created: "2026-05-17"
tags: [testing, architecture]
---

# gen-test-scripts Step 3.5 无条件写入 .ts 残留 + git add -A 放大

## Problem

已通过 `git rm` 删除的 Node.js/Playwright 基础设施文件（20 个 .ts 文件）在后续 fix 任务提交中意外恢复。fix-2 任务本应只修改 2 个 .go 测试文件，实际 commit 包含 169 个文件。

## Root Cause

**三层因果链：**

1. **触发条件**：gen-test-scripts SKILL.md Step 3.5（Shared Infrastructure）硬编码写入 .ts 文件（`helpers.ts`、`playwright.config.ts`、`tsconfig.json`、`auth-setup.ts`），不检查当前 profile 类型。go-test profile 的正确共享基础设施是 `helpers.go` 和 `main_test.go`，但 Step 3.5 仍然写入 TypeScript 模板文件。

2. **残留积累**：每次 /gen-test-scripts 调用（无论 profile 是什么）都会向 `tests/e2e/` 写入 .ts 文件。这些文件不在 .gitignore 中，散落在工作目录成为未跟踪文件。git rm 删除了已跟踪的 .ts，但 Step 3.5 的下一次调用会重新生成它们。

3. **提交放大**：error-fixer agent（error-fixer.md）提交时使用 `git add -A`，将工作目录中所有未跟踪和已修改文件一次性全部 stage。导致 fix-2 的 commit（d7f8a13）包含 169 个文件，其中 20 个 .ts 文件是已删除基础设施被 Step 3.5 重新生成后的残留，加上 47 个 .md、76 个 .go 等其他分支上的未提交文件。

**证据：**
- 对比 `helpers.ts` 内容：删除前（860b7af^）与恢复后（d7f8a13）完全一致（IDENTICAL），证实是模板重新生成而非手动恢复
- 文件时间戳 02:59 位于 fix-1 完成（02:52）和 fix-2 开始（03:02）之间，对应 quality-gate 触发过程中某个 gen-test-scripts 调用

## Solution

1. **Step 3.5 需要感知 profile**：gen-test-scripts 的共享基础设施步骤应根据 profile manifest 决定生成什么文件。go-test profile 生成 `helpers.go` + `main_test.go`，web-playwright profile 才生成 `helpers.ts` + `playwright.config.ts`。
2. **Agent 提交必须精确 stage**：error-fixer 和 task-executor 应使用 `git add <具体文件>` 而非 `git add -A`，只 stage 任务记录中声明的文件。

## Reusable Pattern

**Skill 模板不能假设单一技术栈。** 当 skill 支持多 profile（go-test、web-playwright、pytest 等），所有文件生成步骤都必须读取 profile 配置，不能硬编码某种语言/框架的文件名。共享基础设施（helpers、config）和测试文件（.spec.ts、_test.go）都属于 profile 感知范围。

**Agent 提交范围必须与任务声明一致。** 如果任务记录说"Files Modified: A.go, B.go"，commit 就只能包含 A.go 和 B.go。`git add -A` 是 agent 的反模式——它会拾取工作目录中所有无关的未跟踪文件和未提交变更。

## Example

```
# 错误：Step 3.5 无条件写入 .ts
# gen-test-scripts SKILL.md:
"Generate tests/e2e/helpers.ts from template (if not already present)"
→ 即使 profile=go-test 也会写入 helpers.ts

# 正确：Step 3.5 应根据 profile 决定
if profile.hasCapability("web-ui"):
    generate helpers.ts, playwright.config.ts
elif profile.name == "go-test":
    skip (helpers.go managed separately by go-test strategy)
```

```
# 错误：agent 提交
git add -A && git commit -m "fix: update test assertions"

# 正确：精确 stage
git add tests/e2e/quick_test_slim_cli_test.go tests/e2e/test_scripts_per_type_cli_test.go
git commit -m "fix: update test assertions"
```

## Related Files

- `plugins/forge/skills/gen-test-scripts/SKILL.md` — Step 3.5 共享基础设施逻辑
- `plugins/forge/agents/error-fixer.md` — fix agent 定义
- `tests/e2e/helpers.ts` — 被无条件重新生成的共享基础设施
- `docs/features/auto-behavior-config/tasks/records/fix-2.md` — 只声明修改 2 个 .go 文件的实际任务记录
