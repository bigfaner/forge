---
id: "2"
title: "Implement Annotated Blind Review in Scorer Composition"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
type: "doc"
mainSession: false
---

# 2: Implement Annotated Blind Review in Scorer Composition

## Description

Replace the freeform findings injection mechanism in scorer-composition.md with annotated blind review instructions. The Scorer's composed prompt should NOT include `<injected-freeform-findings>` block; instead, include instructions for interpreting `<!-- pre-revised: {severity} -->` HTML comment markers and performing bias detection (attack density comparison between annotated and unannotated regions). Conditionally deprecate freeform-injection.md by adding `status: deprecated` frontmatter, preserving original content for potential future restoration.

## Reference Files

- `proposal.md#Proposed-Solution` — defines annotated blind review concept and new information flow
- `proposal.md#Design-Decisions` — Decision 2 (annotated blind review design, why not full blind or full traceability), Decision 6 (conditional deprecation of freeform-injection.md)
- `proposal.md#Scope` — 改动文件 table for scorer-composition.md and freeform-injection.md; Non-Functional Requirements (compatibility: only affects type == proposal)
- `proposal.md#Success-Criteria` — SC #2 for scorer prompt composition verification
- `proposal.md#Key-Risks` — annotated blind review false positives (Scorer over-examines pre-revised areas)

## Affected Files

### Create
| File | Description |
|------|-------------|
| _(none)_ | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/rules/scorer-composition.md` | Remove `<injected-freeform-findings>` block from Freeform Findings Injection section; replace with conditional branch (skip injection when `FREEFORM_INJECTION = false`); add `<!-- pre-revised -->` annotation interpretation instructions (~5 lines); add bias detection report template (attack density per annotated/unannotated region, ~5 lines) |
| `plugins/forge/skills/eval/rules/freeform-injection.md` | Add `status: deprecated` frontmatter field; preserve complete injection semantic definition for future restoration |

### Delete
| File | Reason |
|------|--------|
| _(none)_ | |

## Acceptance Criteria

1. **SC #2**: Scorer's composed prompt does NOT contain freeform findings content when pre-revision mode is active
2. Scorer's composed prompt contains `<!-- pre-revised -->` annotation interpretation instructions: focus on whether revision introduced new problems, severity as attention guide, not re-evaluating original issues
3. Scorer prompt includes bias detection report template: attack density recorded separately for annotated and unannotated regions
4. When Scorer's rubric judgment conflicts with pre-revision direction, attack point is annotated with `conflict-with-pre-revision` flag for review
5. `freeform-injection.md` has `status: deprecated` frontmatter with original content fully preserved
6. scorer-composition.md conditional branch: when `FREEFORM_INJECTION = false` (pre-revision mode), skip injection block entirely

## Hard Rules

- Follow forge-distribution.md path conventions (relative paths only, not project root paths)
- Do NOT physically delete freeform-injection.md — conditional deprecation only
- Restore path documented: remove `status: deprecated` frontmatter + remove scorer-composition.md conditional branch = 2 config changes

## Implementation Notes

- **Scorer prompt additions (~5 lines in scorer-composition.md)**:
  ```
  <!-- pre-revised: {severity} --> 标记表示该段落经过 Pre-Revision 修改。
  对标记区域：关注修订是否引入了新问题或遗漏，而非重新评估已修正的原始问题。
  severity 标记供注意力分配参考，不影响评分标准。
  在 eval report 中分别记录标注区域与未标注区域的 attack density，供偏误检测。
  当 Scorer 的 rubric 判断与 pre-revision 修改方向矛盾时，以 rubric 标准为准生成 attack point，标注 conflict-with-pre-revision。
  ```
- **Bias detection threshold**: if annotated region attack density is systematically >= 30% higher than unannotated across >= 2 consecutive evals, trigger "标注偏误告警" — this is an empirical feedback loop, not a code constant
- **scorer-composition.md "Freeform Findings Injection (Phase 0)" section**: convert to conditional branch — when `FREEFORM_INJECTION = false` (set by Task 1's P0.5), skip the injection block; otherwise, follow original injection path from freeform-injection.md
- **Only affects type == proposal**: other eval types are unaffected since freeform review only activates for proposal type
