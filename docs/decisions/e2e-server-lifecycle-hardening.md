# E2E Server Lifecycle 硬化方案

> 决策日期: 2026-05-10
> 来源: `docs/forensics/e2e-server-lifecycle/report.md` (train-recorder 项目)
> 状态: 已实施

## 1. 问题根因

`/run-e2e-tests` 技能定义了手动服务器管理（`just run` → `just probe` → 测试），但 task-executor 跳过了服务器启动步骤，直接运行 `just test-e2e`。Playwright `webServer` 自动接管但只做 TCP 端口检查，Metro bundle 未构建完成就开始执行测试，导致 30/114 UI 测试假阳性失败。

**根因**：技能指令（"Step 1: do X"）无强制约束力——只要存在捷径（Playwright 自动启动），agent 就会绕过。

## 2. 设计原则

1. **复用已有命令** — 不新增 recipe，把服务器生命周期嵌入 `test-e2e` 内部，利用已有的 `run` 和 `probe`
2. **幂等性** — 三层检测（PID 存活 > probe 通过 > 启动服务），多次调用结果一致
3. **配置冲突前置消除** — 在 `gen-test-scripts` 阶段检测并移除 `playwright.config.ts` 中的 `webServer`
4. **probe 语义放宽** — 从"要求 200"改为"非 5xx 即就绪"

## 3. 改动范围

| 文件 | 改动类型 | 说明 |
|------|---------|------|
| `references/justfile-templates/node.just` | 修改 | `test-e2e` 嵌入生命周期 + `probe` 接受非 5xx |
| `references/justfile-templates/go.just` | 修改 | 同上 |
| `references/justfile-templates/rust.just` | 修改 | 同上 |
| `references/justfile-templates/python.just` | 修改 | 同上 |
| `references/justfile-templates/mixed.just` | 修改 | `test-e2e` 多服务生命周期 + `probe` 接受非 5xx |
| `references/justfile-templates/generic.just` | 无变化 | 存根不变（无 `run` 可调用） |
| `skills/gen-test-scripts/SKILL.md` | 修改 | Step 5 新增 webServer 检测 |
| `skills/run-e2e-tests/SKILL.md` | 修改 | Step 1/6 简化 |
| `skills/gen-test-scripts/templates/playwright.config.ts` | 修改 | 添加 no-webServer 注释 |

**不新增任何 recipe，不修改 `init-justfile` 标准目标合约。**

## 4. `probe` 改造（接受非 5xx）

所有模板共享同一个 `probe` 配方，改动一处：

**Before**:

```bash
if curl -sf --max-time 5 "$url" > /dev/null 2>&1; then
    echo "OK: ${label%%:*} ($url)"
else
    echo "FAIL: ${label%%:*} ($url) not responding" >&2
    fail=$((fail+1))
fi
```

`curl -sf` 仅在 HTTP 2xx 时成功。

**After**:

```bash
STATUS=$(curl -s -o /dev/null -w '%{http_code}' --max-time 5 "$url" 2>/dev/null || echo "000")
STATUS=${STATUS:-000}
if [ "$STATUS" != "000" ] && [ "$STATUS" -lt 500 ]; then
    echo "OK: ${label%%:*} ($url)"
else
    echo "FAIL: ${label%%:*} ($url) status=$STATUS" >&2
    fail=$((fail+1))
fi
```

判定逻辑:

| STATUS | 含义 | 判定 |
|--------|------|------|
| `000` | 连接拒绝 / 无响应 | 未就绪 |
| `100-499` | 服务已启动（含 404） | **就绪** |
| `500-599` | 服务端错误 | 未就绪 |

## 5. `test-e2e` 改造（嵌入服务器生命周期）

### 5.1 执行流程

```
just test-e2e
  │
  ├─ Layer 1: PID 存活？→ skip 启动
  ├─ Layer 2: just probe 通过？→ skip 启动（手动启动场景）
  └─ Layer 3: just run &（启动服务，写 PID）
  │
  ├─ 健康检查: 重试 just probe (10×3s=30s)
  │
  └─ npx playwright test（执行测试）
```

### 5.2 单服务模板（node / go / rust / python）

