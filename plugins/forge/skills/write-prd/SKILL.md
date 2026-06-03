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

### Intent Detection

Read the `intent` field from `docs/proposals/<slug>/proposal.md` frontmatter. This determines the PRD format via the Pipeline Configuration table below.

#### Pipeline Configuration

| Intent | PRD Format | User Stories | API Handbook | Test Pipeline | Security Review |
|--------|-----------|-------------|-------------|--------------|----------------|
| `new-feature` | Full | Yes | Yes | Yes | If signal |
| `enhancement` | Simplified (Background + Goals + Test Pipeline, skip User Stories) | No | If signal | Yes | If signal |
| `refactor` | Spec-only | No | If signal | Yes | If signal |
| `cleanup` | Spec-only | No | No | Yes | No |
| `fix` | Spec-only | No | If signal | Yes (reproduce → fix → verify) | No |
| `doc` | Minimal (title + goals + scope only) | No | No | No | No |

**Default**: If `intent` is missing or empty, treat as `new-feature` — full PRD pipeline unchanged.

**Detection**:
```bash
# Read intent from proposal frontmatter
head -20 docs/proposals/<slug>/proposal.md | grep "^intent:"
```

#### Override Signals

Pipeline defaults are determined by intent, but PRD content signals can enable additional pipeline steps. During content generation, detect the following signals within the same LLM call:

| Signal Type | Keywords / Patterns | Override Action |
|-------------|-------------------|-----------------|
| API 变更 | "API", "endpoint", "命令重命名", "接口变更", "breaking change" | Enable API Handbook |
| 用户可见行为 | "用户可见", "UI 变更", "CLI 输出", "新选项" | Enable User Stories |
| 安全相关 | "认证", "授权", "权限", "加密", "token" | Enable Security Review |
| 性能相关 | "性能", "延迟", "吞吐量", "缓存" | Enable Performance Baseline |
| 数据迁移 | "迁移", "schema 变更", "数据格式" | Enable Migration Plan |

