# Eval Report — Iteration 1 (Merged)

**PM Score**: 828/1000
**QA Score**: 720/1000
**Average Score**: 774/1000

## Dimension Breakdown (PM / QA)

| Dimension | PM | QA | Max |
|-----------|-----|-----|-----|
| Background & Goals | 88 | 85 | 100 |
| Flow Diagrams | 135 | 128 | 150 |
| Flow Completeness | 165 | 115 | 200 |
| User Stories | 170 | 166 | 200 |
| Scenario Completeness | 115 | 90 | 150 |
| Edge Case Coverage | 75 | 52 | 100 |
| Scope Clarity | 80 | 84 | 100 |

## Merged Attack Points

1. [Flow Completeness]: 跨组件数据流文档缺失 — "涉及 7+ 组件的跨系统特性，关键数据流（surface 信息从 config.yaml → 任务 frontmatter → run-tests skill 的传递链）无显式文档" — 必须添加跨组件数据流表，覆盖 surface 信息从配置到执行的完整传递链，包括各组件间的数据格式和传递方式

2. [Scenario Completeness]: Exit code 处理与项目已有 exit code 体系不一致 — "exit 1/2/3" for probe failure 均映射到相同的 teardown+abort 行为，忽略了 retryable (exit 1) 和 blocking (exit 2) 的语义区分；HARD-GATE 精确语义未定义 — 对齐 exit code 语义到 BIZ-error-reporting-001（exit 1=retryable, exit 2=blocking），明确 HARD-GATE 与 exit code 的关系

3. [Edge Case Coverage]: 错误路径覆盖严重不足 — 未知 surface 类型、Surface 规则文件缺失、forge surfaces CLI 执行失败、just 版本 < 1.4.0 等关键错误路径均未记录；"覆盖 5 种 surface 类型"未定义未知类型的 fallback — 为每个检测和加载步骤补充错误处理路径，定义未知 surface 类型的显式错误行为

4. [Edge Case Coverage]: run-tests dispatcher 在 surface 信息不可用时缺少 fallback — "优先任务文档 frontmatter → forge surfaces CLI" 无两个来源均失败时的行为定义 — 定义默认编排策略或按 BIZ-error-reporting-002 提供明确错误和恢复提示

5. [blindspot]: Surface-key "用户自定义" 目标与 "5 固定类型" 实现存在张力 — "surface-key 值域从固定枚举迁移为用户自定义" 承诺用户自定义 key，但所有流程和规则仅处理 5 种类型 — 明确 goal 的真实含义：用户自定义 surface-key 名称 + 固定 5 种 surface type，或增加自定义类型扩展机制

6. [blindspot]: forge surfaces CLI 是关键前置依赖但未在 PRD 中声明 — "surface-key-assignment 规则文件路径分类改为 CLI 动态查询"、"forge surfaces <path> longest-prefix-match" — 必须在 Scope 或 Related Changes 中明确声明 forge surfaces CLI 的前置条件状态（已存在/需新建/需扩展）

7. [blindspot]: 混合项目缺少端到端流程描述 — "混合项目 dev 配方接受 surface-key 参数" 仅一句话，无具体的 init-justfile 生成策略和编排表说明 — 必须添加混合项目（如 web+api）的完整生成和编排流程描述

8. [blindspot]: Probe 重试参数分散 — Observability 提到 "[retry 3/30]" 和流程提到 "重试轮询" 但重试次数、间隔、超时未在 PRD 中统一指定 — 将 probe 重试规格整合到 Flow Description 的单一定义点，包含明确参数值
