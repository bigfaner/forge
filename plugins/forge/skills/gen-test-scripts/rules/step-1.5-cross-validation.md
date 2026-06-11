---
name: step-1.5-cross-validation
description: Cross-validate Fact Table (code reconnaissance) against Contract frontmatter anchors, using handbook as authority source. Classify mismatches, generate suggested fixes, and output surface coverage report.
---

# Step 1.5: Cross-Validation

Cross-validate the Fact Table (built in Step 1) against Contract frontmatter anchor fields. This step runs after code reconnaissance and before reading Contract specifications in detail.

## Purpose

Detect mismatches between code reality (Fact Table) and design intent (Contract anchors). The handbook (design document) serves as the authority source for resolving discrepancies.

## 1.5.1 Anchor Field Extraction

For each Contract file in the target Journey, extract frontmatter anchor fields based on the detected surface type:

| Surface | Required anchor fields | Optional anchor fields |
|---------|----------------------|----------------------|
| API | `endpoint`, `method` | `content_type`, `auth_required` |
| CLI | `command` | `subcommand`, `flags`, `aliases` |
| TUI | `command` | `interactive_prompt`, `keybindings` |
| Web | `page` | `route`, `requires_auth`, `layout` |
| Mobile | `screen` | `navigation_path`, `deeplink`, `platform` |

### Extraction Logic

1. Read each Contract file's YAML frontmatter.
2. Check for anchor fields corresponding to the active surface type.
3. If anchor fields are present, record their values for comparison.
4. If anchor fields are absent, mark the Contract as "no anchors" for degradation handling.

## 1.5.2 Fact Table Matching

Map anchor fields to Fact Table entries:

| Anchor field | Fact Table lookup strategy |
|-------------|--------------------------|
| `endpoint` | Search Fact Table for entries containing API route patterns (e.g., `API_ROUTE_*`, `HTTP_*`) |
| `method` | Search Fact Table for HTTP method entries associated with the endpoint |
| `command` | Search Fact Table for CLI command entries (e.g., `CLI_COMMAND_*`) |
| `subcommand` | Search Fact Table for subcommand entries |
| `page` | Search Fact Table for Web page/route entries |
| `screen` | Search Fact Table for Mobile screen entries |

### Comparison Rules

- **Normalized comparison**: Strip trailing slashes, normalize case for HTTP methods (uppercase), normalize whitespace.
- **Parameterized routes**: Convert path parameters to a canonical form for comparison. Both `/teams/:teamId/sub-items/:subId/move` and `/teams/{teamId}/sub-items/{subId}/move` normalize to the same pattern.
- **Partial match**: If Fact Table contains a superset of the anchor's information (e.g., anchor has endpoint but Fact Table has endpoint + method), compare only the overlapping fields.

## 1.5.3 Result Classification

Each anchor-vs-Fact-Table comparison produces one classification:

### High Confidence Match

**Criteria**: Anchor value matches Fact Table value after normalization.

**Action**: No action needed. Log as verified. Include in surface coverage report.

### Low Confidence Mismatch

**Criteria**:
- Anchor value differs from Fact Table value, AND
- Fact Table signal is incomplete or ambiguous (e.g., route discovered via pattern match but handler registration is dynamic, or reconnaissance only found partial matches)

**Action**:
1. Log the mismatch with both values and source citations.
2. Proceed to authority resolution (1.5.4) to determine which value to trust.
3. Do NOT auto-resolve. Present to user for confirmation.

### Cannot Verify

**Criteria**:
- No corresponding Fact Table entry exists for the anchor field, OR
- Anchor field is absent from the Contract frontmatter, OR
- Fact Table entry exists but source citation is `UNKNOWN`

**Action**:
1. If anchor field is present but unverifiable: log as unverifiable, use anchor value with reduced confidence.
2. If anchor field is absent: proceed with Fact Table inference as fallback.
3. Log in surface coverage report as unverifiable.

## 1.5.4 Authority Resolution

When a mismatch is detected, determine the authority source:

### Handbook Lookup

1. Locate the handbook file for the active surface type:
   - API: `docs/features/<slug>/api-handbook.md`
   - CLI/TUI: `docs/features/<slug>/cli-handbook.md`
   - Web: `docs/features/<slug>/page-map.md`
   - Mobile: `docs/features/<slug>/screen-map.md`
2. If handbook exists, extract the authoritative value for the mismatched field.
3. Compare handbook value against both anchor value and Fact Table value.

### Authority Decision Matrix

