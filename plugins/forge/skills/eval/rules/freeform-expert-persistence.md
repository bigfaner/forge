# Expert Persistence, Reuse & Deprecation Rules

Rules for persisting dynamically generated expert profiles to `docs/experts/`, reusing existing experts across evaluations, tracking expert effectiveness, and deprecating underperforming experts.

## Directory Structure

All expert profiles are stored in the user project's `docs/experts/` directory. Each expert is one file:

```
docs/experts/
├── distributed-systems-architect.md
├── test-infrastructure-engineer.md
├── ux-researcher.md
└── ...
```

### File Naming

- Slug derived from the expert's domain: lowercase, hyphens, max 40 characters
- Must be unique within `docs/experts/`
- If a name collision occurs, append a numeric suffix (e.g., `api-designer-2.md`)

### File Format

Each file uses YAML front matter + Markdown body, consistent with the template defined in `experts/freeform/expert-template.md`.

**Required front matter fields:**

```yaml
---
domain: "distributed-systems, consistency, high-concurrency"
background: "10 years building distributed message queues and consensus protocols at scale. Led migration from single-primary to multi-primary replication at a top-tier cloud provider."
review_style: "Identifies hidden coupling by tracing data flows across module boundaries. Focuses on failure mode analysis and operational complexity."
generated_for: "docs/proposals/distributed-consensus/proposal.md"
created_at: "2026-05-23T14:30:00Z"
review_history:
  - proposal: "docs/proposals/distributed-consensus/proposal.md"
    date: "2026-05-23"
    substantive_change: true
    rubric_delta: 45
    attack_points_changed: true
  - proposal: "docs/proposals/raft-optimization/proposal.md"
    date: "2026-05-25"
    substantive_change: false
    rubric_delta: 8
    attack_points_changed: false
deprecated: false
---
```

| Field | Type | Description |
|-------|------|-------------|
| `domain` | string (comma-separated keywords) | 2-5 domain keywords describing the expert's specialization |
| `background` | string | 3-5 sentences of professional background with verifiable specifics |
| `review_style` | string | One paragraph describing review approach and methodology |
| `generated_for` | string (file path) | Path to the proposal that triggered this expert's initial generation |
| `created_at` | string (ISO 8601 timestamp) | When the expert profile was first created |
| `review_history` | array of objects | Record of every evaluation using this expert (see below) |
| `deprecated` | boolean | Whether this expert has been deprecated (default: `false`) |

**review_history entry fields:**

| Field | Type | Description |
|-------|------|-------------|
| `proposal` | string (file path) | Path to the proposal evaluated with this expert |
| `date` | string (ISO 8601 date) | Date of the evaluation |
| `substantive_change` | boolean | Whether using this expert produced a substantive change in rubric outcome |
| `rubric_delta` | integer | Absolute rubric score difference vs. no-injection baseline (1000-point scale) |
| `attack_points_changed` | boolean | Whether the attack points list changed (at least 1 addition/removal/modification) |

The Markdown body follows the template from `experts/freeform/expert-template.md`: persona description, domain keywords, review focus areas, and cross-reference checklist.

## Expert Reuse Matching

Before generating a new expert, the system checks for reusable experts in `docs/experts/`.

### Step 1: Load Candidates

Read all `.md` files from `docs/experts/`. Filter out any file where `deprecated: true`.

### Step 2: Extract Keywords

From each non-deprecated expert, extract the `domain` field keywords:

- Split the `domain` string by comma
- Trim whitespace from each keyword
- Lowercase all keywords
- Result: a set of expert domain keywords per candidate

From the incoming proposal, extract domain keywords:

1. Read the proposal document
2. Extract domain terms using the same protocol as `experts/freeform/expert-inference.md` Step 1 (domain, tech stack, complexity signals, key decisions, risk surface)
3. Normalize: lowercase, trim
4. Result: a set of proposal domain keywords

### Step 3: Compute Overlap Score

For each candidate expert, compute a keyword overlap score:

**Formula:**

```
overlap_score = |expert_keywords ∩ proposal_keywords| / |proposal_keywords ∪ expert_keywords|
```

This is the Jaccard similarity coefficient between the two keyword sets.

**Extended scoring** (when the inference process from `experts/freeform/expert-inference.md` is used, apply its weighted scoring instead):

| Signal | Weight | Description |
|--------|--------|-------------|
| domain keyword match | +2 per match | Each keyword shared between expert domain and proposal domain |
| background relevance | +3 | If expert background directly addresses proposal's tech stack (boolean: check whether >= 2 tech stack terms from the proposal appear in the expert's `background` field) |
| review_style compatibility | +1 | If style matches proposal complexity (analytical for complex proposals, pragmatic for straightforward ones — determined by counting complexity signals in the proposal) |

The Jaccard similarity is the primary metric. The extended weighted scoring from `expert-inference.md` is used when the inference agent performs the matching as part of the full inference flow.

### Step 4: Select Best Match

- Sort candidates by overlap score (descending)
- The top candidate is the **best match**
- **Threshold**: overlap score >= 0.3 (Jaccard) OR weighted score >= 5 (extended scoring). Below threshold, no reuse candidate is proposed.
- If multiple candidates share the same top score, prefer the one with more recent `created_at`

### Step 5: Present to User

Present the best match to the user via AskUserQuestion:

