# Record Format: Gate Category

Task types: `gate`

## JSON Example

```json
{
  "taskId": "20",
  "status": "completed",
  "summary": "Quality gate check passed for feature branch",
  "filesCreated": [],
  "filesModified": [],
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

| Field        | Type     | Required | Description                      |
| ------------ | -------- | -------- | -------------------------------- |
| `gatePassed` | bool     | optional | Whether the gate check passed    |
| `gateChecks` | string[] | optional | Individual gate check results    |

## Rules

- The gate MUST verify all project code compiles before passing.
- Do NOT include `acceptanceCriteria` unless the gate explicitly checks against criteria.
