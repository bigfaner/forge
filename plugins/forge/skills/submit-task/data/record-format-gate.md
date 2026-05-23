# Record Format: Gate Category

Task types: `gate`

## JSON Example

```json
{
	"taskId": "20",
	"status": "completed",
	"summary": "Quality gate check passed for feature branch",
	"gatePassed": true,
	"gateChecks": [
		"All tests passing",
		"No lint errors",
		"Coverage above threshold"
	],
	"notes": "All checks passed"
}
```

## Category-Specific Fields

| Field         | Type     | Required | Description                                    |
| ------------- | -------- | -------- | ---------------------------------------------- |
| `gatePassed`  | bool     | optional | Whether the gate check passed                  |
| `gateChecks`  | string[] | optional | List of individual gate check results          |

## Rules

- Gate tasks produce minimal records — gate checks + overall pass status
- Do NOT include `filesCreated`, `filesModified`, `testsPassed`, `testsFailed`, `coverage`
- Do NOT include `acceptanceCriteria` unless the gate explicitly checks against criteria

## Metrics Collection

Collect from gate checks:
- `gatePassed` — boolean based on overall gate outcome
- `gateChecks` — list each check and its result
