# Record Format: Doc Category

Task types: `doc`, `doc.eval`, `doc.summary`, `doc.consolidate`, `doc.drift`

## JSON Example

```json
{
	"taskId": "7",
	"status": "completed",
	"summary": "Updated submit-task SKILL.md with per-type instructions",
	"filesCreated": [],
	"filesModified": ["plugins/forge/skills/submit-task/SKILL.md"],
	"keyDecisions": ["Added type-specific sections after existing Fields table"],
	"referencedDocs": [
		"docs/proposals/typed-task-records/proposal.md",
		"docs/features/typed-task-records/tasks/1-category-and-recorddata.md"
	],
	"reviewStatus": "final",
	"docMetrics": "Added ~120 lines covering 5 categories with JSON examples and field tables",
	"acceptanceCriteria": [{ "criterion": "JSON examples for all 5 categories", "met": true }],
	"notes": "Followed existing structure; added sections rather than reorganizing"
}
```

## Category-Specific Fields

| Field             | Type     | Required | Description                                      |
| ----------------- | -------- | -------- | ------------------------------------------------ |
| `referencedDocs`  | string[] | optional | List of document paths referenced during work     |
| `reviewStatus`    | string   | optional | Review status (e.g., "draft", "final", "reviewed") |
| `docMetrics`      | string   | optional | Free-text document metrics (word count, sections, etc.) |

## Rules

- Do NOT include `testsPassed`, `testsFailed`, or `coverage` — these are auto-set to "no tests" by the CLI
- Focus on `referencedDocs`, `reviewStatus`, `docMetrics` instead

## Metrics Collection

No test metrics required. Focus on:
- `referencedDocs` — list documents you read/modified
- `reviewStatus` — set to "draft", "final", or "reviewed" as appropriate
- `docMetrics` — describe what was produced (e.g., word count, sections added)
