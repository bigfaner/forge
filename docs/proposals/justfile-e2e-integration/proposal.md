---
created: 2026-04-29
author: faner
status: Draft
---

# Proposal: Justfile 作为全局命令抽象层

## Problem

Skills、agents 和 task 模板中散落着大量原始 shell 命令——既有代码块中的可执行命令，也有文字描述中的命令示例。Agent 执行时必须自行推断工作目录、路径、语言工具链和参数，导致命令不稳定、难以维护。

`init-justfile` 已定义标准目标契约（`just test`、`just build`、`just test-e2e`），但没有任何 skill 或 agent 引用它们。justfile 作为抽象层的价值完全未被利用。

### Evidence

全量审计发现两类问题，覆盖代码块和文字描述：

#### A. E2E 命令

| 位置 | 原始命令/描述 | 类型 |
|------|-------------|------|
| `run-e2e-tests` Step 1 | `cd tests/e2e && npm install` | 代码块 |
| `run-e2e-tests` Step 1 | `npx playwright install chromium` | 代码块 |
| `run-e2e-tests` Step 1 | `npx playwright --version` | 代码块 |
| `run-e2e-tests` Step 1 | "Install dependencies (if `node_modules` doesn't exist):" | 文字描述 |
| `run-e2e-tests` Step 1 | "Install Playwright browser (if not installed):" | 文字描述 |
| `run-e2e-tests` Step 2 | `npx tsx <slug>/cli.spec.ts 2>&1 \| tee ...` | 代码块 |
| `run-e2e-tests` Step 2 | `npx tsx <slug>/api.spec.ts 2>&1 \| tee ...` | 代码块 |
| `run-e2e-tests` Step 2 | `npx tsx <slug>/ui.spec.ts 2>&1 \| tee ...` | 代码块 |
| `run-e2e-tests` Error table | "Run `npx playwright install chromium`, retry" | 文字描述 |
| `run-e2e-tests` Error table | "Run `npm install`, retry" | 文字描述 |
| `gen-test-scripts` Step 4 | `grep -r '// VERIFY:' tests/e2e/<feature>/` | 代码块 |
| `gen-test-scripts` Step 4 | "run `grep -r '// VERIFY:'`" | 文字描述 |
| `gen-test-scripts` Step 5 | `cd tests/e2e && npm install` | 代码块 |
| `fix-bug` Step 3b | `npx tsx <spec-file> 2>&1` | 代码块 |
| `fix-bug` Step 3b | "Run the e2e test — it must fail before the fix:" | 文字描述 |
| `fix-bug` Step 5 | `npx tsx <spec-file> 2>&1` | 代码块 |
| `fix-e2e` template | 无明确重跑命令 | 缺失 |

#### B. 单元测试 / 构建命令

| 位置 | 原始命令/描述 | 类型 |
|------|-------------|------|
| `task-executor` Step 3 | `go build ./... && go vet ./... && go test -race -cover ./...` | 代码块注释 |
| `task-executor` Step 3 | `npm run build && npm test` | 代码块注释 |
| `task-executor` Step 3 | `pytest --cov` | 代码块注释 |
| `task-executor` Step 3 | "Run complete verification suite for your project:" | 文字描述 |
| `error-fixer` Step 4 | 同上三条语言示例 | 代码块注释 |
| `error-fixer` Step 4 | "Run complete verification suite for your project:" | 文字描述 |
| `execute-task` Step 3 | "Run project-specific verification commands." | 文字描述 |
| `fix-bug` Step 2 | `<project-test-command>   # e.g. npm test, go test ./..., pytest` | 代码块 |
| `fix-bug` Step 2 | "Run existing tests to establish baseline" | 文字描述 |
| `fix-bug` Step 3a | `<project-test-command> --testNamePattern "bug:"` | 代码块 |
| `fix-bug` Step 3a | "Run the new test — it must fail before the fix:" | 文字描述 |
| `fix-bug` Step 5 | `<project-test-command>` | 代码块 |
| `fix-bug` Step 5 | "Run the full test suite. All tests must pass." | 文字描述 |
| `run-tasks` Step 5 | `# go test ./... OR npm test OR the testCommand from index.json` | 代码块注释 |
| `run-tasks` Step 5 | "Run project-level full test suite" | 文字描述 |
| `record-task` Metrics | `go test -cover ./changed/package/...` | 代码块 |
| `record-task` Metrics | `npm test -- --coverage --watchAll=false` | 代码块 |
| `record-task` Metrics | `pytest --cov=<module> --cov-report=term-missing` | 代码块 |
| `improve-harness` Step 4.3 | "Run project test suite to ensure nothing broke" | 文字描述 |

