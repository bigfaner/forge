# Record Format: Test Category

Task types: `test.gen-cases`, `test.eval-cases`, `test.gen-scripts`, `test.run`, `test.gen-and-run`, `test.verify-regression`

## JSON Example

```json
{
  "taskId": "12",
  "status": "completed",
  "summary": "Generated test cases for login flow",
  "filesCreated": ["tests/e2e/login/login-flow.spec.ts"],
  "filesModified": [],
  "casesGenerated": 8,
  "casesEvaluated": 8,
  "scriptsCreated": ["tests/e2e/login/login-flow.spec.ts"],
  "testResults": "8 cases generated, 8 passed evaluation criteria",
  "acceptanceCriteria": [{ "criterion": "All login scenarios covered", "met": true }],
  "notes": "Auto-generated via test pipeline"
}
```

## Category-Specific Fields

| Field             | Type     | Required | Description                          |
| ----------------- | -------- | -------- | ------------------------------------ |
| `casesGenerated`  | int      | optional | Number of test cases generated       |
| `casesEvaluated`  | int      | optional | Number of test cases evaluated       |
| `scriptsCreated`  | string[] | optional | Test script files created            |
| `testResults`     | string   | optional | Summary of test execution results    |

## Rules

- Use `testResults` for run outcomes. These are pipeline metrics, distinct from the coding quality gate.
