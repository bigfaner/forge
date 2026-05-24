# Record Format: Eval Category

Task types: `eval.journey`, `eval.contract`

## JSON Example

```json
{
  "taskId": "8",
  "status": "completed",
  "summary": "Evaluated journey quality against rubric",
  "score": 850,
  "findings": ["Journey lacks error handling scenarios", "Missing accessibility test cases"],
  "severity": "major",
  "passed": true,
  "acceptanceCriteria": [{ "criterion": "Score >= 850", "met": true }],
  "notes": "Two major findings, both non-blocking for this iteration"
}
```

## Category-Specific Fields

| Field       | Type     | Required | Description                                         |
|-------------|----------|----------|-----------------------------------------------------|
| `score`     | float    | optional | Eval score (0-1000)                                 |
| `findings`  | string[] | optional | Issues found during evaluation                      |
| `severity`  | string   | optional | Overall severity (critical/major/minor)             |
| `passed`    | bool     | optional | Whether evaluation passed the quality gate          |

## Rules

- Eval tasks are review-type tasks, not test generators. Do NOT use `testsPassed`/`coverage` as primary evidence.
- At least one eval-specific field (findings, severity, or score) must be present for completed status.
- Use `severity` to indicate the worst issue level found: critical, major, or minor.
