# Eval Skills Adversarial Loop Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add adversarial iteration loop (numeric scoring + scorer/reviser subagents) to `eval-design` and `eval-prd`, and refactor `eval-proposal` to use shared `doc-scorer`/`doc-reviser` agents driven by per-skill `rubric.md` files.

**Architecture:** Three eval skills share one orchestration pattern: main session loops calling `doc-scorer` then `doc-reviser` until target score is reached or iterations exhausted. Each skill owns a `rubric.md` that the generic agents read at runtime — no rubric content is hardcoded in agent prompts.

**Tech Stack:** Markdown skill/agent files, Claude Code plugin system (`plugins/zcode/`)

---

## File Map

| Action | Path | Responsibility |
|--------|------|----------------|
| Create | `plugins/zcode/agents/doc-scorer.md` | Generic scorer — reads rubric, scores doc(s), returns structured output |
| Create | `plugins/zcode/agents/doc-reviser.md` | Generic reviser — reads rubric + eval report, overwrites source doc(s) |
| Create | `plugins/zcode/skills/eval-proposal/templates/rubric.md` | Proposal rubric (extracted from proposal-scorer.md) |
| Create | `plugins/zcode/skills/eval-design/templates/rubric.md` | Design rubric (numeric, 100 pts) |
| Create | `plugins/zcode/skills/eval-prd/templates/rubric.md` | PRD rubric (numeric, 100 pts) |
| Update | `plugins/zcode/skills/eval-proposal/SKILL.md` | Point to doc-scorer/doc-reviser + RUBRIC_PATH |
| Update | `plugins/zcode/skills/eval-design/SKILL.md` | Add adversarial loop, numeric scoring |
| Update | `plugins/zcode/skills/eval-prd/SKILL.md` | Add adversarial loop, numeric scoring |
| Update | `plugins/zcode/skills/eval-design/templates/report.md` | Match eval-proposal scorecard format |
| Update | `plugins/zcode/skills/eval-prd/templates/report.md` | Match eval-proposal scorecard format |
| Delete | `plugins/zcode/agents/proposal-scorer.md` | Replaced by doc-scorer |
| Delete | `plugins/zcode/agents/proposal-reviser.md` | Replaced by doc-reviser |

---

## Task 1: Create `doc-scorer` agent

**Files:**
- Create: `plugins/zcode/agents/doc-scorer.md`

- [ ] **Step 1: Create the file**

```
---
name: doc-scorer
description: "Generic document scorer. Reads a rubric file and source documents, scores on 100-point scale, returns structured output the orchestrator parses."
model: sonnet
color: yellow
memory: project
inputs:
  - name: DOC_PATHS
    description: Comma-separated paths to documents to evaluate (skip paths that don't exist)
    required: true
  - name: RUBRIC_PATH
    description: Path to the rubric.md file containing scoring dimensions and criteria
    required: true
  - name: REPORT_PATH
    description: Output path for the evaluation report
    required: true
  - name: ITERATION
    description: Current iteration number (1 = first evaluation)
    required: true
  - name: PREVIOUS_REPORT_PATH
    description: Path to previous iteration's report (only for iteration > 1)
    required: false
---

You are a harsh document evaluator. Score on a 100-point scale. Be critical — find every weakness.

<EXTREMELY-IMPORTANT>
1. You are the ADVERSARY — find flaws, not reasons to be generous
2. Every point deducted must have a concrete reason with a quote from the document
3. Never give full marks unless content is genuinely excellent
4. Return output in the EXACT format specified below — the orchestrator parses it mechanically
</EXTREMELY-IMPORTANT>

## Workflow

### Step 1: Read Inputs

Read each path in `{{DOC_PATHS}}` (comma-separated). Skip any path that does not exist on disk.

Read the rubric at `{{RUBRIC_PATH}}` — it defines scoring dimensions, point allocations, criteria, and the report template path.

If `{{ITERATION}}` > 1, read `{{PREVIOUS_REPORT_PATH}}` to check which issues were addressed.

### Step 2: Score

Apply the rubric to each dimension. Justify every deduction with a specific quote or observation from the document.

<HARD-RULE>
Score independently. Do NOT give credit for "effort" or "improvement from last iteration". Score only what is on the page right now.
</HARD-RULE>

### Step 3: Write Report

The rubric specifies a report template path. Read that template, fill it in, and write to `{{REPORT_PATH}}`.

### Step 4: Return Summary

<HARD-GATE>
Return output in EXACTLY this format. No extra text before or after.
</HARD-GATE>

SCORE: {{total}}/100
DIMENSIONS:
  {{dimension_name}}: {{score}}/{{max}}
  {{dimension_name}}: {{score}}/{{max}}
  ...
ATTACKS:
1. [dimension]: [specific weakness] — [quote from document] — [what must improve]
2. [dimension]: [specific weakness] — [quote from document] — [what must improve]
3. [dimension]: [specific weakness] — [quote from document] — [what must improve]
```

- [ ] **Step 2: Commit**

```bash
git add plugins/zcode/agents/doc-scorer.md
git commit -m "feat(agents): add generic doc-scorer agent"
```

---

## Task 2: Create `doc-reviser` agent

**Files:**
- Create: `plugins/zcode/agents/doc-reviser.md`

- [ ] **Step 1: Create the file**

