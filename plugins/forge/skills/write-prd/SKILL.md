---
name: write-prd
description: Use when user provides requirements or feature requests that need to be formalized into a structured PRD document through collaborative dialogue.
argument-hint: "[feature description or requirements]"
effort: high
---

# Write PRD

## Overview

From vague requirements to clear PRD (Product Requirements Document) through collaborative dialogue.

**Core principle**: Clarify "what to build" and "why" before coding, avoiding directional mistakes.

## Prerequisites

No required artifacts. If a brainstorm proposal exists, use it as optional input:

```bash
ls docs/proposals/<slug>/proposal.md 2>/dev/null  # optional, not blocking
```

<HARD-GATE>
Do NOT write any code, scaffold any project, or take any implementation action until the PRD is finalized and approved. Present the PRD and get user approval first.
</HARD-GATE>

<HARD-RULE>
**No technology selection allowed; constraints are allowed**:

- **Allowed**: Describe non-functional constraints — performance requirements (response time, concurrency), platform requirements (browser, mobile), compatibility, security/compliance. These are business-level requirements.
- **Forbidden**: Mention specific tech stacks — framework names, programming languages, databases, libraries, middleware, architectural patterns (e.g., microservices, event-driven). These are technology selections, left to the `/tech-design` phase.

**Judgment rule**: If the description is about "what effect to achieve" → allowed; if it's about "what tool to implement with" → forbidden.
</HARD-RULE>

## When to Use

**Trigger conditions:**

- User describes a feature/requirement without clear specifications
- User says "I want to..." or "We need..." without details
- Starting a new phase or major feature

**Skip when:**

- Clear task definitions already exist
- Simple bug fix or small tweak

## Process Flow

```
Explore context → Check proposal → Assess scope → Ask questions → Propose approaches → Present PRD sections → Write PRD Spec + User Stories + UI Functions → Create Manifest → Commit
```

## Checklist

1. **Explore project context** — check files, docs, recent commits
2. **Check for existing proposal** — read `docs/proposals/<slug>/proposal.md` if it exists
3. **Assess scope** — determine if request needs decomposition
4. **Ask clarifying questions** — one at a time via AskUserQuestion tool
5. **Propose 2-3 approaches** — with trade-offs and your recommendation
6. **Present PRD sections** — get approval after each section
7. **Write PRD Spec** — save to `docs/features/<slug>/prd/prd-spec.md`
8. **Write User Stories** — save to `docs/features/<slug>/prd/prd-user-stories.md`
9. **Write UI Functions** (mandatory for UI features) — save to `docs/features/<slug>/prd/prd-ui-functions.md`
10. **Create Manifest** — save to `docs/features/<slug>/manifest.md`
11. **Commit** — commit all documents

## Output Documents

| File | Template | Description |
|------|----------|-------------|
| `prd/prd-spec.md` | `templates/prd-spec.md` | Product requirements document with background, goals, scope, flows, functional specs |
| `prd/prd-user-stories.md` | `templates/prd-user-stories.md` | User stories derived from user roles identified in the PRD background |
| `prd/prd-ui-functions.md` | `templates/prd-ui-functions.md` | UI function highlights (requirements level, **mandatory** for features with UI surface) |
| `manifest.md` | `templates/manifest.md` | Feature index and traceability mapping |

## Step 1: Explore Project Context

Before asking questions, understand the current state:

- Check `docs/proposals/<slug>/proposal.md` if a proposal exists — carry forward business context
- Check `docs/features/<slug>/tasks/index.json` for related tasks
- Review recent git commits for related work
- **Read `docs/sitemap/sitemap.json`** if it exists — this is a business asset catalog listing all existing pages and their elements. Use it to understand what pages already exist when the user's requirements involve modifying existing pages.

**Forbidden**: Do not read `ARCHITECTURE.md`, `DECISIONS.md` or other technical docs to guide requirements discussion. Technical constraints do not belong in the PRD. `sitemap.json` is explicitly allowed — it documents business-level page inventory, not technical architecture decisions.

## Step 2: Assess Scope

Evaluate if the request is appropriately scoped:

- If request describes multiple independent subsystems → **Decompose first**
- If single focused feature → **Proceed with questions**

