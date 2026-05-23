# Record Format: Validation Category

Task types: `validation.code`, `validation.ux`

## JSON Example

```json
{
  "taskId": "15",
  "status": "completed",
  "summary": "Validated code quality for auth module",
  "filesCreated": [],
  "filesModified": ["src/auth/middleware.ts"],
  "validationPassed": true,
  "issuesFound": ["Minor: unused import on line 42"],
  "acceptanceCriteria": [{ "criterion": "No critical issues found", "met": true }],
  "notes": "One minor issue auto-fixed during validation"
}
```

## Category-Specific Fields

| Field              | Type     | Required | Description                    |
| ------------------ | -------- | -------- | ------------------------------ |
| `validationPassed` | bool     | optional | Whether validation passed      |
| `issuesFound`      | string[] | optional | Issues found during validation |
