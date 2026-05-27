# Surface: web

> **测试类型参考**：Web surface 的测试类型为 **Web 端到端测试（Web E2E Test）**，通过浏览器自动化验证 DOM 元素可见性 + 用户操作响应 + 页面 URL 变更 + 元素属性值。详见 [测试类型模型](../../../../../docs/reference/test-type-model.md)。

## 编排序列

| 步骤 | 退出码 0 | 退出码 1 | 退出码 2 | 后续动作 |
|------|---------|---------|---------|---------|
| dev | 服务启动成功，等待就绪 | 启动失败（依赖缺失/端口占用） | — | 进入 probe |
| probe | 健康检查通过 | 健康检查超时（服务未就绪） | — | 进入 test |
| test | 测试通过 | 测试失败 | 测试环境异常（需重试） | 进入 teardown |
| teardown | 清理完成 | 清理失败（残留进程） | — | 结束 |

注意事项：
- dev 失败时**不继续**后续步骤，直接 teardown 并退出
- probe 最多重试 3 次，间隔 5 秒；3 次均失败视为退出码 1
- test 退出码 2 允许重跑，skill 应提示用户 "测试环境异常，建议重试"

## 配方调用契约

| 配方名 | just 签名 | 退出码 0 语义 | 退出码 1 语义 |
|--------|----------|--------------|--------------|
| web-dev | `just web-dev` | 开发服务器就绪，监听端口 | 启动失败，stderr 含错误详情 |
| web-probe | `just web-probe` | HTTP 健康检查返回 2xx | 连接拒绝或超时 |
| web-test | `just web-test` | 所有 Web 端到端测试通过 | 至少一个测试失败 |
| web-teardown | `just web-teardown` | 进程终止，端口释放 | 进程残留或清理异常 |
| web | `just web` | 聚合配方：dev→probe→test→teardown 完整流程 | 任一子步骤失败 |

实现约束：
- 每个配方必须支持 `[linux]` 和 `[windows]` 双平台变体
- `web` 聚合配方按编排序列顺序调用子配方，遇到非零退出码立即中断
- `web-teardown` 必须用 `just --dry-run` 验证语法

## journey 过滤策略

| journey 标签 | 匹配规则 | 说明 |
|-------------|---------|------|
| `@web` | 精确匹配 | web surface 的专用 journey |
| `@web-e2e` | 精确匹配 | Web 端到端测试，归入 web surface |
| `@smoke` | 精确匹配 | 冒烟测试，归入 web surface |
| 其他 | 忽略 | 非 web 相关 journey 不由本规则处理 |

## 配方模板（双平台）

```just
# Start web development server
# user-customized
web-dev:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement web-dev (start web dev server)" >&2; exit 1

# user-customized
web-dev:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement web-dev (start web dev server)" >&2; exit 1

# Health check for web server
# user-customized
web-probe:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement web-probe (HTTP health check)" >&2; exit 1

# user-customized
web-probe:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement web-probe (HTTP health check)" >&2; exit 1

# Run Web E2E tests
# user-customized
web-test:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement web-test" >&2; exit 1

# user-customized
web-test:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement web-test" >&2; exit 1


# Clean up web test artifacts
# user-customized
web-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement web-teardown" >&2; exit 1

# user-customized
web-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement web-teardown" >&2; exit 1

# web aggregate: dev -> probe -> test -> teardown
web:
    #!/usr/bin/env bash
    set -euo pipefail
    just web-dev && just web-probe && just web-test; rc=$?; just web-teardown; exit $rc
```

**LLM 指令**：将 TODO 桩替换为从语言模板和 Convention 知识推导出的实际命令。上述桩代码展示了所需的配方结构和双平台属性模式。
