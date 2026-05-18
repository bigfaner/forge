---
name: brainstorm
description: Use when a user has a vague idea or feature request and needs to explore it before formalizing into a PRD. Outputs a structured proposal document.
argument-hint: "[idea or feature description]"
---

# Brainstorm

From vague idea to structured proposal, through relentless collaborative dialogue.

**Core principle**: Relentlessly interview every aspect of the idea until reaching shared understanding. Before investing in a PRD, confirm the problem is worth solving.

<HARD-GATE>
Do NOT write any code or take implementation action. This skill produces a proposal document only.
</HARD-GATE>

<HARD-RULE>
**No technology selection; constraints only.** Describe "what effect to achieve" (performance, platform, security) — not "what tool to implement with" (frameworks, languages, databases, patterns). Technology selection belongs in `/tech-design`.
</HARD-RULE>

## Process Flow

```
Analyze context → Walk the design tree → Propose approaches → Define scope → Write proposal → Commit
```

## Step 1: Analyze Context

Before asking any question, search the codebase for related features, docs, proposals, PRDs, and recent commits. Synthesize findings internally — do NOT show this analysis to the user. It drives informed questioning.

## Step 2: Walk the Design Tree

Interview the user relentlessly about every aspect of the idea until reaching shared understanding. Walk down each branch of the design tree, resolving dependencies between decisions one-by-one. For each question, provide your recommended answer.

**CRITICAL**: Use `AskUserQuestion` tool for ALL questions. One question at a time.

### Core Principles

- **Shared understanding is the termination condition** — keep going until both sides genuinely agree, not just surface agreement
- **Walk the design tree** — resolve parent decisions before sub-branches. If B depends on A, resolve A first
- **Recommend before asking** — provide your recommended answer as the first option
- **Codebase-first answering** — if a question could be answered by exploring the codebase, explore the codebase instead

### Decision Clusters

Three clusters provide direction — traverse freely based on dependencies, not fixed order:

| Cluster       | Drives questions about                                  |
| ------------- | ------------------------------------------------------- |
| **Problem**   | Core problem, affected users, urgency, cost of inaction |
| **Solution**  | Success criteria, must-haves, user workflows            |
| **Challenge** | Simpler alternatives, risks, blind spots                |

Backtrack when a branch reveals an earlier assumption was wrong. Derive questions from findings, not templates — reference concrete facts.

## Step 3: Propose Approaches

Propose 2-3 **business approaches** (not technical implementations). Lead with your recommendation, include trade-offs, and always include "do nothing" as one alternative. Let the user decide.

## Step 4: Define Scope

Propose in-scope and out-of-scope boundaries. Get explicit user agreement. If too large, suggest decomposing.

## Step 5: Write Proposal

Save to `docs/proposals/<slug>/proposal.md` using `templates/proposal.md`.

### Quality Standards

| Section | Standard | Red Flag |
|---------|----------|----------|
| Problem | Specific statement + evidence + urgency | "We need to improve X" |
| Solution | Concrete user-facing behavior | "Build a system that..." |
| Alternatives | Honest trade-offs including "do nothing" | Straw-man alternatives with only pros |
| Scope | Deliverable-level items, bounded | Vague areas, open-ended |
| Risks | 3+ specific risks with actionable mitigations | "We'll handle it" |
| Success Criteria | Measurable, testable, covers all scope | "Works well" or "Better UX" |

<HARD-RULE>
Do NOT commit automatically. Present to the user and wait for explicit approval.
</HARD-RULE>

## Step 6: Commit

```bash
git add docs/proposals/<slug>/
git commit -m "docs: add proposal for <feature-slug>"
```

## Step 7: Adversarial Eval Prompt

After committing, ask via `AskUserQuestion`:

> Run `/eval-proposal` for adversarial evaluation? (default: 900 points / 3 rounds)

- **Yes** → invoke `/eval-proposal` via `Skill` tool
- **Custom** → invoke `/eval-proposal --target X --iterations Y` via `Skill` tool
- **No** → proceed to `/write-prd`
