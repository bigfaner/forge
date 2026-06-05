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
- `contract`: Per-dimension threshold pass/fail table; eval-skipped flag if parse failure after retry. **Anchor Integrity section** (when handbook exists): "### Missing Anchor Fields" table listing Contracts missing required anchor fields, grouped by surface type with columns: Contract File | Missing Field | Expected Value (from handbook). "### Handbook Conflicts" table listing internal handbook inconsistencies: Conflict Type (method-conflict / path-conflict / command-conflict) | Entry A | Entry B | Description. When no handbook found, output "Anchor Integrity: handbook not found — check skipped (backward-compatible)".

Save report to type-specific report path.
