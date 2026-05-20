# Challenge Protocol

Challenge is not a separate step — it is a **mandatory behavior** embedded in every decision point. Each challenge must be grounded in facts (see Fact-Driven Principle). Empty or vague questioning is forbidden.

## Challenge Tools

| Tool | When to Use | Trigger Condition | Termination Condition |
|------|-------------|-------------------|----------------------|
| **5 Whys** | Problem Cluster — drill into root cause | User states a surface-level symptom or vague pain point | Root cause identified (causal chain reaches a fundamental constraint), OR 3 consecutive "why" answers are consistent |
| **XY Problem Detection** | Problem Cluster — detect when user's stated need may not be the real need | User asks for a specific solution rather than describing a problem | User confirms the actual underlying problem, OR user provides clear rationale for why the specific solution is needed |
| **Assumption Flip** | Solution Cluster — validate critical assumptions | User presents a solution that depends on an unverified assumption | Assumption is confirmed by evidence, OR assumption is overturned and solution is adjusted, OR user provides sufficient domain expertise as evidence |
| **Stress Test** | Solution Cluster — expose hidden risks in seemingly perfect solutions | User seems satisfied with a solution without considering edge cases | All identified edge cases are addressed, OR user explicitly accepts residual risk with documented rationale |
| **Occam's Razor** | Both Clusters — meta-principle applied at all times | Multiple competing explanations or solutions coexist | Simplest viable option selected, OR complexity is justified by concrete evidence |

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

## Occam's Razor (Meta-Principle)

Throughout the entire brainstorm, when multiple explanations or approaches coexist, prefer the simplest one that satisfies all known constraints. This is not a separate tool to invoke — it is a standing principle that applies to every decision:

- If a simple explanation covers the observed problem, do not propose complex alternatives without evidence requiring complexity
- If a straightforward approach meets all success criteria, do not add sophistication without justification
- When the user proposes a complex solution, ask: "Is there a simpler way to achieve the same outcome?"
