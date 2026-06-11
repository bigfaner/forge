---
iteration: 0
title: "Pre-Revision (Freeform Findings, Round 2)"
---

# Iteration 0: Pre-Revision Report (Round 2)

## ATTACK_POINTS

- **[high]** `init` zero-discovery lacks diagnostic output, users misjudge as tool failure | quote: "当父目录下存在嵌套分组目录时，扫描结果为零，用户无法区分"工具坏了"和"目录结构不符合预期"" | improvement: Add diagnostic output showing scan count, skipped dirs with reasons, and suggestions
- **[high]** assign "核心上下文" is undefined fuzzy term, field mapping table needed | quote: ""核心上下文"是一个无法实现的模糊术语：继承太少则丢失决策依据，继承全部则因 schema 差异导致字段映射错误" | improvement: Replace with explicit field mapping table
- **[medium]** `.forge-workspace.yaml` lacks field growth constraint mechanism | quote: "没有机制保证这个约束在 Dashboard 和 Wiki 模块开发时不被打破" | improvement: Add governance rule limiting to {schema_version, projects}
- **[medium]** Cache mtime tracking granularity undefined | quote: ""mtime 指纹"是追踪项目根目录？docs/ 下所有文件？还是仅 manifest 文件？" | improvement: Specify tracking per-project docs/ directory max mtime
- **[medium]** Workspace proposal lifecycle lacks Assigned→Done closure | quote: "缺少从 feature 完成回写 workspace proposal 的机制" | improvement: Define Done transition and close command
- **[medium]** Existing manifests without schema_version lose structured info | quote: "无版本号时视为 v0 还是格式未知？" | improvement: Define v0 default with visual distinction in CLI output
- **[medium]** Brainstorm skill needs workspace context mode, scope undeclared | quote: "提案未显式声明哪些技能需要修改、修改范围是什么" | improvement: Add skill modification scope to Constraints section

## BORDERLINE_FINDINGS

(none)

## SKIPPED_FINDINGS

- All 6 suggestions are advisory, directly addressed by corresponding ATTACK_POINTS above

## Rubric

(All dimensions): N/A — Pre-Revision phase
