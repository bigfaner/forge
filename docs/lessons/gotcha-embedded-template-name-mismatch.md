---
created: "2026-05-27"
tags: [architecture, testing, local-dev-deployment]
---

# Embedded 模板文件名与 type 参数不匹配导致静默失败

## Problem

`forge task add --type coding.fix` 生成的任务 ID 为 `disc-1` 而非预期的 `fix-1`，导致 fix task 在 `forge task claim` 排序中因非数值 ID 排在业务任务之后，无法优先执行。

## Root Cause

因果链（3 层）：

1. **表面现象**：fix task ID 为 `disc-1` 而非 `fix-1`，因为 `IDPrefix` 未从模板默认值注入
2. **直接原因**：`tmpl.Get("coding.fix")` 在 embedded FS 中查找 `data/coding.fix.md`（点号），但已安装的 binary 中嵌入的是 `data/coding-fix.md`（连字符），`Get` 返回 error，整个模板分支被跳过，`IDPrefix` 保持空值 fallback 到 `"disc"`
3. **根因**：已安装的 `forge` binary 是旧版本，编译时模板文件名使用连字符；源码已将文件名改为点号（与 `--type` 参数值 `coding.fix` 对齐），但 binary 未重新编译安装

## Solution

重新编译安装 binary：`cd forge-cli && go install .`

## Reusable Pattern

当使用 Go `//go:embed` 嵌入资源文件时：

1. **文件名是隐式 API**：嵌入文件的路径（含文件名）必须与代码中的查找路径精确匹配。重命名嵌入文件后，必须重新编译 binary 才能生效
2. **静默失败是危险信号**：`Get(name)` 返回 error 时代码直接跳过模板分支（`if err == nil { ... }`），不打印 warning。这使得文件名不匹配的问题在运行时完全不可见
3. **CI 验证不能覆盖嵌入资源匹配**：编译测试能通过（新代码编译无错），但不会验证旧 binary 的嵌入资源是否与当前 type 参数匹配
4. **Binary 版本与源码版本必须同步**：修改 `pkg/template/data/` 下的文件名后，必须立即重新编译安装

## Example

```go
// template.go
var templateDefaults = map[string]Defaults{
    "coding.fix": {
        IDPrefix: "fix",  // 只有 tmpl.Get("coding.fix") 成功才会生效
    },
}

// 如果 embedded FS 中文件名是 coding-fix.md（连字符），
// tmpl.Get("coding.fix") 会失败，IDPrefix 永远不会被设置
```

## Related Files

- `forge-cli/pkg/template/template.go` — 模板注册表和 Get/GetDefaults 函数
- `forge-cli/pkg/template/data/coding.fix.md` — 模板文件（当前使用点号）
- `forge-cli/pkg/task/add.go` — generateAutoID 使用 IDPrefix 生成任务 ID
- `forge-cli/internal/cmd/task/claim.go` — claim 排序逻辑（Priority > compareVersionIDs）

## References

- 相关 lesson: `gotcha-fix-task-claim-priority.md`（fix task claim 排序问题）
