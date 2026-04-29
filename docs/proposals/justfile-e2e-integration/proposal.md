---
created: 2026-04-29
author: faner
status: Draft
---

# Proposal: Justfile 作为全局命令抽象层

## Problem

Skills、agents 和 task 模板中散落着大量原始 shell 命令，agent 执行时必须自行推断工作目录、路径、语言工具链和参数，导致命令不稳定、难以维护。

`init-justfile` 已定义标准目标契约（`just test`、`just build`、`just test-e2e`），但没有任何 skill 或 agent 引用它们。justfile 作为抽象层的价值完全未被利用。

### Evidence

全量审计发现两类原始命令问题：

#### A. E2E 命令（散落在 e2e 相关 skills 中）

| 位置 | 原始命令 | 问题 |
|------|---------|------|
| `run-e2e-tests` Step 1 | `cd tests/e2e && npm install` | 工作目录依赖 agent 当前位置 |
| `run-e2e-tests` Step 1 | `npx playwright install chromium` | 未检查是否已安装，每次重复执行 |
| `run-e2e-tests` Step 2 | `npx tsx <slug>/cli.spec.ts 2>&1 \| tee ...` | slug 需 agent 自行替换，路径拼接易错 |
| `run-e2e-tests` Step 2 | `npx tsx <slug>/api.spec.ts 2>&1 \| tee ...` | 同上 |
| `run-e2e-tests` Step 2 | `npx tsx <slug>/ui.spec.ts 2>&1 \| tee ...` | 同上 |
| `run-e2e-tests` Error table | `npx playwright install chromium` | 与 Step 1 重复，无统一入口 |
| `gen-test-scripts` Step 5 | `cd tests/e2e && npm install` | 同上 |
| `gen-test-scripts` Step 4 | `grep -r '// VERIFY:' tests/e2e/<feature>/` | 退出码语义不明确，agent 可能忽略失败 |
| `fix-bug` Step 3b / Step 5 | `npx tsx <spec-file> 2>&1` | 直接调用 tsx，绕过 justfile 抽象 |
| `fix-e2e` template | 无明确重跑命令 | agent 不知道修复后用什么命令验证 |

#### B. 单元测试 / 构建命令（散落在 agents 和 skills 中）

| 位置 | 原始命令 | 问题 |
|------|---------|------|
| `task-executor` Step 3 | `go build ./... && go vet ./... && go test -race -cover ./...` | 语言特定，agent 需判断项目类型 |
| `task-executor` Step 3 | `npm run build && npm test` | 同上 |
| `task-executor` Step 3 | `pytest --cov` | 同上 |
| `error-fixer` Step 4 | 同上三条 | 同上 |
| `fix-bug` Step 2 / Step 5 | `<project-test-command>` | 占位符，agent 需自行推断实际命令 |
| `run-tasks` Step 5 | `# go test ./... OR npm test OR the testCommand from index.json` | 注释形式，不是可执行命令 |
| `record-task` Metrics | `go test -cover ./...` / `npm test -- --coverage` / `pytest --cov` | 三选一，agent 需判断 |

### Urgency

每次 agent 执行任务，都在重新解决"如何调用"的问题。`just test` 和 `just build` 已在标准契约中定义，但没有任何 skill 引用它们——这意味着 justfile 的存在对 agent 行为毫无影响。命令散落在 10+ 个文件中，一旦工具链变化需要同步修改多处。

## Proposed Solution

扩展 justfile 标准目标契约，新增两个 e2e 专用目标，然后将所有 skill/agent 中的原始命令替换为对应的 just 命令。

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

| 原始命令 | 替换为 |
|---------|--------|
| `cd tests/e2e && npm install` | `just e2e-setup` |
| `npx playwright install chromium` | （合并到 `just e2e-setup`，不再单独调用） |
| `npx tsx <slug>/cli.spec.ts 2>&1 \| tee ...` | `just test-e2e --feature <slug>` |
| `npx tsx <slug>/api.spec.ts 2>&1 \| tee ...` | （同上，test-e2e 统一执行所有 spec） |
| `npx tsx <slug>/ui.spec.ts 2>&1 \| tee ...` | （同上） |
| `npx tsx <spec-file> 2>&1`（fix-bug） | `just test-e2e --feature <slug>` |
| `grep -r '// VERIFY:' tests/e2e/<feature>/` | `just e2e-verify --feature <slug>` |

#### 单元测试 / 构建命令

| 原始命令 | 替换为 |
|---------|--------|
| `go test -race ./...` / `npm test` / `pytest` | `just test` |
| `go build ./...` / `npm run build` | `just build` |
| `go build ./... && go vet ./... && go test -race -cover ./...` | `just build && just test` |
| `npm run build && npm test` | `just build && just test` |
| `<project-test-command>` | `just test` |
| `go test -cover ./changed/package/...`（record-task） | `just test` |

### 各文件变更摘要

**`init-justfile` command**：
- 标准目标契约表新增 `e2e-setup`、`e2e-verify` 两行及对应 recipe 模板

**`gen-test-scripts` SKILL.md**：
- Step 4 post-generation check：`grep -r '// VERIFY:'` → `just e2e-verify --feature <slug>`，标注 exit 1 = skill incomplete，不得进入 run-e2e-tests
- Step 5 deps install：`cd tests/e2e && npm install` → `just e2e-setup`

