---
name: brainstorm
description: Use when a user has a vague idea or feature request and needs to explore it before formalizing into a PRD. Outputs a structured proposal document.
---

# Brainstorm

## Overview

From vague idea to structured proposal, explore the problem space through collaborative dialogue.

**Core principle**: Relentlessly interview every aspect of the idea until reaching shared understanding. Before investing in a PRD, confirm the problem is worth solving and the approach is worth investing in.

## Prerequisites

No required artifacts. This is the entry point of the workflow.

<HARD-GATE>
Do NOT write any code or take implementation action. This skill produces a proposal document only.
</HARD-GATE>

<HARD-RULE>
**No technology selection allowed; constraints are allowed**:

- **Allowed**: Describe non-functional constraints — performance requirements (response time, concurrency), platform requirements (browser, mobile), compatibility, security/compliance. These are business-level requirements.
- **Forbidden**: Mention specific tech stacks — framework names, programming languages, databases, libraries, middleware, architectural patterns (e.g., microservices, event-driven). These are technology selections, left to the `/tech-design` phase.

**Judgment rule**: If the description is about "what effect to achieve" → allowed; if it's about "what tool to implement with" → forbidden.
</HARD-RULE>

## When to Use

**Trigger conditions:**

- User describes an idea without clear specifications
- User says "I'm thinking about..." or "What if we..."
- Starting exploration before committing to a feature

**Skip when:**

- Requirements are already clear (use `/write-prd` directly)
- Bug fix or small tweak

## Process Flow

```
Analyze context → Synthesize findings → Walk the design tree (resolve decisions branch-by-branch) → Propose approaches → Define scope → Write proposal → Commit
```

## Checklist

1. **Analyze project context** — extract keywords, grep code, check docs & proposals
2. **Synthesize findings** — summarize what you found; identify gaps, overlaps, and open questions
3. **Walk the design tree** — resolve decisions branch-by-branch via AskUserQuestion, derived from findings not templates
4. **Propose 2-3 approaches** — with trade-offs and your recommendation
5. **Define scope** — get user agreement on in-scope / out-of-scope
6. **Write proposal** — save to `docs/proposals/<slug>/proposal.md`
7. **Commit** — commit the proposal document

## Step 1: Analyze Project Context

Before asking any question, run concrete analysis to understand what already exists. The goal is to ask **informed, targeted** questions — not generic templates.

### 1.1 Extract Keywords

From the user's idea description, extract 3-5 search keywords.

### 1.2 Search Codebase

Run these analyses in parallel:

| Action                   | Tool                                  | Purpose                                          |
| ------------------------ | ------------------------------------- | ------------------------------------------------ |
| Find related features    | `Grep` keywords across codebase       | Is this already implemented?                     |
| Find related docs        | `Glob` `docs/**/*.md`                 | Are there existing decisions or proposals?       |
| Check existing proposals | `Glob` `docs/proposals/**/*.md`       | Has someone already proposed this?               |
| Check existing PRDs      | `Glob` `docs/features/**/prd-spec.md` | Is there a related PRD?                          |
| Review recent commits    | `Bash` `git log --oneline -20`        | Any related recent work?                         |

### 1.3 Synthesize Findings

After analysis, summarize internally:

```
Analysis Brief:
- Found:      [what already exists]
- Gap:        [what's missing that the user's idea would fill]
- Overlap:    [existing features/skills that overlap]
- Open Qs:    [specific things you couldn't determine from code alone]
```

This brief drives Step 2. **Do NOT show the brief to the user** — it's internal. Use it to generate questions that reference concrete facts.

## Step 2: Ask Targeted Questions

Interview the user relentlessly about every aspect of the idea until reaching shared understanding. Walk down each branch of the design tree, resolving dependencies between decisions one-by-one. For each question, provide your recommended answer.

**CRITICAL**: Use `AskUserQuestion` tool for ALL questions.

### Core Principles

- **Shared understanding is the termination condition** — keep going until both sides genuinely agree on each decision, not just surface agreement
- **Walk the design tree** — each question is a branch; resolve parent decisions before exploring sub-branches. If question B depends on question A, resolve A first
- **Recommend before asking** — for each question, provide your recommended answer based on analysis. Present it as the first option so the user can confirm or adjust
- **One question at a time** — never batch
- **Codebase-first answering** — if a question could be answered by exploring the codebase, explore the codebase instead. Never ask what you can discover yourself

### Questioning Techniques

