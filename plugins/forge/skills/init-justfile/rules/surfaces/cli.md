# Surface: cli

> **测试类型参考**：CLI surface 的测试类型为 **CLI 功能测试（CLI Functional Test）**，通过子进程执行验证进程退出码 + stdout/stderr 输出。详见 [测试类型模型](../../../../../docs/reference/test-type-model.md)。

## 编排序列

| 步骤 | 退出码 0 | 退出码 1 | 退出码 2 | 后续动作 |
|------|---------|---------|---------|---------|
| test | 测试通过 | 测试失败 | 测试环境异常（需重试） | 进入 teardown |
| teardown | 清理完成 | 清理失败（残留进程） | — | 结束 |

注意事项：
- **无 dev 步骤**：CLI surface 不启动持久化服务
- **无 probe 步骤**：CLI 工具无需 HTTP 健康检查
- **无聚合配方**：CLI surface 不生成 `cli` 聚合配方
- test 退出码 2 允许重跑，skill 应提示用户 "测试环境异常，建议重试"

## 配方调用契约

| 配方名 | just 签名 | 退出码 0 语义 | 退出码 1 语义 |
|--------|----------|--------------|--------------|
| cli-test | `just cli-test` | 所有 CLI 功能测试通过 | 至少一个测试失败 |
| cli-teardown | `just cli-teardown` | 清理完成 | 清理失败 |

实现约束：
- 每个配方必须支持 `[linux]` 和 `[windows]` 双平台变体
- `cli-teardown` 必须用 `just --dry-run` 验证语法
- **不生成** `cli-dev`、`cli-probe` 或 `cli` 聚合配方

## journey 过滤策略

| journey 标签 | 匹配规则 | 说明 |
|-------------|---------|------|
| `@cli` | 精确匹配 | cli surface 的专用 journey |
| 其他 | 忽略 | 非 cli 相关 journey 不由本规则处理 |

## 配方模板（双平台）

```just
# Run CLI functional tests
# user-customized
cli-test:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement cli-test" >&2; exit 1

# user-customized
cli-test:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement cli-test" >&2; exit 1

# DEPRECATED: removed after v3.2.0 — use cli-test instead
alias test-e2e := cli-test

# Clean up CLI test artifacts
# user-customized
cli-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement cli-teardown" >&2; exit 1

# user-customized
cli-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement cli-teardown" >&2; exit 1
```

**LLM 指令**：将 TODO 桩替换为从语言模板和 Convention 知识推导出的实际命令。上述桩代码展示了所需的配方结构和双平台属性模式。**不生成** `cli-dev`、`cli-probe` 或 `cli` 聚合配方。
