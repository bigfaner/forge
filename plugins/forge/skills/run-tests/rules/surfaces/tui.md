# Surface: tui — 终端功能测试编排

本规则文件定义 run-tests skill 对 tui surface 的终端功能测试编排序列。消费方为 SKILL.md 调度器。

测试类型术语定义参见 `docs/reference/test-type-model.md`。

## 编排序列

| 步骤 | just 配方 | 退出码 0 | 退出码 1 | 退出码 2 | 后续动作 |
|------|----------|---------|---------|---------|---------|
| test | `just tui-test` | 终端功能测试通过 | 终端功能测试失败 | 测试环境异常（需重试） | 进入 teardown |
| teardown | `just tui-teardown` | 清理完成 | 清理失败 | — | 结束 |

注意事项：
- **无 dev 步骤**：TUI surface 不启动持久化服务
- **无 probe 步骤**：TUI 应用无需 HTTP 健康检查
- **无聚合配方**：TUI surface 不执行 `just tui` 聚合配方

## 失败处理

### test 失败

- 退出码 1：执行 teardown，以 exit 1 退出
- 退出码 2（retryable）：执行 teardown，提示用户 "测试环境异常，建议重试"，以 exit 2 退出

### teardown 失败

teardown 失败时记录错误，保留 `.forge/test-state.json` 用于恢复。以当前步骤的退出码退出。

## Suite 名称

测试报告 suite 名称使用 `tui-functional/<journey-name>` 格式。

## Journey 过滤

| 标签 | 匹配规则 |
|------|---------|
| `@tui` | 精确匹配 |

## Per-Journey 执行

TUI surface 的 test 步骤按 journey 逐个执行：

```
for each journey in JOURNEYS:
    just tui-test <journey>
    record results
    on failure: just tui-teardown, exit
just tui-teardown
```

测试配方调用格式为 `just tui-test <journey>`，其中 `<journey>` 是从 `docs/features/<slug>/testing/` 发现的目录名。