```just
[arg("feature", long)]
test-e2e feature="":
    #!/usr/bin/env bash
    set -euo pipefail
    # --- Server Lifecycle (idempotent) ---
    mkdir -p tests/e2e/results
    _root="$(pwd)"
    _pid_file="$_root/tests/e2e/results/.pid-server"
    should_start=false
    # Layer 1: tracked process alive?
    if [ -f "$_pid_file" ] && kill -0 "$(cat "$_pid_file")" 2>/dev/null; then
        should_start=false
    # Layer 2: already responding (manually started)?
    elif just probe > /dev/null 2>&1; then
        should_start=false
    # Layer 3: start
    else
        should_start=true
    fi
    if [ "$should_start" = true ]; then
        just run > /dev/null 2>&1 &
        echo "$!" > "$_pid_file"
        # trap: cleanup on any exit (crash, signal, success)
        trap 'kill "$(cat "$_pid_file")" 2>/dev/null || true; rm -f "$_pid_file"' EXIT
        # Early crash detection: avoid wasting 30s if server fails immediately
        sleep 1
        if ! kill -0 "$(cat "$_pid_file")" 2>/dev/null; then
            echo "e2e: server process exited immediately" >&2
            rm -f "$_pid_file"
            exit 1
        fi
    fi
    if [ -f tests/e2e/config.yaml ]; then
        ready=false
        for i in $(seq 1 10); do
            if just probe > /dev/null 2>&1; then ready=true; break; fi
            sleep 3
        done
        if [ "$ready" = false ]; then
            echo "e2e: health check failed after 30s" >&2
            just probe || true
            exit 1
        fi
    fi
    # --- Run Tests ---
    if [ "{{feature}}" != "" ]; then
        cd tests/e2e && E2E_FEATURE=1 npx playwright test features/{{feature}}/
    else
        if [ ! -d tests/e2e/node_modules ]; then npm install --prefix tests/e2e; fi
        cd tests/e2e && npx playwright test
    fi
```

### 5.3 多服务模板（mixed.just）

利用 `just probe` 的逐服务输出（`OK: frontend` / `FAIL: backend`），只启动失败的服务。

```just
[arg("feature", long)]
test-e2e feature="":
    #!/usr/bin/env bash
    set -euo pipefail
    # --- Server Lifecycle (idempotent, per-service) ---
    mkdir -p tests/e2e/results
    _root="$(pwd)"
    probe_output=$(just probe 2>&1) || true
    if echo "$probe_output" | grep -q "FAIL:"; then
        for svc in frontend backend; do
            pid_file="$_root/tests/e2e/results/.pid-$svc"
            # Layer 1: tracked process alive?
            if [ -f "$pid_file" ] && kill -0 "$(cat "$pid_file")" 2>/dev/null; then continue; fi
            # Layer 2: probe reported this service OK? → skip
            echo "$probe_output" | grep -q "FAIL: $svc" || continue
            # Layer 3: start only the failed service
            just run "$svc" > /dev/null 2>&1 &
            echo "$!" > "$pid_file"
        done
        # trap: cleanup on any exit (moved outside loop)
        trap 'for f in "$_root"/tests/e2e/results/.pid-*; do [ -f "$f" ] && kill "$(cat "$f")" 2>/dev/null || true; rm -f "$f"; done' EXIT
        # Early crash detection: check if any started service died immediately
        sleep 1
        for pid_f in "$_root"/tests/e2e/results/.pid-*; do
            if [ -f "$pid_f" ] && ! kill -0 "$(cat "$pid_f")" 2>/dev/null; then
                echo "e2e: service process exited immediately ($(basename "$pid_f" | sed 's/.pid-//'))" >&2
                rm -f "$pid_f"
                exit 1
            fi
        done
        # Health check (wait for ALL services, 30s timeout)
        ready=false
        for i in $(seq 1 10); do
            if just probe > /dev/null 2>&1; then ready=true; break; fi
            sleep 3
        done
        if [ "$ready" = false ]; then
            echo "e2e: health check failed after 30s" >&2
            just probe || true
            exit 1
        fi
    elif echo "$probe_output" | grep -q "OK:"; then
        : # all healthy, skip startup
    else
        echo "e2e: unexpected probe output: $probe_output" >&2; exit 1
    fi
    # --- Run Tests ---
    if [ "{{feature}}" != "" ]; then
        cd tests/e2e && E2E_FEATURE=1 npx playwright test features/{{feature}}/
    else
        if [ ! -d tests/e2e/node_modules ]; then npm install --prefix tests/e2e; fi
        cd tests/e2e && npx playwright test
    fi
```

