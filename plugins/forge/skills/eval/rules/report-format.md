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
- `harness`: priority improvement table (P0/P1/P2)
- `consistency`: "Files Modified" and "Residual Issues"
- `design`: Breakdown-Readiness gate status

Save report to type-specific report path.