**Detection rules**:
- Content generation and signal matching happen in parallel inference (same LLM call), not sequentially
- Negation handling: skip signals in negative context (e.g., "不涉及 API 变更"). Relies on LLM context understanding, not keyword matching
- Override only adds steps (开启), never removes. Worst case: unnecessary artifact generated, caught in user review
- Multiple signals trigger independently and stack (e.g., both "API" and "性能" → enable both API Handbook and Performance Baseline)
- When an override triggers, generate a comment in the PRD output documenting it (e.g., `<!-- Override: API handbook enabled by signal "接口变更" -->`)
- For `doc` intent: Minimal PRD format has no pipeline steps that can be overridden — override signals become no-op by design

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
Explore context → Check proposal + detect intent → Assess scope → Ask questions → [Intent-specific steps] → Create Manifest → Self-Check → Review & Commit → Adversarial Eval
```

**Intent-specific steps**:

| Intent | Questions Focus | PRD Content | Override Signals |
|--------|----------------|-------------|------------------|
| `new-feature` | User roles, purpose, constraints | Full PRD + User Stories + UI Functions | Detect during content generation |
| `enhancement` | Enhancement target, improvement goals | Background + Goals + Test Pipeline | Detect during content generation |
| `refactor` / `cleanup` / `fix` | Change scope, constraints, verification | Spec-only (3 mandatory fields) | Detect before writing PRD |
| `doc` | Documentation scope | Minimal (title + goals + scope) | No-op (no overridable steps) |

## Checklist

1. **Explore project context** — check files, docs, recent commits
2. **Check proposal + detect intent** — read `docs/proposals/<slug>/proposal.md`, extract `intent` from frontmatter (skip intent detection if `new-feature` or missing)
3. **Assess scope** — determine if request needs decomposition
4. **Ask clarifying questions** — one at a time via AskUserQuestion tool; focus varies by intent (see Process Flow table)
5. **Propose approaches** (`new-feature` only) — 2-3 business approaches with trade-offs and recommendation
6. **Present PRD sections** (`new-feature` only) — get approval after each section
7. **Detect override signals** (all intents except `doc`) — scan content for signals (see Override Signals table); generate `<!-- Override: ... -->` comments for triggered signals
8. **Write PRD** — format determined by intent (see Output Documents); save to `docs/features/<slug>/prd/prd-spec.md`
9. **Write User Stories** (`new-feature` only, or override "用户可见行为") — save to `docs/features/<slug>/prd/prd-user-stories.md`
10. **Write UI Functions** (`new-feature` with UI surface only, or override "用户可见行为") — save to `docs/features/<slug>/prd/prd-ui-functions.md`
11. **Create Manifest** — save to `docs/features/<slug>/manifest.md`
12. **Self-Check** — verify PRD passes checks
13. **Review & Commit** — commit all documents after user approval
14. **Adversarial Eval** — run eval-prd if configured

## Output Documents

| File | Template | Generated For |
|------|----------|---------------|
| `prd/prd-spec.md` | `templates/prd-spec.md` | All intents (format varies, see PRD Format column in Pipeline Configuration) |
| `prd/prd-user-stories.md` | `templates/prd-user-stories.md` | `new-feature`; others only if override "用户可见行为" triggers |
| `prd/prd-ui-functions.md` | `templates/prd-ui-functions.md` | `new-feature` with UI surface; others only if override "用户可见行为" triggers |
| `manifest.md` | `templates/manifest.md` | All intents |

## Step 1: Explore Project Context

Before asking questions, understand the current state:

- Check `docs/proposals/<slug>/proposal.md` if a proposal exists — carry forward business context
- Check `docs/features/<slug>/tasks/index.json` for related tasks
- Review recent git commits for related work
- **Read `docs/sitemap/sitemap.json`** if the project has a `web` surface AND the file exists — check via `forge surfaces --json`; if the output contains a surface with type `web`, read the sitemap. This is a business asset catalog listing all existing pages and their elements. Use it to understand what pages already exist when the user's requirements involve modifying existing pages. **If the project has no `web` surface, skip sitemap reading entirely** — sitemap is a web-specific artifact generated by `/gen-web-sitemap`.

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

<EXTREMELY-IMPORTANT>
**Intent Gate**: If `intent` is `refactor`, `cleanup`, `fix`, `enhancement`, or `doc`, **skip this entire step**. Do NOT generate `prd/prd-user-stories.md`. Proceed directly to Step 9 (Create Manifest).

Exception: If override signal "用户可见行为" is triggered for `refactor`, `cleanup`, or `fix` intent, generate User Stories.

The "As a user / I want / So that" format is semantically empty when there is no new user-observable behavior to describe as a user story.
</EXTREMELY-IMPORTANT>

**Doc-only Gate**: If all In Scope items are non-compilable artifacts (`.md`, `.yaml`, `.json` under `docs/`, `skills/`, etc.), skip this step and note that user stories are not needed for doc-only features. User stories serve gen-journeys → test script generation, which requires testable code.

**How to detect**: Examine the In Scope section of the PRD spec. If every listed item targets non-compilable, non-runnable file paths, the feature is doc-only. If any item involves compilable or runnable files (`.go`, `.ts`, `.py`, `.java`, etc.), proceed with user story generation.

When proceeding, derive user stories from user roles identified in the PRD background. Output to `prd/prd-user-stories.md`.

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

## Step 7A: Write Spec-Only PRD (refactor / cleanup / fix intent)

<EXTREMELY-IMPORTANT>
This step applies **only** when `intent` is `refactor`, `cleanup`, or `fix`. If `intent` is `new-feature`, `enhancement`, or `doc`, skip this step entirely.
</EXTREMELY-IMPORTANT>

When `intent` is `refactor`, `cleanup`, or `fix`, the PRD spec must contain three mandatory fields that provide sufficient information for `/tech-design` without relying on user stories:

### Mandatory Fields

| Field | Description | Example |
|-------|-------------|---------|
| **变更范围 (Change Scope)** | Affected modules, files, or packages. List concrete paths or module names. | `pkg/constants/`, `internal/enum/types.go`, `cmd/convert.go` |
| **约束条件 (Constraints)** | Behavioral invariants that must be preserved during the refactoring. These are the "things that must not break". | "All existing API endpoints return identical responses", "CLI exit codes unchanged", "No new exported symbols" |
| **验证标准 (Verification Criteria)** | Regression acceptance criteria — how to verify the refactoring succeeded without introducing regressions. | "All existing tests pass", "No behavioral change in output of `forge task list`", `gofmt` and `golint` pass |

### Writing Guidelines

- **Focus questions**: In Step 4 (Ask Clarifying Questions), focus questions on:
  1. Which modules/files are in scope for this refactoring?
  2. What behavioral invariants must be preserved?
  3. How should we verify no regressions were introduced?
- **Scope section**: The In Scope / Out of Scope section should map directly to the change scope field
- **No user stories**: Do not generate `prd/prd-user-stories.md` — the three mandatory fields above replace user stories for refactoring
- **No UI functions**: Do not generate `prd/prd-ui-functions.md` — refactoring does not introduce new UI surfaces

### PRD Spec Template Adaptation

When using `templates/prd-spec.md` for a spec-only PRD:
- Replace the "User Stories" reference in the template with the three mandatory fields section
- Omit the Flow Description section (Mermaid diagram) unless the refactoring changes an external flow
- Include the three mandatory fields as a prominent section, e.g.:

```markdown
## Refactoring Specification