- Show the expert profile (front matter summary + persona)
- Show the overlap score and matched keywords
- Provide two options:
  1. **Reuse** — use this existing expert for the freeform review
  2. **Generate new** — proceed to generate a brand new expert profile via `experts/freeform/expert-inference.md`

If the user chooses **Generate new**, proceed with the full inference flow. The newly generated expert is saved alongside existing ones.

If no candidate meets the threshold, skip directly to generating a new expert.

### Interaction with Expert Inference

This reuse matching is the first step in the expert inference process defined in `experts/freeform/expert-inference.md`. The inference prompt's Step 2 ("Check Existing Experts") delegates to these matching rules. When a match is found, the inference prompt presents the reuse option to the user. When no match is found or the user declines reuse, the inference prompt proceeds to Step 3 (generate new expert).

## Quality Tracking

After every evaluation that uses a freeform expert, the system must record whether the expert's contribution produced a substantive change in the rubric outcome.

### When to Record

Record a `review_history` entry after the rubric scoring phase completes (not after the freeform review itself). This ensures the `substantive_change`, `rubric_delta`, and `attack_points_changed` fields reflect actual rubric impact.

### What to Record

Append an entry to the expert's `review_history` array:

```yaml
- proposal: "<path to evaluated proposal>"
  date: "<ISO 8601 date>"
  substantive_change: <true or false>
  rubric_delta: <integer>
  attack_points_changed: <true or false>
```

### Substantive Change Definition

A review produces a **substantive change** when **either** condition is met:

1. **Rubric score delta**: The absolute difference between the rubric score with freeform injection and the rubric score without injection (baseline) is >= 15 points (on the 1000-point scale). This threshold (1.5%) is above the typical variance of LLM repeated runs.
2. **Attack points changed**: The attack points list differs between the injected run and the baseline — at least 1 attack point was added, removed, or modified.

**How to measure**: The baseline is established by running the rubric scorer without freeform injection on the same document. In practice, the baseline is the score from the most recent standard (non-freeform) evaluation of a similar document, or a dedicated baseline run if no prior evaluation exists.

If a baseline is unavailable (first evaluation of this document type with no comparable prior), set `rubric_delta: 0` and `substantive_change: null` (omit or set to `false` after manual review). Do not fabricate a baseline.

### Updating the Expert File

After recording:

1. Read the expert file from `docs/experts/<slug>.md`
2. Parse the YAML front matter
3. Append the new entry to `review_history`
4. Check deprecation rules (see below)
5. If `deprecated` changes to `true`, update the field
6. Write the file back with updated front matter (preserve Markdown body unchanged)

## Auto-Deprecation

An expert is automatically deprecated when it consistently fails to produce substantive changes.

### Trigger Condition

Check after every `review_history` entry is appended. If the **last 3 consecutive entries** all have `substantive_change: false`, set `deprecated: true`.

### Algorithm

```
history = expert.review_history
if length(history) < 3:
    return  // not enough data to deprecate

last_three = history[-3:]
if all(entry.substantive_change == false for entry in last_three):
    expert.deprecated = true
```

### After Deprecation

- Update the expert file's `deprecated` field to `true`
- The expert is excluded from all future reuse matching (filtered out in Step 1 of the reuse flow)
- The expert file is **not deleted** — it remains in `docs/experts/` for historical reference and potential manual reactivation
- Log a message: `"Expert '<slug>' auto-deprecated after 3 consecutive reviews with no substantive change."`

### Reactivation

Auto-deprecated experts can only be reactivated manually (see Manual Deprecation below). There is no automatic reactivation mechanism.

## Manual Deprecation

Users can directly edit any expert file in `docs/experts/` to change its status.

### To Deprecate

Open the expert's `.md` file, change `deprecated: false` to `deprecated: true` in the YAML front matter. The expert will be excluded from reuse matching on the next evaluation run.

### To Reactivate

Open the expert's `.md` file, change `deprecated: true` to `deprecated: false`. The expert will be included in reuse matching again.

### Validation

When loading expert files from `docs/experts/`, skip silently any file that:
- Has invalid or missing YAML front matter
- Is missing required fields (`domain`, `background`, `review_style`, `generated_for`, `created_at`)
- Has `review_history` that is not an array (or missing entirely — treat as empty array)

Do not abort the evaluation if some expert files are malformed. Log a warning and continue with valid experts only.

## Persistence Flow Summary

```
Expert Inference Request
    │
    ├─ Load all non-deprecated experts from docs/experts/
    │   └─ Filter: deprecated == false, valid front matter
    │
    ├─ Compute overlap scores (Jaccard / weighted)
    │
    ├─ Best match >= threshold?
    │   ├─ Yes → Present to user (Reuse / Generate new)
    │   │   ├─ Reuse → Use existing expert profile
    │   │   └─ Generate new → Save new expert to docs/experts/<slug>.md
    │   └─ No → Generate new expert → Save to docs/experts/<slug>.md
    │
    ├─ Freeform review + rubric scoring complete
    │
    └─ Record review_history entry in expert file
        ├─ Compute substantive_change (rubric_delta >= 15 OR attack_points_changed)
        ├─ Append entry to review_history
        └─ Check auto-deprecation (last 3 entries all substantive_change == false)
            └─ If triggered → set deprecated: true
```
