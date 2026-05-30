# Expert Dispatch Table

Expert files are located at `experts/scorer/`.

| type | scorer experts |
|------|---------------|
| `proposal` | `[cto]` |
| `prd` | `[pm, qa]` |
| `design` | `[architect]` |
| `ui-web`, `ui-mobile`, `ui-tui` | `[ux-engineer]` |
| `consistency` | `[editor]` |
| `journey` | `[qa]` |
| `contract` | `[qa]` |
| `validate-code` | `[code-reviewer]` |
| `validate-ux` | `[ux-auditor]` |

Fallback for unmapped types: use the generic inline prompt below (no expert file loaded).

```
You are a senior reviewer evaluating the document at {{DOC_DIR}} against the rubric at {{RUBRIC_PATH}}. Apply the rubric strictly and identify all weaknesses.
```

# Scorer Prompt Composition

Read the scorer protocol at `experts/protocol/scorer-protocol.md`.

Compose the scorer prompt by concatenating sections in the order defined in "Order in final composed prompt" below. Template variables in the scorer protocol: `{{DOC_DIR}}`, `{{RUBRIC_PATH}}`, `{{REPORT_PATH}}`, `{{ITERATION}}`, `{{PREVIOUS_REPORT_PATH}}`.

**Context Injection**: If `CONTEXT_CONTENT` was loaded in Step 1.4, append the following section after the expert content in every composed prompt:

```
<injected-context>
The following project reference material is provided for reality-checking the evaluated document. Use it to detect contradictions, violations, or gaps -- do not evaluate the reference material itself.

{{CONTEXT_CONTENT}}
</injected-context>
```

For unmapped types (not in dispatch table), compose a single prompt using the generic inline fallback above plus the scorer protocol (with variables replaced) plus context injection.

# Scorer Composition — Freeform Integration

Proposal type two-path dispatch:

| Condition | Scorer Behavior |
|-----------|----------------|
| `PRE_REVISION_EXECUTED = true` | Annotated blind review (see below) |
| Unset (Phase 0 degraded / non-proposal) | Standard rubric flow, no freeform block |

## Annotated Blind Review (pre-revision executed)

When pre-revision was executed (`PRE_REVISION_EXECUTED = true`), append annotated blind review instructions. The Scorer sees `<!-- pre-revised: {severity} -->` markers but not freeform findings content.

Append to the composed prompt:

```
<annotated-blind-review>
This document has undergone a Pre-Revision phase. Revised paragraphs are annotated with HTML comment markers `<!-- pre-revised: {severity} -->`.

Annotated blind review rules:
1. `<!-- pre-revised: {severity} -->` markers indicate paragraphs modified during Pre-Revision. For marked regions: focus on whether the revision introduced new issues or omissions, rather than re-evaluating the original corrected problem.
2. The severity marker is for attention allocation only; it does not affect scoring criteria.
3. In the eval report, record attack density separately for annotated and unannotated regions for bias detection. Format:

   **Bias Detection Report**:
   - Annotated regions: N attack points / X paragraphs = density Y
   - Unannotated regions: M attack points / Z paragraphs = density W
   - Ratio (annotated/unannotated): R

4. When the Scorer's rubric judgment contradicts the pre-revision direction (e.g., Scorer believes a paragraph should be deleted but pre-revision just added it), generate the attack point per rubric standards, but tag it with `conflict-with-pre-revision` for review.
</annotated-blind-review>
```

**Order in final composed prompt**:
1. Scorer protocol (with template variables replaced)
2. Expert file content
3. `<injected-context>` block (if `CONTEXT_CONTENT` loaded)
4. `<annotated-blind-review>` block (if `PRE_REVISION_EXECUTED = true`)

# Scorer Agent Inputs

Common inputs for all agents:
- `DOC_DIR` = document directory
- `RUBRIC_PATH` = resolved rubric file
- `ITERATION` = current iteration (1-based)
- `PREVIOUS_REPORT_PATH` = previous report (only if iteration > 1)

Report paths per expert (for multi-expert types, each expert writes to a separate report):
- `REPORT_PATH` = `<doc_dir>/eval/iteration-{{N}}.md` (single-expert)
- `REPORT_PATH` = `<doc_dir>/eval/iteration-{{N}}-{{expert}}.md` (multi-expert)

Type-specific report path overrides:
- `consistency`: `docs/features/<slug>/eval-consistency/eval/iteration-{{N}}.md`
- `proposal`: `docs/proposals/<slug>/eval/iteration-{{N}}.md`
- `validate-code`: `docs/features/<slug>/eval/validate-code.md`
- `validate-ux`: `docs/features/<slug>/eval/validate-ux.md`

Type-specific inputs:
- `ui-*`: add `PRD_PATH` = `docs/features/<slug>/prd/prd-ui-functions.md` (if exists)
- `consistency`: add `SCOPE` = value from `--scope`
- `validate-ux`: add `UX_SNAPSHOT_PATH` = path to generated `ux-snapshot.md`
- `journey`: add `SURFACE_TYPE` = value from `.forge/config.yaml` `surfaces` field; add `SURFACE_RULE_PATH` = gen-journeys skill's `rules/surface-<type>.md` (resolve relative to the gen-journeys skill directory)
- `contract`: add `SURFACE_TYPE` = value from `.forge/config.yaml` `surfaces` field; add `SURFACE_RULE_PATH` = gen-journeys skill's `rules/surface-<type>.md` (resolve relative to the gen-journeys skill directory)

Do NOT pass reviser change summaries to the scorer.

# Multi-Expert Result Merging

**For single-expert types**: extract using robust score extraction:
1. Extract score using regex `/SCORE:\s*(\d+)\/(\d+)/`. If pattern not found, scan the scorer agent's output for the last line matching a `number/number` pattern. If still not found, report error and stop.
2. Per-dimension scores from `DIMENSIONS:` section
3. Attack points from `ATTACKS:` section

**For multi-expert types**:
1. Extract score and attacks from each expert's output
2. **Gate score**: average the total scores across all experts (rounded to nearest integer)
3. **Attack points merge**: LLM-merge attack points from all experts in the main session using this prompt:

```
Merge overlapping attack points from {{N}} expert evaluations. Keep unique attacks from each. Combine duplicates into single attacks preserving the strongest prescription. Do not remove any unique attack. Output the merged list in the same format:

1. [dimension]: [specific weakness] -- [quote from document] -- [what must improve]
```

4. **Write merged report**: Write the merged attacks + averaged scores to `<doc_dir>/eval/iteration-{{N}}-merged.md`. This file serves as `EVAL_REPORT_PATH` for the reviser (Step 4.1). Single-expert types continue using `iteration-{{N}}.md` directly.
