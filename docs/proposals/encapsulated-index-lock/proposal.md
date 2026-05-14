---
created: 2026-05-14
author: faner
status: Draft
---

# Proposal: 封装 index.json 写锁

## Problem

`index.json` 的并发写保护依赖调用方手动配对 `LockFile`/`UnlockFile`。当前只有 `task submit` 一个写入口，但即将新增的 `task add`、`task reorder`、`task block` 等命令都会写入 `index.json`。每个新命令都是潜在的竞态条件。

更深层的问题：**锁的范围是调用方的责任。** 正确的保护必须包裹 LoadIndex → 变更 → SaveIndexAtomic 的完整 read-modify-write 周期。如果未来调用方把 Load 放在锁外（TOCTOU），或在锁内变更但在锁外 Save，保护会被静默绕过。

### Evidence

`internal/cmd/submit.go:88-203` 中，锁和业务逻辑交织 115 行：

```go
lock, err := indexPkg.LockFile(indexPath)    // 手动加锁
if err != nil {
    if errors.Is(err, indexPkg.ErrLockConflict) {
        fmt.Fprintln(os.Stderr, "concurrent write conflict, retry")
        os.Exit(1)
    }
    fmt.Fprintf(os.Stderr, "failed to create lock file: %v\n", err)
    os.Exit(1)
}
defer indexPkg.UnlockFile(lock)              // 手动解锁

idx, err := task.LoadIndex(indexPath)        // read
// ... 80 行变更逻辑 ...
indexPkg.SaveIndexAtomic(indexPath, idx)     // write
```

每个新命令都需要复制这段样板代码，且错误处理必须保持一致。

## Solution

### D1. 新增 `WithLock` 回调函数

在 `pkg/index/lock.go` 中新增：

```go
// WithLock acquires an advisory lock, calls fn, then releases the lock.
// Returns ErrLockConflict if the lock cannot be acquired within the timeout.
func WithLock(indexPath string, fn func() error) error {
    lock, err := LockFile(indexPath)
    if err != nil {
        return err
    }
    defer UnlockFile(lock)
    return fn()
}
```

6 行实现，零外部依赖，零破坏性。

### D2. 为什么不封装 LoadIndex + SaveIndexAtomic

理想方案是：

```go
func UpdateIndex(indexPath string, fn func(*task.TaskIndex) error) error {
    lock, _ := LockFile(indexPath)
    defer UnlockFile(lock)
    idx, _ := task.LoadIndex(indexPath)
    if err := fn(idx); err != nil { return err }
    return SaveIndexAtomic(indexPath, idx)
}
```

但这要求 `pkg/index` 导入 `pkg/task` 的 `TaskIndex` 类型。当前 `pkg/index` 是叶子包（零项目内依赖），引入这个依赖会破坏依赖图的单向性。

`WithLock` 把 LoadIndex/SaveIndexAtomic 留在回调内，保持 `pkg/index` 零依赖。

### D3. 保留 `LockFile`/`UnlockFile` 导出

不立即移除，用于：
- 向后兼容
- 非标准锁模式（如跨多操作的长时间持锁）

用 doc comment 标注推荐使用 `WithLock`。

### D4. 迁移 submit.go

Before → After 对比：

```go
// Before: 手动管理锁
lock, err := indexPkg.LockFile(indexPath)
if err != nil {
    if errors.Is(err, indexPkg.ErrLockConflict) {
        fmt.Fprintln(os.Stderr, "concurrent write conflict, retry")
        os.Exit(1)
    }
    fmt.Fprintf(os.Stderr, "failed to create lock file: %v\n", err)
    os.Exit(1)
}
defer indexPkg.UnlockFile(lock)
idx, err := task.LoadIndex(indexPath)
// ... mutation ...
indexPkg.SaveIndexAtomic(indexPath, idx)

// After: WithLock 回调
err = indexPkg.WithLock(indexPath, func() error {
    idx, err := task.LoadIndex(indexPath)
    if err != nil { return err }
    // ... mutation ...
    return indexPkg.SaveIndexAtomic(indexPath, idx)
})
if err != nil {
    if errors.Is(err, indexPkg.ErrLockConflict) {
        fmt.Fprintln(os.Stderr, "concurrent write conflict, retry")
        os.Exit(1)
    }
    fmt.Fprintf(os.Stderr, "index operation failed: %v\n", err)
    os.Exit(1)
}
```

锁的生命周期从"分散在函数开头和 defer"变为"集中在 WithLock 调用处"。

### D5. 未来命令的标准模式

```go
// task add、task reorder、task block 等命令的标准写法
err = indexPkg.WithLock(indexPath, func() error {
    idx, err := task.LoadIndex(indexPath)
    if err != nil { return err }
    // 命令特定变更
    return indexPkg.SaveIndexAtomic(indexPath, idx)
})
```

一个模式，所有命令。不需要记住 LockFile/UnlockFile 配对。

## Scope

### In Scope

- `pkg/index/lock.go` 新增 `WithLock` 函数
- `internal/cmd/submit.go` 迁移到 `WithLock`
- 新增 `TestWithLock_*` 测试用例
- `LockFile`/`UnlockFile` 添加推荐使用 `WithLock` 的 doc comment

### Out of Scope

- `LockFile`/`UnlockFile` 移除或 unexport（等所有调用方迁移后）
- 封装 LoadIndex/SaveIndexAtomic（受 `pkg/index` 零依赖约束）
- Windows `syscall` → `golang.org/x/sys/windows` 迁移（独立改进）

## Impact

### 改动的文件

| 文件 | 改动 |
|---|---|
| `pkg/index/lock.go` | 新增 `WithLock`（6 行） |
| `internal/cmd/submit.go` | 替换 LockFile/UnlockFile 为 WithLock 回调 |
| `pkg/index/lock_test.go` | 新增 WithLock 测试 |

### 不变的部分

- `LockFile`/`UnlockFile` API 不变（向后兼容）
- `lock_unix.go`/`lock_windows.go` 不变
- `atomic.go` 不变
- 测试策略不变（table-driven）

## Risks

| 风险 | 缓解 |
|---|---|
| 回调内 panic 导致锁泄漏 | `defer UnlockFile(lock)` 确保即使 panic 也会释放 |
| 回调嵌套导致死锁 | 不同 feature 的 indexPath 不同，不会冲突；同 feature 嵌套则是调用方 bug，flock 在同进程内可重入（Unix LOCK_EX 对同 fd 不阻塞） |
| `defer UnlockFile` 的 error 被静默丢弃 | `UnlockFile` 在 `defer` 中调用，其返回值无法传播。若 unlock 失败（fd 无效等极端情况），下一个 `WithLock` 调用会等 5s 超时后返回 `ErrLockConflict`。实践中 flock/LockFileEx 的 unlock 只在 fd 无效时失败，正常路径不会触发。若需防御，可改为 `defer func() { unlockErr = UnlockFile(lock) }()` 将错误提升为返回值，但当前收益不足以覆盖复杂度 |
| 未来需要可配置超时 | 预留 `WithLockTimeout(indexPath, timeout, fn)` 扩展点，当前 5s 对 CLI 足够 |