**`run-e2e-tests` SKILL.md**：
- Step 1 setup：`cd tests/e2e && npm install` + `npx playwright install chromium` → `just e2e-setup`
- Step 2 run specs：三条独立 `npx tsx` 命令 → `just test-e2e --feature <slug>`
- Error handling 表：`npx playwright install chromium` / `npm install` → `just e2e-setup`

**`fix-bug` command**：
- Step 2 / Step 5 的 `<project-test-command>` → `just test`
- Step 3b / Step 5 的 `npx tsx <spec-file> 2>&1` → `just test-e2e --feature <slug>`

**`run-tasks` command**：
- Step 5 Breaking Gate 的注释命令 → `just test`

**`task-executor` agent**：
- Step 3 Full Verification 的语言示例 → `just build && just test`

**`error-fixer` agent**：
- Step 4 Verify 的语言示例 → `just build && just test`

**`record-task` SKILL.md**：
- Metrics Collection 的语言示例 → `just test`

**`breakdown-tasks` 模板 `run-e2e-tests.md`**：
- Implementation Notes 明确写 `just test-e2e --feature <slug>`

**`breakdown-tasks` 模板 `gen-test-scripts.md`**：
- Implementation Notes 明确写 `just e2e-verify --feature <slug>`，exit 1 则任务未完成

**`breakdown-tasks` 模板 `fix-e2e.md`**：
- 修复后验证步骤明确写 `just test-e2e --feature <slug>`

## Alternatives Considered

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | 零改动 | 原始命令继续散落 10+ 文件；justfile 形同虚设；agent 每次重新推断工具链 | Rejected: 维护成本持续累积 |
| 只抽象 e2e 命令 | 解决最高频问题 | 单元测试/构建命令仍是占位符；task-executor 和 error-fixer 仍需判断语言 | Rejected: 只解决一半 |
| 只抽象单元测试/构建命令 | 解决 agent 判断语言的问题 | e2e 的 VERIFY 检查和 deps 安装仍是原始命令 | Rejected: 只解决一半 |
| 全面替换：新增 e2e-setup + e2e-verify，所有 skill/agent 引用 just 命令 | 命令有唯一权威来源；skill 只需知道目标名称；工具链变化只改 justfile | 需修改 11 个文件 | **Selected** |

## Scope

### In Scope

- `init-justfile` command：标准目标契约新增 `e2e-setup`、`e2e-verify` 两个目标及模板实现
- `gen-test-scripts` SKILL.md：Step 4 VERIFY 检查 + Step 5 deps install → just 命令
- `run-e2e-tests` SKILL.md：Step 1 setup + Step 2 run specs + Error table → just 命令
- `fix-bug` command：`<project-test-command>` → `just test`；`npx tsx` → `just test-e2e --feature <slug>`
- `run-tasks` command：Breaking Gate 注释命令 → `just test`
- `task-executor` agent：Step 3 语言示例 → `just build && just test`
- `error-fixer` agent：Step 4 语言示例 → `just build && just test`
- `record-task` SKILL.md：Metrics Collection 语言示例 → `just test`
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
| task-executor / error-fixer 改为 `just build && just test` 后，覆盖率数据来源不明确 | Low | **Low** — record-task 的 coverage 字段可能为 -1.0 | `just test` 的输出包含覆盖率；record-task 从 `just test` 输出解析，与原来从语言命令输出解析等价 |
| agent 在 fix-e2e 任务中仍使用原始命令（习惯路径依赖） | Medium | **Low** — 功能正确但绕过抽象层 | task 模板 Implementation Notes 明确写 just 命令，不留歧义 |

## Success Criteria

- [ ] `init-justfile` command 的标准目标契约表包含 `e2e-setup` 和 `e2e-verify` 两行，且模板中有对应 recipe 实现（验证：`grep -c 'e2e-setup\|e2e-verify' plugins/forge/commands/init-justfile.md` >= 4）
- [ ] `gen-test-scripts` SKILL.md 中不再出现 `npx playwright install` 或 `cd tests/e2e && npm install`，且出现 `just e2e-verify` 和 `just e2e-setup`（验证：各 grep 结果符合预期）
- [ ] `run-e2e-tests` SKILL.md 中不再出现 `npx tsx` 或 `npx playwright install`，且出现 `just e2e-setup` 和 `just test-e2e --feature`
- [ ] `fix-bug` command 中不再出现 `<project-test-command>` 或 `npx tsx <spec-file>`，替换为 `just test` 和 `just test-e2e --feature <slug>`
- [ ] `task-executor` agent 和 `error-fixer` agent 的 Verification 步骤中出现 `just build && just test`，不再出现语言特定命令示例
- [ ] `record-task` SKILL.md 的 Metrics Collection 示例中出现 `just test`，不再出现三种语言命令
- [ ] `breakdown-tasks` 三个 e2e 相关模板（run-e2e-tests.md、gen-test-scripts.md、fix-e2e.md）的 Implementation Notes 中均出现对应 just 命令（验证：各文件 `grep -c 'just ' <template>` >= 1）

## Next Steps

- Proceed to `/write-prd` to formalize requirements
