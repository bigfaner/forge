# Surface: api

> **测试类型参考**：API surface 的测试类型为 **API 功能测试（API Functional Test）**，通过 HTTP 客户端验证 HTTP 状态码 + 响应体 JSON + 响应 Header。

## 编排序列

| 步骤 | 退出码 0 | 退出码 1 | 退出码 2 | 后续动作 |
|------|---------|---------|---------|---------|
| dev | API 服务启动成功，等待就绪 | 启动失败（依赖缺失/端口占用） | — | 进入 probe |
| probe | 健康检查通过（GET /healthz 返回 2xx） | 健康检查超时（服务未就绪） | — | 进入 test |
| test | 测试通过 | 测试失败 | 测试环境异常（需重试） | 进入 teardown |
| teardown | 清理完成 | 清理失败（残留进程） | — | 结束 |

注意事项：
- dev 失败时**不继续**后续步骤，直接 teardown 并退出
- probe 最多重试 3 次，间隔 5 秒；3 次均失败视为退出码 1
- test 退出码 2 允许重跑，skill 应提示用户 "测试环境异常，建议重试"

## 配方调用契约

| 配方名 | just 签名 | 退出码 0 语义 | 退出码 1 语义 |
|--------|----------|--------------|--------------|
| api-dev | `just api-dev` | API 服务器就绪，监听端口 | 启动失败，stderr 含错误详情 |
| api-probe | `just api-probe` | HTTP GET /healthz 返回 2xx | 连接拒绝或超时 |
| api-test | `just api-test [journey]` | 所有 API 功能测试通过 | 至少一个测试失败 |
| api-teardown | `just api-teardown` | 进程终止，端口释放 | 进程残留或清理异常 |
| api | `just api` | 聚合配方：dev→probe→test→teardown 完整流程 | 任一子步骤失败 |

实现约束：
- 每个配方必须支持 `[linux]` 和 `[windows]` 双平台变体
- `api` 聚合配方按编排序列顺序调用子配方，遇到非零退出码立即中断
- `api-teardown` 必须用 `just --dry-run` 验证语法

## journey 过滤策略

| journey 标签 | 匹配规则 | 说明 |
|-------------|---------|------|
| `@api` | 精确匹配 | api surface 的专用 journey |
| 其他 | 忽略 | 非 api 相关 journey 不由本规则处理 |

## 配方模板（双平台）

```just
# Start API development server
# user-customized
api-dev:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement api-dev (start API server)" >&2; exit 1

# user-customized
api-dev:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement api-dev (start API server)" >&2; exit 1

# Health check for API server
# user-customized
api-probe:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement api-probe (HTTP GET /healthz)" >&2; exit 1

# user-customized
api-probe:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement api-probe (HTTP GET /healthz)" >&2; exit 1

# Run API functional tests (optionally filter by journey)
# user-customized
api-test journey='':
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement api-test" >&2; exit 1

# user-customized
api-test journey='':
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement api-test" >&2; exit 1


# Clean up API test artifacts
# user-customized
api-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement api-teardown" >&2; exit 1

# user-customized
api-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement api-teardown" >&2; exit 1

# api aggregate: dev -> probe -> test -> teardown
api:
    #!/usr/bin/env bash
    set -euo pipefail
    just api-dev && just api-probe && just api-test; rc=$?; just api-teardown; exit $rc
```

**LLM 指令**：将 TODO 桩替换为从语言模板和 Convention 知识推导出的实际命令。上述桩代码展示了所需的配方结构和双平台属性模式。