### Urgency

每次 agent 执行任务，都在重新解决"如何调用"的问题。`just test` 和 `just build` 已在标准契约中定义，但没有任何 skill 引用它们——这意味着 justfile 的存在对 agent 行为毫无影响。命令散落在 12 个文件中，一旦工具链变化需要同步修改多处。

## Proposed Solution

扩展 justfile 标准目标契约，新增两个 e2e 专用目标，然后将所有 skill/agent 中的原始命令（代码块和文字描述）统一替换为对应的 just 命令。

### 新增目标（扩展 init-justfile 契约）

**`just e2e-setup`** — 幂等安装（deps + Playwright browser）：

```just
# Install e2e dependencies and Playwright browser (idempotent)
e2e-setup:
    #!/usr/bin/env bash
    [ ! -d tests/e2e/node_modules ] && npm install --prefix tests/e2e || true
    npx --prefix tests/e2e playwright install chromium
```

**`just e2e-verify --feature <slug>`** — 检查未解析的 `// VERIFY:` 标记，有残留则 exit 1：

```just
# Verify no unresolved // VERIFY: markers remain in feature test scripts (exit 1 if any found)
[arg("feature", long)]
e2e-verify feature="":
    #!/usr/bin/env bash
    dir="tests/e2e/{{feature}}"
    if [ -z "{{feature}}" ]; then
        echo "Usage: just e2e-verify --feature <slug>"
        exit 1
    fi
    count=$(grep -r '// VERIFY:' "$dir" 2>/dev/null | wc -l | tr -d ' ')
    if [ "$count" -gt 0 ]; then
        echo "ERROR: $count unresolved // VERIFY: marker(s) in $dir:"
        grep -rn '// VERIFY:' "$dir"
        exit 1
    fi
    echo "OK: no unresolved // VERIFY: markers in $dir"
```

### 完整命令替换映射

#### E2E 命令

| 原始命令/描述 | 替换为 |
|-------------|--------|
| `cd tests/e2e && npm install` | `just e2e-setup` |
| `npx playwright install chromium` | （合并到 `just e2e-setup`，不再单独调用） |
| `npx playwright --version` | （合并到 `just e2e-setup`） |
| "Install dependencies..." / "Install Playwright browser..." | "Run `just e2e-setup`" |
| "Run `npx playwright install chromium`, retry" | "Run `just e2e-setup`, retry" |
| "Run `npm install`, retry" | "Run `just e2e-setup`, retry" |
| `npx tsx <slug>/*.spec.ts 2>&1 \| tee ...`（三条） | `just test-e2e --feature <slug>` |
| `npx tsx <spec-file> 2>&1`（fix-bug 两处） | `just test-e2e --feature <slug>` |
| "Run the e2e test — it must fail before the fix:" | "Run `just test-e2e --feature <slug>`" |
| `grep -r '// VERIFY:' tests/e2e/<feature>/` | `just e2e-verify --feature <slug>` |
| "run `grep -r '// VERIFY:'`" | "run `just e2e-verify --feature <slug>`" |
| fix-e2e template 无重跑命令 | 新增：运行 `just test-e2e --feature <slug>` 验证修复 |

#### 单元测试 / 构建命令