```
---
name: doc-reviser
description: "Generic document reviser. Reads rubric + eval report, rewrites source doc(s) to address attack points. No padding."
model: sonnet
color: cyan
memory: project
inputs:
  - name: DOC_PATHS
    description: Comma-separated paths to documents to revise (overwrite in place)
    required: true
  - name: RUBRIC_PATH
    description: Path to the rubric.md file — used to understand what "good" looks like
    required: true
  - name: EVAL_REPORT_PATH
    description: Path to the evaluation report containing scores and attack points
    required: true
  - name: ATTACK_POINTS
    description: The top 3 attack points from the scorer (newline-separated)
    required: true
---

You are revising document(s) to address specific critique. Improve to score higher, without inflating or padding.

<EXTREMELY-IMPORTANT>
1. Address EACH attack point specifically — do not dodge or wave hands
2. Concise and concrete beats verbose and vague
3. Keep what's already good — only change what the critique targets
4. Maximum 3 rounds of self-review before delivering
</EXTREMELY-IMPORTANT>

## Workflow

### Step 1: Read Inputs

Read each path in `{{DOC_PATHS}}` (comma-separated). Skip any path that does not exist.

Read the rubric at `{{RUBRIC_PATH}}` to understand what a high-scoring document looks like.

Read the evaluation report at `{{EVAL_REPORT_PATH}}`.

<HARD-RULE>
Do NOT skip reading the eval report. The attack points tell you exactly what to fix. Fixing things that scored well wastes the iteration.
</HARD-RULE>

### Step 2: Revise

| Attack Type | Fix Strategy |
|-------------|-------------|
| Vague language | Replace with concrete, quantified statements |
| Missing section | Add real content, not placeholder text |
| Inconsistency | Align scope, solution, and success criteria |
| Weak alternatives | Add honest pros/cons with rationale |
| Unmeasurable criteria | Rewrite as testable, verifiable claims |

<HARD-GATE>
Do NOT add length for the sake of length. Every new sentence must fix a weakness the scorer identified.
</HARD-GATE>

### Step 3: Write & Report

Overwrite each document in `{{DOC_PATHS}}` with the revised content.

Return:

REVISED: {{DOC_PATHS}}
CHANGES:
- [what changed] → [why: which attack point it addresses]
- [what changed] → [why: which attack point it addresses]
- [what changed] → [why: which attack point it addresses]

## Quality Checks

<EXTREMELY-IMPORTANT>
1. Every attack point from the scorer has been addressed
2. No new vague language introduced ("better", "improved", "enhanced" without quantification)
3. Documents are internally consistent after revision
4. Total word count did not increase by more than 30% (padding check)
</EXTREMELY-IMPORTANT>

## Attack Points

{{ATTACK_POINTS}}
```

- [ ] **Step 2: Commit**

```bash
git add plugins/zcode/agents/doc-reviser.md
git commit -m "feat(agents): add generic doc-reviser agent"
```

---

## Task 3: Create proposal rubric

**Files:**
- Create: `plugins/zcode/skills/eval-proposal/templates/rubric.md`

- [ ] **Step 1: Create the file** (extract scoring criteria from `plugins/zcode/agents/proposal-scorer.md`)

```markdown
# Proposal Evaluation Rubric

**Total: 100 points**
**Report template:** `plugins/zcode/skills/eval-proposal/templates/report.md`

## Dimensions

### 1. Problem Definition (20 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Problem stated clearly | 0-7 | Is the core problem unambiguous? Could two readers interpret it differently? |
| Evidence provided | 0-7 | Is there data, user feedback, or concrete examples backing the problem? Not just "we think..." |
| Urgency justified | 0-6 | Why solve this now? What happens if we don't? |

### 2. Solution Clarity (20 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Approach is concrete | 0-7 | Can a reader explain back what will be built? Or is it vague hand-waving? |
| User-facing behavior described | 0-7 | What does the end user experience? Not internals — the observable behavior |
| Distinguishes from alternatives | 0-6 | Is it clear why this approach over others? What's the differentiator? |

### 3. Alternatives Analysis (15 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| At least 2 alternatives listed | 0-5 | Including "do nothing" as a valid alternative |
| Pros/cons for each | 0-5 | Are trade-offs honest? Not straw-man arguments? |
| Rationale for chosen approach | 0-5 | Is the verdict justified against the alternatives? |

### 4. Scope Definition (15 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| In-scope items are concrete | 0-5 | Each item is a deliverable, not a vague area |
| Out-of-scope explicitly listed | 0-5 | Are deferred items named, not just implied? |
| Scope is bounded | 0-5 | Can a team execute this in a defined timeframe? Or is it open-ended? |

### 5. Risk Assessment (15 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Risks identified | 0-5 | At least 3 meaningful risks, not trivial ones |
| Likelihood + impact rated | 0-5 | Is the assessment honest? Not all "low likelihood, high impact"? |
| Mitigations are actionable | 0-5 | Can someone act on the mitigation? Or is it "we'll handle it"? |

### 6. Success Criteria (15 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Criteria are measurable | 0-5 | Can you objectively verify each criterion? "Works well" is not measurable |
| Coverage is complete | 0-5 | Do criteria cover all in-scope items? Any gaps? |
| Criteria are testable | 0-5 | Could you write a test or checklist for each criterion? |

## Deduction Rules

- **Vague language penalty**: -2 per instance of "better", "improved", "enhanced" without quantification
- **Missing section penalty**: 0 points for that dimension
- **Inconsistency penalty**: -3 if scope contradicts solution or success criteria don't cover scope
```

- [ ] **Step 2: Commit**

```bash
git add plugins/zcode/skills/eval-proposal/templates/rubric.md
git commit -m "feat(eval-proposal): extract rubric to rubric.md"
```

---

## Task 4: Create design rubric

**Files:**
- Create: `plugins/zcode/skills/eval-design/templates/rubric.md`

- [ ] **Step 1: Create the file**

