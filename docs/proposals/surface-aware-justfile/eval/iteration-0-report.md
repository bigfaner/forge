---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Pre-Revision Report (Iteration 0)

## ATTACK_POINTS

- **[high]** Step 3a fallback 循环依赖 | quote: "如果 Surface 规则也缺失（无 surface 配置），`test` 配方的生成逻辑是什么？回退到当前行为？那当前行为依赖的就是 `test.execution.run`——但这个字段已经被移除了。这里有一个循环依赖。" | improvement: 定义新的 test 配方生成 fallback 链替代被移除的 test.execution.run

- **[high]** 混合项目 scope 映射规则缺失 | quote: "map 的 key 是用户自定义的目录路径，不是 `frontend`/`backend` 这样的固定枚举值。如果用户的 config 是 `{admin-panel: web, payment-service: api}`，run-tests 应该调用 `just dev admin-panel` 还是 `just dev frontend`？" | improvement: 明确约束 scope 值必须是 surfaces map 的 key

- **[medium]** dev 配方 scope 参数未在 Standard Target Contract 定义 | quote: "当前 init-justfile 的混合项目 scope 参数只定义在 `compile` 和 `run` 的示例中（Step 3b），`dev` 没有明确的 scope 参数规范。" | improvement: 在 Standard Target Contract 和 Recipe Parameter Signatures 中为 dev 增加 [scope] 定义

- **[medium]** surface 信息获取缺少统一优先级规则 | quote: "当两者不一致时（config 说 web，文件信号检测出 api），以哪个为准？init-justfile 用哪个？run-tests 用哪个？" | improvement: 定义统一的 surface 信息源优先级规则

- **[low]** api/web 编排序列完全相同可能过度拆分 | quote: "这两个 surface 在'测试编排序列'和'关键配方'上的唯一差异是 probe 检查的目标不同" | improvement: 考虑合并 api/web 为 service 规则或共享编排模板

- **[high]** config schema 变更范围被低估 | quote: "这不是'常规 Go 开发'——它改变了 Forge CLI 的配置接口契约。如果其他 skill 也需要读取 `test.timeout`，影响面会超出提案预估的 8-12 个任务。" | improvement: 将 config schema 变更描述升级，明确边界和影响面

- **[high]** test.execution 移除后静默忽略已有配置 | quote: "如果有开发者或 AI agent 按照现有文档已经配置了 `test.execution`，移除后会静默忽略这些配置而不报错。提案应该要求在 run-tests 启动时检测到 `test.execution` 节点时输出废弃警告" | improvement: 定义 test.execution 废弃检测和警告行为

- **[medium]** 编排硬编码为固定配方名 | quote: "这看起来简化了，但实际上将'命令可配置性'从 config.yaml 转移到了 justfile 配方体中" | improvement: 讨论这一 trade-off 的显式性差异，论证 justfile 作为唯一抽象层的充分性

- **[medium]** 与 run-tests 现有 surface 感知环境检查重叠 | quote: "run-tests SKILL.md 第 4 步的'Environment Readiness Check'已经通过 `rules/env-check.md` 和 `surface-<type>.md` 做 surface 感知的环境检查了——提案与现有机制之间的重叠没有被讨论。" | improvement: 讨论 run-tests 现有 surface 感知机制与提案的关系

- **[medium]** 语言模板与 surface 规则冲突无仲裁机制 | quote: "提案将这两个步骤分开，但没有讨论当语言模板生成的 `dev` 与 surface 规则指导的 `dev` 发生冲突时如何仲裁。" | improvement: 定义语言 vs surface 的仲裁规则（surface 优先覆盖 test/dev/run/probe）

- **[medium]** journey 过滤策略未定义具体映射 | quote: "风险表中也标记了这个问题的可能性为中、影响为高，但缓解措施只是'Surface 规则记录映射关系'——这个'记录'本身需要具体设计。" | improvement: 为每种 surface 提供 journey 过滤的最小规范和示例

## BORDERLINE_FINDINGS

- 编排硬编码为固定配方名（severity: medium）— 这是设计决策的 trade-off 还是结构性缺陷？归入 structural，因为提案未讨论此 trade-off。

## SKIPPED_FINDINGS

- 7 条 建议： 标记的建议 — 属于 subjective preference，已通过 ATTACK_POINTS 中的事实修正和结构性问题覆盖。

## Classification Audit

| 类别 | 数量 |
|------|------|
| Factual correction | 5 |
| Structural/architectural suggestion | 6 |
| Subjective preference | 7 |

## rubric

(all dimensions): N/A