| 原始命令/描述 | 替换为 |
|-------------|--------|
| `go build ./... && go vet ./... && go test -race -cover ./...` | `just build && just test` |
| `npm run build && npm test` | `just build && just test` |
| `pytest --cov` | `just build && just test` |
| "Run complete verification suite for your project: [语言示例]" | "Run `just build && just test`" |
| "Run project-specific verification commands." | "Run `just build && just test`" |
| `<project-test-command>   # e.g. npm test, go test ./..., pytest` | `just test` |
| `<project-test-command> --testNamePattern "bug:"` | `just test` |
| `<project-test-command>`（Step 5） | `just test` |
| "Run existing tests to establish baseline" | "Run `just test` to establish baseline" |
| "Run the new test — it must fail before the fix:" | "Run `just test`" |
| "Run the full test suite. All tests must pass." | "Run `just build && just test`" |
| `# go test ./... OR npm test OR the testCommand from index.json` | `just test` |
| "Run project-level full test suite" | "Run `just test`" |
| `go test -cover ./...` / `npm test -- --coverage` / `pytest --cov=...` | `just test` |
| "Run project test suite to ensure nothing broke" | "Run `just test`" |

### 各文件变更摘要（12 个文件）

| 文件 | 变更内容 |
|------|---------|
| `init-justfile` command | 标准目标契约新增 `e2e-setup`、`e2e-verify` 两行及 recipe 模板 |
| `gen-test-scripts` SKILL.md | Step 4 VERIFY 检查 + Step 5 deps install → just 命令；文字描述同步 |
| `run-e2e-tests` SKILL.md | Step 1 setup + Step 2 run specs + Error table → just 命令；文字描述同步 |
| `fix-bug` command | `<project-test-command>` → `just test`；`npx tsx` → `just test-e2e --feature <slug>`；文字描述同步 |
| `run-tasks` command | Breaking Gate 注释命令 + 文字描述 → `just test` |
| `task-executor` agent | Step 3 语言示例 + 文字描述 → `just build && just test` |
| `error-fixer` agent | Step 4 语言示例 + 文字描述 → `just build && just test` |
| `execute-task` command | Step 3 文字描述 → `just build && just test` |
| `record-task` SKILL.md | Metrics Collection 三种语言示例 → `just test` |
| `improve-harness` SKILL.md | Step 4.3 文字描述 → `just test` |
| `breakdown-tasks` 模板 `run-e2e-tests.md` | Implementation Notes → `just test-e2e --feature <slug>` |
| `breakdown-tasks` 模板 `gen-test-scripts.md` | Implementation Notes → `just e2e-verify --feature <slug>` |
| `breakdown-tasks` 模板 `fix-e2e.md` | 新增修复后验证步骤 → `just test-e2e --feature <slug>` |

## Alternatives Considered

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | 零改动 | 原始命令继续散落 12 个文件；justfile 形同虚设；agent 每次重新推断工具链 | Rejected: 维护成本持续累积 |
| 只替换代码块命令，保留文字描述 | 改动量小 | 文字描述仍会引导 agent 使用原始命令；效果减半 | Rejected: 不彻底 |
| 只抽象 e2e 命令 | 解决最高频问题 | 单元测试/构建命令仍是占位符；task-executor 和 error-fixer 仍需判断语言 | Rejected: 只解决一半 |
| 全面替换：新增 e2e-setup + e2e-verify，所有 skill/agent 的代码块和文字描述统一引用 just 命令 | 命令有唯一权威来源；skill 只需知道目标名称；工具链变化只改 justfile | 需修改 13 个文件 | **Selected** |

## Scope

### In Scope

- `init-justfile` command：标准目标契约新增 `e2e-setup`、`e2e-verify` 两个目标及模板实现
- `gen-test-scripts` SKILL.md：Step 4 VERIFY 检查 + Step 5 deps install → just 命令（代码块 + 文字描述）
- `run-e2e-tests` SKILL.md：Step 1 setup + Step 2 run specs + Error table → just 命令（代码块 + 文字描述）
- `fix-bug` command：`<project-test-command>` → `just test`；`npx tsx` → `just test-e2e --feature <slug>`（代码块 + 文字描述）
- `run-tasks` command：Breaking Gate → `just test`（代码块注释 + 文字描述）
- `task-executor` agent：Step 3 → `just build && just test`（代码块注释 + 文字描述）
- `error-fixer` agent：Step 4 → `just build && just test`（代码块注释 + 文字描述）
- `execute-task` command：Step 3 文字描述 → `just build && just test`
- `record-task` SKILL.md：Metrics Collection 语言示例 → `just test`
- `improve-harness` SKILL.md：Step 4.3 文字描述 → `just test`
- `breakdown-tasks` 模板（run-e2e-tests.md、gen-test-scripts.md、fix-e2e.md）：Implementation Notes → just 命令

