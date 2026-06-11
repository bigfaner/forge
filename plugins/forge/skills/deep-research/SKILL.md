---
name: deep-research
description: Systematically research a technology or product. Supports single-tech deep dive and multi-candidate comparison. Outputs a structured research report.
argument-hint: "<topic> [--compare <candidates>] [--focus <dimensions>]"
effort: high
---

# Deep Research

From a topic name to a grounded, evidence-based research report — through adaptive multi-source investigation.

**Core principle**: Research is not summarizing search results. It's gathering evidence, cross-referencing sources, contextualizing to the user's specific situation, and producing actionable insight.

<HARD-GATE>
Do NOT write any code or take implementation action. This skill produces a research document only.
</HARD-GATE>

## When to Use

**Trigger:**
- User explicitly requests `/deep-research` or `/deep-research <topic>`
- User asks to investigate, compare, or evaluate a technology/product
- User needs a structured analysis before making a tech selection

**Skip:**
- Quick factual questions (answer directly)
- Recording known decisions (use `/learn`)
- Generating a proposal from a clear idea (use `/brainstorm`)

## Parameters

| Parameter | Default | Description |
|---|---|---|
| `<topic>` | — | Research topic, e.g. "Deno 2.0", "ORM selection" |
| `--compare` | — | Comma-separated candidates for comparison mode, e.g. "drizzle,prisma,postgres.js" |
| `--focus` | — | Comma-separated dimensions to focus on, e.g. "performance,security,ecosystem" |

Parameters serve as **prefill hints** — they reduce the questions asked in Phase 1 but do not override the user's confirmed intent.

## Workflow

```
Phase 1: Clarify needs (AskUserQuestion) → Phase 2: Execute research (adaptive) → Phase 3: Output report
```

### Phase 1: Clarify Needs

Use `AskUserQuestion` to confirm the research scope. Two rounds maximum.

<HARD-RULE>
Use `AskUserQuestion` for all questions. Maximum 2 rounds, 4 questions per round. Parameters prefill answers where applicable — skip questions the user already answered via args.
</HARD-RULE>

**Round 1 (2-3 questions):**

| # | Question | When to skip |
|---|----------|-------------|
| Q1 | Research mode: **single-tech deep dive** or **multi-candidate comparison**? | Skip if `--compare` provided (default to comparison) |
| Q2 | (If comparison) Which candidates to compare? | Skip if `--compare` fully specifies the list |
| Q2' | (If deep dive) What specific aspects are you most interested in? | Always ask |
| Q3 | Briefly describe your goal — what question should this research answer? | Always ask |

**Round 2 (2 questions):**

| # | Question | Format |
|---|----------|--------|
| Q4 | Select research dimensions to cover (multiSelect) — present the dimension set matching the research mode from `rules/research-dimensions.md` (use `single-tech` dimensions for deep dive mode, `comparison` dimensions for comparison mode). Core dimensions (marked with `*` in `rules/research-dimensions.md`) are pre-selected. | AskUserQuestion with multiSelect |
| Q5 | Include project adaptation assessment? (scan current codebase for impact analysis) | AskUserQuestion single select |

After Phase 1, you should have a clear research plan: mode, scope, dimensions, and whether project context is needed.

### Phase 2: Execute Research

Adaptive multi-source research. The agent decides how many search rounds are needed based on information convergence.

<HARD-RULE>
**Convergence rule**: Stop researching a topic when 2 consecutive search rounds yield no substantively new information. Do not cap rounds artificially — trust the convergence signal.
</HARD-RULE>

**Information sources (in priority order):**

1. **Official documentation** — Fetch the technology's official docs, README, getting-started guide
2. **Community & ecosystem** — WebSearch for comparisons, benchmarks, community sentiment, adoption trends
3. **Codebase scan** (if Q5 = yes) — Use Glob/Grep/Agent(Explore) to find related code, dependencies, patterns that would be affected
4. **Mid-research clarification** — If a critical ambiguity blocks progress, ask the user via AskUserQuestion (max 2 mid-research questions)

**Research execution pattern:**

For each dimension confirmed in Phase 1:
1. Search official sources for authoritative information
2. Search community sources for real-world experience and benchmarks
3. Cross-reference claims across sources — flag contradictions
4. Record findings with source URLs for citation

For comparison mode specifically:
- Research each candidate independently before comparing
- Ensure all candidates are evaluated on the same dimensions
- Flag dimensions where information is unavailable for a candidate

### Phase 3: Output Report

1. Read `templates/research-report.md` for the report structure.
2. Write the report to `docs/research/<slug>.md` where `<slug>` is a kebab-case derivation of the topic.
3. Present the report to the user.
4. **Wait for the user to review the report and respond.** Do not proceed to the conversion question until the user has acknowledged or commented on the report content.
5. After the user has reviewed, ask: **"是否要将本次调研结合其他需求，转为提案？如需要，可运行 /brainstorm"**

<HARD-RULE>
Do NOT commit automatically. Present the report and let the user review first.
</HARD-RULE>

## Report Structure

See `templates/research-report.md` for the full template.

## Common Mistakes

| Mistake | Correction |
|---------|------------|
| Summarizing only the first page of search results | Cross-reference multiple sources, seek primary evidence |
| Writing marketing copy instead of analysis | Present trade-offs honestly; every strength has a cost |
| Skipping the official docs in favor of blog posts | Always start with official documentation |
| Researching candidates on different dimensions | Compare on the same axes; note gaps explicitly |
| Producing a generic report detached from the user's context | Anchor findings to the user's stated goal (Q3) |
| Asking more than 2 mid-research clarification questions | If you need more than 2, the Phase 1 clarification was insufficient — acknowledge and do your best |