```markdown
# Design Evaluation Rubric

**Total: 100 points**
**Report template:** `plugins/zcode/skills/eval-design/templates/report.md`

## Required Sections

Mark missing required sections as 0 pts for that dimension:

| Section | Required |
|---------|----------|
| Overview + tech stack | ✓ |
| Architecture (layer + diagram) | ✓ |
| Interfaces | ✓ |
| Data Models | ✓ |
| Error Handling | ✓ |
| Testing Strategy | ✓ |
| Security Considerations | ○ (required if PRD has auth/data requirements) |

## Dimensions

### 1. Architecture Clarity (20 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Layer placement explicit | 0-7 | Does the doc state which layer (API/service/repo/etc.) this belongs to? |
| Component diagram present | 0-7 | Is there an ASCII or text diagram showing components and relationships? |
| Dependencies listed | 0-6 | Are internal modules and external packages named? |

### 2. Interface & Model Definitions (20 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Interface signatures typed | 0-7 | Do all interfaces have typed params and return values (not prose)? |
| Models concrete | 0-7 | Are all model fields named with types and constraints (not just described)? |
| Directly implementable | 0-6 | Can a developer code from this without guessing any types or shapes? |

### 3. Error Handling (15 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Error types defined | 0-5 | Are custom error types or error codes explicitly defined? |
| Propagation strategy clear | 0-5 | Is there a stated strategy for how errors flow between layers? |
| HTTP status codes mapped | 0-5 | If API: are error types mapped to HTTP status codes? |

### 4. Testing Strategy (15 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Per-layer test plan | 0-5 | Does each layer have a stated test approach (unit/integration/e2e)? |
| Coverage target numeric | 0-5 | Is there a numeric coverage target (e.g., 80%)? |
| Test tooling named | 0-5 | Are specific test libraries/frameworks named? |

### 5. Breakdown-Readiness ★ (20 pts — critical gate)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Components enumerable | 0-7 | Can you list and count all components/modules? Or are they described vaguely? |
| Tasks derivable | 0-7 | Does each interface → at least one impl task? Each model → at least one schema task? |
| PRD AC coverage | 0-6 | If PRD exists: are all acceptance criteria addressed somewhere in the design? |

★ This dimension is the direct gate to `/breakdown-tasks`. A score below 12/20 blocks progression.

### 6. Security Considerations (10 pts)

Only scored if PRD has auth, data privacy, or multi-user requirements. Mark N/A (full credit) otherwise.

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Threat model present | 0-5 | Are specific threats identified (not just "we'll add auth")? |
| Mitigations concrete | 0-5 | Is each threat paired with a specific countermeasure? |

## Deduction Rules

- **Missing required section**: 0 pts for that dimension
- **Prose-only (no code/diagram where expected)**: -5 pts from that dimension
- **PRD AC gap**: -3 pts per unaddressed acceptance criterion (from Breakdown-Readiness)
- **Placeholder text ("TBD", "TODO")**: -2 pts per instance
```

- [ ] **Step 2: Commit**

```bash
git add plugins/zcode/skills/eval-design/templates/rubric.md
git commit -m "feat(eval-design): add numeric rubric.md"
```

---

## Task 5: Create PRD rubric

**Files:**
- Create: `plugins/zcode/skills/eval-prd/templates/rubric.md`

- [ ] **Step 1: Create the file**

```markdown
# PRD Evaluation Rubric

**Total: 100 points**
**Report template:** `plugins/zcode/skills/eval-prd/templates/report.md`

## Required Sections (prd-spec.md)

| Section | Required |
|---------|----------|
| 需求背景（原因/对象/人员） | ✓ |
| 需求目标 + 量化指标 | ✓ |
| Scope（In + Out） | ✓ |
| 流程说明 + Mermaid 流程图 | ✓ |
| 功能描述 | ✓ |

## Required Sections (prd-user-stories.md)

| Section | Required |
|---------|----------|
| User Stories | ✓ |
| Acceptance Criteria (Given/When/Then) | ✓ |

## Dimensions

### 1. Background & Goals (20 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Background has three elements (原因/对象/人员) | 0-7 | Are all three present and specific? |
| Goals are quantified | 0-7 | Is there at least one numeric target (%, count, time)? |
| Background and goals are logically consistent | 0-6 | Does the goal follow from the stated problem? |

### 2. Flow Diagrams (20 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Mermaid diagram exists | 0-7 | Is there at least one Mermaid flowchart? |
| Main path complete (start → end) | 0-7 | Does the diagram cover the full happy path? |
| Decision points + error branches covered | 0-6 | Are there diamond nodes and at least one error/exception branch? |

### 3. Functional Specs (20 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Tables complete (list page 7 elements, button 4 elements, form 2 elements) | 0-7 | Are all required table columns filled in? |
| Field descriptions clear | 0-7 | Is each field's purpose, type, and source stated? |
| Validation rules explicit | 0-6 | Are validation rules stated per field/button (not just "validate input")? |

### 4. User Stories (20 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| Coverage: one story per target user | 0-7 | Does every user type from the background section have at least one story? |
| Format correct (As a / I want / So that) | 0-7 | Do all stories follow the format? Are actions concrete (not "manage", "handle")? |
| AC per story (Given/When/Then) | 0-6 | Does every story have at least one AC in Given/When/Then format? |

### 5. Scope Clarity (20 pts)

| Criterion | Points | What to check |
|-----------|--------|---------------|
| In-scope items are concrete deliverables | 0-7 | Each item is a specific feature/screen/API, not a vague area |
| Out-of-scope explicitly lists deferred items | 0-7 | Are deferred items named, not just implied by absence? |
| Scope consistent with functional specs and user stories | 0-6 | Do the in-scope items match what's described in 功能描述 and user stories? |

## Deduction Rules

- **Missing required section**: 0 pts for that dimension
- **Vague language without quantification**: -2 pts per instance ("better UX", "faster", "improved")
- **Inconsistency between sections**: -3 pts per conflict (e.g., scope says X is out but functional spec describes X)
- **Placeholder text ("TBD", "TODO")**: -2 pts per instance
```

- [ ] **Step 2: Commit**

```bash
git add plugins/zcode/skills/eval-prd/templates/rubric.md
git commit -m "feat(eval-prd): add numeric rubric.md"
```

---

## Task 6: Update report templates

**Files:**
- Modify: `plugins/zcode/skills/eval-design/templates/report.md`
- Modify: `plugins/zcode/skills/eval-prd/templates/report.md`

