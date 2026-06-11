---
id: "6"
title: "run-tests 重写 + Journey 隔离"
priority: "P1"
estimated_time: "3h"
dependencies: ["5"]
scope: "backend"
breaking: false
type: "feature"
mainSession: false
---

# 6: run-tests 重写 + Journey 隔离

## Description

重写 run-tests skill：项目声明执行命令，Forge 调度执行 + 结果报告。实现 Journey 隔离机制——每个 Journey 在独立临时工作目录中执行，避免并行运行时互相干扰。

来源：proposal Pipeline 第 4 步、Scope "run-tests skill"和"Journey 隔离"。

## Reference Files
- `docs/proposals/contract-journey-test-model/proposal.md` — Source proposal
- `plugins/forge/skills/run-tests/` — 现有 run-tests 实现
- `forge-cli/internal/cmd/` — 现有 testing 子命令

## Acceptance Criteria

- [ ] run-tests 调用项目声明的测试执行命令（从 `.forge/config.yaml` 的 `test-command` 读取），收集结果并生成报告
- [ ] 每个 Journey 在独立临时目录中执行（路径包含 journey 名称和随机后缀），执行完成后临时目录被清理
- [ ] 3 个 Journey 并行运行时结果与顺序运行完全一致（文件系统状态、退出码、输出内容 diff 为空）
- [ ] Contract 验证失败时，输出包含：失败维度名称、Contract 文件路径、期望值、实际值
- [ ] Contract 测试执行总时间 ≤ 现有集成测试执行时间的 120%（排除 Setup 阶段）

## Hard Rules

- Journey 隔离是硬性要求——不支持共享状态的并行执行
- 执行失败不阻塞其他 Journey 的执行
- 临时目录清理必须在 Journey 执行完成后（无论成功或失败）

## Implementation Notes

- run-tests 是 4 步 pipeline 的最后一步，认知任务最简单（项目声明命令，Forge 调度）
- Journey 隔离通过临时目录实现：每个 Journey 开始前创建临时目录 → 复制必要文件 → 执行测试 → 收集结果 → 清理目录
- 并行安全的关键：文件系统状态完全隔离，无共享文件
