---
created: "{{DATE}}"
sessions: [{{SESSION_IDS}}]
skillsInvolved: [{{SKILL_NAMES}}]
severity: "{{SEVERITY}}"
---

# {{TITLE}}

## Executive Summary

{{EXECUTIVE_SUMMARY}}

## Investigation Scope

| Dimension | Value |
|-----------|-------|
| Sessions analyzed | {{SESSION_COUNT}} |
| Time range | {{EARLIEST}} to {{LATEST}} |
| Skills involved | {{SKILLS_LIST}} |
| Trigger | {{TRIGGER_DESCRIPTION}} |

## Timing Overview

| Session | Duration | Tool Time | Idle* | Top Bottleneck |
|---------|----------|-----------|-------|---------------|
| {{SESSION_ID}} | {{DURATION}} | {{TOTAL_TOOL_MS}} | {{IDLE}} | `{{TOOL}}` ({{MAX}}s) |

*Idle = session duration minus total tool execution time — indicates thinking/waiting.

| Tool | Calls | Total | Avg | Max |
|------|-------|-------|-----|-----|
| `{{TOOL}}` | {{CALL_COUNT}} | {{TOTAL}}s | {{AVG}}s | {{MAX}}s |

## Findings

### Finding 1: {{FINDING_TITLE}}

**Category:** `{{DEVIATION_CATEGORY}}`

**Affected sessions:** {{AFFECTED_SESSION_IDS}}

**Symptom:**
{{OBSERVED_WRONG_BEHAVIOR}}

**Agent reasoning (from thinking block):**
> {{AGENT_REASONING}}

**Expected behavior (from skill definition):**
> {{EXPECTED_BEHAVIOR}}

**Gap:**
{{DEVIATION_ROOT_CAUSE}}

**Causal chain:**
1. **Symptom:** {{OBSERVABLE_WRONG_BEHAVIOR}}
2. **Direct cause:** {{SPECIFIC_ACTION}}
3. **Root cause:** {{ROOT_CAUSE}}

### Finding 2: {{FINDING_TITLE}}

*(repeat structure for each finding)*

## Cross-Session Patterns

{{CROSS_SESSION_PATTERNS}}

| Pattern | Sessions | Category |
|---------|----------|----------|
| {{PATTERN_DESCRIPTION}} | {{SESSION_IDS}} | {{CATEGORY}} |

## Recommendations

| Priority | Action | Target File | Finding |
|----------|--------|-------------|---------|
| P0 | {{WHAT_TO_CHANGE}} | `{{TARGET_FILE}}` | Finding N |
| P1 | {{WHAT_TO_CHANGE}} | `{{TARGET_FILE}}` | Finding N |

## Evidence

Evidence files at: `docs/forensics/{{SLUG}}/evidence/`

| File | Source | Size |
|------|--------|------|
| evidence.json | Main session | ~XX KB |
| evidence-subagent-{{ID}}.json | Subagent {{SUBAGENT_TYPE}} | ~XX KB |