- **Prefer multiple choice** — easier to answer than open-ended
- **Derive questions from findings** — not from templates
- **Reference concrete facts** — "I found X already does Y..." not "Is there something similar?"
- **Skip answered questions** — if the user already stated something, don't re-ask
- **Dig deeper on vagueness** — follow up when answers are unclear

### How to Derive Questions from Findings

| Finding                               | Derive This Question                                                                                                                           |
| ------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| Existing feature does 80% of the idea | "I noticed [X] already handles [Y]. Is your idea extending that, or solving a different aspect?"                                               |
| No existing feature or doc            | "I couldn't find anything related to [topic] in the codebase. Is this a greenfield idea, or does it connect to something I might have missed?" |
| Multiple overlapping features         | "I found [A] and [B] both touch on [area]. Does your idea replace one, integrate both, or work alongside them?"                                |
| Recent commits on related area        | "I see recent work on [X] (commit abc123). Is your idea building on that, or independent?"                                                     |
| Existing proposal or PRD              | "There's already a proposal for [X] at [path]. Is your idea the same direction, a revision, or something different?"                           |

See `examples/context-aware-questions.md` for full worked examples.

### Design Tree

Questioning is a decision tree traversal, not a linear phase sequence. Three decision clusters provide direction for what to ask:

| Cluster      | Drives questions about                                       |
| ------------ | ------------------------------------------------------------ |
| **Problem**  | Core problem, affected users, urgency, cost of inaction      |
| **Solution** | Success criteria, must-haves, user workflows                 |
| **Challenge**| Simpler alternatives, risks, blind spots                     |

Traverse freely based on decision dependencies. Backtrack when a branch reveals an earlier assumption was wrong — a Challenge answer may invalidate a Problem assumption. Walk each branch to shared understanding before moving to the next.

### Fallback to Generic Questions

If analysis yields no useful findings (e.g., greenfield project, empty repo), fall back to the generic question templates in `examples/ask-questions.md`.

## Step 3: Propose Approaches

After understanding the problem and solution direction, propose 2-3 **business approaches** (not technical implementations):

1. **Present options conversationally** with your recommendation
2. **Lead with your recommended option** and explain why
3. **Include trade-offs** for each approach (business impact, user experience, scope)
4. **Always include "do nothing"** (status quo) as one alternative — it forces honest assessment of whether the problem warrants action
5. **Let user make the final decision**

**Forbidden**: Approaches must not involve specific technology selection. Describe "what to build" (features, flows, user experience) not "how to build" (technical implementation). Non-functional constraints (e.g., "response time < 200ms", "must support offline") are allowed.

```
Based on our discussion, I see 3 approaches:

**Approach A (Recommended):** [Business-level description of what to build]
- Business value: ...
- Trade-offs: ...

**Approach B:** [Business-level description]
- Business value: ...
- Trade-offs: ...
```

## Step 4: Define Scope

Propose in-scope and out-of-scope boundaries. Get explicit user agreement.

If scope is too large, suggest decomposing into multiple proposals.

## Step 5: Write Proposal

Save to `docs/proposals/<slug>/proposal.md` using `templates/proposal.md`.

### Quality Standards

Before presenting, verify each section meets these standards:

| Section | Standard | Red Flag |
|---------|----------|----------|
| Problem | Specific statement + evidence + urgency | "We need to improve X" |
| Solution | Concrete user-facing behavior | "Build a system that..." |
| Alternatives | Honest trade-offs including "do nothing" | Straw-man alternatives with only pros |
| Scope | Deliverable-level items, bounded | Vague areas, open-ended |
| Risks | 3+ specific risks with actionable mitigations | "We'll handle it" |
| Success Criteria | Measurable, testable, covers all scope | "Works well" or "Better UX" |

<HARD-RULE>
Do NOT commit the proposal automatically. Present the document to the user for review and wait for explicit approval before committing.
</HARD-RULE>

## Step 6: Review & Commit

1. Present the full proposal content to the user
2. Wait for the user to review and approve (or request changes)
3. Only commit after explicit user approval:

```bash
git add docs/proposals/<slug>/
git commit -m "docs: add proposal for <feature-slug>"
```

## Step 7: Adversarial Eval Prompt

After committing, use `AskUserQuestion` to ask:

> Run `/eval-proposal` for adversarial evaluation? (default: 80 points / 3 rounds)

- **Yes** → invoke `/eval-proposal` via `Skill` tool
- **Custom** → invoke `/eval-proposal --target X --iterations Y` via `Skill` tool
- **No** → proceed to `/write-prd`

## Integration

Works well with:

- `/eval-proposal` — Adversarial evaluation loop after proposal is created
- `/write-prd` — Takes proposal as optional input to formalize into PRD
