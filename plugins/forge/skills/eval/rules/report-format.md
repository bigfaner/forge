# Final Report Format

```
## Eval-{{TYPE}} Complete
**Final Score**: {{SCORE}}/{{SCALE}} (target: {{TARGET}})
**Iterations Used**: {{N}}/{{MAX}}

### Score Progression
| Iteration | Score | Delta |
|-----------|-------|-------|

### Dimension Breakdown (final)
{{from rubric}}

### Outcome
{{"Target reached" / "Target NOT reached -- N iterations exhausted"}}
```

Type-specific additions:
- `consistency`: "Files Modified" and "Residual Issues"
- `design`: Breakdown-Readiness gate status
- `journey`: Per-dimension threshold pass/fail table; eval-skipped flag if parse failure after retry
- `contract`: Per-dimension threshold pass/fail table; eval-skipped flag if parse failure after retry

Save report to type-specific report path.
