---
id: "{{ID}}"
title: "{{TITLE}}"
priority: "P0"
estimated_time: "15min"
dependencies: [{{DEPENDENCIES}}]
status: pending
type: "doc-generation.summary"
mainSession: false
---

# {{ID}}: {{TITLE}}

## Description

Generate a structured summary of all completed tasks in this phase. This summary is read by subsequent phase tasks to maintain cross-phase consistency.

## Instructions

### Step 1: Read all phase task records

Read each record file from `docs/features/<slug>/tasks/records/` whose filename starts with `{{PHASE}}.` and does NOT contain `.summary` (e.g., `{{PHASE}}.1-*.md`, `{{PHASE}}.2-*.md`). Exclude the phase summary's own prior record if one exists.

### Step 2: Extract structured data into the summary field

<HARD-RULE>
The `summary` field in `record.json` MUST follow this exact template. Copy it verbatim and fill in the values. Do NOT restructure, reorder, or omit any section. If a section has no content, write "None."
</HARD-RULE>

```
## Tasks Completed
- {{PHASE}}.1: {{one-line summary from that task's record}}
- {{PHASE}}.2: {{one-line summary from that task's record}}
(...one line per task)

## Key Decisions
- {{each keyDecision from all records, prefixed with task ID: "2.1: ..."}}

## Types & Interfaces Changed
| Name | Change | Affects |
|------|--------|---------|
| {{type/interface name}} | {{added/modified/removed: brief description}} | {{which subsequent tasks care}} |

## Conventions Established
- {{each convention or pattern, prefixed with task ID}}

## Deviations from Design
- {{each deviation from tech-design.md, or "None"}}
```

### Step 3: Populate remaining record.json fields

```json
{
  "taskId": "{{ID}}",
  "status": "completed",
  "summary": "<filled from Step 2 template above>",
  "filesCreated": [],
  "filesModified": [],
  "keyDecisions": ["<list all keyDecisions from all phase records>"],
  "testsPassed": 0,
  "testsFailed": 0,
  "acceptanceCriteria": [
    {"criterion": "All phase task records read and analyzed", "met": true},
    {"criterion": "Summary follows the exact template with all 5 sections", "met": true},
    {"criterion": "Types & Interfaces table lists every changed type", "met": true}
  ]
}
```

## Reference Files

- All phase task records: `docs/features/<slug>/tasks/records/{{PHASE}}.*.md`
- Design reference: `docs/features/<slug>/design/tech-design.md` (Cross-Layer Data Map section)

## Acceptance Criteria

- [ ] All phase task records have been read
- [ ] Summary follows the exact 5-section template above
- [ ] Types & Interfaces Changed table is populated (or "None" if no changes)
- [ ] Record created via `/record-task`

## Implementation Notes

This is a noTest task. No code should be written.
- Proceed directly to generating the summary
- The summary MUST be structured — subsequent phase tasks depend on parsing it
