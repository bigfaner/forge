# Freeform Expert Review Pipeline (Phase 0)

Freeform evaluation pipeline for proposal type — two-tier sequential approval model: domain expert reviews first (Phase 0), findings are routed to Reviser for pre-revision (Phase 0.5), then CTO reviews via rubric with annotated blind review (Steps 2-4 of SKILL.md). For all other types, skip directly to the Expert Dispatch Table in SKILL.md. Phase 0 delegates to subagents via Agent tool, the main session orchestrates.

This phase is executed **by default** when the resolved type is `proposal`.

## P0.1: Expert Reuse Check

Before generating a new expert, check for reusable experts in `docs/experts/`:

1. Load all `.md` files from `docs/experts/`. Filter out files with `deprecated: true` or invalid front matter.
2. Follow the reuse matching rules in `rules/freeform-expert-persistence.md`: extract domain keywords from each candidate and from the proposal, compute Jaccard overlap score.
3. If a candidate meets the threshold (Jaccard >= 0.3 or weighted score >= 5), present the match to the user via `AskUserQuestion` with two options: **Reuse** or **Generate new**.
4. If the user chooses **Reuse**, use that expert profile as `EXPERT_PROFILE` and skip to P0.3.
5. If the user chooses **Generate new**, or no candidate meets the threshold, proceed to P0.2.

## P0.2: Expert Inference

Generate a dynamic expert profile via a `general-purpose` agent using the prompt defined in `experts/freeform/expert-inference.md`:

1. Spawn agent with `model: "sonnet"`, providing `PROPOSAL_PATH` (the proposal document path) and `EXISTING_EXPERTS` (list of current expert file contents from `docs/experts/`).
2. The agent performs domain analysis and generates an expert profile using `experts/freeform/expert-template.md`. Reuse check was already completed in P0.1 — instruct the agent to skip it.
3. The agent returns the generated profile (or indicates reuse).
4. Present the expert profile to the user via `AskUserQuestion` with three options: **Accept**, **Modify**, **Regenerate**. Include the cross-reference report and self-check questions.
5. Handle the modification loop and rejection limits per `experts/freeform/expert-inference.md` (max 3 modification rounds, max 3 consecutive rejections).
6. On acceptance, save the expert profile to `docs/experts/<slug>.md` and set `EXPERT_PROFILE` to the profile content.
7. If the user skips (chooses to skip after rejection limit, or manually aborts), degrade per Phase 0 Degradation Summary table.

**Error: Expert generation failure**: If the inference agent returns incoherent output (missing domain keywords, empty background) or fails entirely, inform the user and offer two options: manually describe an expert direction, or degrade per Phase 0 Degradation Summary table.

## P0.3: Freeform Review

Conduct a freeform narrative review using a `general-purpose` agent:

1. Spawn agent with `model: "sonnet"`, providing `DOC_DIR` and `EXPERT_PROFILE` (from P0.1 reuse or P0.2 generation).
2. The agent reads `experts/freeform/freeform-reviewer.md` (which in turn references `experts/freeform/freeform-review-protocol.md`) and conducts the review.
3. The agent writes the review to `<DOC_DIR>/eval/freeform-review.md` and returns a status summary (`FREEFORM_REVIEW: completed/failed`).

**Error: Freeform review failure**: If the agent returns `FREEFORM_REVIEW: failed`, or the output file is empty/missing, degrade per Phase 0 Degradation Summary table.

## P0.4: Extract Findings

Extract structured findings from the freeform review narrative using a `general-purpose` agent:

1. Read the freeform review from `<DOC_DIR>/eval/freeform-review.md`.
2. Spawn agent with `model: "sonnet"`, providing the extraction prompt from `experts/freeform/extraction-prompt.md` with `{{FREEFORM_REVIEW}}` replaced by the review content.
3. The agent returns a JSON array of findings.
4. Validate the extraction output per `extraction-prompt.md` JSON Validation Rules.
5. Compute hit rate per `extraction-prompt.md` Hit Rate Estimation. Set `HIT_RATE = <computed value>` (used in P0.5c for low-hit-rate annotation).
6. If 0 valid findings remain after validation, degrade per Phase 0 Degradation Summary table.
7. If >= 1 valid finding, set `FREEFORM_FINDINGS = <validated JSON array>` and proceed.

## P0.5: Pre-Revision (Freeform Findings)

After Phase 0 completes with valid findings and before the Scorer cycle starts, execute a Pre-Revision step that routes freeform findings directly to the existing Reviser via a synthetic eval report. This step runs as iteration 0 (separate from the `MAX_ITERATIONS` budget). The Scorer loop runs iterations 1 through MAX_ITERATIONS.

P0.5 executes when P0.4 produces ≥ 1 valid finding, regardless of `MAX_ITERATIONS`. All degradation paths (P0.1–P0.4 failure or P0.5 error) converge to the standard rubric flow.

### P0.5a: BASELINE_SCORE (Informational Metric)

Before pre-revision, obtain an informational baseline score:

1. Spawn a single Scorer subagent (no Reviser, does not consume iteration budget) to evaluate the original proposal.
2. Record the score as `BASELINE_SCORE`.
3. This is informational only — it does not act as a gate or consume an iteration.
4. If the Scorer call fails, set `BASELINE_SCORE = null` and continue. Do not degrade.

### P0.5b: Save Phase 0 Baseline Snapshot

Save a copy of the current proposal document(s) to `<DOC_DIR>/eval/baseline-snapshot/`. This snapshot captures the pre-revision state for the overall-level rollback (see Two-Level Rollback in SKILL.md Step 5).

### P0.5c: Format Findings as ATTACK_POINTS

