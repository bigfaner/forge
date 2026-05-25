# Eval Report — Iteration 2 (Merged)

**PM Score**: 872/1000
**QA Score**: 842/1000
**Average Score**: 857/1000

## Dimension Breakdown (PM / QA)

| Dimension | PM | QA | Max |
|-----------|-----|-----|-----|
| Background & Goals | 89 | 92 | 100 |
| Flow Diagrams | 130 | 135 | 150 |
| Flow Completeness | 180 | 165 | 200 |
| User Stories | 178 | 173 | 200 |
| Scenario Completeness | 127 | 116 | 150 |
| Edge Case Coverage | 79 | 70 | 100 |
| Scope Clarity | 89 | 91 | 100 |

## Merged Attack Points

1. [Scope Clarity]: `forge surfaces` CLI 在 In Scope "需新建" 和 Out of Scope "新增 forge CLI 命令" 之间存在直接矛盾 — In Scope 第59行标记 CLI "需新建"，但 Out of Scope 第75行写"新增 forge CLI 命令" — 必须解决矛盾：将 CLI 纳入 In Scope 或将 In Scope 中的声明改为"已存在"或"独立前置特性"

2. [Flow Completeness]: Goal 3/4 的 surface-key 迁移涉及 9 项变更点但无实施流程描述 — In Scope 列出 prompt.go 重写、Go struct 变更、模板迁移等 9 项变更，Flow Description 部分只描述了 init-justfile 和 run-tests 流程 — 必须添加迁移实施流程简述，包含变更顺序、关键依赖

3. [Flow Completeness]: teardown 作为关键恢复动作自身的失败行为未定义 — 多处将 teardown 作为失败恢复动作但 teardown 本身失败时无退出码或后续行为 — 补充 teardown 失败时的行为（如重试一次后放弃、日志记录、exit code）

4. [Flow Diagrams]: cli/tui 分支在 Mermaid 图中缺少错误分支 — Flow Description 文字说"每步检查退出码"但 Mermaid 图 Test2（cli/tui）直接连 RunEnd 无 exit 1 分支 — 在 Mermaid 图中为 cli/tui test 添加 exit 1 错误路径，或添加注释说明 cli/tui 无 teardown 故直接结束

5. [Edge Case Coverage]: surface-key 命名字符限制未定义 — "配方名带 surface-key 前缀：`dev-<surface-key>`" — just 配方名不允许空格/特殊字符，需定义命名约束或清洗规则（如仅允许 `[a-zA-Z0-9_-]`）

6. [Edge Case Coverage]: config.yaml surfaces 格式错误和 CLI 输出异常等关键错误路径未覆盖 — Error Handling Paths 表覆盖 5 场景但缺少 config.yaml 格式校验失败、forge surfaces CLI 输出格式异常 — 为 config.yaml 校验和 CLI 输出解析补充防御性错误处理

7. [Scenario Completeness]: 混合项目"并行启动"与 just 语法语义不一致 — 第147行 "按依赖序并行启动所有 dev server" 但 just 的 `(dep1 dep2)` 是串行依赖列表 — 澄清并行/串行语义，与 just 实际行为对齐（建议改为"按依赖顺序串行启动"）

8. [User Stories]: Story 4 AC 耦合 Go 函数名实现细节 — prd-user-stories.md 直接引用 `GetSurfaceKey()` 方法名作为验收条件 — 重写为行为验证："Given 旧任务文件含 scope: frontend，When run-tests 读取，Then 按默认编排策略执行"

9. [blindspot]: `.forge/test-state.json` 恢复机制是空壳承诺 — "会话中断后可通过 `.forge/test-state.json` 恢复清理"只此一句无具体流程 — 补充写入时机、格式、恢复步骤，或降级为"现有机制"引用而非新承诺

10. [blindspot]: "编排级配方"概念缺少定义清单 — Flow Description 说"Surface 规则覆盖语言模板的编排级配方（test/dev/run/probe）"但 Scope 和编排表列出 5-6 种含 test-setup/test-teardown — 明确编排级配方的完整清单

11. [blindspot]: 移除 test.execution 对现有用户配置的影响未评估 — "移除 `test.execution` 节点文档"无用户故事描述迁移影响（残留配置被静默忽略还是报错）— 添加简短迁移说明
