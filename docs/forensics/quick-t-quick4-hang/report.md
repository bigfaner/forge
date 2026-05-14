---
date: 2026-05-14
session: dfd6065a-2dc1-45c7-ac50-a720d390bc6e
model: glm-5.1
branch: task-stage-gates
trigger: "用户报告 agent 一直在初始化中"
---

# Forensic Report: Agent 初始化慢

## 症状

用户在 `/quick` pipeline 执行过程中观察到 subagent 长时间"在初始化中"——dispatch 后无可见输出，最终在 T-quick-4 处中断。

## 证据

### 会话概览

- **会话**: `dfd6065a` (2026-05-14 18:07 ~ 20:50, 共 2.7 小时)
- **模型**: glm-5.1
- **Pipeline**: `/quick` → quick-tasks → run-tasks (7 个 agent)

### Agent 执行时间

| Agent | 任务 | 耗时 | 工具调用数 | 工具执行时间 | 推理占比 |
|-------|------|------|-----------|-------------|---------|
| Explore | 代码结构探索 | ~2 min | — | — | — |
| Task 1 | stage-gates 实现 | 8.2 min | — | <1s | ~99% |
| Task 2 | skill 文档更新 | 6.8 min | — | <1s | ~99% |
| T-quick-1 | 生成测试用例 | 3.8 min | — | <1s | ~99% |
| T-quick-2 | 生成测试脚本 | 11.5 min | — | <1s | ~99% |
| T-quick-3 | 执行 e2e 测试 | **53.6 min** | 37 | **8.5s** | **99.97%** |
| T-quick-4 | 毕业测试脚本 | 中断 | 0 | 0 | — |
| **合计** | | **~85 min** | | **~10s** | |

### T-quick-3 详细分析 (最极端案例)

- 会话时长: 53.6 分钟 (3217.6s)
- 工具执行: 8.5s (29 次 Bash + 3 次 Read + 3 次 Skill + 1 次 Write)
- 推理时间: 3209.1s
- **每两次工具调用之间平均推理耗时: ~86 秒**
- Transcript 大小: 191KB, 95 条消息

### T-quick-4 (被中断)

- Transcript 大小: 1KB, 2 条消息
- 工具调用: 0
- 被用户在 dispatch 后立即中断

## 根因分析

### 因果链 (3 层)

1. **症状**: Subagent 看起来"卡在初始化"，长时间无输出
2. **直接原因**: 每个 agent 从 dispatch 到第一个工具调用之间有 1-2 分钟的延迟，后续每个工具调用之间也有 ~86s 延迟
3. **根因**: 模型 (glm-5.1) 推理延迟极高。Agent 执行时间的 95-99% 是模型推理 (thinking token generation)，只有 1-5% 是实际工具执行

### 时间分布图

```
|======== 推理 (99%) ========|工具|======== 推理 (99%) ========|工具|== 推理 ==|
0                          86s  0.2s                          86s  0.3s    86s
```

每个 agent 的实际模式：长推理 → 短工具调用 → 长推理 → 短工具调用，循环往复。

### 累积效应

主会话本身也有推理延迟。T-quick-3 完成于 19:39，T-quick-4 的 dispatch 时间为 20:50 — 中间有 71 分钟的间隔，其中主会话在处理 T-quick-3 的结果、验证状态、claim 下一个任务。这些操作本应只需几秒，但模型推理延迟导致每步都需要数分钟。

## 偏差分类

| 类别 | 适用? | 说明 |
|------|-------|------|
| instruction-gap | 否 | Skill 定义完整，agent 行为正确 |
| context-starvation | 否 | Agent 有足够上下文 |
| trust-without-verify | 否 | 验证步骤完整 |
| wrong-priority | 否 | 优先级正确 |
| scope-creep | 否 | Agent 未超出范围 |
| pipeline-gap | 否 | Pipeline 无缝隙 |
| **模型性能** | **是** | 非行为偏差，是基础设施性能问题 |

## 结论

这不是 agent 行为偏差或 skill 定义问题。Agent 正确执行了所有步骤，产出的代码和测试均通过。问题是模型推理延迟，导致：

1. 每个 agent 平均 14 分钟（实际工具工作 <1s）
2. 7 个 agent 累计 85 分钟
3. 用户在 2+ 小时后中断

## 建议

| 优先级 | 建议 | 预期效果 |
|--------|------|---------|
| P0 | 使用推理更快的模型（如 Claude Sonnet 4.6） | 推理延迟降低 5-10x |
| P1 | 精简 task-executor 的 skill 指令和 prompt 模板，减少上下文大小 | 减少 token 生成量 |
| P2 | 减少 agent 数量：将 pipeline 步骤合并到单个 agent | 减少 dispatch 次数和推理轮数 |
