# Freeform Review Extraction Prompt

Prompt template for extracting structured key findings from freeform review narratives. The extraction output feeds into the rubric scorer as injected attack points.

## Usage

This prompt is sent to an LLM after a freeform review has been completed. The LLM reads the freeform review narrative (identified by `{{FREEFORM_REVIEW}}` placeholder) and produces a JSON array of structured findings.

## Prompt Template

### System Role

```
你是一个分析助手。任务是从自由评审叙事中提取结构化的风险发现。
```

### User Role

```
从以下自由评审叙事中提取所有显式提出的风险点和改进建议。

输出格式：JSON 数组，每个元素包含：
- `summary`: 一句话概括（不超过 50 字）
- `severity`: high / medium / low
- `quote`: 叙事中的原文引用（精确到句）

规则：
1. 仅提取评审者明确表述的风险，不要推断隐含风险
2. 每个风险点独立成条，不合并
3. severity 基于评审者的语气强度判断（如使用「严重」「必须」为 high，使用「建议」「可以考虑」为 low，其余为 medium）
4. quote 必须是原文逐字引用，不得改写或概括

叙事内容：
{{FREEFORM_REVIEW}}
```

## Template Variable

| Variable | Description | Source |
|----------|-------------|--------|
| `{{FREEFORM_REVIEW}}` | The complete freeform review narrative text | Read from `<DOC_DIR>/eval/freeform-review.md` |

## Output Specification

The extraction produces a JSON array. Each element has three required fields:

| Field | Type | Constraint | Description |
|-------|------|-----------|-------------|
| `summary` | string | non-empty, max 50 chars | One-sentence summary of the finding |
| `severity` | enum | `high`, `medium`, `low` | Severity based on reviewer's tone intensity |
| `quote` | string | non-empty | Verbatim quote from the review narrative |

## JSON Validation Rules

After extraction, validate the output against these rules. If any rule fails, see Degradation Logic below.

1. **Non-empty output**: The extraction must produce at least one finding (JSON array with >= 1 element)
2. **Valid JSON**: The raw output must parse as a valid JSON array
3. **Complete fields**: Every element must have all three fields (`summary`, `severity`, `quote`) and each must be a non-empty string
4. **Severity enum**: Every element's `severity` must be one of: `high`, `medium`, `low`

## Hit Rate Estimation

After successful extraction, compute a coarse-grained hit rate to detect partial extraction failure:

**Formula**: `hit_rate = successful_extractions / keyword_paragraphs`

Where:
- `successful_extractions` = count of elements in the validated JSON array
- `keyword_paragraphs` = count of paragraphs in the freeform review that contain at least one of the marker keywords: `风险：`, `问题：`, `建议：`

**Limitations** (must be documented in any report using this metric):

1. **Overestimation of denominator**: In paragraphs where multiple risks are densely packed, a single paragraph counts as one but may contain several distinct findings. This inflates the denominator and depresses the hit rate.
2. **Underestimation of denominator**: When the reviewer describes a risk in prose without using the explicit marker prefix (e.g., describing a concern in a flowing paragraph), the paragraph is not counted. This deflates the denominator and inflates the hit rate.
3. **Purpose**: This metric is designed to trigger alerts, not to provide precise measurement. The 50% threshold is intentionally conservative to catch severe extraction failures while tolerating the metric's inherent noise.

**Threshold**: If `hit_rate < 0.5`, flag as low hit rate.

## Degradation Logic

When the extraction fails, the system degrades gracefully to the standard rubric flow:

| Condition | Action |
|-----------|--------|
| Extraction output is empty (0 elements) | Skip injection. Proceed with standard rubric flow. Inform user: "自由评审未产出有效结构化发现，已降级为标准 rubric 流程。" |
| JSON parse error | Skip injection. Proceed with standard rubric flow. Inform user: "提取产出格式异常，已降级为标准 rubric 流程。" |
| Any element missing required fields | Attempt to filter out invalid elements. If remaining valid elements >= 1, proceed with partial injection (and flag low hit rate if applicable). If 0 valid elements remain, degrade to standard rubric flow. |
| Invalid severity value | Attempt to map common variants (`HIGH` -> `high`, `Medium` -> `medium`). If mapping fails, treat as field validation failure for that element. |
| Hit rate < 50% | Still inject the valid findings, but add annotation to the eval report: "提取命中率低" + attach the complete freeform review narrative for manual review. |

## Output After Validation

On successful extraction (at least 1 valid finding), the validated findings are passed to the Pre-Revision phase (P0.5 in SKILL.md) for formatting as ATTACK_POINTS and routing to the Reviser.

On degradation (0 valid findings), the system proceeds with standard rubric evaluation without any freeform-derived attack points.
