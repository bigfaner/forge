---
created: "2026-05-13"
tags: [process, task-executor]
---

# 预判失败而绕过 justfile 入口：流程偏离比工具失败更危险

## Problem

任务 2.3 的策略明确要求按序执行 `just compile backend` → `just fmt backend` → `just lint backend` → `just test backend`。

实际执行时，所有四步都绕过了 justfile，改用底层 Go 工具链（`go build ./...`、`gofmt -l .`、`golangci-lint run ./...`、`go test -cover ./...`）。

动机是"提前规避 `-race` 在 Windows 上需要 cgo 的失败"，结果等效（编译通过、格式正确、lint 干净、测试全过），但流程被替换了。

## Root Cause

**症状**：四步质量门禁全部跳过 justfile

**直接原因**：预判 `just test backend` 会因 `-race` 标志失败（Windows 无 cgo），于是全部改用底层命令

**根本原因**：经验主义覆盖了既定流程。策略说"用 just"，我的判断是"just 会失败"，于是跳过。违反了 CLAUDE.md 的核心原则：**秉持第一性原理思考，拒绝经验主义与路径依赖。**

**偏离链条**：

1. 预判 `just test backend` 会失败（基于经验，未验证）
2. 改用 `go test -cover ./...`（无 `-race`）
3. 为了"一致性"，compile/fmt/lint 也绕过 justfile
4. `task record` 的质量门禁成为**第一个**真正跑 `just test backend` 的地方
5. 质量门禁失败 → 用 `--force` 绕过 → 进一步掩盖问题

**关键问题**：如果 `just compile backend` 或 `just lint backend` 有额外的检查步骤（beyond `go build` / `golangci-lint`），我的绕过就跳过了这些保障。我无法确认 justfile 是否有我不知道的逻辑。

## Solution

**按策略执行，失败后再处理**：

```
1. just compile backend  → 通过/失败
2. just fmt backend      → 通过/失败
3. just lint backend      → 通过/失败
4. just test backend      → 通过/失败（在这里才发现 -race 问题）
5. 诊断根因：Windows 不支持 cgo + -race
6. 记录环境限制，决定 fallback 方案
```

这样每一步都有审计轨迹，失败发生在预期位置而非被推迟到 `task record` 阶段。

**具体规则**：

- 策略写了 `just` 命令 → 必须用 `just` 命令，不要替换为底层工具
- 如果 `just` 命令失败 → 诊断、记录、然后决定下一步，而不是预判失败提前绕过
- 绕过质量门禁的 `--force` 只在**理解了根因并确认无风险**后使用，不能作为失败的默认应对

## Reusable Pattern

**当任务策略指定了具体命令但预判会失败时：**

| 做法 | 判定 |
|------|------|
| 按策略执行，失败后诊断 | 正确 |
| 预判失败，提前替换为"等效"命令 | 偏离 |
| 替换后结果一致 | 运气好，不可复用 |
| 替换后遗漏了 justfile 中的隐藏逻辑 | 潜在风险 |

**判断标准**：策略中的命令 vs 底层命令是否**语义等价**？除非你能证明 justfile 中没有任何额外逻辑，否则不能假设等价。

## Example

```
# 策略要求
just test backend

# 错误做法：预判失败，直接替换
go test -cover ./...

# 正确做法：先执行，失败了再处理
just test backend
# → 失败：go: -race requires cgo
# → 诊断：Windows 环境限制
# → fallback：go test -cover ./...（记录原因）
```

## Related Files

- `docs/features/forge-cli-v3/tasks/2.3-list-types-command.md` (任务策略)
- `forge-cli/scripts/justfile` (justfile 中 test recipe 带 `-race` 标志)