**逐服务启动判定示例**:

```
场景: frontend 手动启动, backend 未启动

just probe 输出:
  OK: frontend (http://localhost:3456)
  FAIL: backend (http://localhost:8080) status=000

解析:
  frontend → grep "FAIL: frontend" 无匹配 → skip ✓
  backend  → grep "FAIL: backend" 命中 → just run backend &

结果: 只启动 backend, frontend 不受影响（无端口冲突）
```

**并行启动时序**（首次启动，Expo 14s + Backend 5s）:

```
T=0   just probe → 全部 FAIL
      just run frontend & → 后台构建 bundle
      just run backend &  → 后台启动 API
T=0   健康检查开始轮询 just probe
T=5   backend 就绪（probe 中 backend curl 秒过，frontend 未就绪 → probe 仍失败）
T=14  frontend 就绪 → probe 全部通过
T=14  npx playwright test 开始

总耗时 = max(14, 5) = 14s
```

### 5.4 幂等性矩阵

| 场景 | Layer 检测 | 启动 | 健康检查 | 清理 | 测试 |
|------|-----------|------|---------|------|------|
| 首次 `just test-e2e` | Layer 1 ✗, Layer 2 ✗ | `just run &` + 写 PID | 轮询 30s | trap EXIT kill | 执行 |
| 再次 `just test-e2e` | Layer 1 ✓ (PID 存活) | skip | `just probe` 秒过 → break | 无 trap (非本进程启动) | 执行 |
| 服务崩溃后 | Layer 1 ✗ (PID 无效), Layer 2 ✗ | 重启 | 重新轮询 | trap EXIT kill | 执行 |
| 手动 `just run` 后 (单服务) | Layer 1 ✗, Layer 2 ✓ (`just probe` 通过) | skip | skip | 无 trap (非本进程启动) | 执行 |
| 手动启动部分服务 (mixed) | probe 报告 OK/FAIL 逐服务 | 只启动 FAIL 的服务 | 轮询至全部 OK | trap EXIT kill (仅启动的) | 执行 |
| CLI-only 项目 | Layer 2 ✓ (无 config.yaml, probe 成功) | skip | skip (无 config) | 无 trap | 执行 |

**mixed 逐服务启动**: `just probe` 输出包含每个服务的独立状态（`OK: frontend` / `FAIL: backend`）。通过 grep 解析，只对 `FAIL` 的服务执行 `just run $svc`。已手动启动的服务不受影响，无端口冲突。

### 5.5 单服务 vs 多服务对比

| | 单服务 | 多服务 |
|---|--------|--------|
| Layer 2 机制 | `just probe` 整体判定 | 解析 `just probe` 逐服务输出，grep `FAIL: $svc` |
| 启动粒度 | 全部或无 | 只启动 FAIL 的服务 |
| 启动命令 | `just run &` | `just run $svc &`（按需） |
| 健康检查 | `just probe` 重试 | `just probe` 重试 |

## 6. gen-test-scripts 技能改动

在 Step 5 (Ensure Shared Infrastructure) 处理 `playwright.config.ts` 时新增冲突检测:

#### Playwright webServer Audit

When processing `playwright.config.ts` (new copy or existing file):

> **Exception to HARD-RULE "do not modify existing playwright.config.ts"**:
> The webServer audit is the ONLY permitted modification to an existing `playwright.config.ts`.
> All other content must be preserved verbatim.

1. Check for `webServer` key:
   ```bash
   grep -n "webServer" tests/e2e/playwright.config.ts
   ```

2. If `webServer` block found:
   - Remove the entire `webServer: { ... }` block (matching braces)
   - Insert comment: `// Server lifecycle managed by justfile (embedded in test-e2e recipe)`
   - Report: "Removed Playwright webServer config — server lifecycle enforced by test-e2e"

3. If `webServer` not found: no action needed

## 7. run-e2e-tests 技能改动

### Step 1 简化

**Before**: 手动 `just run &` + `echo PID` + `for just probe` 循环

