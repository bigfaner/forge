# Record Format: Coding Category

Task types: `coding.feature`, `coding.enhancement`, `coding.cleanup`, `coding.refactor`, `coding.fix`, `code-quality.simplify`

## JSON Example

```json
{
  "taskId": "3.3.1",
  "status": "completed",
  "summary": "Implemented Button component with accessibility support",
  "filesCreated": ["src/components/Button.tsx"],
  "filesModified": ["src/utils/helpers.ts"],
  "keyDecisions": ["Used compound component pattern for flexibility"],
  "testsPassed": 12,
  "testsFailed": 0,
  "coverage": 85.6,
  "acceptanceCriteria": [{ "criterion": "Button renders all variants", "met": true }],
  "notes": "Tested with screen reader"
}
```

## Category-Specific Fields

| Field          | Type     | Required    | Description                          |
| -------------- | -------- | ----------- | ------------------------------------ |
| `testsPassed`  | int      | conditional | Number of tests passed               |
| `testsFailed`  | int      | conditional | Number of tests failed               |
| `coverage`     | float    | conditional | Coverage percentage from test runner |
| `keyDecisions` | string[] | optional    | Key decisions made during execution  |

> `testsPassed`, `testsFailed`, `coverage`: required when `status: completed`; omit when `status: blocked`.

## Rules

- `testsFailed > 0` with `completed` triggers auto-downgrade to `blocked`.
- All metrics must come from actual test runner output — never guess or default.
- Never write `0.0` for coverage unless the runner reported 0%.
