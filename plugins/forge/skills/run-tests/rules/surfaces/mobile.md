# Surface: mobile — 移动端端到端测试编排

本规则文件定义 run-tests skill 对 mobile surface 的移动端端到端测试编排序列。消费方为 SKILL.md 调度器。

测试类型术语定义参见 `docs/reference/test-type-model.md`。

## 编排序列

| 步骤 | just 配方 | 退出码 0 | 退出码 1 | 退出码 2 | 后续动作 |
|------|----------|---------|---------|---------|---------|
| test-setup | `just mobile-test-setup` | 模拟器就绪，测试环境准备完成 | 模拟器启动失败或环境不可用 | — | 进入 dev |
| dev | `just mobile-dev` | 模拟器运行，应用部署就绪 | 启动失败（模拟器不可用） | — | 进入 probe |
| probe | `just mobile-probe` | Appium 健康检查通过 | Appium 无响应 | — | 进入 test |
| test | `just mobile-test` | 移动端端到端测试通过 | 移动端端到端测试失败 | 测试环境异常（需重试） | 进入 teardown |
| teardown | `just mobile-teardown` | 清理完成（模拟器停止、进程清理） | 清理失败（残留模拟器/进程） | — | 结束 |

## Probe 重试策略

- 最多重试 3 次，间隔 5 秒
- 3 次均失败视为退出码 1（retryable）

## 失败处理

### test-setup 失败

test-setup 失败时直接退出，不执行后续步骤。以 test-setup 的退出码退出。

### dev 失败

dev 退出非零时**不继续**后续步骤，直接执行 teardown 并以 dev 的退出码退出。

### probe 失败（HARD-GATE）

<HARD-GATE>
probe 失败后，在同一编排周期内：
- **禁止**重试 probe（重试由 probe 重试策略在上限内处理，非周期级重试）
- **禁止**重启 dev
- 必须执行 teardown 后退出
</HARD-GATE>

probe 最终失败后：
- 退出码 1（retryable）：执行 teardown，以 exit 1 退出
- 退出码 2（blocking）：执行 teardown，以 exit 2 退出

### test 失败

- 退出码 1：执行 teardown，以 exit 1 退出
- 退出码 2（retryable）：执行 teardown，提示用户 "测试环境异常，建议重试"，以 exit 2 退出

### teardown 失败

teardown 失败时记录错误，保留 `.forge/test-state.json` 用于恢复。以当前步骤的退出码退出。

## Suite 名称

测试报告 suite 名称使用 `mobile-e2e/<journey-name>` 格式。

## Journey 过滤

| 标签 | 匹配规则 |
|------|---------|
| `@mobile` | 精确匹配 |

## Per-Journey 执行

Mobile surface 的 dev/probe 生命周期包裹所有 journey 测试：

```
just mobile-test-setup
just mobile-dev
just mobile-probe (with retry)
for each journey in JOURNEYS:
    just mobile-test <journey>
    record results
    on failure: just mobile-teardown, exit
just mobile-teardown
```

test-setup、dev 和 probe 执行一次，per-journey 循环 test，teardown 执行一次。测试配方调用格式为 `just mobile-test <journey>`，其中 `<journey>` 是从 `docs/features/<slug>/testing/` 发现的目录名。
