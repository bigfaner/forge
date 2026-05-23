# Freeform Finding Injection Rules

Rules for injecting extracted freeform review findings into the rubric scorer prompt. This is the bridge between the freeform expert review (Phase 0) and the rubric scoring phase.

## When to Inject

Injection occurs only when:
1. The freeform review completed successfully (agent returned `FREEFORM_REVIEW: completed`)
2. The extraction process produced at least 1 valid finding (JSON validation passed)
3. The eval type is `proposal` (freeform injection is not applicable to other eval types)

If any condition is not met, skip injection entirely and proceed with standard rubric flow.

## Injection Format

Append the following section **after** the expert content and **after** any existing `<injected-context>` block in the composed scorer prompt:

```markdown
<injected-freeform-findings>
The following findings were extracted from a freeform expert review conducted before rubric scoring. They represent risks and concerns identified by an independent domain expert who reviewed the document without knowledge of this rubric. You MUST address these findings during your evaluation:

{{FORMATTED_FINDINGS}}

Instructions:
- Treat each finding as an attack point to evaluate during scoring
- If a finding maps to an existing rubric dimension, incorporate it into that dimension's evaluation
- If a finding does not map to any rubric dimension, record it as `[beyond-rubric]: [finding summary]` in the ATTACKS section of your report
- If a finding contradicts the direction suggested by the rubric, annotate the relevant dimension with: 「自由评审与 rubric 存在分歧」
</injected-freeform-findings>
```

### Finding Formatting

`{{FORMATTED_FINDINGS}}` is generated from the validated extraction JSON array. Format each finding as:

```
- **[severity]** summary | 原文引用: "quote"
```

Example:
```
- **[high]** 提案未验证分布式一致性场景下的脑裂恢复 | 原文引用: "提案假设分区恢复后状态自动合并，但未讨论脑裂场景下的数据冲突解决策略"
- **[medium]** 回滚计划的触发条件过于模糊 | 原文引用: "回滚条件仅表述为'性能不达标'，缺乏可度量的阈值"
- **[low]** 建议补充竞品分析 | 原文引用: "提案未提及同类工具如 Vitest 的 profile 机制"
```

### Source Identification

All injected content is wrapped in `<injected-freeform-findings>` tags and explicitly labeled as "来自自由专家评审". This ensures the scorer (and anyone reading the composed prompt) can distinguish these attack points from the rubric's own dimensions and from any `<injected-context>` blocks.

## Beyond-Rubric Findings

When the scorer encounters a finding that does not map to any rubric dimension:

1. The scorer records it in the ATTACKS section of the evaluation report
2. Format: `[beyond-rubric]: [finding summary]`
3. Position: at the end of the ATTACKS list, after all rubric-mapped attack points
4. Each `[beyond-rubric]` entry must include `summary` and `quote` sub-fields

Example in scorer output:
```
### ATTACKS
1. [hidden-costs]: 提案未考虑分布式锁的运维成本 -- "锁服务作为共享依赖，其故障影响面未评估" -- 增加锁服务的 SLO 和降级方案
2. [rollback]: 回滚触发条件不可度量 -- "性能不达标"缺乏量化标准 -- 定义具体延迟和吞吐量阈值
3. [beyond-rubric]: 提案的 profile 注册时机影响测试发现机制 -- 原文引用: "profile 在 test runner 初始化后才注册，但测试发现阶段在此之前执行" -- summary: 注册时序问题可能导致 profile 未被测试发现流程感知
```

## Contradiction Annotation

When the scorer identifies that a freeform finding contradicts the direction suggested by the rubric for a given dimension:

1. The scorer annotates the relevant dimension in its evaluation report
2. Annotation text: 「自由评审与 rubric 存在分歧」
3. The scorer provides both perspectives: what the rubric suggests vs. what the freeform expert flagged
4. The scorer does not resolve the contradiction — it surfaces both views for the user to judge

Example:
```
### Dimension: extensibility
Score: 7/10
Comment: The proposal provides a clean plugin interface. However, 自由评审与 rubric 存在分歧：
- Rubric perspective: extensibility is well-designed through the plugin interface
- Freeform expert concern: the plugin lifecycle hooks are insufficient for cross-cutting concerns like observability
```

## Partial Extraction Handling

When the extraction hit rate is below 50% (as computed per `experts/freeform/extraction-prompt.md`):

1. Inject the valid findings normally
2. Add the following annotation to the `<injected-freeform-findings>` block:
   ```
   **注意：提取命中率低。以下仅包含部分发现，完整自由评审叙事见附件。**
   ```
3. Append the complete freeform review narrative text after the annotation
4. In the final eval report, include the annotation: "提取命中率低"

## Degradation Path

When injection is skipped (extraction failed or produced 0 valid findings):

1. No `<injected-freeform-findings>` block is added to the scorer prompt
2. The scorer runs exactly as it would without the freeform phase
3. The eval report includes a note: "自由评审未产出有效结构化发现，已降级为标准 rubric 流程。"
4. The complete freeform review narrative is preserved at `<DOC_DIR>/eval/freeform-review.md` for manual review

## Interaction with Existing Scorer Composition

This injection mechanism integrates with the scorer composition flow defined in `rules/scorer-composition.md`:

1. The standard scorer prompt is composed per the existing rules (protocol + expert + context injection)
2. The `<injected-freeform-findings>` block is appended **after** all existing sections
3. Order in the final composed prompt:
   - Scorer protocol
   - Expert file content
   - `<injected-context>` block (if present)
   - `<injected-freeform-findings>` block (if freeform extraction succeeded)