Both templates are updated to match the eval-proposal scorecard format (ASCII table + Deductions + Attack Points + Previous Issues Check + Verdict).

- [ ] **Step 1: Replace `eval-design/templates/report.md`**

```markdown
---
date: "{{DATE}}"
design_path: "{{DESIGN_PATH}}"
prd_path: "{{PRD_PATH}}"
iteration: "{{ITERATION}}"
target_score: "{{TARGET}}"
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration {{ITERATION}}

**Score: {{SCORE}}/100** (target: {{TARGET}})

\`\`\`
┌─────────────────────────────────────────────────────────────────┐
│                     DESIGN QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Architecture Clarity      │  ___     │  20      │ ✅/⚠️/❌    │
│    Layer placement explicit  │  ___/7   │          │            │
│    Component diagram present │  ___/7   │          │            │
│    Dependencies listed       │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Interface & Model Defs    │  ___     │  20      │ ✅/⚠️/❌    │
│    Interface signatures typed│  ___/7   │          │            │
│    Models concrete           │  ___/7   │          │            │
│    Directly implementable    │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Error Handling            │  ___     │  15      │ ✅/⚠️/❌    │
│    Error types defined       │  ___/5   │          │            │
│    Propagation strategy clear│  ___/5   │          │            │
│    HTTP status codes mapped  │  ___/5   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Testing Strategy          │  ___     │  15      │ ✅/⚠️/❌    │
│    Per-layer test plan       │  ___/5   │          │            │
│    Coverage target numeric   │  ___/5   │          │            │
│    Test tooling named        │  ___/5   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Breakdown-Readiness ★     │  ___     │  20      │ ✅/⚠️/❌    │
│    Components enumerable     │  ___/7   │          │            │
│    Tasks derivable           │  ___/7   │          │            │
│    PRD AC coverage           │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Security Considerations   │  ___     │  10      │ ✅/⚠️/N/A  │
│    Threat model present      │  ___/5   │          │            │
│    Mitigations concrete      │  ___/5   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  ___     │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
\`\`\`

★ Breakdown-Readiness < 12/20 blocks progression to `/breakdown-tasks`

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| <!-- section:line --> | <!-- vague language / missing content / inconsistency --> | <!-- -N pts --> |

---

## Attack Points

### Attack 1: [dimension — specific weakness]

**Where**: <!-- quote from document -->
**Why it's weak**: <!-- concrete critique -->
**What must improve**: <!-- actionable fix -->

### Attack 2: [dimension — specific weakness]

**Where**: <!-- quote from document -->
**Why it's weak**: <!-- concrete critique -->
**What must improve**: <!-- actionable fix -->

### Attack 3: [dimension — specific weakness]

**Where**: <!-- quote from document -->
**Why it's weak**: <!-- concrete critique -->
**What must improve**: <!-- actionable fix -->

---

## Previous Issues Check

<!-- Only for iteration > 1 -->

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| <!-- attack from iter N-1 --> | ✅/❌ | <!-- what changed / what didn't --> |

---

## Verdict

- **Score**: {{SCORE}}/100
- **Target**: {{TARGET}}/100
- **Gap**: {{GAP}} points
- **Breakdown-Readiness**: {{BR_SCORE}}/20 — {{can/cannot proceed to /breakdown-tasks}}
- **Action**: {{Continue to iteration N+1 / Target reached / Iterations exhausted}}
```

- [ ] **Step 2: Replace `eval-prd/templates/report.md`**

```markdown
---
date: "{{DATE}}"
prd_path: "{{PRD_PATH}}"
user_stories_path: "{{USER_STORIES_PATH}}"
iteration: "{{ITERATION}}"
target_score: "{{TARGET}}"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration {{ITERATION}}

**Score: {{SCORE}}/100** (target: {{TARGET}})

\`\`\`
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Background & Goals        │  ___     │  20      │ ✅/⚠️/❌    │
│    Background three elements │  ___/7   │          │            │
│    Goals quantified          │  ___/7   │          │            │
│    Logical consistency       │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Flow Diagrams             │  ___     │  20      │ ✅/⚠️/❌    │
│    Mermaid diagram exists    │  ___/7   │          │            │
│    Main path complete        │  ___/7   │          │            │
│    Decision + error branches │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Functional Specs          │  ___     │  20      │ ✅/⚠️/❌    │
│    Tables complete           │  ___/7   │          │            │
│    Field descriptions clear  │  ___/7   │          │            │
│    Validation rules explicit │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. User Stories              │  ___     │  20      │ ✅/⚠️/❌    │
│    Coverage per user type    │  ___/7   │          │            │
│    Format correct            │  ___/7   │          │            │
│    AC per story              │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Scope Clarity             │  ___     │  20      │ ✅/⚠️/❌    │
│    In-scope concrete         │  ___/7   │          │            │
│    Out-of-scope explicit     │  ___/7   │          │            │
│    Consistent with specs     │  ___/6   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  ___     │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
\`\`\`

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| <!-- section:line --> | <!-- vague language / missing content / inconsistency --> | <!-- -N pts --> |

---

## Attack Points

### Attack 1: [dimension — specific weakness]

**Where**: <!-- quote from document -->
**Why it's weak**: <!-- concrete critique -->
**What must improve**: <!-- actionable fix -->

### Attack 2: [dimension — specific weakness]

**Where**: <!-- quote from document -->
**Why it's weak**: <!-- concrete critique -->
**What must improve**: <!-- actionable fix -->

### Attack 3: [dimension — specific weakness]

**Where**: <!-- quote from document -->
**Why it's weak**: <!-- concrete critique -->
**What must improve**: <!-- actionable fix -->

---

## Previous Issues Check

<!-- Only for iteration > 1 -->

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| <!-- attack from iter N-1 --> | ✅/❌ | <!-- what changed / what didn't --> |

---

## Verdict

- **Score**: {{SCORE}}/100
- **Target**: {{TARGET}}/100
- **Gap**: {{GAP}} points
- **Action**: {{Continue to iteration N+1 / Target reached / Iterations exhausted}}
```