### Change Scope
<!-- List affected modules, files, packages -->

### Constraints (Behavioral Invariants)
<!-- List things that must not change -->

### Verification Criteria
<!-- List regression acceptance criteria -->
```

## Step 8: Write UI Functions (mandatory for UI features)

<EXTREMELY-IMPORTANT>
**Intent Gate**: If `intent` is `refactor`, `cleanup`, `fix`, `enhancement`, or `doc`, **skip this step**. These intents do not introduce new UI surfaces by default.
</EXTREMELY-IMPORTANT>

For features with UI surfaces, create `prd/prd-ui-functions.md` using `templates/prd-ui-functions.md`.
This step is **mandatory** when the feature has any UI surface. Skip only for backend/API/CLI-only features with no UI surface.

Placement rules, Navigation Architecture, and downstream impact rules — see `rules/ui-functions.md`.

## Step 9: Create Manifest

Create `manifest.md` at the feature root using `templates/manifest.md`:
- Fill in PRD entries and summaries
- Replace `{{DATE}}` with today's date in `YYYY-MM-DD` format
- Set status to `prd`
- Include User Stories row only if `prd/prd-user-stories.md` was generated (skip for `refactor`/`cleanup`/`fix`/`enhancement`/`doc` intent unless override signal triggers it)
- Include UI Functions row only if `prd/prd-ui-functions.md` was created

## Step 10: Self-Check

Before presenting to the user, verify the PRD passes all checks in `rules/self-check.md`.

## Step 11: Review & Commit

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

## Step 12: Adversarial Eval Prompt

<EXTREMELY-IMPORTANT>
Eval auto-run check — do NOT use AskUserQuestion when config enables auto-run.

Run the following config check sequence via Bash tool:

```bash
# Eval auto-run check (prd)
EVAL_ENABLED=$(forge config get auto.eval.prd 2>/dev/null)
if [ "$EVAL_ENABLED" = "true" ]; then
  echo "AUTO_RUN"
elif [ "$EVAL_ENABLED" = "false" ]; then
  echo "SKIP"
else
  echo "FALLBACK_ASK"
fi
```

Based on the output:
- **AUTO_RUN** → invoke `/eval-prd` via `Skill` tool (default: 900 points / 3 rounds)
- **SKIP** → skip eval, output "eval-prd 已通过配置跳过", proceed to `/ui-design` (if PRD has UI functions) or `/tech-design`
- **FALLBACK_ASK** → ask via `AskUserQuestion`: "Run `/eval-prd` for adversarial evaluation? (default: 900 points / 3 rounds)"
  - **Yes** → invoke `/eval-prd` via `Skill` tool
  - **Custom** → invoke `/eval-prd --target X --iterations Y` via `Skill` tool
  - **No** → proceed to `/ui-design` (if PRD has UI functions) or `/tech-design`
</EXTREMELY-IMPORTANT>

## Step 13: Knowledge Review

After Step 12 (Adversarial Eval Prompt) completes, run knowledge auto-extraction from the PRD.

Full extraction flow, knowledge type definitions, notable-knowledge heuristics, and deduplication rules — see `rules/knowledge-extraction.md`.

