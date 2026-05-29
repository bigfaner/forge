---
name: brainstorm
description: Use when a user has a vague idea or feature request and needs to explore it before formalizing into a PRD. Outputs a structured proposal document.
argument-hint: "[idea or feature description]"
---

# Brainstorm

From vague idea to structured proposal, through relentless collaborative dialogue.

**Core principle**: Help the user clarify what they truly need through structured dialogue. The goal is shared understanding — not filtering ideas, but making implicit assumptions explicit so the right path becomes obvious. Pseudo-requirements die naturally when the thinking is clear.

<HARD-GATE>
Do NOT write any code or take implementation action. This skill produces a proposal document only.
</HARD-GATE>

<HARD-RULE>
**No technology selection; constraints only.** Describe "what effect to achieve" (performance, platform, security) — not "what tool to implement with" (frameworks, languages, databases, patterns). Technology selection belongs in `/tech-design`.
</HARD-RULE>

## Process Flow

```
Analyze context → Walk the design tree → Propose approaches → Define scope → Infer intent → Write proposal → Commit → Adversarial eval
```

## Step 1: Analyze Context

Before asking any question, search the codebase for related features, docs, proposals, PRDs, and recent commits. Synthesize findings internally — do NOT show this analysis to the user. It drives informed questioning.

## Step 2: Walk the Design Tree

Interview the user about every aspect of the idea until reaching shared understanding. Walk down each branch of the design tree, resolving dependencies between decisions one-by-one. For each question, provide your recommended answer.

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

### Need Gate (embedded at every decision point)

As the design tree reveals specific features or capabilities, apply this clarification check **before diving into how to implement each one**. The gate triggers when a concrete feature crystallizes from the discussion — not on the initial vague idea, but on each actionable feature point that emerges.

**Three checks (in order)**:

1. **Simpler alternative?** — Search for existing tools, commands, pipe compositions, or ecosystem utilities that already solve this. If found → propose the simpler path.
2. **Is this the real need?** — Hypothesize the user's underlying goal (Y). Confirm with the user before challenging: "I understand your core need is [Y] — is that right?" Only after confirmation, assess whether X is the best path to Y.
3. **Why now?** — "What happens if we defer this?" If the cost of delay is low → suggest deferring until a real need emerges.

**If the user overrides the gate** (insists on proceeding despite a simpler alternative or deferral suggestion): accept and record `Challenge Override: user chose to proceed. Reason: <user's reason or "not stated">.` in the proposal. Do NOT keep challenging after an override.

For the full protocol including search strategy, Y-confirmation flow, and timing criteria, see `rules/challenge-protocol.md`.

### Challenge Protocol

Challenge is mandatory at every decision point, not a separate step. For the full challenge toolkit, evidence requirements, and tone guidelines, see `rules/challenge-protocol.md`.

**Key rules**:
- Every challenge must cite one of: codebase facts, logical consistency, or domain common sense
- Challenges must be rationally prudent, not hostile: state observation → present evidence → pose question

## Step 3: Propose Approaches

Propose 2-3 **business approaches** (not technical implementations). Lead with your recommendation, include trade-offs, and always include "do nothing" as one alternative. Let the user decide.

## Step 4: Define Scope

Propose in-scope and out-of-scope boundaries. Get explicit user agreement. If too large, suggest decomposing.

## Step 4.5: Infer Feature Intent

After scope is defined and before writing the proposal, infer the feature intent based on the proposal's **Proposed Solution** and **Scope** content (not just the title).

### Intent Mapping Rules

| Task Type Pattern | Default Intent |
|-------------------|---------------|
| `coding.feature` / `coding.enhancement` | `new-feature` |
| `coding.cleanup` | `cleanup` |
| `coding.refactor` | `refactor` |
| `coding.fix` | Heuristic: does the fix introduce new user-observable behavior? Yes → `new-feature`, internal adjustment only → `refactor` |

For proposals with mixed content (both new behavior and refactoring), determine the **primary** intent by assessing whether the proposal's core goal includes any new externally observable behavior (new API, new CLI command, new output format). If yes → primary intent is `new-feature` (full test pipeline ensures new behavior has test coverage). If the core goal is purely reorganizing internal implementation → `refactor`.

### Confirmation

Use `AskUserQuestion` to present the inferred intent to the user:

- Show the inferred intent value and the reasoning (which mapping rule or heuristic was applied)
- The user can override to any of the three valid values: `new-feature`, `refactor`, `cleanup`
- Once confirmed, the intent value is written into the proposal frontmatter's `intent` field

**Example prompt**: "Based on the proposed solution and scope, I infer this feature's intent as **`new-feature`** (introduces new user-observable behavior). Is this correct, or would you like to override to `refactor` or `cleanup`?"

## Step 5: Write Proposal

Save to `docs/proposals/<slug>/proposal.md` using `templates/proposal.md`.

### SC Consistency Check (mandatory)

After writing **Success Criteria** and **In Scope** sections, execute the SC consistency check defined in `rules/sc-consistency.md`. This is a mandatory step — the proposal must not proceed to Quality Standards without it. The check detects logical contradictions within SC entries (SC-to-SC) and between SC and InScope entries (SC-to-InScope) through clustering and bidirectional satisfiability proof. If conflicts are found, present them to the user and resolve before continuing.

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

<EXTREMELY-IMPORTANT>
Eval auto-run check — do NOT use AskUserQuestion when config enables auto-run.

Run the following config check sequence via Bash tool:

```bash
# Eval auto-run check (proposal)
EVAL_ENABLED=$(forge config get auto.eval.proposal 2>/dev/null)
if [ "$EVAL_ENABLED" = "true" ]; then
  echo "AUTO_RUN"
elif [ "$EVAL_ENABLED" = "false" ]; then
  echo "SKIP"
else
  echo "FALLBACK_ASK"
fi
```

Based on the output:
- **AUTO_RUN** → invoke `/eval-proposal` via `Skill` tool (default: 900 points / 3 rounds)
- **SKIP** → skip eval, output "eval-proposal 已通过配置跳过", proceed to `/write-prd`
- **FALLBACK_ASK** → ask via `AskUserQuestion`: "Run `/eval-proposal` for adversarial evaluation? (default: 900 points / 3 rounds)"
  - **Yes** → invoke `/eval-proposal` via `Skill` tool
  - **Custom** → invoke `/eval-proposal --target X --iterations Y` via `Skill` tool
  - **No** → proceed to `/write-prd`
</EXTREMELY-IMPORTANT>
