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

Two clusters provide direction — traverse based on dependencies, not fixed order. Each cluster has **mandatory** challenge tools embedded within it:

| Cluster       | Drives questions about                                  | Embedded Challenge Tools         |
| ------------- | ------------------------------------------------------- | -------------------------------- |
| **Problem**   | Core problem, affected users, urgency, cost of inaction | **5 Whys** + **XY Problem Detection** |
| **Solution**  | Success criteria, must-haves, user workflows            | **Assumption Flip** + **Stress Test** |

Backtrack when a branch reveals an earlier assumption was wrong. Derive questions from findings, not templates — reference concrete facts.

### Challenge Protocol

Challenge is not a separate step — it is a **mandatory behavior** embedded in every decision point. Each challenge must be grounded in facts (see Fact-Driven Principle below). Empty or vague questioning is forbidden.

#### Challenge Tools

| Tool | When to Use | Trigger Condition | Termination Condition |
|------|-------------|-------------------|----------------------|
| **5 Whys** | Problem Cluster — drill into root cause | User states a surface-level symptom or vague pain point | Root cause identified (causal chain reaches a fundamental constraint), OR 3 consecutive "why" answers are consistent |
| **XY Problem Detection** | Problem Cluster — detect when user's stated need may not be the real need | User asks for a specific solution rather than describing a problem | User confirms the actual underlying problem, OR user provides clear rationale for why the specific solution is needed |
| **Assumption Flip** | Solution Cluster — validate critical assumptions | User presents a solution that depends on an unverified assumption | Assumption is confirmed by evidence, OR assumption is overturned and solution is adjusted, OR user provides sufficient domain expertise as evidence |
| **Stress Test** | Solution Cluster — expose hidden risks in seemingly perfect solutions | User seems satisfied with a solution without considering edge cases | All identified edge cases are addressed, OR user explicitly accepts residual risk with documented rationale |
| **Occam's Razor** | Both Clusters — meta-principle applied at all times | Multiple competing explanations or solutions coexist | Simplest viable option selected, OR complexity is justified by concrete evidence |

#### Fact-Driven Principle

Every challenge must cite one of three evidence types:

1. **Codebase facts** — existing implementations, architecture patterns, API contracts found in the codebase
2. **Logical consistency** — internal contradictions, circular reasoning, or gaps in the user's own argument
3. **Domain common sense** — widely accepted knowledge in the relevant technical or business domain

For greenfield projects (no existing codebase), rely on logical consistency and domain common sense. The absence of code does not excuse challenges from providing evidence.

#### Challenge Tone

Challenges must be **rationally prudent**, not hostile. Every challenge follows this structure:

1. **State the observation** — what was said or assumed
2. **Present the evidence** — cite a specific fact from the three evidence types
3. **Pose the question** — ask what the implication is

Example: "You mentioned caching as the solution for latency. The current p99 latency is 50ms (codebase fact: `src/api/middleware/timer.ts`). At this level, network round-trip dominates — caching may not address the actual bottleneck. What does the latency breakdown show?"

#### Occam's Razor (Meta-Principle)

Throughout the entire brainstorm, when multiple explanations or approaches coexist, prefer the simplest one that satisfies all known constraints. This is not a separate tool to invoke — it is a standing principle that applies to every decision:

- If a simple explanation covers the observed problem, do not propose complex alternatives without evidence requiring complexity
- If a straightforward approach meets all success criteria, do not add sophistication without justification
- When the user proposes a complex solution, ask: "Is there a simpler way to achieve the same outcome?"

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