- [ ] **Step 3: Commit**

```bash
git add plugins/zcode/skills/eval-design/templates/report.md plugins/zcode/skills/eval-prd/templates/report.md
git commit -m "feat(eval-design,eval-prd): update report templates to numeric scorecard format"
```

---

## Task 7: Update `eval-proposal/SKILL.md`

**Files:**
- Modify: `plugins/zcode/skills/eval-proposal/SKILL.md`

Changes: swap `proposal-scorer` → `doc-scorer`, `proposal-reviser` → `doc-reviser`, add `RUBRIC_PATH` input to both agent invocations.

- [ ] **Step 1: Replace the file**

```markdown
---
name: eval-proposal
description: Evaluate a proposal document with 100-point scoring, then run adversarial iterations until target score is met. Main session orchestrates doc-scorer and doc-reviser subagents. Specify target score and max iterations.
---

# Eval Proposal

评估 proposal 文档质量（百分制），通过 doc-scorer / doc-reviser subagent 多轮对抗迭代，直到达到目标分数或耗尽迭代次数。

**架构**：主会话调度循环，独立 subagent 负责评分和修订。

## When to Use

**Trigger:**

- User says yes to adversarial eval prompt after `/brainstorm`
- User provides `/eval-proposal` command
- User wants iterative refinement: `/eval-proposal --target 85 --iterations 5`

**Skip:**

- No proposal document exists (use `/brainstorm` first)
- Requirements are already in PRD form (use `/eval-prd` instead)

## Parameters

| Parameter      | Default | Description                                           |
| -------------- | ------- | ----------------------------------------------------- |
| `--target`     | 80      | Target score (0-100). Loop stops when score >= target |
| `--iterations` | 3       | Max adversarial iterations                            |

Parse from user input. Examples:

- `/eval-proposal` → target=80, iterations=3
- `/eval-proposal --target 90` → target=90, iterations=3
- `/eval-proposal --target 85 --iterations 5` → target=85, iterations=5

## Architecture

```
Main Session (orchestrator)
  │
  ├─ iteration 1:
  │   ├── Agent (doc-scorer)  ──→ score + attack points
  │   ├── score >= target? ──→ yes: stop
  │   └── Agent (doc-reviser) ──→ revised proposal
  │
  ├─ iteration 2:
  │   ├── Agent (doc-scorer)  ──→ score + attack points  (blind to changes)
  │   ├── score >= target? ──→ yes: stop
  │   └── Agent (doc-reviser) ──→ revised proposal
  │
  ├─ ... (loop)
  │
  └─ Final report to user
