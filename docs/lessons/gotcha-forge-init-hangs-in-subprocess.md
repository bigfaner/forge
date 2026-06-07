---
created: "2026-06-07"
tags: [testing, error-handling, local-dev-deployment]
---

# forge init 在子进程中挂起（Windows TUI stdin 检测失效）

## Problem

`forge init` 通过 `exec.Command().CombinedOutput()` 作为子进程调用时，在 Windows 上无限挂起。Go 测试框架在 30 分钟超时后强制终止。影响 `surface-aware-recipe-generation` 和 `test-generation` 两个 journey 的 TC_026 测试用例。

## Root Cause

**Causal chain (3 levels)**:

1. **L1 (症状)**: `forge init` 子进程永不返回，`CombinedOutput()` 阻塞在 `syscall.WaitForSingleObject`
2. **L2 (直接原因)**: `forge init` 的 Step 5 (`runConfigInitIfNeeded`) 和 Step 6 (`runSurfaceConfig`) 使用 `huh` TUI 库（基于 Bubble Tea）进行交互式配置，`huh.NewForm(...).Run()` 阻塞等待 stdin 终端输入
3. **L3 (根因)**: `os.Stdin.Stat()` 的 `os.ModeCharDevice` 检查在 Windows 子进程环境下不可靠——当父进程是 Go 测试进程时，继承的 stdin handle 可能仍被报告为 character device，导致 TUI 守卫条件通过，但实际没有真实终端输入可用

`--skip-just` 只跳过 Step 4（just 安装），**没有 `--non-interactive` 或 `--ci` 标志来跳过 Step 5 和 6**。`huh` 库在 pipe stdin 上不立即报错，而是尝试初始化终端并等待永远不会到来的输入。

## Solution

**短期**（测试端）：在测试中调用 `forge init` 时，设置 `cmd.Stdin = strings.NewReader("")` 强制 stdin 为 pipe，使 `os.ModeCharDevice` 守卫正确检测非交互环境。

**长期**（CLI 端）：为 `forge init` 添加 `--non-interactive` 或 `--ci` 标志，跳过所有 TUI 交互步骤。

## Reusable Pattern

在 Windows 上通过 `exec.Command()` 调用含 TUI 交互的 CLI 工具时：
1. 始终显式设置 `cmd.Stdin`（用 `strings.NewReader("")` 或 `bytes.NewReader(nil)`），**不要**依赖默认的 stdin 继承
2. 不要假设 `os.Stdin.Stat()` + `os.ModeCharDevice` 在 Windows 子进程环境中可靠工作
3. 测试中为每个可能挂起的子进程调用设置 `context.WithTimeout`

## Example

```go
// Bad: stdin 继承自父进程，Windows 上 ModeCharDevice 检测不可靠
cmd := exec.Command(testkit.ForgeBinary, "init", "--skip-just")
cmd.Dir = projectDir
out, err := cmd.CombinedOutput() // may hang forever

// Good: 显式设置 stdin 为空 pipe
cmd := exec.Command(testkit.ForgeBinary, "init", "--skip-just")
cmd.Dir = projectDir
cmd.Stdin = strings.NewReader("") // force non-interactive detection
out, err := cmd.CombinedOutput()
```

## Related Files

- `forge-cli/internal/cmd/init.go` — init 6 步流程编排
- `forge-cli/internal/cmd/init_config.go` — `huh` TUI 配置提示（17 个顺序确认）
- `forge-cli/internal/cmd/init_surfaces.go` — surface 检测 TUI 确认
- `tests/surface-aware-recipe-generation/step2_init_justfile_test.go` — 挂起的测试
- `tests/test-generation/forge_commands_test.go` — 挂起的测试
