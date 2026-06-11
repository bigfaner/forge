# Challenge Protocol

Challenge is not a filter — it is a **clarification behavior** that helps the user see their own assumptions more clearly. Every challenge must be grounded in facts (see Fact-Driven Principle). Empty or vague questioning is forbidden.

## Need Gate (Clarification Checkpoint)

Every feature request passes through this gate (embedded in SKILL.md Step 2, Walk the Design Tree) to clarify what the user truly needs. The gate is not about saying "no" — it's about making sure we're solving the right problem before investing energy in implementation details. It triggers when a concrete feature crystallizes from the discussion — not on the initial vague idea.

### Check 1: Simpler Alternative (Occam's Razor)

Before building anything new, check if existing tools already solve the problem.

**Search strategy** (execute at least 2 of these):
1. **Codebase keyword search** — search for the feature name, related verbs, and domain terms
2. **Pipe composition check** — can existing commands be combined (e.g., `forge proposal | grep <keyword>`)?
3. **Configuration check** — can existing functionality meet the need through settings or parameters?
4. **Ecosystem check** — are standard tools (grep, jq, fzf, sort) or platform features sufficient?

If a simpler alternative exists → propose it. If the user needs a better experience than the alternative provides, that's a valid reason to proceed — but the discussion shifts from "new feature" to "experience upgrade over [existing tool]."

### Check 2: Real Need (XY Detection)

Users often ask for X when they really need Y. Before challenging X, **confirm Y with the user**:

1. Identify what the user asked for (X)
2. Hypothesize the underlying goal (Y)
3. **Ask the user**: "I understand your core need is [Y] — is that right?"
4. Only after Y is confirmed: assess whether X is the best path to Y

Do NOT challenge based on an unconfirmed Y hypothesis — a wrong guess wastes everyone's time.

### Check 3: Timing

"Why now instead of later?"

- If the cost of deferring is low (no blocked users, no architectural window closing) → suggest deferring until real demand emerges
- If there's a concrete deadline, dependency, or compounding cost → the timing is justified

This catches "we'll eventually need this" requirements that are true but premature.

### User Override

If the user insists on proceeding after hearing the gate's findings: **accept immediately**. Record in the proposal:
```
Challenge Override: user chose to proceed despite <finding>.
Reason: <user's reason or "not stated">.
Potential risk: <one-sentence risk if applicable>.
```

Do NOT re-challenge after an override. The user's decision is final.

## Challenge Tools

These tools operate within the Decision Clusters (Problem and Solution) during Walk the Design Tree (SKILL.md Step 2). Need Gate uses Occam's Razor and XY Detection as structured checks (see above) — the table below covers their ongoing use plus additional tools.

| Tool | When to Use | Trigger Condition | Termination Condition |
|------|-------------|-------------------|----------------------|
| **5 Whys** | Problem Cluster — drill into root cause | User states a surface-level symptom or vague pain point | Root cause identified (causal chain reaches a fundamental constraint), OR 3 consecutive "why" answers are consistent |
| **XY Problem Detection** | Problem Cluster — ongoing use after Need Gate | User's stated goal and actual need seem misaligned during discussion | User confirms the actual underlying problem, OR user provides clear rationale for why the specific solution is needed |
| **Assumption Flip** | Solution Cluster — validate critical assumptions | User presents a solution that depends on an unverified assumption | Assumption is confirmed by evidence, OR assumption is overturned and solution is adjusted, OR user provides sufficient domain expertise as evidence |
| **Stress Test** | Solution Cluster — expose hidden risks in seemingly perfect solutions | User seems satisfied with a solution without considering edge cases | All identified edge cases are addressed, OR user explicitly accepts residual risk with documented rationale |

## Fact-Driven Principle

Every challenge must cite one of three evidence types:

1. **Codebase facts** — existing implementations, architecture patterns, API contracts found in the codebase
2. **Logical consistency** — internal contradictions, circular reasoning, or gaps in the user's own argument
3. **Domain common sense** — widely accepted knowledge in the relevant technical or business domain

For greenfield projects (no existing codebase), rely on logical consistency and domain common sense. The absence of code does not excuse challenges from providing evidence.

## Challenge Tone

Challenges must be **rationally prudent**, not hostile. Every challenge follows this structure:

1. **State the observation** — what was said or assumed
2. **Present the evidence** — cite a specific fact from the three evidence types
3. **Pose the question** — ask what the implication is

Example: "You mentioned caching as the solution for latency. The current p99 latency is 50ms (codebase fact: `src/api/middleware/timer.ts`). At this level, network round-trip dominates — caching may not address the actual bottleneck. What does the latency breakdown show?"

## Occam's Razor (Standing Principle)

Beyond the Need Gate's structured check, Occam's Razor applies as a standing principle throughout the brainstorm:

- If a simple explanation covers the observed problem, do not propose complex alternatives without evidence requiring complexity
- If a straightforward approach meets all success criteria, do not add sophistication without justification
- When in doubt, ask: "Is there a simpler way to achieve the same outcome?"
