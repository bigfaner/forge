# Session-start hook 常见启动错误

## 错误一：permission denied

### Problem

运行 hook 脚本时报 `zsh: permission denied`:

```
% ./hooks/session-start
zsh: permission denied: ./hooks/session-start
```

### Root Cause

Hook 文件缺少可执行权限（execute bit 未设置）。

### Solution

```bash
chmod +x hooks/session-start
chmod +x hooks/run-hook.cmd   # 调试脚本同理
```

验证：

```bash
ls -la hooks/session-start
# 应显示 -rwxr-xr-x
```

---

## 错误二：stdout 污染导致 JSON 解析失败

### Problem

Hook 注册后无效果，或 Claude Code 报 hook 返回非法 JSON。

### Root Cause

Hook 脚本在输出 JSON 的同时，还有 `echo` 写入了 stdout：

```bash
printf '{\n  "additional_context": "%s"\n}\n' "$ctx"
echo 'Loaded zcode guide'   # ← 污染 stdout，破坏 JSON
```

Claude Code 将 hook 的 stdout 整体解析为 JSON，任何非 JSON 文本都会导致解析失败。

### Solution

所有日志/调试输出重定向到 stderr：

```bash
echo 'Loaded zcode guide' >&2
```

---

## Key Takeaway

Hook 脚本的两条铁律：

1. **文件必须有执行权限**：`chmod +x` 是前提，缺失时报 `permission denied`。
2. **stdout 只能输出 JSON**：所有 `echo`/`printf` 日志必须加 `>&2`，否则破坏 JSON 解析。

调试 hook 时，用 `hooks/debug session-start` 同时验证这两点。
