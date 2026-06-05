# Iteration 0 Report: Pre-Revision (Freeform Findings)

**Type**: Synthetic Eval Report (Freeform Findings)
**Source**: Freeform Expert Review by Contract Pipeline & Test Specification Architect
**Date**: 2026-06-05

## Classification Audit

### Accepted Findings (Factual/Structural)

- **[high]** 自动修复错误时的缓解措施"保存原始值到注释中可回溯"严重不足 | quote: "当一个错误的自动修复进入 Contract 后，它会影响后续所有 gen-test-scripts 的生成。除非有明确的回滚流程和检测机制，否则'可回溯'只意味着'理论上可以人工恢复'" | improvement: 改为"建议修复"加人工确认流程，至少在初始版本中如此

- **[high]** 设计文档滞后于代码变更时，交叉验证会以过时设计为准破坏正确的 Contract | quote: "在快速迭代的项目中，设计文档经常滞后于代码变更，如果每次代码先行更新了接口但设计文档未同步，交叉验证就会以过时的设计文档为准去'修复'实际上已经正确的 Contract" | improvement: 增加 handbook 新鲜度检查机制，比对 handbook 生成时间戳与 tech-design 最后修改时间

- **[high]** 自动修复是 proposal 中最危险的操作，代码侦察 edge case 可能导致误覆盖 | quote: "如果一个 Contract 的 endpoint 原本是正确的，但由于代码侦察的某个 edge case（比如路由注册使用了动态模式）导致 Fact Table 得到了错误信息，交叉验证就会误判为不匹配" | improvement: 将交叉验证结果分类为高/低/无法验证三个置信度级别，只有高置信度不匹配才自动处理

- **[medium]** 静态分析代码侦察无法覆盖所有路由注册模式，侦察不完整导致误判 | quote: "Fact Table 的代码侦察是静态分析，它无法覆盖所有路由注册模式——如果项目使用了插件系统、动态加载、反射机制等，侦察结果就是不完整的" | improvement: 增加交叉验证结果的分类报告，区分高置信度不匹配、低置信度不匹配、无法验证

- **[medium]** CLI/Web/Mobile 锚点字段定义远比 API 复杂，proposal 对新 handbook 格式一笔带过 | quote: "CLI command 的标识要复杂得多：子命令嵌套、命令别名、参数变体等。page-map 和 screen-map 的挑战更大" | improvement: 为每种 surface 定义锚点字段的完整 schema，包括必填和可选字段

- **[medium]** "全 surface 覆盖"与"可分批实现"存在矛盾，未定义中间状态 | quote: "proposal 在 Scope 中承诺了'全 surface 覆盖'，但在 Key Risks 中又说'可分批实现'。这两者之间存在张力" | improvement: 明确分批实现路线图，第一批先做 API surface

- **[medium]** tech-design 更新后缺乏 handbook 重新生成触发机制 | quote: "当设计文档更新时，是否有机制触发 handbook 的重新生成？如果没有，handbook 就会过期" | improvement: 增加 handbook 新鲜度检查，比对 handbook 与 tech-design 时间戳

- **[medium]** Contract 手动编辑后锚点可能失效 | quote: "用户可能手动编辑 Contract，修改了操作描述但忘记更新 endpoint" | improvement: 在 Contract frontmatter 增加 last_anchor_sync 时间戳字段

- **[medium]** 交叉验证依赖两个不可靠源（Fact Table 和 handbook）的比对 | quote: "当两者都不可靠时，交叉验证就成了两个不可靠源之间的比对" | improvement: 在 proposal 中明确承认这一局限性，并说明降级策略

- **[low]** "无额外网络或 IO 开销"说法不准确 | quote: "交叉验证在 gen-test-scripts Step 1（代码侦察）中执行，无额外网络或 IO 开销" | improvement: 修正为"无显著额外开销"

- **[medium]** 排除代码不存在时的强制交叉验证，导致设计阶段无法捕获 endpoint 冲突 | quote: "在代码还没写之前，如果设计文档中的 endpoint 定义就有冲突，现在没有任何机制能在设计阶段就捕获这个问题" | improvement: 增加设计阶段 handbook 内部一致性检查

- **[low]** 缺少 handbook 时静默降级无提示 | quote: "proposal 没有讨论是否有机制提示用户缺少 handbook 以及建议生成 handbook" | improvement: 在降级路径中增加用户提示

- **[medium]** 部分 surface 有 handbook 时交叉验证只覆盖部分 Contract | quote: "这种不一致可能导致用户产生虚假的安全感——以为所有 surface 都在被验证，但实际上只有 API 被覆盖了" | improvement: 交叉验证报告明确列出已验证和未验证的 surface

### Borderline Findings

- **[low]** Journey 到 Contract 映射中 surface 类型信息来源不明 | quote: "如果这个 surface 类型信息不在 Journey 中而在 Contract 中，那么 anchor 填充时就有一个先有鸡还是先有蛋的问题" | improvement: 确认 surface 类型信息在 Contract 中已有，无需额外改动

### Skipped Findings (Subjective)

None — all findings have factual or structural basis.

## Triage Summary

| Category | Count |
|----------|-------|
| Accepted | 13 |
| Borderline | 1 |
| Skipped | 0 |
| **Total** | **14** |

Triage rate: 100% (14/14)
Accepted + Partially-accepted: 92.9% (13/14)

## Rubric

All dimensions: N/A (freeform findings, not rubric-scored)
