# Record Format: Coding Category

Task types: `coding.feature`, `coding.enhancement`, `coding.cleanup`, `coding.refactor`, `coding.fix`, `code-quality.simplify`

## JSON Example

```json
{
	"taskId": "3.3.1",
	"status": "completed",
	"summary": "What was implemented",
	"filesCreated": ["src/components/Button.tsx"],
	"filesModified": ["src/utils/helpers.ts"],
	"keyDecisions": ["Decision 1"],
	"testsPassed": 12,
	"testsFailed": 0,
	"coverage": 85.6,
	"acceptanceCriteria": [{ "criterion": "Acceptance criterion 1", "met": true }],
	"notes": "Optional observations"
}
```

## Category-Specific Fields

| Field          | Type  | Required | Description                                  |
| -------------- | ----- | -------- | -------------------------------------------- |
| `testsPassed`  | int   | context  | Number of tests passed                       |
| `testsFailed`  | int   | context  | Number of tests failed. >0 with completed = auto-downgrade to blocked |
| `coverage`     | float | context  | Coverage percentage from test runner         |

> **context** = required for `completed` tasks with a `coding.*` type; auto-relaxed for other types.

## Metrics Collection

Collect from the project's test runner. All numeric fields must come from actual output, never guessed or defaulted.

- `coverage` = actual percentage from test runner output. Never write `0.0` unless the runner actually reported 0%.
- Capture metrics from the targeted test runs you performed during task development.
