# Eval Report: Pre-Revision (Freeform Findings)

**Iteration**: 0
**Title**: Pre-Revision (Freeform Findings)
**Source**: Freeform expert review by 测试类型语义精确性审查专家

## ATTACK_POINTS

### Accepted Findings

- **[high]** "API 契约测试"与行业契约测试概念（Pact）产生语义冲突 | quote: "行业中'契约测试'（Contract Testing）有明确且特定的含义：它是指服务消费者与提供者之间关于 API 接口行为的约定的验证，典型框架是 Pact" | improvement: 重新审视"API 契约测试"命名，明确此处的"契约"与 Pact 意义上的消费者契约的区别，或改用不产生歧义的术语

- **[high]** TUI"半黑盒视角"判定标准未定义，与 CLI"黑盒视角"边界不清 | quote: "提案为 TUI 标注了'半黑盒视角'而为 CLI 标注了'黑盒视角'，但从未定义判定标准" | improvement: 定义黑盒/半黑盒的判定标准，或删除该标注避免引入未定义的分类维度

- **[high]** "集成"一词在 CLI/TUI 映射中语义模糊，CLI 测试不集成组件也不集成外部系统 | quote: "CLI 测试编译独立二进制并通过子进程调用，既不是集成多个内部模块（它是黑盒的，根本不接触内部函数），也不是集成外部系统（它只是启动自己的二进制）" | improvement: 重新审视"集成测试"命名，考虑更精确的术语或显式定义此处"集成"的含义

- **[medium]** 将"端到端"保留给 Web 的论证不完整，仅考察了 CLI 一个反例 | quote: "提案的归谬论证只考察了 CLI 这一个反例，就得出'只有 Web 真正端到端'的结论。但这个结论需要更强的支撑" | improvement: 补充对 API 和 Mobile 的端到端性分析，或显式定义"端到端"的充分必要条件

- **[medium]** Mobile 定义为"UI 测试"而非"端到端测试"，与 Web 分类逻辑不一致 | quote: "如果 Web 的测试被称为'端到端测试'是因为它通过浏览器模拟用户操作，那么 Mobile 的测试通过 Maestro 驱动移动端 UI、验证渲染和交互，在逻辑上与 Web 测试是同构的" | improvement: 统一 Web 和 Mobile 的分类逻辑，或给出两者确实不同的理由

- **[medium]** 验证维度列表粒度不一致（CLI/API 是具体输出属性，Web 包含高层概念"状态流转"） | quote: "CLI 的验证维度是'退出码 + stdout/stderr'，API 的验证维度是'HTTP 状态码 + 响应体 + Header'，Web 的验证维度是'UI 渲染 + 用户交互 + 状态流转' — 这三组验证维度的粒度不在同一个层次" | improvement: 统一验证维度的粒度层次

- **[medium]** 分类体系混合了 Surface 和验证视角两个分类维度但未承认 | quote: "提案以 Surface 为一级分类键，但语义定义中又引入了验证视角（'黑盒视角'、'半黑盒视角'）作为二级属性" | improvement: 增加"分类标准声明"，明确各分类维度之间的关系

## BORDERLINE_FINDINGS

（无 borderline findings）

## SKIPPED_FINDINGS

以下 findings 归类为"主观偏好"，标记为 not actionable：

1. 建议增加"分类标准声明" — 归类为主观偏好：属于文档结构建议而非内部不一致
2. 建议将 CLI 集成测试重命名为 CLI 行为测试 — 归类为主观偏好：术语选择有合理空间
3. 建议统一验证维度粒度 — 已通过 ATTACK_POINT 覆盖不一致性本身，具体粒度选择是主观的
4. 建议定义"端到端"精确语义 — 已通过 ATTACK_POINT 覆盖论证不完整性
5. 建议定义半黑盒判定标准或删除 — 已通过 ATTACK_POINT 覆盖
6. 建议增加可扩展性分析 — 归类为主观偏好：属于文档结构建议

## RUBRIC

(All dimensions): N/A

## CLASSIFICATION AUDIT

- Total findings: 14
- Factual correction: 1 (API 契约测试语义冲突)
- Structural suggestion: 6 (论证不完整、分类不一致、维度混合、标注未定义、粒度不一致、术语模糊)
- Subjective preference: 7 (6 条建议 + 1 条可扩展性分析)
