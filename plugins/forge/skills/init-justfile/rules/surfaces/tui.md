# Surface: tui

> **测试类型参考**：TUI surface 的测试类型为 **终端功能测试（Terminal Functional Test）**，通过子进程 + stdin pipe 验证终端渲染输出 + 交互响应序列。

## 编排序列

| 步骤 | 退出码 0 | 退出码 1 | 退出码 2 | 后续动作 |
|------|---------|---------|---------|---------|
| test | 测试通过 | 测试失败 | 测试环境异常（需重试） | 进入 teardown |
| teardown | 清理完成 | 清理失败（残留进程） | — | 结束 |

注意事项：
- **无 dev 步骤**：TUI surface 不启动持久化服务
- **无 probe 步骤**：TUI 应用无需 HTTP 健康检查
- **无聚合配方**：TUI surface 不生成 `tui` 聚合配方
- test 退出码 2 允许重跑，skill 应提示用户 "测试环境异常，建议重试"

## 配方调用契约

| 配方名 | just 签名 | 退出码 0 语义 | 退出码 1 语义 |
|--------|----------|--------------|--------------|
| tui-test | `just tui-test [journey]` | 所有终端功能测试通过 | 至少一个测试失败 |
| tui-teardown | `just tui-teardown` | 清理完成 | 清理失败 |

实现约束：
- 每个配方必须支持 `[linux]` 和 `[windows]` 双平台变体
- `tui-teardown` 必须用 `just --dry-run` 验证语法
- **不生成** `tui-dev`、`tui-probe` 或 `tui` 聚合配方

## journey 过滤策略

| journey 标签 | 匹配规则 | 说明 |
|-------------|---------|------|
| `@tui` | 精确匹配 | tui surface 的专用 journey |
| 其他 | 忽略 | 非 tui 相关 journey 不由本规则处理 |

## 配方模板（双平台）

```just
# Run terminal functional tests (optionally filter by journey)
# user-customized
tui-test journey='':
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement tui-test" >&2; exit 1

# user-customized
tui-test journey='':
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement tui-test" >&2; exit 1


# Clean up TUI test artifacts
# user-customized
tui-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement tui-teardown" >&2; exit 1

# user-customized
tui-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement tui-teardown" >&2; exit 1
```

**LLM 指令**：将 TODO 桩替换为从语言模板和 Convention 知识推导出的实际命令。上述桩代码展示了所需的配方结构和双平台属性模式。**不生成** `tui-dev`、`tui-probe` 或 `tui` 聚合配方。
