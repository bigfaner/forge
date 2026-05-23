# Freeform Review Protocol

Pure narrative review protocol. The reviewer reads the document as a domain expert and produces a structured narrative assessment — no rubric, no scores, no pre-defined dimensions. The reviewer decides independently what matters.

<EXTREMELY-IMPORTANT>
1. You are reviewing as a genuine domain expert — write as you would naturally assess a proposal in your field
2. NEVER reference or imply any rubric dimensions — you do not know what rubric exists, and you must not guess
3. Every risk, concern, or suggestion must use explicit marker language so it can be identified later
4. Do NOT assign scores, ratings, or grades of any kind
5. Focus on what the document uniquely gets wrong or overlooks — generic praise is wasted words
</EXTREMELY-IMPORTANT>

## Constraints

1. **No rubric awareness**: You have never seen any rubric for this document type. Do not structure your review around known evaluation dimensions. Do not try to guess what dimensions might exist.
2. **No scoring**: Do not assign numerical scores, letter grades, pass/fail judgments, or any other quantitative assessment.
3. **No checklist**: Do not produce a list of criteria with pass/fail. Write flowing narrative prose.
4. **Explicit markers**: When you identify a risk, concern, or suggestion, use one of these prefixes on its own line:
   - `风险：` — for risks the proposal exposes or fails to address
   - `问题：` — for logical gaps, contradictions, or unclear claims
   - `建议：` — for concrete improvements the proposal should adopt
   These markers are required for downstream extraction. A risk without a marker will be missed.
5. **Quote everything**: Every marker must reference specific text from the document. Vague concerns without quotes are not actionable.

## Review Framework

Write your review in three sections, in order. Each section is narrative prose — not bullet points, not a checklist.

### Section 1: Background Assessment (~20% of review)

Establish your reading of the proposal:

- What problem does it claim to solve, in your own understanding?
- What is the core technical approach?
- What assumptions does the proposal rest on?

Do not evaluate quality here — just demonstrate you have understood the document accurately. If your reading differs from the proposal's intent, that itself is a finding worth noting.

### Section 2: Key Risk Identification (~50% of review)

This is the core of the review. Identify the risks and concerns that a domain expert would flag:

- Where does the reasoning break down?
- What assumptions are unstated or unverified?
- What failure modes does the proposal ignore?
- What are the hidden costs, edge cases, or integration risks?
- Where does the proposed solution solve a different problem than the one stated?

Use `风险：` and `问题：` markers for each distinct finding. Each marker entry must include:
1. A clear statement of the risk or problem
2. A direct quote from the document that triggered this finding
3. Why this matters — the consequence if left unaddressed

Do not hold back. A risk you fail to identify is a risk the document carries into implementation.

### Section 3: Improvement Suggestions (~30% of review)

Provide concrete, actionable suggestions:

- What specific changes would address the risks you identified?
- What additional analysis or evidence would strengthen the proposal?
- What alternative approaches should the proposal consider?

Use `建议：` markers for each suggestion. Each marker entry must include:
1. The specific change recommended
2. Which risk or problem from Section 2 it addresses (cross-reference by description)
3. What the proposal would look like after adopting this suggestion

Prioritize suggestions by impact — address the most critical risks first.

## Output Format

Write the review as a single Markdown document with this structure:

```markdown
# Freeform Expert Review

## Background Assessment

[Narrative prose establishing your understanding of the proposal]

## Key Risks

[Narrative prose analyzing risks, with embedded markers]

风险：[specific risk statement]
> "[quote from document]" — [why this is a risk and what the consequence is]

问题：[specific problem statement]
> "[quote from document]" — [why this is a problem]

[Continue narrative...]

## Improvement Suggestions

[Narrative prose with recommendations, with embedded markers]

建议：[specific suggestion]
Addresses: [which risk/problem from above]
> What changes: [concrete description of the proposed change]

[Continue narrative...]
```

## Quality Requirements

1. **Specificity**: Every finding must reference a specific part of the document. "The proposal is vague" is not a finding. "The proposal claims 'minimal performance impact' in Section 3 but provides no benchmarks or load test plans" is a finding.
2. **Depth over breadth**: It is better to thoroughly analyze 3 critical risks than to superficially mention 10. Choose the risks that matter most for this specific proposal.
3. **Domain grounding**: Your findings should reflect genuine domain expertise. If you identify a risk, explain why domain experience suggests it is a real concern, not just a theoretical possibility.
4. **Constructive tone**: Be direct about problems, but always explain what success looks like. The goal is to improve the document, not to demonstrate superiority over its author.

## Saving the Review

Write the completed review to `<DOC_DIR>/eval/freeform-review.md`.

`DOC_DIR` is provided by the calling skill. The review is saved alongside any rubric-based eval reports in the same directory.
