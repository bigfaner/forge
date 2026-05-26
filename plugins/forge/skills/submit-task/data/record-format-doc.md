# Record Format: Doc Category

Task types: `doc`, `doc.review`, `doc.summary`, `doc.consolidate`, `doc.drift`

## JSON Example

```json
{
  "taskId": "7",
  "status": "completed",
  "summary": "Updated submit-task SKILL.md with per-type instructions",
  "filesCreated": [],
  "filesModified": ["plugins/forge/skills/submit-task/SKILL.md"],
  "referencedDocs": [
    "docs/proposals/typed-task-records/proposal.md"
  ],
  "reviewStatus": "final",
  "docMetrics": "~120 lines across 5 categories",
  "acceptanceCriteria": [{ "criterion": "JSON examples for all 5 categories", "met": true }],
  "notes": "Followed existing structure"
}
```

## Category-Specific Fields

| Field            | Type     | Required | Description                                        |
| ---------------- | -------- | -------- | -------------------------------------------------- |
| `referencedDocs` | string[] | optional | Document paths referenced during work               |
| `reviewStatus`   | string   | optional | Review status: "draft", "final", or "reviewed"      |
| `docMetrics`     | string   | optional | Scores from evaluation dimensions (e.g. "coverage: 85, consistency: 90") |
