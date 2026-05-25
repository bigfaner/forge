# Surface: api — Run-Tests Orchestration

本规则文件定义 run-tests skill 对 api surface 的编排序列。消费方为 SKILL.md 调度器。

## 编排序列

| 步骤 | just 配方 | 退出码 0 | 退出码 1 | 退出码 2 | 后续动作 |
|------|----------|---------|---------|---------|---------|
| dev | `just api-dev` | API 服务启动成功，等待就绪 | 启动失败（依赖缺失/端口占用） | — | 进入 probe |
| probe | `just api-probe` | 健康检查通过（GET /healthz 返回 2xx） | 健康检查超时（服务未就绪） | — | 进入 test |
| test | `just api-test` | 测试通过 | 测试失败 | 测试环境异常（需重试） | 进入 teardown |
| teardown | `just api-teardown` | 清理完成 | 清理失败（残留进程） | — | 结束 |

## Probe 重试策略

- 最多重试 3 次，间隔 5 秒
- 3 次均失败视为退出码 1（retryable）

## 失败处理

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

## Journey 过滤

| 标签 | 匹配规则 |
|------|---------|
| `@api` | 精确匹配 |
