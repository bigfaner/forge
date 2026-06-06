---
created: "2026-06-06"
tags: [testing, architecture]
---

# gen-test-scripts 生成的测试未遵循项目 surface-first 目录约定

## Problem

项目按 `docs/proposals/e2e-test-conventions/proposal.md` 提案实施了 surface-first 目录结构：

```
tests/api/<feature>/    # API Functional Tests — Vitest
tests/web/<feature>/    # Web E2E Tests — Playwright
tests/infra/            # 静态分析 & 构建检查
```

Convention 文件 `docs/conventions/testing/api/core.md` 和 `docs/conventions/testing/index.md` 都明确写了 `tests/api/<feature>/`。

但 Forge 的 `gen-test-scripts` 生成的测试（如 `item-deletion`、`task-status-transition`、`sub-item-move` 等）放到了 `tests/<journey>/` 根目录，绕过了 surface-first 结构。API 测试文件带有 `_api` 后缀但位于非 surface 目录中。

## Root Cause

1. **L1**: 生成的测试被放入 `tests/<journey>/` 而非 `tests/<surfaceKey>/<journey>/`
2. **L2**: Forge `gen-test-scripts` SKILL.md 第 243 行硬编码了输出目录为 `tests/<journey>/`，且有 HARD-RULE 强制这一行为（第 263 行："Tests go directly to `tests/<journey>/`, NOT to `tests/e2e/features/`"）。当前逻辑不考虑多 surface 场景
3. **L3**: `gen-test-scripts` Step 0 加载 convention 文件仅用于框架选型（imports、assertions、helpers），不读取 convention 中的目录约定。输出路径由 skill 自身规则决定，不受项目 convention 文件控制

## Solution

`gen-test-scripts` 的输出目录应根据 surface 数量自适应：

- **多 surface**（如 `backend=api`, `frontend=web`）：输出到 `tests/<surfaceKey>/<journey>/`
- **单 surface**（如 `app=tui`）：输出到 `tests/<journey>/`

当前 workaround 是生成后手动移动到正确的 surface 目录：

```bash
# 多 surface 项目：移到对应 surface key 下
mv tests/<journey> tests/<surfaceKey>/<journey>
```

## Solution

生成测试后手动移动到正确的 surface 目录：

```bash
# 将 Forge 生成的 API 测试移到正确位置
mv tests/<journey>/*_api*.spec.ts tests/api/<journey>/
# 或整体移动整个 journey 目录到 surface 下
mv tests/<journey> tests/api/<journey>
```

并在 `tests/api/` 中更新引用文件（如有）。

## Reusable Pattern

`gen-test-scripts` 正确的输出目录规则应为：

- **多 surface 项目**（`forge surfaces` 输出多行）：`tests/<surfaceKey>/<journey>/`，其中 `surfaceKey` 是 named surface 的 key（如 `backend`、`frontend`）
- **单 surface 项目**（`forge surfaces` 输出一行）：`tests/<journey>/`，无 surface 层级

判断方法：`forge surfaces` 文本模式下输出行数 > 1 即为多 surface。

当前 workaround：生成后检查 `tests/` 根目录是否出现新的 journey 目录，手动移入正确的 surface 子目录。

## References

- Forge `gen-test-scripts` SKILL.md 第 241-264 行（Output Directory 部分）
- Forge `hooks/guide.md` 第 77 行（"Test files go to `tests/<journey>/` regardless of surface type"）
- 项目提案：`docs/proposals/e2e-test-conventions/proposal.md`
- 项目 Convention：`docs/conventions/testing/api/core.md` 第 12 行