| Handbook | Matches Anchor | Matches Fact Table | Resolution |
|----------|---------------|-------------------|------------|
| Exists | Yes | Yes | No real mismatch (likely normalization issue). Use anchor value. |
| Exists | Yes | No | **Code bug**: handbook and anchor agree, code differs. Flag as code bug. |
| Exists | No | Yes | Anchor is stale. Suggest fix: update anchor to match handbook. |
| Exists | No | No | All three disagree. Present all three values to user. Default to handbook. |
| Missing | N/A | Yes | No handbook to arbitrate. Trust Fact Table. Prompt user to generate handbook. |
| Missing | N/A | No | No handbook, no Fact Table match. Low confidence. Prompt user. |

### Code Bug Report Format

When handbook and anchor agree but code differs:

```
CODE BUG DETECTED
Contract: docs/features/<slug>/testing/<journey>/contracts/step-N.md
Field: method
Handbook (authority): PUT /teams/:teamId/sub-items/:subId/move
Contract anchor: PUT
Code (Fact Table): POST (source: internal/handler/subitem.go:42)
Classification: handbook and Contract agree, code implementation differs.
Action: This is a code bug, not a Contract or test issue. The implementation should be fixed to match the design specification.
```

## 1.5.5 Suggestion Generation

For each mismatch where the handbook provides an authoritative value that differs from the current anchor:

### Diff Generation

1. Construct the proposed frontmatter change.
2. Format as a unified diff:

```diff
--- a/contracts/step-2-move-subitem.md
+++ b/contracts/step-2-move-subitem.md
@@ -2,7 +2,7 @@
 journey: "task-lifecycle"
 step: 2
 step-action: "Move sub-item"
-method: "POST"
+method: "PUT"
 endpoint: "/teams/:teamId/sub-items/:subId/move"
 ---
```

3. Include handbook citation: `Source: api-handbook.md > /teams/:teamId/sub-items/:subId/move`

### User Confirmation Flow

1. Present the diff to the user.
2. Ask: "Apply this anchor fix? [y/n]"
3. If yes: write the updated frontmatter to the Contract file.
4. If no: keep current value, log rejection. Proceed with current anchor.
5. Record the decision in the surface coverage report.

<HARD-RULE>
Suggested fixes require explicit user confirmation before writing to Contract files. No automatic writes.
</HARD-RULE>

## 1.5.6 Surface Coverage Report

After all Contracts in the Journey have been cross-validated, generate a coverage report.

### Report Structure

```
=== Surface Coverage Report ===

Surface: {surface_type}
  Contracts with anchors: {anchored}/{total}
  Cross-validated (high confidence): {count}
  Mismatches detected: {total_mismatches}
    - Low confidence: {count}
    - Cannot verify: {count}
  Code bugs flagged: {count}
  Suggested fixes pending: {count}
  Applied fixes: {count}
  Rejected fixes: {count}

[Repeat for each surface type with Contracts in the Journey]

Surfaces not covered:
  - {surface}: {reason} (no handbook | no contracts | surface not applicable)

Summary: {verified}/{total_anchors} anchors verified, {mismatches} mismatches, {code_bugs} code bugs, {pending} fix(es) pending
```

### Report Requirements

- MUST list every surface type that has Contracts in the Journey
- MUST show anchor coverage ratio per surface
- MUST count results by classification
- MUST explicitly list surfaces that were NOT verified and state why
- MUST provide a one-line summary
- The report is informational only -- it does not block the pipeline

## 1.5.7 Degradation Mode

When cross-validation cannot be performed in full:

| Missing component | Behavior | User prompt |
|------------------|----------|-------------|
| No handbook for surface | Skip cross-validation for that surface. Use Fact Table inference. | "Handbook not found for surface `{surface}`. Cross-validation skipped. Recommend running `/tech-design` to generate handbook." |
| No anchor fields in Contract | Use Fact Table as inference source. No comparison possible. | "Contract `{path}` has no anchor fields. Using Fact Table inference. Consider running `/gen-contracts` to populate anchors from handbook." |
| No handbook AND no anchors | Full degradation. Proceed with Step 1 Fact Table only. | Output both prompts above. |
| Fact Table entry is UNKNOWN | Cannot verify that specific anchor. Classify as "cannot verify". | Include in coverage report as unverifiable. |

Degradation is **non-blocking** and **backward compatible**. The pipeline continues with reduced confidence, using whatever information is available.

### Lesson Scenario Coverage

The degradation mode handles the original lesson scenario (POST vs PUT mismatch) as follows:

1. Fact Table reconnaissance discovers the actual HTTP method from code (PUT).
2. If Contract has `method: POST` anchor -> cross-validation detects mismatch.
3. If api-handbook exists and defines PUT -> authority resolution confirms PUT is correct.
4. Suggestion generated: change anchor from POST to PUT.
5. User confirms -> Contract updated.
6. If api-handbook does NOT exist -> degradation: use Fact Table value (PUT) for test generation, prompt user to generate handbook.

This ensures the lesson scenario is captured regardless of handbook availability.