## Step 3: Ask Clarifying Questions

**CRITICAL**: Use `AskUserQuestion` tool for ALL questions.

### Question Guidelines

- **One question at a time** — never batch questions
- **Prefer multiple choice** — easier to answer than open-ended
- **Focus on understanding**: user roles, purpose, constraints, success criteria
- **Go back when needed** — if something doesn't make sense, clarify

See `examples/ask-questions.md` for concrete examples.

## Step 4: Propose Approaches

After understanding requirements, propose 2-3 **business approaches** (not technical implementations):

1. **Present options conversationally** with your recommendation
2. **Lead with your recommended option** and explain why
3. **Include trade-offs** for each approach (business impact, user experience, scope)

**Forbidden**: Approaches must not involve specific technology selection. Describe different business feature combinations or user flows, not technical implementation paths. Non-functional constraints (e.g., performance, platform, security/compliance) are allowed.

See `examples/propose-approaches.md` for structure and tips.

## Step 5: Present PRD Sections

Present incrementally, getting approval after each section:

| Section | Content | Key Points |
|---------|---------|------------|
| Background | Reason, target users, stakeholders | Must include all three dimensions |
| Goals | Objectives + quantified metrics | Quantify benefits wherever possible |
| Scope | In Scope / Out of Scope | Clear boundaries |
| Flow Description | Business flow + Mermaid diagram | Flowchart is required |
| Functional Specs | Reference to prd-ui-functions.md + Related changes | UI specs in separate file |
| Other Notes | Performance / Data / Monitoring / Security | Non-functional requirements |
| User Stories | As a / I want / So that + AC | Output to separate file |

## Step 6: Write PRD Spec

Use `templates/prd-spec.md` template.

**db-schema determination**: Assess whether the feature involves creating or modifying database tables. If yes → `db-schema: "yes"`. If no → `db-schema: "no"`. When unsure, default to `"yes"` — unnecessary DB design costs less than missing DB design.

**Directory structure:**

```
docs/features/<slug>/
├── manifest.md                # Feature index & traceability
├── prd/
│   ├── prd-spec.md            # PRD Spec
│   ├── prd-user-stories.md    # User Stories
│   └── prd-ui-functions.md    # UI Functions (mandatory for UI features)
├── design/                    # (created by /tech-design)
├── ui/                        # (created by /ui-design)
└── tasks/                     # (created by /breakdown-tasks)
    └── records/
```

## Step 7: Write User Stories

Derive user stories from user roles identified in the PRD background. Output to `prd/prd-user-stories.md`.

```
As a [user role from Background]
I want to [specific action]
So that [concrete benefit/goal]
```

**Coverage rules:**
- Every user type from Background must have at least one story
- Actions must be concrete — not "manage" or "handle" but "create X", "filter by Y"

**Acceptance Criteria** (Given/When/Then) must follow each story. Each AC must be objectively verifiable.

See `examples/user-stories.md` for concrete examples derived from Background roles.

## Step 8: Write UI Functions (mandatory for UI features)

For features with UI surfaces, create `prd/prd-ui-functions.md` using `templates/prd-ui-functions.md`.
This step is **mandatory** when the feature has any UI surface. Skip only for backend/API/CLI-only features with no UI surface.

**Placement rules** (mandatory for every UI Function):

1. Read `docs/sitemap/sitemap.json` to understand existing page inventory
2. For each UI Function, declare its Placement:
   - `new-page` — this function creates a brand new page
   - `existing-page:<route>` — this function adds UI to an existing page (route from sitemap)
3. For `existing-page`, also specify the Position within the page (e.g., "above sub-items table")
4. After all UI Functions are defined, compile the Page Composition summary table at the end of the document

**Validation**: Every UI Function MUST have a Placement section. Missing Placement → error, do not proceed.

**Navigation Architecture platform handling**:

The template contains two conditional navigation sections. Render exactly one based on platform:

- **platform=tui**: Render the **TUI Navigation** section (Keymap, Panel Layout, Modes, Navigation Rules). Omit the Pointer-Driven Navigation section entirely.
- **platform=web|mobile|mini-program|tablet**: Render the **Pointer-Driven Navigation** section (Primary Navigation, Secondary Pages, Navigation Rules). Omit the TUI Navigation section entirely.

