---
created: "<DATE>"
sessions: [<SESSION_IDS>]
skillsInvolved: [<SKILL_NAMES>]
severity: <high|medium|low>
---

# <Title>

## Executive Summary

<1-3 sentences: what was investigated, what was found, what should change>

## Investigation Scope

| Dimension | Value |
|-----------|-------|
| Sessions analyzed | <count> |
| Time range | <earliest> to <latest> |
| Skills involved | <list> |
| Trigger | <what prompted this investigation> |

## Timing Overview

| Session | Duration | Tool Time | Idle* | Top Bottleneck |
|---------|----------|-----------|-------|---------------|
| <session-id-8ch> | <duration> | <totalToolMs> | <idle> | `<tool>` (<max>s) |

*Idle = session duration minus total tool execution time — indicates thinking/waiting.

| Tool | Calls | Total | Avg | Max |
|------|-------|-------|-----|-----|
| `<tool>` | <count> | <total>s | <avg>s | <max>s |

## Findings

### Finding 1: <Title>

**Category:** `<deviation-category>`

**Affected sessions:** <session IDs>

**Symptom:**
<What was observed — the wrong behavior>

**Agent reasoning (from thinking block):**
> <Exact or paraphrased thinking block content showing the agent's decision>

**Expected behavior (from skill definition):**
> <What the skill definition says should happen>

**Gap:**
<Why the agent deviated — the root cause>

**Causal chain:**
1. **Symptom:** <observable wrong behavior>
2. **Direct cause:** <specific action/decision>
3. **Root cause:** <why the agent made that decision>

### Finding 2: <Title>

*(repeat structure for each finding)*

## Cross-Session Patterns

<If multiple sessions show the same deviation, describe the pattern here. Otherwise remove this section.>

| Pattern | Sessions | Category |
|---------|----------|----------|
| <pattern description> | <IDs> | <category> |

## Recommendations

| Priority | Action | Target File | Finding |
|----------|--------|-------------|---------|
| P0 | <what to change> | `path/to/file` | Finding N |
| P1 | <what to change> | `path/to/file` | Finding N |

## Evidence

Evidence files at: `docs/forensics/<slug>/evidence/`

| File | Source | Size |
|------|--------|------|
| evidence.json | Main session | ~XX KB |
| evidence-subagent-<id>.json | Subagent <type> | ~XX KB |