Format each finding from `FREEFORM_FINDINGS` into the ATTACK_POINTS structure expected by the Reviser protocol:

```
- **[severity]** summary | quote: "quote" | improvement: <verb phrase>
```

Each finding is also classified into one of three triage layers:
- **Factual correction** (verifiable defect in original text): direct edit.
- **Structural/architectural suggestion**: edit only when the finding identifies a verifiable internal inconsistency (e.g., two sections make contradictory claims, or a stated constraint is violated by the described architecture); otherwise defer to Scorer cycle. When partially valid (concern is legitimate but proposed solution is not), mark as `partially-accepted`: apply only the non-controversial portion of the edit.
- **Subjective preference**: mark as "not actionable", no edit, but record in iteration-0 report's "Classification Audit" section with classification rationale and original finding summary.

**Low hit-rate annotation**: If `HIT_RATE` (from P0.4 step 5) is < 0.5, add the following annotation to the iteration-0 report (P0.5d) after the ATTACK_POINTS section: `**Note: Low extraction hit rate. The following contains only partial findings; see the full freeform review narrative at <DOC_DIR>/eval/freeform-review.md for context.**`

**Borderline handling**: When a finding does not clearly belong to one layer, the pre-reviser must mark it "borderline" and defer (not silently classify as not actionable). Borderline findings are listed separately in the iteration-0 report for user review.

For error handling, see P0.5 Degradation Summary below.

### P0.5d: Construct Synthetic Eval Report

Build a synthetic eval report that satisfies the Reviser protocol's (`experts/protocol/reviser-protocol.md`) `EVAL_REPORT_PATH` dependency:

```yaml
iteration: 0
title: "Pre-Revision (Freeform Findings)"
ATTACK_POINTS:
  - (formatted findings from P0.5c, excluding subjective-preference findings)
BORDERLINE_FINDINGS:
  - (findings marked "borderline" in P0.5c triage)
SKIPPED_FINDINGS:
  - (subjective-preference findings, marked "not actionable")
rubric:
  (all dimensions): N/A
```

Save to `<DOC_DIR>/eval/iteration-0-report.md`.

This satisfies the Reviser protocol's minimum input format: ATTACK_POINTS-driven revision with rubric data non-participatory (Reviser does not read, compare, or decide based on rubric scores).

### P0.5e: Invoke Reviser (Iteration 0)

Spawn the existing Reviser as a `general-purpose` agent via the Agent tool with `model: "sonnet"`:

- Inputs: `DOC_DIR`, `EVAL_REPORT_PATH` = iteration-0 report, `ATTACK_POINTS` = formatted findings.
- The Reviser follows `experts/protocol/reviser-protocol.md` unchanged — no protocol modification.
- The Reviser applies the triage results from P0.5c and performs edits per the ATTACK_POINTS format.

For error handling, see P0.5 Degradation Summary below.

### P0.5f: Tag Modified Paragraphs

After the Reviser completes its edits, annotate each modified paragraph with:

```html
<!-- pre-revised: {severity} -->
```

Where `{severity}` is the severity of the finding that triggered the edit. These tags:
- Are HTML comments — invisible in rendered output, visible to Scorer when reading the document.
- Enable the Scorer's annotated blind review — the Scorer knows **which areas were changed** but not **why** they were changed.

### P0.5g: Iteration Counter Increment

After successful pre-revision: set `ITERATION = 1` (pre-revision consumed iteration 0). The Scorer loop starts from iteration 1.

Set `PRE_REVISION_EXECUTED = true` — the Scorer uses annotated blind review mode (see `rules/scorer-composition.md`).

After the final report is generated (end of Step 5 in SKILL.md), record the expert's review history per `rules/freeform-expert-persistence.md` quality tracking section, and check auto-deprecation.

### P0.5 Degradation Summary

All Phase 0.5 error paths leave `PRE_REVISION_EXECUTED` unset and degrade to standard rubric flow:

| Error Scenario | Degradation Action | User Notification |
|----------------|-------------------|-------------------|
| Findings formatting failure | Skip pre-revision, enter Scorer directly | "Pre-revision 格式化失败，已跳过。" |
| Pre-reviser returns error | Skip, log warning | "Pre-reviser 执行失败，已跳过。" |
| Empty report produced | Log iteration-0 "no changes", degrade to standard flow | (no notification, standard rubric flow) |
| Format anomaly (Reviser output file exists but is empty, truncated, or contains no ATTACK_POINTS response) | Discard pre-revision, restore from baseline snapshot | "Pre-revision 产出格式异常，已丢弃。" |

## Phase 0 Degradation Summary

All Phase 0 error paths converge to the standard rubric flow. No degradation path interrupts the eval pipeline:

| Error Scenario | Degradation Action | User Notification |
|----------------|-------------------|-------------------|
| Expert generation failure | Skip to standard rubric flow | "专家生成失败，已降级为标准 rubric 流程。" |
| User skips after rejection limit | Skip to standard rubric flow | (user initiated, no extra notification) |
| Freeform review failure (empty/failed) | Skip to standard rubric flow | "自由评审未产出有效发现，已降级为标准 rubric 流程。" |
| Extraction output is empty | Skip to standard rubric flow | "自由评审未产出有效结构化发现，已降级为标准 rubric 流程。" |
| Extraction JSON invalid (0 valid after filtering) | Skip to standard rubric flow | "提取产出格式异常，已降级为标准 rubric 流程。" |
| Partial extraction (hit rate < 50%) | Inject valid findings + annotate | "提取命中率低" annotation in eval report + full narrative attached |

For Phase 0.5 (Pre-Revision) error paths, see P0.5 Degradation Summary above.
