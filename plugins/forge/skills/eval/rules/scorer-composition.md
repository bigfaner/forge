# Expert Dispatch Table

Expert files are located at `experts/scorer/`.

| type | scorer experts |
|------|---------------|
| `proposal` | `[cto]` |
| `prd` | `[pm, qa]` |
| `design` | `[architect]` |
| `ui-web`, `ui-mobile`, `ui-tui` | `[ux-engineer]` |
| `test-cases`, `ui-test-cases`, `tui-test-cases`, `mobile-test-cases`, `api-test-cases`, `cli-test-cases` | `[qa]` |
| `consistency` | `[editor]` |
| `validate-code` | `[code-reviewer]` |
| `validate-ux` | `[ux-auditor]` |

Fallback for unmapped types: use the generic inline prompt below (no expert file loaded).

```
You are a senior reviewer evaluating the document at {{DOC_DIR}} against the rubric at {{RUBRIC_PATH}}. Apply the rubric strictly and identify all weaknesses.
```

# Scorer Prompt Composition

Read the scorer protocol at `experts/protocol/scorer-protocol.md`.

For each expert, compose a scorer prompt by concatenating:
1. Scorer protocol content (with template variables replaced: `{{DOC_DIR}}`, `{{RUBRIC_PATH}}`, `{{REPORT_PATH}}`, `{{ITERATION}}`, `{{PREVIOUS_REPORT_PATH}}`)
2. Expert file content (e.g., `experts/scorer/pm.md`)
3. Context injection (if `CONTEXT_CONTENT` was loaded in Step 1.4 -- see below)

**Context Injection**: If `CONTEXT_CONTENT` was loaded in Step 1.4, append the following section after the expert content in every composed prompt:

```
<injected-context>
The following project reference material is provided for reality-checking the evaluated document. Use it to detect contradictions, violations, or gaps -- do not evaluate the reference material itself.

{{CONTEXT_CONTENT}}
</injected-context>
```

For unmapped types (not in dispatch table), compose a single prompt using the generic inline fallback above plus the scorer protocol (with variables replaced) plus context injection.

# Freeform Findings Injection (Phase 0)

When Phase 0 was executed and produced valid findings (i.e., `FREEFORM_INJECTION = true` in the eval skill), append the `<injected-freeform-findings>` block **after** all existing sections in the composed scorer prompt.

Follow the injection rules in `rules/freeform-injection.md`:
1. Format the validated findings array into `{{FORMATTED_FINDINGS}}` (one line per finding: `- **[severity]** summary | 原文引用: "quote"`)
2. Wrap in `<injected-freeform-findings>` block with the standard header and instructions
3. If `LOW_HIT_RATE = true`, include the partial extraction annotation and append the complete freeform review narrative

**Order in final composed prompt** (when freeform injection is active):
1. Scorer protocol (with template variables replaced)
2. Expert file content
3. `<injected-context>` block (if CONTEXT_CONTENT was loaded)
4. `<injected-freeform-findings>` block (if FREEFORM_INJECTION = true)

When `FREEFORM_INJECTION` is not set (no `--freeform-expert`, or Phase 0 degraded), no freeform block is added — the composed prompt is identical to the standard flow.

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
- `test-cases`, `ui-test-cases`, `tui-test-cases`, `mobile-test-cases`, `api-test-cases`, `cli-test-cases`: add `PRD_FILES` = paths to prd-spec.md and prd-user-stories.md
- `consistency`: add `SCOPE` = value from `--scope`
- `validate-ux`: add `UX_SNAPSHOT_PATH` = path to generated `ux-snapshot.md`

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

`test-cases`, `ui-test-cases`, `tui-test-cases`, `mobile-test-cases`, `api-test-cases`, `cli-test-cases`: If Step Actionability < 200, warn that gen-test-scripts is blocked.