**After**:
```markdown
### Step 1: Setup Environment

Run `just e2e-setup` (idempotent):

```bash
just e2e-setup
mkdir -p tests/e2e/results/
mkdir -p tests/e2e/features/<slug>/results/
rm -f tests/e2e/results/test-results.json
```

Server lifecycle is embedded in `test-e2e`.
Calling `just test-e2e` (Step 3) automatically ensures servers are started and healthy.
No manual server management needed.
```

### Step 6 Teardown

**Before**: 读取 PID 文件 + 逐个 kill + 删除文件

**After**:
```markdown
<HARD-RULE>
**Teardown is mandatory**, even if tests fail:

1. Kill tracked servers: `for f in tests/e2e/results/.pid-*; do [ -f "$f" ] && kill "$(cat "$f")" 2>/dev/null || true; rm -f "$f"; done`
2. Clean up temporary files

Playwright browser instances are automatically closed by the test runner.
</HARD-RULE>
```

## 8. playwright.config.ts 模板改动

在模板末尾添加:

```typescript
// IMPORTANT: Do NOT add a webServer section here.
// Server lifecycle is managed by the test-e2e justfile recipe.
// Adding webServer would conflict and cause false positive test failures
// (TCP port ready ≠ application ready).
```

## 9. 迁移指南

1. **重新生成 justfile**: `/init-justfile` 更新标准配方（boundary marker 内替换）
2. **检查 playwright.config.ts**: 如有 `webServer` 配置则移除
3. **验证**: `just --dry-run test-e2e` 确认无语法错误

向后兼容:
- `just probe` 独立调用：行为从"要求 200"放宽为"接受非 5xx"（更宽松，无破坏性）
- `just test-e2e`：现在自动管理服务器（之前需手动 `just run` + `just probe`）
- `just run`、`just dev` 等：不受影响
- 标准目标合约：不变，无新 target

## 10. 防御纵深

```
代码生成阶段                           运行时
┌───────────────────────────┐    ┌───────────────────────────────────────┐
│ gen-test-scripts Step 5:  │    │ just test-e2e                         │
│ 检测 playwright.config.ts │    │   │                                   │
│ 中的 webServer → 移除     │    │   ├─ Layer 1: PID alive? → skip      │
│                           │    │   ├─ Layer 2: just probe? → skip     │
│ 阻止冲突配置进入代码库     │    │   └─ Layer 3: just run & → start     │
│                           │    │   │                                   │
│                           │    │   ├─ 健康检查: just probe 重试 30s   │
│                           │    │   └─ npx playwright test              │
│                           │    │                                       │
│                           │    │ agent 无法跳过:                       │
│                           │    │ test-e2e 内嵌完整生命周期             │
│                           │    │ 调用即保证服务就绪                    │
└───────────────────────────┘    └───────────────────────────────────────┘
       前置消除                         内嵌自防护
```

## 11. 风险评估

| 风险 | 概率 | 缓解措施 |
|------|------|---------|
| `test-e2e` 配方膨胀，维护成本增加 | 低 | 生命周期逻辑 ~20 行 bash，结构固定（三层检测 + trap + 重试），不随项目变化 |
| `just probe` Layer 2 检查增加 ~2s | 低 | 仅在 PID 不存在时触发；手动启动场景下一次 curl 即过 |
| `kill -0` 在 Windows bash 不支持信号检测 | 中 | Layer 1 失效退化为 Layer 2（probe），功能不受影响，仅性能略降 |
| `kill` 无法终止子进程树（Windows） | 低 | `trap EXIT` 确保退出时清理；测试后重启可彻底清理 |
| 健康检查实际超时 ~80s（非 30s） | 低 | 最坏情况 10×(5s curl timeout + 3s sleep)；正常启动远低于此值 |
| agent 直接调用 `npx playwright test` 绕过 justfile | 中 | 依赖 `run-e2e-tests` 技能约束；Playwright 无 `webServer` 时连接失败会暴露问题 |
| 两个并发 `just test-e2e` 竞争启动 | 低 | PID 文件覆盖导致孤立进程；`trap EXIT` 在退出时清理 |
| CLI-only 项目 | 无 | Layer 2 `just probe` 通过（无 config.yaml → probe 成功）→ skip 启动 |
