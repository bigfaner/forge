---
iteration: 0
title: "Pre-Revision (Freeform Findings)"
---

# Pre-Revision Report (Iteration 0)

## ATTACK_POINTS

1. **[high]** Constraints 与 In Scope 描述自相矛盾 | quote: "修改 `forge-cli/pkg/forgeconfig/` 中的配置结构体（`CopyFiles` → `Includes` + 兼容读取）" 与 In Scope 中 "直接替换，不保留旧字段" 矛盾 | improvement: 统一为"直接替换"策略，删除"兼容读取"字样

2. **[medium]** NFR 遗漏可观测性需求 | quote: "Non-Functional Requirements" 仅列"向后兼容"和"性能"两条 | improvement: 增加可观测性条目，明确日志格式和内容要求

3. **[medium]** SC #2 不可测试 | quote: "已存在时输出明确的提示信息，区分「新建」和「进入」两种路径" | improvement: 指定输出格式和关键词，使 SC 可客观验证

4. **[medium]** SC 覆盖不完整，场景 #5 和 #6 无对应 SC | quote: Key Scenarios #5 (--no-launch) 和 #6 (--interactive) 出现在需求中但无 SC 验证 | improvement: 为这两个场景增加对应的 SC

5. **[medium]** 改动范围估算缺乏依据 | quote: "改动范围较小（2-3 个文件，约 40 行变更）" | improvement: 基于代码库搜索结果重新评估

## BORDERLINE_FINDINGS

- 捆绑变更未论证必要性：两个正交变更的捆绑发布增加了风险表面积，但用户在 brainstorm 阶段已确认两者关联，属于设计决策而非逻辑错误

## SKIPPED_FINDINGS (Subjective preference / user override)

- 配置静默失效风险 + loud deprecation 建议：用户在 brainstorm 阶段明确选择"不兼容也不报错，不留存兼容性代码"。Challenge Override: user chose clean break.
- 拆分为独立变更步骤：用户在 brainstorm 中讨论后选择合并为一个提案
- --reinclude flag：属于未来增强，当前 scope 明确排除
- resume 命令可发现性：次要 UX 问题，非提案核心

## Rubric

(all dimensions): N/A (pre-revision)
