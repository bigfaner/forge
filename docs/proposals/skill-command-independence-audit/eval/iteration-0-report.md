---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Iteration 0: Pre-Revision Report

## ATTACK_POINTS

- **[high]** 耦合图遗漏 gen-contracts 到 gen-journeys 的反向引用 | quote: "gen-contracts 引用 gen-journeys/SKILL.md" is absent from the Evidence section entirely — Evidence 部分列了 6 处跨 skill 引用但遗漏了 gen-contracts SKILL.md 第 58 行对 gen-journeys/SKILL.md "Surface Detection" section 的引用 | improvement: 在 Evidence 部分补充该反向引用，在 Scope 部分将 gen-contracts 的修复描述改为"内联 gen-journeys Surface Detection 相关知识 + 补充双向耦合说明"

- **[high]** Related Skills 删除清单误将 quick-tasks 的 ## Reference Files 归类为 pipeline 无用信息 | quote: "quick-tasks" 列在 "9 个 skill" 中需删除 Related Skills/Integration/References，但 quick-tasks 的 ## Reference Files（第 122 行）是模板占位符 {{REFERENCE_FILES}} 的替换规则说明，不是 pipeline 上下游信息 | improvement: 在 Scope 的 In Scope 中明确 quick-tasks 的 ## Reference Files 不删除，只删除 ## Integration 段落

- **[high]** 删除 gen-contracts 和 gen-journeys 的 ## Reference 段落会违反独立加载原则 | quote: "9 个 skill 的 Related Skills/Integration/References 章节内容均可从正文中隐含推断" — 但 gen-contracts ## Reference（第 267-276 行）包含 Contract、Outcome、Semantic Descriptors 等 6 个概念定义，这些定义在 SKILL.md 正文中被使用但未集中定义 | improvement: 对 gen-contracts 和 gen-journeys 的 ## Reference 段落，标注为"合并到内联知识中作为定义段落"而非"删除"

- **[high]** 所有内联操作缺少"所需内容"的边界定义 | quote: 每个内联操作统一使用 "内联 XXX 所需内容" 措辞，但 journey-contract-model.md 有 184 行、test-isolation.md 有 138 行、test-type-model.md 有 51 行，未说明内联哪些段落 | improvement: 为每个内联操作添加 INJECT/SKIP 段落清单或至少说明"内联至多 N 行"

- **[medium]** extract-design-md 对 ui-design/templates/styles/ 的引用是运行时数据读取而非知识引用 | quote: "extract-design-md: 内联 ui-design/styles 匹配逻辑" — 但实际引用（第 124 行）是 "read the corresponding style file from ui-design/templates/styles/<name>.md"，是运行时数据消费指令 | improvement: 将 extract-design-md 的处理策略改为"在 extract-design-md 内部创建 rules/style-matching.md 包含匹配特征摘要，风格文件保留在 ui-design 中"或标注为类似 forensic 的设计意图豁免

- **[medium]** 15% 行数减少目标与内联增加行数存在数值冲突 | quote: "总行数减少 >= 15%"（即约 1071 行），但内联操作会增加约 150-250 行，精简需压缩约 1200-1400 行才能达标 | improvement: 将目标调整为"净减少 >= 10%"或"总行数不超过 6500 行"

- **[medium]** "功能等价"成功标准不可验证 | quote: "所有 skill/command 修改后功能等价（无行为变更）" — skill 文件是自然语言指导 AI agent 的文档，非可执行程序 | improvement: 替换为结构化检查清单：所有 HARD-RULE/HARD-GATE/EXTREMELY-IMPORTANT/PROHIBITIONS 块计数不变，所有决策表完整保留，所有 Step 序号和流程步骤完整

- **[medium]** Solution 声称"对全部 21 个 skill 和 16 个 command 执行三维度清理"但实际 scope 只涉及约 15 个 skill 和 3 个 command | quote: "对全部 21 个 skill 和 16 个 command 执行三维度清理" vs scope 中列出的具体项 | improvement: 修正 Solution 描述为"对有跨引用、冗余或低价值章节的 skill 和 command 执行清理"

- **[medium]** execute-task 与 run-tasks 的"60-70% 结构重叠"判断不准确 | quote: "execute-task 与 run-tasks command 60-70% 结构重叠" — 实际上两者共享约 20-30 行（claim 格式和 fix-type 表），但核心逻辑完全不同（单任务 vs 循环调度） | improvement: 将重叠描述修正为"共享约 20-30 行接口契约（claim 格式 + fix-type 表），核心逻辑各自独立"

- **[low]** 漂移风险接受缺少量化依据和缓解机制 | quote: "可接受——独立性带来的维护简化大于同步成本" — 无量化依据，内联后将产生约 12 份额外拷贝 | improvement: 添加轻量缓解：对内联段落使用 <!-- INLINE:origin=... --> 标记提供可追溯性

## BORDERLINE_FINDINGS

(none)

## SKIPPED_FINDINGS

- (subjective preference) 建议补全耦合图为双向图格式 — 呈现格式偏好，非结构性缺陷
- (subjective preference) 建议为 extract-design-md 考虑替代方案 — 与 ATTACK_POINT #5 重叠，已在该条中处理

## Rubric

(all dimensions): N/A
