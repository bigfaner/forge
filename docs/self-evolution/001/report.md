# Eval-Plugin Complete

**Final Score**: 913/1000 (target: 900)
**Plugin Version**: 2.16.1
**Iterations Used**: 2/3

### Score Progression

| Iteration | Score | Delta |
|-----------|-------|-------|
| 1 | 893 | - |
| 2 | 913 | +20 |

### Dimension Breakdown (final)

| Dimension | Score | Max |
|-----------|-------|-----|
| 1. Directory-Name Alignment | 40 | 40 |
| 2. Agent Reference Integrity | 100 | 100 |
| 3. Reference Integrity | 80 | 80 |
| 4. Frontmatter Completeness | 110 | 110 |
| 5. Eval Template Convention | 100 | 100 |
| 6. Orchestrator Convention | 40 | 40 |
| 7. Task CLI Alignment | 228 | 240 |
| 8. Hook Wiring Integrity | 70 | 70 |
| 9. Guide Coverage | 40 | 70 |
| 10. Command Metadata | 60 | 60 |
| 11. Plugin Metadata | 40 | 40 |
| 12. Safety Marker Consistency | 35 | 50 |

### Files Modified

| File | Changes |
|------|---------|
| `plugins/forge/commands/git-checkout.md` | Added missing `name` field + `allowed_tools` declaration |
| `plugins/forge/commands/extract-design-md.md` | Added `allowed_tools` + safety marker block |
| `plugins/forge/commands/gen-sitemap.md` | Added `allowed_tools` declaration |
| `plugins/forge/commands/init-forge.md` | Added `allowed_tools` + safety marker block |
| `plugins/forge/commands/init-justfile.md` | Added `allowed_tools` + safety marker block |
| `plugins/forge/commands/record-decision.md` | Added safety marker block |
| `plugins/forge/.claude-plugin/plugin.json` | Added 6 keywords: design, test-generation, forensic, proposal, record, consolidation |
| `plugins/forge/skills/record-task/SKILL.md` | Added Quality Gate Pre-check documentation to Validation Rules |

### Residual Issues

1. **Guide coverage (D9, -30 pts)**: ~10 utility skills/commands have zero mention in guide.md (improve-harness, forensic, eval-harness, eval-consistency, git-commit, simplify-skill, extract-design-md, init-forge, init-justfile, record-task). Adding sections requires editorial judgment on structure and scope.
2. **Schema-code mismatch (D7, -5 pts)**: `index.schema.json` marks `prd`/`design` as required, but Go allows `proposal` as alternative (quick mode). Known acceptable discrepancy per rubric but still a structural inconsistency.
3. **Safety marker coverage (D12, -15 pts)**: Utility commands (init-forge, init-justfile, record-decision, extract-design-md) added safety blocks in iteration 1, but scorer notes they may need stronger markers. Also `fix-bug` dispatch command lacks `EXTREMELY-IMPORTANT` block (has only `HARD-GATE` and `HARD-RULE`).
4. **Guide naming (D7, -2 pts)**: guide.md says "just probe" but actual code uses internal Go function `e2eprobe.ProbeServers()`, not a just recipe.
5. **breakdown-tasks root-ancestor (D7, -5 pts)**: SKILL.md omits CLI's auto-resolution to root ancestor for `--source-task-id`.

### Outcome

Target reached (913 >= 900). 5 residual issues remain — most require human judgment on documentation scope and editorial decisions.