### Out of Scope

- 项目级 justfile（forge 自身的 `justfile`）：只有 `claude`/`claude-c`，不是应用项目的 justfile
- `gen-test-cases` SKILL.md：不涉及命令执行
- `graduate-tests` SKILL.md：毕业流程不涉及 e2e 执行命令
- `gen-sitemap` command：`npx agent-browser` 是专用工具，不属于标准 just 目标范畴
- 现有已生成的 spec 文件：不回溯修改

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| 目标项目未运行 `/init-justfile`，justfile 不存在 | Medium | **High** — agent 调用 `just test` 报错，整个流程中断 | Skills 在调用 just 命令前先检查 `ls justfile`；不存在则提示用户运行 `/init-justfile` |
| `just e2e-verify` 的 `[arg]` 语法要求 just >= 1.50.0 | Low | **Medium** — 旧版 just 解析失败 | `init-justfile` 已有版本检查要求（>= 1.50.0），无需额外处理 |
| `just test-e2e --feature <slug>` 合并执行三个 spec，tee 输出格式变化影响 run-e2e-tests 的结果解析 | Low | **Medium** — pass/fail 统计不准确 | `test-e2e` 逐个执行 spec 并累计 fail 计数，输出格式与单独执行一致；run-e2e-tests 解析逻辑不变 |
| `just build && just test` 替换语言示例后，覆盖率数据来源不明确 | Low | **Low** — record-task 的 coverage 字段可能为 -1.0 | `just test` 输出包含覆盖率；record-task 从 `just test` 输出解析，与原来从语言命令输出解析等价 |
| agent 在 fix-e2e 任务中仍使用原始命令（习惯路径依赖） | Medium | **Low** — 功能正确但绕过抽象层 | task 模板 Implementation Notes 明确写 just 命令，不留歧义 |

## Success Criteria

- [ ] `init-justfile` command 的标准目标契约表包含 `e2e-setup` 和 `e2e-verify` 两行，且模板中有对应 recipe 实现（验证：`grep -c 'e2e-setup\|e2e-verify' plugins/forge/commands/init-justfile.md` >= 4）
- [ ] `gen-test-scripts` SKILL.md 中不再出现 `npx playwright install`、`cd tests/e2e && npm install`、`grep.*VERIFY`（验证：各 grep 结果 = 0），且出现 `just e2e-verify` 和 `just e2e-setup`（>= 1 次）
- [ ] `run-e2e-tests` SKILL.md 中不再出现 `npx tsx`、`npx playwright install`、`npm install`（验证：各 grep 结果 = 0），且出现 `just e2e-setup` 和 `just test-e2e --feature`（>= 1 次）
- [ ] `fix-bug` command 中不再出现 `<project-test-command>` 或 `npx tsx <spec-file>`（验证：grep 结果 = 0），替换为 `just test` 和 `just test-e2e --feature <slug>`
- [ ] `task-executor` agent 和 `error-fixer` agent 的 Verification 步骤中出现 `just build && just test`，不再出现 `go test`、`npm test`、`pytest` 等语言特定命令（验证：各 grep 结果 = 0）
- [ ] `execute-task` command Step 3 出现 `just build && just test`（验证：grep 结果 >= 1）
- [ ] `record-task` SKILL.md 的 Metrics Collection 出现 `just test`，不再出现三种语言命令（验证：`grep -c 'go test\|npm test\|pytest' plugins/forge/skills/record-task/SKILL.md` = 0）
- [ ] `improve-harness` SKILL.md Step 4.3 出现 `just test`（验证：grep 结果 >= 1）
- [ ] `breakdown-tasks` 三个 e2e 相关模板（run-e2e-tests.md、gen-test-scripts.md、fix-e2e.md）的 Implementation Notes 中均出现对应 just 命令（验证：各文件 `grep -c 'just ' <template>` >= 1）

## Next Steps

- Proceed to `/write-prd` to formalize requirements
