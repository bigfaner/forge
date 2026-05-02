# Hook Command Running Unbounded Test Suite Hangs Indefinitely

## Problem

Stop hook (`task all-completed`) 在所有任务完成后触发，被 Claude Code 在 600s 后强制取消。

症状：transcript 中出现 `hook_cancelled`，`durationMs: 600087`。

## Root Cause

**核心问题：Stop hook 在不满足前置条件时仍然执行了完整测试套件。**

`task all-completed` 的设计意图是"所有任务完成后的健康门禁"：依次执行 `just compile` → `just test` → `just test-e2e`。但它**没有检查测试环境是否就绪**：

- `just test` 需要 Go 和 Node.js — ✅ 环境中有
- `just test-e2e` 遍历所有 spec 文件（API/UI/CLI），API/UI spec 需要运行中的应用服务器 — ❌ hook 没有启动 dev server，也没有检查服务器是否在运行，直接跑测试导致挂死

**触发条件链**：

1. 最后一个 task record 完成写入 `.forge/state.json`（`allCompleted=true`）
2. Claude 停止响应 → 触发 Stop hook → `task all-completed`
3. `task all-completed` 读取状态、清除状态、开始执行测试
4. `just test-e2e` 串行运行 30 个 spec，API/UI spec 的 `before()` 尝试连接应用服务器
5. 无运行中的 dev server → 连接挂死 → 单个 spec 永不退出
6. `runCmdCapture()` 无超时 → 整个 hook 进程被卡死
7. 600s 后 Claude Code 强制 kill

**为什么 `node:test` 没有帮助**：`node:test` 没有默认测试超时。`before()` 钩子中的同步阻塞调用（如启动浏览器）会无限等待。

## Solution

**最终方案：在 hook 中按需启动 dev server，跑完关闭。**

**关键认知**：Playwright 不需要单独"启动浏览器"。测试运行时 Playwright 自动启动和管理浏览器实例。只要 `npx playwright install` 已执行过，浏览器就在本地可用。真正缺失的只是运行中的应用服务器。

**实现方式**：`all_completed.go` 在跑 e2e 前后台启动 dev server，用 `just probe` 等待就绪，跑完再关闭：

```go
// 1. 后台启动 dev server（非阻塞）
devCmd := exec.Command("just", "dev")
devCmd.Dir = projectRoot
devCmd.Start()
defer devCmd.Process.Kill()

// 2. 等待就绪（调用 just probe，避免硬编码健康检查细节）
probeCmd := exec.Command("just", "probe")
probeCmd.Dir = projectRoot
if probeCmd.Run() != nil {
    fmt.Fprintln(os.Stderr, "WARNING: dev server not ready, skipping e2e")
    return
}

// 3. 跑全部 e2e spec
runCmdCapture(projectRoot, "just", "test-e2e")

// 4. defer 自动关闭 dev server
```

justfile 增加两个配方：

```just
# 探测 dev server 是否就绪（轮询健康端点，超时 30s）
probe:
    #!/usr/bin/env bash
    set -euo pipefail
    for i in $(seq 1 30); do
        curl -sf http://localhost:8080/health > /dev/null && exit 0
        sleep 1
    done
    echo "ERROR: dev server not ready after 30s" >&2
    exit 1
```

```just
# 每个 spec 加超时防止挂死
test-e2e:
    #!/usr/bin/env bash
    set -euo pipefail
    fail=0
    for spec in $(find tests/e2e -mindepth 2 -name '*.spec.ts'); do
        timeout 30 npx tsx "$spec" || fail=$((fail+1))
    done
    [ "$fail" -eq 0 ]
```

**基础防护**（即使分级执行，这些也是必须的）：

`runCmdCapture` 加超时防止卡死：
```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
defer cancel()
c := exec.CommandContext(ctx, name, args...)
```

**诊断技巧**：transcript JSONL 中搜 `hook_cancelled`，`durationMs ≈ 600000` 说明子进程未退出。

## Key Takeaway

1. **Hook 可以准备完整测试环境** — 后台启动 dev server（`exec.Start()` + `defer Kill()`）、健康检查、跑测试。Playwright 浏览器无需单独启动，测试运行时自动管理
2. **Hook 中的子命令必须有超时** — `exec.CommandContext` + justfile 中 `timeout` 命令，双重保障
3. **`node:test` 没有默认超时** — 写 e2e spec 时必须在 `test()` 或 `before()` 中显式设置 `timeout`，或在 runner 层面用 `timeout` 命令包裹
4. **添加 `just probe` 配方探测服务是否启动** — 在 justfile 中增加专门的健康检查配方（如轮询 `curl localhost:8080/health`），让 hook 代码和手动调试都能用同一套探测逻辑，避免把健康检查细节硬编码在 Go 代码中
5. **`just test-e2e` 要支持调试模式** — 默认静默运行（适合 hook），加 `-v` 参数时输出每个 spec 的详细日志，方便手动排查失败原因
