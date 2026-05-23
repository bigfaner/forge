# Reviser Prompt Composition

Read the reviser protocol at `experts/protocol/reviser-protocol.md`.

Resolve `EVAL_REPORT_PATH`:
- **Single-expert types**: `<doc_dir>/eval/iteration-{{N}}.md`
- **Multi-expert types**: `<doc_dir>/eval/iteration-{{N}}-merged.md` (written in Step 2.3)

Compose the reviser prompt by concatenating:
1. Reviser protocol content (with template variables replaced: `{{DOC_DIR}}`, `{{EVAL_REPORT_PATH}}`)
2. The merged `ATTACK_POINTS` from Step 2.3 (replacing the `{{ATTACK_POINTS}}` placeholder in the protocol)
3. Context injection (if `CONTEXT_CONTENT` was loaded in Step 1.4 -- see below)

**Context Injection**: If `CONTEXT_CONTENT` was loaded in Step 1.4, append the following section after the attack points in the reviser prompt:

```
<injected-context>
The following project reference material is provided for reality-checking the evaluated document. Use it to detect contradictions, violations, or gaps -- do not evaluate the reference material itself.

{{CONTEXT_CONTENT}}
</injected-context>
```

The reviser receives **only** the protocol + merged attacks + optional context. No rubric, no expert file.

# Reviser Type-Specific Constraints

- `consistency`: Do NOT modify `prd/`. Classify attack points by fix target before invoking.
- `journey`: Revise Journey document based on failed dimensions and attack points. Preserve surface-specific `required_outcomes` structure. After reviser completes:
- `contract`: Revise Contract documents based on failed dimensions and attack points. Preserve six-dimension structure and Preconditions mutual exclusivity. After reviser completes:
- `consistency`: re-assemble document bundle
- Increment iteration counter, return to Step 2
