---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Pre-Revision Report (Iteration 0)

## ATTACK_POINTS

- **[high]** CLI 输出语义删除边界不清 | quote: "The phrase 'CLI 输出语义' is dangerously ambiguous. Output field semantics (what SURFACE_KEY means, when it is absent) are not the same as internal implementation details." | improvement: 在提案中增加"删除边界规则"小节，将 CLI 相邻文本分为三类：(1) 指令性操作——保留；(2) 输出契约（字段名+缺失含义）——保留；(3) 行为解释——删除。提供 2-3 个实际文件中的示例。

- **[high]** E-I 去重可能误删跨步骤防护约束 | quote: "a task executor, following the instruction 'remove items that duplicate body text,' may mechanically match keywords and delete these guardrails because they reference concepts that also appear in the body, even though the E-I version encodes a stricter constraint than the body version." | improvement: 修改 E-I 去重成功标准为约束级别审计：每个保留的 E-I 条目满足以下之一：(a) 正文未包含该约束，或 (b) 正文包含该约束但 E-I 版本的强制等级更高。

- **[high]** E-I 去重成功标准不可执行 | quote: "This verification method will produce false positives (flagging items as 'deduplicated' when they actually encode stronger constraints than the body)." | improvement: 同上一条，将关键词匹配改为约束级别审计。

- **[high]** CLI 行为描述删除无法用结构化方法验证 | quote: "There is no regex that distinguishes 'this sentence is a behavioral description that should be deleted' from 'this sentence is an output contract that should be preserved.'" | improvement: 修改 SC-1（grep 验证）为包含三层验证：(1) grep 检查 "What .* Does" section（已定义）；(2) 每个被修改文件的 diff 中无 output contract 丢失（人工 spot-check）；(3) 补充具体 before/after 示例作为校准基准。

- **[medium]** 跨文件冗余作为问题证据但被排除修复范围 | quote: "12 near-verbatim duplications exist between two skill files...the natural moment to address the cross-file duplication is now, while both files are being touched." | improvement: 在 Evidence 中明确标注这些跨文件重复是"已知但有意不在本次修复范围内"，并说明理由（结构独立性优先于文字去重），而非将其作为核心问题证据。

- **[medium]** quick.md fallback 被定性为 bug 但未反驳原始设计意图 | quote: "The proposal treats this as an unambiguous bug, but the source file contains an explicit design justification that the proposal does not acknowledge or refute." | improvement: 在 Evidence 中增加对原始设计理由的分析："quick.md 注释 'This preserves quick mode's streamlined nature' 是有意的 fail-open 设计。然而该设计未区分两种失败场景：(1) 配置不存在（新用户应走确认流程）(2) 配置损坏（应 fail-safe）。两种场景下跳过确认门都不是正确行为。"

- **[medium]** 个别实例证据缺乏具体性（gen-contracts 等） | quote: "No specific section number is cited, no line reference, and no quote from the file." | improvement: 对于 40 个清晰度问题，为每个问题添加文件行号或原文引用。如果位置已知（审计已逐行扫描），包含它们可节省大量重复工作。

## BORDERLINE_FINDINGS

- **[medium]** grep 验证只能捕获标题式描述 | Borderline: 这是对 #4（CLI 行为描述删除验证）的补充说明，不是独立问题。已部分纳入 #4 的修复方案。 |

## SKIPPED_FINDINGS (Subjective Preferences)

- 建议：定义三类 CLI 文本分类规则 — 已纳入 ATTACK_POINT #1
- 建议：约束级别审计替代关键词重叠 — 已纳入 ATTACK_POINT #2/#3
- 建议：承认 quick.md 设计意图 — 已纳入 ATTACK_POINT #6
- 建议：为 40 个清晰度问题提供行号 — 已纳入 ATTACK_POINT #7
- 建议：重新考虑跨文件去重排除决定 — 用户已明确决定保持独立，不修改

## Rubric

(all dimensions): N/A (pre-revision, rubric non-participatory)