```

<EXTREMELY-IMPORTANT>
Scorer and reviser are **independent subagents** defined in `plugins/zcode/agents/`. Do NOT inline their prompts. Invoke them via Agent tool with the correct inputs.
</EXTREMELY-IMPORTANT>

## Step 1: Locate Proposal

Check in order:

1. Path provided by user
2. `docs/proposals/<slug>/proposal.md` — find latest by modification time
3. `Glob` `docs/proposals/*/proposal.md` and list options to user
4. Ask user for path if not found

Determine `<slug>` from path (e.g., `docs/proposals/eval-proposal/proposal.md` → slug is `eval-proposal`).

## Step 2: Invoke Scorer Subagent

Spawn `doc-scorer` agent via **Agent tool** (subagent_type: `zcode:doc-scorer` if registered, otherwise `general-purpose`).

<HARD-RULE>
Pass these inputs to the scorer:
- `DOC_PATHS` = `docs/proposals/<slug>/proposal.md`
- `RUBRIC_PATH` = `plugins/zcode/skills/eval-proposal/templates/rubric.md`
- `REPORT_PATH` = `docs/proposals/<slug>/eval-iteration-{{N}}.md`
- `ITERATION` = current iteration number (1-based)
- `PREVIOUS_REPORT_PATH` = previous iteration report path (only if iteration > 1)
</HARD-RULE>

After the scorer returns, **parse its output in the main session**:

1. Extract `SCORE: X/100` line
2. Extract per-dimension scores from `DIMENSIONS:` section
3. Extract attack points from `ATTACKS:` section
4. Record score in iteration tracker

## Step 3: Decision Gate (Main Session)

<HARD-GATE>
This decision is made in the MAIN SESSION, not delegated to a subagent. The orchestrator (you) controls the loop.
</HARD-GATE>

| Condition                                  | Action                          |
| ------------------------------------------ | ------------------------------- |
| Score >= target                            | Skip to Step 6 (final report)   |
| Score < target AND iterations remaining    | Proceed to Step 4               |
| Score < target AND no iterations remaining | Skip to Step 6 (report failure) |

Report current status to user:

```
Iteration {{N}}/{{MAX}}: scored {{SCORE}}/100 (target: {{TARGET}}). Revision subagent starting...
```

## Step 4: Invoke Reviser Subagent

Spawn `doc-reviser` agent via **Agent tool** (subagent_type: `zcode:doc-reviser` if registered, otherwise `general-purpose`).

<HARD-RULE>
Pass these inputs to the reviser:
- `DOC_PATHS` = `docs/proposals/<slug>/proposal.md`
- `RUBRIC_PATH` = `plugins/zcode/skills/eval-proposal/templates/rubric.md`
- `EVAL_REPORT_PATH` = `docs/proposals/<slug>/eval-iteration-{{N}}.md`
- `ATTACK_POINTS` = the 3 attack points extracted from scorer output
</HARD-RULE>

The reviser will overwrite the proposal file in place.

## Step 5: Loop

Increment iteration counter. Return to Step 2.

<EXTREMELY-IMPORTANT>
The scorer must NEVER be told what changes the reviser made. It evaluates the proposal as-is. Only the `PREVIOUS_REPORT_PATH` input carries forward for "previous issues addressed" checking.
</EXTREMELY-IMPORTANT>

## Step 6: Final Report (Main Session)

When the loop ends, assemble and report to the user:

```
## Eval-Proposal Complete

**Final Score**: {{SCORE}}/100 (target: {{TARGET}})
**Iterations Used**: {{N}}/{{MAX}}

### Score Progression
| Iteration | Score | Delta |
|-----------|-------|-------|
| 1 | {{s1}} | - |
| 2 | {{s2}} | +{{d2}} |
| ... | ... | ... |

### Dimension Breakdown (final)
| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | {{d1}} | 20 |
| Solution Clarity | {{d2}} | 20 |
| Alternatives Analysis | {{d3}} | 15 |
| Scope Definition | {{d4}} | 15 |
| Risk Assessment | {{d5}} | 15 |
| Success Criteria | {{d6}} | 15 |

### Outcome
{{"Target reached" / "Target NOT reached — N iterations exhausted"}}
{{If not reached: "Largest gaps: [dimension names]. Consider manual revision or increasing iterations."}}
```

Save the final report to `docs/proposals/<slug>/eval-report.md`.

## Report Path Convention

| File               | Path                                            |
| ------------------ | ----------------------------------------------- |
| Iteration N report | `docs/proposals/<slug>/eval-iteration-{{N}}.md` |
| Final report       | `docs/proposals/<slug>/eval-report.md`          |

## Related

- `/brainstorm` — Creates or revises the proposal document (runs in main session)
```

- [ ] **Step 2: Commit**

```bash
git add plugins/zcode/skills/eval-proposal/SKILL.md
git commit -m "feat(eval-proposal): switch to doc-scorer/doc-reviser with RUBRIC_PATH"
```

---

## Task 8: Update `eval-design/SKILL.md`

**Files:**
- Modify: `plugins/zcode/skills/eval-design/SKILL.md`

Full rewrite: replace single-pass agent with adversarial loop, add `--target`/`--iterations` params, pass `DOC_PATHS` + `RUBRIC_PATH`.

- [ ] **Step 1: Replace the file**

```markdown
---
name: eval-design
description: Evaluate a tech design document with 100-point scoring, then run adversarial iterations until target score is met. Main session orchestrates doc-scorer and doc-reviser subagents.
---

# Eval Design

评估 tech-design.md 文档质量（百分制），通过 doc-scorer / doc-reviser subagent 多轮对抗迭代，直到达到目标分数或耗尽迭代次数。重点检查能否直接驱动 `/breakdown-tasks`。

**架构**：主会话调度循环，独立 subagent 负责评分和修订。

## When to Use

**Trigger:**
- User asks to "evaluate design" or "check design quality"
- User provides `/eval-design` command
- Before handing off tech-design.md to `/breakdown-tasks`

**Skip:**
- design.md doesn't exist yet (use `/design-tech` first)

## Parameters

| Parameter      | Default | Description                                           |
| -------------- | ------- | ----------------------------------------------------- |
| `--target`     | 80      | Target score (0-100). Loop stops when score >= target |
| `--iterations` | 3       | Max adversarial iterations                            |

## Architecture

```
Main Session (orchestrator)
  │
  ├─ iteration N:
  │   ├── Agent (doc-scorer)  ──→ score + attack points
  │   ├── score >= target? ──→ yes: stop
  │   └── Agent (doc-reviser) ──→ revised design doc(s)
  │
  └─ Final report to user
```

<EXTREMELY-IMPORTANT>
Scorer and reviser are **independent subagents**. Do NOT inline their prompts. Invoke via Agent tool with the correct inputs.
</EXTREMELY-IMPORTANT>

## Step 1: Locate Design Documents

Check in order:
1. Path provided by user
2. Read `docs/features/<current-feature>/manifest.md` → locate design documents
3. Fall back to `design/tech-design.md`, `design/api-handbook.md`, `ui/ui-design.md`
4. Ask user for path if not found

Determine `<feature-slug>` from the path. Build `DOC_PATHS` as a comma-separated list of all design files that exist on disk (tech-design.md, api-handbook.md, ui-design.md).

## Step 2: Invoke Scorer Subagent

Spawn `doc-scorer` agent via **Agent tool** (subagent_type: `zcode:doc-scorer` if registered, otherwise `general-purpose`).

<HARD-RULE>
Pass these inputs to the scorer:
- `DOC_PATHS` = comma-separated paths of existing design files
- `RUBRIC_PATH` = `plugins/zcode/skills/eval-design/templates/rubric.md`
- `REPORT_PATH` = `docs/features/<slug>/design-eval-iteration-{{N}}.md`
- `ITERATION` = current iteration number (1-based)
- `PREVIOUS_REPORT_PATH` = previous iteration report path (only if iteration > 1)
</HARD-RULE>

After the scorer returns, parse its output in the main session:
1. Extract `SCORE: X/100`
2. Extract per-dimension scores from `DIMENSIONS:` section
3. Extract attack points from `ATTACKS:` section
4. Record score in iteration tracker

## Step 3: Decision Gate (Main Session)

<HARD-GATE>
This decision is made in the MAIN SESSION, not delegated to a subagent.
</HARD-GATE>

| Condition                                  | Action                          |
| ------------------------------------------ | ------------------------------- |
| Score >= target                            | Skip to Step 6 (final report)   |
| Score < target AND iterations remaining    | Proceed to Step 4               |
| Score < target AND no iterations remaining | Skip to Step 6 (report failure) |

Report current status to user:

```
Iteration {{N}}/{{MAX}}: scored {{SCORE}}/100 (target: {{TARGET}}). Revision subagent starting...
```

## Step 4: Invoke Reviser Subagent

Spawn `doc-reviser` agent via **Agent tool** (subagent_type: `zcode:doc-reviser` if registered, otherwise `general-purpose`).

<HARD-RULE>
Pass these inputs to the reviser:
- `DOC_PATHS` = same comma-separated paths as scorer
- `RUBRIC_PATH` = `plugins/zcode/skills/eval-design/templates/rubric.md`
- `EVAL_REPORT_PATH` = `docs/features/<slug>/design-eval-iteration-{{N}}.md`
- `ATTACK_POINTS` = the 3 attack points extracted from scorer output
</HARD-RULE>

The reviser will overwrite the design file(s) in place.

## Step 5: Loop

Increment iteration counter. Return to Step 2.

<EXTREMELY-IMPORTANT>
The scorer must NEVER be told what changes the reviser made. Only `PREVIOUS_REPORT_PATH` carries forward.
</EXTREMELY-IMPORTANT>

## Step 6: Final Report (Main Session)

```
## Eval-Design Complete

**Final Score**: {{SCORE}}/100 (target: {{TARGET}})
**Iterations Used**: {{N}}/{{MAX}}

### Score Progression
| Iteration | Score | Delta |
|-----------|-------|-------|
| 1 | {{s1}} | - |
| 2 | {{s2}} | +{{d2}} |

### Dimension Breakdown (final)
| Dimension | Score | Max |
|-----------|-------|-----|
| Architecture Clarity | {{d1}} | 20 |
| Interface & Model Definitions | {{d2}} | 20 |
| Error Handling | {{d3}} | 15 |
| Testing Strategy | {{d4}} | 15 |
| Breakdown-Readiness ★ | {{d5}} | 20 |
| Security Considerations | {{d6}} | 10 |

### Outcome
{{"Target reached" / "Target NOT reached — N iterations exhausted"}}
{{Breakdown-Readiness: {{score}}/20 — can/cannot proceed to /breakdown-tasks}}
{{If not reached: "Largest gaps: [dimension names]. Consider manual revision or increasing iterations."}}
```

Save the final report to `docs/features/<slug>/design-eval.md`.

## Report Path Convention

| File               | Path                                                    |
| ------------------ | ------------------------------------------------------- |
| Iteration N report | `docs/features/<slug>/design-eval-iteration-{{N}}.md`  |
| Final report       | `docs/features/<slug>/design-eval.md`                  |

## Related

- `/design-tech` — Create or revise the design.md
- `/eval-prd` — Evaluate PRD before design starts
- `/breakdown-tasks` — Next step after design passes evaluation
```

- [ ] **Step 2: Commit**

```bash
git add plugins/zcode/skills/eval-design/SKILL.md
git commit -m "feat(eval-design): add adversarial iteration loop with doc-scorer/doc-reviser"
```

---

## Task 9: Update `eval-prd/SKILL.md`

**Files:**
- Modify: `plugins/zcode/skills/eval-prd/SKILL.md`

Full rewrite: replace single-pass agent with adversarial loop, add `--target`/`--iterations` params, pass `DOC_PATHS` + `RUBRIC_PATH`.

- [ ] **Step 1: Replace the file**

```markdown
---
name: eval-prd
description: Evaluate a PRD document with 100-point scoring, then run adversarial iterations until target score is met. Main session orchestrates doc-scorer and doc-reviser subagents.
---

# Eval PRD

评估 PRD 文档质量（百分制），通过 doc-scorer / doc-reviser subagent 多轮对抗迭代，直到达到目标分数或耗尽迭代次数。

**架构**：主会话调度循环，独立 subagent 负责评分和修订。

## Prerequisites

检查上一阶段产物，缺失则中止并提示用户：

| 产物 | 缺失时提示 |
|------|-----------|
| `prd/prd-spec.md` | 先执行 `/write-prd` |
| `prd/prd-user-stories.md` | 先执行 `/write-prd` |

## When to Use

**Trigger:**
- User asks to "evaluate PRD" or "check PRD quality"
- User provides `/eval-prd` command
- Before handing off PRD to `/design-tech` or `/ui-design`

**Skip:**
- PRD doesn't exist yet (use `/write-prd` first)

## Parameters

| Parameter      | Default | Description                                           |
| -------------- | ------- | ----------------------------------------------------- |
| `--target`     | 80      | Target score (0-100). Loop stops when score >= target |
| `--iterations` | 3       | Max adversarial iterations                            |

## Architecture

```
Main Session (orchestrator)
  │
  ├─ iteration N:
  │   ├── Agent (doc-scorer)  ──→ score + attack points
  │   ├── score >= target? ──→ yes: stop
  │   └── Agent (doc-reviser) ──→ revised PRD doc(s)
  │
  └─ Final report to user
```

<EXTREMELY-IMPORTANT>
Scorer and reviser are **independent subagents**. Do NOT inline their prompts. Invoke via Agent tool with the correct inputs.
</EXTREMELY-IMPORTANT>

## Step 1: Locate Documents

Check in order:
1. Path provided by user
2. Read `docs/features/<current-feature>/manifest.md` → locate PRD documents
3. Fall back to `docs/features/<current-feature>/prd/prd-spec.md` + `prd/prd-user-stories.md`
4. Ask user for path if not found

Determine `<feature-slug>` from the path. Build `DOC_PATHS` as a comma-separated list of all PRD files that exist on disk (prd-spec.md, prd-user-stories.md, prd-ui-functions.md).

## Step 2: Invoke Scorer Subagent

Spawn `doc-scorer` agent via **Agent tool** (subagent_type: `zcode:doc-scorer` if registered, otherwise `general-purpose`).

<HARD-RULE>
Pass these inputs to the scorer:
- `DOC_PATHS` = comma-separated paths of existing PRD files
- `RUBRIC_PATH` = `plugins/zcode/skills/eval-prd/templates/rubric.md`
- `REPORT_PATH` = `docs/features/<slug>/prd-eval-iteration-{{N}}.md`
- `ITERATION` = current iteration number (1-based)
- `PREVIOUS_REPORT_PATH` = previous iteration report path (only if iteration > 1)
</HARD-RULE>

After the scorer returns, parse its output in the main session:
1. Extract `SCORE: X/100`
2. Extract per-dimension scores from `DIMENSIONS:` section
3. Extract attack points from `ATTACKS:` section
4. Record score in iteration tracker

## Step 3: Decision Gate (Main Session)

<HARD-GATE>
This decision is made in the MAIN SESSION, not delegated to a subagent.
</HARD-GATE>

| Condition                                  | Action                          |
| ------------------------------------------ | ------------------------------- |
| Score >= target                            | Skip to Step 6 (final report)   |
| Score < target AND iterations remaining    | Proceed to Step 4               |
| Score < target AND no iterations remaining | Skip to Step 6 (report failure) |

Report current status to user:

```
Iteration {{N}}/{{MAX}}: scored {{SCORE}}/100 (target: {{TARGET}}). Revision subagent starting...
```

## Step 4: Invoke Reviser Subagent

Spawn `doc-reviser` agent via **Agent tool** (subagent_type: `zcode:doc-reviser` if registered, otherwise `general-purpose`).

<HARD-RULE>
Pass these inputs to the reviser:
- `DOC_PATHS` = same comma-separated paths as scorer
- `RUBRIC_PATH` = `plugins/zcode/skills/eval-prd/templates/rubric.md`
- `EVAL_REPORT_PATH` = `docs/features/<slug>/prd-eval-iteration-{{N}}.md`
- `ATTACK_POINTS` = the 3 attack points extracted from scorer output
</HARD-RULE>

The reviser will overwrite the PRD file(s) in place.

## Step 5: Loop

Increment iteration counter. Return to Step 2.

<EXTREMELY-IMPORTANT>
The scorer must NEVER be told what changes the reviser made. Only `PREVIOUS_REPORT_PATH` carries forward.
</EXTREMELY-IMPORTANT>

## Step 6: Final Report (Main Session)

```
## Eval-PRD Complete

**Final Score**: {{SCORE}}/100 (target: {{TARGET}})
**Iterations Used**: {{N}}/{{MAX}}

### Score Progression
| Iteration | Score | Delta |
|-----------|-------|-------|
| 1 | {{s1}} | - |
| 2 | {{s2}} | +{{d2}} |

### Dimension Breakdown (final)
| Dimension | Score | Max |
|-----------|-------|-----|
| Background & Goals | {{d1}} | 20 |
| Flow Diagrams | {{d2}} | 20 |
| Functional Specs | {{d3}} | 20 |
| User Stories | {{d4}} | 20 |
| Scope Clarity | {{d5}} | 20 |

### Outcome
{{"Target reached" / "Target NOT reached — N iterations exhausted"}}
{{If not reached: "Largest gaps: [dimension names]. Consider manual revision or increasing iterations."}}
```

Save the final report to `docs/features/<slug>/prd-eval.md`.

## Report Path Convention

| File               | Path                                                 |
| ------------------ | ---------------------------------------------------- |
| Iteration N report | `docs/features/<slug>/prd-eval-iteration-{{N}}.md`  |
| Final report       | `docs/features/<slug>/prd-eval.md`                  |

## Related

- `/write-prd` — Create or revise the PRD
- `/design-tech` — Next step after PRD passes evaluation
- `/ui-design` — Next step (optional) if prd-ui-functions.md exists
- `/breakdown-tasks` — After design docs are finalized
```

- [ ] **Step 2: Commit**

```bash
git add plugins/zcode/skills/eval-prd/SKILL.md
git commit -m "feat(eval-prd): add adversarial iteration loop with doc-scorer/doc-reviser"
```

---

## Task 10: Delete old agents and commit plan

**Files:**
- Delete: `plugins/zcode/agents/proposal-scorer.md`
- Delete: `plugins/zcode/agents/proposal-reviser.md`

- [ ] **Step 1: Delete the old agent files**

```bash
rm plugins/zcode/agents/proposal-scorer.md
rm plugins/zcode/agents/proposal-reviser.md
```

- [ ] **Step 2: Commit**

```bash
git add -A plugins/zcode/agents/
git commit -m "feat(agents): remove proposal-scorer and proposal-reviser, replaced by doc-scorer/doc-reviser"
```

- [ ] **Step 3: Commit the plan file**

```bash
git add docs/superpowers/plans/2026-04-22-eval-skills-adversarial-loop.md
git commit -m "docs(plans): complete implementation plan for eval skills adversarial loop"
```

---

## Self-Review

**Spec coverage:**
- ✅ doc-scorer agent (Task 1)
- ✅ doc-reviser agent (Task 2)
- ✅ proposal rubric extracted (Task 3)
- ✅ design rubric numeric 100pt (Task 4)
- ✅ PRD rubric numeric 100pt (Task 5)
- ✅ report templates updated (Task 6)
- ✅ eval-proposal SKILL.md updated (Task 7)
- ✅ eval-design SKILL.md rewritten (Task 8)
- ✅ eval-prd SKILL.md rewritten (Task 9)
- ✅ old agents deleted (Task 10)

**Input name consistency:**
- doc-scorer inputs: `DOC_PATHS`, `RUBRIC_PATH`, `REPORT_PATH`, `ITERATION`, `PREVIOUS_REPORT_PATH`
- doc-reviser inputs: `DOC_PATHS`, `RUBRIC_PATH`, `EVAL_REPORT_PATH`, `ATTACK_POINTS`
- All three SKILL.md files pass exactly these names ✅

**Rubric → report template dimension alignment:**
- eval-proposal rubric dimensions match report scorecard rows ✅
- eval-design rubric dimensions (Architecture 20, Interface 20, Error 15, Testing 15, Breakdown 20, Security 10) match report scorecard rows ✅
- eval-prd rubric dimensions (Background 20, Flow 20, Functional 20, Stories 20, Scope 20) match report scorecard rows ✅