TUI uses keyboard-driven navigation (keymap + panels + modes) instead of pointer-driven navigation (pages + routes + icons). Do not mix the two patterns.

**Downstream impact**: Placement determines how downstream skills structure their output:
- **`/breakdown-tasks`**: `existing-page` generates a Build task + an Integrate task (wire component into existing page); `new-page` generates a Build task + a Page Assembly task (create page file, register route, compose components)
- **`/gen-test-cases`**: auto-generates integration verification test cases for `existing-page` placements, ensuring the component is visible at the correct position in the target page

## Step 9: Create Manifest

Create `manifest.md` at the feature root using `templates/manifest.md`:
- Fill in PRD entries and summaries
- Set status to `prd`
- Include UI Functions row only if `prd/prd-ui-functions.md` was created

## Step 9.5: Self-Check

Before presenting to the user, verify the PRD passes these checks:

| Check | What to verify |
|-------|----------------|
| Background completeness | Reason + target users + stakeholders all present and specific |
| Goals quantified | At least one numeric target (% , count, time) |
| Flow diagram | Mermaid flowchart with decision points (diamond nodes) and at least one error/exception branch |
| Functional specs | prd-spec.md references prd-ui-functions.md; prd-ui-functions.md tables filled — no placeholder rows |
| User stories | One story per user role, each with Given/When/Then AC |
| Scope consistency | In-scope items match what's described in Functional Specs and user stories |
| No vague language | No "better", "faster", "improved" without quantification |
| Placement completeness | Every UI Function has a Placement section with Mode and target |
| Placement consistency | existing-page routes exist in sitemap.json (if sitemap available) |
| Sitemap availability | If sitemap.json not found, warn: "Sitemap unavailable — existing-page routes cannot be validated. Run /gen-sitemap." |
| Page Composition valid | Page Composition table lists all pages with correct UI Function references |
| db-schema filled | db-schema frontmatter is "yes" or "no" (not empty) |

## Step 10: Review & Commit

<HARD-RULE>
Do NOT commit documents automatically. Present all generated documents to the user for review and wait for explicit approval before committing.
</HARD-RULE>

1. Present the full PRD spec, user stories, and UI functions (if any) to the user
2. Wait for the user to review and approve (or request changes)
3. Only commit after explicit user approval:

```bash
git add docs/features/<slug>/
git commit -m "docs: add PRD for <slug>"
```

## Step 11: Adversarial Eval Prompt

After committing, use `AskUserQuestion` to ask:

> Run `/eval-prd` for adversarial evaluation? (default: 900 points / 3 rounds)

- **Yes** → invoke `/eval-prd` via `Skill` tool
- **Custom** → invoke `/eval-prd --target X --iterations Y` via `Skill` tool
- **No** → proceed to `/ui-design` (if PRD has UI functions) or `/tech-design`

## Step 12: Knowledge Review

After Step 11 (Adversarial Eval Prompt) completes, run knowledge auto-extraction from the PRD:

1. Read `${CLAUDE_SKILL_DIR}/../../references/shared/knowledge-extraction.md` and execute its extraction flow with:
   - `trigger`: `write-prd`
   - `artifacts`: PRD content (`docs/features/<slug>/prd/prd-spec.md`, `docs/features/<slug>/prd/prd-user-stories.md`)

2. The extraction flow handles:
   - Scanning PRD content for notable knowledge (decisions, lessons, conventions, business rules)
   - Focusing on business rules that apply across features, not feature-specific logic
   - Silent exit when no notable knowledge detected — no output, no prompts
   - Presenting extracted knowledge for user confirmation via AskUserQuestion
   - Writing confirmed knowledge to appropriate directories using shared formats

## Integration

Works well with skills:

- `/eval-prd` - Evaluate PRD quality before proceeding to design phase
- `/tech-design` - After PRD passes evaluation, produce technical design document
- `/ui-design` - After PRD passes evaluation, produce UI design spec (if prd-ui-functions.md exists)
- `docs/decisions/` - Record key decisions during PRD creation (category-based directory)
