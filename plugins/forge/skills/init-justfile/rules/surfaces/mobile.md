# Surface: mobile

> **测试类型参考**：Mobile surface 的测试类型为 **移动端端到端测试（Mobile E2E Test）**，通过 Maestro YAML / 设备自动化验证 UI 元素可见性 + 用户操作响应 + 屏幕 ID 变更。

## 编排序列

| 步骤 | 退出码 0 | 退出码 1 | 退出码 2 | 后续动作 |
|------|---------|---------|---------|---------|
| test-setup | 模拟器就绪，测试环境准备完成 | 模拟器启动失败或环境不可用 | — | 进入 dev |
| dev | 模拟器运行，应用部署就绪 | 启动失败（模拟器不可用） | — | 进入 probe |
| probe | Appium 健康检查通过 | Appium 无响应 | — | 进入 test |
| test | 测试通过 | 测试失败 | 测试环境异常（需重试） | 进入 teardown |
| teardown | 清理完成 | 清理失败（残留模拟器/进程） | — | 结束 |

注意事项：
- test-setup 负责模拟器准备，是 mobile surface 的前置步骤；test-setup 失败时直接退出，不继续后续步骤
- dev 失败时**不继续**后续步骤，直接 teardown 并退出
- probe 最多重试 3 次，间隔 5 秒；3 次均失败视为退出码 1
- test 退出码 2 允许重跑，skill 应提示用户 "测试环境异常，建议重试"

## 配方调用契约

| 配方名 | just 签名 | 退出码 0 语义 | 退出码 1 语义 |
|--------|----------|--------------|--------------|
| mobile-test-setup | `just mobile-test-setup` | 模拟器就绪，测试环境准备完成 | 模拟器启动失败，stderr 含错误详情 |
| mobile-dev | `just mobile-dev` | 模拟器运行，应用部署就绪 | 启动失败，stderr 含错误详情 |
| mobile-probe | `just mobile-probe` | Appium 健康检查通过 | Appium 无响应 |
| mobile-test | `just mobile-test` | 所有移动端端到端测试通过 | 至少一个测试失败 |
| mobile-teardown | `just mobile-teardown` | 模拟器停止，进程清理完成 | 残留模拟器或清理异常 |
| mobile | `just mobile` | 聚合配方：test-setup→dev→probe→test→teardown 完整流程 | 任一子步骤失败 |

实现约束：
- 每个配方必须支持 `[linux]` 和 `[windows]` 双平台变体
- `mobile` 聚合配方按编排序列顺序调用子配方，遇到非零退出码立即中断
- `mobile-teardown` 必须用 `just --dry-run` 验证语法

## journey 过滤策略

| journey 标签 | 匹配规则 | 说明 |
|-------------|---------|------|
| `@mobile` | 精确匹配 | mobile surface 的专用 journey |
| 其他 | 忽略 | 非 mobile 相关 journey 不由本规则处理 |

## 配方模板（双平台）

```just
# Prepare emulator and test environment for mobile tests
# user-customized
mobile-test-setup:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-test-setup (prepare emulator)" >&2; exit 1

# user-customized
mobile-test-setup:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-test-setup (prepare emulator)" >&2; exit 1

# Start emulator and deploy app
# user-customized
mobile-dev:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-dev (start emulator + deploy app)" >&2; exit 1

# user-customized
mobile-dev:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-dev (start emulator + deploy app)" >&2; exit 1

# Health check for Appium
# user-customized
mobile-probe:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-probe (Appium health check)" >&2; exit 1

# user-customized
mobile-probe:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-probe (Appium health check)" >&2; exit 1

# Run Mobile E2E tests
# user-customized
mobile-test:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-test" >&2; exit 1

# user-customized
mobile-test:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-test" >&2; exit 1


# Clean up mobile test artifacts
# user-customized
mobile-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-teardown" >&2; exit 1

# user-customized
mobile-teardown:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "TODO: implement mobile-teardown" >&2; exit 1

# mobile aggregate: test-setup -> dev -> probe -> test -> teardown
mobile:
    #!/usr/bin/env bash
    set -euo pipefail
    just mobile-test-setup && just mobile-dev && just mobile-probe && just mobile-test; rc=$?; just mobile-teardown; exit $rc
```

**LLM 指令**：将 TODO 桩替换为从语言模板和 Convention 知识推导出的实际命令。上述桩代码展示了所需的配方结构和双平台属性模式。
